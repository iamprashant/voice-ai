// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package channel_webrtc

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
	"github.com/pion/rtp"
	pionwebrtc "github.com/pion/webrtc/v4"
	"github.com/pion/webrtc/v4/pkg/media"
	internal_audio "github.com/rapidaai/api/assistant-api/internal/audio"
	internal_audio_resampler "github.com/rapidaai/api/assistant-api/internal/audio/resampler"
	webrtc_internal "github.com/rapidaai/api/assistant-api/internal/channel/webrtc/internal"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/protos"
	"google.golang.org/grpc"
)

// ============================================================================
// webrtcStreamer - WebRTC with gRPC signaling
// ============================================================================

// webrtcStreamer implements the Streamer interface using Pion WebRTC
// with gRPC bidirectional stream for signaling instead of WebSocket.
// Audio flows through WebRTC media tracks; gRPC is used for signaling.
type webrtcStreamer struct {
	mu sync.Mutex

	// Core components
	logger     commons.Logger
	config     *webrtc_internal.Config
	grpcStream grpc.BidiStreamingServer[protos.WebTalkRequest, protos.WebTalkResponse]

	// Lifecycle
	ctx    context.Context
	cancel context.CancelFunc

	// Session state
	sessionID string

	// Pion WebRTC
	pc         *pionwebrtc.PeerConnection
	localTrack *pionwebrtc.TrackLocalStaticSample

	// Audio processing (embedded — no separate AudioProcessor)
	resampler           internal_type.AudioResampler
	opusConfig          *protos.AudioConfig // 48kHz mono for Opus/WebRTC
	internalAudioConfig *protos.AudioConfig // 16kHz mono for STT/TTS
	opusCodec           *webrtc_internal.OpusCodec

	// Unified input channel: all downstream-bound messages (gRPC + audio) are pushed here.
	// Recv() simply reads from this channel.
	inputCh     chan internal_type.Stream
	inputBuffer bytes.Buffer
	inputBufMu  sync.Mutex

	// Output: Send() -> resample -> buffer -> runOutputSender -> WebRTC track
	outputBuffer bytes.Buffer
	outputBufMu  sync.Mutex

	errCh chan error

	// Output sender state
	outputStarted bool

	// Audio processing context - cancelled on audio disconnect/reconnect
	audioCtx    context.Context
	audioCancel context.CancelFunc
	audioWg     sync.WaitGroup // Tracks audio goroutines for clean shutdown

	// Connection mode: "text" or "audio"
	currentMode protos.StreamMode
}

// NewWebRTCStreamer creates a new WebRTC streamer with gRPC signaling
func NewWebRTCStreamer(
	ctx context.Context,
	logger commons.Logger,
	grpcStream grpc.BidiStreamingServer[protos.WebTalkRequest, protos.WebTalkResponse],
) (internal_type.Streamer, error) {
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

	s := &webrtcStreamer{
		logger:              logger,
		config:              webrtc_internal.DefaultConfig(),
		grpcStream:          grpcStream,
		ctx:                 streamerCtx,
		cancel:              cancel,
		sessionID:           uuid.New().String(),
		resampler:           resampler,
		opusConfig:          internal_audio.NewLinear48khzMonoAudioConfig(),
		internalAudioConfig: internal_audio.NewLinear16khzMonoAudioConfig(),
		opusCodec:           opusCodec,
		inputCh:             make(chan internal_type.Stream, webrtc_internal.InputChannelSize),
		errCh:               make(chan error, webrtc_internal.ErrorChannelSize),
		currentMode:         protos.StreamMode_STREAM_MODE_TEXT,
	}

	// Start background gRPC reader — pushes all non-signaling messages into inputCh
	go s.runGrpcReader()

	return s, nil
}

// ============================================================================
// Peer Connection Setup
// ============================================================================

// stopAudioProcessing cancels audio goroutines (runOutputSender, readRemoteAudio)
func (s *webrtcStreamer) stopAudioProcessing() {
	s.mu.Lock()
	if s.audioCancel != nil {
		s.audioCancel()
		s.audioCancel = nil
	}
	s.audioCtx = nil
	s.mu.Unlock()

	// Wait for audio goroutines to finish
	s.audioWg.Wait()
}

