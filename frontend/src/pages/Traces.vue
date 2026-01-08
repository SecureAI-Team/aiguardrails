<template>
  <div class="traces-page">
    <div class="page-header">
      <h2>🔍 请求追踪</h2>
      <div class="header-actions">
        <select v-model="filterBlocked" @change="loadTraces">
          <option value="">全部请求</option>
          <option value="true">仅阻断</option>
          <option value="false">仅通过</option>
        </select>
        <button @click="loadTraces" class="btn-outline">🔄 刷新</button>
      </div>
    </div>

    <div class="traces-list">
      <div v-for="trace in traces" :key="trace.id" class="trace-item" :class="{ blocked: trace.blocked, error: trace.status_code >= 400 }">
        <div class="trace-header">
          <span class="trace-method" :class="trace.method.toLowerCase()">{{ trace.method }}</span>
          <span class="trace-path">{{ trace.path }}</span>
          <span class="trace-status" :class="getStatusClass(trace.status_code)">{{ trace.status_code }}</span>
        </div>
        <div class="trace-meta">
          <span>⏱️ {{ trace.duration_ms }}ms</span>
          <span>🪙 {{ trace.input_tokens + trace.output_tokens }} tokens</span>
          <span>📅 {{ formatTime(trace.created_at) }}</span>
          <span v-if="trace.blocked" class="blocked-badge">🛡️ {{ trace.block_reason }}</span>
        </div>
        <div v-if="trace.signals?.length" class="trace-signals">
          <span v-for="sig in trace.signals" :key="sig" class="signal-tag">{{ sig }}</span>
        </div>
        <button @click="viewDetail(trace)" class="btn-link">查看详情 →</button>
      </div>
      <div v-if="traces.length === 0" class="no-traces">暂无请求记录</div>
    </div>

    <!-- Detail Modal -->
    <div v-if="selectedTrace" class="modal-overlay" @click.self="selectedTrace = null">
      <div class="modal detail-modal">
        <h3>请求详情</h3>
        <div class="detail-section">
          <h4>基本信息</h4>
          <div class="detail-row"><span>Trace ID:</span><code>{{ selectedTrace.trace_id }}</code></div>
          <div class="detail-row"><span>状态:</span><span :class="getStatusClass(selectedTrace.status_code)">{{ selectedTrace.status_code }}</span></div>
          <div class="detail-row"><span>耗时:</span>{{ selectedTrace.duration_ms }}ms</div>
          <div class="detail-row"><span>客户端IP:</span>{{ selectedTrace.client_ip }}</div>
        </div>
        <div v-if="selectedTrace.stages" class="detail-section">
          <h4>处理阶段</h4>
          <div v-for="stage in parseStages(selectedTrace.stages)" :key="stage.name" class="stage-item">
            <span class="stage-name">{{ stage.name }}</span>
            <span class="stage-duration">{{ stage.end - stage.start }}ms</span>
            <span :class="'stage-status ' + stage.status">{{ stage.status }}</span>
          </div>
        </div>
        <div v-if="selectedTrace.error" class="detail-section error-section">
          <h4>错误信息</h4>
          <pre>{{ selectedTrace.error }}</pre>
        </div>
        <div class="modal-actions">
          <button @click="selectedTrace = null" class="btn-primary">关闭</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '../services/api'

interface Trace {
  id: string
  trace_id: string
  method: string
  path: string
  status_code: number
  duration_ms: number
  blocked: boolean
  block_reason?: string
  signals?: string[]
  input_tokens: number
  output_tokens: number
  stages?: string
  error?: string
  client_ip?: string
  created_at: string
}

const traces = ref<Trace[]>([])
const filterBlocked = ref('')
const selectedTrace = ref<Trace | null>(null)

onMounted(() => loadTraces())

async function loadTraces() {
  try {
    const query = filterBlocked.value ? `?blocked=${filterBlocked.value}` : ''
    const result = await api.get(`/traces${query}`)
    traces.value = Array.isArray(result) ? result : []
  } catch { 
    traces.value = [] 
  }
}

