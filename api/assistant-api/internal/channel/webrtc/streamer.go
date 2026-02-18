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
	"google.golang.org/protobuf/types/known/timestamppb"
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

	// Disconnect tracking
	closed bool // true once closeWithReason has run

	// Pion WebRTC
	pc         *pionwebrtc.PeerConnection
	localTrack *pionwebrtc.TrackLocalStaticSample
	resampler  internal_type.AudioResampler
	opusCodec  *webrtc_internal.OpusCodec

	// inputCh: all downstream-bound messages (gRPC + decoded audio) funnelled here.
	// recv (non-blocking) -> inputCh -> loop (Recv) -> downstream service
	inputCh chan internal_type.Stream

	inputAudioBuffer     *bytes.Buffer
	inputAudioBufferLock sync.Mutex

	// outputCh: all upstream-bound messages funnelled here to preserve ordering.
	// send (non-blocking) -> outputCh -> loop (runOutputWriter) -> upstream service
	outputCh chan internal_type.Stream

	// outputAudioBuffer accumulates resampled 48kHz PCM and flushes complete
	// 20ms Opus frames into outputCh — mirrors inputAudioBuffer.
	outputAudioBuffer     *bytes.Buffer
	outputAudioBufferLock sync.Mutex

	// flushAudioCh signals runOutputWriter to discard its pending audio queue
	// (used on interruption to silence stale frames immediately).
	flushAudioCh chan struct{}

	// Audio processing context - cancelled on audio disconnect/reconnect
	audioCtx    context.Context
	audioCancel context.CancelFunc
	audioWg     sync.WaitGroup // Tracks audio goroutines for clean shutdown

	currentMode protos.StreamMode
}

// NewWebRTCStreamer creates a new WebRTC streamer with gRPC signaling.
// The streamer owns its own context (derived from context.Background) so that
// cleanup is never short-circuited by the caller's context being cancelled first.
// A separate goroutine watches the caller's context and triggers a graceful close.
func NewWebRTCStreamer(
	ctx context.Context,
	logger commons.Logger,
	grpcStream grpc.BidiStreamingServer[protos.WebTalkRequest, protos.WebTalkResponse],
) (internal_type.Streamer, error) {
	streamerCtx, cancel := context.WithCancel(context.Background())
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
		logger:     logger,
		config:     webrtc_internal.DefaultConfig(),
		grpcStream: grpcStream,
		ctx:        streamerCtx,
		cancel:     cancel,
		sessionID:  uuid.New().String(),
		resampler:  resampler,
		opusCodec:  opusCodec,

		inputCh:  make(chan internal_type.Stream, webrtc_internal.InputChannelSize),
		outputCh: make(chan internal_type.Stream, webrtc_internal.OutputChannelSize),

		inputAudioBuffer:  new(bytes.Buffer),
		outputAudioBuffer: new(bytes.Buffer),
		flushAudioCh:      make(chan struct{}, 1),
		currentMode:       protos.StreamMode_STREAM_MODE_TEXT,
	}

	// Start background loops
	go s.runGrpcReader()   // inputCh feeder
	go s.runOutputWriter() // outputCh consumer

	// Watch the caller's context so a cancelled parent triggers graceful close
	// rather than an abrupt context cancellation mid-cleanup.
	go s.watchCallerContext(ctx)

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
	s.audioWg.Wait()
}

