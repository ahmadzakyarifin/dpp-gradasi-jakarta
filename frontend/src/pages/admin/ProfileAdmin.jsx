import { useState, useEffect } from 'react'
import { useSearchParams } from 'react-router-dom'
import AdminLayout from '../../layouts/AdminLayout'
import { useAuthStore } from '../../store/useAuthStore'
import { apiRequest } from '../../api'

const apiBase = import.meta.env.VITE_API_URL || 'http://127.0.0.1:8080/api/v1'

export default function ProfileAdmin() {
  const { fetchMe, user } = useAuthStore()
  const [searchParams] = useSearchParams()
  const forced = searchParams.get('force') === '1'
  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [photo, setPhoto] = useState(null)
  const [photoPreview, setPhotoPreview] = useState('')

  const [oldPassword, setOldPassword] = useState('')
  const [newPassword, setNewPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [loadingPwd, setLoadingPwd] = useState(false)
  const [message, setMessage] = useState(null)

  const [profileErrors, setProfileErrors] = useState({})
  const [profileTouched, setProfileTouched] = useState({})
  const [pwdErrors, setPwdErrors] = useState({})
  const [pwdTouched, setPwdTouched] = useState({})

  const [verificationToken, setVerificationToken] = useState('')
  const [verifying, setVerifying] = useState(false)
  const [verifyError, setVerifyError] = useState('')

  const validateProfile = (nameVal = name, emailVal = email) => {
    const errors = {}
    if (!nameVal || !nameVal.trim()) {
      errors.name = 'Nama lengkap wajib diisi.'
    }
    if (!emailVal || !emailVal.trim()) {
      errors.email = 'Alamat email wajib diisi.'
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(emailVal)) {
      errors.email = 'Format email tidak valid.'
    }
    return errors
  }

  const validatePwd = (oldVal = oldPassword, newVal = newPassword) => {
    const errors = {}
    if (!oldVal) {
      errors.oldPassword = 'Password lama wajib diisi.'
    }
    if (!newVal) {
      errors.newPassword = 'Password baru wajib diisi.'
    } else if (newVal.length < 6) {
      errors.newPassword = 'Password baru minimal 6 karakter.'
    }
    return errors
  }

  useEffect(() => {
    fetchMe().then(user => {
      if (user) {
        setName(user.name)
        setEmail(user.email)
        setPhotoPreview(user.photo_path ? `${apiBase.replace(/\/api\/v1$/, '')}${user.photo_path}` : '')
      }
    })
  }, [fetchMe])

  const handleUpdateProfile = async (e) => {
    e.preventDefault()
    const errors = validateProfile()
    if (Object.keys(errors).length > 0) {
      setProfileErrors(errors)
      setProfileTouched({ name: true, email: true })
      return
    }
    setProfileErrors({})

    setLoading(true)
    setMessage(null)

    try {
      const res = await apiRequest('/profile', {
        method: 'PUT',
        body: JSON.stringify({ name, email }),
        headers: {
          'Content-Type': 'application/json'
        }
      })
      setProfileTouched({})
      setProfileErrors({})
      let successMsg = res.message || 'Profil berhasil diperbarui.'
      if (res?.data?.user?.email_pending) {
        successMsg = `Perubahan email berhasil disimpan. Silakan periksa email ${res.data.user.email_pending} untuk kode verifikasi (OTP).`
      }
      setMessage({ type: 'success', text: successMsg })
      fetchMe()
    } catch (err) {
      if (err?.data?.errors) {
        const parsedErrors = {}
        err.data.errors.forEach(x => {
          parsedErrors[x.field] = x.message
        })
        setProfileErrors(parsedErrors)
        setProfileTouched(Object.keys(parsedErrors).reduce((acc, k) => ({ ...acc, [k]: true }), {}))
      } else {
        setMessage({ type: 'error', text: err.message || 'Terjadi kesalahan' })
      }
    } finally {
      setLoading(false)
    }
  }

  const handleVerifyEmail = async (e) => {
    e.preventDefault()
    if (!verificationToken.trim()) {
      setVerifyError('Kode OTP verifikasi wajib diisi.')
      return
    }
    setVerifyError('')
    setVerifying(true)
    try {
      const res = await apiRequest('/profile/verify-email', {
        method: 'POST',
        body: JSON.stringify({ token: verificationToken.trim() }),
      })
      setMessage({ type: 'success', text: res.message || 'Email berhasil diverifikasi!' })
      setVerificationToken('')
      fetchMe()
    } catch (err) {
      setVerifyError(err.message || 'Gagal memverifikasi email.')
    } finally {
      setVerifying(false)
    }
  }

  const handleUpdatePassword = async (e) => {
    e.preventDefault()
    const errors = validatePwd()
    if (Object.keys(errors).length > 0) {
      setPwdErrors(errors)
      setPwdTouched({ oldPassword: true, newPassword: true })
      return
    }
    setPwdErrors({})

    setLoadingPwd(true)
    setMessage(null)
    try {
      const res = await apiRequest('/profile/password', {
        method: 'PUT',
        body: JSON.stringify({ old_password: oldPassword, new_password: newPassword }),
      })
      setPwdTouched({})
      setPwdErrors({})
      setMessage({ type: 'success', text: res.message || 'Password berhasil diubah' })
      setOldPassword('')
      setNewPassword('')
    } catch (err) {
      if (err?.data?.errors) {
        const parsedErrors = {}
        err.data.errors.forEach(x => {
          const key = x.field === 'old_password' ? 'oldPassword' : 'newPassword'
          parsedErrors[key] = x.message
        })
        setPwdErrors(parsedErrors)
        setPwdTouched(Object.keys(parsedErrors).reduce((acc, k) => ({ ...acc, [k]: true }), {}))
      } else {
        setMessage({ type: 'error', text: err.message || 'Terjadi kesalahan' })
      }
    } finally {
      setLoadingPwd(false)
    }
  }

  return (
    <AdminLayout title="Profil Saya">
      <div className="max-w-4xl mx-auto space-y-6">

        {forced && (
          <div className="rounded-xl border-2 border-amber-300 bg-amber-50 px-4 py-4 text-sm font-medium text-amber-800 flex items-start gap-3">
            <i className="ph ph-lock-key text-xl mt-0.5" />
            <div>
              <p className="font-bold">Anda menggunakan password default dari email.</p>
              <p className="text-amber-700 mt-0.5">Demi keamanan, silakan ganti password Anda di bawah ini sebelum melanjutkan.</p>
            </div>
          </div>
        )}

        {message && (
          <div className={`rounded-xl border px-4 py-3 text-sm font-medium ${message.type === 'success' ? 'bg-emerald-50 text-emerald-700 border-emerald-200' : 'bg-red-50 text-red-700 border-red-200'}`}>
            {message.text}
          </div>
        )}

        {/* Profile Info Card */}
        <div className="bg-white rounded-2xl shadow-sm border border-gray-200 overflow-hidden">
          <div className="p-6 border-b border-gray-200">
            <h2 className="text-lg font-semibold text-gray-800">Informasi Pribadi</h2>
            <p className="text-sm text-gray-500">Detail informasi profil Anda.</p>
          </div>
          <div className="p-6">
            <form onSubmit={handleUpdateProfile}>
              <div className="flex items-center gap-6 mb-6">
                <img src={photoPreview || `https://ui-avatars.com/api/?name=${encodeURIComponent(name || 'User')}&background=0D8ABC&color=fff`} className="w-20 h-20 rounded-full object-cover border border-gray-200 shadow-sm" alt="Foto" />
                <div>
                  <h3 className="text-sm font-semibold text-gray-800">{name || 'Nama Pengguna'}</h3>
                  <p className="text-xs text-gray-500">{user?.email}</p>
                </div>
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Nama Lengkap <span className="text-red-500">*</span></label>
                  <input
                    type="text"
                    value={name}
                    onChange={(e) => {
                      setName(e.target.value)
                      if (profileTouched.name) {
                        const errs = validateProfile(e.target.value)
                        setProfileErrors(prev => ({ ...prev, name: errs.name }))
                      }
                    }}
                    onBlur={() => {
                      setProfileTouched(prev => ({ ...prev, name: true }))
                      const errs = validateProfile()
                      setProfileErrors(prev => ({ ...prev, name: errs.name }))
                    }}
                    className={`w-full px-3 py-2 border rounded-lg focus:ring-brand-500 focus:border-brand-500 text-sm outline-none transition-colors ${profileTouched.name && profileErrors.name ? 'border-red-400 focus:ring-2 focus:ring-red-100' : 'border-gray-300'}`}
                  />
                  {profileTouched.name && profileErrors.name && (
                    <p className="text-red-500 text-[11px] font-semibold mt-1 flex items-center gap-1">
                      <i className="ph-bold ph-warning-circle text-xs" /> {profileErrors.name}
                    </p>
                  )}
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Alamat Email <span className="text-red-500">*</span></label>
                  <input
                    type="email"
                    value={email}
                    onChange={(e) => {
                      setEmail(e.target.value)
                      if (profileTouched.email) {
                        const errs = validateProfile(name, e.target.value)
                        setProfileErrors(prev => ({ ...prev, email: errs.email }))
                      }
                    }}
                    onBlur={() => {
                      setProfileTouched(prev => ({ ...prev, email: true }))
                      const errs = validateProfile()
                      setProfileErrors(prev => ({ ...prev, email: errs.email }))
                    }}
                    className={`w-full px-3 py-2 border rounded-lg focus:ring-brand-500 focus:border-brand-500 text-sm outline-none transition-colors ${profileTouched.email && profileErrors.email ? 'border-red-400 focus:ring-2 focus:ring-red-100' : 'border-gray-300'}`}
                  />
                  {profileTouched.email && profileErrors.email && (
                    <p className="text-red-500 text-[11px] font-semibold mt-1 flex items-center gap-1">
                      <i className="ph-bold ph-warning-circle text-xs" /> {profileErrors.email}
                    </p>
                  )}
                  {user && email === user.email && user.email_verified_at ? (
                    <p className="text-emerald-600 text-[11px] font-semibold mt-1 flex items-center gap-1">
                      <i className="ph-bold ph-check-circle text-xs" /> Terverifikasi
                    </p>
                  ) : email !== user?.email ? (
                    <p className="text-amber-600 text-[11px] font-semibold mt-1 flex items-center gap-1">
                      <i className="ph-bold ph-warning-circle text-xs" /> Belum Disimpan (Klik Simpan Perubahan)
                    </p>
                  ) : (
                    <p className="text-red-500 text-[11px] font-semibold mt-1 flex items-center gap-1">
                      <i className="ph-bold ph-warning-circle text-xs" /> Belum Terverifikasi
                    </p>
                  )}
                </div>
              </div>
              <div className="flex justify-end">
                <button type="submit" className="px-4 py-2 bg-brand-600 text-white rounded-lg hover:bg-brand-700 transition shadow-sm text-sm font-medium" disabled={loading}>
                  {loading ? 'Menyimpan...' : 'Simpan Perubahan'}
                </button>
              </div>
            </form>

            {user && user.email_pending && (
              <div className="mt-6 p-4 rounded-xl border border-amber-200 bg-amber-50/50">
                <h4 className="text-sm font-bold text-amber-800 flex items-center gap-1.5 mb-1">
                  <i className="ph-bold ph-warning-circle text-base" /> Verifikasi Email Baru
                </h4>
                <p className="text-xs text-amber-700 leading-relaxed mb-3">
                  Kami telah mengirimkan kode verifikasi (OTP) ke <strong>{user.email_pending}</strong>. 
                  Silakan masukkan kode OTP yang Anda terima di bawah ini.
                </p>
                <form onSubmit={handleVerifyEmail} className="flex gap-3 max-w-md">
                  <div className="flex-1">
                    <input
                      type="text"
                      placeholder="Masukkan kode OTP..."
                      value={verificationToken}
                      onChange={(e) => setVerificationToken(e.target.value)}
                      className="w-full px-3 py-1.5 border border-amber-300 rounded-lg text-sm bg-white outline-none focus:ring-2 focus:ring-amber-500/20 focus:border-amber-500"
                    />
                    {verifyError && <p className="text-red-500 text-[10px] font-semibold mt-1">{verifyError}</p>}
                  </div>
                  <button
                    type="submit"
                    disabled={verifying}
                    className="px-4 py-1.5 bg-amber-600 hover:bg-amber-700 text-white rounded-lg text-xs font-semibold transition"
                  >
                    {verifying ? 'Memverifikasi...' : 'Verifikasi'}
                  </button>
                </form>
              </div>
            )}
          </div>
        </div>

        {/* Change Password Card */}
        <div className="bg-white rounded-2xl shadow-sm border border-gray-200 overflow-hidden">
          <div className="p-6 border-b border-gray-200">
            <h2 className="text-lg font-semibold text-gray-800">Ubah Password</h2>
            <p className="text-sm text-gray-500">Pastikan akun Anda aman dengan password yang kuat.</p>
          </div>
          <div className="p-6">
            <form onSubmit={handleUpdatePassword} noValidate>
              <div className="space-y-4 mb-6">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Password Lama <span className="text-red-500">*</span></label>
                  <input
                    type="password"
                    value={oldPassword}
                    onChange={(e) => {
                      setOldPassword(e.target.value)
                      if (pwdTouched.oldPassword) {
                        const errs = validatePwd(e.target.value, newPassword)
                        setPwdErrors(prev => ({ ...prev, oldPassword: errs.oldPassword }))
                      }
                    }}
                    onBlur={() => {
                      setPwdTouched(prev => ({ ...prev, oldPassword: true }))
                      const errs = validatePwd()
                      setPwdErrors(prev => ({ ...prev, oldPassword: errs.oldPassword }))
                    }}
                    className={`w-full px-3 py-2 border rounded-lg focus:ring-brand-500 focus:border-brand-500 text-sm outline-none transition-colors ${pwdTouched.oldPassword && pwdErrors.oldPassword ? 'border-red-400 focus:ring-2 focus:ring-red-100' : 'border-gray-300'}`}
                  />
                  {pwdTouched.oldPassword && pwdErrors.oldPassword && (
                    <p className="text-red-500 text-[11px] font-semibold mt-1 flex items-center gap-1">
                      <i className="ph-bold ph-warning-circle text-xs" /> {pwdErrors.oldPassword}
                    </p>
                  )}
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Password Baru <span className="text-red-500">*</span></label>
                  <input
                    type="password"
                    value={newPassword}
                    onChange={(e) => {
                      setNewPassword(e.target.value)
                      if (pwdTouched.newPassword) {
                        const errs = validatePwd(oldPassword, e.target.value)
                        setPwdErrors(prev => ({ ...prev, newPassword: errs.newPassword }))
                      }
                    }}
                    onBlur={() => {
                      setPwdTouched(prev => ({ ...prev, newPassword: true }))
                      const errs = validatePwd()
                      setPwdErrors(prev => ({ ...prev, newPassword: errs.newPassword }))
                    }}
                    className={`w-full px-3 py-2 border rounded-lg focus:ring-brand-500 focus:border-brand-500 text-sm outline-none transition-colors ${pwdTouched.newPassword && pwdErrors.newPassword ? 'border-red-400 focus:ring-2 focus:ring-red-100' : 'border-gray-300'}`}
                  />
                  {pwdTouched.newPassword && pwdErrors.newPassword && (
                    <p className="text-red-500 text-[11px] font-semibold mt-1 flex items-center gap-1">
                      <i className="ph-bold ph-warning-circle text-xs" /> {pwdErrors.newPassword}
                    </p>
                  )}
                </div>
              </div>
              <div className="flex justify-end">
                <button type="submit" className="px-4 py-2 bg-brand-600 text-white rounded-lg hover:bg-brand-700 transition shadow-sm text-sm font-medium" disabled={loadingPwd}>
                  {loadingPwd ? 'Memproses...' : 'Ubah Password'}
                </button>
              </div>
            </form>
          </div>
        </div>

      </div>
    </AdminLayout>
  )
}
