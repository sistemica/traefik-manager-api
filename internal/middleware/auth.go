// internal/middleware/auth.go
package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sistemica/traefik-manager/internal/logger"
)

// AuthOptions represents options for the authentication middleware
type AuthOptions struct {
	Enabled      bool
	HeaderName   string
	Key          string
	ExcludePaths []string
}

// Auth creates a middleware for API key authentication
func Auth(opts AuthOptions) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip auth if disabled
			if !opts.Enabled {
				return next(c)
			}

			// Skip auth for excluded paths
			path := c.Request().URL.Path
			for _, excludePath := range opts.ExcludePaths {
				if strings.HasPrefix(path, excludePath) {
					return next(c)
				}
			}

			// Check auth header
			apiKey := c.Request().Header.Get(opts.HeaderName)
			if apiKey == "" {
				logger.Warn().Str("path", path).Msg("Missing API key")
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"code":    "unauthorized",
					"message": "API key missing",
				})
			}

			// Validate key
			if apiKey != opts.Key {
				logger.Warn().Str("path", path).Msg("Invalid API key")
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"code":    "unauthorized",
					"message": "Invalid API key",
				})
			}

			return next(c)
		}
	}
}
