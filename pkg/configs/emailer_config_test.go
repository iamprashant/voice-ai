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

func TestEmailerConfig_Provider(t *testing.T) {
	tests := []struct {
		name string
		cfg  EmailerConfig
		want EmailProvider
	}{
		{"SES", EmailerConfig{EmailProvider: "ses"}, SES},
		{"Sendgrid", EmailerConfig{EmailProvider: "sendgrid"}, SENDGRID},
		{"Default", EmailerConfig{EmailProvider: "other"}, SENDGRID},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cfg.Provider(); got != tt.want {
				t.Errorf("EmailerConfig.Provider() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmailerConfig_Validation(t *testing.T) {
	validate := validator.New()
	tests := []struct {
		name    string
		cfg     EmailerConfig
		wantErr bool
	}{
		{"valid", EmailerConfig{EmailProvider: "ses", FromEmail: "test@example.com", FromName: "Test"}, false},
		{"invalid missing provider", EmailerConfig{FromEmail: "test@example.com", FromName: "Test"}, true},
		{"invalid missing from_email", EmailerConfig{EmailProvider: "ses", FromName: "Test"}, true},
		{"invalid missing from_name", EmailerConfig{EmailProvider: "ses", FromEmail: "test@example.com"}, true},
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

func TestEmailerConfig_FromViper(t *testing.T) {
	v := viper.New()
	v.Set("provider", "sendgrid")
	v.Set("from_email", "noreply@example.com")
	v.Set("from_name", "No Reply")
	sendgridKey := "test-key"
	v.Set("sendgrid_key", sendgridKey)

	var cfg EmailerConfig
	err := v.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if cfg.EmailProvider != "sendgrid" {
		t.Errorf("EmailProvider = %v, want sendgrid", cfg.EmailProvider)
	}
	if cfg.FromEmail != "noreply@example.com" {
		t.Errorf("FromEmail = %v, want noreply@example.com", cfg.FromEmail)
	}
	if cfg.FromName != "No Reply" {
		t.Errorf("FromName = %v, want No Reply", cfg.FromName)
	}
	if cfg.SendgridKey == nil || *cfg.SendgridKey != "test-key" {
		t.Errorf("SendgridKey = %v, want test-key", cfg.SendgridKey)
	}
	if cfg.Auth != nil {
		t.Errorf("Auth should be nil")
	}
}

func TestEmailerConfig_Defaults(t *testing.T) {
	v := viper.New()
	// no sets

	var cfg EmailerConfig
	err := v.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if cfg.EmailProvider != "" {
		t.Errorf("EmailProvider = %v, want empty", cfg.EmailProvider)
	}
	if cfg.FromEmail != "" {
		t.Errorf("FromEmail = %v, want empty", cfg.FromEmail)
	}
	if cfg.FromName != "" {
		t.Errorf("FromName = %v, want empty", cfg.FromName)
	}
	if cfg.SendgridKey != nil {
		t.Errorf("SendgridKey should be nil")
	}
	if cfg.Auth != nil {
		t.Errorf("Auth should be nil")
	}
	if cfg.Provider() != SENDGRID {
		t.Errorf("Provider() = %v, want SENDGRID", cfg.Provider())
	}
}
