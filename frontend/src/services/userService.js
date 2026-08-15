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

export const userService = {
  list(params) {
    return apiRequest(`/admin/users${toQuery(params)}`)
  },

  create(payload) {
    return apiRequest('/admin/users', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
  },

  update(id, payload) {
    return apiRequest(`/admin/users/${id}`, {
      method: 'PUT',
      body: JSON.stringify(payload),
    })
  },

  setStatus(id, status) {
    return apiRequest(`/admin/users/${id}/status`, {
      method: 'PUT',
      body: JSON.stringify({ status }),
    })
  },

  remove(id) {
    return apiRequest(`/admin/users/${id}`, {
      method: 'DELETE',
    })
  },

  restore(id) {
    return apiRequest(`/admin/users/${id}/restore`, {
      method: 'POST',
    })
  },

  resendActivation(id) {
    return apiRequest(`/admin/users/${id}/resend-activation`, {
      method: 'POST',
    })
  },

  resetPassword(id, password) {
    return apiRequest(`/admin/users/${id}/password`, {
      method: 'PUT',
      body: JSON.stringify({ password }),
    })
  },

  bulkDelete(ids) {
    return apiRequest('/admin/users/bulk-delete', {
      method: 'POST',
      body: JSON.stringify({ ids }),
    })
  },

  bulkRestore(ids) {
    return apiRequest('/admin/users/bulk-restore', {
      method: 'POST',
      body: JSON.stringify({ ids }),
    })
  },
}
