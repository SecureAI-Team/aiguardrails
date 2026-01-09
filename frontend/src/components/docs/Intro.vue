<template>
  <div class="doc-content">
    <h1>AI GuardRails 简介</h1>
    <p>AI GuardRails 是一个为生成式 AI应用设计的安全护栏系统。它通过拦截和检测 Prompt 以及模型输出，防止 Prompt 注入、数据泄露 (DLP) 和不合规内容的生成。</p>
    
    <h2>核心功能</h2>
    <ul>
      <li><strong>防御 Prompt 注入</strong>: 识别并拦截试图绕过系统指令的恶意输入。</li>
      <li><strong>敏感数据过滤 (DLP)</strong>: 自动检测输出中的身份证号、手机号、密钥等敏感信息。</li>
      <li><strong>RAG 安全检查</strong>: 确保存储和检索的知识库内容符合安全规范。</li>
      <li><strong>多租户架构</strong>: 支持企业级多租户、多应用隔离。</li>
    </ul>

    <div class="tip-box">
      <strong>提示:</strong> 您可以在 <code>/playground</code> 页面实时测试这些防护能力。
    </div>
  </div>
</template>

<style scoped>
.doc-content { line-height: 1.7; color: #334155; }
h1 { color: #0f172a; margin-bottom: 24px; }
h2 { margin-top: 40px; margin-bottom: 16px; color: #1e293b; }
ul { margin-bottom: 24px; padding-left: 20px; }
li { margin-bottom: 8px; }
.tip-box { background: #eff6ff; border-left: 4px solid #3b82f6; padding: 16px; margin: 24px 0; border-radius: 4px; }
</style>