func (s *webrtcStreamer) createPeerConnection() error {
	// Create new audio context and fresh output channel for this connection
	s.mu.Lock()
	s.audioCtx, s.audioCancel = context.WithCancel(s.ctx)
	s.mu.Unlock()

	mediaEngine := &pionwebrtc.MediaEngine{}
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

		// Update mode under lock, then release before any channel operations
		// to avoid holding mu while pushing to outputCh.
		s.mu.Lock()
		switch state {
		case pionwebrtc.PeerConnectionStateConnected:
			s.currentMode = protos.StreamMode_STREAM_MODE_AUDIO
		case pionwebrtc.PeerConnectionStateFailed,
			pionwebrtc.PeerConnectionStateClosed,
			pionwebrtc.PeerConnectionStateDisconnected:
			s.currentMode = protos.StreamMode_STREAM_MODE_TEXT
		}
		s.mu.Unlock()

		// Perform channel / lifecycle operations outside the lock.
		switch state {
		case pionwebrtc.PeerConnectionStateConnected:
			s.sendReady()

		case pionwebrtc.PeerConnectionStateFailed:
			// Connection failed irrecoverably — tear down and notify downstream.
			s.logger.Errorw("WebRTC connection failed, closing session", "session", s.sessionID)
			s.pushDisconnection(protos.ConversationDisconnection_DISCONNECTION_TYPE_USER)

		case pionwebrtc.PeerConnectionStateDisconnected:
			// Transient state — network hiccup, ICE may recover.
			// Only reset audio; do NOT close the gRPC stream/context so the
			// session can continue in text mode or reconnect.
			s.logger.Warnw("WebRTC peer disconnected, resetting audio", "session", s.sessionID)
			s.resetAudioSession()
		}
	})

	// Remote track (incoming audio)
	s.pc.OnTrack(func(track *pionwebrtc.TrackRemote, _ *pionwebrtc.RTPReceiver) {
		if track.Kind() != pionwebrtc.RTPCodecTypeAudio {
			return
		}
		s.logger.Infow("Remote audio track received", "codec", track.Codec().MimeType)
		// Add to WaitGroup before launching goroutine to prevent
		// audioWg.Wait() from racing with audioWg.Add(1).
		s.audioWg.Add(1)
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
	defer s.audioWg.Done()

	s.mu.Lock()
	audioCtx := s.audioCtx
	s.mu.Unlock()

	if audioCtx == nil {
		return
	}

	mimeType := track.Codec().MimeType
	if mimeType != pionwebrtc.MimeTypeOpus {
		s.logger.Errorw("Unsupported codec, only Opus is supported", "codec", mimeType)
		return
	}

	opusDecoder, err := webrtc_internal.NewOpusCodec()
	if err != nil {
		s.logger.Errorw("Failed to create Opus decoder", "error", err)
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
				s.logger.Errorw("Too many consecutive read errors, stopping audio reader", "lastError", err)
				return
			}
			continue
		}
		consecutiveErrors = 0

		pkt := &rtp.Packet{}
		if err := pkt.Unmarshal(buf[:n]); err != nil {
			s.logger.Debugw("Failed to unmarshal RTP packet", "error", err)
			continue
		}
		if len(pkt.Payload) == 0 {
			continue
		}

		// Decode Opus to PCM (48kHz)
		pcm, err := opusDecoder.Decode(pkt.Payload)
		if err != nil {
			s.logger.Debugw("Opus decode failed", "error", err, "payloadSize", len(pkt.Payload))
			continue
		}
		// resample to 16kHz
		resampled, err := s.resampler.Resample(pcm, internal_audio.WEBRTC_AUDIO_CONFIG, internal_audio.RAPIDA_INTERNAL_AUDIO_CONFIG)
		if err != nil {
			s.logger.Debugw("Audio resample failed", "error", err)
			continue
		}

		// Buffer and flush to channel when threshold is reached
		s.bufferAndSendInput(resampled)
	}
}

// bufferAndSendInput accumulates resampled audio and sends it to inputCh
// when the buffer reaches the threshold.
func (s *webrtcStreamer) bufferAndSendInput(audio []byte) {
	s.inputAudioBufferLock.Lock()
	s.inputAudioBuffer.Write(audio)

	if s.inputAudioBuffer.Len() < webrtc_internal.InputBufferThreshold {
		s.inputAudioBufferLock.Unlock()
		return
	}

	audioData := make([]byte, s.inputAudioBuffer.Len())
	s.inputAudioBuffer.Read(audioData)
	s.inputAudioBufferLock.Unlock()

	s.pushInput(&protos.ConversationUserMessage{
		Message: &protos.ConversationUserMessage_Audio{Audio: audioData},
		Time:    timestamppb.Now(),
	})
}

