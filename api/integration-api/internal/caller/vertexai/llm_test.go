// Rapida – Open Source Voice AI Orchestration Platform
// Copyright (C) 2023-2025 Prashant Srivastav <prashant@rapida.ai>
// Licensed under a modified GPL-2.0. See the LICENSE file for details.
package internal_vertexai_callers

import (
	"encoding/json"
	"testing"

	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/protos"
	"github.com/stretchr/testify/assert"
)

func TestBuildHistory_UserMessage(t *testing.T) {
	caller := &largeLanguageCaller{
		VertexAi: VertexAi{
			logger: newTestLogger(),
		},
	}

	allMessages := []*protos.Message{
		{
			Role: "user",
			Message: &protos.Message_User{
				User: &protos.UserMessage{
					Content: "Hello, how are you?",
				},
			},
		},
	}

	instruction, history, lastPart := caller.buildHistory(allMessages)

	// First message becomes the instruction
	assert.NotNil(t, instruction)
	assert.Equal(t, "user", instruction.Role)
	assert.Equal(t, 0, len(history))
	assert.Equal(t, "Hello, how are you?", lastPart.Text)
}

func TestBuildHistory_AssistantMessage_WithContent(t *testing.T) {
	caller := &largeLanguageCaller{
		VertexAi: VertexAi{
			logger: newTestLogger(),
		},
	}

	allMessages := []*protos.Message{
		{
			Role: "assistant",
			Message: &protos.Message_Assistant{
				Assistant: &protos.AssistantMessage{
					Contents: []string{"I'm doing well", "How can I help?"},
				},
			},
		},
	}

	instruction, history, _ := caller.buildHistory(allMessages)

	// First message becomes the instruction
	assert.NotNil(t, instruction)
	assert.Equal(t, "model", instruction.Role)
	assert.Equal(t, 0, len(history))
	assert.Equal(t, 2, len(instruction.Parts))
	assert.Equal(t, "I'm doing well", instruction.Parts[0].Text)
	assert.Equal(t, "How can I help?", instruction.Parts[1].Text)
}

func TestBuildHistory_AssistantMessage_WithToolCall(t *testing.T) {
	caller := &largeLanguageCaller{
		VertexAi: VertexAi{
			logger: newTestLogger(),
		},
	}

	toolArgs := map[string]interface{}{
		"query": "weather in NYC",
	}
	argsJSON, _ := json.Marshal(toolArgs)

	allMessages := []*protos.Message{
		{
			Role: "assistant",
			Message: &protos.Message_Assistant{
				Assistant: &protos.AssistantMessage{
					Contents: []string{"Let me check the weather for you"},
					ToolCalls: []*protos.ToolCall{
						{
							Id: "call_123",
							Function: &protos.FunctionCall{
								Name:      "get_weather",
								Arguments: string(argsJSON),
							},
						},
					},
				},
			},
		},
	}

	instruction, history, _ := caller.buildHistory(allMessages)

	// First message becomes the instruction
	assert.NotNil(t, instruction)
	assert.Equal(t, "model", instruction.Role)
	assert.Equal(t, 0, len(history))
	assert.Equal(t, 2, len(instruction.Parts))

	// Check text content
	assert.Equal(t, "Let me check the weather for you", instruction.Parts[0].Text)

	// Check tool call
	assert.NotNil(t, instruction.Parts[1].FunctionCall)
	assert.Equal(t, "call_123", instruction.Parts[1].FunctionCall.ID)
	assert.Equal(t, "get_weather", instruction.Parts[1].FunctionCall.Name)
	assert.Equal(t, "weather in NYC", instruction.Parts[1].FunctionCall.Args["query"])
}

func TestBuildHistory_SystemMessage(t *testing.T) {
	caller := &largeLanguageCaller{
		VertexAi: VertexAi{
			logger: newTestLogger(),
		},
	}

	allMessages := []*protos.Message{
		{
			Role: "system",
			Message: &protos.Message_System{
				System: &protos.SystemMessage{
					Content: "You are a helpful assistant",
				},
			},
		},
	}

	instruction, history, _ := caller.buildHistory(allMessages)

	// System message becomes the instruction
	assert.NotNil(t, instruction)
	assert.Equal(t, "", instruction.Role) // System message has no role
	assert.Equal(t, 0, len(history))
	assert.Equal(t, "You are a helpful assistant", instruction.Parts[0].Text)
}

