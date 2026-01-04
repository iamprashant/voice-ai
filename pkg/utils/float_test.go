// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package utils

import "testing"

func TestAverageFloat32(t *testing.T) {
	tests := []struct {
		name     string
		input    []float32
		expected float32
	}{
		{"normal case", []float32{1.0, 2.0, 3.0}, 2.0},
		{"single element", []float32{5.0}, 5.0},
		{"empty slice", []float32{}, 0.0},
		{"negative numbers", []float32{-1.0, 1.0}, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AverageFloat32(tt.input)
			if result != tt.expected {
				t.Errorf("expected %f, got %f", tt.expected, result)
			}
		})
	}
}
