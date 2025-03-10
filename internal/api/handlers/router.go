// internal/api/handlers/router.go
package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sistemica/traefik-manager/internal/logger"
	"github.com/sistemica/traefik-manager/internal/models"
	"github.com/sistemica/traefik-manager/internal/store"
)

// RouterHandler handles router-related requests
type RouterHandler struct {
	BaseHandler
}

// NewRouterHandler creates a new RouterHandler
func NewRouterHandler(store store.Store) *RouterHandler {
	return &RouterHandler{
		BaseHandler: NewBaseHandler(store),
	}
}

// List handles the GET /routers endpoint to list all routers
func (h *RouterHandler) List(c echo.Context) error {
	logger.Debug().Msg("Listing routers")

	routers, err := h.Store.ListRouters()
	if err != nil {
		logger.Error().Err(err).Msg("Failed to list routers")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to list routers",
		})
	}

	return c.JSON(http.StatusOK, routers)
}

// Get handles the GET /routers/:id endpoint to get a specific router
func (h *RouterHandler) Get(c echo.Context) error {
	id := c.Param("id")
	logger.Debug().Str("id", id).Msg("Getting router")

	router, err := h.Store.GetRouter(id)
	if err != nil {
		if store.IsNotFound(err) {
			logger.Debug().Str("id", id).Msg("Router not found")
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Router not found",
			})
		}
		logger.Error().Err(err).Str("id", id).Msg("Failed to get router")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get router",
		})
	}

	return c.JSON(http.StatusOK, router)
}

