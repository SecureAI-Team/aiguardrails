<template>
  <div class="page-container">
    <div class="page-header">
      <h2>策略变更历史</h2>
      <div class="header-actions">
        <select v-model="selectedTenantId" @change="loadHistory" class="tenant-select">
          <option value="" disabled>请选择租户</option>
          <option v-for="t in tenants" :key="t.id" :value="t.id">{{ t.name }}</option>
        </select>
        <button @click="loadHistory" class="btn-secondary" :disabled="!selectedTenantId">刷新</button>
      </div>
    </div>

    <div class="timeline-container" v-if="selectedTenantId">
      <div v-for="(item, idx) in history" :key="idx" class="timeline-item">
        <div class="timeline-marker"></div>
        <div class="timeline-content card">
          <div class="item-header">
            <span class="policy-name">{{ item.name }}</span>
            <span class="update-time">{{ formatDate(item.updatedAt) }}</span>
          </div>
          <div class="item-diff">
            <div v-if="item.promptRules?.length" class="diff-section">
              <span class="label">提示词规则:</span>
              <span class="value">{{ item.promptRules.join(', ') }}</span>
            </div>
            <div v-if="item.toolAllowList?.length" class="diff-section">
              <span class="label">工具白名单:</span>
              <span class="value">{{ item.toolAllowList.join(', ') }}</span>
            </div>
            <div v-if="item.ragNamespaces?.length" class="diff-section">
              <span class="label">RAG库:</span>
              <span class="value">{{ item.ragNamespaces.join(', ') }}</span>
            </div>
            <div v-if="item.sensitiveTerms?.length" class="diff-section">
              <span class="label">敏感词:</span>
              <span class="value">{{ item.sensitiveTerms.join(', ') }}</span>
            </div>
          </div>
        </div>
      </div>
      
      <div v-if="history.length === 0" class="empty-timeline">
        该租户暂无策略变更记录
      </div>
    </div>

    <div v-else class="empty-state-large">
      <div class="icon">👈</div>
      <h3>请先选择一个租户</h3>
      <p>选择租户以查看其策略变更历史</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '../services/api'

const tenants = ref<any[]>([])
const selectedTenantId = ref('')
const history = ref<any[]>([])
const loading = ref(false)

async function loadTenants() {
  try {
    tenants.value = await api.listTenants()
  } catch (e) {
    console.error(e)
  }
}

async function loadHistory() {
  if (!selectedTenantId.value) return
  loading.value = true
  try {
    history.value = await api.listPolicyHistory(selectedTenantId.value)
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function formatDate(dateStr?: string) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString()
}

onMounted(loadTenants)
</script>

<style scoped>
.page-container { padding: 24px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
.page-header h2 { margin: 0; font-size: 1.5rem; color: #1e293b; }
.header-actions { display: flex; gap: 12px; }

.tenant-select {
  padding: 8px 12px;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  font-size: 0.95rem;
  min-width: 200px;
}

.timeline-container {
  padding-left: 20px;
  border-left: 2px solid #e2e8f0;
  margin-left: 10px;
}

.timeline-item { position: relative; margin-bottom: 24px; }
.timeline-marker {
  position: absolute;
  left: -27px;
  top: 16px;
  width: 12px;
  height: 12px;
  background: #3b82f6;
  border-radius: 50%;
  border: 2px solid white;
  box-shadow: 0 0 0 2px #3b82f6;
}

.card { background: white; padding: 20px; border-radius: 8px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }
.item-header { display: flex; justify-content: space-between; margin-bottom: 12px; padding-bottom: 12px; border-bottom: 1px solid #f1f5f9; }
.policy-name { font-weight: 600; color: #1e293b; font-size: 1.1rem; }
.update-time { color: #64748b; font-size: 0.9rem; }

.diff-section { margin-bottom: 8px; font-size: 0.95rem; }
.label { color: #64748b; margin-right: 8px; }
.value { color: #334155; }

.btn-secondary { background: white; border: 1px solid #cbd5e1; color: #475569; padding: 8px 16px; border-radius: 6px; cursor: pointer; }

.empty-timeline { color: #94a3b8; padding: 20px 0; }
.empty-state-large { text-align: center; padding: 80px 0; color: #64748b; }
.empty-state-large .icon { font-size: 3rem; margin-bottom: 16px; }
</style>
