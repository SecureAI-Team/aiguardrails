<template>
  <div class="page-container">
    <div class="page-header">
      <div class="header-left">
        <router-link to="/tenants" class="back-link">← 返回租户列表</router-link>
        <h2>租户成员管理</h2>
      </div>
      <button @click="openAddModal" class="btn-primary">+ 添加成员</button>
    </div>

    <div class="table-container">
      <table class="data-table">
        <thead>
          <tr>
            <th>用户名</th>
            <th>显示名</th>
            <th>角色</th>
            <th>加入时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="tu in tenantUsers" :key="tu.id">
            <td class="font-bold">{{ tu.username }}</td>
            <td>
              <div>{{ tu.display_name || '-' }}</div>
              <div v-if="tu.email" class="muted text-sm">{{ tu.email }}</div>
            </td>
            <td>
              <select 
                :value="tu.role" 
                @change="updateRole(tu, ($event.target as HTMLSelectElement).value)"
                class="role-select"
              >
                <option value="tenant_admin">租户管理员</option>
                <option value="tenant_user">普通用户</option>
              </select>
            </td>
            <td class="text-sm font-mono">{{ formatDate(tu.created_at) }}</td>
            <td>
              <button @click="removeUser(tu)" class="btn-xs btn-danger">移除</button>
            </td>
          </tr>
          <tr v-if="tenantUsers.length === 0">
            <td colspan="5" class="empty-state">该租户暂无其他成员</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Add Modal -->
    <div v-if="showAdd" class="modal-overlay" @click.self="showAdd = false">
      <div class="modal">
        <h3>添加成员</h3>
        <p class="modal-desc">将现有平台用户添加到此租户。</p>
        <form @submit.prevent="addUser">
          <div class="form-group">
            <label>选择用户</label>
            <select v-model="addForm.user_id" required class="form-select">
              <option value="" disabled>请选择用户...</option>
              <option v-for="u in availableUsers" :key="u.id" :value="u.id">
                {{ u.username }} <span v-if="u.display_name">({{ u.display_name }})</span>
              </option>
            </select>
            <div class="helper-text">仅显示尚未加入此租户的用户</div>
          </div>
          <div class="form-group">
            <label>分配角色</label>
            <select v-model="addForm.role" class="form-select">
              <option value="tenant_user">普通用户</option>
              <option value="tenant_admin">租户管理员</option>
            </select>
          </div>
          <div class="modal-actions">
            <button type="button" @click="showAdd = false" class="btn-secondary">取消</button>
            <button type="submit" class="btn-primary" :disabled="loading">添加</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, reactive } from 'vue'
import { useRoute } from 'vue-router'
import { api } from '../services/api'

const route = useRoute()
const tenantId = route.params.tenantId as string

const tenantUsers = ref<any[]>([])
const allUsers = ref<any[]>([])
const loading = ref(false)
const showAdd = ref(false)

const addForm = reactive({ user_id: '', role: 'tenant_user' })

// Filter out users already in tenant
const availableUsers = computed(() => {
  const existingIds = new Set(tenantUsers.value.map(tu => tu.user_id))
  return allUsers.value.filter(u => !existingIds.has(u.id))
})

async function loadData() {
  loading.value = true
  try {
    const [usersRes, tenantUsersRes] = await Promise.all([
      api.get('/users'),
      api.get(`/tenants/${tenantId}/users`)
    ])
    allUsers.value = usersRes
    tenantUsers.value = tenantUsersRes
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function openAddModal() {
  addForm.user_id = ''
  addForm.role = 'tenant_user'
  showAdd.value = true
}

async function addUser() {
  if (!addForm.user_id) return
  loading.value = true
  try {
    await api.post(`/tenants/${tenantId}/users`, addForm)
    showAdd.value = false
    await loadData()
  } catch (e) {
    alert('添加失败')
  } finally {
    loading.value = false
  }
}

async function updateRole(tu: any, role: string) {
  try {
    await api.put(`/tenants/${tenantId}/users/${tu.user_id}/role`, { role })
    await loadData()
  } catch (e) {
    alert('更新失败')
  }
}

async function removeUser(tu: any) {
  if (!confirm(`确定要移除成员 "${tu.username}" 吗？他将无法访问此租户。`)) return
  try {
    await api.delete(`/tenants/${tenantId}/users/${tu.user_id}`)
    await loadData()
  } catch (e) {
    alert('移除失败')
  }
}

function formatDate(d?: string) {
  if (!d) return '-'
  return new Date(d).toLocaleString()
}

onMounted(loadData)
</script>

<style scoped>
.page-container { padding: 24px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
.header-left h2 { margin: 8px 0 0; font-size: 1.5rem; color: #1e293b; }
.back-link { color: #64748b; text-decoration: none; font-size: 0.9rem; }
.back-link:hover { color: #2563eb; }

.table-container { background: white; border-radius: 8px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); overflow: hidden; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th, .data-table td { padding: 16px; text-align: left; border-bottom: 1px solid #e2e8f0; }
.data-table th { background: #f8fafc; font-weight: 600; color: #64748b; font-size: 0.875rem; }
.data-table tr:hover { background: #f8fafc; }

.font-bold { font-weight: 600; color: #1e293b; }
.muted { color: #64748b; }
.text-sm { font-size: 0.875rem; }
.font-mono { font-family: monospace; color: #64748b; }

.role-select { padding: 4px 8px; border: 1px solid #cbd5e1; border-radius: 4px; background: transparent; cursor: pointer; }
.role-select:hover { border-color: #2563eb; }

.btn-primary { background: #2563eb; color: white; border: none; padding: 8px 16px; border-radius: 6px; cursor: pointer; }
.btn-secondary { background: white; border: 1px solid #cbd5e1; color: #475569; padding: 8px 16px; border-radius: 6px; cursor: pointer; }
.btn-xs { padding: 4px 10px; font-size: 0.75rem; border: 1px solid #cbd5e1; border-radius: 4px; background: white; cursor: pointer; }
.btn-danger { color: #ef4444; border-color: #fca5a5; }
.btn-danger:hover { background: #fee2e2; }

/* Modal */
.modal-overlay { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 1000; }
.modal { background: white; padding: 24px; border-radius: 12px; width: 400px; box-shadow: 0 20px 25px -5px rgba(0,0,0,0.1); }
.modal h3 { margin: 0 0 8px; font-size: 1.25rem; }
.modal-desc { color: #64748b; font-size: 0.9rem; margin-bottom: 20px; }

.form-group { margin-bottom: 20px; }
.form-group label { display: block; margin-bottom: 8px; font-weight: 500; }
.form-group select { width: 100%; padding: 8px 12px; border: 1px solid #cbd5e1; border-radius: 6px; font-size: 0.95rem; }
.helper-text { font-size: 0.8rem; color: #64748b; margin-top: 4px; }
.modal-actions { display: flex; justify-content: flex-end; gap: 12px; }

.empty-state { text-align: center; padding: 40px; color: #94a3b8; }
</style>
