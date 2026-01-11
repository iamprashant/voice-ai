# Transformer Architecture & Patterns

## Overview Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                    Assistant API Service                      │
└────────────────────┬────────────────────────────────────────┘
                     │
                     │ Uses
                     ▼
┌─────────────────────────────────────────────────────────────┐
│                   Transformer Package                         │
│  ┌──────────────────────────────────────────────────────┐   │
│  │           Core Interfaces (transformer.go)          │   │
│  │  ┌────────────────────────────────────────────────┐ │   │
│  │  │ Transformers[IN, Opts]                         │ │   │
│  │  │  - Initialize()                               │ │   │
│  │  │  - Transform(ctx, IN, Opts)                   │ │   │
│  │  │  - Close(ctx)                                 │ │   │
│  │  └────────────────────────────────────────────────┘ │   │
│  │                       ▲                              │   │
│  │                       │ implements                   │   │
│  │        ┌──────────────┴──────────────┐             │   │
│  │        │                             │             │   │
│  │  ┌─────────────┐          ┌──────────────────┐    │   │
│  │  │STT          │          │TTS               │    │   │
│  │  │Transformer  │          │Transformer       │    │   │
│  │  │             │          │                  │    │   │
│  │  │Input: []byte│          │Input: string     │    │   │
│  │  │Output: via  │          │Output: via       │    │   │
│  │  │OnTranscript │          │OnSpeech, Complete│   │   │
│  │  │callback     │          │callbacks         │    │   │
│  │  └─────────────┘          └──────────────────┘    │   │
│  └──────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
         │                          │
         │                          │
    ┌────▼────────┐        ┌────────▼──────┐
    │ Provider    │        │ Provider      │
    │ Factories   │        │ Registries    │
    └─────────────┘        └───────────────┘
         │
         │ Creates specific implementations
         │
    ┌────┴──────────────────────────────────────────────┐
    │                                                    │
 ┌──▼────┐ ┌──────┐ ┌──────┐ ┌───────┐ ┌────────┐    │
 │Google │ │Azure │ │Deepgr│ │Assem. │ │Elevens │    │
 │       │ │      │ │am    │ │AI     │ │Labs   │    │
 │STT/TTS│ │STT   │ │STT   │ │STT    │ │TTS    │    │
 └───────┘ └──────┘ └──────┘ └───────┘ └────────┘    │
    │        │         │        │         │           │
    └────────┴─────────┴────────┴─────────┘           │
             │                                         │
             │ Network Calls                          │
             │                                         │
             ▼                                         │
    ┌────────────────────────────────────────┐        │
    │  External AI Service APIs              │        │
    │  (Google Cloud, Azure, Deepgram, etc.) │        │
    └────────────────────────────────────────┘        │
```

## Sequence Diagrams

### Speech-to-Text (STT) Flow

```
Client Code          Transformer             Provider API
     │                    │                        │
     │  Create Instance   │                        │
     ├───────────────────►│                        │
     │                    │                        │
     │  Initialize()      │                        │
     ├───────────────────►│                        │
     │                    ├───────────────────────►│
     │                    │  Establish Connection  │
     │                    │◄───────────────────────┤
     │  ✓ initialized     │                        │
     │◄───────────────────┤                        │
     │                    │                        │
     │  Transform(audio)  │                        │
     ├───────────────────►│                        │
     │                    ├───────────────────────►│
     │                    │  Send Audio Data       │
     │                    │                        │
     │                    │◄───────────────────────┤
     │                    │  Transcript Response   │
     │                    │  (via callback)        │
     │  OnTranscript()    │                        │
     │◄───────────────────┤                        │
     │  (repeat)          │                        │
     │                    │                        │
     │  Close()           │                        │
     ├───────────────────►│                        │
     │                    ├───────────────────────►│
     │                    │  Close Connection     │
     │                    │◄───────────────────────┤
     │  ✓ closed         │                        │
     │◄───────────────────┤                        │
```

### Text-to-Speech (TTS) Flow

```
Client Code          Transformer             Provider API
     │                    │                        │
     │  Create Instance   │                        │
     ├───────────────────►│                        │
     │                    │                        │
     │  Initialize()      │                        │
     ├───────────────────►│                        │
     │                    ├───────────────────────►│
     │                    │  Setup Synthesis      │
     │                    │◄───────────────────────┤
     │  ✓ initialized     │                        │
     │◄───────────────────┤                        │
     │                    │                        │
     │  Transform(text)   │                        │
     ├───────────────────►│                        │
     │                    ├───────────────────────►│
     │                    │  Send Synthesis Req.   │
     │                    │                        │
     │                    │◄───────────────────────┤
     │                    │  Audio Chunk           │
     │  OnSpeech()        │  (via callback)        │
     │◄───────────────────┤                        │
     │  (repeat)          │                        │
     │                    │◄───────────────────────┤
     │                    │  More Audio Chunks    │
     │                    │                        │
     │                    │◄───────────────────────┤
     │  OnComplete()      │  End of Audio         │
     │◄───────────────────┤                        │
     │                    │                        │
     │  Close()           │                        │
     ├───────────────────►│                        │
     │                    ├───────────────────────►│
     │                    │  Close Connection     │
     │                    │◄───────────────────────┤
     │  ✓ closed         │                        │
     │◄───────────────────┤                        │
