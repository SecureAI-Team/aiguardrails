<template>
  <div class="page-container">
    <div class="page-header">
      <h2>策略管理</h2>
      <div class="header-actions">
        <router-link to="/policy-history" class="btn-secondary">📜 变更历史</router-link>
        <router-link to="/capabilities" class="btn-secondary">🧩 管理能力</router-link>
        <select v-model="selectedTenantId" @change="onTenantChange" class="tenant-select">
          <option value="" disabled>请选择租户</option>
          <option v-for="t in tenants" :key="t.id" :value="t.id">{{ t.name }}</option>
        </select>
        <button @click="openCreateModal" class="btn-primary" :disabled="!selectedTenantId">
          + 新建策略
        </button>
      </div>
    </div>

    <div class="table-container" v-if="selectedTenantId">
      <table class="data-table">
        <thead>
          <tr>
            <th>策略名称</th>
            <th>提示词规则</th>
            <th>RAG 命名空间</th>
            <th>敏感词库</th>
            <th>创建时间</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="p in policies" :key="p.id">
            <td class="font-bold">{{ p.name }}</td>
            <td>
              <div class="tags">
                <span v-for="r in (p.promptRules || [])" :key="r" class="badge badge-info">{{ r }}</span>
              </div>
            </td>
            <td>
              <div class="tags">
                <span v-for="n in (p.ragNamespaces || [])" :key="n" class="badge badge-warning">{{ n }}</span>
              </div>
            </td>
            <td>{{ (p.sensitiveTerms || []).length }} 个词条</td>
            <td class="text-sm font-mono">{{ new Date().toLocaleDateString() }}</td> <!-- API may not return created_at -->
          </tr>
          <tr v-if="policies.length === 0">
            <td colspan="5" class="empty-state">该租户下暂无策略</td>
          </tr>
        </tbody>
      </table>
    </div>
    
    <div v-else class="empty-state-large">
      <div class="icon">👈</div>
      <h3>请先选择一个租户</h3>
      <p>选择租户以管理其下的安全策略</p>
    </div>

    <!-- Create Modal -->
    <div v-if="showCreateModal" class="modal-overlay" @click.self="showCreateModal = false">
      <div class="modal">
        <h3>新建策略</h3>
        <form @submit.prevent="onCreate">
          <div class="form-group">
            <label>策略名称</label>
            <input v-model="form.name" placeholder="请输入策略名称" required />
          </div>
          
          <div class="form-group">
            <div class="rules-list">
              <label v-for="r in availableRules" :key="r.id" class="checkbox-label" :title="r.description">
                <input type="checkbox" :value="r.id" v-model="form.selectedRules" />
                <span class="rule-name">{{ r.name }}</span>
                <span class="badge badge-xs" :class="getSeverityClass(r.severity)">{{ r.severity }}</span>
              </label>
            </div>
            <div v-if="availableRules.length === 0" class="muted">暂无加载规则或规则库为空</div>
          </div>

          <div class="form-group">
            <label>RAG 命名空间</label>
            <input v-model="form.ragNamespaces" placeholder="doc-1, kb-sales" />
          </div>

          <div class="form-group">
            <label>敏感词库</label>
            <input v-model="form.sensitiveTerms" placeholder="机密, 内部, 薪资" />
          </div>

          <div class="form-group">
            <label>工具能力 (Capabilities)</label>
            <div class="caps-list">
              <label v-for="c in caps" :key="c.id || c.name" class="checkbox-label">
                <input type="checkbox" :value="c.name" v-model="form.selectedCaps" />
                <span>{{ c.name }}</span>
                <span class="muted" v-if="c.tags.length">({{ c.tags.join(', ') }})</span>
              </label>
            </div>
            <div v-if="caps.length === 0" class="muted">暂无可用能力</div>
            <button type="button" @click="loadCaps" class="btn-xs btn-link">刷新能力列表</button>
          </div>

          <div class="modal-actions">
            <button type="button" @click="showCreateModal = false" class="btn-secondary">取消</button>
            <button type="submit" class="btn-primary" :disabled="loading">创建</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api } from '../services/api'

const route = useRoute()
const router = useRouter()

const tenants = ref<any[]>([])
const selectedTenantId = ref('')
const policies = ref<any[]>([])
const loading = ref(false)
const showCreateModal = ref(false)
const caps = ref<any[]>([])
const availableRules = ref<any[]>([])

const form = reactive({
  name: '',
  selectedRules: [] as string[],
  ragNamespaces: '',
  sensitiveTerms: '',
  selectedCaps: [] as string[]
})

async function loadTenants() {
  try {
    const result = await api.listTenants()
    tenants.value = Array.isArray(result) ? result : []
    if (route.query.tenantId) {
      selectedTenantId.value = route.query.tenantId as string
      loadPolicies()
    }
  } catch (e) {
    console.error(e)
    tenants.value = []
  }
}

async function loadPolicies() {
  if (!selectedTenantId.value) return
  loading.value = true
  try {
    const result = await api.listPolicies(selectedTenantId.value)
    policies.value = Array.isArray(result) ? result : []
  } catch (e) {
    console.error(e)
    policies.value = []
  } finally {
    loading.value = false
  }
}

