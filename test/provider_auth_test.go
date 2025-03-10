package integration

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/sistemica/traefik-manager/internal/api/server"
	"github.com/sistemica/traefik-manager/internal/config"
	"github.com/sistemica/traefik-manager/internal/store"
)

// TestProviderAuthIntegration tests the provider authentication functionality
func TestProviderAuthIntegration(t *testing.T) {
	// Skip in CI environment if needed
	if os.Getenv("CI") != "" && os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration tests in CI environment")
	}

	// Create a temporary file for the store
	tmpFile, err := os.CreateTemp("", "provider-auth-test-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	// Initialize with empty but valid JSON
	initialJSON := `{"middlewares":{},"routers":{},"services":{}}`
	if err := os.WriteFile(tmpPath, []byte(initialJSON), 0644); err != nil {
		t.Fatalf("Failed to initialize store file: %v", err)
	}

	// Test cases for different auth configurations
	testCases := []struct {
		name                string
		authEnabled         bool
		providerAuthEnabled bool
		apiKey              string
		providerKey         string
		requestKey          string
		expectedStatusCode  int
	}{
		{
			name:                "No Auth",
			authEnabled:         false,
			providerAuthEnabled: false,
			apiKey:              "",
			providerKey:         "",
			requestKey:          "",
			expectedStatusCode:  http.StatusOK, // Should allow access without auth
		},
		{
			name:                "Global Auth Only - Missing Key",
			authEnabled:         true,
			providerAuthEnabled: false,
			apiKey:              "api-key",
			providerKey:         "",
			requestKey:          "",
			expectedStatusCode:  http.StatusOK, // Provider endpoint should be excluded from global auth
		},
		{
			name:                "Provider Auth Only - Valid Key",
			authEnabled:         false,
			providerAuthEnabled: true,
			apiKey:              "",
			providerKey:         "provider-key",
			requestKey:          "provider-key",
			expectedStatusCode:  http.StatusOK,
		},
		{
			name:                "Provider Auth Only - Invalid Key",
			authEnabled:         false,
			providerAuthEnabled: true,
			apiKey:              "",
			providerKey:         "provider-key",
			requestKey:          "wrong-key",
			expectedStatusCode:  http.StatusUnauthorized,
		},
		{
			name:                "Provider Auth Only - Missing Key",
			authEnabled:         false,
			providerAuthEnabled: true,
			apiKey:              "",
			providerKey:         "provider-key",
			requestKey:          "",
			expectedStatusCode:  http.StatusUnauthorized,
		},
		{
			name:                "Both Auth Enabled - Provider Key",
			authEnabled:         true,
			providerAuthEnabled: true,
			apiKey:              "api-key",
			providerKey:         "provider-key",
			requestKey:          "provider-key",
			expectedStatusCode:  http.StatusOK,
		},
		{
			name:                "Both Auth Enabled - API Key",
			authEnabled:         true,
			providerAuthEnabled: true,
			apiKey:              "api-key",
			providerKey:         "provider-key",
			requestKey:          "api-key",               // Using API key instead of provider key
			expectedStatusCode:  http.StatusUnauthorized, // Should not work - provider auth takes precedence
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Start and end port for tests
			port := 10000 + (tc.expectedStatusCode % 100)

			// Create custom config for this test case
			cfg := &config.Config{
				Server: config.Server{
					Host:         "localhost",
					Port:         port,
					BasePath:     "/api/v1",
					ReadTimeout:  15 * time.Second,
					WriteTimeout: 15 * time.Second,
				},
				Storage: config.Storage{
					FilePath:     tmpPath,
					SaveInterval: 1 * time.Second,
				},
				Provider: config.Provider{
					ProviderPath: "/traefik/provider",
				},
				Logger: config.Logger{
					Level:     "error", // Minimize noise in tests
					Format:    "json",
					UseColors: false,
				},
				Cors: config.Cors{
					AllowedOrigins:   []string{"*"},
					AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
					AllowedHeaders:   []string{"Content-Type", "Authorization", "X-API-Key"},
					AllowCredentials: false,
					MaxAge:           300,
				},
				Auth: config.Auth{
					Enabled:    tc.authEnabled,
					HeaderName: "X-API-Key",
					Key:        tc.apiKey,
				},
			}

			// Add provider auth if enabled
			if tc.providerAuthEnabled {
				cfg.Provider.Auth = &config.Auth{
					Enabled:    true,
					HeaderName: "X-API-Key",
					Key:        tc.providerKey,
				}
			}

			// Create store
			dataStore, err := store.NewFileStore(tmpPath)
			if err != nil {
				t.Fatalf("Failed to create store: %v", err)
			}
			defer dataStore.Close()

			// Create server
			srv := server.New(cfg, dataStore)
			srv.Setup()

			// Get the Echo instance
			e := srv.GetEcho()

			// Create request to provider endpoint
			req := httptest.NewRequest(http.MethodGet, "/traefik/provider", nil)
			if tc.requestKey != "" {
				req.Header.Set("X-API-Key", tc.requestKey)
			}
			rec := httptest.NewRecorder()

			// Serve the request
			e.ServeHTTP(rec, req)

			// Check status code
			if rec.Code != tc.expectedStatusCode {
				t.Fatalf("Expected status code %d, got %d", tc.expectedStatusCode, rec.Code)
			}
		})
	}
}
