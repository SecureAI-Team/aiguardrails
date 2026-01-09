<template>
  <div class="landing-layout">
    <!-- Header -->
    <nav class="nav">
      <div class="logo">
        <router-link to="/landing" class="logo-link">
          <span class="logo-icon">🛡️</span>
          <span class="logo-text">AI GuardRails</span>
        </router-link>
      </div>
      <div class="nav-center">
        <router-link to="/sdks">SDK下载</router-link>
        <router-link to="/playground">API调试</router-link>
        <router-link to="/models">模型目录</router-link>
      </div>
      <div class="nav-links" v-if="isLoggedIn">
        <span class="user-greeting">Hi, {{ username }}</span>
        <router-link to="/dashboard" class="btn-primary">进入控制台</router-link>
      </div>
      <div class="nav-links" v-else>
        <router-link to="/login" class="btn-outline">登录</router-link>
        <router-link to="/login" class="btn-primary">免费试用</router-link>
      </div>
    </nav>

    <!-- Main Content -->
    <main class="main-content">
      <slot></slot>
    </main>

    <!-- Footer -->
    <footer class="footer">
      <div class="footer-grid">
        <div class="footer-brand">
          <div class="logo">
            <span class="logo-icon">🛡️</span>
            <span>AI GuardRails</span>
          </div>
          <p>工业AI应用安全护栏平台</p>
          <p class="copyright">© 2024 AI GuardRails. All rights reserved.</p>
        </div>
        <div class="footer-links">
          <h4>产品</h4>
          <router-link to="/sdks">SDK下载</router-link>
          <router-link to="/playground">API调试</router-link>
          <router-link to="/models">模型目录</router-link>
        </div>
        <div class="footer-links">
          <h4>资源</h4>
          <router-link to="/docs">帮助文档</router-link>
          <router-link to="/api-reference">API参考</router-link>
          <router-link to="/best-practices">最佳实践</router-link>
        </div>
        <div class="footer-links">
          <h4>公司</h4>
          <router-link to="/about">关于我们</router-link>
          <router-link to="/contact">联系我们</router-link>
          <router-link to="/privacy">隐私政策</router-link>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const isLoggedIn = ref(false)
const username = ref('')

function alert(msg: string) {
  window.alert(msg)
}

onMounted(() => {
  const token = localStorage.getItem('token') || localStorage.getItem('auth_token')
  const user = localStorage.getItem('username')
  if (token) {
    isLoggedIn.value = true
    if (user) username.value = user
  }
})
</script>

<style scoped>
.landing-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

/* Navigation */
.nav {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 48px;
  background: rgba(15, 23, 42, 0.95);
  backdrop-filter: blur(10px);
  position: sticky;
  top: 0;
  z-index: 100;
}

.logo-link {
  display: flex;
  align-items: center;
  gap: 10px;
  text-decoration: none;
  color: white;
}

.logo-icon { font-size: 1.8rem; }
.logo-text { font-size: 1.3rem; font-weight: 700; color: white; }

.nav-center {
  display: flex;
  gap: 32px;
}
.nav-center a {
  color: #94a3b8;
  text-decoration: none;
  font-weight: 500;
  transition: color 0.2s;
}
.nav-center a:hover, .nav-center a.router-link-active {
  color: white;
}

.nav-links {
  display: flex;
  gap: 12px;
  align-items: center;
}

.btn-outline {
  padding: 10px 20px;
  border: 1px solid #475569;
  color: white;
  border-radius: 8px;
  text-decoration: none;
  transition: all 0.2s;
}
.btn-outline:hover {
  border-color: #3b82f6;
  color: #3b82f6;
}

.btn-primary {
  padding: 10px 20px;
  background: linear-gradient(90deg, #2563eb, #7c3aed);
  color: white;
  border-radius: 8px;
  text-decoration: none;
  font-weight: 500;
  transition: all 0.2s;
}
.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(37, 99, 235, 0.3);
}

.user-greeting {
  color: #fff;
  font-weight: 500;
  margin-right: 12px;
}

/* Main Content */
.main-content {
  flex: 1;
  background: #f8fafc;
}

/* Footer */
.footer {
  background: #0f172a;
  color: #94a3b8;
  padding: 48px;
}

.footer-grid {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr 1fr;
  gap: 48px;
  max-width: 1200px;
  margin: 0 auto;
}

.footer-brand .logo {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 1.2rem;
  font-weight: 600;
  color: white;
  margin-bottom: 12px;
}
.footer-brand p { margin: 8px 0; }
.footer-brand .copyright { color: #64748b; font-size: 0.85rem; margin-top: 24px; }

.footer-links h4 {
  color: white;
  margin: 0 0 16px;
  font-size: 0.95rem;
}
.footer-links a {
  display: block;
  color: #94a3b8;
  text-decoration: none;
  margin-bottom: 10px;
  font-size: 0.9rem;
  transition: color 0.2s;
}
.footer-links a:hover {
  color: white;
}

@media (max-width: 900px) {
  .nav { padding: 12px 20px; }
  .nav-center { display: none; }
  .footer-grid { grid-template-columns: 1fr; gap: 32px; }
}
</style>
