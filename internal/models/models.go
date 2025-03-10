package models

// ResourceResponse represents a standard response for resource operations
type ResourceResponse struct {
	ID      string `json:"id"`
	Created bool   `json:"created"`
	Updated bool   `json:"updated"`
	Deleted bool   `json:"deleted"`
}

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

// Router represents a Traefik HTTP router
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

// RouterTLS represents TLS configuration for a router
type RouterTLS struct {
	Options      string   `json:"options,omitempty"`
	CertResolver string   `json:"certResolver,omitempty"`
	Domains      []Domain `json:"domains,omitempty"`
}

// Domain represents a TLS domain configuration
type Domain struct {
	Main string   `json:"main"`
	Sans []string `json:"sans,omitempty"`
}

// Observability represents observability configuration for a router
type Observability struct {
	AccessLogs bool `json:"accessLogs,omitempty"`
	Tracing    bool `json:"tracing,omitempty"`
	Metrics    bool `json:"metrics,omitempty"`
}

// Service represents a Traefik service, which can be one of several types
type Service struct {
	ID           string               `json:"id"`
	ServiceType  string               `json:"serviceType,omitempty"`
	URL          string               `json:"url,omitempty"`
	LoadBalancer *LoadBalancerService `json:"loadBalancer,omitempty"`
	Weighted     *WeightedService     `json:"weighted,omitempty"`
	Mirroring    *MirroringService    `json:"mirroring,omitempty"`
	Failover     *FailoverService     `json:"failover,omitempty"`
}

// LoadBalancerService represents a load balancer service configuration
type LoadBalancerService struct {
	Servers            []Server            `json:"servers"`
	HealthCheck        *HealthCheck        `json:"healthCheck,omitempty"`
	PassHostHeader     bool                `json:"passHostHeader,omitempty"`
	ResponseForwarding *ResponseForwarding `json:"responseForwarding,omitempty"`
	ServersTransport   string              `json:"serversTransport,omitempty"`
	Sticky             *Sticky             `json:"sticky,omitempty"`
}

// WeightedService represents a weighted service configuration
type WeightedService struct {
	Services    []WeightedServiceItem `json:"services"`
	Sticky      *Sticky               `json:"sticky,omitempty"`
	HealthCheck *HealthCheck          `json:"healthCheck,omitempty"`
}

// MirroringService represents a mirroring service configuration
type MirroringService struct {
	Service     Service             `json:"service"`
	MirrorBody  bool                `json:"mirrorBody,omitempty"`
	MaxBodySize int                 `json:"maxBodySize,omitempty"`
	Mirrors     []MirrorServiceItem `json:"mirrors"`
	HealthCheck *HealthCheck        `json:"healthCheck,omitempty"`
}

// FailoverService represents a failover service configuration
type FailoverService struct {
	Service     Service      `json:"service"`
	Fallback    Service      `json:"fallback"`
	HealthCheck *HealthCheck `json:"healthCheck,omitempty"`
}

// Server represents a backend server configuration
type Server struct {
	URL          string `json:"url"`
	Weight       int    `json:"weight,omitempty"`
	PreservePath bool   `json:"preservePath,omitempty"`
}

// WeightedServiceItem represents a weighted service configuration item
type WeightedServiceItem struct {
	Name   Service `json:"name"`
	Weight int     `json:"weight"`
}

// MirrorServiceItem represents a mirror service configuration item
type MirrorServiceItem struct {
	Name    Service `json:"name"`
	Percent int     `json:"percent"`
}

// HealthCheck represents health check configuration
type HealthCheck struct {
	Scheme          string            `json:"scheme,omitempty"`
	Mode            string            `json:"mode,omitempty"`
	Path            string            `json:"path,omitempty"`
	Method          string            `json:"method,omitempty"`
	Status          int               `json:"status,omitempty"`
	Port            int               `json:"port,omitempty"`
	Interval        Duration          `json:"interval,omitempty"`
	Timeout         Duration          `json:"timeout,omitempty"`
	Hostname        string            `json:"hostname,omitempty"`
	FollowRedirects bool              `json:"followRedirects,omitempty"`
	Headers         map[string]string `json:"headers,omitempty"`
}

