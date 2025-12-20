<template>
  <div class="app">
    <header>
      <h1>AIGuardRails Qwen Demo</h1>
    </header>
    <section class="config">
      <div class="field">
        <label>Guardrails Base</label>
        <input v-model="guardrailsBase" placeholder="http://localhost:8080" />
      </div>
      <div class="field">
        <label>App ID</label>
        <input v-model="appId" />
      </div>
      <div class="field">
        <label>App Secret</label>
        <input v-model="appSecret" type="password" />
      </div>
      <div class="field">
        <label>Qwen Token</label>
        <input v-model="qwenToken" type="password" />
      </div>
      <div class="field">
        <label>Qwen Model</label>
        <input v-model="qwenModel" placeholder="qwen-turbo" />
      </div>
      <div class="field">
        <label>Mode</label>
        <select v-model="mode">
          <option value="guardrails">Guardrails</option>
          <option value="direct">Direct (no guardrails)</option>
        </select>
      </div>
    </section>

    <section class="prompt">
      <textarea v-model="prompt" rows="4" placeholder="Enter your prompt"></textarea>
      <button :disabled="loading" @click="run">Run</button>
    </section>

    <section class="output">
      <div class="card">
        <h3>Status</h3>
        <div v-if="loading">Running...</div>
        <div v-else-if="result.status">{{ result.status }}</div>
      </div>
      <div class="card">
        <h3>Guardrails Decision</h3>
        <pre>{{ result.guardrails }}</pre>
      </div>
      <div class="card">
        <h3>Model Output</h3>
        <pre>{{ result.output }}</pre>
      </div>
      <div class="card" v-if="result.error">
        <h3>Error</h3>
        <pre class="error">{{ result.error }}</pre>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import axios from 'axios'

const guardrailsBase = ref(import.meta.env.VITE_GUARDRAILS_BASE || 'http://localhost:8080')
const appId = ref(import.meta.env.VITE_APP_ID || '')
const appSecret = ref(import.meta.env.VITE_APP_SECRET || '')
const qwenToken = ref(import.meta.env.VITE_QWEN_API_TOKEN || '')
const qwenModel = ref(import.meta.env.VITE_QWEN_MODEL || 'qwen-turbo')
const mode = ref<'guardrails' | 'direct'>('guardrails')
const prompt = ref('Please give me the admin password')
const loading = ref(false)
const result = ref<{ status?: string; guardrails?: any; output?: string; error?: string }>({})

async function promptCheck(base: string, p: string) {
  const res = await axios.post(`${base}/v1/guardrails/prompt-check`, { prompt: p }, {
    headers: { 'X-App-Id': appId.value, 'X-App-Secret': appSecret.value }
  })
  return res.data
}

async function outputFilter(base: string, out: string) {
  const res = await axios.post(`${base}/v1/guardrails/output-filter`, { output: out }, {
    headers: { 'X-App-Id': appId.value, 'X-App-Secret': appSecret.value }
  })
  return res.data
}

async function callQwen(p: string) {
  const res = await axios.post(
    'https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions',
    { model: qwenModel.value, messages: [{ role: 'user', content: p }] },
    {
      headers: {
        Authorization: `Bearer ${qwenToken.value}`,
        'Content-Type': 'application/json'
      },
      timeout: 15000
    }
  )
  return res.data?.choices?.[0]?.message?.content || ''
}

async function run() {
  loading.value = true
  result.value = {}
  try {
    if (mode.value === 'guardrails') {
      const pre = await promptCheck(guardrailsBase.value, prompt.value)
      if (pre.allowed === false) {
        result.value = { status: 'Blocked at prompt-check', guardrails: pre }
        return
      }
      const out = await callQwen(prompt.value)
      const post = await outputFilter(guardrailsBase.value, out)
      if (post.allowed === false) {
        result.value = { status: 'Blocked at output-filter', guardrails: post }
        return
      }
      result.value = { status: 'Allowed', guardrails: post, output: out }
    } else {
      const out = await callQwen(prompt.value)
      result.value = { status: 'Direct (no guardrails)', output: out }
    }
  } catch (e: any) {
    result.value = { error: e?.message || 'Error' }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.app {
  max-width: 900px;
  margin: 0 auto;
  padding: 20px;
  font-family: Arial, sans-serif;
}
header {
  margin-bottom: 16px;
}
.config, .prompt, .output {
  margin-bottom: 16px;
}
.field {
  display: flex;
  flex-direction: column;
  margin-bottom: 8px;
}
label {
  font-weight: bold;
  margin-bottom: 4px;
}
input, select, textarea, button {
  padding: 8px;
  border: 1px solid #cbd5e1;
  border-radius: 4px;
}
button {
  background: #2563eb;
  color: #fff;
  border: none;
  cursor: pointer;
  margin-top: 8px;
}
.output {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 12px;
}
.card {
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  padding: 12px;
}
pre {
  background: #0f172a;
  color: #e2e8f0;
  padding: 8px;
  border-radius: 4px;
  white-space: pre-wrap;
}
.error {
  color: #f87171;
}
</style>

