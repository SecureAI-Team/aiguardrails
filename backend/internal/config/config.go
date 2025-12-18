package config

import (
	"os"
	"strconv"
	"strings"
)

// Config holds runtime configuration for the platform.
type Config struct {
	HTTPPort       string
	DatabaseURL    string
	RedisURL       string
	RedisNamespace string
	KMSKey         string
	AdminSecretKey string
	QwenSecretKey  string
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
		RedisNamespace: "aiguardrails",
		KMSKey:         "local-dev",
		AdminSecretKey: "ADMIN_TOKEN",
		QwenSecretKey:  "QWEN_API_TOKEN",
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

// FromEnv overlays environment variables onto Default.
func FromEnv() Config {
	cfg := Default()
	if v := os.Getenv("HTTP_PORT"); v != "" {
		cfg.HTTPPort = v
	}
	if v := os.Getenv("DATABASE_URL"); v != "" {
		cfg.DatabaseURL = v
	}
	if v := os.Getenv("REDIS_URL"); v != "" {
		cfg.RedisURL = v
	}
	if v := os.Getenv("REDIS_NAMESPACE"); v != "" {
		cfg.RedisNamespace = v
	}
	if v := os.Getenv("ADMIN_SECRET_KEY"); v != "" {
		cfg.AdminSecretKey = v
	}
	if v := os.Getenv("QWEN_SECRET_KEY"); v != "" {
		cfg.QwenSecretKey = v
	}
	if v := os.Getenv("ADMIN_TOKEN"); v != "" {
		cfg.AdminToken = v
	}
	if v := os.Getenv("ALLOWED_ORIGINS"); v != "" {
		cfg.AllowedOrigins = parseCSV(v)
	}
	if v := os.Getenv("OIDC_ISSUER"); v != "" {
		cfg.OIDCIssuer = v
	}
	if v := os.Getenv("OIDC_AUDIENCE"); v != "" {
		cfg.OIDCAudience = v
	}
	if v := os.Getenv("OIDC_JWKS_URL"); v != "" {
		cfg.OidcJWKSURL = v
	}
	if v := os.Getenv("OIDC_ADMIN_ROLE"); v != "" {
		cfg.OIDCAdminRole = v
	}
	if v := os.Getenv("OIDC_USER_ROLE"); v != "" {
		cfg.OIDCUserRole = v
	}
	if v := os.Getenv("OIDC_TIMEOUT_SEC"); v != "" {
		cfg.OIDCTimeoutSec = atoiDefault(v, cfg.OIDCTimeoutSec)
	}
	if v := os.Getenv("OIDC_CACHE_MIN"); v != "" {
		cfg.OIDCCacheMin = atoiDefault(v, cfg.OIDCCacheMin)
	}
	if v := os.Getenv("QWEN_API_BASE"); v != "" {
		cfg.QwenAPIBase = v
	}
	if v := os.Getenv("QWEN_API_TOKEN"); v != "" {
		cfg.QwenAPIToken = v
	}
	if v := os.Getenv("QWEN_MODEL"); v != "" {
		cfg.QwenModel = v
	}
	if v := os.Getenv("QWEN_TIMEOUT_SEC"); v != "" {
		cfg.QwenTimeoutSec = atoiDefault(v, cfg.QwenTimeoutSec)
	}
	if v := os.Getenv("QWEN_RPS"); v != "" {
		cfg.QwenRPS = atoiDefault(v, cfg.QwenRPS)
	}
	if v := os.Getenv("QWEN_RETRIES"); v != "" {
		cfg.QwenRetries = atoiDefault(v, cfg.QwenRetries)
	}
	if v := os.Getenv("LLM_WORKERS"); v != "" {
		cfg.LLMWorkers = atoiDefault(v, cfg.LLMWorkers)
	}
	if v := os.Getenv("LLM_QUEUE_SIZE"); v != "" {
		cfg.LLMQueueSize = atoiDefault(v, cfg.LLMQueueSize)
	}
	if v := os.Getenv("LLM_CACHE_TTL_MIN"); v != "" {
		cfg.LLMCacheTTLMin = atoiDefault(v, cfg.LLMCacheTTLMin)
	}
	if v := os.Getenv("OUTPUT_MODE"); v != "" {
		cfg.OutputMode = v
	}
	return cfg
}

func parseCSV(v string) []string {
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if s := strings.TrimSpace(p); s != "" {
			out = append(out, s)
		}
	}
	return out
}

func atoiDefault(v string, def int) int {
	if i, err := strconv.Atoi(v); err == nil {
		return i
	}
	return def
}