// Create handles the POST /routers endpoint to create a new router
func (h *RouterHandler) Create(c echo.Context) error {
	logger.Debug().Msg("Creating router")

	// Create a timeout context
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	// First try to parse the request as a map to check the format of 'service'
	var requestData map[string]interface{}
	if err := c.Bind(&requestData); err != nil {
		logger.Warn().Err(err).Msg("Invalid router data")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid router data",
		})
	}

	// Create router object
	var router models.Router

	// Set ID, rule, and entryPoints
	if id, ok := requestData["id"].(string); ok {
		router.ID = id
	}

	if rule, ok := requestData["rule"].(string); ok {
		router.Rule = rule
	}

	if entryPoints, ok := requestData["entryPoints"].([]interface{}); ok {
		router.EntryPoints = make([]string, len(entryPoints))
		for i, ep := range entryPoints {
			if epStr, ok := ep.(string); ok {
				router.EntryPoints[i] = epStr
			}
		}
	}

	// Handle service field - can be either a string (ID) or an object
	serviceField := requestData["service"]
	if serviceID, ok := serviceField.(string); ok {
		// If service is a string, create a Service with just the ID
		router.Service = models.Service{
			ID: serviceID,
		}
	} else if serviceMap, ok := serviceField.(map[string]interface{}); ok {
		// If service is an object, extract the ID
		if serviceID, ok := serviceMap["id"].(string); ok {
			router.Service = models.Service{
				ID: serviceID,
			}
		}
	}

	// Handle middlewares field - can be an array of strings or objects
	if middlewaresField, ok := requestData["middlewares"].([]interface{}); ok {
		router.Middlewares = make([]models.Middleware, 0, len(middlewaresField))

		for _, mw := range middlewaresField {
			if mwID, ok := mw.(string); ok {
				// If middleware is a string, create a Middleware with just the ID
				router.Middlewares = append(router.Middlewares, models.Middleware{
					ID: mwID,
				})
			} else if mwMap, ok := mw.(map[string]interface{}); ok {
				// If middleware is an object, extract the ID
				if mwID, ok := mwMap["id"].(string); ok {
					router.Middlewares = append(router.Middlewares, models.Middleware{
						ID: mwID,
					})
				}
			}
		}
	}

	// Handle Priority
	if priority, ok := requestData["priority"].(float64); ok {
		router.Priority = int(priority)
	}

	// Handle TLS if present
	if tlsField, ok := requestData["tls"].(map[string]interface{}); ok {
		tls := &models.RouterTLS{}

		if options, ok := tlsField["options"].(string); ok {
			tls.Options = options
		}

		if certResolver, ok := tlsField["certResolver"].(string); ok {
			tls.CertResolver = certResolver
		}

		if domainsField, ok := tlsField["domains"].([]interface{}); ok {
			domains := make([]models.Domain, 0, len(domainsField))

			for _, d := range domainsField {
				if domainMap, ok := d.(map[string]interface{}); ok {
					domain := models.Domain{}

					if main, ok := domainMap["main"].(string); ok {
						domain.Main = main
					}

					if sansField, ok := domainMap["sans"].([]interface{}); ok {
						sans := make([]string, 0, len(sansField))
						for _, s := range sansField {
							if san, ok := s.(string); ok {
								sans = append(sans, san)
							}
						}
						domain.Sans = sans
					}

					domains = append(domains, domain)
				}
			}

			tls.Domains = domains
		}

		router.TLS = tls
	}

	// Handle Observability if present
	if obsField, ok := requestData["observability"].(map[string]interface{}); ok {
		obs := &models.Observability{}

		if accessLogs, ok := obsField["accessLogs"].(bool); ok {
			obs.AccessLogs = accessLogs
		}

		if tracing, ok := obsField["tracing"].(bool); ok {
			obs.Tracing = tracing
		}

		if metrics, ok := obsField["metrics"].(bool); ok {
			obs.Metrics = metrics
		}

		router.Observability = obs
	}

	// Validate router
	if router.ID == "" {
		logger.Warn().Msg("Router ID is required")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Router ID is required",
		})
	}

	if router.Rule == "" {
		logger.Warn().Msg("Router rule is required")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Router rule is required",
		})
	}

	// Check if service exists
	if router.Service.ID == "" {
		logger.Warn().Msg("Router service ID is required")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Router service ID is required",
		})
	}

	// Use channels for timeout-safe operations
	serviceExistsChan := make(chan bool, 1)
	serviceErrorChan := make(chan error, 1)

	go func() {
		exists, err := h.Store.ServiceExists(router.Service.ID)
		serviceExistsChan <- exists
		serviceErrorChan <- err
	}()

	// Wait for the service check with timeout
	select {
	case <-ctx.Done():
		logger.Error().Str("id", router.ID).Msg("Timeout checking if service exists")
		return c.JSON(http.StatusGatewayTimeout, map[string]string{
			"error": "Operation timed out while checking service",
		})
	case err := <-serviceErrorChan:
		if err != nil {
			logger.Error().Err(err).Str("service_id", router.Service.ID).Msg("Failed to check if service exists")
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to check if service exists",
			})
		}
	}

	serviceExists := <-serviceExistsChan
	if !serviceExists {
		logger.Warn().Str("service_id", router.Service.ID).Msg("Service not found")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Referenced service does not exist",
		})
	}

	// Check if middlewares exist
	if router.Middlewares != nil {
		for _, middleware := range router.Middlewares {
			middlewareExistsChan := make(chan bool, 1)
			middlewareErrorChan := make(chan error, 1)

			go func(mwID string) {
				exists, err := h.Store.MiddlewareExists(mwID)
				middlewareExistsChan <- exists
				middlewareErrorChan <- err
			}(middleware.ID)

			// Wait for the middleware check with timeout
			select {
			case <-ctx.Done():
				logger.Error().Str("id", router.ID).Msg("Timeout checking if middleware exists")
				return c.JSON(http.StatusGatewayTimeout, map[string]string{
					"error": "Operation timed out while checking middleware",
				})
			case err := <-middlewareErrorChan:
				if err != nil {
					logger.Error().Err(err).Str("middleware_id", middleware.ID).Msg("Failed to check if middleware exists")
					return c.JSON(http.StatusInternalServerError, map[string]string{
						"error": "Failed to check if middleware exists",
					})
				}
			}

			middlewareExists := <-middlewareExistsChan
			if !middlewareExists {
				logger.Warn().Str("middleware_id", middleware.ID).Msg("Middleware not found")
				return c.JSON(http.StatusBadRequest, map[string]string{
					"error": "Referenced middleware does not exist: " + middleware.ID,
				})
			}
		}
	}

	// Check if router already exists
	routerExistsChan := make(chan bool, 1)
	routerErrorChan := make(chan error, 1)

	go func() {
		exists, err := h.Store.RouterExists(router.ID)
		routerExistsChan <- exists
		routerErrorChan <- err
	}()

	// Wait for the router check with timeout
	select {
	case <-ctx.Done():
		logger.Error().Str("id", router.ID).Msg("Timeout checking if router exists")
		return c.JSON(http.StatusGatewayTimeout, map[string]string{
			"error": "Operation timed out while checking router existence",
		})
	case err := <-routerErrorChan:
		if err != nil {
			logger.Error().Err(err).Str("id", router.ID).Msg("Failed to check if router exists")
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to check if router exists",
			})
		}
	}

	routerExists := <-routerExistsChan
	if routerExists {
		logger.Warn().Str("id", router.ID).Msg("Router already exists")
		return c.JSON(http.StatusConflict, map[string]string{
			"error": "Router already exists",
		})
	}

	// Create the router
	createChan := make(chan error, 1)
	go func() {
		createChan <- h.Store.CreateRouter(&router)
	}()

	// Wait for creation with timeout
	select {
	case <-ctx.Done():
		logger.Error().Str("id", router.ID).Msg("Timeout creating router")
		return c.JSON(http.StatusGatewayTimeout, map[string]string{
			"error": "Operation timed out while creating router",
		})
	case err := <-createChan:
		if err != nil {
			logger.Error().Err(err).Str("id", router.ID).Msg("Failed to create router")
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to create router",
			})
		}
	}

	logger.Info().Str("id", router.ID).Msg("Router created")

	// Return success response with 201 Created status
	response := models.ResourceResponse{
		ID:      router.ID,
		Created: true,
		Updated: false,
		Deleted: false,
	}

	return c.JSON(http.StatusCreated, response)
}

