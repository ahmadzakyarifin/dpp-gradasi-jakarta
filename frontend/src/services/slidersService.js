import { apiRequest } from '../api'
import { normalizeImage } from '../utils/normalizeImage'

function toQuery(params = {}) {
  const query = new URLSearchParams()
  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== null && value !== '') {
      query.set(key, value)
    }
  })
  const text = query.toString()
  return text ? `?${text}` : ''
}

export const slidersService = {
  // Publik — hanya slider aktif (tanpa param query, sesuai kontrak)
  list() {
    return apiRequest('/sliders').then(res => {
      if (res?.data?.sliders) res.data.sliders = res.data.sliders.map(normalizeImage)
      return res
    })
  },

  // Admin — semua (active=false) atau aktif saja
  listAdmin(active = false) {
    return apiRequest(`/sliders/admin?active=${active}`).then(res => {
      if (res?.data?.sliders) res.data.sliders = res.data.sliders.map(normalizeImage)
      return res
    })
  },

  getById(id) {
    return apiRequest(`/sliders/${id}`).then(res => {
      if (res?.data) res.data = normalizeImage(res.data)
      return res
    })
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

  restore(id) {
    return apiRequest(`/sliders/${id}/restore`, { method: 'POST' })
  },

  bulkDelete(ids) {
    return apiRequest('/sliders/bulk-delete', {
      method: 'POST',
      body: JSON.stringify({ ids }),
    })
  },

  bulkRestore(ids) {
    return apiRequest('/sliders/bulk-restore', {
      method: 'POST',
      body: JSON.stringify({ ids }),
    })
  },

  reorder(ids) {
    return apiRequest('/sliders/reorder', {
      method: 'PUT',
      body: JSON.stringify({ ids }),
    })
  },
}
