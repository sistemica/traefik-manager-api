package traefik

// DynamicConfig is the top-level structure for Traefik v3.3 dynamic configuration
type DynamicConfig struct {
	HTTP         *HTTPConfiguration `json:"http,omitempty" yaml:"http,omitempty" toml:"http,omitempty"`
	TCP          *TCPConfiguration  `json:"tcp,omitempty" yaml:"tcp,omitempty" toml:"tcp,omitempty"`
	UDP          *UDPConfiguration  `json:"udp,omitempty" yaml:"udp,omitempty" toml:"udp,omitempty"`
	TLS          *TLSConfiguration  `json:"tls,omitempty" yaml:"tls,omitempty" toml:"tls,omitempty"`
	SSLHost      string             `json:"sslHost,omitempty" yaml:"sslHost,omitempty" toml:"sslHost,omitempty"`
	SSLForceHost bool               `json:"sslForceHost,omitempty" yaml:"sslForceHost,omitempty" toml:"sslForceHost,omitempty"`
}

// IPStrategy defines how to extract remote IP
type IPStrategy struct {
	Depth       int      `json:"depth,omitempty" yaml:"depth,omitempty" toml:"depth,omitempty"`
	ExcludedIPs []string `json:"excludedIPs,omitempty" yaml:"excludedIPs,omitempty" toml:"excludedIPs,omitempty"`
	IPv6Subnet  int      `json:"ipv6Subnet,omitempty" yaml:"ipv6Subnet,omitempty" toml:"ipv6Subnet,omitempty"`
}

// SourceCriterion defines how to select a request source
type SourceCriterion struct {
	IPStrategy        *IPStrategy `json:"ipStrategy,omitempty" yaml:"ipStrategy,omitempty" toml:"ipStrategy,omitempty"`
	RequestHeaderName string      `json:"requestHeaderName,omitempty" yaml:"requestHeaderName,omitempty" toml:"requestHeaderName,omitempty"`
	RequestHost       bool        `json:"requestHost,omitempty" yaml:"requestHost,omitempty" toml:"requestHost,omitempty"`
}

// IPAllowListConfig defines the configuration for the IPAllowList middleware
type IPAllowListConfig struct {
	SourceRange      []string    `json:"sourceRange,omitempty" yaml:"sourceRange,omitempty" toml:"sourceRange,omitempty"`
	IPStrategy       *IPStrategy `json:"ipStrategy,omitempty" yaml:"ipStrategy,omitempty" toml:"ipStrategy,omitempty"`
	RejectStatusCode int         `json:"rejectStatusCode,omitempty" yaml:"rejectStatusCode,omitempty" toml:"rejectStatusCode,omitempty"`
}

// IPWhiteListConfig defines the configuration for the IPWhiteList middleware
type IPWhiteListConfig struct {
	SourceRange []string    `json:"sourceRange,omitempty" yaml:"sourceRange,omitempty" toml:"sourceRange,omitempty"`
	IPStrategy  *IPStrategy `json:"ipStrategy,omitempty" yaml:"ipStrategy,omitempty" toml:"ipStrategy,omitempty"`
}

// InFlightReqConfig defines the configuration for the InFlightReq middleware
type InFlightReqConfig struct {
	Amount          int              `json:"amount,omitempty" yaml:"amount,omitempty" toml:"amount,omitempty"`
	SourceCriterion *SourceCriterion `json:"sourceCriterion,omitempty" yaml:"sourceCriterion,omitempty" toml:"sourceCriterion,omitempty"`
}

// PassTLSClientCertConfig defines the configuration for the PassTLSClientCert middleware
type PassTLSClientCertConfig struct {
	PEM  bool            `json:"pem,omitempty" yaml:"pem,omitempty" toml:"pem,omitempty"`
	Info *ClientCertInfo `json:"info,omitempty" yaml:"info,omitempty" toml:"info,omitempty"`
}

// ClientCertInfo defines the certificate info configuration
type ClientCertInfo struct {
	NotAfter     bool               `json:"notAfter,omitempty" yaml:"notAfter,omitempty" toml:"notAfter,omitempty"`
	NotBefore    bool               `json:"notBefore,omitempty" yaml:"notBefore,omitempty" toml:"notBefore,omitempty"`
	Sans         bool               `json:"sans,omitempty" yaml:"sans,omitempty" toml:"sans,omitempty"`
	SerialNumber bool               `json:"serialNumber,omitempty" yaml:"serialNumber,omitempty" toml:"serialNumber,omitempty"`
	Subject      *ClientCertSubject `json:"subject,omitempty" yaml:"subject,omitempty" toml:"subject,omitempty"`
	Issuer       *ClientCertIssuer  `json:"issuer,omitempty" yaml:"issuer,omitempty" toml:"issuer,omitempty"`
}

// ClientCertSubject defines the subject configuration
type ClientCertSubject struct {
	Country            bool `json:"country,omitempty" yaml:"country,omitempty" toml:"country,omitempty"`
	Province           bool `json:"province,omitempty" yaml:"province,omitempty" toml:"province,omitempty"`
	Locality           bool `json:"locality,omitempty" yaml:"locality,omitempty" toml:"locality,omitempty"`
	Organization       bool `json:"organization,omitempty" yaml:"organization,omitempty" toml:"organization,omitempty"`
	OrganizationalUnit bool `json:"organizationalUnit,omitempty" yaml:"organizationalUnit,omitempty" toml:"organizationalUnit,omitempty"`
	CommonName         bool `json:"commonName,omitempty" yaml:"commonName,omitempty" toml:"commonName,omitempty"`
	SerialNumber       bool `json:"serialNumber,omitempty" yaml:"serialNumber,omitempty" toml:"serialNumber,omitempty"`
	DomainComponent    bool `json:"domainComponent,omitempty" yaml:"domainComponent,omitempty" toml:"domainComponent,omitempty"`
}