// Update handles the PUT /routers/:id endpoint to update a router
func (h *RouterHandler) Update(c echo.Context) error {
	id := c.Param("id")
	logger.Debug().Str("id", id).Msg("Updating router")

	// First try to parse the request as a map to check the format of 'service'
	var requestData map[string]interface{}
	if err := c.Bind(&requestData); err != nil {
		logger.Warn().Err(err).Str("id", id).Msg("Invalid router data")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid router data",
		})
	}

	// Create router object
	var router models.Router
	router.ID = id

	// Set rule and entryPoints
	if rule, ok := requestData["rule"].(string); ok {
		router.Rule = rule
	}

	if entryPoints, ok := requestData["entryPoints"].([]interface{}); ok {
		router.EntryPoints = make([]string, len(entryPoints))
		for i, ep := range entryPoints {
			if epStr, ok := ep.(string); ok {
				router.EntryPoints[i] = epStr
			}
		}
	}

	// Handle service field - can be either a string (ID) or an object
	serviceField := requestData["service"]
	if serviceID, ok := serviceField.(string); ok {
		// If service is a string, create a Service with just the ID
		router.Service = models.Service{
			ID: serviceID,
		}
	} else if serviceMap, ok := serviceField.(map[string]interface{}); ok {
		// If service is an object, extract the ID
		if serviceID, ok := serviceMap["id"].(string); ok {
			router.Service = models.Service{
				ID: serviceID,
			}
		}
	}

	// Handle middlewares field - can be an array of strings or objects
	if middlewaresField, ok := requestData["middlewares"].([]interface{}); ok {
		router.Middlewares = make([]models.Middleware, 0, len(middlewaresField))

		for _, mw := range middlewaresField {
			if mwID, ok := mw.(string); ok {
				// If middleware is a string, create a Middleware with just the ID
				router.Middlewares = append(router.Middlewares, models.Middleware{
					ID: mwID,
				})
			} else if mwMap, ok := mw.(map[string]interface{}); ok {
				// If middleware is an object, extract the ID
				if mwID, ok := mwMap["id"].(string); ok {
					router.Middlewares = append(router.Middlewares, models.Middleware{
						ID: mwID,
					})
				}
			}
		}
	}

	// Validate required fields
	if router.Rule == "" {
		logger.Warn().Msg("Router rule is required")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Router rule is required",
		})
	}

	// Update the router with validation in a single operation
	err := h.Store.UpdateRouter(id, &router)
	if err != nil {
		if store.IsNotFound(err) {
			logger.Warn().Str("id", id).Msg("Router not found")
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Router not found",
			})
		}

		logger.Error().Err(err).Str("id", id).Msg("Failed to update router")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	logger.Info().Str("id", id).Msg("Router updated")

	// Return success response
	response := models.ResourceResponse{
		ID:      id,
		Created: false,
		Updated: true,
		Deleted: false,
	}

	return c.JSON(http.StatusOK, response)
}

