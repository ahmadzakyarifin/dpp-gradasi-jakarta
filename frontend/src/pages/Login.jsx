import { useMemo, useState } from 'react'
import { Link, useNavigate, useSearchParams } from 'react-router-dom'
import AuthShell from '../components/auth/AuthShell'
import CaptchaWidget from '../components/auth/CaptchaWidget'
import { authContent } from '../content/authContent'
import { authService } from '../services/authService'
import { useAuthStore } from '../store/useAuthStore'
import { validateEmail } from '../utils/validation'

export default function Login() {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const login = useAuthStore((state) => state.login)

  const captchaEnabled = import.meta.env.VITE_CAPTCHA_ENABLED === 'true' || import.meta.env.VITE_CAPTCHA_ENABLED === true

  const initialView = searchParams.get('view') === 'forgot' ? 'forgot' : 'login'
  const [view, setView] = useState(initialView)
  const [loginForm, setLoginForm] = useState({ email: '', password: '', rememberMe: false })
  const [forgotEmail, setForgotEmail] = useState('')
  const [captchaToken, setCaptchaToken] = useState('')
  const [captchaError, setCaptchaError] = useState(false)
  const [showPassword, setShowPassword] = useState(false)
  const [touched, setTouched] = useState({ email: false, password: false, forgotEmail: false })
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [notice, setNotice] = useState({ type: '', message: '' })

  const content = view === 'forgot' ? authContent.forgot : authContent.login

  const loginErrors = useMemo(() => {
    const next = {}
    const emailError = validateEmail(loginForm.email)
    if (emailError) next.email = emailError
    if (!loginForm.password) next.password = 'Password wajib diisi'
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
  }

  function switchToLogin() {
    setView('login')
    setTouched({ email: false, password: false, forgotEmail: false })
    setNotice({ type: '', message: '' })
  }

  async function handleLogin(event) {
    event.preventDefault()
    setTouched((prev) => ({ ...prev, email: true, password: true }))
    setNotice({ type: '', message: '' })
    setCaptchaError(false)

    if (Object.keys(loginErrors).length > 0) return

    if (captchaEnabled && !captchaToken) {
      setCaptchaError(true)
      setNotice({ type: 'error', message: 'Silakan selesaikan kode CAPTCHA dengan benar.' })
      return
    }

    setIsSubmitting(true)
    try {
      await login({
        ...loginForm,
        captcha_token: captchaToken
      })
      navigate(authContent.adminPath)
    } catch (error) {
      setNotice({ type: 'error', message: error.message || 'Login gagal' })
    } finally {
      setIsSubmitting(false)
    }
  }

  async function handleForgot(event) {
    event.preventDefault()
    setTouched((prev) => ({ ...prev, forgotEmail: true }))
    setNotice({ type: '', message: '' })
    if (Object.keys(forgotErrors).length > 0) return

    setIsSubmitting(true)
    try {
      const response = await authService.forgotPassword(forgotEmail)
      setNotice({
        type: 'success',
        message: response?.message || authContent.forgot.successMessage,
      })
    } catch (error) {
      setNotice({ type: 'error', message: error.message || 'Gagal mengirim link reset password' })
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
            <img src={authContent.logoUrl} alt={authContent.brandName} className="w-full h-full object-contain" />
          </div>
          <h2 className="font-heading text-3xl font-extrabold text-slate-900 mb-2 tracking-tight">{content.title}</h2>
          <p className="text-slate-500 text-sm font-medium">{content.subtitle}</p>
        </div>

        {notice.message && (
          <div className={`mb-5 rounded-xl border px-4 py-3 text-sm font-medium ${notice.type === 'success' ? 'bg-emerald-50 text-emerald-700 border-emerald-200' : 'bg-red-50 text-red-700 border-red-200'}`}>
            {notice.message}
          </div>
        )}

        {view === 'login' ? (
          <form onSubmit={handleLogin} className="w-full flex flex-col gap-5">
            <div>
              <label className="block text-sm font-bold text-slate-700 mb-1.5">{authContent.login.emailLabel}</label>
              <div className="relative group">
                <div className={`absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none transition-colors ${touched.email && loginErrors.email ? 'text-red-400' : 'text-slate-400 group-focus-within:text-brand-600'}`}>
                  <i className="ph-bold ph-envelope-simple text-lg" />
                </div>
                <input
                  type="email"
                  value={loginForm.email}
                  onChange={(event) => setLoginForm((prev) => ({ ...prev, email: event.target.value }))}
                  onBlur={() => setTouched((prev) => ({ ...prev, email: true }))}
                  placeholder={authContent.login.emailPlaceholder}
                  className={fieldClass(touched.email && loginErrors.email)}
                />
              </div>
              {touched.email && loginErrors.email && <p className="text-red-500 text-[11px] font-bold mt-1.5 ml-1">{loginErrors.email}</p>}
            </div>

            <div>
              <div className="flex justify-between items-center mb-1.5">
                <label className="block text-sm font-bold text-slate-700">{authContent.login.passwordLabel}</label>
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
                  value={loginForm.password}
                  onChange={(event) => setLoginForm((prev) => ({ ...prev, password: event.target.value }))}
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

            <button type="submit" disabled={isSubmitting} className="group relative overflow-hidden bg-brand-600 hover:bg-brand-700 text-white w-full py-3.5 rounded-xl text-sm font-bold transition-all duration-300 shadow-sm hover:shadow-md transform hover:-translate-y-0.5 mt-2 flex items-center justify-center gap-2 disabled:opacity-60 disabled:cursor-not-allowed">
              <span>{isSubmitting ? authContent.login.submittingButton : authContent.login.submitButton}</span>
              <i className="ph-bold ph-arrow-right transition-transform duration-300 group-hover:translate-x-1" />
            </button>
          </form>
        ) : (
          <form onSubmit={handleForgot} className="w-full flex flex-col gap-5">
            <div>
              <label className="block text-sm font-bold text-slate-700 mb-1.5">{authContent.forgot.emailLabel}</label>
              <div className="relative group">
                <div className={`absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none transition-colors ${touched.forgotEmail && forgotErrors.forgotEmail ? 'text-red-400' : 'text-slate-400 group-focus-within:text-brand-600'}`}>
                  <i className="ph-bold ph-envelope-simple text-lg" />
                </div>
                <input
                  type="email"
                  value={forgotEmail}
                  onChange={(event) => setForgotEmail(event.target.value)}
                  onBlur={() => setTouched((prev) => ({ ...prev, forgotEmail: true }))}
                  placeholder={authContent.forgot.emailPlaceholder}
                  className={fieldClass(touched.forgotEmail && forgotErrors.forgotEmail)}
                />
              </div>
              {touched.forgotEmail && forgotErrors.forgotEmail && <p className="text-red-500 text-[11px] font-bold mt-1.5 ml-1">{forgotErrors.forgotEmail}</p>}
            </div>

            <button type="submit" disabled={isSubmitting} className="group relative overflow-hidden bg-brand-600 hover:bg-brand-700 text-white w-full py-3.5 rounded-xl text-sm font-bold transition-all duration-300 shadow-sm hover:shadow-md transform hover:-translate-y-0.5 mt-2 flex items-center justify-center gap-2 disabled:opacity-60 disabled:cursor-not-allowed">
              <span>{isSubmitting ? authContent.forgot.submittingButton : authContent.forgot.submitButton}</span>
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
