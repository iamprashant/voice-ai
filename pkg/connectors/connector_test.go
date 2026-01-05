// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package connectors

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVectorSearchOptions_ToJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    VectorSearchOptions
		expected string
		hasError bool
	}{
		{
			name: "valid options",
			input: VectorSearchOptions{
				Alpha:    0.5,
				TopK:     10,
				MinScore: 0.8,
				Source:   []string{"text", "metadata"},
			},
			expected: `{"Alpha":0.5,"TopK":10,"MinScore":0.8,"Source":["text","metadata"]}`,
			hasError: false,
		},
		{
			name:     "empty options",
			input:    VectorSearchOptions{},
			expected: `{"Alpha":0,"TopK":0,"MinScore":0,"Source":null}`,
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.input.ToJSON()
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Parse both to maps for comparison to avoid order issues
				var expectedMap, resultMap map[string]interface{}
				assert.NoError(t, json.Unmarshal([]byte(tt.expected), &expectedMap))
				assert.NoError(t, json.Unmarshal([]byte(result), &resultMap))
				assert.Equal(t, expectedMap, resultMap)
			}
		})
	}
}

func TestNewDefaultVectorSearchOptions(t *testing.T) {
	tests := []struct {
		name     string
		options  []SearchOptions
		expected VectorSearchOptions
	}{
		{
			name:    "no options",
			options: []SearchOptions{},
			expected: VectorSearchOptions{
				TopK:     5,
				MinScore: 0.5,
				Source:   []string{"text", "metadata"},
			},
		},
		{
			name: "with topK option",
			options: []SearchOptions{
				WithTopK(20),
			},
			expected: VectorSearchOptions{
				TopK:     20,
				MinScore: 0.5,
				Source:   []string{"text", "metadata"},
			},
		},
		{
			name: "with multiple options",
			options: []SearchOptions{
				WithTopK(15),
				WithAlpha(0.7),
				WithMinScore(0.9),
				WithSource([]string{"title", "content"}),
			},
			expected: VectorSearchOptions{
				Alpha:    0.7,
				TopK:     15,
				MinScore: 0.9,
				Source:   []string{"title", "content"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewDefaultVectorSearchOptions(tt.options...)
			assert.Equal(t, tt.expected, *result)
		})
	}
}

func TestVectorSearchOptions_WithOptions(t *testing.T) {
	opts := &VectorSearchOptions{
		TopK:     5,
		MinScore: 0.5,
		Source:   []string{"text"},
	}

	result := opts.WithOptions(WithTopK(10), WithMinScore(0.8))

	assert.Equal(t, 10, result.TopK)
	assert.Equal(t, float32(0.8), result.MinScore)
	assert.Equal(t, []string{"text"}, result.Source)
	assert.Equal(t, opts, result) // Should return the same pointer
}

func TestWithTopK(t *testing.T) {
	opts := &VectorSearchOptions{}
	WithTopK(25)(opts)
	assert.Equal(t, 25, opts.TopK)
}

func TestWithAlpha(t *testing.T) {
	opts := &VectorSearchOptions{}
	WithAlpha(0.3)(opts)
	assert.Equal(t, float32(0.3), opts.Alpha)
}

func TestWithMinScore(t *testing.T) {
	opts := &VectorSearchOptions{}
	WithMinScore(0.6)(opts)
	assert.Equal(t, float32(0.6), opts.MinScore)
}

func TestWithSource(t *testing.T) {
	opts := &VectorSearchOptions{}
	WithSource([]string{"field1", "field2"})(opts)
	assert.Equal(t, []string{"field1", "field2"}, opts.Source)
}

func TestVectorSearchOptions_JSONMarshaling(t *testing.T) {
	opts := VectorSearchOptions{
		Alpha:    0.5,
		TopK:     10,
		MinScore: 0.8,
		Source:   []string{"text", "metadata"},
	}

	data, err := json.Marshal(opts)
	assert.NoError(t, err)

	var unmarshaled VectorSearchOptions
	err = json.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, opts, unmarshaled)
}

func TestVectorSearchOptions_EdgeCases(t *testing.T) {
	t.Run("zero values", func(t *testing.T) {
		opts := VectorSearchOptions{}
		jsonStr, err := opts.ToJSON()
		assert.NoError(t, err)
		assert.Contains(t, jsonStr, `"TopK":0`)
		assert.Contains(t, jsonStr, `"MinScore":0`)
		assert.Contains(t, jsonStr, `"Alpha":0`)
	})

	t.Run("large numbers", func(t *testing.T) {
		opts := VectorSearchOptions{
			TopK:     10000,
			MinScore: 0.999,
			Alpha:    1.0,
		}
		jsonStr, err := opts.ToJSON()
		assert.NoError(t, err)
		assert.Contains(t, jsonStr, `"TopK":10000`)
		assert.Contains(t, jsonStr, `"MinScore":0.999`)
		assert.Contains(t, jsonStr, `"Alpha":1`)
	})

	t.Run("empty source", func(t *testing.T) {
		opts := VectorSearchOptions{
			Source: []string{},
		}
		jsonStr, err := opts.ToJSON()
		assert.NoError(t, err)
		assert.Contains(t, jsonStr, `"Source":[]`)
	})

	t.Run("nil source", func(t *testing.T) {
		opts := VectorSearchOptions{
			Source: nil,
		}
		jsonStr, err := opts.ToJSON()
		assert.NoError(t, err)
		assert.Contains(t, jsonStr, `"Source":null`)
	})
}
