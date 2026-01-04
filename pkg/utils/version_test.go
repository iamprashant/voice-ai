// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package utils

import "testing"

func TestGetVersionDefinition(t *testing.T) {
	tests := []struct {
		input    string
		expected *uint64
	}{
		{"", nil},
		{"latest", nil},
		{"vrsn_123", func() *uint64 { v := uint64(123); return &v }()},
		{"vrsn_0", func() *uint64 { v := uint64(0); return &v }()},
		{"invalid", nil},
		{"vrsn_", nil},
		{"vrsn_abc", nil},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := GetVersionDefinition(tt.input)
			if tt.expected == nil {
				if result != nil {
					t.Errorf("expected nil, got %v", result)
				}
			} else {
				if result == nil || *result != *tt.expected {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}
