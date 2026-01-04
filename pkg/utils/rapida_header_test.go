package utils

import "testing"

func TestHeaderConstants(t *testing.T) {
	// Just test that constants are not empty
	if HEADER_API_KEY == "" {
		t.Error("HEADER_API_KEY should not be empty")
	}
	if HEADER_AUTH_KEY == "" {
		t.Error("HEADER_AUTH_KEY should not be empty")
	}
	if HEADER_SOURCE_KEY == "" {
		t.Error("HEADER_SOURCE_KEY should not be empty")
	}
	if HEADER_ENVIRONMENT_KEY == "" {
		t.Error("HEADER_ENVIRONMENT_KEY should not be empty")
	}
	if HEADER_REGION_KEY == "" {
		t.Error("HEADER_REGION_KEY should not be empty")
	}
	// Add more if needed, but since they are just constants, this is sufficient
}