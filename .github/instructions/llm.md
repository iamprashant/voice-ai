# LLM Provider Integration Guide — UI to Backend

## Overview

This guide covers how to add a new Large Language Model (LLM) provider to the Rapida platform. LLM providers are managed by the **integration-api** service and provide chat completion, embedding, and reranking capabilities.

## Architecture

```
UI (React) → Assistant Config → assistant-api → gRPC → integration-api → Provider SDK → LLM Response
```

### Supported Providers

| Provider | Chat | Embedding | Reranking | Streaming |
|---|---|---|---|---|
| **OpenAI** | ✅ | ✅ | — | ✅ Bidirectional |
| **Anthropic** | ✅ | — | — | ✅ Bidirectional |
| **Gemini** | ✅ | ✅ | — | ✅ Bidirectional |
| **Azure (OpenAI/Foundry)** | ✅ | ✅ | — | ✅ Bidirectional |
| **Cohere** | ✅ | ✅ | ✅ | ✅ Bidirectional |
| **Mistral** | ✅ | — | — | ✅ Bidirectional |
| **HuggingFace** | ✅ | ✅ | — | ✅ |
| **Replicate** | ✅ | — | — | ✅ |
| **Vertex AI** | ✅ | ✅ | — | ✅ Bidirectional |
| **Voyage AI** | — | ✅ | ✅ | — |
| **DeepInfra** | ✅ | — | — | ✅ |

## Directory Structure

### Backend — integration-api

```
api/integration-api/
├── api/                                    # Handler layer
│   ├── integration.go                      # integrationApi base struct + PreHook/PostHook/Chat/Embedding helpers
│   ├── chat.go                             # StreamChatBidirectional, Chat helpers
│   ├── embedding.go                        # Embedding handler helper
│   ├── reranking.go                        # Reranking handler helper
│   ├── openai.go                           # openaiIntegrationGRPCApi — implements OpenAiServiceServer
│   ├── anthropic.go                        # anthropicIntegrationGRPCApi
│   ├── gemini.go                           # geminiIntegrationGRPCApi
│   ├── azure.go                            # azureIntegrationGRPCApi
│   ├── cohere.go                           # cohereIntegrationGRPCApi
│   ├── mistral.go                          # mistralIntegrationGRPCApi
│   ├── huggingface.go                      # huggingfaceIntegrationGRPCApi
│   ├── replicate.go                        # replicateIntegrationGRPCApi
│   ├── vertexai.go                         # vertexaiIntegrationGRPCApi
│   ├── voyageai.go                         # voyageaiIntegrationGRPCApi
│   ├── deepinfra.go                        # deepinfraIntegrationGRPCApi
│   ├── audit_logging.go                    # Audit logging gRPC handler
│   └── health/                             # Health checks
├── internal/caller/                        # Caller implementations
│   ├── callers.go                          # Core interfaces: LargeLanguageCaller, EmbeddingCaller, RerankingCaller, Verifier
│   ├── chat_options.go                     # ChatCompletionOptions, ToolDefinition, FunctionDefinition
│   ├── embedding_options.go                # EmbeddingOptions
│   ├── reranking_options.go                # RerankerOptions
│   ├── metrics/                            # Metrics builder for token usage/latency
│   ├── openai/                             # Example implementation
│   │   ├── openai.go                       # OpenAI base struct (credential resolver, client factory)
│   │   ├── llm.go                          # largeLanguageCaller — GetChatCompletion, StreamChatCompletion
│   │   ├── text-embedding.go              # embeddingCaller — GetEmbedding
│   │   ├── verify-credential.go           # Verifier implementation
│   │   └── moderation.go                  # Content moderation
│   ├── anthropic/                          # Anthropic SDK integration
│   ├── azure/                              # Azure OpenAI endpoint
│   ├── cohere/                             # Cohere SDK (chat + embed + rerank)
│   ├── gemini/                             # Google Gemini
│   ├── huggingface/                        # HuggingFace Inference
│   ├── mistral/                            # Mistral AI
│   ├── replicate/                          # Replicate API
│   ├── vertexai/                           # Google Vertex AI
│   └── voyageai/                           # Voyage AI (embedding + rerank)
├── router/
│   └── provider.go                         # ProviderApiRoute — registers all gRPC service servers
├── config/
│   └── config.go                           # IntegrationConfig
└── migrations/                             # SQL migrations
```

### Backend — Protobuf Service Definitions

