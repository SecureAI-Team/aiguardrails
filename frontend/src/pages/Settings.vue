<template>
  <div class="settings-page">
    <h2>🔧 系统设置</h2>

    <div class="settings-section">
      <h3>🤖 安全模型配置 (LLM Security)</h3>
      <div class="config-group">
        <h4>通义千问 (Qwen)</h4>
        <div class="config-row">
          <label>API Endpoint</label>
          <input v-model="config.qwen_endpoint" placeholder="https://dashscope.aliyuncs.com/..." />
        </div>
        <div class="config-row">
          <label>API Key</label>
          <input v-model="config.qwen_api_key" type="password" placeholder="sk-..." />
        </div>
        <div class="config-status">
          <span :class="config.qwen_api_key ? 'enabled' : 'disabled'">
            {{ config.qwen_api_key ? '✅ 已配置' : '❌ 未配置' }}
          </span>
        </div>
      </div>
    </div>

    <div class="settings-section">
      <h3>🔐 社交登录配置</h3>
      
      <div class="config-group">
        <h4>微信登录</h4>
        <div class="config-row">
          <label>AppID</label>
          <input v-model="config.wechat_app_id" placeholder="wx1234567890" />
        </div>
        <div class="config-row">
          <label>AppSecret</label>
          <input v-model="config.wechat_app_secret" type="password" placeholder="******" />
        </div>
        <div class="config-status">
          <span :class="config.wechat_enabled ? 'enabled' : 'disabled'">
            {{ config.wechat_enabled ? '✅ 已启用' : '❌ 未配置' }}
          </span>
        </div>
      </div>

      <div class="config-group">
        <h4>支付宝登录</h4>
        <div class="config-row">
          <label>AppID</label>
          <input v-model="config.alipay_app_id" placeholder="2021001234567890" />
        </div>
        <div class="config-row">
          <label>应用私钥</label>
          <textarea v-model="config.alipay_private_key" placeholder="RSA2私钥" rows="3"></textarea>
        </div>
        <div class="config-row">
          <label>支付宝公钥</label>
          <textarea v-model="config.alipay_public_key" placeholder="支付宝公钥" rows="3"></textarea>
        </div>
        <div class="config-status">
          <span :class="config.alipay_enabled ? 'enabled' : 'disabled'">
            {{ config.alipay_enabled ? '✅ 已启用' : '❌ 未配置' }}
          </span>
        </div>
      </div>
    </div>

    <div class="settings-section">
      <h3>📱 短信服务配置</h3>
      
      <div class="config-group">
        <div class="config-row">
          <label>服务商</label>
          <select v-model="config.sms_provider">
            <option value="">请选择</option>
            <option value="aliyun">阿里云短信</option>
            <option value="tencent">腾讯云短信</option>
          </select>
        </div>
        <div class="config-row">
          <label>AccessKey</label>
          <input v-model="config.sms_access_key" placeholder="AccessKey ID" />
        </div>
        <div class="config-row">
          <label>SecretKey</label>
          <input v-model="config.sms_secret_key" type="password" placeholder="AccessKey Secret" />
        </div>
        <div class="config-row">
          <label>签名名称</label>
          <input v-model="config.sms_sign_name" placeholder="短信签名" />
        </div>
        <div class="config-row">
          <label>模板ID</label>
          <input v-model="config.sms_template_code" placeholder="SMS_123456789" />
        </div>
        <div class="config-status">
          <span :class="config.sms_provider ? 'enabled' : 'disabled'">
            {{ config.sms_provider ? '✅ 已配置 (' + config.sms_provider + ')' : '❌ 未配置' }}
          </span>
        </div>
      </div>
    </div>

    <div class="settings-section">
      <h3>🔑 回调地址</h3>
      <div class="config-group">
        <div class="config-row">
          <label>OAuth回调URL</label>
          <input v-model="config.callback_url" placeholder="https://your-domain.com/callback" />
        </div>
        <p class="hint">请在微信开放平台和支付宝开放平台配置此回调地址</p>
      </div>
    </div>

    <div class="actions">
      <button @click="saveConfig" class="btn-primary" :disabled="saving">
        {{ saving ? '保存中...' : '保存配置' }}
      </button>
      <span v-if="message" :class="'msg ' + msgType">{{ message }}</span>
    </div>

    <div class="env-hint">
      <h4>💡 环境变量配置</h4>
      <p>也可以通过环境变量配置（推荐用于生产环境）：</p>
      <pre>
