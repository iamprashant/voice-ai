// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package commons

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestResponse_JSONMarshal(t *testing.T) {
	resp := Response{
		Code:    200,
		Success: true,
		Data:    map[string]string{"key": "value"},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal Response: %v", err)
	}

	expected := `{"code":200,"success":true,"data":{"key":"value"}}`
	if string(data) != expected {
		t.Errorf("Marshal result = %v, want %v", string(data), expected)
	}
}

func TestResponse_JSONUnmarshal(t *testing.T) {
	jsonStr := `{"code":404,"success":false,"data":{"error":"not found"}}`

	var resp Response
	err := json.Unmarshal([]byte(jsonStr), &resp)
	if err != nil {
		t.Fatalf("Failed to unmarshal Response: %v", err)
	}

	if resp.Code != 404 {
		t.Errorf("Code = %v, want %v", resp.Code, 404)
	}
	if resp.Success != false {
		t.Errorf("Success = %v, want %v", resp.Success, false)
	}
	if data, ok := resp.Data.(map[string]interface{}); ok {
		if data["error"] != "not found" {
			t.Errorf("Data.error = %v, want %v", data["error"], "not found")
		}
	} else {
		t.Errorf("Data is not a map: %v", resp.Data)
	}
}

func TestErrorMessage_JSONMarshal(t *testing.T) {
	errMsg := ErrorMessage{
		Code:    500,
		Message: errors.New("internal server error"),
	}

	data, err := json.Marshal(errMsg)
	if err != nil {
		t.Fatalf("Failed to marshal ErrorMessage: %v", err)
	}

	// Since error marshals to {}, we check the code field
	var unmarshaled map[string]interface{}
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal for verification: %v", err)
	}

	if unmarshaled["code"] != float64(500) {
		t.Errorf("Code = %v, want %v", unmarshaled["code"], 500)
	}
	// Note: Message field marshals to {} since error has no exported fields
	if _, ok := unmarshaled["message"].(map[string]interface{}); !ok {
		t.Errorf("Message should be an empty object")
	}
}

func TestHealthCheck_JSONMarshal(t *testing.T) {
	hc := HealthCheck{
		Healthy: true,
	}

	data, err := json.Marshal(hc)
	if err != nil {
		t.Fatalf("Failed to marshal HealthCheck: %v", err)
	}

	expected := `{"healthy":true}`
	if string(data) != expected {
		t.Errorf("Marshal result = %v, want %v", string(data), expected)
	}
}

func TestHealthCheck_JSONUnmarshal(t *testing.T) {
	jsonStr := `{"healthy":false}`

	var hc HealthCheck
	err := json.Unmarshal([]byte(jsonStr), &hc)
	if err != nil {
		t.Fatalf("Failed to unmarshal HealthCheck: %v", err)
	}

	if hc.Healthy != false {
		t.Errorf("Healthy = %v, want %v", hc.Healthy, false)
	}
}