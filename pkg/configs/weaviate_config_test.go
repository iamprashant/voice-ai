// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package configs

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func TestWeaviateConfig_Validation(t *testing.T) {
	validate := validator.New()
	tests := []struct {
		name    string
		cfg     WeaviateConfig
		wantErr bool
	}{
		{"valid", WeaviateConfig{Host: "localhost", Scheme: "http"}, false},
		{"valid empty", WeaviateConfig{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWeaviateConfig_FromViper(t *testing.T) {
	v := viper.New()
	v.Set("host", "weaviate.example.com")
	v.Set("scheme", "https")
	v.Set("auth.api_key", "weaviate-key")

	var cfg WeaviateConfig
	err := v.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if cfg.Host != "weaviate.example.com" {
		t.Errorf("Host = %v, want weaviate.example.com", cfg.Host)
	}
	if cfg.Scheme != "https" {
		t.Errorf("Scheme = %v, want https", cfg.Scheme)
	}
	if cfg.Auth.ApiKey != "weaviate-key" {
		t.Errorf("Auth.ApiKey = %v, want weaviate-key", cfg.Auth.ApiKey)
	}
}

func TestWeaviateConfig_Defaults(t *testing.T) {
	v := viper.New()
	// no sets

	var cfg WeaviateConfig
	err := v.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if cfg.Host != "" {
		t.Errorf("Host = %v, want empty", cfg.Host)
	}
	if cfg.Scheme != "" {
		t.Errorf("Scheme = %v, want empty", cfg.Scheme)
	}
	if cfg.Auth.ApiKey != "" {
		t.Errorf("Auth.ApiKey = %v, want empty", cfg.Auth.ApiKey)
	}
}
