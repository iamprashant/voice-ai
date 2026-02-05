// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package audiosocket

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

// AudioSocket audio constants (SLIN 16-bit 8kHz)
const (
	chunkDuration           = 20 * time.Millisecond
	slinBytesPerMs          = 16 // 8000 Hz * 2 bytes / 1000
	defaultOptimalFrameSize = 320 // 20ms at 8kHz 16-bit = 320 bytes
	outputChunkSize         = slinBytesPerMs * 20
	inputBufferThreshold    = 32 * 60
	slinSilence             = 0x00
)

// AudioChunk represents a processed audio chunk ready for streaming.
type AudioChunk struct {
	Data     []byte
	Duration time.Duration
}

// AudioProcessor handles audio conversion for AudioSocket (SLIN 8kHz <-> linear16 16kHz).
type AudioProcessor struct {
	logger commons.Logger

	resampler internal_type.AudioResampler

	asteriskConfig   *protos.AudioConfig
	downstreamConfig *protos.AudioConfig

	optimalFrameSize int

	inputBuffer   *bytes.Buffer
	inputBufferMu sync.Mutex

	outputBuffer   *bytes.Buffer
	outputBufferMu sync.Mutex

	onInputAudio  func(audio []byte)
	onOutputChunk func(chunk *AudioChunk) error

	silenceChunk *AudioChunk

	xoffActive bool
	xoffMu     sync.Mutex
}

// NewAudioProcessor creates a new AudioSocket audio processor.
func NewAudioProcessor(logger commons.Logger) (*AudioProcessor, error) {
	resampler, err := internal_audio_resampler.GetResampler(logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create resampler: %w", err)
	}

	p := &AudioProcessor{
		logger:           logger,
		resampler:        resampler,
		asteriskConfig:   internal_audio.NewLinear8khzMonoAudioConfig(),
		downstreamConfig: internal_audio.NewLinear16khzMonoAudioConfig(),
		optimalFrameSize: defaultOptimalFrameSize,
		inputBuffer:      new(bytes.Buffer),
		outputBuffer:     new(bytes.Buffer),
	}

	p.silenceChunk = p.createSilenceChunk()

	return p, nil
}

// SetInputAudioCallback sets the callback for processed input audio.
func (p *AudioProcessor) SetInputAudioCallback(callback func(audio []byte)) {
	p.onInputAudio = callback
}

// SetOutputChunkCallback sets the callback for sending audio chunks.
func (p *AudioProcessor) SetOutputChunkCallback(callback func(chunk *AudioChunk) error) {
	p.onOutputChunk = callback
}

// GetDownstreamConfig returns the downstream audio configuration (16kHz linear16).
func (p *AudioProcessor) GetDownstreamConfig() *protos.AudioConfig {
	return p.downstreamConfig
}

// ProcessInputAudio converts incoming SLIN 8kHz audio to linear16 16kHz.
func (p *AudioProcessor) ProcessInputAudio(audio []byte) error {
	if len(audio) == 0 {
		return nil
	}

	converted, err := p.resampler.Resample(audio, p.asteriskConfig, p.downstreamConfig)
	if err != nil {
		return fmt.Errorf("audio conversion from SLIN 8kHz to 16kHz linear16 failed: %w", err)
	}

	p.bufferAndSendInput(converted)
	return nil
}

func (p *AudioProcessor) bufferAndSendInput(audio []byte) {
	p.inputBufferMu.Lock()
	p.inputBuffer.Write(audio)

	if p.inputBuffer.Len() < inputBufferThreshold {
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

// ClearInputBuffer clears the input audio buffer.
func (p *AudioProcessor) ClearInputBuffer() {
	p.inputBufferMu.Lock()
	p.inputBuffer.Reset()
	p.inputBufferMu.Unlock()
}

// ProcessOutputAudio converts outgoing linear16 16kHz audio to SLIN 8kHz.
func (p *AudioProcessor) ProcessOutputAudio(audio []byte) error {
	if len(audio) == 0 {
		return nil
	}

	converted, err := p.resampler.Resample(audio, p.downstreamConfig, p.asteriskConfig)
	if err != nil {
		return fmt.Errorf("audio conversion to SLIN 8kHz failed: %w", err)
	}

	p.outputBufferMu.Lock()
	p.outputBuffer.Write(converted)
	p.outputBufferMu.Unlock()

	return nil
}

func (p *AudioProcessor) GetNextChunk() *AudioChunk {
	chunkSize := p.optimalFrameSize
	if chunkSize <= 0 {
		chunkSize = outputChunkSize
	}

	chunk := make([]byte, chunkSize)

	p.outputBufferMu.Lock()
	n, _ := p.outputBuffer.Read(chunk)
	p.outputBufferMu.Unlock()

	if n == 0 {
		return nil
	}

	if n < chunkSize {
		for i := n; i < chunkSize; i++ {
			chunk[i] = slinSilence
		}
	}

	return &AudioChunk{
		Data:     chunk,
		Duration: chunkDuration,
	}
}

func (p *AudioProcessor) createSilenceChunk() *AudioChunk {
	chunkSize := p.optimalFrameSize
	if chunkSize <= 0 {
		chunkSize = outputChunkSize
	}

	chunk := make([]byte, chunkSize)
	for i := range chunk {
		chunk[i] = slinSilence
	}

	return &AudioChunk{
		Data:     chunk,
		Duration: chunkDuration,
	}
}

// RunOutputSender continuously sends audio chunks at consistent intervals.
func (p *AudioProcessor) RunOutputSender(ctx context.Context) {
	if p.onOutputChunk == nil {
		p.logger.Error("RunOutputSender called without output callback set")
		return
	}

	nextSendTime := time.Now().Add(chunkDuration)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		now := time.Now()
		if sleepDuration := nextSendTime.Sub(now); sleepDuration > 0 {
			time.Sleep(sleepDuration)
		}

		nextSendTime = nextSendTime.Add(chunkDuration)
		if time.Now().After(nextSendTime) {
			nextSendTime = time.Now().Add(chunkDuration)
		}

		if p.IsXOFF() {
			continue
		}

		chunk := p.GetNextChunk()
		if chunk == nil {
			chunk = p.silenceChunk
		}

		if err := p.onOutputChunk(chunk); err != nil {
			p.logger.Debug("Failed to send audio chunk", "error", err)
		}
	}
}

// ClearOutputBuffer clears the output audio buffer.
func (p *AudioProcessor) ClearOutputBuffer() {
	p.outputBufferMu.Lock()
	p.outputBuffer.Reset()
	p.outputBufferMu.Unlock()
}

// SetXOFF pauses audio output.
func (p *AudioProcessor) SetXOFF() {
	p.xoffMu.Lock()
	p.xoffActive = true
	p.xoffMu.Unlock()
}

// SetXON resumes audio output.
func (p *AudioProcessor) SetXON() {
	p.xoffMu.Lock()
	p.xoffActive = false
	p.xoffMu.Unlock()
}

// IsXOFF returns whether output is paused.
func (p *AudioProcessor) IsXOFF() bool {
	p.xoffMu.Lock()
	defer p.xoffMu.Unlock()
	return p.xoffActive
}
