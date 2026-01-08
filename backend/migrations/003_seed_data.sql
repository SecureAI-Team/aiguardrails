-- AI GuardRails 默认数据初始化
-- 用于填充初始演示数据

-- ============================================
-- 工具能力 (Capabilities) - Agent工具权限
-- 注意: capabilities.id 是 UUID 类型
-- ============================================
INSERT INTO capabilities (id, name, description, tags) VALUES
(gen_random_uuid(), 'web-browsing', '允许Agent访问互联网搜索信息', '["network", "search"]'::jsonb),
(gen_random_uuid(), 'code-execution', '允许Agent执行代码片段', '["compute", "dangerous"]'::jsonb),
(gen_random_uuid(), 'file-read', '允许Agent读取本地文件', '["filesystem", "data"]'::jsonb),
(gen_random_uuid(), 'file-write', '允许Agent写入本地文件', '["filesystem", "dangerous"]'::jsonb),
(gen_random_uuid(), 'api-call', '允许Agent调用外部API', '["network", "integration"]'::jsonb),
(gen_random_uuid(), 'database-query', '允许Agent查询数据库', '["data", "sensitive"]'::jsonb),
(gen_random_uuid(), 'email-send', '允许Agent发送邮件', '["communication"]'::jsonb),
(gen_random_uuid(), 'calendar-access', '允许Agent访问日历', '["pim", "data"]'::jsonb)
ON CONFLICT (name) DO NOTHING;

-- ============================================
-- 规则模板 (Rule Templates) - 预定义的规则模板
-- 注意: 这些模板已在 0006_tenant_rules.sql 中创建
-- 这里添加更多有用的模板
-- ============================================
INSERT INTO rule_templates (id, name, rule_type, description, config_schema, default_config, tags) VALUES
(gen_random_uuid(), 'pii_detection', 'business', 'PII 敏感信息检测 - 检测并保护个人隐私数据',
 '{"type":"object","properties":{"patterns":{"type":"array"},"mask_type":{"type":"string"}}}',
 '{"enabled":true,"patterns":["phone","id_card","email","bank_card"],"mask_type":"partial","responses":{"detected":"检测到敏感信息，已进行脱敏处理"}}',
 ARRAY['pii', 'privacy', 'security']),

(gen_random_uuid(), 'jailbreak_detection', 'business', '越狱攻击检测 - 检测常见的LLM越狱攻击模式',
 '{"type":"object","properties":{"patterns":{"type":"array"},"action":{"type":"string"}}}',
 '{"enabled":true,"patterns":["ignore previous","forget everything","you are now","pretend to be"],"action":"block","responses":{"blocked":"检测到安全威胁，请求已被阻止"}}',
 ARRAY['security', 'jailbreak', 'attack']),

(gen_random_uuid(), 'injection_detection', 'business', '提示注入检测 - 检测提示注入攻击',
 '{"type":"object","properties":{"patterns":{"type":"array"},"action":{"type":"string"}}}',
 '{"enabled":true,"patterns":["]]>","<script>","${","{{"],"action":"block","responses":{"blocked":"检测到注入攻击，请求已被阻止"}}',
 ARRAY['security', 'injection', 'attack']),

(gen_random_uuid(), 'gdpr_compliance', 'business', 'GDPR 欧盟数据保护合规 - 确保符合GDPR要求',
 '{"type":"object","properties":{"regions":{"type":"array"},"data_types":{"type":"array"}}}',
 '{"enabled":true,"regions":["EU"],"data_types":["name","email","address","phone"],"action":"mask","responses":{"violation":"此内容可能违反GDPR规定"}}',
 ARRAY['compliance', 'gdpr', 'eu', 'privacy']),

(gen_random_uuid(), 'pipl_compliance', 'business', 'PIPL 中国个人信息保护法合规',
 '{"type":"object","properties":{"data_types":{"type":"array"}}}',
 '{"enabled":true,"data_types":["id_card","phone","bank_card","address"],"action":"mask","responses":{"violation":"此内容可能违反《个人信息保护法》"}}',
 ARRAY['compliance', 'pipl', 'china', 'privacy']),

(gen_random_uuid(), 'toxic_content', 'business', '有害内容检测 - 检测暴力、色情等有害内容',
 '{"type":"object","properties":{"categories":{"type":"array"}}}',
 '{"enabled":true,"categories":["violence","adult","illegal","hate_speech"],"action":"block","responses":{"blocked":"检测到有害内容，请求已被阻止"}}',
 ARRAY['safety', 'toxic', 'content']),

(gen_random_uuid(), 'rate_limiting', 'permission', '请求速率限制 - 限制API调用频率',
 '{"type":"object","properties":{"max_rpm":{"type":"integer"},"max_tpm":{"type":"integer"}}}',
 '{"enabled":true,"max_rpm":60,"max_tpm":100000,"action":"throttle","responses":{"exceeded":"请求频率过高，请稍后再试"}}',
 ARRAY['rate', 'limit', 'quota'])

ON CONFLICT (name) DO NOTHING;

-- 注意: alert_rules 已在 0009_alert_system.sql 中有默认值
-- 如需添加更多告警规则，取消下面的注释

-- INSERT INTO alert_rules (tenant_id, name, description, event_types, severity_threshold, threshold_count, threshold_window_sec, notify_channels, priority) VALUES
-- (NULL, 'high_token_usage', '高Token消耗告警', ARRAY['high_usage'], 'medium', 1, 3600, ARRAY['webhook'], 30)
-- ON CONFLICT DO NOTHING;