// ClientCertIssuer defines the issuer configuration
type ClientCertIssuer struct {
	Country         bool `json:"country,omitempty" yaml:"country,omitempty" toml:"country,omitempty"`
	Province        bool `json:"province,omitempty" yaml:"province,omitempty" toml:"province,omitempty"`
	Locality        bool `json:"locality,omitempty" yaml:"locality,omitempty" toml:"locality,omitempty"`
	Organization    bool `json:"organization,omitempty" yaml:"organization,omitempty" toml:"organization,omitempty"`
	CommonName      bool `json:"commonName,omitempty" yaml:"commonName,omitempty" toml:"commonName,omitempty"`
	SerialNumber    bool `json:"serialNumber,omitempty" yaml:"serialNumber,omitempty" toml:"serialNumber,omitempty"`
	DomainComponent bool `json:"domainComponent,omitempty" yaml:"domainComponent,omitempty" toml:"domainComponent,omitempty"`
}

// RateLimitConfig defines the configuration for the RateLimit middleware
type RateLimitConfig struct {
	Average         int              `json:"average,omitempty" yaml:"average,omitempty" toml:"average,omitempty"`
	Period          string           `json:"period,omitempty" yaml:"period,omitempty" toml:"period,omitempty"`
	Burst           int              `json:"burst,omitempty" yaml:"burst,omitempty" toml:"burst,omitempty"`
	SourceCriterion *SourceCriterion `json:"sourceCriterion,omitempty" yaml:"sourceCriterion,omitempty" toml:"sourceCriterion,omitempty"`
}

// RedirectRegexConfig defines the configuration for the RedirectRegex middleware
type RedirectRegexConfig struct {
	Regex       string `json:"regex,omitempty" yaml:"regex,omitempty" toml:"regex,omitempty"`
	Replacement string `json:"replacement,omitempty" yaml:"replacement,omitempty" toml:"replacement,omitempty"`
	Permanent   bool   `json:"permanent,omitempty" yaml:"permanent,omitempty" toml:"permanent,omitempty"`
}

// RedirectSchemeConfig defines the configuration for the RedirectScheme middleware
type RedirectSchemeConfig struct {
	Scheme    string `json:"scheme,omitempty" yaml:"scheme,omitempty" toml:"scheme,omitempty"`
	Port      string `json:"port,omitempty" yaml:"port,omitempty" toml:"port,omitempty"`
	Permanent bool   `json:"permanent,omitempty" yaml:"permanent,omitempty" toml:"permanent,omitempty"`
}

// ReplacePathConfig defines the configuration for the ReplacePath middleware
type ReplacePathConfig struct {
	Path string `json:"path,omitempty" yaml:"path,omitempty" toml:"path,omitempty"`
}

// ReplacePathRegexConfig defines the configuration for the ReplacePathRegex middleware
type ReplacePathRegexConfig struct {
	Regex       string `json:"regex,omitempty" yaml:"regex,omitempty" toml:"regex,omitempty"`
	Replacement string `json:"replacement,omitempty" yaml:"replacement,omitempty" toml:"replacement,omitempty"`
}

// RetryConfig defines the configuration for the Retry middleware
type RetryConfig struct {
	Attempts        int    `json:"attempts,omitempty" yaml:"attempts,omitempty" toml:"attempts,omitempty"`
	InitialInterval string `json:"initialInterval,omitempty" yaml:"initialInterval,omitempty" toml:"initialInterval,omitempty"`
}

// StripPrefixConfig defines the configuration for the StripPrefix middleware
type StripPrefixConfig struct {
	Prefixes   []string `json:"prefixes,omitempty" yaml:"prefixes,omitempty" toml:"prefixes,omitempty"`
	ForceSlash bool     `json:"forceSlash,omitempty" yaml:"forceSlash,omitempty" toml:"forceSlash,omitempty"`
}

// StripPrefixRegexConfig defines the configuration for the StripPrefixRegex middleware
type StripPrefixRegexConfig struct {
	Regex []string `json:"regex,omitempty" yaml:"regex,omitempty" toml:"regex,omitempty"`
}

// TCP Structs

// TCPRouter defines a TCP router
type TCPRouter struct {
	EntryPoints []string      `json:"entryPoints,omitempty" yaml:"entryPoints,omitempty" toml:"entryPoints,omitempty"`
	Middlewares []string      `json:"middlewares,omitempty" yaml:"middlewares,omitempty" toml:"middlewares,omitempty"`
	Service     string        `json:"service,omitempty" yaml:"service,omitempty" toml:"service,omitempty"`
	Rule        string        `json:"rule,omitempty" yaml:"rule,omitempty" toml:"rule,omitempty"`
	RuleSyntax  string        `json:"ruleSyntax,omitempty" yaml:"ruleSyntax,omitempty" toml:"ruleSyntax,omitempty"`
	Priority    int           `json:"priority,omitempty" yaml:"priority,omitempty" toml:"priority,omitempty"`
	TLS         *TCPRouterTLS `json:"tls,omitempty" yaml:"tls,omitempty" toml:"tls,omitempty"`
}

// TCPRouterTLS defines the TLS configuration for a TCP router
type TCPRouterTLS struct {
	Passthrough  bool      `json:"passthrough,omitempty" yaml:"passthrough,omitempty" toml:"passthrough,omitempty"`
	Options      string    `json:"options,omitempty" yaml:"options,omitempty" toml:"options,omitempty"`
	CertResolver string    `json:"certResolver,omitempty" yaml:"certResolver,omitempty" toml:"certResolver,omitempty"`
	Domains      []*Domain `json:"domains,omitempty" yaml:"domains,omitempty" toml:"domains,omitempty"`
}

// TCPService defines a TCP service
type TCPService struct {
	LoadBalancer *TCPLoadBalancerService `json:"loadBalancer,omitempty" yaml:"loadBalancer,omitempty" toml:"loadBalancer,omitempty"`
	Weighted     *TCPWeightedService     `json:"weighted,omitempty" yaml:"weighted,omitempty" toml:"weighted,omitempty"`
}

// TCPLoadBalancerService defines a TCP load balancer service
type TCPLoadBalancerService struct {
	ProxyProtocol    *ProxyProtocol `json:"proxyProtocol,omitempty" yaml:"proxyProtocol,omitempty" toml:"proxyProtocol,omitempty"`
	Servers          []TCPServer    `json:"servers,omitempty" yaml:"servers,omitempty" toml:"servers,omitempty"`
	ServersTransport string         `json:"serversTransport,omitempty" yaml:"serversTransport,omitempty" toml:"serversTransport,omitempty"`
	TerminationDelay int            `json:"terminationDelay,omitempty" yaml:"terminationDelay,omitempty" toml:"terminationDelay,omitempty"`
}

