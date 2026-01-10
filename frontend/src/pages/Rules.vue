<template>
  <div class="page-container">
    <div class="page-header">
      <h2>规则库 (Rule Library)</h2>
      <button @click="openCreateModal" class="btn-primary">+ 新建规则</button>
    </div>

    <!-- Filter Bar -->
    <div class="filter-bar">
      <button 
        :class="['filter-btn', activeTab === 'all' ? 'active' : '']"
        @click="activeTab = 'all'"
      >全部规则</button>
      <button 
        :class="['filter-btn', activeTab === 'custom' ? 'active' : '']"
        @click="activeTab = 'custom'"
      >自定义 ({{ customRules.length }})</button>
      <button 
        :class="['filter-btn', activeTab === 'system' ? 'active' : '']"
        @click="activeTab = 'system'"
      >内置 ({{ systemRules.length }})</button>
    </div>

    <div class="table-container">
      <table class="data-table">
        <thead>
          <tr>
            <th>规则名称</th>
            <th>类型</th>
            <th>描述</th>
            <th>来源</th>
            <th style="width: 120px">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="rule in filteredRules" :key="rule.id">
            <td class="font-bold">{{ rule.name }}</td>
            <td>
              <span :class="['badge', getTypeClass(rule.type)]">
                <span v-if="rule.type === 'llm'">🤖 LLM</span>
                <span v-else-if="rule.type === 'opa'">📜 OPA</span>
                <span v-else-if="rule.type === 'keyword'">🚫 Keyword</span>
                <span v-else>{{ rule.type }}</span>
              </span>
            </td>
            <td class="text-desc" :title="rule.description">{{ rule.description || '-' }}</td>
            <td>
              <span v-if="rule.is_system" class="badge badge-system">内置</span>
              <span v-else class="badge badge-custom">自定义</span>
            </td>
            <td class="actions">
              <button v-if="!rule.is_system" @click="deleteRule(rule.id)" class="btn-sm btn-outline btn-danger">删除</button>
              <button v-else disabled class="btn-sm btn-outline disabled">系统锁定</button>
            </td>
          </tr>
          <tr v-if="filteredRules.length === 0">
            <td colspan="5" class="empty-state">
              暂无规则数据
              <button v-if="activeTab !== 'all'" @click="activeTab = 'all'" class="btn-link">查看全部</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Create Modal (Teleported) -->
    <Teleport to="body">
      <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
        <div class="modal modal-lg">
          <div class="modal-header">
            <h3>新建规则</h3>
            <button @click="showModal = false" class="close-btn">&times;</button>
          </div>
          
          <div class="modal-body">
            <div class="form-grid">
              <!-- Left: Meta -->
              <div class="form-col">
                <div class="form-group">
                  <label>规则名称</label>
                  <input v-model="newRule.name" placeholder="例如: 拒绝竞品信息" class="form-input" />
                </div>
                
                <div class="form-group">
                  <label>规则类型</label>
                  <div class="type-selector">
                    <div 
                      v-for="t in ['llm', 'opa', 'keyword']" 
                      :key="t"
                      :class="['type-option', newRule.type === t ? 'active' : '']"
                      @click="newRule.type = t"
                    >
                      <span v-if="t==='llm'">🤖 LLM Security</span>
                      <span v-else-if="t==='opa'">📜 OPA Policy</span>
                      <span v-else>🚫 Keyword List</span>
                    </div>
                  </div>
                </div>

                <div class="form-group">
                  <label>描述</label>
                  <textarea v-model="newRule.description" rows="3" class="form-input" placeholder="简要描述规则用途"></textarea>
                </div>
              </div>

              <!-- Right: Content -->
              <div class="form-col full-height">
                <div class="form-group flex-1">
                  <label>
                    <span v-if="newRule.type === 'llm'">System Prompt 指令</span>
                    <span v-else-if="newRule.type === 'keyword'">敏感词列表 (每行一个)</span>
                    <span v-else>Rego 代码</span>
                  </label>
                  <textarea 
                    v-model="newRule.content" 
                    class="form-input code-editor" 
                    :placeholder="placeholderText"
                  ></textarea>
                </div>
              </div>
            </div>
          </div>

          <div class="modal-actions">
            <button @click="showModal = false" class="btn-secondary">取消</button>
            <button @click="createRule" class="btn-primary" :disabled="loading || !newRule.name || !newRule.content">
              {{ loading ? '创建中...' : '确认创建' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>

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
import { ref, onMounted, computed, reactive } from 'vue'
import { api } from '../services/api'
import AlertModal from '../components/AlertModal.vue'

interface Rule {
  id: string
  name: string
  description: string
  type: string
  content: string
  is_system: boolean
  created_at: string
}

const rules = ref<Rule[]>([])
const loading = ref(false)
const activeTab = ref<'all'|'custom'|'system'>('all')

const showModal = ref(false)
const newRule = reactive({
  name: '',
  description: '',
  type: 'llm',
  content: ''
})

// Alert state
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

const systemRules = computed(() => rules.value.filter(r => r.is_system))
const customRules = computed(() => rules.value.filter(r => !r.is_system))

const filteredRules = computed(() => {
  if (activeTab.value === 'system') return systemRules.value
  if (activeTab.value === 'custom') return customRules.value
  return rules.value
})

const placeholderText = computed(() => {
  if (newRule.type === 'llm') return 'You are a helpful assistant. Please ensure the response does not contain PII...'
  if (newRule.type === 'keyword') return '敏感词1\n敏感词2\nblocked_word'
  return 'package guardrails\n\ndefault allow = true\n...'
})

function getTypeClass(type: string) {
  if (type === 'llm') return 'badge-llm'
  if (type === 'opa') return 'badge-opa'
  if (type === 'keyword') return 'badge-kw'
  return 'badge-gray'
}

async function loadRules() {
  try {
    const list = await api.listRules()
    rules.value = Array.isArray(list) ? list : []
  } catch (e) {
    console.error(e)
    showAlert('加载规则失败', 'error')
  }
}

function openCreateModal() {
  newRule.name = ''
  newRule.description = ''
  newRule.content = ''
  newRule.type = 'llm'
  showModal.value = true
}

async function createRule() {
  loading.value = true
  try {
    await api.createRule(newRule)
    showModal.value = false
    loadRules()
    showAlert('创建成功', 'success')
  } catch (e) {
    showAlert('创建失败', 'error')
  } finally {
    loading.value = false
  }
}

async function deleteRule(id: string) {
  if (!confirm('确认删除此规则？')) return
  try {
    await api.deleteRule(id)
    loadRules()
    showAlert('删除成功', 'success')
  } catch (e) {
    showAlert('删除失败', 'error')
  }
}

onMounted(loadRules)
</script>

<style scoped>
.page-container { padding: 24px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
.page-header h2 { margin: 0; font-size: 1.5rem; color: #1e293b; }

/* Filter Bar */
.filter-bar { display: flex; gap: 8px; margin-bottom: 16px; }
.filter-btn {
  background: white; border: 1px solid #cbd5e1; color: #64748b;
  padding: 6px 16px; border-radius: 20px; cursor: pointer; font-size: 0.9rem; transition: all 0.2s;
}
.filter-btn:hover { background: #f1f5f9; }
.filter-btn.active {
  background: #2563eb; color: white; border-color: #2563eb;
}

/* Table */
.table-container { background: white; border-radius: 8px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); overflow: hidden; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th, .data-table td { padding: 16px; text-align: left; border-bottom: 1px solid #e2e8f0; }
.data-table th { background: #f8fafc; font-weight: 600; color: #64748b; font-size: 0.875rem; }
.data-table tr:hover { background: #f8fafc; }

.font-bold { font-weight: 600; color: #1e293b; }
.text-desc { color: #64748b; font-size: 0.9rem; max-width: 300px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }

/* Badges */
.badge { display: inline-flex; align-items: center; padding: 2px 8px; border-radius: 4px; font-size: 0.75rem; font-weight: 500; }
.badge-llm { background: #e0e7ff; color: #4338ca; }
.badge-opa { background: #dcfce7; color: #15803d; }
.badge-kw { background: #fee2e2; color: #b91c1c; }
.badge-system { background: #f1f5f9; color: #475569; border: 1px solid #cbd5e1; }
.badge-custom { background: #f0f9ff; color: #0369a1; border: 1px solid #bae6fd; }

.actions { display: flex; gap: 8px; }
.btn-sm { padding: 4px 8px; font-size: 0.8rem; border-radius: 4px; cursor: pointer; text-decoration: none; }
.btn-outline { border: 1px solid #cbd5e1; background: white; color: #475569; }
.btn-danger { color: #dc2626; border-color: #fecaca; }
.btn-danger:hover { background: #fef2f2; }
.disabled { opacity: 0.5; cursor: not-allowed; }
.empty-state { text-align: center; padding: 40px; color: #94a3b8; }
.btn-link { background: none; border: none; color: #2563eb; cursor: pointer; text-decoration: underline; }

/* Buttons */
.btn-primary { background: #2563eb; color: white; border: none; padding: 8px 16px; border-radius: 6px; cursor: pointer; font-weight: 500; }
.btn-primary:hover { background: #1d4ed8; }
.btn-primary:disabled { opacity: 0.7; }
.btn-secondary { background: white; border: 1px solid #cbd5e1; color: #475569; padding: 8px 16px; border-radius: 6px; cursor: pointer; }

/* Modal */
.modal-overlay { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 1000; }
.modal { background: white; border-radius: 12px; box-shadow: 0 20px 25px -5px rgba(0,0,0,0.1); width: 400px; display: flex; flex-direction: column; max-height: 90vh; }
.modal-lg { width: 900px; }
.modal-header { padding: 20px; border-bottom: 1px solid #e2e8f0; display: flex; justify-content: space-between; align-items: center; }
.modal-header h3 { margin: 0; font-size: 1.25rem; }
.close-btn { background: none; border: none; font-size: 1.5rem; cursor: pointer; color: #94a3b8; }

.modal-body { padding: 20px; overflow-y: auto; }
.form-grid { display: grid; grid-template-columns: 1fr 1.5fr; gap: 24px; }
.form-col { display: flex; flex-direction: column; gap: 20px; }
.full-height { height: 100%; }

.form-group label { display: block; margin-bottom: 8px; font-weight: 500; font-size: 0.9rem; }
.form-input { width: 100%; padding: 8px 12px; border: 1px solid #cbd5e1; border-radius: 6px; font-size: 0.95rem; }
.code-editor { font-family: monospace; min-height: 300px; background: #f8fafc; line-height: 1.5; }

.type-selector { display: flex; flex-direction: column; gap: 8px; }
.type-option { padding: 10px; border: 1px solid #e2e8f0; border-radius: 8px; cursor: pointer; transition: all 0.2s; }
.type-option:hover { background: #f8fafc; }
.type-option.active { border-color: #2563eb; background: #eff6ff; color: #2563eb; box-shadow: 0 0 0 1px #2563eb; }

.modal-actions { padding: 20px; border-top: 1px solid #e2e8f0; display: flex; justify-content: flex-end; gap: 12px; }
</style>
