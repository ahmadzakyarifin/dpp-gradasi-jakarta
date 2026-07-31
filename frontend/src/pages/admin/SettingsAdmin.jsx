import { useState, useEffect } from 'react'
import AdminLayout from '../../layouts/AdminLayout'
import { settingsService } from '../../services/settingsService'

export default function SettingsAdmin() {
  const [currentTab, setCurrentTab] = useState('profil')
  const [loading, setLoading] = useState(false)
  
  const [formData, setFormData] = useState({
    site_name: '',
    tagline: '',
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
    settingsService.getAdmin()
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

  const handleSubmit = async (e) => {
    e.preventDefault()
    setLoading(true)
    
    // convert about_mission back to array/JSON format
    const payload = {
      ...formData,
      about_mission: JSON.stringify(formData.about_mission.split('\n').map(s => s.trim()).filter(Boolean))
    }

    try {
      const res = await settingsService.update(payload)
      if (res.success) {
        alert('Pengaturan website berhasil disimpan!')
      } else {
        alert('Gagal menyimpan: ' + res.message)
      }
    } catch {
      alert('Terjadi kesalahan')
    } finally {
      setLoading(false)
    }
  }

  return (
    <AdminLayout title="Pengaturan Website">
      <div className="max-w-5xl mx-auto grid grid-cols-1 md:grid-cols-4 gap-6">
        
        {/* Navigation Tabs */}
        <div className="md:col-span-1 bg-white p-4 rounded-2xl border border-gray-200 shadow-sm flex flex-col gap-1">
          <button onClick={() => setCurrentTab('profil')} className={`px-4 py-2 text-left text-sm font-semibold rounded-lg transition ${currentTab === 'profil' ? 'bg-brand-50 text-brand-600' : 'text-slate-600 hover:bg-slate-50'}`}>Profil & Sejarah</button>
          <button onClick={() => setCurrentTab('sambutan')} className={`px-4 py-2 text-left text-sm font-semibold rounded-lg transition ${currentTab === 'sambutan' ? 'bg-brand-50 text-brand-600' : 'text-slate-600 hover:bg-slate-50'}`}>Sambutan Depan</button>
          <button onClick={() => setCurrentTab('kontak')} className={`px-4 py-2 text-left text-sm font-semibold rounded-lg transition ${currentTab === 'kontak' ? 'bg-brand-50 text-brand-600' : 'text-slate-600 hover:bg-slate-50'}`}>Informasi Kontak</button>
          <button onClick={() => setCurrentTab('sosmed')} className={`px-4 py-2 text-left text-sm font-semibold rounded-lg transition ${currentTab === 'sosmed' ? 'bg-brand-50 text-brand-600' : 'text-slate-600 hover:bg-slate-50'}`}>Media Sosial</button>
        </div>

        {/* Form Area */}
        <div className="md:col-span-3 bg-white p-6 rounded-2xl border border-gray-200 shadow-sm">
          <form onSubmit={handleSubmit} className="space-y-6">
            
            {currentTab === 'profil' && (
              <div className="space-y-4">
                <h3 className="font-heading font-bold text-gray-800 text-base border-b pb-2">Profil & Sejarah Organisasi</h3>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Tanggal Terbentuk</label>
                    <input type="text" value={formData.about_formation_date} onChange={e => setFormData({...formData, about_formation_date: e.target.value})} className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                  </div>
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Nomor SK Legalitas</label>
                    <input type="text" value={formData.about_no_sk} onChange={e => setFormData({...formData, about_no_sk: e.target.value})} className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                  </div>
                </div>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Visi Utama</label>
                  <textarea rows="2" value={formData.about_vision} onChange={e => setFormData({...formData, about_vision: e.target.value})} className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Misi Organisasi (Satu Misi Per Baris)</label>
                  <textarea rows="4" value={formData.about_mission} onChange={e => setFormData({...formData, about_mission: e.target.value})} className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Sejarah / Tentang Kami</label>
                  <textarea rows="4" value={formData.history} onChange={e => setFormData({...formData, history: e.target.value})} className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                </div>
              </div>
            )}

            {currentTab === 'sambutan' && (
              <div className="space-y-4">
                <h3 className="font-heading font-bold text-gray-800 text-base border-b pb-2">Sambutan Halaman Depan</h3>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Judul Sambutan</label>
                  <input type="text" value={formData.greeting_title} onChange={e => setFormData({...formData, greeting_title: e.target.value})} className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Sub Judul</label>
                  <input type="text" value={formData.greeting_subtitle} onChange={e => setFormData({...formData, greeting_subtitle: e.target.value})} className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Isi Sambutan</label>
                  <textarea rows="5" value={formData.greeting_content} onChange={e => setFormData({...formData, greeting_content: e.target.value})} className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">URL Gambar / Poster Sambutan</label>
                  <input type="text" value={formData.greeting_image_url} onChange={e => setFormData({...formData, greeting_image_url: e.target.value})} className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                </div>
              </div>
            )}

            {currentTab === 'kontak' && (
              <div className="space-y-4">
                <h3 className="font-heading font-bold text-gray-800 text-base border-b pb-2">Informasi Kontak & Profil Kantor</h3>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Email Resmi</label>
                    <input type="email" value={formData.contact_email} onChange={e => setFormData({...formData, contact_email: e.target.value})} className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                  </div>
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Telepon / WhatsApp</label>
                    <input type="text" value={formData.contact_phone} onChange={e => setFormData({...formData, contact_phone: e.target.value})} className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                  </div>
                </div>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Alamat Kantor</label>
                  <input type="text" value={formData.address} onChange={e => setFormData({...formData, address: e.target.value})} className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Google Maps Embed URL</label>
                  <input type="text" value={formData.maps_embed_url} onChange={e => setFormData({...formData, maps_embed_url: e.target.value})} className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                </div>
              </div>
            )}

            {currentTab === 'sosmed' && (
              <div className="space-y-4">
                <h3 className="font-heading font-bold text-gray-800 text-base border-b pb-2">Social Media & Profil Link</h3>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Facebook URL</label>
                  <input type="text" value={formData.facebook_url} onChange={e => setFormData({...formData, facebook_url: e.target.value})} className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Instagram URL</label>
                  <input type="text" value={formData.instagram_url} onChange={e => setFormData({...formData, instagram_url: e.target.value})} className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">YouTube Channel URL</label>
                  <input type="text" value={formData.youtube_url} onChange={e => setFormData({...formData, youtube_url: e.target.value})} className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Video Profile URL (Youtube Embed)</label>
                  <input type="text" value={formData.video_profile_url} onChange={e => setFormData({...formData, video_profile_url: e.target.value})} className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
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
