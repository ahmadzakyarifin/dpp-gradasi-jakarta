/**
 * Helper parsing error dari apiRequest menjadi:
 * 1. fieldErrors — map { [fieldName]: pesan } untuk ditampilkan inline di bawah input
 * 2. message — pesan umum yang aman ditampilkan (toast/banner)
 * 3. retryAfter — durasi cooldown rate limit (detik)
 *
 * Format error backend (standar helper/response.go + validator):
 *   { success: false, code: 'VALIDATION_ERROR', message: 'validasi gagal',
 *     errors: [{ field: 'title', tag: 'min', param: '5', message: 'title minimal 5 karakter' }] }
 * atau error bisnis: { success: false, code: 'DUPLICATE_TITLE', message: 'Judul sudah digunakan' }
 */
import { useState, useCallback, useEffect } from 'react'

// Peta kode error bisnis → field form yang dituju (jika relevan).
// Generalisasi dari pola yang sudah ada di BeritaAdmin.jsx (DUPLICATE_TITLE → title).
const BUSINESS_CODE_FIELD_MAP = {
  DUPLICATE_TITLE: 'title',
  DUPLICATE_SLUG: 'slug',
  DUPLICATE_EMAIL: 'email',
  DUPLICATE_NAME: 'name',
  DUPLICATE_USERNAME: 'username',
  EMAIL_ALREADY_TAKEN: 'email',
  INVALID_CREDENTIALS: 'password',
  WRONG_PASSWORD: 'current_password',
  PASSWORD_MISMATCH: 'password_confirmation',
  TOKEN_INVALID: 'token',
  TOKEN_EXPIRED: 'token',
  FILE_REQUIRED: 'image',
  FILE_TOO_LARGE: 'image',
  FILE_INVALID_TYPE: 'image',
}

/**
 * Ubah nama field dari backend (snake_case) ke nama field form FE bila berbeda.
 * Default: sama persis. Tambahkan mapping per modul di sini jika payload FE
 * memakai nama berbeda dari JSON body backend.
 */
const FIELD_NAME_MAP = {
  // contoh: 'logo_path': 'logo', // kalau suatu saat form FE pakai nama lain
}

export function parseApiError(err) {
  const fallback = {
    fieldErrors: {},
    message: err?.message || 'Terjadi kesalahan. Silakan coba lagi.',
    code: err?.code || null,
    retryAfter: Number(err?.retryAfter) || 0,
    status: err?.status || null,
  }

  const data = err?.data
  if (!data) return fallback

  const retryAfter = Number(data.retry_after) || Number(err?.retryAfter) || 0

  // 1. Error validasi standar: data.errors = array {field, tag, param, message}
  if (Array.isArray(data.errors) && data.errors.length > 0) {
    const fieldErrors = {}
    for (const e of data.errors) {
      if (e && e.field) {
        const feField = FIELD_NAME_MAP[e.field] || e.field
        // Jangan timpa pesan pertama untuk field yang sama (backend urut per field)
        if (!fieldErrors[feField]) {
          fieldErrors[feField] = e.message || `${e.field} tidak valid`
        }
      }
    }
    return {
      ...fallback,
      fieldErrors,
      retryAfter,
      message: data.message || 'Validasi gagal. Periksa kembali isian form.',
    }
  }

  // 2. Error bisnis dengan code spesifik → arahkan ke field bila ada pemetaan
  const code = data.code || err?.code
  if (code) {
    const targetField = BUSINESS_CODE_FIELD_MAP[code]
    const fieldErrors = targetField ? { [targetField]: data.message } : {}
    return {
      ...fallback,
      fieldErrors,
      retryAfter,
      message: data.message || fallback.message,
    }
  }

  // 3. Fallback: pesan umum dari backend
  return {
    ...fallback,
    retryAfter,
    message: data.message || fallback.message,
  }
}

/**
 * Hook: field errors untuk form admin.
 *   const { fieldErrors, setFieldErrors, clearFieldError } = useFormErrors()
 * Panggil applyError(err) di catch → isi fieldErrors.
 * Panggil clearFieldError('title') di onChange input → error hilang saat user memperbaiki.
 */
export function useFormErrors() {
  const [fieldErrors, setFieldErrors] = useState({})

  const applyError = useCallback((err) => {
    const parsed = parseApiError(err)
    setFieldErrors(parsed.fieldErrors || {})
    return parsed
  }, [])

  const clearFieldError = useCallback((field) => {
    setFieldErrors((prev) => {
      if (!(field in prev)) return prev
      const next = { ...prev }
      delete next[field]
      return next
    })
  }, [])

  const resetFieldErrors = useCallback(() => setFieldErrors({}), [])

  return { fieldErrors, applyError, clearFieldError, resetFieldErrors }
}

/**
 * Hook: countdown rate limit reusable.
 *   const { cooldown, isLimited } = useRateLimitCooldown()
 * Panggil applyError(err) di catch → kalau err.retryAfter > 0 mulai countdown.
 * Selama cooldown > 0, tombol submit sebaiknya disabled & tampilkan "Tunggu Xs".
 */
export function useRateLimitCooldown() {
  const [cooldown, setCooldown] = useState(0)

  useEffect(() => {
    if (cooldown <= 0) return undefined
    const timer = setInterval(() => {
      setCooldown((c) => (c > 1 ? c - 1 : 0))
    }, 1000)
    return () => clearInterval(timer)
  }, [cooldown > 0])

  const applyRateLimit = useCallback((err) => {
    const retryAfter = Number(err?.retryAfter) || Number(err?.data?.retry_after) || 0
    if (retryAfter > 0) setCooldown(retryAfter)
    return retryAfter
  }, [])

  return { cooldown, isLimited: cooldown > 0, applyRateLimit }
}