func (s *webrtcStreamer) createPeerConnection() error {
	// Create new audio context for this connection
	s.mu.Lock()
	s.audioCtx, s.audioCancel = context.WithCancel(s.ctx)
	s.mu.Unlock()

	mediaEngine := &pionwebrtc.MediaEngine{}

	// Opus - primary codec (48kHz mono)
	if err := mediaEngine.RegisterCodec(pionwebrtc.RTPCodecParameters{
		RTPCodecCapability: pionwebrtc.RTPCodecCapability{
			MimeType:    pionwebrtc.MimeTypeOpus,
			ClockRate:   webrtc_internal.OpusSampleRate,
			Channels:    webrtc_internal.OpusChannels,
			SDPFmtpLine: webrtc_internal.OpusSDPFmtpLine,
		},
		PayloadType: webrtc_internal.OpusPayloadType,
	}, pionwebrtc.RTPCodecTypeAudio); err != nil {
		return fmt.Errorf("failed to register Opus codec: %w", err)
	}

	// Interceptors (default includes NACK for audio packet recovery)
	registry := &interceptor.Registry{}
	if err := pionwebrtc.RegisterDefaultInterceptors(mediaEngine, registry); err != nil {
		return fmt.Errorf("failed to register interceptors: %w", err)
	}

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

	s.mu.Lock()
	s.pc = pc
	s.mu.Unlock()

	s.setupPeerEventHandlers()
	return s.createLocalTrack()
}

func (s *webrtcStreamer) setupPeerEventHandlers() {
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
		s.logger.Infow("WebRTC connection state changed", "state", state, "session", s.sessionID)
		s.mu.Lock()
		defer s.mu.Unlock()
		switch state {
		case pionwebrtc.PeerConnectionStateConnected:
			// Mark audio mode as active
			s.currentMode = protos.StreamMode_STREAM_MODE_AUDIO
			if !s.outputStarted {
				s.outputStarted = true
				go s.runOutputSender()
			}
			// Notify client that audio connection is ready
			go func() {
				if err := s.sendReady(); err != nil {
					s.logger.Error("Failed to send READY signal", "error", err)
				}
			}()
		case pionwebrtc.PeerConnectionStateFailed:
			s.logger.Warn("WebRTC connection failed", "session", s.sessionID)
			s.currentMode = protos.StreamMode_STREAM_MODE_TEXT
		case pionwebrtc.PeerConnectionStateClosed:
			s.logger.Infow("WebRTC connection closed", "session", s.sessionID)
			s.currentMode = protos.StreamMode_STREAM_MODE_TEXT
		case pionwebrtc.PeerConnectionStateDisconnected:
			// Temporary disconnect (ICE restart may recover)
			s.logger.Warn("WebRTC connection temporarily disconnected", "session", s.sessionID)
		}
	})

	// Remote track (incoming audio)
	s.pc.OnTrack(func(track *pionwebrtc.TrackRemote, _ *pionwebrtc.RTPReceiver) {
		if track.Kind() != pionwebrtc.RTPCodecTypeAudio {
			return
		}
		s.logger.Infow("Remote audio track received", "codec", track.Codec().MimeType)
		go s.readRemoteAudio(track)
	})
}

func (s *webrtcStreamer) createLocalTrack() error {
	track, err := pionwebrtc.NewTrackLocalStaticSample(
		pionwebrtc.RTPCodecCapability{
			MimeType:  pionwebrtc.MimeTypeOpus,
			ClockRate: webrtc_internal.OpusSampleRate,
			Channels:  webrtc_internal.OpusChannels,
		},
		"audio",
		"rapida-audio",
	)
	if err != nil {
		return fmt.Errorf("failed to create local audio track: %w", err)
	}

	if _, err := s.pc.AddTrack(track); err != nil {
		return fmt.Errorf("failed to add track: %w", err)
	}

	s.mu.Lock()
	s.localTrack = track
	s.mu.Unlock()
	return nil
}

// ============================================================================
// Input Audio: WebRTC track -> decode -> resample -> Recv()
// ============================================================================

