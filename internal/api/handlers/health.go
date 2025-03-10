// internal/api/handlers/health.go
package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sistemica/traefik-manager/internal/logger"
	"github.com/sistemica/traefik-manager/internal/store"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	BaseHandler
	Version   string
	StartTime time.Time
}

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler(store store.Store, version string) *HealthHandler {
	return &HealthHandler{
		BaseHandler: NewBaseHandler(store),
		Version:     version,
		StartTime:   time.Now(),
	}
}

// Check handles the health check endpoint
func (h *HealthHandler) Check(c echo.Context) error {
	logger.Debug().Msg("Health check requested")

	// Calculate uptime
	uptime := time.Since(h.StartTime)

	response := map[string]interface{}{
		"status":  "healthy",
		"version": h.Version,
		"uptime":  uptime.String(),
	}

	return c.JSON(http.StatusOK, response)
}
