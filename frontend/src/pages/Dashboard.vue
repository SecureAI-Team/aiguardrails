<template>
  <div class="dashboard-page">
    <div class="welcome-banner">
      <div class="welcome-text">
        <h1>👋 早安，{{ user.display_name || user.username }}</h1>
        <p>今天是 {{ today }}，系统运行正常。</p>
      </div>
      <div class="welcome-actions">
        <router-link to="/profile" class="btn-glass">个人设置</router-link>
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon bg-blue-100 text-blue-600">🏢</div>
        <div class="stat-content">
          <span class="stat-value">{{ stats.tenants }}</span>
          <span class="stat-label">租户数量</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon bg-purple-100 text-purple-600">📱</div>
        <div class="stat-content">
          <span class="stat-value">{{ stats.apps }}</span>
          <span class="stat-label">应用总数</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon bg-green-100 text-green-600">🛡️</div>
        <div class="stat-content">
          <span class="stat-value">{{ stats.policies }}</span>
          <span class="stat-label">活跃策略</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon bg-orange-100 text-orange-600">🚫</div>
        <div class="stat-content">
          <span class="stat-value">{{ stats.blocked }}</span>
          <span class="stat-label">今日阻断</span>
        </div>
      </div>
    </div>

    <div class="dashboard-content">
      <div class="main-column">
        <!-- Quick Actions -->
        <section class="panel quick-actions-panel">
          <h3>⚡ 快捷操作</h3>
          <div class="action-grid">
            <router-link to="/tenants" class="action-card">
              <span class="action-icon">🏢</span>
              <span class="action-name">管理租户</span>
            </router-link>
            <router-link to="/apps" class="action-card">
              <span class="action-icon">📱</span>
              <span class="action-name">管理应用</span>
            </router-link>
            <router-link to="/policies" class="action-card">
              <span class="action-icon">📋</span>
              <span class="action-name">配置策略</span>
            </router-link>
            <router-link to="/rules" class="action-card">
              <span class="action-icon">📜</span>
              <span class="action-name">规则库</span>
            </router-link>
            <router-link to="/alerts" class="action-card">
              <span class="action-icon">🚨</span>
              <span class="action-name">告警中心</span>
            </router-link>
            <router-link to="/logs" class="action-card">
              <span class="action-icon">📝</span>
              <span class="action-name">审计日志</span>
            </router-link>
          </div>
        </section>

        <!-- Traces Chart (Mock for now, real data points) -->
        <section class="panel">
          <h3>📈 流量趋势 (最近24小时)</h3>
          <div class="chart-box">
             <!-- Simple CSS Bar Chart -->
             <div class="bar-chart">
               <div v-if="chartData.length === 0" class="no-data">暂无流量数据</div>
               <div v-for="d in chartData" :key="d.date" class="bar-col" :title="d.requests + ' Requests'">
                 <!-- Use relative height based on max value -->
                 <div class="bar-fill" :style="{ height: getBarHeight(d.requests) + '%' }"></div>
                 <span class="bar-label">{{ formatDate(d.date) }}</span>
               </div>
             </div>
          </div>
        </section>
      </div>

      <div class="side-column">
        <!-- Recent Audits -->
        <section class="panel">
          <h3>🔔 最近活动</h3>
          <div class="activity-feed">
             <div v-for="log in recentLogs" :key="log.id" class="feed-item">
               <div class="feed-icon" :class="getEventClass(log.event)">●</div>
               <div class="feed-content">
                 <div class="feed-text">
                   <span class="font-medium">{{ log.user_id || 'System' }}</span>
                   {{ log.event }}
                 </div>
                 <div class="feed-time">{{ formatTime(log.created_at) }}</div>
               </div>
             </div>
             <div v-if="recentLogs.length === 0" class="empty-feed">暂无活动</div>
          </div>
          <router-link to="/logs" class="view-all">查看全部 →</router-link>
        </section>

        <!-- System Status -->
        <section class="panel system-status-panel">
          <h3>🖥️ 系统状态</h3>
          <div class="status-items">
            <div class="status-row">
               <span>API 服务</span>
               <span class="badge badge-success">运行中</span>
            </div>
            <div class="status-row">
               <span>数据库</span>
               <span class="badge badge-success">已连接</span>
            </div>
            <div class="status-row">
               <span>Redis</span>
               <span class="badge badge-success">已连接</span>
            </div>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '../services/api'