// ProxyProtocol defines the proxy protocol configuration
type ProxyProtocol struct {
	Version int `json:"version,omitempty" yaml:"version,omitempty" toml:"version,omitempty"`
}

// TCPServer defines a server in a TCP load balancer
type TCPServer struct {
	Address string `json:"address,omitempty" yaml:"address,omitempty" toml:"address,omitempty"`
	TLS     bool   `json:"tls,omitempty" yaml:"tls,omitempty" toml:"tls,omitempty"`
}

// TCPWeightedService defines a TCP weighted service
type TCPWeightedService struct {
	Services []WeightedServiceItem `json:"services,omitempty" yaml:"services,omitempty" toml:"services,omitempty"`
}

// TCPMiddleware defines a TCP middleware
type TCPMiddleware struct {
	IPAllowList  *TCPIPAllowListConfig  `json:"ipAllowList,omitempty" yaml:"ipAllowList,omitempty" toml:"ipAllowList,omitempty"`
	IPWhiteList  *TCPIPWhiteListConfig  `json:"ipWhiteList,omitempty" yaml:"ipWhiteList,omitempty" toml:"ipWhiteList,omitempty"`
	InFlightConn *TCPInFlightConnConfig `json:"inFlightConn,omitempty" yaml:"inFlightConn,omitempty" toml:"inFlightConn,omitempty"`
}

// TCPIPAllowListConfig defines the configuration for the TCP IPAllowList middleware
type TCPIPAllowListConfig struct {
	SourceRange []string `json:"sourceRange,omitempty" yaml:"sourceRange,omitempty" toml:"sourceRange,omitempty"`
}

// TCPIPWhiteListConfig defines the configuration for the TCP IPWhiteList middleware
type TCPIPWhiteListConfig struct {
	SourceRange []string `json:"sourceRange,omitempty" yaml:"sourceRange,omitempty" toml:"sourceRange,omitempty"`
}

// TCPInFlightConnConfig defines the configuration for the TCP InFlightConn middleware
type TCPInFlightConnConfig struct {
	Amount int `json:"amount,omitempty" yaml:"amount,omitempty" toml:"amount,omitempty"`
}

// UDP Structs

// UDPRouter defines a UDP router
type UDPRouter struct {
	EntryPoints []string `json:"entryPoints,omitempty" yaml:"entryPoints,omitempty" toml:"entryPoints,omitempty"`
	Service     string   `json:"service,omitempty" yaml:"service,omitempty" toml:"service,omitempty"`
}

// UDPService defines a UDP service
type UDPService struct {
	LoadBalancer *UDPLoadBalancerService `json:"loadBalancer,omitempty" yaml:"loadBalancer,omitempty" toml:"loadBalancer,omitempty"`
	Weighted     *UDPWeightedService     `json:"weighted,omitempty" yaml:"weighted,omitempty" toml:"weighted,omitempty"`
}

// UDPLoadBalancerService defines a UDP load balancer service
type UDPLoadBalancerService struct {
	Servers []UDPServer `json:"servers,omitempty" yaml:"servers,omitempty" toml:"servers,omitempty"`
}

// UDPServer defines a server in a UDP load balancer
type UDPServer struct {
	Address string `json:"address,omitempty" yaml:"address,omitempty" toml:"address,omitempty"`
}

// UDPWeightedService defines a UDP weighted service
type UDPWeightedService struct {
	Services []WeightedServiceItem `json:"services,omitempty" yaml:"services,omitempty" toml:"services,omitempty"`
}

// TLS Structs

// TLSCertificate defines a TLS certificate
type TLSCertificate struct {
	CertFile string   `json:"certFile,omitempty" yaml:"certFile,omitempty" toml:"certFile,omitempty"`
	KeyFile  string   `json:"keyFile,omitempty" yaml:"keyFile,omitempty" toml:"keyFile,omitempty"`
	Stores   []string `json:"stores,omitempty" yaml:"stores,omitempty" toml:"stores,omitempty"`
}

// TLSOption defines TLS options
type TLSOption struct {
	MinVersion               string      `json:"minVersion,omitempty" yaml:"minVersion,omitempty" toml:"minVersion,omitempty"`
	MaxVersion               string      `json:"maxVersion,omitempty" yaml:"maxVersion,omitempty" toml:"maxVersion,omitempty"`
	CipherSuites             []string    `json:"cipherSuites,omitempty" yaml:"cipherSuites,omitempty" toml:"cipherSuites,omitempty"`
	CurvePreferences         []string    `json:"curvePreferences,omitempty" yaml:"curvePreferences,omitempty" toml:"curvePreferences,omitempty"`
	ClientAuth               *ClientAuth `json:"clientAuth,omitempty" yaml:"clientAuth,omitempty" toml:"clientAuth,omitempty"`
	SNIStrict                bool        `json:"sniStrict,omitempty" yaml:"sniStrict,omitempty" toml:"sniStrict,omitempty"`
	ALPNProtocols            []string    `json:"alpnProtocols,omitempty" yaml:"alpnProtocols,omitempty" toml:"alpnProtocols,omitempty"`
	PreferServerCipherSuites bool        `json:"preferServerCipherSuites,omitempty" yaml:"preferServerCipherSuites,omitempty" toml:"preferServerCipherSuites,omitempty"`
}

// ClientAuth defines the client authentication configuration
type ClientAuth struct {
	CAFiles        []string `json:"caFiles,omitempty" yaml:"caFiles,omitempty" toml:"caFiles,omitempty"`
	ClientAuthType string   `json:"clientAuthType,omitempty" yaml:"clientAuthType,omitempty" toml:"clientAuthType,omitempty"`
}

// TLSStore defines a TLS store
type TLSStore struct {
	DefaultCertificate   *DefaultCertificate   `json:"defaultCertificate,omitempty" yaml:"defaultCertificate,omitempty" toml:"defaultCertificate,omitempty"`
	DefaultGeneratedCert *DefaultGeneratedCert `json:"defaultGeneratedCert,omitempty" yaml:"defaultGeneratedCert,omitempty" toml:"defaultGeneratedCert,omitempty"`
}

