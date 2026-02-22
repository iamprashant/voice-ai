# STT/TTS Transformer Integration Guide — UI to Backend

## Overview

This guide covers how to add a new Speech-to-Text (STT) or Text-to-Speech (TTS) provider to the Rapida platform. Transformers convert between audio and text — STT transcribes user speech, TTS synthesizes agent responses.

## Architecture

```
UI (React) → Deployment Config → assistant-api → Transformer Factory → Provider SDK → Audio/Text Packets
```

### Supported Providers

| Provider | STT | TTS | Transport |
|---|---|---|---|
| **Deepgram** | ✅ | ✅ | WebSocket streaming |
| **Google Cloud Speech** | ✅ | ✅ | gRPC streaming |
| **Azure Speech Services** | ✅ | ✅ | SDK streaming |
| **Cartesia** | ✅ | ✅ | WebSocket streaming |
| **AssemblyAI** | ✅ | — | WebSocket streaming |
| **Rev.ai** | ✅ | ✅ | WebSocket streaming |
| **Sarvam AI** | ✅ | ✅ | REST/WebSocket |
| **ElevenLabs** | — | ✅ | WebSocket streaming |
| **OpenAI** | ✅ | ✅ | REST |
| **AWS** | ✅ | ✅ | SDK streaming |
| **Resemble** | — | ✅ | REST |
| **Speechmatics** | ✅ | — | WebSocket streaming |

## Directory Structure

### Backend

```
api/assistant-api/internal/
├── transformer/                          # Transformer factory + implementations
│   ├── transformer.go                    # Factory: GetSpeechToTextTransformer(), GetTextToSpeechTransformer()
│   ├── transformer_test.go               # Factory tests
│   ├── README.md                         # Detailed implementation guide
│   ├── deepgram/                         # Example provider implementation
│   │   ├── deepgram.go                   # Options struct (credential parsing, config from listen.*/speak.* keys)
│   │   ├── stt.go                        # SpeechToTextTransformer implementation
│   │   ├── tts.go                        # TextToSpeechTransformer implementation
│   │   ├── normalizer.go                 # Provider-specific text normalization
│   │   ├── deepgram_test.go              # Tests
│   │   └── internal/                     # Provider SDK callback wrappers
│   ├── assembly-ai/                      # STT only
│   ├── aws/                              # Polly TTS + Transcribe STT
│   ├── azure/                            # Cognitive Services STT/TTS
│   ├── cartesia/                         # Low-latency STT/TTS
│   ├── elevenlabs/                       # TTS only
│   ├── google/                           # Cloud Speech STT/TTS
│   ├── openai/                           # Whisper STT + TTS
│   ├── resemble/                         # TTS only
│   ├── revai/                            # STT/TTS
│   ├── sarvam/                           # Indian language STT/TTS
│   └── speechmatics/                     # STT only
├── type/
│   ├── transformer.go                    # Base Transformers[IN] interface
│   ├── stt_transformer.go               # SpeechToTextTransformer interface
│   └── tts_transformer.go               # TextToSpeechTransformer interface
└── normalizers/                          # Shared text normalizers (currency, date, URL, numbers, etc.)
```

### UI

```
ui/src/
├── providers/
│   ├── provider.development.json         # Provider registry (featureList: ["stt", "tts"])
│   ├── provider.production.json          # Production provider registry
│   ├── index.ts                          # Provider accessors (SPEECH_TO_TEXT_PROVIDER, TEXT_TO_SPEECH_PROVIDER)
│   ├── deepgram/                         # Example provider metadata
│   │   ├── voices.json                   # TTS voices
│   │   ├── speech-to-text-models.json    # STT models
│   │   └── speech-to-text-languages.json # Supported languages
│   ├── cartesia/                         # voices, models, languages
│   ├── elevenlabs/                       # models, voices, languages
│   ├── google/                           # voice, model, language JSONs
│   ├── azure-speech-service/             # voices, languages
│   ├── sarvam/                           # models, voices, languages
│   └── assemblyai/                       # models, languages
├── app/components/providers/
│   ├── speech-to-text/
│   │   ├── index.tsx                     # Exports SpeechToTextConfigComponent
│   │   ├── provider.tsx                  # Switch statement routing to provider components
│   │   ├── deepgram/                     # Deepgram STT config (model, language, keywords)
│   │   │   ├── index.tsx                 # ConfigureDeepgramSpeechToText component
│   │   │   └── constant.ts              # GetDeepgramDefaultOptions, ValidateDeepgramOptions
│   │   ├── assemblyai/
│   │   ├── azure-speech-service/
│   │   ├── cartesia/
│   │   ├── google-speech-service/
│   │   ├── openai/
│   │   └── sarvam/
│   └── text-to-speech/
│       ├── index.tsx                     # Exports TextToSpeechConfigComponent
│       ├── provider.tsx                  # Switch statement routing to provider components
│       ├── deepgram/                     # Deepgram TTS config (voice)
│       ├── azure-speech-service/
│       ├── cartesia/
│       ├── elevenlabs/
│       ├── google-speech-service/
│       ├── openai/
│       ├── playht/
│       └── sarvam/
```

