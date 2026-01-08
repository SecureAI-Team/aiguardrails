<template>
  <div class="dashboard">
    <header class="dash-header">
      <div class="greeting">
        <h1>👋 欢迎回来, {{ user.display_name || user.username }}</h1>
        <p>{{ today }} | 上次登录: {{ lastLogin }}</p>
      </div>
      <div class="header-actions">
        <router-link to="/profile" class="btn-outline">个人设置</router-link>
      </div>
    </header>

    <!-- Stats Cards -->
    <section class="stats-section">
      <div class="stat-card">
        <div class="stat-icon blue">📊</div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.requests }}</span>
          <span class="stat-label">今日请求</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon red">🚫</div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.blocked }}</span>
          <span class="stat-label">阻断次数</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon yellow">⚠️</div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.alerts }}</span>
          <span class="stat-label">安全告警</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon green">✅</div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.passRate }}%</span>
          <span class="stat-label">通过率</span>
        </div>
      </div>
    </section>

    <div class="dash-grid">
      <!-- Quick Actions -->
      <section class="card quick-actions">
        <h3>⚡ 快捷操作</h3>
        <div class="action-list">
          <router-link to="/" class="action-item">
            <span class="action-icon">🏢</span>
            <span>管理租户</span>
          </router-link>
          <router-link to="/apps" class="action-item">
            <span class="action-icon">📱</span>
            <span>管理应用</span>
          </router-link>
          <router-link to="/policies" class="action-item">
            <span class="action-icon">📋</span>
            <span>配置策略</span>
          </router-link>
          <router-link to="/rules" class="action-item">
            <span class="action-icon">📜</span>
            <span>管理规则</span>
          </router-link>
          <router-link to="/users" class="action-item">
            <span class="action-icon">👥</span>
            <span>用户管理</span>
          </router-link>
          <router-link to="/logs" class="action-item">
            <span class="action-icon">📝</span>
            <span>审计日志</span>
          </router-link>
        </div>
      </section>

      <!-- Recent Events -->
      <section class="card recent-events">
        <h3>🔔 安全事件</h3>
        <div class="event-list">
          <div v-for="event in events" :key="event.id" class="event-item" :class="'severity-' + event.severity">
            <span class="event-icon">{{ severityIcon(event.severity) }}</span>
            <div class="event-content">
              <span class="event-title">{{ event.title }}</span>
              <span class="event-time">{{ event.time }}</span>
            </div>
          </div>
          <div v-if="events.length === 0" class="no-events">
            ✨ 暂无安全事件，系统运行正常
          </div>
        </div>
      </section>

      <!-- Recent Activity -->
      <section class="card recent-activity">
        <h3>📅 最近活动</h3>
        <div class="activity-list">
          <div v-for="act in activities" :key="act.id" class="activity-item">
            <span class="activity-icon">{{ act.icon }}</span>
            <div class="activity-content">
              <span class="activity-title">{{ act.title }}</span>
              <span class="activity-time">{{ act.time }}</span>
            </div>
          </div>
        </div>
      </section>

      <!-- System Status -->
      <section class="card system-status">
        <h3>🖥️ 系统状态</h3>
        <div class="status-list">
          <div class="status-item">
            <span>API服务</span>
            <span class="status-badge online">运行中</span>
          </div>
          <div class="status-item">
            <span>OPA引擎</span>
            <span class="status-badge online">运行中</span>
          </div>
          <div class="status-item">
            <span>数据库</span>
            <span class="status-badge online">正常</span>
          </div>
          <div class="status-item">
            <span>Redis缓存</span>
            <span class="status-badge online">正常</span>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const user = ref({ username: 'admin', display_name: '' })
const today = new Date().toLocaleDateString('zh-CN', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })
const lastLogin = ref('10分钟前')

const stats = ref({
  requests: 1248,
  blocked: 45,
  alerts: 3,
  passRate: 96.4
})

