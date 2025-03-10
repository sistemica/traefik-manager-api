// internal/api/handlers/middleware.go
package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sistemica/traefik-manager/internal/logger"
	"github.com/sistemica/traefik-manager/internal/models"
	"github.com/sistemica/traefik-manager/internal/store"
)

// MiddlewareHandler handles middleware-related requests
type MiddlewareHandler struct {
	BaseHandler
}

// NewMiddlewareHandler creates a new MiddlewareHandler
func NewMiddlewareHandler(store store.Store) *MiddlewareHandler {
	return &MiddlewareHandler{
		BaseHandler: NewBaseHandler(store),
	}
}

// List handles the GET /middlewares endpoint to list all middlewares
func (h *MiddlewareHandler) List(c echo.Context) error {
	logger.Debug().Msg("Listing middlewares")

	middlewares, err := h.Store.ListMiddlewares()
	if err != nil {
		logger.Error().Err(err).Msg("Failed to list middlewares")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to list middlewares",
		})
	}

	return c.JSON(http.StatusOK, middlewares)
}

// Get handles the GET /middlewares/:id endpoint to get a specific middleware
func (h *MiddlewareHandler) Get(c echo.Context) error {
	id := c.Param("id")
	logger.Debug().Str("id", id).Msg("Getting middleware")

	middleware, err := h.Store.GetMiddleware(id)
	if err != nil {
		if store.IsNotFound(err) {
			logger.Debug().Str("id", id).Msg("Middleware not found")
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Middleware not found",
			})
		}
		logger.Error().Err(err).Str("id", id).Msg("Failed to get middleware")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get middleware",
		})
	}

	return c.JSON(http.StatusOK, middleware)
}

// Create handles the POST /middlewares endpoint to create a new middleware
func (h *MiddlewareHandler) Create(c echo.Context) error {
	logger.Debug().Msg("Creating middleware")

	var middleware models.Middleware
	if err := c.Bind(&middleware); err != nil {
		logger.Warn().Err(err).Msg("Invalid middleware data")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid middleware data",
		})
	}

	// Validate middleware
	if middleware.ID == "" {
		logger.Warn().Msg("Middleware ID is required")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Middleware ID is required",
		})
	}

	if middleware.Type == "" {
		logger.Warn().Msg("Middleware type is required")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Middleware type is required",
		})
	}

	// Check if middleware already exists
	exists, err := h.Store.MiddlewareExists(middleware.ID)
	if err != nil {
		logger.Error().Err(err).Str("id", middleware.ID).Msg("Failed to check if middleware exists")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to check if middleware exists",
		})
	}

	if exists {
		logger.Warn().Str("id", middleware.ID).Msg("Middleware already exists")
		return c.JSON(http.StatusConflict, map[string]string{
			"error": "Middleware already exists",
		})
	}

	// Create the middleware
	if err := h.Store.CreateMiddleware(&middleware); err != nil {
		logger.Error().Err(err).Str("id", middleware.ID).Msg("Failed to create middleware")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create middleware",
		})
	}

	logger.Info().Str("id", middleware.ID).Msg("Middleware created")

	// Return success response with 201 Created status
	response := models.ResourceResponse{
		ID:      middleware.ID,
		Created: true,
		Updated: false,
		Deleted: false,
	}

	return c.JSON(http.StatusCreated, response)
}

// Update handles the PUT /middlewares/:id endpoint to update a middleware
func (h *MiddlewareHandler) Update(c echo.Context) error {
	id := c.Param("id")
	logger.Debug().Str("id", id).Msg("Updating middleware")

	var middleware models.Middleware
	if err := c.Bind(&middleware); err != nil {
		logger.Warn().Err(err).Str("id", id).Msg("Invalid middleware data")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid middleware data",
		})
	}

	// Ensure ID in path matches ID in body or set it
	if middleware.ID == "" {
		middleware.ID = id
	} else if middleware.ID != id {
		logger.Warn().Str("id", id).Str("body_id", middleware.ID).Msg("ID mismatch")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "ID in path must match ID in body",
		})
	}

	// Check if middleware exists
	exists, err := h.Store.MiddlewareExists(id)
	if err != nil {
		logger.Error().Err(err).Str("id", id).Msg("Failed to check if middleware exists")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to check if middleware exists",
		})
	}

	if !exists {
		logger.Warn().Str("id", id).Msg("Middleware not found")
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Middleware not found",
		})
	}

	// Update the middleware
	if err := h.Store.UpdateMiddleware(id, &middleware); err != nil {
		logger.Error().Err(err).Str("id", id).Msg("Failed to update middleware")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update middleware",
		})
	}

	logger.Info().Str("id", id).Msg("Middleware updated")

	// Return success response
	response := models.ResourceResponse{
		ID:      id,
		Created: false,
		Updated: true,
		Deleted: false,
	}

	return c.JSON(http.StatusOK, response)
}

// Delete handles the DELETE /middlewares/:id endpoint to delete a middleware
func (h *MiddlewareHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	logger.Debug().Str("id", id).Msg("Deleting middleware")

	// Check if middleware exists
	exists, err := h.Store.MiddlewareExists(id)
	if err != nil {
		logger.Error().Err(err).Str("id", id).Msg("Failed to check if middleware exists")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to check if middleware exists",
		})
	}

	if !exists {
		logger.Warn().Str("id", id).Msg("Middleware not found")
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Middleware not found",
		})
	}

	// Check if middleware is in use
	inUse, usedBy, err := h.Store.MiddlewareInUse(id)
	if err != nil {
		logger.Error().Err(err).Str("id", id).Msg("Failed to check if middleware is in use")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to check if middleware is in use",
		})
	}

	if inUse {
		logger.Warn().Str("id", id).Strs("used_by", usedBy).Msg("Middleware is in use")
		return c.JSON(http.StatusConflict, map[string]interface{}{
			"error":   "Middleware is in use by other resources and cannot be deleted",
			"used_by": usedBy,
		})
	}

	// Delete the middleware
	if err := h.Store.DeleteMiddleware(id); err != nil {
		logger.Error().Err(err).Str("id", id).Msg("Failed to delete middleware")

		if store.IsResourceInUse(err) {
			// This is a fallback for any resource-in-use errors not caught by the explicit check
			return c.JSON(http.StatusConflict, map[string]string{
				"error": "Middleware is in use by other resources and cannot be deleted",
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete middleware",
		})
	}

	logger.Info().Str("id", id).Msg("Middleware deleted")

	// Return success response
	response := models.ResourceResponse{
		ID:      id,
		Created: false,
		Updated: false,
		Deleted: true,
	}

	return c.JSON(http.StatusOK, response)
}
