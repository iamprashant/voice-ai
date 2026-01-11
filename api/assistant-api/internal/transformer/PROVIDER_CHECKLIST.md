# New Provider Implementation Checklist

Use this checklist when adding support for a new speech-to-text or text-to-speech provider.

## Pre-Implementation

- [ ] Understand provider's API documentation
- [ ] Identify supported audio formats and encoding
- [ ] Note required credentials/authentication method
- [ ] Review provider's streaming capabilities (if needed)
- [ ] Plan error handling and retry logic
- [ ] Identify any regional/endpoint variations

## Directory & File Setup

- [ ] Created `api/assistant-api/internal/transformer/myprovider/` directory
- [ ] Created `myprovider/myprovider.go` for configuration
- [ ] Created `myprovider/stt.go` for Speech-to-Text (if supported)
- [ ] Created `myprovider/tts.go` for Text-to-Speech (if supported)
- [ ] Created `myprovider/internal/` directory for internal types
- [ ] Created `myprovider/internal/type.go` for provider-specific structures

## Configuration Implementation (`myprovider.go`)

- [ ] Created `myproviderOption` struct with required fields:
  - [ ] `logger commons.Logger`
  - [ ] `audioConfig *protos.AudioConfig`
  - [ ] `mdlOpts utils.Option`
  - [ ] Provider-specific fields (API key, endpoints, etc.)

- [ ] Implemented `NewMyproviderOption()` constructor:
  - [ ] Extract credentials from vault using `vaultCredential.GetValue().AsMap()`
  - [ ] Validate all required credentials are present
  - [ ] Return error with descriptive message if validation fails
  - [ ] Return configured option struct

- [ ] Implemented `GetEncoding()` helper method:
  - [ ] Map `protos.AudioConfig_LINEAR16` to provider's format
  - [ ] Map `protos.AudioConfig_MuLaw8` to provider's format
  - [ ] Handle unknown formats with sensible default
  - [ ] Log warning if using default encoding

- [ ] Implemented provider-specific helper methods:
  - [ ] Connection string builders
  - [ ] Configuration generators
  - [ ] Format converters

## Speech-to-Text Implementation (`stt.go`)

### Structure Definition
- [ ] Created `myproviderSpeechToText` struct with:
  - [ ] `*myproviderOption` (embedded)
  - [ ] `mu sync.Mutex` for thread safety
  - [ ] `ctx context.Context` and `ctxCancel context.CancelFunc`
  - [ ] `logger commons.Logger`
  - [ ] `options *SpeechToTextInitializeOptions`
  - [ ] Provider client/connection field

### Constructor Function
- [ ] `NewMyproviderSpeechToText()` function signature matches interface
- [ ] Calls `NewMyproviderOption()` to initialize configuration
- [ ] Creates child context with `context.WithCancel(parentCtx)`
- [ ] Returns proper error on configuration failure

### Required Methods

#### `Initialize() error`
- [ ] Establishes connection to provider
- [ ] Sets up any necessary event handlers/callbacks
- [ ] Starts listening goroutine(s) if needed
- [ ] Logs successful initialization
- [ ] Returns error with context on failure

#### `Transform(ctx, audioData, opts) error`
- [ ] Locks mutex before accessing shared state
- [ ] Checks if properly initialized (client != nil)
- [ ] Sends audio data to provider
- [ ] Handles send errors appropriately
- [ ] Returns error if needed
- [ ] Unlocks mutex before returning

#### `Close(ctx context.Context) error`
- [ ] Cancels context: `m.ctxCancel()`
- [ ] Locks mutex
- [ ] Closes provider connection
- [ ] Cleans up goroutines
- [ ] Clears shared state
- [ ] Logs closure
- [ ] Returns any cleanup errors

#### `Name() string`
- [ ] Returns unique identifier (e.g., "myprovider-speech-to-text")
- [ ] Must match service discovery/registration name