// DefaultCertificate defines the default certificate for a TLS store
type DefaultCertificate struct {
	CertFile string `json:"certFile,omitempty" yaml:"certFile,omitempty" toml:"certFile,omitempty"`
	KeyFile  string `json:"keyFile,omitempty" yaml:"keyFile,omitempty" toml:"keyFile,omitempty"`
}

// DefaultGeneratedCert defines the default generated certificate for a TLS store
type DefaultGeneratedCert struct {
	Resolver string          `json:"resolver,omitempty" yaml:"resolver,omitempty" toml:"resolver,omitempty"`
	Domain   *DomainWithSans `json:"domain,omitempty" yaml:"domain,omitempty" toml:"domain,omitempty"`
}

// DomainWithSans defines a domain with optional SANs
type DomainWithSans struct {
	Main string   `json:"main,omitempty" yaml:"main,omitempty" toml:"main,omitempty"`
	Sans []string `json:"sans,omitempty" yaml:"sans,omitempty" toml:"sans,omitempty"`
}

// HTTPConfiguration contains all the HTTP configuration elements
type HTTPConfiguration struct {
	Routers     map[string]*Router     `json:"routers,omitempty" yaml:"routers,omitempty" toml:"routers,omitempty"`
	Services    map[string]*Service    `json:"services,omitempty" yaml:"services,omitempty" toml:"services,omitempty"`
	Middlewares map[string]*Middleware `json:"middlewares,omitempty" yaml:"middlewares,omitempty" toml:"middlewares,omitempty"`
}

// TCPConfiguration contains all the TCP configuration elements
type TCPConfiguration struct {
	Routers     map[string]*TCPRouter     `json:"routers,omitempty" yaml:"routers,omitempty" toml:"routers,omitempty"`
	Services    map[string]*TCPService    `json:"services,omitempty" yaml:"services,omitempty" toml:"services,omitempty"`
	Middlewares map[string]*TCPMiddleware `json:"middlewares,omitempty" yaml:"middlewares,omitempty" toml:"middlewares,omitempty"`
}

// UDPConfiguration contains all the UDP configuration elements
type UDPConfiguration struct {
	Routers  map[string]*UDPRouter  `json:"routers,omitempty" yaml:"routers,omitempty" toml:"routers,omitempty"`
	Services map[string]*UDPService `json:"services,omitempty" yaml:"services,omitempty" toml:"services,omitempty"`
}

// TLSConfiguration contains all the TLS configuration elements
type TLSConfiguration struct {
	Certificates []TLSCertificate      `json:"certificates,omitempty" yaml:"certificates,omitempty" toml:"certificates,omitempty"`
	Options      map[string]*TLSOption `json:"options,omitempty" yaml:"options,omitempty" toml:"options,omitempty"`
	Stores       map[string]*TLSStore  `json:"stores,omitempty" yaml:"stores,omitempty" toml:"stores,omitempty"`
}

// Router defines an HTTP router
type Router struct {
	EntryPoints   []string       `json:"entryPoints,omitempty" yaml:"entryPoints,omitempty" toml:"entryPoints,omitempty"`
	Middlewares   []string       `json:"middlewares,omitempty" yaml:"middlewares,omitempty" toml:"middlewares,omitempty"`
	Service       string         `json:"service,omitempty" yaml:"service,omitempty" toml:"service,omitempty"`
	Rule          string         `json:"rule,omitempty" yaml:"rule,omitempty" toml:"rule,omitempty"`
	Priority      int            `json:"priority,omitempty" yaml:"priority,omitempty" toml:"priority,omitempty"`
	TLS           *RouterTLS     `json:"tls,omitempty" yaml:"tls,omitempty" toml:"tls,omitempty"`
	Observability *Observability `json:"observability,omitempty" yaml:"observability,omitempty" toml:"observability,omitempty"`
}

// RouterTLS defines the TLS configuration for a router
type RouterTLS struct {
	Options      string    `json:"options,omitempty" yaml:"options,omitempty" toml:"options,omitempty"`
	CertResolver string    `json:"certResolver,omitempty" yaml:"certResolver,omitempty" toml:"certResolver,omitempty"`
	Domains      []*Domain `json:"domains,omitempty" yaml:"domains,omitempty" toml:"domains,omitempty"`
}

// Domain defines a domain name with optional SANs
type Domain struct {
	Main string   `json:"main,omitempty" yaml:"main,omitempty" toml:"main,omitempty"`
	Sans []string `json:"sans,omitempty" yaml:"sans,omitempty" toml:"sans,omitempty"`
}

// Observability defines observability parameters for a router
type Observability struct {
	AccessLogs bool `json:"accessLogs,omitempty" yaml:"accessLogs,omitempty" toml:"accessLogs,omitempty"`
	Tracing    bool `json:"tracing,omitempty" yaml:"tracing,omitempty" toml:"tracing,omitempty"`
	Metrics    bool `json:"metrics,omitempty" yaml:"metrics,omitempty" toml:"metrics,omitempty"`
}

// Service defines an HTTP service
type Service struct {
	LoadBalancer *LoadBalancerService `json:"loadBalancer,omitempty" yaml:"loadBalancer,omitempty" toml:"loadBalancer,omitempty"`
	Weighted     *WeightedService     `json:"weighted,omitempty" yaml:"weighted,omitempty" toml:"weighted,omitempty"`
	Mirroring    *MirroringService    `json:"mirroring,omitempty" yaml:"mirroring,omitempty" toml:"mirroring,omitempty"`
	Failover     *FailoverService     `json:"failover,omitempty" yaml:"failover,omitempty" toml:"failover,omitempty"`
	URL          string               `json:"url,omitempty" yaml:"url,omitempty" toml:"url,omitempty"`
}

