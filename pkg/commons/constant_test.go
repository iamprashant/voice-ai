// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package commons

import (
	"testing"

	"github.com/rapidaai/protos"
)

func TestTelemetryIndex(t *testing.T) {
	tests := []struct {
		development bool
		expected    string
	}{
		{true, DEVELOPMENT_TELEMETRY_INDEX},
		{false, TELEMETRY_INDEX},
	}

	for _, tt := range tests {
		result := TelemetryIndex(tt.development)
		if result != tt.expected {
			t.Errorf("TelemetryIndex(%v) = %v, want %v", tt.development, result, tt.expected)
		}
	}
}

func TestEndpointIndex(t *testing.T) {
	tests := []struct {
		development bool
		expected    string
	}{
		{true, DEVELOPMENT_ENDPOINT_INDEX},
		{false, ENDPOINT_INDEX},
	}

	for _, tt := range tests {
		result := EndpointIndex(tt.development)
		if result != tt.expected {
			t.Errorf("EndpointIndex(%v) = %v, want %v", tt.development, result, tt.expected)
		}
	}
}

func TestAssistantIndex(t *testing.T) {
	tests := []struct {
		development bool
		expected    string
	}{
		{true, DEVELOPMENT_ASSISTANT_INDEX},
		{false, ASSISTANT_INDEX},
	}

	for _, tt := range tests {
		result := AssistantIndex(tt.development)
		if result != tt.expected {
			t.Errorf("AssistantIndex(%v) = %v, want %v", tt.development, result, tt.expected)
		}
	}
}

func TestKnowledgeIndex(t *testing.T) {
	tests := []struct {
		development bool
		org         uint64
		prjm        uint64
		kn          uint64
		expected    string
	}{
		{true, 1, 2, 3, "dev__vs__1__2__3"},
		{false, 1, 2, 3, "prod__vs__1__2__3"},
		{true, 0, 0, 0, "dev__vs__0__0__0"},
		{false, 123, 456, 789, "prod__vs__123__456__789"},
	}

	for _, tt := range tests {
		result := KnowledgeIndex(tt.development, tt.org, tt.prjm, tt.kn)
		if result != tt.expected {
			t.Errorf("KnowledgeIndex(%v, %d, %d, %d) = %v, want %v", tt.development, tt.org, tt.prjm, tt.kn, result, tt.expected)
		}
	}
}

func TestResponseContentType_String(t *testing.T) {
	tests := []struct {
		rct      ResponseContentType
		expected string
	}{
		{TEXT_CONTENT, "text"},
		{AUDIO_CONTENT, "audio"},
		{IMAGE_CONTENT, "image"},
		{MULTI_MEDIA_CONTENT, "multi"},
		{ResponseContentType("unknown"), "unknown"},
	}

	for _, tt := range tests {
		result := tt.rct.String()
		if result != tt.expected {
			t.Errorf("ResponseContentType(%v).String() = %v, want %v", tt.rct, result, tt.expected)
		}
	}
}

func TestResponseContentFormat_String(t *testing.T) {
	tests := []struct {
		rcf      ResponseContentFormat
		expected string
	}{
		{TEXT_CONTENT_FORMAT_RAW, "raw"},
		{AUDIO_CONTENT_FORMAT_RAW, "raw"},
		{AUDIO_CONTENT_FORMAT_URL, "url"},
		{IMAGE_CONTENT_FORMAT_RAW, "raw"},
		{IMAGE_CONTENT_FORMAT_URL, "url"},
		{MULTI_MEDIA_CONTENT_FORMAT_RAW, "raw"},
		{MULTI_MEDIA_CONTENT_FORMAT_URL, "url"},
		{ResponseContentFormat("unknown"), "unknown"},
	}

	for _, tt := range tests {
		result := tt.rcf.String()
		if result != tt.expected {
			t.Errorf("ResponseContentFormat(%v).String() = %v, want %v", tt.rcf, result, tt.expected)
		}
	}
}

func TestMessageContent_StringContent(t *testing.T) {
	tests := []struct {
		name     string
		contents []*protos.Content
		expected string
	}{
		{
			name:     "empty contents",
			contents: []*protos.Content{},
			expected: "",
		},
		{
			name: "single text raw content",
			contents: []*protos.Content{
				{
					ContentType:   string(TEXT_CONTENT),
					ContentFormat: string(TEXT_CONTENT_FORMAT_RAW),
					Content:       []byte("Hello World"),
				},
			},
			expected: "Hello World",
		},
		{
			name: "multiple contents with text",
			contents: []*protos.Content{
				{
					ContentType:   string(AUDIO_CONTENT),
					ContentFormat: string(AUDIO_CONTENT_FORMAT_RAW),
					Content:       []byte("audio data"),
				},
				{
					ContentType:   string(TEXT_CONTENT),
					ContentFormat: string(TEXT_CONTENT_FORMAT_RAW),
					Content:       []byte("Text content"),
				},
				{
					ContentType:   string(IMAGE_CONTENT),
					ContentFormat: string(IMAGE_CONTENT_FORMAT_URL),
					Content:       []byte("image url"),
				},
			},
			expected: "Text content",
		},
		{
			name: "text content but not raw format",
			contents: []*protos.Content{
				{
					ContentType:   string(TEXT_CONTENT),
					ContentFormat: string(AUDIO_CONTENT_FORMAT_URL), // not raw
					Content:       []byte("ignored"),
				},
			},
			expected: "",
		},
		{
			name: "multiple text contents",
			contents: []*protos.Content{
				{
					ContentType:   string(TEXT_CONTENT),
					ContentFormat: string(TEXT_CONTENT_FORMAT_RAW),
					Content:       []byte("First "),
				},
				{
					ContentType:   string(TEXT_CONTENT),
					ContentFormat: string(TEXT_CONTENT_FORMAT_RAW),
					Content:       []byte("Second"),
				},
			},
			expected: "First Second",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := &MessageContent{
				Contents: tt.contents,
			}
			result := mc.StringContent()
			if result != tt.expected {
				t.Errorf("StringContent() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestToMessageContent(t *testing.T) {
	original := &protos.Message{
		Role: "user",
		Contents: []*protos.Content{
			{
				ContentType:   string(TEXT_CONTENT),
				ContentFormat: string(TEXT_CONTENT_FORMAT_RAW),
				Content:       []byte("test"),
			},
		},
	}

	result := ToMessageContent(original)

	if result.Role != original.Role {
		t.Errorf("Role = %v, want %v", result.Role, original.Role)
	}

	if len(result.Contents) != len(original.Contents) {
		t.Errorf("Contents length = %d, want %d", len(result.Contents), len(original.Contents))
	}

	// Check that it's a copy, not the same slice
	if &result.Contents == &original.Contents {
		t.Error("Contents should be copied, not the same reference")
	}
}
