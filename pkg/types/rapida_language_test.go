// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package types

import (
	"testing"
)

func TestGetLanguageByName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Language
	}{
		{
			name:  "english",
			input: "en",
			expected: Language{
				Name:     "English",
				ISO639_1: "en",
				ISO639_2: "eng",
			},
		},
		{
			name:  "french",
			input: "fr",
			expected: Language{
				Name:     "French",
				ISO639_1: "fr",
				ISO639_2: "fra",
			},
		},
		{
			name:  "unknown",
			input: "unknown",
			expected: Language{
				Name:     "English",
				ISO639_1: "en",
				ISO639_2: "eng",
			},
		},
		{
			name:  "case insensitive",
			input: "FR",
			expected: Language{
				Name:     "French",
				ISO639_1: "fr",
				ISO639_2: "fra",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetLanguageByName(tt.input)
			if got != tt.expected {
				t.Errorf("GetLanguageByName() = %v, want %v", got, tt.expected)
			}
		})
	}
}