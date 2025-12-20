package guardrails

# Agent tool allowlist enforcement
deny_reason[msg] {
  input.mode == "agent_tool"
  input.allowed_tools != null
  tool := lower(input.tool)
  allowed := [lower(t) | t := input.allowed_tools[_]]
  not tool_allowed(tool, allowed)
  msg := {"allow": false, "reason": "opa_agent_tool_not_allowed", "signals": [tool]}
}

tool_allowed(tool, allowed) {
  allowed[_] == tool
}

# Agent tool category guard (secrets/admin/finance/privileged)
deny_reason[msg] {
  input.mode == "agent_tool"
  input.tool_tags != null
  sensitive := ["secrets", "admin", "finance", "privileged"]
  tag := lower(input.tool_tags[_])
  sensitive[_] == tag
  msg := {"allow": false, "reason": "opa_agent_tool_sensitive", "signals": [tag]}
}

# Agent step budget enforcement
deny_reason[msg] {
  input.mode == "agent_tool"
  input.step != null
  input.max_steps != null
  input.step > input.max_steps
  msg := {"allow": false, "reason": "opa_agent_step_budget", "signals": [sprintf("%d/%d", [input.step, input.max_steps])]}
}

# Agent cross-tenant tool use blocking
deny_reason[msg] {
  input.mode == "agent_tool"
  input.tenantId != null
  input.tool_tenant != null
  lower(input.tool_tenant) != lower(input.tenantId)
  msg := {"allow": false, "reason": "opa_agent_cross_tenant_tool", "signals": [input.tool_tenant, input.tenantId]}
}

# MCP provider allowlist
deny_reason[msg] {
  input.mode == "mcp_call"
  input.allowed_providers != null
  provider := lower(input.provider)
  allowed := [lower(p) | p := input.allowed_providers[_]]
  not provider_allowed(provider, allowed)
  msg := {"allow": false, "reason": "opa_mcp_provider_not_allowed", "signals": [provider]}
}

provider_allowed(provider, allowed) {
  allowed[_] == provider
}

# MCP capability tag denylist
deny_reason[msg] {
  input.mode == "mcp_call"
  input.deny_tags != null
  input.capability_tags != null
  deny := lower(input.deny_tags[_])
  tag := lower(input.capability_tags[_])
  deny == tag
  msg := {"allow": false, "reason": "opa_mcp_capability_denied", "signals": [tag]}
}

# MCP requires sandbox for code execution
deny_reason[msg] {
  input.mode == "mcp_call"
  input.capability_tags != null
  tag := lower(input.capability_tags[_])
  tag == "code_exec"
  not input.sandboxed
  msg := {"allow": false, "reason": "opa_mcp_sandbox_required", "signals": ["code_exec"]}
}

# MCP requires signed tool calls
deny_reason[msg] {
  input.mode == "mcp_call"
  input.require_sign
  not input.signed
  msg := {"allow": false, "reason": "opa_mcp_unsigned_call", "signals": ["signature_missing"]}
}

