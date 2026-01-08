<template>
  <LandingLayout>
    <div class="sdks-page">
      <div class="page-header">
        <h2>📦 SDK下载</h2>
        <p class="subtitle">多语言SDK助您快速集成AI GuardRails</p>
      </div>

      <div class="sdk-grid">
        <div class="sdk-card">
          <div class="sdk-icon">🐹</div>
          <h3>Go SDK</h3>
          <p>适用于Go 1.18+</p>
          <code>go get github.com/aiguardrails/go-sdk</code>
          <div class="sdk-links">
            <a href="#" class="btn-primary">📥 下载</a>
            <a href="#" class="btn-outline">📖 文档</a>
          </div>
        </div>

        <div class="sdk-card">
          <div class="sdk-icon">🐍</div>
          <h3>Python SDK</h3>
          <p>适用于Python 3.8+</p>
          <code>pip install aiguardrails</code>
          <div class="sdk-links">
            <a href="#" class="btn-primary">📥 下载</a>
            <a href="#" class="btn-outline">📖 文档</a>
          </div>
        </div>

        <div class="sdk-card">
          <div class="sdk-icon">🟨</div>
          <h3>Node.js SDK</h3>
          <p>适用于Node.js 16+</p>
          <code>npm install @aiguardrails/sdk</code>
          <div class="sdk-links">
            <a href="#" class="btn-primary">📥 下载</a>
            <a href="#" class="btn-outline">📖 文档</a>
          </div>
        </div>

        <div class="sdk-card">
          <div class="sdk-icon">☕</div>
          <h3>Java SDK</h3>
          <p>适用于Java 11+</p>
          <code>&lt;dependency&gt;aiguardrails-sdk&lt;/dependency&gt;</code>
          <div class="sdk-links">
            <a href="#" class="btn-primary">📥 下载</a>
            <a href="#" class="btn-outline">📖 文档</a>
          </div>
        </div>
      </div>

      <div class="section">
        <h3>📖 快速开始</h3>
        <div class="tabs">
          <button :class="{ active: activeTab === 'go' }" @click="activeTab = 'go'">Go</button>
          <button :class="{ active: activeTab === 'python' }" @click="activeTab = 'python'">Python</button>
          <button :class="{ active: activeTab === 'node' }" @click="activeTab = 'node'">Node.js</button>
        </div>
        <pre v-if="activeTab === 'go'" class="code-block">
import "github.com/aiguardrails/go-sdk"

client := aiguardrails.NewClient("sk_your_api_key")

// 检查提示词安全
result, err := client.PromptCheck(ctx, &aiguardrails.PromptCheckRequest{
    Prompt: "用户输入内容",
})
if result.Blocked {
    log.Println("内容被阻断:", result.Reason)
}</pre>
        <pre v-if="activeTab === 'python'" class="code-block">
from aiguardrails import Client

client = Client(api_key="sk_your_api_key")

# 检查提示词安全
result = client.prompt_check(prompt="用户输入内容")
if result.blocked:
    print(f"内容被阻断: {result.reason}")</pre>
        <pre v-if="activeTab === 'node'" class="code-block">
import { Client } from '@aiguardrails/sdk';

const client = new Client({ apiKey: 'sk_your_api_key' });

// 检查提示词安全
const result = await client.promptCheck({ prompt: '用户输入内容' });
if (result.blocked) {
    console.log('内容被阻断:', result.reason);
}</pre>
      </div>

      <div class="section">
        <h3>🔗 API接口</h3>
        <div class="api-list">
          <div class="api-item">
            <span class="method post">POST</span>
            <span class="path">/v1/guardrails/prompt-check</span>
            <span class="desc">提示词安全检查</span>
          </div>
          <div class="api-item">
            <span class="method post">POST</span>
            <span class="path">/v1/guardrails/output-filter</span>
            <span class="desc">输出内容过滤</span>
          </div>
          <div class="api-item">
            <span class="method post">POST</span>
            <span class="path">/v1/guardrails/rag-check</span>
            <span class="desc">RAG安全检查</span>
          </div>
          <div class="api-item">
            <span class="method get">GET</span>
            <span class="path">/v1/traces</span>
            <span class="desc">请求追踪查询</span>
          </div>
        </div>
      </div>
    </div>
  </LandingLayout>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import LandingLayout from '../components/LandingLayout.vue'

const activeTab = ref('go')
</script>

<style scoped>
.sdks-page { padding: 40px 48px; max-width: 1200px; margin: 0 auto; }
.page-header { margin-bottom: 32px; }
.page-header h2 { margin: 0 0 8px; color: #1e293b; }
.subtitle { color: #64748b; margin: 0; }
.sdk-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(280px, 1fr)); gap: 20px; margin-bottom: 40px; }
.sdk-card { background: white; padding: 24px; border-radius: 12px; box-shadow: 0 2px 10px rgba(0,0,0,0.05); text-align: center; }
.sdk-icon { font-size: 3rem; margin-bottom: 12px; }
.sdk-card h3 { margin: 0 0 8px; }
.sdk-card p { color: #64748b; margin: 0 0 12px; }
.sdk-card code { display: block; background: #f1f5f9; padding: 10px; border-radius: 6px; font-size: 0.85rem; margin-bottom: 16px; }
.sdk-links { display: flex; gap: 12px; justify-content: center; }
.btn-primary { background: linear-gradient(90deg, #2563eb, #7c3aed); color: white; padding: 8px 16px; border-radius: 6px; text-decoration: none; display: inline-block; }
.btn-outline { border: 1px solid #e2e8f0; padding: 8px 16px; border-radius: 6px; text-decoration: none; color: #334155; }
.section { background: white; padding: 24px; border-radius: 12px; margin-bottom: 24px; }
.section h3 { margin: 0 0 16px; }
.tabs { display: flex; gap: 8px; margin-bottom: 16px; }
.tabs button { padding: 8px 20px; border: 1px solid #e2e8f0; background: white; border-radius: 6px; cursor: pointer; }
.tabs button.active { background: #3b82f6; color: white; border-color: #3b82f6; }
.code-block { background: #1e293b; color: #e2e8f0; padding: 20px; border-radius: 8px; overflow-x: auto; font-family: monospace; font-size: 0.9rem; line-height: 1.6; margin: 0; }
.api-list { display: flex; flex-direction: column; gap: 12px; }
.api-item { display: flex; align-items: center; gap: 12px; padding: 12px; background: #f8fafc; border-radius: 6px; }
.method { padding: 4px 12px; border-radius: 4px; font-weight: 600; font-size: 0.75rem; }
.method.post { background: #d1fae5; color: #065f46; }
.method.get { background: #dbeafe; color: #1e40af; }
.path { font-family: monospace; color: #334155; }
.desc { color: #64748b; margin-left: auto; }
</style>
