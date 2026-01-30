// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package integration_client_builders

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"

	gorm_types "github.com/rapidaai/pkg/models/gorm/types"
)

// Note: mockLogger and newTestLogger are defined in chat_input_builder_test.go

func TestNewEmbeddingInputBuilder(t *testing.T) {
	logger := newTestLogger()
	builder := NewEmbeddingInputBuilder(logger)

	assert.NotNil(t, builder, "NewEmbeddingInputBuilder should return a non-nil builder")
	assert.Implements(t, (*InputEmbeddingBuilder)(nil), builder, "builder should implement InputEmbeddingBuilder")
}

func TestEmbeddingInputBuilder_Credential(t *testing.T) {
	logger := newTestLogger()
	builder := NewEmbeddingInputBuilder(logger)

	tests := []struct {
		name     string
		id       uint64
		value    map[string]interface{}
		expected uint64
	}{
		{
			name:     "basic credential",
			id:       12345,
			value:    map[string]interface{}{"api_key": "test-key"},
			expected: 12345,
		},
		{
			name:     "zero id credential",
			id:       0,
			value:    map[string]interface{}{},
			expected: 0,
		},
		{
			name:     "large id credential",
			id:       18446744073709551615, // max uint64
			value:    map[string]interface{}{"key": "value"},
			expected: 18446744073709551615,
		},
		{
			name:     "credential with complex nested value",
			id:       100,
			value:    map[string]interface{}{"config": map[string]interface{}{"endpoint": "https://api.example.com", "timeout": 30}},
			expected: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			structVal, err := structpb.NewStruct(tt.value)
			require.NoError(t, err, "failed to create struct")

			cred := builder.Credential(tt.id, structVal)

			assert.NotNil(t, cred, "credential should not be nil")
			assert.Equal(t, tt.expected, cred.Id, "credential id should match")
			assert.NotNil(t, cred.Value, "credential value should not be nil")
		})
	}
}

func TestEmbeddingInputBuilder_Credential_NilValue(t *testing.T) {
	logger := newTestLogger()
	builder := NewEmbeddingInputBuilder(logger)

	cred := builder.Credential(456, nil)

	assert.NotNil(t, cred, "credential should not be nil")
	assert.Equal(t, uint64(456), cred.Id, "credential id should match")
	assert.Nil(t, cred.Value, "credential value should be nil")
}