func TestBuildHistory_ToolMessage(t *testing.T) {
	caller := &largeLanguageCaller{
		VertexAi: VertexAi{
			logger: newTestLogger(),
		},
	}

	toolResult := map[string]interface{}{
		"temperature": 72,
		"condition":   "sunny",
	}
	resultJSON, _ := json.Marshal(toolResult)

	allMessages := []*protos.Message{
		{
			Role: "tool",
			Message: &protos.Message_Tool{
				Tool: &protos.ToolMessage{
					Tools: []*protos.ToolMessage_Tool{
						{
							Name:    "get_weather",
							Id:      "call_123",
							Content: string(resultJSON),
						},
					},
				},
			},
		},
	}

	instruction, history, _ := caller.buildHistory(allMessages)

	// Tool message becomes the instruction
	assert.NotNil(t, instruction)
	assert.Equal(t, "user", instruction.Role)
	assert.Equal(t, 0, len(history))
	assert.Equal(t, 1, len(instruction.Parts))
	assert.NotNil(t, instruction.Parts[0].FunctionResponse)
	assert.Equal(t, "get_weather", instruction.Parts[0].FunctionResponse.Name)
	assert.Equal(t, "call_123", instruction.Parts[0].FunctionResponse.ID)
	assert.Equal(t, float64(72), instruction.Parts[0].FunctionResponse.Response["temperature"])
	assert.Equal(t, "sunny", instruction.Parts[0].FunctionResponse.Response["condition"])
}

func TestBuildHistory_MixedMessages(t *testing.T) {
	caller := &largeLanguageCaller{
		VertexAi: VertexAi{
			logger: newTestLogger(),
		},
	}

	allMessages := []*protos.Message{
		{
			Role: "system",
			Message: &protos.Message_System{
				System: &protos.SystemMessage{
					Content: "You are a helpful assistant",
				},
			},
		},
		{
			Role: "user",
			Message: &protos.Message_User{
				User: &protos.UserMessage{
					Content: "What's the weather?",
				},
			},
		},
		{
			Role: "assistant",
			Message: &protos.Message_Assistant{
				Assistant: &protos.AssistantMessage{
					Contents: []string{"Let me check"},
				},
			},
		},
	}

	instruction, history, lastPart := caller.buildHistory(allMessages)

	// System message should be returned as instruction (first message)
	assert.NotNil(t, instruction)
	assert.Equal(t, "You are a helpful assistant", instruction.Parts[0].Text)

	// History should contain user and assistant messages
	assert.Equal(t, 2, len(history))
	assert.Equal(t, "user", history[0].Role)
	assert.Equal(t, "model", history[1].Role)
	assert.Equal(t, "Let me check", lastPart.Text)
}

func TestBuildHistory_EmptyMessages(t *testing.T) {
	caller := &largeLanguageCaller{
		VertexAi: VertexAi{
			logger: newTestLogger(),
		},
	}

	allMessages := []*protos.Message{}

	instruction, history, lastPart := caller.buildHistory(allMessages)

	assert.Nil(t, instruction)
	assert.Equal(t, 0, len(history))
	assert.Equal(t, "", lastPart.Text)
}

func TestBuildHistory_MessageWithoutContent(t *testing.T) {
	caller := &largeLanguageCaller{
		VertexAi: VertexAi{
			logger: newTestLogger(),
		},
	}

	allMessages := []*protos.Message{
		{
			Role:    "user",
			Message: nil,
		},
	}

	instruction, history, _ := caller.buildHistory(allMessages)

	assert.Nil(t, instruction)
	assert.Equal(t, 0, len(history))
}

func TestBuildHistory_ModelRole(t *testing.T) {
	caller := &largeLanguageCaller{
		VertexAi: VertexAi{
			logger: newTestLogger(),
		},
	}

	allMessages := []*protos.Message{
		{
			Role: "model",
			Message: &protos.Message_Assistant{
				Assistant: &protos.AssistantMessage{
					Contents: []string{"This is a model response"},
				},
			},
		},
	}

	instruction, history, _ := caller.buildHistory(allMessages)

	// Model message becomes the instruction
	assert.NotNil(t, instruction)
	assert.Equal(t, 0, len(history))
	assert.Equal(t, "model", instruction.Role)
	assert.Equal(t, "This is a model response", instruction.Parts[0].Text)
}

