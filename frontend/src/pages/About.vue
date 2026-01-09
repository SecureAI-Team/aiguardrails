<template>
  <LandingLayout>
    <div class="page-container">
      <div class="hero">
        <h1>关于我们</h1>
        <p class="mission">我们的使命是为企业构建可信赖的 AI 安全基础设施。</p>
      </div>

      <div class="content-section">
        <h2>Who We Are</h2>
        <p>AI GuardRails 成立于 2024 年，由一群来自网络安全和人工智能领域的专家创立。我们致力于解决生成式 AI 应用在落地过程中面临的安全挑战，包括提示词注入、数据泄露、内容合规等问题。</p>
        <p>我们相信，只有建立起坚固的安全护栏，企业才能真正放心地释放 AI 的潜力。</p>
      </div>

      <div class="stats-grid">
        <div class="stat">
          <div class="number">1M+</div>
          <div class="label">每日请求防护</div>
        </div>
        <div class="stat">
          <div class="number">50+</div>
          <div class="label">企业客户</div>
        </div>
        <div class="stat">
          <div class="number">99.9%</div>
          <div class="label">系统可用性</div>
        </div>
      </div>
    </div>
  </LandingLayout>
</template>

<script setup lang="ts">
import LandingLayout from '../components/LandingLayout.vue'
</script>

<style scoped>
.page-container {
  padding: 0 0 60px;
}
.hero {
  background: #0f172a;
  color: white;
  padding: 80px 20px;
  text-align: center;
}
.hero h1 { font-size: 3rem; margin-bottom: 24px; }
.mission { font-size: 1.5rem; color: #94a3b8; max-width: 800px; margin: 0 auto; }

.content-section {
  max-width: 800px;
  margin: 60px auto;
  padding: 0 20px;
  text-align: center;
}
.content-section h2 { font-size: 2rem; color: #1e293b; margin-bottom: 24px; }
.content-section p { font-size: 1.1rem; color: #475569; line-height: 1.8; margin-bottom: 20px; }

.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 40px;
  max-width: 1000px;
  margin: 60px auto;
  text-align: center;
  padding: 0 20px;
}
.stat .number { font-size: 3.5rem; font-weight: 800; color: #2563eb; margin-bottom: 8px; }
.stat .label { color: #64748b; font-size: 1.2rem; }

@media (max-width: 768px) {
  .hero h1 { font-size: 2.5rem; }
  .mission { font-size: 1.2rem; }
  .stats-grid { grid-template-columns: 1fr; gap: 40px; }
}
</style>
