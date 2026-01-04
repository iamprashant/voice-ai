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
)

func TestFunctionCall_MergeArguments(t *testing.T) {
	tests := []struct {
		name     string
		fc       *FunctionCall
		newArgs  *string
		expected *string
	}{
		{
			name: "nil newArgs",
			fc: &FunctionCall{
				Arguments: &[]string{"old"}[0],
			},
			newArgs:  nil,
			expected: &[]string{"old"}[0],
		},
		{
			name: "nil Arguments",
			fc: &FunctionCall{
				Arguments: nil,
			},
			newArgs:  &[]string{"new"}[0],
			expected: &[]string{"new"}[0],
		},
		{
			name: "append",
			fc: &FunctionCall{
				Arguments: &[]string{"old"}[0],
			},
			newArgs:  &[]string{"new"}[0],
			expected: func() *string { s := "oldnew"; return &s }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fc.MergeArguments(tt.newArgs)
			if !reflect.DeepEqual(tt.fc.Arguments, tt.expected) {
				t.Errorf("MergeArguments() = %v, want %v", tt.fc.Arguments, tt.expected)
			}
		})
	}
}

func TestFunctionCall_MergeName(t *testing.T) {
	tests := []struct {
		name     string
		fc       *FunctionCall
		newName  *string
		expected *string
	}{
		{
			name: "nil newName",
			fc: &FunctionCall{
				Name: &[]string{"old"}[0],
			},
			newName:  nil,
			expected: &[]string{"old"}[0],
		},
		{
			name: "nil Name",
			fc: &FunctionCall{
				Name: nil,
			},
			newName:  &[]string{"new"}[0],
			expected: &[]string{"new"}[0],
		},
		{
			name: "append",
			fc: &FunctionCall{
				Name: &[]string{"old"}[0],
			},
			newName:  &[]string{"new"}[0],
			expected: func() *string { s := "oldnew"; return &s }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fc.MergeName(tt.newName)
			if !reflect.DeepEqual(tt.fc.Name, tt.expected) {
				t.Errorf("MergeName() = %v, want %v", tt.fc.Name, tt.expected)
			}
		})
	}
}

