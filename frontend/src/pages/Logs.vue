<template>
  <div class="page-container">
    <div class="page-header">
      <h2>审计日志</h2>
      <div class="filters">
        <input v-model="filters.tenantId" placeholder="租户ID" class="search-input" />
        <input v-model="filters.event" placeholder="事件类型" class="search-input" />
        <button @click="load" class="btn-primary" :disabled="loading">
          {{ loading ? '加载中...' : '刷新' }}
        </button>
      </div>
    </div>

    <div class="table-container">
      <table class="data-table">
        <thead>
          <tr>
            <th>时间</th>
            <th>事件</th>
            <th>租户</th>
            <th>用户/主体</th>
            <th>详情</th>
          </tr>
        </thead>
        <tbody>
          <template v-for="(log, idx) in logs" :key="idx">
            <tr class="log-row" @click="toggleDetail(idx)">
              <td class="whitespace-nowrap">{{ new Date(log.created_at).toLocaleString() }}</td>
              <td>
                <span :class="['badge', getEventColor(log.event)]">{{ log.event }}</span>
              </td>
              <td class="font-mono text-sm">{{ log.tenant_id || '-' }}</td>
              <td>{{ log.user_id || log.actor || '-' }}</td>
              <td>
                <button class="btn-xs btn-link">
                  {{ expandedRows.has(idx) ? '▼ 收起' : '▶ 查看' }}
                </button>
              </td>
            </tr>
            <tr v-if="expandedRows.has(idx)" class="detail-row">
              <td colspan="5">
                <pre class="json-viewer">{{ JSON.stringify(log, null, 2) }}</pre>
              </td>
            </tr>
          </template>
          <tr v-if="logs.length === 0">
            <td colspan="5" class="empty-state">暂无日志</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { api } from '../services/api'

const logs = ref<any[]>([])
const loading = ref(false)
const expandedRows = ref(new Set<number>())

const filters = reactive({
  tenantId: '',
  event: '',
  limit: 100
})

async function load() {
  loading.value = true
  try {
    const result = await api.listAudit(filters.limit, filters.event || undefined, filters.tenantId || undefined)
    logs.value = Array.isArray(result) ? result : []
  } catch (e) {
    console.error(e)
    logs.value = []
  } finally {
    loading.value = false
  }
}

function toggleDetail(idx: number) {
  if (expandedRows.value.has(idx)) {
    expandedRows.value.delete(idx)
  } else {
    expandedRows.value.add(idx)
  }
  // Force reactivity
  expandedRows.value = new Set(expandedRows.value)
}

function getEventColor(event: string) {
  if (event.includes('create')) return 'badge-success'
  if (event.includes('delete') || event.includes('revoke')) return 'badge-danger'
  if (event.includes('update') || event.includes('rotate')) return 'badge-warning'
  return 'badge-info'
}

onMounted(load)
</script>

<style scoped>
.page-container { padding: 24px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
.page-header h2 { margin: 0; font-size: 1.5rem; color: #1e293b; }

.filters { display: flex; gap: 12px; }
.search-input {
  padding: 8px 12px;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  width: 200px;
}

.table-container { background: white; border-radius: 8px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); overflow: hidden; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th, .data-table td { padding: 12px 16px; text-align: left; border-bottom: 1px solid #e2e8f0; }
.data-table th { background: #f8fafc; font-weight: 600; color: #64748b; font-size: 0.875rem; }
.log-row { cursor: pointer; }
.log-row:hover { background: #f8fafc; }

.detail-row td { background: #f1f5f9; padding: 0; }
.json-viewer {
  margin: 0;
  padding: 16px;
  background: #1e293b;
  color: #a5b4fc;
  font-family: monospace;
  font-size: 0.85rem;
  overflow-x: auto;
  white-space: pre-wrap;
  word-break: break-all;
}

.badge { padding: 2px 8px; border-radius: 4px; font-size: 0.75rem; font-weight: 500; }
.badge-success { background: #dcfce7; color: #166534; }
.badge-danger { background: #fee2e2; color: #991b1b; }
.badge-warning { background: #fef9c3; color: #854d0e; }
.badge-info { background: #e0f2fe; color: #075985; }

.btn-primary { background: #2563eb; color: white; border: none; padding: 8px 16px; border-radius: 6px; cursor: pointer; }
.btn-primary:disabled { opacity: 0.7; }
.btn-link { background: none; border: none; color: #2563eb; cursor: pointer; padding: 0; font-size: 0.85rem; }
.btn-xs { font-size: 0.8rem; }

.font-mono { font-family: monospace; color: #64748b; }
.text-sm { font-size: 0.875rem; }
.whitespace-nowrap { white-space: nowrap; }
.empty-state { text-align: center; color: #94a3b8; padding: 40px; }
</style>