// LoadBalancerService defines a load balancer service
type LoadBalancerService struct {
	Servers            []Server            `json:"servers,omitempty" yaml:"servers,omitempty" toml:"servers,omitempty"`
	HealthCheck        *HealthCheck        `json:"healthCheck,omitempty" yaml:"healthCheck,omitempty" toml:"healthCheck,omitempty"`
	PassHostHeader     *bool               `json:"passHostHeader,omitempty" yaml:"passHostHeader,omitempty" toml:"passHostHeader,omitempty"`
	ResponseForwarding *ResponseForwarding `json:"responseForwarding,omitempty" yaml:"responseForwarding,omitempty" toml:"responseForwarding,omitempty"`
	ServersTransport   string              `json:"serversTransport,omitempty" yaml:"serversTransport,omitempty" toml:"serversTransport,omitempty"`
	Sticky             *Sticky             `json:"sticky,omitempty" yaml:"sticky,omitempty" toml:"sticky,omitempty"`
}

// Server defines a server in a load balancer
type Server struct {
	URL          string `json:"url,omitempty" yaml:"url,omitempty" toml:"url,omitempty"`
	Weight       *int   `json:"weight,omitempty" yaml:"weight,omitempty" toml:"weight,omitempty"`
	PreservePath *bool  `json:"preservePath,omitempty" yaml:"preservePath,omitempty" toml:"preservePath,omitempty"`
}

// WeightedService defines a weighted service
type WeightedService struct {
	Services    []WeightedServiceItem `json:"services,omitempty" yaml:"services,omitempty" toml:"services,omitempty"`
	Sticky      *Sticky               `json:"sticky,omitempty" yaml:"sticky,omitempty" toml:"sticky,omitempty"`
	HealthCheck *HealthCheck          `json:"healthCheck,omitempty" yaml:"healthCheck,omitempty" toml:"healthCheck,omitempty"`
}

// WeightedServiceItem defines an item in a weighted service
type WeightedServiceItem struct {
	Name   string `json:"name,omitempty" yaml:"name,omitempty" toml:"name,omitempty"`
	Weight int    `json:"weight,omitempty" yaml:"weight,omitempty" toml:"weight,omitempty"`
}

// MirroringService defines a mirroring service
type MirroringService struct {
	Service     string              `json:"service,omitempty" yaml:"service,omitempty" toml:"service,omitempty"`
	MirrorBody  *bool               `json:"mirrorBody,omitempty" yaml:"mirrorBody,omitempty" toml:"mirrorBody,omitempty"`
	MaxBodySize *int                `json:"maxBodySize,omitempty" yaml:"maxBodySize,omitempty" toml:"maxBodySize,omitempty"`
	Mirrors     []MirrorServiceItem `json:"mirrors,omitempty" yaml:"mirrors,omitempty" toml:"mirrors,omitempty"`
	HealthCheck *HealthCheck        `json:"healthCheck,omitempty" yaml:"healthCheck,omitempty" toml:"healthCheck,omitempty"`
}

// MirrorServiceItem defines an item in a mirroring service
type MirrorServiceItem struct {
	Name    string `json:"name,omitempty" yaml:"name,omitempty" toml:"name,omitempty"`
	Percent int    `json:"percent,omitempty" yaml:"percent,omitempty" toml:"percent,omitempty"`
}

// FailoverService defines a failover service
type FailoverService struct {
	Service     string       `json:"service,omitempty" yaml:"service,omitempty" toml:"service,omitempty"`
	Fallback    string       `json:"fallback,omitempty" yaml:"fallback,omitempty" toml:"fallback,omitempty"`
	HealthCheck *HealthCheck `json:"healthCheck,omitempty" yaml:"healthCheck,omitempty" toml:"healthCheck,omitempty"`
}

// HealthCheck defines a health check
type HealthCheck struct {
	Scheme          string            `json:"scheme,omitempty" yaml:"scheme,omitempty" toml:"scheme,omitempty"`
	Mode            string            `json:"mode,omitempty" yaml:"mode,omitempty" toml:"mode,omitempty"`
	Path            string            `json:"path,omitempty" yaml:"path,omitempty" toml:"path,omitempty"`
	Method          string            `json:"method,omitempty" yaml:"method,omitempty" toml:"method,omitempty"`
	Status          int               `json:"status,omitempty" yaml:"status,omitempty" toml:"status,omitempty"`
	Port            int               `json:"port,omitempty" yaml:"port,omitempty" toml:"port,omitempty"`
	Interval        string            `json:"interval,omitempty" yaml:"interval,omitempty" toml:"interval,omitempty"`
	Timeout         string            `json:"timeout,omitempty" yaml:"timeout,omitempty" toml:"timeout,omitempty"`
	Hostname        string            `json:"hostname,omitempty" yaml:"hostname,omitempty" toml:"hostname,omitempty"`
	FollowRedirects *bool             `json:"followRedirects,omitempty" yaml:"followRedirects,omitempty" toml:"followRedirects,omitempty"`
	Headers         map[string]string `json:"headers,omitempty" yaml:"headers,omitempty" toml:"headers,omitempty"`
}

// ResponseForwarding defines how to forward responses
type ResponseForwarding struct {
	FlushInterval string `json:"flushInterval,omitempty" yaml:"flushInterval,omitempty" toml:"flushInterval,omitempty"`
}

// Sticky defines sticky sessions configuration
type Sticky struct {
	Cookie *StickyCooke `json:"cookie,omitempty" yaml:"cookie,omitempty" toml:"cookie,omitempty"`
}

// StickyCooke defines cookie configuration for sticky sessions
type StickyCooke struct {
	Name     string `json:"name,omitempty" yaml:"name,omitempty" toml:"name,omitempty"`
	Secure   *bool  `json:"secure,omitempty" yaml:"secure,omitempty" toml:"secure,omitempty"`
	HTTPOnly *bool  `json:"httpOnly,omitempty" yaml:"httpOnly,omitempty" toml:"httpOnly,omitempty"`
	SameSite string `json:"sameSite,omitempty" yaml:"sameSite,omitempty" toml:"sameSite,omitempty"`
	MaxAge   int    `json:"maxAge,omitempty" yaml:"maxAge,omitempty" toml:"maxAge,omitempty"`
	Path     string `json:"path,omitempty" yaml:"path,omitempty" toml:"path,omitempty"`
}

