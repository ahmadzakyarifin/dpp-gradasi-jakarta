import { apiRequest } from '../api'

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

export const kontakService = {
  list(params) {
    return apiRequest(`/admin/kontak${toQuery(params)}`)
  },

  getById(id) {
    return apiRequest(`/admin/kontak/${id}`)
  },

  remove(id) {
    return apiRequest(`/admin/kontak/${id}`, {
      method: 'DELETE',
    })
  },

  restore(id) {
    return apiRequest(`/admin/kontak/${id}/restore`, {
      method: 'POST',
    })
  },

  bulkDelete(ids) {
    return apiRequest('/admin/kontak/bulk-delete', {
      method: 'POST',
      body: JSON.stringify({ ids }),
    })
  },

  bulkRestore(ids) {
    return apiRequest('/admin/kontak/bulk-restore', {
      method: 'POST',
      body: JSON.stringify({ ids }),
    })
  },
}
