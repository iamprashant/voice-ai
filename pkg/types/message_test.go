// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package types

import (
	"reflect"
	"testing"
	"time"

	"github.com/rapidaai/pkg/commons"
	lexatic_backend "github.com/rapidaai/protos"
)

func TestOnlyStringContent(t *testing.T) {
	tests := []struct {
		name     string
		contents []*Content
		want     string
	}{
		{
			name:     "empty contents",
			contents: []*Content{},
			want:     "",
		},
		{
			name: "single text content",
			contents: []*Content{
				{
					ContentType:   string(commons.TEXT_CONTENT),
					ContentFormat: string(commons.TEXT_CONTENT_FORMAT_RAW),
					Content:       []byte("hello"),
				},
			},
			want: "hello",
		},
		{
			name: "multiple text contents",
			contents: []*Content{
				{
					ContentType:   string(commons.TEXT_CONTENT),
					ContentFormat: string(commons.TEXT_CONTENT_FORMAT_RAW),
					Content:       []byte("hello"),
				},
				{
					ContentType:   string(commons.TEXT_CONTENT),
					ContentFormat: string(commons.TEXT_CONTENT_FORMAT_RAW),
					Content:       []byte(" world"),
				},
			},
			want: "hello world",
		},
		{
			name: "mixed content types",
			contents: []*Content{
				{
					ContentType:   string(commons.TEXT_CONTENT),
					ContentFormat: string(commons.TEXT_CONTENT_FORMAT_RAW),
					Content:       []byte("text"),
				},
				{
					ContentType:   string(commons.AUDIO_CONTENT),
					ContentFormat: "mp3",
					Content:       []byte("audio data"),
				},
			},
			want: "text",
		},
		{
			name: "non-raw format",
			contents: []*Content{
				{
					ContentType:   string(commons.TEXT_CONTENT),
					ContentFormat: "html",
					Content:       []byte("<p>text</p>"),
				},
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OnlyStringContent(tt.contents); got != tt.want {
				t.Errorf("OnlyStringContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContentString(t *testing.T) {
	tests := []struct {
		name string
		c    *lexatic_backend.Content
		want string
	}{
		{
			name: "text content raw",
			c: &lexatic_backend.Content{
				ContentType:   string(commons.TEXT_CONTENT),
				ContentFormat: string(commons.TEXT_CONTENT_FORMAT_RAW),
				Content:       []byte("content"),
			},
			want: "content",
		},
		{
			name: "non-text content",
			c: &lexatic_backend.Content{
				ContentType:   string(commons.AUDIO_CONTENT),
				ContentFormat: "mp3",
				Content:       []byte("audio"),
			},
			want: "",
		},
		{
			name: "text non-raw",
			c: &lexatic_backend.Content{
				ContentType:   string(commons.TEXT_CONTENT),
				ContentFormat: "html",
				Content:       []byte("<p>text</p>"),
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContentString(tt.c); got != tt.want {
				t.Errorf("ContentString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOnlyStringProtoContent(t *testing.T) {
	tests := []struct {
		name     string
		contents []*lexatic_backend.Content
		want     string
	}{
		{
			name:     "empty",
			contents: []*lexatic_backend.Content{},
			want:     "",
		},
		{
			name: "single text",
			contents: []*lexatic_backend.Content{
				{
					ContentType:   string(commons.TEXT_CONTENT),
					ContentFormat: string(commons.TEXT_CONTENT_FORMAT_RAW),
					Content:       []byte("hello"),
				},
			},
			want: "hello",
		},
		{
			name: "multiple",
			contents: []*lexatic_backend.Content{
				{
					ContentType:   string(commons.TEXT_CONTENT),
					ContentFormat: string(commons.TEXT_CONTENT_FORMAT_RAW),
					Content:       []byte("hello"),
				},
				{
					ContentType:   string(commons.TEXT_CONTENT),
					ContentFormat: string(commons.TEXT_CONTENT_FORMAT_RAW),
					Content:       []byte(" world"),
				},
			},
			want: "hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OnlyStringProtoContent(tt.contents); got != tt.want {
				t.Errorf("OnlyStringProtoContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainsAudioContent(t *testing.T) {
	tests := []struct {
		name     string
		contents []*lexatic_backend.Content
		want     bool
	}{
		{
			name:     "empty",
			contents: []*lexatic_backend.Content{},
			want:     false,
		},
		{
			name: "has audio",
			contents: []*lexatic_backend.Content{
				{
					ContentType: string(commons.TEXT_CONTENT),
				},
				{
					ContentType: string(commons.AUDIO_CONTENT),
				},
			},
			want: true,
		},
		{
			name: "no audio",
			contents: []*lexatic_backend.Content{
				{
					ContentType: string(commons.TEXT_CONTENT),
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsAudioContent(tt.contents); got != tt.want {
				t.Errorf("ContainsAudioContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToMessage(t *testing.T) {
	protoMsg := &lexatic_backend.Message{

		Role: "user",
		Contents: []*lexatic_backend.Content{
			{
				ContentType: string(commons.TEXT_CONTENT),
				Content:     []byte("test"),
			},
		},
	}

	msg := ToMessage(protoMsg)
	if msg == nil {
		t.Errorf("ToMessage() returned nil")
		return
	}
	if msg.Role != "user" {
		t.Errorf("ToMessage() role = %v, want %v", msg.Role, "user")
	}
	if !reflect.DeepEqual(msg.Time, time.Now()) {
		// Time should be set to now, but for test, just check it's not zero
		if msg.Time.IsZero() {
			t.Errorf("ToMessage() time not set")
		}
	}
}

func TestToMessages(t *testing.T) {
	protoMsgs := []*lexatic_backend.Message{
		{
			Role: "user",
		},
		{
			Role: "assistant",
		},
	}

	msgs := ToMessages(protoMsgs)
	if len(msgs) != 2 {
		t.Errorf("ToMessages() length = %v, want %v", len(msgs), 2)
	}
}

func TestToSimpleMessage(t *testing.T) {
	msgs := []*Message{
		{
			Role: "user",
			Contents: []*Content{
				{
					ContentType:   string(commons.TEXT_CONTENT),
					ContentFormat: string(commons.TEXT_CONTENT_FORMAT_RAW),
					Content:       []byte("hello"),
				},
			},
			Time: time.Now(),
		},
		{
			Role: "assistant",
			Contents: []*Content{
				{
					ContentType: string(commons.AUDIO_CONTENT),
				},
			},
			Time: time.Now(),
		},
	}

	simple := ToSimpleMessage(msgs)
	if len(simple) != 1 {
		t.Errorf("ToSimpleMessage() length = %v, want %v", len(simple), 1)
	}
	if simple[0]["role"] != "user" {
		t.Errorf("ToSimpleMessage()[0][role] = %v, want %v", simple[0]["role"], "user")
	}
	if simple[0]["message"] != "hello" {
		t.Errorf("ToSimpleMessage()[0][message] = %v, want %v", simple[0]["message"], "hello")
	}
}
