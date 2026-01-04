// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package type_enums

import "testing"

func TestToConversationDirection(t *testing.T) {
	tests := []struct {
		input    string
		expected ConversationDirection
	}{
		{"inbound", DIRECTION_INBOUND},
		{"outbound", DIRECTION_OUTBOUND},
		{"unknown", DIRECTION_INBOUND},
	}
	for _, tt := range tests {
		result := ToConversationDirection(tt.input)
		if result != tt.expected {
			t.Errorf("ToConversationDirection(%s) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestConversationDirection_String(t *testing.T) {
	if got := DIRECTION_INBOUND.String(); got != "inbound" {
		t.Errorf("String() = %v, want %v", got, "inbound")
	}
	if got := DIRECTION_OUTBOUND.String(); got != "outbound" {
		t.Errorf("String() = %v, want %v", got, "outbound")
	}
}

func TestConversationDirection_MarshalJSON(t *testing.T) {
	got, err := DIRECTION_INBOUND.MarshalJSON()
	if err != nil {
		t.Errorf("MarshalJSON() error = %v", err)
	}
	if string(got) != `"inbound"` {
		t.Errorf("MarshalJSON() = %v, want %v", string(got), `"inbound"`)
	}
}

func TestConversationDirection_Value(t *testing.T) {
	got, err := DIRECTION_INBOUND.Value()
	if err != nil {
		t.Errorf("Value() error = %v", err)
	}
	if got != "inbound" {
		t.Errorf("Value() = %v, want %v", got, "inbound")
	}
}

func TestMessageActor(t *testing.T) {
	if !UserActor.ActingUser() {
		t.Error("UserActor should be acting user")
	}
	if UserActor.ActingAssistant() {
		t.Error("UserActor should not be acting assistant")
	}
	if !AssistantActor.ActingAssistant() {
		t.Error("AssistantActor should be acting assistant")
	}
	if AssistantActor.ActingUser() {
		t.Error("AssistantActor should not be acting user")
	}
}

func TestMessageMode(t *testing.T) {
	if !AudioMode.Audio() {
		t.Error("AudioMode should be audio")
	}
	if AudioMode.Text() {
		t.Error("AudioMode should not be text")
	}
	if !TextMode.Text() {
		t.Error("TextMode should be text")
	}
	if TextMode.Audio() {
		t.Error("TextMode should not be audio")
	}
	if AudioMode.String() != "audio" {
		t.Errorf("AudioMode.String() = %v, want %v", AudioMode.String(), "audio")
	}
	if TextMode.String() != "text" {
		t.Errorf("TextMode.String() = %v, want %v", TextMode.String(), "text")
	}
}

func TestToMessageAction(t *testing.T) {
	tests := []struct {
		input    string
		expected MessageAction
	}{
		{"tool-call", ACTION_TOOL_CALL},
		{"llm-call", ACTION_LLM_CALL},
		{"unknown", ACTION_LLM_CALL},
	}
	for _, tt := range tests {
		result := ToMessageAction(tt.input)
		if result != tt.expected {
			t.Errorf("ToMessageAction(%s) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestMessageAction_String(t *testing.T) {
	if got := ACTION_TOOL_CALL.String(); got != "tool-call" {
		t.Errorf("String() = %v, want %v", got, "tool-call")
	}
}

func TestMessageAction_MarshalJSON(t *testing.T) {
	got, err := ACTION_TOOL_CALL.MarshalJSON()
	if err != nil {
		t.Errorf("MarshalJSON() error = %v", err)
	}
	if string(got) != `"tool-call"` {
		t.Errorf("MarshalJSON() = %v, want %v", string(got), `"tool-call"`)
	}
}

func TestMessageAction_Value(t *testing.T) {
	got, err := ACTION_TOOL_CALL.Value()
	if err != nil {
		t.Errorf("Value() error = %v", err)
	}
	if got != "tool-call" {
		t.Errorf("Value() = %v, want %v", got, "tool-call")
	}
}
