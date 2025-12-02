"""
Copyright (c) 2023-2025 RapidaAI
Author: Prashant Srivastav <prashant@rapida.ai>

Licensed under GPL-2.0 with Rapida Additional Terms.
See LICENSE.md for details or contact sales@rapida.ai for commercial use.
"""
from fastapi import FastAPI, Request
from fastapi.exceptions import RequestValidationError
from pydantic import ValidationError

from app.commons.j_response import JResponse
from app.exceptions.connector_exception import ConnectorClientFailureException
from app.exceptions.rapida_exception import RapidaException


def add_all_exception_handler(app: FastAPI):
    @app.exception_handler(RequestValidationError)
    async def handle_validation_error(
        request: Request, exc: RequestValidationError
    ):  # pylint: disable=unused-argument
        """
        pydantic validation error handler enable
        """
        return JResponse.default_on_error(
            exc=exc,
            error_code=400,
            error_message="validation error for request ensure you have provided all "
            "required fields.",
        )

    @app.exception_handler(ValidationError)
    async def handle_model_validation_error(
        request: Request, exc: ValidationError
    ):  # pylint: disable=unused-argument
        """
        pydantic validation error handler enable
        """
        return JResponse.default_on_error(
            exc=exc,
            error_code=400,
            error_message="validation error for request while parsing the response.",
        )

    @app.exception_handler(ConnectorClientFailureException)
    async def handle_connector_client_error(
        request: Request, exc: ConnectorClientFailureException
    ):  # pylint: disable=unused-argument
        """
        pydantic connector error handler enable
        """
        return JResponse.default_on_error(
            exc=exc,
            error_code=exc.status_code,
            error_message=exc.message,
        )

    @app.exception_handler(RapidaException)
    async def handle_lomotif_error(
        request: Request, exc: RapidaException
    ):  # pylint: disable=unused-argument
        """
        lomotif exception handler enable
        """
        return JResponse.default_on_error(
            exc=exc,
            error_code=exc.status_code,
            error_message=exc.message,
        )
