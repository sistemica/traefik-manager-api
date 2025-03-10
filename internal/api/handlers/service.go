// internal/api/handlers/service.go
package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sistemica/traefik-manager/internal/logger"
	"github.com/sistemica/traefik-manager/internal/models"
	"github.com/sistemica/traefik-manager/internal/store"
)

// ServiceHandler handles service-related requests
type ServiceHandler struct {
	BaseHandler
}

// NewServiceHandler creates a new ServiceHandler
func NewServiceHandler(store store.Store) *ServiceHandler {
	return &ServiceHandler{
		BaseHandler: NewBaseHandler(store),
	}
}

// List handles the GET /services endpoint to list all services
func (h *ServiceHandler) List(c echo.Context) error {
	logger.Debug().Msg("Listing services")

	services, err := h.Store.ListServices()
	if err != nil {
		logger.Error().Err(err).Msg("Failed to list services")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to list services",
		})
	}

	return c.JSON(http.StatusOK, services)
}

// Get handles the GET /services/:id endpoint to get a specific service
func (h *ServiceHandler) Get(c echo.Context) error {
	id := c.Param("id")
	logger.Debug().Str("id", id).Msg("Getting service")

	service, err := h.Store.GetService(id)
	if err != nil {
		if store.IsNotFound(err) {
			logger.Debug().Str("id", id).Msg("Service not found")
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Service not found",
			})
		}
		logger.Error().Err(err).Str("id", id).Msg("Failed to get service")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get service",
		})
	}

	return c.JSON(http.StatusOK, service)
}

// Create handles the POST /services endpoint to create a new service
func (h *ServiceHandler) Create(c echo.Context) error {
	logger.Debug().Msg("Creating service")

	var service models.Service
	if err := c.Bind(&service); err != nil {
		logger.Warn().Err(err).Msg("Invalid service data")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid service data",
		})
	}

	// Validate service
	if service.ID == "" {
		logger.Warn().Msg("Service ID is required")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Service ID is required",
		})
	}

	// Check service type
	if err := validateServiceConfiguration(&service); err != nil {
		logger.Warn().Err(err).Msg("Invalid service configuration")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	// Check if service already exists
	exists, err := h.Store.ServiceExists(service.ID)
	if err != nil {
		logger.Error().Err(err).Str("id", service.ID).Msg("Failed to check if service exists")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to check if service exists",
		})
	}

	if exists {
		logger.Warn().Str("id", service.ID).Msg("Service already exists")
		return c.JSON(http.StatusConflict, map[string]string{
			"error": "Service already exists",
		})
	}

	// Create the service
	if err := h.Store.CreateService(&service); err != nil {
		logger.Error().Err(err).Str("id", service.ID).Msg("Failed to create service")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create service",
		})
	}

	logger.Info().Str("id", service.ID).Msg("Service created")

	// Return success response with 201 Created status
	response := models.ResourceResponse{
		ID:      service.ID,
		Created: true,
		Updated: false,
		Deleted: false,
	}

	return c.JSON(http.StatusCreated, response)
}

// Update handles the PUT /services/:id endpoint to update a service
func (h *ServiceHandler) Update(c echo.Context) error {
	id := c.Param("id")
	logger.Debug().Str("id", id).Msg("Updating service")

	var service models.Service
	if err := c.Bind(&service); err != nil {
		logger.Warn().Err(err).Str("id", id).Msg("Invalid service data")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid service data",
		})
	}

	// Ensure ID in path matches ID in body or set it
	if service.ID == "" {
		service.ID = id
	} else if service.ID != id {
		logger.Warn().Str("id", id).Str("body_id", service.ID).Msg("ID mismatch")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "ID in path must match ID in body",
		})
	}

	// Check service type
	if err := validateServiceConfiguration(&service); err != nil {
		logger.Warn().Err(err).Msg("Invalid service configuration")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	// Check if service exists
	exists, err := h.Store.ServiceExists(id)
	if err != nil {
		logger.Error().Err(err).Str("id", id).Msg("Failed to check if service exists")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to check if service exists",
		})
	}

	if !exists {
		logger.Warn().Str("id", id).Msg("Service not found")
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Service not found",
		})
	}

	// Update the service
	if err := h.Store.UpdateService(id, &service); err != nil {
		logger.Error().Err(err).Str("id", id).Msg("Failed to update service")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update service",
		})
	}

	logger.Info().Str("id", id).Msg("Service updated")

	// Return success response
	response := models.ResourceResponse{
		ID:      id,
		Created: false,
		Updated: true,
		Deleted: false,
	}

	return c.JSON(http.StatusOK, response)
}

// Delete handles the DELETE /services/:id endpoint to delete a service
func (h *ServiceHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	logger.Debug().Str("id", id).Msg("Deleting service")

	// Check if service exists
	exists, err := h.Store.ServiceExists(id)
	if err != nil {
		logger.Error().Err(err).Str("id", id).Msg("Failed to check if service exists")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to check if service exists",
		})
	}

	if !exists {
		logger.Warn().Str("id", id).Msg("Service not found")
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Service not found",
		})
	}

	// Check if service is in use
	inUse, usedBy, err := h.Store.ServiceInUse(id)
	if err != nil {
		logger.Error().Err(err).Str("id", id).Msg("Failed to check if service is in use")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to check if service is in use",
		})
	}

	if inUse {
		logger.Warn().Str("id", id).Strs("used_by", usedBy).Msg("Service is in use")
		return c.JSON(http.StatusConflict, map[string]interface{}{
			"error":   "Service is in use by other resources and cannot be deleted",
			"used_by": usedBy,
		})
	}

	// Delete the service
	if err := h.Store.DeleteService(id); err != nil {
		logger.Error().Err(err).Str("id", id).Msg("Failed to delete service")

		if store.IsResourceInUse(err) {
			// This is a fallback for any resource-in-use errors not caught by the explicit check
			return c.JSON(http.StatusConflict, map[string]string{
				"error": "Service is in use by other resources and cannot be deleted",
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete service",
		})
	}

	logger.Info().Str("id", id).Msg("Service deleted")

	// Return success response
	response := models.ResourceResponse{
		ID:      id,
		Created: false,
		Updated: false,
		Deleted: true,
	}

	return c.JSON(http.StatusOK, response)
}

// validateServiceConfiguration checks if the service has a valid configuration
func validateServiceConfiguration(service *models.Service) error {
	// Simple URL service
	if service.URL != "" {
		return nil
	}

	// Load balancer service
	if service.LoadBalancer != nil {
		if len(service.LoadBalancer.Servers) == 0 {
			return fmt.Errorf("load balancer service must have at least one server")
		}
		for _, server := range service.LoadBalancer.Servers {
			if server.URL == "" {
				return fmt.Errorf("server URL is required for load balancer servers")
			}
		}
		return nil
	}

	// Weighted service
	if service.Weighted != nil {
		if len(service.Weighted.Services) == 0 {
			return fmt.Errorf("weighted service must have at least one service")
		}
		return nil
	}

	// Mirroring service
	if service.Mirroring != nil {
		if service.Mirroring.Service.ID == "" {
			return fmt.Errorf("main service ID is required for mirroring service")
		}
		return nil
	}

	// Failover service
	if service.Failover != nil {
		if service.Failover.Service.ID == "" {
			return fmt.Errorf("main service ID is required for failover service")
		}
		if service.Failover.Fallback.ID == "" {
			return fmt.Errorf("fallback service ID is required for failover service")
		}
		return nil
	}

	return fmt.Errorf("service must have either URL, LoadBalancer, Weighted, Mirroring, or Failover configuration")
}
