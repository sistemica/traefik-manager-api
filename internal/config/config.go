package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	Server   Server
	Storage  Storage
	Provider Provider
	Logger   Logger
	Cors     Cors
	Auth     Auth
}

type Server struct {
	Host string
	Port int
	// Base path for all API endpoints
	BasePath string
	// Read and write timeouts
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Storage struct {
	// Path to the storage file
	FilePath string
	// Debounce interval for saving changes
	SaveInterval time.Duration
}

type Provider struct {
	// Provider endpoint path
	ProviderPath string
	// Auth settings specific to the provider endpoint
	Auth *Auth
}

type Logger struct {
	// Log level (debug, info, warn, error)
	Level string
	// Log format (json, text)
	Format string
	// Path to log file, empty for stdout
	FilePath string
	// Use colors in console output
	UseColors bool
}

type Cors struct {
	// Allowed origins (* for all)
	AllowedOrigins []string
	// Allowed methods
	AllowedMethods []string
	// Allowed headers
	AllowedHeaders []string
	// Allow credentials
	AllowCredentials bool
	// Max age in seconds
	MaxAge int
}

type Auth struct {
	// Enable API key authentication
	Enabled bool
	// API key header name
	HeaderName string
	// API key value
	Key string
}

// LoadConfig loads the application configuration from environment variables
func LoadConfig(envFile string) (*Config, error) {
	// Load .env file if provided
	if envFile != "" {
		// Check if file exists
		if _, err := os.Stat(envFile); err == nil {
			if err := godotenv.Load(envFile); err != nil {
				return nil, fmt.Errorf("error loading .env file: %w", err)
			}
		}
	}

	// Load default .env file if exists in current directory
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			fmt.Printf("Warning: .env file found but could not be loaded: %v\n", err)
		}
	}

	config := &Config{}

	// Server configuration
	config.Server.Host = getEnv("SERVER_HOST", "")
	config.Server.Port = getEnvAsInt("SERVER_PORT", 9000)
	config.Server.BasePath = getEnv("SERVER_BASE_PATH", "/api/v1")
	config.Server.ReadTimeout = getEnvAsDuration("SERVER_READ_TIMEOUT", 15*time.Second)
	config.Server.WriteTimeout = getEnvAsDuration("SERVER_WRITE_TIMEOUT", 15*time.Second)

	// Storage configuration
	// Use system temp directory by default
	defaultStoragePath := filepath.Join(os.TempDir(), "traefik-manager.json")
	config.Storage.FilePath = getEnv("STORAGE_FILE_PATH", defaultStoragePath)

	// Ensure directory exists
	storageDir := filepath.Dir(config.Storage.FilePath)
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}
	config.Storage.SaveInterval = getEnvAsDuration("STORAGE_SAVE_INTERVAL", 5*time.Second)

	// Provider configuration
	config.Provider.ProviderPath = getEnv("PROVIDER_PATH", "/traefik/provider")

	// Provider-specific Auth
	providerAuthEnabled := getEnvAsBool("PROVIDER_AUTH_ENABLED", false)
	if providerAuthEnabled {
		config.Provider.Auth = &Auth{
			Enabled:    providerAuthEnabled,
			HeaderName: getEnv("PROVIDER_AUTH_HEADER_NAME", "X-API-Key"),
			Key:        getEnv("PROVIDER_AUTH_KEY", ""),
		}
	}

	// Logger configuration
	config.Logger.Level = getEnv("LOG_LEVEL", "info")
	config.Logger.Format = getEnv("LOG_FORMAT", "json")
	config.Logger.FilePath = getEnv("LOG_FILE_PATH", "")
	config.Logger.UseColors = getEnvAsBool("LOG_USE_COLORS", true)

	// CORS configuration
	config.Cors.AllowedOrigins = getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"*"})
	config.Cors.AllowedMethods = getEnvAsSlice("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	config.Cors.AllowedHeaders = getEnvAsSlice("CORS_ALLOWED_HEADERS", []string{"Content-Type", "Authorization"})
	config.Cors.AllowCredentials = getEnvAsBool("CORS_ALLOW_CREDENTIALS", false)
	config.Cors.MaxAge = getEnvAsInt("CORS_MAX_AGE", 300)

	// Auth configuration
	config.Auth.Enabled = getEnvAsBool("AUTH_ENABLED", false)
	config.Auth.HeaderName = getEnv("AUTH_HEADER_NAME", "X-API-Key")
	config.Auth.Key = getEnv("AUTH_KEY", "")

	// Validate required configuration
	if config.Auth.Enabled && config.Auth.Key == "" {
		return nil, fmt.Errorf("AUTH_KEY is required when AUTH_ENABLED is true")
	}

	// Validate provider auth if enabled
	if config.Provider.Auth != nil && config.Provider.Auth.Enabled && config.Provider.Auth.Key == "" {
		return nil, fmt.Errorf("PROVIDER_AUTH_KEY is required when PROVIDER_AUTH_ENABLED is true")
	}

	return config, nil
}

// Helper function to get an environment variable with a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Helper function to get an environment variable as int
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		fmt.Printf("Warning: Invalid value for %s, using default: %v\n", key, err)
		return defaultValue
	}
	return value
}

// Helper function to get an environment variable as bool
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		fmt.Printf("Warning: Invalid value for %s, using default: %v\n", key, err)
		return defaultValue
	}
	return value
}

// Helper function to get an environment variable as duration
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		fmt.Printf("Warning: Invalid value for %s, using default: %v\n", key, err)
		return defaultValue
	}
	return value
}

// Helper function to get an environment variable as string slice
func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	return strings.Split(valueStr, ",")
}
