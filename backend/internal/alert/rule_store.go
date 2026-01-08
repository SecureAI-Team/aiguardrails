package alert

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// AlertRule 告警规则
type AlertRule struct {
	ID                 string          `json:"id"`
	TenantID           *string         `json:"tenant_id,omitempty"`
	Name               string          `json:"name"`
	Description        string          `json:"description"`
	EventTypes         []string        `json:"event_types"`
	SeverityThreshold  string          `json:"severity_threshold"`
	ThresholdCount     int             `json:"threshold_count"`
	ThresholdWindowSec int             `json:"threshold_window_sec"`
	NotifyChannels     []string        `json:"notify_channels"`
	NotifyRecipients   json.RawMessage `json:"notify_recipients"`
	CooldownSec        int             `json:"cooldown_sec"`
	Enabled            bool            `json:"enabled"`
	Priority           int             `json:"priority"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

// AlertHistory 告警历史
type AlertHistory struct {
	ID             string          `json:"id"`
	RuleID         *string         `json:"rule_id,omitempty"`
	RuleName       string          `json:"rule_name"`
	TenantID       *string         `json:"tenant_id,omitempty"`
	Severity       string          `json:"severity"`
	Title          string          `json:"title"`
	Message        string          `json:"message"`
	EventData      json.RawMessage `json:"event_data"`
	NotifyStatus   json.RawMessage `json:"notify_status"`
	Acknowledged   bool            `json:"acknowledged"`
	AcknowledgedBy *string         `json:"acknowledged_by,omitempty"`
	AcknowledgedAt *time.Time      `json:"acknowledged_at,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
}

