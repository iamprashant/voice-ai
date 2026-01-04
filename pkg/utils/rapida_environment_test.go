// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package utils

import "testing"

func TestRapidaEnvironment_Get(t *testing.T) {
	tests := []struct {
		env      RapidaEnvironment
		expected string
	}{
		{PRODUCTION, "production"},
		{DEVELOPMENT, "development"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if result := tt.env.Get(); result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestFromEnvironmentStr(t *testing.T) {
	tests := []struct {
		input    string
		expected RapidaEnvironment
	}{
		{"production", PRODUCTION},
		{"PRODUCTION", PRODUCTION},
		{"development", DEVELOPMENT},
		{"DEVELOPMENT", DEVELOPMENT},
		{"invalid", DEVELOPMENT}, // defaults to development
		{"", DEVELOPMENT},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := FromEnvironmentStr(tt.input)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
