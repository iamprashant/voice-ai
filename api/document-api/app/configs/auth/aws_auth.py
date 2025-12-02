"""
Copyright (c) 2023-2025 RapidaAI
Author: Prashant Srivastav <prashant@rapida.ai>

Licensed under GPL-2.0 with Rapida Additional Terms.
See LICENSE.md for details or contact sales@rapida.ai for commercial use.
"""
from typing import Optional

from pydantic import BaseModel, SecretStr, Field


class AWSAuth(BaseModel):
    # aws region
    region: str

    # aws access key
    access_key_id: Optional[SecretStr]

    # aws secret key
    secret_access_key: Optional[SecretStr]

    # if sts get used
    assume_role: Optional[str] = Field(default=None)

    class Config:
        case_sensitive = True
        env_file_encoding = "utf-8"
