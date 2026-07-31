import { useState, useEffect } from 'react'
import AdminLayout from '../../layouts/AdminLayout'
import { kegiatanService } from '../../services/kegiatanService'
import { useDataStore } from '../../store/useDataStore'

const DEFAULT_KEGIATAN = [
  {
    id: 1,
    title: 'Penyaluran Bantuan Kemanusiaan oleh DPP GRADASI',
    slug: 'penyaluran-bantuan-kemanusiaan',
    category: 'Nasional',
    organizer: 'DPP GRADASI',
    event_date: '31 Desember 2025',
    location: 'Jakarta',
    image_url: 'https://gradasi.org/uploads/img/event/1767154719.jpg',
    excerpt: 'Dewan Pimpinan Pusat (DPP) GRADASI turun langsung menyalurkan bantuan kemanusiaan...',
    is_published: true
  },
  {
    id: 2,
    title: 'Pelatihan Digital Marketing UMKM Go Online',
    slug: 'pelatihan-digital-marketing-umkm',
    category: 'Jawa Timur',
    organizer: 'DPD GRADASI Jatim',
    event_date: '15 November 2025',
    location: 'Surabaya',
    image_url: 'https://gradasi.org/uploads/img/event/1767154619.jpg',
    excerpt: 'Program pelatihan intensif bagi pelaku Usaha Mikro Kecil Menengah (UMKM)...',
    is_published: true
  },
  {
    id: 3,
    title: 'Konsolidasi Pengurus DPP & Penyerahan SK Daerah',
    slug: 'konsolidasi-pengurus-dpp-dpd',
    category: 'Lampung',
    organizer: 'DPP GRADASI',
    event_date: '02 Oktober 2025',
    location: 'Bandar Lampung',
    image_url: 'https://gradasi.org/uploads/img/event/1767154397.jpg',
    excerpt: 'Acara konsolidasi pengurus tingkat pusat serta penyerahan Surat Keputusan (SK)...',
    is_published: true
  }
]

