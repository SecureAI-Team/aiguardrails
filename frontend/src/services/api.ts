import axios from 'axios'

const client = axios.create({
  baseURL: import.meta.env.VITE_API_BASE || 'http://localhost:8080',
  headers: {
    'Content-Type': 'application/json',
    'X-Admin-Token': import.meta.env.VITE_ADMIN_TOKEN || ''
  }
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
    const res = await client.get(`/v1/mcp/capabilities${q}`)
    return res.data
  },
  async createCapability(payload: { name: string; description: string; tags: string[] }) {
    const res = await client.post(`/v1/capabilities`, payload)
    return res.data
  },
  async listRules(params: { jurisdiction?: string; regulation?: string; vendor?: string; product?: string }) {
    const search = new URLSearchParams()
    if (params.jurisdiction) search.set('jurisdiction', params.jurisdiction)
    if (params.regulation) search.set('regulation', params.regulation)
    if (params.vendor) search.set('vendor', params.vendor)
    if (params.product) search.set('product', params.product)
    const res = await client.get(`/v1/rules?${search.toString()}`)
    return res.data
  },
  async attachRule(policyId: string, ruleId: string) {
    const res = await client.post(`/v1/policies/${policyId}/rules/${ruleId}`)
    return res.data
  }
}

