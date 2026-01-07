<template>
  <div class="tenant-users-page">
    <div class="header">
      <h2>租户成员管理</h2>
      <button @click="showAdd = true" class="btn-primary">+ 添加成员</button>
    </div>

    <table class="data-table">
      <thead>
        <tr>
          <th>用户名</th>
          <th>显示名</th>
          <th>邮箱</th>
          <th>角色</th>
          <th>加入时间</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="tu in tenantUsers" :key="tu.id">
          <td>{{ tu.username }}</td>
          <td>{{ tu.display_name || '-' }}</td>
          <td>{{ tu.email || '-' }}</td>
          <td>
            <select :value="tu.role" @change="updateRole(tu, ($event.target as HTMLSelectElement).value)">
              <option value="tenant_admin">租户管理员</option>
              <option value="tenant_user">普通用户</option>
            </select>
          </td>
          <td>{{ formatDate(tu.created_at) }}</td>
          <td>
            <button @click="removeUser(tu)" class="btn-sm btn-danger">移除</button>
          </td>
        </tr>
      </tbody>
    </table>

    <!-- Add Modal -->
    <div v-if="showAdd" class="modal-overlay" @click.self="showAdd = false">
      <div class="modal">
        <h3>添加成员</h3>
        <form @submit.prevent="addUser">
          <select v-model="addForm.user_id" required>
            <option value="">选择用户...</option>
            <option v-for="u in availableUsers" :key="u.id" :value="u.id">{{ u.username }}</option>
          </select>
          <select v-model="addForm.role">
            <option value="tenant_user">普通用户</option>
            <option value="tenant_admin">租户管理员</option>
          </select>
          <div class="modal-actions">
            <button type="button" @click="showAdd = false">取消</button>
            <button type="submit" class="btn-primary">添加</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { api } from '../services/api'

const route = useRoute()
const tenantId = route.params.tenantId as string

const tenantUsers = ref<any[]>([])
const availableUsers = ref<any[]>([])
const showAdd = ref(false)
const addForm = ref({ user_id: '', role: 'tenant_user' })

onMounted(async () => {
  await loadTenantUsers()
  await loadAvailableUsers()
})

async function loadTenantUsers() {
  tenantUsers.value = await api.get(`/tenants/${tenantId}/users`)
}

async function loadAvailableUsers() {
  availableUsers.value = await api.get('/users')
}

async function addUser() {
  await api.post(`/tenants/${tenantId}/users`, addForm.value)
  showAdd.value = false
  addForm.value = { user_id: '', role: 'tenant_user' }
  await loadTenantUsers()
}

async function updateRole(tu: any, role: string) {
  await api.put(`/tenants/${tenantId}/users/${tu.user_id}/role`, { role })
  await loadTenantUsers()
}

async function removeUser(tu: any) {
  if (confirm(`确定移除成员 ${tu.username}?`)) {
    await api.delete(`/tenants/${tenantId}/users/${tu.user_id}`)
    await loadTenantUsers()
  }
}

function formatDate(d: string) {
  return d ? new Date(d).toLocaleString() : '-'
}
</script>

<style scoped>
.tenant-users-page { padding: 20px; }
.header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th, .data-table td { padding: 10px; text-align: left; border-bottom: 1px solid #e5e7eb; }
.data-table th { background: #f9fafb; font-weight: 600; }
.data-table select { padding: 4px 8px; border: 1px solid #d1d5db; border-radius: 4px; }
.btn-primary { background: #2563eb; color: white; padding: 8px 16px; border: none; border-radius: 4px; cursor: pointer; }
.btn-sm { padding: 4px 8px; border: 1px solid #d1d5db; border-radius: 4px; background: white; cursor: pointer; }
.btn-danger { color: #dc2626; border-color: #dc2626; }
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; }
.modal { background: white; padding: 24px; border-radius: 8px; min-width: 320px; }
.modal h3 { margin-bottom: 16px; }
.modal form { display: flex; flex-direction: column; gap: 12px; }
.modal select { padding: 8px; border: 1px solid #d1d5db; border-radius: 4px; }
.modal-actions { display: flex; justify-content: flex-end; gap: 8px; margin-top: 8px; }
</style>
