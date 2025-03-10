package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/sistemica/traefik-manager/internal/logger"
)

// Recovery creates a middleware that recovers from panics
func Recovery() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					logEvent := logger.With().
						Str("method", r.Method).
						Str("path", r.URL.Path).
						Str("remote_addr", r.RemoteAddr)

					logger := logEvent.Logger()

					// Get the stack trace
					stackTrace := debug.Stack()

					// Prepare error message
					errorMsg := fmt.Sprintf("panic: %v", rvr)

					// Get request ID if it exists
					requestID := middleware.GetReqID(r.Context())

					if requestID != "" {
						logEvent = logEvent.Str("request_id", requestID)
					}

					logger.Error().
						Str("error", errorMsg).
						Bytes("stack_trace", stackTrace).
						Msg("Panic recovered")

					// Return a 500 error to the client
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)

					// Write error response
					errResponse := `{"error":{"code":"internal_error","message":"An internal server error occurred"}}`
					w.Write([]byte(errResponse))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
