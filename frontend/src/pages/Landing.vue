<template>
  <div class="landing">
    <!-- Hero Section -->
    <header class="hero">
      <nav class="nav">
        <div class="logo">
          <span class="logo-icon">🛡️</span>
          <span class="logo-text">AI GuardRails</span>
        </div>
        <div class="nav-links">
          <router-link to="/login" class="btn-outline">登录</router-link>
          <router-link to="/login" class="btn-primary">免费试用</router-link>
        </div>
      </nav>
      <div class="hero-content">
        <h1>工业AI应用<span class="highlight">安全护栏</span>平台</h1>
        <p class="subtitle">为您的AI应用提供全方位安全保护，确保合规、可控、可信</p>
        <div class="hero-actions">
          <router-link to="/login" class="btn-primary btn-lg">立即体验 →</router-link>
          <a href="#features" class="btn-outline btn-lg">了解更多</a>
        </div>
        <div class="hero-stats">
          <div class="stat"><span class="stat-num">99.9%</span><span class="stat-label">安全拦截率</span></div>
          <div class="stat"><span class="stat-num">50ms</span><span class="stat-label">平均延迟</span></div>
          <div class="stat"><span class="stat-num">100+</span><span class="stat-label">合规规则</span></div>
        </div>
      </div>
    </header>

    <!-- Features Section -->
    <section id="features" class="features">
      <h2>⚡ 核心能力</h2>
      <div class="feature-grid">
        <div class="feature-card">
          <div class="feature-icon">🔥</div>
          <h3>提示词防火墙</h3>
          <p>拦截恶意提示词注入攻击，保护LLM免受越狱攻击</p>
        </div>
        <div class="feature-card">
          <div class="feature-icon">🎯</div>
          <h3>输出内容过滤</h3>
          <p>敏感信息脱敏、PII检测、合规内容审核</p>
        </div>
        <div class="feature-card">
          <div class="feature-icon">🔐</div>
          <h3>工具访问控制</h3>
          <p>Agent工具执行沙箱隔离，权限分级管理</p>
        </div>
        <div class="feature-card">
          <div class="feature-icon">📋</div>
          <h3>法规合规引擎</h3>
          <p>GDPR/PIPL/网络安全法等法规自动检查</p>
        </div>
      </div>
    </section>

    <!-- Advantages Section -->
    <section class="advantages">
      <h2>✨ 平台优势</h2>
      <div class="advantage-list">
        <div class="advantage-item">
          <span class="check">✓</span>
          <div>
            <strong>多租户架构</strong>
            <p>完全隔离的租户空间，支持私有化部署</p>
          </div>
        </div>
        <div class="advantage-item">
          <span class="check">✓</span>
          <div>
            <strong>OPA策略引擎</strong>
            <p>灵活的Rego策略语言，热加载无需重启</p>
          </div>
        </div>
        <div class="advantage-item">
          <span class="check">✓</span>
          <div>
            <strong>实时审计日志</strong>
            <p>全链路追踪，满足安全审计要求</p>
          </div>
        </div>
        <div class="advantage-item">
          <span class="check">✓</span>
          <div>
            <strong>中国法规合规</strong>
            <p>内置网络安全法、PIPL、数据安全法规则</p>
          </div>
        </div>
        <div class="advantage-item">
          <span class="check">✓</span>
          <div>
            <strong>厂商规则定制</strong>
            <p>支持西门子等工业厂商专属助手配置</p>
          </div>
        </div>
        <div class="advantage-item">
          <span class="check">✓</span>
          <div>
            <strong>SDK快速集成</strong>
            <p>Go/Python/Node.js SDK，10分钟完成接入</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Quick Start Section -->
    <section class="quickstart">
      <h2>🚀 快速开始</h2>
      <div class="steps">
        <div class="step">
          <div class="step-num">1</div>
          <div class="step-content">
            <h4>注册账号</h4>
            <p>微信/支付宝/手机号快速注册</p>
          </div>
        </div>
        <div class="step-arrow">→</div>
        <div class="step">
          <div class="step-num">2</div>
          <div class="step-content">
            <h4>创建应用</h4>
            <p>获取API Key和Secret</p>
          </div>
        </div>
        <div class="step-arrow">→</div>
        <div class="step">
          <div class="step-num">3</div>
          <div class="step-content">
            <h4>集成SDK</h4>
            <p>添加几行代码即可</p>
          </div>
        </div>
        <div class="step-arrow">→</div>
        <div class="step">
          <div class="step-num">4</div>
          <div class="step-content">
            <h4>开启防护</h4>
            <p>实时监控AI安全</p>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="cta">
      <h2>立即开始保护您的AI应用</h2>
      <p>免费试用30天，无需绑定信用卡</p>
      <router-link to="/login" class="btn-primary btn-lg">免费注册 →</router-link>
    </section>

    <!-- Footer -->
    <footer class="footer">
      <div class="footer-content">
        <div class="footer-brand">
          <span class="logo-icon">🛡️</span> AI GuardRails
          <p>工业AI应用安全护栏平台</p>
        </div>
        <div class="footer-links">
          <router-link to="/sdks">SDK</router-link>
          <router-link to="/playground">API调试</router-link>
          <router-link to="/models">模型</router-link>
          <router-link to="/login">登录</router-link>
        </div>
      </div>
      <div class="footer-bottom">© 2026 AI GuardRails. All rights reserved.</div>
    </footer>
  </div>
