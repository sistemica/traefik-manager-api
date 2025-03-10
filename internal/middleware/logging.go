package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sistemica/traefik-manager/internal/logger"
)

// Logger creates a middleware that logs HTTP requests
func Logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Start timer
			start := time.Now()

			// Get request information
			req := c.Request()

			// Prepare request logger with context
			reqLoggerCtx := logger.With().
				Str("method", req.Method).
				Str("path", req.URL.Path).
				Str("remote_addr", req.RemoteAddr).
				Str("user_agent", req.UserAgent())

			// Add request ID if it exists
			requestID := req.Header.Get(echo.HeaderXRequestID)
			if requestID == "" {
				requestID = c.Response().Header().Get(echo.HeaderXRequestID)
			}

			if requestID != "" {
				reqLoggerCtx = reqLoggerCtx.Str("request_id", requestID)
			}

			reqLogger := reqLoggerCtx.Logger()

			// Log the request
			reqLogger.Debug().Msg("Request started")

			// Process request
			err := next(c)

			// Request completed
			res := c.Response()

			// Calculate request duration
			duration := time.Since(start)

			// Log the response
			responseLoggerCtx := logger.With().
				Str("method", req.Method).
				Str("path", req.URL.Path).
				Int("status", res.Status).
				Dur("duration", duration).
				Int64("bytes_written", res.Size)

			if requestID != "" {
				responseLoggerCtx = responseLoggerCtx.Str("request_id", requestID)
			}

			responseLogger := responseLoggerCtx.Logger()

			// Log at appropriate level based on status code
			logEvent := responseLogger.Info()
			if err != nil {
				logEvent = responseLogger.Error().Err(err)
			} else if res.Status >= 500 {
				logEvent = responseLogger.Error()
			} else if res.Status >= 400 {
				logEvent = responseLogger.Warn()
			}

			logEvent.Msg("Request completed")

			return err
		}
	}
}
