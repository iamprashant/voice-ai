// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package internal_silero_vad

import (
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	internal_audio "github.com/rapidaai/api/assistant-api/internal/audio"
	internal_audio_resampler "github.com/rapidaai/api/assistant-api/internal/audio/resampler"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/utils"
	"github.com/rapidaai/protos"
	"github.com/streamer45/silero-vad-go/speech"
)

// -----------------------------------------------------------------------------
// Constants
// -----------------------------------------------------------------------------

const (
	// vadName is the identifier for this VAD implementation
	vadName = "silero_vad"

	// Default configuration values
	defaultThreshold            = 0.5
	defaultMinSilenceDurationMs = 100
	defaultSpeechPadMs          = 30

	// Environment variable for model path
	envModelPathKey = "SILERO_MODEL_PATH"

	// Default model filename
	defaultModelFile = "models/silero_vad_20251001.onnx"
)

// -----------------------------------------------------------------------------
// SileroVAD - Voice Activity Detection using Silero
// -----------------------------------------------------------------------------

// SileroVAD implements the Vad interface using the silero-vad-go library.
// It provides thread-safe voice activity detection with automatic cleanup
// on context cancellation.
type SileroVAD struct {
	// Core dependencies
	logger     commons.Logger
	onActivity internal_type.VADCallback

	// Audio processing pipeline
	audioSampler   internal_type.AudioResampler
	audioConverter internal_type.AudioConverter

	// Audio configuration
	inputConfig *protos.AudioConfig // Input audio format from caller
	vadConfig   *protos.AudioConfig // Required format for VAD (16kHz mono)

	// Silero detector (CGO-backed, requires careful lifecycle management)
	detector *speech.Detector

	// Thread-safety for CGO resource protection
	mu           sync.RWMutex
	isTerminated bool
}

// -----------------------------------------------------------------------------
// Constructor
// -----------------------------------------------------------------------------

// NewSileroVAD creates a new SileroVAD instance.
// The VAD will automatically close when the provided context is cancelled,
// ensuring safe cleanup of CGO resources.
func NewSileroVAD(
	ctx context.Context,
	logger commons.Logger,
	inputAudio *protos.AudioConfig,
	callback internal_type.VADCallback,
	options utils.Option,
) (internal_type.Vad, error) {
	// Initialize detector
	detector, err := createDetector(options)
	if err != nil {
		return nil, fmt.Errorf("failed to create silero detector: %w", err)
	}

	// Initialize audio processing pipeline
	resampler, converter, err := createAudioPipeline(logger)
	if err != nil {
		detector.Destroy() // Clean up on failure
		return nil, fmt.Errorf("failed to create audio pipeline: %w", err)
	}

	svad := &SileroVAD{
		logger:         logger,
		onActivity:     callback,
		audioSampler:   resampler,
		audioConverter: converter,
		inputConfig:    inputAudio,
		vadConfig:      internal_audio.NewLinear16khzMonoAudioConfig(),
		detector:       detector,
		isTerminated:   false,
	}

	// Start lifecycle manager for automatic cleanup
	svad.startLifecycleManager(ctx)

	return svad, nil
}

// -----------------------------------------------------------------------------
// Public Interface Methods
// -----------------------------------------------------------------------------

// Name returns the identifier for this VAD implementation.
func (s *SileroVAD) Name() string {
	return vadName
}

// Process analyzes an audio packet for voice activity.
// Returns immediately if the VAD has been terminated.
// Thread-safe for concurrent calls.
func (s *SileroVAD) Process(ctx context.Context, pkt internal_type.UserAudioPacket) error {
	// Early termination check
	if !s.isActive() {
		return nil
	}

	// Prepare audio samples
	samples, err := s.prepareAudioSamples(pkt)
	if err != nil {
		return err
	}

	// Perform detection with CGO safety
	segments, err := s.detectSafely(samples)
	if err != nil {
		return err
	}

	// Notify callback if speech detected
	if len(segments) > 0 {
		s.notifyActivity(segments)
	}

	return nil
}

// Close terminates the VAD and releases all CGO resources.
// Safe to call multiple times; subsequent calls are no-ops.
// Thread-safe.
func (s *SileroVAD) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isTerminated {
		return nil
	}
	s.isTerminated = true

	if s.detector != nil {
		s.detector.Destroy()
		s.detector = nil
	}

	return nil
}