// Middleware defines an HTTP middleware
type Middleware struct {
	AddPrefix         *AddPrefixConfig         `json:"addPrefix,omitempty" yaml:"addPrefix,omitempty" toml:"addPrefix,omitempty"`
	BasicAuth         *BasicAuthConfig         `json:"basicAuth,omitempty" yaml:"basicAuth,omitempty" toml:"basicAuth,omitempty"`
	Buffering         *BufferingConfig         `json:"buffering,omitempty" yaml:"buffering,omitempty" toml:"buffering,omitempty"`
	Chain             *ChainConfig             `json:"chain,omitempty" yaml:"chain,omitempty" toml:"chain,omitempty"`
	CircuitBreaker    *CircuitBreakerConfig    `json:"circuitBreaker,omitempty" yaml:"circuitBreaker,omitempty" toml:"circuitBreaker,omitempty"`
	Compress          *CompressConfig          `json:"compress,omitempty" yaml:"compress,omitempty" toml:"compress,omitempty"`
	ContentType       *ContentTypeConfig       `json:"contentType,omitempty" yaml:"contentType,omitempty" toml:"contentType,omitempty"`
	DigestAuth        *DigestAuthConfig        `json:"digestAuth,omitempty" yaml:"digestAuth,omitempty" toml:"digestAuth,omitempty"`
	Errors            *ErrorsConfig            `json:"errors,omitempty" yaml:"errors,omitempty" toml:"errors,omitempty"`
	ForwardAuth       *ForwardAuthConfig       `json:"forwardAuth,omitempty" yaml:"forwardAuth,omitempty" toml:"forwardAuth,omitempty"`
	GrpcWeb           *GrpcWebConfig           `json:"grpcWeb,omitempty" yaml:"grpcWeb,omitempty" toml:"grpcWeb,omitempty"`
	Headers           *HeadersConfig           `json:"headers,omitempty" yaml:"headers,omitempty" toml:"headers,omitempty"`
	IPAllowList       *IPAllowListConfig       `json:"ipAllowList,omitempty" yaml:"ipAllowList,omitempty" toml:"ipAllowList,omitempty"`
	IPWhiteList       *IPWhiteListConfig       `json:"ipWhiteList,omitempty" yaml:"ipWhiteList,omitempty" toml:"ipWhiteList,omitempty"`
	InFlightReq       *InFlightReqConfig       `json:"inFlightReq,omitempty" yaml:"inFlightReq,omitempty" toml:"inFlightReq,omitempty"`
	PassTLSClientCert *PassTLSClientCertConfig `json:"passTLSClientCert,omitempty" yaml:"passTLSClientCert,omitempty" toml:"passTLSClientCert,omitempty"`
	Plugin            map[string]interface{}   `json:"plugin,omitempty" yaml:"plugin,omitempty" toml:"plugin,omitempty"`
	RateLimit         *RateLimitConfig         `json:"rateLimit,omitempty" yaml:"rateLimit,omitempty" toml:"rateLimit,omitempty"`
	RedirectRegex     *RedirectRegexConfig     `json:"redirectRegex,omitempty" yaml:"redirectRegex,omitempty" toml:"redirectRegex,omitempty"`
	RedirectScheme    *RedirectSchemeConfig    `json:"redirectScheme,omitempty" yaml:"redirectScheme,omitempty" toml:"redirectScheme,omitempty"`
	ReplacePath       *ReplacePathConfig       `json:"replacePath,omitempty" yaml:"replacePath,omitempty" toml:"replacePath,omitempty"`
	ReplacePathRegex  *ReplacePathRegexConfig  `json:"replacePathRegex,omitempty" yaml:"replacePathRegex,omitempty" toml:"replacePathRegex,omitempty"`
	Retry             *RetryConfig             `json:"retry,omitempty" yaml:"retry,omitempty" toml:"retry,omitempty"`
	StripPrefix       *StripPrefixConfig       `json:"stripPrefix,omitempty" yaml:"stripPrefix,omitempty" toml:"stripPrefix,omitempty"`
	StripPrefixRegex  *StripPrefixRegexConfig  `json:"stripPrefixRegex,omitempty" yaml:"stripPrefixRegex,omitempty" toml:"stripPrefixRegex,omitempty"`
}

// AddPrefixConfig defines the configuration for the AddPrefix middleware
type AddPrefixConfig struct {
	Prefix string `json:"prefix,omitempty" yaml:"prefix,omitempty" toml:"prefix,omitempty"`
}

// BasicAuthConfig defines the configuration for the BasicAuth middleware
type BasicAuthConfig struct {
	Users        []string `json:"users,omitempty" yaml:"users,omitempty" toml:"users,omitempty"`
	UsersFile    string   `json:"usersFile,omitempty" yaml:"usersFile,omitempty" toml:"usersFile,omitempty"`
	Realm        string   `json:"realm,omitempty" yaml:"realm,omitempty" toml:"realm,omitempty"`
	RemoveHeader *bool    `json:"removeHeader,omitempty" yaml:"removeHeader,omitempty" toml:"removeHeader,omitempty"`
	HeaderField  string   `json:"headerField,omitempty" yaml:"headerField,omitempty" toml:"headerField,omitempty"`
}

// BufferingConfig defines the configuration for the Buffering middleware
type BufferingConfig struct {
	MaxRequestBodyBytes  int    `json:"maxRequestBodyBytes,omitempty" yaml:"maxRequestBodyBytes,omitempty" toml:"maxRequestBodyBytes,omitempty"`
	MemRequestBodyBytes  int    `json:"memRequestBodyBytes,omitempty" yaml:"memRequestBodyBytes,omitempty" toml:"memRequestBodyBytes,omitempty"`
	MaxResponseBodyBytes int    `json:"maxResponseBodyBytes,omitempty" yaml:"maxResponseBodyBytes,omitempty" toml:"maxResponseBodyBytes,omitempty"`
	MemResponseBodyBytes int    `json:"memResponseBodyBytes,omitempty" yaml:"memResponseBodyBytes,omitempty" toml:"memResponseBodyBytes,omitempty"`
	RetryExpression      string `json:"retryExpression,omitempty" yaml:"retryExpression,omitempty" toml:"retryExpression,omitempty"`
}

