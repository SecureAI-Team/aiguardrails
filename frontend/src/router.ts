import { createRouter, createWebHistory } from 'vue-router'
import Tenants from './pages/Tenants.vue'
import Apps from './pages/Apps.vue'
import Policies from './pages/Policies.vue'
import Logs from './pages/Logs.vue'
import PolicyHistory from './pages/PolicyHistory.vue'
import Capabilities from './pages/Capabilities.vue'

const routes = [
  { path: '/', component: Tenants },
  { path: '/apps', component: Apps },
  { path: '/policies', component: Policies },
  { path: '/policy-history', component: PolicyHistory },
  { path: '/capabilities', component: Capabilities },
  { path: '/logs', component: Logs }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router

