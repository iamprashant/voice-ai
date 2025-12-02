"""
Copyright (c) 2023-2025 RapidaAI
Author: Prashant Srivastav <prashant@rapida.ai>

Licensed under GPL-2.0 with Rapida Additional Terms.
See LICENSE.md for details or contact sales@rapida.ai for commercial use.
"""

from typing import Optional

from pydantic import Field
from pydantic_settings import SettingsConfigDict

from app.configs import ExternalDatasourceModel
from app.configs.auth.aws_auth import AWSAuth


class AssetStoreConfig(ExternalDatasourceModel):
    storage_type: Optional[str]
    storage_path_prefix: Optional[str]
    auth: Optional[AWSAuth] = Field(
        default=None, description="auth information for storage config"
    )
    model_config = SettingsConfigDict(env_file_encoding="utf-8", extra="ignore")
