import { useState, useEffect, useCallback } from 'react'
import AdminLayout from '../../layouts/AdminLayout'
import { kegiatanService } from '../../services/kegiatanService'
import ConfirmDialog from '../../components/admin/ConfirmDialog'
import ToastNotification from '../../components/admin/ToastNotification'

const PAGE_SIZE = 5

export default function KegiatanAdmin() {
  const [items, setItems] = useState([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)
  const [meta, setMeta] = useState({ total_data: 0, total_pages: 1, current_page: 1, limit: 10 })

  const [currentTab, setCurrentTab] = useState('active')
  const [currentPage, setCurrentPage] = useState(1)

  const [searchQuery, setSearchQuery] = useState('')
  const [filterStatus, setFilterStatus] = useState('')
  const [filterSort, setFilterSort] = useState('newest')

  const [selectedItems, setSelectedItems] = useState([])

  const [isFormOpen, setIsFormOpen] = useState(false)
  const [formMode, setFormMode] = useState('create')

  const [categories, setCategories] = useState([])

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

  const loadKegiatan = useCallback(() => {
    setLoading(true)
    const params = {
      page: currentPage,
      limit: PAGE_SIZE,
      search: searchQuery,
      status: currentTab === 'trash' ? 'trashed' : (filterStatus || undefined),
      sort: filterSort || undefined,
    }
    kegiatanService.listAdmin(params)
      .then(res => {
        if (res && res.data) {
          const list = Array.isArray(res.data) ? res.data : (res.data.kegiatan || [])
          setItems(list)
          if (res.data.meta) setMeta(res.data.meta)
        }
      })
      .catch(err => {
        setError(err?.message || 'Gagal memuat kegiatan')
        setItems([])
      })
      .finally(() => setLoading(false))
  }, [currentTab, currentPage, searchQuery, filterStatus, filterSort])

  useEffect(() => {
    loadKegiatan()
  }, [loadKegiatan])

  useEffect(() => {
    setCurrentPage(1)
    setSelectedItems([])
  }, [currentTab, searchQuery, filterStatus, filterSort])

  // Ambil daftar kategori dinamis (dari data yang ada) untuk dropdown
  useEffect(() => {
    kegiatanService.getCategories()
      .then(res => {
        if (res && res.data && Array.isArray(res.data) && res.data.length > 0) {
          setCategories(res.data)
        }
      }).catch(() => {})
  }, [])

  // Server sudah filter & paginate — render langsung
  const paginatedItems = items
  const totalPages = Math.max(1, meta.total_pages || 1)
  const totalData = meta.total_data ?? items.length

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
        image: item.image_url || item.image_path || '',
        excerpt: item.excerpt || '',
        content: item.content || item.excerpt || '',
        isPublished: item.is_published !== false
      })
      // list_admin tidak menyertakan content/tags/gallery — ambil detail penuh dulu
      kegiatanService.detailById(item.id)
        .then(res => {
          if (res?.data) {
            const d = res.data
            setFormData(prev => ({
              ...prev,
              content: d.content ?? prev.content,
              excerpt: d.excerpt ?? prev.excerpt,
              eventDate: d.event_date ?? prev.eventDate,
              location: d.location ?? prev.location,
              organizer: d.organizer ?? prev.organizer,
              image: d.image_path || d.image_url || prev.image,
            }))
          }
        })
        .catch(() => {}) // gagal ambil detail → tetap pakai data list
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
      image_path: formData.image,
      excerpt: formData.excerpt,
      content: formData.content,
      is_published: formData.isPublished
    }

    try {
      if (formMode === 'create') {
        await kegiatanService.create(payload)
        showToast('Kegiatan berhasil dibuat.')
      } else {
        await kegiatanService.update(formData.id, payload)
        showToast('Kegiatan berhasil diperbarui.')
      }
      setIsFormOpen(false)
      loadKegiatan()
    } catch (err) {
      showToast(err.message || 'Gagal menyimpan kegiatan.', 'error')
    }
  }

  function confirmAction(type, id = null) {
    const configs = {
      delete: {
        type: 'danger',
        title: 'Hapus Kegiatan',
        message: 'Kegiatan ini akan dipindahkan ke Sampah. Lanjutkan?',
        action: async () => {
          await kegiatanService.remove(id)
          loadKegiatan()
          showToast('Kegiatan berhasil dihapus.')
        },
      },
      restore: {
        type: 'info',
        title: 'Pulihkan Kegiatan',
        message: 'Kegiatan ini akan dikembalikan dari Sampah. Lanjutkan?',
        action: async () => {
          await kegiatanService.restore(id)
          loadKegiatan()
          showToast('Kegiatan berhasil dipulihkan.')
        },
      },
      bulk_delete: {
        type: 'danger',
        title: 'Hapus Massal',
        message: `Anda akan memindahkan ${selectedItems.length} item ke Sampah. Lanjutkan?`,
        action: async () => {
          await kegiatanService.bulkDelete(selectedItems)
          setSelectedItems([])
          loadKegiatan()
          showToast('Kegiatan berhasil dihapus secara massal.')
        },
      },
      bulk_restore: {
        type: 'info',
        title: 'Pulihkan Massal',
        message: `Anda akan memulihkan ${selectedItems.length} item dari Sampah. Lanjutkan?`,
        action: async () => {
          await kegiatanService.bulkRestore(selectedItems)
          setSelectedItems([])
          loadKegiatan()
          showToast('Kegiatan berhasil dipulihkan secara massal.')
        },
      }
    }

    const cfg = configs[type]
    if (!cfg) return
    setConfirm({ isOpen: true, type: cfg.type, title: cfg.title, message: cfg.message, action: cfg.action })
  }

  async function executeConfirm() {
    const action = confirm.action
    setConfirm(prev => ({ ...prev, isOpen: false, action: null }))
    if (action) {
      try { await action() } catch (err) {
        showToast(err?.message || 'Terjadi kesalahan, coba lagi.', 'error')
      }
    }
  }

  const isAllSelected = paginatedItems.length > 0 && selectedItems.length === paginatedItems.length
  function toggleAll() {
    setSelectedItems(isAllSelected ? [] : paginatedItems.map(i => i.id))
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

  const headerContent = (
    <div className="flex items-center gap-2 w-full max-w-3xl animate-fade-in-up">
      <div className="relative w-full">
        <i className="ph ph-magnifying-glass absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
        <input
          type="text"
          value={searchQuery}
          onChange={e => { setSearchQuery(e.target.value); setCurrentPage(1); }}
          placeholder="Cari kegiatan atau lokasi..."
          className="w-full pl-9 pr-4 py-2 bg-gray-50 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-brand-500 focus:bg-white transition-colors"
        />
      </div>
      <select
        value={filterStatus}
        onChange={e => { setFilterStatus(e.target.value); setCurrentPage(1); }}
        className="shrink-0 bg-gray-50 border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-brand-500 focus:bg-white transition-colors cursor-pointer"
      >
        <option value="">Semua Status</option>
        <option value="published">Terbit</option>
        <option value="draft">Draft</option>
      </select>
      <select
        value={filterSort}
        onChange={e => { setFilterSort(e.target.value); setCurrentPage(1); }}
        className="shrink-0 bg-gray-50 border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-brand-500 focus:bg-white transition-colors cursor-pointer"
      >
        <option value="newest">Terbaru</option>
        <option value="oldest">Terlama</option>
      </select>
      <button
        onClick={resetFilter}
        className="shrink-0 bg-gray-50 border border-gray-200 text-gray-700 hover:bg-gray-100 px-3 py-2 rounded-lg text-sm font-medium flex items-center gap-2 transition-all btn-press"
      >
        <i className="ph ph-arrows-counter-clockwise text-lg" /> Reset
      </button>
    </div>
  )

  return (
    <AdminLayout title="Kelola Kegiatan" headerContent={headerContent}>
      {toast.show && <ToastNotification message={toast.message} type={toast.type} onClose={() => setToast({ ...toast, show: false })} />}
      <ConfirmDialog {...confirm} onClose={() => setConfirm(prev => ({ ...prev, isOpen: false, action: null }))} onConfirm={executeConfirm} />

      <div className="space-y-6 animate-fade-in-up">
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
          <div className="flex bg-white rounded-lg p-1 border border-gray-200 shadow-sm">
            <button
              onClick={() => { setCurrentTab('active'); setCurrentPage(1); }}
              className={`px-4 py-1.5 rounded-md text-sm transition-all duration-200 ${currentTab === 'active' ? 'bg-brand-50 text-brand-600 shadow-sm font-medium' : 'text-gray-500 hover:text-gray-700'}`}
            >
              Aktif & Draft
            </button>
            <button
              onClick={() => { setCurrentTab('trash'); setCurrentPage(1); }}
              className={`px-4 py-1.5 rounded-md text-sm transition-all duration-200 flex items-center gap-2 ${currentTab === 'trash' ? 'bg-red-50 text-red-600 shadow-sm font-medium' : 'text-gray-500 hover:text-gray-700'}`}
            >
              <i className="ph ph-trash" /> Sampah (History)
            </button>
          </div>
          {currentTab === 'active' && (
            <button
              onClick={() => openForm()}
              className="bg-brand-600 hover:bg-brand-700 text-white px-4 py-2 rounded-lg text-sm font-medium flex items-center gap-2 transition-all btn-press shadow-sm"
            >
              <i className="ph ph-plus-circle text-lg" /> Tambah
            </button>
          )}
        </div>

        {/* Bulk Action Bar */}
        {selectedItems.length > 0 && (
          <div className="flex flex-wrap items-center gap-2 bg-brand-50/60 border border-brand-100 rounded-xl px-4 py-2.5">
            <span className="text-sm font-semibold text-brand-700">
              {selectedItems.length} dipilih
            </span>
            <div className="flex gap-2 ml-auto">
              {currentTab === 'trash' ? (
                <button
                  onClick={() => confirmAction('bulk_restore')}
                  className="bg-emerald-600 hover:bg-emerald-700 text-white px-3.5 py-1.5 rounded-lg text-xs font-semibold flex items-center gap-1.5 transition-colors"
                >
                  <i className="ph-bold ph-arrow-counter-clockwise" /> Pulihkan Massal
                </button>
              ) : (
                <button
                  onClick={() => confirmAction('bulk_delete')}
                  className="bg-red-600 hover:bg-red-700 text-white px-3.5 py-1.5 rounded-lg text-xs font-semibold flex items-center gap-1.5 transition-colors"
                >
                  <i className="ph-bold ph-trash" /> Hapus Massal
                </button>
              )}
              <button
                onClick={() => setSelectedItems([])}
                className="bg-white border border-slate-200 text-slate-600 hover:bg-slate-50 px-3.5 py-1.5 rounded-lg text-xs font-semibold transition-colors"
              >
                Batal
              </button>
            </div>
          </div>
        )}

        {/* Data Table */}
        <div className="bg-white rounded-2xl border border-gray-200 overflow-hidden shadow-sm">
          {loading && <div className="py-16 text-center text-slate-500">Memuat kegiatan...</div>}
          {!loading && error && (
            <div className="py-16 text-center text-red-600 font-medium">
              <i className="ph-bold ph-warning-circle text-2xl mb-2 block mx-auto" /> {error}
            </div>
          )}
          {!loading && !error && (
          <>
          <div className="overflow-x-auto">
            <table className="w-full text-left text-sm text-slate-700">
              <thead className="bg-slate-50 border-b border-gray-200 font-semibold text-xs uppercase tracking-wider text-slate-500">
                <tr>
                  <th className="p-4 w-12 text-center">
                    <input type="checkbox" onChange={toggleAll} checked={isAllSelected} className="rounded border-slate-300 text-brand-600 focus:ring-brand-500 cursor-pointer" />
                  </th>
                  <th className="p-4">Kegiatan</th>
                  <th className="p-4">Tanggal & Lokasi</th>
                  <th className="p-4">Kategori</th>
                  <th className="p-4">Status</th>
                  <th className="p-4 text-right">Aksi</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-100">
                {paginatedItems.map(item => (
                  <tr key={item.id} className="hover:bg-slate-50/60 transition admin-row">
                    <td className="p-4 text-center">
                      <input type="checkbox" checked={selectedItems.includes(item.id)} onChange={() => toggleItem(item.id)} className="rounded border-slate-300 text-brand-600 focus:ring-brand-500 cursor-pointer" />
                    </td>
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
                      {currentTab === 'trash' ? (
                        <span className="text-xs font-semibold px-2.5 py-1 rounded-lg bg-red-50 text-red-500">
                          Di Sampah
                        </span>
                      ) : (
                        <span className={`text-xs font-semibold px-2.5 py-1 rounded-lg ${item.is_published ? 'bg-emerald-50 text-emerald-600' : 'bg-slate-100 text-slate-500'}`}>
                          {item.is_published ? 'Terbit' : 'Draft'}
                        </span>
                      )}
                    </td>
                    <td className="p-4 text-right">
                      <div className="flex justify-end gap-2">
                        {currentTab === 'trash' ? (
                          <button onClick={() => confirmAction('restore', item.id)} className="p-2 text-slate-400 hover:text-emerald-600 rounded-lg" title="Pulihkan">
                            <i className="ph ph-arrow-counter-clockwise text-base" /> Pulihkan
                          </button>
                        ) : (
                          <>
                            <button onClick={() => openForm(item)} className="p-2 text-slate-400 hover:text-brand-600 rounded-lg">
                              <i className="ph ph-pencil-simple text-base" />
                            </button>
                            <button onClick={() => confirmAction('delete', item.id)} className="p-2 text-slate-400 hover:text-red-600 rounded-lg">
                              <i className="ph ph-trash text-base" />
                            </button>
                          </>
                        )}
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {/* Pagination */}
          {totalPages > 1 && (
            <div className="flex items-center justify-between px-4 py-3 border-t border-gray-200">
              <span className="text-xs text-slate-500">
                Hal {currentPage} dari {totalPages} · {totalData} data
              </span>
              <div className="flex items-center gap-1.5">
                <button
                  type="button"
                  disabled={currentPage <= 1}
                  onClick={() => setCurrentPage(currentPage - 1)}
                  className="w-8 h-8 flex items-center justify-center rounded-lg border border-slate-200 text-slate-500 hover:border-brand-500 hover:text-brand-600 disabled:opacity-40 transition"
                >
                  <i className="ph-bold ph-caret-left text-sm" />
                </button>
                {Array.from({ length: totalPages }, (_, i) => i + 1).map((n) => (
                  <button
                    key={n}
                    type="button"
                    onClick={() => setCurrentPage(n)}
                    className={`w-8 h-8 flex items-center justify-center rounded-lg text-sm font-semibold transition ${n === currentPage ? 'bg-brand-600 text-white' : 'border border-slate-200 text-slate-500 hover:border-brand-500 hover:text-brand-600'}`}
                  >
                    {n}
                  </button>
                ))}
                <button
                  type="button"
                  disabled={currentPage >= totalPages}
                  onClick={() => setCurrentPage(currentPage + 1)}
                  className="w-8 h-8 flex items-center justify-center rounded-lg border border-slate-200 text-slate-500 hover:border-brand-500 hover:text-brand-600 disabled:opacity-40 transition"
                >
                  <i className="ph-bold ph-caret-right text-sm" />
                </button>
              </div>
            </div>
          )}
          </>
          )}
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
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Kategori</label>
                  <select
                    value={formData.category}
                    onChange={e => {
                      if (e.target.value === '__new__') {
                        const name = window.prompt('Nama kategori baru:')
                        if (name && name.trim()) {
                          const clean = name.trim()
                          if (!categories.includes(clean)) {
                            setCategories(prev => [...prev, clean])
                          }
                          setFormData({ ...formData, category: clean })
                        }
                        return
                      }
                      setFormData({ ...formData, category: e.target.value })
                    }}
                    className="w-full px-3.5 py-2.5 border rounded-xl text-sm outline-none bg-white"
                  >
                    {categories.map(c => <option key={c} value={c}>{c}</option>)}
                    <option value="__new__">+ Buat Kategori Baru...</option>
                  </select>
                </div>
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
              <div>
                <label className="block text-xs font-semibold text-gray-500 mb-1">Konten Lengkap *</label>
                <textarea rows={5} value={formData.content} onChange={e => setFormData({ ...formData, content: e.target.value })} required className="w-full px-3.5 py-2.5 border rounded-xl text-sm outline-none" />
              </div>
              <div>
                <label className="block text-xs font-semibold text-gray-500 mb-2">Status Publikasi</label>
                <div className="flex gap-4">
                  <label className="flex items-center gap-2 text-sm cursor-pointer">
                    <input type="radio" name="isPublished" checked={formData.isPublished} onChange={() => setFormData({ ...formData, isPublished: true })} className="accent-brand-600" />
                    Terbit
                  </label>
                  <label className="flex items-center gap-2 text-sm cursor-pointer">
                    <input type="radio" name="isPublished" checked={!formData.isPublished} onChange={() => setFormData({ ...formData, isPublished: false })} className="accent-slate-400" />
                    Draft
                  </label>
                </div>
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