// -----------------------------------------------------------------------------
// Private Helper Methods - Initialization
// -----------------------------------------------------------------------------

// createDetector initializes the Silero speech detector with configuration.
func createDetector(options utils.Option) (*speech.Detector, error) {
	modelPath := resolveModelPath()
	threshold := resolveThreshold(options)

	config := speech.DetectorConfig{
		ModelPath:            modelPath,
		SampleRate:           16000, // Silero requires 16kHz
		Threshold:            float32(threshold),
		MinSilenceDurationMs: defaultMinSilenceDurationMs,
		SpeechPadMs:          defaultSpeechPadMs,
	}

	return speech.NewDetector(config)
}

// createAudioPipeline initializes the resampler and converter.
func createAudioPipeline(logger commons.Logger) (internal_type.AudioResampler, internal_type.AudioConverter, error) {
	resampler, err := internal_audio_resampler.GetResampler(logger)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get resampler: %w", err)
	}

	converter, err := internal_audio_resampler.GetConverter(logger)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get converter: %w", err)
	}

	return resampler, converter, nil
}

// resolveModelPath determines the ONNX model file path.
func resolveModelPath() string {
	if envPath := os.Getenv(envModelPathKey); envPath != "" {
		return envPath
	}

	_, currentFile, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(currentFile), defaultModelFile)
}

// resolveThreshold extracts threshold from options or returns default.
func resolveThreshold(options utils.Option) float64 {
	if options == nil {
		return defaultThreshold
	}

	if threshold, err := options.GetFloat64("microphone.vad.threshold"); err == nil {
		return threshold
	}

	return defaultThreshold
}

// -----------------------------------------------------------------------------
// Private Helper Methods - Lifecycle
// -----------------------------------------------------------------------------

// startLifecycleManager spawns a goroutine that closes the VAD
// when the context is cancelled.
func (s *SileroVAD) startLifecycleManager(ctx context.Context) {
	go func() {
		<-ctx.Done()
		s.Close()
	}()
}

// isActive checks if the VAD is still operational.
// Thread-safe.
func (s *SileroVAD) isActive() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return !s.isTerminated && s.detector != nil
}

// -----------------------------------------------------------------------------
// Private Helper Methods - Audio Processing
// -----------------------------------------------------------------------------

// prepareAudioSamples resamples and converts audio to the format required by Silero.
func (s *SileroVAD) prepareAudioSamples(pkt internal_type.UserAudioPacket) ([]float32, error) {
	// Resample to 16kHz mono
	resampled, err := s.audioSampler.Resample(pkt.Audio, s.inputConfig, s.vadConfig)
	if err != nil {
		s.logger.Debugf("Resampling failed: %+v", err)
		return nil, fmt.Errorf("resampling failed: %w", err)
	}

	// Convert to float32 samples
	samples, err := s.audioConverter.ConvertToFloat32Samples(resampled, s.vadConfig)
	if err != nil {
		s.logger.Debugf("Sample conversion failed: %+v", err)
		return nil, fmt.Errorf("sample conversion failed: %w", err)
	}

	return samples, nil
}

// detectSafely performs voice activity detection with CGO resource protection.
// Acquires read lock to prevent Close() from destroying detector during detection.
func (s *SileroVAD) detectSafely(samples []float32) ([]speech.Segment, error) {
	s.mu.RLock()
	if s.isTerminated || s.detector == nil {
		s.mu.RUnlock()
		return nil, nil
	}
	detector := s.detector
	s.mu.RUnlock()

	segments, err := detector.Detect(samples)
	if err != nil {
		return nil, fmt.Errorf("detection failed: %w", err)
	}

	return segments, nil
}

// notifyActivity calculates speech boundaries and invokes the callback.
func (s *SileroVAD) notifyActivity(segments []speech.Segment) {
	minStart := math.MaxFloat64
	maxEnd := -math.MaxFloat64

	for _, seg := range segments {
		start := float64(seg.SpeechStartAt)
		end := float64(seg.SpeechEndAt)

		if start < minStart {
			minStart = start
		}
		if end > maxEnd {
			maxEnd = end
		}
	}

	s.onActivity(internal_type.InterruptionPacket{
		Source:  internal_type.InterruptionSourceVad,
		StartAt: minStart,
		EndAt:   maxEnd,
	})
}
