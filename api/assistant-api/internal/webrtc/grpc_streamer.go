// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package internal_webrtc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pion/interceptor"
	"github.com/pion/interceptor/pkg/intervalpli"
	"github.com/pion/rtp"
	pionwebrtc "github.com/pion/webrtc/v4"
	"github.com/pion/webrtc/v4/pkg/media"
	internal_audio "github.com/rapidaai/api/assistant-api/internal/audio"
	internal_audio_resampler "github.com/rapidaai/api/assistant-api/internal/audio/resampler"
	internal_streamers "github.com/rapidaai/api/assistant-api/internal/streamers"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	webrtc_internal "github.com/rapidaai/api/assistant-api/internal/webrtc/internal"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/protos"
)

// ============================================================================
// GrpcStreamer - WebRTC with gRPC signaling
// ============================================================================

// GrpcStreamer implements the Streamer interface using Pion WebRTC
// with gRPC bidirectional stream for signaling instead of WebSocket.
// Audio flows through WebRTC media tracks; gRPC is used for signaling.
type GrpcStreamer struct {
	mu sync.Mutex

	// Core components
	logger     commons.Logger
	config     *webrtc_internal.Config
	grpcStream internal_streamers.Streamer

	// Lifecycle
	ctx    context.Context
	cancel context.CancelFunc

	// Session state
	sessionID string

	// Pion WebRTC
	pc         *pionwebrtc.PeerConnection
	localTrack *pionwebrtc.TrackLocalStaticSample
	opusCodec  *webrtc_internal.OpusCodec

	// Single channel for all inputs to downstream
	inputCh chan *protos.AssistantTalkInput
	errCh   chan error

	// Buffer for incoming audio accumulation
	inputBuffer   *bytes.Buffer
	inputBufferMu sync.Mutex

	// Resampler
	resampler    internal_type.AudioResampler
	opusConfig   *protos.AudioConfig
	sttTtsConfig *protos.AudioConfig

	// Output audio queue
	outputBuffer   *bytes.Buffer
	outputBufferMu sync.Mutex
	outputStarted  bool
}

// NewGrpcStreamer creates a new WebRTC streamer with gRPC signaling
func NewGrpcStreamer(
	ctx context.Context,
	logger commons.Logger,
	grpcStream internal_streamers.Streamer,
) (internal_streamers.Streamer, error) {
	streamerCtx, cancel := context.WithCancel(ctx)

	resampler, err := internal_audio_resampler.GetResampler(logger)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create resampler: %w", err)
	}

	opusCodec, err := webrtc_internal.NewOpusCodec()
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create Opus codec: %w", err)
	}

	s := &GrpcStreamer{
		logger:     logger,
		config:     webrtc_internal.DefaultConfig(),
		grpcStream: grpcStream,
		ctx:        streamerCtx,
		cancel:     cancel,

		sessionID: uuid.New().String(),

		inputCh: make(chan *protos.AssistantTalkInput, 100),
		errCh:   make(chan error, 1),

		inputBuffer: new(bytes.Buffer),

		resampler:    resampler,
		opusConfig:   internal_audio.NewLinear48khzMonoAudioConfig(),
		sttTtsConfig: internal_audio.NewLinear16khzMonoAudioConfig(),
		opusCodec:    opusCodec,

		outputBuffer: new(bytes.Buffer),
	}

	// Create peer connection
	if err := s.createPeerConnection(); err != nil {
		cancel()
		return nil, err
	}

	// Initiate WebRTC handshake
	if err := s.initiateWebRTCHandshake(); err != nil {
		cancel()
		s.pc.Close()
		return nil, fmt.Errorf("failed to initiate WebRTC handshake: %w", err)
	}

	// Start gRPC message reader
	go s.readGrpcMessages()

	return s, nil
}

// ============================================================================
// Peer Connection Setup (same as WebSocket version)
// ============================================================================