```
protos/
├── integration-api.pb.go                   # Shared types (ChatRequest, ChatResponse, etc.)
├── integration-api_grpc.pb.go              # Shared service definitions
├── artifacts/                              # Proto source files (git submodule)
│   ├── integration-api.proto               # Common messages
│   ├── openai-integration.proto            # OpenAiService
│   ├── anthropic-integration.proto         # AnthropicService
│   ├── gemini-integration.proto            # GeminiService
│   └── ...                                 # One proto per provider
```

### UI

```
ui/src/
├── providers/
│   ├── provider.development.json           # Provider registry (featureList: ["text"])
│   ├── index.ts                            # TEXT_PROVIDERS filter + model accessors
│   ├── openai/                             # (no text-models.json — uses inline constants)
│   ├── gemini/
│   │   └── text-models.json                # Gemini model list
│   ├── vertexai/
│   │   └── models.json                     # Vertex AI model list
│   └── azure-foundry/
│       └── text-models.json                # Azure Foundry model list
├── app/components/providers/text/          # Text (LLM) provider config components
│   ├── index.tsx                           # TextProviderConfigComponent + defaults/validation switches
│   ├── openai/                             # OpenAI model config
│   │   ├── index.tsx                       # ConfigureOpenaiTextProviderModel
│   │   └── constants.ts                    # GetOpenaiTextProviderDefaultOptions, models list
│   ├── anthropic/
│   │   ├── index.tsx
│   │   └── constants.ts
│   ├── gemini/
│   │   ├── index.tsx
│   │   └── constants.ts
│   ├── azure-foundry/
│   │   ├── index.tsx
│   │   └── constants.ts
│   ├── cohere/
│   │   ├── index.tsx
│   │   └── constants.ts
│   └── vertexai/
│       ├── index.tsx
│       └── constants.ts
```

## Core Interfaces

### LargeLanguageCaller

```go
// api/integration-api/internal/caller/callers.go
type LargeLanguageCaller interface {
    GetChatCompletion(
        ctx context.Context,
        allMessages []*protos.Message,
        options *ChatCompletionOptions,
    ) (*protos.Message, []*protos.Metric, error)

    StreamChatCompletion(
        ctx context.Context,
        allMessages []*protos.Message,
        options *ChatCompletionOptions,
        onStream func(rID string, msg *protos.Message) error,
        onMetrics func(rID string, msg *protos.Message, mtrx []*protos.Metric) error,
        onError func(rID string, err error),
    ) error
}
```

### EmbeddingCaller

```go
type EmbeddingCaller interface {
    GetEmbedding(
        ctx context.Context,
        content map[int32]string,
        options *EmbeddingOptions,
    ) ([]*protos.Embedding, []*protos.Metric, error)
}
```

### RerankingCaller

```go
type RerankingCaller interface {
    GetReranking(
        ctx context.Context,
        query string,
        content map[int32]string,
        options *RerankerOptions,
    ) ([]*protos.Reranking, []*protos.Metric, error)
}
```

### Verifier

```go
type Verifier interface {
    CredentialVerifier(
        ctx context.Context,
        options *CredentialVerifierOptions,
    ) (*string, error)
}
```

## Handler Pattern

All LLM providers follow the same handler pattern — embed `integrationApi` and delegate to shared methods:

```go
// api/integration-api/api/<provider>.go

type newproviderIntegrationApi struct {
    integrationApi                     // Embeds base (logger, config, storage, auditService)
}

type newproviderIntegrationGRPCApi struct {
    newproviderIntegrationApi          // Embeds REST layer
}

// Constructor returns the protobuf service server interface
func NewProviderGRPC(config *config.IntegrationConfig, logger commons.Logger,
    postgres connectors.PostgresConnector) protos.NewProviderServiceServer {
    return &newproviderIntegrationGRPCApi{
        newproviderIntegrationApi{
            integrationApi: NewInegrationApi(config, logger, postgres),
        },
    }
}
```

### Shared Base Methods

The `integrationApi` base struct provides:

| Method | Purpose |
|---|---|
| `Chat(ctx, request, providerName, caller)` | Non-streaming chat completion |
| `StreamChatBidirectional(ctx, providerName, callerFactory, stream)` | Bidirectional streaming (persistent connection) |
| `Embedding(ctx, request, providerName, caller)` | Embedding generation |
| `Reranking(ctx, request, providerName, caller)` | Reranking |
| `PreHook(ctx, auth, request, requestId, name)` | Audit logging before request (saves request to S3) |
| `PostHook(ctx, auth, request, requestId, name)` | Audit logging after response (saves response + metrics to S3) |
| `RequestId()` | Generates unique Snowflake ID for audit |

