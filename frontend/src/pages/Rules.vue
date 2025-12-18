<template>
  <div>
    <h2>Rule Library</h2>
    <form @submit.prevent="load">
      <input v-model="jurisdiction" placeholder="Jurisdiction (e.g., EU)" />
      <input v-model="regulation" placeholder="Regulation (e.g., GDPR)" />
      <input v-model="vendor" placeholder="Vendor (e.g., Siemens)" />
      <input v-model="product" placeholder="Product" />
      <button :disabled="loading">Filter</button>
    </form>
    <div class="card" v-for="r in rules" :key="r.id">
      <div><strong>{{ r.name }}</strong> ({{ r.regulation }} / {{ r.jurisdiction }})</div>
      <div class="muted">Vendor: {{ r.vendor || 'N/A' }} | Product: {{ r.product || 'N/A' }}</div>
      <div class="muted">Severity: {{ r.severity }} | Category: {{ r.category }}</div>
      <div class="muted">Tags: {{ (r.tags || []).join(', ') }}</div>
      <div>{{ r.description }}</div>
      <div class="muted">Remediation: {{ r.remediation }}</div>
      <div class="muted">Refs: {{ (r.references || []).join(', ') }}</div>
      <div class="attach">
        <input v-model="policyId" placeholder="Policy ID to attach" />
        <button @click="attach(r.id)" :disabled="loading || !policyId">Attach</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '../services/api'

const jurisdiction = ref('')
const regulation = ref('')
const vendor = ref('')
const product = ref('')
const policyId = ref('')
const rules = ref<any[]>([])
const loading = ref(false)

async function load() {
  loading.value = true
  rules.value = await api.listRules({
    jurisdiction: jurisdiction.value,
    regulation: regulation.value,
    vendor: vendor.value,
    product: product.value
  })
  loading.value = false
}

async function attach(ruleId: string) {
  if (!policyId.value) return
  loading.value = true
  await api.attachRule(policyId.value, ruleId)
  loading.value = false
}

onMounted(load)
</script>

<style scoped>
.muted {
  color: #64748b;
  font-size: 12px;
}
.card {
  margin-bottom: 12px;
}
.attach {
  display: flex;
  gap: 8px;
  margin-top: 8px;
}
</style>