const user = ref({ username: 'admin', display_name: '' })
const today = new Date().toLocaleDateString('zh-CN', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })

const stats = ref({
  tenants: 0,
  apps: 0,
  policies: 0,
  blocked: 0
})

const recentLogs = ref<any[]>([])
const chartData = ref<any[]>([])

function getBarHeight(val: number) {
  if (!chartData.value.length) return 0
  const max = Math.max(...chartData.value.map(d => d.requests || 0)) || 1
  return Math.max(5, (val / max) * 100) // min 5% height
}

function formatDate(ds: string) {
  const d = new Date(ds)
  return `${d.getMonth() + 1}/${d.getDate()}`
}

onMounted(async () => {
  loadUser()
  loadStats()
})

function loadUser() {
  const token = localStorage.getItem('auth_token')
  if (token) {
    try {
      const payload = JSON.parse(atob(token.split('.')[1]))
      user.value = { 
        username: payload.sub || 'admin',
        display_name: payload.display_name || ''
      }
    } catch {}
  }
}

async function loadStats() {
  try {
    const [tenantsResult, logsResult, summaryRes] = await Promise.all([
      api.listTenants().catch(() => []),
      api.listAudit(5).catch(() => []),
      api.get('/usage/summary?days=7').catch(() => [])
    ])
    
    const tenants = Array.isArray(tenantsResult) ? tenantsResult : []
    const logs = Array.isArray(logsResult) ? logsResult : []
    const summary = Array.isArray(summaryRes) ? summaryRes : []

    if (summary.length > 0) {
      // Sort chronologically
      chartData.value = summary.sort((a: any, b: any) => new Date(a.date).getTime() - new Date(b.date).getTime())
    } else {
      chartData.value = []
    }
    
    stats.value.tenants = tenants.length
    
    // Fetch apps for each tenant (capped to first 5 tenants to avoid flooding)
    let totalApps = 0
    let totalPolicies = 0
    for (const t of tenants.slice(0, 5)) {
      try {
        const appsResult = await api.listApps(t.id)
        const apps = Array.isArray(appsResult) ? appsResult : []
        totalApps += apps.length
        
        const policiesResult = await api.listPolicies(t.id)
        const policies = Array.isArray(policiesResult) ? policiesResult : []
        totalPolicies += policies.length
      } catch {}
    }
    stats.value.apps = totalApps
    stats.value.policies = totalPolicies
    
    recentLogs.value = logs
    
    // Attempt to get trace stats if available
    try {
      const traces = await api.get('/traces?blocked=true&limit=10')
      if (Array.isArray(traces)) {
         stats.value.blocked = traces.length
      }
    } catch {}
    
  } catch (e) {
    console.error('Failed to load dashboard stats', e)
  }
}

function getEventClass(event: string) {
  if (event.includes('delete') || event.includes('block')) return 'text-red-500'
  if (event.includes('create')) return 'text-green-500'
  return 'text-blue-500'
}

function formatTime(ts: string) {
  if (!ts) return ''
  const d = new Date(ts)
  const now = new Date()
  const diff = (now.getTime() - d.getTime()) / 1000
  if (diff < 60) return '刚刚'
  if (diff < 3600) return Math.floor(diff / 60) + '分钟前'
  return d.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
}
</script>

<style scoped>
.dashboard-page { padding: 24px; max-width: 1400px; margin: 0 auto; }

