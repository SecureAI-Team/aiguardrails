<template>
  <div class="page-container">
    <div class="page-header">
      <h2>规则库</h2>
    </div>

    <div class="filter-bar">
      <input v-model="filters.jurisdiction" placeholder="司法管辖区 (如 EU)" class="filter-input" />
      <input v-model="filters.regulation" placeholder="法规 (如 GDPR)" class="filter-input" />
      <input v-model="filters.tag" placeholder="标签" class="filter-input" />
      <select v-model="filters.severity" class="filter-select">
        <option value="">全部严重性</option>
        <option value="low">低</option>
        <option value="medium">中</option>
        <option value="high">高</option>
        <option value="critical">严重</option>
      </select>
      <button @click="load" class="btn-primary" :disabled="loading">搜索</button>
    </div>

    <div class="table-container">
      <table class="data-table">
        <thead>
          <tr>
            <th>规则名称</th>
            <th>法规/管辖区</th>
            <th>厂商/产品</th>
            <th>严重性</th>
            <th>处置动作</th>
            <th>描述</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="r in rules" :key="r.id">
            <td class="font-bold">{{ r.name }}</td>
            <td>
              <div v-if="r.regulation || r.jurisdiction">
                <span v-if="r.regulation" class="badge badge-info">{{ r.regulation }}</span>
                <span v-if="r.jurisdiction" class="badge badge-warning">{{ r.jurisdiction }}</span>
              </div>
              <span v-else>-</span>
            </td>
            <td>
              <div v-if="r.vendor">{{ r.vendor }} {{ r.product }}</div>
              <span v-else>-</span>
            </td>
            <td>
              <span :class="['badge', getSeverityClass(r.severity)]">{{ r.severity }}</span>
            </td>
            <td>
              <span :class="['badge', r.decision === 'block' ? 'badge-danger' : 'badge-success']">{{ r.decision || 'block' }}</span>
            </td>
            <td class="desc-cell" :title="r.description">{{ r.description }}</td>
          </tr>
          <tr v-if="rules.length === 0">
            <td colspan="6" class="empty-state">未找到匹配的规则</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { api } from '../services/api'

const rules = ref<any[]>([])
const loading = ref(false)

const filters = reactive({
  jurisdiction: '',
  regulation: '',
  vendor: '',
  product: '',
  tag: '',
  severity: '',
  decision: ''
})

async function load() {
  loading.value = true
  try {
    const result = await api.listRules(filters)
    rules.value = Array.isArray(result) ? result : []
  } catch (e) {
    console.error(e)
    rules.value = []
  } finally {
    loading.value = false
  }
}

function getSeverityClass(severity: string) {
  switch (severity) {
    case 'critical': return 'badge-danger'
    case 'high': return 'badge-danger'
    case 'medium': return 'badge-warning'
    case 'low': return 'badge-info'
    default: return 'badge-secondary'
  }
}

onMounted(load)
</script>

<style scoped>
.page-container { padding: 24px; }
.page-header { margin-bottom: 24px; }
.page-header h2 { margin: 0; font-size: 1.5rem; color: #1e293b; }

.filter-bar { display: flex; flex-wrap: wrap; gap: 12px; margin-bottom: 20px; background: white; padding: 16px; border-radius: 8px; box-shadow: 0 1px 2px rgba(0,0,0,0.05); }
.filter-input, .filter-select { padding: 8px 12px; border: 1px solid #cbd5e1; border-radius: 6px; font-size: 0.95rem; }
.filter-input { width: 180px; }

.table-container { background: white; border-radius: 8px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); overflow: hidden; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th, .data-table td { padding: 16px; text-align: left; border-bottom: 1px solid #e2e8f0; }
.data-table th { background: #f8fafc; font-weight: 600; color: #64748b; font-size: 0.875rem; }
.data-table tr:hover { background: #f8fafc; }

.font-bold { font-weight: 600; color: #1e293b; }
.desc-cell { max-width: 300px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; color: #64748b; }

.badge { padding: 2px 6px; border-radius: 4px; font-size: 0.75rem; font-weight: 500; margin-right: 4px; display: inline-block; }
.badge-info { background: #e0f2fe; color: #075985; }
.badge-warning { background: #fef9c3; color: #854d0e; }
.badge-danger { background: #fee2e2; color: #991b1b; }
.badge-success { background: #dcfce7; color: #166534; }
.badge-secondary { background: #f1f5f9; color: #64748b; }

.btn-primary { background: #2563eb; color: white; border: none; padding: 8px 16px; border-radius: 6px; cursor: pointer; }
.btn-primary:hover { background: #1d4ed8; }
.btn-primary:disabled { opacity: 0.7; }

.empty-state { text-align: center; padding: 40px; color: #94a3b8; }
</style>
