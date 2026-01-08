<template>
  <div class="page-container">
    <div class="page-header">
      <h2>平台用户管理</h2>
      <div class="header-actions">
        <select v-model="filterRole" @change="loadUsers" class="filter-select">
          <option value="">所有角色</option>
          <option value="platform_admin">平台管理员</option>
          <option value="tenant_admin">租户管理员</option>
          <option value="tenant_user">普通用户</option>
        </select>
        <button @click="openCreateModal" class="btn-primary">+ 新建用户</button>
      </div>
    </div>

    <div class="table-container">
      <table class="data-table">
        <thead>
          <tr>
            <th>用户名</th>
            <th>显示名</th>
            <th>角色</th>
            <th>状态</th>
            <th>最后登录</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in users" :key="user.id">
            <td class="font-bold">{{ user.username }}</td>
            <td>
              <div v-if="user.display_name">{{ user.display_name }}</div>
              <div v-if="user.email" class="muted text-sm">{{ user.email }}</div>
            </td>
            <td><span :class="['badge', getRoleBadge(user.role)]">{{ roleLabel(user.role) }}</span></td>
            <td>
              <span :class="['badge', user.status === 'active' ? 'badge-success' : 'badge-danger']">
                {{ user.status === 'active' ? '启用' : '禁用' }}
              </span>
            </td>
            <td class="text-sm">{{ formatDate(user.last_login_at) }}</td>
            <td class="actions">
              <button @click="editUser(user)" class="btn-xs btn-outline">编辑</button>
              <button @click="resetPwd(user)" class="btn-xs btn-outline">重置密码</button>
              <button @click="deleteUser(user)" class="btn-xs btn-danger">删除</button>
            </td>
          </tr>
          <tr v-if="users.length === 0">
            <td colspan="6" class="empty-state">暂无用户</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Create Modal -->
    <div v-if="showCreate" class="modal-overlay" @click.self="showCreate = false">
      <div class="modal">
        <h3>新建用户</h3>
        <form @submit.prevent="createUser">
          <div class="form-group">
            <label>用户名</label>
            <input v-model="form.username" placeholder="登录用户名" required :disabled="loading" />
          </div>
          <div class="form-group">
            <label>初始密码</label>
            <input v-model="form.password" type="password" placeholder="设置初始密码" required :disabled="loading" />
          </div>
          <div class="form-group">
            <label>角色</label>
            <select v-model="form.role" class="form-select">
              <option value="tenant_user">普通用户</option>
              <option value="tenant_admin">租户管理员</option>
              <option value="platform_admin">平台管理员</option>
            </select>
          </div>
          <div class="form-group">
            <label>显示名 (可选)</label>
            <input v-model="form.display_name" placeholder="真实姓名或昵称" />
          </div>
          <div class="form-group">
            <label>邮箱 (可选)</label>
            <input v-model="form.email" type="email" placeholder="contact@example.com" />
          </div>
          
          <div class="modal-actions">
            <button type="button" @click="showCreate = false" class="btn-secondary">取消</button>
            <button type="submit" class="btn-primary" :disabled="loading">创建</button>
          </div>
        </form>
      </div>
    </div>

    <!-- Edit Modal -->
    <div v-if="showEdit" class="modal-overlay" @click.self="showEdit = false">
      <div class="modal">
        <h3>编辑用户</h3>
        <form @submit.prevent="updateUser">
          <div class="form-group">
            <label>显示名</label>
            <input v-model="editForm.display_name" />
          </div>
          <div class="form-group">
            <label>邮箱</label>
            <input v-model="editForm.email" />
          </div>
          <div class="form-group">
            <label>角色</label>
            <select v-model="editForm.role" class="form-select">
              <option value="tenant_user">普通用户</option>
              <option value="tenant_admin">租户管理员</option>
              <option value="platform_admin">平台管理员</option>
            </select>
          </div>
          <div class="form-group">
            <label>状态</label>
            <select v-model="editForm.status" class="form-select">
              <option value="active">启用</option>
              <option value="inactive">禁用</option>
            </select>
          </div>
          <div class="modal-actions">
            <button type="button" @click="showEdit = false" class="btn-secondary">取消</button>
            <button type="submit" class="btn-primary" :disabled="loading">保存此变更</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { api } from '../services/api'

const users = ref<any[]>([])
const loading = ref(false)
const filterRole = ref('')
const showCreate = ref(false)
const showEdit = ref(false)

const form = reactive({ username: '', password: '', email: '', display_name: '', role: 'tenant_user' })
const editForm = reactive({ id: '', email: '', display_name: '', role: '', status: '' })