// ChainConfig defines the configuration for the Chain middleware
type ChainConfig struct {
	Middlewares []string `json:"middlewares,omitempty" yaml:"middlewares,omitempty" toml:"middlewares,omitempty"`
}

// CircuitBreakerConfig defines the configuration for the CircuitBreaker middleware
type CircuitBreakerConfig struct {
	Expression       string `json:"expression,omitempty" yaml:"expression,omitempty" toml:"expression,omitempty"`
	CheckPeriod      string `json:"checkPeriod,omitempty" yaml:"checkPeriod,omitempty" toml:"checkPeriod,omitempty"`
	FallbackDuration string `json:"fallbackDuration,omitempty" yaml:"fallbackDuration,omitempty" toml:"fallbackDuration,omitempty"`
	RecoveryDuration string `json:"recoveryDuration,omitempty" yaml:"recoveryDuration,omitempty" toml:"recoveryDuration,omitempty"`
	ResponseCode     int    `json:"responseCode,omitempty" yaml:"responseCode,omitempty" toml:"responseCode,omitempty"`
}

// CompressConfig defines the configuration for the Compress middleware
type CompressConfig struct {
	ExcludedContentTypes []string `json:"excludedContentTypes,omitempty" yaml:"excludedContentTypes,omitempty" toml:"excludedContentTypes,omitempty"`
	IncludedContentTypes []string `json:"includedContentTypes,omitempty" yaml:"includedContentTypes,omitempty" toml:"includedContentTypes,omitempty"`
	MinResponseBodyBytes int      `json:"minResponseBodyBytes,omitempty" yaml:"minResponseBodyBytes,omitempty" toml:"minResponseBodyBytes,omitempty"`
	Encodings            []string `json:"encodings,omitempty" yaml:"encodings,omitempty" toml:"encodings,omitempty"`
	DefaultEncoding      string   `json:"defaultEncoding,omitempty" yaml:"defaultEncoding,omitempty" toml:"defaultEncoding,omitempty"`
}

// ContentTypeConfig defines the configuration for the ContentType middleware
type ContentTypeConfig struct {
	AutoDetect bool `json:"autoDetect,omitempty" yaml:"autoDetect,omitempty" toml:"autoDetect,omitempty"`
}

// DigestAuthConfig defines the configuration for the DigestAuth middleware
type DigestAuthConfig struct {
	Users        []string `json:"users,omitempty" yaml:"users,omitempty" toml:"users,omitempty"`
	UsersFile    string   `json:"usersFile,omitempty" yaml:"usersFile,omitempty" toml:"usersFile,omitempty"`
	RemoveHeader bool     `json:"removeHeader,omitempty" yaml:"removeHeader,omitempty" toml:"removeHeader,omitempty"`
	Realm        string   `json:"realm,omitempty" yaml:"realm,omitempty" toml:"realm,omitempty"`
	HeaderField  string   `json:"headerField,omitempty" yaml:"headerField,omitempty" toml:"headerField,omitempty"`
}

// ErrorsConfig defines the configuration for the Errors middleware
type ErrorsConfig struct {
	Status  []string `json:"status,omitempty" yaml:"status,omitempty" toml:"status,omitempty"`
	Service string   `json:"service,omitempty" yaml:"service,omitempty" toml:"service,omitempty"`
	Query   string   `json:"query,omitempty" yaml:"query,omitempty" toml:"query,omitempty"`
}

// ForwardAuthConfig defines the configuration for the ForwardAuth middleware
type ForwardAuthConfig struct {
	Address                  string          `json:"address,omitempty" yaml:"address,omitempty" toml:"address,omitempty"`
	TLS                      *ForwardAuthTLS `json:"tls,omitempty" yaml:"tls,omitempty" toml:"tls,omitempty"`
	TrustForwardHeader       bool            `json:"trustForwardHeader,omitempty" yaml:"trustForwardHeader,omitempty" toml:"trustForwardHeader,omitempty"`
	AuthResponseHeaders      []string        `json:"authResponseHeaders,omitempty" yaml:"authResponseHeaders,omitempty" toml:"authResponseHeaders,omitempty"`
	AuthResponseHeadersRegex string          `json:"authResponseHeadersRegex,omitempty" yaml:"authResponseHeadersRegex,omitempty" toml:"authResponseHeadersRegex,omitempty"`
	AuthRequestHeaders       []string        `json:"authRequestHeaders,omitempty" yaml:"authRequestHeaders,omitempty" toml:"authRequestHeaders,omitempty"`
	AddAuthCookiesToResponse []string        `json:"addAuthCookiesToResponse,omitempty" yaml:"addAuthCookiesToResponse,omitempty" toml:"addAuthCookiesToResponse,omitempty"`
	HeaderField              string          `json:"headerField,omitempty" yaml:"headerField,omitempty" toml:"headerField,omitempty"`
	ForwardBody              bool            `json:"forwardBody,omitempty" yaml:"forwardBody,omitempty" toml:"forwardBody,omitempty"`
	MaxBodySize              int             `json:"maxBodySize,omitempty" yaml:"maxBodySize,omitempty" toml:"maxBodySize,omitempty"`
	PreserveLocationHeader   bool            `json:"preserveLocationHeader,omitempty" yaml:"preserveLocationHeader,omitempty" toml:"preserveLocationHeader,omitempty"`
}

// ForwardAuthTLS defines the TLS configuration for ForwardAuth
type ForwardAuthTLS struct {
	CA                 string `json:"ca,omitempty" yaml:"ca,omitempty" toml:"ca,omitempty"`
	Cert               string `json:"cert,omitempty" yaml:"cert,omitempty" toml:"cert,omitempty"`
	Key                string `json:"key,omitempty" yaml:"key,omitempty" toml:"key,omitempty"`
	InsecureSkipVerify bool   `json:"insecureSkipVerify,omitempty" yaml:"insecureSkipVerify,omitempty" toml:"insecureSkipVerify,omitempty"`
	CAOptional         bool   `json:"caOptional,omitempty" yaml:"caOptional,omitempty" toml:"caOptional,omitempty"`
}

// GrpcWebConfig defines the configuration for the GrpcWeb middleware
type GrpcWebConfig struct {
	AllowOrigins []string `json:"allowOrigins,omitempty" yaml:"allowOrigins,omitempty" toml:"allowOrigins,omitempty"`
}

