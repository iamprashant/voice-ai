// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package utils

import (
	"encoding/json"
	"testing"
)

func TestRapidaSource_Get(t *testing.T) {
	tests := []struct {
		source   RapidaSource
		expected string
	}{
		{WebPlugin, "web-plugin"},
		{Debugger, "debugger"},
		{SDK, "sdk"},
		{PhoneCall, "phone-call"},
		{Whatsapp, "whatsapp"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if result := tt.source.Get(); result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestFromSourceStr(t *testing.T) {
	tests := []struct {
		input    string
		expected RapidaSource
	}{
		{"web-plugin", WebPlugin},
		{"WEB-PLUGIN", WebPlugin},
		{"debugger", Debugger},
		{"DEBUGGER", Debugger},
		{"sdk", SDK},
		{"SDK", SDK},
		{"phone-call", PhoneCall},
		{"PHONE-CALL", PhoneCall},
		{"whatsapp", Whatsapp},
		{"WHATSAPP", Whatsapp},
		{"invalid", Debugger}, // defaults to debugger
		{"", Debugger},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := FromSourceStr(tt.input)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestRapidaSource_MarshalJSON(t *testing.T) {
	source := WebPlugin
	data, err := source.MarshalJSON()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var unmarshaled string
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("error unmarshaling: %v", err)
	}

	if unmarshaled != "web-plugin" {
		t.Errorf("expected 'web-plugin', got %s", unmarshaled)
	}
}