WECHAT_APP_ID=your_wechat_app_id
WECHAT_APP_SECRET=your_wechat_app_secret
ALIPAY_APP_ID=your_alipay_app_id
ALIPAY_PRIVATE_KEY=your_private_key
SMS_PROVIDER=aliyun
SMS_ACCESS_KEY=your_access_key
SMS_SECRET_KEY=your_secret_key
SMS_SIGN_NAME=your_sign_name
SMS_TEMPLATE_CODE=SMS_123456789
      </pre>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '../services/api'

const config = ref({
  qwen_api_key: '',
  qwen_endpoint: '',
  wechat_app_id: '',
  wechat_app_secret: '',
  wechat_enabled: false,
  alipay_app_id: '',
  alipay_private_key: '',
  alipay_public_key: '',
  alipay_enabled: false,
  sms_provider: '',
  sms_access_key: '',
  sms_secret_key: '',
  sms_sign_name: '',
  sms_template_code: '',
  callback_url: window.location.origin + '/callback'
})

const saving = ref(false)
const message = ref('')
const msgType = ref('success')

onMounted(async () => {
  try {
    const res = await api.get('/settings') // Backend settings
    if (res.data) {
      Object.assign(config.value, res.data)
    }
    // Also load local storage for social auth (legacy demo)
    const saved = localStorage.getItem('social_auth_config')
    if (saved) {
      try {
        const parsed = JSON.parse(saved)
        Object.assign(config.value, parsed)
      } catch {}
    }
  } catch (e) {
    console.error("Failed to load settings", e)
  }
})

async function saveConfig() {
  saving.value = true
  message.value = ''
  
  try {
    // 1. Save Backend Settings
    await api.post('/settings', {
      settings: {
        qwen_api_key: config.value.qwen_api_key,
        qwen_endpoint: config.value.qwen_endpoint
      }
    })

    // 2. Save Local Storage (Social Auth Demo)
    localStorage.setItem('social_auth_config', JSON.stringify(config.value))
    
    message.value = '配置已保存'
    msgType.value = 'success'
  } catch (e) {
    message.value = '保存失败'
    msgType.value = 'error'
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.settings-page { max-width: 800px; margin: 0 auto; padding: 20px; }
.settings-page h2 { margin-bottom: 24px; color: #1e293b; }
.settings-section { background: white; border-radius: 12px; padding: 24px; margin-bottom: 24px; box-shadow: 0 2px 10px rgba(0,0,0,0.05); }
.settings-section h3 { margin: 0 0 20px; color: #374151; font-size: 1.1rem; }
.config-group { background: #f8fafc; padding: 20px; border-radius: 8px; margin-bottom: 16px; }
.config-group h4 { margin: 0 0 16px; color: #475569; }
.config-row { display: flex; align-items: center; margin-bottom: 12px; }
.config-row label { width: 120px; color: #64748b; font-weight: 500; }
.config-row input, .config-row select, .config-row textarea { flex: 1; padding: 10px; border: 1px solid #e2e8f0; border-radius: 6px; font-size: 0.95rem; }
.config-row textarea { font-family: monospace; font-size: 0.85rem; }
.config-status { margin-top: 12px; padding-top: 12px; border-top: 1px solid #e2e8f0; }
.config-status .enabled { color: #16a34a; }
.config-status .disabled { color: #9ca3af; }
.hint { color: #64748b; font-size: 0.85rem; margin-top: 8px; }
.actions { display: flex; align-items: center; gap: 16px; margin-bottom: 24px; }
.btn-primary { background: linear-gradient(90deg, #2563eb, #7c3aed); color: white; padding: 12px 24px; border: none; border-radius: 8px; cursor: pointer; font-size: 1rem; }
.btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
.msg { font-size: 0.9rem; }
.msg.success { color: #16a34a; }
.msg.error { color: #dc2626; }
.env-hint { background: #f1f5f9; padding: 20px; border-radius: 8px; }
.env-hint h4 { margin: 0 0 12px; color: #475569; }
.env-hint p { color: #64748b; margin-bottom: 12px; }
.env-hint pre { background: #1e293b; color: #e2e8f0; padding: 16px; border-radius: 6px; overflow-x: auto; font-size: 0.85rem; }
</style>
