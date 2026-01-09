package server

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"aiguardrails/internal/usage"
)

// registerUsageRoutes 注册用量统计路由
func (s *Server) registerUsageRoutes(r chi.Router) {
	// 用量统计
	r.Get("/usage/stats", s.getUsageStats)
	r.Get("/usage/summary", s.getUsageSummary)
	r.Get("/usage/overview", s.getUsageOverview)

	// API Keys
	r.Get("/apikeys", s.listAPIKeys)
	r.Post("/apikeys", s.createAPIKey)
	r.Delete("/apikeys/{id}", s.revokeAPIKey)

	// 配额
	r.Get("/quota", s.getQuota)
	r.Put("/quota", s.updateQuota)
}

func (s *Server) getUsageStats(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")
	if startDate == "" {
		startDate = time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	stats, err := s.usageStore.GetStats(tenantID, nil, startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusOK, stats)
}

func (s *Server) getUsageSummary(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	days := 7
	if d := r.URL.Query().Get("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil {
			days = parsed
		}
	}

	summary, err := s.usageStore.GetDailySummary(tenantID, days)
	if err != nil {
		log.Printf("getUsageSummary error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusOK, summary)
}

func (s *Server) getUsageOverview(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")

	// 获取今日和本月统计
	today := time.Now().Format("2006-01-02")
	monthStart := time.Now().Format("2006-01") + "-01"

	todayStats, _ := s.usageStore.GetStats(tenantID, nil, today, today)
	monthStats, _ := s.usageStore.GetStats(tenantID, nil, monthStart, today)

	var todayRequests, todayErrors, todayBlocked, todayTokens int64
	var monthRequests, monthTokens int64

	for _, st := range todayStats {
		todayRequests += st.RequestCount
		todayErrors += st.ErrorCount
		todayBlocked += st.BlockedCount
		todayTokens += st.InputTokens + st.OutputTokens
	}
	for _, st := range monthStats {
		monthRequests += st.RequestCount
		monthTokens += st.InputTokens + st.OutputTokens
	}

	// 获取配额
	quota, _ := s.usageStore.GetQuota(tenantID)
	var quotaPercent float64
	if quota != nil && quota.DailyRequestLimit != nil && *quota.DailyRequestLimit > 0 {
		quotaPercent = float64(todayRequests) / float64(*quota.DailyRequestLimit) * 100
	}

	overview := map[string]interface{}{
		"today": map[string]interface{}{
			"requests": todayRequests,
			"errors":   todayErrors,
			"blocked":  todayBlocked,
			"tokens":   todayTokens,
		},
		"month": map[string]interface{}{
			"requests": monthRequests,
			"tokens":   monthTokens,
		},
		"quota_percent": quotaPercent,
	}
	s.writeJSON(w, http.StatusOK, overview)
}

func (s *Server) listAPIKeys(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")

	keys, err := s.usageStore.ListAPIKeys(tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusOK, keys)
}

type createKeyRequest struct {
	Name         string   `json:"name"`
	TenantID     string   `json:"tenant_id"`
	AppID        *string  `json:"app_id,omitempty"`
	Scopes       []string `json:"scopes"`
	RateLimitRPM *int     `json:"rate_limit_rpm,omitempty"`
}

func (s *Server) createAPIKey(w http.ResponseWriter, r *http.Request) {
	var req createKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.TenantID == "" {
		// Do not default, allow empty for platform admin? No, creating API Key requires tenantID.
		// If we are platform admin, we should specify tenantID.
		// If we are tenant admin, auth middleware should verify.
		// For now, let's remove "default". If empty, DB will complain UUID invalid? No, DB constraints.
		// Actually, if req.TenantID is empty here, and we pass it to Store, it will fail if column is NOT NULL.
		// Let's leave it empty and let DB constraints handle it, or check for error.
		// Better to not set default "default".
	}

	// 生成API Key
	keyBytes := make([]byte, 32)
	_, _ = rand.Read(keyBytes)
	fullKey := "sk_" + hex.EncodeToString(keyBytes)

	key := &usage.APIKey{
		TenantID:     req.TenantID,
		AppID:        req.AppID,
		Name:         req.Name,
		Scopes:       req.Scopes,
		RateLimitRPM: req.RateLimitRPM,
		Enabled:      true,
	}

	if err := s.usageStore.CreateAPIKey(key, fullKey); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回完整Key（仅此一次）
	s.writeJSON(w, http.StatusCreated, map[string]interface{}{
		"id":         key.ID,
		"key":        fullKey,
		"key_prefix": key.KeyPrefix,
		"name":       key.Name,
		"message":    "请妥善保管此密钥，它只会显示一次",
	})
}

func (s *Server) revokeAPIKey(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := s.usageStore.RevokeAPIKey(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusOK, map[string]string{"status": "revoked"})
}

func (s *Server) getQuota(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")

	quota, err := s.usageStore.GetQuota(tenantID)
	if err != nil {
		// 返回默认配额
		s.writeJSON(w, http.StatusOK, map[string]interface{}{
			"tenant_id":              tenantID,
			"daily_request_limit":    nil,
			"monthly_request_limit":  nil,
			"current_daily_requests": 0,
		})
		return
	}
	s.writeJSON(w, http.StatusOK, quota)
}

func (s *Server) updateQuota(w http.ResponseWriter, r *http.Request) {
	var q usage.QuotaConfig
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if q.TenantID == "" {
		// q.TenantID = "default"
	}

	if err := s.usageStore.UpdateQuota(q.TenantID, &q); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, http.StatusOK, q)
}