func TestBuildHistory_MultipleToolMessages(t *testing.T) {
	caller := &largeLanguageCaller{
		VertexAi: VertexAi{
			logger: newTestLogger(),
		},
	}

	result1 := map[string]interface{}{"status": "ok", "value": 100}
	result1JSON, _ := json.Marshal(result1)

	result2 := map[string]interface{}{"status": "completed", "value": 200}
	result2JSON, _ := json.Marshal(result2)

	allMessages := []*protos.Message{
		{
			Role: "tool",
			Message: &protos.Message_Tool{
				Tool: &protos.ToolMessage{
					Tools: []*protos.ToolMessage_Tool{
						{
							Name:    "operation1",
							Id:      "call_1",
							Content: string(result1JSON),
						},
						{
							Name:    "operation2",
							Id:      "call_2",
							Content: string(result2JSON),
						},
					},
				},
			},
		},
	}

	instruction, history, _ := caller.buildHistory(allMessages)

	// Tool message becomes the instruction
	assert.NotNil(t, instruction)
	assert.Equal(t, 0, len(history))
	assert.Equal(t, 2, len(instruction.Parts))
	assert.Equal(t, "operation1", instruction.Parts[0].FunctionResponse.Name)
	assert.Equal(t, "operation2", instruction.Parts[1].FunctionResponse.Name)
}

func TestBuildHistory_InvalidToolJSONHandling(t *testing.T) {
	caller := &largeLanguageCaller{
		VertexAi: VertexAi{
			logger: newTestLogger(),
		},
	}

	allMessages := []*protos.Message{
		{
			Role: "tool",
			Message: &protos.Message_Tool{
				Tool: &protos.ToolMessage{
					Tools: []*protos.ToolMessage_Tool{
						{
							Name:    "operation1",
							Id:      "call_1",
							Content: "invalid json {{{",
						},
					},
				},
			},
		},
	}

	instruction, history, _ := caller.buildHistory(allMessages)

	// Tool message becomes the instruction
	assert.NotNil(t, instruction)
	assert.Equal(t, 0, len(history))
	// Should still create function response with empty map for invalid JSON
	assert.Equal(t, "operation1", instruction.Parts[0].FunctionResponse.Name)
	assert.Equal(t, 0, len(instruction.Parts[0].FunctionResponse.Response))
}

func TestBuildHistory_InvalidToolCallJSONHandling(t *testing.T) {
	caller := &largeLanguageCaller{
		VertexAi: VertexAi{
			logger: newTestLogger(),
		},
	}

	allMessages := []*protos.Message{
		{
			Role: "assistant",
			Message: &protos.Message_Assistant{
				Assistant: &protos.AssistantMessage{
					Contents: []string{},
					ToolCalls: []*protos.ToolCall{
						{
							Id: "call_123",
							Function: &protos.FunctionCall{
								Name:      "get_weather",
								Arguments: "invalid json {{{",
							},
						},
					},
				},
			},
		},
	}

	instruction, history, _ := caller.buildHistory(allMessages)

	// Assistant message becomes the instruction
	assert.NotNil(t, instruction)
	assert.Equal(t, 0, len(history))
	assert.NotNil(t, instruction.Parts[0].FunctionCall)
	// Should still create function call with empty map for invalid JSON
	assert.Equal(t, 0, len(instruction.Parts[0].FunctionCall.Args))
}

func TestBuildHistory_LastPartExtraction(t *testing.T) {
	caller := &largeLanguageCaller{
		VertexAi: VertexAi{
			logger: newTestLogger(),
		},
	}

	allMessages := []*protos.Message{
		{
			Role: "user",
			Message: &protos.Message_User{
				User: &protos.UserMessage{
					Content: "First message",
				},
			},
		},
		{
			Role: "assistant",
			Message: &protos.Message_Assistant{
				Assistant: &protos.AssistantMessage{
					Contents: []string{"Second message"},
				},
			},
		},
		{
			Role: "user",
			Message: &protos.Message_User{
				User: &protos.UserMessage{
					Content: "Third message",
				},
			},
		},
	}

	instruction, history, lastPart := caller.buildHistory(allMessages)

	// First message becomes the instruction
	assert.NotNil(t, instruction)
	assert.Equal(t, "First message", instruction.Parts[0].Text)
	assert.Equal(t, 2, len(history))
	// Last part should be the first part of the last message
	assert.Equal(t, "Third message", lastPart.Text)
}

// Helper function to create a test logger
func newTestLogger() commons.Logger {
	lgr, _ := commons.NewApplicationLogger()
	return lgr
}
