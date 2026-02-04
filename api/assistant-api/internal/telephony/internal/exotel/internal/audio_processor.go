// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package internal_exotel

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	internal_audio "github.com/rapidaai/api/assistant-api/internal/audio"
	internal_audio_resampler "github.com/rapidaai/api/assistant-api/internal/audio/resampler"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/protos"
)

// Exotel audio constants (linear16 8kHz)
const (
	// Standard chunk duration for telephony (20ms)
	ChunkDuration = 20 * time.Millisecond

	// Linear16 8kHz: 16 bytes per ms (16-bit mono, 8000 samples/sec)
	Linear8kHzBytesPerMs = 16

	// Output chunk size: 20ms at 8kHz linear16 = 320 bytes
	OutputChunkSize = Linear8kHzBytesPerMs * 20

	// Input buffer threshold: 60ms at 16kHz linear16 = 1920 bytes
	InputBufferThreshold = 32 * 60
)

// AudioChunk represents a processed audio chunk ready for streaming
type AudioChunk struct {
	Data     []byte
	Duration time.Duration
}

// AudioProcessor handles audio conversion for Exotel (linear16 8kHz <-> linear16 16kHz)
type AudioProcessor struct {
	logger commons.Logger

	// Resampler for sample rate conversion
	resampler internal_type.AudioResampler

	// Audio configs
	exotelConfig     *protos.AudioConfig // linear16 8kHz for Exotel
	downstreamConfig *protos.AudioConfig // linear16 16kHz for STT/TTS

	// Input buffer for accumulating incoming audio (converted to 16kHz)
	inputBuffer   *bytes.Buffer
	inputBufferMu sync.Mutex

	// Output buffer for audio to be sent to Exotel (converted to 8kHz)
	outputBuffer   *bytes.Buffer
	outputBufferMu sync.Mutex

	// Callback for processed input audio (to send to downstream)
	onInputAudio func(audio []byte)

	// Callback for sending audio chunk to Exotel
	onOutputChunk func(chunk *AudioChunk) error

	// Pre-created silence chunk
	silenceChunk *AudioChunk
}

// NewAudioProcessor creates a new Exotel audio processor
func NewAudioProcessor(logger commons.Logger) (*AudioProcessor, error) {
	resampler, err := internal_audio_resampler.GetResampler(logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create resampler: %w", err)
	}

	p := &AudioProcessor{
		logger:           logger,
		resampler:        resampler,
		exotelConfig:     internal_audio.NewLinear8khzMonoAudioConfig(),
		downstreamConfig: internal_audio.NewLinear16khzMonoAudioConfig(),
		inputBuffer:      new(bytes.Buffer),
		outputBuffer:     new(bytes.Buffer),
	}

	// Pre-create silence chunk (all zeros for linear16)
	p.silenceChunk = p.createSilenceChunk()

	return p, nil
}

// SetInputAudioCallback sets the callback for processed input audio
func (p *AudioProcessor) SetInputAudioCallback(callback func(audio []byte)) {
	p.onInputAudio = callback
}

// SetOutputChunkCallback sets the callback for sending audio chunks to Exotel
func (p *AudioProcessor) SetOutputChunkCallback(callback func(chunk *AudioChunk) error) {
	p.onOutputChunk = callback
}

// GetDownstreamConfig returns the downstream audio configuration (16kHz linear16)
func (p *AudioProcessor) GetDownstreamConfig() *protos.AudioConfig {
	return p.downstreamConfig
}

// ============================================================================
// Input Audio Processing (from Exotel linear16 8kHz -> downstream linear16 16kHz)
// ============================================================================

// ProcessInputAudio converts incoming linear16 8kHz audio to linear16 16kHz
func (p *AudioProcessor) ProcessInputAudio(audio []byte) error {
	if len(audio) == 0 {
		return nil
	}

	// Convert from linear16 8kHz to linear16 16kHz
	converted, err := p.resampler.Resample(audio, p.exotelConfig, p.downstreamConfig)
	if err != nil {
		return fmt.Errorf("audio conversion to 16kHz failed: %w", err)
	}

	// Buffer and send when threshold reached
	p.bufferAndSendInput(converted)
	return nil
}

