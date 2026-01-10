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
            <!-- Modal Header -->
            <div class="px-6 py-4 border-b border-gray-100 flex justify-between items-center">
              <div>
                <h3 class="text-xl font-bold text-gray-900">新建规则</h3>
                <p class="mt-1 text-sm text-gray-500">配置安全护栏规则以拦截风险内容。</p>
              </div>
              <button @click="showModal = false" class="rounded-full p-1 bg-gray-50 hover:bg-gray-100 text-gray-400 hover:text-gray-500 transition">
                <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            
            <div class="px-6 py-6 bg-gray-50/30 flex-1 overflow-hidden flex flex-col">
              <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 h-full">
                <!-- Left Column: Settings -->
                <div class="lg:col-span-4 space-y-6 overflow-y-auto pr-2">
                  <div>
                    <label class="block text-sm font-semibold text-gray-700 mb-2">规则名称 <span class="text-red-500">*</span></label>
                    <input
                      v-model="newRule.name"
                      type="text"
                      class="block w-full rounded-lg border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm px-4 py-2.5 transition"
                      placeholder="例如: 拒绝竞品信息"
                    />
                  </div>

                  <div>
                    <label class="block text-sm font-semibold text-gray-700 mb-2">规则类型</label>
                    <div class="grid grid-cols-1 gap-3">
                      <!-- LLM Type -->
                      <div 
                        @click="newRule.type = 'llm'"
                        :class="[
                          newRule.type === 'llm' ? 'border-indigo-500 ring-1 ring-indigo-500 bg-indigo-50/50' : 'border-gray-200 hover:border-gray-300 bg-white shadow-sm',
                          'cursor-pointer relative flex items-center px-4 py-3 border rounded-xl transition-all duration-200'
                        ]"
                      >
                         <div class="flex-shrink-0 h-10 w-10 flex items-center justify-center rounded-lg" :class="newRule.type === 'llm' ? 'bg-indigo-100 text-indigo-600' : 'bg-gray-100 text-gray-500'">
                           <span class="text-xl">🤖</span>
                         </div>
                         <div class="ml-3">
                           <p class="text-sm font-medium text-gray-900">LLM Security</p>
                           <p class="text-xs text-gray-500">语义检测</p>
                         </div>
                      </div>

                      <!-- OPA Type -->
                      <div 
                        @click="newRule.type = 'opa'"
                        :class="[
                          newRule.type === 'opa' ? 'border-green-500 ring-1 ring-green-500 bg-green-50/50' : 'border-gray-200 hover:border-gray-300 bg-white shadow-sm',
                          'cursor-pointer relative flex items-center px-4 py-3 border rounded-xl transition-all duration-200'
                        ]"
                      >
                         <div class="flex-shrink-0 h-10 w-10 flex items-center justify-center rounded-lg" :class="newRule.type === 'opa' ? 'bg-green-100 text-green-600' : 'bg-gray-100 text-gray-500'">
                           <span class="text-xl">📜</span>
                         </div>
                         <div class="ml-3">
                           <p class="text-sm font-medium text-gray-900">OPA Policy</p>
                           <p class="text-xs text-gray-500">逻辑代码</p>
                         </div>
                      </div>

                      <!-- Keyword Type -->
                      <div 
                        @click="newRule.type = 'keyword'"
                        :class="[
                          newRule.type === 'keyword' ? 'border-red-500 ring-1 ring-red-500 bg-red-50/50' : 'border-gray-200 hover:border-gray-300 bg-white shadow-sm',
                          'cursor-pointer relative flex items-center px-4 py-3 border rounded-xl transition-all duration-200'
                        ]"
                      >
                         <div class="flex-shrink-0 h-10 w-10 flex items-center justify-center rounded-lg" :class="newRule.type === 'keyword' ? 'bg-red-100 text-red-600' : 'bg-gray-100 text-gray-500'">
                           <span class="text-xl">🚫</span>
                         </div>
                         <div class="ml-3">
                           <p class="text-sm font-medium text-gray-900">Keyword List</p>
                           <p class="text-xs text-gray-500">敏感词库</p>
                         </div>
                      </div>
                    </div>
                  </div>

                  <div>
                    <label class="block text-sm font-semibold text-gray-700 mb-2">描述</label>
                    <textarea
                      v-model="newRule.description"
                      rows="3"
                      class="block w-full rounded-lg border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm px-4 py-2"
                      placeholder="简要描述规则的用途..."
                    ></textarea>
                  </div>
                </div>

                <!-- Right Column: Content -->
                <div class="lg:col-span-8 flex flex-col h-full overflow-hidden">
                  <div class="flex items-center justify-between mb-2">
                    <label class="block text-sm font-semibold text-gray-700">
                      <span v-if="newRule.type === 'llm'">System Prompt 指令</span>
                      <span v-else-if="newRule.type === 'keyword'">敏感词列表 (每行一个)</span>
                      <span v-else>Rego 策略代码</span>
                       <span class="text-red-500">*</span>
                    </label>
                    <span class="inline-flex items-center rounded-full bg-blue-50 px-2.5 py-0.5 text-xs font-medium text-blue-700 border border-blue-100">
                       {{ newRule.type === 'llm' ? 'Natural Language' : newRule.type === 'keyword' ? 'Line separated' : 'Rego' }}
                    </span>
                  </div>
                  
                  <div class="relative flex-1 rounded-lg border border-gray-300 shadow-sm overflow-hidden focus-within:ring-1 focus-within:ring-indigo-500 focus-within:border-indigo-500 bg-white">
                    <textarea
                      v-model="newRule.content"
                      class="block w-full h-full resize-none border-0 p-4 font-mono text-sm leading-relaxed focus:ring-0"
                      :placeholder="placeholderText"
                      spellcheck="false"
                    ></textarea>
                  </div>
                  
                  <!-- Tip Box -->
                  <div class="mt-4 rounded-lg bg-gray-50 p-3 border border-gray-200 flex items-start text-xs text-gray-600">
                    <svg class="h-4 w-4 text-gray-400 mr-2 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                       <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    <span v-if="newRule.type === 'llm'">该指令将嵌入到 System Prompt 中，用于指导大模型进行安全拦截。请使用清晰的自然语言描述。</span>
                    <span v-else-if="newRule.type === 'keyword'">输入需要拦截的敏感词汇，每行一个。支持精确匹配。</span>
                    <span v-else>使用 Open Policy Agent (OPA) 的 Rego 语言编写复杂策略。必须包含 'default allow' 和 'deny' 规则。</span>
                  </div>
                </div>
              </div>
            </div>

            <div class="bg-gray-50 px-6 py-4 flex justify-end space-x-3 border-t border-gray-200">
              <button
                @click="showModal = false"
                class="px-5 py-2.5 bg-white border border-gray-300 rounded-lg text-gray-700 font-medium hover:bg-gray-50 transition shadow-sm"
              >
                取消
              </button>
              <button
                @click="createRule"
                :disabled="loading || !newRule.name || !newRule.content"
                class="px-5 py-2.5 bg-indigo-600 rounded-lg text-white font-medium hover:bg-indigo-700 transition shadow-md disabled:opacity-50 disabled:cursor-not-allowed flex items-center"
              >
                <svg v-if="loading" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
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
  } catch (e: any) {
    console.error(e)
    const msg = e.response?.data?.error || e.message || '未知错误'
    showAlert('创建失败: ' + msg, 'error')
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
.modal { 
  background: white; 
  border-radius: 12px; 
  box-shadow: 0 20px 25px -5px rgba(0,0,0,0.1); 
  width: 90%; 
  max-width: 900px; 
  display: flex; 
  flex-direction: column; 
  max-height: 90vh; 
  margin: 0 20px;
}
.modal-lg { max-width: 900px; }
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
