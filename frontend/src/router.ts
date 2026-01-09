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
import Landing from './pages/Landing.vue'
import Dashboard from './pages/Dashboard.vue'
import Settings from './pages/Settings.vue'
import AlertRules from './pages/AlertRules.vue'
import AlertCenter from './pages/AlertCenter.vue'
import Stats from './pages/Stats.vue'
import ApiKeys from './pages/ApiKeys.vue'
import Traces from './pages/Traces.vue'
import SDKs from './pages/SDKs.vue'
import Playground from './pages/Playground.vue'
import Models from './pages/Models.vue'
import Organizations from './pages/Organizations.vue'
import Docs from './pages/Docs.vue'
import APIReference from './pages/APIReference.vue'
import BestPractices from './pages/BestPractices.vue'
import About from './pages/About.vue'
import Contact from './pages/Contact.vue'
import Privacy from './pages/Privacy.vue'

const routes = [
  { path: '/', component: Landing },
  { path: '/landing', component: Landing },
  { path: '/login', component: Login },
  { path: '/dashboard', component: Dashboard },
  { path: '/tenants', component: Tenants },
  { path: '/apps', component: Apps },
  { path: '/policies', component: Policies },
  { path: '/policy-history', component: PolicyHistory },
  { path: '/rules', component: Rules },
  { path: '/capabilities', component: Capabilities },
  { path: '/logs', component: Logs },
  { path: '/users', component: Users },
  { path: '/tenants/:tenantId/users', component: TenantUsers },
  { path: '/profile', component: Profile },
  { path: '/settings', component: Settings },
  { path: '/alerts', component: AlertCenter },
  { path: '/alerts/rules', component: AlertRules },
  { path: '/stats', component: Stats },
  { path: '/apikeys', component: ApiKeys },
  { path: '/traces', component: Traces },
  { path: '/sdks', component: SDKs },
  { path: '/playground', component: Playground },
  { path: '/models', component: Models },
  { path: '/orgs', component: Organizations },
  { path: '/docs', component: Docs },
  { path: '/api-reference', component: APIReference },
  { path: '/best-practices', component: BestPractices },
  { path: '/about', component: About },
  { path: '/contact', component: Contact },
  { path: '/privacy', component: Privacy }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

const publicPages = ['/', '/landing', '/login', '/sdks', '/playground', '/models', '/docs', '/api-reference', '/best-practices', '/about', '/contact', '/privacy']

router.beforeEach((to, from, next) => {
  const publicPages = ['/', '/landing', '/login', '/sdks', '/playground', '/models', '/docs', '/api-reference', '/best-practices', '/about', '/contact', '/privacy']
  const authRequired = !publicPages.includes(to.path)
  const token = localStorage.getItem('token') || localStorage.getItem('auth_token')

  if (authRequired && !token) {
    return next('/login')
  }

  next()
})

export default router
