// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package internal_sip_telephony

import (
	"context"
	"io"
	"sync"
	"sync/atomic"
	"time"

	internal_audio "github.com/rapidaai/api/assistant-api/internal/audio"
	callcontext "github.com/rapidaai/api/assistant-api/internal/callcontext"
	internal_telephony_base "github.com/rapidaai/api/assistant-api/internal/channel/telephony/internal/base"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	sip_infra "github.com/rapidaai/api/assistant-api/sip/infra"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/protos"
)

// Streamer constants
const (
	// Audio buffer size threshold in milliseconds
	audioBufferThresholdMs = 60
	// RTP packet interval in milliseconds
	packetIntervalMs = 20
	// rtpLogInterval is the number of packets between periodic log entries
	rtpLogInterval = 50
)

// Streamer implements the TelephonyStreamer interface using native SIP signaling and RTP.
// No WebSocket needed — uses sipgo for signaling, RTP/UDP for audio.
//
// SIP manages its own context (derived from a parent session context) and
// concurrency primitives (RWMutex + atomic.Bool) because its lifecycle is
// tied to the SIP session rather than a simple background context.
type Streamer struct {
	internal_telephony_base.BaseTelephonyStreamer

	mu     sync.RWMutex
	closed atomic.Bool

	config     *sip_infra.Config
	session    *sip_infra.Session
	server     *sip_infra.Server
	rtpHandler *sip_infra.RTPHandler

	codec *sip_infra.Codec

	// SIP uses its own context derived from the session/parent context,
	// overriding the BaseStreamer context.
	ctx    context.Context
	cancel context.CancelFunc

	inputBuffer []byte
	configSent  atomic.Bool
}

// NewStreamer creates a SIP streamer.
//
// When sipSession is non-nil (inbound path), the streamer reuses the session's
// existing RTP handler from the global SIP server — no new server is created.
//
// When sipSession is nil (outbound / standalone path), a dedicated SIP server
// is spun up with its own RTP port pool and event handlers.
func NewStreamer(ctx context.Context,
	config *sip_infra.Config,
	logger commons.Logger,
	sipSession *sip_infra.Session,
	cc *callcontext.CallContext,
	vaultCred *protos.VaultCredential,
) (internal_type.Streamer, error) {
	streamerCtx, cancel := context.WithCancel(ctx)

	// Default codec — overridden below when an existing session carries a
	// negotiated codec.
	pcmu := sip_infra.CodecPCMU
	codec := &pcmu

	s := &Streamer{
		BaseTelephonyStreamer: internal_telephony_base.NewBaseTelephonyStreamer(
			logger, cc, vaultCred,
			internal_telephony_base.WithSourceAudioConfig(internal_audio.NewMulaw8khzMonoAudioConfig()),
		),
		config: config,
		codec:  codec,
		ctx:    streamerCtx,
		cancel: cancel,
	}

	// --- Inbound: reuse existing session's RTP handler ---
	if sipSession != nil {
		rtpHandler := sipSession.GetRTPHandler()
		if rtpHandler == nil {
			cancel()
			return nil, sip_infra.NewSIPError("NewStreamer", sipSession.GetCallID(), "session has no RTP handler", sip_infra.ErrRTPNotInitialized)
		}

		if negotiated := sipSession.GetNegotiatedCodec(); negotiated != nil {
			s.codec = negotiated
		}

		s.session = sipSession
		s.rtpHandler = rtpHandler

		go s.forwardIncomingAudio()

		localIP, localPort := rtpHandler.LocalAddr()
		logger.Infow("SIP streamer created (inbound)",
			"call_id", sipSession.GetCallID(),
			"codec", s.codec.Name,
			"rtp_port", localPort,
			"local_ip", localIP)

		return s, nil
	}

	// --- Outbound / standalone: spin up a dedicated SIP server ---
	// Address is the bind address (0.0.0.0 to listen on all interfaces).
	// config.Server is the remote SIP server — it must NOT be used as the
	// bind address or the RTP socket will fail to bind to a non-local IP.
	listenConfig := &sip_infra.ListenConfig{
		Address:   "0.0.0.0",
		Port:      config.Port,
		Transport: config.Transport,
	}

	tenantConfig := config
	configResolver := func(reqCtx *sip_infra.SIPRequestContext) (*sip_infra.InviteResult, error) {
		return &sip_infra.InviteResult{
			Config:      tenantConfig,
			ShouldAllow: true,
		}, nil
	}

	server, err := sip_infra.NewServer(streamerCtx, &sip_infra.ServerConfig{
		ListenConfig:      listenConfig,
		ConfigResolver:    configResolver,
		Logger:            logger,
		RTPPortRangeStart: config.RTPPortRangeStart,
		RTPPortRangeEnd:   config.RTPPortRangeEnd,
	})
	if err != nil {
		cancel()
		return nil, err
	}
	s.server = server

	server.SetOnInvite(s.handleInvite)
	server.SetOnBye(s.handleBye)
	server.SetOnError(s.handleError)

	if err := server.Start(); err != nil {
		cancel()
		return nil, err
	}

	logger.Infow("SIP streamer created (outbound)")
	return s, nil
}

