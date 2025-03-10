// internal/api/routes/routes.go
package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/sistemica/traefik-manager/internal/api/handlers"
	"github.com/sistemica/traefik-manager/internal/config"
	"github.com/sistemica/traefik-manager/internal/store"
)

// RegisterRoutes sets up all API routes
func RegisterRoutes(e *echo.Echo, s store.Store, basePath string, cfg *config.Config) {
	// Create handler instances
	middlewareHandler := handlers.NewMiddlewareHandler(s)
	routerHandler := handlers.NewRouterHandler(s)
	serviceHandler := handlers.NewServiceHandler(s)
	healthHandler := handlers.NewHealthHandler(s, "1.0.0")

	// API group with base path
	api := e.Group(basePath)

	// Health check - always public
	api.GET("/health", healthHandler.Check)

	// Traefik provider endpoint - with custom auth
	providerAuth := cfg.Provider.Auth
	if providerAuth == nil && cfg.Auth.Enabled {
		// If global auth is enabled but no specific provider auth,
		// the provider endpoint is public (excluded from auth)
		providerHandler := handlers.NewProviderHandler(s)
		e.GET(cfg.Provider.ProviderPath, providerHandler.GetConfig)
	} else {
		// Either provider-specific auth or no auth at all
		providerHandlerWithAuth := handlers.NewProviderHandlerWithAuth(s, providerAuth)
		e.GET(cfg.Provider.ProviderPath, providerHandlerWithAuth.GetConfigWithAuth)
	}

	// Middlewares
	middlewares := api.Group("/middlewares")
	middlewares.GET("", middlewareHandler.List)
	middlewares.POST("", middlewareHandler.Create)
	middlewares.GET("/:id", middlewareHandler.Get)
	middlewares.PUT("/:id", middlewareHandler.Update)
	middlewares.DELETE("/:id", middlewareHandler.Delete)

	// Routers
	routers := api.Group("/routers")
	routers.GET("", routerHandler.List)
	routers.POST("", routerHandler.Create)
	routers.GET("/:id", routerHandler.Get)
	routers.PUT("/:id", routerHandler.Update)
	routers.DELETE("/:id", routerHandler.Delete)

	// Services
	services := api.Group("/services")
	services.GET("", serviceHandler.List)
	services.POST("", serviceHandler.Create)
	services.GET("/:id", serviceHandler.Get)
	services.PUT("/:id", serviceHandler.Update)
	services.DELETE("/:id", serviceHandler.Delete)
}
