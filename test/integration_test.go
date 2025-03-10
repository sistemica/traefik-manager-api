package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/sistemica/traefik-manager/internal/api/server"
	"github.com/sistemica/traefik-manager/internal/config"
	"github.com/sistemica/traefik-manager/internal/models"
	"github.com/sistemica/traefik-manager/internal/store"
)

// TestAPIIntegration runs integration tests against the API
func TestAPIIntegration(t *testing.T) {
	// Skip in CI environment if needed
	if os.Getenv("CI") != "" && os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration tests in CI environment")
	}

	// Create temporary file for store
	// Find where you create the temporary file:
	tmpFile, err := os.CreateTemp("", "traefik-manager-integration-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	// After this, add this code to initialize the file with empty valid JSON:
	// Initialize the file with empty but valid JSON
	initialJSON := `{"middlewares":{},"routers":{},"services":{}}`
	if err := os.WriteFile(tmpPath, []byte(initialJSON), 0644); err != nil {
		t.Fatalf("Failed to initialize store file: %v", err)
	}

	// Create test configuration
	cfg := &config.Config{
		Server: config.Server{
			Host:         "localhost",
			Port:         9999, // Use a different port for testing
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
			Level:     "info",
			Format:    "text",
			UseColors: true,
		},
		Cors: config.Cors{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Content-Type", "Authorization"},
			AllowCredentials: false,
			MaxAge:           300,
		},
		Auth: config.Auth{
			Enabled:    false,
			HeaderName: "X-API-Key",
			Key:        "",
		},
	}

	// Create store
	dataStore, err := store.NewFileStore(tmpPath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// Create server
	srv := server.New(cfg, dataStore)
	srv.Setup()

	// Start server in a goroutine
	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			t.Errorf("Server error: %v", err)
		}
	}()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Base URL for API requests
	baseURL := fmt.Sprintf("http://%s:%d%s", cfg.Server.Host, cfg.Server.Port, cfg.Server.BasePath)

	// Create HTTP client
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Test health endpoint
	t.Run("Health Check", func(t *testing.T) {
		resp, err := client.Get(fmt.Sprintf("http://%s:%d/api/v1/health", cfg.Server.Host, cfg.Server.Port))
		if err != nil {
			t.Fatalf("Failed to call health endpoint: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}

		var health map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if health["status"] != "healthy" {
			t.Fatalf("Expected status 'healthy', got '%v'", health["status"])
		}
	})

	// Test middleware endpoints
	t.Run("Middleware CRUD", func(t *testing.T) {
		middlewareID := "test-middleware"
		middlewareType := "redirectScheme"

		// Create middleware
		createBody := map[string]interface{}{
			"id":   middlewareID,
			"type": middlewareType,
			"config": map[string]interface{}{
				"scheme":    "https",
				"permanent": true,
			},
		}
		createJSON, _ := json.Marshal(createBody)

		resp, err := client.Post(
			fmt.Sprintf("%s/middlewares", baseURL),
			"application/json",
			bytes.NewBuffer(createJSON),
		)
		if err != nil {
			t.Fatalf("Failed to create middleware: %v", err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("Expected status code %d, got %d", http.StatusCreated, resp.StatusCode)
		}

		// Get middleware
		resp, err = client.Get(fmt.Sprintf("%s/middlewares/%s", baseURL, middlewareID))
		if err != nil {
			t.Fatalf("Failed to get middleware: %v", err)
		}

		var middleware models.Middleware
		if err := json.NewDecoder(resp.Body).Decode(&middleware); err != nil {
			t.Fatalf("Failed to decode middleware: %v", err)
		}
		resp.Body.Close()

		if middleware.ID != middlewareID {
			t.Fatalf("Expected middleware ID '%s', got '%s'", middlewareID, middleware.ID)
		}
		if middleware.Type != middlewareType {
			t.Fatalf("Expected middleware type '%s', got '%s'", middlewareType, middleware.Type)
		}

		// List middlewares
		resp, err = client.Get(fmt.Sprintf("%s/middlewares", baseURL))
		if err != nil {
			t.Fatalf("Failed to list middlewares: %v", err)
		}

		var middlewares []models.Middleware
		if err := json.NewDecoder(resp.Body).Decode(&middlewares); err != nil {
			t.Fatalf("Failed to decode middlewares: %v", err)
		}
		resp.Body.Close()

		if len(middlewares) != 1 {
			t.Fatalf("Expected 1 middleware, got %d", len(middlewares))
		}

		// Update middleware
		updateBody := map[string]interface{}{
			"type": middlewareType,
			"config": map[string]interface{}{
				"scheme":    "http",
				"permanent": false,
			},
		}
		updateJSON, _ := json.Marshal(updateBody)

		req, _ := http.NewRequest(
			http.MethodPut,
			fmt.Sprintf("%s/middlewares/%s", baseURL, middlewareID),
			bytes.NewBuffer(updateJSON),
		)
		req.Header.Set("Content-Type", "application/json")

		resp, err = client.Do(req)
		if err != nil {
			t.Fatalf("Failed to update middleware: %v", err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}

		// Get updated middleware
		resp, err = client.Get(fmt.Sprintf("%s/middlewares/%s", baseURL, middlewareID))
		if err != nil {
			t.Fatalf("Failed to get updated middleware: %v", err)
		}

		if err := json.NewDecoder(resp.Body).Decode(&middleware); err != nil {
			t.Fatalf("Failed to decode updated middleware: %v", err)
		}
		resp.Body.Close()

		config, ok := middleware.Config.(map[string]interface{})
		if !ok {
			t.Fatalf("Failed to cast config to map[string]interface{}")
		}

		if config["scheme"] != "http" {
			t.Fatalf("Expected scheme 'http', got '%v'", config["scheme"])
		}

		// Delete middleware
		req, _ = http.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("%s/middlewares/%s", baseURL, middlewareID),
			nil,
		)

		resp, err = client.Do(req)
		if err != nil {
			t.Fatalf("Failed to delete middleware: %v", err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}

		// Verify middleware was deleted
		resp, err = client.Get(fmt.Sprintf("%s/middlewares/%s", baseURL, middlewareID))
		if err != nil {
			t.Fatalf("Failed to get deleted middleware: %v", err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("Expected status code %d, got %d", http.StatusNotFound, resp.StatusCode)
		}
	})

	// Test end-to-end service and router creation
	t.Run("Service and Router E2E", func(t *testing.T) {
		// Create service
		serviceID := "test-service"
		serviceBody := map[string]interface{}{
			"id":  serviceID,
			"url": "http://test-service:8080",
		}
		serviceJSON, _ := json.Marshal(serviceBody)

		resp, err := client.Post(
			fmt.Sprintf("%s/services", baseURL),
			"application/json",
			bytes.NewBuffer(serviceJSON),
		)
		if err != nil {
			t.Fatalf("Failed to create service: %v", err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status code %d, got %d: %s", http.StatusCreated, resp.StatusCode, string(bodyBytes))
		}

		// Create router referencing the service
		routerID := "test-router"
		routerBody := map[string]interface{}{
			"id":          routerID,
			"rule":        "Host(`test.example.com`)",
			"entryPoints": []string{"web"},
			"service": map[string]interface{}{
				"id": serviceID,
			},
		}
		routerJSON, _ := json.Marshal(routerBody)

		resp, err = client.Post(
			fmt.Sprintf("%s/routers", baseURL),
			"application/json",
			bytes.NewBuffer(routerJSON),
		)
		if err != nil {
			t.Fatalf("Failed to create router: %v", err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status code %d, got %d: %s", http.StatusCreated, resp.StatusCode, string(bodyBytes))
		}

		// Test provider endpoint
		resp, err = client.Get(fmt.Sprintf("http://%s:%d%s", cfg.Server.Host, cfg.Server.Port, cfg.Provider.ProviderPath))
		if err != nil {
			t.Fatalf("Failed to get provider config: %v", err)
		}

		var providerConfig map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&providerConfig); err != nil {
			t.Fatalf("Failed to decode provider config: %v", err)
		}
		resp.Body.Close()

		// Verify provider config contains our router and service
		ht, ok := providerConfig["http"].(map[string]interface{})
		if !ok {
			t.Fatalf("Provider config missing HTTP section")
		}

		routers, ok := ht["routers"].(map[string]interface{})
		if !ok {
			t.Fatalf("Provider config missing routers section")
		}

		if _, ok := routers[routerID]; !ok {
			t.Fatalf("Provider config missing router '%s'", routerID)
		}

		services, ok := ht["services"].(map[string]interface{})
		if !ok {
			t.Fatalf("Provider config missing services section")
		}

		if _, ok := services[serviceID]; !ok {
			t.Fatalf("Provider config missing service '%s'", serviceID)
		}

		// Clean up
		// Delete router
		req, _ := http.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("%s/routers/%s", baseURL, routerID),
			nil,
		)

		resp, err = client.Do(req)
		if err != nil {
			t.Fatalf("Failed to delete router: %v", err)
		}
		resp.Body.Close()

		// Delete service
		req, _ = http.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("%s/services/%s", baseURL, serviceID),
			nil,
		)

		resp, err = client.Do(req)
		if err != nil {
			t.Fatalf("Failed to delete service: %v", err)
		}
		resp.Body.Close()
	})

	// Shutdown server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		t.Fatalf("Failed to shutdown server: %v", err)
	}
}
