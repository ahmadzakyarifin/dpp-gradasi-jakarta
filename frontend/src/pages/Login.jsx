import { useEffect, useMemo, useState } from 'react'
import { Link, useNavigate, useSearchParams } from 'react-router-dom'
import AuthShell from '../components/auth/AuthShell'
import { useAuthStore } from '../store/useAuthStore'
import CaptchaWidget from '../components/auth/CaptchaWidget'
import { authContent } from '../content/authContent'
import { authService } from '../services/authService'
import { validateEmail, validatePassword } from '../utils/validation'
import { useSettings } from '../context/useSettings'
import { resolveAssetUrl } from '../utils/assetUrl'

export default function Login() {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const login = useAuthStore((state) => state.login)
  // Status CAPTCHA dari backend (single source of truth)
  const { settings } = useSettings()
  const captchaEnabled = !!settings?.captcha_enabled

  const initialView = searchParams.get('view') === 'forgot' ? 'forgot' : 'login'
  const [view, setView] = useState(initialView)
  const [loginForm, setLoginForm] = useState({ email: '', password: '', rememberMe: false })
  const [forgotEmail, setForgotEmail] = useState('')
  const [captchaToken, setCaptchaToken] = useState('')
  const [forgotCaptchaToken, setForgotCaptchaToken] = useState('')
  const [captchaError, setCaptchaError] = useState(false)
  const [forgotCaptchaError, setForgotCaptchaError] = useState(false)
  const [showPassword, setShowPassword] = useState(false)
  const [touched, setTouched] = useState({ email: false, password: false, forgotEmail: false })
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [notice, setNotice] = useState({ type: '', message: '' })
  const [cooldown, setCooldown] = useState(0)

  // Reset form setiap kali halaman login dimuat (termasuk setelah logout)
  useEffect(() => {
    const savedEmail = localStorage.getItem('remembered_email')
    setLoginForm({ email: savedEmail || '', password: '', rememberMe: !!savedEmail })
    setCaptchaToken('')
    setCaptchaError(false)

    if (searchParams.get('expired') === '1') {
      setNotice({ type: 'error', message: 'Sesi login telah berakhir. Silakan login kembali.' })
    }
  }, [searchParams])

  // Countdown saat rate limit (429) — tampilkan "coba lagi dalam X detik" & disable tombol
  useEffect(() => {
    if (cooldown <= 0) return
    const timer = setInterval(() => setCooldown((s) => s - 1), 1000)
    return () => clearInterval(timer)
  }, [cooldown > 0]) // eslint-disable-line react-hooks/exhaustive-deps

  const content = view === 'forgot' ? authContent.forgot : authContent.login

  const loginErrors = useMemo(() => {
    const next = {}
    const emailError = validateEmail(loginForm.email)
    if (emailError) next.email = emailError
    const passwordError = validatePassword(loginForm.password, 'Password')
    if (passwordError) next.password = passwordError
    return next
  }, [loginForm])

  const forgotErrors = useMemo(() => {
    const next = {}
    const emailError = validateEmail(forgotEmail)
    if (emailError) next.forgotEmail = emailError
    return next
  }, [forgotEmail])

  function switchToForgot() {
    setView('forgot')
    setTouched({ email: false, password: false, forgotEmail: false })
    setNotice({ type: '', message: '' })
    setCooldown(0)
    setForgotCaptchaToken('')
    setForgotCaptchaError(false)
  }

  function switchToLogin() {
    setView('login')
    setTouched({ email: false, password: false, forgotEmail: false })
    setNotice({ type: '', message: '' })
    setCooldown(0)
    setCaptchaToken('')
    setCaptchaError(false)
  }

  async function handleLogin(event) {
    event.preventDefault()
    setTouched((prev) => ({ ...prev, email: true, password: true }))
    setNotice({ type: '', message: '' })
    setCaptchaError(false)

    if (Object.keys(loginErrors).length > 0) return

    // Resolusi token CAPTCHA: state React (callback onVerify) atau fallback
    // hidden input cf-turnstile-response yang ditulis langsung oleh Cloudflare.
    let token = captchaToken
    if (captchaEnabled && !token) {
      // Tunggu widget selesai render (maks ~2 detik) — user yang klik submit
      // terlalu cepat sebelum Turnstile selesai load akan tetap berhasil.
      for (let i = 0; i < 10 && !token; i++) {
        await new Promise((r) => setTimeout(r, 200))
        token = document.querySelector('input[name="cf-turnstile-response"]')?.value || ''
      }
    }
    // JANGAN hard-block submit di sini: backend yang memutuskan apakah CAPTCHA
    // wajib (CAPTCHA_ENABLED=true + environment production). Kalau token kosong
    // tapi backend tidak butuh (mode dev / captcha off), login tetap jalan.
    // Kalau backend butuh dan token kosong, ia akan balas AUTH_CAPTCHA_REQUIRED
    // yang ditampilkan di notice + reset widget.

    setIsSubmitting(true)
    try {
      await login({
        ...loginForm,
        captcha_token: token
      })

      // Simpan/hapus email berdasarkan rememberMe
      if (loginForm.rememberMe) {
        localStorage.setItem('remembered_email', loginForm.email)
      } else {
        localStorage.removeItem('remembered_email')
      }

      // Kalau user wajib ganti password default (login pertama), arahkan ke halaman profil
      const user = useAuthStore.getState().user
      if (user?.must_change_password) {
        navigate('/admin/profile?force=1')
      } else {
        navigate(authContent.adminPath)
      }
    } catch (error) {
      // CAPTCHA diminta backend tapi token kosong/tidak valid → reset widget
      // supaya user bisa coba sekali lagi.
      if (error?.code === 'AUTH_CAPTCHA_REQUIRED' || error?.data?.code === 'AUTH_CAPTCHA_REQUIRED') {
        setCaptchaError(true)
        setNotice({ type: 'error', message: 'Silakan selesaikan verifikasi CAPTCHA terlebih dahulu.' })
        return
      }
      if (error.retryAfter > 0) {
        setCooldown(error.retryAfter)
        setNotice({ type: 'error', message: 'Terlalu banyak percobaan. Silakan tunggu sebelum mencoba lagi.' })
      } else {
        setNotice({ type: 'error', message: error.message || 'Login gagal' })
      }
      // Reset inputan email & password (biar aman dari keylogger / salah input berulang)
      setLoginForm((prev) => ({ ...prev, email: '', password: '' }))
      setTouched({ email: false, password: false, forgotEmail: false })
      setCaptchaToken('')
      setCaptchaError(false)
    } finally {
      setIsSubmitting(false)
    }
  }

  async function handleForgot(event) {
    event.preventDefault()
    setTouched((prev) => ({ ...prev, forgotEmail: true }))
    setNotice({ type: '', message: '' })
    setForgotCaptchaError(false)
    if (Object.keys(forgotErrors).length > 0) return

    if (captchaEnabled && !forgotCaptchaToken) {
      // JANGAN hard-block: backend yang memutuskan (dev mode = skip captcha).
      // Kalau backend butuh captcha, ia balas AUTH_CAPTCHA_REQUIRED.
      setForgotCaptchaError(false)
    }

    setIsSubmitting(true)
    try {
      const response = await authService.forgotPassword(forgotEmail, forgotCaptchaToken)
      setNotice({
        type: 'success',
        message: response?.message || authContent.forgot.successMessage,
      })
      setForgotEmail('')
      setForgotCaptchaToken('')
      setForgotCaptchaError(false)
      setTouched({ email: false, password: false, forgotEmail: false })
    } catch (error) {
      if (error.retryAfter > 0) {
        setCooldown(error.retryAfter)
        setNotice({ type: 'error', message: 'Terlalu banyak permintaan. Silakan tunggu sebelum mencoba lagi.' })
      } else {
        setNotice({ type: 'error', message: error.message || 'Gagal mengirim link reset password' })
      }
    } finally {
      setIsSubmitting(false)
    }
  }

  const fieldClass = (hasError) =>
    `w-full pl-11 pr-4 py-3 bg-slate-50 border rounded-xl focus:outline-none focus:bg-white transition-all text-sm font-medium shadow-sm group-hover:border-slate-300 ${
      hasError
        ? 'border-red-400 focus:ring-2 focus:ring-red-500/20 focus:border-red-500 text-red-900 placeholder-red-300'
        : 'border-slate-200 focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 text-slate-700 placeholder-slate-400'
    }`

  return (
    <AuthShell variant="login">
      <Link to={authContent.homePath} className="md:hidden absolute top-6 left-6 inline-flex items-center gap-2 text-slate-500 hover:text-brand-600 transition group text-sm font-semibold bg-slate-50 px-4 py-2 rounded-full border border-slate-200 shadow-sm">
        <i className="ph-bold ph-arrow-left group-hover:-translate-x-1 transition-transform" />
        Beranda
      </Link>

      <div className="w-full max-w-md mt-12 md:mt-0">
        <div className="mb-10 text-center md:text-left">
          <div className="md:hidden inline-flex w-16 h-16 bg-white rounded-xl p-2 shadow-lg mb-6 border border-slate-100">
            <img src={settings?.logo_path ? resolveAssetUrl(settings.logo_path) : ""} alt={settings?.site_name || ""} className="w-full h-full object-contain" />
          </div>
          <h2 className="font-heading text-3xl font-extrabold text-slate-900 mb-2 tracking-tight">{content.title}</h2>
          <p className="text-slate-500 text-sm font-medium">{content.subtitle}</p>
        </div>

        {notice.message && (
          <div className={`mb-5 rounded-xl border px-4 py-3 text-sm font-medium ${notice.type === 'success' ? 'bg-emerald-50 text-emerald-700 border-emerald-200' : 'bg-red-50 text-red-700 border-red-200'}`}>
            {notice.message}
            {cooldown > 0 && (
              <span className="block mt-1 font-bold tabular-nums">
                Coba lagi dalam {cooldown} detik.
              </span>
            )}
          </div>
        )}

        {view === 'login' ? (
          <form onSubmit={handleLogin} autoComplete="off" className="w-full flex flex-col gap-5">
            <div>
              <label className="block text-sm font-bold text-slate-700 mb-1.5">{authContent.login.emailLabel} <span className="text-red-500">*</span></label>
              <div className="relative group">
                <div className={`absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none transition-colors ${touched.email && loginErrors.email ? 'text-red-400' : 'text-slate-400 group-focus-within:text-brand-600'}`}>
                  <i className="ph-bold ph-envelope-simple text-lg" />
                </div>
                <input
                  type="email"
                  autoComplete="off"
                  value={loginForm.email}
                  onChange={(event) => {
                    setLoginForm((prev) => ({ ...prev, email: event.target.value }))
                    setTouched((prev) => ({ ...prev, email: true }))
                  }}
                  onBlur={() => setTouched((prev) => ({ ...prev, email: true }))}
                  placeholder={authContent.login.emailPlaceholder}
                  className={fieldClass(touched.email && loginErrors.email)}
                />
              </div>
              {touched.email && loginErrors.email && <p className="text-red-500 text-[11px] font-bold mt-1.5 ml-1">{loginErrors.email}</p>}
            </div>

            <div>
              <div className="flex justify-between items-center mb-1.5">
                <label className="block text-sm font-bold text-slate-700">{authContent.login.passwordLabel} <span className="text-red-500">*</span></label>
                <button type="button" onClick={switchToForgot} className="text-[11px] font-bold text-brand-600 hover:text-brand-800 transition uppercase tracking-wider">
                  {authContent.login.forgotButton}
                </button>
              </div>
              <div className="relative group">
                <div className={`absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none transition-colors ${touched.password && loginErrors.password ? 'text-red-400' : 'text-slate-400 group-focus-within:text-brand-600'}`}>
                  <i className="ph-bold ph-lock-key text-lg" />
                </div>
                <input
                  type={showPassword ? 'text' : 'password'}
                  autoComplete="new-password"
                  value={loginForm.password}
                  onChange={(event) => {
                    setLoginForm((prev) => ({ ...prev, password: event.target.value }))
                    setTouched((prev) => ({ ...prev, password: true }))
                  }}
                  onBlur={() => setTouched((prev) => ({ ...prev, password: true }))}
                  placeholder={authContent.login.passwordPlaceholder}
                  className={fieldClass(touched.password && loginErrors.password).replace('pr-4', 'pr-12')}
                />
                <button type="button" onClick={() => setShowPassword((prev) => !prev)} className="absolute inset-y-0 right-0 pr-3.5 flex items-center text-slate-400 hover:text-brand-600 transition-colors">
                  <i className={`ph-bold text-lg ${showPassword ? 'ph-eye-slash' : 'ph-eye'}`} />
                </button>
              </div>
              {touched.password && loginErrors.password && <p className="text-red-500 text-[11px] font-bold mt-1.5 ml-1">{loginErrors.password}</p>}
            </div>

            {captchaEnabled && (
              <CaptchaWidget onVerify={(token) => setCaptchaToken(token)} hasError={captchaError} />
            )}

            <label className="flex items-center gap-2 text-sm text-slate-600 font-medium">
              <input
                type="checkbox"
                checked={loginForm.rememberMe}
                onChange={(event) => setLoginForm((prev) => ({ ...prev, rememberMe: event.target.checked }))}
                className="w-4 h-4 rounded border-slate-300 text-brand-600 focus:ring-brand-500"
              />
              Ingat saya
            </label>

            <button type="submit" disabled={isSubmitting || cooldown > 0} className="group relative overflow-hidden bg-brand-600 hover:bg-brand-700 text-white w-full py-3.5 rounded-xl text-sm font-bold transition-all duration-300 shadow-sm hover:shadow-md transform hover:-translate-y-0.5 mt-2 flex items-center justify-center gap-2 disabled:opacity-60 disabled:cursor-not-allowed">
              <span>{isSubmitting ? authContent.login.submittingButton : cooldown > 0 ? `Tunggu ${cooldown}s` : authContent.login.submitButton}</span>
              <i className="ph-bold ph-arrow-right transition-transform duration-300 group-hover:translate-x-1" />
            </button>
          </form>
        ) : (
          <form onSubmit={handleForgot} autoComplete="off" className="w-full flex flex-col gap-5">
            <div>
              <label className="block text-sm font-bold text-slate-700 mb-1.5">{authContent.forgot.emailLabel} <span className="text-red-500">*</span></label>
              <div className="relative group">
                <div className={`absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none transition-colors ${touched.forgotEmail && forgotErrors.forgotEmail ? 'text-red-400' : 'text-slate-400 group-focus-within:text-brand-600'}`}>
                  <i className="ph-bold ph-envelope-simple text-lg" />
                </div>
                <input
                  type="email"
                  autoComplete="off"
                  value={forgotEmail}
                  onChange={(event) => {
                    setForgotEmail(event.target.value)
                    setTouched((prev) => ({ ...prev, forgotEmail: true }))
                  }}
                  onBlur={() => setTouched((prev) => ({ ...prev, forgotEmail: true }))}
                  placeholder={authContent.forgot.emailPlaceholder}
                  className={fieldClass(touched.forgotEmail && forgotErrors.forgotEmail)}
                />
              </div>
              {touched.forgotEmail && forgotErrors.forgotEmail && <p className="text-red-500 text-[11px] font-bold mt-1.5 ml-1">{forgotErrors.forgotEmail}</p>}
            </div>

            {captchaEnabled && (
              <CaptchaWidget onVerify={(token) => setForgotCaptchaToken(token)} hasError={forgotCaptchaError} />
            )}

            <button type="submit" disabled={isSubmitting || cooldown > 0} className="group relative overflow-hidden bg-brand-600 hover:bg-brand-700 text-white w-full py-3.5 rounded-xl text-sm font-bold transition-all duration-300 shadow-sm hover:shadow-md transform hover:-translate-y-0.5 mt-2 flex items-center justify-center gap-2 disabled:opacity-60 disabled:cursor-not-allowed">
              <span>{isSubmitting ? authContent.forgot.submittingButton : cooldown > 0 ? `Tunggu ${cooldown}s` : authContent.forgot.submitButton}</span>
              <i className="ph-bold ph-paper-plane-right transition-transform duration-300 group-hover:translate-x-1 group-hover:-translate-y-1" />
            </button>

            <div className="text-center md:text-left mt-2">
              <button type="button" onClick={switchToLogin} className="text-sm font-bold text-slate-500 hover:text-brand-600 transition inline-flex items-center gap-1.5">
                <i className="ph-bold ph-arrow-left" /> {authContent.forgot.backButton}
              </button>
            </div>
          </form>
        )}
      </div>
    </AuthShell>
  )
}
