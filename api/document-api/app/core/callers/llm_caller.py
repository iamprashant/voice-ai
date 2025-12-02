"""
Copyright (c) 2023-2025 RapidaAI
Author: Prashant Srivastav <prashant@rapida.ai>

Licensed under GPL-2.0 with Rapida Additional Terms.
See LICENSE.md for details or contact sales@rapida.ai for commercial use.
"""

import logging
from typing import Dict

from app.bridges.artifacts.protos.integration_api_pb2 import EmbeddingRequest
from app.bridges.bridge_factory import get_me_integration_client
from app.bridges.internals.integration_bridge import IntegrationBridge
from app.configs.internal_service_config import InternalServiceConfig
from app.exceptions import RapidaException

_log = logging.getLogger("app.core.callers.llm_caller")


class LLMCaller:
    integration_client: IntegrationBridge

    def __init__(self, cfg: InternalServiceConfig):
        self.integration_client = get_me_integration_client(cfg.integration_host)

    async def get_embedding(
        self, auth: str, provider_name: str, request: EmbeddingRequest
    ) -> Dict:
        try:
            return await self.integration_client.get_embedding(
                auth, provider_name, request
            )
        except RapidaException as err:
            _log.debug(f"Error while creating embedding from LLM {err}")
            raise err
