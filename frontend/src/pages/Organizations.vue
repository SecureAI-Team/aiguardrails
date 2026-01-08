<template>
  <div class="orgs-page">
    <div class="page-header">
      <h2>🏢 组织管理</h2>
      <button @click="showCreate = true" class="btn-primary">+ 创建组织</button>
    </div>

    <div class="orgs-grid">
      <div v-for="o in orgs" :key="o.id" class="org-card" @click="selectOrg(o)">
        <div class="org-logo">{{ o.name.charAt(0) }}</div>
        <div class="org-info">
          <h3>{{ o.name }}</h3>
          <p class="org-slug">{{ o.slug }}</p>
          <p class="org-desc">{{ o.description }}</p>
        </div>
        <div class="org-badges">
          <span v-if="o.sso_enabled" class="badge sso">SSO</span>
          <span class="badge status">{{ o.status }}</span>
        </div>
      </div>
    </div>

    <!-- Create Modal -->
    <div v-if="showCreate" class="modal-overlay" @click.self="showCreate = false">
      <div class="modal">
        <h3>创建组织</h3>
        <div class="form-group">
          <label>组织名称</label>
          <input v-model="form.name" placeholder="如: 我的公司" />
        </div>
        <div class="form-group">
          <label>唯一标识 (URL)</label>
          <input v-model="form.slug" placeholder="my-company" />
        </div>
        <div class="form-group">
          <label>描述</label>
          <textarea v-model="form.description" rows="3" placeholder="组织描述"></textarea>
        </div>
        <div class="modal-actions">
          <button @click="showCreate = false" class="btn-outline">取消</button>
          <button @click="createOrg" class="btn-primary">创建</button>
        </div>
      </div>
    </div>

    <!-- Org Detail -->
    <div v-if="selectedOrg" class="modal-overlay" @click.self="selectedOrg = null">
      <div class="modal detail-modal">
        <div class="modal-header">
          <h3>{{ selectedOrg.name }}</h3>
          <button @click="selectedOrg = null" class="btn-close">×</button>
        </div>
        
        <div class="tabs">
          <button :class="{ active: activeTab === 'teams' }" @click="activeTab = 'teams'">团队</button>
          <button :class="{ active: activeTab === 'members' }" @click="activeTab = 'members'">成员</button>
          <button :class="{ active: activeTab === 'whitelist' }" @click="activeTab = 'whitelist'">IP白名单</button>
          <button :class="{ active: activeTab === 'sso' }" @click="activeTab = 'sso'">SSO配置</button>
        </div>

        <div v-if="activeTab === 'teams'" class="tab-content">
          <div class="section-header">
            <h4>团队列表</h4>
            <button @click="showAddTeam = true" class="btn-sm">+ 新建团队</button>
          </div>
          <div v-for="t in teams" :key="t.id" class="list-item">
            <span class="item-name">{{ t.name }}</span>
            <span class="item-role">{{ t.default_role }}</span>
          </div>
        </div>

        <div v-if="activeTab === 'members'" class="tab-content">
          <div class="section-header">
            <h4>成员列表</h4>
            <button @click="showAddMember = true" class="btn-sm">+ 邀请成员</button>
          </div>
          <div v-for="m in members" :key="m.id" class="list-item">
            <span class="item-name">{{ m.username || m.user_id }}</span>
            <span class="item-role" :class="m.role">{{ m.role }}</span>
            <span class="item-status">{{ m.status }}</span>
          </div>
        </div>

        <div v-if="activeTab === 'whitelist'" class="tab-content">
          <div class="section-header">
            <h4>IP白名单</h4>
            <button @click="showAddIP = true" class="btn-sm">+ 添加IP</button>
          </div>
          <div v-for="w in whitelist" :key="w.id" class="list-item">
            <span class="item-ip">{{ w.ip_address || w.ip_cidr }}</span>
            <span class="item-desc">{{ w.description }}</span>
            <button @click="deleteIP(w.id)" class="btn-delete">删除</button>
          </div>
          <p v-if="whitelist.length === 0" class="empty-hint">未设置IP白名单，允许所有IP访问</p>
        </div>

        <div v-if="activeTab === 'sso'" class="tab-content">
          <div class="sso-config">
            <div class="form-group">
              <label>SSO类型</label>
              <select v-model="ssoConfig.provider">
                <option value="">未启用</option>
                <option value="saml">SAML 2.0</option>
                <option value="oidc">OIDC</option>
              </select>
            </div>
            <div v-if="ssoConfig.provider === 'saml'" class="form-group">
              <label>SAML元数据URL</label>
              <input v-model="ssoConfig.metadata_url" placeholder="https://idp.example.com/metadata" />
            </div>
            <button @click="saveSSO" class="btn-primary">保存SSO配置</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { api } from '../services/api'

interface Organization {
  id: string
  name: string
  slug: string
  description: string
  sso_enabled: boolean
  status: string
}

const orgs = ref<Organization[]>([])
const selectedOrg = ref<Organization | null>(null)
const showCreate = ref(false)
const activeTab = ref('teams')
const form = ref({ name: '', slug: '', description: '' })
const teams = ref<any[]>([])
const members = ref<any[]>([])
const whitelist = ref<any[]>([])
const ssoConfig = ref({ provider: '', metadata_url: '' })
const showAddTeam = ref(false)
const showAddMember = ref(false)
const showAddIP = ref(false)

