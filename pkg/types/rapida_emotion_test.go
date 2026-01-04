// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package types

import (
	"testing"
)

func TestGetEmotionByName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Emotion
	}{
		{"angry", "angry", Emotion{"anger", "high"}},
		{"frustrated", "frustrated", Emotion{"anger", "low"}},
		{"happy", "happy", Emotion{"positivity", "high"}},
		{"content", "content", Emotion{"positivity", "low"}},
		{"sad", "sad", Emotion{"sadness", "high"}},
		{"curious", "curious", Emotion{"curiosity", ""}},
		{"neutral", "unknown", Emotion{"neutral", ""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetEmotionByName(tt.input)
			if got != tt.expected {
				t.Errorf("GetEmotionByName(%s) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}
