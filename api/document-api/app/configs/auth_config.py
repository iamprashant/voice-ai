"""
Copyright (c) 2023-2025 RapidaAI
Author: Prashant Srivastav <prashant@rapida.ai>

Licensed under GPL-2.0 with Rapida Additional Terms.
See LICENSE.md for details or contact sales@rapida.ai for commercial use.
"""

from typing import Set, Optional, Union

from pydantic import BaseModel, SecretStr, Field


class TokenAuthenticationConfig(BaseModel):
    # strict authentication
    strict: bool = False

    header_key: str = "Authorization"

    auth_host: Optional[str] = Field(
        default=None, description="auth host that will get called when authentication"
    )


class JwtAuthenticationConfig(BaseModel):
    """
    JWT configuration Template
    """

    strict: bool = False
    # secret key
    # print safe
    secret_key: SecretStr
    # algorithms for jwt
    algorithms: Set = ["HS256"]

    # header keys
    header_key: str = "Authorization"

    class Config:
        case_sensitive = True
        env_file_encoding = "utf-8"


class AuthenticationConfig(BaseModel):
    #
    type: Optional[str]

    config: Optional[Union[TokenAuthenticationConfig, JwtAuthenticationConfig]]
