# Telephony Integration Guide — UI to Backend

## Overview

This guide covers how to add a new telephony provider to the Rapida platform. Telephony enables voice agents to make and receive phone calls through providers like Twilio, Vonage, Exotel, Asterisk, and native SIP.

## Architecture

```
UI (React) → gRPC → assistant-api → SIP Server / Telephony Providers → RTP Audio → Voice Pipeline
```

### Supported Providers

| Provider | Transport | Direction | Description |
|---|---|---|---|
| **Twilio** | WebSocket | Inbound + Outbound | Cloud telephony via webhook + WebSocket media |
| **Vonage** | WebSocket | Inbound + Outbound | Cloud telephony via webhook + WebSocket media |
| **Exotel** | WebSocket | Inbound + Outbound | Cloud telephony via webhook + WebSocket media |
| **Asterisk** | AudioSocket (TCP) | Inbound + Outbound | PBX via AudioSocket protocol |
| **SIP** | Native SIP/RTP | Inbound + Outbound | Direct SIP trunk integration |

## Directory Structure

```
api/assistant-api/
├── api/talk/                              # REST/gRPC handlers for calls
│   ├── talk.go                            # ConversationApi + ConversationGrpcApi
│   ├── inbound_call.go                    # CallReciever (webhook), CallTalkerByContext (WS upgrade)
│   ├── outbound_call.go                   # CreatePhoneCall gRPC handler
│   ├── whatsapp.go                        # WhatsApp receiver
│   ├── get.go                             # GetAllConversation, GetAllMessage queries
│   └── metric.go                          # CreateMessageMetric, CreateConversationMetric
├── api/assistant-deployment/              # Deployment configuration handlers
│   ├── assistant_deployment.go            # Base deployment handler struct
│   ├── create_assistant_phone_deployment.go  # Create phone deployment
│   └── get_assistant_phone_deployment.go     # Get phone deployment
├── internal/channel/telephony/            # Telephony transport layer
│   ├── telephony.go                       # Telephony struct — streamer factory
│   ├── inbound.go                         # InboundDispatcher — call receive + Redis context
│   ├── outbound.go                        # OutboundDispatcher — REST API call placement
│   └── internal/                          # Provider-specific streamer implementations
├── internal/callcontext/                  # Redis-backed call context store
│   ├── store.go                           # Atomic get-and-delete via Lua, 5-min TTL
│   └── types.go                           # CallContext struct
├── sip/                                   # Native SIP server
│   ├── sip.go                             # SIPEngine — multi-tenant orchestrator
│   └── infra/                             # SIP infrastructure
│       ├── server.go                      # SIP signaling server (~1799 lines)
│       ├── session.go                     # Per-call session state machine
│       ├── auth.go                        # SIP URI credential parsing
│       ├── rtp.go                         # RTP audio handler (UDP)
│       ├── sdp.go                         # SDP generation/parsing, codec negotiation
│       ├── rtp_port_allocator.go          # Redis-backed distributed port allocation
│       └── types.go                       # Config, Transport, CallState, SessionInfo
└── socket/socket.go                       # AudioSocket server (Asterisk)

ui/src/
├── app/pages/assistant/actions/create-deployment/phone/  # Phone deployment config page
├── app/components/providers/telephony/    # Provider-specific config components
│   ├── index.tsx                          # TelephonyProvider — main provider selector
│   ├── twilio/index.tsx                   # Twilio config (credential + phone)
│   ├── sip/index.tsx                      # SIP config (credential + caller ID)
│   ├── asterisk/index.tsx                 # Asterisk config
│   ├── vonage/index.tsx                   # Vonage config
│   └── exotel/index.tsx                   # Exotel config
└── providers/                             # Provider metadata
    └── provider.development.json          # Provider registry (featureList: ["telephony"])
```

## End-to-End Flow

### 1. UI — Phone Deployment Configuration

The phone deployment page collects:
- **Telephony provider** — dropdown filtered by `featureList.includes('telephony')`
- **Vault credential** — `rapida.credential_id` links to stored provider secrets
- **Provider-specific params** — phone number, caller ID, etc.
- **Experience config** — greeting, error message, idle timeout (30s default), max session (300s default)
- **Audio input** (STT provider) and **output** (TTS provider)

On submit → gRPC `CreateAssistantPhoneDeployment` with `AssistantPhoneDeployment` proto.

### 2. Deployment Entity

```
AssistantPhoneDeployment
├── AssistantDeploymentBehavior
│   ├── Greeting (string)
│   ├── Mistake (string — error message)
│   ├── IdealTimeout (int — seconds)
│   ├── IdealTimeoutBackoff (float)
│   └── MaxSessionDuration (int — seconds)
├── AssistantDeploymentTelephony
│   ├── TelephonyProvider (string — "twilio", "sip", etc.)
│   └── TelephonyOption[] (key-value pairs — credential_id, phone, etc.)
├── InputAudio (provider code + options for STT)
└── OutputAudio (provider code + options for TTS)
```