## Core Interfaces

### Base Transformer Interface

```go
// api/assistant-api/internal/type/transformer.go
type Transformers[IN any] interface {
    Initialize() error
    Transform(context.Context, IN) error
    Close(context.Context) error
}
```

### SpeechToTextTransformer

```go
// api/assistant-api/internal/type/stt_transformer.go
type SpeechToTextTransformer interface {
    Name() string
    Transformers[UserAudioPacket]
}
```

- **Input**: `UserAudioPacket` (raw PCM audio bytes from user microphone/phone)
- **Output**: Calls `onPacket(TranscriptPacket)` callback with transcribed text
- **Lifecycle**: `Initialize()` opens connection → `Transform()` streams audio → `Close()` disconnects

### TextToSpeechTransformer

```go
// api/assistant-api/internal/type/tts_transformer.go
type TextToSpeechTransformer interface {
    Name() string
    Transformers[LLMPacket]
}
```

- **Input**: `LLMPacket` (text chunks from LLM response)
- **Output**: Calls `onPacket(AudioPacket)` callback with synthesized PCM audio
- **Lifecycle**: `Initialize()` opens connection → `Transform()` sends text → `Close()` disconnects

## Factory Pattern

The factory in `transformer.go` uses a switch on `AudioTransformer` constants:

```go
// Constants for provider codes — MUST match UI provider codes
const (
    DEEPGRAM              AudioTransformer = "deepgram"
    GOOGLE_SPEECH_SERVICE AudioTransformer = "google-speech-service"
    AZURE_SPEECH_SERVICE  AudioTransformer = "azure-speech-service"
    CARTESIA              AudioTransformer = "cartesia"
    REVAI                 AudioTransformer = "revai"
    SARVAM                AudioTransformer = "sarvamai"
    ELEVENLABS            AudioTransformer = "elevenlabs"
    ASSEMBLYAI            AudioTransformer = "assemblyai"
)

// Factory functions — instantiate transformers by provider code
func GetTextToSpeechTransformer(ctx, logger, provider, credential, onPacket, opts) (TextToSpeechTransformer, error)
func GetSpeechToTextTransformer(ctx, logger, provider, credential, onPacket, opts) (SpeechToTextTransformer, error)
```

## Parameter Key Conventions

Options are passed via `utils.Option` (a `map[string]interface{}`), with provider-specific keys:

| Prefix | Domain | Examples |
|---|---|---|
| `listen.*` | STT parameters | `listen.language`, `listen.model`, `listen.smart_format`, `listen.filler_words`, `listen.vad_events`, `listen.endpointing`, `listen.keyword` |
| `speak.*` | TTS parameters | `speak.voice.id`, `speak.model`, `speak.language`, `speak.speed`, `speak.emotion` |
| `microphone.*` | Audio pipeline | `microphone.eos.timeout`, `microphone.eos.provider`, `microphone.denoising.provider`, `microphone.vad.provider`, `microphone.vad.threshold` |

These keys are set by the UI deployment configuration and stored in the assistant deployment entity.

## Adding a New STT/TTS Provider — Step by Step

### Step 1: Backend — Create Provider Package

Create `api/assistant-api/internal/transformer/<provider>/` with:

#### `<provider>.go` — Options/Config

