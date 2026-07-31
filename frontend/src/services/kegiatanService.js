import { apiRequest } from '../api'

function toQuery(params = {}) {
  const query = new URLSearchParams()
  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== null && value !== '') query.set(key, value)
  })
  const text = query.toString()
  return text ? `?${text}` : ''
}

export const kegiatanService = {
  list(params) {
    return apiRequest(`/kegiatan${toQuery(params)}`)
  },

  detailBySlug(slug) {
    return apiRequest(`/kegiatan/${encodeURIComponent(slug)}`)
  },

  listAdmin(params) {
    return apiRequest(`/kegiatan/admin${toQuery(params)}`)
  },

  detailById(id) {
    return apiRequest(`/kegiatan/id/${id}`)
  },

  create(payload) {
    return apiRequest('/kegiatan', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
  },

  update(id, payload) {
    return apiRequest(`/kegiatan/${id}`, {
      method: 'PUT',
      body: JSON.stringify(payload),
    })
  },

  remove(id) {
    return apiRequest(`/kegiatan/${id}`, { method: 'DELETE' })
  },

  restore(id) {
    return apiRequest(`/kegiatan/${id}/restore`, { method: 'POST' })
  },

  bulkDelete(ids) {
    return apiRequest('/kegiatan/bulk-delete', {
      method: 'POST',
      body: JSON.stringify({ ids }),
    })
  },

  bulkRestore(ids) {
    return apiRequest('/kegiatan/bulk-restore', {
      method: 'POST',
      body: JSON.stringify({ ids }),
    })
  },
}
