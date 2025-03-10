package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	// Create a temp directory for testing
	tempDir := t.TempDir()

	// Test default configuration
	t.Run("Default Config", func(t *testing.T) {
		// Clear environment variables that might affect the test
		os.Unsetenv("SERVER_HOST")
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("SERVER_BASE_PATH")
		os.Unsetenv("STORAGE_FILE_PATH")
		os.Unsetenv("PROVIDER_PATH")
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("AUTH_ENABLED")
		os.Unsetenv("PROVIDER_AUTH_ENABLED")

		cfg, err := LoadConfig("")
		if err != nil {
			t.Fatalf("Failed to load default config: %v", err)
		}

		// Verify default values
		if cfg.Server.Port != 9000 {
			t.Errorf("Expected default port 9000, got %d", cfg.Server.Port)
		}

		if cfg.Server.BasePath != "/api/v1" {
			t.Errorf("Expected default base path '/api/v1', got '%s'", cfg.Server.BasePath)
		}

		if cfg.Provider.ProviderPath != "/traefik/provider" {
			t.Errorf("Expected default provider path '/traefik/provider', got '%s'", cfg.Provider.ProviderPath)
		}

		if cfg.Logger.Level != "info" {
			t.Errorf("Expected default log level 'info', got '%s'", cfg.Logger.Level)
		}

		if cfg.Auth.Enabled {
			t.Errorf("Expected auth to be disabled by default")
		}

		// Verify storage path is in system temp directory
		if !strings.HasPrefix(cfg.Storage.FilePath, os.TempDir()) {
			t.Errorf("Expected storage path to be in temp directory, got '%s'", cfg.Storage.FilePath)
		}
	})

	// Test configuration from environment variables
	t.Run("Environment Config", func(t *testing.T) {
		// Set up a temporary file path for storage
		tmpFilePath := filepath.Join(tempDir, "test-data.json")

		// Set environment variables
		os.Setenv("SERVER_HOST", "localhost")
		os.Setenv("SERVER_PORT", "8080")
		os.Setenv("SERVER_BASE_PATH", "/api/v2")
		os.Setenv("STORAGE_FILE_PATH", tmpFilePath)
		os.Setenv("STORAGE_SAVE_INTERVAL", "10s")
		os.Setenv("PROVIDER_PATH", "/custom/provider")
		os.Setenv("LOG_LEVEL", "debug")
		os.Setenv("LOG_FORMAT", "text")
		os.Setenv("AUTH_ENABLED", "true")
		os.Setenv("AUTH_KEY", "test-api-key")

		defer func() {
			// Clean up environment variables after test
			os.Unsetenv("SERVER_HOST")
			os.Unsetenv("SERVER_PORT")
			os.Unsetenv("SERVER_BASE_PATH")
			os.Unsetenv("STORAGE_FILE_PATH")
			os.Unsetenv("STORAGE_SAVE_INTERVAL")
			os.Unsetenv("PROVIDER_PATH")
			os.Unsetenv("LOG_LEVEL")
			os.Unsetenv("LOG_FORMAT")
			os.Unsetenv("AUTH_ENABLED")
			os.Unsetenv("AUTH_KEY")
		}()

		cfg, err := LoadConfig("")
		if err != nil {
			t.Fatalf("Failed to load config from environment: %v", err)
		}

		// Verify values from environment
		if cfg.Server.Host != "localhost" {
			t.Errorf("Expected host 'localhost', got '%s'", cfg.Server.Host)
		}

		if cfg.Server.Port != 8080 {
			t.Errorf("Expected port 8080, got %d", cfg.Server.Port)
		}

		if cfg.Server.BasePath != "/api/v2" {
			t.Errorf("Expected base path '/api/v2', got '%s'", cfg.Server.BasePath)
		}

		if cfg.Storage.FilePath != tmpFilePath {
			t.Errorf("Expected storage file path '%s', got '%s'", tmpFilePath, cfg.Storage.FilePath)
		}

		if cfg.Storage.SaveInterval != 10*time.Second {
			t.Errorf("Expected save interval '10s', got '%v'", cfg.Storage.SaveInterval)
		}

		if cfg.Provider.ProviderPath != "/custom/provider" {
			t.Errorf("Expected provider path '/custom/provider', got '%s'", cfg.Provider.ProviderPath)
		}

		if cfg.Logger.Level != "debug" {
			t.Errorf("Expected log level 'debug', got '%s'", cfg.Logger.Level)
		}

		if cfg.Logger.Format != "text" {
			t.Errorf("Expected log format 'text', got '%s'", cfg.Logger.Format)
		}

		if !cfg.Auth.Enabled {
			t.Errorf("Expected auth to be enabled")
		}

		if cfg.Auth.Key != "test-api-key" {
			t.Errorf("Expected auth key 'test-api-key', got '%s'", cfg.Auth.Key)
		}
	})

	// Test provider-specific auth configuration
	t.Run("Provider Auth Config", func(t *testing.T) {
		// Clear previous environment
		os.Unsetenv("SERVER_HOST")
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("SERVER_BASE_PATH")
		os.Unsetenv("STORAGE_FILE_PATH")
		os.Unsetenv("STORAGE_SAVE_INTERVAL")
		os.Unsetenv("PROVIDER_PATH")
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("LOG_FORMAT")
		os.Unsetenv("AUTH_ENABLED")
		os.Unsetenv("AUTH_KEY")

		// Set provider auth environment variables
		os.Setenv("PROVIDER_AUTH_ENABLED", "true")
		os.Setenv("PROVIDER_AUTH_HEADER_NAME", "X-Provider-Key")
		os.Setenv("PROVIDER_AUTH_KEY", "provider-key-123")

		defer func() {
			// Clean up environment variables after test
			os.Unsetenv("PROVIDER_AUTH_ENABLED")
			os.Unsetenv("PROVIDER_AUTH_HEADER_NAME")
			os.Unsetenv("PROVIDER_AUTH_KEY")
		}()

		cfg, err := LoadConfig("")
		if err != nil {
			t.Fatalf("Failed to load config with provider auth: %v", err)
		}

		// Verify provider auth configuration
		if cfg.Provider.Auth == nil {
			t.Fatalf("Expected provider auth to be non-nil")
		}

		if !cfg.Provider.Auth.Enabled {
			t.Errorf("Expected provider auth to be enabled")
		}

		if cfg.Provider.Auth.HeaderName != "X-Provider-Key" {
			t.Errorf("Expected provider auth header name 'X-Provider-Key', got '%s'", cfg.Provider.Auth.HeaderName)
		}

		if cfg.Provider.Auth.Key != "provider-key-123" {
			t.Errorf("Expected provider auth key 'provider-key-123', got '%s'", cfg.Provider.Auth.Key)
		}
	})

	// Test invalid configuration validation
	t.Run("Invalid Config Validation", func(t *testing.T) {
		// Set invalid configuration (auth enabled but no key)
		os.Setenv("AUTH_ENABLED", "true")
		os.Setenv("AUTH_KEY", "")

		defer func() {
			os.Unsetenv("AUTH_ENABLED")
			os.Unsetenv("AUTH_KEY")
		}()

		_, err := LoadConfig("")
		if err == nil {
			t.Fatalf("Expected error for invalid auth config, but got nil")
		}
	})

	// Test invalid provider auth configuration validation
	t.Run("Invalid Provider Auth Config Validation", func(t *testing.T) {
		// Set invalid configuration (provider auth enabled but no key)
		os.Setenv("PROVIDER_AUTH_ENABLED", "true")
		os.Setenv("PROVIDER_AUTH_KEY", "")

		defer func() {
			os.Unsetenv("PROVIDER_AUTH_ENABLED")
			os.Unsetenv("PROVIDER_AUTH_KEY")
		}()

		_, err := LoadConfig("")
		if err == nil {
			t.Fatalf("Expected error for invalid provider auth config, but got nil")
		}
	})

	// Test helper functions
	t.Run("Helper Functions", func(t *testing.T) {
		// Test getEnvAsInt
		os.Setenv("TEST_INT", "123")
		defer os.Unsetenv("TEST_INT")
		if getEnvAsInt("TEST_INT", 0) != 123 {
			t.Errorf("getEnvAsInt failed to parse '123' as int")
		}

		// Test with invalid int
		os.Setenv("TEST_INVALID_INT", "abc")
		defer os.Unsetenv("TEST_INVALID_INT")
		if getEnvAsInt("TEST_INVALID_INT", 42) != 42 {
			t.Errorf("getEnvAsInt should return default for invalid int")
		}

		// Test getEnvAsBool
		os.Setenv("TEST_BOOL_TRUE", "true")
		defer os.Unsetenv("TEST_BOOL_TRUE")
		if !getEnvAsBool("TEST_BOOL_TRUE", false) {
			t.Errorf("getEnvAsBool failed to parse 'true' as bool")
		}

		os.Setenv("TEST_BOOL_FALSE", "false")
		defer os.Unsetenv("TEST_BOOL_FALSE")
		if getEnvAsBool("TEST_BOOL_FALSE", true) {
			t.Errorf("getEnvAsBool failed to parse 'false' as bool")
		}

		// Test getEnvAsDuration
		os.Setenv("TEST_DURATION", "5s")
		defer os.Unsetenv("TEST_DURATION")
		if getEnvAsDuration("TEST_DURATION", time.Second) != 5*time.Second {
			t.Errorf("getEnvAsDuration failed to parse '5s' as duration")
		}

		// Test getEnvAsSlice
		os.Setenv("TEST_SLICE", "a,b,c")
		defer os.Unsetenv("TEST_SLICE")
		slice := getEnvAsSlice("TEST_SLICE", []string{})
		if len(slice) != 3 || slice[0] != "a" || slice[1] != "b" || slice[2] != "c" {
			t.Errorf("getEnvAsSlice failed to parse 'a,b,c' as slice")
		}
	})
}
