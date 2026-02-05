// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package internal_asterisk

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

// Asterisk audio constants (ulaw 8kHz)
const (
	// Standard chunk duration for telephony (20ms)
	ChunkDuration = 20 * time.Millisecond

	// Ulaw 8kHz: 8 bytes per ms (8-bit mono, 8000 samples/sec)
	UlawBytesPerMs = 8

	// Default optimal frame size from Asterisk (typically 160 bytes = 20ms)
	DefaultOptimalFrameSize = 160

	// Output chunk size: 20ms at 8kHz ulaw = 160 bytes
	OutputChunkSize = UlawBytesPerMs * 20

	// Input buffer threshold: 60ms at 16kHz linear16 = 1920 bytes
	InputBufferThreshold = 32 * 60

	// Ulaw silence value (0xFF represents silence in ulaw)
	UlawSilence = 0xFF
)

// AudioChunk represents a processed audio chunk ready for streaming
type AudioChunk struct {
	Data     []byte
	Duration time.Duration
}

// AudioProcessor handles audio conversion for Asterisk (ulaw 8kHz <-> linear16 16kHz)
type AudioProcessor struct {
	logger commons.Logger

	// Resampler for format and sample rate conversion
	resampler internal_type.AudioResampler

	// Audio configs
	asteriskConfig   *protos.AudioConfig // ulaw 8kHz for Asterisk
	downstreamConfig *protos.AudioConfig // linear16 16kHz for STT/TTS

	// Optimal frame size from Asterisk
	optimalFrameSize int

	// Input buffer for accumulating incoming audio (converted to 16kHz)
	inputBuffer   *bytes.Buffer
	inputBufferMu sync.Mutex

	// Output buffer for audio to be sent to Asterisk (converted to ulaw 8kHz)
	outputBuffer   *bytes.Buffer
	outputBufferMu sync.Mutex

	// Callback for processed input audio (to send to downstream)
	onInputAudio func(audio []byte)

	// Callback for sending audio chunk to Asterisk
	onOutputChunk func(chunk *AudioChunk) error

	// Pre-created silence chunk
	silenceChunk *AudioChunk

	// Flow control
	xoffActive bool
	xoffMu     sync.Mutex
}

// NewAudioProcessor creates a new Asterisk audio processor
func NewAudioProcessor(logger commons.Logger) (*AudioProcessor, error) {
	resampler, err := internal_audio_resampler.GetResampler(logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create resampler: %w", err)
	}

	p := &AudioProcessor{
		logger:           logger,
		resampler:        resampler,
		asteriskConfig:   internal_audio.NewMulaw8khzMonoAudioConfig(),
		downstreamConfig: internal_audio.NewLinear16khzMonoAudioConfig(),
		optimalFrameSize: DefaultOptimalFrameSize,
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

// SetOutputChunkCallback sets the callback for sending audio chunks to Asterisk
func (p *AudioProcessor) SetOutputChunkCallback(callback func(chunk *AudioChunk) error) {
	p.onOutputChunk = callback
}

// GetDownstreamConfig returns the downstream audio configuration (16kHz linear16)
func (p *AudioProcessor) GetDownstreamConfig() *protos.AudioConfig {
	return p.downstreamConfig
}

// SetOptimalFrameSize sets the optimal frame size received from Asterisk
func (p *AudioProcessor) SetOptimalFrameSize(size int) {
	if size > 0 {
		p.optimalFrameSize = size
	}
}

// GetOptimalFrameSize returns the current optimal frame size
func (p *AudioProcessor) GetOptimalFrameSize() int {
	return p.optimalFrameSize
}

// ============================================================================
// Flow Control (XOFF/XON)
// ============================================================================

// SetXOFF pauses audio output (flow control)
func (p *AudioProcessor) SetXOFF() {
	p.xoffMu.Lock()
	p.xoffActive = true
	p.xoffMu.Unlock()
}

// SetXON resumes audio output (flow control)
func (p *AudioProcessor) SetXON() {
	p.xoffMu.Lock()
	p.xoffActive = false
	p.xoffMu.Unlock()
}

// IsXOFF returns whether output is paused
func (p *AudioProcessor) IsXOFF() bool {
	p.xoffMu.Lock()
	defer p.xoffMu.Unlock()
	return p.xoffActive
}

// ============================================================================
// Input Audio Processing (from Asterisk ulaw 8kHz -> downstream linear16 16kHz)
// ============================================================================

// ProcessInputAudio converts incoming ulaw 8kHz audio to linear16 16kHz
func (p *AudioProcessor) ProcessInputAudio(audio []byte) error {
	if len(audio) == 0 {
		return nil
	}

	// Convert from ulaw 8kHz to linear16 16kHz
	converted, err := p.resampler.Resample(audio, p.asteriskConfig, p.downstreamConfig)
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
// Output Audio Processing (from downstream linear16 16kHz -> Asterisk ulaw 8kHz)
// ============================================================================

// ProcessOutputAudio converts outgoing linear16 16kHz audio to ulaw 8kHz
func (p *AudioProcessor) ProcessOutputAudio(audio []byte) error {
	if len(audio) == 0 {
		return nil
	}

	// Convert from linear16 16kHz to ulaw 8kHz
	converted, err := p.resampler.Resample(audio, p.downstreamConfig, p.asteriskConfig)
	if err != nil {
		return fmt.Errorf("audio conversion to ulaw 8kHz failed: %w", err)
	}

	p.outputBufferMu.Lock()
	p.outputBuffer.Write(converted)
	p.outputBufferMu.Unlock()

	return nil
}

// GetNextChunk retrieves the next audio chunk from the output buffer
func (p *AudioProcessor) GetNextChunk() *AudioChunk {
	chunkSize := p.optimalFrameSize
	if chunkSize <= 0 {
		chunkSize = OutputChunkSize
	}

	chunk := make([]byte, chunkSize)

	p.outputBufferMu.Lock()
	n, _ := p.outputBuffer.Read(chunk)
	p.outputBufferMu.Unlock()

	if n == 0 {
		return nil
	}

	// Pad with ulaw silence if chunk is not full
	if n < chunkSize {
		for i := n; i < chunkSize; i++ {
			chunk[i] = UlawSilence
		}
	}

	return &AudioChunk{
		Data:     chunk,
		Duration: ChunkDuration,
	}
}

// createSilenceChunk creates a ulaw silence chunk
func (p *AudioProcessor) createSilenceChunk() *AudioChunk {
	chunkSize := p.optimalFrameSize
	if chunkSize <= 0 {
		chunkSize = OutputChunkSize
	}

	chunk := make([]byte, chunkSize)
	for i := range chunk {
		chunk[i] = UlawSilence
	}

	return &AudioChunk{
		Data:     chunk,
		Duration: ChunkDuration,
	}
}

// RunOutputSender continuously sends audio chunks at consistent intervals
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

		// Check flow control
		if p.IsXOFF() {
			continue
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
