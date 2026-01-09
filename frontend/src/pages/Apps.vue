<template>
  <div class="page-container">
    <div class="page-header">
      <h2>应用管理</h2>
      <div class="header-actions">
        <select v-model="selectedTenantId" @change="onTenantChange" class="tenant-select">
          <option value="" disabled>请选择租户</option>
          <option v-for="t in tenants" :key="t.id" :value="t.id">{{ t.name || '未命名' }} ({{ (t.id || '').substring(0,8) }})</option>
        </select>
        <button @click="openCreateModal" class="btn-primary" :disabled="!selectedTenantId">
          + 新建应用
        </button>
      </div>
    </div>

    <div class="table-container" v-if="selectedTenantId">
      <table class="data-table">
        <thead>
          <tr>
            <th>应用名称</th>
            <th>App ID</th>
            <th>API Secret</th>
            <th>配额/小时</th>
            <th>状态</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="app in apps" :key="app.id">
            <td class="font-bold">{{ app.name }}</td>
            <td class="font-mono text-sm">{{ app.id }}</td>
            <td class="font-mono text-sm">
              <div class="secret-cell">
                <span v-if="visibleSecrets[app.id]">{{ app.api_key_hash ? '(已隐藏)' : app.api_key }}</span>
                <span v-else>••••••••••••••••</span>
                <!-- 这里实际上后端返回的是明文Key吗？
                     根据之前的代码，createApp返回Key，listApps 可能不返回 Key 或者返回 Hash?
                     查看 api.go，如果是 listApps，通常不返回 Secret，只返回 Prefix/Hash。
                     如果是刚创建返回 Secret。
                     这里假设列表返回的是掩码或Hash，只有创建时显示一次。
                     或者如果有 Rotate 功能，会生成新的。
                -->
                <!-- 修正：listApps 应该不返回完整 Key。这里只显示状态或部分 -->
                <span class="badge-key">已配置</span>
              </div>
            </td>
            <td>{{ app.quota_per_hr }}</td>
            <td>
              <span class="badge badge-success" v-if="!app.is_revoked">正常</span>
              <span class="badge badge-danger" v-else>已吊销</span>
            </td>
            <td class="actions">
              <button @click="onRotate(app)" class="btn-sm btn-outline" :disabled="app.is_revoked">重置密钥</button>
              <button @click="onRevoke(app)" class="btn-sm btn-danger" :disabled="app.is_revoked">吊销</button>
            </td>
          </tr>
          <tr v-if="apps.length === 0">
            <td colspan="6" class="empty-state">该租户下暂无应用</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-else class="empty-state-large">
      <div class="icon">👈</div>
      <h3>请先选择一个租户</h3>
      <p>选择租户以管理其下的应用</p>
    </div>

    <!-- Create Modal -->
    <div v-if="showCreateModal" class="modal-overlay" @click.self="showCreateModal = false">
      <div class="modal">
        <h3>新建应用</h3>
        <form @submit.prevent="onCreate">
          <div class="form-group">
            <label>应用名称</label>
            <input v-model="form.name" placeholder="请输入应用名称" required />
          </div>
          <div class="form-group">
            <label>API调用配额 (次/小时)</label>
            <input v-model.number="form.quota" type="number" placeholder="1000" min="1" required />
          </div>
          <div class="modal-actions">
            <button type="button" @click="showCreateModal = false" class="btn-secondary">取消</button>
            <button type="submit" class="btn-primary" :disabled="loading">创建</button>
          </div>
        </form>
      </div>
    </div>

    <!-- Secret Display Modal (Created/Rotated) -->
    <div v-if="showSecretModal" class="modal-overlay">
      <div class="modal">
        <h3 class="text-success">🎉 {{ isRotate ? '密钥重置成功' : '应用创建成功' }}</h3>
        <p>请立即保存您的 API Key，因为它不会再次显示。</p>
        <div class="key-display">
          <div class="label">App ID</div>
          <div class="value">{{ currentApp.id }}</div>
        </div>
        <div class="key-display">
          <div class="label">API Key</div>
          <div class="value highlight">{{ currentKey }}</div>
          <button @click="copyKey" class="btn-sm btn-outline copy-btn">{{ copyText }}</button>
        </div>
        <div class="modal-actions">
          <button @click="showSecretModal = false" class="btn-primary">我已保存</button>
        </div>
      </div>
    </div>

    <!-- Confirm Modal -->
    <ConfirmModal 
      :is-open="showConfirmModal"
      :title="confirmTitle"
      :message="confirmMessage"
      :danger="confirmDanger"
      @confirm="onConfirmAction"
      @cancel="showConfirmModal = false"
    />

    <!-- Alert Modal -->
    <AlertModal
      :is-open="showAlertModal"
      :title="alertTitle"
      :message="alertMessage"
      :type="alertType"
      @close="showAlertModal = false"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api } from '../services/api'
