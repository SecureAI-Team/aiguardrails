<template>
  <div class="doc-content">
    <h1>开发指南</h1>

    <h2>配置 RAG 检查</h2>
    <p>
      Retrieval-Augmented Generation (RAG) 容易受到上下文注入攻击。
      要启用 RAG 检查，请在策略中配置允许的 <code>RAG Namespaces</code>。
    </p>
    <div class="code-block">
      <pre><code>// 示例：仅允许访问 finance 和 HR 知识库
"rag_namespaces": ["finance", "hr_policy"]</code></pre>
    </div>

    <h2>集成 SDK</h2>
    <p>我们提供了 Node.js 和 Python SDK 以简化集成。</p>
    <div class="code-block">
      <pre><code>npm install aiguardrails
// or
pip install aiguardrails</code></pre>
    </div>
    <p>初始化 SDK 时，传入 <code>app_id</code> 和 <code>app_secret</code> 即可自动处理认证和重试。</p>

    <h2>查看告警</h2>
    <p>
      当拦截发生时，系统会记录审计日志。如果配置了告警通道（如 Email 或 Webhook），
      管理员会立即收到通知。请定期检查 <strong>告警中心</strong> 以调整策略的误报率。
    </p>
  </div>
</template>

<style scoped>
.doc-content { line-height: 1.7; color: #334155; }
h1 { color: #0f172a; margin-bottom: 24px; }
h2 { margin-top: 40px; margin-bottom: 16px; color: #1e293b; }
.code-block { background: #1e293b; color: #f1f5f9; padding: 20px; border-radius: 8px; overflow-x: auto; margin: 20px 0; font-family: 'Fira Code', monospace; }
</style>
