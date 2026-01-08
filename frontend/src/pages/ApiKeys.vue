<template>
  <div class="apikeys-page">
    <div class="page-header">
      <h2>🔑 API密钥管理</h2>
      <button @click="showCreate = true" class="btn-primary">+ 创建密钥</button>
    </div>

    <div class="info-card">
      <p>💡 API密钥用于调用AI GuardRails API。请妥善保管，密钥只在创建时显示一次。</p>
    </div>

    <div class="keys-list">
      <div v-for="key in keys" :key="key.id" class="key-card">
        <div class="key-header">
          <span class="key-name">{{ key.name }}</span>
          <span :class="'key-status ' + (key.enabled ? 'active' : 'disabled')">
            {{ key.enabled ? '有效' : '已禁用' }}
          </span>
        </div>
        <div class="key-preview">
          <code>{{ key.key_prefix }}••••••••••••••••</code>
        </div>
        <div class="key-meta">
          <span>创建: {{ formatDate(key.created_at) }}</span>
          <span v-if="key.last_used_at">最后使用: {{ formatDate(key.last_used_at) }}</span>
          <span v-if="key.rate_limit_rpm">限速: {{ key.rate_limit_rpm }} 次/分</span>
        </div>
        <div class="key-actions">
          <button @click="revokeKey(key.id)" class="btn-danger">吊销</button>
        </div>
      </div>
      <div v-if="keys.length === 0" class="no-keys">
        暂无API密钥，点击上方按钮创建
      </div>
    </div>

    <!-- Create Modal -->
    <div v-if="showCreate" class="modal-overlay" @click.self="closeModal">
      <div class="modal">
        <h3>创建API密钥</h3>
        <div class="form-group">
          <label>名称</label>
          <input v-model="form.name" placeholder="如: 生产环境密钥" />
        </div>
        <div class="form-group">
          <label>速率限制 (次/分钟)</label>
          <input v-model.number="form.rate_limit_rpm" type="number" placeholder="留空为不限" />
        </div>
        <div class="form-group">
          <label>权限范围</label>
          <div class="scope-checks">
            <label><input type="checkbox" v-model="form.scopes.read" /> 读取</label>
            <label><input type="checkbox" v-model="form.scopes.write" /> 写入</label>
            <label><input type="checkbox" v-model="form.scopes.admin" /> 管理</label>
          </div>
        </div>
        <div class="modal-actions">
          <button @click="closeModal" class="btn-outline">取消</button>
          <button @click="createKey" class="btn-primary">创建</button>
        </div>
      </div>
    </div>

    <!-- Show Key Modal -->
    <div v-if="newKey" class="modal-overlay">
      <div class="modal">
        <h3>🔐 密钥已创建</h3>
        <div class="warning-box">
          ⚠️ 这是密钥唯一一次显示，请立即复制保存！
        </div>
        <div class="new-key-display">
          <code>{{ newKey }}</code>
          <button @click="copyKey" class="btn-copy">📋 复制</button>
        </div>
        <div class="modal-actions">
          <button @click="newKey = ''; loadKeys()" class="btn-primary">我已保存</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '../services/api'

interface APIKey {
  id: string
  name: string
  key_prefix: string
  enabled: boolean
  rate_limit_rpm?: number
  last_used_at?: string
  created_at: string
}

const keys = ref<APIKey[]>([])
const showCreate = ref(false)
const newKey = ref('')
const form = ref({
  name: '',
  rate_limit_rpm: null as number | null,
  scopes: { read: true, write: true, admin: false }
})

onMounted(() => loadKeys())

async function loadKeys() {
  try {
    keys.value = await api.get('/apikeys')
  } catch { keys.value = [] }
}

async function createKey() {
  const scopes = Object.entries(form.value.scopes).filter(([, v]) => v).map(([k]) => k)
  try {
    const result = await api.post('/apikeys', {
      name: form.value.name,
      scopes: scopes,
      rate_limit_rpm: form.value.rate_limit_rpm
    })
    newKey.value = result.key
    showCreate.value = false
  } catch {}
}

async function revokeKey(id: string) {
  if (!confirm('确定要吊销此密钥吗？吊销后无法恢复。')) return
  try {
    await api.delete(`/apikeys/${id}`)
    loadKeys()
  } catch {}
}

function copyKey() {
  navigator.clipboard.writeText(newKey.value)
  alert('已复制到剪贴板')
}

function closeModal() {
  showCreate.value = false
  form.value = { name: '', rate_limit_rpm: null, scopes: { read: true, write: true, admin: false } }
}

function formatDate(ts: string) {
  return new Date(ts).toLocaleString('zh-CN')
}
</script>

<style scoped>
.apikeys-page { padding: 20px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
.page-header h2 { margin: 0; }
.btn-primary { background: linear-gradient(90deg, #2563eb, #7c3aed); color: white; padding: 10px 20px; border: none; border-radius: 8px; cursor: pointer; }
.info-card { background: #eff6ff; padding: 16px; border-radius: 8px; margin-bottom: 24px; color: #1e40af; }
.keys-list { display: flex; flex-direction: column; gap: 16px; }
.key-card { background: white; padding: 20px; border-radius: 12px; box-shadow: 0 2px 10px rgba(0,0,0,0.05); }
.key-header { display: flex; justify-content: space-between; margin-bottom: 12px; }
.key-name { font-weight: 600; font-size: 1.1rem; }
.key-status { padding: 4px 12px; border-radius: 20px; font-size: 0.8rem; }
.key-status.active { background: #d1fae5; color: #065f46; }
.key-status.disabled { background: #fee2e2; color: #991b1b; }
.key-preview { background: #f1f5f9; padding: 12px; border-radius: 6px; margin-bottom: 12px; }
.key-preview code { font-family: monospace; font-size: 0.9rem; color: #334155; }
.key-meta { display: flex; gap: 20px; color: #64748b; font-size: 0.85rem; margin-bottom: 12px; }
.key-actions { display: flex; gap: 8px; }
.btn-danger { padding: 6px 16px; background: #ef4444; color: white; border: none; border-radius: 6px; cursor: pointer; }
.btn-outline { padding: 8px 16px; border: 1px solid #e2e8f0; background: white; border-radius: 6px; cursor: pointer; }
.no-keys { text-align: center; padding: 60px; color: #64748b; }
.modal-overlay { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 1000; }
.modal { background: white; padding: 24px; border-radius: 12px; width: 450px; }
.modal h3 { margin: 0 0 20px; }
.form-group { margin-bottom: 16px; }
.form-group label { display: block; margin-bottom: 6px; color: #374151; font-weight: 500; }
.form-group input { width: 100%; padding: 10px; border: 1px solid #e2e8f0; border-radius: 6px; box-sizing: border-box; }
.scope-checks { display: flex; gap: 20px; }
.scope-checks label { display: flex; align-items: center; gap: 6px; cursor: pointer; }
.modal-actions { display: flex; justify-content: flex-end; gap: 12px; margin-top: 24px; }
.warning-box { background: #fef3c7; padding: 12px; border-radius: 6px; color: #92400e; margin-bottom: 16px; }
.new-key-display { display: flex; gap: 12px; align-items: center; background: #f1f5f9; padding: 12px; border-radius: 6px; }
.new-key-display code { flex: 1; font-family: monospace; font-size: 0.85rem; word-break: break-all; }
.btn-copy { padding: 6px 12px; background: #3b82f6; color: white; border: none; border-radius: 4px; cursor: pointer; }
</style>