import ConfirmModal from '../components/ConfirmModal.vue'
import AlertModal from '../components/AlertModal.vue'

const route = useRoute()
const router = useRouter()

const tenants = ref<any[]>([])
const selectedTenantId = ref('')
const apps = ref<any[]>([])
const visibleSecrets = reactive<Record<string, boolean>>({})
const loading = ref(false)

const showCreateModal = ref(false)
const showSecretModal = ref(false)
const currentApp = ref<any>({})
const currentKey = ref('')
const copyText = ref('复制')
const isRotate = ref(false)

// Forms
const form = reactive({
  name: '',
  quota: 1000
})

// Modals State
const showConfirmModal = ref(false)
const confirmTitle = ref('')
const confirmMessage = ref('')
const confirmDanger = ref(false)
const pendingAction = ref<() => Promise<void> | void>(() => {})

const showAlertModal = ref(false)
const alertTitle = ref('')
const alertMessage = ref('')
const alertType = ref('info')

function showAlert(msg: string, type = 'info', title = '提示') {
  alertMessage.value = msg
  alertType.value = type
  alertTitle.value = title
  showAlertModal.value = true
}

function showConfirm(title: string, msg: string, action: () => Promise<void> | void, danger = false) {
  confirmTitle.value = title
  confirmMessage.value = msg
  pendingAction.value = action
  confirmDanger.value = danger
  showConfirmModal.value = true
}

async function onConfirmAction() {
  showConfirmModal.value = false
  if (pendingAction.value) {
    await pendingAction.value()
  }
}

async function loadTenants() {
  try {
    const result = await api.listTenants()
    tenants.value = Array.isArray(result) ? result : []
    const qTenant = route.query.tenantId as string
    if (qTenant && tenants.value.some(t => t.id === qTenant)) {
      selectedTenantId.value = qTenant
      loadApps()
    }
  } catch (e) {
    console.error(e)
    tenants.value = []
  }
}

async function loadApps() {
  if (!selectedTenantId.value) return
  loading.value = true
  try {
    const result = await api.listApps(selectedTenantId.value)
    apps.value = Array.isArray(result) ? result : []
  } catch (e) {
    console.error(e)
    apps.value = []
  } finally {
    loading.value = false
  }
}

function onTenantChange() {
  router.push({ query: { tenantId: selectedTenantId.value } })
  loadApps()
}

function openCreateModal() {
  form.name = ''
  form.quota = 1000
  showCreateModal.value = true
}

async function onCreate() {
  if (!form.name) return
  loading.value = true
  try {
    const res = await api.createApp(selectedTenantId.value, form.name, form.quota)
    showCreateModal.value = false
    
    // Show Secret
    currentApp.value = res.app || res 
    currentKey.value = res.api_key || 'N/A'
    isRotate.value = false
    showSecretModal.value = true
    
    await loadApps()
  } catch (e) {
    showAlert('创建失败', 'error')
  } finally {
    loading.value = false
  }
}

function onRotate(app: any) {
  showConfirm(
    '重置密钥',
    `确定要重置应用 "${app.name}" 的密钥吗？旧密钥将立即失效。`,
    async () => {
      try {
        // Fix: Use generic /apps/{id}/rotate endpoint
        const res = await api.post(`/apps/${app.id}/rotate`, {}) 
        currentApp.value = app
        currentKey.value = res.api_key
        isRotate.value = true
        showSecretModal.value = true
      } catch (e: any) {
        showAlert('重置失败: ' + (e.response?.data?.error || e.message), 'error')
      }
    },
    true
  )
}

