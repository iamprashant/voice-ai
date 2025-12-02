"""
Copyright (c) 2023-2025 RapidaAI
Author: Prashant Srivastav <prashant@rapida.ai>

Licensed under GPL-2.0 with Rapida Additional Terms.
See LICENSE.md for details or contact sales@rapida.ai for commercial use.
"""

import logging

from fastapi import FastAPI

from app.bridges.bridge_factory import service_grpc_client
from app.bridges.internals.auth_bridge import AuthBridge
from app.config import ApplicationSettings, get_settings
from app.middlewares.contextual_logging_middleware import ContextualLoggingMiddleware
from app.middlewares.cors_middleware import CORSMiddleware
from app.middlewares.jwt_authorization_middleware import JwtAuthorizationMiddleware
from app.middlewares.request_logger_middleware import RequestLoggerMiddleware
from app.middlewares.token_authorization_middleware import TokenAuthorizationMiddleware

_log = logging.getLogger("app.middlewares")


def add_all_enabled_middleware(app: FastAPI, setting: ApplicationSettings):
    """
    Adding enabled apm middleware to service
    :param app: fastApi app
    :param setting: _Setting, settings of app
    """
    if (
        get_settings().authentication_config
        and get_settings().authentication_config.type == "jwt"
    ):
        app.add_middleware(
            JwtAuthorizationMiddleware,
            config=get_settings().authentication_config.config,
        )

    if (
        get_settings().authentication_config
        and get_settings().authentication_config.type == "token"
    ):
        app.add_middleware(
            TokenAuthorizationMiddleware,
            user_info_resolver=service_grpc_client(
                bridge=AuthBridge,
                service_url=get_settings().authentication_config.config.auth_host,
            ).authorize,
        )

    # if there are any allowed origins, add middleware for cors with all configured cors settings
    app.add_middleware(CORSMiddleware, settings=setting)

    # add all default middleware
    app.add_middleware(RequestLoggerMiddleware)

    # adding contextual middleware
    app.add_middleware(ContextualLoggingMiddleware, settings=setting)