func TestMessage_GetLanguage(t *testing.T) {
	tests := []struct {
		name     string
		meta     map[string]interface{}
		expected Language
	}{
		{
			name: "with language",
			meta: map[string]interface{}{
				"language_iso_1": "fr",
			},
			expected: GetLanguageByName("fr"),
		},
		{
			name:     "without language",
			meta:     nil,
			expected: GetLanguageByName("en"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &Message{
				Meta: tt.meta,
			}
			got := msg.GetLanguage()
			if got != tt.expected {
				t.Errorf("GetLanguage() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestMessage_AddMetadata(t *testing.T) {
	msg := &Message{}
	msg.AddMetadata("key1", "value1")
	if msg.Meta == nil {
		t.Errorf("AddMetadata() did not initialize Meta")
		return
	}
	if msg.Meta["key1"] != "value1" {
		t.Errorf("AddMetadata() = %v, want %v", msg.Meta["key1"], "value1")
	}
}

func TestMessage_WithMetadata(t *testing.T) {
	msg := &Message{}
	result := msg.WithMetadata(map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	})
	if result != msg {
		t.Errorf("WithMetadata() did not return self")
	}
	if msg.Meta["key1"] != "value1" {
		t.Errorf("WithMetadata() key1 = %v, want %v", msg.Meta["key1"], "value1")
	}
}

func TestMessage_GetId(t *testing.T) {
	msg := &Message{Id: "test-id"}
	if got := msg.GetId(); got != "test-id" {
		t.Errorf("GetId() = %v, want %v", got, "test-id")
	}
}

func TestMessage_GetContents(t *testing.T) {
	contents := []*Content{{}}
	msg := &Message{Contents: contents}
	if got := msg.GetContents(); !reflect.DeepEqual(got, contents) {
		t.Errorf("GetContents() = %v, want %v", got, contents)
	}
}

func TestMessage_GetToolCalls(t *testing.T) {
	toolCalls := []*ToolCall{{}}
	msg := &Message{ToolCalls: toolCalls}
	if got := msg.GetToolCalls(); !reflect.DeepEqual(got, toolCalls) {
		t.Errorf("GetToolCalls() = %v, want %v", got, toolCalls)
	}
}

func TestMessage_GetRole(t *testing.T) {
	msg := &Message{Role: "user"}
	if got := msg.GetRole(); got != "user" {
		t.Errorf("GetRole() = %v, want %v", got, "user")
	}
}

func TestMessage_GetTime(t *testing.T) {
	now := time.Now()
	msg := &Message{Time: now}
	got := msg.GetTime()
	// Should be RFC3339 formatted
	expected := now.UTC().Format(time.RFC3339)
	if got != expected {
		t.Errorf("GetTime() = %v, want %v", got, expected)
	}
}

func TestMessage_ToProto(t *testing.T) {
	msg := &Message{
		Id:   "id",
		Role: "role",
	}
	proto := msg.ToProto()
	if proto.Role != "role" {
		t.Errorf("ToProto() Role = %v, want %v", proto.Role, "role")
	}
}

func TestMessage_String(t *testing.T) {
	tests := []struct {
		name     string
		contents []*Content
		expected string
	}{
		{
			name:     "empty",
			contents: []*Content{},
			expected: "",
		},
		{
			name: "single text",
			contents: []*Content{
				{
					ContentType:   string(commons.TEXT_CONTENT),
					ContentFormat: string(commons.TEXT_CONTENT_FORMAT_RAW),
					Content:       []byte("hello"),
				},
			},
			expected: "hello",
		},
		{
			name: "multiple text",
			contents: []*Content{
				{
					ContentType:   string(commons.TEXT_CONTENT),
					ContentFormat: string(commons.TEXT_CONTENT_FORMAT_RAW),
					Content:       []byte("hello"),
				},
				{
					ContentType:   string(commons.TEXT_CONTENT),
					ContentFormat: string(commons.TEXT_CONTENT_FORMAT_RAW),
					Content:       []byte("world"),
				},
			},
			expected: "hello world",
		},
		{
			name: "with space",
			contents: []*Content{
				{
					ContentType:   string(commons.TEXT_CONTENT),
					ContentFormat: string(commons.TEXT_CONTENT_FORMAT_RAW),
					Content:       []byte("hello"),
				},
				{
					ContentType:   string(commons.TEXT_CONTENT),
					ContentFormat: string(commons.TEXT_CONTENT_FORMAT_RAW),
					Content:       []byte("world"),
				},
			},
			expected: "hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &Message{Contents: tt.contents}
			if got := msg.String(); got != tt.expected {
				t.Errorf("String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNewMessage(t *testing.T) {
	role := "user"
	content := []*Content{{Content: []byte("test")}}
	msg := NewMessage(role, content...)
	if msg.Role != role {
		t.Errorf("NewMessage() Role = %v, want %v", msg.Role, role)
	}
	if len(msg.Contents) != 1 {
		t.Errorf("NewMessage() Contents length = %v, want %v", len(msg.Contents), 1)
	}
	if msg.Id == "" {
		t.Errorf("NewMessage() Id not set")
	}
	if msg.Time.IsZero() {
		t.Errorf("NewMessage() Time not set")
	}
}

func TestMessage_MergeContent(t *testing.T) {
	msg := &Message{
		Contents: []*Content{{Content: []byte("old")}},
	}
	newContent := []*Content{{Content: []byte("new")}}
	result := msg.MergeContent(newContent...)
	if result != msg {
		t.Errorf("MergeContent() did not return self")
	}
	if len(msg.Contents) != 2 {
		t.Errorf("MergeContent() length = %v, want %v", len(msg.Contents), 2)
	}
}