async function loadUsers() {
  loading.value = true
  try {
    const params = new URLSearchParams()
    if (filterRole.value) params.append('role', filterRole.value)
    users.value = await api.get(`/users?${params}`)
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function openCreateModal() {
  form.username = ''
  form.password = ''
  form.email = ''
  form.display_name = ''
  form.role = 'tenant_user'
  showCreate.value = true
}

async function createUser() {
  loading.value = true
  try {
    await api.post('/users', form)
    showCreate.value = false
    await loadUsers()
  } catch (e: any) {
    alert('创建失败: ' + (e.response?.data?.error || e.message))
  } finally {
    loading.value = false
  }
}

function editUser(user: any) {
  editForm.id = user.id
  editForm.email = user.email
  editForm.display_name = user.display_name
  editForm.role = user.role
  editForm.status = user.status
  showEdit.value = true
}

async function updateUser() {
  loading.value = true
  try {
    await api.put(`/users/${editForm.id}`, {
      email: editForm.email,
      display_name: editForm.display_name,
      role: editForm.role,
      status: editForm.status
    })
    showEdit.value = false
    await loadUsers()
  } catch (e: any) {
    alert('更新失败: ' + (e.response?.data?.error || e.message))
  } finally {
    loading.value = false
  }
}

async function resetPwd(user: any) {
  const pwd = prompt(`请输入用户 ${user.username} 的新密码:`)
  if (!pwd) return
  try {
    await api.post(`/users/${user.id}/password`, { new_password: pwd })
    alert('密码已重置成功')
  } catch (e: any) {
    alert('重置失败')
  }
}

async function deleteUser(user: any) {
  if (!confirm(`确定要永久删除用户 "${user.username}" 吗？此操作不可恢复。`)) return
  try {
    await api.delete(`/users/${user.id}`)
    await loadUsers()
  } catch(e) {
    alert('删除失败')
  }
}

function roleLabel(role: string) {
  const map: Record<string, string> = {
    platform_admin: '平台管理员',
    tenant_admin: '租户管理员',
    tenant_user: '普通用户'
  }
  return map[role] || role
}

function getRoleBadge(role: string) {
  if (role === 'platform_admin') return 'badge-purple'
  if (role === 'tenant_admin') return 'badge-blue'
  return 'badge-gray'
}

function formatDate(d?: string) {
  if (!d) return '-'
  return new Date(d).toLocaleString()
}

onMounted(loadUsers)
</script>

<style scoped>
.page-container { padding: 24px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
.page-header h2 { margin: 0; font-size: 1.5rem; color: #1e293b; }
.header-actions { display: flex; gap: 12px; }

.filter-select { padding: 8px 12px; border: 1px solid #cbd5e1; border-radius: 6px; }

.table-container { background: white; border-radius: 8px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); overflow: hidden; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th, .data-table td { padding: 16px; text-align: left; border-bottom: 1px solid #e2e8f0; }
.data-table th { background: #f8fafc; font-weight: 600; color: #64748b; font-size: 0.875rem; }
.data-table tr:hover { background: #f8fafc; }

.font-bold { font-weight: 600; color: #1e293b; }
.muted { color: #64748b; }
.text-sm { font-size: 0.875rem; }

.badge { padding: 2px 8px; border-radius: 12px; font-size: 0.75rem; font-weight: 500; }
.badge-purple { background: #f3e8ff; color: #7e22ce; }
.badge-blue { background: #dbeafe; color: #1e40af; }
.badge-gray { background: #f1f5f9; color: #475569; }
.badge-success { background: #dcfce7; color: #166534; }
.badge-danger { background: #fee2e2; color: #991b1b; }

.btn-primary { background: #2563eb; color: white; border: none; padding: 8px 16px; border-radius: 6px; cursor: pointer; }
.btn-secondary { background: white; border: 1px solid #cbd5e1; color: #475569; padding: 8px 16px; border-radius: 6px; cursor: pointer; }

.actions { display: flex; gap: 6px; }
.btn-xs { padding: 4px 8px; font-size: 0.75rem; border-radius: 4px; cursor: pointer; }
.btn-outline { background: white; border: 1px solid #cbd5e1; color: #475569; }
.btn-outline:hover { border-color: #2563eb; color: #2563eb; }
.btn-danger { background: white; border: 1px solid #fca5a5; color: #ef4444; }
.btn-danger:hover { background: #fee2e2; }

/* Modal */
.modal-overlay { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 1000; }
.modal { background: white; padding: 24px; border-radius: 12px; width: 400px; box-shadow: 0 20px 25px -5px rgba(0,0,0,0.1); }
.modal h3 { margin: 0 0 20px; font-size: 1.25rem; }

.form-group { margin-bottom: 16px; }
.form-group label { display: block; margin-bottom: 6px; font-weight: 500; font-size: 0.9rem; }
.form-group input, .form-group select { width: 100%; padding: 8px 12px; border: 1px solid #cbd5e1; border-radius: 6px; font-size: 0.95rem; }
.modal-actions { display: flex; justify-content: flex-end; gap: 12px; margin-top: 24px; }
.empty-state { text-align: center; padding: 40px; color: #94a3b8; }
</style>