function onRevoke(app: any) {
  showConfirm(
    '吊销应用',
    `确定要吊销应用 "${app.name}" 吗？此操作不可逆。`,
    async () => {
      try {
        await api.post(`/apps/${app.id}/revoke`, {})
        await loadApps()
        showAlert('应用已吊销', 'success')
      } catch (e: any) {
        showAlert('吊销失败: ' + (e.response?.data?.error || e.message), 'error')
      }
    },
    true
  )
}

function copyKey() {
  navigator.clipboard.writeText(currentKey.value)
  copyText.value = '已复制!'
  setTimeout(() => copyText.value = '复制', 2000)
}

onMounted(() => {
  loadTenants()
})
</script>

<style scoped>
.page-container { padding: 24px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
.page-header h2 { margin: 0; font-size: 1.5rem; color: #1e293b; }
.header-actions { display: flex; gap: 12px; }

.tenant-select {
  padding: 8px 12px;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  font-size: 0.95rem;
  min-width: 200px;
}

.table-container { background: white; border-radius: 8px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); overflow: hidden; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th, .data-table td { padding: 16px; text-align: left; border-bottom: 1px solid #e2e8f0; }
.data-table th { background: #f8fafc; font-weight: 600; color: #64748b; font-size: 0.875rem; }
.data-table tr:hover { background: #f8fafc; }

.font-bold { font-weight: 600; color: #1e293b; }
.font-mono { font-family: 'SFMono-Regular', Consolas, monospace; color: #64748b; }
.text-sm { font-size: 0.875rem; }
.badge-key { background: #f1f5f9; color: #64748b; padding: 2px 6px; border-radius: 4px; font-size: 0.75rem; }

.badge { padding: 4px 8px; border-radius: 4px; font-size: 0.75rem; font-weight: 500; }
.badge-success { background: #dcfce7; color: #166534; }
.badge-danger { background: #fee2e2; color: #991b1b; }

.btn-primary { background: #2563eb; color: white; border: none; padding: 8px 16px; border-radius: 6px; cursor: pointer; }
.btn-primary:hover { background: #1d4ed8; }
.btn-primary:disabled { opacity: 0.7; cursor: not-allowed; }

.btn-secondary { background: white; border: 1px solid #cbd5e1; color: #475569; padding: 8px 16px; border-radius: 6px; cursor: pointer; }

.btn-sm { padding: 4px 10px; font-size: 0.8rem; border-radius: 4px; cursor: pointer; margin-right: 6px; }
.btn-outline { border: 1px solid #cbd5e1; background: white; color: #475569; }
.btn-outline:hover { border-color: #2563eb; color: #2563eb; }
.btn-danger { border: 1px solid #fca5a5; background: white; color: #ef4444; }
.btn-danger:hover { background: #fee2e2; }

/* Modal */
.modal-overlay { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 1000; }
.modal { background: white; padding: 24px; border-radius: 12px; width: 450px; box-shadow: 0 20px 25px -5px rgba(0,0,0,0.1); }
.modal h3 { margin: 0 0 20px; font-size: 1.25rem; }
.text-success { color: #16a34a; }

.form-group { margin-bottom: 20px; }
.form-group label { display: block; margin-bottom: 8px; font-weight: 500; }
.form-group input { width: 100%; padding: 8px; border: 1px solid #cbd5e1; border-radius: 6px; }
.modal-actions { display: flex; justify-content: flex-end; gap: 12px; margin-top: 24px; }

.key-display { margin-bottom: 16px; background: #f8fafc; padding: 12px; border-radius: 8px; border: 1px solid #e2e8f0; }
.key-display .label { font-size: 0.75rem; color: #64748b; margin-bottom: 4px; text-transform: uppercase; letter-spacing: 0.05em; }
.key-display .value { font-family: monospace; font-size: 1.1rem; color: #1e293b; word-break: break-all; }
.key-display .highlight { color: #2563eb; font-weight: 600; }
.copy-btn { float: right; margin-top: -30px; }

.empty-state { text-align: center; padding: 40px; color: #94a3b8; }
.empty-state-large { text-align: center; padding: 80px 0; color: #64748b; }
.empty-state-large .icon { font-size: 3rem; margin-bottom: 16px; }
</style>
