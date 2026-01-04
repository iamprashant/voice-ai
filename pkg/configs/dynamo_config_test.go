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

func TestDynamoConfig_Validation(t *testing.T) {
	validate := validator.New()
	tests := []struct {
		name    string
		cfg     DynamoConfig
		wantErr bool
	}{
		{"valid", DynamoConfig{Auth: AwsConfig{Region: "us-east-1"}, MaxRetries: 3}, false},
		{"invalid missing max_retries", DynamoConfig{Auth: AwsConfig{Region: "us-east-1"}}, true},
		{"invalid missing auth region", DynamoConfig{MaxRetries: 3}, true},
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

func TestDynamoConfig_FromViper(t *testing.T) {
	v := viper.New()
	v.Set("auth.region", "us-west-2")
	v.Set("auth.assume_role", "test-role")
	v.Set("max_retries", 5)

	var cfg DynamoConfig
	err := v.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if cfg.Auth.Region != "us-west-2" {
		t.Errorf("Auth.Region = %v, want us-west-2", cfg.Auth.Region)
	}
	if cfg.Auth.AssumeRole != "test-role" {
		t.Errorf("Auth.AssumeRole = %v, want test-role", cfg.Auth.AssumeRole)
	}
	if cfg.MaxRetries != 5 {
		t.Errorf("MaxRetries = %v, want 5", cfg.MaxRetries)
	}
}

func TestDynamoConfig_Defaults(t *testing.T) {
	v := viper.New()
	// no sets

	var cfg DynamoConfig
	err := v.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if cfg.Auth.Region != "" {
		t.Errorf("Auth.Region = %v, want empty", cfg.Auth.Region)
	}
	if cfg.Auth.AssumeRole != "" {
		t.Errorf("Auth.AssumeRole = %v, want empty", cfg.Auth.AssumeRole)
	}
	if cfg.Auth.AccessKeyId != "" {
		t.Errorf("Auth.AccessKeyId = %v, want empty", cfg.Auth.AccessKeyId)
	}
	if cfg.Auth.SecretKey != "" {
		t.Errorf("Auth.SecretKey = %v, want empty", cfg.Auth.SecretKey)
	}
	if cfg.MaxRetries != 0 {
		t.Errorf("MaxRetries = %v, want 0", cfg.MaxRetries)
	}
}
