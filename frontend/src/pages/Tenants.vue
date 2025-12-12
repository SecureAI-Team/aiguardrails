<template>
  <div>
    <h2>Tenants</h2>
    <form @submit.prevent="onCreate">
      <input v-model="name" placeholder="Tenant name" />
      <button :disabled="loading">Create</button>
    </form>
    <div class="card" v-for="t in tenants" :key="t.id">
      <div><strong>{{ t.name }}</strong></div>
      <div class="muted">ID: {{ t.id }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { api } from '../services/api'

const tenants = ref<any[]>([])
const name = ref('')
const loading = ref(false)

async function load() {
  tenants.value = await api.listTenants()
}

async function onCreate() {
  if (!name.value) return
  loading.value = true
  await api.createTenant(name.value)
  name.value = ''
  await load()
  loading.value = false
}

onMounted(load)
</script>

<style scoped>
.muted {
  color: #64748b;
  font-size: 12px;
}
</style>

