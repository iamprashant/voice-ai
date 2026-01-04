// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package token_tiktoken_calculators

import (
	"testing"

	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/tokens"
	"github.com/rapidaai/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestNewTikTokenCostCalculator(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()
	model := "gpt-3.5-turbo"

	calculator := NewTikTokenCostCalculator(logger, model)

	assert.NotNil(t, calculator)
	assert.IsType(t, &tikTokenCostCalculator{}, calculator)
}

func TestTikTokenCostCalculator_Token_EmptyMessages(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()
	calculator := NewTikTokenCostCalculator(logger, "gpt-3.5-turbo")

	in := []*types.Message{}
	out := &types.Message{
		Role: "assistant",
		Contents: []*types.Content{
			{
				ContentType:   "text",
				ContentFormat: "raw",
				Content:       []byte("Hello world"),
			},
		},
	}

	metrics := calculator.Token(in, out)

	assert.Len(t, metrics, 3)
	assert.Equal(t, "INPUT_TOKEN", metrics[0].Name)
	assert.Equal(t, "OUTPUT_TOKEN", metrics[1].Name)
	assert.Equal(t, "TOTAL_TOKEN", metrics[2].Name)
}

func TestTikTokenCostCalculator_Token_SingleInputMessage(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()
	calculator := NewTikTokenCostCalculator(logger, "gpt-3.5-turbo")

	in := []*types.Message{
		{
			Role: "user",
			Contents: []*types.Content{
				{
					ContentType:   "text",
					ContentFormat: "raw",
					Content:       []byte("Hello"),
				},
			},
		},
	}
	out := &types.Message{
		Role: "assistant",
		Contents: []*types.Content{
			{
				ContentType:   "text",
				ContentFormat: "raw",
				Content:       []byte("Hi there"),
			},
		},
	}

	metrics := calculator.Token(in, out)

	assert.Len(t, metrics, 3)
	assert.Equal(t, "INPUT_TOKEN", metrics[0].Name)
	assert.Equal(t, "OUTPUT_TOKEN", metrics[1].Name)
	assert.Equal(t, "TOTAL_TOKEN", metrics[2].Name)

	// Values should be positive integers
	inputTokens := metrics[0].Value
	outputTokens := metrics[1].Value
	totalTokens := metrics[2].Value

	assert.NotEqual(t, "0", inputTokens)
	assert.NotEqual(t, "0", outputTokens)
	assert.NotEqual(t, "0", totalTokens)
}

func TestTikTokenCostCalculator_Token_MultipleInputMessages(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()
	calculator := NewTikTokenCostCalculator(logger, "gpt-3.5-turbo")

	in := []*types.Message{
		{
			Role: "system",
			Contents: []*types.Content{
				{
					ContentType:   "text",
					ContentFormat: "raw",
					Content:       []byte("You are a helpful assistant"),
				},
			},
		},
		{
			Role: "user",
			Contents: []*types.Content{
				{
					ContentType:   "text",
					ContentFormat: "raw",
					Content:       []byte("Hello"),
				},
			},
		},
	}
	out := &types.Message{
		Role: "assistant",
		Contents: []*types.Content{
			{
				ContentType:   "text",
				ContentFormat: "raw",
				Content:       []byte("Hi there! How can I help you?"),
			},
		},
	}

	metrics := calculator.Token(in, out)

	assert.Len(t, metrics, 3)
	assert.Equal(t, "INPUT_TOKEN", metrics[0].Name)
	assert.Equal(t, "OUTPUT_TOKEN", metrics[1].Name)
	assert.Equal(t, "TOTAL_TOKEN", metrics[2].Name)
}

