import { apiRequest } from '../api'

export const settingsService = {
  get() {
    return apiRequest('/settings')
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

  uploadSign1(file) {
    const formData = new FormData()
    formData.append('image', file)
    return apiRequest('/admin/settings/sign1', {
      method: 'POST',
      body: formData,
    })
  },

  uploadSign2(file) {
    const formData = new FormData()
    formData.append('image', file)
    return apiRequest('/admin/settings/sign2', {
      method: 'POST',
      body: formData,
    })
  },

  uploadGreetingImage(file) {
    const formData = new FormData()
    formData.append('image', file)
    return apiRequest('/admin/settings/greeting-image', {
      method: 'POST',
      body: formData,
    })
  },
}
