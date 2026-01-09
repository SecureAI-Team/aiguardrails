<template>
  <LandingLayout>
    <div class="page-container">
      <div class="header">
        <h1>最佳实战</h1>
        <p>构建安全、可靠的生成式 AI 应用的前沿指南</p>
      </div>

      <div class="article-grid">
        <article class="article-card">
          <div class="badge">安全</div>
          <h2>防止 Prompt Injection 的 5 个关键策略</h2>
          <p>深入分析提示词注入攻击原理，以及如何利用 AI GuardRails 的多层防护机制进行有效拦截。</p>
          <a href="#" class="read-more">阅读全文 →</a>
        </article>

        <article class="article-card">
          <div class="badge badge-perf">性能</div>
          <h2>在 RAG 架构中优化 Guardrails 延迟</h2>
          <p>如何在不牺牲安全性的前提下，通过缓存、并行检测和策略优化，将延迟降低至毫秒级。</p>
          <a href="#" class="read-more">阅读全文 →</a>
        </article>

        <article class="article-card">
          <div class="badge badge-comp">合规</div>
          <h2>处理 PII 数据：自动化脱敏的最佳实践</h2>
          <p>如何配置敏感数据识别规则，确保在数据发送给 LLM 之前完成自动脱敏，满足 GDPR/CCPA 要求。</p>
          <a href="#" class="read-more">阅读全文 →</a>
        </article>
      </div>
    </div>
  </LandingLayout>
</template>

<script setup lang="ts">
import LandingLayout from '../components/LandingLayout.vue'
</script>

<style scoped>
.page-container {
  padding: 60px 20px;
  max-width: 1000px;
  margin: 0 auto;
}
.header {
  text-align: center;
  margin-bottom: 60px;
}
.header h1 { font-size: 2.5rem; color: #1e293b; margin-bottom: 16px; }
.header p { color: #64748b; font-size: 1.2rem; }

.article-grid {
  display: grid;
  gap: 40px;
}

.article-card {
  background: white;
  padding: 40px;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  transition: all 0.2s;
  position: relative;
}
.article-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 10px 30px rgba(0,0,0,0.05);
}

.badge {
  display: inline-block;
  background: #dbeafe;
  color: #2563eb;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 0.85rem;
  font-weight: 600;
  margin-bottom: 16px;
}
.badge-perf { background: #dcfce7; color: #16a34a; }
.badge-comp { background: #fef9c3; color: #ca8a04; }

.article-card h2 { margin: 0 0 12px; color: #0f172a; font-size: 1.5rem; }
.article-card p { color: #64748b; line-height: 1.6; margin-bottom: 24px; }
.read-more { color: #2563eb; text-decoration: none; font-weight: 500; font-size: 1.1rem; }
.read-more:hover { text-decoration: underline; }
</style>
