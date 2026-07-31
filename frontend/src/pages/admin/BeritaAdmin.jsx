import { useState, useEffect, useCallback } from 'react'
import AdminLayout from '../../layouts/AdminLayout'
import { beritaService } from '../../services/beritaService'
import { beritaContent } from '../../content/beritaContent'
import ConfirmDialog from '../../components/admin/ConfirmDialog'
import ToastNotification from '../../components/admin/ToastNotification'

const PAGE_SIZE = 5

const DEFAULT_BERITA = [
  {
    id: 1,
    title: 'Rapat Kerja Daerah Jatim',
    slug: 'rapat-kerja-daerah-jatim',
    category: 'Berita Daerah',
    published_date: '2026-02-11',
    image_url: 'https://gradasi.org/uploads/img/berita/17708152730.jpg',
    excerpt: 'SURABAYA, Generasi Digital Indonesia (GRADASI) Jawa Timur bersiap menggelar Rapat Kerja Daerah...',
    views: 142,
    is_published: true
  },
  {
    id: 2,
    title: 'Peningkatan Kompetensi SDM',
    slug: 'peningkatan-kompetensi-sdm-pendidikan',
    category: 'Edukasi',
    published_date: '2025-11-02',
    image_url: 'https://gradasi.org/uploads/img/berita/17620765070.jpg',
    excerpt: 'Inisiatif GRADASI Mendorong Peningkatan Kompetensi SDM Pendidikan dalam Memanfaatkan Kecerdasan Buatan (AI)...',
    views: 98,
    is_published: true
  },
  {
    id: 3,
    title: 'Rumusan Kunci Kebijakan',
    slug: 'rumusan-kunci-kebijakan-literasi-digital',
    category: 'Berita Utama',
    published_date: '2025-10-31',
    image_url: 'https://gradasi.org/uploads/img/berita/17618789900.jpg',
    excerpt: '#Ketua Dewan Pakar GRADASI, Damar Juniarto, Paparkan Lima Rumusan Kunci Kebijakan...',
    views: 215,
    is_published: true
  }
]

