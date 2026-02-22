# Voice Agent Pipeline — Architecture & Implementation Guide

## Overview

The voice agent pipeline is the core of the assistant-api service. It orchestrates real-time voice conversations by routing audio/text packets through a chain of components: STT, VAD, LLM, TTS, and back to the client.

## High-Level Flow

```
Client (WebRTC/SIP/gRPC/WebPlugin)
  ↕ BaseStreamer (buffered I/O channels, audio resampling, 20ms frame output)
    ↕ Talk() loop — infinite recv loop dispatching to OnPacket()
      ↕ OnPacket() — THE CENTRAL PACKET ROUTER
        ↓ Audio Path: Denoiser → VAD → STT → EndOfSpeech → LLM Executor → TextAggregator → TTS → Output
        ↓ Text Path: EndOfSpeech (or direct) → LLM Executor → TextAggregator → TTS (audio) / Direct (text)
```

## Directory Structure

```
api/assistant-api/internal/
├── adapters/                     # Request adapters — the orchestration layer
│   ├── request.go                # Entry point: creates Talker from request context
│   ├── customizers/
│   │   └── messaging.go          # State machine (Interrupt → LLMGenerating → LLMGenerated)
│   └── internal/
│       ├── generic.go            # Core generic adapter struct holding all components
│       ├── assistant_generic.go  # Assistant loading & configuration
│       ├── behaviors_generic.go  # Greeting, idle timeout, max session, error handling
│       ├── callback_generic.go   # OnPacket() — central packet router (~493 lines)
│       ├── hook_generic.go       # Webhook & analysis lifecycle hooks
│       ├── io.go                 # I/O notification (send packets to client)
│       ├── knowledge_generic.go  # RAG/knowledge retrieval (hybrid/semantic/text search)
│       ├── log_generic.go        # Webhook & tool execution logging
│       ├── session_generic.go    # Connect()/Disconnect() — session lifecycle
│       └── talking.go            # Talk() — main conversation loop
├── agent/                        # AI execution layer
│   ├── embedding/                # Query embedding for RAG
│   │   ├── embeddings.go         # Embedding interface
│   │   └── query_embedding.go    # Integration-api gRPC embedding client
│   ├── executor/                 # LLM execution
│   │   ├── executor.go           # AssistantExecutor interface
│   │   ├── llm/
│   │   │   ├── llm.go            # LLM executor factory (MODEL/AGENTKIT/WEBSOCKET)
│   │   │   └── internal/
│   │   │       ├── model/        # MODEL — persistent gRPC stream to integration-api
│   │   │       ├── agentkit/     # AGENTKIT — external gRPC agent server
│   │   │       └── websocket/    # WEBSOCKET — external WebSocket agent
│   │   └── tool/
│   │       ├── tool.go           # ToolExecutor — local + MCP tool execution
│   │       └── internal/         # Tool implementations
│   └── reranker/                 # Result reranking
│       ├── reranking.go          # Reranker interface
│       └── text_reranking.go     # Integration-api gRPC reranking client
├── aggregator/text/              # Text stream aggregation (sentence assembly)
├── audio/                        # Audio config, recorder, resampler
├── callcontext/                  # Redis-backed call context store (5-min TTL)
├── capturers/                    # S3 audio/text capture for recording
├── channel/                      # Transport layer
│   ├── base/base_streamer.go     # Transport-agnostic buffered I/O (20ms frames)
│   ├── grpc/streamer.go          # gRPC bidirectional streaming
│   ├── telephony/                # SIP/WebSocket/AudioSocket telephony
│   └── webrtc/                   # WebRTC + Pion (Opus 48kHz ↔ PCM 16kHz)
├── denoiser/                     # Audio noise reduction (Krisp/RNNoise)
├── end_of_speech/                # Silence-based end-of-speech detection
├── normalizers/                  # Text normalization pipeline (URL, currency, date, etc.)
├── telemetry/                    # OpenTelemetry-style voice agent tracing
├── transformer/                  # STT/TTS provider adapters (12 providers)
├── type/                         # Core interfaces (16 files)
└── vad/                          # Voice Activity Detection (Silero)
```

