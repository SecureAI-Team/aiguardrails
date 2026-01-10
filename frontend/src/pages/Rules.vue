<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Rule Library</h1>
        <p class="mt-1 text-sm text-gray-500">
          Manage reuseable guardrail rules (OPA Policy or LLM Check).
        </p>
      </div>
      <button
        @click="openCreateModal"
        class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700"
      >
        <span class="mr-2">+</span> Create Rule
      </button>
    </div>

    <!-- Rule List -->
    <div class="bg-white shadow overflow-hidden sm:rounded-lg">
      <ul role="list" class="divide-y divide-gray-200">
        <li v-for="rule in rules" :key="rule.id" class="px-6 py-4">
          <div class="flex items-center justify-between">
            <div class="flex-1 min-w-0">
              <div class="flex items-center space-x-3">
                <span class="text-sm font-medium text-indigo-600 truncate">{{ rule.name }}</span>
                <span
                  v-if="rule.type === 'llm'"
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-purple-100 text-purple-800"
                >
                  <span class="mr-1">🤖</span> LLM Security
                </span>
                <span
                  v-else
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800"
                >
                  <span class="mr-1">📜</span> OPA Policy
                </span>
                <span
                  v-if="rule.is_system"
                  class="inline-flex items-center px-2.5 py-0.5 rounded-md text-xs font-medium bg-gray-100 text-gray-800"
                >
                  System
                </span>
              </div>
              <p class="mt-1 text-sm text-gray-500 truncate">{{ rule.description }}</p>
            </div>
            <div class="flex items-center space-x-4">
              <span class="text-xs text-gray-400">ID: {{ rule.id }}</span>
              <button
                v-if="!rule.is_system"
                @click="deleteRule(rule.id)"
                class="text-red-600 hover:text-red-900 text-sm"
              >
                Delete
              </button>
            </div>
          </div>
        </li>
      </ul>
      <div v-if="rules.length === 0" class="px-6 py-12 text-center text-gray-500">
        No rules found. Create one to get started.
      </div>
    </div>

    <!-- Create Modal -->
    <div
      v-if="showModal"
      class="fixed inset-0 bg-gray-500 bg-opacity-75 flex items-center justify-center p-4 z-50"
    >
      <div class="bg-white rounded-lg max-w-2xl w-full p-6 space-y-4">
        <h3 class="text-lg font-medium text-gray-900">Create New Rule</h3>
        
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700">Name</label>
            <input
              v-model="newRule.name"
              type="text"
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
              placeholder="e.g. Block Personal Info"
            />
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700">Type</label>
            <select
              v-model="newRule.type"
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
            >
              <option value="llm">LLM Security (AI Check)</option>
              <option value="opa">OPA Policy (Rego)</option>
            </select>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700">Description</label>
            <input
              v-model="newRule.description"
              type="text"
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700">
              {{ newRule.type === 'llm' ? 'Safety Instruction (Prompt)' : 'Rego Policy Code' }}
            </label>
            <textarea
              v-model="newRule.content"
              rows="6"
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 font-mono text-sm"
              :placeholder="newRule.type === 'llm' ? 'e.g. Detect and block any content related to ...' : 'package guardrails\n\ndeny[msg] { ... }'"
            ></textarea>
            <p class="mt-1 text-xs text-gray-500">
              {{ newRule.type === 'llm' ? 'This instruction will be sent to the Safety LLM.' : 'Standard OPA Rego code.' }}
            </p>
          </div>
        </div>

        <div class="flex justify-end space-x-3 mt-6">
          <button
            @click="showModal = false"
            class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
          >
            Cancel
          </button>
          <button
            @click="createRule"
            :disabled="loading"
            class="px-4 py-2 text-sm font-medium text-white bg-indigo-600 rounded-md hover:bg-indigo-700 disabled:opacity-50"
          >
            {{ loading ? 'Creating...' : 'Create Rule' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { api } from '../services/api';

interface Rule {
  id: string;
  name: string;
  description: string;
  type: 'opa' | 'llm';
  content: string;
  is_system: boolean;
}

const rules = ref<Rule[]>([]);
const showModal = ref(false);
const loading = ref(false);

const newRule = ref({
  name: '',
  description: '',
  type: 'llm',
  content: ''
});

onMounted(() => {
  fetchRules();
});

async function fetchRules() {
  try {
    const res = await api.get('/rules');
    rules.value = Array.isArray(res.data) ? res.data : [];
  } catch (e) {
    console.error(e);
  }
}

function openCreateModal() {
  newRule.value = { name: '', description: '', type: 'llm', content: '' };
  showModal.value = true;
}

async function createRule() {
  if (!newRule.value.name || !newRule.value.content) return;
  loading.value = true;
  try {
    await api.post('/rules', newRule.value);
    showModal.value = false;
    await fetchRules();
  } catch (e) {
    alert('Failed to create rule');
  } finally {
    loading.value = false;
  }
}

async function deleteRule(id: string) {
  if (!confirm('Are you sure?')) return;
  try {
    await api.delete(`/rules/${id}`);
    await fetchRules();
  } catch (e) {
    alert('Failed to delete rule');
  }
}
</script>
