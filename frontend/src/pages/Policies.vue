<template>
  <div>
    <h2>Policies</h2>
    <form @submit.prevent="onCreate">
      <input v-model="tenantId" placeholder="Tenant ID" />
      <input v-model="name" placeholder="Policy name" />
      <input v-model="promptRules" placeholder="Prompt rules (comma)" />
    <input v-model="toolAllow" placeholder="Tool allowlist (comma)" />
      <input v-model="ragNamespaces" placeholder="RAG namespaces (comma)" />
    <input v-model="sensitiveTerms" placeholder="Sensitive terms (comma)" />
    <div class="caps">
      <div class="caps-header">
        <span>Capabilities</span>
        <input v-model="capTag" placeholder="Filter tag" />
        <button type="button" @click="loadCaps">Filter</button>
      </div>
      <div class="caps-list">
        <label v-for="c in caps" :key="c.id || c.name">
          <input type="checkbox" :value="c.name" v-model="selectedCaps" />
          {{ c.name }} ({{ (c.tags || []).join(', ') }})
        </label>
      </div>
    </div>
    <div class="caps-selected" v-if="selectedCaps.length">
      <strong>Selected capabilities:</strong>
      <span class="chip" v-for="c in selectedCaps" :key="c">{{ c }}</span>
    </div>
      <button :disabled="loading">Create</button>
    </form>
    <button @click="load" :disabled="!tenantId">Load policies</button>
    <div class="card" v-for="p in policies" :key="p.id">
      <div><strong>{{ p.name }}</strong></div>
      <div class="muted">Prompt rules: {{ p.promptRules?.join(', ') }}</div>
      <div class="muted">Tools: {{ p.toolAllowList?.join(', ') }}</div>
      <div class="muted">Namespaces: {{ p.ragNamespaces?.join(', ') }}</div>
      <div class="muted">Sensitive: {{ p.sensitiveTerms?.join(', ') }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { api } from '../services/api'

const tenantId = ref('')
const name = ref('')
const promptRules = ref('')
const toolAllow = ref('')
const ragNamespaces = ref('')
const sensitiveTerms = ref('')
const policies = ref<any[]>([])
const loading = ref(false)
const caps = ref<any[]>([])
const capTag = ref('')
const selectedCaps = ref<string[]>([])

async function load() {
  if (!tenantId.value) return
  policies.value = await api.listPolicies(tenantId.value)
}

async function loadCaps() {
  caps.value = await api.listCapabilities(capTag.value || undefined)
}

async function onCreate() {
  if (!tenantId.value || !name.value) return
  loading.value = true
  await api.createPolicy(tenantId.value, {
    name: name.value,
    prompt_rules: promptRules.value.split(',').map((s) => s.trim()).filter(Boolean),
    tool_allowlist: [
      ...toolAllow.value.split(',').map((s) => s.trim()).filter(Boolean),
      ...selectedCaps.value
    ],
    rag_namespaces: ragNamespaces.value.split(',').map((s) => s.trim()).filter(Boolean),
    sensitive_terms: sensitiveTerms.value.split(',').map((s) => s.trim()).filter(Boolean)
  })
  name.value = ''
  await load()
  loading.value = false
}

loadCaps()
</script>

<style scoped>
.muted {
  color: #64748b;
  font-size: 12px;
}
.caps {
  padding: 8px 0;
}
.caps-header {
  display: flex;
  gap: 8px;
  align-items: center;
  margin-bottom: 6px;
}
.caps-list {
  display: grid;
  gap: 4px;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
}
.caps-selected {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin: 6px 0 10px;
  align-items: center;
}
.chip {
  background: #e2e8f0;
  padding: 4px 8px;
  border-radius: 10px;
  font-size: 12px;
}
</style>

