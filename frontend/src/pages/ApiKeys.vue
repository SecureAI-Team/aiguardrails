<template>
  <div class="page-container">
    <div class="page-header">
      <h2>🔑 API 密钥概览</h2>
      <div class="header-actions">
        <button @click="loadData" class="btn-outline">🔄 刷新</button>
        <router-link to="/apps" class="btn-primary">+ 创建应用密钥</router-link>
      </div>
    </div>

    <div class="info-card">
      <p>💡 此页面汇总了所有租户下应用的 API 密钥。如需创建新密钥，请前往 <strong>应用管理</strong> 页面。</p>
    </div>

    <div class="table-container">
      <table class="data-table">
        <thead>
          <tr>
            <th>应用名称</th>
            <th>所属租户</th>
            <th>API Key 前缀</th>
            <th>状态</th>
            <th>创建时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="key in allKeys" :key="key.id">
            <td class="font-bold">{{ key.app_name }}</td>
            <td><span class="badge badge-gray">{{ key.tenant_name }}</span></td>
            <td class="font-mono">{{ key.key_prefix }}••••••••</td>
            <td>
              <span :class="['badge', key.is_active ? 'badge-success' : 'badge-danger']">
                {{ key.is_active ? '有效' : '已禁用' }}
              </span>
            </td>
            <td class="text-sm text-gray">{{ formatDate(key.created_at) }}</td>
            <td>
              <router-link :to="{ path: '/apps', query: { tenantId: key.tenant_id } }" class="btn-xs btn-link">
                管理
              </router-link>
            </td>
          </tr>
          <tr v-if="loading">
            <td colspan="6" class="text-center py-4">加载中...</td>
          </tr>
          <tr v-if="!loading && allKeys.length === 0">
            <td colspan="6" class="empty-state">暂无应用密钥</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '../services/api'

interface AggregatedKey {
  id: string
  app_name: string
  tenant_id: string
  tenant_name: string
  key_prefix: string
  is_active: boolean
  created_at: string
}

const allKeys = ref<AggregatedKey[]>([])
const loading = ref(false)

async function loadData() {
  loading.value = true
  allKeys.value = []
  try {
    const tenants = await api.listTenants()
    
    // Fetch apps for all tenants in parallel
    const promises = tenants.map(async (t: any) => {
      try {
        const apps = await api.listApps(t.id)
        return apps.map((app: any) => ({
          id: app.id,
          app_name: app.name,
          tenant_id: t.id,
          tenant_name: t.name,
          key_prefix: app.api_key_prefix || (app.api_key ? app.api_key.substring(0, 8) : 'sk-'),
          is_active: app.is_active !== false, // Assume true if undefined
          created_at: app.created_at
        }))
      } catch {
        return []
      }
    })
    
    const results = await Promise.all(promises)
    allKeys.value = results.flat().sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
    
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function formatDate(d?: string) {
  if (!d) return '-'
  return new Date(d).toLocaleString()
}

onMounted(loadData)
</script>

<style scoped>
.page-container { padding: 24px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
.page-header h2 { margin: 0; font-size: 1.5rem; color: #1e293b; }
.header-actions { display: flex; gap: 12px; }

.info-card { background: #eff6ff; padding: 12px 16px; border-radius: 8px; margin-bottom: 24px; color: #1e40af; font-size: 0.95rem; border: 1px solid #dbeafe; }

.table-container { background: white; border-radius: 8px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); overflow: hidden; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th, .data-table td { padding: 16px; text-align: left; border-bottom: 1px solid #e2e8f0; }
.data-table th { background: #f8fafc; font-weight: 600; color: #64748b; font-size: 0.875rem; }
.data-table tr:hover { background: #f8fafc; }

.font-bold { font-weight: 600; color: #1e293b; }
.font-mono { font-family: monospace; color: #64748b; }
.text-gray { color: #64748b; }
.text-center { text-align: center; }
.py-4 { padding-top: 16px; padding-bottom: 16px; }

.badge { padding: 2px 8px; border-radius: 12px; font-size: 0.75rem; font-weight: 500; }
.badge-gray { background: #f1f5f9; color: #475569; }
.badge-success { background: #dcfce7; color: #166534; }
.badge-danger { background: #fee2e2; color: #991b1b; }

.btn-primary { background: #2563eb; color: white; border: none; padding: 8px 16px; border-radius: 6px; cursor: pointer; text-decoration: none; display: inline-flex; align-items: center; }
.btn-outline { background: white; border: 1px solid #cbd5e1; color: #475569; padding: 8px 16px; border-radius: 6px; cursor: pointer; display: inline-flex; align-items: center; }
.btn-link { background: none; border: none; color: #2563eb; cursor: pointer; text-decoration: underline; padding: 0; }
.btn-xs { font-size: 0.85rem; }

.empty-state { text-align: center; padding: 40px; color: #94a3b8; }
</style>