onMounted(() => loadOrgs())

watch(selectedOrg, (org) => {
  if (org) {
    loadTeams(org.id)
    loadMembers(org.id)
    loadWhitelist(org.id)
  }
})

async function loadOrgs() {
  try { 
    const result = await api.get('/orgs')
    orgs.value = Array.isArray(result) ? result : []
  } catch { 
    orgs.value = [] 
  }
}

async function createOrg() {
  try {
    await api.post('/orgs', form.value)
    showCreate.value = false
    form.value = { name: '', slug: '', description: '' }
    loadOrgs()
  } catch {}
}

function selectOrg(org: Organization) {
  selectedOrg.value = org
  activeTab.value = 'teams'
}

async function loadTeams(orgId: string) {
  try { teams.value = await api.get(`/orgs/${orgId}/teams`) } catch { teams.value = [] }
}

async function loadMembers(orgId: string) {
  try { members.value = await api.get(`/orgs/${orgId}/members`) } catch { members.value = [] }
}

async function loadWhitelist(orgId: string) {
  try { whitelist.value = await api.get(`/orgs/${orgId}/whitelist`) } catch { whitelist.value = [] }
}

async function deleteIP(id: string) {
  if (!confirm('确定删除此IP?')) return
  try {
    await api.delete(`/whitelist/${id}`)
    if (selectedOrg.value) loadWhitelist(selectedOrg.value.id)
  } catch {}
}

async function saveSSO() {
  alert('SSO配置已保存(演示)')
}
</script>

<style scoped>
.orgs-page { padding: 20px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
.page-header h2 { margin: 0; }
.btn-primary { background: linear-gradient(90deg, #2563eb, #7c3aed); color: white; padding: 10px 20px; border: none; border-radius: 8px; cursor: pointer; }
.orgs-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 20px; }
.org-card { background: white; padding: 20px; border-radius: 12px; box-shadow: 0 2px 10px rgba(0,0,0,0.05); cursor: pointer; display: flex; gap: 16px; align-items: center; transition: transform 0.2s; }
.org-card:hover { transform: translateY(-2px); }
.org-logo { width: 50px; height: 50px; background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%); color: white; border-radius: 12px; display: flex; align-items: center; justify-content: center; font-size: 1.5rem; font-weight: 700; }
.org-info { flex: 1; }
.org-info h3 { margin: 0 0 4px; }
.org-slug { color: #64748b; font-family: monospace; margin: 0 0 4px; font-size: 0.85rem; }
.org-desc { color: #64748b; margin: 0; font-size: 0.9rem; }
.org-badges { display: flex; flex-direction: column; gap: 4px; }
.badge { padding: 4px 10px; border-radius: 20px; font-size: 0.75rem; text-align: center; }
.badge.sso { background: #dbeafe; color: #1e40af; }
.badge.status { background: #d1fae5; color: #065f46; }
.modal-overlay { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 1000; }
.modal { background: white; padding: 24px; border-radius: 12px; width: 500px; max-height: 80vh; overflow-y: auto; }
.detail-modal { width: 700px; }
.modal-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.modal-header h3 { margin: 0; }
.btn-close { background: none; border: none; font-size: 1.5rem; cursor: pointer; color: #64748b; }
.form-group { margin-bottom: 16px; }
.form-group label { display: block; margin-bottom: 6px; color: #374151; font-weight: 500; }
.form-group input, .form-group textarea, .form-group select { width: 100%; padding: 10px; border: 1px solid #e2e8f0; border-radius: 6px; box-sizing: border-box; }
.modal-actions { display: flex; justify-content: flex-end; gap: 12px; margin-top: 24px; }
.btn-outline { padding: 10px 20px; border: 1px solid #e2e8f0; background: white; border-radius: 8px; cursor: pointer; }
.tabs { display: flex; gap: 8px; margin-bottom: 20px; border-bottom: 1px solid #e2e8f0; padding-bottom: 12px; }
.tabs button { padding: 8px 16px; border: none; background: none; cursor: pointer; color: #64748b; border-radius: 6px; }
.tabs button.active { background: #3b82f6; color: white; }
.tab-content { min-height: 200px; }
.section-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.section-header h4 { margin: 0; }
.btn-sm { padding: 6px 12px; background: #3b82f6; color: white; border: none; border-radius: 6px; cursor: pointer; font-size: 0.85rem; }
.list-item { display: flex; align-items: center; gap: 12px; padding: 12px; background: #f8fafc; border-radius: 6px; margin-bottom: 8px; }
.item-name { flex: 1; font-weight: 500; }
.item-role { padding: 4px 10px; background: #e0f2fe; color: #0369a1; border-radius: 20px; font-size: 0.75rem; }
.item-role.owner { background: #fef3c7; color: #92400e; }
.item-role.admin { background: #ede9fe; color: #5b21b6; }
.item-status { color: #64748b; font-size: 0.85rem; }
.item-ip { font-family: monospace; }
.item-desc { flex: 1; color: #64748b; }
.btn-delete { padding: 4px 10px; background: #fee2e2; color: #991b1b; border: none; border-radius: 4px; cursor: pointer; font-size: 0.8rem; }
.empty-hint { color: #64748b; text-align: center; padding: 20px; }
.sso-config { padding: 16px; background: #f8fafc; border-radius: 8px; }
</style>