// bufferAndSendOutput accumulates resampled 48kHz PCM and flushes consistent
// 20ms frames into outputCh as ConversationAssistantMessage_Audio messages.
// Opus encoding happens later in runOutputWriter at pacing time.
//
// audio received -> push to outputAudioBuffer -> check size ->
// push ConversationAssistantMessage_Audio -> outputCh
func (s *webrtcStreamer) bufferAndSendOutput(audio48kHz []byte) {
	s.outputAudioBufferLock.Lock()
	s.outputAudioBuffer.Write(audio48kHz)

	if s.outputAudioBuffer.Len() < webrtc_internal.OutputBufferThreshold {
		s.outputAudioBufferLock.Unlock()
		return
	}

	// Flush as many complete 20ms PCM frames as possible.
	for s.outputAudioBuffer.Len() >= webrtc_internal.OpusFrameBytes {
		frame := make([]byte, webrtc_internal.OpusFrameBytes)
		s.outputAudioBuffer.Read(frame)
		s.outputAudioBufferLock.Unlock()

		// Push raw PCM frame; Opus encoding is deferred to runOutputWriter.
		s.pushOutput(&protos.ConversationAssistantMessage{
			Message: &protos.ConversationAssistantMessage_Audio{Audio: frame},
			Time:    timestamppb.Now(),
		})

		s.outputAudioBufferLock.Lock()
	}
	s.outputAudioBufferLock.Unlock()
}

// runOutputWriter is the single output loop:
//
//	outputCh -> loop (process) -> upstream service
//
// All outbound messages flow through outputCh to preserve ordering.
// Raw proto types and pre-built *WebTalkResponse (signaling) are accepted.
// The writer wraps raw types into WebTalkResponse before sending to gRPC.
//
//   - ConversationAssistantMessage_Audio → queue raw PCM → Opus-encode → WebRTC track
//     (paced at 20ms real-time intervals to smooth TTS bursts)
//   - *protos.WebTalkResponse (signaling) → send directly to gRPC
//   - All other raw types → wrap in WebTalkResponse → send to gRPC
//
// Runs for the lifetime of the streamer (exits when ctx is cancelled).
func (s *webrtcStreamer) runOutputWriter() {
	ticker := time.NewTicker(time.Duration(webrtc_internal.OutputPaceInterval) * time.Millisecond)
	defer ticker.Stop()

	// pendingAudio holds raw 20ms PCM frames waiting for the next tick.
	var pendingAudio [][]byte

	for {
		select {
		case <-s.ctx.Done():
			return

		case <-s.flushAudioCh:
			// Interruption: discard all queued audio immediately.
			pendingAudio = pendingAudio[:0]

		case <-ticker.C:
			// Encode and send one paced audio frame per tick (20ms real-time).
			if len(pendingAudio) > 0 {
				encoded, err := s.opusCodec.Encode(pendingAudio[0])
				if err != nil {
					s.logger.Debugw("Opus encode failed", "error", err)
				} else {
					s.writeAudioFrame(encoded)
				}
				pendingAudio = pendingAudio[1:]
			}

		case msg := <-s.outputCh:
			// Assistant audio → queue raw PCM for paced Opus encoding.
			if m, ok := msg.(*protos.ConversationAssistantMessage); ok {
				if audio, ok := m.Message.(*protos.ConversationAssistantMessage_Audio); ok {
					pendingAudio = append(pendingAudio, audio.Audio)
					continue
				}
			}

			// Wrap raw types in WebTalkResponse and send to gRPC.
			if resp := s.buildGRPCResponse(msg); resp != nil {
				s.dispatchOutput(resp)
			}
		}
	}
}

// buildGRPCResponse wraps a raw proto type into a WebTalkResponse for gRPC.
// Pre-built *WebTalkResponse (e.g. signaling) are passed through as-is.
func (s *webrtcStreamer) buildGRPCResponse(msg internal_type.Stream) *protos.WebTalkResponse {
	resp := &protos.WebTalkResponse{Code: 200, Success: true}
	switch m := msg.(type) {
	case *protos.ConversationAssistantMessage:
		resp.Data = &protos.WebTalkResponse_Assistant{Assistant: m}
	case *protos.ConversationConfiguration:
		resp.Data = &protos.WebTalkResponse_Configuration{Configuration: m}
	case *protos.ConversationInitialization:
		resp.Data = &protos.WebTalkResponse_Initialization{Initialization: m}
	case *protos.ConversationUserMessage:
		resp.Data = &protos.WebTalkResponse_User{User: m}
	case *protos.ConversationInterruption:
		resp.Data = &protos.WebTalkResponse_Interruption{Interruption: m}
	case *protos.ConversationDirective:
		resp.Data = &protos.WebTalkResponse_Directive{Directive: m}
	case *protos.ConversationError:
		resp.Data = &protos.WebTalkResponse_Error{Error: m}
	case *protos.ConversationMetadata:
		resp.Data = &protos.WebTalkResponse_Metadata{Metadata: m}
	case *protos.ConversationMetric:
		resp.Data = &protos.WebTalkResponse_Metric{Metric: m}
	case *protos.ServerSignaling:
		resp.Data = &protos.WebTalkResponse_Signaling{Signaling: m}
	default:
		s.logger.Warnw("Unknown output message type, skipping", "type", fmt.Sprintf("%T", msg))
		return nil
	}
	return resp
}

