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

export const pengurusService = {
  // Publik — daftar pengurus aktif (tanpa auth); filter level/provinsi/kabupaten
  list(params = {}) {
    return apiRequest(`/pengurus${toQuery(params)}`).then(res => {
      if (Array.isArray(res?.data)) {
        res.data = res.data.map(normalizeImage)
      } else if (res?.data && Array.isArray(res.data.data)) {
        res.data.data = res.data.data.map(normalizeImage)
      }
      return res
    })
  },

  // Publik — daftar wilayah (provinsi + kabupaten) untuk dropdown filter
  regions() {
    return apiRequest('/pengurus/regions')
  },

  listAdmin(params = {}) {
    return apiRequest(`/admin/pengurus${toQuery(params)}`).then(res => {
      if (Array.isArray(res?.data)) {
        res.data = res.data.map(normalizeImage)
      } else if (res?.data && Array.isArray(res.data.data)) {
        res.data.data = res.data.data.map(normalizeImage)
      }
      return res
    })
  },

  // Backend pakai multipart/form-data (c.ShouldBind) — jangan JSON.stringify
  // `image` = File object (upload), `image_url` = path lama saat edit tanpa ganti foto
  create(payload) {
    const form = new FormData()
    Object.entries(payload).forEach(([key, value]) => {
      if (value === undefined || value === null) return
      if (key === 'id') return
      if (key === 'image') {
        if (value instanceof File) form.append('image', value)
        return
      }
      if (key === 'image_url' && payload.image instanceof File) return // ganti image_url via upload
      form.append(key, value)
    })
    return apiRequest('/admin/pengurus', {
      method: 'POST',
      body: form,
    })
  },

  update(id, payload) {
    const form = new FormData()
    Object.entries(payload).forEach(([key, value]) => {
      if (value === undefined || value === null) return
      if (key === 'id') return
      if (key === 'image') {
        if (value instanceof File) form.append('image', value)
        return
      }
      if (key === 'image_url' && payload.image instanceof File) return
      form.append(key, value)
    })
    return apiRequest(`/admin/pengurus/${id}`, {
      method: 'PUT',
      body: form,
    })
  },

  remove(id) {
    return apiRequest(`/admin/pengurus/${id}`, { method: 'DELETE' })
  },

  restore(id) {
    return apiRequest(`/admin/pengurus/${id}/restore`, { method: 'POST' })
  },

  bulkDelete(ids) {
    return apiRequest('/admin/pengurus/bulk-delete', {
      method: 'POST',
      body: JSON.stringify({ ids }),
    })
  },

  bulkRestore(ids) {
    return apiRequest('/admin/pengurus/bulk-restore', {
      method: 'POST',
      body: JSON.stringify({ ids }),
    })
  },

  reorder(ids) {
    return apiRequest('/admin/pengurus/reorder', {
      method: 'PUT',
      body: JSON.stringify({ ids }),
    })
  },
}
