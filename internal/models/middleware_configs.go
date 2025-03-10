package models

// MiddlewareConfig is an interface for the various middleware configurations
type MiddlewareConfig interface{}

// AddPrefixConfig represents the configuration for the AddPrefix middleware
type AddPrefixConfig struct {
	Prefix string `json:"prefix"`
}

// BasicAuthConfig represents the configuration for the BasicAuth middleware
type BasicAuthConfig struct {
	Users        []string `json:"users,omitempty"`
	UsersFile    string   `json:"usersFile,omitempty"`
	Realm        string   `json:"realm,omitempty"`
	RemoveHeader bool     `json:"removeHeader,omitempty"`
	HeaderField  string   `json:"headerField,omitempty"`
}

// BufferingConfig represents the configuration for the Buffering middleware
type BufferingConfig struct {
	MaxRequestBodyBytes  int    `json:"maxRequestBodyBytes,omitempty"`
	MemRequestBodyBytes  int    `json:"memRequestBodyBytes,omitempty"`
	MaxResponseBodyBytes int    `json:"maxResponseBodyBytes,omitempty"`
	MemResponseBodyBytes int    `json:"memResponseBodyBytes,omitempty"`
	RetryExpression      string `json:"retryExpression,omitempty"`
}

// ChainConfig represents the configuration for the Chain middleware
type ChainConfig struct {
	Middlewares []Middleware `json:"middlewares"`
}

// CircuitBreakerConfig represents the configuration for the CircuitBreaker middleware
type CircuitBreakerConfig struct {
	Expression       string   `json:"expression"`
	CheckPeriod      Duration `json:"checkPeriod,omitempty"`
	FallbackDuration Duration `json:"fallbackDuration,omitempty"`
	RecoveryDuration Duration `json:"recoveryDuration,omitempty"`
	ResponseCode     int      `json:"responseCode,omitempty"`
}

// CompressConfig represents the configuration for the Compress middleware
type CompressConfig struct {
	ExcludedContentTypes []string `json:"excludedContentTypes,omitempty"`
	IncludedContentTypes []string `json:"includedContentTypes,omitempty"`
	MinResponseBodyBytes int      `json:"minResponseBodyBytes,omitempty"`
	Encodings            []string `json:"encodings,omitempty"`
	DefaultEncoding      string   `json:"defaultEncoding,omitempty"`
}

// ContentTypeConfig represents the configuration for the ContentType middleware
type ContentTypeConfig struct {
	AutoDetect bool `json:"autoDetect,omitempty"`
}

// DigestAuthConfig represents the configuration for the DigestAuth middleware
type DigestAuthConfig struct {
	Users        []string `json:"users,omitempty"`
	UsersFile    string   `json:"usersFile,omitempty"`
	RemoveHeader bool     `json:"removeHeader,omitempty"`
	Realm        string   `json:"realm,omitempty"`
	HeaderField  string   `json:"headerField,omitempty"`
}

// ErrorsConfig represents the configuration for the Errors middleware
type ErrorsConfig struct {
	Status  []string `json:"status"`
	Service Service  `json:"service"`
	Query   string   `json:"query,omitempty"`
}

// ForwardAuthConfig represents the configuration for the ForwardAuth middleware
type ForwardAuthConfig struct {
	Address                  string          `json:"address"`
	TLS                      *ForwardAuthTLS `json:"tls,omitempty"`
	TrustForwardHeader       bool            `json:"trustForwardHeader,omitempty"`
	AuthResponseHeaders      []string        `json:"authResponseHeaders,omitempty"`
	AuthResponseHeadersRegex string          `json:"authResponseHeadersRegex,omitempty"`
	AuthRequestHeaders       []string        `json:"authRequestHeaders,omitempty"`
	AddAuthCookiesToResponse []string        `json:"addAuthCookiesToResponse,omitempty"`
	HeaderField              string          `json:"headerField,omitempty"`
	ForwardBody              bool            `json:"forwardBody,omitempty"`
	MaxBodySize              int             `json:"maxBodySize,omitempty"`
	PreserveLocationHeader   bool            `json:"preserveLocationHeader,omitempty"`
}

