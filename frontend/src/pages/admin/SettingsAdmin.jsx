import { useState, useEffect } from 'react'
import AdminLayout from '../../layouts/AdminLayout'
import { settingsService } from '../../services/settingsService'
import { useSettings } from '../../context/useSettings'
import { resolveAssetUrl } from '../../utils/assetUrl'
import { useFormErrors, useRateLimitCooldown } from '../../utils/parseApiError'

export default function SettingsAdmin() {
  const { refresh } = useSettings()
  const [currentTab, setCurrentTab] = useState('profil')
  const [loading, setLoading] = useState(false)
  const [formErrors, setFormErrors] = useState({})
  const [touched, setTouched] = useState({})

  const [formData, setFormData] = useState({
    site_name: '',
    tagline: '',
    logo_path: '',
    contact_email: '',
    contact_phone: '',
    address: '',
    maps_embed_url: '',
    facebook_url: '',
    instagram_url: '',
    youtube_url: '',
    video_profile_path: '',
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
    greeting_image_path: '',
    greeting_sign_name: '',
    greeting_sign_subtitle: '',
    log_retention_days: 30
  })

  const validateForm = (data = formData) => {
    const errors = {}
    if (!data.site_name || !data.site_name.trim()) {
      errors.site_name = 'Nama website / organisasi wajib diisi.'
    }
    if (data.contact_email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(data.contact_email)) {
      errors.contact_email = 'Format email resmi tidak valid.'
    }
    return errors
  }

  // Error backend: field errors inline + countdown rate limit
  const { fieldErrors, applyError, clearFieldError, resetFieldErrors } = useFormErrors()
  const { cooldown, isLimited, applyRateLimit } = useRateLimitCooldown()
  const [logoUploading, setLogoUploading] = useState(false)
  const [logoPreview, setLogoPreview] = useState(null)
  const [saved, setSaved] = useState(false)

  useEffect(() => {
    settingsService.get()
      .then(res => {
        if (res.success && res.data) {
          const data = res.data
          let missionArray = []
          try {
            const arr = typeof data.about_mission === 'string' ? JSON.parse(data.about_mission) : data.about_mission
            if (Array.isArray(arr)) {
              missionArray = arr.map(s => s || '')
            }
          } catch {
            missionArray = data.about_mission ? [data.about_mission] : []
          }
          if (missionArray.length === 0) {
            missionArray = ['']
          }
          setFormData({ ...data, about_mission: missionArray })
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
        setFormData(prev => ({ ...prev, logo_path: res.data.logo_path }))
        showToast('Logo berhasil diunggah!')
        refresh()
      } else {
        showToast(res.message || 'Gagal mengunggah logo.', true)
      }
    } catch (err) {
      const parsed = applyError(err)
      applyRateLimit(err)
      if (Object.keys(parsed.fieldErrors).length === 0) {
        showToast(parsed.message || 'Terjadi kesalahan saat mengunggah logo.', true)
      }
    } finally {
      setLogoUploading(false)
    }
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    const errors = validateForm()
    if (Object.keys(errors).length > 0) {
      setFormErrors(errors)
      setTouched(Object.keys(errors).reduce((acc, k) => ({ ...acc, [k]: true }), {}))
      // Pindahkan tab ke tab yang bermasalah agar user bisa melihat error
      if (errors.site_name) {
        setCurrentTab('profil')
      } else if (errors.contact_email) {
        setCurrentTab('kontak')
      }
      return
    }
    setFormErrors({})
    resetFieldErrors()

    setLoading(true)

    // Buang field meta yang tidak boleh dikirim (backend: id/updated_by harus string/null,
    // dan itu bukan kontrak update) + field hasil komputasi backend (captcha_*, bukan
    // kolom settings) + konversi about_mission ke array of string
    const { id: _id, created_at: _ca, updated_at: _ua, updated_by: _ub, captcha_enabled: _ce, captcha_site_key: _csk, ...rest } = formData
    const payload = {
      ...rest,
      about_mission: (rest.about_mission || []).map(s => s.trim()).filter(Boolean)
    }

    try {
      const res = await settingsService.update(payload)
      if (res.success) {
        setTouched({})
        setFormErrors({})
        showToast('Pengaturan website berhasil disimpan!')
        refresh()
      } else {
        showToast('Gagal menyimpan: ' + res.message, true)
      }
    } catch (err) {
      const parsed = applyError(err)
      applyRateLimit(err)
      setFormErrors(prev => ({ ...prev, ...parsed.fieldErrors }))
      setTouched(prev => ({ ...prev, ...Object.keys(parsed.fieldErrors).reduce((acc, k) => ({ ...acc, [k]: true }), {}) }))
      // Redirect tab if error is in a field from a specific tab
      if (parsed.fieldErrors.site_name) {
        setCurrentTab('profil')
      } else if (parsed.fieldErrors.contact_email) {
        setCurrentTab('kontak')
      }
      if (Object.keys(parsed.fieldErrors).length === 0) {
        showToast(parsed.message || 'Terjadi kesalahan saat menyimpan pengaturan', true)
      }
    } finally {
      setLoading(false)
    }
  }

  const inputCls = "w-full px-3 py-2 border border-slate-300 rounded-xl focus:ring-brand-500 focus:border-brand-500 text-sm bg-white outline-none transition-colors"

  return (
    <AdminLayout title="Pengaturan Website">
      <div className="max-w-5xl mx-auto grid grid-cols-1 md:grid-cols-4 gap-6 animate-fade-in-up">

        {/* Navigation Tabs */}
         <div className="md:col-span-1 bg-white p-4 rounded-2xl border border-gray-200 shadow-sm flex flex-col gap-1 h-fit">
          <button onClick={() => setCurrentTab('profil')} type="button" className={`px-4 py-2 text-left text-sm font-semibold rounded-lg transition-all duration-200 btn-press ${currentTab === 'profil' ? 'bg-brand-50 text-brand-600 font-medium' : 'text-slate-600 hover:bg-slate-50'}`}>Profil & Sejarah</button>
          <button onClick={() => setCurrentTab('sambutan')} type="button" className={`px-4 py-2 text-left text-sm font-semibold rounded-lg transition-all duration-200 btn-press ${currentTab === 'sambutan' ? 'bg-brand-50 text-brand-600 font-medium' : 'text-slate-600 hover:bg-slate-50'}`}>Sambutan Depan</button>
          <button onClick={() => setCurrentTab('kontak')} type="button" className={`px-4 py-2 text-left text-sm font-semibold rounded-lg transition-all duration-200 btn-press ${currentTab === 'kontak' ? 'bg-brand-50 text-brand-600 font-medium' : 'text-slate-600 hover:bg-slate-50'}`}>Informasi Kontak</button>
          <button onClick={() => setCurrentTab('sosmed')} type="button" className={`px-4 py-2 text-left text-sm font-semibold rounded-lg transition-all duration-200 btn-press ${currentTab === 'sosmed' ? 'bg-brand-50 text-brand-600 font-medium' : 'text-slate-600 hover:bg-slate-50'}`}>Media Sosial</button>
          <button onClick={() => setCurrentTab('logo')} type="button" className={`px-4 py-2 text-left text-sm font-semibold rounded-lg transition-all duration-200 btn-press ${currentTab === 'logo' ? 'bg-brand-50 text-brand-600 font-medium' : 'text-slate-600 hover:bg-slate-50'}`}>Logo Website</button>
          <button onClick={() => setCurrentTab('log')} type="button" className={`px-4 py-2 text-left text-sm font-semibold rounded-lg transition-all duration-200 btn-press ${currentTab === 'log' ? 'bg-brand-50 text-brand-600 font-medium' : 'text-slate-600 hover:bg-slate-50'}`}>Pembersihan Log</button>
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
                      src={logoPreview || resolveAssetUrl(formData.logo_path) || 'https://via.placeholder.com/160'}
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
                    {formData.logo_path && (
                      <span className="text-[11px] text-gray-400 break-all">Path saat ini: {formData.logo_path}</span>
                    )}
                  </div>
                </div>
              </div>
            )}

            {currentTab === 'log' && (
              <div className="space-y-4">
                <h3 className="font-heading font-bold text-gray-800 text-base border-b pb-2">Auto Clear Log Aktivitas</h3>
                <p className="text-xs text-gray-500 leading-relaxed">Log aktivitas super admin & admin yang telah melewati jangka waktu di bawah ini akan dihapus secara otomatis secara berkala setiap 24 jam demi menghemat kapasitas penyimpanan database.</p>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Masa Simpan / Retensi Log Aktivitas</label>
                  <select
                    value={formData.log_retention_days}
                    onChange={e => setFormData({ ...formData, log_retention_days: Number(e.target.value) })}
                    className={inputCls}
                  >
                    <option value={7}>7 Hari (1 Minggu)</option>
                    <option value={30}>30 Hari (1 Bulan)</option>
                    <option value={90}>90 Hari (3 Bulan)</option>
                    <option value={180}>180 Hari (6 Bulan)</option>
                    <option value={365}>365 Hari (1 Tahun)</option>
                    <option value={0}>Selamanya (Jangan Hapus Log)</option>
                  </select>
                </div>
              </div>
            )}

            {currentTab === 'profil' && (
              <div className="space-y-4">
                <h3 className="font-heading font-bold text-gray-800 text-base border-b pb-2">Profil & Sejarah Organisasi</h3>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Nama Website / Organisasi <span className="text-red-500">*</span></label>
                    <input
                      type="text"
                      value={formData.site_name}
                      onChange={e => {
                        setFormData({...formData, site_name: e.target.value})
                        clearFieldError('site_name')
                        if (touched.site_name) {
                          const errs = validateForm({...formData, site_name: e.target.value})
                          setFormErrors(prev => ({ ...prev, site_name: errs.site_name }))
                        }
                      }}
                      onBlur={() => {
                        setTouched(prev => ({ ...prev, site_name: true }))
                        const errs = validateForm()
                        setFormErrors(prev => ({ ...prev, site_name: errs.site_name }))
                      }}
                      className={`${inputCls} ${(touched.site_name && formErrors.site_name) || fieldErrors.site_name ? 'border-red-400 focus:ring-2 focus:ring-red-100' : ''}`}
                    />
                    {((touched.site_name && formErrors.site_name) || fieldErrors.site_name) && (
                      <p className="text-red-500 text-[11px] font-semibold mt-1 flex items-center gap-1">
                        <i className="ph-bold ph-warning-circle text-xs" /> {formErrors.site_name || fieldErrors.site_name}
                      </p>
                    )}
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
                  <textarea rows="3" value={formData.about_vision} onChange={e => setFormData({...formData, about_vision: e.target.value})} className={`${inputCls} overflow-y-auto resize-y min-h-[80px]`} />
                </div>
                <div>
                  <div className="flex items-center justify-between mb-2">
                    <label className="block text-xs font-semibold text-gray-500">Misi Organisasi</label>
                    <button
                      type="button"
                      onClick={() => {
                        setFormData(prev => ({
                          ...prev,
                          about_mission: [...(prev.about_mission || []), '']
                        }))
                      }}
                      className="inline-flex items-center gap-1 text-xs font-bold text-brand-600 hover:text-brand-700 transition"
                    >
                      <i className="ph-bold ph-plus-circle text-sm" /> Tambah Misi
                    </button>
                  </div>
                  
                  <div className="max-h-56 overflow-y-auto space-y-2.5 pr-2 border border-slate-200/60 rounded-xl p-3 bg-slate-50/50">
                    {(formData.about_mission || []).length === 0 ? (
                      <p className="text-xs text-gray-400 text-center py-2">Belum ada misi. Klik Tambah Misi.</p>
                    ) : (
                      (formData.about_mission || []).map((misi, index) => (
                        <div key={index} className="flex gap-2 items-center animate-scale-up">
                          <span className="text-xs font-bold text-gray-400 w-6 text-center">{index + 1}</span>
                          <input
                            type="text"
                            value={misi}
                            placeholder={`Masukkan misi ke-${index + 1}...`}
                            onChange={e => {
                              const newMissions = [...formData.about_mission]
                              newMissions[index] = e.target.value
                              setFormData({ ...formData, about_mission: newMissions })
                            }}
                            className={inputCls}
                          />
                          <button
                            type="button"
                            disabled={formData.about_mission.length <= 1}
                            onClick={() => {
                              const newMissions = formData.about_mission.filter((_, i) => i !== index)
                              setFormData({ ...formData, about_mission: newMissions })
                            }}
                            className="p-2 text-gray-400 hover:text-red-500 disabled:opacity-30 disabled:cursor-not-allowed hover:bg-red-50 rounded-xl transition"
                            title="Hapus Misi"
                          >
                            <i className="ph-bold ph-trash text-base" />
                          </button>
                        </div>
                      ))
                    )}
                  </div>
                </div>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Sejarah / Tentang Kami</label>
                  <textarea rows="6" value={formData.history} onChange={e => setFormData({...formData, history: e.target.value})} className={`${inputCls} overflow-y-auto resize-y min-h-[120px]`} />
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
                  <textarea rows="5" value={formData.greeting_content} onChange={e => setFormData({...formData, greeting_content: e.target.value})} className={`${inputCls} overflow-y-auto resize-y min-h-[120px]`} />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">URL Gambar / Poster Sambutan</label>
                  <input type="text" value={formData.greeting_image_path} onChange={e => setFormData({...formData, greeting_image_path: e.target.value})} className={inputCls} />
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Nama Organisasi / Pengirim di Tanda Tangan (Signature)</label>
                    <input 
                      type="text" 
                      placeholder="Contoh: GRADASI (Kosongkan untuk default nama website)" 
                      value={formData.greeting_sign_name || ''} 
                      onChange={e => setFormData({...formData, greeting_sign_name: e.target.value})} 
                      className={inputCls} 
                    />
                  </div>
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Nama Penandatangan (Signature)</label>
                    <input 
                      type="text" 
                      placeholder="Contoh: Upi Asmaradhana & Junaidi, S.Kom (Kosongkan untuk default data pengurus)" 
                      value={formData.greeting_sign_subtitle || ''} 
                      onChange={e => setFormData({...formData, greeting_sign_subtitle: e.target.value})} 
                      className={inputCls} 
                    />
                  </div>
                </div>
              </div>
            )}

            {currentTab === 'kontak' && (
              <div className="space-y-4">
                <h3 className="font-heading font-bold text-gray-800 text-base border-b pb-2">Informasi Kontak & Profil Kantor</h3>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Email Resmi <span className="text-gray-400 font-normal">(opsional)</span></label>
                    <input
                      type="email"
                      value={formData.contact_email}
                      onChange={e => {
                        setFormData({...formData, contact_email: e.target.value})
                        clearFieldError('contact_email')
                        if (touched.contact_email) {
                          const errs = validateForm({...formData, contact_email: e.target.value})
                          setFormErrors(prev => ({ ...prev, contact_email: errs.contact_email }))
                        }
                      }}
                      onBlur={() => {
                        setTouched(prev => ({ ...prev, contact_email: true }))
                        const errs = validateForm()
                        setFormErrors(prev => ({ ...prev, contact_email: errs.contact_email }))
                      }}
                      className={`${inputCls} ${(touched.contact_email && formErrors.contact_email) || fieldErrors.contact_email ? 'border-red-400 focus:ring-2 focus:ring-red-100' : ''}`}
                    />
                    {((touched.contact_email && formErrors.contact_email) || fieldErrors.contact_email) && (
                      <p className="text-red-500 text-[11px] font-semibold mt-1 flex items-center gap-1">
                        <i className="ph-bold ph-warning-circle text-xs" /> {formErrors.contact_email || fieldErrors.contact_email}
                      </p>
                    )}
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
                  <input type="text" value={formData.video_profile_path} onChange={e => setFormData({...formData, video_profile_path: e.target.value})} className={inputCls} />
                </div>
              </div>
            )}

            <div className="flex justify-end pt-4 border-t border-slate-100">
              <button type="submit" disabled={loading} className="px-6 py-2.5 bg-brand-600 text-white rounded-xl hover:bg-brand-700 text-sm font-semibold transition btn-press">
                {loading ? 'Menyimpan...' : 'Simpan Semua Pengaturan'}
              </button>
            </div>
          </form>
        </div>

      </div>

      {/* Floating Toast Notification */}
      {saved && (
        <div className={`fixed bottom-4 right-4 z-50 flex items-center p-4 rounded-xl shadow-xl text-white transition-opacity duration-300 ${
          saved.isError ? 'bg-red-500' : 'bg-emerald-500'
        }`}>
          <i className={`text-xl mr-2 ph ${saved.isError ? 'ph-warning-circle' : 'ph-check-circle'}`} />
          <span className="text-sm font-semibold">{saved.msg}</span>
        </div>
      )}
    </AdminLayout>
  )
}
