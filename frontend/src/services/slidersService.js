import { apiRequest } from '../api'

export const slidersService = {
  list(activeOnly = true) {
    return apiRequest(`/sliders?active_only=${activeOnly}`)
  },

  create(payload) {
    return apiRequest('/sliders', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
  },

  update(id, payload) {
    return apiRequest(`/sliders/${id}`, {
      method: 'PUT',
      body: JSON.stringify(payload),
    })
  },

  remove(id) {
    return apiRequest(`/sliders/${id}`, { method: 'DELETE' })
  },
}
