import { apiRequest } from '../api'

function toQuery(params = {}) {
  const query = new URLSearchParams()
  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== null && value !== '') query.set(key, value)
  })
  const text = query.toString()
  return text ? `?${text}` : ''
}

export const beritaService = {
  list(params) {
    return apiRequest(`/berita${toQuery(params)}`)
  },

  detailBySlug(slug) {
    return apiRequest(`/berita/${encodeURIComponent(slug)}`)
  },

  listAdmin(params) {
    return apiRequest(`/berita/admin${toQuery(params)}`)
  },

  detailById(id) {
    return apiRequest(`/berita/id/${id}`)
  },

  create(payload) {
    return apiRequest('/berita', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
  },

  update(id, payload) {
    return apiRequest(`/berita/${id}`, {
      method: 'PUT',
      body: JSON.stringify(payload),
    })
  },

  remove(id) {
    return apiRequest(`/berita/${id}`, { method: 'DELETE' })
  },

  restore(id) {
    return apiRequest(`/berita/${id}/restore`, { method: 'POST' })
  },

  bulkDelete(ids) {
    return apiRequest('/berita/bulk-delete', {
      method: 'POST',
      body: JSON.stringify({ ids }),
    })
  },

  bulkRestore(ids) {
    return apiRequest('/berita/bulk-restore', {
      method: 'POST',
      body: JSON.stringify({ ids }),
    })
  },
}
