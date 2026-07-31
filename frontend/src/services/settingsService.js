import { apiRequest } from '../api'

export const settingsService = {
  get() {
    return apiRequest('/settings')
  },

  getAdmin() {
    return apiRequest('/admin/settings')
  },

  update(payload) {
    return apiRequest('/admin/settings', {
      method: 'PUT',
      body: JSON.stringify(payload),
    })
  },
}
