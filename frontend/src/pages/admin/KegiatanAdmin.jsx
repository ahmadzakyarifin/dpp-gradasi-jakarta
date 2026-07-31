import { useState, useEffect } from 'react'
import AdminLayout from '../../layouts/AdminLayout'
import { kegiatanService } from '../../services/kegiatanService'

export default function KegiatanAdmin() {
  const [items, setItems] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  
  const [currentTab, setCurrentTab] = useState('active') // active, trash
  const [currentPage, setCurrentPage] = useState(1)
  const pageSize = 5

  const [searchQuery, setSearchQuery] = useState('')
  const [filterStatus, setFilterStatus] = useState('') // all, published, draft
  const [filterSort, setFilterSort] = useState('newest')

  const [isFormOpen, setIsFormOpen] = useState(false)
  const [formMode, setFormMode] = useState('create')
  
  const [formData, setFormData] = useState({
    id: null,
    title: '',
    category: 'Kegiatan',
    organizer: 'DPP GRADASI',
    author: 'Super Admin',
    eventDate: '',
    location: '',
    image: '',
    content: '',
    tags: '',
    isPublished: true,
    gallery: [] // array of {image_url, caption}
  })

  // Temporary state for adding a gallery item in the form modal
  const [newGalleryUrl, setNewGalleryUrl] = useState('')
  const [newGalleryCaption, setNewGalleryCaption] = useState('')

  const loadKegiatan = () => {
    setLoading(true)
    const serviceCall = currentTab === 'active' 
      ? kegiatanService.list({ active_only: true }) 
      : kegiatanService.list({ trash_only: true })
      
    serviceCall
      .then(res => {
        if (res.success && res.data) {
          setItems(res.data)
        } else {
          setError('Gagal memuat kegiatan')
        }
      })
      .catch(() => setError('Kesalahan koneksi ke server'))
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    loadKegiatan()
  }, [currentTab])

  // Filter & Sort
  const filteredItems = items.filter(item => {
    const matchesSearch = item.title.toLowerCase().includes(searchQuery.toLowerCase()) || 
                          item.location.toLowerCase().includes(searchQuery.toLowerCase())
    
    let matchesStatus = true
    if (filterStatus === 'published' && !item.is_published) matchesStatus = false
    if (filterStatus === 'draft' && item.is_published) matchesStatus = false

    return matchesSearch && matchesStatus
  }).sort((a, b) => {
    if (filterSort === 'newest') return new Date(b.created_at) - new Date(a.created_at)
    if (filterSort === 'oldest') return new Date(a.created_at) - new Date(b.created_at)
    return 0
  })

  const paginatedItems = filteredItems.slice((currentPage - 1) * pageSize, currentPage * pageSize)
  const totalPages = Math.ceil(filteredItems.length / pageSize) || 1

  const openForm = (item = null) => {
    if (item) {
      setFormMode('edit')
      setFormData({
        id: item.id,
        title: item.title,
        category: item.category || 'Kegiatan',
        organizer: item.organizer || 'DPP GRADASI',
        author: item.author_name || 'Super Admin',
        eventDate: item.event_date || '',
        location: item.location || '',
        image: item.image_url || '',
        content: item.content || '',
        tags: Array.isArray(item.tags) ? item.tags.join(', ') : (item.tags || ''),
        isPublished: item.is_published,
        gallery: item.gallery || []
      })
    } else {
      setFormMode('create')
      setFormData({
        id: null,
        title: '',
        category: 'Kegiatan',
        organizer: 'DPP GRADASI',
        author: 'Super Admin',
        eventDate: new Date().toISOString().slice(0, 10),
        location: '',
        image: '',
        content: '',
        tags: '',
        isPublished: true,
        gallery: []
      })
    }
    setNewGalleryUrl('')
    setNewGalleryCaption('')
    setIsFormOpen(true)
  }

  const handleAddGalleryItem = () => {
    if (!newGalleryUrl) return
    const newItem = {
      image_url: newGalleryUrl,
      caption: newGalleryCaption
    }
    setFormData({
      ...formData,
      gallery: [...formData.gallery, newItem]
    })
    setNewGalleryUrl('')
    setNewGalleryCaption('')
  }

  const handleRemoveGalleryItem = (idx) => {
    setFormData({
      ...formData,
      gallery: formData.gallery.filter((_, i) => i !== idx)
    })
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    
    // Prepare payload matching Go backend
    const payload = {
      title: formData.title,
      category: formData.category,
      organizer: formData.organizer,
      author: formData.author,
      event_date: formData.eventDate,
      location: formData.location,
      image_url: formData.image,
      content: formData.content,
      tags: formData.tags,
      is_published: formData.isPublished,
      gallery_json: JSON.stringify(formData.gallery)
    }

    try {
      if (formMode === 'create') {
        await kegiatanService.create(payload)
      } else {
        await kegiatanService.update(formData.id, payload)
      }
      setIsFormOpen(false)
      loadKegiatan()
      alert('Kegiatan berhasil disimpan!')
    } catch {
      alert('Gagal menyimpan kegiatan')
    }
  }

  const handleTogglePublish = async (item) => {
    const updatedStatus = !item.is_published
    try {
      await kegiatanService.update(item.id, {
        title: item.title,
        category: item.category,
        organizer: item.organizer,
        event_date: item.event_date,
        location: item.location,
        image_url: item.image_url,
        content: item.content,
        is_published: updatedStatus
      })
      loadKegiatan()
      alert(updatedStatus ? 'Kegiatan diterbitkan!' : 'Kegiatan dijadikan draft.')
    } catch {
      alert('Gagal memperbarui status publikasi')
    }
  }

  const handleDelete = async (id) => {
    if (window.confirm('Pindahkan kegiatan ini ke sampah?')) {
      try {
        await kegiatanService.remove(id)
        loadKegiatan()
        alert('Kegiatan berhasil dipindahkan ke sampah.')
      } catch {
        alert('Gagal menghapus kegiatan')
      }
    }
  }

  const handleRestore = async (id) => {
    try {
      await kegiatanService.restore(id)
      loadKegiatan()
      alert('Kegiatan berhasil dipulihkan.')
    } catch {
      alert('Gagal memulihkan kegiatan')
    }
  }

  return (
    <AdminLayout title="Kelola Kegiatan">
      <div className="space-y-6">
        
        {/* Navigation Tabs */}
        <div className="flex justify-between items-center bg-white p-4 rounded-xl border border-gray-200">
          <div className="flex gap-2">
            <button 
              onClick={() => { setCurrentTab('active'); setCurrentPage(1); }}
              className={`px-4 py-2 text-xs font-semibold rounded-lg ${currentTab === 'active' ? 'bg-brand-600 text-white' : 'bg-gray-100 text-slate-700'}`}
            >
              Kegiatan Aktif
            </button>
            <button 
              onClick={() => { setCurrentTab('trash'); setCurrentPage(1); }}
              className={`px-4 py-2 text-xs font-semibold rounded-lg ${currentTab === 'trash' ? 'bg-brand-600 text-white' : 'bg-gray-100 text-slate-700'}`}
            >
              Sampah
            </button>
          </div>
          {currentTab === 'active' && (
            <button 
              onClick={() => openForm()}
              className="px-4 py-2 bg-brand-600 text-white rounded-lg hover:bg-brand-700 text-xs font-semibold"
            >
              + Tambah Kegiatan Baru
            </button>
          )}
        </div>

        {/* Filter Bar */}
        <div className="flex flex-col md:flex-row gap-4 justify-between items-center bg-white p-4 rounded-xl border border-gray-200">
          <input 
            type="text" 
            placeholder="Cari kegiatan atau lokasi..." 
            value={searchQuery}
            onChange={e => { setSearchQuery(e.target.value); setCurrentPage(1); }}
            className="w-full md:max-w-md px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none"
          />
          <div className="flex gap-2 w-full md:w-auto">
            <select 
              value={filterStatus}
              onChange={e => { setFilterStatus(e.target.value); setCurrentPage(1); }}
              className="px-3 py-2 border border-gray-300 rounded-lg text-sm text-slate-600 outline-none w-full sm:w-auto"
            >
              <option value="">Semua Status</option>
              <option value="published">Diterbitkan</option>
              <option value="draft">Draft</option>
            </select>
            <select 
              value={filterSort}
              onChange={e => { setFilterSort(e.target.value); setCurrentPage(1); }}
              className="px-3 py-2 border border-gray-300 rounded-lg text-sm text-slate-600 outline-none w-full sm:w-auto"
            >
              <option value="newest">Terbaru</option>
              <option value="oldest">Terlama</option>
            </select>
          </div>
        </div>

        {loading && <div className="text-slate-500 py-10 text-center">Memuat kegiatan...</div>}
        {error && <div className="text-red-600 py-10 text-center font-medium">{error}</div>}

        {!loading && !error && (
          <div className="bg-white rounded-2xl border border-gray-200 overflow-hidden shadow-sm">
            <table className="w-full text-left text-sm text-slate-700">
              <thead className="bg-slate-50 border-b border-gray-200 font-semibold">
                <tr>
                  <th className="p-4">Judul Kegiatan</th>
                  <th className="p-4">Tanggal Event</th>
                  <th className="p-4">Lokasi</th>
                  <th className="p-4">Status</th>
                  <th className="p-4">Aksi</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-100">
                {paginatedItems.map(item => (
                  <tr key={item.id} className="hover:bg-slate-50/50">
                    <td className="p-4 font-medium text-slate-900">{item.title}</td>
                    <td className="p-4">{item.event_date}</td>
                    <td className="p-4">{item.location}</td>
                    <td className="p-4">
                      <span className={`px-2.5 py-0.5 rounded-full text-xs font-semibold ${item.is_published ? 'bg-emerald-50 text-emerald-700' : 'bg-amber-50 text-amber-700'}`}>
                        {item.is_published ? 'Diterbitkan' : 'Draft'}
                      </span>
                    </td>
                    <td className="p-4 flex gap-2">
                      {currentTab === 'active' ? (
                        <>
                          <button onClick={() => handleTogglePublish(item)} className="text-xs text-brand-600 font-semibold hover:underline">
                            {item.is_published ? 'Jadikan Draft' : 'Terbitkan'}
                          </button>
                          <button onClick={() => openForm(item)} className="text-xs text-slate-600 font-semibold hover:underline">Edit</button>
                          <button onClick={() => handleDelete(item.id)} className="text-xs text-red-600 font-semibold hover:underline">Hapus</button>
                        </>
                      ) : (
                        <button onClick={() => handleRestore(item.id)} className="text-xs text-emerald-600 font-semibold hover:underline">Pulihkan</button>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        {/* Modal Form */}
        {isFormOpen && (
          <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
            <div className="bg-white rounded-2xl p-6 max-w-2xl w-full shadow-2xl max-h-[90vh] overflow-y-auto space-y-4">
              <h3 className="font-heading font-bold text-lg text-slate-900">{formMode === 'create' ? 'Tambah Kegiatan Baru' : 'Edit Kegiatan'}</h3>
              
              <form onSubmit={handleSubmit} className="space-y-4">
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Judul Kegiatan</label>
                  <input type="text" value={formData.title} onChange={e => setFormData({...formData, title: e.target.value})} required className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Kategori</label>
                    <input type="text" value={formData.category} onChange={e => setFormData({...formData, category: e.target.value})} className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                  </div>
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Penyelenggara</label>
                    <input type="text" value={formData.organizer} onChange={e => setFormData({...formData, organizer: e.target.value})} className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                  </div>
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Tanggal Event</label>
                    <input type="text" value={formData.eventDate} onChange={e => setFormData({...formData, eventDate: e.target.value})} placeholder="Misal: 10 Okt 2026" className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                  </div>
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Lokasi</label>
                    <input type="text" value={formData.location} onChange={e => setFormData({...formData, location: e.target.value})} placeholder="Gedung XYZ, Jakarta" className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                  </div>
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">URL Gambar Utama</label>
                    <input type="text" value={formData.image} onChange={e => setFormData({...formData, image: e.target.value})} placeholder="https://example.com/cover.jpg" className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                  </div>
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Tags (Pisahkan dengan koma)</label>
                    <input type="text" value={formData.tags} onChange={e => setFormData({...formData, tags: e.target.value})} placeholder="webinar, it, umkm" className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                  </div>
                </div>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Detail Konten Kegiatan (HTML atau Teks)</label>
                  <textarea rows="4" value={formData.content} onChange={e => setFormData({...formData, content: e.target.value})} className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none font-mono" />
                </div>

                {/* Gallery Manage */}
                <div className="border border-slate-200 p-4 rounded-xl space-y-2">
                  <span className="block text-xs font-bold text-slate-700">Galeri Foto Kegiatan</span>
                  
                  {/* Gallery List */}
                  {formData.gallery.length > 0 && (
                    <div className="grid grid-cols-2 gap-2 max-h-[150px] overflow-y-auto pb-2 border-b border-slate-100">
                      {formData.gallery.map((img, idx) => (
                        <div key={idx} className="flex items-center gap-2 bg-slate-50 border p-1 rounded-lg">
                          <img src={img.image_url} alt="" className="w-10 h-10 object-cover rounded-md" />
                          <div className="flex-grow min-w-0">
                            <p className="text-[10px] text-slate-500 truncate">{img.caption || 'Tanpa Keterangan'}</p>
                          </div>
                          <button type="button" onClick={() => handleRemoveGalleryItem(idx)} className="text-red-500 text-xs px-2 hover:bg-red-50 rounded">✕</button>
                        </div>
                      ))}
                    </div>
                  )}

                  {/* Add Gallery Item Form */}
                  <div className="flex gap-2">
                    <input type="text" value={newGalleryUrl} onChange={e => setNewGalleryUrl(e.target.value)} placeholder="URL Foto Galeri" className="flex-grow px-3 py-1.5 border border-gray-300 rounded-lg text-xs outline-none" />
                    <input type="text" value={newGalleryCaption} onChange={e => setNewGalleryCaption(e.target.value)} placeholder="Keterangan" className="px-3 py-1.5 border border-gray-300 rounded-lg text-xs outline-none w-28" />
                    <button type="button" onClick={handleAddGalleryItem} className="px-3 py-1.5 bg-slate-100 text-slate-700 rounded-lg text-xs font-semibold hover:bg-slate-200">Tambah</button>
                  </div>
                </div>

                <label className="flex items-center gap-2">
                  <input type="checkbox" checked={formData.isPublished} onChange={e => setFormData({...formData, isPublished: e.target.checked})} />
                  <span className="text-xs text-gray-600 font-semibold">Terbitkan Kegiatan (Public)</span>
                </label>

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
