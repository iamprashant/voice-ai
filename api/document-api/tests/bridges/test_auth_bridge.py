"""
Tests for auth_bridge.py

Covers:
- scope_authorize uses 'ScopeAuthorize' RPC (not the old 'Invoke')
- scope_authorize uses ScopeAuthorizeRequest (not AuthorizeRequest)
- scope_authorize returns ScopedAuthenticationResponse (not AuthenticateResponse)
- x-api-key metadata header
- authorize still works correctly
"""
import pytest
from unittest.mock import AsyncMock

from app.bridges.artifacts.protos import web_api_pb2
from app.bridges.internals.auth_bridge import AuthBridge
from app.exceptions.bridges_exceptions import BridgeException


def make_bridge() -> AuthBridge:
    bridge = AuthBridge.__new__(AuthBridge)
    bridge.service_url = "localhost:9001"
    return bridge


class TestScopeAuthorize:
    """Tests for the scope_authorize method after proto upgrade."""

    @pytest.mark.asyncio
    async def test_uses_scope_authorize_rpc_not_invoke(self):
        """Regression: old code called attr='Invoke'; must now be 'ScopeAuthorize'."""
        bridge = make_bridge()
        mock_fetch = AsyncMock(return_value={"success": True, "code": 200, "data": {}})
        bridge.fetch = mock_fetch

        await bridge.scope_authorize("api-key")

        attr = mock_fetch.call_args.kwargs["attr"]
        assert attr == "ScopeAuthorize", (
            f"Expected attr='ScopeAuthorize', got '{attr}'. "
            "The old broken value was 'Invoke'."
        )

    @pytest.mark.asyncio
    async def test_does_not_use_invoke_rpc(self):
        bridge = make_bridge()
        bridge.fetch = AsyncMock(return_value={"success": True, "code": 200, "data": {}})

        await bridge.scope_authorize("api-key")

        assert mock_fetch.call_args.kwargs["attr"] != "Invoke" \
            if (mock_fetch := bridge.fetch) else True

    @pytest.mark.asyncio
    async def test_uses_scope_authorize_request_message(self):
        """Regression: old code used AuthorizeRequest; must now use ScopeAuthorizeRequest."""
        bridge = make_bridge()
        mock_fetch = AsyncMock(return_value={"success": True, "code": 200, "data": {}})
        bridge.fetch = mock_fetch

        await bridge.scope_authorize("api-key")

        message_type = mock_fetch.call_args.kwargs["message_type"]
        assert isinstance(message_type, web_api_pb2.ScopeAuthorizeRequest), (
            f"Expected ScopeAuthorizeRequest, got {type(message_type).__name__}. "
            "The old broken value was AuthorizeRequest."
        )

    @pytest.mark.asyncio
    async def test_does_not_use_authorize_request(self):
        bridge = make_bridge()
        mock_fetch = AsyncMock(return_value={"success": True, "code": 200, "data": {}})
        bridge.fetch = mock_fetch

        await bridge.scope_authorize("api-key")

        message_type = mock_fetch.call_args.kwargs["message_type"]
        assert not isinstance(message_type, web_api_pb2.AuthorizeRequest), (
            "scope_authorize must not use AuthorizeRequest (old broken type)"
        )

    @pytest.mark.asyncio
    async def test_returns_scoped_authentication_response(self):
        """Regression: old code returned AuthenticateResponse; must now return ScopedAuthenticationResponse."""
        bridge = make_bridge()
        bridge.fetch = AsyncMock(return_value={"success": True, "code": 200, "data": {}})

        result = await bridge.scope_authorize("api-key")

        assert isinstance(result, web_api_pb2.ScopedAuthenticationResponse), (
            f"Expected ScopedAuthenticationResponse, got {type(result).__name__}. "
            "The old broken return type was AuthenticateResponse."
        )

    @pytest.mark.asyncio
    async def test_does_not_return_authenticate_response(self):
        bridge = make_bridge()
        bridge.fetch = AsyncMock(return_value={"success": True, "code": 200, "data": {}})

        result = await bridge.scope_authorize("api-key")

        assert not isinstance(result, web_api_pb2.AuthenticateResponse), (
            "scope_authorize must not return AuthenticateResponse (old broken type)"
        )

    @pytest.mark.asyncio
    async def test_api_key_set_in_x_api_key_header(self):
        bridge = make_bridge()
        mock_fetch = AsyncMock(return_value={"success": True, "code": 200, "data": {}})
        bridge.fetch = mock_fetch

        await bridge.scope_authorize("my-secret-api-key")

        metadata = {k: v for k, v in mock_fetch.call_args.kwargs["metadata"]}
        assert metadata.get("x-api-key") == "my-secret-api-key"

    @pytest.mark.asyncio
    async def test_preserving_proto_field_name_is_true(self):
        bridge = make_bridge()
        mock_fetch = AsyncMock(return_value={"success": True, "code": 200, "data": {}})
        bridge.fetch = mock_fetch

        await bridge.scope_authorize("tok")

        assert mock_fetch.call_args.kwargs.get("preserving_proto_field_name") is True


class TestAuthorize:
    """Tests for the existing authorize method (must remain unaffected)."""

    @pytest.mark.asyncio
    async def test_authorize_uses_authorize_rpc(self):
        bridge = make_bridge()
        mock_fetch = AsyncMock(
            return_value={"success": True, "code": 200, "data": {}}
        )
        bridge.fetch = mock_fetch

        await bridge.authorize("Bearer token", "user-123")

        assert mock_fetch.call_args.kwargs["attr"] == "Authorize"

    @pytest.mark.asyncio
    async def test_authorize_uses_authorize_request(self):
        bridge = make_bridge()
        mock_fetch = AsyncMock(
            return_value={"success": True, "code": 200, "data": {}}
        )
        bridge.fetch = mock_fetch

        await bridge.authorize("Bearer token", "user-123")

        msg = mock_fetch.call_args.kwargs["message_type"]
        assert isinstance(msg, web_api_pb2.AuthorizeRequest)

    @pytest.mark.asyncio
    async def test_authorize_sets_authorization_and_user_id_headers(self):
        bridge = make_bridge()
        mock_fetch = AsyncMock(
            return_value={"success": True, "code": 200, "data": {}}
        )
        bridge.fetch = mock_fetch

        await bridge.authorize("Bearer my-token", "uid-456")

        metadata = {k: v for k, v in mock_fetch.call_args.kwargs["metadata"]}
        assert metadata.get("authorization") == "Bearer my-token"
        assert metadata.get("x-auth-id") == "uid-456"

    @pytest.mark.asyncio
    async def test_authorize_failed_response_raises_bridge_exception(self):
        bridge = make_bridge()
        bridge.fetch = AsyncMock(return_value={"success": False})

        with pytest.raises(BridgeException) as exc_info:
            await bridge.authorize("Bearer token", "user-123")

        assert exc_info.value.bridge_name == "web"


class TestProtoMessageTypesAvailable:

    def test_scope_authorize_request_exists_in_proto(self):
        assert hasattr(web_api_pb2, "ScopeAuthorizeRequest")

    def test_scoped_authentication_response_exists_in_proto(self):
        assert hasattr(web_api_pb2, "ScopedAuthenticationResponse")

    def test_scope_authorize_request_instantiable(self):
        req = web_api_pb2.ScopeAuthorizeRequest()
        assert req is not None

    def test_scoped_authentication_response_instantiable(self):
        resp = web_api_pb2.ScopedAuthenticationResponse()
        assert resp is not None
