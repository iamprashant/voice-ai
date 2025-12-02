"""
Copyright (c) 2023-2025 RapidaAI
Author: Prashant Srivastav <prashant@rapida.ai>

Licensed under GPL-2.0 with Rapida Additional Terms.
See LICENSE.md for details or contact sales@rapida.ai for commercial use.
"""
from typing import Optional

from pydantic import BaseModel


class ExternalDatasourceModel(BaseModel):
    """
    Datasource model provide scope to config which are more related to datasource
    """

    # do we need to place in docker container
    # or place as external service
    # default false
    dockerize: Optional[bool] = False
