import { useState, useEffect, useCallback } from 'react'
import AdminLayout from '../../layouts/AdminLayout'
import { beritaService } from '../../services/beritaService'
import { beritaContent } from '../../content/beritaContent'
import ConfirmDialog from '../../components/admin/ConfirmDialog'
import ToastNotification from '../../components/admin/ToastNotification'

const PAGE_SIZE = 5

export default function BeritaAdmin() {
  const [items, setItems] = useState([])
  const [meta, setMeta] = useState({ current_page: 1, limit: PAGE_SIZE, total_data: 0, total_pages: 1 })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)

  const [currentTab, setCurrentTab] = useState('active')
  const [currentPage, setCurrentPage] = useState(1)
  const [searchQuery, setSearchQuery] = useState('')
  const [filterStatus, setFilterStatus] = useState('')
  const [filterSort, setFilterSort] = useState('newest')

  const [selectedItems, setSelectedItems] = useState([])

  const [categories, setCategories] = useState(beritaContent.categories)

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
  const [formErrors, setFormErrors] = useState({})
  const [touched, setTouched] = useState({})

  // Validasi custom — Bahasa Indonesia, per-field
  const validateForm = useCallback((data = formData) => {
    const errors = {}
    if (!data.title || !data.title.trim()) {
      errors.title = 'Judul berita wajib diisi.'
    } else if (data.title.trim().length < 5) {
      errors.title = 'Judul minimal 5 karakter.'
    } else if (data.title.trim().length > 300) {
      errors.title = 'Judul maksimal 300 karakter.'
    }
    if (!data.category || !data.category.trim()) {
      errors.category = 'Kategori wajib dipilih.'
    }
    if (!data.published_date) {
      errors.published_date = 'Tanggal terbit wajib diisi.'
    }
    if (!data.content || !data.content.trim()) {
      errors.content = 'Konten lengkap wajib diisi.'
    }
    return errors
  }, [formData])

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
    setLoading(true)
    setError(null)
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
        setItems(list)
        setMeta(res.data.meta || { current_page: currentPage, limit: PAGE_SIZE, total_data: list.length, total_pages: Math.ceil(list.length / PAGE_SIZE) || 1 })
      }
    } catch (err) {
      setError(err?.message || 'Gagal memuat berita')
      setItems([])
    } finally {
      setLoading(false)
    }
  }, [currentPage, currentTab, searchQuery, filterStatus, filterSort])

  useEffect(() => {
    loadBerita()
  }, [loadBerita])

  // Ambil daftar kategori dinamis (dari data yang ada) untuk dropdown
  useEffect(() => {
    beritaService.getCategories()
      .then(res => {
        if (res && res.data && Array.isArray(res.data) && res.data.length > 0) {
          setCategories(res.data)
        }
      }).catch(() => {})
  }, [])

  useEffect(() => {
    setCurrentPage(1)
    setSelectedItems([])
  }, [currentTab, searchQuery, filterStatus, filterSort])

  function openForm(item = null) {
    setFormErrors({})
    setTouched({})
    if (item) {
      setFormMode('edit')
      setFormData({
        id: item.id,
        title: item.title,
        category: item.category,
        published_date: item.published_date,
        image_url: item.image_url || '',
        excerpt: item.excerpt || '',
        content: item.content || '',
        tags: item.tags || '',
        is_published: item.is_published,
      })
    } else {
      setFormMode('create')
      setFormData({
        id: null,
        title: '',
        category: categories[0] || 'Berita Nasional',
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
    const errors = validateForm()
    if (Object.keys(errors).length > 0) {
      setFormErrors(errors)
      setTouched(Object.keys(errors).reduce((acc, k) => ({ ...acc, [k]: true }), {}))
      return
    }
    setFormErrors({})
    setFormLoading(true)
    try {
      if (formMode === 'create') {
        await beritaService.create(formData)
        showToast('Berita berhasil dibuat.')
      } else {
        await beritaService.update(formData.id, formData)
        showToast('Berita berhasil diperbarui.')
      }
      setIsFormOpen(false)
      await loadBerita()
    } catch (err) {
      // Tangkap DUPLICATE_TITLE dari BE → tampilkan di field title
      if (err?.code === 'DUPLICATE_TITLE' || (err?.data?.errors && err.data.errors.some(x => x.field === 'title' && x.tag === 'duplicate_title'))) {
        setFormErrors(prev => ({ ...prev, title: 'Judul sudah digunakan, gunakan judul lain.' }))
        setTouched(prev => ({ ...prev, title: true }))
      } else {
        showToast(err.message || 'Gagal menyimpan berita.', 'error')
      }
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
          await beritaService.remove(id)
          await loadBerita()
          showToast('Berita berhasil dihapus.')
        },
      },
      restore: {
        type: 'info',
        title: 'Pulihkan Berita',
        message: 'Berita ini akan dikembalikan dari Sampah. Lanjutkan?',
        action: async () => {
          await beritaService.restore(id)
          await loadBerita()
          showToast('Berita berhasil dipulihkan.')
        },
      },
      toggle_publish: {
        type: extraData ? 'warning' : 'info',
        title: extraData ? 'Jadikan Draft' : 'Terbitkan Berita',
        message: extraData ? 'Berita ini akan diubah menjadi draft. Lanjutkan?' : 'Berita ini akan diterbitkan. Lanjutkan?',
        action: async () => {
          try {
            await beritaService.update(id, { is_published: !extraData })
            setItems(prev => prev.map(i => i.id === id ? { ...i, is_published: !extraData } : i))
            showToast(extraData ? 'Berita dijadikan draft.' : 'Berita berhasil diterbitkan!')
          } catch (err) {
            showToast(err?.message || 'Gagal mengubah status berita.', 'error')
          }
        },
      },
      bulk_delete: {
        type: 'danger',
        title: 'Hapus Massal',
        message: `Anda akan memindahkan ${selectedItems.length} item ke Sampah. Lanjutkan?`,
        action: async () => {
          await beritaService.bulkDelete(selectedItems)
          setSelectedItems([])
          await loadBerita()
          showToast('Berita berhasil dihapus secara massal.')
        },
      },
      bulk_restore: {
        type: 'info',
        title: 'Pulihkan Massal',
        message: `Anda akan memulihkan ${selectedItems.length} item dari Sampah. Lanjutkan?`,
        action: async () => {
          await beritaService.bulkRestore(selectedItems)
          setSelectedItems([])
          await loadBerita()
          showToast('Berita berhasil dipulihkan secara massal.')
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
      try { await action() } catch {}
    }
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

  const paginatedItems = items
  const totalPages = meta.total_pages || 1

  return (
    <AdminLayout title="Kelola Berita">
      {toast.show && <ToastNotification message={toast.message} type={toast.type} onClose={() => setToast({ ...toast, show: false })} />}
      <ConfirmDialog {...confirm} onClose={() => setConfirm(prev => ({ ...prev, isOpen: false, action: null }))} onConfirm={executeConfirm} />

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
        <div className="flex flex-col md:flex-row gap-3 items-stretch md:items-center bg-white p-4 rounded-xl border border-slate-200 shadow-sm">
          <div className="relative flex-1 min-w-0">
            <i className="ph-bold ph-magnifying-glass absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400" />
            <input
              type="text"
              value={searchQuery}
              onChange={e => setSearchQuery(e.target.value)}
              placeholder={beritaContent.admin.searchPlaceholder}
              className="w-full pl-10 pr-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl text-sm text-slate-700 placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 focus:bg-white transition-all"
            />
          </div>

          <div className="flex gap-2 shrink-0">
            <select
              value={filterStatus}
              onChange={e => setFilterStatus(e.target.value)}
              className="px-3.5 py-2.5 bg-slate-50 border border-slate-200 rounded-xl text-sm text-slate-600 outline-none focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 cursor-pointer transition-all hover:bg-slate-100"
            >
              <option value="">{beritaContent.admin.allStatus}</option>
              <option value="published">Public</option>
              <option value="draft">Draft</option>
            </select>
            <button
              onClick={resetFilter}
              className="bg-slate-50 border border-slate-200 text-slate-600 hover:bg-slate-100 hover:text-slate-800 px-3.5 py-2.5 rounded-xl text-sm font-semibold flex items-center gap-2 transition-colors"
            >
              <i className="ph-bold ph-arrows-counter-clockwise" /> Reset
            </button>
          </div>
        </div>

        {/* BULK ACTION BAR */}
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

        {/* DATA TABLE */}
        <div className="bg-white rounded-2xl border border-slate-200 overflow-hidden shadow-sm">
          {loading && <div className="py-16 text-center text-slate-500">Memuat berita...</div>}
          {!loading && error && (
            <div className="py-16 text-center text-red-600 font-medium">
              <i className="ph-bold ph-warning-circle text-2xl mb-2 block mx-auto" /> {error}
            </div>
          )}
          {!loading && !error && (
          <>
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
                      {currentTab === 'trash' ? (
                        <span className="inline-flex items-center gap-1.5 text-xs font-semibold text-red-500">
                          <i className="ph-bold ph-trash" /> Di Sampah
                        </span>
                      ) : (
                        <button onClick={() => confirmAction('toggle_publish', item.id, item.is_published)} className="inline-flex items-center gap-2 cursor-pointer">
                          <div className={`relative w-9 h-5 rounded-full transition-colors ${item.is_published ? 'bg-brand-500' : 'bg-slate-200'}`}>
                            <div className={`absolute top-[2px] left-[2px] w-4 h-4 bg-white rounded-full transition-transform ${item.is_published ? 'translate-x-4' : 'translate-x-0'}`} />
                          </div>
                          <span className={`text-xs font-semibold ${item.is_published ? 'text-brand-600' : 'text-slate-400'}`}>
                            {item.is_published ? 'Published' : 'Draft'}
                          </span>
                        </button>
                      )}
                    </td>
                    <td className="p-4 text-right">
                      <div className="flex items-center justify-end gap-1">
                        {currentTab === 'trash' ? (
                          <button onClick={() => confirmAction('restore', item.id)} className="p-2 text-slate-400 hover:text-emerald-600 rounded-lg" title="Pulihkan">
                            <i className="ph-bold ph-arrow-counter-clockwise text-base" />
                          </button>
                        ) : (
                          <>
                            <button onClick={() => openForm(item)} className="p-2 text-slate-400 hover:text-brand-600 rounded-lg">
                              <i className="ph-bold ph-pencil-simple text-base" />
                            </button>
                            <button onClick={() => confirmAction('delete', item.id)} className="p-2 text-slate-400 hover:text-red-600 rounded-lg">
                              <i className="ph-bold ph-trash text-base" />
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

          {/* PAGINATION (server-driven) */}
          {totalPages > 1 && (
            <div className="flex items-center justify-between px-4 py-3 border-t border-slate-200">
              <span className="text-xs text-slate-500">
                Hal {currentPage} dari {totalPages} · {meta.total_data || 0} data
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
            <form onSubmit={handleSubmit} noValidate className="p-6 space-y-4 overflow-y-auto">
              <div>
                <label className="block text-xs font-semibold text-slate-500 mb-1">Judul Berita *</label>
                <input
                  type="text"
                  value={formData.title}
                  onChange={e => {
                    setFormData({ ...formData, title: e.target.value })
                    if (touched.title) {
                      const errs = validateForm({ ...formData, title: e.target.value })
                      setFormErrors(prev => ({ ...prev, title: errs.title }))
                    }
                  }}
                  onBlur={() => {
                    setTouched(prev => ({ ...prev, title: true }))
                    const errs = validateForm()
                    setFormErrors(prev => ({ ...prev, title: errs.title }))
                  }}
                  className={`w-full px-3.5 py-2.5 border rounded-xl text-sm outline-none transition-colors ${touched.title && formErrors.title ? 'border-red-400 focus:ring-2 focus:ring-red-100' : 'border-slate-300 focus:border-brand-500 focus:ring-2 focus:ring-brand-100'}`}
                />
                {touched.title && formErrors.title && (
                  <p className="text-red-500 text-[11px] font-semibold mt-1.5 flex items-center gap-1">
                    <i className="ph-bold ph-warning-circle text-xs" /> {formErrors.title}
                  </p>
                )}
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-semibold text-slate-500 mb-1">Kategori *</label>
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
                          setFormData(prev => ({ ...prev, category: clean }))
                          if (touched.category) {
                            const errs = validateForm({ ...formData, category: clean })
                            setFormErrors(prev => ({ ...prev, category: errs.category }))
                          }
                        }
                        return
                      }
                      const val = e.target.value
                      setFormData(prev => ({ ...prev, category: val }))
                      if (touched.category) {
                        const errs = validateForm({ ...formData, category: val })
                        setFormErrors(prev => ({ ...prev, category: errs.category }))
                      }
                    }}
                    onBlur={() => {
                      setTouched(prev => ({ ...prev, category: true }))
                      const errs = validateForm()
                      setFormErrors(prev => ({ ...prev, category: errs.category }))
                    }}
                    className={`w-full px-3.5 py-2.5 border rounded-xl text-sm outline-none bg-white transition-colors ${touched.category && formErrors.category ? 'border-red-400 focus:ring-2 focus:ring-red-100' : 'border-slate-300 focus:border-brand-500 focus:ring-2 focus:ring-brand-100'}`}
                  >
                    <option value="">Pilih Kategori</option>
                    {categories.map(c => <option key={c} value={c}>{c}</option>)}
                    <option value="__new__">+ Buat Kategori Baru...</option>
                  </select>
                  {touched.category && formErrors.category && (
                    <p className="text-red-500 text-[11px] font-semibold mt-1.5 flex items-center gap-1">
                      <i className="ph-bold ph-warning-circle text-xs" /> {formErrors.category}
                    </p>
                  )}
                </div>
                <div>
                  <label className="block text-xs font-semibold text-slate-500 mb-1">Tanggal Terbit *</label>
                  <input
                    type="date"
                    value={formData.published_date}
                    onChange={e => {
                      setFormData({ ...formData, published_date: e.target.value })
                      if (touched.published_date) {
                        const errs = validateForm({ ...formData, published_date: e.target.value })
                        setFormErrors(prev => ({ ...prev, published_date: errs.published_date }))
                      }
                    }}
                    onBlur={() => {
                      setTouched(prev => ({ ...prev, published_date: true }))
                      const errs = validateForm()
                      setFormErrors(prev => ({ ...prev, published_date: errs.published_date }))
                    }}
                    className={`w-full px-3.5 py-2.5 border rounded-xl text-sm outline-none transition-colors ${touched.published_date && formErrors.published_date ? 'border-red-400 focus:ring-2 focus:ring-red-100' : 'border-slate-300 focus:border-brand-500 focus:ring-2 focus:ring-brand-100'}`}
                  />
                  {touched.published_date && formErrors.published_date && (
                    <p className="text-red-500 text-[11px] font-semibold mt-1.5 flex items-center gap-1">
                      <i className="ph-bold ph-warning-circle text-xs" /> {formErrors.published_date}
                    </p>
                  )}
                </div>
              </div>
              <div>
                <label className="block text-xs font-semibold text-slate-500 mb-1">URL Gambar</label>
                <input type="text" value={formData.image_url} onChange={e => setFormData({ ...formData, image_url: e.target.value })} className="w-full px-3.5 py-2.5 border border-slate-300 rounded-xl text-sm outline-none focus:border-brand-500 focus:ring-2 focus:ring-brand-100 transition-colors" />
              </div>
              <div>
                <label className="block text-xs font-semibold text-slate-500 mb-1">Ringkasan (Excerpt)</label>
                <textarea rows={2} value={formData.excerpt} onChange={e => setFormData({ ...formData, excerpt: e.target.value })} className="w-full px-3.5 py-2.5 border border-slate-300 rounded-xl text-sm outline-none focus:border-brand-500 focus:ring-2 focus:ring-brand-100 transition-colors" />
              </div>
              <div>
                <label className="block text-xs font-semibold text-slate-500 mb-1">Konten Lengkap *</label>
                <textarea
                  rows={6}
                  value={formData.content}
                  onChange={e => {
                    setFormData({ ...formData, content: e.target.value })
                    if (touched.content) {
                      const errs = validateForm({ ...formData, content: e.target.value })
                      setFormErrors(prev => ({ ...prev, content: errs.content }))
                    }
                  }}
                  onBlur={() => {
                    setTouched(prev => ({ ...prev, content: true }))
                    const errs = validateForm()
                    setFormErrors(prev => ({ ...prev, content: errs.content }))
                  }}
                  placeholder="Isi berita lengkap di sini..."
                  className={`w-full px-3.5 py-2.5 border rounded-xl text-sm outline-none transition-colors ${touched.content && formErrors.content ? 'border-red-400 focus:ring-2 focus:ring-red-100' : 'border-slate-300 focus:border-brand-500 focus:ring-2 focus:ring-brand-100'}`}
                />
                {touched.content && formErrors.content && (
                  <p className="text-red-500 text-[11px] font-semibold mt-1.5 flex items-center gap-1">
                    <i className="ph-bold ph-warning-circle text-xs" /> {formErrors.content}
                  </p>
                )}
              </div>
              <div>
                <label className="block text-xs font-semibold text-slate-500 mb-1">Tags (pisahkan dengan koma)</label>
                <input type="text" value={formData.tags} onChange={e => setFormData({ ...formData, tags: e.target.value })} placeholder="teknologi, pendidikan, digital" className="w-full px-3.5 py-2.5 border border-slate-300 rounded-xl text-sm outline-none focus:border-brand-500 focus:ring-2 focus:ring-brand-100 transition-colors" />
              </div>
              <div>
                <label className="flex items-center gap-2 text-sm cursor-pointer">
                  <input type="checkbox" checked={formData.is_published} onChange={e => setFormData({ ...formData, is_published: e.target.checked })} className="accent-brand-600" />
                  Terbitkan langsung (Published)
                </label>
              </div>
              <div className="flex justify-end gap-2 pt-4 border-t">
                <button type="button" onClick={() => setIsFormOpen(false)} disabled={formLoading} className="px-4 py-2 border rounded-xl text-sm font-semibold hover:bg-slate-50 disabled:opacity-50 transition-colors">Batal</button>
                <button type="submit" disabled={formLoading} className="px-5 py-2 bg-brand-600 text-white rounded-xl text-sm font-semibold hover:bg-brand-700 disabled:opacity-60 disabled:cursor-not-allowed transition-colors flex items-center gap-2">
                  {formLoading && <i className="ph-bold ph-circle-notch animate-spin text-sm" />}
                  {formLoading ? 'Menyimpan...' : 'Simpan'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </AdminLayout>
  )
}
