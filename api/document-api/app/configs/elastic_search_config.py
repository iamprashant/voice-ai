"""
Copyright (c) 2023-2025 RapidaAI
Author: Prashant Srivastav <prashant@rapida.ai>

Licensed under GPL-2.0 with Rapida Additional Terms.
See LICENSE.md for details or contact sales@rapida.ai for commercial use.
"""

from typing import Optional, Union

from pydantic import Field

from app.configs import ExternalDatasourceModel
from app.configs.auth.aws_auth import AWSAuth
from app.configs.auth.basic_auth import BasicAuth


class ElasticSearchConfig(ExternalDatasourceModel):
    """
    Elastic search configuration template
    """

    # Host
    host: str

    # Port of running elastic search
    port: Optional[int]

    # authentication of elastic search node
    auth: Optional[Union[BasicAuth, AWSAuth]] = Field(
        default=None, description="auth information for elastic search config"
    )

    # default schema is https can be override from env
    scheme: str = "https"

    # max number of connections
    max_connection: int = 5

    class Config:
        env_nested_delimiter = "__"
        env_file_encoding = "utf-8"
