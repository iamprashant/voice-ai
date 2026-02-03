// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package webrtc_internal

// Opus audio constants (WebRTC standard: 48kHz)
const (
	OpusSampleRate       = 48000
	OpusFrameDuration    = 20                   // milliseconds
	OpusFrameBytes       = 1920                 // 960 samples * 2 bytes (20ms at 48kHz)
	MaxOutputBufferBytes = OpusFrameBytes * 500 // ~10 seconds, warning threshold
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
