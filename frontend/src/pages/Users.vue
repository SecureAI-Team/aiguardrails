<template>
  <div class="users-page">
    <div class="header">
      <h2>用户管理</h2>
      <button @click="showCreate = true" class="btn-primary">+ 新建用户</button>
    </div>

    <div class="filters">
      <select v-model="filterRole" @change="loadUsers">
        <option value="">所有角色</option>
        <option value="platform_admin">平台管理员</option>
        <option value="tenant_admin">租户管理员</option>
        <option value="tenant_user">普通用户</option>
      </select>
      <select v-model="filterStatus" @change="loadUsers">
        <option value="">所有状态</option>
        <option value="active">启用</option>
        <option value="inactive">禁用</option>
      </select>
    </div>

    <table class="data-table">
      <thead>
        <tr>
          <th>用户名</th>
          <th>显示名</th>
          <th>邮箱</th>
          <th>角色</th>
          <th>状态</th>
          <th>最后登录</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="user in users" :key="user.id">
          <td>{{ user.username }}</td>
          <td>{{ user.display_name || '-' }}</td>
          <td>{{ user.email || '-' }}</td>
          <td><span :class="'role-' + user.role">{{ roleLabel(user.role) }}</span></td>
          <td><span :class="'status-' + user.status">{{ user.status }}</span></td>
          <td>{{ formatDate(user.last_login_at) }}</td>
          <td>
            <button @click="editUser(user)" class="btn-sm">编辑</button>
            <button @click="resetPwd(user)" class="btn-sm">重置密码</button>
            <button @click="deleteUser(user)" class="btn-sm btn-danger">删除</button>
          </td>
        </tr>
      </tbody>
    </table>

    <!-- Create Modal -->
    <div v-if="showCreate" class="modal-overlay" @click.self="showCreate = false">
      <div class="modal">
        <h3>新建用户</h3>
        <form @submit.prevent="createUser">
          <input v-model="form.username" placeholder="用户名" required />
          <input v-model="form.password" type="password" placeholder="密码" required />
          <input v-model="form.email" type="email" placeholder="邮箱" />
          <input v-model="form.display_name" placeholder="显示名" />
          <select v-model="form.role">
            <option value="tenant_user">普通用户</option>
            <option value="tenant_admin">租户管理员</option>
            <option value="platform_admin">平台管理员</option>
          </select>
          <div class="modal-actions">
            <button type="button" @click="showCreate = false">取消</button>
            <button type="submit" class="btn-primary">创建</button>
          </div>
        </form>
      </div>
    </div>

    <!-- Edit Modal -->
    <div v-if="showEdit" class="modal-overlay" @click.self="showEdit = false">
      <div class="modal">
        <h3>编辑用户</h3>
        <form @submit.prevent="updateUser">
          <input v-model="editForm.email" type="email" placeholder="邮箱" />
          <input v-model="editForm.display_name" placeholder="显示名" />
          <select v-model="editForm.role">
            <option value="tenant_user">普通用户</option>
            <option value="tenant_admin">租户管理员</option>
            <option value="platform_admin">平台管理员</option>
          </select>
          <select v-model="editForm.status">
            <option value="active">启用</option>
            <option value="inactive">禁用</option>
          </select>
          <div class="modal-actions">
            <button type="button" @click="showEdit = false">取消</button>
            <button type="submit" class="btn-primary">保存</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '../services/api'

const users = ref<any[]>([])
const filterRole = ref('')
const filterStatus = ref('')
const showCreate = ref(false)
const showEdit = ref(false)
const form = ref({ username: '', password: '', email: '', display_name: '', role: 'tenant_user' })
const editForm = ref({ id: '', email: '', display_name: '', role: '', status: '' })

onMounted(() => loadUsers())

async function loadUsers() {
  const params = new URLSearchParams()
  if (filterRole.value) params.append('role', filterRole.value)
  if (filterStatus.value) params.append('status', filterStatus.value)
  users.value = await api.get(`/users?${params}`)
}

async function createUser() {
  await api.post('/users', form.value)
  showCreate.value = false
  form.value = { username: '', password: '', email: '', display_name: '', role: 'tenant_user' }
  await loadUsers()
}

function editUser(user: any) {
  editForm.value = { id: user.id, email: user.email, display_name: user.display_name, role: user.role, status: user.status }
  showEdit.value = true
}

async function updateUser() {
  await api.put(`/users/${editForm.value.id}`, editForm.value)
  showEdit.value = false
  await loadUsers()
}

async function resetPwd(user: any) {
  const pwd = prompt('请输入新密码:')
  if (pwd) {
    await api.post(`/users/${user.id}/password`, { new_password: pwd })
    alert('密码已重置')
  }
}

async function deleteUser(user: any) {
  if (confirm(`确定删除用户 ${user.username}?`)) {
    await api.delete(`/users/${user.id}`)
    await loadUsers()
  }
}

function roleLabel(role: string) {
  const labels: Record<string, string> = {
    platform_admin: '平台管理员',
    tenant_admin: '租户管理员',
    tenant_user: '普通用户'
  }
  return labels[role] || role
}

function formatDate(d: string) {
  return d ? new Date(d).toLocaleString() : '-'
}
</script>

<style scoped>
.users-page { padding: 20px; }
.header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.filters { display: flex; gap: 10px; margin-bottom: 15px; }
.filters select { padding: 6px 12px; border: 1px solid #d1d5db; border-radius: 4px; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th, .data-table td { padding: 10px; text-align: left; border-bottom: 1px solid #e5e7eb; }
.data-table th { background: #f9fafb; font-weight: 600; }
.btn-primary { background: #2563eb; color: white; padding: 8px 16px; border: none; border-radius: 4px; cursor: pointer; }
.btn-sm { padding: 4px 8px; margin-right: 4px; border: 1px solid #d1d5db; border-radius: 4px; background: white; cursor: pointer; }
.btn-danger { color: #dc2626; border-color: #dc2626; }
.role-platform_admin { color: #7c3aed; font-weight: 600; }
.role-tenant_admin { color: #2563eb; }
.status-active { color: #16a34a; }
.status-inactive { color: #9ca3af; }
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; }
.modal { background: white; padding: 24px; border-radius: 8px; min-width: 360px; }
.modal h3 { margin-bottom: 16px; }
.modal form { display: flex; flex-direction: column; gap: 12px; }
.modal input, .modal select { padding: 8px; border: 1px solid #d1d5db; border-radius: 4px; }
.modal-actions { display: flex; justify-content: flex-end; gap: 8px; margin-top: 8px; }
</style>
