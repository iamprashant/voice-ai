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

func TestNewRerankingInputBuilder(t *testing.T) {
	logger := newTestLogger()
	builder := NewRerankingInputBuilder(logger)

	assert.NotNil(t, builder, "NewRerankingInputBuilder should return a non-nil builder")
	assert.Implements(t, (*InputRerankingBuilder)(nil), builder, "builder should implement InputRerankingBuilder")
}

func TestRerankingInputBuilder_Credential(t *testing.T) {
	logger := newTestLogger()
	builder := NewRerankingInputBuilder(logger)

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
			name:     "credential with cohere config",
			id:       200,
			value:    map[string]interface{}{"api_key": "cohere-key", "model": "rerank-english-v3.0"},
			expected: 200,
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

func TestRerankingInputBuilder_Credential_NilValue(t *testing.T) {
	logger := newTestLogger()
	builder := NewRerankingInputBuilder(logger)

	cred := builder.Credential(789, nil)

	assert.NotNil(t, cred, "credential should not be nil")
	assert.Equal(t, uint64(789), cred.Id, "credential id should match")
	assert.Nil(t, cred.Value, "credential value should be nil")
}

func TestRerankingInputBuilder_Reranking(t *testing.T) {
	logger := newTestLogger()
	builder := NewRerankingInputBuilder(logger)

	t.Run("basic reranking request", func(t *testing.T) {
		structVal, _ := structpb.NewStruct(map[string]interface{}{"api_key": "test"})
		credential := builder.Credential(1, structVal)

		modelOpts := make(map[string]*anypb.Any)
		topNVal, _ := structpb.NewValue(5)
		anyTopN, _ := anypb.New(topNVal)
		modelOpts["top_n"] = anyTopN

		contents := map[int32]string{
			0: "Document about machine learning",
			1: "Document about deep learning",
			2: "Document about natural language processing",
		}

		additionalData := map[string]string{
			"query": "What is deep learning?",
		}

		request := builder.Reranking(credential, modelOpts, additionalData, contents)

		assert.NotNil(t, request, "request should not be nil")
		assert.NotNil(t, request.Credential, "credential should not be nil")
		assert.Equal(t, uint64(1), request.Credential.Id, "credential id should match")
		assert.NotNil(t, request.ModelParameters, "model parameters should not be nil")
		assert.Len(t, request.Content, 3, "should have three contents")
		assert.Equal(t, "What is deep learning?", request.AdditionalData["query"], "query should match")
	})

	t.Run("reranking request with nil model opts", func(t *testing.T) {
		structVal, _ := structpb.NewStruct(map[string]interface{}{"api_key": "test"})
		credential := builder.Credential(1, structVal)

		contents := map[int32]string{
			0: "Text to rerank",
		}

		request := builder.Reranking(credential, nil, nil, contents)

		assert.NotNil(t, request, "request should not be nil")
		assert.Nil(t, request.ModelParameters, "model parameters should be nil")
		assert.Nil(t, request.AdditionalData, "additional data should be nil")
	})

	t.Run("reranking request with empty contents", func(t *testing.T) {
		structVal, _ := structpb.NewStruct(map[string]interface{}{"api_key": "test"})
		credential := builder.Credential(1, structVal)

		contents := map[int32]string{}

		request := builder.Reranking(credential, nil, nil, contents)

		assert.NotNil(t, request, "request should not be nil")
		assert.Empty(t, request.Content, "content should be empty")
	})

	t.Run("reranking request with many documents", func(t *testing.T) {
		structVal, _ := structpb.NewStruct(map[string]interface{}{"api_key": "test"})
		credential := builder.Credential(1, structVal)

		// Simulate a large batch of documents for reranking
		contents := make(map[int32]string)
		for i := int32(0); i < 50; i++ {
			contents[i] = "Document content " + string(rune('A'+i%26))
		}

		request := builder.Reranking(credential, nil, nil, contents)

		assert.NotNil(t, request, "request should not be nil")
		assert.Len(t, request.Content, 50, "should have 50 contents")
	})

	t.Run("reranking request with sparse indices", func(t *testing.T) {
		structVal, _ := structpb.NewStruct(map[string]interface{}{"api_key": "test"})
		credential := builder.Credential(1, structVal)

		contents := map[int32]string{
			0:   "First document",
			10:  "Tenth document",
			100: "Hundredth document",
		}

		request := builder.Reranking(credential, nil, nil, contents)

		assert.NotNil(t, request, "request should not be nil")
		assert.Len(t, request.Content, 3, "should have three contents")
		assert.Equal(t, "First document", request.Content[0], "index 0 should match")
		assert.Equal(t, "Tenth document", request.Content[10], "index 10 should match")
		assert.Equal(t, "Hundredth document", request.Content[100], "index 100 should match")
	})

	t.Run("reranking request with multilingual content", func(t *testing.T) {
		structVal, _ := structpb.NewStruct(map[string]interface{}{"api_key": "test"})
		credential := builder.Credential(1, structVal)

		contents := map[int32]string{
			0: "English document about artificial intelligence",
			1: "文档关于人工智能",                                 // Chinese
			2: "Document sur l'intelligence artificielle", // French
			3: "Documento sobre inteligencia artificial",  // Spanish
		}

		request := builder.Reranking(credential, nil, nil, contents)

		assert.NotNil(t, request, "request should not be nil")
		assert.Len(t, request.Content, 4, "should have four contents")
		assert.Equal(t, "文档关于人工智能", request.Content[1], "Chinese content should be preserved")
	})
}

func TestRerankingInputBuilder_Arguments(t *testing.T) {
	logger := newTestLogger()
	builder := NewRerankingInputBuilder(logger).(*rerankingInputBuilder)

	t.Run("merge with existing variables", func(t *testing.T) {
		variables := []*gorm_types.PromptVariable{
			{Name: "model", DefaultValue: "rerank-english-v2.0"},
			{Name: "top_n", DefaultValue: "10"},
		}

		modelVal, _ := structpb.NewValue("rerank-multilingual-v3.0")
		modelAny, _ := anypb.New(modelVal)

		arguments := map[string]*anypb.Any{
			"model": modelAny,
		}

		result := builder.Arguments(variables, arguments)

		// The anypb value gets converted to a map with protobuf metadata
		assert.NotNil(t, result["model"], "model should be overridden")
		assert.Equal(t, "10", result["top_n"], "top_n should use default")
	})

	t.Run("empty variables", func(t *testing.T) {
		variables := []*gorm_types.PromptVariable{}

		topNVal, _ := structpb.NewValue(5)
		topNAny, _ := anypb.New(topNVal)

		arguments := map[string]*anypb.Any{
			"top_n": topNAny,
		}

		result := builder.Arguments(variables, arguments)

		assert.NotNil(t, result["top_n"], "top_n should be set from arguments")
	})

	t.Run("nil arguments", func(t *testing.T) {
		variables := []*gorm_types.PromptVariable{
			{Name: "model", DefaultValue: "default-rerank-model"},
		}

		result := builder.Arguments(variables, nil)

		assert.Equal(t, "default-rerank-model", result["model"], "should use default value")
	})

	t.Run("both empty", func(t *testing.T) {
		variables := []*gorm_types.PromptVariable{}
		arguments := map[string]*anypb.Any{}

		result := builder.Arguments(variables, arguments)

		assert.Empty(t, result, "result should be empty")
	})

	t.Run("override multiple variables", func(t *testing.T) {
		variables := []*gorm_types.PromptVariable{
			{Name: "model", DefaultValue: "default-model"},
			{Name: "top_n", DefaultValue: "10"},
			{Name: "return_documents", DefaultValue: "true"},
		}

		modelVal, _ := structpb.NewValue("new-model")
		modelAny, _ := anypb.New(modelVal)
		topNVal, _ := structpb.NewValue(5)
		topNAny, _ := anypb.New(topNVal)

		arguments := map[string]*anypb.Any{
			"model": modelAny,
			"top_n": topNAny,
		}

		result := builder.Arguments(variables, arguments)

		// The anypb values get converted to maps with protobuf metadata
		assert.NotNil(t, result["model"], "model should be overridden")
		assert.NotNil(t, result["top_n"], "top_n should be overridden")
		assert.Equal(t, "true", result["return_documents"], "return_documents should use default")
	})
}

func TestRerankingInputBuilder_Options(t *testing.T) {
	logger := newTestLogger()
	builder := NewRerankingInputBuilder(logger)

	t.Run("add options to nil map", func(t *testing.T) {
		opts := map[string]interface{}{
			"model": "rerank-english-v3.0",
			"top_n": 5,
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
			"top_n": 10,
		}

		result := builder.Options(opts, existingOpts)

		assert.Len(t, result, 2, "should have two options")
		assert.NotNil(t, result["model"], "existing key should be preserved")
		assert.NotNil(t, result["top_n"], "new key should be added")
	})

	t.Run("override existing option", func(t *testing.T) {
		existingVal, _ := structpb.NewValue(10)
		existingAny, _ := anypb.New(existingVal)

		existingOpts := map[string]*anypb.Any{
			"top_n": existingAny,
		}

		opts := map[string]interface{}{
			"top_n": 5,
		}

		result := builder.Options(opts, existingOpts)

		assert.Len(t, result, 1, "should have one option")
	})

	t.Run("empty opts preserves existing", func(t *testing.T) {
		existingVal, _ := structpb.NewValue("existing-model")
		existingAny, _ := anypb.New(existingVal)

		existingOpts := map[string]*anypb.Any{
			"model": existingAny,
		}

		opts := map[string]interface{}{}

		result := builder.Options(opts, existingOpts)

		assert.Len(t, result, 1, "should preserve existing options")
	})

	t.Run("reranking specific options", func(t *testing.T) {
		opts := map[string]interface{}{
			"model":            "rerank-english-v3.0",
			"top_n":            10,
			"return_documents": true,
			"max_chunks":       200,
		}

		result := builder.Options(opts, nil)

		assert.Len(t, result, 4, "should handle all reranking options")
	})

	t.Run("nil opts and nil existing", func(t *testing.T) {
		result := builder.Options(nil, nil)

		assert.NotNil(t, result, "result should not be nil")
		assert.Empty(t, result, "result should be empty")
	})
}

func TestRerankingInputBuilder_Integration(t *testing.T) {
	logger := newTestLogger()
	builder := NewRerankingInputBuilder(logger)

	t.Run("full reranking request workflow for search results", func(t *testing.T) {
		// Create credential
		structVal, _ := structpb.NewStruct(map[string]interface{}{
			"api_key": "cohere-test-key",
		})
		credential := builder.Credential(42, structVal)

		// Create model options
		modelOpts := builder.Options(map[string]interface{}{
			"model.name":       "rerank-english-v3.0",
			"model.top_n":      5,
			"model.max_chunks": 200,
		}, nil)

		// Simulate search results to rerank
		searchResults := map[int32]string{
			0: "Deep learning is a subset of machine learning that uses neural networks with multiple layers.",
			1: "Machine learning algorithms can learn from and make predictions on data.",
			2: "Artificial intelligence is the broader concept of machines being able to carry out tasks in a smart way.",
			3: "Neural networks are computing systems inspired by biological neural networks in animal brains.",
			4: "Supervised learning uses labeled datasets to train algorithms to classify data.",
			5: "Reinforcement learning trains agents to make decisions by rewarding desired behaviors.",
			6: "Transfer learning allows models trained on one task to be used for another related task.",
			7: "Computer vision enables machines to interpret and make decisions based on visual data.",
			8: "Natural language processing helps computers understand, interpret, and generate human language.",
			9: "Convolutional neural networks are particularly effective for image recognition tasks.",
		}

		// Create additional data with the search query
		additionalData := map[string]string{
			"query":      "What is deep learning and how does it relate to neural networks?",
			"session_id": "session-xyz",
		}

		// Build the request
		request := builder.Reranking(credential, modelOpts, additionalData, searchResults)

		// Assertions
		assert.NotNil(t, request)
		assert.Equal(t, uint64(42), request.Credential.Id)
		assert.Len(t, request.Content, 10)
		assert.Equal(t, "What is deep learning and how does it relate to neural networks?", request.AdditionalData["query"])
		assert.NotNil(t, request.ModelParameters["model.name"])
		assert.NotNil(t, request.ModelParameters["model.top_n"])
	})

	t.Run("reranking request for RAG pipeline", func(t *testing.T) {
		structVal, _ := structpb.NewStruct(map[string]interface{}{"api_key": "test"})
		credential := builder.Credential(1, structVal)

		// Simulate retrieved document chunks from a vector database
		chunks := map[int32]string{
			0: "The company was founded in 2020 by a team of AI researchers from Stanford University.",
			1: "Our mission is to democratize AI and make it accessible to everyone.",
			2: "The product uses state-of-the-art language models for natural language understanding.",
			3: "We have raised $50 million in Series A funding led by top-tier venture capital firms.",
			4: "The team consists of over 50 engineers and researchers from leading tech companies.",
		}

		modelOpts := builder.Options(map[string]interface{}{
			"model.name":       "rerank-multilingual-v3.0",
			"model.top_n":      3,
			"return_documents": true,
		}, nil)

		additionalData := map[string]string{
			"query":        "When was the company founded and by whom?",
			"knowledge_id": "kb-12345",
		}

		request := builder.Reranking(credential, modelOpts, additionalData, chunks)

		assert.NotNil(t, request)
		assert.Len(t, request.Content, 5)
		assert.Equal(t, "When was the company founded and by whom?", request.AdditionalData["query"])
		assert.Equal(t, "kb-12345", request.AdditionalData["knowledge_id"])
	})
}

// Benchmark tests
func BenchmarkRerankingInputBuilder_Reranking(b *testing.B) {
	logger := newTestLogger()
	builder := NewRerankingInputBuilder(logger)

	structVal, _ := structpb.NewStruct(map[string]interface{}{"api_key": "test"})
	credential := builder.Credential(1, structVal)

	contents := make(map[int32]string)
	for i := int32(0); i < 100; i++ {
		contents[i] = "Document content for benchmarking"
	}

	additionalData := map[string]string{
		"query": "benchmark query",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builder.Reranking(credential, nil, additionalData, contents)
	}
}

func BenchmarkRerankingInputBuilder_Options(b *testing.B) {
	logger := newTestLogger()
	builder := NewRerankingInputBuilder(logger)

	opts := map[string]interface{}{
		"model":            "rerank-english-v3.0",
		"top_n":            10,
		"return_documents": true,
		"max_chunks":       200,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builder.Options(opts, nil)
	}
}
