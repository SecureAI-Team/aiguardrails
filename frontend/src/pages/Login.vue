<template>
  <div class="login-page">
    <div class="login-container">
      <!-- Left: Branding -->
      <div class="login-branding">
        <div class="brand-content">
          <router-link to="/" class="brand-logo">
            <span class="logo-icon">🛡️</span>
            <span>AI GuardRails</span>
          </router-link>
          <h1>工业AI应用安全护栏平台</h1>
          <p>为您的AI应用提供全方位安全保护</p>
          <div class="brand-features">
            <div class="feature"><span>✓</span> 提示词注入防护</div>
            <div class="feature"><span>✓</span> 敏感数据脱敏</div>
            <div class="feature"><span>✓</span> 法规合规检查</div>
            <div class="feature"><span>✓</span> 实时审计日志</div>
          </div>
        </div>
      </div>
      
      <!-- Right: Form -->
      <div class="login-form-area">
        <div class="login-card">
          <h2>欢迎回来</h2>
          <p class="login-subtitle">登录您的账号继续使用</p>

          <!-- Tab Switch -->
          <div class="tabs">
            <button :class="{ active: tab === 'password' }" @click="tab = 'password'">
              <span class="tab-icon">🔑</span> 账号密码
            </button>
            <button :class="{ active: tab === 'phone' }" @click="tab = 'phone'">
              <span class="tab-icon">📱</span> 手机验证
            </button>
          </div>

          <!-- Password Login -->
          <form v-if="tab === 'password'" @submit.prevent="onPasswordLogin" class="login-form">
            <div class="form-group">
              <label>用户名</label>
              <input v-model="username" placeholder="请输入用户名" required />
            </div>
            <div class="form-group">
              <label>密码</label>
              <input v-model="password" placeholder="请输入密码" type="password" required />
            </div>
            <button type="submit" :disabled="loading" class="btn-primary">
              {{ loading ? '登录中...' : '登录' }}
            </button>
          </form>

          <!-- Phone Login -->
          <form v-if="tab === 'phone'" @submit.prevent="onPhoneLogin" class="login-form">
            <div class="form-group">
              <label>手机号</label>
              <input v-model="phone" placeholder="请输入手机号" maxlength="11" required />
            </div>
            <div class="form-group">
              <label>验证码</label>
              <div class="code-input">
                <input v-model="code" placeholder="请输入验证码" maxlength="6" required />
                <button type="button" @click="sendCode" :disabled="countdown > 0" class="btn-code">
                  {{ countdown > 0 ? countdown + 's' : '获取验证码' }}
                </button>
              </div>
            </div>
            <button type="submit" :disabled="loading" class="btn-primary">
              {{ loading ? '验证中...' : '登录 / 自动注册' }}
            </button>
          </form>

          <!-- Social Login -->
          <div class="divider"><span>第三方登录</span></div>
          <div class="social-buttons">
            <button @click="loginWeChat" class="social-btn wechat" title="微信登录">
              <svg viewBox="0 0 24 24" width="20" height="20" fill="#07C160">
                <path d="M8.691 2.188C3.891 2.188 0 5.476 0 9.53c0 2.212 1.17 4.203 3.002 5.55a.59.59 0 0 1 .213.665l-.39 1.48c-.019.07-.048.141-.048.213 0 .163.13.295.29.295a.32.32 0 0 0 .168-.054l1.903-1.114a.86.86 0 0 1 .717-.098 10.16 10.16 0 0 0 2.837.403c.276 0 .543-.027.811-.05-.857-2.578.157-5.506 2.625-7.194 2.468-1.688 5.867-1.962 8.63-.695C19.357 5.127 15.455 2.188 8.69 2.188zm12.618 7.343c-4.045 0-7.309 2.635-7.309 5.884 0 3.248 3.264 5.883 7.309 5.883.748 0 1.479-.095 2.172-.273a.68.68 0 0 1 .55.077l1.467.858a.24.24 0 0 0 .13.041.224.224 0 0 0 .219-.226c0-.055-.023-.11-.037-.163l-.3-1.138a.46.46 0 0 1 .163-.51c1.406-1.037 2.303-2.567 2.303-4.203 0-3.249-3.264-5.883-7.309-5.883z"/>
              </svg>
            </button>
            <button @click="loginAlipay" class="social-btn alipay" title="支付宝登录">
              <svg viewBox="0 0 24 24" width="20" height="20" fill="#1677FF">
                <path d="M21.422 15.358c-1.389-.56-6.12-2.47-7.992-3.093.678-1.559 1.168-3.346 1.335-5.265h-4.457V5.651h5.269V4.401h-5.269V1.5h-2.54v2.901H2.5v1.25h5.268V7h-4.457v1.25h9.926c-.142 1.397-.491 2.673-.981 3.809-1.594-.556-3.282-.868-4.869-.868-3.218 0-5.387 1.333-5.387 3.485 0 2.152 2.169 3.485 5.387 3.485 2.576 0 4.917-.868 6.648-2.295 1.639.868 5.636 2.62 7.016 3.207.357.152.733-.155.612-.525zm-14.035 1.103c-2.152 0-3.485-.868-3.485-2.152s1.333-2.152 3.485-2.152c1.333 0 2.737.265 4.07.715-1.233 2.152-2.737 3.589-4.07 3.589z"/>
              </svg>
            </button>
          </div>

          <div v-if="error" class="error">{{ error }}</div>

          <div class="login-footer">
            <router-link to="/landing">← 返回首页</router-link>
            <span class="separator">|</span>
            <router-link to="/playground">在线体验</router-link>
          </div>
        </div>
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
    localStorage.setItem('token', res.token)
    localStorage.setItem('username', username.value)
    await router.push('/dashboard')
  } catch (e: any) {
    error.value = e?.response?.data || e?.message || '登录失败，请检查用户名和密码'
  } finally {
    loading.value = false
  }
}

