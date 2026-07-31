import { useState, useEffect } from 'react'
import AdminLayout from '../../layouts/AdminLayout'
import { useAuthStore } from '../../store/useAuthStore'

export default function ProfileAdmin() {
  const { token, fetchMe } = useAuthStore()
  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [photo, setPhoto] = useState(null)
  const [photoPreview, setPhotoPreview] = useState('')
  
  const [oldPassword, setOldPassword] = useState('')
  const [newPassword, setNewPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [loadingPwd, setLoadingPwd] = useState(false)

  useEffect(() => {
    fetchMe().then(user => {
      if (user) {
        setName(user.name)
        setEmail(user.email)
        setPhotoPreview(user.photo_path ? `http://127.0.0.1:8080${user.photo_path}` : '')
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
    const fd = new FormData()
    fd.append('name', name)
    fd.append('email', email)
    if (photo) fd.append('photo', photo)

    try {
      const response = await fetch('http://127.0.0.1:8080/api/v1/profile', {
        method: 'PUT',
        headers: { Authorization: `Bearer ${token}` },
        body: fd
      })
      const res = await response.json()
      if (res.success) {
        alert(res.message || 'Profil berhasil diperbarui')
        fetchMe()
      } else {
        alert('Gagal: ' + res.message)
      }
    } catch {
      alert('Terjadi kesalahan')
    } finally {
      setLoading(false)
    }
  }

  const handleUpdatePassword = async (e) => {
    e.preventDefault()
    setLoadingPwd(true)
    try {
      const response = await fetch('http://127.0.0.1:8080/api/v1/profile/password', {
        method: 'PUT',
        headers: {
          Authorization: `Bearer ${token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ old_password: oldPassword, new_password: newPassword })
      })
      const res = await response.json()
      if (res.success) {
        alert('Password berhasil diubah')
        setOldPassword('')
        setNewPassword('')
      } else {
        alert('Gagal: ' + res.message)
      }
    } catch {
      alert('Terjadi kesalahan')
    } finally {
      setLoadingPwd(false)
    }
  }

  return (
    <AdminLayout title="Profil Saya">
      <div className="max-w-4xl mx-auto space-y-6">
        
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