### Callback Handling
- [ ] Checks if `m.options` is not nil before calling
- [ ] Checks if specific callback function is not nil
- [ ] Passes correct parameters to callback:
  - [ ] `transcript string`: transcribed text
  - [ ] `confidence float64`: confidence score (0-1)
  - [ ] `language string`: detected language code
  - [ ] `isCompleted bool`: whether transcript is final
- [ ] Handles callback errors (logs and returns appropriately)

### Goroutine Management
- [ ] Listening goroutine monitors `ctx.Done()` channel
- [ ] Exits cleanly when context is cancelled
- [ ] Properly handles connection drops/errors
- [ ] No goroutine leaks on Close()

### Error Handling
- [ ] All errors logged with provider prefix (e.g., "myprovider-stt:")
- [ ] Errors wrapped with context: `fmt.Errorf("message: %w", err)`
- [ ] Descriptive error messages
- [ ] Distinguishes between setup errors and runtime errors

## Text-to-Speech Implementation (`tts.go`)

Follow same structure as STT, but:

### Transform Method
- [ ] Takes `text string` input instead of `[]byte`
- [ ] Sends synthesis request to provider
- [ ] Manages `contextId` for tracking synthesis sessions

### Callback Handling
- [ ] Calls `OnSpeech(contextId, audioData)` when audio chunks arrive
- [ ] Calls `OnComplete(contextId)` when synthesis finishes
- [ ] Handles callback errors appropriately

### Additional Considerations
- [ ] Manages multiple concurrent synthesis requests (if needed)
- [ ] Properly tracks context IDs for each request
- [ ] Handles incomplete synthesis gracefully

## Testing Implementation

### Unit Test File (`stt_test.go` or `tts_test.go`)
- [ ] Created test file in provider directory
- [ ] Tests follow Go testing conventions

### Test Coverage

#### Initialization Tests
- [ ] `TestNewMyproviderOption_Success`: Valid credentials
- [ ] `TestNewMyproviderOption_MissingCredential`: Missing required credential
- [ ] `TestNewMyproviderOption_InvalidCredential`: Wrong credential type

#### Transformer Creation Tests
- [ ] `TestNewMyproviderSpeechToText_Success`: Successful creation
- [ ] `TestNewMyproviderSpeechToText_InvalidConfig`: Invalid configuration

#### Lifecycle Tests
- [ ] `TestInitialize_Success`: Connection established
- [ ] `TestInitialize_Failure`: Handle connection errors
- [ ] `TestClose_Cleanup`: Proper resource cleanup
- [ ] `TestClose_GoroutineCleanup`: No goroutine leaks

#### Functionality Tests
- [ ] `TestTransform_NotInitialized`: Error when not initialized
- [ ] `TestTransform_Success`: Audio sent successfully
- [ ] `TestTransform_WithCallback`: Callback invoked correctly
- [ ] `TestTransform_CallbackError`: Errors handled

#### Concurrency Tests
- [ ] `TestConcurrentTransform`: Multiple simultaneous calls
- [ ] `TestConcurrentTransformAndClose`: Close during Transform
- [ ] Run with `go test -race` to detect race conditions

### Mock/Integration Tests
- [ ] Mock provider client for unit tests (if needed)
- [ ] Integration tests with real provider (optional, mark as skippable)
- [ ] Test with actual audio samples (if possible)

## Code Quality

### Naming Conventions
- [ ] Package name: `internal_transformer_myprovider`
- [ ] Internal subpackage: `myprovider_internal`
- [ ] Function names are descriptive: `NewMyprovider...`
- [ ] Receiver names are short (single letter or abbreviation)
- [ ] Private struct fields are unexported (lowercase)

### Documentation
- [ ] Copyright header on all files:
  ```go
  // Copyright (c) 2023-2025 RapidaAI
  // Author: Prashant Srivastav <prashant@rapida.ai>
  //
  // Licensed under GPL-2.0 with Rapida Additional Terms.
  // See LICENSE.md or contact sales@rapida.ai for commercial usage.
  ```
- [ ] Package-level godoc comment
- [ ] Function godoc comments for all exported functions
- [ ] Inline comments for complex logic
- [ ] No commented-out code (unless explaining implementation)

