// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package types

import (
	"testing"
)

func TestEvent_SetValue(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected map[string]interface{}
	}{
		{"string", "test", map[string]interface{}{"value": "test"}},
		{"int", 123, map[string]interface{}{"value": "123"}},
		{"float", 1.23, map[string]interface{}{"value": "1.230000"}},
		{"bool", true, map[string]interface{}{"value": "true"}},
		{"nil", nil, map[string]interface{}{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Event{}
			err := e.SetValue(tt.value)
			if err != nil {
				t.Errorf("SetValue() error = %v", err)
				return
			}
			if len(e.Payload) != len(tt.expected) {
				t.Errorf("SetValue() payload = %v, want %v", e.Payload, tt.expected)
			}
			for k, v := range tt.expected {
				if e.Payload[k] != v {
					t.Errorf("SetValue() payload[%s] = %v, want %v", k, e.Payload[k], v)
				}
			}
		})
	}
}

func TestNewEvent(t *testing.T) {
	e := NewEvent("test", "value")
	if e.EventType != "test" {
		t.Errorf("NewEvent() EventType = %v, want %v", e.EventType, "test")
	}
	if e.Payload["value"] != "value" {
		t.Errorf("NewEvent() Payload = %v, want %v", e.Payload, map[string]interface{}{"value": "value"})
	}
}
