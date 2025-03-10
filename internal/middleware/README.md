# Middleware Package

This package provides middleware components for the Traefik Manager HTTP API. These middleware components enhance the API with features like authentication, logging, and error recovery.

## Available Middleware

### Authentication (auth.go)
This package provides API key authentication middleware for the Traefik Manager application.

#### Overview

The authentication middleware validates requests by checking for a valid API key in the request headers. It supports:

- API key-based authentication
- Path-based exclusions for public endpoints
- Easy configuration through options

#### Usage

```go
import (
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/sistemica/traefik-manager/internal/api/middleware"
)

func main() {
    // Create router
    r := chi.NewRouter()

    // Set up authentication
    r.Use(middleware.Auth(middleware.AuthOptions{
        Enabled:      true,
        HeaderName:   "X-API-Key",
        Key:          "your-secret-api-key",
        ExcludePaths: []string{"/health", "/traefik/provider"},
    }))

    // Define routes
    r.Get("/health", healthHandler)
    r.Get("/api/resources", resourceHandler)
    
    // Start server
    http.ListenAndServe(":9000", r)
}
```

#### Options

The authentication middleware supports the following options:

- `Enabled` - Enable or disable authentication
- `HeaderName` - Name of the header containing the API key
- `Key` - Expected API key value
- `ExcludePaths` - Paths that should be accessible without authentication

#### Example

```go
// Setup authentication middleware
authMiddleware := middleware.Auth(middleware.AuthOptions{
    Enabled:    cfg.Auth.Enabled,
    HeaderName: cfg.Auth.HeaderName,
    Key:        cfg.Auth.Key,
    ExcludePaths: []string{
        "/health",
        "/traefik/provider",
        "/api/v1/health",
    },
})

// Apply middleware to router
r.Use(authMiddleware)
```

#### Best Practices

1. **Generate Strong API Keys**: Use a secure random generator to create API keys
2. **Secure Key Transmission**: Always transmit the API key over HTTPS
3. **Key Rotation**: Implement a mechanism to periodically rotate API keys
4. **Minimal Exclusions**: Only exclude paths that truly need to be public
5. **Environment Variables**: Store API keys in environment variables, not in code

### Logging (logging.go)

Request logging middleware that logs HTTP requests with context.

**Features:**
- Structured logging of request details
- Response status and timing information
- Support for request ID tracing

**Usage:**
```go
// Apply logging middleware
r.Use(middleware.Logging())
```

### Recovery (recovery.go)

Panic recovery middleware that prevents crashes from unhandled panics.

**Features:**
- Catches and logs panics
- Returns appropriate error responses
- Prevents application crashes

**Usage:**
```go
// Apply recovery middleware
r.Use(middleware.Recovery())
```

## Middleware Order

The recommended order for applying middleware is:

1. Recovery (first, to catch panics in other middleware)
2. Logging (to log all requests, including those that fail auth)
3. Authentication (after logging but before business logic)

Example:
```go
r := chi.NewRouter()
r.Use(middleware.Recovery())
r.Use(middleware.Logging())
r.Use(middleware.Auth(authOptions))
```

## Best Practices

### Authentication

- Use strong, randomly generated API keys
- Store keys securely (environment variables, secrets manager)
- Transmit keys only over HTTPS
- Implement key rotation mechanism
- Limit public endpoints to only what's necessary

### Logging

- Be mindful of sensitive data in logs
- Use request IDs for request tracing
- Don't log large request bodies
- Configure appropriate log levels based on environment

### Recovery

- Use recovery middleware in all environments
- Provide meaningful error messages to clients
- Log full error context for debugging