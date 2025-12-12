<template>
  <div>
    <h2>Policy History</h2>
    <form @submit.prevent="load">
      <input v-model="tenantId" placeholder="Tenant ID" />
      <button :disabled="loading">Load History</button>
    </form>
    <div class="card" v-for="p in history" :key="p.id + p.updatedAt">
      <div><strong>{{ p.name }}</strong> — {{ p.updatedAt }}</div>
      <div class="muted">Prompt rules: {{ (p.promptRules || []).join(', ') }}</div>
      <div class="muted">Tools: {{ (p.toolAllowList || []).join(', ') }}</div>
      <div class="muted">Namespaces: {{ (p.ragNamespaces || []).join(', ') }}</div>
      <div class="muted">Sensitive: {{ (p.sensitiveTerms || []).join(', ') }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { api } from '../services/api'

const tenantId = ref('')
const history = ref<any[]>([])
const loading = ref(false)

async function load() {
  if (!tenantId.value) return
  loading.value = true
  history.value = await api.listPolicyHistory(tenantId.value)
  loading.value = false
}
</script>

<style scoped>
.muted {
  color: #64748b;
  font-size: 12px;
}
</style>

