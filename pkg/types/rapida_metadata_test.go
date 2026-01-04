// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package types

import (
	"testing"
)

func TestMetadata_SetValue(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{"string", "test", "test"},
		{"int", 123, "123"},
		{"float", 1.23, "1.230000"},
		{"bool", true, "true"},
		{"nil", nil, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metadata{}
			err := m.SetValue(tt.value)
			if err != nil {
				t.Errorf("SetValue() error = %v", err)
				return
			}
			if m.Value != tt.expected {
				t.Errorf("SetValue() Value = %v, want %v", m.Value, tt.expected)
			}
		})
	}
}

func TestNewMetadata(t *testing.T) {
	m := NewMetadata("key", "value")
	if m.Key != "key" {
		t.Errorf("NewMetadata() Key = %v, want %v", m.Key, "key")
	}
	if m.Value != "value" {
		t.Errorf("NewMetadata() Value = %v, want %v", m.Value, "value")
	}
}

func TestNewMetadataList(t *testing.T) {
	data := map[string]interface{}{"key1": "value1", "key2": 2}
	list := NewMetadataList(data)
	if len(list) != 2 {
		t.Errorf("NewMetadataList() length = %v, want %v", len(list), 2)
	}
	// Check if keys are present
	keys := make(map[string]bool)
	for _, m := range list {
		keys[m.Key] = true
		if m.Key == "key1" && m.Value != "value1" {
			t.Errorf("NewMetadataList() key1 value = %v, want %v", m.Value, "value1")
		}
		if m.Key == "key2" && m.Value != "2" {
			t.Errorf("NewMetadataList() key2 value = %v, want %v", m.Value, "2")
		}
	}
	if !keys["key1"] || !keys["key2"] {
		t.Errorf("NewMetadataList() keys = %v, want key1 and key2", keys)
	}
}
