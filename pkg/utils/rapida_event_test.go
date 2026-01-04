// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package utils

import "testing"

func TestRapidaEvent_Get(t *testing.T) {
	tests := []struct {
		event    RapidaEvent
		expected string
	}{
		{TalkPause, "talk.onPause"},
		{TalkInterruption, "talk.onInterrupt"},
		{TalkTranscript, "talk.onTranscript"},
		{TalkStart, "talk.onStart"},
		{TalkComplete, "talk.onComplete"},
		{TalkGeneration, "talk.onGeneration"},
		{TalkCompleteGeneration, "talk.onCompleteGeneration"},
		{TalkStartConversation, "talk.onStartConversation"},
		{TalkCompleteConversation, "talk.onCompleteConversation"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if result := tt.event.Get(); result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestAssistantServerEvent_Get(t *testing.T) {
	event := AssistantInitiated
	expected := "conversation.initiated"
	if result := event.Get(); result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestAssistantWebhookEvent_Get(t *testing.T) {
	tests := []struct {
		event    AssistantWebhookEvent
		expected string
	}{
		{MessageReceived, "message.received"},
		{MessageSent, "message.sent"},
		{ConversationBegin, "conversation.begin"},
		{ConversationResume, "conversation.resume"},
		{ConversationCompleted, "conversation.completed"},
		{ConversationFailed, "conversation.failed"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if result := tt.event.Get(); result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}
