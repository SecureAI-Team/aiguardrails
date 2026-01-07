package guardrails

# ===== 中国网络安全法 =====
cn_cybersec_terms = ["跨境", "cross-border", "境外存储", "data export", "数据出境", "境外传输", "出境审批"]

deny_reason[msg] {
  lp := lower(input.prompt)
  term := cn_cybersec_terms[_]
  contains(lp, term)
  msg := {"allow": false, "reason": "cn_cybersec_data_export", "signals": [term, "regulation:cybersecurity_law", "severity:high", "decision:block"]}
}

deny_reason[msg] {
  lo := lower(input.output)
  term := cn_cybersec_terms[_]
  contains(lo, term)
  msg := {"allow": false, "reason": "cn_cybersec_output_block", "signals": [term, "regulation:cybersecurity_law", "severity:high"]}
}

# ===== 中国个人信息保护法 (PIPL) =====
cn_personal_info_terms = [
  "身份证", "身份证号", "idcard",
  "手机号", "电话号码", "phone number",
  "银行卡", "银行账号", "bank account",
  "人脸", "面部识别", "facial recognition",
  "指纹", "生物特征", "biometric",
  "住址", "家庭地址", "home address",
  "户籍", "籍贯"
]

deny_reason[msg] {
  lp := lower(input.prompt)
  term := cn_personal_info_terms[_]
  contains(lp, term)
  msg := {"allow": false, "reason": "cn_pipl_personal_info", "signals": [term, "regulation:pipl", "severity:high", "decision:block"]}
}

deny_reason[msg] {
  lo := lower(input.output)
  term := cn_personal_info_terms[_]
  contains(lo, term)
  msg := {"allow": false, "reason": "cn_pipl_output_block", "signals": [term, "regulation:pipl", "severity:high", "decision:redact"]}
}

# ===== 中国数据安全法 =====
cn_data_security_terms = ["核心数据", "重要数据", "数据分级", "数据安全评估", "关键信息基础设施"]

deny_reason[msg] {
  lp := lower(input.prompt)
  term := cn_data_security_terms[_]
  contains(lp, term)
  msg := {"allow": false, "reason": "cn_data_security_law", "signals": [term, "regulation:data_security_law", "severity:high", "decision:mark"]}
}

# ===== 中国算法推荐规定 =====
cn_algorithm_terms = ["算法推荐", "个性化推送", "信息茧房", "大数据杀熟", "差别定价"]

deny_reason[msg] {
  lp := lower(input.prompt)
  term := cn_algorithm_terms[_]
  contains(lp, term)
  msg := {"allow": false, "reason": "cn_algorithm_regulation", "signals": [term, "regulation:algorithm_recommendation", "severity:medium", "decision:mark"]}
}

# ===== 中国生成式AI管理规定 =====
cn_genai_blocked_content = [
  "颠覆国家政权", "分裂国家",
  "恐怖主义", "极端主义",
  "虚假信息", "谣言",
  "歧视", "侮辱"
]

deny_reason[msg] {
  lp := lower(input.prompt)
  content := cn_genai_blocked_content[_]
  contains(lp, content)
  msg := {"allow": false, "reason": "cn_genai_blocked_content", "signals": [content, "regulation:genai_regulation", "severity:critical", "decision:block"]}
}

deny_reason[msg] {
  lo := lower(input.output)
  content := cn_genai_blocked_content[_]
  contains(lo, content)
  msg := {"allow": false, "reason": "cn_genai_output_block", "signals": [content, "regulation:genai_regulation", "severity:critical", "decision:block"]}
}
