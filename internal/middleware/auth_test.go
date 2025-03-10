package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestAuth(t *testing.T) {
	// Create echo instance
	e := echo.New()

	// Create a simple handler for testing
	testHandler := func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	}

	// Test authentication disabled
	t.Run("Auth Disabled", func(t *testing.T) {
		// Create middleware with authentication disabled
		authMiddleware := Auth(AuthOptions{
			Enabled:      false,
			HeaderName:   "X-API-Key",
			Key:          "test-key",
			ExcludePaths: []string{},
		})

		// Create request
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Test middleware
		handler := authMiddleware(testHandler)
		err := handler(c)

		// Check results
		if err != nil {
			t.Fatalf("Authentication middleware returned error when disabled: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, rec.Code)
		}

		if rec.Body.String() != "success" {
			t.Fatalf("Expected body 'success', got '%s'", rec.Body.String())
		}
	})

	// Test authentication enabled but missing header
	t.Run("Missing API Key", func(t *testing.T) {
		// Create middleware with authentication enabled
		authMiddleware := Auth(AuthOptions{
			Enabled:      true,
			HeaderName:   "X-API-Key",
			Key:          "test-key",
			ExcludePaths: []string{},
		})

		// Create request without API key
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Test middleware
		handler := authMiddleware(testHandler)
		err := handler(c)

		// Check results
		if err != nil {
			t.Fatalf("Authentication middleware should handle errors internally: %v", err)
		}

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("Expected status code %d, got %d", http.StatusUnauthorized, rec.Code)
		}
	})

	// Test authentication enabled with invalid key
	t.Run("Invalid API Key", func(t *testing.T) {
		// Create middleware with authentication enabled
		authMiddleware := Auth(AuthOptions{
			Enabled:      true,
			HeaderName:   "X-API-Key",
			Key:          "test-key",
			ExcludePaths: []string{},
		})

		// Create request with invalid API key
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-API-Key", "wrong-key")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Test middleware
		handler := authMiddleware(testHandler)
		err := handler(c)

		// Check results
		if err != nil {
			t.Fatalf("Authentication middleware should handle errors internally: %v", err)
		}

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("Expected status code %d, got %d", http.StatusUnauthorized, rec.Code)
		}
	})

	// Test authentication enabled with valid key
	t.Run("Valid API Key", func(t *testing.T) {
		// Create middleware with authentication enabled
		authMiddleware := Auth(AuthOptions{
			Enabled:      true,
			HeaderName:   "X-API-Key",
			Key:          "test-key",
			ExcludePaths: []string{},
		})

		// Create request with valid API key
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-API-Key", "test-key")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Test middleware
		handler := authMiddleware(testHandler)
		err := handler(c)

		// Check results
		if err != nil {
			t.Fatalf("Authentication middleware returned error with valid key: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, rec.Code)
		}

		if rec.Body.String() != "success" {
			t.Fatalf("Expected body 'success', got '%s'", rec.Body.String())
		}
	})

	// Test authentication with excluded path
	t.Run("Excluded Path", func(t *testing.T) {
		// Create middleware with authentication enabled and excluded path
		authMiddleware := Auth(AuthOptions{
			Enabled:      true,
			HeaderName:   "X-API-Key",
			Key:          "test-key",
			ExcludePaths: []string{"/health", "/public"},
		})

		// Create request to excluded path without API key
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/health")

		// Test middleware
		handler := authMiddleware(testHandler)
		err := handler(c)

		// Check results
		if err != nil {
			t.Fatalf("Authentication middleware returned error for excluded path: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, rec.Code)
		}

		if rec.Body.String() != "success" {
			t.Fatalf("Expected body 'success', got '%s'", rec.Body.String())
		}
	})

	// Test authentication with path that starts with excluded prefix
	t.Run("Path With Excluded Prefix", func(t *testing.T) {
		// Create middleware with authentication enabled and excluded path prefix
		authMiddleware := Auth(AuthOptions{
			Enabled:      true,
			HeaderName:   "X-API-Key",
			Key:          "test-key",
			ExcludePaths: []string{"/public"},
		})

		// Create request to path with excluded prefix without API key
		req := httptest.NewRequest(http.MethodGet, "/public/assets/image.jpg", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/public/assets/image.jpg")

		// Test middleware
		handler := authMiddleware(testHandler)
		err := handler(c)

		// Check results
		if err != nil {
			t.Fatalf("Authentication middleware returned error for path with excluded prefix: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, rec.Code)
		}

		if rec.Body.String() != "success" {
			t.Fatalf("Expected body 'success', got '%s'", rec.Body.String())
		}
	})
}
