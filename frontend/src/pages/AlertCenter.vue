<template>
  <div class="alert-center">
    <div class="page-header">
      <h2>🚨 告警中心</h2>
      <div class="header-actions">
        <select v-model="filterSeverity" @change="loadAlerts">
          <option value="">全部等级</option>
          <option value="critical">严重</option>
          <option value="high">高</option>
          <option value="medium">中</option>
          <option value="low">低</option>
        </select>
        <button @click="loadAlerts" class="btn-outline">🔄 刷新</button>
      </div>
    </div>

    <div class="stats-row">
      <div class="stat-card critical">
        <span class="stat-num">{{ countBySeverity('critical') }}</span>
        <span class="stat-label">严重</span>
      </div>
      <div class="stat-card high">
        <span class="stat-num">{{ countBySeverity('high') }}</span>
        <span class="stat-label">高</span>
      </div>
      <div class="stat-card medium">
        <span class="stat-num">{{ countBySeverity('medium') }}</span>
        <span class="stat-label">中</span>
      </div>
      <div class="stat-card low">
        <span class="stat-num">{{ countBySeverity('low') }}</span>
        <span class="stat-label">低</span>
      </div>
    </div>

    <div class="alerts-list">
      <div v-for="alert in alerts" :key="alert.id" class="alert-item" :class="'severity-' + alert.severity">
        <div class="alert-header">
          <span class="alert-icon">{{ severityIcon(alert.severity) }}</span>
          <div class="alert-info">
            <span class="alert-title">{{ alert.title }}</span>
            <span class="alert-meta">{{ alert.rule_name }} · {{ formatTime(alert.created_at) }}</span>
          </div>
          <span :class="'severity-badge severity-' + alert.severity">{{ alert.severity }}</span>
        </div>
        <p class="alert-message">{{ alert.message }}</p>
        <div class="alert-footer">
          <div class="notify-status">
            <span v-for="(status, channel) in parseStatus(alert.notify_status)" :key="channel" 
                  :class="'status-tag ' + (status === 'sent' ? 'success' : 'failed')">
              {{ channel }}: {{ status === 'sent' ? '✓' : '✗' }}
            </span>
          </div>
          <div class="alert-actions">
            <button v-if="!alert.acknowledged" @click="acknowledge(alert.id)" class="btn-ack">
              ✓ 确认处理
            </button>
            <span v-else class="acked">已确认 {{ formatTime(alert.acknowledged_at) }}</span>
          </div>
        </div>
      </div>
      <div v-if="alerts.length === 0" class="no-alerts">
        ✅ 暂无告警，系统运行正常
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { api } from '../services/api'

interface AlertHistory {
  id: string
  rule_name: string
  title: string
  message: string
  severity: string
  notify_status: string
  acknowledged: boolean
  acknowledged_at?: string
  created_at: string
}

const alerts = ref<AlertHistory[]>([])
const filterSeverity = ref('')

onMounted(() => loadAlerts())

async function loadAlerts() {
  try {
    const query = filterSeverity.value ? `?severity=${filterSeverity.value}` : ''
    alerts.value = await api.get(`/alerts/history${query}`)
  } catch { alerts.value = [] }
}

function countBySeverity(sev: string) {
  return alerts.value.filter(a => a.severity === sev).length
}

function severityIcon(sev: string) {
  const icons: Record<string, string> = { critical: '🔴', high: '🟠', medium: '🟡', low: '🟢' }
  return icons[sev] || '⚪'
}

function formatTime(ts: string | undefined): string {
  if (!ts) return ''
  const d = new Date(ts)
  const now = new Date()
  const diff = (now.getTime() - d.getTime()) / 1000
  if (diff < 60) return '刚刚'
  if (diff < 3600) return Math.floor(diff / 60) + '分钟前'
  if (diff < 86400) return Math.floor(diff / 3600) + '小时前'
  return d.toLocaleString('zh-CN')
}

function parseStatus(status: string): Record<string, string> {
  try { return JSON.parse(status) } catch { return {} }
}

async function acknowledge(id: string) {
  try {
    await api.post(`/alerts/history/${id}/ack`, {})
    loadAlerts()
  } catch {}
}
</script>

<style scoped>
.alert-center { padding: 20px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
.page-header h2 { margin: 0; }
.header-actions { display: flex; gap: 12px; }
.header-actions select { padding: 8px 16px; border: 1px solid #e2e8f0; border-radius: 6px; }
.btn-outline { padding: 8px 16px; border: 1px solid #e2e8f0; background: white; border-radius: 6px; cursor: pointer; }

.stats-row { display: flex; gap: 16px; margin-bottom: 24px; }
.stat-card { flex: 1; padding: 20px; border-radius: 12px; text-align: center; }
.stat-card.critical { background: #fef2f2; }
.stat-card.high { background: #fff7ed; }
.stat-card.medium { background: #fefce8; }
.stat-card.low { background: #f0fdf4; }
.stat-num { font-size: 2rem; font-weight: 700; display: block; }
.stat-label { color: #64748b; }

.alerts-list { display: flex; flex-direction: column; gap: 12px; }
.alert-item { background: white; border-radius: 12px; padding: 16px; box-shadow: 0 2px 10px rgba(0,0,0,0.05); border-left: 4px solid; }
.alert-item.severity-critical { border-left-color: #ef4444; }
.alert-item.severity-high { border-left-color: #f97316; }
.alert-item.severity-medium { border-left-color: #eab308; }
.alert-item.severity-low { border-left-color: #22c55e; }
.alert-header { display: flex; align-items: flex-start; gap: 12px; margin-bottom: 8px; }
.alert-icon { font-size: 1.5rem; }
.alert-info { flex: 1; }
.alert-title { display: block; font-weight: 600; color: #1e293b; }
.alert-meta { color: #94a3b8; font-size: 0.85rem; }
.severity-badge { padding: 4px 12px; border-radius: 20px; font-size: 0.75rem; font-weight: 600; }
.severity-badge.severity-critical { background: #fee2e2; color: #991b1b; }
.severity-badge.severity-high { background: #ffedd5; color: #9a3412; }
.severity-badge.severity-medium { background: #fef9c3; color: #854d0e; }
.severity-badge.severity-low { background: #dcfce7; color: #166534; }
.alert-message { color: #475569; margin: 8px 0; }
.alert-footer { display: flex; justify-content: space-between; align-items: center; margin-top: 12px; }
.notify-status { display: flex; gap: 8px; }
.status-tag { padding: 2px 8px; border-radius: 4px; font-size: 0.75rem; }
.status-tag.success { background: #d1fae5; color: #065f46; }
.status-tag.failed { background: #fee2e2; color: #991b1b; }
.btn-ack { padding: 6px 16px; background: #10b981; color: white; border: none; border-radius: 6px; cursor: pointer; }
.acked { color: #10b981; font-size: 0.85rem; }
.no-alerts { text-align: center; padding: 60px; color: #64748b; background: white; border-radius: 12px; }
</style>