// NotificationChannel 通知渠道
type NotificationChannel struct {
	ID          string          `json:"id"`
	TenantID    *string         `json:"tenant_id,omitempty"`
	ChannelType string          `json:"channel_type"`
	Name        string          `json:"name"`
	Config      json.RawMessage `json:"config"`
	Enabled     bool            `json:"enabled"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// RuleStore 规则存储
type RuleStore struct {
	db *sql.DB
}

// NewRuleStore 创建规则存储
func NewRuleStore(db *sql.DB) *RuleStore {
	return &RuleStore{db: db}
}

// ListRules 列出告警规则
func (s *RuleStore) ListRules(tenantID *string, enabledOnly bool) ([]AlertRule, error) {
	query := `SELECT id, tenant_id, name, description, event_types, severity_threshold, 
		threshold_count, threshold_window_sec, notify_channels, notify_recipients,
		cooldown_sec, enabled, priority, created_at, updated_at
		FROM alert_rules WHERE (tenant_id IS NULL OR tenant_id = $1)`
	if enabledOnly {
		query += " AND enabled = true"
	}
	query += " ORDER BY priority DESC, created_at DESC"

	var tid interface{}
	if tenantID != nil {
		tid = *tenantID
	}
	rows, err := s.db.Query(query, tid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []AlertRule
	for rows.Next() {
		var r AlertRule
		err := rows.Scan(&r.ID, &r.TenantID, &r.Name, &r.Description,
			pq.Array(&r.EventTypes), &r.SeverityThreshold,
			&r.ThresholdCount, &r.ThresholdWindowSec,
			pq.Array(&r.NotifyChannels), &r.NotifyRecipients,
			&r.CooldownSec, &r.Enabled, &r.Priority, &r.CreatedAt, &r.UpdatedAt)
		if err != nil {
			return nil, err
		}
		rules = append(rules, r)
	}
	return rules, nil
}

// CreateRule 创建告警规则
func (s *RuleStore) CreateRule(rule *AlertRule) error {
	rule.ID = uuid.NewString()
	now := time.Now().UTC()
	rule.CreatedAt = now
	rule.UpdatedAt = now

	_, err := s.db.Exec(`INSERT INTO alert_rules 
		(id, tenant_id, name, description, event_types, severity_threshold, 
		threshold_count, threshold_window_sec, notify_channels, notify_recipients,
		cooldown_sec, enabled, priority, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`,
		rule.ID, rule.TenantID, rule.Name, rule.Description,
		pq.Array(rule.EventTypes), rule.SeverityThreshold,
		rule.ThresholdCount, rule.ThresholdWindowSec,
		pq.Array(rule.NotifyChannels), rule.NotifyRecipients,
		rule.CooldownSec, rule.Enabled, rule.Priority, rule.CreatedAt, rule.UpdatedAt)
	return err
}

// UpdateRule 更新告警规则
func (s *RuleStore) UpdateRule(id string, rule *AlertRule) error {
	rule.UpdatedAt = time.Now().UTC()
	_, err := s.db.Exec(`UPDATE alert_rules SET 
		name=$2, description=$3, event_types=$4, severity_threshold=$5,
		threshold_count=$6, threshold_window_sec=$7, notify_channels=$8, notify_recipients=$9,
		cooldown_sec=$10, enabled=$11, priority=$12, updated_at=$13
		WHERE id=$1`,
		id, rule.Name, rule.Description, pq.Array(rule.EventTypes), rule.SeverityThreshold,
		rule.ThresholdCount, rule.ThresholdWindowSec, pq.Array(rule.NotifyChannels), rule.NotifyRecipients,
		rule.CooldownSec, rule.Enabled, rule.Priority, rule.UpdatedAt)
	return err
}

// DeleteRule 删除告警规则
func (s *RuleStore) DeleteRule(id string) error {
	_, err := s.db.Exec(`DELETE FROM alert_rules WHERE id=$1`, id)
	return err
}

// SaveHistory 保存告警历史
func (s *RuleStore) SaveHistory(h *AlertHistory) error {
	h.ID = uuid.NewString()
	h.CreatedAt = time.Now().UTC()
	_, err := s.db.Exec(`INSERT INTO alert_history 
		(id, rule_id, rule_name, tenant_id, severity, title, message, event_data, notify_status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		h.ID, h.RuleID, h.RuleName, h.TenantID, h.Severity, h.Title, h.Message, h.EventData, h.NotifyStatus, h.CreatedAt)
	return err
}

// ListHistory 列出告警历史
func (s *RuleStore) ListHistory(tenantID *string, severity string, acknowledged *bool, limit int) ([]AlertHistory, error) {
	query := `SELECT id, rule_id, rule_name, tenant_id, severity, title, message, event_data, 
		notify_status, acknowledged, acknowledged_by, acknowledged_at, created_at
		FROM alert_history WHERE 1=1`
	args := []interface{}{}
	argIdx := 1

	if tenantID != nil {
		query += ` AND (tenant_id IS NULL OR tenant_id = $` + string(rune('0'+argIdx)) + `)`
		args = append(args, *tenantID)
		argIdx++
	}
	if severity != "" {
		query += ` AND severity = $` + string(rune('0'+argIdx))
		args = append(args, severity)
		argIdx++
	}
	if acknowledged != nil {
		query += ` AND acknowledged = $` + string(rune('0'+argIdx))
		args = append(args, *acknowledged)
		argIdx++
	}
	query += ` ORDER BY created_at DESC LIMIT $` + string(rune('0'+argIdx))
	args = append(args, limit)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []AlertHistory
	for rows.Next() {
		var h AlertHistory
		err := rows.Scan(&h.ID, &h.RuleID, &h.RuleName, &h.TenantID, &h.Severity, &h.Title, &h.Message,
			&h.EventData, &h.NotifyStatus, &h.Acknowledged, &h.AcknowledgedBy, &h.AcknowledgedAt, &h.CreatedAt)
		if err != nil {
			return nil, err
		}
		alerts = append(alerts, h)
	}
	return alerts, nil
}

// AcknowledgeAlert 确认告警
func (s *RuleStore) AcknowledgeAlert(id, userID string) error {
	now := time.Now().UTC()
	_, err := s.db.Exec(`UPDATE alert_history SET acknowledged=true, acknowledged_by=$2, acknowledged_at=$3 WHERE id=$1`,
		id, userID, now)
	return err
}

// SaveChannel 保存通知渠道
func (s *RuleStore) SaveChannel(c *NotificationChannel) error {
	c.ID = uuid.NewString()
	now := time.Now().UTC()
	c.CreatedAt = now
	c.UpdatedAt = now
	_, err := s.db.Exec(`INSERT INTO notification_channels 
		(id, tenant_id, channel_type, name, config, enabled, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		c.ID, c.TenantID, c.ChannelType, c.Name, c.Config, c.Enabled, c.CreatedAt, c.UpdatedAt)
	return err
}

// ListChannels 列出通知渠道
func (s *RuleStore) ListChannels(tenantID *string) ([]NotificationChannel, error) {
	query := `SELECT id, tenant_id, channel_type, name, config, enabled, created_at, updated_at
		FROM notification_channels WHERE tenant_id IS NULL OR tenant_id = $1`
	var tid interface{}
	if tenantID != nil {
		tid = *tenantID
	}
	rows, err := s.db.Query(query, tid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var channels []NotificationChannel
	for rows.Next() {
		var c NotificationChannel
		err := rows.Scan(&c.ID, &c.TenantID, &c.ChannelType, &c.Name, &c.Config, &c.Enabled, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		channels = append(channels, c)
	}
	return channels, nil
}
