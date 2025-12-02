"""
Copyright (c) 2023-2025 RapidaAI
Author: Prashant Srivastav <prashant@rapida.ai>

Licensed under GPL-2.0 with Rapida Additional Terms.
See LICENSE.md for details or contact sales@rapida.ai for commercial use.
"""

from pydantic import BaseModel, SecretStr


class BasicAuth(BaseModel):
    """
    Basic Auth support of all connectors
     - username and password connect
    """

    #  username or user
    user: str

    # password string
    # print safe
    password: SecretStr = None