async function sendCode() {
  if (phone.value.length !== 11) {
    error.value = '请输入正确的11位手机号'
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
    error.value = e?.response?.data || '发送失败，请稍后重试'
  }
}

async function onPhoneLogin() {
  loading.value = true
  error.value = ''
  try {
    const res = await client.post('/v1/auth/sms/verify', { phone: phone.value, code: code.value })
    localStorage.setItem('token', res.data.token)
    localStorage.setItem('username', phone.value)
    await router.push('/dashboard')
  } catch (e: any) {
    error.value = e?.response?.data || '验证失败，请检查验证码'
  } finally {
    loading.value = false
  }
}

async function loginWeChat() {
  try {
    const res = await client.get('/v1/auth/oauth/wechat/url')
    window.location.href = res.data.url
  } catch (e: any) {
    error.value = '微信登录暂不可用'
  }
}

async function loginAlipay() {
  try {
    const res = await client.get('/v1/auth/oauth/alipay/url')
    window.location.href = res.data.url
  } catch (e: any) {
    error.value = '支付宝登录暂不可用'
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  background: #f8fafc;
}
.login-container {
  display: flex;
  min-height: 100vh;
}

/* Branding Side */
.login-branding {
  flex: 1;
  background: linear-gradient(135deg, #0f172a 0%, #1e3a5f 50%, #2563eb 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 60px;
}
.brand-content {
  max-width: 400px;
}
.brand-logo {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 1.5rem;
  font-weight: 700;
  color: white;
  text-decoration: none;
  margin-bottom: 40px;
}
.logo-icon { font-size: 2.5rem; }
.brand-content h1 {
  font-size: 2.2rem;
  font-weight: 700;
  margin-bottom: 16px;
  line-height: 1.3;
}
.brand-content > p {
  color: #94a3b8;
  font-size: 1.1rem;
  margin-bottom: 40px;
}
.brand-features {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.feature {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 1rem;
  color: #e2e8f0;
}
.feature span {
  color: #22c55e;
  font-weight: bold;
}

/* Form Side */
.login-form-area {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
}
.login-card {
  background: white;
  padding: 48px;
  border-radius: 20px;
  width: 420px;
  box-shadow: 0 20px 60px rgba(0,0,0,0.08);
}
.login-card h2 {
  margin: 0 0 8px;
  font-size: 1.8rem;
  color: #1e293b;
}
.login-subtitle {
  color: #64748b;
  margin: 0 0 32px;
}

/* Tabs */
.tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 28px;
}
.tabs button {
  flex: 1;
  padding: 14px;
  border: 2px solid #e2e8f0;
  background: white;
  border-radius: 10px;
  cursor: pointer;
  color: #64748b;
  font-size: 0.95rem;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  transition: all 0.2s;
}
.tabs button:hover {
  border-color: #3b82f6;
}
.tabs button.active {
  background: linear-gradient(90deg, #2563eb, #7c3aed);
  color: white;
  border-color: transparent;
}
.tab-icon { font-size: 1.1rem; }

/* Form */
.login-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}
.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.form-group label {
  font-size: 0.9rem;
  font-weight: 500;
  color: #374151;
}
.form-group input {
  padding: 14px 16px;
  border: 2px solid #e2e8f0;
  border-radius: 10px;
  font-size: 1rem;
  transition: border-color 0.2s;
}
.form-group input:focus {
  outline: none;
  border-color: #3b82f6;
}
.code-input {
  display: flex;
  gap: 10px;
}
.code-input input {
  flex: 1;
}
.btn-code {
  padding: 14px 18px;
  background: #f1f5f9;
  border: 2px solid #e2e8f0;
  border-radius: 10px;
  cursor: pointer;
  white-space: nowrap;
  font-size: 0.9rem;
  color: #3b82f6;
  font-weight: 500;
  transition: all 0.2s;
}
.btn-code:hover:not(:disabled) {
  background: #e0f2fe;
}
.btn-code:disabled {
  color: #94a3b8;
  cursor: not-allowed;
}
.btn-primary {
  padding: 16px;
  background: linear-gradient(90deg, #2563eb, #7c3aed);
  color: white;
  border: none;
  border-radius: 10px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}
.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 10px 30px rgba(37,99,235,0.3);
}
.btn-primary:disabled {
  opacity: 0.7;
  cursor: not-allowed;
  transform: none;
}

/* Divider */
.divider {
  display: flex;
  align-items: center;
  margin: 28px 0;
  color: #94a3b8;
  font-size: 0.85rem;
}
.divider::before, .divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: #e2e8f0;
}
.divider span {
  padding: 0 16px;
}

/* Social Buttons */
.social-buttons {
  display: flex;
  gap: 12px;
  justify-content: center;
}
.social-btn {
  width: 44px;
  height: 44px;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  background: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}
.social-btn .icon { font-size: 1.2rem; }
.social-btn.wechat:hover { background: #f0fff4; border-color: #22c55e; }
.social-btn.alipay:hover { background: #eff6ff; border-color: #3b82f6; }

/* Error */
.error {
  color: #dc2626;
  text-align: center;
  margin-top: 20px;
  padding: 12px;
  background: #fef2f2;
  border-radius: 8px;
  font-size: 0.9rem;
}

/* Footer */
.login-footer {
  text-align: center;
  margin-top: 28px;
  padding-top: 20px;
  border-top: 1px solid #e2e8f0;
}
.login-footer a {
  color: #64748b;
  text-decoration: none;
  font-size: 0.9rem;
}
.login-footer a:hover {
  color: #3b82f6;
}
.separator {
  margin: 0 12px;
  color: #e2e8f0;
}

@media (max-width: 900px) {
  .login-container { flex-direction: column; }
  .login-branding { padding: 40px; min-height: auto; }
  .login-form-area { padding: 20px; }
  .login-card { width: 100%; max-width: 420px; }
}
</style>
