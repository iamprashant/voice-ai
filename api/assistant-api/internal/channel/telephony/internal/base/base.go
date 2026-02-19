// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package internal_telephony_base

import (
	"encoding/base64"

	internal_audio "github.com/rapidaai/api/assistant-api/internal/audio"
	internal_audio_resampler "github.com/rapidaai/api/assistant-api/internal/audio/resampler"
	channel_base "github.com/rapidaai/api/assistant-api/internal/channel/base"
	internal_assistant_entity "github.com/rapidaai/api/assistant-api/internal/entity/assistants"
	internal_conversation_entity "github.com/rapidaai/api/assistant-api/internal/entity/conversations"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/utils"
	"github.com/rapidaai/protos"
)

// TelephonyOption configures a BaseTelephonyStreamer.
type TelephonyOption func(*telephonyConfig)

type telephonyConfig struct {
	// sourceAudioConfig is the native audio format received from the telephony
	// provider. Defaults to RAPIDA_AUDIO_CONFIG (linear16 16kHz) if nil.
	sourceAudioConfig *protos.AudioConfig

	// baseOpts are forwarded to channel_base.NewBaseStreamer.
	baseOpts []channel_base.Option
}

// WithSourceAudioConfig sets the native audio format of the telephony provider.
// This drives automatic derivation of input/output buffer thresholds in BaseStreamer.
// Pass nil to default to RAPIDA_AUDIO_CONFIG (linear16 16kHz — no resampling).
func WithSourceAudioConfig(cfg *protos.AudioConfig) TelephonyOption {
	return func(c *telephonyConfig) { c.sourceAudioConfig = cfg }
}

// WithBaseOption appends one or more channel_base.Option to the underlying
// BaseStreamer configuration. Use this for advanced overrides (channel sizes,
// explicit thresholds, etc.).
func WithBaseOption(opts ...channel_base.Option) TelephonyOption {
	return func(c *telephonyConfig) { c.baseOpts = append(c.baseOpts, opts...) }
}

// internal rapida audio config
var RAPIDA_AUDIO_CONFIG = internal_audio.NewLinear16khzMonoAudioConfig()

// ============================================================================
// BaseTelephonyStreamer — telephony-specific base that embeds BaseStreamer
// ============================================================================

// BaseTelephonyStreamer embeds channel_base.BaseStreamer for common buffer,
// channel, and lifecycle management. It adds telephony-specific concerns:
// entity references, audio resampler, base64 encoder, and vault credentials.
//
// Concrete telephony streamers (Twilio, Exotel, Vonage, SIP, Asterisk) embed
// this struct and only implement transport-specific I/O logic.
type BaseTelephonyStreamer struct {
	channel_base.BaseStreamer

	// Telephony-specific fields
	assistant             *internal_assistant_entity.Assistant
	assistantConversation *internal_conversation_entity.AssistantConversation

	resampler       internal_type.AudioResampler
	encoder         *base64.Encoding
	vaultCredential *protos.VaultCredential

	// sourceAudioConfig is the native audio format received from the telephony
	// provider (e.g. µ-law 8kHz for Twilio, linear16 16kHz for Vonage).
	// CreateVoiceRequest uses this to resample input audio to the internal
	// Rapida format (linear16 16kHz) before sending downstream.
	sourceAudioConfig *protos.AudioConfig
}