func (s *GrpcStreamer) createPeerConnection() error {
	mediaEngine := &pionwebrtc.MediaEngine{}

	// Opus - primary codec
	if err := mediaEngine.RegisterCodec(pionwebrtc.RTPCodecParameters{
		RTPCodecCapability: pionwebrtc.RTPCodecCapability{
			MimeType:    pionwebrtc.MimeTypeOpus,
			ClockRate:   48000,
			Channels:    2,
			SDPFmtpLine: "minptime=10;useinbandfec=1;stereo=0;sprop-stereo=0",
		},
		PayloadType: 111,
	}, pionwebrtc.RTPCodecTypeAudio); err != nil {
		return fmt.Errorf("failed to register Opus: %w", err)
	}

	// Note: Only Opus codec registered. PCMU/PCMA removed for simplicity.

	// Interceptors
	registry := &interceptor.Registry{}
	if err := pionwebrtc.RegisterDefaultInterceptors(mediaEngine, registry); err != nil {
		return fmt.Errorf("failed to register interceptors: %w", err)
	}
	pli, err := intervalpli.NewReceiverInterceptor()
	if err != nil {
		return fmt.Errorf("failed to create PLI interceptor: %w", err)
	}
	registry.Add(pli)

	api := pionwebrtc.NewAPI(
		pionwebrtc.WithMediaEngine(mediaEngine),
		pionwebrtc.WithInterceptorRegistry(registry),
	)

	iceServers := make([]pionwebrtc.ICEServer, len(s.config.ICEServers))
	for i, srv := range s.config.ICEServers {
		iceServers[i] = pionwebrtc.ICEServer{
			URLs:       srv.URLs,
			Username:   srv.Username,
			Credential: srv.Credential,
		}
	}

	pcConfig := pionwebrtc.Configuration{ICEServers: iceServers}
	if s.config.ICETransportPolicy == "relay" {
		pcConfig.ICETransportPolicy = pionwebrtc.ICETransportPolicyRelay
	}

	pc, err := api.NewPeerConnection(pcConfig)
	if err != nil {
		return fmt.Errorf("failed to create peer connection: %w", err)
	}
	s.pc = pc

	s.setupPeerEventHandlers()
	return s.createLocalTrack()
}

func (s *GrpcStreamer) setupPeerEventHandlers() {
	// ICE candidates - send via gRPC using clean proto types
	s.pc.OnICECandidate(func(c *pionwebrtc.ICECandidate) {
		if c == nil {
			return
		}
		cJSON := c.ToJSON()
		ice := &webrtc_internal.ICECandidate{Candidate: cJSON.Candidate}
		if cJSON.SDPMid != nil {
			ice.SDPMid = *cJSON.SDPMid
		}
		if cJSON.SDPMLineIndex != nil {
			ice.SDPMLineIndex = int(*cJSON.SDPMLineIndex)
		}
		if cJSON.UsernameFragment != nil {
			ice.UsernameFragment = *cJSON.UsernameFragment
		}
		s.sendICECandidate(ice)
	})

	// Connection state
	s.pc.OnConnectionStateChange(func(state pionwebrtc.PeerConnectionState) {
		s.logger.Info("WebRTC connection state", "state", state.String(), "session", s.sessionID)

		s.mu.Lock()
		defer s.mu.Unlock()

		if state == pionwebrtc.PeerConnectionStateConnected && !s.outputStarted {
			s.outputStarted = true
			go s.runOutputSender()
		}
	})

	// Remote track (incoming audio)
	s.pc.OnTrack(func(track *pionwebrtc.TrackRemote, _ *pionwebrtc.RTPReceiver) {
		if track.Kind() != pionwebrtc.RTPCodecTypeAudio {
			return
		}
		s.logger.Info("Remote audio track received", "codec", track.Codec().MimeType)
		go s.readRemoteAudio(track)
	})
}

func (s *GrpcStreamer) createLocalTrack() error {
	track, err := pionwebrtc.NewTrackLocalStaticSample(
		pionwebrtc.RTPCodecCapability{
			MimeType:  pionwebrtc.MimeTypeOpus,
			ClockRate: webrtc_internal.OpusSampleRate,
			Channels:  2,
		},
		"audio",
		"rapida-voice-ai",
	)
	if err != nil {
		return fmt.Errorf("failed to create Opus track: %w", err)
	}

	if _, err := s.pc.AddTrack(track); err != nil {
		return fmt.Errorf("failed to add track: %w", err)
	}

	s.localTrack = track
	return nil
}

// ============================================================================
// Audio Processing (same as WebSocket version)
// ============================================================================