func TestEmbeddingInputBuilder_Embedding(t *testing.T) {
	logger := newTestLogger()
	builder := NewEmbeddingInputBuilder(logger)

	t.Run("basic embedding request", func(t *testing.T) {
		structVal, _ := structpb.NewStruct(map[string]interface{}{"api_key": "test"})
		credential := builder.Credential(1, structVal)

		modelOpts := make(map[string]*anypb.Any)
		dimVal, _ := structpb.NewValue(1536)
		anyDim, _ := anypb.New(dimVal)
		modelOpts["dimensions"] = anyDim

		contents := map[int32]string{
			0: "Hello, world!",
			1: "How are you?",
		}

		additionalData := map[string]string{
			"trace_id": "trace-123",
		}

		request := builder.Embedding(credential, modelOpts, additionalData, contents)

		assert.NotNil(t, request, "request should not be nil")
		assert.NotNil(t, request.Credential, "credential should not be nil")
		assert.Equal(t, uint64(1), request.Credential.Id, "credential id should match")
		assert.NotNil(t, request.ModelParameters, "model parameters should not be nil")
		assert.Len(t, request.Content, 2, "should have two contents")
		assert.Equal(t, "Hello, world!", request.Content[0], "first content should match")
		assert.Equal(t, "How are you?", request.Content[1], "second content should match")
		assert.Equal(t, "trace-123", request.AdditionalData["trace_id"], "additional data should match")
	})

	t.Run("embedding request with nil model opts", func(t *testing.T) {
		structVal, _ := structpb.NewStruct(map[string]interface{}{"api_key": "test"})
		credential := builder.Credential(1, structVal)

		contents := map[int32]string{
			0: "Text to embed",
		}

		request := builder.Embedding(credential, nil, nil, contents)

		assert.NotNil(t, request, "request should not be nil")
		assert.Nil(t, request.ModelParameters, "model parameters should be nil")
		assert.Nil(t, request.AdditionalData, "additional data should be nil")
	})

	t.Run("embedding request with empty contents", func(t *testing.T) {
		structVal, _ := structpb.NewStruct(map[string]interface{}{"api_key": "test"})
		credential := builder.Credential(1, structVal)

		contents := map[int32]string{}

		request := builder.Embedding(credential, nil, nil, contents)

		assert.NotNil(t, request, "request should not be nil")
		assert.Empty(t, request.Content, "content should be empty")
	})

	t.Run("embedding request with large content batch", func(t *testing.T) {
		structVal, _ := structpb.NewStruct(map[string]interface{}{"api_key": "test"})
		credential := builder.Credential(1, structVal)

		contents := make(map[int32]string)
		for i := int32(0); i < 100; i++ {
			contents[i] = "Text content " + string(rune('0'+i%10))
		}

		request := builder.Embedding(credential, nil, nil, contents)

		assert.NotNil(t, request, "request should not be nil")
		assert.Len(t, request.Content, 100, "should have 100 contents")
	})

	t.Run("embedding request with sparse indices", func(t *testing.T) {
		structVal, _ := structpb.NewStruct(map[string]interface{}{"api_key": "test"})
		credential := builder.Credential(1, structVal)

		contents := map[int32]string{
			0:   "First",
			5:   "Fifth",
			100: "Hundredth",
		}

		request := builder.Embedding(credential, nil, nil, contents)

		assert.NotNil(t, request, "request should not be nil")
		assert.Len(t, request.Content, 3, "should have three contents")
		assert.Equal(t, "First", request.Content[0], "index 0 should match")
		assert.Equal(t, "Fifth", request.Content[5], "index 5 should match")
		assert.Equal(t, "Hundredth", request.Content[100], "index 100 should match")
	})

	t.Run("embedding request with unicode content", func(t *testing.T) {
		structVal, _ := structpb.NewStruct(map[string]interface{}{"api_key": "test"})
		credential := builder.Credential(1, structVal)

		contents := map[int32]string{
			0: "Hello 你好 مرحبا 🌍",
			1: "Émoji test: 🚀💻🎉",
		}

		request := builder.Embedding(credential, nil, nil, contents)

		assert.NotNil(t, request, "request should not be nil")
		assert.Equal(t, "Hello 你好 مرحبا 🌍", request.Content[0], "unicode content should be preserved")
		assert.Equal(t, "Émoji test: 🚀💻🎉", request.Content[1], "emoji content should be preserved")
	})
}

func TestEmbeddingInputBuilder_Arguments(t *testing.T) {
	logger := newTestLogger()
	builder := NewEmbeddingInputBuilder(logger).(*embeddingInputBuilder)

	t.Run("merge with existing variables", func(t *testing.T) {
		variables := []*gorm_types.PromptVariable{
			{Name: "model", DefaultValue: "text-embedding-ada-002"},
			{Name: "dimensions", DefaultValue: "1536"},
		}

		modelVal, _ := structpb.NewValue("text-embedding-3-large")
		modelAny, _ := anypb.New(modelVal)

		arguments := map[string]*anypb.Any{
			"model": modelAny,
		}

		result := builder.Arguments(variables, arguments)

		// The anypb value gets converted to a map with protobuf metadata
		assert.NotNil(t, result["model"], "model should be overridden")
		assert.Equal(t, "1536", result["dimensions"], "dimensions should use default")
	})

	t.Run("empty variables", func(t *testing.T) {
		variables := []*gorm_types.PromptVariable{}

		dimVal, _ := structpb.NewValue(768)
		dimAny, _ := anypb.New(dimVal)

		arguments := map[string]*anypb.Any{
			"dimensions": dimAny,
		}

		result := builder.Arguments(variables, arguments)

		assert.NotNil(t, result["dimensions"], "dimensions should be set from arguments")
	})

	t.Run("nil arguments", func(t *testing.T) {
		variables := []*gorm_types.PromptVariable{
			{Name: "model", DefaultValue: "default-model"},
		}

		result := builder.Arguments(variables, nil)

		assert.Equal(t, "default-model", result["model"], "should use default value")
	})

	t.Run("both empty", func(t *testing.T) {
		variables := []*gorm_types.PromptVariable{}
		arguments := map[string]*anypb.Any{}

		result := builder.Arguments(variables, arguments)

		assert.Empty(t, result, "result should be empty")
	})
}

