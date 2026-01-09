<template>
  <div class="page-container">
    <div class="page-header">
      <h2>租户管理</h2>
      <button @click="showCreateModal = true" class="btn-primary">+ 新建租户</button>
    </div>

    <div class="table-container">
      <table class="data-table">
        <thead>
          <tr>
            <th>租户名称</th>
            <th>ID</th>
            <th>创建时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="t in tenants" :key="t.id">
            <td class="font-bold">{{ t.name }}</td>
            <td class="font-mono text-sm">{{ t.id }}</td>
            <td>{{ new Date(t.created_at || Date.now()).toLocaleString() }}</td>
            <td class="actions">
              <router-link :to="{ path: '/apps', query: { tenantId: t.id } }" class="btn-sm btn-outline">应用管理</router-link>
              <router-link :to="{ path: '/policies', query: { tenantId: t.id } }" class="btn-sm btn-outline">策略配置</router-link>
              <router-link :to="`/tenants/${t.id}/users`" class="btn-sm btn-outline">用户管理</router-link>
            </td>
          </tr>
          <tr v-if="tenants.length === 0">
            <td colspan="4" class="empty-state">暂无租户，请点击右上角创建</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Create Modal -->
    <div v-if="showCreateModal" class="modal-overlay" @click.self="showCreateModal = false">
      <div class="modal">
        <h3>新建租户</h3>
        <form @submit.prevent="onCreate">
          <div class="form-group">
            <label>租户名称</label>
            <input v-model="form.name" placeholder="请输入租户名称" required :disabled="loading" />
          </div>
          <div class="modal-actions">
            <button type="button" @click="showCreateModal = false" class="btn-secondary">取消</button>
            <button type="submit" class="btn-primary" :disabled="loading">
              {{ loading ? '创建中...' : '创建' }}
            </button>
          </div>
        </form>
      </div>
    </div>
    <!-- Alert Modal -->
    <AlertModal
      :is-open="showAlertModal"
      :title="alertTitle"
      :message="alertMessage"
      :type="alertType"
      @close="showAlertModal = false"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, reactive } from 'vue'
import { api } from '../services/api'
import AlertModal from '../components/AlertModal.vue'

const tenants = ref<any[]>([])
const loading = ref(false)
const showCreateModal = ref(false)
const form = reactive({
  name: ''
})

const showAlertModal = ref(false)
const alertTitle = ref('')
const alertMessage = ref('')
const alertType = ref('info')

function showAlert(msg: string, type = 'info', title = '提示') {
  alertMessage.value = msg
  alertType.value = type
  alertTitle.value = title
  showAlertModal.value = true
}

async function load() {
  try {
    const result = await api.listTenants()
    tenants.value = Array.isArray(result) ? result : []
  } catch (e) {
    console.error(e)
    tenants.value = []
  }
}

async function onCreate() {
  if (!form.name) return
  loading.value = true
  try {
    await api.createTenant(form.name)
    form.name = ''
    showCreateModal.value = false
    await load()
  } catch (e) {
    showAlert('创建失败', 'error')
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<style scoped>
/* ... existing styles ... */
.page-container {
  padding: 24px;
}
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}
.page-header h2 {
  margin: 0;
  font-size: 1.5rem;
  color: #1e293b;
}

/* Table Styles */
.table-container {
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
  overflow: hidden;
}
.data-table {
  width: 100%;
  border-collapse: collapse;
}
.data-table th, .data-table td {
  padding: 16px;
  text-align: left;
  border-bottom: 1px solid #e2e8f0;
}
.data-table th {
  background: #f8fafc;
  font-weight: 600;
  color: #64748b;
  font-size: 0.875rem;
}
.data-table tr:hover {
  background: #f8fafc;
}
.font-bold { font-weight: 600; color: #1e293b; }
.font-mono { font-family: monospace; color: #64748b; }
.text-sm { font-size: 0.875rem; }
.empty-state { text-align: center; padding: 40px; color: #94a3b8; }
.actions { display: flex; gap: 8px; }

/* Buttons */
.btn-primary {
  background: #2563eb;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 500;
}
.btn-primary:hover { background: #1d4ed8; }
.btn-primary:disabled { opacity: 0.7; cursor: not-allowed; }

.btn-secondary {
  background: white;
  border: 1px solid #cbd5e1;
  color: #475569;
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
}
.btn-secondary:hover { background: #f1f5f9; }

.btn-sm {
  padding: 4px 8px;
  font-size: 0.8rem;
  text-decoration: none;
  border-radius: 4px;
  cursor: pointer;
}
.btn-outline {
  border: 1px solid #cbd5e1;
  color: #475569;
  background: white;
}
.btn-outline:hover {
  border-color: #2563eb;
  color: #2563eb;
}

/* Modal */
.modal-overlay {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}
.modal {
  background: white;
  padding: 24px;
  border-radius: 12px;
  width: 400px;
  box-shadow: 0 20px 25px -5px rgba(0,0,0,0.1);
}
.modal h3 { margin: 0 0 20px; font-size: 1.25rem; }
.form-group { margin-bottom: 20px; }
.form-group label { display: block; margin-bottom: 8px; font-weight: 500; font-size: 0.9rem; }
.form-group input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  font-size: 1rem;
}
.modal-actions { display: flex; justify-content: flex-end; gap: 12px; }
</style>
