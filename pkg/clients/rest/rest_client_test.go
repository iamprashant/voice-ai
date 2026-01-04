// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/rapidaai/config"
	"github.com/rapidaai/pkg/commons"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRestClient(t *testing.T) {
	logger, _ := commons.NewApplicationLogger()
	cfg := &config.AppConfig{}
	baseURL := "https://api.example.com"

	client := NewRestClient(logger, cfg, baseURL)

	assert.NotNil(t, client)
	assert.Equal(t, logger, client.logger)
	assert.Equal(t, cfg, client.cfg)
	assert.Equal(t, baseURL, client.BaseURL)
	assert.NotNil(t, client.HTTPClient)
	assert.Equal(t, 60*time.Second, client.HTTPClient.Timeout)
	assert.NotNil(t, client.Headers)
}

func TestNewRestClientWithConfig(t *testing.T) {
	baseURL := "https://api.example.com"
	defaultHeaders := map[string]string{
		"Authorization": "Bearer token",
		"User-Agent":    "test-client",
	}
	timeout := uint32(30)

	client := NewRestClientWithConfig(baseURL, defaultHeaders, timeout)

	assert.NotNil(t, client)
	assert.Equal(t, baseURL, client.BaseURL)
	assert.NotNil(t, client.HTTPClient)
	assert.Equal(t, 30*time.Second, client.HTTPClient.Timeout)
	assert.Equal(t, defaultHeaders, client.Headers)
	assert.Nil(t, client.logger)
	assert.Nil(t, client.cfg)
}

func TestRestClient_Get_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/test", r.URL.Path)
		assert.Equal(t, "value1", r.URL.Query().Get("param1"))
		assert.Equal(t, "value2", r.URL.Query().Get("param2"))
		assert.Equal(t, "Bearer token", r.Header.Get("Authorization"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "success"})
	}))
	defer server.Close()

	client := NewRestClientWithConfig(server.URL, map[string]string{"Authorization": "Bearer token"}, 10)
	params := map[string]interface{}{
		"param1": "value1",
		"param2": "value2",
	}
	headers := map[string]string{"X-Custom": "custom-value"}

	resp, err := client.Get(context.Background(), "/test", params, headers)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "200 OK", resp.Status)
	assert.Contains(t, resp.Headers.Get("Content-Type"), "application/json")

	var result map[string]string
	err = json.Unmarshal(resp.Body, &result)
	require.NoError(t, err)
	assert.Equal(t, "success", result["message"])
}

func TestRestClient_Post_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/create", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer token", r.Header.Get("Authorization"))

		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		assert.Equal(t, "test-value", body["key"])

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": "123"})
	}))
	defer server.Close()

	client := NewRestClientWithConfig(server.URL, map[string]string{"Authorization": "Bearer token"}, 10)
	body := map[string]string{"key": "test-value"}
	headers := map[string]string{"X-Custom": "custom"}

	resp, err := client.Post(context.Background(), "/create", body, headers)

	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var result map[string]string
	err = json.Unmarshal(resp.Body, &result)
	require.NoError(t, err)
	assert.Equal(t, "123", result["id"])
}

