<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-header">
        <span class="logo">🛡️</span>
        <h2>登录 AI GuardRails</h2>
        <p>工业AI应用安全护栏平台</p>
      </div>

      <!-- Tab Switch -->
      <div class="tabs">
        <button :class="{ active: tab === 'password' }" @click="tab = 'password'">账号密码</button>
        <button :class="{ active: tab === 'phone' }" @click="tab = 'phone'">手机验证码</button>
      </div>

      <!-- Password Login -->
      <form v-if="tab === 'password'" @submit.prevent="onPasswordLogin" class="login-form">
        <input v-model="username" placeholder="用户名" required />
        <input v-model="password" placeholder="密码" type="password" required />
        <button type="submit" :disabled="loading" class="btn-primary">登录</button>
      </form>

      <!-- Phone Login -->
      <form v-if="tab === 'phone'" @submit.prevent="onPhoneLogin" class="login-form">
        <input v-model="phone" placeholder="手机号" maxlength="11" required />
        <div class="code-input">
          <input v-model="code" placeholder="验证码" maxlength="6" required />
          <button type="button" @click="sendCode" :disabled="countdown > 0" class="btn-outline">
            {{ countdown > 0 ? countdown + 's' : '获取验证码' }}
          </button>
        </div>
        <button type="submit" :disabled="loading" class="btn-primary">登录 / 注册</button>
      </form>

      <!-- Social Login -->
      <div class="divider"><span>其他登录方式</span></div>
      <div class="social-buttons">
        <button @click="loginWeChat" class="social-btn wechat">
          <span class="icon">💬</span> 微信登录
        </button>
        <button @click="loginAlipay" class="social-btn alipay">
          <span class="icon">💰</span> 支付宝登录
        </button>
      </div>

      <div v-if="error" class="error">{{ error }}</div>

      <div class="login-footer">
        <router-link to="/landing">← 返回首页</router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { api, client } from '../services/api'

const router = useRouter()
const tab = ref<'password' | 'phone'>('password')
const username = ref('')
const password = ref('')
const phone = ref('')
const code = ref('')
const loading = ref(false)
const error = ref('')
const countdown = ref(0)

async function onPasswordLogin() {
  loading.value = true
  error.value = ''
  try {
    const res = await api.login(username.value, password.value)
    localStorage.setItem('auth_token', res.token)
    await router.push('/dashboard')
  } catch (e: any) {
    error.value = e?.response?.data || e?.message || '登录失败'
  } finally {
    loading.value = false
  }
}

async function sendCode() {
  if (phone.value.length !== 11) {
    error.value = '请输入正确的手机号'
    return
  }
  try {
    await client.post('/v1/auth/sms/send', { phone: phone.value })
    countdown.value = 60
    const timer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0) clearInterval(timer)
    }, 1000)
  } catch (e: any) {
    error.value = e?.response?.data || '发送失败'
  }
}

async function onPhoneLogin() {
  loading.value = true
  error.value = ''
  try {
    const res = await client.post('/v1/auth/sms/verify', { phone: phone.value, code: code.value })
    localStorage.setItem('auth_token', res.data.token)
    await router.push('/dashboard')
  } catch (e: any) {
    error.value = e?.response?.data || '验证失败'
  } finally {
    loading.value = false
  }
}

async function loginWeChat() {
  try {
    const res = await client.get('/v1/auth/oauth/wechat/url')
    window.location.href = res.data.url
  } catch (e: any) {
    error.value = '获取微信授权失败'
  }
}

async function loginAlipay() {
  try {
    const res = await client.get('/v1/auth/oauth/alipay/url')
    window.location.href = res.data.url
  } catch (e: any) {
    error.value = '获取支付宝授权失败'
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #0f172a 0%, #1e3a5f 100%);
}
.login-card {
  background: white;
  padding: 40px;
  border-radius: 16px;
  width: 400px;
  box-shadow: 0 20px 60px rgba(0,0,0,0.3);
}
.login-header {
  text-align: center;
  margin-bottom: 30px;
}
.logo {
  font-size: 3rem;
  display: block;
  margin-bottom: 12px;
}
.login-header h2 {
  margin: 0;
  color: #1e293b;
}
.login-header p {
  color: #64748b;
  margin-top: 4px;
}
.tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 24px;
}
.tabs button {
  flex: 1;
  padding: 12px;
  border: 1px solid #e2e8f0;
  background: white;
  border-radius: 8px;
  cursor: pointer;
  color: #64748b;
}
.tabs button.active {
  background: #2563eb;
  color: white;
  border-color: #2563eb;
}
.login-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.login-form input {
  padding: 14px;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  font-size: 1rem;
}
.code-input {
  display: flex;
  gap: 8px;
}
.code-input input {
  flex: 1;
}
.btn-primary {
  padding: 14px;
  background: linear-gradient(90deg, #2563eb, #7c3aed);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  cursor: pointer;
}
.btn-outline {
  padding: 14px 20px;
  border: 1px solid #2563eb;
  background: white;
  color: #2563eb;
  border-radius: 8px;
  cursor: pointer;
  white-space: nowrap;
}
.divider {
  display: flex;
  align-items: center;
  margin: 24px 0;
  color: #94a3b8;
  font-size: 0.875rem;
}
.divider::before, .divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: #e2e8f0;
}
.divider span {
  padding: 0 12px;
}
.social-buttons {
  display: flex;
  gap: 12px;
}
.social-btn {
  flex: 1;
  padding: 12px;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  background: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}
.social-btn.wechat:hover { background: #f0fff4; border-color: #22c55e; }
.social-btn.alipay:hover { background: #eff6ff; border-color: #3b82f6; }
.error {
  color: #dc2626;
  text-align: center;
  margin-top: 16px;
  padding: 10px;
  background: #fef2f2;
  border-radius: 6px;
}
.login-footer {
  text-align: center;
  margin-top: 20px;
}
.login-footer a {
  color: #64748b;
  text-decoration: none;
}
</style>
