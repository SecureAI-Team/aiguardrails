package guardrails

# ===== 厂商白名单/黑名单 =====

# 允许的厂商（可通过input.allowed_vendors覆盖）
default_allowed_vendors = ["siemens", "西门子"]

# 竞品厂商黑名单
blocked_vendors = [
  "abb",
  "rockwell", "allen-bradley", "ab plc",
  "schneider", "施耐德", "modicon",
  "mitsubishi", "三菱", "melsec",
  "omron", "欧姆龙",
  "ge", "general electric",
  "honeywell", "霍尼韦尔",
  "beckhoff", "倍福",
  "delta", "台达",
  "keyence", "基恩士",
  "fanuc", "发那科"
]

# 竞品产品型号黑名单
blocked_products = [
  # Rockwell/AB
  "compactlogix", "controllogix", "micrologix", "slc500",
  # ABB
  "ac800", "ac500", "pm5", "freelance",
  # Schneider
  "m340", "m580", "quantum", "unity pro",
  # Mitsubishi
  "fx5u", "fx3u", "q series", "iq-r", "iq-f", "gx works",
  # Omron
  "cp1", "cj2", "nx1", "nj", "sysmac",
  # Other
  "twincat", "codesys"
]

# 允许的西门子产品关键词
siemens_products = [
  "s7-1500", "s7-1200", "s7-300", "s7-400", "s7-200",
  "tia portal", "step 7", "wincc", "simatic",
  "sinumerik", "sinamics", "simotion",
  "logo!", "et 200", "profinet", "profibus",
  "scalance", "sitop", "sirius"
]

# ===== 厂商限制规则 =====

# 检查提示词是否包含竞品厂商
deny_reason[msg] {
  input.mode == "vendor_check"
  lp := lower(input.prompt)
  vendor := blocked_vendors[_]
  contains(lp, vendor)
  msg := {
    "allow": false, 
    "reason": "competitor_vendor_blocked", 
    "signals": [vendor, "category:vendor_restriction", "decision:block"],
    "response": "抱歉，我是西门子产品专属助手，无法回答其他厂商产品的问题。如需西门子产品帮助，请告诉我。"
  }
}

# 检查提示词是否包含竞品产品
deny_reason[msg] {
  input.mode == "vendor_check"
  lp := lower(input.prompt)
  product := blocked_products[_]
  contains(lp, product)
  msg := {
    "allow": false, 
    "reason": "competitor_product_blocked", 
    "signals": [product, "category:vendor_restriction", "decision:block"],
    "response": "抱歉，我无法提供该产品的信息。我可以帮助您了解西门子的同类产品。"
  }
}

# 检查输出是否包含竞品信息（需要脱敏）
deny_reason[msg] {
  input.mode == "vendor_check"
  lo := lower(input.output)
  vendor := blocked_vendors[_]
  contains(lo, vendor)
  msg := {
    "allow": false, 
    "reason": "competitor_in_output", 
    "signals": [vendor, "category:vendor_restriction", "decision:redact"]
  }
}

deny_reason[msg] {
  input.mode == "vendor_check"
  lo := lower(input.output)
  product := blocked_products[_]
  contains(lo, product)
  msg := {
    "allow": false, 
    "reason": "competitor_product_in_output", 
    "signals": [product, "category:vendor_restriction", "decision:redact"]
  }
}

# ===== 西门子产品识别 =====

# 检查是否为西门子相关问题
is_siemens_question {
  lp := lower(input.prompt)
  product := siemens_products[_]
  contains(lp, product)
}

is_siemens_question {
  lp := lower(input.prompt)
  vendor := default_allowed_vendors[_]
  contains(lp, vendor)
}

# 通用问题（不涉及任何厂商）也允许
is_generic_question {
  lp := lower(input.prompt)
  not contains_any_vendor(lp)
}

contains_any_vendor(text) {
  vendor := blocked_vendors[_]
  contains(text, vendor)
}

contains_any_vendor(text) {
  vendor := default_allowed_vendors[_]
  contains(text, vendor)
}

# ===== 产品推荐引导 =====

# 当用户询问竞品时，可以推荐西门子替代品
product_alternatives = {
  "compactlogix": "S7-1500 Compact CPU",
  "controllogix": "S7-1500",
  "m340": "S7-1200",
  "m580": "S7-1500",
  "fx5u": "S7-1200",
  "q series": "S7-1500",
  "cp1": "S7-1200",
  "cj2": "S7-1500"
}

# ===== 领域边界控制 =====

# 允许的领域
allowed_domains = ["plc", "scada", "hmi", "motion", "drive", "automation", "industrial", "工业自动化", "plc编程", "控制系统"]

# 超出领域的话题
out_of_scope_topics = [
  "股票", "投资理财", "基金",
  "医疗诊断", "处方", "用药建议",
  "法律咨询", "诉讼",
  "政治", "选举"
]

deny_reason[msg] {
  input.mode == "domain_check"
  lp := lower(input.prompt)
  topic := out_of_scope_topics[_]
  contains(lp, topic)
  msg := {
    "allow": false, 
    "reason": "out_of_scope_topic", 
    "signals": [topic, "category:domain_boundary", "decision:deflect"],
    "response": "这个问题超出了我的专业范围。我专注于工业自动化和西门子产品技术支持。"
  }
}
