import { createRouter, createWebHistory } from 'vue-router'
import Tenants from './pages/Tenants.vue'
import Apps from './pages/Apps.vue'
import Policies from './pages/Policies.vue'
import Logs from './pages/Logs.vue'
import PolicyHistory from './pages/PolicyHistory.vue'
import Capabilities from './pages/Capabilities.vue'
import Rules from './pages/Rules.vue'
import Login from './pages/Login.vue'
import Users from './pages/Users.vue'
import TenantUsers from './pages/TenantUsers.vue'
import Profile from './pages/Profile.vue'

const routes = [
  { path: '/login', component: Login },
  { path: '/', component: Tenants },
  { path: '/apps', component: Apps },
  { path: '/policies', component: Policies },
  { path: '/policy-history', component: PolicyHistory },
  { path: '/rules', component: Rules },
  { path: '/capabilities', component: Capabilities },
  { path: '/logs', component: Logs },
  { path: '/users', component: Users },
  { path: '/tenants/:tenantId/users', component: TenantUsers },
  { path: '/profile', component: Profile }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
