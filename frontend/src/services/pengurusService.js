import { apiRequest } from '../api'

export const pengurusService = {
  list() {
    return apiRequest('/pengurus')
  },

  create(payload) {
    return apiRequest('/pengurus', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
  },

  update(id, payload) {
    return apiRequest(`/pengurus/${id}`, {
      method: 'PUT',
      body: JSON.stringify(payload),
    })
  },

  remove(id) {
    return apiRequest(`/pengurus/${id}`, { method: 'DELETE' })
  },
}
