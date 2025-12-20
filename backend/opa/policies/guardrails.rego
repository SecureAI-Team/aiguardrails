package guardrails

default allow = true

# Deny if prompt or output contains banned keywords
banned = ["password", "secret", "admin password", "ssh key"]

deny_reason[msg] {
  lp := lower(input.prompt)
  ban := banned[_]
  contains(lp, ban)
  msg := {"allow": false, "reason": "opa_block_prompt", "signals": [ban]}
}

deny_reason[msg] {
  lo := lower(input.output)
  ban := banned[_]
  contains(lo, ban)
  msg := {"allow": false, "reason": "opa_block_output", "signals": [ban]}
}

allow = false {
  r := deny_reason[_]
  r.allow == false
}

