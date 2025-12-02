"""
Copyright (c) 2023-2025 RapidaAI
Author: Prashant Srivastav <prashant@rapida.ai>

Licensed under GPL-2.0 with Rapida Additional Terms.
See LICENSE.md for details or contact sales@rapida.ai for commercial use.
"""
import logging
import time
from collections import OrderedDict
from typing import Dict, Literal, Optional

import botocore.exceptions
from aiobotocore.session import AioSession
from types_aiobotocore_sts.client import STSClient
from types_aiobotocore_sts.type_defs import AssumeRoleResponseTypeDef

from app.configs.auth.aws_auth import AWSAuth
from app.connectors.aws import AWSConnector
from app.exceptions.connector_exception import ConnectorClientFailureException
from app.observabilities import within_span

_log = logging.getLogger("app.connectors.aws.sts_connector")


class STSConnector(AWSConnector):
    """
    Simple token service
    Provide an interface getting token from assuming tole
    """

    # for the purpose of anyone who wants to use this
    # connector directly and implement their own way to operate
    connection: Optional[AioSession] = None

    # credentials store
    __credentials_store: OrderedDict = OrderedDict()

    # aging credentials max in second
    __max_age_in_store: int = 600

    def __init__(self, aws_config: AWSAuth):
        super().__init__(aws_config)

    async def connect(self):
        """
        Connect to AWS S3
        """
        if self.connection:
            _log.debug(f"Already connected to {self.name()}")
            return

        self.connection = self.get_session()

    async def disconnect(self):
        """Not needed for s3"""
        raise NotImplementedError

    async def is_connected(self, bucket_name: str) -> bool:
        """check if bucket exists"""
        raise NotImplementedError

    def name(self) -> Literal["sts"]:
        """
        return sts literal value
        """
        return "sts"

    async def _operate(self, operation: str, **kwargs):
        """
        Operate on s3 client
        :type operation: str
        :param operation
        """
        try:
            # connect / create session for s3 before any operation
            await self.connect()
            async with self.connection.create_client(self.name()) as client:

                # for better suggestions added hint from AIO
                client: STSClient = client
                return await getattr(client, operation)(**kwargs)
        except botocore.exceptions.ConnectionError as boto_error:
            _log.error(f"Failed to connect for {self.name()} . {str(boto_error)}")
            raise ConnectorClientFailureException(
                connector_name=self.name(), message=str(boto_error)
            )
        except botocore.exceptions.ClientError as error:
            _log.error(f"Failed to do {operation} from {self.name()}. {str(error)}")
            error_message = str(error)
            if (
                error.response["Error"]["Code"]
                == "STS.Client.exceptions.RegionDisabledException"
            ):
                error_message = "Illegal region for sts."

            raise ConnectorClientFailureException(
                connector_name=self.name(), message=error_message
            )
        except Exception as err:
            _log.error(
                f"Failed to do the operation {operation} from {self.name()}. {str(err)}",
                exc_info=True,
            )
            raise ConnectorClientFailureException(
                connector_name=self.name(), message=str(err)
            )

    async def get_temporary_credentials(self, session_name: str, **kwargs) -> Dict:
        """
        get_credentials for given role arn
        :param session_name
        :param kwargs:
        """
        _log.debug("Check if the credentials is there in __credentials store")

        store_credentials = self.__credentials_store.get(session_name, None)
        if store_credentials:
            _log.debug("Found value in credentials store.")

            store_credentials_age = time.time() - store_credentials[1]
            if store_credentials_age < self.__max_age_in_store:
                _log.debug(
                    f"Credentials is still valid with time {store_credentials_age}"
                )

                return store_credentials[0]
            else:
                del self.__credentials_store[session_name]
                _log.debug(f"Credentials expired, deleting the key {session_name}")

        _log.debug(f"Requested sts aws credentials for session name {session_name}")
        if not self._aws_config.assume_role:

            _log.debug("Role not set, returning frozen credentials")

            aio_frozen_credential: Dict = await self.get_credentials(session_name)
            self.__credentials_store[session_name] = (
                aio_frozen_credential,
                time.time(),
            )
            return aio_frozen_credential

        with within_span(
            name="STS assuming_role",
            span_type="external",
            span_subtype="aws_sts",
            span_action=f"assume_role {self._aws_config.assume_role}",
        ):
            _log.debug("Role is set, returning temporary credentials.")
            role_definition: AssumeRoleResponseTypeDef = await self._operate(
                "assume_role",
                RoleArn=self._aws_config.assume_role,
                RoleSessionName=session_name,
                **kwargs,
            )
            sts_temporary_credential: Dict = {
                "access_key": role_definition["Credentials"]["AccessKeyId"],
                "secret_key": role_definition["Credentials"]["SecretAccessKey"],
                "token": role_definition["Credentials"]["SessionToken"],
                "region": self.connection.get_config_variable("region"),
            }
            self.__credentials_store[session_name] = (
                sts_temporary_credential,
                time.time(),
            )
            return sts_temporary_credential
