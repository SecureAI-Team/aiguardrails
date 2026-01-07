package policy

import (
	"encoding/json"
	"time"
)

// TenantRuleType 租户规则类型
type TenantRuleType string

const (
	RuleTypeBusiness   TenantRuleType = "business"   // 业务规则
	RuleTypePermission TenantRuleType = "permission" // 权限规则
)

// TenantRule 租户规则
type TenantRule struct {
	ID          string          `json:"id"`
	TenantID    string          `json:"tenant_id"`
	RuleType    TenantRuleType  `json:"rule_type"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Config      json.RawMessage `json:"config"`
	Enabled     bool            `json:"enabled"`
	Priority    int             `json:"priority"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	CreatedBy   string          `json:"created_by,omitempty"`
}

// RuleTemplate 规则模板
type RuleTemplate struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	RuleType      TenantRuleType  `json:"rule_type"`
	Description   string          `json:"description"`
	ConfigSchema  json.RawMessage `json:"config_schema"`
	DefaultConfig json.RawMessage `json:"default_config"`
	Tags          []string        `json:"tags"`
	CreatedAt     time.Time       `json:"created_at"`
}

// VendorRuleConfig 厂商规则配置
type VendorRuleConfig struct {
	Enabled         bool              `json:"enabled"`
	Mode            string            `json:"mode"` // exclusive | inclusive
	AllowedVendors  []string          `json:"allowed_vendors"`
	BlockedVendors  []string          `json:"blocked_vendors"`
	AllowedProducts []string          `json:"allowed_products,omitempty"`
	BlockedProducts []string          `json:"blocked_products,omitempty"`
	Responses       map[string]string `json:"responses"`
}

// DomainRuleConfig 领域边界配置
type DomainRuleConfig struct {
	Enabled       bool              `json:"enabled"`
	AllowedTopics []string          `json:"allowed_topics"`
	BlockedTopics []string          `json:"blocked_topics"`
	Responses     map[string]string `json:"responses"`
}

// PermissionRuleConfig 权限规则配置
type PermissionRuleConfig struct {
	Roles           map[string]RoleConfig     `json:"roles"`
	ToolPermissions map[string]ToolPermConfig `json:"tool_permissions,omitempty"`
}

// RoleConfig 角色配置
type RoleConfig struct {
	Level       int      `json:"level"`
	Name        string   `json:"name,omitempty"`
	Permissions []string `json:"permissions"`
}

// ToolPermConfig 工具权限配置
type ToolPermConfig struct {
	MinLevel             int  `json:"min_level"`
	RequiresConfirmation bool `json:"requires_confirmation,omitempty"`
	RequiresMFA          bool `json:"requires_mfa,omitempty"`
}

// ParseVendorConfig 解析厂商规则配置
func (r *TenantRule) ParseVendorConfig() (*VendorRuleConfig, error) {
	var cfg VendorRuleConfig
	if err := json.Unmarshal(r.Config, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// ParseDomainConfig 解析领域边界配置
func (r *TenantRule) ParseDomainConfig() (*DomainRuleConfig, error) {
	var cfg DomainRuleConfig
	if err := json.Unmarshal(r.Config, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// ParsePermissionConfig 解析权限规则配置
func (r *TenantRule) ParsePermissionConfig() (*PermissionRuleConfig, error) {
	var cfg PermissionRuleConfig
	if err := json.Unmarshal(r.Config, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// ToOPAInput 转换为OPA输入格式
func (r *TenantRule) ToOPAInput() map[string]interface{} {
	var configMap map[string]interface{}
	_ = json.Unmarshal(r.Config, &configMap)
	return map[string]interface{}{
		"rule_id":   r.ID,
		"rule_type": string(r.RuleType),
		"name":      r.Name,
		"enabled":   r.Enabled,
		"priority":  r.Priority,
		"config":    configMap,
	}
}
