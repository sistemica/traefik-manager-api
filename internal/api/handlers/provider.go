// internal/api/handlers/provider.go
package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sistemica/traefik-manager/internal/config"
	"github.com/sistemica/traefik-manager/internal/logger"
	"github.com/sistemica/traefik-manager/internal/models"
	"github.com/sistemica/traefik-manager/internal/store"
	"github.com/sistemica/traefik-manager/internal/traefik"
)

// ProviderHandler handles Traefik provider endpoint requests
type ProviderHandler struct {
	BaseHandler
}

// NewProviderHandler creates a new ProviderHandler
func NewProviderHandler(store store.Store) *ProviderHandler {
	return &ProviderHandler{
		BaseHandler: NewBaseHandler(store),
	}
}

// ProviderHandlerWithAuth is an extension of ProviderHandler that includes auth settings
type ProviderHandlerWithAuth struct {
	BaseHandler
	AuthConfig *config.Auth
}

// NewProviderHandlerWithAuth creates a new ProviderHandler with auth settings
func NewProviderHandlerWithAuth(store store.Store, authConfig *config.Auth) *ProviderHandlerWithAuth {
	return &ProviderHandlerWithAuth{
		BaseHandler: NewBaseHandler(store),
		AuthConfig:  authConfig,
	}
}

// GetConfigWithAuth handles the provider endpoint with direct auth check
func (h *ProviderHandlerWithAuth) GetConfigWithAuth(c echo.Context) error {
	// Handle authentication if enabled
	if h.AuthConfig != nil && h.AuthConfig.Enabled {
		// Check API key
		apiKey := c.Request().Header.Get(h.AuthConfig.HeaderName)
		if apiKey == "" {
			logger.Warn().Str("path", c.Request().URL.Path).Msg("Missing API key")
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "API key missing",
			})
		}

		if apiKey != h.AuthConfig.Key {
			logger.Warn().Str("path", c.Request().URL.Path).Msg("Invalid API key")
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Invalid API key",
			})
		}
	}

	// After authentication succeeds or if auth is disabled, serve the configuration
	logger.Debug().Msg("Traefik requesting configuration")

	// Get all resources from store
	routers, err := h.Store.ListRouters()
	if err != nil {
		logger.Error().Err(err).Msg("Failed to list routers")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve configuration",
		})
	}

	services, err := h.Store.ListServices()
	if err != nil {
		logger.Error().Err(err).Msg("Failed to list services")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve configuration",
		})
	}

	middlewares, err := h.Store.ListMiddlewares()
	if err != nil {
		logger.Error().Err(err).Msg("Failed to list middlewares")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve configuration",
		})
	}

	// Convert to Traefik configuration
	config := convertToTraefikConfig(routers, services, middlewares)

	logger.Debug().Int("routers", len(routers)).Int("services", len(services)).Int("middlewares", len(middlewares)).Msg("Configuration served to Traefik")

	return c.JSON(http.StatusOK, config)
}

// GetConfig handles the provider endpoint that Traefik polls for configuration
func (h *ProviderHandler) GetConfig(c echo.Context) error {
	logger.Debug().Msg("Traefik requesting configuration")

	// Get all resources from store
	routers, err := h.Store.ListRouters()
	if err != nil {
		logger.Error().Err(err).Msg("Failed to list routers")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve configuration",
		})
	}

	services, err := h.Store.ListServices()
	if err != nil {
		logger.Error().Err(err).Msg("Failed to list services")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve configuration",
		})
	}

	middlewares, err := h.Store.ListMiddlewares()
	if err != nil {
		logger.Error().Err(err).Msg("Failed to list middlewares")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve configuration",
		})
	}

	// Convert to Traefik configuration
	config := convertToTraefikConfig(routers, services, middlewares)

	logger.Debug().Int("routers", len(routers)).Int("services", len(services)).Int("middlewares", len(middlewares)).Msg("Configuration served to Traefik")

	return c.JSON(http.StatusOK, config)
}

