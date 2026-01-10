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

    <!-- Create Modal -->
    <Teleport to="body">
      <div
        v-if="showModal"
        class="fixed inset-0 z-[100] overflow-y-auto"
        role="dialog"
        aria-modal="true"
      >
        <!-- Overlay -->
        <div class="fixed inset-0 bg-gray-900 bg-opacity-50 transition-opacity backdrop-blur-sm" @click="showModal = false"></div>

        <div class="flex min-h-full items-center justify-center p-4 text-center sm:p-0">
          <div class="relative transform overflow-hidden rounded-2xl bg-white text-left shadow-2xl transition-all sm:my-8 sm:w-full sm:max-w-4xl border border-gray-100">
            
            <!-- Modal Header -->
            <div class="bg-white border-b border-gray-100 px-6 py-4 flex justify-between items-center sm:px-8">
              <div>
                <h3 class="text-xl font-bold text-gray-900">创建新规则</h3>
                <p class="mt-1 text-sm text-gray-500">配置安全护栏规则以拦截风险内容。</p>
              </div>
              <button @click="showModal = false" class="rounded-full p-1 bg-gray-50 hover:bg-gray-100 text-gray-400 hover:text-gray-500 transition">
                <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            
            <div class="px-6 py-6 sm:px-8 bg-gray-50/50"> <!-- Added slight bg -->
              <div class="grid grid-cols-1 lg:grid-cols-12 gap-8">
                <!-- Left Column: Settings (5 cols) -->
                <div class="lg:col-span-5 space-y-6">
                  <div>
                    <label class="block text-sm font-semibold text-gray-700 mb-2">规则名称</label>
                    <input
                      v-model="newRule.name"
                      type="text"
                      class="block w-full rounded-lg border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm px-4 py-2.5 transition ease-in-out duration-150"
                      placeholder="例如: 拒绝竞品信息"
                    />
                  </div>

                  <div>
                    <label class="block text-sm font-semibold text-gray-700 mb-2">规则类型</label>
                    <div class="space-y-3">
                      <!-- LLM Type -->
                      <div 
                        @click="newRule.type = 'llm'"
                        :class="[
                          newRule.type === 'llm' ? 'border-indigo-500 ring-1 ring-indigo-500 bg-indigo-50/50' : 'border-gray-200 hover:border-gray-300 bg-white shadow-sm',
                          'cursor-pointer relative flex items-center px-4 py-3 border rounded-xl transition-all duration-200'
                        ]"
                      >
                         <div class="flex-shrink-0 h-10 w-10 flex items-center justify-center rounded-lg" :class="newRule.type === 'llm' ? 'bg-indigo-100 text-indigo-600' : 'bg-gray-100 text-gray-500'">
                           <span class="text-xl">🤖</span>
                         </div>
                         <div class="ml-3">
                           <p class="text-sm font-medium text-gray-900">LLM Security</p>
                           <p class="text-xs text-gray-500">语义模型检测</p>
                         </div>
                         <div v-if="newRule.type === 'llm'" class="absolute right-4 top-1/2 -translate-y-1/2 text-indigo-500">
                           <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" /></svg>
                         </div>
                      </div>

                      <!-- OPA Type -->
                      <div 
                        @click="newRule.type = 'opa'"
                        :class="[
                          newRule.type === 'opa' ? 'border-green-500 ring-1 ring-green-500 bg-green-50/50' : 'border-gray-200 hover:border-gray-300 bg-white shadow-sm',
                          'cursor-pointer relative flex items-center px-4 py-3 border rounded-xl transition-all duration-200'
                        ]"
                      >
                         <div class="flex-shrink-0 h-10 w-10 flex items-center justify-center rounded-lg" :class="newRule.type === 'opa' ? 'bg-green-100 text-green-600' : 'bg-gray-100 text-gray-500'">
                           <span class="text-xl">📜</span>
                         </div>
                         <div class="ml-3">
                           <p class="text-sm font-medium text-gray-900">OPA Policy</p>
                           <p class="text-xs text-gray-500">Rego 代码规则</p>
                         </div>
                         <div v-if="newRule.type === 'opa'" class="absolute right-4 top-1/2 -translate-y-1/2 text-green-500">
                           <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" /></svg>
                         </div>
                      </div>

                      <!-- Keyword Type -->
                      <div 
                        @click="newRule.type = 'keyword'"
                        :class="[
                          newRule.type === 'keyword' ? 'border-red-500 ring-1 ring-red-500 bg-red-50/50' : 'border-gray-200 hover:border-gray-300 bg-white shadow-sm',
                          'cursor-pointer relative flex items-center px-4 py-3 border rounded-xl transition-all duration-200'
                        ]"
                      >
                         <div class="flex-shrink-0 h-10 w-10 flex items-center justify-center rounded-lg" :class="newRule.type === 'keyword' ? 'bg-red-100 text-red-600' : 'bg-gray-100 text-gray-500'">
                           <span class="text-xl">🚫</span>
                         </div>
                         <div class="ml-3">
                           <p class="text-sm font-medium text-gray-900">Keyword List</p>
                           <p class="text-xs text-gray-500">敏感词列表</p>
                         </div>
                         <div v-if="newRule.type === 'keyword'" class="absolute right-4 top-1/2 -translate-y-1/2 text-red-500">
                           <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" /></svg>
                         </div>
                      </div>
                    </div>
                  </div>

                  <div>
                    <label class="block text-sm font-semibold text-gray-700 mb-2">描述</label>
                    <textarea
                      v-model="newRule.description"
                      rows="3"
                      class="block w-full rounded-lg border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm px-4 py-2"
                      placeholder="简要描述规则的用途..."
                    ></textarea>
                  </div>
                </div>

                <!-- Right Column: Content (7 cols) -->
                <div class="lg:col-span-7 flex flex-col h-full">
                  <div class="flex items-center justify-between mb-2">
                    <label class="block text-sm font-semibold text-gray-700">
                      <span v-if="newRule.type === 'llm'">安全指令 (System Prompt Instruction)</span>
                      <span v-else-if="newRule.type === 'keyword'">敏感词列表 (Blocked Keywords)</span>
                      <span v-else>Rego 策略代码</span>
                    </label>
                    <span class="inline-flex items-center rounded-full bg-gray-100 px-2.5 py-0.5 text-xs font-medium text-gray-800">
                       {{ newRule.type === 'llm' ? 'Natural Language' : newRule.type === 'keyword' ? 'Line separated' : 'Rego' }}
                    </span>
                  </div>
                  
                  <div class="relative flex-1">
                    <textarea
                      v-model="newRule.content"
                      class="block w-full h-full min-h-[300px] rounded-lg border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 font-mono text-sm leading-relaxed p-4 bg-white"
                      :placeholder="placeholderText"
                    ></textarea>
                  </div>
                  
                  <!-- Tip Box -->
                  <div class="mt-4 rounded-lg bg-blue-50 p-4 border border-blue-100 flex items-start">
                    <div class="flex-shrink-0">
                      <svg class="h-5 w-5 text-blue-400" viewBox="0 0 20 20" fill="currentColor">
                        <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
                      </svg>
                    </div>
                    <div class="ml-3 flex-1 md:flex md:justify-between">
                      <p class="text-sm text-blue-700">
                        <span v-if="newRule.type === 'llm'">此指令将作为 System Prompt 发送给 Qwen 模型。请清晰描述拦截逻辑。</span>
                        <span v-else-if="newRule.type === 'keyword'">输入需要拦截的词汇，通过换行分隔。</span>
                        <span v-else>编写 Rego 代码以定义 OPA 策略。</span>
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div class="bg-white border-t border-gray-100 px-6 py-4 flex justify-end space-x-3 sm:px-8">
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
    </Teleport>
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
