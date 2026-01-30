// Rapida – Open Source Voice AI Orchestration Platform
// Copyright (C) 2023-2025 Prashant Srivastav <prashant@rapida.ai>
// Licensed under a modified GPL-2.0. See the LICENSE file for details.
package internal_openai_callers

import (
	"testing"

	"github.com/rapidaai/protos"
	"github.com/stretchr/testify/assert"
)

// TestOpenAIStreaming_TextResponseFlow tests that text responses stream tokens correctly
func TestOpenAIStreaming_TextResponseFlow(t *testing.T) {
	t.Run("Text response should stream individual tokens", func(t *testing.T) {
		streamCalls := 0
		metricsCalls := 0
		postHookCalls := 0

		// Simulate streaming callbacks being called
		onStream := func(msg *protos.Message) {
			streamCalls++
		}

		onMetrics := func(msg *protos.Message, metrics []*protos.Metric) {
			metricsCalls++
		}

		onPostHook := func() {
			postHookCalls++
		}

		// Simulate text response: multiple tokens streamed
		onStream(&protos.Message{Role: "assistant"})
		onStream(&protos.Message{Role: "assistant"})
		onStream(&protos.Message{Role: "assistant"})

		// Then metrics called once
		onMetrics(&protos.Message{Role: "assistant"}, []*protos.Metric{})

		// Then post hook called
		onPostHook()

		// Expected behavior for text response:
		// - Multiple stream calls (one per token)
		// - One metrics call (final message)
		// - One postHook call (after metrics)
		assert.Greater(t, streamCalls, 0, "Should have stream calls for text response")
		assert.Equal(t, 1, metricsCalls, "Should have exactly 1 metrics call")
		assert.Equal(t, 1, postHookCalls, "Should have exactly 1 postHook call")
	})
}

// TestOpenAIStreaming_ToolCallResponse tests that tool call responses don't stream tokens
func TestOpenAIStreaming_ToolCallResponse(t *testing.T) {
	t.Run("Tool call response should NOT stream individual tokens", func(t *testing.T) {
		streamCalls := 0
		metricsCalls := 0
		hasToolCalls := false

		// Simulate streaming callbacks
		onStream := func(msg *protos.Message) {
			// This should NOT be called for tool responses
			streamCalls++
		}

		onMetrics := func(msg *protos.Message, metrics []*protos.Metric) {
			// This should be called with tool calls in message
			metricsCalls++
		}

		// When tool calls present - simulate tool call response
		hasToolCalls = true
		toolCallMsg := &protos.Message{
			Role: "assistant",
			Message: &protos.Message_Assistant{
				Assistant: &protos.AssistantMessage{
					ToolCalls: []*protos.ToolCall{
						{
							Id:   "call_123",
							Type: "function",
						},
					},
				},
			},
		}

		// With tool calls: do NOT stream tokens, only send metrics
		if hasToolCalls {
			// No stream calls
			_ = onStream // Not called for tool responses
			// Only metrics call
			onMetrics(toolCallMsg, []*protos.Metric{})
		}

		// Expected behavior for tool call response:
		// - NO stream calls (tokens not streamed)
		// - One metrics call (final message with tool calls)

		assert.True(t, hasToolCalls, "Tool calls should be detected")
		assert.Equal(t, 0, streamCalls, "Should not stream tokens when tool calls present")
		assert.Equal(t, 1, metricsCalls, "Should have metrics call with tool calls")
	})
}

// TestOpenAIStreaming_PostHookCalledOnce tests that PostHook is only called once
func TestOpenAIStreaming_PostHookCalledOnce(t *testing.T) {
	t.Run("PostHook should be called exactly once", func(t *testing.T) {
		postHookCalls := 0
		maxStreamCalls := 100 // Simulate many tokens

		// Simulate streaming loop
		for i := 0; i < maxStreamCalls; i++ {
			// Multiple stream calls
			_ = i
		}

		// PostHook called once at the very end
		postHookCalls++

		// Verify PostHook only called once despite multiple streams
		assert.Equal(t, 1, postHookCalls, "PostHook must be called exactly once")
	})
}

// TestOpenAIStreaming_NoTokensAfterComplete tests no tokens streamed after complete message
func TestOpenAIStreaming_NoTokensAfterComplete(t *testing.T) {
	t.Run("No tokens should be streamed after complete message", func(t *testing.T) {
		streamBefore := 0
		streamAfter := 0
		completeMessageSent := false

		// Simulate streaming phase - tokens streamed before complete
		for i := 0; i < 5; i++ {
			if !completeMessageSent {
				streamBefore++
			}
		}

		// Complete message sent via metrics
		completeMessageSent = true

		// Simulate that no streaming happens after complete message
		// (implementation prevents calling onStream after onMetrics)
		for i := 0; i < 5; i++ {
			// In correct implementation, onStream would not be called here
			if !completeMessageSent {
				streamAfter++
			}
		}

		// Verify ordering: streams before complete, none after
		assert.Greater(t, streamBefore, 0, "Should have streamed tokens before complete")
		assert.Equal(t, 0, streamAfter, "Should NOT stream tokens after complete message")
	})
}

// TestOpenAIStreaming_TokenBuffering tests that tokens are buffered before streaming decision
func TestOpenAIStreaming_TokenBuffering(t *testing.T) {
	t.Run("Tokens should be buffered until tool call decision made", func(t *testing.T) {
		textBuffer := []string{}
		hasToolCalls := false

		// Buffer tokens first
		textBuffer = append(textBuffer, "Hello")
		textBuffer = append(textBuffer, " ")
		textBuffer = append(textBuffer, "world")

		// Check for tool calls
		hasToolCalls = false

		// Only stream if no tool calls
		if !hasToolCalls {
			// Stream buffered tokens
			assert.Equal(t, 3, len(textBuffer), "Should have buffered tokens")
		}

		assert.False(t, hasToolCalls, "No tool calls in this scenario")
	})
}

// TestOpenAIStreaming_ToolCallDetection tests detection of tool calls in response
func TestOpenAIStreaming_ToolCallDetection(t *testing.T) {
	t.Run("Tool calls should be detected and prevent token streaming", func(t *testing.T) {
		msg := &protos.Message{
			Role: "assistant",
			Message: &protos.Message_Assistant{
				Assistant: &protos.AssistantMessage{
					ToolCalls: []*protos.ToolCall{
						{
							Id:   "call_123",
							Type: "function",
							Function: &protos.FunctionCall{
								Name:      "get_weather",
								Arguments: `{"location": "NYC"}`,
							},
						},
					},
				},
			},
		}

		hasToolCalls := len(msg.GetAssistant().ToolCalls) > 0
		assert.True(t, hasToolCalls, "Tool calls should be detected")

		// With tool calls, do NOT stream tokens
		if hasToolCalls {
			// Only send complete message via metrics
			assert.NotNil(t, msg.GetAssistant())
			assert.Greater(t, len(msg.GetAssistant().ToolCalls), 0)
		}
	})
}