```go
package internal_transformer_newprovider

import (
    "fmt"
    "github.com/rapidaai/pkg/commons"
    "github.com/rapidaai/pkg/utils"
    "github.com/rapidaai/protos"
)

type newproviderOption struct {
    apiKey  string
    logger  commons.Logger
    mdlOpts utils.Option
}

func NewProviderOption(
    logger commons.Logger,
    vaultCredential *protos.VaultCredential,
    opts utils.Option,
) (*newproviderOption, error) {
    cx, ok := vaultCredential.GetValue().AsMap()["key"]
    if !ok {
        return nil, fmt.Errorf("illegal vault config: missing API key")
    }
    return &newproviderOption{
        apiKey:  cx.(string),
        logger:  logger,
        mdlOpts: opts,
    }, nil
}

func (o *newproviderOption) GetEncoding() string {
    return "linear16"
}

// Parse listen.* options into provider-specific STT config
func (o *newproviderOption) SpeechToTextOptions() *ProviderSTTConfig {
    cfg := &ProviderSTTConfig{
        Language: "en-US",
        Model:    "default",
    }
    if lang, err := o.mdlOpts.GetString("listen.language"); err == nil {
        cfg.Language = lang
    }
    if model, err := o.mdlOpts.GetString("listen.model"); err == nil {
        cfg.Model = model
    }
    return cfg
}
```

#### `stt.go` — Speech-to-Text Implementation

```go
package internal_transformer_newprovider

import (
    "context"
    "sync"
    internal_type "github.com/rapidaai/api/assistant-api/internal/type"
    "github.com/rapidaai/pkg/commons"
    "github.com/rapidaai/pkg/utils"
    "github.com/rapidaai/protos"
)

type newproviderSTT struct {
    *newproviderOption
    mu        sync.Mutex
    ctx       context.Context
    ctxCancel context.CancelFunc
    logger    commons.Logger
    onPacket  func(pkt ...internal_type.Packet) error
    // provider-specific client
}

func (*newproviderSTT) Name() string {
    return "newprovider-speech-to-text"
}

func NewProviderSpeechToText(
    ctx context.Context,
    logger commons.Logger,
    vaultCredential *protos.VaultCredential,
    onPacket func(pkt ...internal_type.Packet) error,
    opts utils.Option,
) (internal_type.SpeechToTextTransformer, error) {
    provOpts, err := NewProviderOption(logger, vaultCredential, opts)
    if err != nil {
        return nil, err
    }
    ct, ctxCancel := context.WithCancel(ctx)
    return &newproviderSTT{
        ctx:               ct,
        ctxCancel:         ctxCancel,
        logger:            logger,
        newproviderOption: provOpts,
        onPacket:          onPacket,
    }, nil
}

func (s *newproviderSTT) Initialize() error {
    // Open WebSocket/gRPC connection to provider
    // Set up callback to call s.onPacket() with TranscriptPacket on transcription results
    return nil
}

func (s *newproviderSTT) Transform(ctx context.Context, audio internal_type.UserAudioPacket) error {
    // Send audio bytes to provider for transcription
    return nil
}

func (s *newproviderSTT) Close(ctx context.Context) error {
    s.ctxCancel()
    // Close provider connection
    return nil
}
```

#### `tts.go` — Text-to-Speech Implementation

```go
package internal_transformer_newprovider

import (
    "context"
    "sync"
    internal_type "github.com/rapidaai/api/assistant-api/internal/type"
    "github.com/rapidaai/pkg/commons"
    "github.com/rapidaai/pkg/utils"
    "github.com/rapidaai/protos"
)

type newproviderTTS struct {
    *newproviderOption
    mu        sync.Mutex
    ctx       context.Context
    ctxCancel context.CancelFunc
    logger    commons.Logger
    onPacket  func(pkt ...internal_type.Packet) error
}

func (*newproviderTTS) Name() string {
    return "newprovider-text-to-speech"
}

func NewProviderTextToSpeech(
    ctx context.Context,
    logger commons.Logger,
    vaultCredential *protos.VaultCredential,
    onPacket func(pkt ...internal_type.Packet) error,
    opts utils.Option,
) (internal_type.TextToSpeechTransformer, error) {
    provOpts, err := NewProviderOption(logger, vaultCredential, opts)
    if err != nil {
        return nil, err
    }
    ct, ctxCancel := context.WithCancel(ctx)
    return &newproviderTTS{
        ctx:               ct,
        ctxCancel:         ctxCancel,
        logger:            logger,
        newproviderOption: provOpts,
        onPacket:          onPacket,
    }, nil
}

func (t *newproviderTTS) Initialize() error {
    // Open connection to TTS provider
    return nil
}

func (t *newproviderTTS) Transform(ctx context.Context, packet internal_type.LLMPacket) error {
    // Send text to provider, receive audio, call t.onPacket() with AudioPacket
    return nil
}

func (t *newproviderTTS) Close(ctx context.Context) error {
    t.ctxCancel()
    return nil
}
```

#### `normalizer.go` — Optional Text Normalization

