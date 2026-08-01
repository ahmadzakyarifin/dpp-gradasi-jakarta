import { useState, useEffect } from 'react'
import AdminLayout from '../../layouts/AdminLayout'
import { settingsService } from '../../services/settingsService'
import { useSettings } from '../../context/SettingsContext'
import { resolveAssetUrl } from '../../utils/assetUrl'

export default function SettingsAdmin() {
  const { refresh } = useSettings()
  const [currentTab, setCurrentTab] = useState('profil')
  const [loading, setLoading] = useState(false)
  const [logoUploading, setLogoUploading] = useState(false)
  const [logoPreview, setLogoPreview] = useState(null)
  const [saved, setSaved] = useState(false)

  const [formData, setFormData] = useState({
    site_name: '',
    tagline: '',
    logo_url: '',
    contact_email: '',
    contact_phone: '',
    address: '',
    maps_embed_url: '',
    facebook_url: '',
    instagram_url: '',
    youtube_url: '',
    video_profile_url: '',
    history: '',
    about_tutorial: '',
    about_formation_date: '',
    about_no_sk: '',
    about_vision: '',
    about_mission: '',
    greeting_title: '',
    greeting_subtitle: '',
    greeting_date: '',
    greeting_content: '',
    greeting_image_url: ''
  })

  useEffect(() => {
    settingsService.get()
      .then(res => {
        if (res.success && res.data) {
          const data = res.data
          // convert about_mission from array/JSON to simple newline string for textarea
          let missionText = ''
          try {
            const arr = typeof data.about_mission === 'string' ? JSON.parse(data.about_mission) : data.about_mission
            if (Array.isArray(arr)) {
              missionText = arr.join('\n')
            }
          } catch {
            missionText = data.about_mission || ''
          }
          setFormData({ ...data, about_mission: missionText })
        }
      }).catch(() => {})
  }, [])

  const showToast = (msg, isError = false) => {
    setSaved({ msg, isError })
    setTimeout(() => setSaved(null), 4000)
  }

  const handleLogoChange = (e) => {
    const file = e.target.files?.[0]
    if (!file) return
    // Validasi client-side: ukuran & mime
    if (file.size > 2 * 1024 * 1024) {
      showToast('Ukuran file maksimal 2MB.', true)
      e.target.value = ''
      return
    }
    if (!['image/png', 'image/jpeg', 'image/webp'].includes(file.type)) {
      showToast('Format logo harus PNG, JPG, atau WEBP.', true)
      e.target.value = ''
      return
    }
    setLogoPreview(URL.createObjectURL(file))
    handleLogoUpload(file)
  }

  const handleLogoUpload = async (file) => {
    setLogoUploading(true)
    try {
      const res = await settingsService.uploadLogo(file)
      if (res.success && res.data) {
        setFormData(prev => ({ ...prev, logo_url: res.data.logo_url }))
        showToast('Logo berhasil diunggah!')
        refresh()
      } else {
        showToast(res.message || 'Gagal mengunggah logo.', true)
      }
    } catch (err) {
      showToast(err.message || 'Terjadi kesalahan saat mengunggah logo.', true)
    } finally {
      setLogoUploading(false)
    }
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setLoading(true)

    // Buang field meta yang tidak boleh dikirim (backend: id/updated_by harus string/null,
    // dan itu bukan kontrak update) + konversi about_mission ke array of string
    const { id, created_at, updated_at, updated_by, ...rest } = formData
    const payload = {
      ...rest,
      about_mission: (rest.about_mission || '')
        .split('\n').map(s => s.trim()).filter(Boolean)
    }

    try {
      const res = await settingsService.update(payload)
      if (res.success) {
        showToast('Pengaturan website berhasil disimpan!')
        refresh()
      } else {
        showToast('Gagal menyimpan: ' + res.message, true)
      }
    } catch (err) {
      showToast(err?.message || 'Terjadi kesalahan saat menyimpan pengaturan', true)
    } finally {
      setLoading(false)
    }
  }

  const inputCls = "w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none"

  return (
    <AdminLayout title="Pengaturan Website">
      <div className="max-w-5xl mx-auto grid grid-cols-1 md:grid-cols-4 gap-6">

        {saved && (
          <div className={`md:col-span-4 px-4 py-3 rounded-xl text-sm font-bold flex items-center gap-2 ${saved.isError ? 'bg-rose-50 text-rose-700 border border-rose-200' : 'bg-emerald-50 text-emerald-700 border border-emerald-200'}`}>
            <i className={`ph-bold ${saved.isError ? 'ph-warning-circle' : 'ph-check-circle'} text-lg`} /> {saved.msg}
          </div>
        )}

        {/* Navigation Tabs */}
        <div className="md:col-span-1 bg-white p-4 rounded-2xl border border-gray-200 shadow-sm flex flex-col gap-1">
          <button onClick={() => setCurrentTab('profil')} className={`px-4 py-2 text-left text-sm font-semibold rounded-lg transition ${currentTab === 'profil' ? 'bg-brand-50 text-brand-600' : 'text-slate-600 hover:bg-slate-50'}`}>Profil & Sejarah</button>
          <button onClick={() => setCurrentTab('sambutan')} className={`px-4 py-2 text-left text-sm font-semibold rounded-lg transition ${currentTab === 'sambutan' ? 'bg-brand-50 text-brand-600' : 'text-slate-600 hover:bg-slate-50'}`}>Sambutan Depan</button>
          <button onClick={() => setCurrentTab('kontak')} className={`px-4 py-2 text-left text-sm font-semibold rounded-lg transition ${currentTab === 'kontak' ? 'bg-brand-50 text-brand-600' : 'text-slate-600 hover:bg-slate-50'}`}>Informasi Kontak</button>
          <button onClick={() => setCurrentTab('sosmed')} className={`px-4 py-2 text-left text-sm font-semibold rounded-lg transition ${currentTab === 'sosmed' ? 'bg-brand-50 text-brand-600' : 'text-slate-600 hover:bg-slate-50'}`}>Media Sosial</button>
          <button onClick={() => setCurrentTab('logo')} className={`px-4 py-2 text-left text-sm font-semibold rounded-lg transition ${currentTab === 'logo' ? 'bg-brand-50 text-brand-600' : 'text-slate-600 hover:bg-slate-50'}`}>Logo Website</button>
        </div>

        {/* Form Area */}
        <div className="md:col-span-3 bg-white p-6 rounded-2xl border border-gray-200 shadow-sm">
          <form onSubmit={handleSubmit} className="space-y-6">

            {currentTab === 'logo' && (
              <div className="space-y-4">
                <h3 className="font-heading font-bold text-gray-800 text-base border-b pb-2">Logo Website</h3>
                <p className="text-xs text-gray-500">Logo tampil di header, footer, dan section Visi & Misi halaman publik. Maksimal 2MB dengan format PNG, JPG, atau WEBP.</p>
                <div className="flex items-center gap-6">
                  <div className="w-40 h-40 bg-slate-50 border border-gray-200 rounded-2xl flex items-center justify-center overflow-hidden">
                    <img
                      src={logoPreview || resolveAssetUrl(formData.logo_url) || 'https://via.placeholder.com/160'}
                      alt="Preview Logo"
                      className="w-full h-full object-contain p-2"
                    />
                  </div>
                  <div className="flex flex-col gap-3">
                    <label className="inline-flex items-center gap-2 px-5 py-2.5 bg-brand-600 text-white rounded-lg hover:bg-brand-700 text-sm font-semibold cursor-pointer transition">
                      <i className="ph-bold ph-upload-simple" />
                      {logoUploading ? 'Mengunggah...' : 'Pilih & Upload Logo Baru'}
                      <input type="file" accept="image/png,image/jpeg,image/webp" onChange={handleLogoChange} className="hidden" disabled={logoUploading} />
                    </label>
                    <span className="text-[11px] text-gray-400">PNG / JPG / WEBP · maks 2MB</span>
                    {formData.logo_url && (
                      <span className="text-[11px] text-gray-400 break-all">Path saat ini: {formData.logo_url}</span>
                    )}
                  </div>
                </div>
              </div>
            )}

            {currentTab === 'profil' && (
              <div className="space-y-4">
                <h3 className="font-heading font-bold text-gray-800 text-base border-b pb-2">Profil & Sejarah Organisasi</h3>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Nama Website / Organisasi</label>
                    <input type="text" value={formData.site_name} onChange={e => setFormData({...formData, site_name: e.target.value})} className={inputCls} />
                  </div>
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Tagline</label>
                    <input type="text" value={formData.tagline} onChange={e => setFormData({...formData, tagline: e.target.value})} className={inputCls} />
                  </div>
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Tanggal Terbentuk</label>
                    <input type="text" value={formData.about_formation_date} onChange={e => setFormData({...formData, about_formation_date: e.target.value})} className={inputCls} />
                  </div>
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Nomor SK Legalitas</label>
                    <input type="text" value={formData.about_no_sk} onChange={e => setFormData({...formData, about_no_sk: e.target.value})} className={inputCls} />
                  </div>
                </div>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Visi Utama</label>
                  <textarea rows="2" value={formData.about_vision} onChange={e => setFormData({...formData, about_vision: e.target.value})} className={inputCls} />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Misi Organisasi (Satu Misi Per Baris)</label>
                  <textarea rows="4" value={formData.about_mission} onChange={e => setFormData({...formData, about_mission: e.target.value})} className={inputCls} />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Sejarah / Tentang Kami</label>
                  <textarea rows="4" value={formData.history} onChange={e => setFormData({...formData, history: e.target.value})} className={inputCls} />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Selayang Pandang (Paragraf Kedua)</label>
                  <textarea rows="3" value={formData.about_tutorial} onChange={e => setFormData({...formData, about_tutorial: e.target.value})} className={inputCls} />
                </div>
              </div>
            )}

            {currentTab === 'sambutan' && (
              <div className="space-y-4">
                <h3 className="font-heading font-bold text-gray-800 text-base border-b pb-2">Sambutan Halaman Depan</h3>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Judul Sambutan</label>
                  <input type="text" value={formData.greeting_title} onChange={e => setFormData({...formData, greeting_title: e.target.value})} className={inputCls} />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Sub Judul</label>
                  <input type="text" value={formData.greeting_subtitle} onChange={e => setFormData({...formData, greeting_subtitle: e.target.value})} className={inputCls} />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Tanggal Sambutan</label>
                  <input type="text" value={formData.greeting_date} onChange={e => setFormData({...formData, greeting_date: e.target.value})} className={inputCls} />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Isi Sambutan</label>
                  <textarea rows="5" value={formData.greeting_content} onChange={e => setFormData({...formData, greeting_content: e.target.value})} className={inputCls} />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">URL Gambar / Poster Sambutan</label>
                  <input type="text" value={formData.greeting_image_url} onChange={e => setFormData({...formData, greeting_image_url: e.target.value})} className={inputCls} />
                </div>
              </div>
            )}

            {currentTab === 'kontak' && (
              <div className="space-y-4">
                <h3 className="font-heading font-bold text-gray-800 text-base border-b pb-2">Informasi Kontak & Profil Kantor</h3>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Email Resmi</label>
                    <input type="email" value={formData.contact_email} onChange={e => setFormData({...formData, contact_email: e.target.value})} className={inputCls} />
                  </div>
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Telepon / WhatsApp</label>
                    <input type="text" value={formData.contact_phone} onChange={e => setFormData({...formData, contact_phone: e.target.value})} className={inputCls} />
                  </div>
                </div>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Alamat Kantor</label>
                  <input type="text" value={formData.address} onChange={e => setFormData({...formData, address: e.target.value})} className={inputCls} />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Google Maps Embed URL</label>
                  <input type="text" value={formData.maps_embed_url} onChange={e => setFormData({...formData, maps_embed_url: e.target.value})} className={inputCls} />
                </div>
              </div>
            )}

            {currentTab === 'sosmed' && (
              <div className="space-y-4">
                <h3 className="font-heading font-bold text-gray-800 text-base border-b pb-2">Social Media & Profil Link</h3>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Facebook URL</label>
                  <input type="text" value={formData.facebook_url} onChange={e => setFormData({...formData, facebook_url: e.target.value})} className={inputCls} />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Instagram URL</label>
                  <input type="text" value={formData.instagram_url} onChange={e => setFormData({...formData, instagram_url: e.target.value})} className={inputCls} />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">YouTube Channel URL</label>
                  <input type="text" value={formData.youtube_url} onChange={e => setFormData({...formData, youtube_url: e.target.value})} className={inputCls} />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Video Profile URL (Youtube Embed)</label>
                  <input type="text" value={formData.video_profile_url} onChange={e => setFormData({...formData, video_profile_url: e.target.value})} className={inputCls} />
                </div>
              </div>
            )}

            <div className="flex justify-end pt-4 border-t border-slate-100">
              <button type="submit" disabled={loading} className="px-6 py-2.5 bg-brand-600 text-white rounded-lg hover:bg-brand-700 text-sm font-semibold">
                {loading ? 'Menyimpan...' : 'Simpan Semua Pengaturan'}
              </button>
            </div>
          </form>
        </div>

      </div>
    </AdminLayout>
  )
}
