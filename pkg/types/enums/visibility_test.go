// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package type_enums

import "testing"

func TestToRecordVisibility(t *testing.T) {
	tests := []struct {
		input    string
		expected RecordVisibility
	}{
		{"public", RECORD_PUBLIC},
		{"private", RECORD_PRIVATE},
		{"unknown", RECORD_PRIVATE},
	}
	for _, tt := range tests {
		result := ToRecordVisibility(tt.input)
		if result != tt.expected {
			t.Errorf("ToRecordVisibility(%s) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestRecordVisibility_String(t *testing.T) {
	if got := RECORD_PUBLIC.String(); got != "PUBLIC" {
		t.Errorf("String() = %v, want %v", got, "PUBLIC")
	}
}

func TestRecordVisibility_MarshalJSON(t *testing.T) {
	got, err := RECORD_PUBLIC.MarshalJSON()
	if err != nil {
		t.Errorf("MarshalJSON() error = %v", err)
	}
	if string(got) != `"PUBLIC"` {
		t.Errorf("MarshalJSON() = %v, want %v", string(got), `"PUBLIC"`)
	}
}

func TestRecordVisibility_Value(t *testing.T) {
	got, err := RECORD_PUBLIC.Value()
	if err != nil {
		t.Errorf("Value() error = %v", err)
	}
	if got != "PUBLIC" {
		t.Errorf("Value() = %v, want %v", got, "PUBLIC")
	}
}