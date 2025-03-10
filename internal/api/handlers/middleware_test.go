package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/sistemica/traefik-manager/internal/models"
	"github.com/sistemica/traefik-manager/internal/store"
)

// MockStore implements the store.Store interface for testing
type MockStore struct {
	middlewares map[string]models.Middleware
	services    map[string]models.Service
	routers     map[string]models.Router
}

// NewMockStore creates a new mock store for testing
func NewMockStore() *MockStore {
	return &MockStore{
		middlewares: make(map[string]models.Middleware),
		services:    make(map[string]models.Service),
		routers:     make(map[string]models.Router),
	}
}

// Make sure all these methods are defined in your mockstore:

// Middleware methods
func (m *MockStore) ListMiddlewares() ([]models.Middleware, error) {
	result := make([]models.Middleware, 0, len(m.middlewares))
	for _, mw := range m.middlewares {
		result = append(result, mw)
	}
	return result, nil
}

func (m *MockStore) GetMiddleware(id string) (*models.Middleware, error) {
	mw, exists := m.middlewares[id]
	if !exists {
		return nil, store.ErrNotFound
	}
	return &mw, nil
}

func (m *MockStore) CreateMiddleware(middleware *models.Middleware) error {
	if _, exists := m.middlewares[middleware.ID]; exists {
		return store.ErrAlreadyExists
	}
	m.middlewares[middleware.ID] = *middleware
	return nil
}

func (m *MockStore) UpdateMiddleware(id string, middleware *models.Middleware) error {
	if _, exists := m.middlewares[id]; !exists {
		return store.ErrNotFound
	}
	m.middlewares[id] = *middleware
	return nil
}

func (m *MockStore) DeleteMiddleware(id string) error {
	if _, exists := m.middlewares[id]; !exists {
		return store.ErrNotFound
	}

	// Check if in use by any router
	inUse, _, err := m.MiddlewareInUse(id)
	if err != nil {
		return err
	}
	if inUse {
		return store.ErrResourceInUse
	}

	delete(m.middlewares, id)
	return nil
}

func (m *MockStore) MiddlewareExists(id string) (bool, error) {
	_, exists := m.middlewares[id]
	return exists, nil
}

func (m *MockStore) MiddlewareInUse(id string) (bool, []string, error) {
	usedBy := []string{}
	for routerID, router := range m.routers {
		for _, mw := range router.Middlewares {
			if mw.ID == id {
				usedBy = append(usedBy, "router:"+routerID)
				break
			}
		}
	}
	return len(usedBy) > 0, usedBy, nil
}

// Service methods (minimal implementation for tests)
func (m *MockStore) ListServices() ([]models.Service, error) {
	result := make([]models.Service, 0, len(m.services))
	for _, svc := range m.services {
		result = append(result, svc)
	}
	return result, nil
}

func (m *MockStore) GetService(id string) (*models.Service, error) {
	svc, exists := m.services[id]
	if !exists {
		return nil, store.ErrNotFound
	}
	return &svc, nil
}

func (m *MockStore) CreateService(service *models.Service) error {
	if _, exists := m.services[service.ID]; exists {
		return store.ErrAlreadyExists
	}
	m.services[service.ID] = *service
	return nil
}

func (m *MockStore) UpdateService(id string, service *models.Service) error {
	if _, exists := m.services[id]; !exists {
		return store.ErrNotFound
	}
	m.services[id] = *service
	return nil
}

func (m *MockStore) DeleteService(id string) error {
	if _, exists := m.services[id]; !exists {
		return store.ErrNotFound
	}

	// Check if in use by any router
	inUse, _, err := m.ServiceInUse(id)
	if err != nil {
		return err
	}
	if inUse {
		return store.ErrResourceInUse
	}

	delete(m.services, id)
	return nil
}

func (m *MockStore) ServiceExists(id string) (bool, error) {
	_, exists := m.services[id]
	return exists, nil
}

func (m *MockStore) ServiceInUse(id string) (bool, []string, error) {
	usedBy := []string{}
	for routerID, router := range m.routers {
		if router.Service.ID == id {
			usedBy = append(usedBy, "router:"+routerID)
		}
	}
	return len(usedBy) > 0, usedBy, nil
}

// Router methods (minimal implementation for tests)
func (m *MockStore) ListRouters() ([]models.Router, error) {
	result := make([]models.Router, 0, len(m.routers))
	for _, router := range m.routers {
		result = append(result, router)
	}
	return result, nil
}

func (m *MockStore) GetRouter(id string) (*models.Router, error) {
	router, exists := m.routers[id]
	if !exists {
		return nil, store.ErrNotFound
	}
	return &router, nil
}

