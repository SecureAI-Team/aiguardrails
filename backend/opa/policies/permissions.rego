package guardrails

# ===== 角色定义 =====

roles = {
  "platform_admin": {
    "level": 4,
    "name": "平台管理员",
    "permissions": ["*"]
  },
  "tenant_admin": {
    "level": 3,
    "name": "租户管理员",
    "permissions": ["manage_policy", "manage_apps", "use_tools", "chat", "view_audit", "industrial_control"]
  },
  "operator": {
    "level": 2,
    "name": "操作员",
    "permissions": ["use_tools", "chat", "read_status"]
  },
  "viewer": {
    "level": 1,
    "name": "只读用户",
    "permissions": ["chat", "read_status"]
  }
}

# ===== 工具权限等级 =====

tool_permissions = {
  # 只读操作
  "read_status": {"min_level": 1, "requires_confirmation": false},
  "diagnostic": {"min_level": 2, "requires_confirmation": false},
  "upload_logs": {"min_level": 2, "requires_confirmation": false},
  
  # 写入操作
  "set_parameter": {"min_level": 2, "requires_confirmation": true},
  "configure": {"min_level": 2, "requires_confirmation": true},
  
  # 控制操作
  "start_stop": {"min_level": 3, "requires_confirmation": true},
  "change_program": {"min_level": 3, "requires_confirmation": true},
  
  # 危险操作
  "firmware_update": {"min_level": 4, "requires_confirmation": true},
  "safety_override": {"min_level": 4, "requires_confirmation": true, "requires_mfa": true},
  "reset_factory": {"min_level": 4, "requires_confirmation": true, "requires_mfa": true}
}

# ===== 权限检查函数 =====

# 检查角色是否有特定权限
has_permission(role, perm) {
  roles[role].permissions[_] == perm
}

has_permission(role, perm) {
  roles[role].permissions[_] == "*"
}

# 获取角色等级
get_role_level(role) = level {
  level := roles[role].level
}

get_role_level(role) = 0 {
  not roles[role]
}

# ===== 权限检查规则 =====

# 通用权限检查
deny_reason[msg] {
  input.mode == "permission_check"
  role := input.role
  permission := input.required_permission
  not has_permission(role, permission)
  msg := {
    "allow": false,
    "reason": "permission_denied",
    "signals": [role, permission, "category:permission"],
    "response": sprintf("您的角色 [%s] 没有 [%s] 权限", [roles[role].name, permission])
  }
}

# 工具权限检查
deny_reason[msg] {
  input.mode == "tool_permission"
  role := input.role
  tool := input.tool
  tool_perm := tool_permissions[tool]
  role_level := get_role_level(role)
  role_level < tool_perm.min_level
  msg := {
    "allow": false,
    "reason": "tool_permission_denied",
    "signals": [role, tool, sprintf("requires_level_%d", [tool_perm.min_level])],
    "response": sprintf("您的权限级别不足以执行 [%s] 操作", [tool])
  }
}

# 需要确认的操作
deny_reason[msg] {
  input.mode == "tool_permission"
  tool := input.tool
  tool_perm := tool_permissions[tool]
  tool_perm.requires_confirmation == true
  not input.confirmed
  msg := {
    "allow": false,
    "reason": "confirmation_required",
    "signals": [tool, "requires_confirmation"],
    "decision": "confirm",
    "response": sprintf("执行 [%s] 操作需要您的确认", [tool])
  }
}

# 需要MFA的操作
deny_reason[msg] {
  input.mode == "tool_permission"
  tool := input.tool
  tool_perm := tool_permissions[tool]
  tool_perm.requires_mfa == true
  not input.mfa_verified
  msg := {
    "allow": false,
    "reason": "mfa_required",
    "signals": [tool, "requires_mfa"],
    "decision": "mfa",
    "response": sprintf("执行 [%s] 操作需要多因素认证", [tool])
  }
}

# ===== 工业控制特殊权限 =====

# 工业控制命令需要tenant_admin及以上
deny_reason[msg] {
  input.mode == "industrial_command"
  role := input.role
  role_level := get_role_level(role)
  role_level < 3
  msg := {
    "allow": false,
    "reason": "industrial_command_requires_admin",
    "signals": [role, "requires_tenant_admin"],
    "response": "工业控制命令需要管理员权限"
  }
}

# 安全覆盖仅限platform_admin
deny_reason[msg] {
  input.mode == "safety_override"
  role := input.role
  role != "platform_admin"
  msg := {
    "allow": false,
    "reason": "safety_override_requires_platform_admin",
    "signals": [role, "requires_platform_admin"],
    "response": "安全覆盖操作仅限平台管理员"
  }
}

# 工业控制二次确认
deny_reason[msg] {
  input.mode == "industrial_command"
  role := input.role
  role_level := get_role_level(role)
  role_level >= 3
  not input.confirmed
  msg := {
    "allow": false,
    "reason": "industrial_confirmation_required",
    "signals": ["industrial_command", "requires_confirmation"],
    "decision": "confirm",
    "response": "工业控制命令需要二次确认，请确认您要执行此操作"
  }
}

# ===== 租户隔离 =====

# 跨租户操作检查
deny_reason[msg] {
  input.mode == "tenant_check"
  input.target_tenant != null
  input.user_tenant != null
  input.target_tenant != input.user_tenant
  role := input.role
  role != "platform_admin"
  msg := {
    "allow": false,
    "reason": "cross_tenant_denied",
    "signals": [input.user_tenant, input.target_tenant],
    "response": "您无法访问其他租户的资源"
  }
}

# ===== 时间窗口控制 =====

# 工作时间检查（可选）
deny_reason[msg] {
  input.mode == "industrial_command"
  input.enforce_work_hours == true
  input.current_hour != null
  input.current_hour < 8
  msg := {
    "allow": false,
    "reason": "outside_work_hours",
    "signals": ["work_hours", sprintf("current:%d", [input.current_hour])],
    "response": "工业控制命令仅在工作时间(8:00-20:00)执行"
  }
}

deny_reason[msg] {
  input.mode == "industrial_command"
  input.enforce_work_hours == true
  input.current_hour != null
  input.current_hour > 20
  msg := {
    "allow": false,
    "reason": "outside_work_hours",
    "signals": ["work_hours", sprintf("current:%d", [input.current_hour])],
    "response": "工业控制命令仅在工作时间(8:00-20:00)执行"
  }
}