// bufferAndSendInput buffers input audio and sends when threshold is reached
func (p *AudioProcessor) bufferAndSendInput(audio []byte) {
	p.inputBufferMu.Lock()
	p.inputBuffer.Write(audio)

	if p.inputBuffer.Len() < InputBufferThreshold {
		p.inputBufferMu.Unlock()
		return
	}

	audioData := make([]byte, p.inputBuffer.Len())
	p.inputBuffer.Read(audioData)
	p.inputBufferMu.Unlock()

	if p.onInputAudio != nil {
		p.onInputAudio(audioData)
	}
}

// ClearInputBuffer clears the input audio buffer
func (p *AudioProcessor) ClearInputBuffer() {
	p.inputBufferMu.Lock()
	p.inputBuffer.Reset()
	p.inputBufferMu.Unlock()
}

// ============================================================================
// Output Audio Processing (from downstream linear16 16kHz -> Exotel linear16 8kHz)
// ============================================================================

// ProcessOutputAudio converts outgoing linear16 16kHz audio to linear16 8kHz
func (p *AudioProcessor) ProcessOutputAudio(audio []byte) error {
	if len(audio) == 0 {
		return nil
	}

	// Convert from linear16 16kHz to linear16 8kHz
	converted, err := p.resampler.Resample(audio, p.downstreamConfig, p.exotelConfig)
	if err != nil {
		return fmt.Errorf("audio conversion to 8kHz failed: %w", err)
	}

	p.outputBufferMu.Lock()
	p.outputBuffer.Write(converted)
	p.outputBufferMu.Unlock()

	return nil
}

// GetNextChunk retrieves the next audio chunk from the output buffer
func (p *AudioProcessor) GetNextChunk() *AudioChunk {
	chunk := make([]byte, OutputChunkSize)

	p.outputBufferMu.Lock()
	n, _ := p.outputBuffer.Read(chunk)
	p.outputBufferMu.Unlock()

	if n == 0 {
		return nil
	}

	// Pad with silence (zeros) if chunk is not full
	// Linear16 silence is already zeros, no need to fill

	return &AudioChunk{
		Data:     chunk,
		Duration: ChunkDuration,
	}
}

// createSilenceChunk creates a linear16 silence chunk (all zeros)
func (p *AudioProcessor) createSilenceChunk() *AudioChunk {
	return &AudioChunk{
		Data:     make([]byte, OutputChunkSize),
		Duration: ChunkDuration,
	}
}

// RunOutputSender continuously sends audio chunks at consistent 20ms intervals
func (p *AudioProcessor) RunOutputSender(ctx context.Context) {
	if p.onOutputChunk == nil {
		p.logger.Error("RunOutputSender called without output callback set")
		return
	}

	nextSendTime := time.Now().Add(ChunkDuration)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		// Wait until next send time with precision
		now := time.Now()
		if sleepDuration := nextSendTime.Sub(now); sleepDuration > 0 {
			time.Sleep(sleepDuration)
		}

		// Schedule next send immediately to minimize drift
		nextSendTime = nextSendTime.Add(ChunkDuration)

		// If we've fallen behind, reset timing
		if time.Now().After(nextSendTime) {
			nextSendTime = time.Now().Add(ChunkDuration)
		}

		// Get audio chunk or use silence
		chunk := p.GetNextChunk()
		if chunk == nil {
			chunk = p.silenceChunk
		}

		// Send chunk via callback
		if err := p.onOutputChunk(chunk); err != nil {
			p.logger.Debug("Failed to send audio chunk", "error", err)
		}
	}
}

// ClearOutputBuffer clears the output audio buffer
func (p *AudioProcessor) ClearOutputBuffer() {
	p.outputBufferMu.Lock()
	p.outputBuffer.Reset()
	p.outputBufferMu.Unlock()
}
