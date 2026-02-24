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
	"time"

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

// chunk is a recorded audio fragment placed at a specific position on the
// timeline. ByteOffset is the byte position relative to Start().
type chunk struct {
	ByteOffset int
	Data       []byte
	Track      int // trackUser or trackSystem
}

const (
	trackUser   = 0
	trackSystem = 1
)

type audioRecorder struct {
	logger    commons.Logger
	mu        sync.Mutex
	startTime time.Time
	started   bool
	chunks    []chunk
	// Per-track cursor: the byte position just past the last written byte on
	// each track. For user track wall-clock placement is used. For system
	// (TTS) track the cursor paces audio at the playback rate — only the
	// first chunk after a gap uses wall-clock to anchor position.
	cursor [2]int
	// clock is injectable for testing; defaults to time.Now.
	clock func() time.Time
}

func NewDefaultAudioRecorder(logger commons.Logger) (internal_type.Recorder, error) {
	return &audioRecorder{
		logger: logger,
		clock:  time.Now,
	}, nil
}

// Start begins the recording session. Both tracks share this start time.
// Audio is placed on the timeline based on when it arrives relative to
// this moment.
func (r *audioRecorder) Start() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.startTime = r.clock()
	r.started = true
}

func bytesPerSecond() int {
	return int(audioConfig.SampleRate) * int(audioConfig.Channels) * AudioBytesPerSample
}

// durationBytes converts a wall-clock duration to a frame-aligned byte count.
func durationBytes(d time.Duration) int {
	raw := int(d.Seconds() * float64(bytesPerSecond()))
	frameSize := AudioBytesPerSample * int(audioConfig.Channels)
	return (raw / frameSize) * frameSize
}

// Record places audio on the appropriate track at the current wall-clock
// position. Each chunk is positioned based on WHEN it arrives, not just
// appended. Both tracks share the same timeline (Start → Persist).
func (r *audioRecorder) Record(ctx context.Context, p internal_type.Packet) error {
	switch vl := p.(type) {
	case internal_type.UserAudioPacket:
		return r.push(vl.Audio, trackUser)
	case internal_type.TextToSpeechAudioPacket:
		return r.push(vl.AudioChunk, trackSystem)
	}
	return nil
}

func (r *audioRecorder) push(data []byte, track int) error {
	if len(data) == 0 {
		return nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	// Compute wall-clock byte offset.
	wallOffset := 0
	if r.started {
		wallOffset = durationBytes(r.clock().Sub(r.startTime))
	}

	var offset int
	switch track {
	case trackUser:
		// User (mic) audio: wall-clock placement. Mic delivers at real-time
		// rate, so wall-clock offset is the correct timeline position.
		offset = wallOffset
		if r.cursor[track] > offset {
			offset = r.cursor[track]
		}

	case trackSystem:
		// TTS audio: PACING. TTS delivers audio in bursts (faster than
		// real-time). We pace it at the playback rate on the timeline:
		//
		//   - First chunk after silence (cursor <= wallOffset): anchor at
		//     wall-clock offset.
		//   - Subsequent burst chunks (cursor > wallOffset): place at cursor
		//     so audio is continuous at the playback rate with no gaps.
		//
		// This avoids "audio breaking" that happened when wall-clock was used
		// for every chunk — variable delivery timing introduced tiny gaps or
		// overlaps between TTS chunks.
		if r.cursor[track] > wallOffset {
			// Burst continuation: pace from cursor.
			offset = r.cursor[track]
		} else {
			// New TTS segment: anchor at wall-clock.
			offset = wallOffset
		}
	}

	// Copy to avoid caller mutations.
	buf := make([]byte, len(data))
	copy(buf, data)

	r.chunks = append(r.chunks, chunk{
		ByteOffset: offset,
		Data:       buf,
		Track:      track,
	})
	// Advance cursor past this chunk.
	r.cursor[track] = offset + len(buf)
	return nil
}

// Persist renders two WAV files — one per track. Both WAVs span the full
// session duration (Start → Persist). Audio chunks are placed at their
// recorded timeline positions; gaps are silence.
func (r *audioRecorder) Persist() ([]byte, []byte, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.chunks) == 0 {
		return nil, nil, fmt.Errorf("no audio chunks to persist")
	}

	// Total session duration in bytes.
	sessionBytes := 0
	if r.started {
		sessionBytes = durationBytes(r.clock().Sub(r.startTime))
	}

	// Determine the minimum buffer size: max(sessionBytes, furthest chunk end).
	totalLen := sessionBytes
	for _, c := range r.chunks {
		end := c.ByteOffset + len(c.Data)
		if end > totalLen {
			totalLen = end
		}
	}

	// Allocate zero-filled (silence) buffers for each track.
	userPCM := make([]byte, totalLen)
	systemPCM := make([]byte, totalLen)

	// Paint each chunk onto its track buffer.
	userAudioBytes := 0
	systemAudioBytes := 0
	for _, c := range r.chunks {
		var dst []byte
		if c.Track == trackUser {
			dst = userPCM
			userAudioBytes += len(c.Data)
		} else {
			dst = systemPCM
			systemAudioBytes += len(c.Data)
		}
		copy(dst[c.ByteOffset:], c.Data)
	}

	r.logger.Info(fmt.Sprintf(
		"Audio persist: userAudio=%d (%.2fs), systemAudio=%d (%.2fs), totalLen=%d (%.2fs), chunks=%d",
		userAudioBytes, float64(userAudioBytes)/float64(bytesPerSecond()),
		systemAudioBytes, float64(systemAudioBytes)/float64(bytesPerSecond()),
		totalLen, float64(totalLen)/float64(bytesPerSecond()),
		len(r.chunks),
	))

	userWAV, _ := createWAVFile(userPCM)
	systemWAV, _ := createWAVFile(systemPCM)
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
