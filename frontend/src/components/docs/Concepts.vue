<template>
  <div class="doc-content">
    <h1>核心概念</h1>
    
    <h2>租户 (Tenant)</h2>
    <p>
      租户是系统的最高层级隔离单位，通常对应一个公司或组织。
      一个租户下可以创建多个应用和策略。租户之间的数据完全隔离。
    </p>

    <h2>应用 (App)</h2>
    <p>
      应用代表具体的业务接入点，例如 "客服机器人" 或 "内部问答系统"。
      每个应用拥有独立的 <strong>App ID</strong> 和 <strong>API Secret</strong>，用于接口认证。
      您还可以为每个应用设置独立的速率限制 (Quota)。
    </p>

    <h2>策略 (Policy)</h2>
    <p>
      策略是一组安全规则的集合。您可以定义：
    </p>
    <ul>
      <li><strong>Prompt 规则</strong>: 关键词屏蔽、正则匹配。</li>
      <li><strong>输出过滤</strong>: 敏感词库、PII (个人隐私信息) 检测类型。</li>
      <li><strong>RAG 命名空间</strong>: 允许访问的知识库范围。</li>
    </ul>
    
    <h2>Guardrails 接口</h2>
    <p>系统主要提供三类防护接口：</p>
    <ul>
      <li><code>/prompt-check</code>: 在发送给 LLM 之前检查用户输入。</li>
      <li><code>/output-filter</code>: 在返回给用户之前检查 LLM 输出。</li>
      <li><code>/rag-check</code>: 检索增强生成的上下文安全检查。</li>
    </ul>
  </div>
</template>

<style scoped>
.doc-content { line-height: 1.7; color: #334155; }
h1 { color: #0f172a; margin-bottom: 24px; }
h2 { margin-top: 40px; margin-bottom: 16px; color: #1e293b; }
ul { margin-bottom: 24px; padding-left: 20px; }
li { margin-bottom: 8px; }
</style>
