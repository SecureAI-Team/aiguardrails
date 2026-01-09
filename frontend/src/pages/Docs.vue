<template>
  <LandingLayout>
    <div class="page-container">
      <div class="content-wrapper">
        <h1 class="page-title">帮助文档</h1>
        <p class="page-subtitle">快速上手 AI GuardRails，构建安全的生成式 AI 应用</p>

        <div class="grid-cards">
          <div class="card">
            <h3>🚀 快速开始</h3>
            <p>了解如何集成 SDK 并发送第一个请求。</p>
            <a href="#" @click.prevent="alert('快速指南正在编写中，请参考 SDK 文档')" class="link">查看指南 →</a>
          </div>
          <div class="card">
            <h3>🛡️ 核心概念</h3>
            <p>深入理解策略、规则、租户与应用架构。</p>
            <a href="#" @click.prevent="alert('核心概念文档即将上线')" class="link">阅读更多 →</a>
          </div>
          <div class="card">
            <h3>🔧 SDK 集成</h3>
            <p>Go, Python, Java, Node.js SDK 详细文档。</p>
            <router-link to="/sdks" class="link">下载 SDK →</router-link>
          </div>
          <div class="card">
            <h3>❓ 常见问题</h3>
            <p>关于 API 配额、计费及部署的常见疑问。</p>
            <a href="#" @click.prevent="alert('FAQ 页面建设中')" class="link">查看 FAQ →</a>
          </div>
        </div>
      </div>
    </div>
  </LandingLayout>
</template>

<script setup lang="ts">
import LandingLayout from '../components/LandingLayout.vue'

function alert(msg: string) {
  window.alert(msg)
}
</script>

<style scoped>
.page-container {
  padding: 60px 20px;
  max-width: 1200px;
  margin: 0 auto;
}
.content-wrapper {
  text-align: center;
}
.page-title {
  font-size: 2.5rem;
  color: #1e293b;
  margin-bottom: 16px;
}
.page-subtitle {
  font-size: 1.2rem;
  color: #64748b;
  margin-bottom: 60px;
}
.grid-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 30px;
  text-align: left;
}
.card {
  background: white;
  padding: 30px;
  border-radius: 12px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  border: 1px solid #e2e8f0;
  transition: transform 0.2s;
}
.card:hover {
  transform: translateY(-5px);
}
.card h3 {
  margin: 0 0 12px;
  color: #0f172a;
}
.card p {
  color: #64748b;
  line-height: 1.6;
  margin-bottom: 20px;
}
.link {
  color: #2563eb;
  text-decoration: none;
  font-weight: 500;
}
.link:hover {
  text-decoration: underline;
}
</style>