// ForwardAuthTLS represents TLS configuration for ForwardAuth
type ForwardAuthTLS struct {
	CA                 string `json:"ca,omitempty"`
	Cert               string `json:"cert,omitempty"`
	Key                string `json:"key,omitempty"`
	InsecureSkipVerify bool   `json:"insecureSkipVerify,omitempty"`
	CAOptional         bool   `json:"caOptional,omitempty"`
}

// GrpcWebConfig represents the configuration for the GrpcWeb middleware
type GrpcWebConfig struct {
	AllowOrigins []string `json:"allowOrigins,omitempty"`
}

// HeadersConfig represents the configuration for the Headers middleware
type HeadersConfig struct {
	CustomRequestHeaders              map[string]string `json:"customRequestHeaders,omitempty"`
	CustomResponseHeaders             map[string]string `json:"customResponseHeaders,omitempty"`
	AccessControlAllowCredentials     bool              `json:"accessControlAllowCredentials,omitempty"`
	AccessControlAllowHeaders         []string          `json:"accessControlAllowHeaders,omitempty"`
	AccessControlAllowMethods         []string          `json:"accessControlAllowMethods,omitempty"`
	AccessControlAllowOriginList      []string          `json:"accessControlAllowOriginList,omitempty"`
	AccessControlAllowOriginListRegex []string          `json:"accessControlAllowOriginListRegex,omitempty"`
	AccessControlExposeHeaders        []string          `json:"accessControlExposeHeaders,omitempty"`
	AccessControlMaxAge               int               `json:"accessControlMaxAge,omitempty"`
	AddVaryHeader                     bool              `json:"addVaryHeader,omitempty"`
	AllowedHosts                      []string          `json:"allowedHosts,omitempty"`
	HostsProxyHeaders                 []string          `json:"hostsProxyHeaders,omitempty"`
	SSLProxyHeaders                   map[string]string `json:"sslProxyHeaders,omitempty"`
	STSSeconds                        int               `json:"stsSeconds,omitempty"`
	STSIncludeSubdomains              bool              `json:"stsIncludeSubdomains,omitempty"`
	STSPreload                        bool              `json:"stsPreload,omitempty"`
	ForceSTSHeader                    bool              `json:"forceSTSHeader,omitempty"`
	FrameDeny                         bool              `json:"frameDeny,omitempty"`
	CustomFrameOptionsValue           string            `json:"customFrameOptionsValue,omitempty"`
	ContentTypeNosniff                bool              `json:"contentTypeNosniff,omitempty"`
	BrowserXSSFilter                  bool              `json:"browserXssFilter,omitempty"`
	CustomBrowserXSSValue             string            `json:"customBrowserXSSValue,omitempty"`
	ContentSecurityPolicy             string            `json:"contentSecurityPolicy,omitempty"`
	ContentSecurityPolicyReportOnly   string            `json:"contentSecurityPolicyReportOnly,omitempty"`
	PublicKey                         string            `json:"publicKey,omitempty"`
	ReferrerPolicy                    string            `json:"referrerPolicy,omitempty"`
	PermissionsPolicy                 string            `json:"permissionsPolicy,omitempty"`
	IsDevelopment                     bool              `json:"isDevelopment,omitempty"`
	FeaturePolicy                     string            `json:"featurePolicy,omitempty"`
	SSLRedirect                       bool              `json:"sslRedirect,omitempty"`
	SSLTemporaryRedirect              bool              `json:"sslTemporaryRedirect,omitempty"`
	SSLHost                           string            `json:"sslHost,omitempty"`
	SSLForceHost                      bool              `json:"sslForceHost,omitempty"`
}

// IPStrategy represents IP strategy configuration
type IPStrategy struct {
	Depth       int      `json:"depth,omitempty"`
	ExcludedIPs []string `json:"excludedIPs,omitempty"`
	IPv6Subnet  int      `json:"ipv6Subnet,omitempty"`
}

// IPAllowListConfig represents the configuration for the IPAllowList middleware
type IPAllowListConfig struct {
	SourceRange      []string    `json:"sourceRange"`
	IPStrategy       *IPStrategy `json:"ipStrategy,omitempty"`
	RejectStatusCode int         `json:"rejectStatusCode,omitempty"`
}

// IPWhiteListConfig represents the configuration for the IPWhiteList middleware
type IPWhiteListConfig struct {
	SourceRange []string    `json:"sourceRange"`
	IPStrategy  *IPStrategy `json:"ipStrategy,omitempty"`
}