// dispatchOutput sends a WebTalkResponse directly to the gRPC stream.
func (s *webrtcStreamer) dispatchOutput(resp *protos.WebTalkResponse) {
	if err := s.grpcStream.Send(resp); err != nil {
		s.logger.Errorw("Failed to send gRPC response", "error", err)
	}
}

// writeAudioFrame writes an encoded Opus frame to the WebRTC local track.
func (s *webrtcStreamer) writeAudioFrame(data []byte) {
	s.mu.Lock()
	track := s.localTrack
	s.mu.Unlock()

	if track == nil {
		return
	}
	if err := track.WriteSample(media.Sample{
		Data:     data,
		Duration: webrtc_internal.OpusFrameDuration * time.Millisecond,
	}); err != nil {
		s.logger.Debugw("Failed to write sample to track", "error", err)
	}
}

// ============================================================================
// Signaling helpers
// ============================================================================

// sendConfig sends WebRTC configuration (ICE servers, codec info) to client via outputCh.
func (s *webrtcStreamer) sendConfig() {
	iceServers := make([]*protos.ICEServer, len(s.config.ICEServers))
	for i, srv := range s.config.ICEServers {
		iceServers[i] = &protos.ICEServer{
			Urls:       srv.URLs,
			Username:   srv.Username,
			Credential: srv.Credential,
		}
	}

	s.pushOutput(&protos.ServerSignaling{
		SessionId: s.sessionID,
		Message: &protos.ServerSignaling_Config{
			Config: &protos.WebRTCConfig{
				IceServers: iceServers,
				AudioCodec: "opus",
				SampleRate: int32(webrtc_internal.OpusSampleRate),
			},
		},
	},
	)
}

// sendOffer sends SDP offer to client via outputCh.
func (s *webrtcStreamer) sendOffer(sdp string) {
	s.pushOutput(&protos.ServerSignaling{
		SessionId: s.sessionID,
		Message: &protos.ServerSignaling_Sdp{
			Sdp: &protos.WebRTCSDP{
				Type: protos.WebRTCSDP_OFFER,
				Sdp:  sdp,
			},
		},
	})
}

// sendICECandidate sends ICE candidate to client via outputCh.
func (s *webrtcStreamer) sendICECandidate(ice *webrtc_internal.ICECandidate) {
	s.pushOutput(&protos.ServerSignaling{
		SessionId: s.sessionID,
		Message: &protos.ServerSignaling_IceCandidate{
			IceCandidate: &protos.ICECandidate{
				Candidate:        ice.Candidate,
				SdpMid:           ice.SDPMid,
				SdpMLineIndex:    int32(ice.SDPMLineIndex),
				UsernameFragment: ice.UsernameFragment,
			},
		},
	})
}

// sendReady sends ready signal to client via outputCh.
func (s *webrtcStreamer) sendReady() {
	s.pushOutput(&protos.ServerSignaling{
		SessionId: s.sessionID,
		Message:   &protos.ServerSignaling_Ready{Ready: true},
	})
}

