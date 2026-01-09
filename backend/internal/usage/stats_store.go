package usage

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// UsageStats 用量统计
type UsageStats struct {
	ID           string    `json:"id"`
	TenantID     string    `json:"tenant_id"`
	AppID        *string   `json:"app_id,omitempty"`
	Date         string    `json:"date"`
	Hour         *int      `json:"hour,omitempty"`
	RequestCount int64     `json:"request_count"`
	SuccessCount int64     `json:"success_count"`
	ErrorCount   int64     `json:"error_count"`
	BlockedCount int64     `json:"blocked_count"`
	LatencyAvg   int       `json:"latency_avg"`
	LatencyP50   int       `json:"latency_p50"`
	LatencyP99   int       `json:"latency_p99"`
	InputTokens  int64     `json:"input_tokens"`
	OutputTokens int64     `json:"output_tokens"`
	CreatedAt    time.Time `json:"created_at"`
}

// APIKey API密钥
type APIKey struct {
	ID           string     `json:"id"`
	TenantID     string     `json:"tenant_id"`
	AppID        *string    `json:"app_id,omitempty"`
	Name         string     `json:"name"`
	KeyPrefix    string     `json:"key_prefix"`
	Scopes       []string   `json:"scopes"`
	IPWhitelist  []string   `json:"ip_whitelist"`
	RateLimitRPM *int       `json:"rate_limit_rpm,omitempty"`
	RateLimitRPD *int       `json:"rate_limit_rpd,omitempty"`
	Enabled      bool       `json:"enabled"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
	LastUsedAt   *time.Time `json:"last_used_at,omitempty"`
	CreatedBy    *string    `json:"created_by,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	RevokedAt    *time.Time `json:"revoked_at,omitempty"`
}

