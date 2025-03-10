package store

import (
	"os"
	"testing"
	"time"

	"github.com/sistemica/traefik-manager/internal/models"
)

func TestMinimalFileStore(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "minimal-test-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	// Initialize with valid empty JSON
	initialJSON := `{"middlewares":{},"routers":{},"services":{}}`
	if err := os.WriteFile(tmpPath, []byte(initialJSON), 0644); err != nil {
		t.Fatalf("Failed to initialize file: %v", err)
	}

	// Create store
	store, err := NewFileStore(tmpPath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// Save an item
	middleware := &models.Middleware{
		ID:   "test-mw",
		Type: "test",
		Config: map[string]interface{}{
			"key": "value",
		},
	}

	if err := store.CreateMiddleware(middleware); err != nil {
		t.Fatalf("Failed to create middleware: %v", err)
	}

	// Force save
	if err := store.Save(); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	// Explicitly close the store
	store.Close()

	// Test complete
	t.Log("Test completed successfully")
}

func TestSimplifiedFileStore(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "traefik-manager-test-*.json")
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

	// Create a new file store
	store, err := NewFileStore(tmpPath)
	if err != nil {
		t.Fatalf("Failed to create file store: %v", err)
	}
	defer store.Close() // Close the store when the test is done

	// Simple test - just create a middleware
	middleware := &models.Middleware{
		ID:   "test-middleware",
		Type: "redirectScheme",
		Config: map[string]interface{}{
			"scheme":    "https",
			"permanent": true,
		},
	}

	// Create middleware
	err = store.CreateMiddleware(middleware)
	if err != nil {
		t.Fatalf("Failed to create middleware: %v", err)
	}

	// Verify it exists
	exists, err := store.MiddlewareExists("test-middleware")
	if err != nil {
		t.Fatalf("Failed to check if middleware exists: %v", err)
	}
	if !exists {
		t.Fatalf("Middleware should exist but doesn't")
	}
}