### 3. Inbound Call Flow

#### Path A — WebSocket-based (Twilio, Vonage, Exotel)

```
1. Provider sends HTTP webhook → CallReciever handler
2. InboundDispatcher.HandleReceiveCall():
   - Creates conversation record in DB
   - Saves CallContext to Redis (5-min TTL)
   - Returns contextID in webhook response
3. Provider opens WebSocket → CallTalkerByContext handler:
   - Upgrades HTTP to WebSocket
   - InboundDispatcher.ResolveCallSessionByContext() — atomic Redis GET+DELETE
   - Telephony.NewStreamer() creates provider-specific streamer
   - Creates Talker → talker.Talk() (blocking for call duration)
```

#### Path B — Native SIP

```
1. External SIP device sends INVITE → SIP Server handleInvite()
2. Middleware chain: CredentialMiddleware → authMiddleware → assistantMiddleware → vaultConfigResolver
3. Parses SDP, negotiates codec (PCMU/PCMA/G722), allocates RTP port from Redis pool
4. Sends 100 Trying → 180 Ringing → 200 OK with SDP answer
5. Callback to SIPEngine.handleInvite():
   - Resolves assistant from SIP URI: sip:{assistantID}:{apiKey}@host
   - Creates conversation, builds CallContext
   - Launches startCall() goroutine → creates SIP streamer + Talker
```

#### Path C — AudioSocket (Asterisk)

```
1. Asterisk opens TCP → AudioSocket server
2. Reads UUID frame (16-byte contextID)
3. InboundDispatcher.ResolveCallSessionByContext()
4. Creates Asterisk streamer + Talker
```

### 4. Outbound Call Flow

```
1. SDK/API calls CreatePhoneCall gRPC:
   - Loads assistant, validates phone deployment exists
   - Creates conversation, saves CallContext to Redis (status="queued")
   - Fires go outboundDispatcher.Dispatch() asynchronously
   - Returns immediately to caller

2. OutboundDispatcher.Dispatch():
   - Resolves call context from Redis
   - Loads assistant + vault credential
   - For Twilio/Vonage/Exotel: places call via provider REST API (webhook-based)
   - For SIP: calls SIPEngine.PlaceOutboundCall():
     - Allocates RTP port, creates RTP handler
     - Generates SDP offer, sends SIP INVITE
     - Handles digest auth (401/407 challenges)
     - On 200 OK: parse answer SDP → start RTP → send ACK
     - Notifies onInvite → starts conversation
```

## SIP Infrastructure Details

### Middleware Chain
```
CredentialMiddleware → authMiddleware → assistantMiddleware → vaultConfigResolver
```
Each middleware enriches the `SIPRequestContext` with parsed credentials, authenticated principal, loaded assistant, and vault-resolved SIP config.

### Session State Machine
```
Initializing → Ringing → Connected → OnHold → Ending → Ended/Failed
```

### SIP Methods Handled
INVITE, ACK, BYE, CANCEL, REGISTER, OPTIONS, UPDATE, INFO, NOTIFY, REFER (declined), SUBSCRIBE (489), MESSAGE, unknown.

### RTP Details
- Dual UDP sockets: receive (unconnected `ReadFromUDP`) + send (connected via `DialUDP`)
- 20ms packet interval, G.711 codecs (PCMU silence=0xFF, PCMA silence=0xD5)
- `sendInitialSilence` for NAT/firewall RTP path punching
- Auto-detect remote address from first received packet
- Codec hot-swap on re-INVITE/UPDATE

### Credential Resolution (Vault)
```go
// Vault fields for SIP
sip_uri         // "sip:192.168.1.5:5060" — parsed for server+port
sip_server      // Explicit server (overrides sip_uri)
sip_username    // Auth username
sip_password    // Auth password
sip_realm       // SIP realm
sip_domain      // SIP domain
```

## Adding a New Telephony Provider

### Step 1: UI — Add Provider Metadata

In `ui/src/providers/provider.development.json` (and `provider.production.json`), add:

```json
{
  "code": "newprovider",
  "name": "New Provider",
  "description": "New telephony provider for voice calls",
  "image": "newprovider.svg",
  "featureList": ["telephony"],
  "configurations": [
    { "name": "key", "type": "string", "label": "API Key" },
    { "name": "secret", "type": "string", "label": "API Secret" }
  ]
}
```

The `configurations` array defines credential fields stored in Vault.

### Step 2: UI — Create Telephony Config Component

Create `ui/src/app/components/providers/telephony/newprovider/index.tsx`:

