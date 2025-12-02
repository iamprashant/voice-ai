"""
Copyright (c) 2023-2025 RapidaAI
Author: Prashant Srivastav <prashant@rapida.ai>

Licensed under GPL-2.0 with Rapida Additional Terms.
See LICENSE.md for details or contact sales@rapida.ai for commercial use.
"""

from pydantic import BaseModel


class InternalServiceConfig(BaseModel):
    web_host: str
    integration_host: str
    endpoint_host: str
    assistant_host: str

    class Config:
        # For secret key
        env_file_encoding = "utf-8"
