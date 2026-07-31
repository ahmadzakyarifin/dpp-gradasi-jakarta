import { useState, useEffect, useCallback } from 'react'
import AdminLayout from '../../layouts/AdminLayout'
import { beritaService } from '../../services/beritaService'
import { beritaContent } from '../../content/beritaContent'
import ConfirmDialog from '../../components/admin/ConfirmDialog'
import ToastNotification from '../../components/admin/ToastNotification'

const PAGE_SIZE = 5

export default function BeritaAdmin() {
  // --- Data & Loading ---
  const [items, setItems] = useState([])
  const [meta, setMeta] = useState({ current_page: 1, limit: PAGE_SIZE, total_data: 0, total_pages: 1 })
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  // --- Tab, Filter, Search, Sort ---
  const [currentTab, setCurrentTab] = useState('active')
  const [currentPage, setCurrentPage] = useState(1)
  const [searchQuery, setSearchQuery] = useState('')
  const [filterStatus, setFilterStatus] = useState('')
  const [filterSort, setFilterSort] = useState('newest')

  // --- Selection ---
  const [selectedItems, setSelectedItems] = useState([])

  // --- Form ---
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

  // --- Confirm Dialog ---
  const [confirm, setConfirm] = useState({
    isOpen: false,
    type: 'danger',
    title: '',
    message: '',
    action: null,
  })

  // --- Toast ---
  const [toast, setToast] = useState({ show: false, message: '', type: 'success' })

  const showToast = useCallback((message, type = 'success') => {
    setToast({ show: true, message, type })
  }, [])

  // ==========================================================================
  // DATA LOADING
  // ==========================================================================
  const loadBerita = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const params = {
        page: currentPage,
        limit: PAGE_SIZE,
        sort: filterSort,
      }
      if (searchQuery.trim()) params.search = searchQuery.trim()

      if (currentTab === 'trash') {
        params.status = 'trashed'
      } else if (filterStatus) {
        params.status = filterStatus
      }

      const res = await beritaService.listAdmin(params)
      if (res.success && res.data) {
        setItems(res.data.berita || [])
        setMeta(res.data.meta || { current_page: 1, limit: PAGE_SIZE, total_data: 0, total_pages: 1 })
      } else {
        setError('Gagal memuat data berita.')
      }
    } catch (err) {
      setError(err.message || 'Kesalahan koneksi ke server.')
    } finally {
      setLoading(false)
    }
  }, [currentPage, currentTab, searchQuery, filterStatus, filterSort])

  useEffect(() => {
    loadBerita()
  }, [loadBerita])

  // Reset page when tab, filter, or search changes
  useEffect(() => {
    setCurrentPage(1)
    setSelectedItems([])
  }, [currentTab, searchQuery, filterStatus, filterSort])

  // ==========================================================================
  // FORM
  // ==========================================================================
  function openForm(item = null) {
    if (item) {
      setFormMode('edit')
      // Need full detail for editing content
      beritaService.detailById(item.id).then(res => {
        if (res.success && res.data) {
          const d = res.data
          setFormData({
            id: d.id,
            title: d.title || '',
            category: d.category || 'Berita Nasional',
            published_date: d.published_date || new Date().toISOString().slice(0, 10),
            image_url: d.image_url || '',
            excerpt: d.excerpt || '',
            content: d.content || '',
            tags: Array.isArray(d.tags) ? d.tags.join(', ') : (d.tags || ''),
            is_published: d.is_published !== false,
          })
        }
      }).catch(() => {
        showToast('Gagal memuat detail berita untuk diedit.', 'error')
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
      const payload = {
        title: formData.title,
        category: formData.category,
        published_date: formData.published_date,
        image_url: formData.image_url,
        excerpt: formData.excerpt,
        content: formData.content,
        tags: formData.tags,
        is_published: formData.is_published,
      }

      if (formMode === 'create') {
        await beritaService.create(payload)
        showToast('Berita berhasil ditambahkan!')
      } else {
        await beritaService.update(formData.id, payload)
        showToast('Berita berhasil diperbarui!')
      }
      setIsFormOpen(false)
      loadBerita()
    } catch (err) {
      showToast(err.message || 'Gagal menyimpan berita.', 'error')
    } finally {
      setFormLoading(false)
    }
  }

  // ==========================================================================
  // ACTIONS
  // ==========================================================================
  function confirmAction(type, id = null, extraData = null) {
    const configs = {
      delete: {
        type: 'danger',
        title: 'Hapus Berita',
        message: 'Berita ini akan dipindahkan ke Sampah (History). Anda masih bisa memulihkannya nanti.',
        action: async () => {
          await beritaService.remove(id)
          showToast('Berita berhasil dipindahkan ke sampah.')
          loadBerita()
        },
      },
      restore: {
        type: 'success',
        title: 'Pulihkan Berita',
        message: 'Berita ini akan dikembalikan ke daftar aktif.',
        action: async () => {
          await beritaService.restore(id)
          showToast('Berita berhasil dipulihkan.')
          loadBerita()
        },
      },
      toggle_publish: {
        type: extraData ? 'warning' : 'info',
        title: extraData ? 'Jadikan Draft' : 'Terbitkan Berita',
        message: extraData
          ? 'Berita ini akan ditarik dari publik dan diubah menjadi draft. Lanjutkan?'
          : 'Berita ini akan bisa dilihat oleh publik di website. Lanjutkan?',
        action: async () => {
          // Fetch full detail first to get content for update
          const detail = await beritaService.detailById(id)
          if (detail.success && detail.data) {
            const d = detail.data
            await beritaService.update(id, {
              title: d.title,
              category: d.category,
              published_date: d.published_date,
              image_url: d.image_url,
              excerpt: d.excerpt,
              content: d.content,
              tags: Array.isArray(d.tags) ? d.tags.join(',') : (d.tags || ''),
              is_published: !extraData,
            })
          }
          showToast(extraData ? 'Berita dijadikan draft.' : 'Berita berhasil diterbitkan!')
          loadBerita()
        },
      },
      bulk_delete: {
        type: 'danger',
        title: 'Hapus Massal',
        message: `Anda akan memindahkan ${selectedItems.length} item ke Sampah. Lanjutkan?`,
        action: async () => {
          await beritaService.bulkDelete(selectedItems)
          setSelectedItems([])
          showToast('Berita berhasil dihapus secara massal.')
          loadBerita()
        },
      },
      bulk_restore: {
        type: 'success',
        title: 'Pulihkan Massal',
        message: `Anda akan memulihkan ${selectedItems.length} item kembali ke daftar aktif. Lanjutkan?`,
        action: async () => {
          await beritaService.bulkRestore(selectedItems)
          setSelectedItems([])
          showToast('Berita berhasil dipulihkan secara massal.')
          loadBerita()
        },
      },
    }

    const cfg = configs[type]
    if (!cfg) return

    setConfirm({
      isOpen: true,
      type: cfg.type,
      title: cfg.title,
      message: cfg.message,
      action: cfg.action,
    })
  }

  async function executeConfirm() {
    setConfirm(prev => ({ ...prev, isOpen: false }))
    try {
      await confirm.action?.()
    } catch (err) {
      showToast(err.message || 'Aksi gagal dilakukan.', 'error')
    }
  }

  // ==========================================================================
  // SELECTION HELPERS
  // ==========================================================================
  const isAllSelected = items.length > 0 && selectedItems.length === items.length

  function toggleAll() {
    if (isAllSelected) {
      setSelectedItems([])
    } else {
      setSelectedItems(items.map(i => i.id))
    }
  }

  function toggleItem(id) {
    setSelectedItems(prev =>
      prev.includes(id) ? prev.filter(x => x !== id) : [...prev, id]
    )
  }

  // ==========================================================================
  // UTILITIES
  // ==========================================================================
  function copyLink(id) {
    const url = `${window.location.origin}/berita/${id}`
    navigator.clipboard.writeText(url).then(() => {
      showToast('Link berhasil disalin!')
    }).catch(() => {
      showToast('Gagal menyalin link.', 'error')
    })
  }

  function resetFilter() {
    setSearchQuery('')
    setFilterStatus('')
    setFilterSort('newest')
    setCurrentPage(1)
    showToast('Filter direset.', 'info')
  }

  function formatDate(dateStr) {
    if (!dateStr) return '-'
    try {
      return new Date(dateStr).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' })
    } catch {
      return dateStr
    }
  }

  // Pagination helpers
  const startIndex = (meta.current_page - 1) * meta.limit + 1
  const endIndex = Math.min(meta.current_page * meta.limit, meta.total_data)

  return (
    <AdminLayout title="Kelola Berita">
      <div className="space-y-5">

        {/* ================================================================ */}
        {/* HEADER: Tabs + Add Button */}
        {/* ================================================================ */}
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
          <div className="flex bg-white rounded-xl p-1 border border-slate-200 shadow-sm">
            <button
              onClick={() => setCurrentTab('active')}
              className={`px-4 py-2 rounded-lg text-sm font-semibold transition-all ${
                currentTab === 'active'
                  ? 'bg-brand-50 text-brand-600 shadow-sm'
                  : 'text-slate-500 hover:text-slate-700'
              }`}
            >
              {beritaContent.admin.activeTab}
            </button>
            <button
              onClick={() => setCurrentTab('trash')}
              className={`px-4 py-2 rounded-lg text-sm font-semibold transition-all flex items-center gap-1.5 ${
                currentTab === 'trash'
                  ? 'bg-red-50 text-red-600 shadow-sm'
                  : 'text-slate-500 hover:text-slate-700'
              }`}
            >
              <i className="ph-bold ph-trash" /> {beritaContent.admin.trashTab}
            </button>
          </div>

          {currentTab === 'active' && (
            <button
              onClick={() => openForm()}
              className="bg-brand-600 hover:bg-brand-700 text-white px-4 py-2.5 rounded-xl text-sm font-semibold flex items-center gap-2 transition-colors shadow-sm hover:shadow-md"
            >
              <i className="ph-bold ph-plus-circle text-lg" /> {beritaContent.admin.add}
            </button>
          )}
        </div>

        {/* ================================================================ */}
        {/* FILTER BAR */}
        {/* ================================================================ */}
        <div className="flex flex-col md:flex-row gap-3 items-center bg-white p-4 rounded-xl border border-slate-200 shadow-sm">
          {/* Search */}
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

          {/* Filters */}
          <div className="flex gap-2 w-full md:w-auto flex-wrap sm:flex-nowrap">
            {currentTab === 'active' && (
              <select
                value={filterStatus}
                onChange={e => setFilterStatus(e.target.value)}
                className="flex-1 sm:flex-initial px-3 py-2.5 bg-slate-50 border border-slate-200 rounded-xl text-sm text-slate-600 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 focus:bg-white transition-colors"
              >
                <option value="">{beritaContent.admin.allStatus}</option>
                <option value="published">{beritaContent.admin.published}</option>
                <option value="draft">{beritaContent.admin.draft}</option>
              </select>
            )}
            <select
              value={filterSort}
              onChange={e => setFilterSort(e.target.value)}
              className="flex-1 sm:flex-initial px-3 py-2.5 bg-slate-50 border border-slate-200 rounded-xl text-sm text-slate-600 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 focus:bg-white transition-colors"
            >
              <option value="newest">{beritaContent.admin.newest}</option>
              <option value="oldest">{beritaContent.admin.oldest}</option>
              <option value="most_viewed">Terpopuler</option>
            </select>
            <button
              onClick={resetFilter}
              className="flex-shrink-0 bg-slate-50 border border-slate-200 text-slate-600 hover:bg-slate-100 px-3.5 py-2.5 rounded-xl text-sm font-semibold flex items-center gap-2 transition-colors"
            >
              <i className="ph-bold ph-arrows-counter-clockwise" /> {beritaContent.admin.reset}
            </button>
          </div>
        </div>

        {/* ================================================================ */}
        {/* BULK ACTIONS BAR */}
        {/* ================================================================ */}
        {selectedItems.length > 0 && (
          <div className="bg-indigo-50 border border-indigo-100 rounded-xl p-3.5 flex items-center justify-between shadow-sm animate-[slideDown_200ms_ease-out]">
            <span className="text-sm text-indigo-800 font-semibold">
              {selectedItems.length} item terpilih
            </span>
            <div className="flex gap-2">
              {currentTab === 'active' ? (
                <button
                  onClick={() => confirmAction('bulk_delete')}
                  className="bg-red-100 text-red-700 hover:bg-red-200 px-3.5 py-2 rounded-lg text-xs font-semibold flex items-center gap-1.5 transition-colors"
                >
                  <i className="ph-bold ph-trash" /> Hapus Massal
                </button>
              ) : (
                <button
                  onClick={() => confirmAction('bulk_restore')}
                  className="bg-emerald-100 text-emerald-700 hover:bg-emerald-200 px-3.5 py-2 rounded-lg text-xs font-semibold flex items-center gap-1.5 transition-colors"
                >
                  <i className="ph-bold ph-arrow-counter-clockwise" /> Pulihkan Massal
                </button>
              )}
            </div>
          </div>
        )}

        {/* ================================================================ */}
        {/* DATA TABLE */}
        {/* ================================================================ */}
        {loading && (
          <div className="py-16 text-center text-slate-500 font-medium">
            <i className="ph-bold ph-spinner animate-spin text-2xl mb-2 block text-brand-500" />
            Memuat berita...
          </div>
        )}
        {error && (
          <div className="py-16 text-center">
            <i className="ph-bold ph-warning-circle text-3xl text-red-400 mb-2 block" />
            <p className="text-red-600 font-medium text-sm">{error}</p>
            <button onClick={loadBerita} className="mt-3 text-sm text-brand-600 font-semibold hover:underline">
              Coba Lagi
            </button>
          </div>
        )}

        {!loading && !error && (
          <div className="bg-white rounded-2xl border border-slate-200 overflow-hidden shadow-sm">
            <div className="overflow-x-auto">
              <table className="w-full text-left border-collapse">
                <thead>
                  <tr className="bg-slate-50 border-b border-slate-200 text-[11px] uppercase tracking-wider text-slate-500 font-semibold">
                    <th className="p-4 w-12 text-center">
                      <input
                        type="checkbox"
                        onChange={toggleAll}
                        checked={isAllSelected}
                        className="rounded border-slate-300 text-brand-600 focus:ring-brand-500 cursor-pointer"
                      />
                    </th>
                    <th className="p-4">Judul & Kategori</th>
                    <th className="p-4">Tanggal</th>
                    <th className="p-4">Status</th>
                    <th className="p-4 text-right">Aksi</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100 text-sm text-slate-700">
                  {items.length === 0 ? (
                    <tr>
                      <td colSpan={5} className="py-16 text-center">
                        <i className="ph-bold ph-folder-open text-4xl text-slate-300 mb-2 block" />
                        <p className="text-slate-500 text-sm">{beritaContent.admin.empty}</p>
                      </td>
                    </tr>
                  ) : (
                    items.map(item => (
                      <tr key={item.id} className="hover:bg-slate-50/70 transition-colors group">
                        {/* Checkbox */}
                        <td className="p-4 text-center">
                          <input
                            type="checkbox"
                            checked={selectedItems.includes(item.id)}
                            onChange={() => toggleItem(item.id)}
                            className="rounded border-slate-300 text-brand-600 focus:ring-brand-500 cursor-pointer"
                          />
                        </td>

                        {/* Title + Image + Category + Views */}
                        <td className="p-4">
                          <div className="flex items-start gap-3">
                            {item.image_url ? (
                              <img
                                src={item.image_url}
                                alt=""
                                className="w-16 h-12 rounded-lg object-cover border border-slate-200 shrink-0"
                              />
                            ) : (
                              <div className="w-16 h-12 rounded-lg bg-slate-100 border border-slate-200 shrink-0 flex items-center justify-center">
                                <i className="ph-bold ph-image text-slate-300 text-lg" />
                              </div>
                            )}
                            <div className="min-w-0">
                              <p className="font-semibold text-slate-900 leading-snug line-clamp-1">{item.title}</p>
                              <div className="flex items-center gap-2 mt-1">
                                <span className="inline-block bg-slate-100 text-slate-600 text-[10px] font-semibold px-2 py-0.5 rounded-full">
                                  {item.category}
                                </span>
                                <span className="text-xs text-slate-400 flex items-center gap-1">
                                  <i className="ph-bold ph-eye text-[10px]" /> {item.views || 0}
                                </span>
                              </div>
                            </div>
                          </div>
                        </td>

                        {/* Date */}
                        <td className="p-4 text-slate-500 text-sm whitespace-nowrap">
                          {formatDate(item.published_date || item.created_at)}
                        </td>

                        {/* Status */}
                        <td className="p-4">
                          {currentTab === 'active' ? (
                            <button
                              onClick={() => confirmAction('toggle_publish', item.id, item.is_published)}
                              className="inline-flex items-center gap-2 cursor-pointer group/toggle"
                            >
                              {/* Toggle Switch */}
                              <div className={`relative w-9 h-5 rounded-full transition-colors ${
                                item.is_published ? 'bg-brand-500' : 'bg-slate-200'
                              }`}>
                                <div className={`absolute top-[2px] left-[2px] w-4 h-4 bg-white rounded-full border transition-transform shadow-sm ${
                                  item.is_published
                                    ? 'translate-x-4 border-brand-400'
                                    : 'translate-x-0 border-slate-300'
                                }`} />
                              </div>
                              <span className={`text-xs font-semibold ${
                                item.is_published ? 'text-brand-600' : 'text-slate-400'
                              }`}>
                                {item.is_published ? 'Published' : 'Draft'}
                              </span>
                            </button>
                          ) : (
                            <span className="inline-flex items-center gap-1 bg-red-50 text-red-600 text-xs px-2.5 py-1 rounded-lg font-semibold border border-red-100">
                              <i className="ph-bold ph-trash text-[10px]" /> Terhapus
                            </span>
                          )}
                        </td>

                        {/* Actions */}
                        <td className="p-4 text-right">
                          <div className="flex items-center justify-end gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                            {currentTab === 'active' ? (
                              <>
                                <button
                                  onClick={() => openForm(item)}
                                  className="p-2 text-slate-400 hover:text-brand-600 hover:bg-brand-50 rounded-lg transition-colors"
                                  title="Edit"
                                >
                                  <i className="ph-bold ph-pencil-simple text-base" />
                                </button>
                                <button
                                  onClick={() => copyLink(item.id)}
                                  className="p-2 text-slate-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
                                  title="Salin Tautan"
                                >
                                  <i className="ph-bold ph-link text-base" />
                                </button>
                                <button
                                  onClick={() => confirmAction('delete', item.id)}
                                  className="p-2 text-slate-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                                  title="Hapus"
                                >
                                  <i className="ph-bold ph-trash text-base" />
                                </button>
                              </>
                            ) : (
                              <button
                                onClick={() => confirmAction('restore', item.id)}
                                className="p-2 text-slate-400 hover:text-emerald-600 hover:bg-emerald-50 rounded-lg transition-colors"
                                title="Pulihkan"
                              >
                                <i className="ph-bold ph-arrow-counter-clockwise text-base" />
                              </button>
                            )}
                          </div>
                        </td>
                      </tr>
                    ))
                  )}
                </tbody>
              </table>
            </div>

            {/* Pagination Footer */}
            {meta.total_data > 0 && (
              <div className="bg-slate-50 border-t border-slate-200 px-5 py-3.5 flex items-center justify-between rounded-b-2xl">
                <div className="text-sm text-slate-500">
                  Menampilkan{' '}
                  <span className="font-semibold text-slate-700">{startIndex}</span>{' '}
                  sampai{' '}
                  <span className="font-semibold text-slate-700">{endIndex}</span>{' '}
                  dari{' '}
                  <span className="font-semibold text-slate-700">{meta.total_data}</span>{' '}
                  hasil
                </div>
                <div className="flex gap-1">
                  <button
                    onClick={() => setCurrentPage(p => Math.max(1, p - 1))}
                    disabled={currentPage <= 1}
                    className="px-3 py-1.5 border border-slate-200 bg-white text-slate-500 rounded-lg hover:bg-slate-50 disabled:opacity-40 text-sm font-medium transition-colors"
                  >
                    Prev
                  </button>
                  {Array.from({ length: meta.total_pages }, (_, i) => i + 1).map(p => (
                    <button
                      key={p}
                      onClick={() => setCurrentPage(p)}
                      className={`px-3 py-1.5 border rounded-lg text-sm font-medium transition-colors ${
                        currentPage === p
                          ? 'border-brand-500 bg-brand-50 text-brand-600'
                          : 'border-slate-200 bg-white text-slate-500 hover:bg-slate-50'
                      }`}
                    >
                      {p}
                    </button>
                  ))}
                  <button
                    onClick={() => setCurrentPage(p => Math.min(meta.total_pages, p + 1))}
                    disabled={currentPage >= meta.total_pages}
                    className="px-3 py-1.5 border border-slate-200 bg-white text-slate-500 rounded-lg hover:bg-slate-50 disabled:opacity-40 text-sm font-medium transition-colors"
                  >
                    Next
                  </button>
                </div>
              </div>
            )}
          </div>
        )}
      </div>

      {/* ================================================================== */}
      {/* FORM MODAL */}
      {/* ================================================================== */}
      {isFormOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
          <div className="fixed inset-0 bg-black/40 backdrop-blur-[2px]" onClick={() => setIsFormOpen(false)} />
          <div className="relative bg-white rounded-2xl shadow-2xl max-w-3xl w-full max-h-[90vh] flex flex-col animate-[scaleIn_200ms_ease-out] overflow-hidden">
            {/* Header */}
            <div className="border-b border-slate-200 px-6 py-4 flex items-center justify-between shrink-0">
              <h3 className="font-heading font-bold text-slate-900 text-lg">
                {formMode === 'create' ? 'Tambah Berita Baru' : 'Edit Berita'}
              </h3>
              <button onClick={() => setIsFormOpen(false)} className="text-slate-400 hover:text-slate-600 transition-colors p-1">
                <i className="ph-bold ph-x text-lg" />
              </button>
            </div>

            {/* Form Body */}
            <form onSubmit={handleSubmit} className="flex-1 overflow-y-auto">
              <div className="px-6 py-5 space-y-5">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
                  {/* Title */}
                  <div className="md:col-span-2">
                    <label className="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1.5">Judul Berita *</label>
                    <input
                      type="text"
                      value={formData.title}
                      onChange={e => setFormData(prev => ({ ...prev, title: e.target.value }))}
                      required
                      className="w-full px-3.5 py-2.5 border border-slate-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 transition-colors"
                    />
                  </div>

                  {/* Category */}
                  <div>
                    <label className="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1.5">Kategori *</label>
                    <select
                      value={formData.category}
                      onChange={e => setFormData(prev => ({ ...prev, category: e.target.value }))}
                      required
                      className="w-full px-3.5 py-2.5 border border-slate-200 rounded-xl text-sm bg-white focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 transition-colors"
                    >
                      {beritaContent.categories.map(cat => (
                        <option key={cat} value={cat}>{cat}</option>
                      ))}
                    </select>
                  </div>

                  {/* Published Date */}
                  <div>
                    <label className="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1.5">Tanggal Terbit *</label>
                    <input
                      type="date"
                      value={formData.published_date}
                      onChange={e => setFormData(prev => ({ ...prev, published_date: e.target.value }))}
                      required
                      className="w-full px-3.5 py-2.5 border border-slate-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 transition-colors"
                    />
                  </div>

                  {/* Tags */}
                  <div>
                    <label className="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1.5">Tags (pisahkan koma)</label>
                    <input
                      type="text"
                      value={formData.tags}
                      onChange={e => setFormData(prev => ({ ...prev, tags: e.target.value }))}
                      placeholder="digital, literasi, gradasi"
                      className="w-full px-3.5 py-2.5 border border-slate-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 transition-colors"
                    />
                  </div>

                  {/* Image Preview + URL */}
                  <div>
                    <label className="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1.5">Gambar Utama (URL)</label>
                    <input
                      type="url"
                      value={formData.image_url}
                      onChange={e => setFormData(prev => ({ ...prev, image_url: e.target.value }))}
                      placeholder="https://example.com/cover.jpg"
                      className="w-full px-3.5 py-2.5 border border-slate-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 transition-colors"
                    />
                  </div>

                  {/* Excerpt */}
                  <div className="md:col-span-2">
                    <label className="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1.5">Ringkasan (Excerpt)</label>
                    <textarea
                      rows={2}
                      value={formData.excerpt}
                      onChange={e => setFormData(prev => ({ ...prev, excerpt: e.target.value }))}
                      className="w-full px-3.5 py-2.5 border border-slate-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 transition-colors resize-none"
                    />
                  </div>

                  {/* Content */}
                  <div className="md:col-span-2">
                    <label className="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1.5">Konten Lengkap *</label>
                    <textarea
                      rows={8}
                      value={formData.content}
                      onChange={e => setFormData(prev => ({ ...prev, content: e.target.value }))}
                      required
                      placeholder="Format HTML atau teks..."
                      className="w-full px-3.5 py-2.5 border border-slate-200 rounded-xl text-sm font-mono focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 transition-colors resize-y"
                    />
                  </div>

                  {/* Publish Checkbox */}
                  <div className="md:col-span-2">
                    <label className="flex items-center gap-2.5 cursor-pointer">
                      <input
                        type="checkbox"
                        checked={formData.is_published}
                        onChange={e => setFormData(prev => ({ ...prev, is_published: e.target.checked }))}
                        className="rounded border-slate-300 text-brand-600 focus:ring-brand-500 w-4 h-4"
                      />
                      <span className="text-sm font-semibold text-slate-700">Langsung Terbitkan (Publish)</span>
                    </label>
                  </div>
                </div>
              </div>

              {/* Footer */}
              <div className="bg-slate-50 border-t border-slate-200 px-6 py-4 flex justify-end gap-2.5 shrink-0">
                <button
                  type="button"
                  onClick={() => setIsFormOpen(false)}
                  className="px-5 py-2.5 bg-white border border-slate-200 rounded-xl text-sm font-semibold text-slate-600 hover:bg-slate-50 transition-colors shadow-sm"
                >
                  {beritaContent.admin.cancel}
                </button>
                <button
                  type="submit"
                  disabled={formLoading}
                  className="px-5 py-2.5 bg-brand-600 hover:bg-brand-700 text-white rounded-xl text-sm font-semibold transition-colors shadow-sm disabled:opacity-60 disabled:cursor-not-allowed flex items-center gap-2"
                >
                  {formLoading && <i className="ph-bold ph-spinner animate-spin" />}
                  {beritaContent.admin.save}
                </button>
              </div>
            </form>
          </div>

          <style>{`
            @keyframes scaleIn { from { opacity: 0; transform: scale(0.95); } to { opacity: 1; transform: scale(1); } }
          `}</style>
        </div>
      )}

      {/* ================================================================== */}
      {/* CONFIRM DIALOG */}
      {/* ================================================================== */}
      <ConfirmDialog
        isOpen={confirm.isOpen}
        title={confirm.title}
        message={confirm.message}
        type={confirm.type}
        onConfirm={executeConfirm}
        onCancel={() => setConfirm(prev => ({ ...prev, isOpen: false }))}
      />

      {/* ================================================================== */}
      {/* TOAST NOTIFICATION */}
      {/* ================================================================== */}
      <ToastNotification
        show={toast.show}
        message={toast.message}
        type={toast.type}
        onClose={() => setToast(prev => ({ ...prev, show: false }))}
      />

      <style>{`
        @keyframes slideDown { from { opacity: 0; transform: translateY(-8px); } to { opacity: 1; transform: translateY(0); } }
      `}</style>
    </AdminLayout>
  )
}
