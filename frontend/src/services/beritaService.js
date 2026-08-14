import { apiRequest } from '../api'
import { normalizeImage } from '../utils/normalizeImage'

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
    return apiRequest(`/berita${toQuery(params)}`).then(res => {
      if (res?.data?.berita) res.data.berita = res.data.berita.map(normalizeImage)
      return res
    })
  },

  detailBySlug(slug) {
    return apiRequest(`/berita/${encodeURIComponent(slug)}`).then(res => {
      if (res?.data) res.data = normalizeImage(res.data)
      return res
    })
  },

  getCategories() {
    return apiRequest('/berita/categories')
  },

  listAdmin(params) {
    return apiRequest(`/berita/admin${toQuery(params)}`).then(res => {
      if (res?.data?.berita) res.data.berita = res.data.berita.map(normalizeImage)
      return res
    })
  },

  detailById(id) {
    return apiRequest(`/berita/id/${id}`).then(res => {
      if (res?.data) res.data = normalizeImage(res.data)
      return res
    })
  },

  create(payload) {
    return apiRequest('/berita', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
  },

  uploadImage(file) {
    const fd = new FormData()
    fd.append('image', file)
    return apiRequest('/berita/upload-image', {
      method: 'POST',
      body: fd,
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

  renameCategory(oldName, newName) {
    return apiRequest('/berita/categories', {
      method: 'PUT',
      body: JSON.stringify({ old_name: oldName, new_name: newName }),
    })
  },

  deleteCategory(name) {
    return apiRequest(`/berita/categories/${encodeURIComponent(name)}`, {
      method: 'DELETE',
    })
  },
}
