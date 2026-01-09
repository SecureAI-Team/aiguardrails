<template>
  <LandingLayout>
    <div class="page-container">
      <div class="docs-layout">
        <!-- Sidebar Navigation -->
        <aside class="sidebar">
          <div class="sidebar-header">
            <h3>文档目录</h3>
          </div>
          <nav>
            <div 
              v-for="(item, key) in sections" 
              :key="key"
              class="nav-item"
              :class="{ active: currentSection === key }"
              @click="currentSection = key"
            >
              {{ item.title }}
            </div>
          </nav>
          <div class="sidebar-footer">
            <router-link to="/api-reference" class="nav-item is-link">
              API 参考手册 ↗
            </router-link>
            <a href="https://github.com/helper-ai-team/aiguardrails" target="_blank" class="nav-item is-link">
              GitHub 仓库 ↗
            </a>
          </div>
        </aside>

        <!-- Main Content -->
        <main class="content-area">
          <div class="content-header">
            <h2 class="section-title">{{ sections[currentSection].title }}</h2>
          </div>
          <div class="content-body">
            <component :is="sections[currentSection].component" />
          </div>
        </main>
      </div>
    </div>
  </LandingLayout>
</template>

<script setup lang="ts">
import { ref, shallowRef } from 'vue'
import LandingLayout from '../components/LandingLayout.vue'
import Intro from '../components/docs/Intro.vue'
import GettingStarted from '../components/docs/GettingStarted.vue'
import Concepts from '../components/docs/Concepts.vue'
import Guides from '../components/docs/Guides.vue'

const currentSection = ref('intro')

const sections = {
  intro: { title: '产品简介', component: Intro },
  gettingStarted: { title: '快速开始', component: GettingStarted },
  concepts: { title: '核心概念', component: Concepts },
  guides: { title: '开发指南', component: Guides }
}
</script>

<style scoped>
.page-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0;
  height: calc(100vh - 64px); /* Subtract header height roughly */
  overflow: hidden;
}

.docs-layout {
  display: flex;
  height: 100%;
  background: white;
}

/* Sidebar */
.sidebar {
  width: 280px;
  background: #f8fafc;
  border-right: 1px solid #e2e8f0;
  display: flex;
  flex-direction: column;
  padding: 24px 0;
}

.sidebar-header {
  padding: 0 24px 20px;
  border-bottom: 1px solid #e2e8f0;
  margin-bottom: 12px;
}

.sidebar-header h3 {
  margin: 0;
  font-size: 1.1rem;
  color: #0f172a;
  font-weight: 600;
}

.nav-item {
  padding: 10px 24px;
  cursor: pointer;
  color: #64748b;
  font-size: 0.95rem;
  transition: all 0.2s;
  display: block;
  text-decoration: none;
}

.nav-item:hover {
  color: #2563eb;
  background: #eff6ff;
}

.nav-item.active {
  color: #2563eb;
  background: #eff6ff;
  border-right: 3px solid #2563eb;
  font-weight: 500;
}

.nav-item.is-link {
  color: #475569;
  font-size: 0.9rem;
}

.sidebar-footer {
  margin-top: auto;
  padding-top: 20px;
  border-top: 1px solid #e2e8f0;
}

/* Content */
.content-area {
  flex: 1;
  overflow-y: auto;
  padding: 40px 60px;
}

.content-header {
  margin-bottom: 40px;
  padding-bottom: 20px;
  border-bottom: 1px solid #e2e8f0;
}

.section-title {
  font-size: 2rem;
  color: #1e293b;
  margin: 0;
}

.content-body {
  max-width: 800px;
}
</style>
