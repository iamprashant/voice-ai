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

func TestAssetStoreConfig_Type(t *testing.T) {
	tests := []struct {
		name string
		cfg  AssetStoreConfig
		want StorageType
	}{
		{"S3", AssetStoreConfig{StorageType: "s3"}, S3},
		{"Local", AssetStoreConfig{StorageType: "local"}, LOCAL},
		{"CDN", AssetStoreConfig{StorageType: "cdn"}, LOCAL},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cfg.Type(); got != tt.want {
				t.Errorf("AssetStoreConfig.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAssetStoreConfig_Validation(t *testing.T) {
	validate := validator.New()
	tests := []struct {
		name    string
		cfg     AssetStoreConfig
		wantErr bool
	}{
		{"valid s3", AssetStoreConfig{StorageType: "s3"}, false},
		{"valid local", AssetStoreConfig{StorageType: "local"}, false},
		{"invalid empty", AssetStoreConfig{}, true},
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

func TestAssetStoreConfig_FromViper(t *testing.T) {
	v := viper.New()
	v.Set("storage_type", "s3")
	v.Set("storage_path_prefix", "test_prefix")

	var cfg AssetStoreConfig
	err := v.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if cfg.StorageType != "s3" {
		t.Errorf("StorageType = %v, want s3", cfg.StorageType)
	}
	if cfg.StoragePathPrefix != "test_prefix" {
		t.Errorf("StoragePathPrefix = %v, want test_prefix", cfg.StoragePathPrefix)
	}
	if cfg.Auth != nil {
		t.Errorf("Auth should be nil")
	}
}

func TestAssetStoreConfig_Defaults(t *testing.T) {
	v := viper.New()
	// no sets, so defaults

	var cfg AssetStoreConfig
	err := v.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if cfg.StorageType != "" {
		t.Errorf("StorageType = %v, want empty", cfg.StorageType)
	}
	if cfg.StoragePathPrefix != "" {
		t.Errorf("StoragePathPrefix = %v, want empty", cfg.StoragePathPrefix)
	}
	if cfg.Type() != LOCAL {
		t.Errorf("Type() = %v, want LOCAL", cfg.Type())
	}
	if cfg.Auth != nil {
		t.Errorf("Auth should be nil")
	}
}