// NewBaseTelephonyStreamer creates a new BaseTelephonyStreamer with telephony-specific
// entity references. Use TelephonyOption values to configure the source audio
// format and any BaseStreamer overrides.
//
// By default the source audio config is RAPIDA_AUDIO_CONFIG (linear16 16kHz),
// and input/output thresholds are automatically derived from that config
// (60 ms input, 20 ms output frames). Concrete streamers only need to provide
// WithSourceAudioConfig to declare their native format — everything else is
// handled by the BaseStreamer defaults.
//
// Example:
//
//	base := NewBaseTelephonyStreamer(logger, assistant, conv, vlt,
//	    WithSourceAudioConfig(audio.NewMulaw8khzMonoAudioConfig()),
//	)
func NewBaseTelephonyStreamer(
	logger commons.Logger,
	assistant *internal_assistant_entity.Assistant,
	assistantConversation *internal_conversation_entity.AssistantConversation,
	vlt *protos.VaultCredential,
	opts ...TelephonyOption,
) BaseTelephonyStreamer {
	tc := telephonyConfig{}
	for _, opt := range opts {
		opt(&tc)
	}

	sourceAudioCfg := tc.sourceAudioConfig
	if sourceAudioCfg == nil {
		sourceAudioCfg = RAPIDA_AUDIO_CONFIG
	}

	// Build base options: derive thresholds from the source audio config,
	// then allow caller to override via WithBaseOption.
	baseOpts := []channel_base.Option{
		channel_base.WithInputAudioConfig(sourceAudioCfg),
		channel_base.WithOutputAudioConfig(sourceAudioCfg),
	}
	baseOpts = append(baseOpts, tc.baseOpts...)

	resampler, _ := internal_audio_resampler.GetResampler(logger)
	return BaseTelephonyStreamer{
		BaseStreamer:          channel_base.NewBaseStreamer(logger, baseOpts...),
		assistant:             assistant,
		resampler:             resampler,
		assistantConversation: assistantConversation,
		encoder:               base64.StdEncoding,
		vaultCredential:       vlt,
		sourceAudioConfig:     sourceAudioCfg,
	}
}

// ============================================================================
// Telephony helpers
// ============================================================================

// CreateVoiceRequest resamples raw audio from the provider's native format
// to the internal Rapida format (linear16 16kHz) and wraps it in a
// ConversationUserMessage for downstream processing.
func (base *BaseTelephonyStreamer) CreateVoiceRequest(audioData []byte) *protos.ConversationUserMessage {
	resampled, err := base.resampler.Resample(audioData, base.sourceAudioConfig, RAPIDA_AUDIO_CONFIG)
	if err != nil {
		base.Logger.Warnw("Failed to resample input audio, forwarding raw bytes",
			"error", err.Error(),
			"source_format", base.sourceAudioConfig.GetAudioFormat(),
			"source_rate", base.sourceAudioConfig.GetSampleRate(),
		)
		resampled = audioData
	}
	return &protos.ConversationUserMessage{
		Message: &protos.ConversationUserMessage_Audio{
			Audio: resampled,
		},
	}
}

// GetAssistantDefinition returns the protobuf assistant definition.
func (base *BaseTelephonyStreamer) GetAssistantDefinition() *protos.AssistantDefinition {
	return &protos.AssistantDefinition{
		AssistantId: base.assistant.Id,
		Version:     utils.GetVersionString(base.assistant.AssistantProviderId),
	}
}

// GetConversationId returns the conversation ID.
func (base *BaseTelephonyStreamer) GetConversationId() uint64 {
	return base.assistantConversation.Id
}

// Encoder returns the base64 encoder used by the streamer.
func (base *BaseTelephonyStreamer) Encoder() *base64.Encoding {
	return base.encoder
}

// Credential returns the vault credential associated with the streamer.
func (base *BaseTelephonyStreamer) Credential() *protos.VaultCredential {
	return base.vaultCredential
}

// VaultCredential returns the vault credential associated with the streamer.
func (base *BaseTelephonyStreamer) VaultCredential() *protos.VaultCredential {
	return base.vaultCredential
}

// Resampler returns the audio resampler.
func (base *BaseTelephonyStreamer) Resampler() internal_type.AudioResampler {
	return base.resampler
}

// SourceAudioConfig returns the native audio format of the telephony provider.
func (base *BaseTelephonyStreamer) SourceAudioConfig() *protos.AudioConfig {
	return base.sourceAudioConfig
}

// CreateConnectionRequest builds the initial ConversationInitialization message.
func (base *BaseTelephonyStreamer) CreateConnectionRequest() *protos.ConversationInitialization {
	return &protos.ConversationInitialization{
		AssistantConversationId: base.GetConversationId(),
		Assistant:               base.GetAssistantDefinition(),
		StreamMode:              protos.StreamMode_STREAM_MODE_AUDIO,
	}
}

// GetAssistatntConversation returns the assistant conversation entity.
func (base *BaseTelephonyStreamer) GetAssistatntConversation() *internal_conversation_entity.AssistantConversation {
	return base.assistantConversation
}
