"""
Tests for bridge_factory.py

Covers:
- get_me_user_service_client removed (P4 dead code)
- AuthBridge import removed
- get_me_integration_client / get_me_vault_service_client accept URL strings
- ValueError on empty / None URL
"""
import pytest

import app.bridges.bridge_factory as factory
from app.bridges.bridge_factory import (
    get_me_integration_client,
    get_me_vault_service_client,
)
from app.bridges.internals.integration_bridge import IntegrationBridge
from app.bridges.internals.vault_bridge import VaultBridge


class TestDeadCodeRemoved:

    def test_get_me_user_service_client_removed(self):
        """P4: get_me_user_service_client referenced settings.lomotif_user_service_url
        which doesn't exist; the function must be deleted."""
        assert not hasattr(factory, "get_me_user_service_client"), (
            "get_me_user_service_client is dead code and must be removed from bridge_factory"
        )

    def test_auth_bridge_not_imported_at_module_level(self):
        """After removing get_me_user_service_client, AuthBridge import is unused."""
        assert not hasattr(factory, "AuthBridge"), (
            "AuthBridge was imported only for get_me_user_service_client; "
            "it must be removed with that function"
        )


class TestGetMeIntegrationClient:

    def test_returns_integration_bridge_instance(self):
        client = get_me_integration_client("localhost:9004")
        assert isinstance(client, IntegrationBridge)

    def test_sets_service_url(self):
        client = get_me_integration_client("my-host:9004")
        assert client.service_url == "my-host:9004"

    def test_raises_on_empty_string(self):
        with pytest.raises(ValueError):
            get_me_integration_client("")

    def test_raises_on_none(self):
        with pytest.raises(ValueError):
            get_me_integration_client(None)


class TestGetMeVaultServiceClient:

    def test_returns_vault_bridge_instance(self):
        client = get_me_vault_service_client("localhost:9001")
        assert isinstance(client, VaultBridge)

    def test_sets_service_url(self):
        client = get_me_vault_service_client("web-host:9001")
        assert client.service_url == "web-host:9001"

    def test_raises_on_empty_string(self):
        with pytest.raises(ValueError):
            get_me_vault_service_client("")

    def test_raises_on_none(self):
        with pytest.raises(ValueError):
            get_me_vault_service_client(None)


class TestFactoryFunctionsExist:
    """Verify the expected public API of bridge_factory is intact."""

    def test_get_me_integration_client_exists(self):
        assert hasattr(factory, "get_me_integration_client")

    def test_get_me_vault_service_client_exists(self):
        assert hasattr(factory, "get_me_vault_service_client")

    def test_get_vault_service_client_fastapi_dep_exists(self):
        assert hasattr(factory, "get_vault_service_client")

    def test_get_integration_client_fastapi_dep_exists(self):
        assert hasattr(factory, "get_integration_client")
