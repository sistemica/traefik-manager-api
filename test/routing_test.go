// Modified test/routing_test.go with fixes for the type issues
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
	"github.com/sistemica/traefik-manager/internal/store"
)

func TestComplexRouting(t *testing.T) {
	// Skip in CI environment if needed
	if os.Getenv("CI") != "" && os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration tests in CI environment")
	}

	// Create temporary file for store
	tmpFile, err := os.CreateTemp("", "complex-routing-test-*.json")
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

	// Use a unique port for this test
	port := 9990

	// Create test configuration
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
	defer dataStore.Close()

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
	providerURL := fmt.Sprintf("http://%s:%d%s", cfg.Server.Host, cfg.Server.Port, cfg.Provider.ProviderPath)

	// Create HTTP client
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Run the complex workflow test
	t.Run("Complex Routing Workflow", func(t *testing.T) {
		// Step 1: Create necessary middleware components
		middlewares := []map[string]interface{}{
			{
				"id":   "https-redirect",
				"type": "redirectScheme",
				"config": map[string]interface{}{
					"scheme":    "https",
					"permanent": true,
				},
			},
			{
				"id":   "strip-api",
				"type": "stripPrefix",
				"config": map[string]interface{}{
					"prefixes": []string{"/api"},
				},
			},
			{
				"id":   "admin-auth",
				"type": "basicAuth",
				"config": map[string]interface{}{
					"users": []string{"admin:$apr1$H6uskkkW$IgXLP6ewTrSuBkTrqE8wj/"},
					"realm": "Admin Area",
				},
			},
			{
				"id":   "rate-limit",
				"type": "rateLimit",
				"config": map[string]interface{}{
					"average": 100,
					"burst":   50,
				},
			},
		}

		for _, mw := range middlewares {
			middlewareJSON, _ := json.Marshal(mw)
			resp, err := client.Post(
				fmt.Sprintf("%s/middlewares", baseURL),
				"application/json",
				bytes.NewBuffer(middlewareJSON),
			)
			if err != nil {
				t.Fatalf("Failed to create middleware %s: %v", mw["id"], err)
			}
			resp.Body.Close()

			if resp.StatusCode != http.StatusCreated {
				t.Fatalf("Expected status code %d for middleware %s, got %d",
					http.StatusCreated, mw["id"], resp.StatusCode)
			}
		}

		// Step 2: Create services
		services := []map[string]interface{}{
			{
				"id":  "web-service",
				"url": "http://web-backend:8080",
			},
			{
				"id": "api-service",
				"loadBalancer": map[string]interface{}{
					"servers": []map[string]interface{}{
						{"url": "http://api-server1:8000"},
						{"url": "http://api-server2:8000"},
					},
					"healthCheck": map[string]interface{}{
						"path":     "/health",
						"interval": "10s",
					},
				},
			},
			{
				"id":  "admin-service",
				"url": "http://admin-backend:8080",
			},
		}

		for _, svc := range services {
			serviceJSON, _ := json.Marshal(svc)
			resp, err := client.Post(
				fmt.Sprintf("%s/services", baseURL),
				"application/json",
				bytes.NewBuffer(serviceJSON),
			)
			if err != nil {
				t.Fatalf("Failed to create service %s: %v", svc["id"], err)
			}
			resp.Body.Close()

			if resp.StatusCode != http.StatusCreated {
				t.Fatalf("Expected status code %d for service %s, got %d",
					http.StatusCreated, svc["id"], resp.StatusCode)
			}
		}

		// Step 3: Create routers with different configurations
		routers := []map[string]interface{}{
			{
				"id":          "web-http",
				"rule":        "Host(`example.com`)",
				"entryPoints": []string{"web"},
				"service":     "web-service",
				"middlewares": []string{"https-redirect"},
			},
			{
				"id":          "web-https",
				"rule":        "Host(`example.com`)",
				"entryPoints": []string{"websecure"},
				"service":     "web-service",
				"tls": map[string]interface{}{
					"certResolver": "default",
				},
			},
			{
				"id":          "api-router",
				"rule":        "Host(`api.example.com`) && PathPrefix(`/api`)",
				"entryPoints": []string{"websecure"},
				"service":     "api-service",
				"middlewares": []string{"strip-api", "rate-limit"},
				"tls":         map[string]interface{}{},
			},
			{
				"id":          "admin-router",
				"rule":        "Host(`admin.example.com`)",
				"entryPoints": []string{"websecure"},
				"service":     "admin-service",
				"middlewares": []string{"admin-auth"},
				"tls": map[string]interface{}{
					"certResolver": "default",
				},
			},
		}

		for _, rtr := range routers {
			routerJSON, _ := json.Marshal(rtr)
			resp, err := client.Post(
				fmt.Sprintf("%s/routers", baseURL),
				"application/json",
				bytes.NewBuffer(routerJSON),
			)
			if err != nil {
				t.Fatalf("Failed to create router %s: %v", rtr["id"], err)
			}
			resp.Body.Close()

			if resp.StatusCode != http.StatusCreated {
				body, _ := io.ReadAll(resp.Body)
				t.Fatalf("Expected status code %d for router %s, got %d: %s",
					http.StatusCreated, rtr["id"], resp.StatusCode, string(body))
			}
		}

		// Step 4: Verify the provider endpoint includes all configurations
		resp, err := client.Get(providerURL)
		if err != nil {
			t.Fatalf("Failed to get provider config: %v", err)
		}

		var providerConfig map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&providerConfig); err != nil {
			t.Fatalf("Failed to decode provider config: %v", err)
		}
		resp.Body.Close()

		// Verify HTTP section exists
		httpSection, ok := providerConfig["http"].(map[string]interface{})
		if !ok {
			t.Fatalf("HTTP section missing in provider config")
		}

		// Verify all components are present
		verifyConfigCount(t, httpSection, "routers", 4)
		verifyConfigCount(t, httpSection, "services", 3)
		verifyConfigCount(t, httpSection, "middlewares", 4)

		// Verify specific router configuration
		routersSection, ok := httpSection["routers"].(map[string]interface{})
		if !ok {
			t.Fatalf("Routers section missing in provider config")
		}

		apiRouter, ok := routersSection["api-router"].(map[string]interface{})
		if !ok {
			t.Fatalf("api-router missing in provider config")
		}

		// Check router rule
		if rule, ok := apiRouter["rule"].(string); !ok || rule != "Host(`api.example.com`) && PathPrefix(`/api`)" {
			t.Fatalf("Incorrect rule for api-router: %v", apiRouter["rule"])
		}

		// Check router middlewares
		apiMiddlewares, ok := apiRouter["middlewares"].([]interface{})
		if !ok || len(apiMiddlewares) != 2 {
			t.Fatalf("Expected 2 middlewares for api-router, got: %v", apiRouter["middlewares"])
		}

		// Step 5: Clean up all resources
		t.Log("Cleaning up resources...")

		// Create string slices for the IDs to make the loop simpler
		routerIDs := []string{"web-http", "web-https", "api-router", "admin-router"}
		serviceIDs := []string{"web-service", "api-service", "admin-service"}
		middlewareIDs := []string{"https-redirect", "strip-api", "admin-auth", "rate-limit"}

		for _, id := range routerIDs {
			req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/routers/%s", baseURL, id), nil)
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Failed to delete router %s: %v", id, err)
			}
			resp.Body.Close()
		}

		for _, id := range serviceIDs {
			req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/services/%s", baseURL, id), nil)
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Failed to delete service %s: %v", id, err)
			}
			resp.Body.Close()
		}

		for _, id := range middlewareIDs {
			req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/middlewares/%s", baseURL, id), nil)
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Failed to delete middleware %s: %v", id, err)
			}
			resp.Body.Close()
		}
	})

	// Shutdown server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		t.Fatalf("Failed to shutdown server: %v", err)
	}
}

// Helper function to verify config counts
func verifyConfigCount(t *testing.T, httpSection map[string]interface{}, section string, expectedCount int) {
	sectionData, ok := httpSection[section].(map[string]interface{})
	if !ok {
		t.Fatalf("%s section missing in provider config", section)
	}

	if len(sectionData) != expectedCount {
		t.Fatalf("Expected %d %s, got %d", expectedCount, section, len(sectionData))
	}
}
