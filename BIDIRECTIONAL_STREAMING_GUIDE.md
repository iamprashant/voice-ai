# Bidirectional Streaming Chat Implementation Guide

## Overview

The streaming chat functionality has been optimized from **unidirectional (server-streaming)** to **bidirectional gRPC streaming**. This enables a persistent connection where clients can send multiple chat requests without reconnecting, significantly improving efficiency and real-time communication capabilities.

## Architecture Changes

### Before (Unidirectional)

```
Client          │          Server
    │                          │
    ├─► ChatRequest ──────────► │
    │                          │ Process
    │ ◄──────── ChatResponse ◄─┤
    │ ◄──────── ChatResponse ◄─┤
    │ ◄──────── ChatResponse ◄─┤
    │                      (Close)
    │                          │
```

### After (Bidirectional)

```
Client          │          Server
    │                          │
    ├─► ChatRequest ──────────► │
    │  ◄──────── ChatResponse ◄─┤
    │                    (Process)
    ├─► ChatRequest ──────────► │ ◄─ NEW: Another request
    │  ◄──────── ChatResponse ◄─┤     without reconnecting
    │  ◄──────── ChatResponse ◄─┤
    │                          │
    │ (Persistent connection)  │
    │                          │
```

## Benefits

1. **Single Persistent Connection**: No reconnection overhead for multiple messages
2. **Real-time Bidirectional Communication**: True request/response streaming
3. **Improved Efficiency**: Reduced latency and network overhead
4. **Better User Experience**: Seamless multi-turn conversations
5. **Server Resource Efficiency**: Single connection per client instead of multiple short-lived connections

## Proto Changes

### File: `protos/artifacts/integration-api.proto`

Updated all services' StreamChat RPC definitions:

```protobuf
// OLD - Unidirectional (Server Streaming)
rpc StreamChat(ChatRequest) returns (stream ChatResponse);

// NEW - Bidirectional
rpc StreamChat(stream ChatRequest) returns (stream ChatResponse);
```

**Services Updated:**

- OpenAiService
- AzureService
- GeminiService
- VertexAiService
- ReplicateService
- AnthropicService
- CohereService
- MistralService

## Go Implementation Changes

### New Method: `StreamChatBidirectional`

Added in `api/integration-api/api/chat.go`:

```go
func (iApi *integrationApi) StreamChatBidirectional(
    context context.Context,
    providerName string,
    callerFactory func(*protos.Credential) internal_callers.LargeLanguageCaller,
    stream grpc.BidiStreamingServer[protos.ChatRequest, protos.ChatResponse],
) error
```

**Key Features:**

- Authenticates client once at the beginning
- Maintains persistent connection
- Receives multiple `ChatRequest` messages via `stream.Recv()`
- Sends responses via `stream.Send()`
- Handles EOF when client closes stream gracefully
- Generates unique request ID for each message
- Supports per-request authentication via credentials in each ChatRequest

### Flow

```go
1. Authenticate client (once)
2. Loop:
    a. Receive ChatRequest from stream
    b. Extract credential from request
    c. Create LLM caller with credential
    d. Process StreamChatCompletion
    e. Send responses back through stream
    f. Repeat until client closes (EOF)
3. Return nil on graceful close
```

## Updated Provider Implementations

All providers have been updated to use the new bidirectional streaming signature:

**Files Modified:**

- `api/integration-api/api/anthropic.go`
- `api/integration-api/api/openai.go`
- `api/integration-api/api/azure.go`
- `api/integration-api/api/cohere.go`
- `api/integration-api/api/gemini.go`
- `api/integration-api/api/mistral.go`
- `api/integration-api/api/vertexai.go`
- `api/integration-api/api/replicate.go`

**Example:**

```go
// OLD Implementation
func (api *openaiIntegrationGRPCApi) StreamChat(
    irRequest *integration_api.ChatRequest,
    stream integration_api.OpenAiService_StreamChatServer,
) error {
    return api.integrationApi.StreamChat(
        irRequest,
        stream.Context(),
        "OPENAI",
        internal_openai_callers.NewLargeLanguageCaller(api.logger, irRequest.GetCredential()),
        stream.Send,
    )
}

// NEW Implementation
func (oiGRPC *openaiIntegrationGRPCApi) StreamChat(
    stream integration_api.OpenAiService_StreamChatServer,
) error {
    oiGRPC.logger.Debugf("Bidirectional stream chat opened for openai")
    return oiGRPC.integrationApi.StreamChatBidirectional(
        stream.Context(),
        "OPENAI",
        func(cred *integration_api.Credential) internal_callers.LargeLanguageCaller {
            return internal_openai_callers.NewLargeLanguageCaller(oiGRPC.logger, cred)
        },
        stream,
    )
}
```

## Client Usage Example

### Old Approach (Unidirectional)

```go
// Unary request, streaming response
request := &protos.ChatRequest{...}
stream, err := client.StreamChat(ctx, request)
// Process responses
for {
    resp, err := stream.Recv()
    if err == io.EOF {
        break
    }
    // Handle response
}
// For next conversation, need new connection
```

### New Approach (Bidirectional)

```go
// Bidirectional stream
stream, err := client.StreamChat(ctx)

// Send first request
stream.Send(&protos.ChatRequest{...})

// Receive responses
for {
    resp, err := stream.Recv()
    if err == io.EOF {
        break
    }
    // Handle response
}

// Send another request on SAME connection
stream.Send(&protos.ChatRequest{...}) // Conversation continues!

// Receive responses for second request
for {
    resp, err := stream.Recv()
    // ...
}

// Close when done
stream.CloseSend()
```

## Error Handling

### Server-Side

- Errors during stream reception are caught and logged
- Malformed/nil requests are skipped (continue to next)
- Processing errors are sent back to client without closing stream
- Stream remains open to accept further requests

### Client-Side

- `io.EOF` indicates graceful stream closure
- Other errors indicate connection issues
- Can handle errors per-message and continue

## Testing Recommendations

1. **Single Message**: Send one request, receive response
2. **Multiple Messages**: Send multiple requests on same connection
3. **Concurrent Requests**: Send requests before all responses received
4. **Error Cases**: Test with invalid credentials, malformed requests
5. **Connection Termination**: Test client-initiated and server-initiated closes
6. **Different Providers**: Test bidirectional streaming with each LLM provider

## Migration Notes

- **Backwards Compatibility**: Old synchronous `Chat()` RPC remains unchanged
- **Existing Code**: Clients using old streaming approach will need updates
- **Performance**: No performance loss; actually improved due to connection persistence
- **Load Testing**: Recommend testing with higher concurrency due to persistent connections

## Future Enhancements

1. **Connection Pooling**: Implement client-side connection pooling
2. **Keepalive**: Add gRPC keepalive configuration
3. **Rate Limiting**: Per-stream rate limiting
4. **Metrics**: Add connection duration and message throughput metrics
5. **Graceful Degradation**: Fallback to unary RPC if bidirectional not supported

## Configuration

The bidirectional streaming is now the default for all `StreamChat` RPCs. No additional configuration required.

## Support

For issues or questions about bidirectional streaming implementation, refer to:

- gRPC Documentation: https://grpc.io/docs/guides/performance-best-practices/
- Go gRPC: https://pkg.go.dev/google.golang.org/grpc
- Project Documentation: See copilot-instructions.md