```

## Component Architecture

### Provider Implementation Structure

```
myprovider/
│
├── myprovider.go
│   └── myproviderOption (configuration holder)
│       ├── NewMyproviderOption()
│       ├── GetEncoding()
│       └── Helper methods (connection strings, etc.)
│
├── stt.go
│   ├── myproviderSpeechToText (implements SpeechToTextTransformer)
│   │   ├── NewMyproviderSpeechToText()
│   │   ├── Initialize()
│   │   ├── Transform()
│   │   ├── Close()
│   │   ├── Name()
│   │   └── speechToTextCallback() (private, handles streaming)
│   └── Goroutine: Listening for provider responses
│
├── tts.go
│   ├── myproviderTextToSpeech (implements TextToSpeechTransformer)
│   │   ├── NewMyproviderTextToSpeech()
│   │   ├── Initialize()
│   │   ├── Transform()
│   │   ├── Close()
│   │   ├── Name()
│   │   └── textToSpeechCallback() (private, handles streaming)
│   └── Goroutine: Listening for provider responses
│
├── internal/
│   └── type.go
│       └── Provider-specific message types
│           ├── TranscriptMessage
│           ├── SynthesisResponse
│           └── etc.
│
└── [optional]
    ├── stt_test.go
    ├── tts_test.go
    └── README.md
```

## Concurrency Model

### Thread Safety Design

```
┌─────────────────────────────────────────┐
│   myproviderSpeechToText (shared)       │
│                                          │
│  ┌────────────────────────────────────┐ │
│  │ Mutex (mu)                          │ │
│  │ ┌──────────────────────────────┐  │ │
│  │ │ Protected Resources:         │  │ │
│  │ │  - client (websocket, gRPC)  │  │ │
│  │ │  - stream (provider stream)  │  │ │
│  │ │  - contextId (current ID)    │  │ │
│  │ └──────────────────────────────┘  │ │
│  └────────────────────────────────────┘ │
│                                          │
│  Unprotected (read-only after init):    │
│  - logger                                │
│  - options                               │
│  - audioConfig                           │
│  - ctx, ctxCancel                        │
└─────────────────────────────────────────┘

Access Pattern:
  Transform():
    1. mu.Lock()
    2. Get reference: client := m.client
    3. mu.Unlock()
    4. Use client (outside lock)
    
  Close():
    1. ctxCancel() - signal goroutines
    2. mu.Lock()
    3. Cleanup: m.client = nil
    4. mu.Unlock()
    
  Goroutine:
    for {
        select {
        case <-ctx.Done(): // Respond to cancellation
            return
        }
        // Process messages
    }
```

## Error Handling Flow

```
Provider Operation
    │
    ├─► Success Path
    │   └─► Call Callback with Result
    │       └─► Handle Callback Error
    │           ├─► Log Error
    │           └─► Continue or Return
    │
    ├─► Network Error
    │   ├─► Log with Context: "myprovider-stt: error: %v"
    │   ├─► Wrap Error: fmt.Errorf("failed to X: %w", err)
    │   └─► Return to Caller
    │
    ├─► Configuration Error
    │   ├─► Validate Early (in New...())
    │   ├─► Clear Error Message
    │   └─► Return Error
    │
    └─► Context Cancellation
        ├─► Clean Exit from Goroutine
        ├─► Release Resources
        └─► No Error Returned
```

## Callback Delivery Pattern

### Safe Callback Invocation

```go
// Pattern: Always check before calling optional callbacks
if m.options != nil {
    if m.options.OnTranscript != nil {
        err := m.options.OnTranscript(transcript, confidence, language, isFinal)
        if err != nil {
            m.logger.Errorf("myprovider-stt: callback error: %v", err)
            // Decide: continue, retry, or abort based on context
        }
    }
}
```

### Callback Context

```
Transcript Callback Parameters:
  - transcript: Full or partial transcribed text
  - confidence: [0.0, 1.0] confidence score
  - language: Language code (e.g., "en", "en-US")
  - isCompleted: true if final, false if interim

Speech Callback Parameters:
  - contextId: Unique ID for synthesis request
  - audioData: []byte chunk of synthesized audio

Completion Callback Parameters:
  - contextId: Unique ID for synthesis request
  - (signals end of synthesis)
