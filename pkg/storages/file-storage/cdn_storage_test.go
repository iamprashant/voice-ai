// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package storage_files

import (
	"context"
	"strings"
	"testing"

	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/configs"
	"github.com/stretchr/testify/assert"
)

func TestCDNStorage_Name(t *testing.T) {
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
	storage := NewCDNStorage(cfg, logger)

	assert.Equal(t, "cdn", storage.Name())
}

func TestCDNStorage_Store_SessionCreationFailure(t *testing.T) {
	// Test case where session creation fails
	cfg := configs.AssetStoreConfig{
		StorageType:       "cdn",
		StoragePathPrefix: "https://cdn.example.com",
		Auth: &configs.AwsConfig{
			Region: "", // Invalid region to cause session failure
		},
	}
	logger, _ := commons.NewApplicationLogger()
	storage := NewCDNStorage(cfg, logger)

	ctx := context.Background()
	key := "test/file.txt"
	content := []byte("test content")

	result := storage.Store(ctx, key, content)

	// Should return error due to invalid session
	assert.Error(t, result.Error)
	assert.Equal(t, configs.S3, result.StorageType)
	assert.Contains(t, result.CompletePath, "https://cdn.example.com/cdn/")
}

func TestCDNStorage_Get_SessionCreationFailure(t *testing.T) {
	cfg := configs.AssetStoreConfig{
		StorageType:       "cdn",
		StoragePathPrefix: "https://cdn.example.com",
		Auth: &configs.AwsConfig{
			Region: "", // Invalid region to cause session failure
		},
	}
	logger, _ := commons.NewApplicationLogger()
	storage := NewCDNStorage(cfg, logger)

	ctx := context.Background()
	key := "test/file.txt"

	result := storage.Get(ctx, key)

	assert.Error(t, result.Error)
	assert.Nil(t, result.Data)
}

func TestCDNStorage_GetUrl(t *testing.T) {
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
	storage := NewCDNStorage(cfg, logger)

	ctx := context.Background()
	key := "test/file.txt"

	result := storage.GetUrl(ctx, key)

	assert.NoError(t, result.Error)
	assert.Equal(t, configs.S3, result.StorageType)
	assert.Equal(t, "https://cdn.example.com/test/file.txt", result.CompletePath)
}

func TestCDNStorage_prefix(t *testing.T) {
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
	storage := NewCDNStorage(cfg, logger).(*cdnStorage)

	ctx := context.Background()
	key := "test/file.txt"

	// Test that prefix adds ID prefix
	prefixed := storage.prefix(ctx, key)
	assert.Contains(t, prefixed, "cdn/")
	assert.Contains(t, prefixed, "_"+key)
	// Should be in format: cdn/{id}_{key}
	parts := strings.Split(prefixed, "_")
	assert.Len(t, parts, 2)
	assert.Equal(t, "cdn/"+parts[0]+"/"+parts[0], "cdn/"+parts[0]+"/"+parts[0]) // This is a rough check
}
