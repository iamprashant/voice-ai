"""
Copyright (c) 2023-2025 RapidaAI
Author: Prashant Srivastav <prashant@rapida.ai>

Licensed under GPL-2.0 with Rapida Additional Terms.
See LICENSE.md for details or contact sales@rapida.ai for commercial use.
"""

from typing import Optional

from pydantic import Field

from app.configs import ExternalDatasourceModel
from app.configs.auth.basic_auth import BasicAuth


class RedisConfig(ExternalDatasourceModel):
    """
    Redis config template
    - defined all required or optional parameter which will be needed to create redis connection
    - organize with pydantic lib to get loaded from .env
    """

    # Host of redis instance
    host: str

    # redis port
    port: int

    # db of redis
    db: int = 0

    # maximum number of connection at given point of time.
    max_connection: int = 5

    # charset
    charset: str = "utf-8"

    # decode_responses
    decode_responses: bool = True

    # support only basic auth
    auth: Optional[BasicAuth] = Field(
        default=None, description="authentication information for redis"
    )

    # do we need to place in docker container
    # or place as external service
    dockerize: Optional[bool] = True

    class Config:
        env_nested_delimiter = "__"
        env_file_encoding = "utf-8"