### Code Standards
- [ ] Consistent indentation (4 spaces or tabs - check existing code)
- [ ] Error messages start with lowercase (unless constant)
- [ ] No unused imports
- [ ] No exported fields that should be private
- [ ] Proper interface satisfaction (compile check)

### Concurrency Safety
- [ ] Mutex protects all shared state
- [ ] No double-locking
- [ ] Lock held for minimal duration
- [ ] No blocking operations while holding lock
- [ ] Context passed correctly to goroutines

## Credential Management

- [ ] Documented required vault credentials
- [ ] Validated all credentials in constructor
- [ ] Clear error messages for missing credentials
- [ ] No hardcoded credentials (check with `grep -r "YOUR_API_KEY"`)
- [ ] No credentials in logs

## Audio Format Support

- [ ] Documented supported audio formats
- [ ] Implemented `GetEncoding()` for all supported formats
- [ ] Default to safe format if unknown
- [ ] Handle sample rate correctly
- [ ] Handle channel configuration (mono/stereo)

## Configuration Options

- [ ] Documented all `mdlOpts` keys the provider supports
- [ ] Safe access with error handling: `if val, err := mdlOpts.GetString(...); err == nil`
- [ ] Sensible defaults for optional parameters
- [ ] Log when using defaults

## Integration

- [ ] Provider registered in service factory (if applicable)
- [ ] Can be discovered by service discovery mechanism
- [ ] Documented provider credentials structure
- [ ] Examples provided for vault configuration
- [ ] Added to provider list in documentation

## Documentation

- [ ] Created `myprovider/README.md` with:
  - [ ] Provider name and link
  - [ ] Supported features (STT/TTS)
  - [ ] Audio formats supported
  - [ ] Configuration requirements
  - [ ] Example vault configuration
  - [ ] Known limitations or gotchas
  - [ ] Cost/quota information (if applicable)

- [ ] Updated main `transformer/README.md` with provider reference
- [ ] Added provider to architecture documentation if needed
- [ ] Examples of using the provider in integration code

## Security Review

- [ ] No hardcoded secrets
- [ ] Credentials extracted from vault only
- [ ] API keys not logged
- [ ] No credentials in error messages
- [ ] Sensitive data cleared on Close()
- [ ] HTTPS/TLS used for all connections (if applicable)

## Performance Considerations

- [ ] Streaming support implemented (not buffering entire audio)
- [ ] Efficient error handling (no unnecessary retries)
- [ ] Proper timeout configuration (not infinite waits)
- [ ] Resource cleanup to prevent leaks
- [ ] No unnecessary allocations in hot paths

## Edge Cases & Error Handling

- [ ] Network timeout handling
- [ ] Connection drop handling
- [ ] Malformed response handling
- [ ] Empty audio/text handling
- [ ] Concurrent close and transform
- [ ] Context cancellation during various states
- [ ] Provider rate limiting/quota handling
- [ ] Invalid credential handling

## Final Verification

- [ ] `go build ./...` succeeds
- [ ] `go test ./...` passes (with custom test timeouts if needed)
- [ ] `go test -race ./...` shows no race conditions
- [ ] `golint ./...` shows no linting issues (if applicable)
- [ ] `go fmt ./...` formatting is correct
- [ ] No TODO/FIXME comments left
- [ ] All functions have godoc comments
- [ ] Code follows repository style guide
- [ ] Reviewed by another team member
- [ ] Works with real provider API (manual test)

## Go-Live Checklist

- [ ] Provider documented in public documentation
- [ ] Example configuration provided
- [ ] Monitoring/alerting configured if needed
- [ ] Support documentation written
- [ ] Team trained on new provider
- [ ] Changelog/release notes updated
- [ ] Backward compatibility verified (if replacing existing provider)
- [ ] Rollout plan in place (if needed)

---

**Date Created:** ___________  
**Provider Name:** ___________  
**Implemented By:** ___________  
**Reviewed By:** ___________  
**Approved For Production:** ___________
