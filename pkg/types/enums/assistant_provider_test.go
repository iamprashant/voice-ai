// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package type_enums

import "testing"

func TestToAssistantProvider(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected AssistantProvider
	}{
		{"AGENTKIT", "AGENTKIT", AGENTKIT},
		{"WEBSOCKET", "WEBSOCKET", WEBSOCKET},
		{"MODEL", "MODEL", MODEL},
		{"default", "unknown", MODEL},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToAssistantProvider(tt.input)
			if result != tt.expected {
				t.Errorf("ToAssistantProvider(%s) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestAssistantProvider_String(t *testing.T) {
	tests := []struct {
		provider AssistantProvider
		expected string
	}{
		{AGENTKIT, "AGENTKIT"},
		{WEBSOCKET, "WEBSOCKET"},
		{MODEL, "MODEL"},
	}
	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.provider.String(); got != tt.expected {
				t.Errorf("String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAssistantProvider_MarshalJSON(t *testing.T) {
	tests := []struct {
		provider AssistantProvider
		expected string
	}{
		{AGENTKIT, `"AGENTKIT"`},
		{WEBSOCKET, `"WEBSOCKET"`},
		{MODEL, `"MODEL"`},
	}
	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			got, err := tt.provider.MarshalJSON()
			if err != nil {
				t.Errorf("MarshalJSON() error = %v", err)
				return
			}
			if string(got) != tt.expected {
				t.Errorf("MarshalJSON() = %v, want %v", string(got), tt.expected)
			}
		})
	}
}

func TestAssistantProvider_Value(t *testing.T) {
	tests := []struct {
		provider AssistantProvider
		expected string
	}{
		{AGENTKIT, "AGENTKIT"},
		{WEBSOCKET, "WEBSOCKET"},
		{MODEL, "MODEL"},
	}
	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			got, err := tt.provider.Value()
			if err != nil {
				t.Errorf("Value() error = %v", err)
				return
			}
			if got != tt.expected {
				t.Errorf("Value() = %v, want %v", got, tt.expected)
			}
		})
	}
}
