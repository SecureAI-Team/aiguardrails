<template>
  <div class="stats-page">
    <div class="page-header">
      <h2>📊 用量统计</h2>
      <div class="header-actions">
        <select v-model="timeRange" @change="loadData">
          <option value="7">最近7天</option>
          <option value="30">最近30天</option>
          <option value="90">最近90天</option>
        </select>
      </div>
    </div>

    <div class="overview-cards">
      <div class="stat-card">
        <div class="stat-icon">📡</div>
        <div class="stat-content">
          <span class="stat-value">{{ overview.today?.requests || 0 }}</span>
          <span class="stat-label">今日请求</span>
        </div>
      </div>
      <div class="stat-card success">
        <div class="stat-icon">✅</div>
        <div class="stat-content">
          <span class="stat-value">{{ successRate }}%</span>
          <span class="stat-label">成功率</span>
        </div>
      </div>
      <div class="stat-card warning">
        <div class="stat-icon">🛡️</div>
        <div class="stat-content">
          <span class="stat-value">{{ overview.today?.blocked || 0 }}</span>
          <span class="stat-label">今日阻断</span>
        </div>
      </div>
      <div class="stat-card info">
        <div class="stat-icon">🪙</div>
        <div class="stat-content">
          <span class="stat-value">{{ formatNumber(overview.today?.tokens || 0) }}</span>
          <span class="stat-label">今日Token</span>
        </div>
      </div>
    </div>

    <div class="quota-section" v-if="overview.quota_percent > 0">
      <h3>配额使用</h3>
      <div class="quota-bar">
        <div class="quota-fill" :style="{ width: Math.min(overview.quota_percent, 100) + '%' }"
             :class="{ warning: overview.quota_percent > 80, danger: overview.quota_percent > 95 }"></div>
      </div>
      <span class="quota-text">{{ overview.quota_percent.toFixed(1) }}% 已使用</span>
    </div>

    <div class="chart-section">
      <h3>请求趋势</h3>
      <div class="chart-container">
        <div class="chart-bars">
          <div v-for="day in dailyData" :key="day.date" class="chart-bar-group">
            <div class="bar-container">
              <div class="bar success" :style="{ height: getBarHeight(day.success) + 'px' }"></div>
              <div class="bar error" :style="{ height: getBarHeight(day.errors) + 'px' }"></div>
              <div class="bar blocked" :style="{ height: getBarHeight(day.blocked) + 'px' }"></div>
            </div>
            <span class="bar-label">{{ formatDate(day.date) }}</span>
          </div>
        </div>
        <div class="chart-legend">
          <span class="legend-item"><span class="dot success"></span> 成功</span>
          <span class="legend-item"><span class="dot error"></span> 错误</span>
          <span class="legend-item"><span class="dot blocked"></span> 阻断</span>
        </div>
      </div>
    </div>

    <div class="summary-section">
      <h3>本月统计</h3>
      <div class="summary-grid">
        <div class="summary-item">
          <span class="summary-value">{{ formatNumber(overview.month?.requests || 0) }}</span>
          <span class="summary-label">总请求数</span>
        </div>
        <div class="summary-item">
          <span class="summary-value">{{ formatNumber(overview.month?.tokens || 0) }}</span>
          <span class="summary-label">Token消耗</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { api } from '../services/api'

const timeRange = ref(7)
const overview = ref<any>({})
const dailyData = ref<any[]>([])

onMounted(() => loadData())

async function loadData() {
  try {
    overview.value = await api.get('/usage/overview')
    dailyData.value = await api.get(`/usage/summary?days=${timeRange.value}`)
  } catch {}
}

const successRate = computed(() => {
  const today = overview.value.today
  if (!today || today.requests === 0) return 100
  return ((today.requests - today.errors) / today.requests * 100).toFixed(1)
})

function formatNumber(n: number): string {
  if (n >= 1000000) return (n / 1000000).toFixed(1) + 'M'
  if (n >= 1000) return (n / 1000).toFixed(1) + 'K'
  return n.toString()
}

function formatDate(dateStr: string): string {
  const d = new Date(dateStr)
  return `${d.getMonth() + 1}/${d.getDate()}`
}

function getBarHeight(value: number): number {
  const maxVal = Math.max(...dailyData.value.map(d => d.requests || 1))
  return Math.max(4, (value / maxVal) * 120)
}
</script>

<style scoped>
.stats-page { padding: 20px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
.page-header h2 { margin: 0; }
.header-actions select { padding: 8px 16px; border: 1px solid #e2e8f0; border-radius: 6px; }

.overview-cards { display: grid; grid-template-columns: repeat(4, 1fr); gap: 16px; margin-bottom: 24px; }
.stat-card { background: white; padding: 20px; border-radius: 12px; display: flex; align-items: center; gap: 16px; box-shadow: 0 2px 10px rgba(0,0,0,0.05); }
.stat-card.success { background: linear-gradient(135deg, #ecfdf5 0%, #d1fae5 100%); }
.stat-card.warning { background: linear-gradient(135deg, #fffbeb 0%, #fef3c7 100%); }
.stat-card.info { background: linear-gradient(135deg, #eff6ff 0%, #dbeafe 100%); }
.stat-icon { font-size: 2rem; }
.stat-content { display: flex; flex-direction: column; }
.stat-value { font-size: 1.8rem; font-weight: 700; color: #1e293b; }
.stat-label { color: #64748b; font-size: 0.9rem; }

.quota-section { background: white; padding: 20px; border-radius: 12px; margin-bottom: 24px; }
.quota-section h3 { margin: 0 0 12px; color: #374151; }
.quota-bar { background: #e2e8f0; height: 12px; border-radius: 6px; overflow: hidden; }
.quota-fill { height: 100%; background: #10b981; transition: width 0.3s; }
.quota-fill.warning { background: #f59e0b; }
.quota-fill.danger { background: #ef4444; }
.quota-text { color: #64748b; font-size: 0.85rem; margin-top: 8px; display: block; }

.chart-section { background: white; padding: 20px; border-radius: 12px; margin-bottom: 24px; }
.chart-section h3 { margin: 0 0 16px; color: #374151; }
.chart-container { padding: 16px 0; }
.chart-bars { display: flex; gap: 8px; align-items: flex-end; height: 150px; overflow-x: auto; }
.chart-bar-group { display: flex; flex-direction: column; align-items: center; min-width: 40px; }
.bar-container { display: flex; gap: 2px; align-items: flex-end; }
.bar { width: 10px; border-radius: 2px 2px 0 0; }
.bar.success { background: #10b981; }
.bar.error { background: #ef4444; }
.bar.blocked { background: #f59e0b; }
.bar-label { font-size: 0.7rem; color: #94a3b8; margin-top: 4px; }
.chart-legend { display: flex; gap: 16px; margin-top: 16px; justify-content: center; }
.legend-item { display: flex; align-items: center; gap: 6px; font-size: 0.85rem; color: #64748b; }
.dot { width: 10px; height: 10px; border-radius: 2px; }
.dot.success { background: #10b981; }
.dot.error { background: #ef4444; }
.dot.blocked { background: #f59e0b; }

.summary-section { background: white; padding: 20px; border-radius: 12px; }
.summary-section h3 { margin: 0 0 16px; color: #374151; }
.summary-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 16px; }
.summary-item { text-align: center; padding: 16px; background: #f8fafc; border-radius: 8px; }
.summary-value { display: block; font-size: 1.5rem; font-weight: 700; color: #1e293b; }
.summary-label { color: #64748b; font-size: 0.9rem; }
</style>