func (s *GrpcStreamer) readRemoteAudio(track *pionwebrtc.TrackRemote) {
	buf := make([]byte, 1500)
	mimeType := track.Codec().MimeType

	// Only Opus is supported
	if mimeType != pionwebrtc.MimeTypeOpus {
		s.logger.Error("Unsupported codec, only Opus is supported", "codec", mimeType)
		return
	}

	opusDecoder, err := webrtc_internal.NewOpusCodec()
	if err != nil {
		s.logger.Error("Failed to create Opus decoder", "error", err)
		return
	}

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
		}

		n, _, err := track.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			continue
		}

		pkt := &rtp.Packet{}
		if err := pkt.Unmarshal(buf[:n]); err != nil {
			continue
		}

		if len(pkt.Payload) == 0 {
			continue
		}

		pcm, err := opusDecoder.Decode(pkt.Payload)
		if err != nil {
			continue
		}

		resampled, err := s.resampler.Resample(pcm, s.opusConfig, s.sttTtsConfig)
		if err != nil {
			continue
		}

		s.bufferAndSendAudio(resampled)
	}
}

// ============================================================================
// gRPC Signaling - Using clean proto types
// ============================================================================

// sendConfig sends WebRTC configuration (ICE servers, codec info) to client
func (s *GrpcStreamer) sendConfig() error {
	iceServers := make([]*protos.ICEServer, len(s.config.ICEServers))
	for i, srv := range s.config.ICEServers {
		iceServers[i] = &protos.ICEServer{
			Urls:       srv.URLs,
			Username:   srv.Username,
			Credential: srv.Credential,
		}
	}

	return s.grpcStream.Send(&protos.AssistantTalkOutput{
		Code:    200,
		Success: true,
		Data: &protos.AssistantTalkOutput_Signaling{
			Signaling: &protos.ServerSignaling{
				SessionId: s.sessionID,
				Message: &protos.ServerSignaling_Config{
					Config: &protos.WebRTCConfig{
						IceServers: iceServers,
						AudioCodec: "opus",
						SampleRate: int32(webrtc_internal.OpusSampleRate),
					},
				},
			},
		},
	})
}

// sendOffer sends SDP offer to client
func (s *GrpcStreamer) sendOffer(sdp string) error {
	return s.grpcStream.Send(&protos.AssistantTalkOutput{
		Code:    200,
		Success: true,
		Data: &protos.AssistantTalkOutput_Signaling{
			Signaling: &protos.ServerSignaling{
				SessionId: s.sessionID,
				Message: &protos.ServerSignaling_Sdp{
					Sdp: &protos.WebRTCSDP{
						Type: protos.WebRTCSDP_OFFER,
						Sdp:  sdp,
					},
				},
			},
		},
	})
}

// sendICECandidate sends ICE candidate to client
func (s *GrpcStreamer) sendICECandidate(ice *webrtc_internal.ICECandidate) error {
	return s.grpcStream.Send(&protos.AssistantTalkOutput{
		Code:    200,
		Success: true,
		Data: &protos.AssistantTalkOutput_Signaling{
			Signaling: &protos.ServerSignaling{
				SessionId: s.sessionID,
				Message: &protos.ServerSignaling_IceCandidate{
					IceCandidate: &protos.ICECandidate{
						Candidate:        ice.Candidate,
						SdpMid:           ice.SDPMid,
						SdpMLineIndex:    int32(ice.SDPMLineIndex),
						UsernameFragment: ice.UsernameFragment,
					},
				},
			},
		},
	})
}

// sendClear sends clear/interrupt signal to client
func (s *GrpcStreamer) sendClear() error {
	return s.grpcStream.Send(&protos.AssistantTalkOutput{
		Code:    200,
		Success: true,
		Data: &protos.AssistantTalkOutput_Signaling{
			Signaling: &protos.ServerSignaling{
				SessionId: s.sessionID,
				Message:   &protos.ServerSignaling_Clear{Clear: true},
			},
		},
	})
}

// ============================================================================
// Streamer Interface Implementation
// ============================================================================

func (s *GrpcStreamer) Context() context.Context {
	return s.ctx
}

// readGrpcMessages reads from gRPC stream and routes messages
func (s *GrpcStreamer) readGrpcMessages() {
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
		}

		msg, err := s.grpcStream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				s.sendError(io.EOF)
				return
			}
			s.logger.Error("gRPC read error", "error", err)
			s.sendError(err)
			return
		}

		s.handleGrpcMessage(msg)
	}
}