func TestEmbeddingInputBuilder_Options(t *testing.T) {
	logger := newTestLogger()
	builder := NewEmbeddingInputBuilder(logger)

	t.Run("add options to nil map", func(t *testing.T) {
		opts := map[string]interface{}{
			"model":      "text-embedding-3-small",
			"dimensions": 1536,
		}

		result := builder.Options(opts, nil)

		assert.NotNil(t, result, "result should not be nil")
		assert.Len(t, result, 2, "should have two options")
	})

	t.Run("add options to existing map", func(t *testing.T) {
		existingVal, _ := structpb.NewValue("existing-model")
		existingAny, _ := anypb.New(existingVal)

		existingOpts := map[string]*anypb.Any{
			"model": existingAny,
		}

		opts := map[string]interface{}{
			"dimensions": 768,
		}

		result := builder.Options(opts, existingOpts)

		assert.Len(t, result, 2, "should have two options")
		assert.NotNil(t, result["model"], "existing key should be preserved")
		assert.NotNil(t, result["dimensions"], "new key should be added")
	})

	t.Run("override existing option", func(t *testing.T) {
		existingVal, _ := structpb.NewValue(1536)
		existingAny, _ := anypb.New(existingVal)

		existingOpts := map[string]*anypb.Any{
			"dimensions": existingAny,
		}

		opts := map[string]interface{}{
			"dimensions": 768,
		}

		result := builder.Options(opts, existingOpts)

		assert.Len(t, result, 1, "should have one option")
	})

	t.Run("empty opts preserves existing", func(t *testing.T) {
		existingVal, _ := structpb.NewValue("existing")
		existingAny, _ := anypb.New(existingVal)

		existingOpts := map[string]*anypb.Any{
			"model": existingAny,
		}

		opts := map[string]interface{}{}

		result := builder.Options(opts, existingOpts)

		assert.Len(t, result, 1, "should preserve existing options")
	})

	t.Run("various value types", func(t *testing.T) {
		opts := map[string]interface{}{
			"string_val": "text-embedding-3-small",
			"int_val":    1536,
			"float_val":  0.5,
			"bool_val":   true,
			"list_val":   []interface{}{"a", "b"},
		}

		result := builder.Options(opts, nil)

		assert.Len(t, result, 5, "should handle all value types")
	})

	t.Run("nil opts and nil existing", func(t *testing.T) {
		result := builder.Options(nil, nil)

		assert.NotNil(t, result, "result should not be nil")
		assert.Empty(t, result, "result should be empty")
	})
}

func TestEmbeddingInputBuilder_Integration(t *testing.T) {
	logger := newTestLogger()
	builder := NewEmbeddingInputBuilder(logger)

	t.Run("full embedding request workflow", func(t *testing.T) {
		// Create credential
		structVal, _ := structpb.NewStruct(map[string]interface{}{
			"api_key": "sk-test-key",
		})
		credential := builder.Credential(99, structVal)

		// Create model options
		modelOpts := builder.Options(map[string]interface{}{
			"model.name":       "text-embedding-3-large",
			"model.dimensions": 3072,
		}, nil)

		// Create contents
		contents := map[int32]string{
			0: "The quick brown fox jumps over the lazy dog.",
			1: "Machine learning is a subset of artificial intelligence.",
			2: "Natural language processing enables computers to understand human language.",
		}

		// Create additional data
		additionalData := map[string]string{
			"batch_id": "batch-001",
			"source":   "document-parser",
		}

		// Build the request
		request := builder.Embedding(credential, modelOpts, additionalData, contents)

		// Assertions
		assert.NotNil(t, request)
		assert.Equal(t, uint64(99), request.Credential.Id)
		assert.Len(t, request.Content, 3)
		assert.Equal(t, "batch-001", request.AdditionalData["batch_id"])
		assert.NotNil(t, request.ModelParameters["model.name"])
	})

	t.Run("embedding request with document chunks", func(t *testing.T) {
		structVal, _ := structpb.NewStruct(map[string]interface{}{"api_key": "test"})
		credential := builder.Credential(1, structVal)

		// Simulate document chunks
		chunks := []string{
			"Chapter 1: Introduction to AI. Artificial intelligence (AI) is the simulation of human intelligence processes by computer systems.",
			"Chapter 2: Machine Learning. Machine learning is a method of data analysis that automates analytical model building.",
			"Chapter 3: Deep Learning. Deep learning is a subset of machine learning based on artificial neural networks.",
		}

		contents := make(map[int32]string)
		for i, chunk := range chunks {
			contents[int32(i)] = chunk
		}

		modelOpts := builder.Options(map[string]interface{}{
			"model.name": "text-embedding-ada-002",
		}, nil)

		request := builder.Embedding(credential, modelOpts, nil, contents)

		assert.NotNil(t, request)
		assert.Len(t, request.Content, 3)
		for i, chunk := range chunks {
			assert.Equal(t, chunk, request.Content[int32(i)])
		}
	})
}
