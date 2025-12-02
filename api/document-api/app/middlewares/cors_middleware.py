"""
Copyright (c) 2023-2025 RapidaAI
Author: Prashant Srivastav <prashant@rapida.ai>

Licensed under GPL-2.0 with Rapida Additional Terms.
See LICENSE.md for details or contact sales@rapida.ai for commercial use.
"""
import logging

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware as _CORSMiddleware

from app.config import ApplicationSettings

_log = logging.getLogger("app.middlewares.cors_middleware")


class CORSMiddleware(_CORSMiddleware):
    """
    CORS middleware for service
    Extension of starlette.middleware.cors.CORSMiddleware with default parameters
    """

    def __init__(self, app: FastAPI, settings: ApplicationSettings):
        super().__init__(
            app=app,
            allow_origins=settings.cors_allow_origins,
            allow_credentials=settings.cors_allow_credentials,
            allow_methods=settings.cors_allow_methods,
            allow_headers=settings.cors_allow_headers,
        )
