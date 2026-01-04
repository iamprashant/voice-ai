// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestInitConfig(t *testing.T) {
	// Mock Environment setup
	envPath := filepath.Join(os.TempDir(), ".endpointtest.env")
	err := os.WriteFile(envPath, []byte(`
		SERVICE_NAME="endpoint-api"
		HOST="0.0.0.0"
		PORT=9005
		LOG_LEVEL="debug"
		SECRET="rpd_pks"
		ENV="development"
		
		POSTGRES__HOST="localhost"
		POSTGRES__DB_NAME="endpoint_db"
		POSTGRES__AUTH__USER="rapida_user"
		POSTGRES__AUTH__PASSWORD="rapida_db_password"
		POSTGRES__PORT=5432
		POSTGRES__MAX_OPEN_CONNECTION=10
		POSTGRES__MAX_IDEAL_CONNECTION=10
		POSTGRES__SSL_MODE="disable"

		REDIS__HOST="localhost"
		REDIS__PORT="6379"
		REDIS__MAX_CONNECTION=5

		ASSET_STORE__STORAGE_TYPE="local"
		ASSET_STORE__STORAGE_PATH_PREFIX=${HOME}/rapida-data/assets/endpoint
		
		INTEGRATION_HOST=localhost:9004
		ENDPOINT_HOST=localhost:9005
		ASSISTANT_HOST=localhost:9007
		WEB_HOST=localhost:9001
		DOCUMENT_HOST=http://localhost:9010
		UI_HOST=http://localhost:3000
	`), 0644)
	if err != nil {
		t.Fatalf("Failed to create mock env file: %v", err)
	}
	defer os.Remove(envPath) // Clean up after test

	os.Setenv("ENV_PATH", envPath)
	defer os.Unsetenv("ENV_PATH") // Reset ENV_PATH after test

	// Test initializing configuration
	vConfig, err := InitConfig()
	if err != nil {
		t.Fatalf("InitConfig returned an error: %v", err)
	}
	if vConfig == nil {
		t.Fatalf("vConfig is nil")
	}

	// Verify parameters
	if vConfig.GetString("SERVICE_NAME") != "endpoint-api" {
		t.Errorf("Expected SERVICE_NAME to be 'endpoint-api', but got %v", vConfig.GetString("SERVICE_NAME"))
	}
	if vConfig.GetString("POSTGRES__DB_NAME") != "endpoint_db" {
		t.Errorf("Expected POSTGRES__DB_NAME to be 'endpoint_db', but got %v", vConfig.GetString("POSTGRES__DB_NAME"))
	}

	// Add more parameter validations here as necessary...
}

// ... existing code ...

func TestGetApplicationConfig(t *testing.T) {
	vConfig := viper.NewWithOptions(viper.KeyDelimiter("__"))
	vConfig.Set("SERVICE_NAME", "endpoint-api")
	vConfig.Set("HOST", "0.0.0.0")
	vConfig.Set("PORT", 9005)
	vConfig.Set("LOG_LEVEL", "debug")
	vConfig.Set("SECRET", "rpd_pks")
	vConfig.Set("ENV", "development")

	vConfig.Set("POSTGRES__HOST", "localhost")
	vConfig.Set("POSTGRES__DB_NAME", "endpoint_db")
	vConfig.Set("POSTGRES__AUTH__USER", "rapida_user")
	vConfig.Set("POSTGRES__AUTH__PASSWORD", "rapida_db_password")
	vConfig.Set("POSTGRES__PORT", 5432)
	vConfig.Set("POSTGRES__MAX_OPEN_CONNECTION", 10)
	vConfig.Set("POSTGRES__MAX_IDEAL_CONNECTION", 10)
	vConfig.Set("POSTGRES__SSL_MODE", "disable")

	vConfig.Set("REDIS__HOST", "localhost")
	vConfig.Set("REDIS__PORT", "6379")
	vConfig.Set("REDIS__MAX_CONNECTION", 5)

	vConfig.Set("ASSET_STORE__STORAGE_TYPE", "local")
	vConfig.Set("ASSET_STORE__STORAGE_PATH_PREFIX", os.Getenv("HOME")+"/rapida-data/assets/endpoint")

	vConfig.Set("INTEGRATION_HOST", "localhost:9004")
	vConfig.Set("ENDPOINT_HOST", "localhost:9005")
	vConfig.Set("ASSISTANT_HOST", "localhost:9007")
	vConfig.Set("WEB_HOST", "localhost:9001")
	vConfig.Set("DOCUMENT_HOST", "http://localhost:9010")
	vConfig.Set("UI_HOST", "http://localhost:3000")

	appConfig, err := GetApplicationConfig(vConfig)
	if err != nil {
		t.Fatalf("GetApplicationConfig returned an error: %v", err)
	}
	if appConfig == nil {
		t.Fatalf("appConfig is nil")
	}

	// Validate new configurations
	if appConfig.PostgresConfig.DBName != "endpoint_db" {
		t.Errorf("Expected PostgresConfig.DBName to be 'endpoint_db', but got %v", appConfig.PostgresConfig.DBName)
	}
	if appConfig.AssetStoreConfig.StorageType != "local" {
		t.Errorf("Expected AssetStoreConfig.StorageType to be 'local', but got %v", appConfig.AssetStoreConfig.StorageType)
	}
	if appConfig.RedisConfig.Host != "localhost" || appConfig.RedisConfig.Port != 6379 {
		t.Errorf("Redis Config mismatch: Host=%v, Port=%v", appConfig.RedisConfig.Host, appConfig.RedisConfig.Port)
	}

}