// convertToTraefikConfig converts internal models to Traefik's dynamic configuration
func convertToTraefikConfig(routers []models.Router, services []models.Service, middlewares []models.Middleware) *traefik.DynamicConfig {
	// Initialize Traefik dynamic config
	config := &traefik.DynamicConfig{
		HTTP: &traefik.HTTPConfiguration{
			Routers:     make(map[string]*traefik.Router),
			Services:    make(map[string]*traefik.Service),
			Middlewares: make(map[string]*traefik.Middleware),
		},
	}

	// Convert middlewares
	for _, mw := range middlewares {
		traefikMw := convertMiddleware(mw)
		config.HTTP.Middlewares[mw.ID] = traefikMw
	}

	// Convert services
	for _, svc := range services {
		traefikSvc := convertService(svc)
		config.HTTP.Services[svc.ID] = traefikSvc
	}

	// Convert routers
	for _, router := range routers {
		traefikRouter := convertRouter(router)
		config.HTTP.Routers[router.ID] = traefikRouter
	}

	return config
}

// convertRouter converts a models.Router to a traefik.Router
func convertRouter(router models.Router) *traefik.Router {
	traefikRouter := &traefik.Router{
		EntryPoints: router.EntryPoints,
		Rule:        router.Rule,
		Priority:    router.Priority,
		Service:     router.Service.ID,
	}

	// Convert middlewares (just need the IDs)
	if router.Middlewares != nil && len(router.Middlewares) > 0 {
		traefikRouter.Middlewares = make([]string, len(router.Middlewares))
		for i, mw := range router.Middlewares {
			traefikRouter.Middlewares[i] = mw.ID
		}
	}

	// Convert TLS configuration if present
	if router.TLS != nil {
		traefikRouter.TLS = &traefik.RouterTLS{
			Options:      router.TLS.Options,
			CertResolver: router.TLS.CertResolver,
		}

		// Convert domains if present
		if router.TLS.Domains != nil && len(router.TLS.Domains) > 0 {
			traefikRouter.TLS.Domains = make([]*traefik.Domain, len(router.TLS.Domains))
			for i, domain := range router.TLS.Domains {
				traefikRouter.TLS.Domains[i] = &traefik.Domain{
					Main: domain.Main,
					Sans: domain.Sans,
				}
			}
		}
	}

	// Convert observability settings if present
	if router.Observability != nil {
		traefikRouter.Observability = &traefik.Observability{
			AccessLogs: router.Observability.AccessLogs,
			Tracing:    router.Observability.Tracing,
			Metrics:    router.Observability.Metrics,
		}
	}

	return traefikRouter
}