// sendError sends error to errCh
func (s *GrpcStreamer) sendError(err error) {
	select {
	case s.errCh <- err:
	default:
	}
}

// sendInput sends input to inputCh
func (s *GrpcStreamer) sendInput(input *protos.AssistantTalkInput) {
	select {
	case s.inputCh <- input:
	case <-s.ctx.Done():
	}
}

// sendConfigUpstream sends configuration to upstream immediately
func (s *GrpcStreamer) sendConfigUpstream(config *protos.ConversationConfiguration) {
	audioConfig := internal_audio.NewLinear16khzMonoAudioConfig()
	config.InputConfig = &protos.StreamConfig{Audio: audioConfig}
	config.OutputConfig = &protos.StreamConfig{Audio: audioConfig}
	s.sendInput(&protos.AssistantTalkInput{
		Request: &protos.AssistantTalkInput_Configuration{
			Configuration: config,
		},
	})
}

// bufferAndSendAudio buffers audio and sends when threshold is reached
func (s *GrpcStreamer) bufferAndSendAudio(audio []byte) {
	const bufferThreshold = 32 * 60 // 60ms at 16kHz

	s.inputBufferMu.Lock()
	s.inputBuffer.Write(audio)

	if s.inputBuffer.Len() < bufferThreshold {
		s.inputBufferMu.Unlock()
		return
	}

	audioData := make([]byte, s.inputBuffer.Len())
	s.inputBuffer.Read(audioData)
	s.inputBufferMu.Unlock()

	s.sendInput(&protos.AssistantTalkInput{
		Request: &protos.AssistantTalkInput_Message{
			Message: &protos.ConversationUserMessage{
				Message: &protos.ConversationUserMessage_Audio{Audio: audioData},
			},
		},
	})
}

// handleGrpcMessage processes incoming gRPC message
func (s *GrpcStreamer) handleGrpcMessage(msg *protos.AssistantTalkInput) {
	switch msg.GetRequest().(type) {
	case *protos.AssistantTalkInput_Message:
		s.sendInput(msg)
		return
	case *protos.AssistantTalkInput_Signaling:
		s.handleClientSignaling(msg.GetSignaling())
		return
	case *protos.AssistantTalkInput_Configuration:
		s.handleGrpcConnect(msg.GetConfiguration())
	default:
		s.logger.Warn("Unknown gRPC message type received")
	}
}

// Recv receives the next input - simple channel read
func (s *GrpcStreamer) Recv() (*protos.AssistantTalkInput, error) {
	select {
	case <-s.ctx.Done():
		return nil, io.EOF
	case err := <-s.errCh:
		return nil, err
	case input := <-s.inputCh:
		return input, nil
	}
}

// handleClientSignaling processes client WebRTC signaling messages
func (s *GrpcStreamer) handleClientSignaling(signaling *protos.ClientSignaling) {
	switch msg := signaling.GetMessage().(type) {
	case *protos.ClientSignaling_Sdp:
		if msg.Sdp.GetType() == protos.WebRTCSDP_ANSWER {
			if err := s.pc.SetRemoteDescription(pionwebrtc.SessionDescription{
				Type: pionwebrtc.SDPTypeAnswer,
				SDP:  msg.Sdp.GetSdp(),
			}); err != nil {
				s.logger.Error("Failed to set remote description", "error", err)
			}
		}

	case *protos.ClientSignaling_IceCandidate:
		ice := msg.IceCandidate
		idx := uint16(ice.GetSdpMLineIndex())
		sdpMid := ice.GetSdpMid()
		usernameFragment := ice.GetUsernameFragment()
		if err := s.pc.AddICECandidate(pionwebrtc.ICECandidateInit{
			Candidate:        ice.GetCandidate(),
			SDPMid:           &sdpMid,
			SDPMLineIndex:    &idx,
			UsernameFragment: &usernameFragment,
		}); err != nil {
			s.logger.Error("Failed to add ICE candidate", "error", err)
		}

	case *protos.ClientSignaling_Disconnect:
		if msg.Disconnect {
			s.sendError(io.EOF)
			s.Close()
		}
	}
}

func (s *GrpcStreamer) handleGrpcConnect(config *protos.ConversationConfiguration) {
	s.sendConfigUpstream(config)
}