## Parameter Key Conventions

LLM model parameters use `model.*` prefix, passed via `ModelParameter` map:

| Key | Type | Description |
|---|---|---|
| `model.name` | string | Model identifier (e.g., "gpt-4o", "claude-3-5-sonnet") |
| `model.temperature` | float | Sampling temperature (0.0–2.0) |
| `model.max_tokens` | int | Maximum output tokens |
| `model.top_p` | float | Nucleus sampling |
| `model.frequency_penalty` | float | Frequency penalty |
| `model.presence_penalty` | float | Presence penalty |
| `model.stop_sequences` | []string | Stop sequences |
| `model.response_format` | string | Output format ("json", "text") |

## Adding a New LLM Provider — Step by Step

### Step 1: Protobuf — Define Service (if needed)

Most providers can reuse the shared proto messages (`ChatRequest`, `ChatResponse`, etc.). If you need a new proto service:

In `protos/artifacts/`, create or extend a proto file:

```protobuf
service NewProviderService {
    rpc Chat(ChatRequest) returns (ChatResponse);
    rpc StreamChat(stream ChatRequest) returns (stream ChatResponse);
    rpc Embedding(EmbeddingRequest) returns (EmbeddingResponse);
    rpc VerifyCredential(VerifyCredentialRequest) returns (VerifyCredentialResponse);
}
```

Generate Go code: `make proto-gen` or `buf generate`.

### Step 2: Backend — Create Caller Package

Create `api/integration-api/internal/caller/newprovider/`:

#### `newprovider.go` — Base Struct + Credential

```go
package internal_newprovider_callers

import (
    "errors"
    internal_callers "github.com/rapidaai/api/integration-api/internal/caller"
    "github.com/rapidaai/pkg/commons"
    "github.com/rapidaai/protos"
)

type NewProvider struct {
    logger     commons.Logger
    credential internal_callers.CredentialResolver
}

var (
    API_KEY = "key"
    API_URL = "url"
)

func newProvider(logger commons.Logger, credential *protos.Credential) NewProvider {
    _credential := credential.GetValue().AsMap()
    return NewProvider{
        logger: logger,
        credential: func() map[string]interface{} {
            return _credential
        },
    }
}

func (np *NewProvider) GetClient() (*ProviderClient, error) {
    credentials := np.credential()
    apiKey, ok := credentials[API_KEY]
    if !ok {
        return nil, errors.New("unable to resolve the credential")
    }
    // Initialize and return provider SDK client
    return NewProviderClient(apiKey.(string)), nil
}
```

#### `llm.go` — LargeLanguageCaller

```go
package internal_newprovider_callers

import (
    "context"
    "time"
    internal_callers "github.com/rapidaai/api/integration-api/internal/caller"
    internal_caller_metrics "github.com/rapidaai/api/integration-api/internal/caller/metrics"
    "github.com/rapidaai/pkg/commons"
    "github.com/rapidaai/pkg/utils"
    "github.com/rapidaai/protos"
)

type largeLanguageCaller struct {
    NewProvider
}

func NewLargeLanguageCaller(logger commons.Logger, credential *protos.Credential) internal_callers.LargeLanguageCaller {
    return &largeLanguageCaller{
        NewProvider: newProvider(logger, credential),
    }
}

func (llc *largeLanguageCaller) GetChatCompletion(
    ctx context.Context,
    allMessages []*protos.Message,
    options *internal_callers.ChatCompletionOptions,
) (*protos.Message, []*protos.Metric, error) {
    startTime := time.Now()

    // 1. Get client
    client, err := llc.GetClient()
    if err != nil {
        return nil, nil, err
    }

    // 2. Convert protos.Message to provider-specific format
    // 3. Parse model.* parameters from options.ModelParameter
    modelName, _ := utils.AnyToString(options.ModelParameter["model.name"])

    // 4. Call provider API
    // 5. Convert response back to protos.Message

    // 6. Build metrics
    metrics := internal_caller_metrics.NewMetricsBuilder().
        WithLatency(time.Since(startTime)).
        WithModel(modelName).
        WithTokenUsage(inputTokens, outputTokens).
        Build()

    // 7. Call pre/post hooks for audit
    if options.PreHook != nil {
        options.PreHook(map[string]interface{}{"messages": allMessages})
    }
    if options.PostHook != nil {
        options.PostHook(map[string]interface{}{"response": responseMsg}, metrics)
    }

    return responseMsg, metrics, nil
}

func (llc *largeLanguageCaller) StreamChatCompletion(
    ctx context.Context,
    allMessages []*protos.Message,
    options *internal_callers.ChatCompletionOptions,
    onStream func(rID string, msg *protos.Message) error,
    onMetrics func(rID string, msg *protos.Message, mtrx []*protos.Metric) error,
    onError func(rID string, err error),
) error {
    // Similar to GetChatCompletion but:
    // 1. Use provider's streaming API
    // 2. For each chunk, call onStream(requestID, partialMessage)
    // 3. On completion, call onMetrics(requestID, finalMessage, metrics)
    // 4. On error, call onError(requestID, err)
    return nil
}
```