```

## State Machine: STT Transformer Lifecycle

```
    ┌──────────────┐
    │   Created    │
    │              │
    └──────┬───────┘
           │ NewMyproviderSpeechToText()
           │ Returns instance, not connected
           ▼
    ┌──────────────┐
    │Initialized? │ No
    │   (ready)   │────► Ready to Initialize
    └──────┬───────┘
           │ Initialize() called
           ▼
    ┌──────────────┐
    │Initializing  │
    │   (setup)    │
    └──────┬───────┘
           │ Connection established
           ▼
    ┌──────────────┐
    │Running       │
    │              │─────────────────┐
    └──────┬───────┘                │
           │                         │
    Transform()                  Close()
    │                            │
    ├─► Send Audio ◄─────────────┼─► Cancel Context
    │                            │   Stop Goroutines
    └─────────────────┬──────────┘   Close Connection
                      │              
                      ▼              
            Callback triggered       
            (OnTranscript)           
                      │              
            ┌─────────┴─────────┐   
            │                   │   
        Interim            Final   
        Result             Result  
        (loop)             (may continue)

    At any point: Close() → Closed state
```

## Data Flow: Audio Processing

```
Audio Input ([]byte)
    │
    ▼
Buffering/Streaming
    │
    ▼
Provider Encoding
    │
    ├─► LINEAR16 → pcm_s16le / linear16 / etc.
    │
    └─► MuLaw8 → mulaw / pcm_mulaw / etc.
    │
    ▼
Send to Provider API
    │
    ├─► WebSocket: stream.Write(data)
    ├─► gRPC: stream.Send(&request)
    └─► REST: HTTP POST with audio
    │
    ▼
Provider Processing
    │
    ▼
Transcript Response
    │
    ▼
Parse/Deserialize
    │
    ▼
Extract Transcript, Confidence, Language
    │
    ▼
Call OnTranscript Callback
    │
    ▼
Client Application
```

## Configuration Hierarchy

```
┌─────────────────────────────────────────┐
│ Runtime Configuration                   │
│ (passed to NewMyprovider...Stt)         │
│ ┌─────────────────────────────────────┐ │
│ │ VaultCredential                     │ │
│ │ ├─ key                              │ │
│ │ ├─ project_id                       │ │
│ │ ├─ endpoint                         │ │
│ │ └─ [provider-specific]              │ │
│ │                                      │ │
│ │ AudioConfig                         │ │
│ │ ├─ audioFormat (LINEAR16/MuLaw8)   │ │
│ │ ├─ sampleRate (8000/16000/etc.)    │ │
│ │ └─ channels                         │ │
│ │                                      │ │
│ │ ModelOptions (utils.Option)         │ │
│ │ ├─ listen.language                  │ │
│ │ ├─ listen.model                     │ │
│ │ ├─ speak.voice.id                   │ │
│ │ ├─ speak.language                   │ │
│ │ └─ [provider-specific options]      │ │
│ └─────────────────────────────────────┘ │
│                                          │
│ Combined in myproviderOption:           │
│ ├─ apiKey (from vault)                 │
│ ├─ audioConfig                         │
│ ├─ mdlOpts                             │
│ └─ [connection config]                 │
└─────────────────────────────────────────┘
```

## Context Flow

```
Parent Context (from caller)
    │
    ▼
context.WithCancel(parentCtx)
    │
    ├─► Used as m.ctx
    │
    ├─► Passed to listening goroutines
    │   └─► Monitors <-m.ctx.Done()
    │
    ├─► Passed to goroutines for cleanup
    │
    └─► Cancelled on Close()
        └─► ctxCancel() signal
            └─► All goroutines receive Done() signal
                └─► Graceful exit
```

## Resource Lifecycle

```
New Instance
    │
    ├─ Memory allocated
    ├─ Fields initialized
    └─ No external resources yet
    │
    ▼
Initialize()
    │
    ├─ Network connection established
    ├─ Listening goroutines spawned
    ├─ Event handlers registered
    └─ State ready for Transform()
    │
    ▼
Transform() × N
    │
    ├─ Audio sent to provider
    ├─ Callbacks invoked with results
    └─ Resources reused (no allocation)
    │
    ▼
Close()
    │
    ├─ Context cancelled
    ├─ Goroutines signaled to exit
    ├─ Network connection closed
    ├─ Event handlers unregistered
    ├─ Resources released
    └─ Safe to discard instance

Cleanup Guarantee: Resources released even if Close() errors
```

---

## Best Practices Summary

1. **Mutex**: Always protect shared state (client, stream, etc.)
2. **Context**: Cancel context before cleanup, monitor in goroutines
3. **Callbacks**: Check nil, handle errors, log failures
4. **Errors**: Log + wrap with context, never silent failures
5. **Logging**: Prefix with provider name for traceability
6. **Resources**: Clean up in Close() even on error
7. **Credentials**: Only from vault, never hardcoded
8. **Goroutines**: Exit cleanly on context cancellation

---

See [DEVELOPMENT.md](DEVELOPMENT.md) for implementation details and code examples.
