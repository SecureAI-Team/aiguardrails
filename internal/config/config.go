package config

// Config holds runtime configuration for the platform.
type Config struct {
	HTTPPort       string
	DatabaseURL    string
	RedisURL       string
	KMSKey         string
	AllowedOrigins []string
	AdminToken     string
	OIDCIssuer     string
	OIDCAudience   string
	OidcJWKSURL    string
	OIDCUserRole   string
	OIDCAdminRole  string
	OIDCTimeoutSec int
	OIDCCacheMin   int
	QwenAPIBase    string
	QwenAPIToken   string
	QwenModel      string
	QwenTimeoutSec int
	LLMWorkers     int
	LLMQueueSize   int
	OutputMode     string // "block" or "mark"
	QwenRPS        int
	QwenRetries    int
	LLMCacheTTLMin int
}

// Default returns a minimal runnable configuration for local development.
func Default() Config {
	return Config{
		HTTPPort:       "8080",
		DatabaseURL:    "postgres://localhost:5432/aiguardrails",
		RedisURL:       "redis://localhost:6379",
		KMSKey:         "local-dev",
		AllowedOrigins: []string{"*"},
		AdminToken:     "changeme-admin",
		OIDCIssuer:     "http://localhost/issuer",
		OIDCAudience:   "aiguardrails",
		OidcJWKSURL:    "http://localhost/jwks",
		OIDCUserRole:   "tenant_user",
		OIDCAdminRole:  "tenant_admin",
		OIDCTimeoutSec: 5,
		OIDCCacheMin:   10,
		QwenAPIBase:    "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-moderation",
		QwenAPIToken:   "",
		QwenModel:      "qwen-moderation",
		QwenTimeoutSec: 8,
		LLMWorkers:     2,
		LLMQueueSize:   64,
		OutputMode:     "mark",
		QwenRPS:        5,
		QwenRetries:    2,
		LLMCacheTTLMin: 10,
	}
}