func (s *Streamer) handleInvite(session *sip_infra.Session, fromURI, toURI string) error {
	s.mu.Lock()
	s.session = session
	codec := s.codec
	s.mu.Unlock()

	// Allocate RTP port from the server's port pool
	rtpPort, err := s.server.AllocateRTPPort()
	if err != nil {
		return sip_infra.NewSIPError("handleInvite", session.GetCallID(), "no RTP ports available", err)
	}

	// Initialize RTP handler for audio
	rtpHandler, err := sip_infra.NewRTPHandler(s.ctx, &sip_infra.RTPConfig{
		LocalIP:     s.config.Server,
		LocalPort:   rtpPort,
		PayloadType: codec.PayloadType,
		ClockRate:   codec.ClockRate,
		Logger:      s.Logger,
	})
	if err != nil {
		s.server.ReleaseRTPPort(rtpPort)
		return sip_infra.NewSIPError("handleInvite", session.GetCallID(), "failed to create RTP handler", err)
	}

	s.mu.Lock()
	s.rtpHandler = rtpHandler
	s.mu.Unlock()

	// Update session with local RTP address
	localIP, localPort := rtpHandler.LocalAddr()
	session.SetLocalRTP(localIP, localPort)
	session.SetNegotiatedCodec(codec.Name, int(codec.ClockRate))
	session.SetRTPHandler(rtpHandler)

	// Start RTP processing
	rtpHandler.Start()

	// Start audio forwarding
	go s.forwardIncomingAudio()

	s.Logger.Infow("SIP call established",
		"call_id", session.GetCallID(),
		"from", fromURI,
		"to", toURI,
		"codec", codec.Name)

	return nil
}

func (s *Streamer) handleBye(session *sip_infra.Session) error {
	s.Logger.Infow("BYE received, closing streamer", "call_id", session.GetCallID())
	return s.Close()
}

func (s *Streamer) handleError(session *sip_infra.Session, err error) {
	s.Logger.Error("SIP error occurred",
		"call_id", session.GetCallID(),
		"error", err)
}

func (s *Streamer) forwardIncomingAudio() {
	s.mu.RLock()
	rtpHandler := s.rtpHandler
	s.mu.RUnlock()

	if rtpHandler == nil {
		s.Logger.Error("forwardIncomingAudio: RTP handler is nil")
		return
	}

	s.Logger.Debug("forwardIncomingAudio: Started listening for RTP audio")
	packetCount := 0

	for {
		select {
		case <-s.ctx.Done():
			s.Logger.Debug("forwardIncomingAudio: Context cancelled", "total_packets", packetCount)
			return
		case audioData, ok := <-rtpHandler.AudioIn():
			if !ok {
				s.Logger.Debug("forwardIncomingAudio: Audio channel closed", "total_packets", packetCount)
				return
			}
			packetCount++
			s.mu.Lock()
			s.inputBuffer = append(s.inputBuffer, audioData...)
			bufLen := len(s.inputBuffer)
			s.mu.Unlock()

			// Log periodically (every 50 packets = ~1 second)
			if packetCount%rtpLogInterval == 1 {
				s.Logger.Debug("forwardIncomingAudio: Buffered audio",
					"packet_count", packetCount,
					"buffer_size", bufLen,
					"chunk_size", len(audioData))
			}
		}
	}
}

