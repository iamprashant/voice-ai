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

func TestBatchConfig_IsLocal(t *testing.T) {
	tests := []struct {
		name string
		cfg  BatchConfig
		want bool
	}{
		{"AWS", BatchConfig{BatchType: "aws"}, false},
		{"Local", BatchConfig{BatchType: "local"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cfg.IsLocal(); got != tt.want {
				t.Errorf("BatchConfig.IsLocal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBatchConfig_Validation(t *testing.T) {
	validate := validator.New()
	tests := []struct {
		name    string
		cfg     BatchConfig
		wantErr bool
	}{
		{"valid", BatchConfig{BatchType: "aws"}, false},
		{"valid empty", BatchConfig{}, false},
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

func TestBatchConfig_FromViper(t *testing.T) {
	v := viper.New()
	v.Set("batch_type", "aws")
	v.Set("batch_script", "script.sh")

	var cfg BatchConfig
	err := v.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if cfg.BatchType != "aws" {
		t.Errorf("BatchType = %v, want aws", cfg.BatchType)
	}
	if cfg.BatchScript != "script.sh" {
		t.Errorf("BatchScript = %v, want script.sh", cfg.BatchScript)
	}
	if cfg.Auth != nil {
		t.Errorf("Auth should be nil")
	}
}

func TestBatchConfig_Defaults(t *testing.T) {
	v := viper.New()
	// no sets

	var cfg BatchConfig
	err := v.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if cfg.BatchType != "" {
		t.Errorf("BatchType = %v, want empty", cfg.BatchType)
	}
	if cfg.BatchScript != "" {
		t.Errorf("BatchScript = %v, want empty", cfg.BatchScript)
	}
	if cfg.Auth != nil {
		t.Errorf("Auth should be nil")
	}
	if !cfg.IsLocal() {
		t.Errorf("IsLocal() = false, want true for default")
	}
}