#### `verify-credential.go` — Credential Verification

```go
package internal_newprovider_callers

import (
    "context"
    internal_callers "github.com/rapidaai/api/integration-api/internal/caller"
    "github.com/rapidaai/pkg/commons"
    "github.com/rapidaai/protos"
)

type verifyCredentialCaller struct {
    NewProvider
}

func NewVerifyCredentialCaller(logger commons.Logger, credential *protos.Credential) internal_callers.Verifier {
    return &verifyCredentialCaller{
        NewProvider: newProvider(logger, credential),
    }
}

func (v *verifyCredentialCaller) CredentialVerifier(
    ctx context.Context,
    options *internal_callers.CredentialVerifierOptions,
) (*string, error) {
    // Make a lightweight API call (e.g., list models) to verify credentials
    client, err := v.GetClient()
    if err != nil {
        return nil, err
    }
    // Call a simple endpoint to test credential validity
    result := "verified"
    return &result, nil
}
```

#### `embedding.go` — Optional Embedding Support

```go
func NewEmbeddingCaller(logger commons.Logger, credential *protos.Credential) internal_callers.EmbeddingCaller {
    return &embeddingCaller{NewProvider: newProvider(logger, credential)}
}
```

### Step 3: Backend — Create API Handler

Create `api/integration-api/api/newprovider.go`:

```go
package integration_api

import (
    "context"
    config "github.com/rapidaai/api/integration-api/config"
    internal_callers "github.com/rapidaai/api/integration-api/internal/caller"
    internal_newprovider_callers "github.com/rapidaai/api/integration-api/internal/caller/newprovider"
    "github.com/rapidaai/pkg/commons"
    "github.com/rapidaai/pkg/connectors"
    integration_api "github.com/rapidaai/protos"
)

type newproviderIntegrationApi struct {
    integrationApi
}

type newproviderIntegrationGRPCApi struct {
    newproviderIntegrationApi
}

func NewProviderGRPC(config *config.IntegrationConfig, logger commons.Logger,
    postgres connectors.PostgresConnector) integration_api.NewProviderServiceServer {
    return &newproviderIntegrationGRPCApi{
        newproviderIntegrationApi{
            integrationApi: NewInegrationApi(config, logger, postgres),
        },
    }
}

// Chat implements protos.NewProviderServiceServer.
func (np *newproviderIntegrationGRPCApi) Chat(c context.Context,
    irRequest *integration_api.ChatRequest) (*integration_api.ChatResponse, error) {
    return np.integrationApi.Chat(c, irRequest, "NEWPROVIDER",
        internal_newprovider_callers.NewLargeLanguageCaller(np.logger, irRequest.GetCredential()))
}

// StreamChat implements protos.NewProviderServiceServer (bidirectional streaming).
func (np *newproviderIntegrationGRPCApi) StreamChat(
    stream integration_api.NewProviderService_StreamChatServer) error {
    return np.integrationApi.StreamChatBidirectional(
        stream.Context(),
        "NEWPROVIDER",
        func(cred *integration_api.Credential) internal_callers.LargeLanguageCaller {
            return internal_newprovider_callers.NewLargeLanguageCaller(np.logger, cred)
        },
        stream,
    )
}

// Embedding implements protos.NewProviderServiceServer (optional).
func (np *newproviderIntegrationGRPCApi) Embedding(c context.Context,
    irRequest *integration_api.EmbeddingRequest) (*integration_api.EmbeddingResponse, error) {
    return np.integrationApi.Embedding(c, irRequest, "NEWPROVIDER",
        internal_newprovider_callers.NewEmbeddingCaller(np.logger, irRequest.GetCredential()))
}

// VerifyCredential implements protos.NewProviderServiceServer.
func (np *newproviderIntegrationGRPCApi) VerifyCredential(c context.Context,
    irRequest *integration_api.VerifyCredentialRequest) (*integration_api.VerifyCredentialResponse, error) {
    caller := internal_newprovider_callers.NewVerifyCredentialCaller(np.logger, irRequest.GetCredential())
    st, err := caller.CredentialVerifier(c, &internal_callers.CredentialVerifierOptions{})
    if err != nil {
        return &integration_api.VerifyCredentialResponse{Status: false}, nil
    }
    return &integration_api.VerifyCredentialResponse{Status: true, Message: *st}, nil
}
```

