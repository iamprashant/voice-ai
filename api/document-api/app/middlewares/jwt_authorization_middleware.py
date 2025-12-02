"""
Copyright (c) 2023-2025 RapidaAI
Author: Prashant Srivastav <prashant@rapida.ai>

Licensed under GPL-2.0 with Rapida Additional Terms.
See LICENSE.md for details or contact sales@rapida.ai for commercial use.
"""

import logging
from typing import Dict, Optional, Tuple

import jwt
from fastapi import FastAPI
from jwt import DecodeError
from starlette.authentication import AuthenticationBackend
from starlette.middleware.authentication import AuthenticationMiddleware
from starlette.requests import HTTPConnection

from app.configs.auth_config import JwtAuthenticationConfig
from app.exceptions.authentication_exception import (
    AuthenticationException,
    InvalidAuthorizationTokenException,
    MissingAuthorizationKeyException,
)
from app.middlewares.auth.user import (
    AnonymousUser,
    InternalAuthenticatedUser,
    User,
)

_log = logging.getLogger("app.middlewares.jwt_authorization_middleware")


class JwtAuthorizationMiddleware(AuthenticationMiddleware):
    """
    Authorize user for request using jwt token
    """

    def __init__(self, app: FastAPI, config: JwtAuthenticationConfig):
        super().__init__(backend=JwtAuthBackend(config=config), app=app)


class JwtAuthBackend(AuthenticationBackend):
    """
    starlette custom authentication backend to authenticate user using jwt.
    """

    def __init__(self, config: JwtAuthenticationConfig):
        self.config = config

    async def authenticate(self, conn: HTTPConnection) -> Tuple[bool, Optional[User]]:
        """
        Authenticate user from given jwt token
        :param conn:
        :return:
        """
        try:
            authorization: str = conn.headers.get(self.config.header_key)
            if not authorization:
                raise MissingAuthorizationKeyException(auth_type="JWT")
            payload: Dict = jwt.decode(
                authorization,
                self.config.secret_key.get_secret_value(),
                algorithms=self.config.algorithms,
            )

            if not payload or not payload.get("userId"):
                raise InvalidAuthorizationTokenException("invalid token payload.")
            return True, InternalAuthenticatedUser.parse_obj(payload)
        except DecodeError as err:
            _log.debug(f"Authentication Exception while decoding token: {err}")
            raise InvalidAuthorizationTokenException(
                f"unable to decode given token. {err}"
            )
        except AuthenticationException as ex:
            _log.debug(f"Authentication Exception while authorizing: {ex}")
            if self.config.strict:
                raise ex
            return False, AnonymousUser()
