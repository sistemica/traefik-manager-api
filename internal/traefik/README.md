# Traefik Package

This package contains the Traefik data models and mapping utilities for the Traefik Manager application. It serves as a bridge between the application's internal models and Traefik's expected configuration format.

## Purpose

The Traefik package has two main purposes:

1. **Model Definition**: Define data structures that exactly match Traefik's expected dynamic configuration format
2. **Model Conversion**: Provide mapping functionality to convert between internal models and Traefik-compatible models

This separation allows the application to:
- Present a more user-friendly API to clients
- Add convenience features not directly supported by Traefik
- Validate configurations before they are passed to Traefik
- Handle version compatibility with different Traefik releases

## Key Components

### DynamicConfig

The top-level structure for Traefik's dynamic configuration:

```go
type DynamicConfig struct {
    HTTP         *HTTPConfiguration `json:"http,omitempty"`
    TCP          *TCPConfiguration  `json:"tcp,omitempty"`
    UDP          *UDPConfiguration  `json:"udp,omitempty"`
    TLS          *TLSConfiguration  `json:"tls,omitempty"`
    SSLHost      string             `json:"sslHost,omitempty"`
    SSLForceHost bool               `json:"sslForceHost,omitempty"`
}
```

### HTTPConfiguration

Contains all HTTP configuration elements:

```go
type HTTPConfiguration struct {
    Routers     map[string]*Router     `json:"routers,omitempty"`
    Services    map[string]*Service    `json:"services,omitempty"`
    Middlewares map[string]*Middleware `json:"middlewares,omitempty"`
}
```

### Router, Service, Middleware

Core components that match Traefik's expected format:

```go
type Router struct {
    EntryPoints   []string       `json:"entryPoints,omitempty"`
    Middlewares   []string       `json:"middlewares,omitempty"`
    Service       string         `json:"service,omitempty"`
    Rule          string         `json:"rule,omitempty"`
    Priority      int            `json:"priority,omitempty"`
    TLS           *RouterTLS     `json:"tls,omitempty"`
    Observability *Observability `json:"observability,omitempty"`
}
```

Note that unlike our internal models, Traefik models use strings to reference other components (e.g., middleware IDs) rather than embedding the complete objects.

## Model Mapping

The package provides functionality to convert between internal models and Traefik models:

### Internal → Traefik

When converting from internal models to Traefik models:

1. Convert internal Router objects to Traefik Router objects
   - Extract middleware IDs as strings instead of embedding middleware objects
   - Map service reference to string ID
   - Copy rule, entryPoints, etc.

2. Convert internal Service objects to Traefik Service objects
   - Handle different service types (LoadBalancer, Weighted, etc.)
   - Convert embedded objects to Traefik format
   - Ensure proper pointer handling for optional fields

3. Convert internal Middleware objects to Traefik Middleware objects
   - Map each middleware type to the appropriate Traefik configuration
   - Handle special cases for complex middleware types

## Version Compatibility

This package is designed to work with Traefik v3.3. Key compatibility considerations:

- The structure matches Traefik v3.3's expected format
- Field names and types are aligned with Traefik v3.3's requirements
- Optional fields are properly implemented with pointers
- Boolean values use pointers when appropriate to distinguish between false and unset

## Usage

This package is primarily used in the `ProviderHandler` to convert internal models to Traefik format:

```go
// Convert to Traefik configuration
config := convertToTraefikConfig(routers, services, middlewares)

// Return configuration in Traefik's expected format
return c.JSON(http.StatusOK, config)
```

When extending the application to support new Traefik features:

1. Add the new fields/structures to the corresponding Traefik models
2. Update the conversion functions to handle the new fields
3. Test with the target Traefik version to ensure compatibility