```tsx
export default function ConfigureNewproviderTelephony({ parameters, onChangeParameter }: ProviderComponentProps) {
    return (
        <>
            <CredentialDropdown
                value={parameters?.['rapida.credential_id']}
                onChange={(val) => onChangeParameter('rapida.credential_id', val)}
                provider="newprovider"
            />
            <Input
                label="Phone Number"
                value={parameters?.['phone']}
                onChange={(val) => onChangeParameter('phone', val)}
            />
        </>
    );
}
```

### Step 3: UI — Register in Telephony Provider Switch

In `ui/src/app/components/providers/telephony/index.tsx`, add the new provider case to the switch statement that renders provider-specific config.

### Step 4: Backend — Add Inbound Webhook Handler

In `api/assistant-api/api/talk/inbound_call.go`, add routes for the new provider's webhook format. The provider will call back with call metadata, and you need to:
1. Parse the provider's webhook payload
2. Call `InboundDispatcher.HandleReceiveCall()` with extracted call info
3. Return a provider-specific response (e.g., TwiML for Twilio, NCCO for Vonage)

### Step 5: Backend — Create Provider Streamer

In `api/assistant-api/internal/channel/telephony/internal/`, create a new streamer that:
1. Implements the `Streamer` interface (`Send()`, `Recv()`, `Close()`)
2. Handles the provider's media transport (WebSocket, TCP, etc.)
3. Converts between the provider's audio format and Rapida's internal PCM format

### Step 6: Backend — Register in Telephony Factory

In `api/assistant-api/internal/channel/telephony/telephony.go`, add a case in `NewStreamer()` for the new provider code.

### Step 7: Backend — Add Outbound Call Support

In `api/assistant-api/internal/channel/telephony/outbound.go`, add logic in `Dispatch()` to:
1. Call the provider's REST API to initiate an outbound call
2. The provider then calls back via webhook → same inbound flow

## Data Flow Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│ UI: Configure Phone Deployment                                   │
│   Select provider │ Link vault credential │ Set phone number     │
│   Configure STT/TTS │ Set greeting + timeouts                    │
└──────────────────────┬──────────────────────────────────────────┘
                       │ gRPC: CreateAssistantPhoneDeployment
                       ▼
┌──────────────────────────────────────────────────────────────────┐
│ Deployment Service → PostgreSQL                                  │
│   Stores: behavior + telephony provider + audio configs          │
└──────────────────────────────────────────────────────────────────┘

═══════════ INBOUND CALL ═══════════

Provider webhook ──HTTP──► CallReciever ──► InboundDispatcher
                                             │ Create conversation
                                             │ Save CallContext → Redis
                                             ▼
Provider media ───WS/AudioSocket/SIP──► Streamer ──► Talker.Talk()
                                                       │
                                              ┌────────┴────────┐
                                              │ STT → LLM → TTS │
                                              └─────────────────┘

═══════════ OUTBOUND CALL ═══════════

SDK gRPC ──► CreatePhoneCall ──► Save CallContext → Redis
                                  │ go Dispatch()
                       ┌──────────┴──────────┐
                       │ Twilio/Vonage/Exotel │  SIP/Asterisk
                       │   REST API call      │  SIP INVITE
                       │   → webhook back     │  → RTP setup
                       └──────────┬──────────┘  → startCall
                                  ▼
                          Streamer + Talker.Talk()

═══════════ NATIVE SIP ═══════════

External SIP ──INVITE──► SIP Server (port 4573)
  │                         │ Middleware chain
  │                         │ SDP negotiate + RTP allocate
  │                         │ 100→180→200 OK
  │◄────────────────────────┘
  │──RTP audio──► RTPHandler ──► SIP Streamer ──► Talker.Talk()
  │◄──RTP audio──┘                                  │
                                           ┌────────┴────────┐
                                           │ STT → LLM → TTS │
                                           └─────────────────┘
```

## Checklist for New Telephony Provider

- [ ] Add provider entry to `ui/src/providers/provider.development.json` with `"telephony"` in `featureList`
- [ ] Add provider entry to `ui/src/providers/provider.production.json`
- [ ] Create UI config component at `ui/src/app/components/providers/telephony/<provider>/index.tsx`
- [ ] Register in telephony provider switch in `ui/src/app/components/providers/telephony/index.tsx`
- [ ] Add inbound webhook handler in `api/assistant-api/api/talk/inbound_call.go`
- [ ] Add inbound webhook route in `api/assistant-api/router/assistant.go`
- [ ] Create provider streamer in `api/assistant-api/internal/channel/telephony/internal/`
- [ ] Register streamer in `api/assistant-api/internal/channel/telephony/telephony.go`
- [ ] Add outbound dispatch logic in `api/assistant-api/internal/channel/telephony/outbound.go`
- [ ] Add vault credential configuration fields (provider `configurations` in JSON)
- [ ] Test inbound + outbound call flows end-to-end
