# Transformer Package - Quick Reference

## What is a Transformer?

A Transformer is a provider-agnostic adapter that converts audio between different formats or processes speech/text using various AI services (Google, Azure, Deepgram, etc.).

## Quick Start: Adding a New Provider (5 Steps)

### 1️⃣ Create Directory
```bash
mkdir -p api/assistant-api/internal/transformer/myprovider/internal
```

### 2️⃣ Create Configuration (`myprovider/myprovider.go`)
```go
package internal_transformer_myprovider

type myproviderOption struct {
    logger      commons.Logger
    apiKey      string
    audioConfig *protos.AudioConfig
    mdlOpts     utils.Option
}

func NewMyproviderOption(...) (*myproviderOption, error) {
    // Extract credentials from vault
    // Validate configuration
    // Return new instance
}
```

### 3️⃣ Implement STT (`myprovider/stt.go`)
```go
type myproviderSpeechToText struct {
    *myproviderOption
    // Add client, context, mutex, logger, options
}

func (m *myproviderSpeechToText) Initialize() error { /* Setup connection */ }
func (m *myproviderSpeechToText) Transform(ctx context.Context, audioData []byte, opts *SpeechToTextOption) error { /* Send audio */ }
func (m *myproviderSpeechToText) Close(ctx context.Context) error { /* Cleanup */ }
func (m *myproviderSpeechToText) Name() string { return "myprovider-speech-to-text" }
```

### 4️⃣ Implement TTS (`myprovider/tts.go`)
```go
type myproviderTextToSpeech struct {
    *myproviderOption
    // Add client, context, mutex, logger, options
}

func (m *myproviderTextToSpeech) Initialize() error { /* Setup connection */ }
func (m *myproviderTextToSpeech) Transform(ctx context.Context, text string, opts *TextToSpeechOption) error { /* Send text */ }
func (m *myproviderTextToSpeech) Close(ctx context.Context) error { /* Cleanup */ }
func (m *myproviderTextToSpeech) Name() string { return "myprovider-text-to-speech" }
```

### 5️⃣ Add Tests
```go
// myprovider/stt_test.go & tts_test.go
// Test initialization, transformation, and cleanup
```

---

## Key Interfaces

### Speech-to-Text Transformer
**Input:** `[]byte` (audio data)  
**Output:** Calls `OnTranscript(transcript, confidence, language, isCompleted)`

### Text-to-Speech Transformer
**Input:** `string` (text)  
**Output:** Calls `OnSpeech(contextId, audioData)` and `OnComplete(contextId)`

---

## Common Patterns

### Get Configuration from Vault
```go
credMap := vaultCredential.GetValue().AsMap()
apiKey, ok := credMap["api_key"]
if !ok {
    return nil, fmt.Errorf("api_key not found in vault")
}
```

### Thread-Safe Client Access
```go
m.mu.Lock()
client := m.client
m.mu.Unlock()

// Use client
```

### Call Callback with Error Handling
```go
if m.options.OnTranscript != nil {
    if err := m.options.OnTranscript(text, conf, lang, true); err != nil {
        m.logger.Errorf("callback error: %v", err)
    }
}
```

### Proper Logging
```go
m.logger.Debugf("myprovider-stt: connection established")
m.logger.Errorf("myprovider-stt: failed to send: %v", err)
```

### Context Cancellation in Goroutine
```go
go func(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        default:
            // Process messages
        }
    }
}(m.ctx)
```

---

## Credential Configuration

Credentials are stored in vault with provider-specific keys:

### Example Vault Structure
```json
{
    "key": "api-key-value",
    "project_id": "my-project",
    "subscription_key": "sub-key",
    "endpoint": "https://api.example.com",
    "service_account_key": "json-string"
}
```

### Extract in Code
```go
keyValue, ok := credMap["key"]
if !ok {
    return nil, fmt.Errorf("key not found")
}
key := keyValue.(string)
```

---

## Model Options Configuration

Providers support dynamic configuration through `ModelOptions`:

### STT Options Example
```go
if language, err := m.mdlOpts.GetString("listen.language"); err == nil {
    // Use language
}
if model, err := m.mdlOpts.GetString("listen.model"); err == nil {
    // Use model
}
```

### TTS Options Example
```go
if voiceId, err := m.mdlOpts.GetString("speak.voice.id"); err == nil {
    // Use voice ID
}
if speed, err := m.mdlOpts.GetString("speak.speed"); err == nil {
    // Use speed
}
```

---

## Testing Checklist

- [ ] Unit test for NewMyproviderOption
- [ ] Unit test for STT Initialize/Transform/Close
- [ ] Unit test for TTS Initialize/Transform/Close
- [ ] Test callback invocation with correct parameters
- [ ] Test callback error handling
- [ ] Test context cancellation cleanup
- [ ] Test nil callback safety
- [ ] Integration test with real provider (if applicable)

---

## Existing Providers Reference

### Google Cloud
- **Files:** `google/google.go`, `google/stt.go`, `google/tts.go`
- **Strength:** Well-structured, proper error handling, streaming support
- **Use for:** Architecture reference

### Deepgram
- **Files:** `deepgram/deepgram.go`, `deepgram/stt.go`
- **Strength:** WebSocket streaming, callback pattern
- **Use for:** WebSocket implementation patterns

### Azure
- **Files:** `azure/azure.go`, `azure/stt.go`, `azure/tts.go`
- **Strength:** Event-based callbacks, error handling
- **Use for:** Event-driven patterns

### AssemblyAI
- **Files:** `assembly-ai/assemblyai.go`, `assembly-ai/stt.go`
- **Strength:** WebSocket configuration
- **Use for:** WebSocket query parameter patterns

---

## Troubleshooting

### Issue: "Connection not initialized"
→ Check if `Initialize()` was called and returned without error

### Issue: Callbacks not triggered
→ Verify callback is not nil before calling  
→ Check if listening goroutine is running  
→ Review error logs in callback processing

### Issue: Memory leaks
→ Verify `Close()` cancels context  
→ Ensure all goroutines listen to context.Done()  
→ Check for closed channels being written to

### Issue: Race condition errors
→ Add `go test -race` to your testing  
→ Ensure mutex protects all shared state access  
→ Never hold mutex across blocking operations

---

## File Structure Template

```
transformer/
├── myprovider/
│   ├── myprovider.go       # Configuration and helper methods
│   ├── stt.go              # Speech-to-Text implementation
│   ├── tts.go              # Text-to-Speech implementation
│   ├── stt_test.go         # STT unit tests
│   ├── tts_test.go         # TTS unit tests
│   ├── internal/
│   │   └── type.go         # Provider-specific types
│   └── README.md           # Provider-specific documentation
```

---

## Code Template

See [DEVELOPMENT.md](DEVELOPMENT.md) for complete implementation template with detailed explanations.

---

## Contact & Support

For questions about transformer development:
1. Review existing implementations (Google, Deepgram, Azure)
2. Check [DEVELOPMENT.md](DEVELOPMENT.md) for detailed guide
3. Refer to test files for usage examples