func (m *MockStore) CreateRouter(router *models.Router) error {
	if _, exists := m.routers[router.ID]; exists {
		return store.ErrAlreadyExists
	}

	// Validate service exists
	if _, exists := m.services[router.Service.ID]; !exists {
		return store.ErrNotFound
	}

	// Validate middlewares exist
	for _, mw := range router.Middlewares {
		if _, exists := m.middlewares[mw.ID]; !exists {
			return store.ErrNotFound
		}
	}

	m.routers[router.ID] = *router
	return nil
}

func (m *MockStore) UpdateRouter(id string, router *models.Router) error {
	if _, exists := m.routers[id]; !exists {
		return store.ErrNotFound
	}

	// Validate service exists
	if _, exists := m.services[router.Service.ID]; !exists {
		return store.ErrNotFound
	}

	// Validate middlewares exist
	for _, mw := range router.Middlewares {
		if _, exists := m.middlewares[mw.ID]; !exists {
			return store.ErrNotFound
		}
	}

	m.routers[id] = *router
	return nil
}

func (m *MockStore) DeleteRouter(id string) error {
	if _, exists := m.routers[id]; !exists {
		return store.ErrNotFound
	}
	delete(m.routers, id)
	return nil
}

func (m *MockStore) RouterExists(id string) (bool, error) {
	_, exists := m.routers[id]
	return exists, nil
}

func (m *MockStore) RouterInUse(id string) (bool, []string, error) {
	// Routers don't have dependencies
	return false, nil, nil
}

// Persistence methods (no-op for mock)
func (m *MockStore) Save() error {
	return nil
}

func (m *MockStore) Load() error {
	return nil
}

func (m *MockStore) Close() {
	// No-op for mock store
}