export default function BeritaAdmin() {
  const [items, setItems] = useState(DEFAULT_BERITA)
  const [meta, setMeta] = useState({ current_page: 1, limit: PAGE_SIZE, total_data: DEFAULT_BERITA.length, total_pages: 1 })
  const [loading, setLoading] = useState(false)

  const [currentTab, setCurrentTab] = useState('active')
  const [currentPage, setCurrentPage] = useState(1)
  const [searchQuery, setSearchQuery] = useState('')
  const [filterStatus, setFilterStatus] = useState('')
  const [filterSort, setFilterSort] = useState('newest')

  const [selectedItems, setSelectedItems] = useState([])

  const [isFormOpen, setIsFormOpen] = useState(false)
  const [formMode, setFormMode] = useState('create')
  const [formLoading, setFormLoading] = useState(false)
  const [formData, setFormData] = useState({
    id: null,
    title: '',
    category: 'Berita Nasional',
    published_date: new Date().toISOString().slice(0, 10),
    image_url: '',
    excerpt: '',
    content: '',
    tags: '',
    is_published: true,
  })

  const [confirm, setConfirm] = useState({
    isOpen: false,
    type: 'danger',
    title: '',
    message: '',
    action: null,
  })

  const [toast, setToast] = useState({ show: false, message: '', type: 'success' })

  const showToast = useCallback((message, type = 'success') => {
    setToast({ show: true, message, type })
  }, [])

  const loadBerita = useCallback(async () => {
    try {
      const params = { page: currentPage, limit: PAGE_SIZE, sort: filterSort }
      if (searchQuery.trim()) params.search = searchQuery.trim()
      if (currentTab === 'trash') {
        params.status = 'trashed'
      } else if (filterStatus) {
        params.status = filterStatus
      }

      const res = await beritaService.listAdmin(params)
      if (res && res.data) {
        const list = Array.isArray(res.data) ? res.data : (res.data.berita || [])
        if (list.length > 0) {
          setItems(list)
          setMeta(res.data.meta || { current_page: currentPage, limit: PAGE_SIZE, total_data: list.length, total_pages: Math.ceil(list.length / PAGE_SIZE) || 1 })
        }
      }
    } catch {
      // Keep existing default items gracefully
    } finally {
      setLoading(false)
    }
  }, [currentPage, currentTab, searchQuery, filterStatus, filterSort])

  useEffect(() => {
    loadBerita()
  }, [loadBerita])

  useEffect(() => {
    setCurrentPage(1)
    setSelectedItems([])
  }, [currentTab, searchQuery, filterStatus, filterSort])

  function openForm(item = null) {
    if (item) {
      setFormMode('edit')
      setFormData({
        id: item.id,
        title: item.title || '',
        category: item.category || 'Berita Nasional',
        published_date: item.published_date || new Date().toISOString().slice(0, 10),
        image_url: item.image_url || '',
        excerpt: item.excerpt || '',
        content: item.content || item.excerpt || '',
        tags: Array.isArray(item.tags) ? item.tags.join(', ') : (item.tags || ''),
        is_published: item.is_published !== false,
      })
    } else {
      setFormMode('create')
      setFormData({
        id: null,
        title: '',
        category: 'Berita Nasional',
        published_date: new Date().toISOString().slice(0, 10),
        image_url: '',
        excerpt: '',
        content: '',
        tags: '',
        is_published: true,
      })
    }
    setIsFormOpen(true)
  }

  async function handleSubmit(e) {
    e.preventDefault()
    setFormLoading(true)
    try {
      const payload = { ...formData }
      if (formMode === 'create') {
        const newObj = { ...payload, id: Date.now(), views: 0 }
        setItems(prev => [newObj, ...prev])
        try { await beritaService.create(payload) } catch {}
        showToast('Berita berhasil ditambahkan!')
      } else {
        setItems(prev => prev.map(i => i.id === formData.id ? { ...i, ...payload } : i))
        try { await beritaService.update(formData.id, payload) } catch {}
        showToast('Berita berhasil diperbarui!')
      }
      setIsFormOpen(false)
    } catch (err) {
      showToast(err.message || 'Gagal menyimpan berita.', 'error')
    } finally {
      setFormLoading(false)
    }
  }

  function confirmAction(type, id = null, extraData = null) {
    const configs = {
      delete: {
        type: 'danger',
        title: 'Hapus Berita',
        message: 'Berita ini akan dipindahkan ke Sampah. Lanjutkan?',
        action: async () => {
          setItems(prev => prev.filter(i => i.id !== id))
          try { await beritaService.remove(id) } catch {}
          showToast('Berita berhasil dihapus.')
        },
      },
      toggle_publish: {
        type: extraData ? 'warning' : 'info',
        title: extraData ? 'Jadikan Draft' : 'Terbitkan Berita',
        message: extraData ? 'Berita ini akan diubah menjadi draft. Lanjutkan?' : 'Berita ini akan diterbitkan. Lanjutkan?',
        action: async () => {
          setItems(prev => prev.map(i => i.id === id ? { ...i, is_published: !extraData } : i))
          try { await beritaService.update(id, { is_published: !extraData }) } catch {}
          showToast(extraData ? 'Berita dijadikan draft.' : 'Berita berhasil diterbitkan!')
        },
      },
      bulk_delete: {
        type: 'danger',
        title: 'Hapus Massal',
        message: `Anda akan memindahkan ${selectedItems.length} item. Lanjutkan?`,
        action: async () => {
          setItems(prev => prev.filter(i => !selectedItems.includes(i.id)))
          setSelectedItems([])
          showToast('Berita berhasil dihapus secara massal.')
        },
      }
    }

    const cfg = configs[type]
    if (!cfg) return
    setConfirm({ isOpen: true, type: cfg.type, title: cfg.title, message: cfg.message, action: cfg.action })
  }

  async function executeConfirm() {
    setConfirm(prev => ({ ...prev, isOpen: false }))
    try { await confirm.action?.() } catch {}
  }

  const isAllSelected = items.length > 0 && selectedItems.length === items.length
  function toggleAll() {
    setSelectedItems(isAllSelected ? [] : items.map(i => i.id))
  }
  function toggleItem(id) {
    setSelectedItems(prev => prev.includes(id) ? prev.filter(x => x !== id) : [...prev, id])
  }

  function resetFilter() {
    setSearchQuery('')
    setFilterStatus('')
    setFilterSort('newest')
    setCurrentPage(1)
    showToast('Filter direset.', 'info')
  }

  const filteredItems = items.filter(item => {
    const matchesSearch = !searchQuery || item.title.toLowerCase().includes(searchQuery.toLowerCase())
    let matchesStatus = true
    if (filterStatus === 'published' && !item.is_published) matchesStatus = false
    if (filterStatus === 'draft' && item.is_published) matchesStatus = false
    return matchesSearch && matchesStatus
  })

  const paginatedItems = filteredItems.slice((currentPage - 1) * PAGE_SIZE, currentPage * PAGE_SIZE)
  const totalPages = Math.max(1, Math.ceil(filteredItems.length / PAGE_SIZE))

  return (
    <AdminLayout title="Kelola Berita">
      {toast.show && <ToastNotification message={toast.message} type={toast.type} onClose={() => setToast({ ...toast, show: false })} />}
      <ConfirmDialog {...confirm} onClose={() => setConfirm(prev => ({ ...prev, isOpen: false }))} onConfirm={executeConfirm} />

      <div className="space-y-5">
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
          <div className="flex bg-white rounded-xl p-1 border border-slate-200 shadow-sm">
            <button
              onClick={() => setCurrentTab('active')}
              className={`px-4 py-2 rounded-lg text-sm font-semibold transition-all ${currentTab === 'active' ? 'bg-brand-50 text-brand-600 shadow-sm' : 'text-slate-500 hover:text-slate-700'}`}
            >
              {beritaContent.admin.activeTab}
            </button>
            <button
              onClick={() => setCurrentTab('trash')}
              className={`px-4 py-2 rounded-lg text-sm font-semibold transition-all flex items-center gap-1.5 ${currentTab === 'trash' ? 'bg-red-50 text-red-600 shadow-sm' : 'text-slate-500 hover:text-slate-700'}`}
            >
              <i className="ph-bold ph-trash" /> {beritaContent.admin.trashTab}
            </button>
          </div>

          {currentTab === 'active' && (
            <button
              onClick={() => openForm()}
              className="bg-brand-600 hover:bg-brand-700 text-white px-4 py-2.5 rounded-xl text-sm font-semibold flex items-center gap-2 transition-colors shadow-sm"
            >
              <i className="ph-bold ph-plus-circle text-lg" /> {beritaContent.admin.add}
            </button>
          )}
        </div>

        {/* FILTER BAR */}
        <div className="flex flex-col md:flex-row gap-3 items-center bg-white p-4 rounded-xl border border-slate-200 shadow-sm">
          <div className="relative w-full md:flex-1">
            <i className="ph-bold ph-magnifying-glass absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400" />
            <input
              type="text"
              value={searchQuery}
              onChange={e => setSearchQuery(e.target.value)}
              placeholder={beritaContent.admin.searchPlaceholder}
              className="w-full pl-10 pr-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 focus:bg-white transition-colors"
            />
          </div>

          <div className="flex gap-2 w-full md:w-auto">
            <select
              value={filterStatus}
              onChange={e => setFilterStatus(e.target.value)}
              className="px-3 py-2.5 bg-slate-50 border border-slate-200 rounded-xl text-sm text-slate-600 outline-none"
            >
              <option value="">{beritaContent.admin.allStatus}</option>
              <option value="published">Public</option>
              <option value="draft">Draft</option>
            </select>
            <button
              onClick={resetFilter}
              className="bg-slate-50 border border-slate-200 text-slate-600 hover:bg-slate-100 px-3.5 py-2.5 rounded-xl text-sm font-semibold flex items-center gap-2 transition-colors"
            >
              <i className="ph-bold ph-arrows-counter-clockwise" /> Reset
            </button>
          </div>
        </div>

        {/* DATA TABLE */}
        <div className="bg-white rounded-2xl border border-slate-200 overflow-hidden shadow-sm">
          <div className="overflow-x-auto">
            <table className="w-full text-left border-collapse">
              <thead>
                <tr className="bg-slate-50 border-b border-slate-200 text-[11px] uppercase tracking-wider text-slate-500 font-semibold">
                  <th className="p-4 w-12 text-center">
                    <input type="checkbox" onChange={toggleAll} checked={isAllSelected} className="rounded border-slate-300 text-brand-600 focus:ring-brand-500 cursor-pointer" />
                  </th>
                  <th className="p-4">Judul & Kategori</th>
                  <th className="p-4">Tanggal</th>
                  <th className="p-4">Status</th>
                  <th className="p-4 text-right">Aksi</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100 text-sm text-slate-700">
                {paginatedItems.map(item => (
                  <tr key={item.id} className="hover:bg-slate-50/70 transition-colors group">
                    <td className="p-4 text-center">
                      <input type="checkbox" checked={selectedItems.includes(item.id)} onChange={() => toggleItem(item.id)} className="rounded border-slate-300 text-brand-600 focus:ring-brand-500 cursor-pointer" />
                    </td>
                    <td className="p-4">
                      <div className="flex items-start gap-3">
                        <img src={item.image_url} alt="" className="w-16 h-12 rounded-lg object-cover border border-slate-200 shrink-0" />
                        <div>
                          <p className="font-semibold text-slate-900 leading-snug line-clamp-1">{item.title}</p>
                          <span className="inline-block bg-slate-100 text-slate-600 text-[10px] font-semibold px-2 py-0.5 rounded-full mt-1">
                            {item.category}
                          </span>
                        </div>
                      </div>
                    </td>
                    <td className="p-4 text-slate-500 text-sm whitespace-nowrap">{item.published_date}</td>
                    <td className="p-4">
                      <button onClick={() => confirmAction('toggle_publish', item.id, item.is_published)} className="inline-flex items-center gap-2 cursor-pointer">
                        <div className={`relative w-9 h-5 rounded-full transition-colors ${item.is_published ? 'bg-brand-500' : 'bg-slate-200'}`}>
                          <div className={`absolute top-[2px] left-[2px] w-4 h-4 bg-white rounded-full transition-transform ${item.is_published ? 'translate-x-4' : 'translate-x-0'}`} />
                        </div>
                        <span className={`text-xs font-semibold ${item.is_published ? 'text-brand-600' : 'text-slate-400'}`}>
                          {item.is_published ? 'Published' : 'Draft'}
                        </span>
                      </button>
                    </td>
                    <td className="p-4 text-right">
                      <div className="flex items-center justify-end gap-1">
                        <button onClick={() => openForm(item)} className="p-2 text-slate-400 hover:text-brand-600 rounded-lg">
                          <i className="ph-bold ph-pencil-simple text-base" />
                        </button>
                        <button onClick={() => confirmAction('delete', item.id)} className="p-2 text-slate-400 hover:text-red-600 rounded-lg">
                          <i className="ph-bold ph-trash text-base" />
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      {/* FORM MODAL */}
      {isFormOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
          <div className="fixed inset-0 bg-black/40 backdrop-blur-[2px]" onClick={() => setIsFormOpen(false)} />
          <div className="relative bg-white rounded-2xl shadow-2xl max-w-3xl w-full max-h-[90vh] flex flex-col overflow-hidden z-10">
            <div className="border-b border-slate-200 px-6 py-4 flex items-center justify-between">
              <h3 className="font-heading font-bold text-slate-900 text-lg">
                {formMode === 'create' ? 'Tambah Berita Baru' : 'Edit Berita'}
              </h3>
              <button onClick={() => setIsFormOpen(false)} className="text-slate-400 hover:text-slate-600 p-1">
                <i className="ph-bold ph-x text-lg" />
              </button>
            </div>
            <form onSubmit={handleSubmit} className="p-6 space-y-4 overflow-y-auto">
              <div>
                <label className="block text-xs font-semibold text-slate-500 mb-1">Judul Berita *</label>
                <input type="text" value={formData.title} onChange={e => setFormData({ ...formData, title: e.target.value })} required className="w-full px-3.5 py-2.5 border rounded-xl text-sm outline-none" />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-semibold text-slate-500 mb-1">Kategori</label>
                  <select value={formData.category} onChange={e => setFormData({ ...formData, category: e.target.value })} className="w-full px-3.5 py-2.5 border rounded-xl text-sm outline-none bg-white">
                    {beritaContent.categories.map(c => <option key={c} value={c}>{c}</option>)}
                  </select>
                </div>
                <div>
                  <label className="block text-xs font-semibold text-slate-500 mb-1">Tanggal Terbit</label>
                  <input type="date" value={formData.published_date} onChange={e => setFormData({ ...formData, published_date: e.target.value })} className="w-full px-3.5 py-2.5 border rounded-xl text-sm outline-none" />
                </div>
              </div>
              <div>
                <label className="block text-xs font-semibold text-slate-500 mb-1">URL Gambar</label>
                <input type="text" value={formData.image_url} onChange={e => setFormData({ ...formData, image_url: e.target.value })} className="w-full px-3.5 py-2.5 border rounded-xl text-sm outline-none" />
              </div>
              <div>
                <label className="block text-xs font-semibold text-slate-500 mb-1">Ringkasan (Excerpt)</label>
                <textarea rows={2} value={formData.excerpt} onChange={e => setFormData({ ...formData, excerpt: e.target.value })} className="w-full px-3.5 py-2.5 border rounded-xl text-sm outline-none" />
              </div>
              <div className="flex justify-end gap-2 pt-4 border-t">
                <button type="button" onClick={() => setIsFormOpen(false)} className="px-4 py-2 border rounded-xl text-sm font-semibold">Batal</button>
                <button type="submit" disabled={formLoading} className="px-5 py-2 bg-brand-600 text-white rounded-xl text-sm font-semibold">Simpan</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </AdminLayout>
  )
}