.welcome-banner { 
  background: linear-gradient(135deg, #2563eb 0%, #4f46e5 100%); 
  padding: 30px; 
  border-radius: 16px; 
  color: white;
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
  box-shadow: 0 4px 15px rgba(37, 99, 235, 0.2);
}
.welcome-text h1 { margin: 0 0 8px; font-size: 1.8rem; }
.welcome-text p { margin: 0; opacity: 0.9; }
.btn-glass {
  background: rgba(255,255,255,0.2);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255,255,255,0.3);
  color: white;
  padding: 10px 20px;
  border-radius: 8px;
  text-decoration: none;
  font-weight: 500;
  transition: all 0.2s;
}
.btn-glass:hover { background: rgba(255,255,255,0.3); }

.stats-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 24px; margin-bottom: 32px; }
.stat-card { background: white; padding: 24px; border-radius: 12px; display: flex; align-items: center; gap: 20px; box-shadow: 0 2px 8px rgba(0,0,0,0.05); }
.stat-icon { width: 56px; height: 56px; border-radius: 16px; display: flex; align-items: center; justify-content: center; font-size: 1.75rem; }
.stat-content { display: flex; flex-direction: column; }
.stat-value { font-size: 2rem; font-weight: 700; color: #1e293b; line-height: 1.2; }
.stat-label { color: #64748b; font-size: 0.9rem; }

.dashboard-content { display: grid; grid-template-columns: 2fr 1fr; gap: 24px; }
.panel { background: white; border-radius: 12px; padding: 24px; box-shadow: 0 2px 10px rgba(0,0,0,0.05); margin-bottom: 24px; }
.panel h3 { margin: 0 0 20px; color: #1e293b; font-size: 1.1rem; }

.action-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 16px; }
.action-card { 
  display: flex; flex-direction: column; align-items: center; justify-content: center;
  padding: 24px; background: #f8fafc; border-radius: 12px; text-decoration: none; color: #334155;
  transition: all 0.2s; border: 1px solid transparent;
}
.action-card:hover { background: white; border-color: #e2e8f0; transform: translateY(-2px); box-shadow: 0 4px 12px rgba(0,0,0,0.05); }
.action-icon { font-size: 2rem; margin-bottom: 12px; }
.action-name { font-weight: 500; }

.chart-box { height: 200px; display: flex; align-items: flex-end; }
.bar-chart { display: flex; width: 100%; height: 100%; align-items: flex-end; justify-content: space-between; gap: 8px; }
.bar-col { flex: 1; display: flex; flex-direction: column; align-items: center; height: 100%; justify-content: flex-end; }
.bar-fill { width: 100%; background: #3b82f6; border-radius: 4px 4px 0 0; opacity: 0.8; transition: height 0.5s; min-height: 4px; }
.bar-label { margin-top: 8px; font-size: 0.75rem; color: #94a3b8; }

.activity-feed { display: flex; flex-direction: column; gap: 16px; }
.feed-item { display: flex; gap: 12px; align-items: flex-start; padding-bottom: 16px; border-bottom: 1px solid #f1f5f9; }
.feed-item:last-child { border-bottom: none; }
.feed-icon { font-size: 0.8rem; margin-top: 4px; }
.feed-text { font-size: 0.95rem; color: #334155; margin-bottom: 4px; }
.feed-time { font-size: 0.8rem; color: #94a3b8; }
.view-all { display: block; text-align: center; margin-top: 16px; color: #2563eb; text-decoration: none; font-size: 0.9rem; }

.status-items { display: flex; flex-direction: column; gap: 12px; }
.status-row { display: flex; justify-content: space-between; align-items: center; padding: 12px; background: #f8fafc; border-radius: 8px; }
.badge { padding: 4px 10px; border-radius: 12px; font-size: 0.75rem; font-weight: 600; }
.badge-success { background: #dcfce7; color: #166534; }

/* Utilities */
.bg-blue-100 { background: #dbeafe; } .text-blue-600 { color: #2563eb; }
.bg-purple-100 { background: #f3e8ff; } .text-purple-600 { color: #9333ea; }
.bg-green-100 { background: #dcfce7; } .text-green-600 { color: #16a34a; }
.bg-orange-100 { background: #ffedd5; } .text-orange-600 { color: #ea580c; }
.text-red-500 { color: #ef4444; } .text-blue-500 { color: #3b82f6; } .text-green-500 { color: #10b981; }
</style>