func TestRestClient_Put_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/update/123", r.URL.Path)

		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		assert.Equal(t, "updated-value", body["key"])

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
	}))
	defer server.Close()

	client := NewRestClientWithConfig(server.URL, nil, 10)
	body := map[string]string{"key": "updated-value"}

	resp, err := client.Put(context.Background(), "/update/123", body, nil)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRestClient_Patch_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		assert.Equal(t, "/patch/123", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"patched": "true"})
	}))
	defer server.Close()

	client := NewRestClientWithConfig(server.URL, nil, 10)

	resp, err := client.Patch(context.Background(), "/patch/123", nil, nil)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRestClient_Delete_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/delete/123", r.URL.Path)

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewRestClientWithConfig(server.URL, nil, 10)

	resp, err := client.Delete(context.Background(), "/delete/123", nil)

	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestRestClient_Request_CustomMethod(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "CUSTOM", r.Method)
		assert.Equal(t, "/custom", r.URL.Path)

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewRestClientWithConfig(server.URL, nil, 10)

	resp, err := client.Request(context.Background(), "CUSTOM", "/custom", nil, nil, nil)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRestClient_Request_WithParamsAndHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/search", r.URL.Path)
		assert.Equal(t, "golang", r.URL.Query().Get("q"))
		assert.Equal(t, "10", r.URL.Query().Get("limit"))
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.Equal(t, "Bearer token", r.Header.Get("Authorization"))
		assert.Equal(t, "custom-value", r.Header.Get("X-Custom"))

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewRestClientWithConfig(server.URL, map[string]string{"Authorization": "Bearer token"}, 10)
	params := map[string]interface{}{
		"q":     "golang",
		"limit": 10,
	}
	headers := map[string]string{
		"Accept":   "application/json",
		"X-Custom": "custom-value",
	}

	resp, err := client.Request(context.Background(), http.MethodGet, "/search", nil, params, headers)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRestClient_Request_BodyMarshalingError(t *testing.T) {
	client := NewRestClientWithConfig("https://api.example.com", nil, 10)

	// Use a channel which cannot be marshaled to JSON
	body := make(chan int)

	_, err := client.Request(context.Background(), http.MethodPost, "/test", body, nil, nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to marshal request body")
}

func TestRestClient_Request_InvalidURL(t *testing.T) {
	client := NewRestClientWithConfig("http://invalid url", nil, 10)

	_, err := client.Request(context.Background(), http.MethodGet, "/test", nil, nil, nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to build URL")
}

func TestRestClient_Request_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
	}))
	defer server.Close()

	client := NewRestClientWithConfig(server.URL, nil, 10)

	resp, err := client.Request(context.Background(), http.MethodGet, "/error", nil, nil, nil)

	require.NoError(t, err) // HTTP errors are not returned as Go errors
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestRestClient_Request_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // Sleep longer than timeout
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewRestClientWithConfig(server.URL, nil, 1) // 1 second timeout

	_, err := client.Request(context.Background(), http.MethodGet, "/timeout", nil, nil, nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to execute request")
}

func TestRestClient_Request_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewRestClientWithConfig(server.URL, nil, 10)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := client.Request(ctx, http.MethodGet, "/cancel", nil, nil, nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to execute request")
}

func TestRestClient_buildURL_Success(t *testing.T) {
	client := NewRestClientWithConfig("https://api.example.com", nil, 10)

	tests := []struct {
		name     string
		endpoint string
		params   map[string]interface{}
		expected string
	}{
		{
			name:     "no params",
			endpoint: "/users",
			params:   nil,
			expected: "https://api.example.com/users",
		},
		{
			name:     "with params",
			endpoint: "/users",
			params:   map[string]interface{}{"page": 1, "limit": 10},
			expected: "https://api.example.com/users?limit=10&page=1",
		},
		{
			name:     "endpoint with query",
			endpoint: "/users?sort=name",
			params:   map[string]interface{}{"page": 1},
			expected: "https://api.example.com/users?page=1&sort=name",
		},
		{
			name:     "full URL endpoint",
			endpoint: "https://other.api.com/users",
			params:   nil,
			expected: "https://other.api.com/users",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := client.buildURL(tt.endpoint, tt.params)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, url)
		})
	}
}

func TestRestClient_buildURL_InvalidBaseURL(t *testing.T) {
	client := &RestClient{BaseURL: "http://invalid url"}

	_, err := client.buildURL("/test", nil)

	assert.Error(t, err)
}

func TestRestClient_buildURL_InvalidEndpoint(t *testing.T) {
	client := NewRestClientWithConfig("https://api.example.com", nil, 10)

	_, err := client.buildURL("http://invalid url", nil)

	assert.Error(t, err)
}

func TestRestClient_SetDefaultHeader(t *testing.T) {
	client := NewRestClientWithConfig("https://api.example.com", nil, 10)

	client.SetDefaultHeader("Authorization", "Bearer token")
	client.SetDefaultHeader("User-Agent", "test-client")

	assert.Equal(t, "Bearer token", client.Headers["Authorization"])
	assert.Equal(t, "test-client", client.Headers["User-Agent"])
}

