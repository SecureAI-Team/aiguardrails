import axios from 'axios'

export const client = axios.create({
  baseURL: '',
  headers: {
    'Content-Type': 'application/json'
  }
})

client.interceptors.request.use((config) => {
  // Support both 'auth_token' (new) and 'token' (old) for compatibility
  const token = localStorage.getItem('auth_token') || localStorage.getItem('token')
  if (token) {
    config.headers['Authorization'] = `Bearer ${token}`
  } else {
    const adminToken = import.meta.env.VITE_ADMIN_TOKEN || ''
    if (adminToken) {
      config.headers['X-Admin-Token'] = adminToken
    }
  }
  return config
})

export const api = {
  async listTenants() {
    const res = await client.get('/v1/tenants')
    return res.data
  },
  async createTenant(name: string) {
    const res = await client.post('/v1/tenants', { name })
    return res.data
  },
  async listApps(tenantId: string) {
    const res = await client.get(`/v1/tenants/${tenantId}/apps`)
    return res.data
  },
  async createApp(tenantId: string, name: string, quotaPerHr: number) {
    const res = await client.post(`/v1/tenants/${tenantId}/apps`, { name, quota_per_hr: quotaPerHr })
    return res.data
  },
  async listPolicies(tenantId: string) {
    const res = await client.get(`/v1/tenants/${tenantId}/policies`)
    return res.data
  },
  async createPolicy(tenantId: string, payload: any) {
    const res = await client.post(`/v1/tenants/${tenantId}/policies`, payload)
    return res.data
  },
  async updatePolicy(tenantId: string, policyId: string, payload: any) {
    const res = await client.put(`/v1/tenants/${tenantId}/policies/${policyId}`, payload)
    return res.data
  },
  async deletePolicy(tenantId: string, policyId: string) {
    const res = await client.delete(`/v1/tenants/${tenantId}/policies/${policyId}`)
    return res.data
  },
  async listPolicyHistory(tenantId: string) {
    const res = await client.get(`/v1/tenants/${tenantId}/policies/history`)
    return res.data
  },
  async listAudit(limit = 100, event?: string, tenantId?: string) {
    const params = new URLSearchParams()
    params.set('limit', String(limit))
    if (event) params.set('event', event)
    if (tenantId) params.set('tenant_id', tenantId)
    const res = await client.get(`/v1/audit?${params.toString()}`)
    return res.data
  },
  async listCapabilities(tag?: string) {
    const q = tag ? `?tag=${encodeURIComponent(tag)}` : ''
    const res = await client.get(`/v1/capabilities${q}`)
    return res.data
  },
  async createCapability(payload: { name: string; description: string; tags: string[] }) {
    const res = await client.post(`/v1/capabilities`, payload)
    return res.data
  },
  async listRules(params: { jurisdiction?: string; regulation?: string; vendor?: string; product?: string; severity?: string; decision?: string; tag?: string } = {}) {
    const search = new URLSearchParams()
    if (params.jurisdiction) search.set('jurisdiction', params.jurisdiction)
    if (params.regulation) search.set('regulation', params.regulation)
    if (params.vendor) search.set('vendor', params.vendor)
    if (params.product) search.set('product', params.product)
    if (params.severity) search.set('severity', params.severity)
    if (params.decision) search.set('decision', params.decision)
    if (params.tag) search.set('tag', params.tag)
    const res = await client.get(`/v1/rules?${search.toString()}`)
    return res.data
  },
  async createRule(payload: any) {
    const res = await client.post('/v1/rules', payload)
    return res.data
  },
  async updateRule(id: string, payload: any) {
    const res = await client.put(`/v1/rules/${id}`, payload)
    return res.data
  },
  async deleteRule(id: string) {
    const res = await client.delete(`/v1/rules/${id}`)
    return res.data
  },
  async attachRule(policyId: string, ruleId: string) {
    const res = await client.post(`/v1/policies/${policyId}/rules/${ruleId}`)
    return res.data
  },
  async login(username: string, password: string) {
    const res = await client.post('/v1/auth/login', { username, password })
    return res.data
  },
  // User management
  async get(path: string) {
    const res = await client.get(`/v1${path}`)
    return res.data
  },
  async post(path: string, data: any) {
    const res = await client.post(`/v1${path}`, data)
    return res.data
  },
  async put(path: string, data: any) {
    const res = await client.put(`/v1${path}`, data)
    return res.data
  },
  async delete(path: string) {
    const res = await client.delete(`/v1${path}`)
    return res.data
  },
  // Tenant rules
  async listTenantRules(tenantId: string, type?: string) {
    const q = type ? `?type=${type}` : ''
    const res = await client.get(`/v1/tenants/${tenantId}/rules${q}`)
    return res.data
  },
  async createTenantRule(tenantId: string, rule: any) {
    const res = await client.post(`/v1/tenants/${tenantId}/rules`, rule)
    return res.data
  },
  async updateTenantRule(tenantId: string, ruleId: string, rule: any) {
    const res = await client.put(`/v1/tenants/${tenantId}/rules/${ruleId}`, rule)
    return res.data
  },
  async deleteTenantRule(tenantId: string, ruleId: string) {
    const res = await client.delete(`/v1/tenants/${tenantId}/rules/${ruleId}`)
    return res.data
  },
  async listRuleTemplates(type?: string) {
    const q = type ? `?type=${type}` : ''
    const res = await client.get(`/v1/rules/templates${q}`)
    return res.data
  }
}

