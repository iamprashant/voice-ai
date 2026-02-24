"""
Tests for integration_bridge.py

Covers:
- Provider stub mapping (google → GeminiServiceStub, vertex-ai → VertexAiServiceStub)
- Regression: GoogleServiceStub must not be used
- Unsupported provider raises BridgeException
- RPC attr, metadata, and response handling
"""
import pytest
from unittest.mock import AsyncMock

from app.bridges.artifacts.protos import integration_api_pb2_grpc
from app.bridges.artifacts.protos.integration_api_pb2 import EmbeddingRequest
from app.bridges.internals.integration_bridge import IntegrationBridge
from app.exceptions.bridges_exceptions import BridgeException

VALID_RESPONSE = {"success": True, "data": [], "metrics": []}
FAIL_RESPONSE = {"success": False, "data": [], "metrics": []}


def make_bridge() -> IntegrationBridge:
    bridge = IntegrationBridge.__new__(IntegrationBridge)
    bridge.service_url = "localhost:9004"
    return bridge


class TestProviderStubMapping:
    """Each provider name must map to the correct gRPC stub class."""

    @pytest.mark.asyncio
    @pytest.mark.parametrize(
        "provider_name, expected_stub",
        [
            ("openai", integration_api_pb2_grpc.OpenAiServiceStub),
            ("cohere", integration_api_pb2_grpc.CohereServiceStub),
            ("voyageai", integration_api_pb2_grpc.VoyageAiServiceStub),
            ("bedrock", integration_api_pb2_grpc.BedrockServiceStub),
            ("azure-openai", integration_api_pb2_grpc.AzureServiceStub),
            ("google", integration_api_pb2_grpc.GeminiServiceStub),
            ("vertex-ai", integration_api_pb2_grpc.VertexAiServiceStub),
        ],
    )
    async def test_provider_uses_correct_stub(self, provider_name, expected_stub):
        bridge = make_bridge()
        mock_fetch = AsyncMock(return_value=VALID_RESPONSE)
        bridge.fetch = mock_fetch

        await bridge.get_embedding(
            auth_token="tok",
            provider_name=provider_name,
            request=EmbeddingRequest(),
        )

        stub_used = mock_fetch.call_args.kwargs["stub"]
        assert stub_used is expected_stub, (
            f"Provider '{provider_name}' should map to {expected_stub.__name__}, "
            f"but got {stub_used.__name__}"
        )

    @pytest.mark.asyncio
    async def test_google_uses_gemini_not_google_stub(self):
        """Regression: GoogleServiceStub was removed; 'google' must route to GeminiServiceStub."""
        bridge = make_bridge()
        mock_fetch = AsyncMock(return_value=VALID_RESPONSE)
        bridge.fetch = mock_fetch

        await bridge.get_embedding("tok", "google", EmbeddingRequest())

        stub_used = mock_fetch.call_args.kwargs["stub"]
        assert stub_used is integration_api_pb2_grpc.GeminiServiceStub
        # Confirm the old stub does not exist
        old_stub = getattr(integration_api_pb2_grpc, "GoogleServiceStub", None)
        assert stub_used is not old_stub

    @pytest.mark.asyncio
    async def test_unsupported_provider_raises_bridge_exception(self):
        bridge = make_bridge()
        mock_fetch = AsyncMock(return_value=VALID_RESPONSE)
        bridge.fetch = mock_fetch

        with pytest.raises(BridgeException) as exc_info:
            await bridge.get_embedding("tok", "unknown-llm", EmbeddingRequest())

        assert exc_info.value.bridge_name == "integration"
        mock_fetch.assert_not_called()

    @pytest.mark.asyncio
    async def test_empty_provider_raises_bridge_exception(self):
        bridge = make_bridge()
        bridge.fetch = AsyncMock(return_value=VALID_RESPONSE)

        with pytest.raises(BridgeException):
            await bridge.get_embedding("tok", "", EmbeddingRequest())


class TestIntegrationBridgeRpcBehaviour:

    @pytest.mark.asyncio
    async def test_rpc_attr_is_always_embedding(self):
        bridge = make_bridge()
        mock_fetch = AsyncMock(return_value=VALID_RESPONSE)
        bridge.fetch = mock_fetch

        await bridge.get_embedding("tok", "openai", EmbeddingRequest())

        assert mock_fetch.call_args.kwargs["attr"] == "Embedding"

    @pytest.mark.asyncio
    async def test_auth_token_placed_in_metadata(self):
        bridge = make_bridge()
        mock_fetch = AsyncMock(return_value=VALID_RESPONSE)
        bridge.fetch = mock_fetch

        await bridge.get_embedding("secret-token", "openai", EmbeddingRequest())

        metadata = {k: v for k, v in mock_fetch.call_args.kwargs["metadata"]}
        assert metadata.get("x-internal-service-key") == "secret-token"

    @pytest.mark.asyncio
    async def test_successful_response_returns_data_and_metrics_tuple(self):
        bridge = make_bridge()
        bridge.fetch = AsyncMock(return_value=VALID_RESPONSE)

        result = await bridge.get_embedding("tok", "openai", EmbeddingRequest())

        assert isinstance(result, tuple)
        assert len(result) == 2
        data, metrics = result
        # Proto repeated fields are containers; verify they are not scalars
        assert data is not None
        assert metrics is not None

    @pytest.mark.asyncio
    async def test_failed_response_raises_bridge_exception(self):
        bridge = make_bridge()
        bridge.fetch = AsyncMock(return_value=FAIL_RESPONSE)

        with pytest.raises(BridgeException) as exc_info:
            await bridge.get_embedding("tok", "openai", EmbeddingRequest())

        assert exc_info.value.bridge_name == "integration"

    @pytest.mark.asyncio
    async def test_request_passed_as_message_type(self):
        bridge = make_bridge()
        mock_fetch = AsyncMock(return_value=VALID_RESPONSE)
        bridge.fetch = mock_fetch

        req = EmbeddingRequest()
        await bridge.get_embedding("tok", "openai", req)

        assert mock_fetch.call_args.kwargs["message_type"] is req


class TestProtoModuleIntegrity:

    def test_google_service_stub_removed_from_proto(self):
        """GoogleServiceStub was renamed to GeminiServiceStub in the new proto version."""
        assert not hasattr(integration_api_pb2_grpc, "GoogleServiceStub"), (
            "GoogleServiceStub still exists; it should have been removed "
            "when proto was upgraded to use GeminiServiceStub"
        )

    def test_gemini_service_stub_exists(self):
        assert hasattr(integration_api_pb2_grpc, "GeminiServiceStub")

    def test_vertex_ai_service_stub_exists(self):
        assert hasattr(integration_api_pb2_grpc, "VertexAiServiceStub")

    def test_all_expected_stubs_present(self):
        required = [
            "OpenAiServiceStub",
            "CohereServiceStub",
            "VoyageAiServiceStub",
            "BedrockServiceStub",
            "AzureServiceStub",
            "GeminiServiceStub",
            "VertexAiServiceStub",
        ]
        for stub_name in required:
            assert hasattr(integration_api_pb2_grpc, stub_name), (
                f"{stub_name} missing from integration_api_pb2_grpc"
            )