function getStatusClass(code: number) {
  if (code >= 500) return 'status-5xx'
  if (code >= 400) return 'status-4xx'
  if (code >= 300) return 'status-3xx'
  return 'status-2xx'
}

function formatTime(ts: string) {
  return new Date(ts).toLocaleString('zh-CN')
}

function parseStages(stages: string) {
  try { return JSON.parse(stages) } catch { return [] }
}

async function viewDetail(trace: Trace) {
  try {
    selectedTrace.value = await api.get(`/traces/${trace.id}`)
  } catch { selectedTrace.value = trace }
}
</script>

<style scoped>
.traces-page { padding: 20px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
.page-header h2 { margin: 0; }
.header-actions { display: flex; gap: 12px; }
.header-actions select { padding: 8px 16px; border: 1px solid #e2e8f0; border-radius: 6px; }
.btn-outline { padding: 8px 16px; border: 1px solid #e2e8f0; background: white; border-radius: 6px; cursor: pointer; }
.traces-list { display: flex; flex-direction: column; gap: 12px; }
.trace-item { background: white; padding: 16px; border-radius: 12px; box-shadow: 0 2px 10px rgba(0,0,0,0.05); border-left: 4px solid #10b981; }
.trace-item.blocked { border-left-color: #f59e0b; }
.trace-item.error { border-left-color: #ef4444; }
.trace-header { display: flex; align-items: center; gap: 12px; margin-bottom: 8px; }
.trace-method { padding: 4px 12px; border-radius: 4px; font-weight: 600; font-size: 0.8rem; }
.trace-method.get { background: #dbeafe; color: #1e40af; }
.trace-method.post { background: #d1fae5; color: #065f46; }
.trace-method.put { background: #fef3c7; color: #92400e; }
.trace-method.delete { background: #fee2e2; color: #991b1b; }
.trace-path { font-family: monospace; color: #334155; flex: 1; }
.trace-status { font-weight: 600; }
.status-2xx { color: #10b981; }
.status-3xx { color: #3b82f6; }
.status-4xx { color: #f59e0b; }
.status-5xx { color: #ef4444; }
.trace-meta { display: flex; gap: 16px; color: #64748b; font-size: 0.85rem; margin-bottom: 8px; }
.blocked-badge { color: #f59e0b; }
.trace-signals { display: flex; gap: 8px; margin-bottom: 8px; }
.signal-tag { background: #fee2e2; color: #991b1b; padding: 2px 8px; border-radius: 4px; font-size: 0.75rem; }
.btn-link { background: none; border: none; color: #3b82f6; cursor: pointer; padding: 0; }
.no-traces { text-align: center; padding: 60px; color: #64748b; }
.modal-overlay { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 1000; }
.modal { background: white; padding: 24px; border-radius: 12px; max-height: 80vh; overflow-y: auto; }
.detail-modal { width: 600px; }
.detail-section { margin-bottom: 20px; }
.detail-section h4 { margin: 0 0 12px; color: #374151; }
.detail-row { display: flex; justify-content: space-between; padding: 8px 0; border-bottom: 1px solid #f1f5f9; }
.detail-row code { font-family: monospace; font-size: 0.85rem; }
.stage-item { display: flex; justify-content: space-between; padding: 8px; background: #f8fafc; margin-bottom: 4px; border-radius: 4px; }
.stage-name { font-weight: 500; }
.stage-status.success { color: #10b981; }
.stage-status.error { color: #ef4444; }
.error-section pre { background: #fef2f2; padding: 12px; border-radius: 6px; color: #991b1b; overflow-x: auto; }
.modal-actions { display: flex; justify-content: flex-end; margin-top: 24px; }
.btn-primary { background: linear-gradient(90deg, #2563eb, #7c3aed); color: white; padding: 10px 20px; border: none; border-radius: 8px; cursor: pointer; }
</style>
