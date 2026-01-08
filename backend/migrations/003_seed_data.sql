-- AI GuardRails 默认数据初始化
-- 运行此脚本以填充初始演示数据

-- ============================================
-- 规则库 (Rules) - 合规规则和安全规则
-- ============================================
INSERT INTO rules (id, name, description, type, severity, decision, jurisdiction, regulation, vendor, product, rego_module, created_at) VALUES
-- 法规合规规则
('rule-gdpr-pii-01', 'GDPR 个人数据脱敏', '检测并脱敏欧盟公民的个人身份信息', 'output_filter', 'high', 'mask', 'EU', 'GDPR', NULL, NULL, 'package gdpr_pii

default allow = true

deny[msg] {
    contains(input.output, "passport")
    msg := "检测到护照信息，需要脱敏"
}', NOW()),

('rule-pipl-china-01', 'PIPL 个人信息保护', '检测中国公民敏感个人信息', 'output_filter', 'high', 'mask', 'CN', 'PIPL', NULL, NULL, 'package pipl

default allow = true

deny[msg] {
    regex.match(`\d{18}|\d{17}X`, input.output)
    msg := "检测到身份证号，需要脱敏"
}', NOW()),

('rule-csl-china-01', 'CSL 网络安全合规', '网络安全法合规检查', 'prompt_check', 'critical', 'block', 'CN', 'CSL', NULL, NULL, 'package csl

default allow = true

deny[msg] {
    contains(lower(input.prompt), "vpn")
    msg := "检测到敏感网络工具讨论"
}', NOW()),

-- 安全规则
('rule-jailbreak-01', '越狱攻击检测', '检测常见的LLM越狱攻击模式', 'prompt_check', 'critical', 'block', NULL, NULL, NULL, NULL, 'package jailbreak

default allow = true

deny[msg] {
    contains(lower(input.prompt), "ignore previous instructions")
    msg := "检测到越狱攻击尝试"
}

deny[msg] {
    contains(lower(input.prompt), "你是")
    contains(lower(input.prompt), "假装")
    msg := "检测到角色扮演越狱"
}', NOW()),

('rule-injection-01', '提示注入检测', '检测提示注入攻击', 'prompt_check', 'high', 'block', NULL, NULL, NULL, NULL, 'package injection

default allow = true

deny[msg] {
    contains(input.prompt, "]]>")
    msg := "检测到XML注入尝试"
}

deny[msg] {
    contains(input.prompt, "<script>")
    msg := "检测到脚本注入"
}', NOW()),

('rule-toxic-01', '有害内容检测', '检测有害、非法内容生成请求', 'prompt_check', 'critical', 'block', NULL, NULL, NULL, NULL, 'package toxic

default allow = true

deny[msg] {
    keywords := ["制作炸弹", "制毒", "黑客攻击"]
    some keyword
    contains(lower(input.prompt), keywords[keyword])
    msg := "检测到有害内容请求"
}', NOW()),

('rule-pii-detect-01', 'PII 敏感信息检测', '检测并保护个人身份信息', 'output_filter', 'medium', 'mask', NULL, NULL, NULL, NULL, 'package pii

default allow = true

sensitive_patterns := [
    `\b\d{3}-\d{4}-\d{4}\b`,
    `\b\d{11}\b`,
    `[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`
]

deny[msg] {
    some pattern
    regex.match(sensitive_patterns[pattern], input.output)
    msg := sprintf("检测到敏感信息模式: %v", [pattern])
}', NOW()),

-- 厂商特定规则
('rule-openai-01', 'OpenAI 使用策略', 'OpenAI API 使用合规检查', 'prompt_check', 'medium', 'warn', NULL, NULL, 'OpenAI', 'GPT-4', 'package openai_policy

default allow = true

warn[msg] {
    count(input.prompt) > 32000
    msg := "提示词超过GPT-4推荐长度"
}', NOW()),

('rule-anthropic-01', 'Anthropic 安全策略', 'Claude 模型安全使用检查', 'prompt_check', 'medium', 'warn', NULL, NULL, 'Anthropic', 'Claude', 'package anthropic_policy

default allow = true

warn[msg] {
    contains(input.prompt, "Human:")
    contains(input.prompt, "Assistant:")
    msg := "检测到可能的对话格式注入"
}', NOW())

ON CONFLICT (id) DO NOTHING;

-- ============================================
-- 告警规则 (Alert Rules) - 默认告警配置
-- ============================================
INSERT INTO alert_rules (id, name, description, event_types, severity_threshold, threshold_count, threshold_window_sec, notify_channels, cooldown_sec, enabled, created_at) VALUES
('alert-critical-block', '严重阻断告警', '当发生严重等级的内容阻断时立即告警', ARRAY['prompt_blocked', 'output_blocked'], 'critical', 1, 60, ARRAY['wecom', 'email'], 300, true, NOW()),
('alert-high-frequency', '高频攻击告警', '5分钟内发生10次以上阻断触发告警', ARRAY['prompt_blocked'], 'high', 10, 300, ARRAY['wecom'], 600, true, NOW()),
('alert-pii-leak', 'PII泄露风险告警', '检测到PII信息泄露风险时告警', ARRAY['pii_detected', 'sensitive_data_masked'], 'medium', 3, 600, ARRAY['email'], 1800, true, NOW()),
('alert-jailbreak', '越狱攻击告警', '检测到越狱攻击尝试时告警', ARRAY['jailbreak_detected'], 'critical', 1, 60, ARRAY['wecom', 'dingtalk'], 300, true, NOW())
ON CONFLICT (id) DO NOTHING;

-- ============================================
-- 工具能力 (Capabilities) - Agent工具权限
-- ============================================
INSERT INTO capabilities (id, name, description, tags, created_at) VALUES
('cap-web-browse', 'web-browsing', '允许Agent访问互联网搜索信息', ARRAY['network', 'search'], NOW()),
('cap-code-exec', 'code-execution', '允许Agent执行代码片段', ARRAY['compute', 'dangerous'], NOW()),
('cap-file-read', 'file-read', '允许Agent读取本地文件', ARRAY['filesystem', 'data'], NOW()),
('cap-file-write', 'file-write', '允许Agent写入本地文件', ARRAY['filesystem', 'dangerous'], NOW()),
('cap-api-call', 'api-call', '允许Agent调用外部API', ARRAY['network', 'integration'], NOW()),
('cap-db-query', 'database-query', '允许Agent查询数据库', ARRAY['data', 'sensitive'], NOW()),
('cap-email-send', 'email-send', '允许Agent发送邮件', ARRAY['communication'], NOW()),
('cap-calendar', 'calendar-access', '允许Agent访问日历', ARRAY['pim', 'data'], NOW())
ON CONFLICT (id) DO NOTHING;

-- ============================================
-- 演示租户和默认策略
-- ============================================
-- 注意: 这部分通常由应用自动创建，这里仅作参考
-- INSERT INTO tenants ... 
-- INSERT INTO policies ...
