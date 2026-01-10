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
              <button v-if="!rule.is_system" @click="openEditModal(rule)" class="btn-sm btn-outline">编辑</button>
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

    <!-- Create/Edit Modal -->
    <Teleport to="body">
      <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
        <div class="modal">
            <div class="modal-header">
              <h3>{{ isEdit ? '编辑规则' : '新建规则' }}</h3>
              <button @click="showModal = false" class="close-btn">&times;</button>
            </div>
            
            <div class="modal-body">
                <form @submit.prevent="submitRule">
                  <!-- Name -->
                  <div class="form-group">
                    <label>规则名称 <span class="required">*</span></label>
                    <input v-model="newRule.name" type="text" placeholder="例如: 拒绝竞品信息" required class="form-input" />
                  </div>

                  <!-- Type -->
                  <div class="form-group">
                    <label>规则类型</label>
                    <div class="type-selector">
                      <div 
                        @click="newRule.type = 'llm'"
                        :class="['type-option', newRule.type === 'llm' ? 'active' : '']"
                      >
                        <span class="type-icon">🤖</span>
                        <div class="type-info">
                          <div class="type-title">LLM Security</div>
                          <div class="type-desc">语义检测</div>
                        </div>
                      </div>
                      
                      <div 
                        @click="newRule.type = 'opa'"
                        :class="['type-option', newRule.type === 'opa' ? 'active' : '']"
                      >
                        <span class="type-icon">📜</span>
                        <div class="type-info">
                          <div class="type-title">OPA Policy</div>
                          <div class="type-desc">逻辑代码</div>
                        </div>
                      </div>

                      <div 
                        @click="newRule.type = 'keyword'"
                        :class="['type-option', newRule.type === 'keyword' ? 'active' : '']"
                      >
                        <span class="type-icon">🚫</span>
                        <div class="type-info">
                          <div class="type-title">Keyword List</div>
                          <div class="type-desc">敏感词库</div>
                        </div>
                      </div>
                    </div>
                  </div>

                  <!-- Description -->
                  <div class="form-group">
                    <label>描述</label>
                    <textarea v-model="newRule.description" rows="2" placeholder="简要描述规则的用途..." class="form-input"></textarea>
                  </div>

                  <!-- Content -->
                  <div class="form-group full-height">
                    <label>
                      <span v-if="newRule.type === 'llm'">System Prompt 指令</span>
                      <span v-else-if="newRule.type === 'keyword'">敏感词列表 (每行一个)</span>
                      <span v-else>Rego 策略代码</span>
                      <span class="required">*</span>
                    </label>
                    <textarea
                      v-model="newRule.content"
                      class="form-input code-editor"
                      :placeholder="placeholderText"
                      spellcheck="false"
                      required
                    ></textarea>
                    <div class="tip-box">
                      <span v-if="newRule.type === 'llm'">该指令将嵌入到 System Prompt 中，用于指导大模型进行安全拦截。请使用清晰的自然语言描述。</span>
                      <span v-else-if="newRule.type === 'keyword'">输入需要拦截的敏感词汇，每行一个。此类规则同时适用于输入 Prompt 和输出 Response 检测。</span>
                      <span v-else>使用 Open Policy Agent (OPA) 的 Rego 语言编写复杂策略。必须包含 'default allow' 和 'deny' 规则。</span>
                    </div>
                  </div>
                
                  <!-- Actions -->
                  <div class="modal-actions">
                    <button type="button" @click="showModal = false" class="btn-secondary">取消</button>
                    <button type="submit" class="btn-primary" :disabled="loading || !newRule.name || !newRule.content">
                      {{ loading ? '处理中...' : (isEdit ? '确认保存' : '确认创建') }}
                    </button>
                  </div>
                </form>
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
const isEdit = ref(false)
const editingId = ref('')

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
  isEdit.value = false
  editingId.value = ''
  newRule.name = ''
  newRule.description = ''
  newRule.content = ''
  newRule.type = 'llm'
  showModal.value = true
}

function openEditModal(rule: Rule) {
  isEdit.value = true
  editingId.value = rule.id
  newRule.name = rule.name
  newRule.description = rule.description
  newRule.content = rule.content
  newRule.type = rule.type
  showModal.value = true
}

async function submitRule() {
  loading.value = true
  try {
    if (isEdit.value) {
      await api.updateRule(editingId.value, newRule)
      showAlert('更新成功', 'success')
    } else {
      await api.createRule(newRule)
      showAlert('创建成功', 'success')
    }
    showModal.value = false
    loadRules()
  } catch (e: any) {
    console.error(e)
    const msg = e.response?.data?.error || e.message || '未知错误'
    showAlert((isEdit.value ? '更新' : '创建') + '失败: ' + msg, 'error')
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
  width: 600px; /* Reduced width for simpler look */
  max-width: 95%; 
  display: flex; 
  flex-direction: column; 
  max-height: 90vh; 
}

.modal-header { padding: 20px; border-bottom: 1px solid #e2e8f0; display: flex; justify-content: space-between; align-items: center; }
.modal-header h3 { margin: 0; font-size: 1.25rem; }
.close-btn { background: none; border: none; font-size: 1.5rem; cursor: pointer; color: #94a3b8; }

.modal-body { padding: 24px; overflow-y: auto; display: flex; flex-direction: column; gap: 20px; }

.form-group label { display: block; margin-bottom: 8px; font-weight: 500; font-size: 0.9rem; color: #374151; }
.required { color: #ef4444; }
.form-input { width: 100%; padding: 10px 12px; border: 1px solid #cbd5e1; border-radius: 6px; font-size: 0.95rem; transition: border 0.2s; box-sizing: border-box; }
.form-input:focus { border-color: #2563eb; outline: none; ring: 1px #2563eb; }

.code-editor { font-family: monospace; min-height: 200px; background: #f8fafc; line-height: 1.5; }

/* Type Selector */
.type-selector { display: flex; gap: 12px; }
.type-option { 
  flex: 1; padding: 12px; border: 1px solid #e2e8f0; border-radius: 8px; 
  cursor: pointer; transition: all 0.2s; display: flex; align-items: center; gap: 10px;
}
.type-option:hover { background: #f8fafc; border-color: #cbd5e1; }
.type-option.active { border-color: #2563eb; background: #eff6ff; color: #2563eb; box-shadow: 0 0 0 1px #2563eb; }
.type-icon { font-size: 1.5rem; }
.type-info { display: flex; flex-direction: column; }
.type-title { font-weight: 600; font-size: 0.9rem; }
.type-desc { font-size: 0.75rem; color: inherit; opacity: 0.8; }

.tip-box { margin-top: 8px; font-size: 0.85rem; color: #64748b; background: #f8fafc; padding: 10px; border-radius: 6px; }

.modal-actions { display: flex; justify-content: flex-end; gap: 12px; margin-top: 20px; }
</style>
