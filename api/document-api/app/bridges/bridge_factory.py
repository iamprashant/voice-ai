"""
Copyright (c) 2023-2025 RapidaAI
Author: Prashant Srivastav <prashant@rapida.ai>

Licensed under GPL-2.0 with Rapida Additional Terms.
See LICENSE.md for details or contact sales@rapida.ai for commercial use.
"""

import logging
from typing import Type, TypeVar

from fastapi import Depends

from app.bridges import GRPCBridge
from app.bridges.internals.auth_bridge import (
    AuthBridge,
)
from app.bridges.internals.integration_bridge import IntegrationBridge
from app.bridges.internals.vault_bridge import VaultBridge
from app.config import ApplicationSettings, get_settings

GRPC_T = TypeVar("GRPC_T", bound=GRPCBridge)
_log = logging.getLogger("app.bridges.bridge_factory")


def service_grpc_client(bridge: Type[GRPC_T], service_url: str) -> GRPC_T:
    """
    Initialize service client which extends grpc client implementations and return the instance
    :param bridge: class context
    :param service_url: url -> base url
    :return: client wrapper to initiate any rpc call
    :rtype: RPCBridge -> UserBridge
    """

    # url to establish the connections with the grpc server
    if not service_url:
        raise ValueError(
            "Configuration error: service_url is not set for dependable bridge."
        )
    return bridge(service_url=service_url)


def get_me_user_service_client(
        settings: ApplicationSettings = Depends(get_settings),
) -> AuthBridge:
    """
    Dependable function returns user Bridge
    :param settings: application settings to get user_service url
    :return: user service client wrapper to initiate any rpc call
    :rtype: RPCBridge -> UserBridge
    """
    return service_grpc_client(
        bridge=AuthBridge, service_url=settings.lomotif_user_service_url
    )


def get_vault_service_client(
        settings: ApplicationSettings = Depends(get_settings),
) -> VaultBridge:
    return get_me_vault_service_client(settings.internal_service.web_host)


def get_me_vault_service_client(vault_service_url: str) -> VaultBridge:
    """

    Args:
        vault_service_url:

    Returns:

    """
    return service_grpc_client(bridge=VaultBridge, service_url=vault_service_url)


def get_integration_client(
        settings: ApplicationSettings = Depends(get_settings),
) -> IntegrationBridge:
    return get_me_integration_client(settings.internal_service.integration_host)


def get_me_integration_client(integration_service_url: str) -> IntegrationBridge:
    """

    Args:
        integration_service_url:

    Returns:

    """
    return service_grpc_client(
        bridge=IntegrationBridge, service_url=integration_service_url
    )
