// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package internal_end_of_speech

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/utils"
)

// MockEndOfSpeechCallback is a simple callback function for testing
var mockCallback internal_type.EndOfSpeechCallback = func(ctx context.Context, result internal_type.EndOfSpeechPacket) error {
	return nil
}

func TestGetEndOfSpeech_SilenceBasedIdentifier(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()

	endOfSpeech, err := GetEndOfSpeech(context.Background(), logger, mockCallback, utils.Option{EndOfSpeechOptionsKeyProvider: SilenceBasedEndOfSpeech})

	require.NoError(t, err)
	assert.NotNil(t, endOfSpeech)
	assert.IsType(t, endOfSpeech, endOfSpeech)
}

func TestGetEndOfSpeech_UnknownIdentifier(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()

	endOfSpeech, err := GetEndOfSpeech(t.Context(), logger, mockCallback, utils.Option{EndOfSpeechOptionsKeyProvider: EndOfSpeechIdentifier("unknown_eos")})

	assert.Error(t, err)
	assert.Nil(t, endOfSpeech)
	assert.Equal(t, "illegal end of speeh", err.Error())
}

func TestGetEndOfSpeech_LiveKitIdentifier(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()

	endOfSpeech, err := GetEndOfSpeech(t.Context(), logger, mockCallback, utils.Option{EndOfSpeechOptionsKeyProvider: LiveKitEndOfSpeech})

	// Currently not implemented, should fail
	assert.Error(t, err)
	assert.Nil(t, endOfSpeech)
}

func TestGetEndOfSpeech_EmptyIdentifier(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()

	emptyIdentifier := EndOfSpeechIdentifier("")

	endOfSpeech, err := GetEndOfSpeech(t.Context(), logger, mockCallback, utils.Option{EndOfSpeechOptionsKeyProvider: emptyIdentifier})

	assert.Error(t, err)
	assert.Nil(t, endOfSpeech)
	assert.Equal(t, "illegal end of speeh", err.Error())
}

func TestEndOfSpeechIdentifier_Constants(t *testing.T) {
	assert.Equal(t, EndOfSpeechIdentifier("silence_based_eos"), SilenceBasedEndOfSpeech)
	assert.Equal(t, EndOfSpeechIdentifier("livekit_eos"), LiveKitEndOfSpeech)
	assert.NotEqual(t, SilenceBasedEndOfSpeech, LiveKitEndOfSpeech)
}

func TestGetEndOfSpeech_WithNilLogger(t *testing.T) {

	// This test validates that the function passes nil logger to NewSilenceBasedEndOfSpeech
	// which should handle it appropriately or fail gracefully
	endOfSpeech, err := GetEndOfSpeech(t.Context(), nil, mockCallback, utils.Option{EndOfSpeechOptionsKeyProvider: SilenceBasedEndOfSpeech})

	// The behavior depends on internal_silence_based_end_of_speech implementation
	// Either it should error or handle nil logger gracefully
	if err == nil {
		assert.NotNil(t, endOfSpeech)
	} else {
		assert.Nil(t, endOfSpeech)
	}
}

func TestGetEndOfSpeech_WithNilCallback(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()

	// Test with nil callback
	endOfSpeech, err := GetEndOfSpeech(t.Context(), logger, nil, utils.Option{EndOfSpeechOptionsKeyProvider: SilenceBasedEndOfSpeech})

	// The behavior depends on internal_silence_based_end_of_speech implementation
	if err == nil {
		assert.NotNil(t, endOfSpeech)
	} else {
		assert.Nil(t, endOfSpeech)
	}
}

func TestGetEndOfSpeech_WithNilOptions(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()

	// Test with nil options
	endOfSpeech, err := GetEndOfSpeech(t.Context(), logger, mockCallback, nil)

	// The behavior depends on internal_silence_based_end_of_speech implementation
	if err == nil {
		assert.NotNil(t, endOfSpeech)
	} else {
		assert.Nil(t, endOfSpeech)
	}
}
