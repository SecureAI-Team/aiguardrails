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
          <h2>{{ isRegister ? '创建新账号' : '欢迎回来' }}</h2>
          <p class="login-subtitle">
            {{ isRegister ? '注册即刻开启AI安全之旅' : '登录您的账号继续使用' }}
          </p>

          <form @submit.prevent="handleSubmit" class="login-form">
            <div class="form-group">
              <label>用户名</label>
              <input v-model="form.username" placeholder="请输入用户名" required />
            </div>
            <div class="form-group">
              <label>密码</label>
              <input v-model="form.password" placeholder="请输入密码" type="password" required />
            </div>
            
            <button type="submit" :disabled="loading" class="btn-primary">
              {{ loading ? '处理中...' : (isRegister ? '立即注册' : '登录') }}
            </button>
          </form>

          <div v-if="error" class="error">{{ error }}</div>

          <div class="toggle-area">
            <span v-if="!isRegister">
              还没有账号？ <a href="#" @click.prevent="isRegister = true">现在注册</a>
            </span>
            <span v-else>
              已有账号？ <a href="#" @click.prevent="isRegister = false">立即登录</a>
            </span>
          </div>

          <!-- Divider for other methods -->
          <div class="divider"><span>其他登录方式</span></div>
          <div class="social-buttons">
             <button class="social-btn" title="手机验证码登录" @click="alert('暂未开放手机登录')">📱</button>
             <button class="social-btn" title="微信登录" @click="alert('暂未开放微信登录')">💬</button>
          </div>

          <div class="login-footer">
            <router-link to="/landing">← 返回首页</router-link>
            <span class="separator">|</span>
            <a href="#" @click.prevent>忘记密码?</a>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../services/api'

const router = useRouter()
const isRegister = ref(false)
const loading = ref(false)
const error = ref('')

const form = reactive({
  username: '',
  password: ''
})

function alert(msg: string) {
  window.alert(msg)
}

async function handleSubmit() {
  loading.value = true
  error.value = ''
  try {
    let res
    if (isRegister.value) {
      // Register
      res = await api.post('/auth/register', form)
    } else {
      // Login
      res = await api.login(form.username, form.password)
    }
    
    // Both return token
    localStorage.setItem('auth_token', res.token)
    localStorage.setItem('username', form.username)
    await router.push('/dashboard')
  } catch (e: any) {
    const msg = e?.response?.data || e?.message || '操作失败'
    if (msg.includes('conflict')) error.value = '用户名已存在'
    else if (msg.includes('unauthorized')) error.value = '用户名或密码错误'
    else error.value = msg
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page { min-height: 100vh; background: #f8fafc; }
.login-container { display: flex; min-height: 100vh; }

/* Branding Side */
.login-branding { flex: 1; background: linear-gradient(135deg, #0f172a 0%, #1e3a5f 50%, #2563eb 100%); color: white; display: flex; align-items: center; justify-content: center; padding: 60px; }
.brand-content { max-width: 400px; }
.brand-logo { display: flex; align-items: center; gap: 12px; font-size: 1.5rem; font-weight: 700; color: white; text-decoration: none; margin-bottom: 40px; }
.brand-content h1 { font-size: 2.2rem; font-weight: 700; margin-bottom: 16px; line-height: 1.3; }
.brand-content > p { color: #94a3b8; font-size: 1.1rem; margin-bottom: 40px; }
.brand-features { display: flex; flex-direction: column; gap: 16px; }
.feature { display: flex; align-items: center; gap: 12px; font-size: 1rem; color: #e2e8f0; }
.feature span { color: #22c55e; font-weight: bold; }

/* Form Side */
.login-form-area { flex: 1; display: flex; align-items: center; justify-content: center; padding: 40px; }
.login-card { background: white; padding: 48px; border-radius: 20px; width: 420px; box-shadow: 0 20px 60px rgba(0,0,0,0.08); }
.login-card h2 { margin: 0 0 8px; font-size: 1.8rem; color: #1e293b; }
.login-subtitle { color: #64748b; margin: 0 0 32px; }

.login-form { display: flex; flex-direction: column; gap: 20px; }
.form-group { display: flex; flex-direction: column; gap: 8px; }
.form-group label { font-size: 0.9rem; font-weight: 500; color: #374151; }
.form-group input { padding: 14px 16px; border: 2px solid #e2e8f0; border-radius: 10px; font-size: 1rem; transition: border-color 0.2s; }
.form-group input:focus { outline: none; border-color: #3b82f6; }

.btn-primary { padding: 16px; background: linear-gradient(90deg, #2563eb, #7c3aed); color: white; border: none; border-radius: 10px; font-size: 1rem; font-weight: 600; cursor: pointer; transition: all 0.2s; margin-top: 10px; }
.btn-primary:hover:not(:disabled) { transform: translateY(-2px); box-shadow: 0 10px 30px rgba(37,99,235,0.3); }
.btn-primary:disabled { opacity: 0.7; cursor: not-allowed; }

.toggle-area { text-align: center; margin-top: 24px; color: #64748b; font-size: 0.95rem; }
.toggle-area a { color: #2563eb; text-decoration: none; font-weight: 600; }

.divider { display: flex; align-items: center; margin: 28px 0; color: #94a3b8; font-size: 0.85rem; }
.divider::before, .divider::after { content: ''; flex: 1; height: 1px; background: #e2e8f0; }
.divider span { padding: 0 16px; }

.social-buttons { display: flex; gap: 12px; justify-content: center; }
.social-btn { width: 44px; height: 44px; border: 1px solid #e2e8f0; border-radius: 10px; background: white; cursor: pointer; display: flex; align-items: center; justify-content: center; font-size: 1.2rem; transition: all 0.2s; }
.social-btn:hover { background: #f8fafc; border-color: #cbd5e1; }

.error { color: #dc2626; text-align: center; margin-top: 20px; padding: 12px; background: #fef2f2; border-radius: 8px; font-size: 0.9rem; }

.login-footer { text-align: center; margin-top: 28px; padding-top: 20px; border-top: 1px solid #e2e8f0; }
.login-footer a { color: #64748b; text-decoration: none; font-size: 0.9rem; }
.login-footer a:hover { color: #3b82f6; }
.separator { margin: 0 12px; color: #e2e8f0; }

@media (max-width: 900px) {
  .login-container { flex-direction: column; }
  .login-branding { padding: 40px; min-height: auto; }
  .login-form-area { padding: 20px; }
  .login-card { width: 100%; max-width: 420px; }
}
</style>
