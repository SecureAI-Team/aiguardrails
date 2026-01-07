package policy

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrRuleNotFound  = errors.New("rule not found")
	ErrRuleDuplicate = errors.New("rule with same name already exists")
)

// TenantRuleStore 租户规则存储
type TenantRuleStore struct {
	db *sql.DB
}

// NewTenantRuleStore 创建存储实例
func NewTenantRuleStore(db *sql.DB) *TenantRuleStore {
	return &TenantRuleStore{db: db}
}

// Create 创建规则
func (s *TenantRuleStore) Create(rule TenantRule) (*TenantRule, error) {
	rule.ID = uuid.NewString()
	rule.CreatedAt = time.Now().UTC()
	rule.UpdatedAt = rule.CreatedAt

	_, err := s.db.Exec(`
		INSERT INTO tenant_rules (id, tenant_id, rule_type, name, description, config, enabled, priority, created_at, updated_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		rule.ID, rule.TenantID, rule.RuleType, rule.Name, rule.Description,
		rule.Config, rule.Enabled, rule.Priority, rule.CreatedAt, rule.UpdatedAt, rule.CreatedBy)
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// Update 更新规则
func (s *TenantRuleStore) Update(rule TenantRule) (*TenantRule, error) {
	rule.UpdatedAt = time.Now().UTC()
	result, err := s.db.Exec(`
		UPDATE tenant_rules 
		SET name=$1, description=$2, config=$3, enabled=$4, priority=$5, updated_at=$6
		WHERE id=$7 AND tenant_id=$8`,
		rule.Name, rule.Description, rule.Config, rule.Enabled, rule.Priority, rule.UpdatedAt,
		rule.ID, rule.TenantID)
	if err != nil {
		return nil, err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return nil, ErrRuleNotFound
	}
	return &rule, nil
}

// Delete 删除规则
func (s *TenantRuleStore) Delete(tenantID, ruleID string) error {
	result, err := s.db.Exec(`DELETE FROM tenant_rules WHERE id=$1 AND tenant_id=$2`, ruleID, tenantID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrRuleNotFound
	}
	return nil
}

// Get 获取单个规则
func (s *TenantRuleStore) Get(tenantID, ruleID string) (*TenantRule, error) {
	row := s.db.QueryRow(`
		SELECT id, tenant_id, rule_type, name, description, config, enabled, priority, created_at, updated_at, created_by
		FROM tenant_rules WHERE id=$1 AND tenant_id=$2`, ruleID, tenantID)
	return s.scanRule(row)
}

// List 列出租户所有规则
func (s *TenantRuleStore) List(tenantID string, ruleType TenantRuleType) ([]TenantRule, error) {
	query := `
		SELECT id, tenant_id, rule_type, name, description, config, enabled, priority, created_at, updated_at, created_by
		FROM tenant_rules WHERE tenant_id=$1`
	args := []interface{}{tenantID}

	if ruleType != "" {
		query += ` AND rule_type=$2`
		args = append(args, ruleType)
	}
	query += ` ORDER BY priority DESC, created_at ASC`

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []TenantRule
	for rows.Next() {
		rule, err := s.scanRuleFromRows(rows)
		if err != nil {
			return nil, err
		}
		rules = append(rules, *rule)
	}
	return rules, nil
}

// ListEnabled 列出启用的规则
func (s *TenantRuleStore) ListEnabled(tenantID string, ruleType TenantRuleType) ([]TenantRule, error) {
	query := `
		SELECT id, tenant_id, rule_type, name, description, config, enabled, priority, created_at, updated_at, created_by
		FROM tenant_rules WHERE tenant_id=$1 AND enabled=true`
	args := []interface{}{tenantID}

	if ruleType != "" {
		query += ` AND rule_type=$2`
		args = append(args, ruleType)
	}
	query += ` ORDER BY priority DESC`

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []TenantRule
	for rows.Next() {
		rule, err := s.scanRuleFromRows(rows)
		if err != nil {
			return nil, err
		}
		rules = append(rules, *rule)
	}
	return rules, nil
}

// ListTemplates 列出规则模板
func (s *TenantRuleStore) ListTemplates(ruleType TenantRuleType) ([]RuleTemplate, error) {
	query := `SELECT id, name, rule_type, description, config_schema, default_config, tags, created_at FROM rule_templates`
	args := []interface{}{}
	if ruleType != "" {
		query += ` WHERE rule_type=$1`
		args = append(args, ruleType)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []RuleTemplate
	for rows.Next() {
		var t RuleTemplate
		var tags []byte
		if err := rows.Scan(&t.ID, &t.Name, &t.RuleType, &t.Description, &t.ConfigSchema, &t.DefaultConfig, &tags, &t.CreatedAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(tags, &t.Tags)
		templates = append(templates, t)
	}
	return templates, nil
}

// GetTemplate 获取单个模板
func (s *TenantRuleStore) GetTemplate(name string) (*RuleTemplate, error) {
	row := s.db.QueryRow(`SELECT id, name, rule_type, description, config_schema, default_config, tags, created_at FROM rule_templates WHERE name=$1`, name)
	var t RuleTemplate
	var tags []byte
	if err := row.Scan(&t.ID, &t.Name, &t.RuleType, &t.Description, &t.ConfigSchema, &t.DefaultConfig, &tags, &t.CreatedAt); err != nil {
		return nil, err
	}
	_ = json.Unmarshal(tags, &t.Tags)
	return &t, nil
}

func (s *TenantRuleStore) scanRule(row *sql.Row) (*TenantRule, error) {
	var r TenantRule
	var desc sql.NullString
	var createdBy sql.NullString
	if err := row.Scan(&r.ID, &r.TenantID, &r.RuleType, &r.Name, &desc, &r.Config, &r.Enabled, &r.Priority, &r.CreatedAt, &r.UpdatedAt, &createdBy); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRuleNotFound
		}
		return nil, err
	}
	r.Description = desc.String
	r.CreatedBy = createdBy.String
	return &r, nil
}

func (s *TenantRuleStore) scanRuleFromRows(rows *sql.Rows) (*TenantRule, error) {
	var r TenantRule
	var desc sql.NullString
	var createdBy sql.NullString
	if err := rows.Scan(&r.ID, &r.TenantID, &r.RuleType, &r.Name, &desc, &r.Config, &r.Enabled, &r.Priority, &r.CreatedAt, &r.UpdatedAt, &createdBy); err != nil {
		return nil, err
	}
	r.Description = desc.String
	r.CreatedBy = createdBy.String
	return &r, nil
}