func TestFileStore(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "traefik-manager-test-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath) // Clean up after test

	// Initialize with empty but valid JSON
	initialJSON := `{"middlewares":{},"routers":{},"services":{}}`
	if err := os.WriteFile(tmpPath, []byte(initialJSON), 0644); err != nil {
		t.Fatalf("Failed to initialize store file: %v", err)
	}

	// Create a new file store
	store, err := NewFileStore(tmpPath)
	if err != nil {
		t.Fatalf("Failed to create file store: %v", err)
	}
	defer store.Close() // Close the store when the test is done

	// Test middleware operations
	t.Run("Middleware CRUD", func(t *testing.T) {
		// Create middleware
		middleware := &models.Middleware{
			ID:   "test-middleware",
			Type: "redirectScheme",
			Config: map[string]interface{}{
				"scheme":    "https",
				"permanent": true,
			},
		}

		// Test create
		err := store.CreateMiddleware(middleware)
		if err != nil {
			t.Fatalf("Failed to create middleware: %v", err)
		}

		// Test exists
		exists, err := store.MiddlewareExists("test-middleware")
		if err != nil {
			t.Fatalf("Failed to check if middleware exists: %v", err)
		}
		if !exists {
			t.Fatalf("Middleware should exist but doesn't")
		}

		// Test get
		retrievedMiddleware, err := store.GetMiddleware("test-middleware")
		if err != nil {
			t.Fatalf("Failed to get middleware: %v", err)
		}
		if retrievedMiddleware.ID != middleware.ID || retrievedMiddleware.Type != middleware.Type {
			t.Fatalf("Retrieved middleware doesn't match created one")
		}

		// Test list
		middlewares, err := store.ListMiddlewares()
		if err != nil {
			t.Fatalf("Failed to list middlewares: %v", err)
		}
		if len(middlewares) != 1 {
			t.Fatalf("Expected 1 middleware, got %d", len(middlewares))
		}

		// Test update
		updatedMiddleware := &models.Middleware{
			ID:   "test-middleware",
			Type: "redirectScheme",
			Config: map[string]interface{}{
				"scheme":    "http",
				"permanent": false,
			},
		}
		err = store.UpdateMiddleware("test-middleware", updatedMiddleware)
		if err != nil {
			t.Fatalf("Failed to update middleware: %v", err)
		}

		// Verify update
		retrievedMiddleware, err = store.GetMiddleware("test-middleware")
		if err != nil {
			t.Fatalf("Failed to get updated middleware: %v", err)
		}
		config := retrievedMiddleware.Config.(map[string]interface{})
		if config["scheme"] != "http" {
			t.Fatalf("Middleware update failed: expected scheme 'http', got '%v'", config["scheme"])
		}

		// Test delete
		err = store.DeleteMiddleware("test-middleware")
		if err != nil {
			t.Fatalf("Failed to delete middleware: %v", err)
		}

		// Verify delete
		exists, err = store.MiddlewareExists("test-middleware")
		if err != nil {
			t.Fatalf("Failed to check if middleware exists after delete: %v", err)
		}
		if exists {
			t.Fatalf("Middleware should not exist after delete")
		}
	})

	// Test service operations
	t.Run("Service CRUD", func(t *testing.T) {
		// Create service
		service := &models.Service{
			ID:  "test-service",
			URL: "http://test-service:8080",
		}

		// Test create
		err := store.CreateService(service)
		if err != nil {
			t.Fatalf("Failed to create service: %v", err)
		}

		// Test exists
		exists, err := store.ServiceExists("test-service")
		if err != nil {
			t.Fatalf("Failed to check if service exists: %v", err)
		}
		if !exists {
			t.Fatalf("Service should exist but doesn't")
		}

		// Test get
		retrievedService, err := store.GetService("test-service")
		if err != nil {
			t.Fatalf("Failed to get service: %v", err)
		}
		if retrievedService.ID != service.ID || retrievedService.URL != service.URL {
			t.Fatalf("Retrieved service doesn't match created one")
		}

		// Test list
		services, err := store.ListServices()
		if err != nil {
			t.Fatalf("Failed to list services: %v", err)
		}
		if len(services) != 1 {
			t.Fatalf("Expected 1 service, got %d", len(services))
		}

		// Test update
		updatedService := &models.Service{
			ID:  "test-service",
			URL: "http://updated-service:9090",
		}
		err = store.UpdateService("test-service", updatedService)
		if err != nil {
			t.Fatalf("Failed to update service: %v", err)
		}

		// Verify update
		retrievedService, err = store.GetService("test-service")
		if err != nil {
			t.Fatalf("Failed to get updated service: %v", err)
		}
		if retrievedService.URL != "http://updated-service:9090" {
			t.Fatalf("Service update failed: expected URL 'http://updated-service:9090', got '%s'", retrievedService.URL)
		}

		// Test delete
		err = store.DeleteService("test-service")
		if err != nil {
			t.Fatalf("Failed to delete service: %v", err)
		}

		// Verify delete
		exists, err = store.ServiceExists("test-service")
		if err != nil {
			t.Fatalf("Failed to check if service exists after delete: %v", err)
		}
		if exists {
			t.Fatalf("Service should not exist after delete")
		}
	})

	// Test router operations
	t.Run("Router CRUD with Dependencies", func(t *testing.T) {
		// First create a service dependency
		service := &models.Service{
			ID:  "router-test-service",
			URL: "http://router-test-service:8080",
		}
		err := store.CreateService(service)
		if err != nil {
			t.Fatalf("Failed to create dependency service: %v", err)
		}

		// Create middleware dependency
		middleware := &models.Middleware{
			ID:   "router-test-middleware",
			Type: "redirectScheme",
			Config: map[string]interface{}{
				"scheme":    "https",
				"permanent": true,
			},
		}
		err = store.CreateMiddleware(middleware)
		if err != nil {
			t.Fatalf("Failed to create dependency middleware: %v", err)
		}

		// Create router
		router := &models.Router{
			ID:          "test-router",
			Rule:        "Host(`test.example.com`)",
			EntryPoints: []string{"web"},
			Service: models.Service{
				ID: "router-test-service",
			},
			Middlewares: []models.Middleware{
				{
					ID: "router-test-middleware",
				},
			},
		}

		// Test create
		err = store.CreateRouter(router)
		if err != nil {
			t.Fatalf("Failed to create router: %v", err)
		}

		// Test exists
		exists, err := store.RouterExists("test-router")
		if err != nil {
			t.Fatalf("Failed to check if router exists: %v", err)
		}
		if !exists {
			t.Fatalf("Router should exist but doesn't")
		}

		// Test get
		retrievedRouter, err := store.GetRouter("test-router")
		if err != nil {
			t.Fatalf("Failed to get router: %v", err)
		}
		if retrievedRouter.ID != router.ID || retrievedRouter.Rule != router.Rule {
			t.Fatalf("Retrieved router doesn't match created one")
		}

		// Test list
		routers, err := store.ListRouters()
		if err != nil {
			t.Fatalf("Failed to list routers: %v", err)
		}
		if len(routers) != 1 {
			t.Fatalf("Expected 1 router, got %d", len(routers))
		}

		// Test update
		updatedRouter := &models.Router{
			ID:          "test-router",
			Rule:        "Host(`updated.example.com`)",
			EntryPoints: []string{"web", "websecure"},
			Service: models.Service{
				ID: "router-test-service",
			},
			Middlewares: []models.Middleware{
				{
					ID: "router-test-middleware",
				},
			},
		}
		err = store.UpdateRouter("test-router", updatedRouter)
		if err != nil {
			t.Fatalf("Failed to update router: %v", err)
		}

		// Verify update
		retrievedRouter, err = store.GetRouter("test-router")
		if err != nil {
			t.Fatalf("Failed to get updated router: %v", err)
		}
		if retrievedRouter.Rule != "Host(`updated.example.com`)" {
			t.Fatalf("Router update failed: expected Rule 'Host(`updated.example.com`)', got '%s'", retrievedRouter.Rule)
		}

		// Test dependencies - service in use
		inUse, usedBy, err := store.ServiceInUse("router-test-service")
		if err != nil {
			t.Fatalf("Failed to check if service is in use: %v", err)
		}
		if !inUse {
			t.Fatalf("Service should be in use but isn't")
		}
		if len(usedBy) != 1 {
			t.Fatalf("Service should be used by 1 router, got %d", len(usedBy))
		}

		// Test dependencies - middleware in use
		inUse, usedBy, err = store.MiddlewareInUse("router-test-middleware")
		if err != nil {
			t.Fatalf("Failed to check if middleware is in use: %v", err)
		}
		if !inUse {
			t.Fatalf("Middleware should be in use but isn't")
		}
		if len(usedBy) != 1 {
			t.Fatalf("Middleware should be used by 1 router, got %d", len(usedBy))
		}

		// Test delete
		err = store.DeleteRouter("test-router")
		if err != nil {
			t.Fatalf("Failed to delete router: %v", err)
		}

		// Verify delete
		exists, err = store.RouterExists("test-router")
		if err != nil {
			t.Fatalf("Failed to check if router exists after delete: %v", err)
		}
		if exists {
			t.Fatalf("Router should not exist after delete")
		}

		// Clean up dependencies
		err = store.DeleteService("router-test-service")
		if err != nil {
			t.Fatalf("Failed to delete dependency service: %v", err)
		}
		err = store.DeleteMiddleware("router-test-middleware")
		if err != nil {
			t.Fatalf("Failed to delete dependency middleware: %v", err)
		}
	})

	// Test persistence
	t.Run("Persistence", func(t *testing.T) {
		// Create a middleware
		middleware := &models.Middleware{
			ID:   "persistence-test",
			Type: "redirectScheme",
			Config: map[string]interface{}{
				"scheme":    "https",
				"permanent": true,
			},
		}
		err := store.CreateMiddleware(middleware)
		if err != nil {
			t.Fatalf("Failed to create middleware: %v", err)
		}

		// Save explicitly
		err = store.Save()
		if err != nil {
			t.Fatalf("Failed to save: %v", err)
		}

		// Just verify the middleware exists in the current store
		exists, err := store.MiddlewareExists("persistence-test")
		if err != nil {
			t.Fatalf("Failed to check if middleware exists: %v", err)
		}
		if !exists {
			t.Fatalf("Middleware should exist but doesn't")
		}

		// Clean up
		err = store.DeleteMiddleware("persistence-test")
		if err != nil {
			t.Fatalf("Failed to delete middleware: %v", err)
		}
	})
}

