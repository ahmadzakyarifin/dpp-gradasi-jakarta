import { useState, useEffect } from 'react'
import { useSearchParams } from 'react-router-dom'
import AdminLayout from '../../layouts/AdminLayout'
import { useAuthStore } from '../../store/useAuthStore'
import { apiRequest } from '../../api'

const apiBase = import.meta.env.VITE_API_URL || 'http://127.0.0.1:8080/api/v1'

export default function ProfileAdmin() {
  const { fetchMe } = useAuthStore()
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

  useEffect(() => {
    fetchMe().then(user => {
      if (user) {
        setName(user.name)
        setEmail(user.email)
        setPhotoPreview(user.photo_path ? `${apiBase.replace(/\/api\/v1$/, '')}${user.photo_path}` : '')
      }
    })
  }, [fetchMe])

  const handleFileChange = (e) => {
    const file = e.target.files[0]
    if (file) {
      setPhoto(file)
      setPhotoPreview(URL.createObjectURL(file))
    }
  }

  const handleUpdateProfile = async (e) => {
    e.preventDefault()
    setLoading(true)
    setMessage(null)
    const fd = new FormData()
    fd.append('name', name)
    fd.append('email', email)
    if (photo) fd.append('photo', photo)

    try {
      const res = await apiRequest('/profile', {
        method: 'PUT',
        body: fd,
        headers: {}, // apiRequest handles Authorization; FormData sets boundary
      })
      setMessage({ type: 'success', text: res.message || 'Profil berhasil diperbarui' })
      if (res.data?.email && res.data.email !== email) {
        // email berubah → perlu verifikasi
        setMessage({ type: 'success', text: res.message })
      }
      fetchMe()
    } catch (err) {
      setMessage({ type: 'error', text: err.message || 'Terjadi kesalahan' })
    } finally {
      setLoading(false)
    }
  }

  const handleUpdatePassword = async (e) => {
    e.preventDefault()
    setLoadingPwd(true)
    setMessage(null)
    try {
      const res = await apiRequest('/profile/password', {
        method: 'PUT',
        body: JSON.stringify({ old_password: oldPassword, new_password: newPassword }),
      })
      setMessage({ type: 'success', text: res.message || 'Password berhasil diubah' })
      setOldPassword('')
      setNewPassword('')
    } catch (err) {
      setMessage({ type: 'error', text: err.message || 'Terjadi kesalahan' })
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
            <p className="text-sm text-gray-500">Perbarui nama, email, dan foto profil Anda.</p>
          </div>
          <div className="p-6">
            <form onSubmit={handleUpdateProfile}>
              <div className="flex items-center gap-6 mb-6">
                <img src={photoPreview || 'https://ui-avatars.com/api/?name=Admin&background=0D8ABC&color=fff'} className="w-20 h-20 rounded-full object-cover border border-gray-200" alt="Foto" />
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">Foto Profil</label>
                  <input type="file" onChange={handleFileChange} accept="image/*" className="text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-brand-50 file:text-brand-700 hover:file:bg-brand-100 cursor-pointer" />
                </div>
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Nama Lengkap</label>
                  <input type="text" value={name} onChange={(e) => setName(e.target.value)} required className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-brand-500 focus:border-brand-500 text-sm outline-none" />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Alamat Email</label>
                  <input type="email" value={email} onChange={(e) => setEmail(e.target.value)} required className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-brand-500 focus:border-brand-500 text-sm outline-none" />
                </div>
              </div>
              <div className="flex justify-end">
                <button type="submit" className="px-4 py-2 bg-brand-600 text-white rounded-lg hover:bg-brand-700 transition shadow-sm text-sm font-medium" disabled={loading}>
                  {loading ? 'Menyimpan...' : 'Simpan Perubahan'}
                </button>
              </div>
            </form>
          </div>
        </div>

        {/* Change Password Card */}
        <div className="bg-white rounded-2xl shadow-sm border border-gray-200 overflow-hidden">
          <div className="p-6 border-b border-gray-200">
            <h2 className="text-lg font-semibold text-gray-800">Ubah Password</h2>
            <p className="text-sm text-gray-500">Pastikan akun Anda aman dengan password yang kuat.</p>
          </div>
          <div className="p-6">
            <form onSubmit={handleUpdatePassword}>
              <div className="space-y-4 mb-6">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Password Lama</label>
                  <input type="password" value={oldPassword} onChange={(e) => setOldPassword(e.target.value)} required className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-brand-500 focus:border-brand-500 text-sm outline-none" />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Password Baru</label>
                  <input type="password" value={newPassword} onChange={(e) => setNewPassword(e.target.value)} required minLength={6} className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-brand-500 focus:border-brand-500 text-sm outline-none" />
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
