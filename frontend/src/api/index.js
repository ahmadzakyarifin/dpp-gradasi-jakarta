const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://127.0.0.1:8080/api/v1'

export async function apiRequest(path, options = {}) {
  const token = localStorage.getItem('access_token')
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
