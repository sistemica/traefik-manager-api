// internal/api/server/server.go
package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sistemica/traefik-manager/internal/api/routes"
	"github.com/sistemica/traefik-manager/internal/config"
	"github.com/sistemica/traefik-manager/internal/logger"
	customMiddleware "github.com/sistemica/traefik-manager/internal/middleware"
	"github.com/sistemica/traefik-manager/internal/store"
)

// Server represents the HTTP server
type Server struct {
	echo       *echo.Echo
	config     *config.Config
	store      store.Store
	httpServer *http.Server
}

// New creates a new server instance
func New(cfg *config.Config, store store.Store) *Server {
	e := echo.New()
	e.HideBanner = true

	return &Server{
		echo:   e,
		config: cfg,
		store:  store,
	}
}

// Setup configures the server
func (s *Server) Setup() {
	// Setup middleware
	s.echo.Use(middleware.RequestID())
	s.echo.Use(customMiddleware.Logger())
	s.echo.Use(middleware.Recover())
	s.echo.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 30 * time.Second,
	}))

	// Setup CORS
	s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     s.config.Cors.AllowedOrigins,
		AllowMethods:     s.config.Cors.AllowedMethods,
		AllowHeaders:     s.config.Cors.AllowedHeaders,
		AllowCredentials: s.config.Cors.AllowCredentials,
		MaxAge:           int(s.config.Cors.MaxAge),
	}))

	// Determine excluded paths based on config
	excludedPaths := []string{"/health"}

	// If global auth is enabled, exclude provider path
	if s.config.Auth.Enabled {
		excludedPaths = append(excludedPaths, s.config.Provider.ProviderPath)
	}

	// API authentication middleware (for all other endpoints)
	if s.config.Auth.Enabled {
		s.echo.Use(customMiddleware.Auth(customMiddleware.AuthOptions{
			Enabled:      s.config.Auth.Enabled,
			HeaderName:   s.config.Auth.HeaderName,
			Key:          s.config.Auth.Key,
			ExcludePaths: excludedPaths,
		}))
	}

	// Register routes with all configuration context
	routes.RegisterRoutes(s.echo, s.store, s.config.Server.BasePath, s.config)

	// Configure HTTP server
	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port),
		ReadTimeout:  s.config.Server.ReadTimeout,
		WriteTimeout: s.config.Server.WriteTimeout,
	}
}

// Start starts the server
func (s *Server) Start() error {
	logger.Info().Str("address", s.httpServer.Addr).Msg("Starting server")
	return s.echo.StartServer(s.httpServer)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	logger.Info().Msg("Shutting down server")
	return s.echo.Shutdown(ctx)
}

// GetEcho returns the echo instance for testing
func (s *Server) GetEcho() *echo.Echo {
	return s.echo
}