Provider-specific text normalization (e.g., handling currency symbols, numbers, URLs for TTS). Can use shared normalizers from `api/assistant-api/internal/normalizers/`.

### Step 2: Backend — Register in Factory

In `api/assistant-api/internal/transformer/transformer.go`:

1. **Add import**:
```go
internal_transformer_newprovider "github.com/rapidaai/api/assistant-api/internal/transformer/newprovider"
```

2. **Add constant**:
```go
const (
    // ... existing constants
    NEWPROVIDER AudioTransformer = "newprovider"
)
```

3. **Add TTS case** in `GetTextToSpeechTransformer()`:
```go
case NEWPROVIDER:
    return internal_transformer_newprovider.NewProviderTextToSpeech(ctx, logger, credential, onPacket, opts)
```

4. **Add STT case** in `GetSpeechToTextTransformer()`:
```go
case NEWPROVIDER:
    return internal_transformer_newprovider.NewProviderSpeechToText(ctx, logger, credential, onPacket, opts)
```

### Step 3: UI — Add Provider Metadata JSON

In `ui/src/providers/provider.development.json` (and `provider.production.json`), add:

```json
{
    "code": "newprovider",
    "name": "New Provider",
    "description": "Description of the provider",
    "image": "newprovider.svg",
    "featureList": ["stt", "tts"],
    "configurations": [
        { "name": "key", "type": "string", "label": "API Key" }
    ]
}
```

- `"stt"` → appears in `SPEECH_TO_TEXT_PROVIDER` dropdown
- `"tts"` → appears in `TEXT_TO_SPEECH_PROVIDER` dropdown
- `configurations` → defines credential fields stored in Vault

### Step 4: UI — Create Provider Data Files

Create `ui/src/providers/newprovider/`:

**`voices.json`** (for TTS):
```json
[
    { "id": "voice-1", "name": "Voice One", "language": "en-US", "gender": "female" },
    { "id": "voice-2", "name": "Voice Two", "language": "en-US", "gender": "male" }
]
```

**`speech-to-text-models.json`** (for STT):
```json
[
    { "id": "default", "name": "Default Model", "description": "General purpose STT" },
    { "id": "enhanced", "name": "Enhanced Model", "description": "Higher accuracy" }
]
```

**`languages.json`**:
```json
[
    { "id": "en-US", "name": "English (US)" },
    { "id": "es-ES", "name": "Spanish (Spain)" }
]
```

### Step 5: UI — Register Provider Accessors

In `ui/src/providers/index.ts`, add accessor functions:

```typescript
export const NEWPROVIDER_VOICE = () => {
    return require('./newprovider/voices.json');
};

export const NEWPROVIDER_SPEECH_TO_TEXT_MODEL = () => {
    return require('./newprovider/speech-to-text-models.json');
};

export const NEWPROVIDER_LANGUAGE = () => {
    return require('./newprovider/languages.json');
};
```

### Step 6: UI — Create STT Config Component

Create `ui/src/app/components/providers/speech-to-text/newprovider/`:

**`constant.ts`**:
```typescript
import { Metadata } from '@rapidaai/react';

export const GetNewproviderDefaultOptions = (parameters: Metadata[]): Metadata[] => {
    const defaults = [
        { key: 'listen.language', value: 'en-US' },
        { key: 'listen.model', value: 'default' },
    ];
    const existingKeys = new Set(parameters.map(p => p.getKey()));
    const newParams = defaults
        .filter(d => !existingKeys.has(d.key))
        .map(d => { const m = new Metadata(); m.setKey(d.key); m.setValue(d.value); return m; });
    return [...parameters, ...newParams];
};

export const ValidateNewproviderOptions = (parameters: Metadata[]): string | undefined => {
    const model = parameters.find(p => p.getKey() === 'listen.model');
    if (!model?.getValue()) return 'Model is required';
    return undefined;
};
```

**`index.tsx`**:
```tsx
import { FC } from 'react';
import { Dropdown } from '@/app/components/dropdown';
import { NEWPROVIDER_SPEECH_TO_TEXT_MODEL, NEWPROVIDER_LANGUAGE } from '@/providers';

interface Props {
    parameters: Metadata[];
    onParameterChange: (key: string, value: string) => void;
}

export const ConfigureNewproviderSpeechToText: FC<Props> = ({ parameters, onParameterChange }) => {
    return (
        <>
            <Dropdown
                label="Model"
                options={NEWPROVIDER_SPEECH_TO_TEXT_MODEL()}
                value={parameters.find(p => p.getKey() === 'listen.model')?.getValue()}
                onChange={(val) => onParameterChange('listen.model', val)}
            />
            <Dropdown
                label="Language"
                options={NEWPROVIDER_LANGUAGE()}
                value={parameters.find(p => p.getKey() === 'listen.language')?.getValue()}
                onChange={(val) => onParameterChange('listen.language', val)}
            />
        </>
    );
};
```

