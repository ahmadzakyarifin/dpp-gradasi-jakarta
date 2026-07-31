import { useState, useEffect } from 'react'
import AdminLayout from '../../layouts/AdminLayout'
import { beritaService } from '../../services/beritaService'

export default function BeritaAdmin() {
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
    category: 'Berita Nasional',
    image_url: '',
    excerpt: '',
    content: '',
    tags: '',
    is_published: true
  })

  const loadBerita = () => {
    setLoading(true)
    const serviceCall = currentTab === 'active' 
      ? beritaService.list({ active_only: true }) 
      : beritaService.list({ trash_only: true }) // Adjust parameter based on backend contract
      
    serviceCall
      .then(res => {
        if (res.success && res.data) {
          // If backend returns data wrapped in an object or array
          const beritaList = Array.isArray(res.data) ? res.data : (res.data.berita || [])
          setItems(beritaList)
        } else {
          setError('Gagal memuat berita')
        }
      })
      .catch(() => setError('Kesalahan koneksi ke server'))
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    loadBerita()
  }, [currentTab])

  // Filter & Sort
  const filteredItems = items.filter(item => {
    const matchesSearch = item.title.toLowerCase().includes(searchQuery.toLowerCase())
    
    let matchesStatus = true
    if (filterStatus === 'published' && !item.is_published) matchesStatus = false
    if (filterStatus === 'draft' && item.is_published) matchesStatus = false

    return matchesSearch && matchesStatus
  }).sort((a, b) => {
    if (filterSort === 'newest') return new Date(b.created_at) - new Date(a.created_at)
    if (filterSort === 'oldest') return new Date(a.created_at) - new Date(b.created_at)
    if (filterSort === 'most_viewed') return (b.views || 0) - (a.views || 0)
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
        category: item.category || 'Berita Nasional',
        image_url: item.image_url || '',
        excerpt: item.excerpt || '',
        content: item.content || '',
        tags: Array.isArray(item.tags) ? item.tags.join(', ') : (item.tags || ''),
        is_published: item.is_published
      })
    } else {
      setFormMode('create')
      setFormData({
        id: null,
        title: '',
        category: 'Berita Nasional',
        image_url: '',
        excerpt: '',
        content: '',
        tags: '',
        is_published: true
      })
    }
    setIsFormOpen(true)
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    try {
      if (formMode === 'create') {
        await beritaService.create(formData)
      } else {
        await beritaService.update(formData.id, formData)
      }
      setIsFormOpen(false)
      loadBerita()
      alert('Berita berhasil disimpan!')
    } catch {
      alert('Gagal menyimpan berita')
    }
  }

  const handleTogglePublish = async (item) => {
    const updatedStatus = !item.is_published
    try {
      await beritaService.update(item.id, {
        ...item,
        is_published: updatedStatus
      })
      loadBerita()
      alert(updatedStatus ? 'Berita diterbitkan!' : 'Berita dijadikan draft.')
    } catch {
      alert('Gagal memperbarui status publikasi')
    }
  }

  const handleDelete = async (id) => {
    if (window.confirm('Pindahkan berita ini ke sampah?')) {
      try {
        await beritaService.remove(id)
        loadBerita()
        alert('Berita berhasil dipindahkan ke sampah.')
      } catch {
        alert('Gagal menghapus berita')
      }
    }
  }

  const handleRestore = async (id) => {
    try {
      await beritaService.restore(id)
      loadBerita()
      alert('Berita berhasil dipulihkan.')
    } catch {
      alert('Gagal memulihkan berita')
    }
  }

  return (
    <AdminLayout title="Kelola Berita">
      <div className="space-y-6">
        
        {/* Navigation Tabs */}
        <div className="flex justify-between items-center bg-white p-4 rounded-xl border border-gray-200">
          <div className="flex gap-2">
            <button 
              onClick={() => { setCurrentTab('active'); setCurrentPage(1); }}
              className={`px-4 py-2 text-xs font-semibold rounded-lg ${currentTab === 'active' ? 'bg-brand-600 text-white' : 'bg-gray-100 text-slate-700'}`}
            >
              Berita Aktif
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
              + Tambah Berita Baru
            </button>
          )}
        </div>

        {/* Filter Bar */}
        <div className="flex flex-col md:flex-row gap-4 justify-between items-center bg-white p-4 rounded-xl border border-gray-200">
          <input 
            type="text" 
            placeholder="Cari berita..." 
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
              <option value="published">Diterbitkan (Public)</option>
              <option value="draft">Draft</option>
            </select>
            <select 
              value={filterSort}
              onChange={e => { setFilterSort(e.target.value); setCurrentPage(1); }}
              className="px-3 py-2 border border-gray-300 rounded-lg text-sm text-slate-600 outline-none w-full sm:w-auto"
            >
              <option value="newest">Terbaru</option>
              <option value="oldest">Terlama</option>
              <option value="most_viewed">Terpopuler</option>
            </select>
          </div>
        </div>

        {loading && <div className="text-slate-500 py-10 text-center">Memuat berita...</div>}
        {error && <div className="text-red-600 py-10 text-center font-medium">{error}</div>}

        {!loading && !error && (
          <div className="bg-white rounded-2xl border border-gray-200 overflow-hidden shadow-sm">
            <table className="w-full text-left text-sm text-slate-700">
              <thead className="bg-slate-50 border-b border-gray-200 font-semibold">
                <tr>
                  <th className="p-4">Judul Berita</th>
                  <th className="p-4">Kategori</th>
                  <th className="p-4">Pembaca (Views)</th>
                  <th className="p-4">Status</th>
                  <th className="p-4">Aksi</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-100">
                {paginatedItems.map(item => (
                  <tr key={item.id} className="hover:bg-slate-50/50">
                    <td className="p-4 font-medium text-slate-900">{item.title}</td>
                    <td className="p-4">{item.category}</td>
                    <td className="p-4">{item.views || 0}</td>
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

        {/* Form Modal */}
        {isFormOpen && (
          <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
            <div className="bg-white rounded-2xl p-6 max-w-2xl w-full shadow-2xl max-h-[90vh] overflow-y-auto space-y-4">
              <h3 className="font-heading font-bold text-lg text-slate-900">{formMode === 'create' ? 'Tambah Berita Baru' : 'Edit Berita'}</h3>
              
              <form onSubmit={handleSubmit} className="space-y-4">
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Judul Berita</label>
                  <input type="text" value={formData.title} onChange={e => setFormData({...formData, title: e.target.value})} required className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Kategori</label>
                    <input type="text" value={formData.category} onChange={e => setFormData({...formData, category: e.target.value})} className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                  </div>
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Tags (Pisahkan dengan koma)</label>
                    <input type="text" value={formData.tags} onChange={e => setFormData({...formData, tags: e.target.value})} placeholder="digital, literasi, gradasi" className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                  </div>
                </div>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">URL Gambar Cover</label>
                  <input type="text" value={formData.image_url} onChange={e => setFormData({...formData, image_url: e.target.value})} placeholder="https://example.com/cover.jpg" className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Ringkasan (Excerpt)</label>
                  <textarea rows="2" value={formData.excerpt} onChange={e => setFormData({...formData, excerpt: e.target.value})} className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Isi Berita Lengkap (HTML atau Teks)</label>
                  <textarea rows="6" value={formData.content} onChange={e => setFormData({...formData, content: e.target.value})} className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none font-mono" />
                </div>
                <label className="flex items-center gap-2">
                  <input type="checkbox" checked={formData.is_published} onChange={e => setFormData({...formData, is_published: e.target.checked})} />
                  <span className="text-xs text-gray-600 font-semibold">Terbitkan Langsung (Public)</span>
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