// SourceCriterion represents source criterion configuration
type SourceCriterion struct {
	IPStrategy        *IPStrategy `json:"ipStrategy,omitempty"`
	RequestHeaderName string      `json:"requestHeaderName,omitempty"`
	RequestHost       bool        `json:"requestHost,omitempty"`
}

// InFlightReqConfig represents the configuration for the InFlightReq middleware
type InFlightReqConfig struct {
	Amount          int              `json:"amount"`
	SourceCriterion *SourceCriterion `json:"sourceCriterion,omitempty"`
}

// PassTLSClientCertConfig represents the configuration for the PassTLSClientCert middleware
type PassTLSClientCertConfig struct {
	PEM  bool            `json:"pem,omitempty"`
	Info *ClientCertInfo `json:"info,omitempty"`
}

// ClientCertInfo represents client certificate information configuration
type ClientCertInfo struct {
	NotAfter     bool               `json:"notAfter,omitempty"`
	NotBefore    bool               `json:"notBefore,omitempty"`
	Sans         bool               `json:"sans,omitempty"`
	SerialNumber bool               `json:"serialNumber,omitempty"`
	Subject      *ClientCertSubject `json:"subject,omitempty"`
	Issuer       *ClientCertIssuer  `json:"issuer,omitempty"`
}

// ClientCertSubject represents client certificate subject configuration
type ClientCertSubject struct {
	Country            bool `json:"country,omitempty"`
	Province           bool `json:"province,omitempty"`
	Locality           bool `json:"locality,omitempty"`
	Organization       bool `json:"organization,omitempty"`
	OrganizationalUnit bool `json:"organizationalUnit,omitempty"`
	CommonName         bool `json:"commonName,omitempty"`
	SerialNumber       bool `json:"serialNumber,omitempty"`
	DomainComponent    bool `json:"domainComponent,omitempty"`
}

// ClientCertIssuer represents client certificate issuer configuration
type ClientCertIssuer struct {
	Country         bool `json:"country,omitempty"`
	Province        bool `json:"province,omitempty"`
	Locality        bool `json:"locality,omitempty"`
	Organization    bool `json:"organization,omitempty"`
	CommonName      bool `json:"commonName,omitempty"`
	SerialNumber    bool `json:"serialNumber,omitempty"`
	DomainComponent bool `json:"domainComponent,omitempty"`
}

// PluginConfig represents the configuration for a Plugin middleware
type PluginConfig map[string]interface{}

// RateLimitConfig represents the configuration for the RateLimit middleware
type RateLimitConfig struct {
	Average         int              `json:"average"`
	Period          Duration         `json:"period,omitempty"`
	Burst           int              `json:"burst,omitempty"`
	SourceCriterion *SourceCriterion `json:"sourceCriterion,omitempty"`
}

// RedirectRegexConfig represents the configuration for the RedirectRegex middleware
type RedirectRegexConfig struct {
	Regex       string `json:"regex"`
	Replacement string `json:"replacement"`
	Permanent   bool   `json:"permanent,omitempty"`
}

// RedirectSchemeConfig represents the configuration for the RedirectScheme middleware
type RedirectSchemeConfig struct {
	Scheme    string `json:"scheme"`
	Port      string `json:"port,omitempty"`
	Permanent bool   `json:"permanent,omitempty"`
}

// ReplacePathConfig represents the configuration for the ReplacePath middleware
type ReplacePathConfig struct {
	Path string `json:"path"`
}

// ReplacePathRegexConfig represents the configuration for the ReplacePathRegex middleware
type ReplacePathRegexConfig struct {
	Regex       string `json:"regex"`
	Replacement string `json:"replacement"`
}

// RetryConfig represents the configuration for the Retry middleware
type RetryConfig struct {
	Attempts        int      `json:"attempts"`
	InitialInterval Duration `json:"initialInterval,omitempty"`
}

// StripPrefixConfig represents the configuration for the StripPrefix middleware
type StripPrefixConfig struct {
	Prefixes   []string `json:"prefixes"`
	ForceSlash bool     `json:"forceSlash,omitempty"`
}

// StripPrefixRegexConfig represents the configuration for the StripPrefixRegex middleware
type StripPrefixRegexConfig struct {
	Regex []string `json:"regex"`
}