// HeadersConfig defines the configuration for the Headers middleware
type HeadersConfig struct {
	CustomRequestHeaders              map[string]string `json:"customRequestHeaders,omitempty" yaml:"customRequestHeaders,omitempty" toml:"customRequestHeaders,omitempty"`
	CustomResponseHeaders             map[string]string `json:"customResponseHeaders,omitempty" yaml:"customResponseHeaders,omitempty" toml:"customResponseHeaders,omitempty"`
	AccessControlAllowCredentials     bool              `json:"accessControlAllowCredentials,omitempty" yaml:"accessControlAllowCredentials,omitempty" toml:"accessControlAllowCredentials,omitempty"`
	AccessControlAllowHeaders         []string          `json:"accessControlAllowHeaders,omitempty" yaml:"accessControlAllowHeaders,omitempty" toml:"accessControlAllowHeaders,omitempty"`
	AccessControlAllowMethods         []string          `json:"accessControlAllowMethods,omitempty" yaml:"accessControlAllowMethods,omitempty" toml:"accessControlAllowMethods,omitempty"`
	AccessControlAllowOriginList      []string          `json:"accessControlAllowOriginList,omitempty" yaml:"accessControlAllowOriginList,omitempty" toml:"accessControlAllowOriginList,omitempty"`
	AccessControlAllowOriginListRegex []string          `json:"accessControlAllowOriginListRegex,omitempty" yaml:"accessControlAllowOriginListRegex,omitempty" toml:"accessControlAllowOriginListRegex,omitempty"`
	AccessControlExposeHeaders        []string          `json:"accessControlExposeHeaders,omitempty" yaml:"accessControlExposeHeaders,omitempty" toml:"accessControlExposeHeaders,omitempty"`
	AccessControlMaxAge               int               `json:"accessControlMaxAge,omitempty" yaml:"accessControlMaxAge,omitempty" toml:"accessControlMaxAge,omitempty"`
	AddVaryHeader                     bool              `json:"addVaryHeader,omitempty" yaml:"addVaryHeader,omitempty" toml:"addVaryHeader,omitempty"`
	AllowedHosts                      []string          `json:"allowedHosts,omitempty" yaml:"allowedHosts,omitempty" toml:"allowedHosts,omitempty"`
	HostsProxyHeaders                 []string          `json:"hostsProxyHeaders,omitempty" yaml:"hostsProxyHeaders,omitempty" toml:"hostsProxyHeaders,omitempty"`
	SSLProxyHeaders                   map[string]string `json:"sslProxyHeaders,omitempty" yaml:"sslProxyHeaders,omitempty" toml:"sslProxyHeaders,omitempty"`
	STSSeconds                        int               `json:"stsSeconds,omitempty" yaml:"stsSeconds,omitempty" toml:"stsSeconds,omitempty"`
	STSIncludeSubdomains              bool              `json:"stsIncludeSubdomains,omitempty" yaml:"stsIncludeSubdomains,omitempty" toml:"stsIncludeSubdomains,omitempty"`
	STSPreload                        bool              `json:"stsPreload,omitempty" yaml:"stsPreload,omitempty" toml:"stsPreload,omitempty"`
	ForceSTSHeader                    bool              `json:"forceSTSHeader,omitempty" yaml:"forceSTSHeader,omitempty" toml:"forceSTSHeader,omitempty"`
	FrameDeny                         bool              `json:"frameDeny,omitempty" yaml:"frameDeny,omitempty" toml:"frameDeny,omitempty"`
	CustomFrameOptionsValue           string            `json:"customFrameOptionsValue,omitempty" yaml:"customFrameOptionsValue,omitempty" toml:"customFrameOptionsValue,omitempty"`
	ContentTypeNosniff                bool              `json:"contentTypeNosniff,omitempty" yaml:"contentTypeNosniff,omitempty" toml:"contentTypeNosniff,omitempty"`
	BrowserXSSFilter                  bool              `json:"browserXssFilter,omitempty" yaml:"browserXssFilter,omitempty" toml:"browserXssFilter,omitempty"`
	CustomBrowserXSSValue             string            `json:"customBrowserXSSValue,omitempty" yaml:"customBrowserXSSValue,omitempty" toml:"customBrowserXSSValue,omitempty"`
	ContentSecurityPolicy             string            `json:"contentSecurityPolicy,omitempty" yaml:"contentSecurityPolicy,omitempty" toml:"contentSecurityPolicy,omitempty"`
	ContentSecurityPolicyReportOnly   string            `json:"contentSecurityPolicyReportOnly,omitempty" yaml:"contentSecurityPolicyReportOnly,omitempty" toml:"contentSecurityPolicyReportOnly,omitempty"`
	PublicKey                         string            `json:"publicKey,omitempty" yaml:"publicKey,omitempty" toml:"publicKey,omitempty"`
	ReferrerPolicy                    string            `json:"referrerPolicy,omitempty" yaml:"referrerPolicy,omitempty" toml:"referrerPolicy,omitempty"`
	PermissionsPolicy                 string            `json:"permissionsPolicy,omitempty" yaml:"permissionsPolicy,omitempty" toml:"permissionsPolicy,omitempty"`
	IsDevelopment                     bool              `json:"isDevelopment,omitempty" yaml:"isDevelopment,omitempty" toml:"isDevelopment,omitempty"`
	FeaturePolicy                     string            `json:"featurePolicy,omitempty" yaml:"featurePolicy,omitempty" toml:"featurePolicy,omitempty"`
	SSLRedirect                       bool              `json:"sslRedirect,omitempty" yaml:"sslRedirect,omitempty" toml:"sslRedirect,omitempty"`
	SSLTemporaryRedirect              bool              `json:"sslTemporaryRedirect,omitempty" yaml:"sslTemporaryRedirect,omitempty" toml:"sslTemporaryRedirect,omitempty"`
	SSLHost                           string            `json:"sslHost,omitempty" yaml:"sslHost,omitempty" toml:"sslHost,omitempty"`
	SSLForceHost                      bool              `json:"sslForceHost,omitempty" yaml:"sslForceHost,omitempty" toml:"sslForceHost,omitempty"`
}
