package tracing

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// RequestTrace 请求追踪
type RequestTrace struct {
	ID           string          `json:"id"`
	TraceID      string          `json:"trace_id"`
	SpanID       string          `json:"span_id,omitempty"`
	ParentSpanID string          `json:"parent_span_id,omitempty"`
	TenantID     string          `json:"tenant_id"`
	AppID        *string         `json:"app_id,omitempty"`
	Method       string          `json:"method"`
	Path         string          `json:"path"`
	QueryParams  json.RawMessage `json:"query_params,omitempty"`
	Headers      json.RawMessage `json:"headers,omitempty"`
	RequestBody  string          `json:"request_body,omitempty"`
	StatusCode   int             `json:"status_code"`
	ResponseBody string          `json:"response_body,omitempty"`
	StartTime    time.Time       `json:"start_time"`
	EndTime      *time.Time      `json:"end_time,omitempty"`
	DurationMs   int             `json:"duration_ms"`
	Stages       json.RawMessage `json:"stages,omitempty"`
	Blocked      bool            `json:"blocked"`
	BlockReason  string          `json:"block_reason,omitempty"`
	Signals      []string        `json:"signals,omitempty"`
	InputTokens  int             `json:"input_tokens"`
	OutputTokens int             `json:"output_tokens"`
	Error        string          `json:"error,omitempty"`
	UserAgent    string          `json:"user_agent,omitempty"`
	ClientIP     string          `json:"client_ip,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
}

// Stage 处理阶段
type Stage struct {
	Name   string `json:"name"`
	Start  int64  `json:"start"`
	End    int64  `json:"end"`
	Status string `json:"status"`
}

// CostRecord 成本记录
type CostRecord struct {
	ID              string          `json:"id"`
	TenantID        string          `json:"tenant_id"`
	AppID           *string         `json:"app_id,omitempty"`
	Date            string          `json:"date"`
	RequestCount    int64           `json:"request_count"`
	InputTokens     int64           `json:"input_tokens"`
	OutputTokens    int64           `json:"output_tokens"`
	InputCostCents  int64           `json:"input_cost_cents"`
	OutputCostCents int64           `json:"output_cost_cents"`
	TotalCostCents  int64           `json:"total_cost_cents"`
	PricingConfig   json.RawMessage `json:"pricing_config,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
}

// ModelInfo 模型信息
type ModelInfo struct {
	ID              string   `json:"id"`
	Provider        string   `json:"provider"`
	ModelID         string   `json:"model_id"`
	DisplayName     string   `json:"display_name"`
	Description     string   `json:"description"`
	Capabilities    []string `json:"capabilities"`
	ContextWindow   int      `json:"context_window"`
	MaxOutputTokens int      `json:"max_output_tokens"`
	InputPricePerM  float64  `json:"input_price_per_m"`
	OutputPricePerM float64  `json:"output_price_per_m"`
	Enabled         bool     `json:"enabled"`
	Deprecated      bool     `json:"deprecated"`
}

