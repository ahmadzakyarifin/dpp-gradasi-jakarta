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

export const kegiatanService = {
  list(params) {
    return apiRequest(`/kegiatan${toQuery(params)}`).then(res => {
      if (res?.data?.kegiatan) res.data.kegiatan = res.data.kegiatan.map(normalizeImage)
      return res
    })
  },

  detailBySlug(slug) {
    return apiRequest(`/kegiatan/${encodeURIComponent(slug)}`).then(res => {
      if (res?.data) res.data = normalizeImage(res.data)
      return res
    })
  },

  getCategories() {
    return apiRequest('/kegiatan/categories')
  },

  listAdmin(params) {
    return apiRequest(`/kegiatan/admin${toQuery(params)}`).then(res => {
      if (res?.data?.kegiatan) res.data.kegiatan = res.data.kegiatan.map(normalizeImage)
      return res
    })
  },

  detailById(id) {
    return apiRequest(`/kegiatan/id/${id}`).then(res => {
      if (res?.data) res.data = normalizeImage(res.data)
      return res
    })
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

  uploadImage(file) {
    const fd = new FormData()
    fd.append('image', file)
    return apiRequest('/kegiatan/upload-image', {
      method: 'POST',
      body: fd,
    })
  },

  removeGallery(galleryId) {
    return apiRequest(`/kegiatan/gallery/${galleryId}`, { method: 'DELETE' })
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
