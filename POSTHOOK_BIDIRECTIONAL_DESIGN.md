# PostHook Bidirectional Design Change

## Overview

Updated the PostHook callback timing in all 6 LLM providers to support proper bidirectional streaming. PostHook is now called for **each message end** instead of just at the end of the connection context.

## Change Summary

### Previous Design (Problematic for Bidirectional)

```
PreHook → onStream(tokens) → PostHook → onMetrics
                            ↑
                    Called BEFORE metrics
```

### New Design (Correct for Bidirectional)

```
PreHook → onStream(tokens) → onMetrics → PostHook
                                         ↑
                              Called AFTER metrics
```

## Why This Matters

In bidirectional streaming, multiple message exchanges occur on a single persistent connection:

1. **Client sends message** → LLM streams response → **Message 1 complete**
2. Client processes response
3. **Client sends new message** → LLM streams response → **Message 2 complete**
4. ... and so on

Each message completion should trigger PostHook so that:

- Hooks have access to complete metrics for that message
- Multiple messages in a bidirectional context each get their own PostHook call
- The callback order is consistent and predictable

## Implementation Details

### Callback Lifecycle for Each Message

```
1. PreHook(request_data)              # Before processing
2. onStream(token_1)                  # First token
3. onStream(token_2)                  # More tokens
4. ...
5. onStream(token_N)                  # Last token (if no tool calls)
6. onMetrics(complete_message, metrics) # Message complete with metrics
7. PostHook(result, metrics)          # Hook called AFTER metrics
```

### For Tool Call Responses

Tool calls don't stream individual tokens, so:

```
1. PreHook(request_data)              # Before processing
2. (NO onStream calls)                # Skip token streaming
3. onMetrics(complete_message, metrics) # Message with tool calls
4. PostHook(result, metrics)          # Hook called AFTER metrics
```

## Files Modified

All 6 providers updated with consistent pattern:

### 1. OpenAI Provider

**File**: `api/integration-api/internal/caller/openai/llm.go`

- Moved PostHook call from before to after onMetrics
- No functional change, only callback order

### 2. Azure OpenAI Provider

**File**: `api/integration-api/internal/caller/azure/llm.go`

- Identical change to OpenAI provider

### 3. Anthropic Provider

**File**: `api/integration-api/internal/caller/anthropic/llm.go`

- Moved PostHook call to after onMetrics in MessageStopEvent
- Ensures metrics are available when hook is called

### 4. Gemini Provider

**File**: `api/integration-api/internal/caller/gemini/llm.go`

- Removed PostHook call before metrics processing
- Added PostHook call after onMetrics
- Same timing guarantee as other providers

### 5. VertexAI Provider

**File**: `api/integration-api/internal/caller/vertexai/llm.go`

- Mirrors Gemini implementation
- Moved PostHook to after onMetrics

### 6. Cohere Provider

**File**: `api/integration-api/internal/caller/cohere/llm.go`

- Moved PostHook from before token streaming to after onMetrics
- Ensures consistent callback order

## Bidirectional Streaming Flow

In `StreamChatBidirectional` (api/integration-api/api/chat.go):

```go
for {
    // Receive next message from client
    irRequest, err := stream.Recv()
    if err == io.EOF {
        return nil  // Client closed stream
    }

    // Create a new LLM caller for this message
    llmCaller := callerFactory(irRequest.GetCredential())

    // Process this message's streaming with all callbacks
    err = llmCaller.StreamChatCompletion(
        stream.Context(),
        irRequest.GetConversations(),
        options,  // Includes fresh PreHook and PostHook
        onStream,
        onMetrics,
        onError,
    )

    // Continue to next message in loop (don't close stream)
    // PostHook already called for this message
}
```

Each iteration:

1. Receives a new ChatRequest from client
2. Creates callbacks including PostHook
3. Processes streaming for that message
4. PostHook is called when that message completes
5. Loop continues for next message

## Benefits

1. **Proper Message Lifecycle**: Each message has a complete lifecycle with metrics available when hook runs
2. **Bidirectional Support**: Works correctly for persistent connections with multiple message exchanges
3. **Consistent Ordering**: All providers follow same callback order (stream → metrics → hook)
4. **Hook Access to Metrics**: PostHook can now access complete metrics for that message
5. **Predictable Behavior**: Clear contract for when hooks are called relative to metrics

## Testing

All unit tests pass with the new callback order:

```
✅ TestOpenAIStreaming_TextResponseFlow
✅ TestOpenAIStreaming_ToolCallResponse
✅ TestOpenAIStreaming_PostHookCalledOnce
✅ TestOpenAIStreaming_NoTokensAfterComplete
✅ TestOpenAIStreaming_TokenBuffering
✅ TestOpenAIStreaming_ToolCallDetection
```

## Migration Notes

This is **non-breaking** for:

- Existing unary streaming (single message)
- Stateless providers
- Callback signatures remain unchanged

The change only affects:

- **Timing** of when PostHook is called
- **Guarantees** about metrics being available when hook runs

Consumers relying on PostHook must ensure they handle the new callback order where metrics are available.

## Code Example: PostHook Usage

### Old Behavior (metrics not yet collected)

```go
PostHook: func(data map[string]interface{}, metrics []*protos.Metric) {
    // metrics might not be complete
    log.Println("Request processed")
}
```

### New Behavior (metrics guaranteed available)

```go
PostHook: func(data map[string]interface{}, metrics []*protos.Metric) {
    // metrics are now complete and available
    for _, m := range metrics {
        log.Printf("Token count: %d", m.Value)
    }
    // Can safely access and log complete metrics
}
```

## Future Enhancements

With this proper callback ordering, future enhancements become possible:

- Per-message metrics aggregation
- Per-message billing/logging
- Per-message performance tracking
- Message-specific hooks or callbacks
- Multi-turn conversation tracking

## Rollout Plan

1. **Testing**: Unit tests verify all 6 providers work correctly ✅
2. **Integration**: Integration tests verify bidirectional flow
3. **Staging**: Deploy to staging for end-to-end testing
4. **Production**: Gradual rollout with monitoring

## Questions & Support

For questions about this design change:

- Review `StreamChatBidirectional` in `api/integration-api/api/chat.go`
- Check provider implementations for callback order
- Review unit tests for usage examples