### Step 7: UI — Register in STT Provider Switch

In `ui/src/app/components/providers/speech-to-text/provider.tsx`:

1. **Import** the new component:
```typescript
import {
    ConfigureNewproviderSpeechToText,
    GetNewproviderDefaultOptions,
    ValidateNewproviderOptions,
} from '@/app/components/providers/speech-to-text/newprovider';
```

2. **Add case** to `GetDefaultSpeechToTextIfInvalid`:
```typescript
case 'newprovider':
    return GetNewproviderDefaultOptions(parameters);
```

3. **Add case** to `ValidateSpeechToTextIfInvalid`:
```typescript
case 'newprovider':
    return ValidateNewproviderOptions(parameters);
```

4. **Add case** to `SpeechToTextConfigComponent`:
```typescript
case 'newprovider':
    return (
        <ConfigureNewproviderSpeechToText
            parameters={parameters}
            onParameterChange={onChangeParameter}
        />
    );
```

### Step 8: UI — Create TTS Config Component (if applicable)

Same pattern as STT but in `ui/src/app/components/providers/text-to-speech/newprovider/`. Register in `text-to-speech/provider.tsx` switch statements. TTS typically uses `speak.*` parameter keys.

## Data Flow Diagram

```
┌───────────────────────────────────────────────────────────────────┐
│ UI: Deployment Config                                             │
│  Select STT provider → Configure listen.* params → Link vault    │
│  Select TTS provider → Configure speak.* params → Link vault     │
└──────────────────────┬────────────────────────────────────────────┘
                       │ gRPC: CreateDeployment
                       ▼
┌───────────────────────────────────────────────────────────────────┐
│ Deployment Entity → PostgreSQL                                    │
│  InputAudio: { provider: "deepgram", options: [...listen.*...] }  │
│  OutputAudio: { provider: "elevenlabs", options: [...speak.*...] }│
└───────────────────────────────────────────────────────────────────┘

═══════════ AT CALL TIME ═══════════

Talker.Talk()
    │
    ├── GetSpeechToTextTransformer("deepgram", credential, opts)
    │       → deepgramSTT.Initialize() (opens WebSocket)
    │
    └── GetTextToSpeechTransformer("elevenlabs", credential, opts)
            → elevenlabsTTS.Initialize() (opens WebSocket)

User speaks → PCM audio
    │
    ├── VAD filters silence
    ├── Denoiser cleans audio
    └── STT.Transform(audioPacket)
            │
            └── onPacket(TranscriptPacket) → LLM processing
                                                │
                                                └── TTS.Transform(llmPacket)
                                                        │
                                                        └── onPacket(AudioPacket) → Streamer → User hears
```

## Checklist for New STT/TTS Provider

### Backend
- [ ] Create `api/assistant-api/internal/transformer/<provider>/` directory
- [ ] Implement `<provider>.go` — options struct, credential parsing, config from `listen.*`/`speak.*` keys
- [ ] Implement `stt.go` — `SpeechToTextTransformer` (if STT supported)
- [ ] Implement `tts.go` — `TextToSpeechTransformer` (if TTS supported)
- [ ] Implement `normalizer.go` — optional text normalization for TTS
- [ ] Add import + constant + switch case in `transformer.go` factory
- [ ] Add tests in `<provider>_test.go`

### UI — Provider Metadata
- [ ] Add provider entry to `provider.development.json` with `"stt"` and/or `"tts"` in `featureList`
- [ ] Add provider entry to `provider.production.json`
- [ ] Create `ui/src/providers/<provider>/` with JSON files (voices, models, languages)
- [ ] Add accessor functions in `ui/src/providers/index.ts`

### UI — Config Components
- [ ] Create STT config component at `ui/src/app/components/providers/speech-to-text/<provider>/`
- [ ] Create TTS config component at `ui/src/app/components/providers/text-to-speech/<provider>/`
- [ ] Register in `speech-to-text/provider.tsx` — 3 switch statements (defaults, validation, component)
- [ ] Register in `text-to-speech/provider.tsx` — 3 switch statements (defaults, validation, component)
- [ ] Ensure provider `code` in JSON matches `AudioTransformer` constant in backend