// convertService converts a models.Service to a traefik.Service
func convertService(service models.Service) *traefik.Service {
	traefikService := &traefik.Service{}

	// Handle simple URL-based service
	if service.URL != "" {
		// For URL-based services, create a LoadBalancer with a single server
		traefikService.LoadBalancer = &traefik.LoadBalancerService{
			Servers: []traefik.Server{
				{
					URL: service.URL,
				},
			},
		}
		return traefikService
	}

	// Handle LoadBalancer service
	if service.LoadBalancer != nil {
		loadBalancer := &traefik.LoadBalancerService{
			ServersTransport: service.LoadBalancer.ServersTransport,
		}

		// Convert servers
		if service.LoadBalancer.Servers != nil {
			loadBalancer.Servers = make([]traefik.Server, len(service.LoadBalancer.Servers))
			for i, server := range service.LoadBalancer.Servers {
				preservePath := false
				if server.PreservePath {
					preservePath = true
				}

				weight := 1
				if server.Weight > 0 {
					weight = server.Weight
				}

				loadBalancer.Servers[i] = traefik.Server{
					URL:          server.URL,
					Weight:       &weight,
					PreservePath: &preservePath,
				}
			}
		}

		// Convert health check if present
		if service.LoadBalancer.HealthCheck != nil {
			loadBalancer.HealthCheck = convertHealthCheck(service.LoadBalancer.HealthCheck)
		}

		// Convert sticky configuration if present
		if service.LoadBalancer.Sticky != nil {
			loadBalancer.Sticky = convertSticky(service.LoadBalancer.Sticky)
		}

		// Convert passHostHeader if specified
		passHostHeader := true // Default value in Traefik
		loadBalancer.PassHostHeader = &passHostHeader

		// Convert response forwarding if present
		if service.LoadBalancer.ResponseForwarding != nil {
			loadBalancer.ResponseForwarding = &traefik.ResponseForwarding{
				FlushInterval: string(service.LoadBalancer.ResponseForwarding.FlushInterval),
			}
		}

		traefikService.LoadBalancer = loadBalancer
		return traefikService
	}

	// Handle Weighted service
	if service.Weighted != nil {
		weighted := &traefik.WeightedService{}

		// Convert services
		if service.Weighted.Services != nil {
			weighted.Services = make([]traefik.WeightedServiceItem, len(service.Weighted.Services))
			for i, item := range service.Weighted.Services {
				weighted.Services[i] = traefik.WeightedServiceItem{
					Name:   item.Name.ID,
					Weight: item.Weight,
				}
			}
		}

		// Convert health check if present
		if service.Weighted.HealthCheck != nil {
			weighted.HealthCheck = convertHealthCheck(service.Weighted.HealthCheck)
		}

		// Convert sticky configuration if present
		if service.Weighted.Sticky != nil {
			weighted.Sticky = convertSticky(service.Weighted.Sticky)
		}

		traefikService.Weighted = weighted
		return traefikService
	}

	// Handle Mirroring service
	if service.Mirroring != nil {
		mirroring := &traefik.MirroringService{
			Service: service.Mirroring.Service.ID,
		}

		// Set mirrorBody default or value
		mirrorBody := true
		if service.Mirroring.MirrorBody {
			mirrorBody = true
		}
		mirroring.MirrorBody = &mirrorBody

		// Set maxBodySize if present
		if service.Mirroring.MaxBodySize > 0 {
			maxBodySize := service.Mirroring.MaxBodySize
			mirroring.MaxBodySize = &maxBodySize
		}

		// Convert mirrors
		if service.Mirroring.Mirrors != nil {
			mirroring.Mirrors = make([]traefik.MirrorServiceItem, len(service.Mirroring.Mirrors))
			for i, mirror := range service.Mirroring.Mirrors {
				mirroring.Mirrors[i] = traefik.MirrorServiceItem{
					Name:    mirror.Name.ID,
					Percent: mirror.Percent,
				}
			}
		}

		// Convert health check if present
		if service.Mirroring.HealthCheck != nil {
			mirroring.HealthCheck = convertHealthCheck(service.Mirroring.HealthCheck)
		}

		traefikService.Mirroring = mirroring
		return traefikService
	}

	// Handle Failover service
	if service.Failover != nil {
		failover := &traefik.FailoverService{
			Service:  service.Failover.Service.ID,
			Fallback: service.Failover.Fallback.ID,
		}

		// Convert health check if present
		if service.Failover.HealthCheck != nil {
			failover.HealthCheck = convertHealthCheck(service.Failover.HealthCheck)
		}

		traefikService.Failover = failover
		return traefikService
	}

	return traefikService
}