// TestMiddlewareHandler tests the middleware handler
func TestMiddlewareHandler(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create mock store
	mockStore := NewMockStore()

	// Create handler
	handler := NewMiddlewareHandler(mockStore)

	// Test List when empty
	t.Run("List Empty", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/middlewares", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if err := handler.List(c); err != nil {
			t.Fatalf("Failed to list middlewares: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, rec.Code)
		}

		var middlewares []models.Middleware
		err := json.Unmarshal(rec.Body.Bytes(), &middlewares)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if len(middlewares) != 0 {
			t.Fatalf("Expected empty middleware list, got %d items", len(middlewares))
		}
	})

	// Test Create
	t.Run("Create", func(t *testing.T) {
		reqBody := `{
			"id": "test-middleware",
			"type": "redirectScheme",
			"config": {
				"scheme": "https",
				"permanent": true
			}
		}`

		req := httptest.NewRequest(http.MethodPost, "/api/v1/middlewares", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if err := handler.Create(c); err != nil {
			t.Fatalf("Failed to create middleware: %v", err)
		}

		if rec.Code != http.StatusCreated {
			t.Fatalf("Expected status code %d, got %d", http.StatusCreated, rec.Code)
		}

		// Verify middleware was created in store
		exists, err := mockStore.MiddlewareExists("test-middleware")
		if err != nil {
			t.Fatalf("Failed to check if middleware exists: %v", err)
		}
		if !exists {
			t.Fatalf("Middleware should exist after creation")
		}
	})

	// Test Get
	t.Run("Get", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/middlewares/:id")
		c.SetParamNames("id")
		c.SetParamValues("test-middleware")

		if err := handler.Get(c); err != nil {
			t.Fatalf("Failed to get middleware: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, rec.Code)
		}

		var middleware models.Middleware
		err := json.Unmarshal(rec.Body.Bytes(), &middleware)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if middleware.ID != "test-middleware" {
			t.Fatalf("Expected middleware ID 'test-middleware', got '%s'", middleware.ID)
		}

		if middleware.Type != "redirectScheme" {
			t.Fatalf("Expected middleware type 'redirectScheme', got '%s'", middleware.Type)
		}
	})

	// Test Get Non-Existent
	t.Run("Get Non-Existent", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/middlewares/:id")
		c.SetParamNames("id")
		c.SetParamValues("non-existent")

		if err := handler.Get(c); err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusNotFound {
			t.Fatalf("Expected status code %d, got %d", http.StatusNotFound, rec.Code)
		}
	})

	// Test List
	t.Run("List After Create", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/middlewares", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if err := handler.List(c); err != nil {
			t.Fatalf("Failed to list middlewares: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, rec.Code)
		}

		var middlewares []models.Middleware
		err := json.Unmarshal(rec.Body.Bytes(), &middlewares)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if len(middlewares) != 1 {
			t.Fatalf("Expected 1 middleware, got %d", len(middlewares))
		}

		if middlewares[0].ID != "test-middleware" {
			t.Fatalf("Expected middleware ID 'test-middleware', got '%s'", middlewares[0].ID)
		}
	})

	// Test Update
	t.Run("Update", func(t *testing.T) {
		reqBody := `{
			"type": "redirectScheme",
			"config": {
				"scheme": "http",
				"permanent": false
			}
		}`

		req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/middlewares/:id")
		c.SetParamNames("id")
		c.SetParamValues("test-middleware")

		if err := handler.Update(c); err != nil {
			t.Fatalf("Failed to update middleware: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, rec.Code)
		}

		// Verify middleware was updated in store
		middleware, err := mockStore.GetMiddleware("test-middleware")
		if err != nil {
			t.Fatalf("Failed to get middleware: %v", err)
		}

		config, ok := middleware.Config.(map[string]interface{})
		if !ok {
			t.Fatalf("Failed to cast config to map[string]interface{}")
		}

		if scheme, ok := config["scheme"].(string); !ok || scheme != "http" {
			t.Fatalf("Expected scheme 'http', got '%v'", config["scheme"])
		}

		if permanent, ok := config["permanent"].(bool); !ok || permanent != false {
			t.Fatalf("Expected permanent 'false', got '%v'", config["permanent"])
		}
	})

	// Test dependency check with router
	t.Run("Dependency Check", func(t *testing.T) {
		// Create a service first (needed for router)
		service := &models.Service{
			ID:  "dependency-service",
			URL: "http://test:8080",
		}
		err := mockStore.CreateService(service)
		if err != nil {
			t.Fatalf("Failed to create service: %v", err)
		}

		// Create a router that uses our middleware
		router := &models.Router{
			ID:          "dependency-router",
			Rule:        "Host(`test.example.com`)",
			EntryPoints: []string{"web"},
			Service: models.Service{
				ID: "dependency-service",
			},
			Middlewares: []models.Middleware{
				{
					ID: "test-middleware",
				},
			},
		}
		err = mockStore.CreateRouter(router)
		if err != nil {
			t.Fatalf("Failed to create router: %v", err)
		}

		// Try to delete the middleware (should fail)
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/middlewares/:id")
		c.SetParamNames("id")
		c.SetParamValues("test-middleware")

		if err := handler.Delete(c); err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusConflict {
			t.Fatalf("Expected status code %d, got %d", http.StatusConflict, rec.Code)
		}

		// Delete the router
		err = mockStore.DeleteRouter("dependency-router")
		if err != nil {
			t.Fatalf("Failed to delete router: %v", err)
		}
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/middlewares/:id")
		c.SetParamNames("id")
		c.SetParamValues("test-middleware")

		if err := handler.Delete(c); err != nil {
			t.Fatalf("Failed to delete middleware: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, rec.Code)
		}

		// Verify middleware was deleted from store
		exists, err := mockStore.MiddlewareExists("test-middleware")
		if err != nil {
			t.Fatalf("Failed to check if middleware exists: %v", err)
		}
		if exists {
			t.Fatalf("Middleware should not exist after deletion")
		}
	})

	// Test Create with invalid data
	t.Run("Create Invalid", func(t *testing.T) {
		reqBody := `{
			"id": "",
			"type": "redirectScheme",
			"config": {
				"scheme": "https",
				"permanent": true
			}
		}`

		req := httptest.NewRequest(http.MethodPost, "/api/v1/middlewares", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if err := handler.Create(c); err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("Expected status code %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	// Test Create with duplicate ID
	t.Run("Create Duplicate", func(t *testing.T) {
		// First create a middleware
		middleware := &models.Middleware{
			ID:   "duplicate-middleware",
			Type: "redirectScheme",
			Config: map[string]interface{}{
				"scheme":    "https",
				"permanent": true,
			},
		}
		err := mockStore.CreateMiddleware(middleware)
		if err != nil {
			t.Fatalf("Failed to create middleware: %v", err)
		}

		// Try to create another middleware with the same ID
		reqBody := `{
			"id": "duplicate-middleware",
			"type": "redirectScheme",
			"config": {
				"scheme": "https",
				"permanent": true
			}
		}`

		req := httptest.NewRequest(http.MethodPost, "/api/v1/middlewares", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if err := handler.Create(c); err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusConflict {
			t.Fatalf("Expected status code %d, got %d", http.StatusConflict, rec.Code)
		}

		// Clean up
		err = mockStore.DeleteMiddleware("duplicate-middleware")
		if err != nil {
			t.Fatalf("Failed to delete middleware: %v", err)
		}
	})
}
