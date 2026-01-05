// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package connectors

import (
	"context"
	"testing"

	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/rapidaai/pkg/commons"
	configs "github.com/rapidaai/pkg/configs"
	"github.com/stretchr/testify/assert"
)

// Local config struct for testing - removed since we can use configs package

func TestNewOpenSearchConnector(t *testing.T) {
	// Requires external config and logger packages
	t.Skip("Requires external packages - integration test")
}

func TestOpenSearchConnector_Name(t *testing.T) {
	tests := []struct {
		name     string
		schema   string
		host     string
		expected string
	}{
		{
			name:     "http schema",
			schema:   "http",
			host:     "localhost",
			expected: "ES http://localhost",
		},
		{
			name:     "https schema",
			schema:   "https",
			host:     "opensearch.example.com",
			expected: "ES https://opensearch.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			connector := &openSearchConnector{
				cfg: &configs.OpenSearchConfig{Schema: tt.schema, Host: tt.host},
			}
			result := connector.Name()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestOpenSearchConnector_IsConnected(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()
	connector := &openSearchConnector{
		logger: logger,
	}

	ctx := context.Background()
	// Without connection, this would fail, but implementation tries to call Info
	// For unit test, we test the logic path
	t.Run("without connection", func(t *testing.T) {
		connector.Connection = nil
		result := connector.IsConnected(ctx)
		assert.False(t, result)
	})

	t.Run("with mock connection", func(t *testing.T) {
		// Mock client would be needed, but for simplicity, test the return path
		t.Skip("Requires mock OpenSearch client - integration test")
	})
}

func TestOpenSearchConnector_Disconnect(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()
	mockClient := &opensearch.Client{}
	connector := &openSearchConnector{
		Connection: mockClient,
		logger:     logger,
	}

	ctx := context.Background()
	err := connector.Disconnect(ctx)

	assert.NoError(t, err)
	assert.Nil(t, connector.Connection)
}

// Note: Search methods (VectorSearch, HybridSearch, TextSearch, Search, SearchWithCount)
// require a real OpenSearch connection and are tested as integration tests.
// Persist, Update, Bulk also require connection.

func TestOpenSearchConnector_Connect_ErrorHandling(t *testing.T) {
	// Test with invalid config
	t.Skip("Requires external packages - integration test")
}

func TestOpenSearchConnector_EdgeCases(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()
	t.Run("disconnect with nil connection", func(t *testing.T) {
		port := 9200
		cfg := &configs.OpenSearchConfig{
			Schema: "http",
			Host:   "localhost",
			Port:   &port,
		}
		connector := &openSearchConnector{
			cfg:        cfg,
			logger:     logger,
			Connection: nil,
		}

		ctx := context.Background()
		err := connector.Disconnect(ctx)
		assert.NoError(t, err)
		assert.Nil(t, connector.Connection)
	})

	t.Run("name with nil port", func(t *testing.T) {
		cfg := &configs.OpenSearchConfig{
			Schema: "http",
			Host:   "localhost",
			Port:   nil,
		}
		connector := &openSearchConnector{cfg: cfg}
		result := connector.Name()
		assert.Equal(t, "ES http://localhost", result)
	})
}

// Test SearchResponse methods
func TestSearchResponse_Error(t *testing.T) {
	tests := []struct {
		name     string
		response SearchResponse
		expected error
	}{
		{
			name: "with error",
			response: SearchResponse{
				OpenSearchResponse: OpenSearchResponse{
					Err: assert.AnError,
				},
			},
			expected: assert.AnError,
		},
		{
			name: "without error",
			response: SearchResponse{
				OpenSearchResponse: OpenSearchResponse{
					Err: nil,
				},
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.response.Error()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSearchResponseWithCount_Error(t *testing.T) {
	response := SearchResponseWithCount{
		OpenSearchResponse: OpenSearchResponse{
			Err: assert.AnError,
		},
	}
	result := response.Error()
	assert.Equal(t, assert.AnError, result)
}