// sendClear sends clear/interrupt signal to client via outputCh.
func (s *webrtcStreamer) sendClear() {
	s.pushOutput(&protos.ServerSignaling{
		SessionId: s.sessionID,
		Message:   &protos.ServerSignaling_Clear{Clear: true},
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
// by background goroutines. Shutdown is signalled by a ConversationDisconnection
// message through inputCh, which the Talk loop handles to trigger Disconnect().
// No context select here — Close() pushes ConversationDisconnection first and
// cancels the context afterwards, so a competing select could skip the message.
func (s *webrtcStreamer) Recv() (internal_type.Stream, error) {
	select {
	case msg, ok := <-s.inputCh:
		if !ok {
			return nil, io.EOF
		}
		return msg, nil
	case <-s.ctx.Done():
		return nil, io.EOF
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
				s.pushDisconnection(protos.ConversationDisconnection_DISCONNECTION_TYPE_USER)
			} else {
				s.pushDisconnection(protos.ConversationDisconnection_DISCONNECTION_TYPE_USER)
			}
			return
		}
		switch msg.GetRequest().(type) {
		case *protos.WebTalkRequest_Initialization:
			s.pushInput(msg.GetInitialization())
			s.handleConfigurationMessage(msg.GetInitialization().GetStreamMode())
		case *protos.WebTalkRequest_Configuration:
			s.pushInput(msg.GetConfiguration())
			s.handleConfigurationMessage(msg.GetConfiguration().GetStreamMode())
		case *protos.WebTalkRequest_Message:
			s.pushInput(msg.GetMessage())
		case *protos.WebTalkRequest_Metadata:
			s.pushInput(msg.GetMetadata())
		case *protos.WebTalkRequest_Metric:
			s.pushInput(msg.GetMetric())
		case *protos.WebTalkRequest_Disconnection:
			s.pushInput(msg.GetDisconnection())
		case *protos.WebTalkRequest_Signaling:
			s.handleClientSignaling(msg.GetSignaling())
		default:
			s.logger.Warnw("Unknown message type", "type", fmt.Sprintf("%T", msg.GetRequest()))
		}
	}
}

// pushInput sends a message to the unified input channel (non-blocking).
// recv (non-blocking) -> inputCh
// Safe to call after inputCh is closed — the send is guarded by the closed flag.
func (s *webrtcStreamer) pushInput(msg internal_type.Stream) {
	select {
	case s.inputCh <- msg:
	default:
		s.logger.Warnw("Input channel full, dropping message",
			"type", fmt.Sprintf("%T", msg))
	}
}

// pushOutput sends a message to the unified output channel (non-blocking).
// send (non-blocking) -> outputCh
func (s *webrtcStreamer) pushOutput(msg internal_type.Stream) {
	select {
	case s.outputCh <- msg:
	default:
		s.logger.Warnw("Output channel full, dropping message",
			"type", fmt.Sprintf("%T", msg))
	}
}

// handleConfigurationMessage processes transport mode changes.
// Switching text <-> audio only changes I/O transport - it does NOT create a new session.
func (s *webrtcStreamer) handleConfigurationMessage(mode protos.StreamMode) {
	s.mu.Lock()
	currentMode := s.currentMode
	s.mu.Unlock()

	if mode == currentMode {
		return
	}

	switch mode {
	case protos.StreamMode_STREAM_MODE_AUDIO:
		if err := s.setupAudioAndHandshake(); err != nil {
			s.logger.Errorw("Failed to setup audio", "error", err)
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
				s.logger.Warnw("Received SDP answer but peer connection is nil, ignoring")
				return
			}
			if err := pc.SetRemoteDescription(pionwebrtc.SessionDescription{
				Type: pionwebrtc.SDPTypeAnswer,
				SDP:  msg.Sdp.GetSdp(),
			}); err != nil {
				s.logger.Errorw("Failed to set remote description", "error", err)
			}
		}

	case *protos.ClientSignaling_IceCandidate:
		if pc == nil {
			s.logger.Warnw("Received ICE candidate but peer connection is nil, ignoring")
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
			s.logger.Errorw("Failed to add ICE candidate", "error", err)
		}

	case *protos.ClientSignaling_Disconnect:
		if msg.Disconnect {
			s.pushDisconnection(protos.ConversationDisconnection_DISCONNECTION_TYPE_USER)
		}
	}
}

