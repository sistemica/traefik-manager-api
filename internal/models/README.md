# Models Package

This package contains the internal data models for the Traefik Manager application. These models define the structure of all configuration components used by the application.

## Core Models

### ResourceResponse

`ResourceResponse` represents a standard response for resource operations (create, update, delete):

```go
type ResourceResponse struct {
    ID      string `json:"id"`
    Created bool   `json:"created"`
    Updated bool   `json:"updated"`
    Deleted bool   `json:"deleted"`
}
```

### ErrorResponse

`ErrorResponse` represents a standard error response:

```go
type ErrorResponse struct {
    Error   string `json:"error"`
    Details string `json:"details,omitempty"`
}
```

## Traefik Configuration Components

### Router

`Router` represents a Traefik HTTP router, defining how incoming requests are routed to services:

```go
type Router struct {
    ID            string         `json:"id"`
    EntryPoints   []string       `json:"entryPoints,omitempty"`
    Middlewares   []Middleware   `json:"middlewares,omitempty"`
    Service       Service        `json:"service"`
    Rule          string         `json:"rule"`
    RuleSyntax    string         `json:"ruleSyntax,omitempty"`
    Priority      int            `json:"priority,omitempty"`
    TLS           *RouterTLS     `json:"tls,omitempty"`
    Observability *Observability `json:"observability,omitempty"`
}
```

Key relationships:
- A router references one service via `Service`
- A router can reference multiple middlewares via `Middlewares`
- TLS and observability configurations are optional

### Service

`Service` represents a backend service that Traefik routes traffic to:

```go
type Service struct {
    ID           string               `json:"id"`
    ServiceType  string               `json:"serviceType,omitempty"`
    URL          string               `json:"url,omitempty"`
    LoadBalancer *LoadBalancerService `json:"loadBalancer,omitempty"`
    Weighted     *WeightedService     `json:"weighted,omitempty"`
    Mirroring    *MirroringService    `json:"mirroring,omitempty"`
    Failover     *FailoverService     `json:"failover,omitempty"`
}
```

Service types:
- Simple service: defined by a URL
- LoadBalancer: distributes traffic across multiple backend servers
- Weighted: distributes traffic based on server weights
- Mirroring: sends traffic to a main service and mirrors a percentage to others
- Failover: routes to a fallback service when the main one is down

### Middleware

`Middleware` represents a Traefik middleware that processes requests before they reach services:

```go
type Middleware struct {
    ID     string           `json:"id"`
    Type   string           `json:"type"`
    Config MiddlewareConfig `json:"config"`
}
```

The `MiddlewareConfig` interface is implemented by various concrete types in `middleware_configs.go`.

## Dynamic Configuration

`DynamicConfig` represents the complete dynamic configuration for Traefik:

```go
type DynamicConfig struct {
    HTTPRouters     map[string]Router     `json:"httpRouters,omitempty"`
    HTTPServices    map[string]Service    `json:"httpServices,omitempty"`
    HTTPMiddlewares map[string]Middleware `json:"httpMiddlewares,omitempty"`
    TLSCertificates []TLSCertificate      `json:"tlsCertificates,omitempty"`
    TLSOptions      map[string]TLSOption  `json:"tlsOptions,omitempty"`
    TLSStores       map[string]TLSStore   `json:"tlsStores,omitempty"`
}
```

## Validation Rules

The following validation rules apply to these models:

1. Router IDs, Service IDs, and Middleware IDs must be unique
2. Router rules must be valid according to Traefik's rule syntax
3. Referenced services and middlewares must exist
4. Each service must have a valid configuration (URL or one of the service types)
5. Middleware configurations must be valid for their respective types

## Usage

Models are used throughout the application:
- In the API layer for request/response data
- In the storage layer to persist configurations
- In the Traefik provider to generate dynamic configurations

When creating new models or extending existing ones, ensure they remain compatible with Traefik's expected formats while providing a more user-friendly API for clients.