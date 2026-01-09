<template>
  <LandingLayout>
    <div class="page-container">
      <div class="content">
        <h1>隐私政策</h1>
        <p class="update-time">最后更新日期：2024年1月1日</p>

        <section>
          <h2>1. 信息收集</h2>
          <p>作为安全服务提供商，我们会收集您在使用 AI GuardRails 服务过程中产生的数据，包括但不限于 API 调用日志、策略配置信息以及用户账号信息。对于 Prompt 内容，我们仅在进行安全检测时进行瞬时处理，除非您开启了特定的审计日志留存功能，否则我们不会持久化存储您的原始内容。</p>
        </section>

        <section>
          <h2>2. 数据使用</h2>
          <p>我们收集的数据主要用于：</p>
          <ul>
            <li>提供安全检测与防护服务</li>
            <li>改进我们的检测模型效果</li>
            <li>生成用量统计与计费依据</li>
            <li>发送服务变更通知</li>
          </ul>
        </section>

        <section>
          <h2>3. 数据安全</h2>
          <p>我们采用行业标准的安全措施来保护您的数据，包括传输加密 (TLS 1.3)、静态数据加密 (AES-256) 以及严格的访问控制机制 (RBAC)。</p>
        </section>

        <section>
          <h2>4. 第三方共享</h2>
          <p>我们承诺不会将您的数据出售给第三方。仅在法律要求或为了提供您所请求的服务（如短信通知）时，我们才会与合作伙伴共享必要的最少化数据。</p>
        </section>
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
  background: #f8fafc;
  min-height: 80vh;
}
.content {
  max-width: 800px;
  margin: 0 auto;
  background: white;
  padding: 60px;
  border-radius: 12px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05);
}

h1 { font-size: 2.5rem; color: #1e293b; margin-bottom: 8px; }
.update-time { color: #64748b; margin-bottom: 40px; }

section { margin-bottom: 40px; }
h2 { font-size: 1.5rem; color: #0f172a; margin-bottom: 16px; }
p { color: #475569; line-height: 1.8; font-size: 1.1rem; margin-bottom: 16px; }
ul { margin-left: 20px; color: #475569; line-height: 1.8; font-size: 1.1rem; }
li { margin-bottom: 8px; }

@media (max-width: 768px) {
  .content { padding: 30px; }
}
</style>
