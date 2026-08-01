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

  uploadLogo(file) {
    const formData = new FormData()
    formData.append('logo', file)
    return apiRequest('/admin/settings/logo', {
      method: 'POST',
      body: formData,
    })
  },
}