// ResponseForwarding represents response forwarding configuration
type ResponseForwarding struct {
	FlushInterval Duration `json:"flushInterval,omitempty"`
}

// Sticky represents sticky session configuration
type Sticky struct {
	Cookie *StickyCookie `json:"cookie,omitempty"`
}

// StickyCookie represents sticky cookie configuration
type StickyCookie struct {
	Name     string `json:"name"`
	Secure   bool   `json:"secure,omitempty"`
	HTTPOnly bool   `json:"httpOnly,omitempty"`
	SameSite string `json:"sameSite,omitempty"`
	MaxAge   int    `json:"maxAge,omitempty"`
	Path     string `json:"path,omitempty"`
}

// Duration represents a time duration string that can be unmarshaled from JSON
type Duration string

// Middleware represents a Traefik middleware configuration
type Middleware struct {
	ID     string           `json:"id"`
	Type   string           `json:"type"`
	Config MiddlewareConfig `json:"config"`
}

// DynamicConfig represents a dynamic configuration for Traefik
type DynamicConfig struct {
	HTTPRouters     map[string]Router     `json:"httpRouters,omitempty"`
	HTTPServices    map[string]Service    `json:"httpServices,omitempty"`
	HTTPMiddlewares map[string]Middleware `json:"httpMiddlewares,omitempty"`
	TLSCertificates []TLSCertificate      `json:"tlsCertificates,omitempty"`
	TLSOptions      map[string]TLSOption  `json:"tlsOptions,omitempty"`
	TLSStores       map[string]TLSStore   `json:"tlsStores,omitempty"`
}

// ProxyProtocol represents proxy protocol configuration
type ProxyProtocol struct {
	Version int `json:"version,omitempty"`
}

// TLSCertificate represents a TLS certificate configuration
type TLSCertificate struct {
	CertFile string   `json:"certFile"`
	KeyFile  string   `json:"keyFile"`
	Stores   []string `json:"stores,omitempty"`
}

// TLSOption represents TLS options configuration
type TLSOption struct {
	MinVersion               string      `json:"minVersion,omitempty"`
	MaxVersion               string      `json:"maxVersion,omitempty"`
	CipherSuites             []string    `json:"cipherSuites,omitempty"`
	CurvePreferences         []string    `json:"curvePreferences,omitempty"`
	ClientAuth               *ClientAuth `json:"clientAuth,omitempty"`
	SNIStrict                bool        `json:"sniStrict,omitempty"`
	ALPNProtocols            []string    `json:"alpnProtocols,omitempty"`
	PreferServerCipherSuites bool        `json:"preferServerCipherSuites,omitempty"`
}

// ClientAuth represents client authentication configuration
type ClientAuth struct {
	CAFiles        []string `json:"caFiles,omitempty"`
	ClientAuthType string   `json:"clientAuthType,omitempty"`
}

// TLSStore represents a TLS store configuration
type TLSStore struct {
	DefaultCertificate   *DefaultCertificate   `json:"defaultCertificate,omitempty"`
	DefaultGeneratedCert *DefaultGeneratedCert `json:"defaultGeneratedCert,omitempty"`
}

// DefaultCertificate represents a default certificate configuration
type DefaultCertificate struct {
	CertFile string `json:"certFile"`
	KeyFile  string `json:"keyFile"`
}

// DefaultGeneratedCert represents a default generated certificate configuration
type DefaultGeneratedCert struct {
	Resolver string         `json:"resolver"`
	Domain   DomainWithSans `json:"domain"`
}

// DomainWithSans represents a domain with SANS configuration
type DomainWithSans struct {
	Main string   `json:"main"`
	Sans []string `json:"sans,omitempty"`
}