// Delete handles the DELETE /routers/:id endpoint to delete a router
func (h *RouterHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	logger.Debug().Str("id", id).Msg("Deleting router")

	// Create a timeout context
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	// Check if router exists
	routerExistsChan := make(chan bool, 1)
	routerErrorChan := make(chan error, 1)

	go func() {
		exists, err := h.Store.RouterExists(id)
		routerExistsChan <- exists
		routerErrorChan <- err
	}()

	// Wait for the router check with timeout
	select {
	case <-ctx.Done():
		logger.Error().Str("id", id).Msg("Timeout checking if router exists")
		return c.JSON(http.StatusGatewayTimeout, map[string]string{
			"error": "Operation timed out while checking router existence",
		})
	case err := <-routerErrorChan:
		if err != nil {
			logger.Error().Err(err).Str("id", id).Msg("Failed to check if router exists")
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to check if router exists",
			})
		}
	}

	routerExists := <-routerExistsChan
	if !routerExists {
		logger.Warn().Str("id", id).Msg("Router not found")
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Router not found",
		})
	}

	// Check if router is in use
	routerInUseChan := make(chan bool, 1)
	routerUsedByChan := make(chan []string, 1)
	routerInUseErrorChan := make(chan error, 1)

	go func() {
		inUse, usedBy, err := h.Store.RouterInUse(id)
		routerInUseChan <- inUse
		routerUsedByChan <- usedBy
		routerInUseErrorChan <- err
	}()

	// Wait for the in-use check with timeout
	select {
	case <-ctx.Done():
		logger.Error().Str("id", id).Msg("Timeout checking if router is in use")
		return c.JSON(http.StatusGatewayTimeout, map[string]string{
			"error": "Operation timed out while checking router dependencies",
		})
	case err := <-routerInUseErrorChan:
		if err != nil {
			logger.Error().Err(err).Str("id", id).Msg("Failed to check if router is in use")
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to check if router is in use",
			})
		}
	}

	routerInUse := <-routerInUseChan
	routerUsedBy := <-routerUsedByChan

	if routerInUse {
		logger.Warn().Str("id", id).Strs("used_by", routerUsedBy).Msg("Router is in use")
		return c.JSON(http.StatusConflict, map[string]interface{}{
			"error":   "Router is in use by other resources and cannot be deleted",
			"used_by": routerUsedBy,
		})
	}

	// Delete the router
	deleteChan := make(chan error, 1)
	go func() {
		deleteChan <- h.Store.DeleteRouter(id)
	}()

	// Wait for deletion with timeout
	select {
	case <-ctx.Done():
		logger.Error().Str("id", id).Msg("Timeout deleting router")
		return c.JSON(http.StatusGatewayTimeout, map[string]string{
			"error": "Operation timed out while deleting router",
		})
	case err := <-deleteChan:
		if err != nil {
			logger.Error().Err(err).Str("id", id).Msg("Failed to delete router")

			if store.IsResourceInUse(err) {
				// This is a fallback for any resource-in-use errors not caught by the explicit check
				return c.JSON(http.StatusConflict, map[string]string{
					"error": "Router is in use by other resources and cannot be deleted",
				})
			}

			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to delete router",
			})
		}
	}

	logger.Info().Str("id", id).Msg("Router deleted")

	// Return success response
	response := models.ResourceResponse{
		ID:      id,
		Created: false,
		Updated: false,
		Deleted: true,
	}

	return c.JSON(http.StatusOK, response)
}