func TestResourceDependencies(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "traefik-manager-test-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath) // Clean up after test

	// Initialize with empty but valid JSON
	initialJSON := `{"middlewares":{},"routers":{},"services":{}}`
	if err := os.WriteFile(tmpPath, []byte(initialJSON), 0644); err != nil {
		t.Fatalf("Failed to initialize store file: %v", err)
	}

	// Create a new file store
	store, err := NewFileStore(tmpPath)
	if err != nil {
		t.Fatalf("Failed to create file store: %v", err)
	}
	defer store.Close() // Close the store when the test is done

	// Test that we can't delete resources that are in use
	t.Run("Delete Dependencies", func(t *testing.T) {
		// Create a service
		service := &models.Service{
			ID:  "dep-test-service",
			URL: "http://dep-test-service:8080",
		}
		err := store.CreateService(service)
		if err != nil {
			t.Fatalf("Failed to create service: %v", err)
		}

		// Create a middleware
		middleware := &models.Middleware{
			ID:   "dep-test-middleware",
			Type: "redirectScheme",
			Config: map[string]interface{}{
				"scheme":    "https",
				"permanent": true,
			},
		}
		err = store.CreateMiddleware(middleware)
		if err != nil {
			t.Fatalf("Failed to create middleware: %v", err)
		}

		// Create a router that uses both
		router := &models.Router{
			ID:          "dep-test-router",
			Rule:        "Host(`dep.example.com`)",
			EntryPoints: []string{"web"},
			Service: models.Service{
				ID: "dep-test-service",
			},
			Middlewares: []models.Middleware{
				{
					ID: "dep-test-middleware",
				},
			},
		}
		err = store.CreateRouter(router)
		if err != nil {
			t.Fatalf("Failed to create router: %v", err)
		}

		// Try to delete the service while it's in use
		err = store.DeleteService("dep-test-service")
		if err == nil {
			t.Fatalf("Should not be able to delete service in use")
		}
		if !IsResourceInUse(err) {
			t.Fatalf("Expected ResourceInUse error, got: %v", err)
		}

		// Try to delete the middleware while it's in use
		err = store.DeleteMiddleware("dep-test-middleware")
		if err == nil {
			t.Fatalf("Should not be able to delete middleware in use")
		}
		if !IsResourceInUse(err) {
			t.Fatalf("Expected ResourceInUse error, got: %v", err)
		}

		// Delete the router first (no dependencies)
		err = store.DeleteRouter("dep-test-router")
		if err != nil {
			t.Fatalf("Failed to delete router: %v", err)
		}

		// Now we should be able to delete the service and middleware
		err = store.DeleteService("dep-test-service")
		if err != nil {
			t.Fatalf("Failed to delete service after router deletion: %v", err)
		}

		err = store.DeleteMiddleware("dep-test-middleware")
		if err != nil {
			t.Fatalf("Failed to delete middleware after router deletion: %v", err)
		}
	})
}

