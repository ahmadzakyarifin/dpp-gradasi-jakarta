// Helper untuk resolve URL file upload lokal (path relatif /uploads/...) vs URL absolut.
// Backend menyimpan path relatif (mis. /uploads/settings/123.png) — FE harus prepend base URL.
// Base URL default mengikuti API base di src/api/index.js.
const BASE_URL = import.meta.env.VITE_API_URL?.startsWith('http')
  ? new URL(import.meta.env.VITE_API_URL).origin
  : '' // Gunakan path relatif (same-origin) jika menggunakan proxy (/api/v1)

export function resolveAssetUrl(url) {
  if (!url) return ''
  if (url.startsWith('http://') || url.startsWith('https://') || url.startsWith('data:')) return url
  if (url.startsWith('/')) return `${BASE_URL}${url}`
  return url
}