const events = ref([
  { id: 1, title: '检测到提示注入攻击尝试', severity: 'high', time: '10分钟前' },
  { id: 2, title: '竞品厂商查询被阻断', severity: 'medium', time: '1小时前' },
  { id: 3, title: '敏感信息脱敏处理', severity: 'low', time: '3小时前' }
])

const activities = ref([
  { id: 1, icon: '📋', title: '策略「工业安全」已更新', time: '2小时前' },
  { id: 2, icon: '👤', title: '新用户 engineer01 加入', time: '5小时前' },
  { id: 3, icon: '📱', title: '应用「PLC助手」密钥已轮转', time: '昨天' }
])

onMounted(() => {
  const token = localStorage.getItem('auth_token')
  if (token) {
    try {
      const payload = JSON.parse(atob(token.split('.')[1]))
      user.value.username = payload.sub || 'admin'
    } catch {}
  }
})

function severityIcon(severity: string) {
  const icons: Record<string, string> = { high: '🔴', medium: '🟡', low: '🟢' }
  return icons[severity] || '⚪'
}
</script>

<style scoped>
.dashboard { padding: 30px; background: #f1f5f9; min-height: 100vh; }

.dash-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 30px; }
.greeting h1 { font-size: 1.75rem; margin: 0; color: #1e293b; }
.greeting p { color: #64748b; margin-top: 4px; }
.btn-outline { border: 1px solid #cbd5e1; padding: 10px 20px; border-radius: 8px; text-decoration: none; color: #475569; }

.stats-section { display: grid; grid-template-columns: repeat(4, 1fr); gap: 20px; margin-bottom: 30px; }
.stat-card { background: white; padding: 24px; border-radius: 12px; display: flex; align-items: center; gap: 16px; box-shadow: 0 2px 10px rgba(0,0,0,0.04); }
.stat-icon { width: 50px; height: 50px; border-radius: 12px; display: flex; align-items: center; justify-content: center; font-size: 1.5rem; }
.stat-icon.blue { background: #eff6ff; }
.stat-icon.red { background: #fef2f2; }
.stat-icon.yellow { background: #fffbeb; }
.stat-icon.green { background: #f0fdf4; }
.stat-value { font-size: 1.75rem; font-weight: 700; color: #1e293b; display: block; }
.stat-label { color: #64748b; font-size: 0.875rem; }

.dash-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 24px; }
.card { background: white; padding: 24px; border-radius: 12px; box-shadow: 0 2px 10px rgba(0,0,0,0.04); }
.card h3 { margin: 0 0 20px; font-size: 1.1rem; color: #1e293b; }

.action-list { display: grid; grid-template-columns: repeat(3, 1fr); gap: 12px; }
.action-item { display: flex; flex-direction: column; align-items: center; gap: 8px; padding: 20px; background: #f8fafc; border-radius: 10px; text-decoration: none; color: #1e293b; transition: background 0.2s; }
.action-item:hover { background: #e2e8f0; }
.action-icon { font-size: 1.75rem; }

.event-list, .activity-list { display: flex; flex-direction: column; gap: 12px; }
.event-item, .activity-item { display: flex; gap: 12px; padding: 12px; background: #f8fafc; border-radius: 8px; }
.event-item.severity-high { border-left: 3px solid #ef4444; }
.event-item.severity-medium { border-left: 3px solid #f59e0b; }
.event-item.severity-low { border-left: 3px solid #10b981; }
.event-content, .activity-content { flex: 1; }
.event-title, .activity-title { display: block; font-weight: 500; color: #1e293b; }
.event-time, .activity-time { font-size: 0.8rem; color: #94a3b8; }
.no-events { text-align: center; color: #64748b; padding: 30px; }

.status-list { display: flex; flex-direction: column; gap: 12px; }
.status-item { display: flex; justify-content: space-between; padding: 12px; background: #f8fafc; border-radius: 8px; }
.status-badge { padding: 4px 12px; border-radius: 20px; font-size: 0.8rem; font-weight: 500; }
.status-badge.online { background: #d1fae5; color: #065f46; }
</style>
