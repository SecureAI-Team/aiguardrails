<template>
  <div class="models-page">
    <div class="page-header">
      <h2>🤖 模型目录</h2>
      <p class="subtitle">支持的AI模型及定价信息</p>
    </div>

    <div class="filter-bar">
      <select v-model="filterProvider" @change="loadModels">
        <option value="">全部供应商</option>
        <option value="qwen">通义千问</option>
        <option value="openai">OpenAI</option>
        <option value="anthropic">Anthropic</option>
      </select>
    </div>

    <div class="models-grid">
      <div v-for="model in models" :key="model.id" class="model-card">
        <div class="model-header">
          <span class="provider-badge" :class="model.provider">{{ providerName(model.provider) }}</span>
          <span v-if="model.deprecated" class="deprecated-badge">已弃用</span>
        </div>
        <h3>{{ model.display_name }}</h3>
        <p class="model-id">{{ model.model_id }}</p>
        <p class="model-desc">{{ model.description }}</p>
        <div class="model-specs">
          <div class="spec-item">
            <span class="spec-label">上下文窗口</span>
            <span class="spec-value">{{ formatNumber(model.context_window) }} tokens</span>
          </div>
          <div class="spec-item">
            <span class="spec-label">能力</span>
            <span class="spec-value">
              <span v-for="cap in model.capabilities" :key="cap" class="cap-tag">{{ cap }}</span>
            </span>
          </div>
        </div>
        <div class="model-pricing">
          <div class="price-item">
            <span>输入</span>
            <span class="price">¥{{ model.input_price_per_m }}/M tokens</span>
          </div>
          <div class="price-item">
            <span>输出</span>
            <span class="price">¥{{ model.output_price_per_m }}/M tokens</span>
          </div>
        </div>
      </div>
    </div>

    <div v-if="models.length === 0" class="no-models">
      暂无模型数据
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '../services/api'

interface Model {
  id: string
  provider: string
  model_id: string
  display_name: string
  description: string
  capabilities: string[]
  context_window: number
  input_price_per_m: number
  output_price_per_m: number
  enabled: boolean
  deprecated: boolean
}

const models = ref<Model[]>([])
const filterProvider = ref('')

onMounted(() => loadModels())

async function loadModels() {
  try {
    const query = filterProvider.value ? `?provider=${filterProvider.value}` : ''
    models.value = await api.get(`/models${query}`)
  } catch { models.value = [] }
}

function providerName(p: string) {
  const names: Record<string, string> = {
    qwen: '通义千问',
    openai: 'OpenAI',
    anthropic: 'Anthropic',
    baidu: '百度文心'
  }
  return names[p] || p
}

function formatNumber(n: number) {
  if (n >= 1000000) return (n / 1000000).toFixed(0) + 'M'
  if (n >= 1000) return (n / 1000).toFixed(0) + 'K'
  return n.toString()
}
</script>

<style scoped>
.models-page { padding: 20px; max-width: 1200px; margin: 0 auto; }
.page-header { margin-bottom: 24px; }
.page-header h2 { margin: 0 0 8px; }
.subtitle { color: #64748b; margin: 0; }
.filter-bar { margin-bottom: 24px; }
.filter-bar select { padding: 10px 16px; border: 1px solid #e2e8f0; border-radius: 6px; }
.models-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(320px, 1fr)); gap: 20px; }
.model-card { background: white; padding: 24px; border-radius: 12px; box-shadow: 0 2px 10px rgba(0,0,0,0.05); }
.model-header { display: flex; justify-content: space-between; margin-bottom: 12px; }
.provider-badge { padding: 4px 12px; border-radius: 20px; font-size: 0.8rem; font-weight: 500; }
.provider-badge.qwen { background: #fef3c7; color: #92400e; }
.provider-badge.openai { background: #d1fae5; color: #065f46; }
.provider-badge.anthropic { background: #ede9fe; color: #5b21b6; }
.deprecated-badge { background: #fee2e2; color: #991b1b; padding: 4px 12px; border-radius: 20px; font-size: 0.8rem; }
.model-card h3 { margin: 0 0 4px; color: #1e293b; }
.model-id { font-family: monospace; color: #64748b; font-size: 0.85rem; margin: 0 0 8px; }
.model-desc { color: #475569; margin: 0 0 16px; font-size: 0.9rem; }
.model-specs { margin-bottom: 16px; }
.spec-item { display: flex; justify-content: space-between; padding: 8px 0; border-bottom: 1px solid #f1f5f9; }
.spec-label { color: #64748b; }
.spec-value { color: #1e293b; }
.cap-tag { background: #e0f2fe; color: #0369a1; padding: 2px 8px; border-radius: 4px; font-size: 0.75rem; margin-left: 4px; }
.model-pricing { background: #f8fafc; padding: 16px; border-radius: 8px; }
.price-item { display: flex; justify-content: space-between; margin-bottom: 8px; }
.price-item:last-child { margin-bottom: 0; }
.price { font-weight: 600; color: #059669; }
.no-models { text-align: center; padding: 60px; color: #64748b; }
</style>