## Core Components

### 1. Talk Loop (`talking.go`)

The infinite receive loop. Dispatches protobuf messages from the client:

| Message Type | Action |
|---|---|
| `ConversationInitialization` | `Connect()` — session creation, component init |
| `ConversationUserMessage` | `OnPacket()` with `UserAudioPacket` or `UserTextPacket` |
| `ConversationConfiguration` | Reconfigure STT/TTS/mode mid-call |
| `ConversationDisconnection` | Store disconnect reason |
| EOF | `Disconnect()` |

### 2. Session Lifecycle (`session_generic.go`)

**Connect** (called on `ConversationInitialization`):
1. Authenticate → Fetch assistant → Create/resume conversation in DB
2. Concurrent init via `errgroup`: LLM executor, TTS + behavior (greeting/idle/session timers), text aggregator
3. Background init: STT, recorder, end-of-speech, metrics, client info, webhooks

**Disconnect** (5 phases):
1. Close STT + EOS and TTS + aggregator concurrently
2. Fire `OnEndConversation` hooks (webhooks + analyses)
3. Persist recording to S3
4. End tracing span
5. Export telemetry to OpenSearch, close executor, stop timers

### 3. Central Packet Router — `OnPacket()` (`callback_generic.go`)

The ~493-line switch statement that routes **all** pipeline packets. This is the heart of the agent:

| Packet | Action |
|---|---|
| `UserAudioPacket` | Denoise → Record → VAD → STT |
| `UserTextPacket` | Interrupt → Assign context ID → EndOfSpeech |
| `SpeechToTextPacket` | Assign context ID → EndOfSpeech analysis |
| `EndOfSpeechPacket` | Stop idle timer → LLMGenerating state → Create message → **Execute LLM** |
| `InterimEndOfSpeechPacket` | Notify user with incomplete transcript |
| `LLMResponseDeltaPacket` | Check stale context → Send to text aggregator |
| `LLMResponseDonePacket` | Check stale → Start idle timer → LLMGenerated → Create message → Aggregator |
| `StaticPacket` | Create message → Execute → Aggregator (for greetings/static responses) |
| `InterruptionPacket` | Reset idle timer → EOS → Interrupt all providers → State transition |
| `TextToSpeechAudioPacket` | Reset idle timer → Notify audio chunk → Record |
| `TextToSpeechEndPacket` | Notify completion |
| `DirectivePacket` | END_CONVERSATION notification |
| `ConversationMetricPacket` / `MetadataPacket` | Async persistence |

### 4. State Machine — Messaging (`messaging.go`)

States: `Unknown(1)` → `Interrupt(6)` → `Interrupted(7)` → `LLMGenerating(8)` → `LLMGenerated(5)`

Key rules:
- Each interruption generates a **new UUID context ID** — enables stale-context detection
- LLM deltas check `ContextId != messaging.GetID()` to discard stale responses
- Modes: `TextMode` / `AudioMode` — determines whether TTS is invoked

### 5. Audio Pipeline Components

| Component | Interface File | Implementation | Purpose |
|---|---|---|---|
| **Denoiser** | `type/denoiser.go` | Krisp / RNNoise | `Denoise(ctx, []byte) → ([]byte, float64, error)` |
| **VAD** | `type/vad.go` | Silero | `Process(ctx, UserAudioPacket)` — emits `InterruptionPacket` |
| **STT** | `type/stt_transformer.go` | 12 providers | `Transform(ctx, UserAudioPacket)` → emits `SpeechToTextPacket` |
| **EndOfSpeech** | `type/end_of_speech.go` | Silence-based | `Analyze(ctx, Packet)` → emits `EndOfSpeechPacket` |
| **TextAggregator** | `type/aggregator.go` | Sentence assembly | `Aggregate(ctx, ...LLMPacket)` + `Result() <-chan Packet` |
| **TTS** | `type/tts_transformer.go` | 12 providers | `Transform(ctx, LLMPacket)` → emits `TextToSpeechAudioPacket` |
| **Recorder** | `type/recorder.go` | S3 capturer | `Record(ctx, Packet)` + `Persist() → ([]byte, []byte)` |
| **Resampler** | `type/resampler.go` | Audio converter | Sample rate/channel/format conversion |
| **Normalizer** | `type/normalizer.go` | Pipeline | URL, currency, date, time, number, symbol normalizers |