func TestTikTokenCostCalculator_Token_GPT35Turbo0613(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()
	calculator := NewTikTokenCostCalculator(logger, "gpt-3.5-turbo-0613")

	in := []*types.Message{
		{
			Role: "user",
			Contents: []*types.Content{
				{
					ContentType:   "text",
					ContentFormat: "raw",
					Content:       []byte("Test"),
				},
			},
		},
	}
	out := &types.Message{
		Role: "assistant",
		Contents: []*types.Content{
			{
				ContentType:   "text",
				ContentFormat: "raw",
				Content:       []byte("Response"),
			},
		},
	}

	metrics := calculator.Token(in, out)

	assert.Len(t, metrics, 3)
	assert.Equal(t, "INPUT_TOKEN", metrics[0].Name)
	assert.Equal(t, "OUTPUT_TOKEN", metrics[1].Name)
	assert.Equal(t, "TOTAL_TOKEN", metrics[2].Name)
}

func TestTikTokenCostCalculator_Token_GPT4(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()
	calculator := NewTikTokenCostCalculator(logger, "gpt-4")

	in := []*types.Message{
		{
			Role: "user",
			Contents: []*types.Content{
				{
					ContentType:   "text",
					ContentFormat: "raw",
					Content:       []byte("Test GPT-4"),
				},
			},
		},
	}
	out := &types.Message{
		Role: "assistant",
		Contents: []*types.Content{
			{
				ContentType:   "text",
				ContentFormat: "raw",
				Content:       []byte("GPT-4 response"),
			},
		},
	}

	metrics := calculator.Token(in, out)

	assert.Len(t, metrics, 3)
	assert.Equal(t, "INPUT_TOKEN", metrics[0].Name)
	assert.Equal(t, "OUTPUT_TOKEN", metrics[1].Name)
	assert.Equal(t, "TOTAL_TOKEN", metrics[2].Name)
}

func TestTikTokenCostCalculator_Token_GPT35Turbo0301(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()
	calculator := NewTikTokenCostCalculator(logger, "gpt-3.5-turbo-0301")

	in := []*types.Message{
		{
			Role: "user",
			Contents: []*types.Content{
				{
					ContentType:   "text",
					ContentFormat: "raw",
					Content:       []byte("Test"),
				},
			},
		},
	}
	out := &types.Message{
		Role: "assistant",
		Contents: []*types.Content{
			{
				ContentType:   "text",
				ContentFormat: "raw",
				Content:       []byte("Response"),
			},
		},
	}

	metrics := calculator.Token(in, out)

	assert.Len(t, metrics, 3)
	assert.Equal(t, "INPUT_TOKEN", metrics[0].Name)
	assert.Equal(t, "OUTPUT_TOKEN", metrics[1].Name)
	assert.Equal(t, "TOTAL_TOKEN", metrics[2].Name)
}

func TestTikTokenCostCalculator_Token_UnknownModel(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()
	calculator := NewTikTokenCostCalculator(logger, "unknown-model")

	in := []*types.Message{
		{
			Role: "user",
			Contents: []*types.Content{
				{
					ContentType:   "text",
					ContentFormat: "raw",
					Content:       []byte("Test"),
				},
			},
		},
	}
	out := &types.Message{
		Role: "assistant",
		Contents: []*types.Content{
			{
				ContentType:   "text",
				ContentFormat: "raw",
				Content:       []byte("Response"),
			},
		},
	}

	metrics := calculator.Token(in, out)

	assert.Len(t, metrics, 3)
	assert.Equal(t, "INPUT_TOKEN", metrics[0].Name)
	assert.Equal(t, "OUTPUT_TOKEN", metrics[1].Name)
	assert.Equal(t, "TOTAL_TOKEN", metrics[2].Name)
	assert.Equal(t, "0", metrics[0].Value) // Should be 0 for unknown model
	assert.Equal(t, "0", metrics[1].Value)
	assert.Equal(t, "0", metrics[2].Value)
}

