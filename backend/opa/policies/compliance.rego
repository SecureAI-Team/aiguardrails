package guardrails

# Deny prompts mentioning GDPR clause request for personal data.
deny_reason[msg] {
  lp := lower(input.prompt)
  contains(lp, "gdpr")
  contains(lp, "personal data")
  msg := {"allow": false, "reason": "opa_gdpr_block", "signals": ["gdpr_personal_data"]}
}

# Deny outputs that leak pii markers.
deny_reason[msg] {
  lo := lower(input.output)
  contains(lo, "pii")
  msg := {"allow": false, "reason": "opa_gdpr_output_block", "signals": ["pii_marker"]}
}

# Broad PII keywords
pii_terms = ["ssn", "passport", "bank account", "iban", "swift", "medical record", "personal data", "id number"]

deny_reason[msg] {
  lp := lower(input.prompt)
  term := pii_terms[_]
  contains(lp, term)
  msg := {"allow": false, "reason": "opa_pii_prompt", "signals": [term]}
}

deny_reason[msg] {
  lo := lower(input.output)
  term := pii_terms[_]
  contains(lo, term)
  msg := {"allow": false, "reason": "opa_pii_output", "signals": [term]}
}

# HIPAA / PHI markers
hipaa_terms = ["hipaa", "phi", "patient", "mrn", "medical record", "diagnosis", "treatment plan"]

deny_reason[msg] {
  lp := lower(input.prompt)
  term := hipaa_terms[_]
  contains(lp, term)
  msg := {"allow": false, "reason": "opa_hipaa_prompt", "signals": [term, "decision:block"]}
}

deny_reason[msg] {
  lo := lower(input.output)
  term := hipaa_terms[_]
  contains(lo, term)
  msg := {"allow": false, "reason": "opa_hipaa_output", "signals": [term, "decision:block"]}
}

# PCI-DSS markers
pci_terms = ["card number", "credit card", "cvv", "cvc", "pan", "cardholder"]

deny_reason[msg] {
  lp := lower(input.prompt)
  term := pci_terms[_]
  contains(lp, term)
  msg := {"allow": false, "reason": "opa_pci_prompt", "signals": [term, "decision:block"]}
}

deny_reason[msg] {
  lo := lower(input.output)
  term := pci_terms[_]
  contains(lo, term)
  msg := {"allow": false, "reason": "opa_pci_output", "signals": [term, "decision:block"]}
}

# Data sovereignty markers -> prefer mark/route review
sovereignty_terms = ["cross-border", "cross border", "store in us", "store in eu", "china data", "data residency", "data sovereignty"]

deny_reason[msg] {
  lp := lower(input.prompt)
  term := sovereignty_terms[_]
  contains(lp, term)
  msg := {"allow": false, "reason": "opa_sovereignty_prompt", "signals": [term, "decision:mark"]}
}

deny_reason[msg] {
  lo := lower(input.output)
  term := sovereignty_terms[_]
  contains(lo, term)
  msg := {"allow": false, "reason": "opa_sovereignty_output", "signals": [term, "decision:mark"]}
}

# EU AI Act high-risk mention + safety critical override
deny_reason[msg] {
  lp := lower(input.prompt)
  contains(lp, "high-risk")
  contains(lp, "override")
  msg := {"allow": false, "reason": "opa_ai_act_override", "signals": ["ai_act_override"]}
}