### 6. LLM Executors (`agent/executor/`)

Three executor types selected by `AssistantProvider` enum:

#### MODEL (`agent/executor/llm/internal/model/`)
- Opens a **persistent bidirectional gRPC stream** to integration-api
- `Initialize()`: Fetch credential from vault + init tool executor concurrently → open `StreamChat` → start `listen()` goroutine
- `Execute()`: Send user text/static packets via stream with full conversation history
- `listen()`: Read streaming responses → dispatch:
  - Has metrics → `LLMResponseDonePacket` + execute tool calls if present
  - No metrics → `LLMResponseDeltaPacket` (streaming delta)
  - Error → `LLMErrorPacket`
- Tool loop: `executeToolCalls()` → `toolExecutor.ExecuteAll()` → re-send via `chat()` with tool results

#### AGENTKIT (`agent/executor/llm/internal/agentkit/`)
- Connects to **external gRPC server** (user's custom agent) via `protos.AgentKitClient.Talk()`
- Supports TLS with custom certificates or insecure mode
- Sends `ConversationInitialization` → user messages as `TalkInput`
- Receives `TalkOutput`: text deltas/completions, interruptions, directives
- No local tool execution — the external agent handles everything

#### WEBSOCKET (`agent/executor/llm/internal/websocket/`)
- Connects to **external WebSocket server** (JSON protocol)
- Message types: `configuration`, `user_message`, `stream`, `complete`, `tool_call`, `interruption`, `close`, `ping`/`pong`
- Handles streaming text chunks + completion with metrics

### 7. Tool Execution (`agent/executor/tool/`)

`ToolExecutor` manages **local** and **MCP** tools:

**Local tools:**
- `knowledge_retrieval` — RAG search against knowledge bases
- `api_request` — HTTP calls to external APIs
- `endpoint_request` — Invoke Rapida endpoints
- `end_of_conversation` — Terminate conversation

**MCP tools:** External MCP servers, dynamically discovered via `ListTools()`.

`ExecuteAll()` runs all tool calls **concurrently** via goroutines.

### 8. RAG/Knowledge Retrieval (`knowledge_generic.go`)

Three search methods:
- `hybrid-search` — Embedding + keyword (via integration-api → OpenSearch `HybridSearch()`)
- `semantic-search` — Embedding only (via `VectorSearch()`)
- `text-search` — Keyword only (via `TextSearch()`)

Flow: Query → Embedding (integration-api gRPC) → OpenSearch → Reranking (integration-api gRPC) → Context

### 9. Transport Layer — BaseStreamer (`channel/base/base_streamer.go`)

Transport-agnostic buffered I/O:
- `InputCh` / `OutputCh` channels with non-blocking push
- Input: Accumulates + resamples audio → flushes at threshold
- Output: Accumulates TTS audio → flushes fixed **20ms frames** via `sync.Pool` frame reuse
- `ClearInputBuffer()` / `ClearOutputBuffer()` for interruption handling
- Extended by WebRTC, telephony, and gRPC streamers

### 10. Behavior System (`behaviors_generic.go`)

- **Greeting**: Templated initial message sent via `OnPacket(StaticPacket)`
- **Idle Timeout**: Timer with backoff retries, can auto-end conversation
- **Max Session Duration**: Hard time limit on conversation
- **Error Handling**: Default or configured error message

### 11. Webhook & Analysis Hooks (`hook_generic.go`)

Lifecycle events: `OnBeginConversation`, `OnResumeConversation`, `OnErrorConversation`, `OnEndConversation`

- **Analysis**: Post-conversation endpoint invocation → stores results as metadata
- **Webhooks**: HTTP calls with retry logic + structured argument building

## Packet Flow Diagram (Audio Mode)

```
User speaks → BaseStreamer.BufferAndSendInput()
  → InputCh → Talk() recv loop
    → OnPacket(UserAudioPacket)
      → Denoiser.Denoise() → OnPacket(denoised UserAudioPacket)
        → Recorder.Record()
        → VAD.Process() [may emit InterruptionPacket]
        → STT.Transform()
          → OnPacket(SpeechToTextPacket)
            → EndOfSpeech.Analyze()
              → OnPacket(EndOfSpeechPacket)
                → assistantExecutor.Execute()
                  → gRPC stream to integration-api
                    → OnPacket(LLMResponseDeltaPacket) [streaming]
                      → TextAggregator.Aggregate()
                        → Result() channel → callSpeaking()
                          → TTS.Transform()
                            → OnPacket(TextToSpeechAudioPacket)
                              → Notify() → streamer.Send()
                                → OutputCh → BaseStreamer.BufferAndSendOutput()
                                  → 20ms frames → Client
                    → OnPacket(LLMResponseDonePacket) [final]
                      → Create message → TextAggregator (flush)
```

## Key Design Decisions

1. **Single `OnPacket()` dispatch**: All pipeline stages communicate through a single event bus, enabling loose coupling and easy interruption handling
2. **Context ID staleness checks**: Each interruption generates a new UUID; LLM deltas check `ContextId != messaging.GetID()` to discard stale responses
3. **Persistent bidirectional streams**: The MODEL executor maintains a persistent gRPC stream for the entire conversation duration (not per-turn)
4. **Concurrent initialization**: `errgroup` throughout for parallel setup of STT, TTS, VAD, denoiser, credentials, tools
5. **Three executor abstractions**: MODEL, AGENTKIT, WEBSOCKET — unified behind `AssistantExecutor` interface
6. **20ms output framing**: `BaseStreamer` buffers TTS audio and outputs precisely 20ms frames for smooth playback
7. **sync.Pool frame reuse**: Hot-path output frames are pooled to reduce GC pressure
8. **Redis-backed call context**: Atomic get-and-delete via Lua scripts with 5-minute TTL for telephony session handoff

## Core Interfaces (`type/` directory)

| Interface | File | Key Methods |
|---|---|---|
| `Communication` | `communication.go` | `Auth()`, `Source()`, `Tracer()`, `Assistant()`, `Conversation()`, `GetBehavior()`, `GetHistories()` |
| `Callback` | `callback.go` | LLM response callbacks |
| `Transformers[IN]` | `transformer.go` | `Initialize()`, `Transform(ctx, IN)`, `Close(ctx)` |
| `SpeechToTextTransformer` | `stt_transformer.go` | Extends `Transformers[UserAudioPacket]` |
| `TextToSpeechTransformer` | `tts_transformer.go` | Extends `Transformers[LLMPacket]` |
| `VAD` | `vad.go` | `Process(ctx, UserAudioPacket)` |
| `EndOfSpeech` | `end_of_speech.go` | `Analyze(ctx, Packet)` |
| `Streamer` | `streamer.go` | `Send()`, `Recv()`, `Close()` |
| `Aggregator` | `aggregator.go` | `Aggregate(ctx, ...LLMPacket)`, `Result() <-chan Packet` |
| `Normalizer` | `normalizer.go` | `Normalize(text string) string` |
| `Denoiser` | `denoiser.go` | `Denoise(ctx, []byte) ([]byte, float64, error)` |
| `Resampler` | `resampler.go` | Sample rate/channel/format conversion |
| `Recorder` | `recorder.go` | `Record(ctx, Packet)`, `Persist() ([]byte, []byte)` |
