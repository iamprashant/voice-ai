"""
Copyright (c) 2023-2025 RapidaAI
Author: Prashant Srivastav <prashant@rapida.ai>

Licensed under GPL-2.0 with Rapida Additional Terms.
See LICENSE.md for details or contact sales@rapida.ai for commercial use.
"""


class RapidaException(Exception):
    """
    Lomotif common exceptions class.
    Exception should be controlled by error code.
    """

    def __init__(
        self,
        status_code: int,
        message: str,
        error_code: int,
        error_prefix: str = "RAPIDA",
        service_code: str = "KN_API",
    ):
        """
        :param status_code: http status code
        :param error_prefix: error prefix
        :param service_code: service code <should be unique identifiable>
        """
        super().__init__(self)
        self.status_code = status_code
        self.message = message
        self.error_code = f"{error_prefix}_{service_code}_{error_code}"

    def __str__(self):
        return f"{self.error_code} - {self.message}"
