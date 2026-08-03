import { useEffect, useRef, useState } from 'react'
import { useSettings } from '../../context/useSettings'

// Cloudflare Turnstile Widget
// - Sumber kebenaran status CAPTCHA = config backend (GET /settings), bukan env VITE_*.
// - Load script dari challenges.cloudflare.com (hanya sekali)
// - Render widget ke container; onVerify memberikan token ke parent
const SCRIPT_SRC = 'https://challenges.cloudflare.com/turnstile/v0/api.js?render=explicit'
const SCRIPT_TIMEOUT_MS = 5000

export default function CaptchaWidget({ onVerify, hasError }) {
  const containerRef = useRef(null)
  const [scriptLoaded, setScriptLoaded] = useState(false)
  const [scriptFailed, setScriptFailed] = useState(false)
  const [widgetId, setWidgetId] = useState(null)
  const [expired, setExpired] = useState(false)

  // Status CAPTCHA dari backend (single source of truth — tidak ada lagi
  // kemungkinan FE aktif tapi BE tidak aktif atau sebaliknya).
  const { settings } = useSettings()
  const enabled = !!settings?.captcha_enabled
  const siteKey = settings?.captcha_site_key || ''

  // Load Turnstile script sekali + timeout: kalau gagal dimuat (adblocker/
  // jaringan), tampilkan pesan eksplisit, jangan biarkan widget kosong diam-diam.
  useEffect(() => {
    if (!enabled || !siteKey) return
    if (window.turnstile) {
      setScriptLoaded(true)
      return
    }
    let timedOut = false
    const timer = setTimeout(() => {
      if (!window.turnstile) {
        timedOut = true
        setScriptFailed(true)
      }
    }, SCRIPT_TIMEOUT_MS)
    const script = document.createElement('script')
    script.src = SCRIPT_SRC
    script.async = true
    script.onload = () => {
      clearTimeout(timer)
      setScriptLoaded(true)
    }
    script.onerror = () => {
      clearTimeout(timer)
      timedOut = true
      setScriptFailed(true)
    }
    document.head.appendChild(script)
    return () => {
      clearTimeout(timer)
      // eslint-disable-next-line no-unused-expressions
      timedOut
    }
  }, [enabled, siteKey])

  // Render widget saat script siap
  useEffect(() => {
    if (!enabled || !siteKey || !scriptLoaded || !containerRef.current) return
    if (window.turnstile) {
      const id = window.turnstile.render(containerRef.current, {
        sitekey: siteKey,
        callback: (token) => {
          setExpired(false)
          onVerify?.(token)
        },
        'expired-callback': () => {
          setExpired(true)
          onVerify?.('')
        },
        'error-callback': () => {
          setExpired(true)
          onVerify?.('')
        },
        theme: 'light',
      })
      setWidgetId(id)
    }
    return () => {
      if (widgetId && window.turnstile) {
        window.turnstile.remove(widgetId)
      }
    }
  }, [enabled, siteKey, scriptLoaded]) // eslint-disable-line react-hooks/exhaustive-deps

  // Reset widget saat hasError berubah jadi true (retry)
  useEffect(() => {
    if (hasError && widgetId && window.turnstile) {
      window.turnstile.reset(widgetId)
      onVerify?.('')
    }
  }, [hasError, widgetId]) // eslint-disable-line react-hooks/exhaustive-deps

  if (!enabled || !siteKey) return null

  return (
    <div className="space-y-2">
      <label className="block text-xs font-semibold text-slate-600 uppercase tracking-wider">
        Verifikasi Keamanan
      </label>
      <div
        ref={containerRef}
        className={`inline-block rounded-xl overflow-hidden border transition ${
          hasError || expired ? 'border-red-400 ring-2 ring-red-500/20' : 'border-slate-200'
        }`}
      />
      {scriptFailed && (
        <p className="text-red-500 text-[11px] font-bold mt-1.5 ml-1">
          CAPTCHA gagal dimuat. Periksa koneksi atau nonaktifkan adblocker, lalu refresh halaman.
        </p>
      )}
      {expired && !scriptFailed && (
        <p className="text-red-500 text-[11px] font-bold mt-1.5 ml-1">
          Kode verifikasi kedaluwarsa. Silakan selesaikan kembali.
        </p>
      )}
    </div>
  )
}