// readRemoteAudio reads from the WebRTC remote track, decodes Opus to PCM,
// resamples from 48kHz to 16kHz, and pushes onto inputAudioCh for Recv().
func (s *webrtcStreamer) readRemoteAudio(track *pionwebrtc.TrackRemote) {
	s.audioWg.Add(1)
	defer s.audioWg.Done()

	s.mu.Lock()
	audioCtx := s.audioCtx
	s.mu.Unlock()

	if audioCtx == nil {
		return
	}

	mimeType := track.Codec().MimeType
	if mimeType != pionwebrtc.MimeTypeOpus {
		s.logger.Error("Unsupported codec, only Opus is supported", "codec", mimeType)
		return
	}

	opusDecoder, err := webrtc_internal.NewOpusCodec()
	if err != nil {
		s.logger.Error("Failed to create Opus decoder", "error", err)
		return
	}

	buf := make([]byte, webrtc_internal.RTPBufferSize)
	consecutiveErrors := 0

	for {
		select {
		case <-audioCtx.Done():
			return
		default:
		}

		n, _, err := track.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			consecutiveErrors++
			if consecutiveErrors >= webrtc_internal.MaxConsecutiveErrors {
				s.logger.Error("Too many consecutive read errors, stopping audio reader", "lastError", err)
				return
			}
			continue
		}
		consecutiveErrors = 0

		pkt := &rtp.Packet{}
		if err := pkt.Unmarshal(buf[:n]); err != nil {
			s.logger.Debug("Failed to unmarshal RTP packet", "error", err)
			continue
		}
		if len(pkt.Payload) == 0 {
			continue
		}

		// Decode Opus to PCM (48kHz)
		pcm, err := opusDecoder.Decode(pkt.Payload)
		if err != nil {
			s.logger.Debug("Opus decode failed", "error", err)
			continue
		}

		// Resample from 48kHz to 16kHz for downstream STT
		resampled, err := s.resampler.Resample(pcm, s.opusConfig, s.internalAudioConfig)
		if err != nil {
			s.logger.Debug("Audio resample failed", "error", err)
			continue
		}

		// Buffer and flush to channel when threshold is reached
		s.bufferAndSendInput(resampled)
	}
}

// bufferAndSendInput accumulates resampled audio and sends it to inputAudioCh
// when the buffer reaches the threshold.
func (s *webrtcStreamer) bufferAndSendInput(audio []byte) {
	s.inputBufMu.Lock()
	s.inputBuffer.Write(audio)

	if s.inputBuffer.Len() < webrtc_internal.InputBufferThreshold {
		s.inputBufMu.Unlock()
		return
	}

	audioData := make([]byte, s.inputBuffer.Len())
	s.inputBuffer.Read(audioData)
	s.inputBufMu.Unlock()

	msg := &protos.ConversationUserMessage{
		Message: &protos.ConversationUserMessage_Audio{Audio: audioData},
	}
	// Non-blocking send - drop if channel is full to avoid blocking the reader
	select {
	case s.inputCh <- msg:
	default:
		s.logger.Debug("Input channel full, dropping audio frame")
	}
}

// ============================================================================
// Output Audio: Send() -> resample -> buffer -> runOutputSender -> WebRTC track
// ============================================================================

// sendAudio resamples 16kHz TTS audio to 48kHz and buffers it for the output sender.
func (s *webrtcStreamer) sendAudio(audio []byte) error {
	if len(audio) == 0 {
		return nil
	}

	// Resample from 16kHz to 48kHz for Opus/WebRTC
	audio48kHz, err := s.resampler.Resample(audio, s.internalAudioConfig, s.opusConfig)
	if err != nil {
		s.logger.Error("Resample to 48kHz failed", "error", err)
		return err
	}

	s.outputBufMu.Lock()
	s.outputBuffer.Write(audio48kHz)
	s.outputBufMu.Unlock()
	return nil
}

