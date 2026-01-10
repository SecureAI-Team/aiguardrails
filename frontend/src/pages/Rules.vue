<template>
  <div class="h-full flex flex-col bg-gray-50">
    <!-- Header -->
    <header class="bg-white shadow">
      <div class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8 flex justify-between items-center">
        <div>
          <h1 class="text-3xl font-bold text-gray-900 tracking-tight">规则库 (Rule Library)</h1>
          <p class="mt-2 text-sm text-gray-500">
            管理 AI 安全护栏规则，使用内置模版或自定义专属规则。
          </p>
        </div>
        <button
          @click="openCreateModal"
          class="inline-flex items-center px-5 py-2.5 border border-transparent shadow-sm text-sm font-medium rounded-lg text-white bg-indigo-600 hover:bg-indigo-700 transition"
        >
          <span class="mr-2 text-lg">+</span> 新建规则
        </button>
      </div>
    </header>

    <!-- Main Content -->
    <main class="flex-1 overflow-y-auto p-4 sm:p-8">
      <div class="max-w-7xl mx-auto">
        
        <!-- Tabs -->
        <div class="mb-8 border-b border-gray-200">
          <nav class="-mb-px flex space-x-8" aria-label="Tabs">
            <button
              v-for="tab in tabs"
              :key="tab.name"
              @click="currentTab = tab.name"
              :class="[
                currentTab === tab.name
                  ? 'border-indigo-500 text-indigo-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300',
                'whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm transition-colors'
              ]"
            >
              {{ tab.label }}
              <span class="ml-2 py-0.5 px-2.5 rounded-full text-xs font-medium bg-gray-100 text-gray-600">
                {{ tab.count }}
              </span>
            </button>
          </nav>
        </div>

        <!-- Rules Grid -->
        <div v-if="filteredRules.length > 0" class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
          <div
            v-for="rule in filteredRules"
            :key="rule.id"
            class="bg-white rounded-xl shadow-sm hover:shadow-md transition duration-200 border border-gray-100 flex flex-col overflow-hidden"
          >
            <!-- Badge stripe -->
            <div :class="[
              rule.type === 'llm' ? 'bg-purple-500' : 
              rule.type === 'keyword' ? 'bg-red-500' :
              'bg-green-500', 
              'h-1.5 w-full']"></div>
            
            <div class="p-6 flex-1 flex flex-col">
              <div class="flex items-start justify-between mb-4">
                <span
                  class="inline-flex items-center px-2.5 py-0.5 rounded-md text-xs font-medium bg-opacity-10"
                  :class="[
                    rule.type === 'llm' ? 'bg-purple-50 text-purple-700' : 
                    rule.type === 'keyword' ? 'bg-red-50 text-red-700' :
                    'bg-green-50 text-green-700'
                  ]"
                >
                  <span class="mr-1.5" v-if="rule.type === 'llm'">🤖</span>
                  <span class="mr-1.5" v-else-if="rule.type === 'keyword'">🚫</span>
                  <span class="mr-1.5" v-else>📜</span>
                  
                  <span v-if="rule.type === 'llm'">LLM Security</span>
                  <span v-else-if="rule.type === 'keyword'">Keyword List</span>
                  <span v-else>OPA Policy</span>
                </span>
                <span
                  v-if="rule.is_system"
                  class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-600 border border-gray-200"
                >
                  系统内置
                </span>
              </div>

              <h3 class="text-lg font-bold text-gray-900 mb-2 line-clamp-1" :title="rule.name">
                {{ rule.name }}
              </h3>
              <p class="text-sm text-gray-500 mb-4 line-clamp-3 flex-1" :title="rule.description">
                {{ rule.description }}
              </p>

              <!-- Footer with ID and Actions -->
              <div class="pt-4 border-t border-gray-50 flex items-center justify-between mt-auto">
                <code class="text-xs text-gray-400 bg-gray-50 px-1.5 py-0.5 rounded select-all">
                  {{ rule.id.substring(0, 8) }}...
                </code>
                
                <div class="flex space-x-2">
                   <!-- <button
                    v-if="!rule.is_system"
                    @click="editRule(rule)"
                    class="text-gray-400 hover:text-indigo-600 transition"
                    title="Edit (Coming Soon)"
                  >
                    ✏️
                  </button> -->
                  <button
                    v-if="!rule.is_system"
                    @click="deleteRule(rule.id)"
                    class="text-gray-400 hover:text-red-600 transition p-1 rounded-md hover:bg-red-50"
                    title="Delete"
                  >
                    🗑️
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Empty State -->
        <div v-else class="text-center py-20 bg-white rounded-xl border-2 border-dashed border-gray-300">
          <div class="text-5xl mb-4">📭</div>
          <h3 class="mt-2 text-sm font-medium text-gray-900">暂无规则</h3>
          <p class="mt-1 text-sm text-gray-500">
            {{ currentTab === 'builtin' ? '没有找到内置规则。' : '还没有创建自定义规则。' }}
          </p>
          <div v-if="currentTab === 'custom'" class="mt-6">
            <button
              @click="openCreateModal"
              class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700"
            >
              立即创建
            </button>
          </div>
        </div>
      </div>
    </main>

    <!-- Create Modal (Better Design) -->
    <div
      v-if="showModal"
      class="fixed inset-0 z-50 overflow-y-auto"
      aria-labelledby="modal-title"
      role="dialog"
      aria-modal="true"
    >
      <div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <!-- Overlay -->
        <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" aria-hidden="true" @click="showModal = false"></div>

        <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>

        <div class="inline-block align-bottom bg-white rounded-2xl text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-4xl sm:w-full">
          <div class="bg-indigo-600 px-6 py-4 flex justify-between items-center">
            <h3 class="text-lg leading-6 font-semibold text-white" id="modal-title">创建新规则</h3>
            <button @click="showModal = false" class="text-indigo-200 hover:text-white transition">
              <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          
          <div class="px-6 py-6 sm:p-8">
            <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
              <!-- Left Column: Settings -->
              <div class="space-y-6">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">规则名称</label>
                  <input
                    v-model="newRule.name"
                    type="text"
                    class="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                    placeholder="例如: 拒绝竞品信息"
                  />
                </div>

                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">规则类型</label>
                  <div class="grid grid-cols-3 gap-3">
                    <div 
                      @click="newRule.type = 'llm'"
                      :class="[
                        newRule.type === 'llm' ? 'border-indigo-500 ring-2 ring-indigo-200 bg-indigo-50' : 'border-gray-300 hover:bg-gray-50',
                        'cursor-pointer border rounded-lg p-2 flex flex-col items-center justify-center text-center transition-all'
                      ]"
                    >
                      <span class="text-xl mb-1">🤖</span>
                      <span class="font-medium text-xs text-gray-900">LLM Security</span>
                    </div>
                    <div 
                      @click="newRule.type = 'opa'"
                      :class="[
                        newRule.type === 'opa' ? 'border-green-500 ring-2 ring-green-200 bg-green-50' : 'border-gray-300 hover:bg-gray-50',
                        'cursor-pointer border rounded-lg p-2 flex flex-col items-center justify-center text-center transition-all'
                      ]"
                    >
                      <span class="text-xl mb-1">📜</span>
                      <span class="font-medium text-xs text-gray-900">OPA Policy</span>
                    </div>
                    <div 
                      @click="newRule.type = 'keyword'"
                      :class="[
                        newRule.type === 'keyword' ? 'border-red-500 ring-2 ring-red-200 bg-red-50' : 'border-gray-300 hover:bg-gray-50',
                        'cursor-pointer border rounded-lg p-2 flex flex-col items-center justify-center text-center transition-all'
                      ]"
                    >
                      <span class="text-xl mb-1">🚫</span>
                      <span class="font-medium text-xs text-gray-900">Keyword List</span>
                    </div>
                  </div>
                </div>

                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">描述</label>
                  <textarea
                    v-model="newRule.description"
                    rows="3"
                    class="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                    placeholder="简要描述规则的用途..."
                  ></textarea>
                </div>
              </div>

              <!-- Right Column: Content -->
              <div class="space-y-6">
                <div class="flex items-center justify-between">
                  <label class="block text-sm font-medium text-gray-700">
                    <span v-if="newRule.type === 'llm'">安全指令 (System Prompt Instruction)</span>
                    <span v-else-if="newRule.type === 'keyword'">敏感词列表 (Blocked Keywords)</span>
                    <span v-else>Rego 策略代码</span>
                  </label>
                  <span class="text-xs text-gray-400">
                     <span v-if="newRule.type === 'llm'">自然语言</span>
                     <span v-else-if="newRule.type === 'keyword'">每行一个词</span>
                     <span v-else>Rego Language</span>
                  </span>
                </div>
                
                <div class="relative">
                  <textarea
                    v-model="newRule.content"
                    rows="12"
                    class="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 font-mono text-sm leading-relaxed"
                    :placeholder="placeholderText"
                  ></textarea>
                  <div class="absolute bottom-2 right-2 text-xs text-gray-400 bg-white px-1">
                    {{ newRule.content.length }} chars
                  </div>
                </div>
                
                <div class="bg-gray-50 p-3 rounded-md border border-gray-200">
                  <h4 class="text-xs font-semibold text-gray-700 mb-1">🔔 提示</h4>
                  <p v-if="newRule.type === 'llm'" class="text-xs text-gray-500">
                    此指令将作为 System Prompt 的一部分发送给 Qwen 安全模型。请清晰描述需要拦截的场景。
                  </p>
                  <p v-else-if="newRule.type === 'keyword'" class="text-xs text-gray-500">
                    输入需要拦截的敏感词或短语，每行一个。
                  </p>
                  <p v-else class="text-xs text-gray-500">
                    OPA 规则需要编写合法的 Rego 代码。通常用于结构化数据的精确匹配。
                  </p>
                </div>
              </div>
            </div>
          </div>

          <div class="bg-gray-50 px-6 py-4 flex justify-end space-x-3 rounded-b-2xl">
            <button
              @click="showModal = false"
              class="px-5 py-2.5 bg-white border border-gray-300 rounded-lg text-gray-700 font-medium hover:bg-gray-50 transition shadow-sm"
            >
              取消
            </button>
            <button
              @click="createRule"
              :disabled="loading || !newRule.name || !newRule.content"
              class="px-5 py-2.5 bg-indigo-600 rounded-lg text-white font-medium hover:bg-indigo-700 transition shadow-md disabled:opacity-50 disabled:cursor-not-allowed flex items-center"
            >
              <svg v-if="loading" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              {{ loading ? '创建中...' : '确认创建' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { api } from '../services/api';

interface Rule {
  id: string;
  name: string;
  description: string;
  type: 'opa' | 'llm' | 'keyword';
  content: string;
  is_system: boolean;
}

const rules = ref<Rule[]>([]);
const showModal = ref(false);
const loading = ref(false);
const currentTab = ref('builtin');

const newRule = ref({
  name: '',
  description: '',
  type: 'llm',
  content: ''
});

const builtInRules = computed(() => rules.value.filter(r => r.is_system));
const customRules = computed(() => rules.value.filter(r => !r.is_system));

const filteredRules = computed(() => {
  return currentTab.value === 'builtin' ? builtInRules.value : customRules.value;
});

const tabs = computed(() => [
  { name: 'builtin', label: 'Preset Library (内置)', count: builtInRules.value.length },
  { name: 'custom', label: 'Custom Rules (自定义)', count: customRules.value.length }
]);

const placeholderText = computed(() => {
  if (newRule.value.type === 'llm') {
    return 'You are a helpful assistant. Please ensure the response does not contain any personal identifiable information...';
  } else if (newRule.value.type === 'keyword') {
    return '敏感词1\n敏感词2\nblocked_word\n...';
  }
  return 'package guardrails\n\ndefault allow = true\n\ndeny[msg] {\n  input.prompt == "fail"\n  msg := "prompt blocked"\n}';
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
    rules.value = [];
  }
}

function openCreateModal() {
  newRule.value = { name: '', description: '', type: 'llm', content: '' };
  showModal.value = true;
  // Automatically switch to custom tab if creating
  currentTab.value = 'custom';
}

async function createRule() {
  if (!newRule.value.name || !newRule.value.content) return;
  loading.value = true;
  try {
    await api.post('/rules', newRule.value);
    showModal.value = false;
    await fetchRules();
    currentTab.value = 'custom';
  } catch (e) {
    alert('Failed to create rule');
  } finally {
    loading.value = false;
  }
}

async function deleteRule(id: string) {
  if (!confirm('确定要删除此规则吗？')) return;
  try {
    await api.delete(`/rules/${id}`);
    await fetchRules();
  } catch (e) {
    alert('Failed to delete rule');
  }
}
</script>
