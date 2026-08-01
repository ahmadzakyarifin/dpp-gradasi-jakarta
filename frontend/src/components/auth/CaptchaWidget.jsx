import { useEffect, useRef, useState } from 'react'

// Cloudflare Turnstile Widget
// - Load script dari challenges.cloudflare.com (hanya sekali)
// - Render widget ke container; onVerify memberikan token ke parent
// - Mode development: VITE_CAPTCHA_ENABLED=false => widget tidak dirender
//   (backend juga bypass ketika CAPTCHA_ENABLED=false)
const SCRIPT_SRC = 'https://challenges.cloudflare.com/turnstile/v0/api.js?render=explicit'

export default function CaptchaWidget({ onVerify, hasError }) {
  const containerRef = useRef(null)
  const [scriptLoaded, setScriptLoaded] = useState(false)
  const [widgetId, setWidgetId] = useState(null)
  const [expired, setExpired] = useState(false)

  const siteKey = import.meta.env.VITE_CAPTCHA_SITE_KEY || ''
  const enabled = import.meta.env.VITE_CAPTCHA_ENABLED === 'true' || import.meta.env.VITE_CAPTCHA_ENABLED === true

  // Load Turnstile script sekali
  useEffect(() => {
    if (!enabled || !siteKey) return
    if (window.turnstile) {
      setScriptLoaded(true)
      return
    }
    const script = document.createElement('script')
    script.src = SCRIPT_SRC
    script.async = true
    script.onload = () => setScriptLoaded(true)
    document.head.appendChild(script)
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
      {expired && (
        <p className="text-red-500 text-[11px] font-bold mt-1.5 ml-1">
          Kode verifikasi kedaluwarsa. Silakan selesaikan kembali.
        </p>
      )}
    </div>
  )
}
