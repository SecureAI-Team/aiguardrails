<template>
  <LandingLayout>
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
  </LandingLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '../services/api'
import LandingLayout from '../components/LandingLayout.vue'

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

// Default models data
const defaultModels: Model[] = [
  { id: '1', provider: 'qwen', model_id: 'qwen-max', display_name: '通义千问 Max', description: '阿里云最强大的语言模型，适合复杂推理和创作任务', capabilities: ['chat', 'reasoning'], context_window: 32000, input_price_per_m: 40, output_price_per_m: 120, enabled: true, deprecated: false },
  { id: '2', provider: 'qwen', model_id: 'qwen-plus', display_name: '通义千问 Plus', description: '性价比优秀的通用模型，适合日常对话和写作', capabilities: ['chat'], context_window: 128000, input_price_per_m: 4, output_price_per_m: 12, enabled: true, deprecated: false },
  { id: '3', provider: 'qwen', model_id: 'qwen-turbo', display_name: '通义千问 Turbo', description: '快速响应的轻量模型，适合简单任务', capabilities: ['chat'], context_window: 8000, input_price_per_m: 2, output_price_per_m: 6, enabled: true, deprecated: false },
  { id: '4', provider: 'openai', model_id: 'gpt-4o', display_name: 'GPT-4o', description: 'OpenAI 最新旗舰多模态模型', capabilities: ['chat', 'vision', 'reasoning'], context_window: 128000, input_price_per_m: 25, output_price_per_m: 100, enabled: true, deprecated: false },
  { id: '5', provider: 'openai', model_id: 'gpt-4-turbo', display_name: 'GPT-4 Turbo', description: 'GPT-4 高性价比版本', capabilities: ['chat', 'vision'], context_window: 128000, input_price_per_m: 100, output_price_per_m: 300, enabled: true, deprecated: false },
  { id: '6', provider: 'openai', model_id: 'gpt-3.5-turbo', display_name: 'GPT-3.5 Turbo', description: '快速经济的对话模型', capabilities: ['chat'], context_window: 16000, input_price_per_m: 5, output_price_per_m: 15, enabled: true, deprecated: false },
  { id: '7', provider: 'anthropic', model_id: 'claude-3-opus', display_name: 'Claude 3 Opus', description: 'Anthropic 最强大模型，适合复杂分析任务', capabilities: ['chat', 'vision', 'reasoning'], context_window: 200000, input_price_per_m: 150, output_price_per_m: 750, enabled: true, deprecated: false },
  { id: '8', provider: 'anthropic', model_id: 'claude-3-sonnet', display_name: 'Claude 3 Sonnet', description: '平衡性能和成本的多模态模型', capabilities: ['chat', 'vision'], context_window: 200000, input_price_per_m: 30, output_price_per_m: 150, enabled: true, deprecated: false },
  { id: '9', provider: 'anthropic', model_id: 'claude-3-haiku', display_name: 'Claude 3 Haiku', description: '快速高效的轻量模型', capabilities: ['chat', 'vision'], context_window: 200000, input_price_per_m: 2.5, output_price_per_m: 12.5, enabled: true, deprecated: false },
]

async function loadModels() {
  try {
    const query = filterProvider.value ? `?provider=${filterProvider.value}` : ''
    const result = await api.get(`/models${query}`)
    const apiModels = Array.isArray(result) ? result : []
    if (apiModels.length === 0) {
      models.value = filterProvider.value 
        ? defaultModels.filter(m => m.provider === filterProvider.value)
        : defaultModels
    } else {
      models.value = apiModels
    }
  } catch { 
    models.value = filterProvider.value 
      ? defaultModels.filter(m => m.provider === filterProvider.value)
      : defaultModels
  }
}

function providerName(p: string) {
  const names: Record<string, string> = { qwen: '通义千问', openai: 'OpenAI', anthropic: 'Anthropic', baidu: '百度文心' }
  return names[p] || p
}

function formatNumber(n: number) {
  if (n >= 1000000) return (n / 1000000).toFixed(0) + 'M'
  if (n >= 1000) return (n / 1000).toFixed(0) + 'K'
  return n.toString()
}
</script>

<style scoped>
.models-page { padding: 40px 48px; max-width: 1200px; margin: 0 auto; }
.page-header { margin-bottom: 24px; }
.page-header h2 { margin: 0 0 8px; color: #1e293b; }
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
