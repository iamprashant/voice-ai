// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package internal_twilio

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

// Twilio audio constants (mulaw 8kHz)
const (
	// Standard chunk duration for telephony (20ms)
	ChunkDuration = 20 * time.Millisecond

	// Mulaw 8kHz: 8 bytes per ms (8-bit mono, 8000 samples/sec)
	MulawBytesPerMs = 8

	// Output chunk size: 20ms at 8kHz mulaw = 160 bytes
	OutputChunkSize = MulawBytesPerMs * 20

	// Input buffer threshold: 60ms at 16kHz linear16 = 1920 bytes
	InputBufferThreshold = 32 * 60

	// Mulaw silence value (0x7F or 0xFF represents silence)
	MulawSilence = 0xFF
)

// AudioChunk represents a processed audio chunk ready for streaming
type AudioChunk struct {
	Data     []byte
	Duration time.Duration
}

// AudioProcessor handles audio conversion for Twilio (mulaw 8kHz <-> linear16 16kHz)
type AudioProcessor struct {
	logger commons.Logger

	// Resampler for format and sample rate conversion
	resampler internal_type.AudioResampler

	// Audio configs
	twilioConfig     *protos.AudioConfig // mulaw 8kHz for Twilio
	downstreamConfig *protos.AudioConfig // linear16 16kHz for STT/TTS

	// Input buffer for accumulating incoming audio (converted to 16kHz)
	inputBuffer   *bytes.Buffer
	inputBufferMu sync.Mutex

	// Output buffer for audio to be sent to Twilio (converted to mulaw 8kHz)
	outputBuffer   *bytes.Buffer
	outputBufferMu sync.Mutex

	// Callback for processed input audio (to send to downstream)
	onInputAudio func(audio []byte)

	// Callback for sending audio chunk to Twilio
	onOutputChunk func(chunk *AudioChunk) error

	// Pre-created silence chunk
	silenceChunk *AudioChunk
}

// NewAudioProcessor creates a new Twilio audio processor
func NewAudioProcessor(logger commons.Logger) (*AudioProcessor, error) {
	resampler, err := internal_audio_resampler.GetResampler(logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create resampler: %w", err)
	}

	p := &AudioProcessor{
		logger:           logger,
		resampler:        resampler,
		twilioConfig:     internal_audio.NewMulaw8khzMonoAudioConfig(),
		downstreamConfig: internal_audio.NewLinear16khzMonoAudioConfig(),
		inputBuffer:      new(bytes.Buffer),
		outputBuffer:     new(bytes.Buffer),
	}

	// Pre-create silence chunk
	p.silenceChunk = p.createSilenceChunk()

	return p, nil
}

// SetInputAudioCallback sets the callback for processed input audio
func (p *AudioProcessor) SetInputAudioCallback(callback func(audio []byte)) {
	p.onInputAudio = callback
}

// SetOutputChunkCallback sets the callback for sending audio chunks to Twilio
func (p *AudioProcessor) SetOutputChunkCallback(callback func(chunk *AudioChunk) error) {
	p.onOutputChunk = callback
}

// GetDownstreamConfig returns the downstream audio configuration (16kHz linear16)
func (p *AudioProcessor) GetDownstreamConfig() *protos.AudioConfig {
	return p.downstreamConfig
}

// ============================================================================
// Input Audio Processing (from Twilio mulaw 8kHz -> downstream linear16 16kHz)
// ============================================================================

// ProcessInputAudio converts incoming mulaw 8kHz audio to linear16 16kHz
func (p *AudioProcessor) ProcessInputAudio(audio []byte) error {
	if len(audio) == 0 {
		return nil
	}

	// Convert from mulaw 8kHz to linear16 16kHz
	converted, err := p.resampler.Resample(audio, p.twilioConfig, p.downstreamConfig)
	if err != nil {
		return fmt.Errorf("audio conversion to 16kHz linear16 failed: %w", err)
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
// Output Audio Processing (from downstream linear16 16kHz -> Twilio mulaw 8kHz)
// ============================================================================

// ProcessOutputAudio converts outgoing linear16 16kHz audio to mulaw 8kHz
func (p *AudioProcessor) ProcessOutputAudio(audio []byte) error {
	if len(audio) == 0 {
		return nil
	}

	// Convert from linear16 16kHz to mulaw 8kHz
	converted, err := p.resampler.Resample(audio, p.downstreamConfig, p.twilioConfig)
	if err != nil {
		return fmt.Errorf("audio conversion to mulaw 8kHz failed: %w", err)
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

	// Pad with mulaw silence if chunk is not full
	if n < OutputChunkSize {
		for i := n; i < OutputChunkSize; i++ {
			chunk[i] = MulawSilence
		}
	}

	return &AudioChunk{
		Data:     chunk,
		Duration: ChunkDuration,
	}
}

// createSilenceChunk creates a mulaw silence chunk
func (p *AudioProcessor) createSilenceChunk() *AudioChunk {
	chunk := make([]byte, OutputChunkSize)
	for i := range chunk {
		chunk[i] = MulawSilence
	}
	return &AudioChunk{
		Data:     chunk,
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