</template>

<script setup lang="ts">
</script>

<style scoped>
.landing { font-family: 'Inter', -apple-system, sans-serif; color: #1e293b; }

/* Hero */
.hero { background: linear-gradient(135deg, #0f172a 0%, #1e3a5f 100%); color: white; padding: 0 0 80px; }
.nav { display: flex; justify-content: space-between; align-items: center; padding: 20px 60px; }
.logo { display: flex; align-items: center; gap: 10px; font-size: 1.5rem; font-weight: 700; }
.logo-icon { font-size: 2rem; }
.nav-links { display: flex; gap: 12px; }
.hero-content { text-align: center; padding: 60px 20px; max-width: 800px; margin: 0 auto; }
.hero h1 { font-size: 3.5rem; font-weight: 800; margin-bottom: 20px; line-height: 1.2; }
.highlight { background: linear-gradient(90deg, #3b82f6, #8b5cf6); -webkit-background-clip: text; -webkit-text-fill-color: transparent; }
.subtitle { font-size: 1.25rem; color: #94a3b8; margin-bottom: 40px; }
.hero-actions { display: flex; gap: 16px; justify-content: center; margin-bottom: 60px; }
.hero-stats { display: flex; justify-content: center; gap: 60px; }
.stat { text-align: center; }
.stat-num { display: block; font-size: 2.5rem; font-weight: 700; color: #3b82f6; }
.stat-label { font-size: 0.875rem; color: #94a3b8; }

/* Buttons */
.btn-primary { background: linear-gradient(90deg, #2563eb, #7c3aed); color: white; padding: 12px 24px; border-radius: 8px; text-decoration: none; font-weight: 600; transition: transform 0.2s; }
.btn-primary:hover { transform: translateY(-2px); }
.btn-outline { border: 1px solid rgba(255,255,255,0.3); color: white; padding: 12px 24px; border-radius: 8px; text-decoration: none; }
.btn-lg { padding: 16px 32px; font-size: 1.1rem; }

/* Features */
.features { padding: 80px 60px; background: #f8fafc; }
.features h2 { text-align: center; font-size: 2rem; margin-bottom: 50px; }
.feature-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 30px; max-width: 1200px; margin: 0 auto; }
.feature-card { background: white; padding: 30px; border-radius: 12px; box-shadow: 0 4px 20px rgba(0,0,0,0.05); text-align: center; transition: transform 0.2s; }
.feature-card:hover { transform: translateY(-5px); }
.feature-icon { font-size: 3rem; margin-bottom: 16px; }
.feature-card h3 { font-size: 1.25rem; margin-bottom: 12px; }
.feature-card p { color: #64748b; font-size: 0.95rem; }

/* Advantages */
.advantages { padding: 80px 60px; max-width: 1000px; margin: 0 auto; }
.advantages h2 { text-align: center; font-size: 2rem; margin-bottom: 50px; }
.advantage-list { display: grid; grid-template-columns: repeat(2, 1fr); gap: 24px; }
.advantage-item { display: flex; gap: 16px; padding: 20px; background: #f8fafc; border-radius: 10px; }
.check { color: #10b981; font-size: 1.5rem; font-weight: bold; }
.advantage-item strong { font-size: 1.1rem; display: block; margin-bottom: 4px; }
.advantage-item p { color: #64748b; font-size: 0.9rem; margin: 0; }

/* Quick Start */
.quickstart { padding: 80px 60px; background: #f8fafc; }
.quickstart h2 { text-align: center; font-size: 2rem; margin-bottom: 50px; }
.steps { display: flex; justify-content: center; align-items: center; gap: 20px; flex-wrap: wrap; }
.step { background: white; padding: 30px; border-radius: 12px; text-align: center; width: 180px; box-shadow: 0 4px 20px rgba(0,0,0,0.05); }
.step-num { width: 40px; height: 40px; background: linear-gradient(90deg, #2563eb, #7c3aed); color: white; border-radius: 50%; display: flex; align-items: center; justify-content: center; margin: 0 auto 16px; font-weight: 700; }
.step h4 { margin: 0 0 8px; }
.step p { color: #64748b; font-size: 0.85rem; margin: 0; }
.step-arrow { font-size: 1.5rem; color: #cbd5e1; }

/* CTA */
.cta { padding: 80px 60px; background: linear-gradient(135deg, #0f172a 0%, #1e3a5f 100%); color: white; text-align: center; }
.cta h2 { font-size: 2rem; margin-bottom: 16px; }
.cta p { color: #94a3b8; margin-bottom: 30px; }

/* Footer */
.footer { padding: 40px 60px; background: #0f172a; color: #94a3b8; }
.footer-content { display: flex; justify-content: space-between; align-items: center; margin-bottom: 30px; }
.footer-brand { font-size: 1.2rem; color: white; }
.footer-brand p { font-size: 0.85rem; color: #64748b; margin-top: 4px; }
.footer-links { display: flex; gap: 24px; }
.footer-links a { color: #94a3b8; text-decoration: none; }
.footer-links a:hover { color: white; }
.footer-bottom { text-align: center; font-size: 0.85rem; border-top: 1px solid #1e293b; padding-top: 20px; }
</style>