// pushDisconnection pushes disconnect metrics followed by a ConversationDisconnection
// into inputCh. FIFO ordering guarantees the Talk loop processes metrics before
// the disconnection signal, eliminating the EOF race condition.
func (s *webrtcStreamer) pushDisconnection(reason protos.ConversationDisconnection_DisconnectionType) {
	s.mu.Lock()
	alreadyClosed := s.closed
	s.closed = true
	s.mu.Unlock()
	if alreadyClosed {
		return
	}

	s.pushInput(&protos.ConversationDisconnection{
		Type: reason,
		Time: timestamppb.Now(),
	})
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

// initiateWebRTCHandshake sends config and creates/sends SDP offer via outputCh.
func (s *webrtcStreamer) initiateWebRTCHandshake() error {
	s.sendConfig()

	offer, err := s.createAndSetLocalOffer()
	if err != nil {
		return err
	}

	s.sendOffer(offer.SDP)
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

// Send pushes output to the client via the unified output channel.
// All messages (audio and non-audio) flow through outputCh to preserve ordering.
// send (non-blocking) -> outputCh -> loop (runOutputWriter) -> upstream service
func (s *webrtcStreamer) Send(response internal_type.Stream) error {
	switch data := response.(type) {
	case *protos.ConversationAssistantMessage:
		switch content := data.Message.(type) {
		case *protos.ConversationAssistantMessage_Audio:
			audio48kHz, err := s.resampler.Resample(content.Audio, internal_audio.RAPIDA_INTERNAL_AUDIO_CONFIG, internal_audio.WEBRTC_AUDIO_CONFIG)
			if err != nil {
				return err
			}
			s.bufferAndSendOutput(audio48kHz)
			return nil
		case *protos.ConversationAssistantMessage_Text:
			s.pushOutput(data)
		}
	case *protos.ConversationConfiguration:
		s.pushOutput(data)
	case *protos.ConversationInitialization:
		s.pushOutput(data)
	case *protos.ConversationUserMessage:
		s.pushOutput(data)
	case *protos.ConversationInterruption:
		if data.Type == protos.ConversationInterruption_INTERRUPTION_TYPE_WORD {
			s.clearInputBuffer()
			s.clearOutputBuffer()
			s.sendClear()
		}
		s.pushOutput(data)
	case *protos.ConversationDirective:
		s.pushOutput(data)
		if data.GetType() == protos.ConversationDirective_END_CONVERSATION {
			s.pushDisconnection(protos.ConversationDisconnection_DISCONNECTION_TYPE_TOOL)
		}
	case *protos.ConversationError:
		s.pushOutput(data)
	case *protos.ConversationMetadata:
		s.pushOutput(data)
	case *protos.ConversationDisconnection:
		s.pushOutput(data)
	case *protos.ConversationMetric:
		s.pushOutput(data)
	}
	return nil
}

// ============================================================================
// Buffer helpers
// ============================================================================

func (s *webrtcStreamer) clearInputBuffer() {
	s.inputAudioBufferLock.Lock()
	s.inputAudioBuffer.Reset()
	s.inputAudioBufferLock.Unlock()
	for {
		select {
		case <-s.inputCh:
		default:
			return
		}
	}
}

func (s *webrtcStreamer) clearOutputBuffer() {
	// 1. Reset the PCM accumulation buffer so no new frames are encoded.
	s.outputAudioBufferLock.Lock()
	s.outputAudioBuffer.Reset()
	s.outputAudioBufferLock.Unlock()

	// 2. Signal the output writer to flush its local pending audio queue first,
	//    before draining outputCh, to prevent the writer from dequeuing a message
	//    between drain and signal.
	select {
	case s.flushAudioCh <- struct{}{}:
	default:
	}

	// 3. Drain the output channel (pending audio + gRPC messages).
	for {
		select {
		case <-s.outputCh:
		default:
			return
		}
	}
}

// ============================================================================
// Lifecycle
// ============================================================================

// watchCallerContext monitors the caller's context and triggers a graceful
// close when it is cancelled, ensuring cleanup is never short-circuited.
func (s *webrtcStreamer) watchCallerContext(callerCtx context.Context) {
	select {
	case <-callerCtx.Done():
		s.logger.Infow("Caller context cancelled, closing streamer gracefully", "session", s.sessionID)
		s.Close()
	case <-s.ctx.Done():
		// Streamer already closed on its own, nothing to do.
	}
}

// Close closes the WebRTC connection and releases all resources.
// It is idempotent — safe to call from multiple goroutines or multiple times.
// pushDisconnection handles the closed flag and idempotency; if it has already
// been called (e.g. from runGrpcReader or a client disconnect signal), the
// duplicate push is a no-op.
func (s *webrtcStreamer) Close() error {
	// Push disconnection signal into inputCh so the Talk loop exits cleanly.
	// pushDisconnection is idempotent (checks+sets s.closed under lock).
	s.pushDisconnection(protos.ConversationDisconnection_DISCONNECTION_TYPE_USER)

	// Tear down audio goroutines first (they depend on audioCtx).
	s.stopAudioProcessing()

	// Close the peer connection and nil out resources.
	s.mu.Lock()
	if s.pc != nil {
		s.pc.Close()
		s.pc = nil
	}
	s.localTrack = nil
	s.mu.Unlock()

	// Cancel the streamer-wide context last so that Recv() can still
	// drain inputCh before the context fires.
	s.cancel()
	return nil
}
