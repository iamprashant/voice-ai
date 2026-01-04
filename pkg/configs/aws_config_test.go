// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package configs

import (
	"reflect"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func TestAwsConfig_ToMap(t *testing.T) {
	cfg := AwsConfig{Region: "us-east-1", AssumeRole: "role", AccessKeyId: "key", SecretKey: "secret"}
	want := map[string]interface{}{"region": "us-east-1", "assume_role": "role", "access_key_id": "key", "secret_key": "secret"}
	if got := cfg.ToMap(); !reflect.DeepEqual(got, want) {
		t.Errorf("AwsConfig.ToMap() = %v, want %v", got, want)
	}
}

func TestAwsConfig_Validation(t *testing.T) {
	validate := validator.New()
	tests := []struct {
		name    string
		cfg     AwsConfig
		wantErr bool
	}{
		{"valid", AwsConfig{Region: "us-east-1"}, false},
		{"invalid empty region", AwsConfig{}, true},
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

func TestAwsConfig_FromViper(t *testing.T) {
	v := viper.New()
	v.Set("region", "us-west-2")
	v.Set("assume_role", "test-role")
	v.Set("access_key_id", "test-key")
	v.Set("secret_key", "test-secret")

	var cfg AwsConfig
	err := v.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if cfg.Region != "us-west-2" {
		t.Errorf("Region = %v, want us-west-2", cfg.Region)
	}
	if cfg.AssumeRole != "test-role" {
		t.Errorf("AssumeRole = %v, want test-role", cfg.AssumeRole)
	}
	if cfg.AccessKeyId != "test-key" {
		t.Errorf("AccessKeyId = %v, want test-key", cfg.AccessKeyId)
	}
	if cfg.SecretKey != "test-secret" {
		t.Errorf("SecretKey = %v, want test-secret", cfg.SecretKey)
	}
}

func TestAwsConfig_Defaults(t *testing.T) {
	v := viper.New()
	// no sets

	var cfg AwsConfig
	err := v.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if cfg.Region != "" {
		t.Errorf("Region = %v, want empty", cfg.Region)
	}
	if cfg.AssumeRole != "" {
		t.Errorf("AssumeRole = %v, want empty", cfg.AssumeRole)
	}
	if cfg.AccessKeyId != "" {
		t.Errorf("AccessKeyId = %v, want empty", cfg.AccessKeyId)
	}
	if cfg.SecretKey != "" {
		t.Errorf("SecretKey = %v, want empty", cfg.SecretKey)
	}
}
