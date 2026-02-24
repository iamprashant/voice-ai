"""
Tests for vault_bridge.py

Covers:
- get_credential uses correct RPC attr
- auth token placed in x-internal-service-key metadata
- successful response returns VaultCredential
- failed response raises BridgeException
"""
import pytest
from unittest.mock import AsyncMock

from app.bridges.artifacts.protos.vault_api_pb2 import VaultCredential
from app.bridges.internals.vault_bridge import VaultBridge
from app.exceptions.bridges_exceptions import BridgeException

SUCCESS_RESPONSE = {"success": True, "data": {"id": 42, "name": "my-key", "value": {}}}
FAIL_RESPONSE = {"success": False}


def make_bridge() -> VaultBridge:
    bridge = VaultBridge.__new__(VaultBridge)
    bridge.service_url = "localhost:9001"
    return bridge


class TestVaultBridgeGetCredential:

    @pytest.mark.asyncio
    async def test_calls_get_credential_rpc(self):
        bridge = make_bridge()
        mock_fetch = AsyncMock(return_value=SUCCESS_RESPONSE)
        bridge.fetch = mock_fetch

        await bridge.get_credential(auth_token="token", crendential_id=42)

        assert mock_fetch.call_args.kwargs["attr"] == "GetCredential"

    @pytest.mark.asyncio
    async def test_auth_token_in_internal_service_key_header(self):
        bridge = make_bridge()
        mock_fetch = AsyncMock(return_value=SUCCESS_RESPONSE)
        bridge.fetch = mock_fetch

        await bridge.get_credential(auth_token="my-jwt-token", crendential_id=1)

        metadata = {k: v for k, v in mock_fetch.call_args.kwargs["metadata"]}
        assert metadata.get("x-internal-service-key") == "my-jwt-token"

    @pytest.mark.asyncio
    async def test_returns_vault_credential_on_success(self):
        bridge = make_bridge()
        bridge.fetch = AsyncMock(return_value=SUCCESS_RESPONSE)

        result = await bridge.get_credential("token", 42)

        assert isinstance(result, VaultCredential)

    @pytest.mark.asyncio
    async def test_failed_response_raises_bridge_exception(self):
        bridge = make_bridge()
        bridge.fetch = AsyncMock(return_value=FAIL_RESPONSE)

        with pytest.raises(BridgeException) as exc_info:
            await bridge.get_credential("token", 1)

        assert exc_info.value.bridge_name == "vault"

    @pytest.mark.asyncio
    async def test_vault_id_set_on_request(self):
        bridge = make_bridge()
        mock_fetch = AsyncMock(return_value=SUCCESS_RESPONSE)
        bridge.fetch = mock_fetch

        await bridge.get_credential("token", 99)

        message = mock_fetch.call_args.kwargs["message_type"]
        assert message.vaultId == 99

    @pytest.mark.asyncio
    async def test_different_credential_ids(self):
        for cred_id in [1, 100, 99999]:
            bridge = make_bridge()
            mock_fetch = AsyncMock(return_value=SUCCESS_RESPONSE)
            bridge.fetch = mock_fetch

            await bridge.get_credential("tok", cred_id)

            msg = mock_fetch.call_args.kwargs["message_type"]
            assert msg.vaultId == cred_id
