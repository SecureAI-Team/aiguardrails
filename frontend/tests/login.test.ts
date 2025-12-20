import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import Login from '../src/pages/Login.vue'

vi.mock('../src/services/api', () => ({
  api: {
    login: vi.fn().mockResolvedValue({ token: 'tkn' })
  }
}))

vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: vi.fn()
  })
}))

describe('Login page', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it('stores token on success', async () => {
    const wrapper = mount(Login)
    await wrapper.find('input[placeholder="Username"]').setValue('u')
    await wrapper.find('input[placeholder="Password"]').setValue('p')
    await wrapper.find('form').trigger('submit.prevent')
    expect(localStorage.getItem('auth_token')).toBe('tkn')
  })
})