// ExportJob 导出任务
type ExportJob struct {
	ID          string          `json:"id"`
	TenantID    string          `json:"tenant_id"`
	Type        string          `json:"type"`
	Status      string          `json:"status"`
	Filters     json.RawMessage `json:"filters,omitempty"`
	FilePath    string          `json:"file_path,omitempty"`
	FileSize    int64           `json:"file_size,omitempty"`
	RowCount    int             `json:"row_count,omitempty"`
	Error       string          `json:"error,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	CompletedAt *time.Time      `json:"completed_at,omitempty"`
}

// Store 追踪存储
type Store struct {
	db *sql.DB
}

// NewStore 创建追踪存储
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// SaveTrace 保存追踪
func (s *Store) SaveTrace(t *RequestTrace) error {
	t.ID = uuid.NewString()
	t.CreatedAt = time.Now().UTC()

	_, err := s.db.Exec(`INSERT INTO request_traces 
		(id, trace_id, span_id, parent_span_id, tenant_id, app_id, method, path, query_params, headers, request_body,
		status_code, response_body, start_time, end_time, duration_ms, stages, blocked, block_reason, signals,
		input_tokens, output_tokens, error, user_agent, client_ip, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26)`,
		t.ID, t.TraceID, t.SpanID, t.ParentSpanID, t.TenantID, t.AppID, t.Method, t.Path, t.QueryParams, t.Headers, t.RequestBody,
		t.StatusCode, t.ResponseBody, t.StartTime, t.EndTime, t.DurationMs, t.Stages, t.Blocked, t.BlockReason, pq.Array(t.Signals),
		t.InputTokens, t.OutputTokens, t.Error, t.UserAgent, t.ClientIP, t.CreatedAt)
	return err
}

// ListTraces 列出追踪
func (s *Store) ListTraces(tenantID string, blocked *bool, statusCode *int, limit int) ([]RequestTrace, error) {
	query := `SELECT id, trace_id, span_id, tenant_id, app_id, method, path, status_code, 
		start_time, duration_ms, blocked, block_reason, signals, input_tokens, output_tokens, error, created_at
		FROM request_traces WHERE 1=1`
	args := []interface{}{}
	argIdx := 1

	if tenantID != "" {
		query += ` AND tenant_id = $` + string(rune('0'+argIdx))
		args = append(args, tenantID)
		argIdx++
	}

	if blocked != nil {
		query += ` AND blocked = $` + string(rune('0'+argIdx))
		args = append(args, *blocked)
		argIdx++
	}
	if statusCode != nil {
		query += ` AND status_code = $` + string(rune('0'+argIdx))
		args = append(args, *statusCode)
		argIdx++
	}
	query += ` ORDER BY created_at DESC LIMIT $` + string(rune('0'+argIdx))
	args = append(args, limit)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var traces []RequestTrace
	for rows.Next() {
		var t RequestTrace
		err := rows.Scan(&t.ID, &t.TraceID, &t.SpanID, &t.TenantID, &t.AppID, &t.Method, &t.Path, &t.StatusCode,
			&t.StartTime, &t.DurationMs, &t.Blocked, &t.BlockReason, pq.Array(&t.Signals), &t.InputTokens, &t.OutputTokens, &t.Error, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		traces = append(traces, t)
	}
	return traces, nil
}

// GetTrace 获取追踪详情
func (s *Store) GetTrace(id string) (*RequestTrace, error) {
	var t RequestTrace
	err := s.db.QueryRow(`SELECT id, trace_id, span_id, parent_span_id, tenant_id, app_id, method, path, 
		query_params, headers, request_body, status_code, response_body, start_time, end_time, duration_ms,
		stages, blocked, block_reason, signals, input_tokens, output_tokens, error, user_agent, client_ip, created_at
		FROM request_traces WHERE id = $1`, id).
		Scan(&t.ID, &t.TraceID, &t.SpanID, &t.ParentSpanID, &t.TenantID, &t.AppID, &t.Method, &t.Path,
			&t.QueryParams, &t.Headers, &t.RequestBody, &t.StatusCode, &t.ResponseBody, &t.StartTime, &t.EndTime, &t.DurationMs,
			&t.Stages, &t.Blocked, &t.BlockReason, pq.Array(&t.Signals), &t.InputTokens, &t.OutputTokens, &t.Error, &t.UserAgent, &t.ClientIP, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// ListModels 列出模型
func (s *Store) ListModels(provider string) ([]ModelInfo, error) {
	query := `SELECT id, provider, model_id, display_name, description, capabilities, 
		context_window, COALESCE(max_output_tokens, 0), COALESCE(input_price_per_m, 0), COALESCE(output_price_per_m, 0), enabled, deprecated
		FROM model_catalog WHERE enabled = true`
	args := []interface{}{}
	if provider != "" {
		query += " AND provider = $1"
		args = append(args, provider)
	}
	query += " ORDER BY provider, display_name"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []ModelInfo
	for rows.Next() {
		var m ModelInfo
		err := rows.Scan(&m.ID, &m.Provider, &m.ModelID, &m.DisplayName, &m.Description, pq.Array(&m.Capabilities),
			&m.ContextWindow, &m.MaxOutputTokens, &m.InputPricePerM, &m.OutputPricePerM, &m.Enabled, &m.Deprecated)
		if err != nil {
			return nil, err
		}
		models = append(models, m)
	}
	return models, nil
}

// GetCostSummary 获取成本汇总
func (s *Store) GetCostSummary(tenantID string, startDate, endDate string) ([]CostRecord, error) {
	rows, err := s.db.Query(`SELECT id, tenant_id, app_id, date, request_count, input_tokens, output_tokens,
		input_cost_cents, output_cost_cents, total_cost_cents, created_at
		FROM cost_records WHERE tenant_id = $1 AND date >= $2 AND date <= $3 ORDER BY date DESC`,
		tenantID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []CostRecord
	for rows.Next() {
		var r CostRecord
		err := rows.Scan(&r.ID, &r.TenantID, &r.AppID, &r.Date, &r.RequestCount, &r.InputTokens, &r.OutputTokens,
			&r.InputCostCents, &r.OutputCostCents, &r.TotalCostCents, &r.CreatedAt)
		if err != nil {
			return nil, err
		}
		records = append(records, r)
	}
	return records, nil
}

// CreateExportJob 创建导出任务
func (s *Store) CreateExportJob(tenantID, jobType string, filters json.RawMessage) (*ExportJob, error) {
	job := &ExportJob{
		ID:        uuid.NewString(),
		TenantID:  tenantID,
		Type:      jobType,
		Status:    "pending",
		Filters:   filters,
		CreatedAt: time.Now().UTC(),
	}
	_, err := s.db.Exec(`INSERT INTO export_jobs (id, tenant_id, type, status, filters, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`, job.ID, job.TenantID, job.Type, job.Status, job.Filters, job.CreatedAt)
	if err != nil {
		return nil, err
	}
	return job, nil
}

// ListExportJobs 列出导出任务
func (s *Store) ListExportJobs(tenantID string) ([]ExportJob, error) {
	rows, err := s.db.Query(`SELECT id, tenant_id, type, status, filters, file_path, file_size, row_count, error, created_at, completed_at
		FROM export_jobs WHERE tenant_id = $1 ORDER BY created_at DESC LIMIT 50`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []ExportJob
	for rows.Next() {
		var j ExportJob
		err := rows.Scan(&j.ID, &j.TenantID, &j.Type, &j.Status, &j.Filters, &j.FilePath, &j.FileSize, &j.RowCount, &j.Error, &j.CreatedAt, &j.CompletedAt)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	return jobs, nil
}