// runOutputSender continuously reads from the output buffer, encodes to Opus,
// and writes samples to the WebRTC local track at consistent 20ms intervals.
func (s *webrtcStreamer) runOutputSender() {
	s.audioWg.Add(1)
	defer s.audioWg.Done()

	s.mu.Lock()
	audioCtx := s.audioCtx
	localTrack := s.localTrack
	s.mu.Unlock()

	if audioCtx == nil || localTrack == nil {
		s.logger.Error("runOutputSender called with nil audioCtx or localTrack")
		return
	}

	chunkDuration := webrtc_internal.OpusFrameDuration * time.Millisecond
	nextSendTime := time.Now().Add(chunkDuration)

	// Pre-encode a silence chunk
	silenceData := s.encodeSilence()

	for {
		select {
		case <-audioCtx.Done():
			return
		default:
		}

		// Wait until next send time
		now := time.Now()
		if sleepDuration := nextSendTime.Sub(now); sleepDuration > 0 {
			time.Sleep(sleepDuration)
		}

		// Schedule next send immediately to minimize drift
		nextSendTime = nextSendTime.Add(chunkDuration)
		if time.Now().After(nextSendTime) {
			nextSendTime = time.Now().Add(chunkDuration)
		}

		// Read one frame from the output buffer
		chunk := make([]byte, webrtc_internal.OpusFrameBytes)
		s.outputBufMu.Lock()
		n, _ := s.outputBuffer.Read(chunk)
		s.outputBufMu.Unlock()

		var encoded []byte
		if n > 0 {
			// Pad with silence if the chunk is not full
			if n < webrtc_internal.OpusFrameBytes {
				for i := n; i < webrtc_internal.OpusFrameBytes; i++ {
					chunk[i] = 0
				}
			}
			var err error
			encoded, err = s.opusCodec.Encode(chunk)
			if err != nil {
				s.logger.Debug("Opus encode failed", "error", err)
				encoded = silenceData
			}
		} else {
			encoded = silenceData
		}

		if err := localTrack.WriteSample(media.Sample{
			Data:     encoded,
			Duration: chunkDuration,
		}); err != nil {
			s.logger.Debug("Failed to write sample to track", "error", err)
		}
	}
}

// encodeSilence returns a pre-encoded Opus silence frame.
func (s *webrtcStreamer) encodeSilence() []byte {
	silence := make([]byte, webrtc_internal.OpusFrameBytes)
	encoded, err := s.opusCodec.Encode(silence)
	if err != nil {
		s.logger.Error("Failed to encode silence chunk", "error", err)
		return silence // raw fallback (shouldn't happen)
	}
	return encoded
}

// ============================================================================
// Signaling helpers
// ============================================================================

