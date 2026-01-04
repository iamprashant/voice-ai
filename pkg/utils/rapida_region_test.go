// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package utils

import "testing"

func TestRapidaRegion_Get(t *testing.T) {
	tests := []struct {
		region   RapidaRegion
		expected string
	}{
		{AP, "ap"},
		{US, "us"},
		{EU, "eu"},
		{ALL, "all"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if result := tt.region.Get(); result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestFromRegionStr(t *testing.T) {
	tests := []struct {
		input    string
		expected RapidaRegion
	}{
		{"ap", AP},
		{"AP", AP},
		{"us", US},
		{"US", US},
		{"eu", EU},
		{"EU", EU},
		{"all", ALL},
		{"ALL", ALL},
		{"invalid", ALL}, // defaults to all
		{"", ALL},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := FromRegionStr(tt.input)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
