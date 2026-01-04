// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package storage_files

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/configs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalFileStorage_Name(t *testing.T) {
	cfg := configs.AssetStoreConfig{
		StorageType:       "local",
		StoragePathPrefix: "/tmp/test",
	}
	logger, _ := commons.NewApplicationLogger()
	storage := NewLocalFileStorage(cfg, logger)

	assert.Equal(t, "local", storage.Name())
}

func TestLocalFileStorage_Store(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "local_storage_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	cfg := configs.AssetStoreConfig{
		StorageType:       "local",
		StoragePathPrefix: tempDir,
	}
	logger, _ := commons.NewApplicationLogger()
	storage := NewLocalFileStorage(cfg, logger)

	ctx := context.Background()
	key := "test/file.txt"
	content := []byte("Hello, World!")

	result := storage.Store(ctx, key, content)

	assert.NoError(t, result.Error)
	assert.Equal(t, configs.LOCAL, result.StorageType)
	assert.Equal(t, filepath.Join(tempDir, key), result.CompletePath)

	// Verify file was created
	filePath := filepath.Join(tempDir, key)
	assert.FileExists(t, filePath)

	// Verify content
	storedContent, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Equal(t, content, storedContent)
}

func TestLocalFileStorage_Store_NestedDirectories(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "local_storage_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	cfg := configs.AssetStoreConfig{
		StorageType:       "local",
		StoragePathPrefix: tempDir,
	}
	logger, _ := commons.NewApplicationLogger()
	storage := NewLocalFileStorage(cfg, logger)

	ctx := context.Background()
	key := "deep/nested/path/file.txt"
	content := []byte("Nested file content")

	result := storage.Store(ctx, key, content)

	assert.NoError(t, result.Error)
	assert.Equal(t, configs.LOCAL, result.StorageType)
	assert.Equal(t, filepath.Join(tempDir, key), result.CompletePath)

	// Verify file was created
	filePath := filepath.Join(tempDir, key)
	assert.FileExists(t, filePath)

	// Verify content
	storedContent, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Equal(t, content, storedContent)
}

func TestLocalFileStorage_Get(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "local_storage_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	cfg := configs.AssetStoreConfig{
		StorageType:       "local",
		StoragePathPrefix: tempDir,
	}
	logger, _ := commons.NewApplicationLogger()
	storage := NewLocalFileStorage(cfg, logger)

	ctx := context.Background()
	key := "test/file.txt"
	content := []byte("Hello, World!")

	// First store the file
	storeResult := storage.Store(ctx, key, content)
	require.NoError(t, storeResult.Error)

	// Now get it back
	getResult := storage.Get(ctx, key)

	assert.NoError(t, getResult.Error)
	assert.Equal(t, content, getResult.Data)
}

func TestLocalFileStorage_Get_FileNotExists(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "local_storage_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	cfg := configs.AssetStoreConfig{
		StorageType:       "local",
		StoragePathPrefix: tempDir,
	}
	logger, _ := commons.NewApplicationLogger()
	storage := NewLocalFileStorage(cfg, logger)

	ctx := context.Background()
	key := "nonexistent/file.txt"

	getResult := storage.Get(ctx, key)

	assert.Error(t, getResult.Error)
	assert.Nil(t, getResult.Data)
}

func TestLocalFileStorage_GetUrl(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "local_storage_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	cfg := configs.AssetStoreConfig{
		StorageType:       "local",
		StoragePathPrefix: tempDir,
	}
	logger, _ := commons.NewApplicationLogger()
	storage := NewLocalFileStorage(cfg, logger)

	ctx := context.Background()
	key := "test/file.txt"

	result := storage.GetUrl(ctx, key)

	assert.NoError(t, result.Error)
	assert.Equal(t, configs.LOCAL, result.StorageType)
	expectedPath := filepath.Join("file://", tempDir, key)
	assert.Equal(t, expectedPath, result.CompletePath)
}
