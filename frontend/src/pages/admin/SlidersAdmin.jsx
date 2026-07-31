import { useState, useEffect } from 'react'
import AdminLayout from '../../layouts/AdminLayout'
import { slidersService } from '../../services/slidersService'

export default function SlidersAdmin() {
  const [items, setItems] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  
  const [isFormOpen, setIsFormOpen] = useState(false)
  const [formMode, setFormMode] = useState('create')
  
  const [formData, setFormData] = useState({
    id: null,
    title: '',
    subtitle: '',
    tag: '',
    image_url: '',
    link_url: '',
    sort_order: 1,
    event_date: '',
    location: '',
    is_new: false,
    is_active: true
  })

  const loadSliders = () => {
    setLoading(true)
    slidersService.list(false)
      .then(res => {
        if (res.success && res.data && res.data.sliders) {
          setItems(res.data.sliders)
        } else {
          setError('Gagal memuat sliders')
        }
      })
      .catch(() => setError('Kesalahan koneksi ke server'))
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    loadSliders()
  }, [])

  const openForm = (item = null) => {
    if (item) {
      setFormMode('edit')
      setFormData({
        id: item.id,
        title: item.title,
        subtitle: item.subtitle || '',
        tag: item.tag || '',
        image_url: item.image_url,
        link_url: item.link_url || '',
        sort_order: item.sort_order,
        event_date: item.event_date || '',
        location: item.location || '',
        is_new: item.is_new,
        is_active: item.is_active
      })
    } else {
      setFormMode('create')
      setFormData({
        id: null,
        title: '',
        subtitle: '',
        tag: '',
        image_url: '',
        link_url: '',
        sort_order: items.length + 1,
        event_date: '',
        location: '',
        is_new: false,
        is_active: true
      })
    }
    setIsFormOpen(true)
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    try {
      if (formMode === 'create') {
        await slidersService.create(formData)
      } else {
        await slidersService.update(formData.id, formData)
      }
      setIsFormOpen(false)
      loadSliders()
      alert('Data slider berhasil disimpan!')
    } catch {
      alert('Gagal menyimpan data slider')
    }
  }

  const handleDelete = async (id) => {
    if (window.confirm('Yakin ingin menghapus slider ini?')) {
      try {
        await slidersService.remove(id)
        loadSliders()
        alert('Slider berhasil dihapus!')
      } catch {
        alert('Gagal menghapus slider')
      }
    }
  }

  return (
    <AdminLayout title="Kelola Banner Sliders">
      <div className="space-y-6">
        <div className="flex justify-between items-center bg-white p-4 rounded-xl border border-gray-200">
          <h2 className="text-sm font-semibold text-gray-600">Daftar Sliders Utama</h2>
          <button 
            onClick={() => openForm()} 
            className="px-4 py-2 bg-brand-600 text-white rounded-lg hover:bg-brand-700 text-xs font-semibold"
          >
            + Tambah Slider Baru
          </button>
        </div>

        {loading && <div className="text-slate-500 py-10 text-center">Memuat sliders...</div>}
        {error && <div className="text-red-600 py-10 text-center font-medium">{error}</div>}

        {!loading && !error && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {items.map(item => (
              <div key={item.id} className="bg-white rounded-2xl border border-gray-200 overflow-hidden shadow-sm flex flex-col justify-between">
                <div className="h-44 relative bg-gray-100">
                  <img src={item.image_url.startsWith('http') ? item.image_url : `http://127.0.0.1:8080${item.image_url}`} alt={item.title} className="w-full h-full object-cover" />
                </div>
                <div className="p-5 flex-grow space-y-2">
                  <span className="text-[10px] font-bold text-brand-600 uppercase tracking-widest block">{item.tag || 'SLIDER'}</span>
                  <h3 className="font-heading font-bold text-slate-800 text-lg line-clamp-1">{item.title}</h3>
                  <p className="text-slate-500 text-xs line-clamp-2">{item.subtitle}</p>
                </div>
                <div className="p-4 border-t border-gray-100 flex justify-end gap-2 bg-gray-50/50">
                  <button onClick={() => openForm(item)} className="px-3 py-1.5 bg-gray-100 text-slate-700 rounded-lg hover:bg-gray-200 text-xs font-semibold">Edit</button>
                  <button onClick={() => handleDelete(item.id)} className="px-3 py-1.5 bg-red-50 text-red-700 rounded-lg hover:bg-red-100 text-xs font-semibold">Hapus</button>
                </div>
              </div>
            ))}
          </div>
        )}

        {/* Modal Form */}
        {isFormOpen && (
          <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
            <div className="bg-white rounded-2xl p-6 max-w-lg w-full shadow-2xl max-h-[90vh] overflow-y-auto space-y-4">
              <h3 className="font-heading font-bold text-lg text-slate-900">{formMode === 'create' ? 'Tambah Slider Baru' : 'Edit Slider'}</h3>
              
              <form onSubmit={handleSubmit} className="space-y-4">
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Judul Utama</label>
                  <input type="text" value={formData.title} onChange={e => setFormData({...formData, title: e.target.value})} required className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Sub Judul / Deskripsi</label>
                  <input type="text" value={formData.subtitle} onChange={e => setFormData({...formData, subtitle: e.target.value})} className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Tag / Kategori Badge</label>
                    <input type="text" value={formData.tag} onChange={e => setFormData({...formData, tag: e.target.value})} className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                  </div>
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Urutan Sortir</label>
                    <input type="number" value={formData.sort_order} onChange={e => setFormData({...formData, sort_order: parseInt(e.target.value)})} required className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                  </div>
                </div>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">URL Gambar Banner</label>
                  <input type="text" value={formData.image_url} onChange={e => setFormData({...formData, image_url: e.target.value})} required className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Link URL Aksi</label>
                  <input type="text" value={formData.link_url} onChange={e => setFormData({...formData, link_url: e.target.value})} className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Tanggal Kegiatan (Opsional)</label>
                    <input type="text" value={formData.event_date} onChange={e => setFormData({...formData, event_date: e.target.value})} className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                  </div>
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Lokasi Kegiatan (Opsional)</label>
                    <input type="text" value={formData.location} onChange={e => setFormData({...formData, location: e.target.value})} className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                  </div>
                </div>
                <div className="flex gap-4">
                  <label className="flex items-center gap-2">
                    <input type="checkbox" checked={formData.is_new} onChange={e => setFormData({...formData, is_new: e.target.checked})} />
                    <span className="text-xs text-gray-600">Terbaru (New) Badge</span>
                  </label>
                  <label className="flex items-center gap-2">
                    <input type="checkbox" checked={formData.is_active} onChange={e => setFormData({...formData, is_active: e.target.checked})} />
                    <span className="text-xs text-gray-600">Tampilkan / Aktif</span>
                  </label>
                </div>
                
                <div className="flex justify-end gap-2 pt-4 border-t border-gray-100">
                  <button type="button" onClick={() => setIsFormOpen(false)} className="px-4 py-2 border border-gray-300 rounded-lg text-xs font-semibold text-gray-700">Batal</button>
                  <button type="submit" className="px-4 py-2 bg-brand-600 text-white rounded-lg hover:bg-brand-700 text-xs font-semibold">Simpan</button>
                </div>
              </form>
            </div>
          </div>
        )}
      </div>
    </AdminLayout>
  )
}
