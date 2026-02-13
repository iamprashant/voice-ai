// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package webrtc_internal

// DisconnectReason describes why a WebRTC session was torn down.
type DisconnectReason string

const (
	DisconnectReasonNormal           DisconnectReason = "normal"            // clean end-of-conversation directive
	DisconnectReasonTool             DisconnectReason = "tool"              // tool directive with "end_conversation" signal
	DisconnectReasonClientDisconnect DisconnectReason = "client_disconnect" // client sent disconnect signal
	DisconnectReasonConnectionFailed DisconnectReason = "connection_failed" // ICE/DTLS failure
	DisconnectReasonPeerDisconnected DisconnectReason = "peer_disconnected" // transient network loss (not recovered)
	DisconnectReasonGRPCClosed       DisconnectReason = "grpc_closed"       // gRPC stream closed (EOF)
	DisconnectReasonGRPCError        DisconnectReason = "grpc_error"        // gRPC stream error
	DisconnectReasonContextCancelled DisconnectReason = "context_cancelled" // parent context cancelled
	DisconnectReasonUnknown          DisconnectReason = "unknown"
)

// Metric name constants for WebRTC streamer-level metrics.
const (
	MetricStatus           = "STATUS"
	MetricSessionDuration  = "TIME_TAKEN"
	MetricDisconnectReason = "DISCONNECT_REASON"
	MetricAudioDuration    = "AUDIO_DURATION_MS"
	MetricSessionID        = "WEBRTC_SESSION_ID"
)

// Opus audio constants (WebRTC standard: 48kHz)
const (
	OpusSampleRate    = 48000
	OpusFrameDuration = 20   // milliseconds
	OpusFrameBytes    = 1920 // 960 samples * 2 bytes (20ms at 48kHz)
	OpusChannels      = 2    // Opus RTP always signals 2 encoding channels (opus/48000/2) per RFC 7587, even for mono voice
	OpusPayloadType   = 111  // Standard dynamic payload type for Opus
	OpusSDPFmtpLine   = "minptime=10;useinbandfec=1;stereo=0;sprop-stereo=0"
)

// Channel and buffer sizes
const (
	InputChannelSize     = 500  // Buffered channel for incoming messages (~10s of 20ms audio frames)
	OutputChannelSize    = 1500 // Buffered channel for outgoing messages (~30s of 20ms audio frames)
	ErrorChannelSize     = 1    // Error channel buffer
	RTPBufferSize        = 1500 // Max RTP packet size (MTU)
	MaxConsecutiveErrors = 50   // Max read errors before stopping
	InputBufferThreshold = 3200 // 100ms at 16kHz (32 bytes/ms * 100ms) — larger chunks = fewer channel writes

	// OutputBufferThreshold triggers flushing accumulated 48kHz PCM into
	// Opus-encoded frames. Must be >= OpusFrameBytes so at least one full
	// 20ms frame can be encoded per flush.
	OutputBufferThreshold = OpusFrameBytes // flush every complete 20ms frame

	// OutputPaceInterval is the real-time interval between sending consecutive
	// audio frames to the WebRTC track. Matches OpusFrameDuration so that
	// TTS bursts are smoothed to playback rate rather than flooding the client.
	OutputPaceInterval = OpusFrameDuration // milliseconds (20ms per frame)
)

// Config holds WebRTC configuration
type Config struct {
	ICEServers         []ICEServer
	ICETransportPolicy string // "all" or "relay"
}

// ICEServer represents a STUN/TURN server
type ICEServer struct {
	URLs       []string
	Username   string
	Credential string
}

// DefaultConfig returns default WebRTC configuration
func DefaultConfig() *Config {
	return &Config{
		ICEServers: []ICEServer{
			{URLs: []string{"stun:stun.l.google.com:19302"}},
			{URLs: []string{"stun:stun1.l.google.com:19302"}},
		},
		ICETransportPolicy: "all",
	}
}

// ICECandidate represents an ICE candidate for signaling
type ICECandidate struct {
	Candidate        string
	SDPMid           string
	SDPMLineIndex    int
	UsernameFragment string
}