### Step 4: Backend — Register in Router

In `api/integration-api/router/provider.go`, add:

```go
protos.RegisterNewProviderServiceServer(S, integrationApi.NewProviderGRPC(Cfg, Logger, Postgres))
```

### Step 5: Backend — Add integration-api Client (for assistant-api)

If assistant-api needs to call the new provider for voice conversations, update the integration client in `pkg/clients/integration/` to add the new provider's gRPC client methods.

### Step 6: UI — Add Provider Metadata

In `ui/src/providers/provider.development.json` (and `provider.production.json`):

```json
{
    "code": "newprovider",
    "name": "New Provider",
    "description": "New LLM provider",
    "image": "newprovider.svg",
    "featureList": ["text", "external"],
    "configurations": [
        { "name": "key", "type": "string", "label": "API Key" }
    ]
}
```

- `"text"` → appears in text/LLM provider dropdown (`TEXT_PROVIDERS`)
- `"external"` → appears in external integrations page (`INTEGRATION_PROVIDER`)
- `configurations` → defines credential fields stored in Vault

### Step 7: UI — Create Model Metadata

Create `ui/src/providers/newprovider/text-models.json`:

```json
[
    {
        "id": "newprovider-model-v1",
        "name": "New Provider Model v1",
        "description": "General purpose model",
        "contextWindow": 128000,
        "maxOutputTokens": 4096
    },
    {
        "id": "newprovider-model-v2",
        "name": "New Provider Model v2",
        "description": "Advanced model",
        "contextWindow": 200000,
        "maxOutputTokens": 8192
    }
]
```

Add accessor in `ui/src/providers/index.ts`:

```typescript
export const NEWPROVIDER_MODEL = () => {
    return require('./newprovider/text-models.json');
};
```

### Step 8: UI — Create Text Provider Config Component

Create `ui/src/app/components/providers/text/newprovider/`:

**`constants.ts`**:
```typescript
import { Metadata } from '@rapidaai/react';

export const NewProviderModels = [
    { id: 'newprovider-model-v1', name: 'Model v1' },
    { id: 'newprovider-model-v2', name: 'Model v2' },
];

export const GetNewProviderTextProviderDefaultOptions = (parameters: Metadata[]): Metadata[] => {
    const defaults = [
        { key: 'model.name', value: 'newprovider-model-v1' },
        { key: 'model.temperature', value: '0.7' },
        { key: 'model.max_tokens', value: '4096' },
    ];
    const existingKeys = new Set(parameters.map(p => p.getKey()));
    const newParams = defaults
        .filter(d => !existingKeys.has(d.key))
        .map(d => { const m = new Metadata(); m.setKey(d.key); m.setValue(d.value); return m; });
    return [...parameters, ...newParams];
};

export const ValidateNewProviderTextProviderDefaultOptions = (parameters: Metadata[]): string | undefined => {
    const model = parameters.find(p => p.getKey() === 'model.name');
    if (!model?.getValue()) return 'Model is required';
    return undefined;
};
```

**`index.tsx`**:
```tsx
import { FC } from 'react';
import { Dropdown } from '@/app/components/dropdown';
import { NewProviderModels } from './constants';

export const ConfigureNewProviderTextProviderModel: FC<ProviderComponentProps> = ({
    parameters,
    onChangeParameter,
}) => {
    return (
        <>
            <Dropdown
                label="Model"
                options={NewProviderModels}
                value={parameters.find(p => p.getKey() === 'model.name')?.getValue()}
                onChange={(val) => onChangeParameter('model.name', val)}
            />
            {/* Temperature, max_tokens, etc. — use shared form components */}
        </>
    );
};
```

### Step 9: UI — Register in Text Provider Switch

In `ui/src/app/components/providers/text/index.tsx`:

