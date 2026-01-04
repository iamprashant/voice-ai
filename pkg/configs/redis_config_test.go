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

func TestRedisConfig_Validation(t *testing.T) {
	validate := validator.New()
	tests := []struct {
		name    string
		cfg     RedisConfig
		wantErr bool
	}{
		{"valid", RedisConfig{Host: "localhost", Port: 6379, MaxConnection: 10}, false},
		{"invalid missing host", RedisConfig{Port: 6379, MaxConnection: 10}, true},
		{"invalid missing port", RedisConfig{Host: "localhost", MaxConnection: 10}, true},
		{"invalid missing max_connection", RedisConfig{Host: "localhost", Port: 6379}, true},
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

func TestRedisConfig_FromViper(t *testing.T) {
	v := viper.New()
	v.Set("host", "redis.example.com")
	v.Set("port", 6380)
	v.Set("db", 1)
	v.Set("max_connection", 50)
	v.Set("auth.user", "redisuser")
	v.Set("auth.password", "redispass")
	v.Set("insecure_skip_verify", true)

	var cfg RedisConfig
	err := v.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if cfg.Host != "redis.example.com" {
		t.Errorf("Host = %v, want redis.example.com", cfg.Host)
	}
	if cfg.Port != 6380 {
		t.Errorf("Port = %v, want 6380", cfg.Port)
	}
	if cfg.Db != 1 {
		t.Errorf("Db = %v, want 1", cfg.Db)
	}
	if cfg.MaxConnection != 50 {
		t.Errorf("MaxConnection = %v, want 50", cfg.MaxConnection)
	}
	if cfg.Auth.User != "redisuser" {
		t.Errorf("Auth.User = %v, want redisuser", cfg.Auth.User)
	}
	if cfg.Auth.Password != "redispass" {
		t.Errorf("Auth.Password = %v, want redispass", cfg.Auth.Password)
	}
	if !cfg.InsecureSkipVerify {
		t.Errorf("InsecureSkipVerify = %v, want true", cfg.InsecureSkipVerify)
	}
}

func TestRedisConfig_Defaults(t *testing.T) {
	v := viper.New()
	// no sets

	var cfg RedisConfig
	err := v.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if cfg.Host != "" {
		t.Errorf("Host = %v, want empty", cfg.Host)
	}
	if cfg.Port != 0 {
		t.Errorf("Port = %v, want 0", cfg.Port)
	}
	if cfg.Db != 0 {
		t.Errorf("Db = %v, want 0", cfg.Db)
	}
	if cfg.MaxConnection != 0 {
		t.Errorf("MaxConnection = %v, want 0", cfg.MaxConnection)
	}
	if cfg.Auth.User != "" {
		t.Errorf("Auth.User = %v, want empty", cfg.Auth.User)
	}
	if cfg.Auth.Password != "" {
		t.Errorf("Auth.Password = %v, want empty", cfg.Auth.Password)
	}
	if cfg.InsecureSkipVerify {
		t.Errorf("InsecureSkipVerify = %v, want false", cfg.InsecureSkipVerify)
	}
}
