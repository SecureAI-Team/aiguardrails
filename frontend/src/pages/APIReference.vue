<template>
  <LandingLayout>
    <div class="page-container">
      <div class="header">
        <h1>API 参考</h1>
        <p>Base URL: <code class="bg-gray">https://api.aiguardrails.com/v1</code></p>
      </div>

      <div class="section">
        <h2>Authentication</h2>
        <p>所有请求需在 Header 中包含认证 Token。</p>
        <div class="code-block">
          <code>Authorization: Bearer &lt;YOUR_API_TOKEN&gt;</code>
        </div>
      </div>

      <div class="section">
        <h2>Endpoints</h2>
        <table class="api-table">
          <thead>
            <tr>
              <th>Method</th>
              <th>Endpoint</th>
              <th>Description</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td><span class="badge get">GET</span></td>
              <td>/tenants</td>
              <td>获取租户列表</td>
            </tr>
            <tr>
              <td><span class="badge post">POST</span></td>
              <td>/tenants/{id}/policies</td>
              <td>创建安全策略</td>
            </tr>
            <tr>
              <td><span class="badge post">POST</span></td>
              <td>/guardrails/prompt-check</td>
              <td>检测 Prompt 安全性</td>
            </tr>
            <tr>
              <td><span class="badge post">POST</span></td>
              <td>/guardrails/output-filter</td>
              <td>过滤模型输出</td>
            </tr>
          </tbody>
        </table>
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
  margin-bottom: 60px;
  text-align: center;
}
.header h1 { font-size: 2.5rem; margin-bottom: 16px; color: #1e293b; }
.bg-gray { background: #e2e8f0; padding: 4px 8px; border-radius: 4px; font-family: monospace; }
.section { margin-bottom: 50px; background: white; padding: 40px; border-radius: 12px; border: 1px solid #e2e8f0; }
.section h2 { margin-top: 0; color: #0f172a; border-bottom: 1px solid #e2e8f0; padding-bottom: 16px; margin-bottom: 24px; }
.code-block { background: #1e293b; color: #e2e8f0; padding: 20px; border-radius: 8px; font-family: monospace; }
.api-table { width: 100%; border-collapse: collapse; }
.api-table th, .api-table td { padding: 16px; text-align: left; border-bottom: 1px solid #e2e8f0; }
.api-table th { background: #f8fafc; font-weight: 600; color: #475569; }
.badge { padding: 4px 8px; border-radius: 4px; font-size: 0.8rem; font-weight: bold; }
.get { background: #dbfafe; color: #0284c7; }
.post { background: #dcfce7; color: #16a34a; }
</style>
