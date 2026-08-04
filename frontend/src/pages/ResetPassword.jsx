import { useEffect, useMemo, useState } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import AuthShell from '../components/auth/AuthShell'
import { authContent } from '../content/authContent'
import { authService } from '../services/authService'
import { validatePassword } from '../utils/validation'

export default function ResetPassword() {
  const [searchParams] = useSearchParams()
  const token = searchParams.get('token') || ''
  const email = searchParams.get('email') || ''

  const [status, setStatus] = useState('checking')
  const [errorMessage, setErrorMessage] = useState('')
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [showPass, setShowPass] = useState(false)
  const [showConfirm, setShowConfirm] = useState(false)
  const [touched, setTouched] = useState({ password: false, confirmPassword: false })
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [submitError, setSubmitError] = useState('')
  const [cooldown, setCooldown] = useState(0)

  useEffect(() => {
    if (cooldown <= 0) return
    const timer = setInterval(() => setCooldown((s) => s - 1), 1000)
    return () => clearInterval(timer)
  }, [cooldown])

  useEffect(() => {
    async function checkToken() {
      if (!token) {
        setStatus('expired')
        setErrorMessage('Tautan tidak ditemukan. Silakan periksa kembali email Anda.')
        return
      }

      try {
        await authService.validateResetToken(token)
        setStatus('valid')
      } catch (resetError) {
        try {
          await authService.validateActivationToken(token)
          setStatus('valid')
        } catch {
          setStatus('expired')
          setErrorMessage(resetError.message || authContent.reset.expiredMessage)
        }
      }
    }

    checkToken()
  }, [token])

  const errors = useMemo(() => {
    const next = {}
    const passwordError = validatePassword(password, 'Password baru')
    if (passwordError) next.password = passwordError
    if (!confirmPassword) next.confirmPassword = 'Konfirmasi password wajib diisi.'
    if (confirmPassword && confirmPassword !== password) next.confirmPassword = 'Konfirmasi password tidak cocok.'
    return next
  }, [password, confirmPassword])

  async function handleSubmit(event) {
    event.preventDefault()
    setTouched({ password: true, confirmPassword: true })
    setSubmitError('')
    if (Object.keys(errors).length > 0) return

    const payload = {
      token,
      password,
      confirm_password: confirmPassword,
    }

    setIsSubmitting(true)
    try {
      await authService.resetPassword(payload)
      setStatus('success')
    } catch (resetError) {
      try {
        await authService.activateAccount(payload)
        setStatus('success')
      } catch (activationError) {
        const rateLimitErr = activationError.retryAfter > 0 ? activationError : (resetError.retryAfter > 0 ? resetError : null)
        if (rateLimitErr) {
          setCooldown(rateLimitErr.retryAfter)
          setSubmitError('Terlalu banyak percobaan. Silakan tunggu sebelum mencoba lagi.')
        } else {
          setSubmitError(activationError.message || resetError.message || 'Gagal menyimpan password baru.')
        }
      }
    } finally {
      setIsSubmitting(false)
    }
  }

  const inputClass = (hasError, success = false) => {
    if (success) return 'border-emerald-500 bg-emerald-50/20 focus:ring-2 focus:ring-emerald-500/20 text-emerald-900'
    if (hasError) return 'border-red-400 bg-red-50/20 focus:ring-2 focus:ring-red-500/20 text-red-900 placeholder-red-300'
    return 'bg-slate-50 border-slate-200 focus:bg-white focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 text-slate-700 placeholder-slate-400'
  }

  return (
    <AuthShell variant="reset">
      <div className="w-full max-w-md">
        {status === 'checking' && (
          <div className="text-center py-12 flex flex-col items-center justify-center gap-4">
            <div className="w-16 h-16 border-4 border-brand-200 border-t-brand-600 rounded-full animate-spin" />
            <h3 className="font-heading text-xl font-bold text-slate-800">{authContent.reset.checkingTitle}</h3>
            <p className="text-slate-500 text-sm">{authContent.reset.checkingMessage}</p>
          </div>
        )}

        {status === 'expired' && (
          <div className="flex flex-col gap-6">
            <div className="bg-red-50 border border-red-200 rounded-3xl p-8 text-center flex flex-col items-center shadow-sm">
              <div className="w-16 h-16 bg-red-100 text-red-600 rounded-2xl flex items-center justify-center mb-4 shadow-sm">
                <i className="ph-bold ph-clock-afternoon text-3xl" />
              </div>
              <h2 className="font-heading text-2xl font-extrabold text-slate-900 mb-2">{authContent.reset.expiredTitle}</h2>
              <p className="text-slate-600 text-sm leading-relaxed mb-4">{errorMessage || authContent.reset.expiredMessage}</p>
              <div className="w-full bg-white/80 p-3.5 rounded-xl border border-red-100 text-left text-xs text-slate-500 mb-2">
                <span className="font-bold text-red-600">Saran Keamanan:</span> {authContent.reset.expiredTip}
              </div>
            </div>

            <div className="flex flex-col gap-3">
              <Link to="/login?view=forgot" className="w-full py-3.5 bg-brand-600 hover:bg-brand-700 text-white rounded-xl text-sm font-bold text-center transition-all duration-300 shadow-sm hover:shadow flex items-center justify-center gap-2">
                <i className="ph-bold ph-paper-plane-right" />
                <span>{authContent.reset.requestNewLink}</span>
              </Link>
              <Link to="/login" className="w-full py-3 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-xl text-sm font-semibold text-center transition">
                {authContent.reset.backToLogin}
              </Link>
            </div>
          </div>
        )}

        {status === 'valid' && (
          <div className="flex flex-col gap-6">
            <div className="mb-2 text-center md:text-left">
              <div className="inline-flex w-12 h-12 bg-brand-50 text-brand-600 rounded-xl items-center justify-center mb-4 border border-brand-100">
                <i className="ph-bold ph-key text-xl" />
              </div>
              <h2 className="font-heading text-3xl font-extrabold text-slate-900 mb-2 tracking-tight">{authContent.reset.formTitle}</h2>
              <p className="text-slate-500 text-sm font-medium">
                {authContent.reset.formSubtitle} <span className="font-bold text-slate-800">{email || 'Anda'}</span>.
              </p>
            </div>

            {submitError && <div className="rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm font-medium text-red-700">{submitError}</div>}

            <form onSubmit={handleSubmit} className="flex flex-col gap-5">
              <div>
                <label className="block text-sm font-bold text-slate-700 mb-1.5">{authContent.reset.passwordLabel} <span className="text-red-500">*</span></label>
                <div className="relative group">
                  <div className={`absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none transition-colors ${touched.password && errors.password ? 'text-red-400' : 'text-slate-400 group-focus-within:text-brand-600'}`}>
                    <i className="ph-bold ph-lock-key text-lg" />
                  </div>
                  <input
                    type={showPass ? 'text' : 'password'}
                    value={password}
                    onChange={(event) => setPassword(event.target.value)}
                    onBlur={() => setTouched((prev) => ({ ...prev, password: true }))}
                    placeholder={authContent.reset.passwordPlaceholder}
                    className={`w-full pl-11 pr-12 py-3 bg-slate-50 border rounded-xl focus:outline-none focus:bg-white transition-all text-sm font-medium shadow-sm group-hover:border-slate-300 ${inputClass(touched.password && errors.password)}`}
                  />
                  <button type="button" onClick={() => setShowPass((prev) => !prev)} className="absolute inset-y-0 right-0 pr-3.5 flex items-center text-slate-400 hover:text-brand-600 transition-colors">
                    <i className={`ph-bold text-lg ${showPass ? 'ph-eye-slash' : 'ph-eye'}`} />
                  </button>
                </div>
                {touched.password && errors.password && <p className="text-red-500 text-[11px] font-bold mt-1.5 ml-1">{errors.password}</p>}
              </div>

              <div>
                <label className="block text-sm font-bold text-slate-700 mb-1.5">{authContent.reset.confirmPasswordLabel} <span className="text-red-500">*</span></label>
                <div className="relative group">
                  <div className={`absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none transition-colors ${confirmPassword ? (confirmPassword === password ? 'text-emerald-600' : 'text-red-400') : 'text-slate-400 group-focus-within:text-brand-600'}`}>
                    <i className={`ph-bold ${confirmPassword ? (confirmPassword === password ? 'ph-check-circle' : 'ph-x-circle') : 'ph-check-square-offset'}`} />
                  </div>
                  <input
                    type={showConfirm ? 'text' : 'password'}
                    value={confirmPassword}
                    onChange={(event) => setConfirmPassword(event.target.value)}
                    onBlur={() => setTouched((prev) => ({ ...prev, confirmPassword: true }))}
                    placeholder={authContent.reset.confirmPasswordPlaceholder}
                    className={`w-full pl-11 pr-12 py-3 border rounded-xl focus:outline-none transition-all text-sm font-medium shadow-sm ${inputClass(errors.confirmPassword && confirmPassword, confirmPassword && confirmPassword === password)}`}
                  />
                  <button type="button" onClick={() => setShowConfirm((prev) => !prev)} className="absolute inset-y-0 right-0 pr-3.5 flex items-center text-slate-400 hover:text-brand-600 transition-colors">
                    <i className={`ph-bold text-lg ${showConfirm ? 'ph-eye-slash' : 'ph-eye'}`} />
                  </button>
                </div>
                {confirmPassword && confirmPassword === password && <p className="text-emerald-600 text-[11px] font-bold mt-1.5 ml-1 flex items-center gap-1"><i className="ph-bold ph-check-circle text-xs" /> Password cocok</p>}
                {confirmPassword && confirmPassword !== password && errors.confirmPassword && <p className="text-red-500 text-[11px] font-bold mt-1.5 ml-1 flex items-center gap-1"><i className="ph-bold ph-x-circle text-xs" /> {errors.confirmPassword}</p>}
              </div>

              <button type="submit" disabled={isSubmitting || cooldown > 0} className="group relative overflow-hidden bg-brand-600 hover:bg-brand-700 text-white w-full py-3.5 rounded-xl text-sm font-bold transition-all duration-300 shadow-sm hover:shadow-md transform hover:-translate-y-0.5 mt-2 flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed">
                <span>{isSubmitting ? authContent.reset.submittingButton : cooldown > 0 ? `Tunggu ${cooldown}s` : authContent.reset.submitButton}</span>
                <i className="ph-bold ph-floppy-disk transition-transform group-hover:scale-110" />
              </button>
            </form>
          </div>
        )}

        {status === 'success' && (
          <div className="bg-emerald-50 border border-emerald-200 rounded-3xl p-8 text-center flex flex-col items-center shadow-sm">
            <div className="w-16 h-16 bg-emerald-100 text-emerald-600 rounded-2xl flex items-center justify-center mb-4 shadow-sm animate-bounce">
              <i className="ph-bold ph-check-circle text-3xl" />
            </div>
            <h2 className="font-heading text-2xl font-extrabold text-slate-900 mb-2">{authContent.reset.successTitle}</h2>
            <p className="text-slate-600 text-sm leading-relaxed mb-6">{authContent.reset.successMessage}</p>
            <Link to="/login" className="w-full py-3.5 bg-emerald-600 hover:bg-emerald-700 text-white rounded-xl text-sm font-bold text-center transition-all duration-300 shadow-sm hover:shadow flex items-center justify-center gap-2">
              <span>{authContent.reset.loginButton}</span>
              <i className="ph-bold ph-arrow-right" />
            </Link>
          </div>
        )}
      </div>
    </AuthShell>
  )
}
