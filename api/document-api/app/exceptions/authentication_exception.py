"""
Copyright (c) 2023-2025 RapidaAI
Author: Prashant Srivastav <prashant@rapida.ai>

Licensed under GPL-2.0 with Rapida Additional Terms.
See LICENSE.md for details or contact sales@rapida.ai for commercial use.
"""

from app.exceptions.rapida_exception import RapidaException


class AuthenticationException(RapidaException):
    def __init__(
        self,
        message: str,
        auth_type: str,
        status_code: int = 400,
        error_code: int = 1000,
    ):
        super().__init__(
            status_code=status_code,
            message=f"{auth_type}: {message}",
            error_code=error_code,
        )


class MissingAuthorizationKeyException(AuthenticationException):
    error_code: int = 3001
    status_code: int = 400

    def __init__(self, auth_type: str, message: str = "Missing Authorization Key"):
        super().__init__(
            message=message,
            error_code=self.error_code,
            auth_type=auth_type,
            status_code=self.status_code,
        )


class InvalidAuthorizationTokenException(RapidaException):
    error_code: int = 3002
    status_code: int = 401

    def __init__(
        self,
        message: str = "Invalid Authorization Token | "
        "Invalid Signature Error | "
        "Invalid Key Error | "
        "Signature has expired",
    ):
        super().__init__(
            message=message,
            error_code=self.error_code,
            status_code=self.status_code,
        )
