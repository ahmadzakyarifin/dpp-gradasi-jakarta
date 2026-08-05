const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://127.0.0.1:8080/api/v1'

// State untuk mencegah race-condition saat banyak request 401 bersamaan:
// hanya satu proses refresh yang berjalan, sisanya menunggu promise yang sama.
let refreshPromise = null

async function refreshAccessToken() {
  // Refresh token dikirim otomatis oleh browser via HttpOnly cookie (path /api/v1/auth).
  // Tidak perlu menyertakan body atau header Authorization.
  const res = await fetch(`${API_BASE_URL}/auth/refresh`, {
    method: 'POST',
    credentials: 'include',
    headers: { Accept: 'application/json' },
  })
  const data = await res.json().catch(() => null)
  if (!res.ok || !data?.data?.access_token) {
    throw new Error(data?.message || 'Sesi berakhir, silakan login kembali.')
  }
  return data.data.access_token
}

function getAccessToken() {
  return localStorage.getItem('access_token')
}

export async function apiRequest(path, options = {}) {
  const token = getAccessToken()
  const headers = {
    Accept: 'application/json',
    ...(options.body instanceof FormData ? {} : { 'Content-Type': 'application/json' }),
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
    ...(options.headers || {}),
  }

  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...options,
    headers,
    credentials: 'include',
  })

  let data = null
  try {
    data = await response.json()
  } catch {
    data = null
  }

  // 401 pada endpoint auth/refresh atau login → jangan auto-refresh (hindari loop)
  const isAuthRefreshPath = path.includes('/auth/refresh')
  const isLoginPath = path.includes('/auth/login')

  if (response.status === 401 && !isAuthRefreshPath && !isLoginPath && getAccessToken()) {
    // Coba refresh token sekali. Jika berhasil, ulangi request asli dengan token baru.
    try {
      if (!refreshPromise) {
        refreshPromise = refreshAccessToken()
          .then((newToken) => {
            localStorage.setItem('access_token', newToken)
            return newToken
          })
          .finally(() => {
            refreshPromise = null
          })
      }
      const newToken = await refreshPromise
      const retryHeaders = {
        ...headers,
        Authorization: `Bearer ${newToken}`,
      }
      const retryResponse = await fetch(`${API_BASE_URL}${path}`, {
        ...options,
        headers: retryHeaders,
        credentials: 'include',
      })
      let retryData = null
      try {
        retryData = await retryResponse.json()
      } catch {
        retryData = null
      }
      if (!retryResponse.ok || retryData?.success === false) {
        const error = new Error(retryData?.message || 'Terjadi kesalahan pada server')
        error.status = retryResponse.status
        error.code = retryData?.code
        error.data = retryData
        error.retryAfter = retryData?.retry_after || Number(retryResponse.headers.get('Retry-After')) || 0
        throw error
      }
      return retryData
    } catch (refreshErr) {
      // Refresh gagal → sesi benar-benar berakhir, bersihkan token lokal.
      localStorage.removeItem('access_token')
      localStorage.removeItem('user_role')
      window.location.href = '/login?expired=1'
      const error = new Error(refreshErr?.message || 'Sesi berakhir, silakan login kembali.')
      error.status = 401
      error.code = 'AUTH_SESSION_EXPIRED'
      error.data = refreshErr?.data || null
      throw error
    }
  }

  if (!response.ok || data?.success === false) {
    const error = new Error(data?.message || 'Terjadi kesalahan pada server')
    error.status = response.status
    error.code = data?.code
    error.data = data
    // retry_after (detik) dari response rate limit — dipakai UI untuk countdown
    error.retryAfter = data?.retry_after || Number(response.headers.get('Retry-After')) || 0
    throw error
  }

  return data
}

export default apiRequest
