<template>
  <div class="alert-rules-page">
    <div class="page-header">
      <h2>🔔 告警规则配置</h2>
      <button @click="showCreate = true" class="btn-primary">+ 新建规则</button>
    </div>

    <div class="rules-list">
      <div v-for="rule in rules" :key="rule.id" class="rule-card" :class="{ disabled: !rule.enabled }">
        <div class="rule-header">
          <span class="rule-name">{{ rule.name }}</span>
          <span :class="'severity severity-' + rule.severity_threshold">{{ rule.severity_threshold }}</span>
        </div>
        <p class="rule-desc">{{ rule.description }}</p>
        <div class="rule-config">
          <span class="config-item">📊 阈值: {{ rule.threshold_count }}次 / {{ rule.threshold_window_sec }}秒</span>
          <span class="config-item">⏰ 冷却: {{ rule.cooldown_sec }}秒</span>
        </div>
        <div class="rule-channels">
          <span v-for="ch in rule.notify_channels" :key="ch" class="channel-tag">{{ channelLabel(ch) }}</span>
        </div>
        <div class="rule-actions">
          <button @click="toggleRule(rule)" :class="rule.enabled ? 'btn-warn' : 'btn-success'">
            {{ rule.enabled ? '禁用' : '启用' }}
          </button>
          <button @click="editRule(rule)" class="btn-outline">编辑</button>
          <button @click="deleteRule(rule.id)" class="btn-danger">删除</button>
        </div>
      </div>
    </div>

    <!-- Create/Edit Modal -->
    <div v-if="showCreate || editingRule" class="modal-overlay" @click.self="closeModal">
      <div class="modal">
        <h3>{{ editingRule ? '编辑规则' : '创建规则' }}</h3>
        <div class="form-group">
          <label>规则名称</label>
          <input v-model="form.name" placeholder="如：暴力登录检测" />
        </div>
        <div class="form-group">
          <label>描述</label>
          <input v-model="form.description" placeholder="规则描述" />
        </div>
        <div class="form-group">
          <label>事件类型</label>
          <input v-model="form.event_types_str" placeholder="多个用逗号分隔" />
        </div>
        <div class="form-row">
          <div class="form-group">
            <label>触发等级</label>
            <select v-model="form.severity_threshold">
              <option value="critical">严重 (Critical)</option>
              <option value="high">高 (High)</option>
              <option value="medium">中 (Medium)</option>
              <option value="low">低 (Low)</option>
            </select>
          </div>
          <div class="form-group">
            <label>触发次数</label>
            <input v-model.number="form.threshold_count" type="number" min="1" />
          </div>
        </div>
        <div class="form-row">
          <div class="form-group">
            <label>时间窗口(秒)</label>
            <input v-model.number="form.threshold_window_sec" type="number" min="1" />
          </div>
          <div class="form-group">
            <label>冷却时间(秒)</label>
            <input v-model.number="form.cooldown_sec" type="number" min="0" />
          </div>
        </div>
        <div class="form-group">
          <label>通知渠道</label>
          <div class="channel-checks">
            <label><input type="checkbox" v-model="form.channels.sms" /> 短信</label>
            <label><input type="checkbox" v-model="form.channels.wechat" /> 微信</label>
            <label><input type="checkbox" v-model="form.channels.wecom" /> 企业微信</label>
            <label><input type="checkbox" v-model="form.channels.dingtalk" /> 钉钉</label>
            <label><input type="checkbox" v-model="form.channels.webhook" /> Webhook</label>
          </div>
        </div>
        <div class="modal-actions">
          <button @click="closeModal" class="btn-outline">取消</button>
          <button @click="saveRule" class="btn-primary">保存</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '../services/api'

interface AlertRule {
  id: string
  name: string
  description: string
  event_types: string[]
  severity_threshold: string
  threshold_count: number
  threshold_window_sec: number
  notify_channels: string[]
  cooldown_sec: number
  enabled: boolean
}

const rules = ref<AlertRule[]>([])
const showCreate = ref(false)
const editingRule = ref<AlertRule | null>(null)

const form = ref({
  name: '',
  description: '',
  event_types_str: '',
  severity_threshold: 'high',
  threshold_count: 1,
  threshold_window_sec: 60,
  cooldown_sec: 300,
  channels: { sms: false, wechat: false, wecom: false, dingtalk: false, webhook: false }
})

onMounted(() => loadRules())

async function loadRules() {
  try {
    const result = await api.get('/alerts/rules')
    rules.value = Array.isArray(result) ? result : []
  } catch { 
    rules.value = [] 
  }
}

function channelLabel(ch: string) {
  const labels: Record<string, string> = {
    sms: '📱 短信', wechat: '💬 微信', wecom: '🏢 企业微信',
    dingtalk: '💎 钉钉', webhook: '🔗 Webhook', email: '📧 邮件'
  }
  return labels[ch] || ch
}

