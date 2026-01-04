// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package type_enums

import "testing"

func TestToRecordState(t *testing.T) {
	tests := []struct {
		input    string
		expected RecordState
	}{
		{"ACTIVE", RECORD_ACTIVE},
		{"INACTIVE", RECORD_INACTIVE},
		{"unknown", RECORD_INACTIVE},
	}
	for _, tt := range tests {
		result := ToRecordState(tt.input)
		if result != tt.expected {
			t.Errorf("ToRecordState(%s) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestRecordState_String(t *testing.T) {
	if got := RECORD_ACTIVE.String(); got != "ACTIVE" {
		t.Errorf("String() = %v, want %v", got, "ACTIVE")
	}
}

func TestRecordState_MarshalJSON(t *testing.T) {
	got, err := RECORD_ACTIVE.MarshalJSON()
	if err != nil {
		t.Errorf("MarshalJSON() error = %v", err)
	}
	if string(got) != `"ACTIVE"` {
		t.Errorf("MarshalJSON() = %v, want %v", string(got), `"ACTIVE"`)
	}
}

func TestRecordState_Value(t *testing.T) {
	got, err := RECORD_ACTIVE.Value()
	if err != nil {
		t.Errorf("Value() error = %v", err)
	}
	if got != "ACTIVE" {
		t.Errorf("Value() = %v, want %v", got, "ACTIVE")
	}
}