1. **Import**:
```typescript
import { ConfigureNewProviderTextProviderModel } from '@/app/components/providers/text/newprovider';
import {
    GetNewProviderTextProviderDefaultOptions,
    ValidateNewProviderTextProviderDefaultOptions,
} from '@/app/components/providers/text/newprovider/constants';
```

2. **Add case** to `GetDefaultTextProviderConfigIfInvalid`:
```typescript
case 'newprovider':
    return GetNewProviderTextProviderDefaultOptions(parameters);
```

3. **Add case** to validation switch

4. **Add case** to render switch:
```typescript
case 'newprovider':
    return (
        <ConfigureNewProviderTextProviderModel
            parameters={parameters}
            onChangeParameter={onChangeParameter}
        />
    );
```

## Data Flow Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│ UI: Assistant Configuration                                      │
│  Select LLM provider → Configure model.* params → Link vault    │
└──────────────────────┬──────────────────────────────────────────┘
                       │ gRPC: CreateAssistant / UpdateAssistant
                       ▼
┌─────────────────────────────────────────────────────────────────┐
│ assistant-api → PostgreSQL                                       │
│  AssistantProviderModel: provider code, model params, credential │
└──────────────────────┬──────────────────────────────────────────┘
                       │ At conversation time
                       ▼
┌─────────────────────────────────────────────────────────────────┐
│ assistant-api Agent Executor                                     │
│  Builds ChatRequest with messages, tools, model params           │
│  ──gRPC──► integration-api                                       │
└──────────────────────┬──────────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────────┐
│ integration-api Handler                                          │
│  1. Authenticate request (gRPC interceptor)                      │
│  2. PreHook → audit log request to S3                            │
│  3. NewLargeLanguageCaller(credential) → resolve API key         │
│  4. StreamChatCompletion(messages, options, callbacks)            │
│     │                                                            │
│     └── Provider SDK → LLM API                                   │
│         ├── onStream(chunk) → gRPC stream → assistant-api        │
│         ├── onMetrics(usage) → audit + response                  │
│         └── onError(err) → error handling                        │
│  5. PostHook → audit log response + metrics to S3                │
└─────────────────────────────────────────────────────────────────┘
```

## Bidirectional Streaming

The `StreamChatBidirectional` method provides persistent connections for voice conversations:

```
assistant-api ──gRPC bidi stream──► integration-api
    │                                     │
    │── ChatRequest (msg 1) ─────────────►│
    │◄── ChatResponse (stream) ───────────│
    │                                     │
    │── ChatRequest (msg 2) ─────────────►│
    │◄── ChatResponse (stream) ───────────│
    │                                     │
    │── EOF ─────────────────────────────►│
```

This avoids reconnection overhead for multi-turn voice conversations.

## Checklist for New LLM Provider

### Protobuf
- [ ] Define or reuse proto service in `protos/artifacts/`
- [ ] Generate Go code (`make proto-gen` or `buf generate`)

### Backend — Caller
- [ ] Create `api/integration-api/internal/caller/<provider>/` directory
- [ ] Implement `<provider>.go` — base struct, credential resolver, client factory
- [ ] Implement `llm.go` — `LargeLanguageCaller` (`GetChatCompletion`, `StreamChatCompletion`)
- [ ] Implement `verify-credential.go` — `Verifier`
- [ ] Implement `embedding.go` — `EmbeddingCaller` (optional)
- [ ] Implement reranking caller (optional, only Cohere/Voyage support this)

### Backend — Handler + Router
- [ ] Create `api/integration-api/api/<provider>.go` — gRPC handler embedding `integrationApi`
- [ ] Implement `Chat`, `StreamChat`, `Embedding`, `VerifyCredential` methods
- [ ] Register in `api/integration-api/router/provider.go` — `protos.Register<Provider>ServiceServer(...)`

### Backend — Integration Client
- [ ] Update `pkg/clients/integration/` if assistant-api needs to call new provider

### UI — Provider Metadata
- [ ] Add provider entry to `provider.development.json` with `"text"` in `featureList`
- [ ] Add provider entry to `provider.production.json`
- [ ] Create `ui/src/providers/<provider>/text-models.json`
- [ ] Add model accessor function in `ui/src/providers/index.ts`

### UI — Config Component
- [ ] Create `ui/src/app/components/providers/text/<provider>/` with `index.tsx` + `constants.ts`
- [ ] Register in `ui/src/app/components/providers/text/index.tsx` — 3 switch statements (defaults, validation, component)
- [ ] Ensure provider `code` in JSON matches protobuf service registration
