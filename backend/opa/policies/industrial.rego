package guardrails

# Block industrial vendor/tool sensitive prompts
banned_vendors = ["siemens", "abb", "ge"]

# Action buckets to grade risk
read_actions = ["read status", "diagnostic", "monitor", "upload logs"]
write_actions = ["write", "setpoint", "configure", "change parameter", "change speed", "change ip", "change ladder", "change program", "drive parameter"]
stop_actions = ["shutdown", "emergency stop", "stop plc", "halt plc"]
override_actions = ["override", "motion override", "override e-stop", "bypass interlock", "disable safety", "reset safety", "safety relay", "safety interlock"]
firmware_actions = ["upload firmware", "download firmware", "firmware"]

# PLC-specific dangerous commands (covered by buckets but explicit list)
plc_actions = [
  "force coil",
  "write coil",
  "force output",
  "write output",
  "bypass interlock",
  "disable safety",
  "reset safety",
  "upload firmware",
  "download firmware",
  "change program",
  "change ladder",
  "change speed",
  "change ip",
  "stop plc",
  "start plc",
  "halt plc",
  "override e-stop"
]

simulation_mode {
  input.simulation == true
}
simulation_mode {
  input.mode != null
  lower(input.mode) == "simulation"
}

deny_reason[msg] {
  lp := lower(input.prompt)
  v := banned_vendors[_]
  act := stop_actions[_]
  contains(lp, v)
  contains(lp, act)
  not simulation_mode
  msg := {"allow": false, "reason": "opa_industrial_prompt_stop", "signals": [v, act, "severity:critical", "decision:block"]}
}

deny_reason[msg] {
  lp := lower(input.prompt)
  v := banned_vendors[_]
  act := override_actions[_]
  contains(lp, v)
  contains(lp, act)
  not simulation_mode
  msg := {"allow": false, "reason": "opa_industrial_prompt_override", "signals": [v, act, "severity:critical", "decision:block"]}
}

deny_reason[msg] {
  lp := lower(input.prompt)
  v := banned_vendors[_]
  act := firmware_actions[_]
  contains(lp, v)
  contains(lp, act)
  not simulation_mode
  msg := {"allow": false, "reason": "opa_industrial_prompt_firmware", "signals": [v, act, "severity:high", "decision:block"]}
}

deny_reason[msg] {
  lp := lower(input.prompt)
  v := banned_vendors[_]
  act := write_actions[_]
  contains(lp, v)
  contains(lp, act)
  not simulation_mode
  msg := {"allow": false, "reason": "opa_industrial_prompt_write", "signals": [v, act, "severity:high", "decision:block"]}
}

deny_reason[msg] {
  lp := lower(input.prompt)
  v := banned_vendors[_]
  act := read_actions[_]
  contains(lp, v)
  contains(lp, act)
  not simulation_mode
  msg := {"allow": false, "reason": "opa_industrial_prompt_read", "signals": [v, act, "severity:medium", "decision:mark"]}
}

deny_reason[msg] {
  lo := lower(input.output)
  v := banned_vendors[_]
  act := stop_actions[_]
  contains(lo, v)
  contains(lo, act)
  not simulation_mode
  msg := {"allow": false, "reason": "opa_industrial_output_stop", "signals": [v, act, "severity:critical", "decision:block"]}
}

deny_reason[msg] {
  lo := lower(input.output)
  v := banned_vendors[_]
  act := override_actions[_]
  contains(lo, v)
  contains(lo, act)
  not simulation_mode
  msg := {"allow": false, "reason": "opa_industrial_output_override", "signals": [v, act, "severity:critical", "decision:block"]}
}

deny_reason[msg] {
  lo := lower(input.output)
  v := banned_vendors[_]
  act := firmware_actions[_]
  contains(lo, v)
  contains(lo, act)
  not simulation_mode
  msg := {"allow": false, "reason": "opa_industrial_output_firmware", "signals": [v, act, "severity:high", "decision:block"]}
}

deny_reason[msg] {
  lo := lower(input.output)
  v := banned_vendors[_]
  act := write_actions[_]
  contains(lo, v)
  contains(lo, act)
  not simulation_mode
  msg := {"allow": false, "reason": "opa_industrial_output_write", "signals": [v, act, "severity:high", "decision:block"]}
}

deny_reason[msg] {
  lo := lower(input.output)
  v := banned_vendors[_]
  act := read_actions[_]
  contains(lo, v)
  contains(lo, act)
  not simulation_mode
  msg := {"allow": false, "reason": "opa_industrial_output_read", "signals": [v, act, "severity:medium", "decision:mark"]}
}

# Block PLC dangerous commands even without explicit vendor mention
deny_reason[msg] {
  lp := lower(input.prompt)
  plc := plc_actions[_]
  contains(lp, plc)
  not simulation_mode
  msg := {"allow": false, "reason": "opa_plc_prompt", "signals": [plc, "severity:critical", "decision:block"]}
}

deny_reason[msg] {
  lo := lower(input.output)
  plc := plc_actions[_]
  contains(lo, plc)
  not simulation_mode
  msg := {"allow": false, "reason": "opa_plc_output", "signals": [plc, "severity:critical", "decision:block"]}
}

