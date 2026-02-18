// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package internal_recorder

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"sync"

	internal_audio "github.com/rapidaai/api/assistant-api/internal/audio"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	"github.com/rapidaai/pkg/commons"
)

const (
	AudioBytesPerSample = 2  // LINEAR16 → 2 bytes per sample
	AudioBitsPerSample  = 16 // LINEAR16 → 16 bits per sample
	AudioPCMFormat      = 1  // WAV PCM format tag
)

var audioConfig = internal_audio.RAPIDA_INTERNAL_AUDIO_CONFIG

type trackKind int

const (
	trackUser trackKind = iota
	trackSystem
)

type chunk struct {
	Data  []byte
	Track trackKind
}

type audioRecorder struct {
	logger commons.Logger
	mu     sync.Mutex
	chunks []chunk
}

func NewDefaultAudioRecorder(logger commons.Logger) (internal_type.Recorder, error) {
	return &audioRecorder{logger: logger}, nil
}

// Record appends audio data. Mutex-guarded so safe to call from goroutines.
// TTS delivers chunks sequentially, so arrival order is preserved.
func (r *audioRecorder) Record(ctx context.Context, p internal_type.Packet) error {
	switch vl := p.(type) {
	case internal_type.UserAudioPacket:
		return r.push(vl.Audio, trackUser)
	case internal_type.TextToSpeechAudioPacket:
		return r.push(vl.AudioChunk, trackSystem)
	}
	return nil
}

func (r *audioRecorder) push(data []byte, track trackKind) error {
	if len(data) == 0 {
		return nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	buf := make([]byte, len(data))
	copy(buf, data)
	r.chunks = append(r.chunks, chunk{Data: buf, Track: track})
	return nil
}

func bytesPerSecond() int {
	return int(audioConfig.SampleRate) * int(audioConfig.Channels) * AudioBytesPerSample
}

// Persist renders two time-aligned WAV files (user + system).
// Consecutive same-track chunks are merged; the other track gets silence.
func (r *audioRecorder) Persist() ([]byte, []byte, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.chunks) == 0 {
		return nil, nil, fmt.Errorf("no audio chunks to persist")
	}

	// Build two PCM buffers in one pass over chunks.
	var userPCM, systemPCM bytes.Buffer
	for _, c := range r.chunks {
		silence := make([]byte, len(c.Data))
		switch c.Track {
		case trackUser:
			userPCM.Write(c.Data)
			systemPCM.Write(silence)
		case trackSystem:
			userPCM.Write(silence)
			systemPCM.Write(c.Data)
		}
	}

	r.logger.Info(fmt.Sprintf("Audio timeline: totalBytes=%d, duration=%.2fs, chunks=%d",
		userPCM.Len(), float64(userPCM.Len())/float64(bytesPerSecond()), len(r.chunks)))

	userWAV, _ := createWAVFile(userPCM.Bytes())
	systemWAV, _ := createWAVFile(systemPCM.Bytes())
	return userWAV, systemWAV, nil
}

func createWAVFile(pcmData []byte) ([]byte, error) {
	var buf bytes.Buffer
	sampleRate := audioConfig.SampleRate
	channels := audioConfig.Channels
	bps := int(sampleRate) * int(channels) * AudioBytesPerSample

	buf.Write([]byte("RIFF"))
	binary.Write(&buf, binary.LittleEndian, uint32(36+len(pcmData)))
	buf.Write([]byte("WAVE"))

	buf.Write([]byte("fmt "))
	binary.Write(&buf, binary.LittleEndian, uint32(16))
	binary.Write(&buf, binary.LittleEndian, uint16(AudioPCMFormat))
	binary.Write(&buf, binary.LittleEndian, uint16(channels))
	binary.Write(&buf, binary.LittleEndian, uint32(sampleRate))
	binary.Write(&buf, binary.LittleEndian, uint32(bps))
	binary.Write(&buf, binary.LittleEndian, uint16(AudioBytesPerSample))
	binary.Write(&buf, binary.LittleEndian, uint16(AudioBitsPerSample))

	// data chunk
	buf.Write([]byte("data"))
	binary.Write(&buf, binary.LittleEndian, uint32(len(pcmData)))
	buf.Write(pcmData)

	return buf.Bytes(), nil
}
