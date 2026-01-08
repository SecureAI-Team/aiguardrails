<template>
  <div class="page-container">
    <div class="page-header">
      <h2>工具能力管理 (Capabilities)</h2>
      <div class="header-actions">
        <div class="search-box">
          <input v-model="filterTag" placeholder="按标签筛选" @keyup.enter="load" class="search-input" />
          <button @click="load" class="btn-secondary">🔍</button>
        </div>
        <button @click="openCreateModal" class="btn-primary">+ 新建能力</button>
      </div>
    </div>

    <div class="table-container">
      <table class="data-table">
        <thead>
          <tr>
            <th>能力名称</th>
            <th>描述</th>
            <th>标签</th>
            <th>创建时间</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="c in caps" :key="c.id || c.name">
            <td class="font-bold">{{ c.name }}</td>
            <td class="desc-cell">{{ c.description || '-' }}</td>
            <td>
              <div class="tags">
                <span v-for="tag in (c.tags || [])" :key="tag" class="badge badge-info">{{ tag }}</span>
              </div>
            </td>
            <td class="text-sm font-mono">{{ formatDate(c.created_at) }}</td>
          </tr>
          <tr v-if="caps.length === 0">
            <td colspan="4" class="empty-state">暂无工具能力</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Create Modal -->
    <div v-if="showCreateModal" class="modal-overlay" @click.self="showCreateModal = false">
      <div class="modal">
        <h3>新建工具能力</h3>
        <form @submit.prevent="onCreate">
          <div class="form-group">
            <label>能力名称</label>
            <input v-model="form.name" placeholder="请输入能力名称 (如: 'web-browsing')" required />
            <div class="helper-text">通常对应 Agent 调用的工具函数名</div>
          </div>
          <div class="form-group">
            <label>描述</label>
            <textarea v-model="form.description" placeholder="描述该能力的作用"></textarea>
          </div>
          <div class="form-group">
            <label>标签</label>
            <input v-model="form.tagsInput" placeholder="输入标签，用逗号分隔 (如: 'network, search')" />
          </div>
          <div class="modal-actions">
            <button type="button" @click="showCreateModal = false" class="btn-secondary">取消</button>
            <button type="submit" class="btn-primary" :disabled="loading">创建</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { api } from '../services/api'

const caps = ref<any[]>([])
const loading = ref(false)
const filterTag = ref('')
const showCreateModal = ref(false)

const form = reactive({
  name: '',
  description: '',
  tagsInput: ''
})

async function load() {
  loading.value = true
  try {
    caps.value = await api.listCapabilities(filterTag.value || undefined)
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function openCreateModal() {
  form.name = ''
  form.description = ''
  form.tagsInput = ''
  showCreateModal.value = true
}

async function onCreate() {
  if (!form.name) return
  loading.value = true
  try {
    await api.createCapability({
      name: form.name,
      description: form.description,
      tags: form.tagsInput.split(',').map(t => t.trim()).filter(Boolean)
    })
    showCreateModal.value = false
    await load()
  } catch (e) {
    alert('创建失败')
  } finally {
    loading.value = false
  }
}

function formatDate(dateStr?: string) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString()
}

onMounted(load)
</script>

<style scoped>
.page-container { padding: 24px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
.page-header h2 { margin: 0; font-size: 1.5rem; color: #1e293b; }
.header-actions { display: flex; gap: 12px; }

.search-box { display: flex; gap: 4px; }
.search-input { padding: 8px 12px; border: 1px solid #cbd5e1; border-radius: 6px 0 0 6px; border-right: none; }
.search-input + button { border-radius: 0 6px 6px 0; }

.table-container { background: white; border-radius: 8px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); overflow: hidden; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th, .data-table td { padding: 16px; text-align: left; border-bottom: 1px solid #e2e8f0; }
.data-table th { background: #f8fafc; font-weight: 600; color: #64748b; font-size: 0.875rem; }
.data-table tr:hover { background: #f8fafc; }

.font-bold { font-weight: 600; color: #1e293b; }
.desc-cell { color: #64748b; max-width: 400px; }
.text-sm { font-size: 0.875rem; }
.font-mono { font-family: monospace; color: #64748b; }

.badge { padding: 2px 6px; border-radius: 4px; font-size: 0.75rem; margin-right: 4px; display: inline-block; }
.badge-info { background: #e0f2fe; color: #075985; }

.btn-primary { background: #2563eb; color: white; border: none; padding: 8px 16px; border-radius: 6px; cursor: pointer; }
.btn-secondary { background: white; border: 1px solid #cbd5e1; color: #475569; padding: 8px 16px; border-radius: 6px; cursor: pointer; }

/* Modal */
.modal-overlay { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 1000; }
.modal { background: white; padding: 24px; border-radius: 12px; width: 450px; box-shadow: 0 20px 25px -5px rgba(0,0,0,0.1); }
.modal h3 { margin: 0 0 20px; font-size: 1.25rem; }

.form-group { margin-bottom: 20px; }
.form-group label { display: block; margin-bottom: 8px; font-weight: 500; }
.form-group input, .form-group textarea { width: 100%; padding: 8px; border: 1px solid #cbd5e1; border-radius: 6px; font-family: inherit; }
.form-group textarea { height: 80px; resize: vertical; }
.helper-text { font-size: 0.8rem; color: #64748b; margin-top: 4px; }
.modal-actions { display: flex; justify-content: flex-end; gap: 12px; }
.empty-state { text-align: center; padding: 40px; color: #94a3b8; }
</style>
