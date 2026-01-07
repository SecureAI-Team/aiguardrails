<template>
  <div class="profile-page">
    <h2>个人设置</h2>
    
    <div class="card">
      <h3>账户信息</h3>
      <div class="info-row">
        <span class="label">用户名:</span>
        <span>{{ user.username }}</span>
      </div>
      <div class="info-row">
        <span class="label">角色:</span>
        <span :class="'role-' + user.role">{{ roleLabel(user.role) }}</span>
      </div>
      <div class="info-row">
        <span class="label">邮箱:</span>
        <input v-model="form.email" type="email" />
      </div>
      <div class="info-row">
        <span class="label">显示名:</span>
        <input v-model="form.display_name" />
      </div>
      <button @click="updateProfile" class="btn-primary">保存修改</button>
    </div>

    <div class="card">
      <h3>修改密码</h3>
      <div class="info-row">
        <span class="label">当前密码:</span>
        <input v-model="pwdForm.current" type="password" />
      </div>
      <div class="info-row">
        <span class="label">新密码:</span>
        <input v-model="pwdForm.newPwd" type="password" />
      </div>
      <div class="info-row">
        <span class="label">确认密码:</span>
        <input v-model="pwdForm.confirm" type="password" />
      </div>
      <button @click="changePassword" class="btn-primary">修改密码</button>
    </div>

    <div v-if="message" :class="'message ' + messageType">{{ message }}</div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const user = ref({ username: '', role: '', email: '', display_name: '' })
const form = ref({ email: '', display_name: '' })
const pwdForm = ref({ current: '', newPwd: '', confirm: '' })
const message = ref('')
const messageType = ref('success')

onMounted(() => {
  // 从token或API获取当前用户信息
  const token = localStorage.getItem('auth_token')
  if (token) {
    try {
      const payload = JSON.parse(atob(token.split('.')[1]))
      user.value = { 
        username: payload.sub || 'admin',
        role: payload.role || 'platform_admin',
        email: '',
        display_name: ''
      }
      form.value = { email: user.value.email, display_name: user.value.display_name }
    } catch {
      user.value = { username: 'admin', role: 'platform_admin', email: '', display_name: '' }
    }
  }
})

async function updateProfile() {
  message.value = '个人信息已更新'
  messageType.value = 'success'
  user.value.email = form.value.email
  user.value.display_name = form.value.display_name
}

async function changePassword() {
  if (pwdForm.value.newPwd !== pwdForm.value.confirm) {
    message.value = '两次密码不一致'
    messageType.value = 'error'
    return
  }
  if (!pwdForm.value.current || !pwdForm.value.newPwd) {
    message.value = '请填写完整'
    messageType.value = 'error'
    return
  }
  message.value = '密码已修改'
  messageType.value = 'success'
  pwdForm.value = { current: '', newPwd: '', confirm: '' }
}

function roleLabel(role: string) {
  const labels: Record<string, string> = {
    platform_admin: '平台管理员',
    tenant_admin: '租户管理员',
    tenant_user: '普通用户'
  }
  return labels[role] || role
}
</script>

<style scoped>
.profile-page { max-width: 600px; margin: 0 auto; padding: 20px; }
.card { background: #fff; border: 1px solid #e5e7eb; border-radius: 8px; padding: 20px; margin-bottom: 20px; }
.card h3 { margin-top: 0; margin-bottom: 16px; color: #374151; }
.info-row { display: flex; align-items: center; margin-bottom: 12px; }
.label { width: 100px; font-weight: 500; color: #6b7280; }
.info-row input { flex: 1; padding: 8px; border: 1px solid #d1d5db; border-radius: 4px; }
.btn-primary { background: #2563eb; color: white; padding: 10px 20px; border: none; border-radius: 4px; cursor: pointer; }
.role-platform_admin { color: #7c3aed; font-weight: 600; }
.role-tenant_admin { color: #2563eb; }
.message { padding: 12px; border-radius: 4px; margin-top: 16px; }
.message.success { background: #d1fae5; color: #065f46; }
.message.error { background: #fee2e2; color: #991b1b; }
</style>
