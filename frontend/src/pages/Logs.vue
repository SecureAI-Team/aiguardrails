<template>
  <div>
    <h2>Audit Logs</h2>
    <form @submit.prevent="load">
      <input v-model="eventLike" placeholder="Event contains" />
      <input v-model="tenantId" placeholder="Tenant ID" />
      <input v-model.number="limit" type="number" min="1" max="500" />
      <button :disabled="loading">Refresh</button>
    </form>
    <div class="card" v-for="(e, idx) in events" :key="idx">
      <div><strong>{{ e.event || e.type }}</strong></div>
      <div class="muted">{{ e.created_at }}</div>
      <pre class="code">{{ JSON.stringify(e, null, 2) }}</pre>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '../services/api'

const events = ref<any[]>([])
const loading = ref(false)
const eventLike = ref('')
const tenantId = ref('')
const limit = ref(100)

async function load() {
  loading.value = true
  events.value = await api.listAudit(limit.value, eventLike.value || undefined, tenantId.value || undefined)
  loading.value = false
}

onMounted(load)
</script>

<style scoped>
.muted {
  color: #64748b;
  font-size: 12px;
}
.code {
  background: #0f172a;
  color: #e2e8f0;
  padding: 8px;
  border-radius: 4px;
  overflow-x: auto;
}
</style>