// QuotaConfig 配额配置
type QuotaConfig struct {
	ID                     string    `json:"id"`
	TenantID               string    `json:"tenant_id"`
	DailyRequestLimit      *int64    `json:"daily_request_limit,omitempty"`
	MonthlyRequestLimit    *int64    `json:"monthly_request_limit,omitempty"`
	DailyTokenLimit        *int64    `json:"daily_token_limit,omitempty"`
	MonthlyTokenLimit      *int64    `json:"monthly_token_limit,omitempty"`
	CurrentDailyRequests   int64     `json:"current_daily_requests"`
	CurrentMonthlyRequests int64     `json:"current_monthly_requests"`
	CurrentDailyTokens     int64     `json:"current_daily_tokens"`
	CurrentMonthlyTokens   int64     `json:"current_monthly_tokens"`
	AlertThreshold         int       `json:"alert_threshold_percent"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

// UsageStore 用量存储
type UsageStore struct {
	db *sql.DB
}

// NewUsageStore 创建用量存储
func NewUsageStore(db *sql.DB) *UsageStore {
	return &UsageStore{db: db}
}

// RecordUsage 记录用量
func (s *UsageStore) RecordUsage(tenantID, appID string, success bool, blocked bool, latencyMs, inputTokens, outputTokens int) error {
	date := time.Now().Format("2006-01-02")
	hour := time.Now().Hour()

	var successInc, errorInc, blockedInc int
	if success {
		successInc = 1
	} else {
		errorInc = 1
	}
	if blocked {
		blockedInc = 1
	}

	_, err := s.db.Exec(`
		INSERT INTO api_usage_stats (tenant_id, app_id, date, hour, request_count, success_count, error_count, blocked_count, latency_sum, input_tokens, output_tokens)
		VALUES ($1, $2, $3, $4, 1, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (tenant_id, app_id, date, hour) DO UPDATE SET
			request_count = api_usage_stats.request_count + 1,
			success_count = api_usage_stats.success_count + $5,
			error_count = api_usage_stats.error_count + $6,
			blocked_count = api_usage_stats.blocked_count + $7,
			latency_sum = api_usage_stats.latency_sum + $8,
			input_tokens = api_usage_stats.input_tokens + $9,
			output_tokens = api_usage_stats.output_tokens + $10,
			updated_at = NOW()`,
		tenantID, appID, date, hour, successInc, errorInc, blockedInc, latencyMs, inputTokens, outputTokens)
	return err
}

// GetStats 获取统计数据
func (s *UsageStore) GetStats(tenantID string, appID *string, startDate, endDate string) ([]UsageStats, error) {
	query := `SELECT id, tenant_id, app_id, date, hour, request_count, success_count, error_count, blocked_count,
		CASE WHEN request_count > 0 THEN latency_sum / request_count ELSE 0 END as latency_avg,
		latency_p50, latency_p99, input_tokens, output_tokens, created_at
		FROM api_usage_stats WHERE tenant_id = $1 AND date >= $2 AND date <= $3`
	args := []interface{}{tenantID, startDate, endDate}

	if appID != nil {
		query += " AND app_id = $4"
		args = append(args, *appID)
	}
	query += " ORDER BY date DESC, hour DESC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []UsageStats
	for rows.Next() {
		var st UsageStats
		err := rows.Scan(&st.ID, &st.TenantID, &st.AppID, &st.Date, &st.Hour,
			&st.RequestCount, &st.SuccessCount, &st.ErrorCount, &st.BlockedCount,
			&st.LatencyAvg, &st.LatencyP50, &st.LatencyP99, &st.InputTokens, &st.OutputTokens, &st.CreatedAt)
		if err != nil {
			return nil, err
		}
		stats = append(stats, st)
	}
	return stats, nil
}

// GetDailySummary 获取日汇总
	startDate := time.Now().AddDate(0, 0, -days).Format("2006-01-02")
	rows, err := s.db.Query(`
		SELECT date, SUM(request_count) as requests, SUM(success_count) as success,
			SUM(error_count) as errors, SUM(blocked_count) as blocked,
			SUM(input_tokens) as input_tokens, SUM(output_tokens) as output_tokens
		FROM api_usage_stats WHERE tenant_id = $1 AND date >= $2
		GROUP BY date ORDER BY date DESC`, tenantID, startDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var date string
		var requests, success, errors, blocked, inputTokens, outputTokens int64
		if err := rows.Scan(&date, &requests, &success, &errors, &blocked, &inputTokens, &outputTokens); err != nil {
			return nil, err
		}
		result = append(result, map[string]interface{}{
			"date": date, "requests": requests, "success": success, "errors": errors,
			"blocked": blocked, "input_tokens": inputTokens, "output_tokens": outputTokens,
		})
	}
	return result, nil
}

// CreateAPIKey 创建API Key
func (s *UsageStore) CreateAPIKey(key *APIKey, fullKey string) error {
	key.ID = uuid.NewString()
	key.KeyPrefix = fullKey[:8]
	key.CreatedAt = time.Now().UTC()

	hash := sha256.Sum256([]byte(fullKey))
	keyHash := hex.EncodeToString(hash[:])

	_, err := s.db.Exec(`INSERT INTO api_keys 
		(id, tenant_id, app_id, name, key_prefix, key_hash, scopes, ip_whitelist, rate_limit_rpm, rate_limit_rpd, enabled, expires_at, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`,
		key.ID, key.TenantID, key.AppID, key.Name, key.KeyPrefix, keyHash,
		pq.Array(key.Scopes), pq.Array(key.IPWhitelist),
		key.RateLimitRPM, key.RateLimitRPD, key.Enabled, key.ExpiresAt, key.CreatedBy, key.CreatedAt)
	return err
}

// ListAPIKeys 列出API Keys
func (s *UsageStore) ListAPIKeys(tenantID string) ([]APIKey, error) {
	rows, err := s.db.Query(`SELECT id, tenant_id, app_id, name, key_prefix, scopes, ip_whitelist, 
		rate_limit_rpm, rate_limit_rpd, enabled, expires_at, last_used_at, created_by, created_at, revoked_at
		FROM api_keys WHERE tenant_id = $1 AND revoked_at IS NULL ORDER BY created_at DESC`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []APIKey
	for rows.Next() {
		var k APIKey
		err := rows.Scan(&k.ID, &k.TenantID, &k.AppID, &k.Name, &k.KeyPrefix,
			pq.Array(&k.Scopes), pq.Array(&k.IPWhitelist),
			&k.RateLimitRPM, &k.RateLimitRPD, &k.Enabled, &k.ExpiresAt, &k.LastUsedAt, &k.CreatedBy, &k.CreatedAt, &k.RevokedAt)
		if err != nil {
			return nil, err
		}
		keys = append(keys, k)
	}
	return keys, nil
}

// RevokeAPIKey 吊销API Key
func (s *UsageStore) RevokeAPIKey(id string) error {
	_, err := s.db.Exec(`UPDATE api_keys SET revoked_at = NOW(), enabled = false WHERE id = $1`, id)
	return err
}

// ValidateAPIKey 验证API Key
func (s *UsageStore) ValidateAPIKey(fullKey string) (*APIKey, error) {
	hash := sha256.Sum256([]byte(fullKey))
	keyHash := hex.EncodeToString(hash[:])

	var k APIKey
	err := s.db.QueryRow(`SELECT id, tenant_id, app_id, name, key_prefix, scopes, ip_whitelist,
		rate_limit_rpm, rate_limit_rpd, enabled, expires_at, created_at
		FROM api_keys WHERE key_hash = $1 AND revoked_at IS NULL AND enabled = true`, keyHash).
		Scan(&k.ID, &k.TenantID, &k.AppID, &k.Name, &k.KeyPrefix,
			pq.Array(&k.Scopes), pq.Array(&k.IPWhitelist),
			&k.RateLimitRPM, &k.RateLimitRPD, &k.Enabled, &k.ExpiresAt, &k.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Update last used
	_, _ = s.db.Exec(`UPDATE api_keys SET last_used_at = NOW() WHERE id = $1`, k.ID)
	return &k, nil
}

// GetQuota 获取配额
func (s *UsageStore) GetQuota(tenantID string) (*QuotaConfig, error) {
	var q QuotaConfig
	err := s.db.QueryRow(`SELECT id, tenant_id, daily_request_limit, monthly_request_limit,
		daily_token_limit, monthly_token_limit, current_daily_requests, current_monthly_requests,
		current_daily_tokens, current_monthly_tokens, alert_threshold_percent, created_at, updated_at
		FROM quota_config WHERE tenant_id = $1`, tenantID).
		Scan(&q.ID, &q.TenantID, &q.DailyRequestLimit, &q.MonthlyRequestLimit,
			&q.DailyTokenLimit, &q.MonthlyTokenLimit, &q.CurrentDailyRequests, &q.CurrentMonthlyRequests,
			&q.CurrentDailyTokens, &q.CurrentMonthlyTokens, &q.AlertThreshold, &q.CreatedAt, &q.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &q, nil
}

// UpdateQuota 更新配额配置
func (s *UsageStore) UpdateQuota(tenantID string, q *QuotaConfig) error {
	_, err := s.db.Exec(`INSERT INTO quota_config (tenant_id, daily_request_limit, monthly_request_limit, daily_token_limit, monthly_token_limit, alert_threshold_percent)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (tenant_id) DO UPDATE SET
			daily_request_limit = $2, monthly_request_limit = $3,
			daily_token_limit = $4, monthly_token_limit = $5,
			alert_threshold_percent = $6, updated_at = NOW()`,
		tenantID, q.DailyRequestLimit, q.MonthlyRequestLimit, q.DailyTokenLimit, q.MonthlyTokenLimit, q.AlertThreshold)
	return err
}
