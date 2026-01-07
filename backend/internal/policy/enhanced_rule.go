package policy

import "time"

// RuleCategory 规则类别
type RuleCategory string

const (
	RuleCategoryLegal      RuleCategory = "legal"      // 法规合规
	RuleCategoryBusiness   RuleCategory = "business"   // 业务规则
	RuleCategoryPermission RuleCategory = "permission" // 操作权限
	RuleCategorySafety     RuleCategory = "safety"     // 安全规则
	RuleCategoryVendor     RuleCategory = "vendor"     // 厂商限制
)

// RuleSeverity 规则严重度
type RuleSeverity string

const (
	SeverityLow      RuleSeverity = "low"
	SeverityMedium   RuleSeverity = "medium"
	SeverityHigh     RuleSeverity = "high"
	SeverityCritical RuleSeverity = "critical"
)

// RuleDecision 规则决策
type RuleDecision string

const (
	DecisionAllow   RuleDecision = "allow"   // 允许
	DecisionBlock   RuleDecision = "block"   // 阻断
	DecisionMark    RuleDecision = "mark"    // 标记人工审核
	DecisionRedact  RuleDecision = "redact"  // 脱敏
	DecisionDeflect RuleDecision = "deflect" // 转移话题
	DecisionConfirm RuleDecision = "confirm" // 需要确认
)

// EnhancedRule 增强规则结构
type EnhancedRule struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Category     RuleCategory    `json:"category"`
	Jurisdiction string          `json:"jurisdiction"` // CN, EU, US, Global
	Regulation   string          `json:"regulation"`
	Vendor       string          `json:"vendor"`
	Product      string          `json:"product"`
	Version      string          `json:"version"`
	Severity     RuleSeverity    `json:"severity"`
	Decision     RuleDecision    `json:"decision"`
	Priority     int             `json:"priority"` // 规则优先级，数字越大优先级越高
	Enabled      bool            `json:"enabled"`
	Conditions   []RuleCondition `json:"conditions"`
	Actions      []RuleAction    `json:"actions"`
	Tags         []string        `json:"tags"`
	References   []string        `json:"references"`
	Description  string          `json:"description"`
	Remediation  string          `json:"remediation"`
	ResponseText string          `json:"response_text,omitempty"` // 自定义响应文本
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// RuleCondition 规则条件
type RuleCondition struct {
	Field         string   `json:"field"`            // prompt | output | tool | role | vendor
	Operator      string   `json:"operator"`         // contains | equals | regex | not_contains | in | not_in
	Value         string   `json:"value"`            // 条件值
	Values        []string `json:"values,omitempty"` // 多值条件
	CaseSensitive bool     `json:"case_sensitive"`
}

// RuleAction 规则动作
type RuleAction struct {
	Type          string `json:"type"`           // block | mark | redact | deflect | confirm | notify | log
	Response      string `json:"response"`       // 可选自定义响应
	Notify        bool   `json:"notify"`         // 是否触发告警
	NotifyLevel   string `json:"notify_level"`   // 告警级别: info | warning | critical
	RedactPattern string `json:"redact_pattern"` // 脱敏模式
}

// VendorConfig 厂商配置
type VendorConfig struct {
	Enabled         bool              `json:"enabled"`
	Mode            string            `json:"mode"` // strict | loose
	AllowedVendors  []string          `json:"allowed_vendors"`
	AllowedProducts []string          `json:"allowed_products"`
	BlockedVendors  []string          `json:"blocked_vendors"`
	BlockedProducts []string          `json:"blocked_products"`
	Responses       map[string]string `json:"responses"`    // 自定义响应
	Alternatives    map[string]string `json:"alternatives"` // 竞品替代推荐
}

// PermissionConfig 权限配置
type PermissionConfig struct {
	Roles           map[string]RoleDefinition `json:"roles"`
	ToolPermissions map[string]ToolPermission `json:"tool_permissions"`
}

// RoleDefinition 角色定义
type RoleDefinition struct {
	Level       int      `json:"level"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

// ToolPermission 工具权限
type ToolPermission struct {
	MinLevel             int  `json:"min_level"`
	RequiresConfirmation bool `json:"requires_confirmation"`
	RequiresMFA          bool `json:"requires_mfa"`
	WorkHoursOnly        bool `json:"work_hours_only"`
}

// RuleEvaluationResult 规则评估结果
type RuleEvaluationResult struct {
	RuleID       string       `json:"rule_id"`
	RuleName     string       `json:"rule_name"`
	Matched      bool         `json:"matched"`
	Decision     RuleDecision `json:"decision"`
	Reason       string       `json:"reason"`
	Signals      []string     `json:"signals"`
	Response     string       `json:"response,omitempty"`
	Severity     RuleSeverity `json:"severity"`
	ShouldNotify bool         `json:"should_notify"`
}

// RuleSet 规则集
type RuleSet struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Version     string         `json:"version"`
	Rules       []EnhancedRule `json:"rules"`
	Enabled     bool           `json:"enabled"`
	Priority    int            `json:"priority"`
}
