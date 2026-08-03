import { apiRequest } from '../api'

function toQuery(params = {}) {
  const query = new URLSearchParams()
  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== null && value !== '') query.set(key, value)
  })
  const text = query.toString()
  return text ? `?${text}` : ''
}

// CRUD role — backend: /api/v1/roles (super_admin, admin)
export const roleService = {
  list(params = {}) {
    return apiRequest(`/roles${toQuery(params)}`)
  },

  getById(id) {
    return apiRequest(`/roles/${id}`)
  },

  create(payload) {
    return apiRequest('/roles', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
  },

  update(id, payload) {
    return apiRequest(`/roles/${id}`, {
      method: 'PUT',
      body: JSON.stringify(payload),
    })
  },

  remove(id) {
    return apiRequest(`/roles/${id}`, { method: 'DELETE' })
  },

  restore(id) {
    return apiRequest(`/roles/${id}/restore`, { method: 'PATCH' })
  },

  updateStatus(id, isActive) {
    return apiRequest(`/roles/${id}/status`, {
      method: 'PATCH',
      body: JSON.stringify({ is_active: isActive }),
    })
  },

  bulkDelete(ids) {
    return apiRequest('/roles/bulk-delete', {
      method: 'POST',
      body: JSON.stringify({ ids }),
    })
  },

  bulkRestore(ids) {
    return apiRequest('/roles/bulk-restore', {
      method: 'PATCH',
      body: JSON.stringify({ ids }),
    })
  },

  dependencyInfo(id) {
    return apiRequest(`/roles/${id}/dependency-info`)
  },

  checkUnique(field, value, excludeId) {
    return apiRequest(`/roles/check-unique${toQuery({ field, value, exclude_id: excludeId })}`)
  },
}