func TestTikTokenCostCalculator_Token_GPT35TurboVariant(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()
	calculator := NewTikTokenCostCalculator(logger, "gpt-3.5-turbo-1106")

	in := []*types.Message{
		{
			Role: "user",
			Contents: []*types.Content{
				{
					ContentType:   "text",
					ContentFormat: "raw",
					Content:       []byte("Test"),
				},
			},
		},
	}
	out := &types.Message{
		Role: "assistant",
		Contents: []*types.Content{
			{
				ContentType:   "text",
				ContentFormat: "raw",
				Content:       []byte("Response"),
			},
		},
	}

	metrics := calculator.Token(in, out)

	assert.Len(t, metrics, 3)
	assert.Equal(t, "INPUT_TOKEN", metrics[0].Name)
	assert.Equal(t, "OUTPUT_TOKEN", metrics[1].Name)
	assert.Equal(t, "TOTAL_TOKEN", metrics[2].Name)
}

func TestTikTokenCostCalculator_Token_GPT4Variant(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()
	calculator := NewTikTokenCostCalculator(logger, "gpt-4-turbo")

	in := []*types.Message{
		{
			Role: "user",
			Contents: []*types.Content{
				{
					ContentType:   "text",
					ContentFormat: "raw",
					Content:       []byte("Test"),
				},
			},
		},
	}
	out := &types.Message{
		Role: "assistant",
		Contents: []*types.Content{
			{
				ContentType:   "text",
				ContentFormat: "raw",
				Content:       []byte("Response"),
			},
		},
	}

	metrics := calculator.Token(in, out)

	assert.Len(t, metrics, 3)
	assert.Equal(t, "INPUT_TOKEN", metrics[0].Name)
	assert.Equal(t, "OUTPUT_TOKEN", metrics[1].Name)
	assert.Equal(t, "TOTAL_TOKEN", metrics[2].Name)
}

func TestTikTokenCostCalculator_Token_EmptyContent(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()
	calculator := NewTikTokenCostCalculator(logger, "gpt-3.5-turbo")

	in := []*types.Message{
		{
			Role:     "user",
			Contents: []*types.Content{}, // Empty contents
		},
	}
	out := &types.Message{
		Role:     "assistant",
		Contents: []*types.Content{}, // Empty contents
	}

	metrics := calculator.Token(in, out)

	assert.Len(t, metrics, 3)
	assert.Equal(t, "INPUT_TOKEN", metrics[0].Name)
	assert.Equal(t, "OUTPUT_TOKEN", metrics[1].Name)
	assert.Equal(t, "TOTAL_TOKEN", metrics[2].Name)
}

func TestTikTokenCostCalculator_Token_NonTextContent(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()
	calculator := NewTikTokenCostCalculator(logger, "gpt-3.5-turbo")

	in := []*types.Message{
		{
			Role: "user",
			Contents: []*types.Content{
				{
					ContentType:   "image", // Non-text content
					ContentFormat: "url",
					Content:       []byte("http://example.com/image.jpg"),
				},
			},
		},
	}
	out := &types.Message{
		Role: "assistant",
		Contents: []*types.Content{
			{
				ContentType:   "text",
				ContentFormat: "raw",
				Content:       []byte("I see an image"),
			},
		},
	}

	metrics := calculator.Token(in, out)

	assert.Len(t, metrics, 3)
	assert.Equal(t, "INPUT_TOKEN", metrics[0].Name)
	assert.Equal(t, "OUTPUT_TOKEN", metrics[1].Name)
	assert.Equal(t, "TOTAL_TOKEN", metrics[2].Name)
}

func TestTikTokenCostCalculator_Token_NilMessages(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()
	calculator := NewTikTokenCostCalculator(logger, "gpt-3.5-turbo")

	// Test with nil input messages
	in := []*types.Message(nil)
	out := &types.Message{
		Role: "assistant",
		Contents: []*types.Content{
			{
				ContentType:   "text",
				ContentFormat: "raw",
				Content:       []byte("Response"),
			},
		},
	}

	metrics := calculator.Token(in, out)

	assert.Len(t, metrics, 3)
	assert.Equal(t, "INPUT_TOKEN", metrics[0].Name)
	assert.Equal(t, "OUTPUT_TOKEN", metrics[1].Name)
	assert.Equal(t, "TOTAL_TOKEN", metrics[2].Name)
}