export default function KegiatanAdmin() {
  const [items, setItems] = useState(DEFAULT_KEGIATAN)
  const [loading, setLoading] = useState(false)
  
  const [currentTab, setCurrentTab] = useState('active')
  const [currentPage, setCurrentPage] = useState(1)
  const pageSize = 5

  const [searchQuery, setSearchQuery] = useState('')
  const [filterStatus, setFilterStatus] = useState('')
  const [filterSort, setFilterSort] = useState('newest')

  const [isFormOpen, setIsFormOpen] = useState(false)
  const [formMode, setFormMode] = useState('create')
  
  const [formData, setFormData] = useState({
    id: null,
    title: '',
    category: 'Kegiatan',
    organizer: 'DPP GRADASI',
    eventDate: '',
    location: '',
    image: '',
    excerpt: '',
    content: '',
    isPublished: true
  })

  const loadKegiatan = () => {
    const serviceCall = currentTab === 'active' 
      ? kegiatanService.list({ active_only: true }) 
      : kegiatanService.list({ trash_only: true })
      
    serviceCall
      .then(res => {
        if (res && res.data) {
          const list = Array.isArray(res.data) ? res.data : (res.data.kegiatan || [])
          if (list.length > 0) setItems(list)
        }
      })
      .catch(() => {})
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    loadKegiatan()
  }, [currentTab])

  const filteredItems = items.filter(item => {
    const matchesSearch = !searchQuery || item.title.toLowerCase().includes(searchQuery.toLowerCase()) || 
                          (item.location || '').toLowerCase().includes(searchQuery.toLowerCase())
    let matchesStatus = true
    if (filterStatus === 'published' && !item.is_published) matchesStatus = false
    if (filterStatus === 'draft' && item.is_published) matchesStatus = false
    return matchesSearch && matchesStatus
  })

  const paginatedItems = filteredItems.slice((currentPage - 1) * pageSize, currentPage * pageSize)
  const totalPages = Math.max(1, Math.ceil(filteredItems.length / pageSize))

  const openForm = (item = null) => {
    if (item) {
      setFormMode('edit')
      setFormData({
        id: item.id,
        title: item.title,
        category: item.category || 'Kegiatan',
        organizer: item.organizer || 'DPP GRADASI',
        eventDate: item.event_date || '',
        location: item.location || '',
        image: item.image_url || '',
        excerpt: item.excerpt || '',
        content: item.content || item.excerpt || '',
        isPublished: item.is_published !== false
      })
    } else {
      setFormMode('create')
      setFormData({
        id: null,
        title: '',
        category: 'Kegiatan',
        organizer: 'DPP GRADASI',
        eventDate: '',
        location: '',
        image: '',
        excerpt: '',
        content: '',
        isPublished: true
      })
    }
    setIsFormOpen(true)
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    const payload = {
      title: formData.title,
      category: formData.category,
      organizer: formData.organizer,
      event_date: formData.eventDate,
      location: formData.location,
      image_url: formData.image,
      excerpt: formData.excerpt,
      content: formData.content,
      is_published: formData.isPublished
    }

    if (formMode === 'create') {
      const newObj = { ...payload, id: Date.now() }
      setItems(prev => [newObj, ...prev])
      useDataStore.getState().addKegiatan(newObj)
      try { await kegiatanService.create(payload) } catch {}
    } else {
      setItems(prev => prev.map(i => i.id === formData.id ? { ...i, ...payload } : i))
      useDataStore.getState().updateKegiatan(formData.id, payload)
      try { await kegiatanService.update(formData.id, payload) } catch {}
    }
    setIsFormOpen(false)
  }

  const handleDelete = async (id) => {
    if (window.confirm('Yakin ingin menghapus kegiatan ini?')) {
      setItems(prev => prev.filter(i => i.id !== id))
      useDataStore.getState().deleteKegiatan(id)
      try { await kegiatanService.remove(id) } catch {}
    }
  }

  return (
    <AdminLayout title="Kelola Kegiatan">
      <div className="space-y-6">
        <div className="flex justify-between items-center bg-white p-4 rounded-xl border border-gray-200 shadow-sm">
          <div className="flex gap-2">
            <button 
              onClick={() => { setCurrentTab('active'); setCurrentPage(1); }}
              className={`px-4 py-2 text-xs font-semibold rounded-lg transition ${currentTab === 'active' ? 'bg-brand-600 text-white' : 'bg-gray-100 text-slate-700'}`}
            >
              Kegiatan Aktif
            </button>
            <button 
              onClick={() => { setCurrentTab('trash'); setCurrentPage(1); }}
              className={`px-4 py-2 text-xs font-semibold rounded-lg transition ${currentTab === 'trash' ? 'bg-red-50 text-red-600' : 'bg-gray-100 text-slate-700'}`}
            >
              Sampah
            </button>
          </div>
          {currentTab === 'active' && (
            <button 
              onClick={() => openForm()}
              className="px-4 py-2.5 bg-brand-600 text-white rounded-xl hover:bg-brand-700 text-xs font-semibold flex items-center gap-2 shadow-sm"
            >
              + Tambah Kegiatan Baru
            </button>
          )}
        </div>

        {/* Filter Bar */}
        <div className="flex flex-col md:flex-row gap-4 justify-between items-center bg-white p-4 rounded-xl border border-gray-200 shadow-sm">
          <input 
            type="text" 
            placeholder="Cari kegiatan atau lokasi..." 
            value={searchQuery}
            onChange={e => { setSearchQuery(e.target.value); setCurrentPage(1); }}
            className="w-full md:max-w-md px-3.5 py-2.5 bg-slate-50 border border-gray-200 rounded-xl text-sm outline-none"
          />
          <div className="flex gap-2 w-full md:w-auto">
            <select 
              value={filterStatus}
              onChange={e => { setFilterStatus(e.target.value); setCurrentPage(1); }}
              className="px-3.5 py-2.5 bg-slate-50 border border-gray-200 rounded-xl text-sm text-slate-600 outline-none"
            >
              <option value="">Semua Status</option>
              <option value="published">Terbit</option>
              <option value="draft">Draft</option>
            </select>
          </div>
        </div>

        {/* Data Table */}
        <div className="bg-white rounded-2xl border border-gray-200 overflow-hidden shadow-sm">
          <table className="w-full text-left text-sm text-slate-700">
            <thead className="bg-slate-50 border-b border-gray-200 font-semibold text-xs uppercase tracking-wider text-slate-500">
              <tr>
                <th className="p-4">Kegiatan</th>
                <th className="p-4">Tanggal & Lokasi</th>
                <th className="p-4">Kategori</th>
                <th className="p-4">Status</th>
                <th className="p-4 text-right">Aksi</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100">
              {paginatedItems.map(item => (
                <tr key={item.id} className="hover:bg-slate-50/60 transition">
                  <td className="p-4">
                    <div className="flex items-center gap-3">
                      <img src={item.image_url} alt="" className="w-16 h-12 rounded-lg object-cover border border-slate-200 shrink-0" />
                      <p className="font-bold text-slate-900 line-clamp-1">{item.title}</p>
                    </div>
                  </td>
                  <td className="p-4 text-slate-500 text-xs">
                    <p className="font-semibold text-slate-700">{item.event_date}</p>
                    <p>{item.location}</p>
                  </td>
                  <td className="p-4">
                    <span className="bg-brand-50 text-brand-700 text-[10px] font-bold px-2.5 py-1 rounded-full uppercase">
                      {item.category}
                    </span>
                  </td>
                  <td className="p-4">
                    <span className={`text-xs font-semibold px-2.5 py-1 rounded-lg ${item.is_published ? 'bg-emerald-50 text-emerald-600' : 'bg-slate-100 text-slate-500'}`}>
                      {item.is_published ? 'Terbit' : 'Draft'}
                    </span>
                  </td>
                  <td className="p-4 text-right">
                    <div className="flex justify-end gap-2">
                      <button onClick={() => openForm(item)} className="p-2 text-slate-400 hover:text-brand-600 rounded-lg">Edit</button>
                      <button onClick={() => handleDelete(item.id)} className="p-2 text-slate-400 hover:text-red-600 rounded-lg">Hapus</button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* Form Modal */}
      {isFormOpen && (
        <div className="fixed inset-0 bg-black/40 backdrop-blur-[2px] flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl p-6 max-w-lg w-full shadow-2xl space-y-4 max-h-[90vh] overflow-y-auto">
            <h3 className="font-heading font-bold text-lg text-slate-900">{formMode === 'create' ? 'Tambah Kegiatan Baru' : 'Edit Kegiatan'}</h3>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-xs font-semibold text-gray-500 mb-1">Judul Kegiatan *</label>
                <input type="text" value={formData.title} onChange={e => setFormData({ ...formData, title: e.target.value })} required className="w-full px-3.5 py-2.5 border rounded-xl text-sm outline-none" />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Tanggal Event</label>
                  <input type="text" value={formData.eventDate} onChange={e => setFormData({ ...formData, eventDate: e.target.value })} placeholder="31 Desember 2025" className="w-full px-3.5 py-2.5 border rounded-xl text-sm outline-none" />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Lokasi</label>
                  <input type="text" value={formData.location} onChange={e => setFormData({ ...formData, location: e.target.value })} placeholder="Jakarta" className="w-full px-3.5 py-2.5 border rounded-xl text-sm outline-none" />
                </div>
              </div>
              <div>
                <label className="block text-xs font-semibold text-gray-500 mb-1">URL Gambar</label>
                <input type="text" value={formData.image} onChange={e => setFormData({ ...formData, image: e.target.value })} className="w-full px-3.5 py-2.5 border rounded-xl text-sm outline-none" />
              </div>
              <div>
                <label className="block text-xs font-semibold text-gray-500 mb-1">Ringkasan</label>
                <textarea rows={3} value={formData.excerpt} onChange={e => setFormData({ ...formData, excerpt: e.target.value })} className="w-full px-3.5 py-2.5 border rounded-xl text-sm outline-none" />
              </div>
              <div className="flex justify-end gap-2 pt-4 border-t">
                <button type="button" onClick={() => setIsFormOpen(false)} className="px-4 py-2 border rounded-xl text-sm font-semibold">Batal</button>
                <button type="submit" className="px-5 py-2 bg-brand-600 text-white rounded-xl text-sm font-semibold">Simpan</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </AdminLayout>
  )
}
