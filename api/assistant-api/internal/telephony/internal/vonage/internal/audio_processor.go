// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package internal_vonage

import (
	"bytes"
	"context"
	"sync"
	"time"

	internal_audio "github.com/rapidaai/api/assistant-api/internal/audio"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/protos"
)

// Vonage audio constants (linear16 16kHz - same as downstream)
const (
	// Standard chunk duration for telephony (20ms)
	ChunkDuration = 20 * time.Millisecond

	// Linear16 16kHz: 32 bytes per ms (16-bit mono, 16000 samples/sec)
	Linear16BytesPerMs = 32

	// Output chunk size: 20ms at 16kHz linear16 = 640 bytes
	OutputChunkSize = Linear16BytesPerMs * 20

	// Input buffer threshold: 60ms at 16kHz linear16 = 1920 bytes
	InputBufferThreshold = Linear16BytesPerMs * 60
)

// AudioChunk represents a processed audio chunk ready for streaming
type AudioChunk struct {
	Data     []byte
	Duration time.Duration
}

// AudioProcessor handles audio for Vonage (linear16 16kHz - no conversion needed)
type AudioProcessor struct {
	logger commons.Logger

	// Audio config (same for Vonage and downstream)
	audioConfig *protos.AudioConfig // linear16 16kHz

	// Input buffer for accumulating incoming audio
	inputBuffer   *bytes.Buffer
	inputBufferMu sync.Mutex

	// Output buffer for audio to be sent to Vonage
	outputBuffer   *bytes.Buffer
	outputBufferMu sync.Mutex

	// Callback for processed input audio (to send to downstream)
	onInputAudio func(audio []byte)

	// Callback for sending audio chunk to Vonage
	onOutputChunk func(chunk *AudioChunk) error

	// Pre-created silence chunk
	silenceChunk *AudioChunk
}

// NewAudioProcessor creates a new Vonage audio processor
func NewAudioProcessor(logger commons.Logger) (*AudioProcessor, error) {
	p := &AudioProcessor{
		logger:       logger,
		audioConfig:  internal_audio.NewLinear16khzMonoAudioConfig(),
		inputBuffer:  new(bytes.Buffer),
		outputBuffer: new(bytes.Buffer),
	}

	// Pre-create silence chunk (all zeros for linear16)
	p.silenceChunk = p.createSilenceChunk()

	return p, nil
}

// SetInputAudioCallback sets the callback for processed input audio
func (p *AudioProcessor) SetInputAudioCallback(callback func(audio []byte)) {
	p.onInputAudio = callback
}

// SetOutputChunkCallback sets the callback for sending audio chunks to Vonage
func (p *AudioProcessor) SetOutputChunkCallback(callback func(chunk *AudioChunk) error) {
	p.onOutputChunk = callback
}

// GetDownstreamConfig returns the downstream audio configuration (16kHz linear16)
func (p *AudioProcessor) GetDownstreamConfig() *protos.AudioConfig {
	return p.audioConfig
}

// ============================================================================
// Input Audio Processing (from Vonage linear16 16kHz -> downstream - no conversion)
// ============================================================================

// ProcessInputAudio buffers incoming linear16 16kHz audio (no conversion needed)
func (p *AudioProcessor) ProcessInputAudio(audio []byte) error {
	if len(audio) == 0 {
		return nil
	}

	// No conversion needed - Vonage uses same format as downstream
	p.bufferAndSendInput(audio)
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
// Output Audio Processing (from downstream linear16 16kHz -> Vonage - no conversion)
// ============================================================================

// ProcessOutputAudio buffers outgoing linear16 16kHz audio (no conversion needed)
func (p *AudioProcessor) ProcessOutputAudio(audio []byte) error {
	if len(audio) == 0 {
		return nil
	}

	// No conversion needed - downstream uses same format as Vonage
	p.outputBufferMu.Lock()
	p.outputBuffer.Write(audio)
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
		Data:     chunk[:n],
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
