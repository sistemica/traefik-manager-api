package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/sistemica/traefik-manager/internal/models"
	"github.com/sistemica/traefik-manager/internal/traefik"
)

// TestProviderHandler tests the provider handler functionality
func TestProviderHandler(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create mock store and set up test data
	mockStore := NewMockStore()
	setupTestData(t, mockStore)

	// Create provider handler
	handler := NewProviderHandler(mockStore)

	// Test GET provider endpoint
	t.Run("Get Provider Config", func(t *testing.T) {
		// Create request
		req := httptest.NewRequest(http.MethodGet, "/traefik/provider", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Call the handler
		if err := handler.GetConfig(c); err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		// Check status code
		if rec.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, rec.Code)
		}

		// Parse response
		var config traefik.DynamicConfig
		if err := json.Unmarshal(rec.Body.Bytes(), &config); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		// Verify router configuration
		if config.HTTP == nil {
			t.Fatalf("HTTP configuration missing in response")
		}

		// Verify routers
		if len(config.HTTP.Routers) != 2 {
			t.Fatalf("Expected 2 routers, got %d", len(config.HTTP.Routers))
		}

		testRouter, exists := config.HTTP.Routers["test-router"]
		if !exists {
			t.Fatalf("test-router not found in provider config")
		}

		if testRouter.Rule != "Host(`test.example.com`)" {
			t.Fatalf("Expected router rule 'Host(`test.example.com`)', got '%s'", testRouter.Rule)
		}

		if testRouter.Service != "test-service" {
			t.Fatalf("Expected router service 'test-service', got '%s'", testRouter.Service)
		}

		// Verify TLS settings
		if testRouter.TLS == nil || testRouter.TLS.CertResolver != "default" {
			t.Fatalf("TLS configuration incorrect for test-router-tls")
		}

		// Verify middlewares
		if len(testRouter.Middlewares) != 1 || testRouter.Middlewares[0] != "test-middleware" {
			t.Fatalf("Middlewares configuration incorrect")
		}

		// Verify services
		if len(config.HTTP.Services) != 2 {
			t.Fatalf("Expected 2 services, got %d", len(config.HTTP.Services))
		}

		testService, exists := config.HTTP.Services["test-service"]
		if !exists {
			t.Fatalf("test-service not found in provider config")
		}

		if testService.LoadBalancer == nil {
			t.Fatalf("LoadBalancer configuration missing for test-service")
		}

		if len(testService.LoadBalancer.Servers) != 1 {
			t.Fatalf("Expected 1 server in load balancer, got %d", len(testService.LoadBalancer.Servers))
		}

		if testService.LoadBalancer.Servers[0].URL != "http://test-backend:8080" {
			t.Fatalf("Expected server URL 'http://test-backend:8080', got '%s'", testService.LoadBalancer.Servers[0].URL)
		}

		// Verify middlewares
		if len(config.HTTP.Middlewares) != 2 {
			t.Fatalf("Expected 2 middlewares, got %d", len(config.HTTP.Middlewares))
		}

		testMiddleware, exists := config.HTTP.Middlewares["test-middleware"]
		if !exists {
			t.Fatalf("test-middleware not found in provider config")
		}

		if testMiddleware.RedirectScheme == nil {
			t.Fatalf("RedirectScheme configuration missing for test-middleware")
		}

		if testMiddleware.RedirectScheme.Scheme != "https" {
			t.Fatalf("Expected redirect scheme 'https', got '%s'", testMiddleware.RedirectScheme.Scheme)
		}
	})
}

// setupTestData adds test data to the mock store
func setupTestData(t *testing.T, mockStore *MockStore) {
	// Create test middleware
	middleware1 := &models.Middleware{
		ID:   "test-middleware",
		Type: "redirectScheme",
		Config: map[string]interface{}{
			"scheme":    "https",
			"permanent": true,
		},
	}

	middleware2 := &models.Middleware{
		ID:   "strip-prefix",
		Type: "stripPrefix",
		Config: map[string]interface{}{
			"prefixes": []string{"/api"},
		},
	}

	if err := mockStore.CreateMiddleware(middleware1); err != nil {
		t.Fatalf("Failed to create test middleware: %v", err)
	}

	if err := mockStore.CreateMiddleware(middleware2); err != nil {
		t.Fatalf("Failed to create strip-prefix middleware: %v", err)
	}

	// Create test service
	service1 := &models.Service{
		ID:  "test-service",
		URL: "http://test-backend:8080",
	}

	service2 := &models.Service{
		ID: "load-balanced-service",
		LoadBalancer: &models.LoadBalancerService{
			Servers: []models.Server{
				{URL: "http://server1:8080"},
				{URL: "http://server2:8080"},
			},
			HealthCheck: &models.HealthCheck{
				Path:     "/health",
				Interval: "10s",
			},
		},
	}

	if err := mockStore.CreateService(service1); err != nil {
		t.Fatalf("Failed to create test service: %v", err)
	}

	if err := mockStore.CreateService(service2); err != nil {
		t.Fatalf("Failed to create load-balanced service: %v", err)
	}

	// Create test router
	router1 := &models.Router{
		ID:          "test-router",
		Rule:        "Host(`test.example.com`)",
		EntryPoints: []string{"web", "websecure"},
		Service: models.Service{
			ID: "test-service",
		},
		Middlewares: []models.Middleware{
			{
				ID: "test-middleware",
			},
		},
		TLS: &models.RouterTLS{
			CertResolver: "default",
		},
	}

	router2 := &models.Router{
		ID:          "api-router",
		Rule:        "Host(`api.example.com`) && PathPrefix(`/api`)",
		EntryPoints: []string{"websecure"},
		Service: models.Service{
			ID: "load-balanced-service",
		},
		Middlewares: []models.Middleware{
			{
				ID: "strip-prefix",
			},
		},
	}

	if err := mockStore.CreateRouter(router1); err != nil {
		t.Fatalf("Failed to create test router: %v", err)
	}

	if err := mockStore.CreateRouter(router2); err != nil {
		t.Fatalf("Failed to create API router: %v", err)
	}
}