// sendConfig sends WebRTC configuration (ICE servers, codec info) to client
func (s *webrtcStreamer) sendConfig() error {
	iceServers := make([]*protos.ICEServer, len(s.config.ICEServers))
	for i, srv := range s.config.ICEServers {
		iceServers[i] = &protos.ICEServer{
			Urls:       srv.URLs,
			Username:   srv.Username,
			Credential: srv.Credential,
		}
	}

	return s.grpcStream.Send(&protos.WebTalkResponse{
		Code:    200,
		Success: true,
		Data: &protos.WebTalkResponse_Signaling{
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
func (s *webrtcStreamer) sendOffer(sdp string) error {
	return s.grpcStream.Send(&protos.WebTalkResponse{
		Code:    200,
		Success: true,
		Data: &protos.WebTalkResponse_Signaling{
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
func (s *webrtcStreamer) sendICECandidate(ice *webrtc_internal.ICECandidate) error {
	return s.grpcStream.Send(&protos.WebTalkResponse{
		Code:    200,
		Success: true,
		Data: &protos.WebTalkResponse_Signaling{
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

// sendReady sends ready signal to client
func (s *webrtcStreamer) sendReady() error {
	return s.grpcStream.Send(&protos.WebTalkResponse{
		Code:    200,
		Success: true,
		Data: &protos.WebTalkResponse_Signaling{
			Signaling: &protos.ServerSignaling{
				SessionId: s.sessionID,
				Message:   &protos.ServerSignaling_Ready{Ready: true},
			},
		},
	})
}

// sendClear sends clear/interrupt signal to client
func (s *webrtcStreamer) sendClear() error {
	return s.grpcStream.Send(&protos.WebTalkResponse{
		Code:    200,
		Success: true,
		Data: &protos.WebTalkResponse_Signaling{
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

func (s *webrtcStreamer) Context() context.Context {
	return s.ctx
}

// Recv reads the next downstream-bound message from the unified input channel.
// Both gRPC messages and decoded WebRTC audio are fed into the same channel
// by background goroutines.
func (s *webrtcStreamer) Recv() (internal_type.Stream, error) {
	select {
	case <-s.ctx.Done():
		return nil, io.EOF
	case msg, ok := <-s.inputCh:
		if !ok {
			return nil, io.EOF
		}
		return msg, nil
	case err := <-s.errCh:
		return nil, err
	}
}

// runGrpcReader reads from the gRPC stream in a loop and pushes
// non-signaling messages into inputCh. Signaling is handled internally.
// Runs until the gRPC stream closes or the context is cancelled.
func (s *webrtcStreamer) runGrpcReader() {
	for {
		msg, err := s.grpcStream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				s.sendError(io.EOF)
			} else {
				s.sendError(fmt.Errorf("failed to receive gRPC message: %w", err))
			}
			return
		}
		switch msg.GetRequest().(type) {
		case *protos.WebTalkRequest_Initialization:
			s.logger.Debug("Received initialization %+v", msg.GetInitialization())
			s.pushInput(msg.GetInitialization())
		case *protos.WebTalkRequest_Configuration:
			s.logger.Debug("Received configuration %+v", msg.GetConfiguration())
			s.handleConfigurationMessage(msg.GetConfiguration())
			s.pushInput(msg.GetConfiguration())
		case *protos.WebTalkRequest_Message:
			s.logger.Debug("Received message %+v", msg.GetMessage())
			s.pushInput(msg.GetMessage())
		case *protos.WebTalkRequest_Signaling:
			s.logger.Debugw("Received signaling", "session", s.sessionID)
			s.handleClientSignaling(msg.GetSignaling())
			// Handled internally — not forwarded downstream
		default:
			s.logger.Warn("Unknown message type", "type", fmt.Sprintf("%T", msg.GetRequest()))
		}
	}
}

// pushInput sends a message to the unified input channel (non-blocking).
func (s *webrtcStreamer) pushInput(msg internal_type.Stream) {
	select {
	case s.inputCh <- msg:
	default:
		s.logger.Warn("Input channel full, dropping message")
	}
}

// sendError sends error to errCh
func (s *webrtcStreamer) sendError(err error) {
	select {
	case s.errCh <- err:
	default:
		s.logger.Warn("Error channel full, dropping error", "error", err)
	}
}

// handleConfigurationMessage processes transport mode changes.
// Switching text <-> audio only changes I/O transport - it does NOT create a new session.
func (s *webrtcStreamer) handleConfigurationMessage(config *protos.ConversationConfiguration) {
	s.mu.Lock()
	currentMode := s.currentMode
	s.mu.Unlock()

	if config.GetStreamMode() == currentMode {
		return
	}

	switch config.GetStreamMode() {
	case protos.StreamMode_STREAM_MODE_AUDIO:
		if err := s.setupAudioAndHandshake(); err != nil {
			s.logger.Error("Failed to setup audio", "error", err)
			s.resetAudioSession()
		}
	case protos.StreamMode_STREAM_MODE_TEXT:
		s.resetAudioSession()
	}
}

// handleClientSignaling processes client WebRTC signaling messages
func (s *webrtcStreamer) handleClientSignaling(signaling *protos.ClientSignaling) {
	s.mu.Lock()
	pc := s.pc
	s.mu.Unlock()

	switch msg := signaling.GetMessage().(type) {
	case *protos.ClientSignaling_Sdp:
		if msg.Sdp.GetType() == protos.WebRTCSDP_ANSWER {
			if pc == nil {
				s.logger.Warn("Received SDP answer but peer connection is nil, ignoring")
				return
			}
			if err := pc.SetRemoteDescription(pionwebrtc.SessionDescription{
				Type: pionwebrtc.SDPTypeAnswer,
				SDP:  msg.Sdp.GetSdp(),
			}); err != nil {
				s.logger.Error("Failed to set remote description", "error", err)
			}
		}

	case *protos.ClientSignaling_IceCandidate:
		if pc == nil {
			s.logger.Warn("Received ICE candidate but peer connection is nil, ignoring")
			return
		}
		ice := msg.IceCandidate
		idx := uint16(ice.GetSdpMLineIndex())
		sdpMid := ice.GetSdpMid()
		usernameFragment := ice.GetUsernameFragment()
		if err := pc.AddICECandidate(pionwebrtc.ICECandidateInit{
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

func (s *webrtcStreamer) resetAudioSession() {
	s.stopAudioProcessing()
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.pc != nil {
		s.pc.Close()
		s.pc = nil
	}
	s.localTrack = nil
	s.outputStarted = false
	s.currentMode = protos.StreamMode_STREAM_MODE_TEXT
}

// setupAudioAndHandshake tears down any stale peer connection, creates a fresh
// one, and initiates the WebRTC handshake (config -> offer -> answer -> ICE).
func (s *webrtcStreamer) setupAudioAndHandshake() error {
	// Always start fresh
	s.mu.Lock()
	if s.pc != nil {
		s.pc.Close()
		s.pc = nil
		s.localTrack = nil
	}
	s.mu.Unlock()

	if err := s.createPeerConnection(); err != nil {
		return fmt.Errorf("failed to create peer connection: %w", err)
	}

	return s.initiateWebRTCHandshake()
}

// initiateWebRTCHandshake sends config and creates/sends SDP offer.
func (s *webrtcStreamer) initiateWebRTCHandshake() error {
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
func (s *webrtcStreamer) createAndSetLocalOffer() (*pionwebrtc.SessionDescription, error) {
	offer, err := s.pc.CreateOffer(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create offer: %w", err)
	}

	if err := s.pc.SetLocalDescription(offer); err != nil {
		return nil, fmt.Errorf("failed to set local description: %w", err)
	}

	return &offer, nil
}

// ============================================================================
// Send - output to client
// ============================================================================

// Send sends output to the client
func (s *webrtcStreamer) Send(response internal_type.Stream) error {
	switch data := response.(type) {
	case *protos.ConversationAssistantMessage:
		switch content := data.Message.(type) {
		case *protos.ConversationAssistantMessage_Audio:
			return s.sendAudio(content.Audio)
		case *protos.ConversationAssistantMessage_Text:
			return s.grpcStream.Send(&protos.WebTalkResponse{
				Code:    200,
				Success: true,
				Data:    &protos.WebTalkResponse_Assistant{Assistant: data},
			})
		}
	case *protos.ConversationConfiguration:
		return s.grpcStream.Send(&protos.WebTalkResponse{
			Code:    200,
			Success: true,
			Data:    &protos.WebTalkResponse_Configuration{Configuration: data},
		})
	case *protos.ConversationInitialization:
		return s.grpcStream.Send(&protos.WebTalkResponse{
			Code:    200,
			Success: true,
			Data:    &protos.WebTalkResponse_Initialization{Initialization: data},
		})
	case *protos.ConversationUserMessage:
		return s.grpcStream.Send(&protos.WebTalkResponse{
			Code:    200,
			Success: true,
			Data:    &protos.WebTalkResponse_User{User: data},
		})

	case *protos.ConversationInterruption:
		if data.Type == protos.ConversationInterruption_INTERRUPTION_TYPE_WORD {
			s.clearInputBuffer()
			s.clearOutputBuffer()
			s.sendClear()
			return s.grpcStream.Send(&protos.WebTalkResponse{
				Code:    200,
				Success: true,
				Data:    &protos.WebTalkResponse_Interruption{Interruption: data},
			})
		}

	case *protos.ConversationDirective:
		if err := s.grpcStream.Send(&protos.WebTalkResponse{
			Code:    200,
			Success: true,
			Data:    &protos.WebTalkResponse_Directive{Directive: data},
		}); err != nil {
			s.logger.Error("Failed to send directive", "error", err)
			return err
		}
		if data.GetType() == protos.ConversationDirective_END_CONVERSATION {
			s.logger.Infow("Ending conversation", "session", s.sessionID)
			return s.Close()
		}
		return nil
	case *protos.ConversationError:
		return s.grpcStream.Send(&protos.WebTalkResponse{
			Code:    200,
			Success: true,
			Data:    &protos.WebTalkResponse_Error{Error: data},
		})
	}
	return nil
}

// ============================================================================
// Buffer helpers
// ============================================================================

func (s *webrtcStreamer) clearInputBuffer() {
	s.inputBufMu.Lock()
	s.inputBuffer.Reset()
	s.inputBufMu.Unlock()

	// Drain audio messages from the input channel
	for {
		select {
		case <-s.inputCh:
		default:
			return
		}
	}
}

func (s *webrtcStreamer) clearOutputBuffer() {
	s.outputBufMu.Lock()
	s.outputBuffer.Reset()
	s.outputBufMu.Unlock()
}

// ============================================================================
// Lifecycle
// ============================================================================

// Close closes the WebRTC connection and releases all resources
func (s *webrtcStreamer) Close() error {
	// Stop audio processing goroutines first
	s.stopAudioProcessing()

	// Cancel main context
	s.cancel()

	s.mu.Lock()
	if s.pc != nil {
		s.pc.Close()
		s.pc = nil
	}
	s.localTrack = nil
	s.mu.Unlock()

	return nil
}
