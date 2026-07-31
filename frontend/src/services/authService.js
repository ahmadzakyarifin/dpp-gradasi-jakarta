import { apiRequest } from '../api'

export const authService = {
  async login(payload) {
    return apiRequest('/auth/login', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
  },

  async forgotPassword(email) {
    return apiRequest('/auth/forgot-password', {
      method: 'POST',
      body: JSON.stringify({ email }),
    })
  },

  async validateResetToken(token) {
    return apiRequest(`/auth/validate-reset-token?token=${encodeURIComponent(token)}`)
  },

  async validateActivationToken(token) {
    return apiRequest(`/auth/validate-activation-token?token=${encodeURIComponent(token)}`)
  },

  async resetPassword(payload) {
    return apiRequest('/auth/reset-password', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
  },

  async activateAccount(payload) {
    return apiRequest('/auth/activate-account', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
  },

  async me() {
    return apiRequest('/auth/me')
  },

  async logout() {
    return apiRequest('/auth/logout', { method: 'POST' })
  },
}
