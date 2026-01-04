// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package storage_files

import (
	"testing"

	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/configs"
	"github.com/stretchr/testify/assert"
)

func TestNewStorage_S3(t *testing.T) {
	cfg := configs.AssetStoreConfig{
		StorageType:       "s3",
		StoragePathPrefix: "test-bucket",
		Auth: &configs.AwsConfig{
			Region:      "us-east-1",
			AccessKeyId: "test-key",
			SecretKey:   "test-secret",
		},
	}
	logger, _ := commons.NewApplicationLogger()

	storage := NewStorage(cfg, logger)

	assert.NotNil(t, storage)
	assert.Equal(t, "aws", storage.Name())
}

func TestNewStorage_Local(t *testing.T) {
	cfg := configs.AssetStoreConfig{
		StorageType:       "local",
		StoragePathPrefix: "/tmp/test",
	}
	logger, _ := commons.NewApplicationLogger()

	storage := NewStorage(cfg, logger)

	assert.NotNil(t, storage)
	assert.Equal(t, "local", storage.Name())
}

func TestNewStorage_CDN(t *testing.T) {
	cfg := configs.AssetStoreConfig{
		StorageType:       "cdn",
		StoragePathPrefix: "https://cdn.example.com",
		Auth: &configs.AwsConfig{
			Region:      "us-east-1",
			AccessKeyId: "test-key",
			SecretKey:   "test-secret",
		},
	}
	logger, _ := commons.NewApplicationLogger()

	storage := NewStorage(cfg, logger)

	assert.NotNil(t, storage)
	assert.Equal(t, "cdn", storage.Name())
}

func TestNewStorage_UnsupportedType(t *testing.T) {
	cfg := configs.AssetStoreConfig{
		StorageType:       "unsupported",
		StoragePathPrefix: "/tmp/test",
	}
	logger, _ := commons.NewApplicationLogger()

	storage := NewStorage(cfg, logger)

	assert.NotNil(t, storage)
	// Should default to local storage for unsupported types
	assert.Equal(t, "local", storage.Name())
}