function editRule(rule: AlertRule) {
  editingRule.value = rule
  form.value = {
    name: rule.name,
    description: rule.description,
    event_types_str: rule.event_types.join(', '),
    severity_threshold: rule.severity_threshold,
    threshold_count: rule.threshold_count,
    threshold_window_sec: rule.threshold_window_sec,
    cooldown_sec: rule.cooldown_sec,
    channels: {
      sms: rule.notify_channels.includes('sms'),
      wechat: rule.notify_channels.includes('wechat'),
      wecom: rule.notify_channels.includes('wecom'),
      dingtalk: rule.notify_channels.includes('dingtalk'),
      webhook: rule.notify_channels.includes('webhook')
    }
  }
}

async function toggleRule(rule: AlertRule) {
  try {
    await api.put(`/alerts/rules/${rule.id}`, { ...rule, enabled: !rule.enabled })
    loadRules()
  } catch {}
}

async function deleteRule(id: string) {
  if (!confirm('确定删除此规则?')) return
  try {
    await api.delete(`/alerts/rules/${id}`)
    loadRules()
  } catch {}
}

async function saveRule() {
  const channels = Object.entries(form.value.channels).filter(([, v]) => v).map(([k]) => k)
  const data = {
    name: form.value.name,
    description: form.value.description,
    event_types: form.value.event_types_str.split(',').map(s => s.trim()).filter(Boolean),
    severity_threshold: form.value.severity_threshold,
    threshold_count: form.value.threshold_count,
    threshold_window_sec: form.value.threshold_window_sec,
    notify_channels: channels,
    cooldown_sec: form.value.cooldown_sec,
    enabled: true
  }
  try {
    if (editingRule.value) {
      await api.put(`/alerts/rules/${editingRule.value.id}`, data)
    } else {
      await api.post('/alerts/rules', data)
    }
    closeModal()
    loadRules()
  } catch {}
}

function closeModal() {
  showCreate.value = false
  editingRule.value = null
  form.value = {
    name: '', description: '', event_types_str: '', severity_threshold: 'high',
    threshold_count: 1, threshold_window_sec: 60, cooldown_sec: 300,
    channels: { sms: false, wechat: false, wecom: false, dingtalk: false, webhook: false }
  }
}
</script>

<style scoped>
.alert-rules-page { padding: 20px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
.page-header h2 { margin: 0; }
.btn-primary { background: linear-gradient(90deg, #2563eb, #7c3aed); color: white; padding: 10px 20px; border: none; border-radius: 8px; cursor: pointer; }
.rules-list { display: grid; gap: 16px; }
.rule-card { background: white; padding: 20px; border-radius: 12px; box-shadow: 0 2px 10px rgba(0,0,0,0.05); }
.rule-card.disabled { opacity: 0.6; }
.rule-header { display: flex; justify-content: space-between; margin-bottom: 8px; }
.rule-name { font-weight: 600; font-size: 1.1rem; }
.severity { padding: 4px 12px; border-radius: 20px; font-size: 0.8rem; font-weight: 500; }
.severity-critical { background: #fee2e2; color: #991b1b; }
.severity-high { background: #fef3c7; color: #92400e; }
.severity-medium { background: #dbeafe; color: #1e40af; }
.severity-low { background: #f3f4f6; color: #4b5563; }
.rule-desc { color: #64748b; margin-bottom: 12px; }
.rule-config { display: flex; gap: 16px; margin-bottom: 12px; }
.config-item { color: #64748b; font-size: 0.9rem; }
.rule-channels { display: flex; gap: 8px; flex-wrap: wrap; margin-bottom: 12px; }
.channel-tag { background: #f1f5f9; padding: 4px 10px; border-radius: 4px; font-size: 0.85rem; }
.rule-actions { display: flex; gap: 8px; }
.btn-outline { padding: 6px 16px; border: 1px solid #e2e8f0; background: white; border-radius: 6px; cursor: pointer; }
.btn-warn { padding: 6px 16px; background: #f59e0b; color: white; border: none; border-radius: 6px; cursor: pointer; }
.btn-success { padding: 6px 16px; background: #10b981; color: white; border: none; border-radius: 6px; cursor: pointer; }
.btn-danger { padding: 6px 16px; background: #ef4444; color: white; border: none; border-radius: 6px; cursor: pointer; }
.modal-overlay { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 1000; }
.modal { background: white; padding: 24px; border-radius: 12px; width: 500px; max-height: 80vh; overflow-y: auto; }
.modal h3 { margin: 0 0 20px; }
.form-group { margin-bottom: 16px; }
.form-group label { display: block; margin-bottom: 6px; color: #374151; font-weight: 500; }
.form-group input, .form-group select { width: 100%; padding: 10px; border: 1px solid #e2e8f0; border-radius: 6px; }
.form-row { display: flex; gap: 16px; }
.form-row .form-group { flex: 1; }
.channel-checks { display: flex; gap: 16px; flex-wrap: wrap; }
.channel-checks label { display: flex; align-items: center; gap: 6px; cursor: pointer; }
.modal-actions { display: flex; justify-content: flex-end; gap: 12px; margin-top: 24px; }
</style>
