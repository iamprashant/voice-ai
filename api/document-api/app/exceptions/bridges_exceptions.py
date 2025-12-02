"""
Copyright (c) 2023-2025 RapidaAI
Author: Prashant Srivastav <prashant@rapida.ai>

Licensed under GPL-2.0 with Rapida Additional Terms.
See LICENSE.md for details or contact sales@rapida.ai for commercial use.
"""
from app.exceptions.rapida_exception import RapidaException


class BridgeException(RapidaException):
    """
    Bridge Error -> all bridge class or function should raise BridgeException or child errors
    """

    def __init__(self, message: str, bridge_name: str, error_code: int = 1000):
        self.bridge_name = bridge_name
        self.status_code = 400
        self.message = f"{bridge_name}: {message}"
        self.error_code = error_code


class BridgeClientException(BridgeException):
    """
    BridgeClientException is wrapper for all client exception raised by aiohttp client.
    """

    #
    # Error code for BridgeException translate to internal service failure
    error_code = 1001

    def __init__(self, message: str, bridge_name: str):
        super().__init__(
            message=message, error_code=self.error_code, bridge_name=bridge_name
        )


class BridgeInternalFailureException(BridgeException):
    """
    BridgeInternalFailureException is wrapper for all internal exception raised by internal service.
    """

    #
    # Error code for BridgeException translate to internal service failure
    error_code = 1002

    def __init__(self, message: str, bridge_name: str):
        super().__init__(
            message=message, error_code=self.error_code, bridge_name=bridge_name
        )
