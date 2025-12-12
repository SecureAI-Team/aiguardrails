<template>
  <div>
    <h2>Capabilities</h2>
    <form @submit.prevent="onCreate">
      <input v-model="name" placeholder="Name" />
      <input v-model="description" placeholder="Description" />
      <input v-model="tagsInput" placeholder="Tags (comma)" />
      <button :disabled="loading">Create</button>
    </form>
    <form @submit.prevent="load">
      <input v-model="filterTag" placeholder="Filter tag" />
      <button :disabled="loading">Filter</button>
    </form>
    <div class="card" v-for="c in caps" :key="c.id || c.name">
      <div><strong>{{ c.name }}</strong></div>
      <div class="muted">{{ c.description }}</div>
      <div class="muted">Tags: {{ (c.tags || []).join(', ') }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '../services/api'

const name = ref('')
const description = ref('')
const tagsInput = ref('')
const filterTag = ref('')
const caps = ref<any[]>([])
const loading = ref(false)

async function load() {
  loading.value = true
  caps.value = await api.listCapabilities(filterTag.value || undefined)
  loading.value = false
}

async function onCreate() {
  loading.value = true
  await api.createCapability({
    name: name.value,
    description: description.value,
    tags: tagsInput.value.split(',').map((t) => t.trim()).filter(Boolean)
  })
  name.value = ''
  description.value = ''
  tagsInput.value = ''
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