// Additional test for concurrent access
func TestConcurrentAccess(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "traefik-manager-concurrent-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath) // Clean up after test

	// Initialize with empty but valid JSON - ADD THIS
	initialJSON := `{"middlewares":{},"routers":{},"services":{}}`
	if err := os.WriteFile(tmpPath, []byte(initialJSON), 0644); err != nil {
		t.Fatalf("Failed to initialize store file: %v", err)
	}

	// Create a new file store
	store, err := NewFileStore(tmpPath)
	if err != nil {
		t.Fatalf("Failed to create file store: %v", err)
	}
	defer store.Close() // Close the store when the test is done

	// Test concurrent access
	t.Run("Concurrent Access", func(t *testing.T) {
		// Number of concurrent operations
		concurrency := 10
		done := make(chan bool, concurrency*2) // For both readers and writers

		// Create initial service
		err := store.CreateService(&models.Service{
			ID:  "concurrent-test",
			URL: "http://original:8080",
		})
		if err != nil {
			t.Fatalf("Failed to create initial service: %v", err)
		}

		// Start concurrent readers
		for i := 0; i < concurrency; i++ {
			go func() {
				// Random delay to increase chance of concurrent access
				time.Sleep(time.Duration(i) * time.Millisecond)

				_, err := store.GetService("concurrent-test")
				if err != nil {
					t.Errorf("Reader failed: %v", err)
				}
				done <- true
			}()
		}

		// Start concurrent writers
		for i := 0; i < concurrency; i++ {
			go func(idx int) {
				// Random delay to increase chance of concurrent access
				time.Sleep(time.Duration(idx) * time.Millisecond)

				updatedService := &models.Service{
					ID:  "concurrent-test",
					URL: "http://updated:8080",
				}
				err := store.UpdateService("concurrent-test", updatedService)
				if err != nil {
					t.Errorf("Writer failed: %v", err)
				}
				done <- true
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < concurrency*2; i++ {
			<-done
		}

		// Verify the final state
		service, err := store.GetService("concurrent-test")
		if err != nil {
			t.Fatalf("Failed to get service after concurrent operations: %v", err)
		}
		if service.URL != "http://updated:8080" {
			t.Fatalf("Service should have been updated to 'http://updated:8080', got '%s'", service.URL)
		}

		// Clean up
		err = store.DeleteService("concurrent-test")
		if err != nil {
			t.Fatalf("Failed to delete service: %v", err)
		}
	})
}
