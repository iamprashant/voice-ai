"""
Copyright (c) 2023-2025 RapidaAI
Author: Prashant Srivastav <prashant@rapida.ai>

Licensed under GPL-2.0 with Rapida Additional Terms.
See LICENSE.md for details or contact sales@rapida.ai for commercial use.
"""
from typing import Union

from starlette.requests import Request

from app.config import ApplicationSettings
from app.exceptions.connector_exception import (
    ConnectorNotThereException,
)
from app.storage.storage import Storage


def attach_storage(setting: ApplicationSettings) -> Union[Storage, None]:
    if setting.storage:
        return Storage(setting.storage)


async def get_me_storage(request) -> Storage:
    """
    Return elastic search connection wrapper class from request context
    :param request: request context
    :return: :class:`ElasticSearchConnector`.
    """
    key = "storage"
    try:
        if isinstance(request, Request):
            return request.state.datasource[key]
        return request.state["datasource"][key]
    except KeyError:
        raise ConnectorNotThereException(key, f"{key} is not enable in env.")


async def get_storage(request: Request) -> Storage:
    """
    Return elastic search connection wrapper class from request context
    :param request: request context
    :return: :class:`ElasticSearchConnector`.
    """
    return await get_me_storage(request)
