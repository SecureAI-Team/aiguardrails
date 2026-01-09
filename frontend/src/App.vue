<template>
  <div id="app">
    <!-- 公开页面：Landing/Login 不显示Console导航 -->
    <template v-if="isPublicPage">
      <RouterView />
    </template>
    
    <!-- 已登录：显示Console布局 -->
    <template v-else>
      <header class="topbar">
        <h1>AI Security Console</h1>
        <nav>
          <RouterLink to="/dashboard">仪表板</RouterLink>
          <RouterLink to="/tenants">租户</RouterLink>
          <RouterLink to="/apps">应用</RouterLink>
          <RouterLink to="/policies">策略</RouterLink>
          <RouterLink to="/rules">规则</RouterLink>
          <RouterLink to="/alerts">告警</RouterLink>
          <RouterLink to="/stats">统计</RouterLink>
          <RouterLink to="/apikeys">密钥</RouterLink>
          <RouterLink to="/traces">追踪</RouterLink>
          <RouterLink to="/users">用户</RouterLink>
          <RouterLink to="/orgs">组织</RouterLink>
          <RouterLink to="/logs">日志</RouterLink>
        </nav>
        <div class="user-menu">
          <span v-if="username">{{ username }}</span>
          <button @click="logout">退出</button>
        </div>
      </header>
      <main class="content">
        <RouterView />
      </main>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const username = ref('')

// 公开页面列表（不需要登录）
const publicPages = ['/', '/landing', '/login', '/sdks', '/playground', '/models', '/docs', '/api-reference', '/best-practices', '/about', '/contact', '/privacy']

const isPublicPage = computed(() => {
  return publicPages.includes(route.path)
})

function checkAuth() {
  const token = localStorage.getItem('token') || localStorage.getItem('auth_token')
  const user = localStorage.getItem('username')
  if (user) username.value = user
  
  // 如果未登录且不是公开页面，跳转到登录
  if (!token && !publicPages.includes(route.path)) {
    router.push('/login')
  }
}

watch(() => route.path, () => {
  checkAuth()
})

onMounted(() => {
  checkAuth()
})

function logout() {
  localStorage.removeItem('token')
  localStorage.removeItem('auth_token')
  localStorage.removeItem('username')
  router.push('/landing')
}
</script>

<style>
body {
  margin: 0;
  font-family: 'Inter', Arial, sans-serif;
  background: #f6f8fb;
}
.topbar {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px 20px;
  background: #0f172a;
  color: #fff;
}
.topbar h1 {
  margin: 0;
  font-size: 1.2rem;
  margin-right: 24px;
}
.topbar nav {
  flex: 1;
  display: flex;
  gap: 4px;
}
.topbar nav a {
  color: #94a3b8;
  padding: 6px 12px;
  text-decoration: none;
  border-radius: 6px;
  font-size: 0.9rem;
  transition: all 0.2s;
}
.topbar nav a:hover {
  background: rgba(255,255,255,0.1);
  color: #fff;
}
.topbar nav a.router-link-active {
  background: #3b82f6;
  color: #fff;
}
.user-menu {
  display: flex;
  align-items: center;
  gap: 12px;
}
.user-menu span {
  color: #94a3b8;
  font-size: 0.9rem;
}
.user-menu button {
  padding: 6px 12px;
  background: transparent;
  border: 1px solid #475569;
  color: #94a3b8;
  border-radius: 6px;
  cursor: pointer;
}
.user-menu button:hover {
  background: #ef4444;
  border-color: #ef4444;
  color: #fff;
}
.content {
  padding: 20px;
  min-height: calc(100vh - 60px);
}
</style>


