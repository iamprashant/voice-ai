"""
Copyright (c) 2023-2025 RapidaAI
Author: Prashant Srivastav <prashant@rapida.ai>

Licensed under GPL-2.0 with Rapida Additional Terms.
See LICENSE.md for details or contact sales@rapida.ai for commercial use.
"""
from typing import Dict

from fastapi import APIRouter, Depends, Request

from app.connectors import Connector
from app.connectors.connector_factory import get_all_connectors

H1 = APIRouter()


@H1.get("/healthz/")
async def health(request: Request):
    """Health check enabled"""
    return {"healthy": True}


@H1.get("/readiness/")
async def readiness(request: Request, connectors: Dict = Depends(get_all_connectors)):
    """rediness enabled"""
    connections_status: Dict = {}
    for key in connectors:
        connection: Connector = connectors[key]
        connections_status[key] = {"is_connected": await connection.is_connected()}
    return connections_status