func TestTikTokenCostCalculator_Token_LongMessages(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()
	calculator := NewTikTokenCostCalculator(logger, "gpt-3.5-turbo")

	longText := make([]byte, 10000) // 10KB of text
	for i := range longText {
		longText[i] = 'a'
	}

	in := []*types.Message{
		{
			Role: "user",
			Contents: []*types.Content{
				{
					ContentType:   "text",
					ContentFormat: "raw",
					Content:       longText,
				},
			},
		},
	}
	out := &types.Message{
		Role: "assistant",
		Contents: []*types.Content{
			{
				ContentType:   "text",
				ContentFormat: "raw",
				Content:       longText,
			},
		},
	}

	metrics := calculator.Token(in, out)

	assert.Len(t, metrics, 3)
	assert.Equal(t, "INPUT_TOKEN", metrics[0].Name)
	assert.Equal(t, "OUTPUT_TOKEN", metrics[1].Name)
	assert.Equal(t, "TOTAL_TOKEN", metrics[2].Name)

	// Long messages should result in higher token counts
	inputTokens := metrics[0].Value
	outputTokens := metrics[1].Value
	assert.NotEqual(t, "0", inputTokens)
	assert.NotEqual(t, "0", outputTokens)
}

func TestTikTokenCostCalculator_Token_SpecialCharacters(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()
	calculator := NewTikTokenCostCalculator(logger, "gpt-3.5-turbo")

	in := []*types.Message{
		{
			Role: "user",
			Contents: []*types.Content{
				{
					ContentType:   "text",
					ContentFormat: "raw",
					Content:       []byte("Hello 🌍 with émojis and spëcial chärs!"),
				},
			},
		},
	}
	out := &types.Message{
		Role: "assistant",
		Contents: []*types.Content{
			{
				ContentType:   "text",
				ContentFormat: "raw",
				Content:       []byte("Hi there! 👋"),
			},
		},
	}

	metrics := calculator.Token(in, out)

	assert.Len(t, metrics, 3)
	assert.Equal(t, "INPUT_TOKEN", metrics[0].Name)
	assert.Equal(t, "OUTPUT_TOKEN", metrics[1].Name)
	assert.Equal(t, "TOTAL_TOKEN", metrics[2].Name)
}

func TestTikTokenCostCalculator_Token_MixedContentTypes(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()
	calculator := NewTikTokenCostCalculator(logger, "gpt-3.5-turbo")

	in := []*types.Message{
		{
			Role: "user",
			Contents: []*types.Content{
				{
					ContentType:   "text",
					ContentFormat: "raw",
					Content:       []byte("Hello"),
				},
				{
					ContentType:   "image",
					ContentFormat: "url",
					Content:       []byte("http://example.com/image.jpg"),
				},
				{
					ContentType:   "text",
					ContentFormat: "raw",
					Content:       []byte(" world"),
				},
			},
		},
	}
	out := &types.Message{
		Role: "assistant",
		Contents: []*types.Content{
			{
				ContentType:   "text",
				ContentFormat: "raw",
				Content:       []byte("Hi there!"),
			},
		},
	}

	metrics := calculator.Token(in, out)

	assert.Len(t, metrics, 3)
	assert.Equal(t, "INPUT_TOKEN", metrics[0].Name)
	assert.Equal(t, "OUTPUT_TOKEN", metrics[1].Name)
	assert.Equal(t, "TOTAL_TOKEN", metrics[2].Name)
}

func TestTikTokenCostCalculator_ImplementsTokenCalculatorInterface(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()
	calculator := NewTikTokenCostCalculator(logger, "gpt-3.5-turbo")

	var _ tokens.TokenCalculator = calculator
}