func (s *Streamer) Context() context.Context {
	return s.ctx
}

func (s *Streamer) Recv() (internal_type.Stream, error) {
	if s.closed.Load() {
		return nil, io.EOF
	}

	// Send connection/config request on first call
	if s.configSent.CompareAndSwap(false, true) {
		return s.CreateConnectionRequest(), nil
	}

	// Calculate buffer threshold based on codec sample rate
	sampleRate := int(s.codec.ClockRate)
	bufferSizeThreshold := sampleRate * audioBufferThresholdMs / 1000

	// Block until we have enough audio data or context is cancelled
	for {
		select {
		case <-s.ctx.Done():
			return nil, io.EOF
		default:
		}

		s.mu.RLock()
		session := s.session
		s.mu.RUnlock()

		// Check if session is active
		if session == nil || !session.IsActive() {
			return nil, io.EOF
		}

		s.mu.Lock()
		if len(s.inputBuffer) >= bufferSizeThreshold {
			audioData := make([]byte, bufferSizeThreshold)
			copy(audioData, s.inputBuffer[:bufferSizeThreshold])
			s.inputBuffer = s.inputBuffer[bufferSizeThreshold:]
			s.mu.Unlock()

			// Resample from native µ-law 8kHz to linear16 16kHz for downstream
			return s.CreateVoiceRequest(audioData), nil
		}
		s.mu.Unlock()

		// Wait a bit before checking again (20ms = typical RTP packet interval)
		select {
		case <-s.ctx.Done():
			return nil, io.EOF
		case <-time.After(packetIntervalMs * time.Millisecond):
			// Continue polling
		}
	}
}

// createConnectionRequest creates the initial connection request for the talker
// func (s *Streamer) createConnectionRequest() (internal_type.Stream, error) {
// 	inputConfig, outputConfig := s.GetAudioConfig()
// 	return &protos.AssistantTalkInput{
// 		Request: &protos.AssistantTalkInput_Configuration{
// 			Configuration: &protos.ConversationConfiguration{
// 				AssistantConversationId: s.conversation.Id,
// 				Assistant: &protos.AssistantDefinition{
// 					AssistantId: s.assistant.Id,
// 					Version:     "latest",
// 				},
// 				InputConfig:  &protos.StreamConfig{Audio: inputConfig},
// 				OutputConfig: &protos.StreamConfig{Audio: outputConfig},
// 			},
// 		},
// 	}, nil
// }

func (s *Streamer) Send(response internal_type.Stream) error {
	if s.closed.Load() {
		return sip_infra.ErrSessionClosed
	}

	switch data := response.(type) {
	case *protos.ConversationAssistantMessage:
		switch content := data.Message.(type) {
		case *protos.ConversationAssistantMessage_Audio:
			s.Logger.Debug("Send: Received audio output from assistant", "audio_size", len(content.Audio))
			return s.sendAudio(content.Audio)
		}
	case *protos.ConversationInterruption:
		s.Logger.Debug("Send: Received interruption", "type", data.Type)
		if data.Type == protos.ConversationInterruption_INTERRUPTION_TYPE_WORD {
			return s.handleInterruption()
		}
	case *protos.ConversationDirective:
		s.Logger.Debug("Send: Received directive", "type", data.GetType())
		if data.GetType() == protos.ConversationDirective_END_CONVERSATION {
			return s.Close()
		}
	}
	return nil
}

