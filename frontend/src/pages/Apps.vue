<template>
  <div>
    <h2>Apps</h2>
    <form @submit.prevent="onCreate">
      <input v-model="tenantId" placeholder="Tenant ID" />
      <input v-model="name" placeholder="App name" />
      <input v-model.number="quota" type="number" placeholder="Quota/hr" />
      <button :disabled="loading">Create</button>
    </form>
    <button @click="load" :disabled="!tenantId">Load apps</button>
    <div class="card" v-for="a in apps" :key="a.id">
      <div><strong>{{ a.name }}</strong></div>
      <div class="muted">App ID: {{ a.id }}</div>
      <div class="muted">API Key: {{ a.apiKey }}</div>
      <div class="muted">Quota/hr: {{ a.quotaPerHr }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { api } from '../services/api'

const tenantId = ref('')
const name = ref('')
const quota = ref(1000)
const apps = ref<any[]>([])
const loading = ref(false)

async function load() {
  if (!tenantId.value) return
  apps.value = await api.listApps(tenantId.value)
}

async function onCreate() {
  if (!tenantId.value || !name.value) return
  loading.value = true
  await api.createApp(tenantId.value, name.value, quota.value || 0)
  name.value = ''
  await load()
  loading.value = false
}
</script>

<style scoped>
.muted {
  color: #64748b;
  font-size: 12px;
}
</style>