func TestRestClient_RemoveDefaultHeader(t *testing.T) {
	client := NewRestClientWithConfig("https://api.example.com", map[string]string{
		"Auth": "token",
		"Key":  "value",
	}, 10)

	client.RemoveDefaultHeader("Auth")

	assert.NotContains(t, client.Headers, "Auth")
	assert.Contains(t, client.Headers, "Key")
	assert.Equal(t, "value", client.Headers["Key"])
}

func TestAPIResponse_ToJSON(t *testing.T) {
	resp := &APIResponse{
		StatusCode: 200,
		Status:     "OK",
		Body:       []byte(`{"message": "test"}`),
		Headers:    make(http.Header),
	}

	jsonData, err := resp.ToJSON()

	require.NoError(t, err)
	// Just verify it's valid JSON
	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	require.NoError(t, err)
	assert.Equal(t, float64(200), result["status_code"])
	assert.Equal(t, "OK", result["status"])
}

func TestAPIResponse_ToMap(t *testing.T) {
	resp := &APIResponse{
		Body: []byte(`{"key": "value", "number": 42}`),
	}

	result, err := resp.ToMap()

	require.NoError(t, err)
	assert.Equal(t, "value", result["key"])
	assert.Equal(t, float64(42), result["number"])
}

func TestAPIResponse_ToMap_InvalidJSON(t *testing.T) {
	resp := &APIResponse{
		Body: []byte(`invalid json`),
	}

	_, err := resp.ToMap()

	assert.Error(t, err)
}

func TestAPIResponse_ToString(t *testing.T) {
	body := `{"message": "hello world"}`
	resp := &APIResponse{
		Body: []byte(body),
	}

	result := resp.ToString()

	assert.Equal(t, body, result)
}

func TestRestClient_Request_ContentTypeHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		if r.Method == http.MethodPost {
			assert.Equal(t, "application/json", contentType)
		} else {
			assert.Empty(t, contentType)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewRestClientWithConfig(server.URL, nil, 10)

	// POST should set Content-Type
	resp, err := client.Post(context.Background(), "/test", map[string]string{"key": "value"}, nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// GET should not set Content-Type
	resp, err = client.Get(context.Background(), "/test", nil, nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRestClient_Request_CustomContentType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/xml", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewRestClientWithConfig(server.URL, nil, 10)
	headers := map[string]string{"Content-Type": "application/xml"}

	resp, err := client.Post(context.Background(), "/test", "<xml></xml>", headers)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRestClient_Request_EmptyBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := make([]byte, 0)
		r.Body.Read(body)
		assert.Empty(t, body)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewRestClientWithConfig(server.URL, nil, 10)

	resp, err := client.Post(context.Background(), "/test", nil, nil)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRestClient_Request_LargeResponse(t *testing.T) {
	largeData := strings.Repeat("a", 1024*1024) // 1MB response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(largeData))
	}))
	defer server.Close()

	client := NewRestClientWithConfig(server.URL, nil, 10)

	resp, err := client.Get(context.Background(), "/large", nil, nil)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, largeData, string(resp.Body))
}

func TestRestClient_Request_Redirect(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/redirect" {
			w.Header().Set("Location", "/final")
			w.WriteHeader(http.StatusFound)
			return
		}
		if r.URL.Path == "/final" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"redirected": "true"})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewRestClientWithConfig(server.URL, nil, 10)

	resp, err := client.Get(context.Background(), "/redirect", nil, nil)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]string
	err = json.Unmarshal(resp.Body, &result)
	require.NoError(t, err)
	assert.Equal(t, "true", result["redirected"])
}

func TestRestClient_Request_QueryParamEncoding(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "hello world", r.URL.Query().Get("message"))
		assert.Equal(t, "a+b=c&d=e", r.URL.Query().Get("complex"))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewRestClientWithConfig(server.URL, nil, 10)
	params := map[string]interface{}{
		"message": "hello world",
		"complex": "a+b=c&d=e",
	}

	resp, err := client.Get(context.Background(), "/test", params, nil)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRestClient_Request_NonJSONBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := make([]byte, 100)
		n, _ := r.Body.Read(body)
		assert.Equal(t, `"plain text body"`, string(body[:n])) // JSON encoded
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewRestClientWithConfig(server.URL, nil, 10)
	headers := map[string]string{"Content-Type": "text/plain"}

	resp, err := client.Post(context.Background(), "/test", "plain text body", headers)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