func (s *Streamer) sendAudio(audioData []byte) error {
	s.mu.RLock()
	rtpHandler := s.rtpHandler
	codec := s.codec
	s.mu.RUnlock()

	if rtpHandler == nil {
		return sip_infra.ErrRTPNotInitialized
	}

	// TTS always produces μ-law (PCMU) audio. If the negotiated codec is
	// PCMA (A-law), transcode each sample before handing it to the RTP handler.
	outData := audioData
	if codec != nil && codec.Name == "PCMA" {
		outData = mulawToAlaw(audioData)
	}

	select {
	case rtpHandler.AudioOut() <- outData:
		return nil
	case <-s.ctx.Done():
		return s.ctx.Err()
	default:
		s.Logger.Warn("sendAudio: RTP output channel full, dropping audio", "size", len(audioData))
		return nil
	}
}

func (s *Streamer) handleInterruption() error {
	s.mu.Lock()
	s.inputBuffer = nil // Clear input buffer on interruption
	s.mu.Unlock()

	s.Logger.Debug("Handled interruption, cleared audio buffers")
	return nil
}

// GetAudioConfig returns the audio configuration for this streamer.
// The TTS engine always produces μ-law 8kHz audio (the protobuf AudioFormat
// enum has no A-law value). When PCMA is negotiated, sendAudio() transcodes
// the μ-law samples to A-law before writing to the RTP handler.
func (s *Streamer) GetAudioConfig() (*protos.AudioConfig, *protos.AudioConfig) {
	// Always request μ-law from the TTS engine — it's the only G.711
	// variant the protobuf AudioFormat supports. Transcoding to PCMA
	// (if negotiated) happens in sendAudio().
	inputConfig := internal_audio.NewMulaw8khzMonoAudioConfig()
	outputConfig := internal_audio.NewMulaw8khzMonoAudioConfig()
	return inputConfig, outputConfig
}

// Close closes the streamer and releases all resources
func (s *Streamer) Close() error {
	if !s.closed.CompareAndSwap(false, true) {
		return nil // Already closed
	}
	// Cancel context first
	s.cancel()

	s.mu.Lock()
	rtpHandler := s.rtpHandler
	server := s.server
	session := s.session
	s.rtpHandler = nil
	s.server = nil
	s.session = nil
	s.inputBuffer = nil
	s.mu.Unlock()

	// Stop RTP handler
	if rtpHandler != nil {
		if err := rtpHandler.Stop(); err != nil {
			s.Logger.Warn("Error stopping RTP handler", "error", err)
		}
	}

	// Stop server
	if server != nil {
		server.Stop()
	}

	// End session
	if session != nil {
		session.End()
	}

	s.Logger.Infow("SIP streamer closed")
	return nil
}

// IsClosed returns whether the streamer has been closed
func (s *Streamer) IsClosed() bool {
	return s.closed.Load()
}

// GetSession returns the underlying SIP session
func (s *Streamer) GetSession() *sip_infra.Session {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.session
}

// mulawToAlaw converts μ-law (PCMU) encoded audio samples to A-law (PCMA).
// Uses the ITU G.711 standard conversion table (μ-law → linear → A-law).
// This is needed because the TTS engine always produces μ-law audio (the
// protobuf AudioFormat enum has no A-law variant), but the remote side may
// have negotiated PCMA.
func mulawToAlaw(in []byte) []byte {
	out := make([]byte, len(in))
	for i, b := range in {
		out[i] = mulaw2alaw[b]
	}
	return out
}

