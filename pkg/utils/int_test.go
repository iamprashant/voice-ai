// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package utils

import "testing"

func TestMaxUint64(t *testing.T) {
	tests := []struct {
		name     string
		a, b     uint64
		expected uint64
	}{
		{"a > b", 10, 5, 10},
		{"a < b", 5, 10, 10},
		{"equal", 5, 5, 5},
		{"zero", 0, 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaxUint64(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestMinUint64(t *testing.T) {
	tests := []struct {
		name     string
		a, b     uint64
		expected uint64
	}{
		{"a > b", 10, 5, 5},
		{"a < b", 5, 10, 5},
		{"equal", 5, 5, 5},
		{"zero", 0, 1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MinUint64(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}