// convertMiddleware converts a models.Middleware to a traefik.Middleware
func convertMiddleware(middleware models.Middleware) *traefik.Middleware {
	traefikMiddleware := &traefik.Middleware{}

	// Set the appropriate middleware configuration based on type
	switch middleware.Type {
	case "redirectScheme":
		// Try to convert the config to the correct type
		// The issue might be that config is stored as a generic map[string]interface{} in JSON
		configMap, ok := middleware.Config.(map[string]interface{})
		if ok {
			schemeConfig := &traefik.RedirectSchemeConfig{}

			// Extract scheme
			if scheme, ok := configMap["scheme"].(string); ok {
				schemeConfig.Scheme = scheme
			}

			// Extract permanent flag
			if permanent, ok := configMap["permanent"].(bool); ok {
				schemeConfig.Permanent = permanent
			}

			traefikMiddleware.RedirectScheme = schemeConfig
		}

	case "addPrefix":
		configMap, ok := middleware.Config.(map[string]interface{})
		if ok {
			addPrefixConfig := &traefik.AddPrefixConfig{}

			if prefix, ok := configMap["prefix"].(string); ok {
				addPrefixConfig.Prefix = prefix
			}

			traefikMiddleware.AddPrefix = addPrefixConfig
		}

	case "basicAuth":
		configMap, ok := middleware.Config.(map[string]interface{})
		if ok {
			basicAuthConfig := &traefik.BasicAuthConfig{}

			if users, ok := configMap["users"].([]interface{}); ok {
				strUsers := make([]string, 0, len(users))
				for _, u := range users {
					if strUser, ok := u.(string); ok {
						strUsers = append(strUsers, strUser)
					}
				}
				basicAuthConfig.Users = strUsers
			}

			if usersFile, ok := configMap["usersFile"].(string); ok {
				basicAuthConfig.UsersFile = usersFile
			}

			if realm, ok := configMap["realm"].(string); ok {
				basicAuthConfig.Realm = realm
			}

			if removeHeader, ok := configMap["removeHeader"].(bool); ok {
				basicAuthConfig.RemoveHeader = &removeHeader
			}

			if headerField, ok := configMap["headerField"].(string); ok {
				basicAuthConfig.HeaderField = headerField
			}

			traefikMiddleware.BasicAuth = basicAuthConfig
		}

	case "stripPrefix":
		configMap, ok := middleware.Config.(map[string]interface{})
		if ok {
			stripPrefixConfig := &traefik.StripPrefixConfig{}

			if prefixes, ok := configMap["prefixes"].([]interface{}); ok {
				strPrefixes := make([]string, 0, len(prefixes))
				for _, p := range prefixes {
					if strPrefix, ok := p.(string); ok {
						strPrefixes = append(strPrefixes, strPrefix)
					}
				}
				stripPrefixConfig.Prefixes = strPrefixes
			}

			if forceSlash, ok := configMap["forceSlash"].(bool); ok {
				stripPrefixConfig.ForceSlash = forceSlash
			}

			traefikMiddleware.StripPrefix = stripPrefixConfig
		}

		// Add more middleware types as needed...
	}

	return traefikMiddleware
}

// convertHealthCheck converts a models.HealthCheck to a traefik.HealthCheck
func convertHealthCheck(healthCheck *models.HealthCheck) *traefik.HealthCheck {
	if healthCheck == nil {
		return nil
	}

	traefikHealthCheck := &traefik.HealthCheck{
		Scheme:   healthCheck.Scheme,
		Mode:     healthCheck.Mode,
		Path:     healthCheck.Path,
		Method:   healthCheck.Method,
		Status:   healthCheck.Status,
		Port:     healthCheck.Port,
		Interval: string(healthCheck.Interval),
		Timeout:  string(healthCheck.Timeout),
		Hostname: healthCheck.Hostname,
		Headers:  healthCheck.Headers,
	}

	// Convert followRedirects if specified
	followRedirects := false
	if healthCheck.FollowRedirects {
		followRedirects = true
	}
	traefikHealthCheck.FollowRedirects = &followRedirects

	return traefikHealthCheck
}

// convertSticky converts a models.Sticky to a traefik.Sticky
func convertSticky(sticky *models.Sticky) *traefik.Sticky {
	if sticky == nil || sticky.Cookie == nil {
		return nil
	}

	traefikSticky := &traefik.Sticky{
		Cookie: &traefik.StickyCooke{
			Name:     sticky.Cookie.Name,
			Path:     sticky.Cookie.Path,
			SameSite: sticky.Cookie.SameSite,
			MaxAge:   sticky.Cookie.MaxAge,
		},
	}

	// Convert secure if specified
	secure := false
	if sticky.Cookie.Secure {
		secure = true
	}
	traefikSticky.Cookie.Secure = &secure

	// Convert httpOnly if specified
	httpOnly := false
	if sticky.Cookie.HTTPOnly {
		httpOnly = true
	}
	traefikSticky.Cookie.HTTPOnly = &httpOnly

	return traefikSticky
}