// mulaw2alaw is the direct μ-law → A-law conversion lookup table (ITU G.711).
// Index: μ-law byte value (0-255). Value: corresponding A-law byte value.
var mulaw2alaw = [256]byte{
	0x29, 0x2A, 0x27, 0x28, 0x2D, 0x2E, 0x2B, 0x2C, // 0-7
	0x21, 0x22, 0x1F, 0x20, 0x25, 0x26, 0x23, 0x24, // 8-15
	0x39, 0x3A, 0x37, 0x38, 0x3D, 0x3E, 0x3B, 0x3C, // 16-23
	0x31, 0x32, 0x2F, 0x30, 0x35, 0x36, 0x33, 0x34, // 24-31
	0x0A, 0x0B, 0x08, 0x09, 0x0E, 0x0F, 0x0C, 0x0D, // 32-39
	0x02, 0x03, 0x00, 0x01, 0x06, 0x07, 0x04, 0x05, // 40-47
	0x1A, 0x1B, 0x18, 0x19, 0x1E, 0x1F, 0x1C, 0x1D, // 48-55
	0x12, 0x13, 0x10, 0x11, 0x16, 0x17, 0x14, 0x15, // 56-63
	0x62, 0x63, 0x60, 0x61, 0x66, 0x67, 0x64, 0x65, // 64-71
	0x5D, 0x5D, 0x5C, 0x5C, 0x5F, 0x5F, 0x5E, 0x5E, // 72-79
	0x74, 0x76, 0x70, 0x72, 0x7C, 0x7E, 0x78, 0x7A, // 80-87
	0x6A, 0x6B, 0x68, 0x69, 0x6E, 0x6F, 0x6C, 0x6D, // 88-95
	0x48, 0x49, 0x46, 0x47, 0x4C, 0x4D, 0x4A, 0x4B, // 96-103
	0x40, 0x41, 0x3F, 0x3F, 0x44, 0x45, 0x42, 0x43, // 104-111
	0x56, 0x57, 0x54, 0x55, 0x5A, 0x5B, 0x58, 0x59, // 112-119
	0x4F, 0x4F, 0x4E, 0x4E, 0x52, 0x53, 0x50, 0x51, // 120-127
	0xA9, 0xAA, 0xA7, 0xA8, 0xAD, 0xAE, 0xAB, 0xAC, // 128-135
	0xA1, 0xA2, 0x9F, 0xA0, 0xA5, 0xA6, 0xA3, 0xA4, // 136-143
	0xB9, 0xBA, 0xB7, 0xB8, 0xBD, 0xBE, 0xBB, 0xBC, // 144-151
	0xB1, 0xB2, 0xAF, 0xB0, 0xB5, 0xB6, 0xB3, 0xB4, // 152-159
	0x8A, 0x8B, 0x88, 0x89, 0x8E, 0x8F, 0x8C, 0x8D, // 160-167
	0x82, 0x83, 0x80, 0x81, 0x86, 0x87, 0x84, 0x85, // 168-175
	0x9A, 0x9B, 0x98, 0x99, 0x9E, 0x9F, 0x9C, 0x9D, // 176-183
	0x92, 0x93, 0x90, 0x91, 0x96, 0x97, 0x94, 0x95, // 184-191
	0xE2, 0xE3, 0xE0, 0xE1, 0xE6, 0xE7, 0xE4, 0xE5, // 192-199
	0xDD, 0xDD, 0xDC, 0xDC, 0xDF, 0xDF, 0xDE, 0xDE, // 200-207
	0xF4, 0xF6, 0xF0, 0xF2, 0xFC, 0xFE, 0xF8, 0xFA, // 208-215
	0xEA, 0xEB, 0xE8, 0xE9, 0xEE, 0xEF, 0xEC, 0xED, // 216-223
	0xC8, 0xC9, 0xC6, 0xC7, 0xCC, 0xCD, 0xCA, 0xCB, // 224-231
	0xC0, 0xC1, 0xBF, 0xBF, 0xC4, 0xC5, 0xC2, 0xC3, // 232-239
	0xD6, 0xD7, 0xD4, 0xD5, 0xDA, 0xDB, 0xD8, 0xD9, // 240-247
	0xCF, 0xCF, 0xCE, 0xCE, 0xD2, 0xD3, 0xD0, 0xD1, // 248-255
}
