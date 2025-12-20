import { describe, it, expect } from 'vitest'
import { client } from '../src/services/api'

describe('api interceptor', () => {
  it('adds bearer token when available', async () => {
    localStorage.setItem('auth_token', 'abc123')
    const cfg: any = { headers: {} }
    const out = await (client.interceptors.request as any).handlers[0].fulfilled(cfg)
    expect(out.headers.Authorization).toBe('Bearer abc123')
    localStorage.removeItem('auth_token')
  })
})