async function loadCaps() {
  try {
    const result = await api.listCapabilities()
    caps.value = Array.isArray(result) ? result : []
  } catch (e) {
    console.error(e)
    caps.value = []
  }
}

async function loadRules() {
  try {
    const result = await api.listRules({})
    availableRules.value = Array.isArray(result) ? result : []
  } catch (e) {
    console.error(e)
    availableRules.value = []
  }
}

function onTenantChange() {
  router.push({ query: { tenantId: selectedTenantId.value } })
  loadPolicies()
}

function openCreateModal() {
  form.name = ''
  form.selectedRules = []
  form.ragNamespaces = ''
  form.sensitiveTerms = ''
  form.selectedCaps = []
  showCreateModal.value = true
  if (caps.value.length === 0) loadCaps()
  if (availableRules.value.length === 0) loadRules()
}

async function onCreate() {
  if (!form.name) return
  loading.value = true
  try {
    await api.createPolicy(selectedTenantId.value, {
      name: form.name,
      prompt_rules: form.selectedRules,
      rag_namespaces: form.ragNamespaces.split(',').map(s => s.trim()).filter(Boolean),
      sensitive_terms: form.sensitiveTerms.split(',').map(s => s.trim()).filter(Boolean),
      tool_allowlist: form.selectedCaps
    })
    showCreateModal.value = false
    await loadPolicies()
  } catch (e: any) {
    alert('创建失败: ' + (e.response?.data?.error || e.message))
  } finally {
    loading.value = false
  }
}

function getSeverityClass(severity: string) {
  switch (severity) {
    case 'critical': return 'badge-danger'
    case 'high': return 'badge-danger'
    case 'medium': return 'badge-warning'
    case 'low': return 'badge-info'
    default: return 'badge-secondary'
  }
}

onMounted(() => {
  loadTenants()
})
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

.table-container { background: white; border-radius: 8px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); overflow: hidden; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th, .data-table td { padding: 16px; text-align: left; border-bottom: 1px solid #e2e8f0; }
.data-table th { background: #f8fafc; font-weight: 600; color: #64748b; font-size: 0.875rem; }
.data-table tr:hover { background: #f8fafc; }

.font-bold { font-weight: 600; color: #1e293b; }
.text-sm { font-size: 0.875rem; }
.font-mono { font-family: monospace; color: #64748b; }

.badge { padding: 2px 6px; border-radius: 4px; font-size: 0.75rem; margin-right: 4px; display: inline-block; margin-bottom: 2px; }
.badge-info { background: #e0f2fe; color: #075985; }
.badge-warning { background: #fef9c3; color: #854d0e; }
.badge-danger { background: #fee2e2; color: #991b1b; }
.badge-secondary { background: #f1f5f9; color: #64748b; }
.badge-xs { font-size: 0.7rem; padding: 1px 4px; }

.btn-primary { background: #2563eb; color: white; border: none; padding: 8px 16px; border-radius: 6px; cursor: pointer; text-decoration: none; display: inline-flex; align-items: center; }
.btn-secondary { background: white; border: 1px solid #cbd5e1; color: #475569; padding: 8px 16px; border-radius: 6px; cursor: pointer; text-decoration: none; display: inline-flex; align-items: center; }
.btn-xs { padding: 2px 6px; font-size: 0.75rem; }
.btn-link { background: none; border: none; color: #2563eb; cursor: pointer; text-decoration: underline; }

.empty-state { text-align: center; padding: 40px; color: #94a3b8; }
.empty-state-large { text-align: center; padding: 80px 0; color: #64748b; }
.empty-state-large .icon { font-size: 3rem; margin-bottom: 16px; }

/* Modal */
.modal-overlay { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 1000; }
.modal { background: white; padding: 24px; border-radius: 12px; width: 600px; max-height: 90vh; overflow-y: auto; box-shadow: 0 20px 25px -5px rgba(0,0,0,0.1); }
.modal h3 { margin: 0 0 20px; font-size: 1.25rem; }

.form-group { margin-bottom: 20px; }
.form-group label { display: block; margin-bottom: 8px; font-weight: 500; }
.form-group input, .form-group textarea { width: 100%; padding: 8px; border: 1px solid #cbd5e1; border-radius: 6px; font-family: inherit; }
.form-group textarea { height: 80px; resize: vertical; }
.helper-text { font-size: 0.8rem; color: #64748b; margin-top: 4px; }
.modal-actions { display: flex; justify-content: flex-end; gap: 12px; }

.caps-list, .rules-list { display: grid; grid-template-columns: repeat(2, 1fr); gap: 8px; max-height: 200px; overflow-y: auto; border: 1px solid #e2e8f0; padding: 8px; border-radius: 6px; }
.checkbox-label { display: flex; gap: 8px; align-items: center; font-size: 0.9rem; margin-bottom: 0; cursor: pointer; }
.checkbox-label:hover { background: #f8fafc; }
.rule-name { flex: 1; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.muted { color: #94a3b8; }
</style>