// initiateWebRTCHandshake sends config and creates/sends SDP offer.
func (s *GrpcStreamer) initiateWebRTCHandshake() error {
	if err := s.sendConfig(); err != nil {
		return fmt.Errorf("failed to send config: %w", err)
	}

	offer, err := s.createAndSetLocalOffer()
	if err != nil {
		return err
	}

	if err := s.sendOffer(offer.SDP); err != nil {
		return fmt.Errorf("failed to send offer: %w", err)
	}
	return nil
}

// createAndSetLocalOffer creates SDP offer and sets it as local description.
func (s *GrpcStreamer) createAndSetLocalOffer() (*pionwebrtc.SessionDescription, error) {
	offer, err := s.pc.CreateOffer(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create offer: %w", err)
	}

	if err := s.pc.SetLocalDescription(offer); err != nil {
		return nil, fmt.Errorf("failed to set local description: %w", err)
	}

	return &offer, nil
}

// Send sends output to the client
func (s *GrpcStreamer) Send(response *protos.AssistantTalkOutput) error {
	switch data := response.GetData().(type) {
	case *protos.AssistantTalkOutput_Assistant:
		switch content := data.Assistant.Message.(type) {
		case *protos.ConversationAssistantMessage_Audio:
			return s.sendAudio(content.Audio)
		case *protos.ConversationAssistantMessage_Text:
			// Send text via gRPC
			return s.grpcStream.Send(&protos.AssistantTalkOutput{
				Code:    200,
				Success: true,
				Data:    response.GetData(),
			})
		}

	case *protos.AssistantTalkOutput_Configuration:
		return s.grpcStream.Send(response)

	case *protos.AssistantTalkOutput_User:
		return s.grpcStream.Send(response)

	case *protos.AssistantTalkOutput_Interruption:
		if data.Interruption.Type == protos.ConversationInterruption_INTERRUPTION_TYPE_WORD {
			s.inputBufferMu.Lock()
			s.inputBuffer.Reset()
			s.inputBufferMu.Unlock()
			s.clearOutputBuffer()
			// Send clear signal via WebRTC signaling
			return s.sendClear()
		}

	case *protos.AssistantTalkOutput_Directive:
		s.grpcStream.Send(response)
		if data.Directive.GetType() == protos.ConversationDirective_END_CONVERSATION {
			return s.Close()
		}
		return nil
	case *protos.AssistantTalkOutput_Error:
		return s.grpcStream.Send(response)
	}
	return nil
}

func (s *GrpcStreamer) sendAudio(audio []byte) error {
	if len(audio) == 0 {
		return nil
	}

	audio48kHz, err := s.resampler.Resample(audio, s.sttTtsConfig, s.opusConfig)
	if err != nil {
		s.logger.Error("Resample to 48kHz failed", "error", err)
		return err
	}

	s.outputBufferMu.Lock()
	s.outputBuffer.Write(audio48kHz)
	s.outputBufferMu.Unlock()

	return nil
}

func (s *GrpcStreamer) runOutputSender() {
	ticker := time.NewTicker(webrtc_internal.OpusFrameDuration * time.Millisecond)
	defer ticker.Stop()

	chunk := make([]byte, webrtc_internal.OpusFrameBytes)

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.outputBufferMu.Lock()
			n, _ := s.outputBuffer.Read(chunk)
			s.outputBufferMu.Unlock()

			if n == 0 {
				continue
			}

			if n < webrtc_internal.OpusFrameBytes {
				for i := n; i < webrtc_internal.OpusFrameBytes; i++ {
					chunk[i] = 0
				}
			}

			encoded, err := s.opusCodec.Encode(chunk)
			if err != nil {
				continue
			}

			if s.localTrack != nil {
				s.localTrack.WriteSample(media.Sample{
					Data:     encoded,
					Duration: webrtc_internal.OpusFrameDuration * time.Millisecond,
				})
			}
		}
	}
}

func (s *GrpcStreamer) clearOutputBuffer() {
	s.outputBufferMu.Lock()
	s.outputBuffer.Reset()
	s.outputBufferMu.Unlock()
}

// Close closes the WebRTC connection
func (s *GrpcStreamer) Close() error {
	s.cancel()

	if s.pc != nil {
		s.pc.Close()
	}

	return nil
}
