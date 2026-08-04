import { useState, useEffect, useCallback } from 'react'
import AdminLayout from '../../layouts/AdminLayout'
import { slidersService } from '../../services/slidersService'
import ConfirmDialog from '../../components/admin/ConfirmDialog'
import ToastNotification from '../../components/admin/ToastNotification'
import { resolveAssetUrl } from '../../utils/assetUrl'
import { useFormErrors, useRateLimitCooldown } from '../../utils/parseApiError'

const PAGE_SIZE = 5

export default function SlidersAdmin() {
  const [items, setItems] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  const [currentTab, setCurrentTab] = useState('active') // active | trash
  const [selectedItems, setSelectedItems] = useState([])
  const [searchQuery, setSearchQuery] = useState('')
  const [filterActive, setFilterActive] = useState('')
  const [currentPage, setCurrentPage] = useState(1)

  const [isFormOpen, setIsFormOpen] = useState(false)
  const [formMode, setFormMode] = useState('create')
  const [formData, setFormData] = useState({
    id: null,
    title: '',
    subtitle: '',
    tag: '',
    image_path: '',
    link_url: '',
    sort_order: 1,
    event_date: '',
    location: '',
    is_new: false,
    is_active: true,
  })

  const [formErrors, setFormErrors] = useState({})
  const [touched, setTouched] = useState({})

  const validateForm = useCallback((data = formData) => {
    const errors = {}
    if (!data.title || !data.title.trim()) {
      errors.title = 'Judul utama wajib diisi.'
    }
    if (data.sort_order === undefined || data.sort_order === null || String(data.sort_order).trim() === '') {
      errors.sort_order = 'Urutan wajib diisi.'
    }
    if (!data.image_path || !data.image_path.trim()) {
      errors.image_path = 'Gambar slider (URL) wajib diisi.'
    }
    return errors
  }, [formData])

  const [confirm, setConfirm] = useState({ isOpen: false, type: 'danger', title: '', message: '', action: null })
  const [toast, setToast] = useState({ show: false, message: '', type: 'success' })
  // Error backend: pesan error dari helper + countdown rate limit
  const { fieldErrors, applyError, clearFieldError, resetFieldErrors } = useFormErrors()
  const { cooldown, isLimited, applyRateLimit } = useRateLimitCooldown()

  const showToast = useCallback((message, type = 'success') => {
    setToast({ show: true, message, type })
    setTimeout(() => setToast(prev => ({ ...prev, show: false })), 3000)
  }, [])

  const loadSliders = useCallback(() => {
    setLoading(true)
    setError(null)
    slidersService.listAdmin(false)
      .then(res => {
        if (res && res.data) {
          const list = Array.isArray(res.data) ? res.data : (res.data.sliders || [])
          setItems(list)
        } else {
          setError('Gagal memuat sliders')
        }
      })
      .catch(() => setError('Kesalahan koneksi ke server'))
      .finally(() => setLoading(false))
  }, [])

  useEffect(() => {
    loadSliders()
  }, [loadSliders])

  // Reset halaman & pilihan saat tab / filter / pencarian berubah
  useEffect(() => {
    setCurrentPage(1)
    setSelectedItems([])
  }, [currentTab, searchQuery, filterActive])

  // --- Filtering (mirip sliders.html) ---
  const filteredItems = items
    .filter(item => {
      const inTab = currentTab === 'trash' ? !item.is_active : item.is_active
      if (!inTab) return false
      if (searchQuery) {
        const q = searchQuery.toLowerCase()
        const hay = `${item.title || ''} ${item.subtitle || ''}`.toLowerCase()
        if (!hay.includes(q)) return false
      }
      if (filterActive) {
        const wantActive = filterActive === 'active'
          if (item.is_active !== wantActive) return false
      }
      return true
    })
    .sort((a, b) => (a.sort_order ?? 0) - (b.sort_order ?? 0))

  const totalPages = Math.ceil(filteredItems.length / PAGE_SIZE) || 1
  const paginatedItems = filteredItems.slice((currentPage - 1) * PAGE_SIZE, currentPage * PAGE_SIZE)
  const pageStart = filteredItems.length === 0 ? 0 : (currentPage - 1) * PAGE_SIZE + 1
  const pageEnd = Math.min(currentPage * PAGE_SIZE, filteredItems.length)
  const isAllSelected = filteredItems.length > 0 && selectedItems.length === filteredItems.length

  function toggleAll() {
    setSelectedItems(isAllSelected ? [] : filteredItems.map(i => i.id))
  }
  function toggleOne(id) {
    setSelectedItems(prev => prev.includes(id) ? prev.filter(x => x !== id) : [...prev, id])
  }
  function resetFilter() {
    setSearchQuery('')
    setFilterActive('')
    setCurrentPage(1)
    showToast('Filter direset.', 'success')
  }

  // --- Form ---
  const openForm = (item = null) => {
    setFormErrors({})
    setTouched({})
    resetFieldErrors()
    if (item) {
      setFormMode('edit')
      setFormData({
        id: item.id,
        title: item.title,
        subtitle: item.subtitle || '',
        tag: item.tag || '',
        image_path: item.image_path || item.image_url || '',
        link_url: item.link_url || '',
        sort_order: item.sort_order,
        event_date: item.event_date || '',
        location: item.location || '',
        is_new: item.is_new,
        is_active: item.is_active,
      })
    } else {
      setFormMode('create')
      setFormData({
        id: null,
        title: '',
        subtitle: '',
        tag: '',
        image_path: '',
        link_url: '',
        sort_order: items.length + 1,
        event_date: '',
        location: '',
        is_new: false,
        is_active: true,
      })
    }
    setIsFormOpen(true)
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    const errors = validateForm()
    if (Object.keys(errors).length > 0) {
      setFormErrors(errors)
      setTouched(Object.keys(errors).reduce((acc, k) => ({ ...acc, [k]: true }), {}))
      return
    }
    setFormErrors({})
    resetFieldErrors()

    try {
      if (formMode === 'create') {
        await slidersService.create(formData)
        showToast('Slider berhasil ditambahkan.')
      } else {
        await slidersService.update(formData.id, formData)
        showToast('Slider berhasil diperbarui.')
      }
      setIsFormOpen(false)
      loadSliders()
    } catch (err) {
      const parsed = applyError(err)
      applyRateLimit(err)
      setFormErrors(prev => ({ ...prev, ...parsed.fieldErrors }))
      setTouched(prev => ({ ...prev, ...Object.keys(parsed.fieldErrors).reduce((acc, k) => ({ ...acc, [k]: true }), {}) }))
      if (Object.keys(parsed.fieldErrors).length === 0) {
        showToast(parsed.message || 'Gagal menyimpan slider', 'error')
      }
    }
  }

  // Update urutan langsung dari input number di tabel (mirip sliders.html)
  const handleSortChange = async (item, value) => {
    const newOrder = Number(value)
    if (!Number.isFinite(newOrder) || newOrder === item.sort_order) return
    try {
      await slidersService.update(item.id, { ...item, sort_order: newOrder })
      showToast('Urutan slider diperbarui.')
      loadSliders()
    } catch (err) {
      showToast(err.message || 'Gagal mengubah urutan', 'error')
    }
  }

  const confirmAction = (type, item = null) => {
    const title = item ? item.title : ''
    const configs = {
      delete: {
        type: 'danger',
        title: 'Hapus Slider',
        message: `Anda akan memindahkan "${title}" ke Sampah. Lanjutkan?`,
        action: async () => {
          await slidersService.remove(item.id)
          showToast('Slider berhasil dihapus.')
          loadSliders()
        },
      },
      restore: {
        type: 'info',
        title: 'Pulihkan Slider',
        message: `Anda akan memulihkan "${title}" dari Sampah. Lanjutkan?`,
        action: async () => {
          await slidersService.restore(item.id)
          showToast('Slider berhasil dipulihkan.')
          loadSliders()
        },
      },
      toggle_publish: {
        type: 'info',
        title: item.is_active ? 'Nonaktifkan Slider' : 'Aktifkan Slider',
        message: `Anda akan ${item.is_active ? 'menonaktifkan' : 'mengaktifkan'} "${title}". Lanjutkan?`,
        action: async () => {
          await slidersService.update(item.id, { ...item, is_active: !item.is_active })
          showToast(item.is_active ? 'Slider dinonaktifkan.' : 'Slider diaktifkan.')
          loadSliders()
        },
      },
      bulk_delete: {
        type: 'danger',
        title: 'Hapus Massal',
        message: `Anda akan memindahkan ${selectedItems.length} item ke Sampah. Lanjutkan?`,
        action: async () => {
          await slidersService.bulkDelete(selectedItems)
          setSelectedItems([])
          showToast('Slider berhasil dihapus massal.')
          loadSliders()
        },
      },
      bulk_restore: {
        type: 'info',
        title: 'Pulihkan Massal',
        message: `Anda akan memulihkan ${selectedItems.length} item dari Sampah. Lanjutkan?`,
        action: async () => {
          await slidersService.bulkRestore(selectedItems)
          setSelectedItems([])
          showToast('Slider berhasil dipulihkan massal.')
          loadSliders()
        },
      },
    }
    const cfg = configs[type]
    if (!cfg) return
    setConfirm({ isOpen: true, type: cfg.type, title: cfg.title, message: cfg.message, action: cfg.action })
  }

  async function executeConfirm() {
    const action = confirm.action
    setConfirm(prev => ({ ...prev, isOpen: false, action: null }))
    if (action) {
      try { await action() } catch (err) { showToast(err.message || 'Terjadi kesalahan', 'error') }
    }
  }

  const headerContent = (
    <div className="flex items-center gap-2 w-full max-w-2xl animate-fade-in-up">
      <div className="relative w-full">
        <i className="ph ph-magnifying-glass absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
        <input
          type="text"
          value={searchQuery}
          onChange={e => setSearchQuery(e.target.value)}
          placeholder="Cari slider..."
          className="w-full pl-9 pr-4 py-2 bg-gray-50 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-brand-500 focus:bg-white transition-colors"
        />
      </div>
      <select
        value={filterActive}
        onChange={e => setFilterActive(e.target.value)}
        className="shrink-0 bg-gray-50 border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-brand-500 focus:bg-white transition-colors cursor-pointer"
      >
        <option value="">Semua Status</option>
        <option value="active">Aktif</option>
        <option value="inactive">Non-aktif</option>
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
    <AdminLayout title="Manajemen Sliders" headerContent={headerContent}>
      <ConfirmDialog {...confirm} onClose={() => setConfirm(prev => ({ ...prev, isOpen: false, action: null }))} onConfirm={executeConfirm} />
      <ToastNotification show={toast.show} message={toast.message} type={toast.type} />

      <div className="space-y-6 animate-fade-in-up">
        {/* Tabs + Tambah */}
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
          <div className="flex bg-white rounded-lg p-1 border border-gray-200 shadow-sm">
            <button
              onClick={() => setCurrentTab('active')}
              className={`px-4 py-1.5 rounded-md text-sm transition-all duration-200 ${currentTab === 'active' ? 'bg-brand-50 text-brand-600 shadow-sm font-medium' : 'text-gray-500 hover:text-gray-700'}`}
            >
              Daftar Aktif
            </button>
            <button
              onClick={() => setCurrentTab('trash')}
              className={`px-4 py-1.5 rounded-md text-sm transition-all duration-200 flex items-center gap-2 ${currentTab === 'trash' ? 'bg-red-50 text-red-600 shadow-sm font-medium' : 'text-gray-500 hover:text-gray-700'}`}
            >
              <i className="ph ph-trash" /> Sampah (History)
            </button>
          </div>
          {currentTab === 'active' && (
            <button
              onClick={() => openForm()}
              className="shrink-0 bg-brand-600 hover:bg-brand-700 text-white px-4 py-2 rounded-lg text-sm font-medium flex items-center gap-2 transition-all btn-press shadow-sm"
            >
              <i className="ph ph-plus-circle text-lg" /> Tambah
            </button>
          )}
        </div>

        {/* Bulk Actions Bar */}
        {selectedItems.length > 0 && (
          <div className="bg-indigo-50 border border-indigo-100 rounded-lg p-3 flex items-center justify-between shadow-sm">
            <span className="text-sm text-indigo-800 font-medium">{selectedItems.length} item terpilih</span>
            <div className="flex gap-2">
              {currentTab === 'trash' ? (
                <button
                  onClick={() => confirmAction('bulk_restore')}
                  className="bg-emerald-100 text-emerald-700 hover:bg-emerald-200 px-3 py-1.5 rounded-md text-sm font-medium flex items-center gap-2 transition-colors"
                >
                  <i className="ph-bold ph-arrow-counter-clockwise" /> Pulihkan Massal
                </button>
              ) : (
                <button
                  onClick={() => confirmAction('bulk_delete')}
                  className="bg-red-100 text-red-700 hover:bg-red-200 px-3 py-1.5 rounded-md text-sm font-medium flex items-center gap-2 transition-colors"
                >
                  <i className="ph-bold ph-trash" /> Hapus Massal
                </button>
              )}
            </div>
          </div>
        )}

        {loading && <div className="text-slate-500 py-10 text-center">Memuat sliders...</div>}
        {error && <div className="text-red-600 py-10 text-center font-medium">{error}</div>}

        {!loading && !error && (
          <div className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
            <div className="overflow-x-auto">
              <table className="w-full text-left border-collapse">
                <thead>
                  <tr className="bg-gray-50 border-b border-gray-200 text-xs uppercase tracking-wider text-gray-500 font-semibold">
                    <th className="p-4 w-12 text-center">
                      <input type="checkbox" checked={isAllSelected} onChange={toggleAll} className="rounded border-gray-300 text-brand-600 focus:ring-brand-500 accent-brand-600" />
                    </th>
                    <th className="p-4">Tampilan / Judul</th>
                    <th className="p-4 w-24">Urutan</th>
                    <th className="p-4">Status</th>
                    <th className="p-4 text-right">Aksi</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-100 text-sm text-gray-700">
                  {paginatedItems.map(item => (
                    <tr key={item.id} className="hover:bg-gray-50 transition-colors group admin-row">
                      <td className="p-4 text-center">
                        <input
                          type="checkbox"
                          checked={selectedItems.includes(item.id)}
                          onChange={() => toggleOne(item.id)}
                          className="rounded border-gray-300 text-brand-600 focus:ring-brand-500 accent-brand-600"
                        />
                      </td>
                      <td className="p-4">
                        <div className="flex items-start gap-4">
                          {item.image_url ? (
                            <img src={resolveAssetUrl(item.image_url)} alt={item.title} className="w-32 h-16 rounded object-cover border border-gray-200 shrink-0" />
                          ) : (
                            <div className="w-32 h-16 rounded bg-gray-100 border border-gray-200 shrink-0 flex items-center justify-center">
                              <i className="ph ph-image text-gray-300 text-2xl" />
                            </div>
                          )}
                          <div>
                            <p className="font-medium text-gray-900 leading-snug">{item.title}</p>
                            {item.subtitle && <p className="text-xs text-gray-500 mt-1">{item.subtitle}</p>}
                            {item.is_new && (
                              <span className="inline-block mt-1 bg-brand-50 text-brand-600 text-[10px] px-2 py-0.5 rounded-full font-medium">NEW</span>
                            )}
                          </div>
                        </div>
                      </td>
                      <td className="p-4">
                        <div className="w-16">
                          <input
                            type="number"
                            defaultValue={item.sort_order}
                            onBlur={e => handleSortChange(item, e.target.value)}
                            onKeyDown={e => { if (e.key === 'Enter') e.target.blur() }}
                            className="w-full px-2 py-1 text-sm border border-gray-300 rounded focus:ring-brand-500 focus:border-brand-500 text-center outline-none"
                          />
                        </div>
                      </td>
                      <td className="p-4">
                        {currentTab === 'active' ? (
                          <button
                            onClick={() => confirmAction('toggle_publish', item)}
                            className="relative inline-flex items-center cursor-pointer"
                            title={item.is_active ? 'Nonaktifkan' : 'Aktifkan'}
                          >
                            <span className={`w-9 h-5 rounded-full transition-colors relative ${item.is_active ? 'bg-brand-500' : 'bg-gray-200'}`}>
                              <span className={`absolute top-[2px] left-[2px] h-4 w-4 bg-white border rounded-full transition-transform ${item.is_active ? 'translate-x-4 border-white' : 'border-gray-300'}`} />
                            </span>
                            <span className={`ml-2 text-xs font-medium ${item.is_active ? 'text-brand-600' : 'text-gray-400'}`}>
                              {item.is_active ? 'Aktif' : 'Non-aktif'}
                            </span>
                          </button>
                        ) : (
                          <span className="inline-flex items-center gap-1 bg-red-50 text-red-600 text-xs px-2 py-1 rounded-md font-medium border border-red-100">
                            <i className="ph ph-trash" /> Terhapus
                          </span>
                        )}
                      </td>
                      <td className="p-4 text-right">
                        <div className="flex items-center justify-end gap-2">
                          {currentTab === 'active' ? (
                            <>
                              <button onClick={() => openForm(item)} className="p-1.5 text-gray-500 hover:text-brand-600 hover:bg-brand-50 rounded" title="Edit">
                                <i className="ph ph-pencil-simple text-lg" />
                              </button>
                              <button onClick={() => confirmAction('delete', item)} className="p-1.5 text-gray-500 hover:text-red-600 hover:bg-red-50 rounded" title="Hapus (Soft Delete)">
                                <i className="ph ph-trash text-lg" />
                              </button>
                            </>
                          ) : (
                            <button onClick={() => confirmAction('restore', item)} className="p-1.5 text-gray-500 hover:text-emerald-600 hover:bg-emerald-50 rounded" title="Pulihkan">
                              <i className="ph ph-arrow-counter-clockwise text-lg" />
                            </button>
                          )}
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
 
              {/* Empty State */}
              {filteredItems.length === 0 && (
                <div className="py-16 text-center text-slate-500 flex flex-col items-center justify-center">
                  <i className="ph ph-image text-gray-300 text-5xl mb-4" />
                  <p className="font-medium text-gray-500">Tidak ada data slider untuk ditampilkan.</p>
                </div>
              )}
            </div>

            {/* Pagination */}
            {filteredItems.length > 0 && (
              <div className="bg-gray-50 border-t border-gray-200 px-4 py-3 flex items-center justify-between sm:px-6 rounded-b-xl">
                <div className="text-sm text-gray-500">
                  Menampilkan <span className="font-medium text-gray-900">{pageStart}</span> sampai <span className="font-medium text-gray-900">{pageEnd}</span> dari <span className="font-medium text-gray-900">{filteredItems.length}</span> hasil
                </div>
                <div className="flex gap-1">
                  <button
                    onClick={() => setCurrentPage(p => Math.max(1, p - 1))}
                    disabled={currentPage === 1}
                    className="px-3 py-1 border border-gray-200 bg-white text-gray-500 rounded hover:bg-gray-50 disabled:opacity-50 text-sm"
                  >
                    Prev
                  </button>
                  {Array.from({ length: totalPages }, (_, i) => i + 1).map(p => (
                    <button
                      key={p}
                      onClick={() => setCurrentPage(p)}
                      className={`px-3 py-1 border rounded text-sm ${currentPage === p ? 'border-brand-500 bg-brand-50 text-brand-600 font-medium' : 'border-gray-200 bg-white text-gray-500 hover:bg-gray-50'}`}
                    >
                      {p}
                    </button>
                  ))}
                  <button
                    onClick={() => setCurrentPage(p => Math.min(totalPages, p + 1))}
                    disabled={currentPage === totalPages}
                    className="px-3 py-1 border border-gray-200 bg-white text-gray-500 rounded hover:bg-gray-50 disabled:opacity-50 text-sm"
                  >
                    Next
                  </button>
                </div>
              </div>
            )}
          </div>
        )}
      </div>

      {/* FORM MODAL (Create/Edit) */}
      {isFormOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
          <div className="fixed inset-0 bg-black/40 backdrop-blur-[2px]" onClick={() => setIsFormOpen(false)} />
          <div className="relative bg-white rounded-2xl shadow-2xl max-w-3xl w-full max-h-[90vh] flex flex-col overflow-hidden z-10">
            <div className="border-b border-slate-200 px-6 py-4 flex items-center justify-between">
              <h3 className="font-heading font-bold text-slate-900 text-lg">
                {formMode === 'create' ? 'Tambah Slider' : 'Edit Slider'}
              </h3>
              <button onClick={() => setIsFormOpen(false)} className="text-slate-400 hover:text-slate-600 p-1">
                <i className="ph-bold ph-x text-lg" />
              </button>
            </div>
            <form onSubmit={handleSubmit} noValidate className="p-6 overflow-y-auto max-h-[calc(90vh-120px)] space-y-5">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-5 animate-fade-in-up">
                <div className="md:col-span-2">
                  <label className="block text-sm font-medium text-gray-700 mb-1">Judul Utama <span className="text-red-500">*</span></label>
                  <input
                    type="text"
                    value={formData.title}
                    onChange={e => {
                      setFormData({ ...formData, title: e.target.value })
                      clearFieldError('title')
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
                    className={`w-full px-3 py-2 border rounded-lg focus:ring-brand-500 focus:border-brand-500 text-sm outline-none transition-colors ${touched.title && formErrors.title ? 'border-red-400 focus:ring-2 focus:ring-red-100' : 'border-gray-300'}`}
                  />
                  {touched.title && formErrors.title && (
                    <p className="text-red-500 text-[11px] font-semibold mt-1 flex items-center gap-1">
                      <i className="ph-bold ph-warning-circle text-xs" /> {formErrors.title}
                    </p>
                  )}
                </div>
                <div className="md:col-span-2">
                  <label className="block text-sm font-medium text-gray-700 mb-1">Sub-judul <span className="text-gray-400 font-normal">(opsional)</span></label>
                  <input
                    type="text"
                    value={formData.subtitle}
                    onChange={e => setFormData({ ...formData, subtitle: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-brand-500 focus:border-brand-500 text-sm outline-none"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Badge Tag <span className="text-gray-400 font-normal">(opsional)</span></label>
                  <input
                    type="text"
                    value={formData.tag}
                    onChange={e => setFormData({ ...formData, tag: e.target.value })}
                    placeholder="Misal: Webinar, Event"
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-brand-500 focus:border-brand-500 text-sm outline-none"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Urutan (Sort Order) <span className="text-red-500">*</span></label>
                  <input
                    type="number"
                    value={formData.sort_order}
                    onChange={e => {
                      const val = e.target.value
                      setFormData({ ...formData, sort_order: val === '' ? '' : Number(val) })
                      clearFieldError('sort_order')
                      if (touched.sort_order) {
                        const errs = validateForm({ ...formData, sort_order: val })
                        setFormErrors(prev => ({ ...prev, sort_order: errs.sort_order }))
                      }
                    }}
                    onBlur={() => {
                      setTouched(prev => ({ ...prev, sort_order: true }))
                      const errs = validateForm()
                      setFormErrors(prev => ({ ...prev, sort_order: errs.sort_order }))
                    }}
                    className={`w-full px-3 py-2 border rounded-lg focus:ring-brand-500 focus:border-brand-500 text-sm outline-none transition-colors ${touched.sort_order && formErrors.sort_order ? 'border-red-400 focus:ring-2 focus:ring-red-100' : 'border-gray-300'}`}
                  />
                  {touched.sort_order && formErrors.sort_order && (
                    <p className="text-red-500 text-[11px] font-semibold mt-1 flex items-center gap-1">
                      <i className="ph-bold ph-warning-circle text-xs" /> {formErrors.sort_order}
                    </p>
                  )}
                </div>
                <div className="md:col-span-2">
                  <label className="block text-sm font-medium text-gray-700 mb-1">Gambar Slider (URL) <span className="text-red-500">*</span></label>
                  <input
                    type="url"
                    value={formData.image_path}
                    onChange={e => {
                      setFormData({ ...formData, image_path: e.target.value })
                      clearFieldError('image_path')
                      if (touched.image_path) {
                        const errs = validateForm({ ...formData, image_path: e.target.value })
                        setFormErrors(prev => ({ ...prev, image_path: errs.image_path }))
                      }
                    }}
                    onBlur={() => {
                      setTouched(prev => ({ ...prev, image_path: true }))
                      const errs = validateForm()
                      setFormErrors(prev => ({ ...prev, image_path: errs.image_path }))
                    }}
                    placeholder="https://example.com/banner.jpg"
                    className={`w-full px-3 py-2 border rounded-lg focus:ring-brand-500 focus:border-brand-500 text-sm outline-none transition-colors ${touched.image_path && formErrors.image_path ? 'border-red-400 focus:ring-2 focus:ring-red-100' : 'border-gray-300'}`}
                  />
                  {touched.image_path && formErrors.image_path && (
                    <p className="text-red-500 text-[11px] font-semibold mt-1 flex items-center gap-1">
                      <i className="ph-bold ph-warning-circle text-xs" /> {formErrors.image_path}
                    </p>
                  )}
                  <p className="text-xs text-gray-500 mt-1">Rekomendasi ukuran: 1920x600 px (Rasio lebar)</p>
                </div>
                <div className="md:col-span-2">
                  <label className="block text-sm font-medium text-gray-700 mb-1">Link URL (Jika diklik) <span className="text-gray-400 font-normal">(opsional)</span></label>
                  <input
                    type="url"
                    value={formData.link_url}
                    onChange={e => setFormData({ ...formData, link_url: e.target.value })}
                    placeholder="https://..."
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-brand-500 focus:border-brand-500 text-sm outline-none"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Tanggal Kegiatan <span className="text-gray-400 font-normal">(opsional di slider)</span></label>
                  <input
                    type="text"
                    value={formData.event_date}
                    onChange={e => setFormData({ ...formData, event_date: e.target.value })}
                    placeholder="20 Okt 2024"
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-brand-500 focus:border-brand-500 text-sm outline-none"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Lokasi <span className="text-gray-400 font-normal">(opsional di slider)</span></label>
                  <input
                    type="text"
                    value={formData.location}
                    onChange={e => setFormData({ ...formData, location: e.target.value })}
                    placeholder="Jakarta"
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-brand-500 focus:border-brand-500 text-sm outline-none"
                  />
                </div>
                <div className="md:col-span-2 flex gap-6">
                  <label className="flex items-center gap-2 cursor-pointer font-medium text-gray-700">
                    <input
                      type="checkbox"
                      checked={formData.is_new}
                      onChange={e => setFormData({ ...formData, is_new: e.target.checked })}
                      className="rounded border-gray-300 text-brand-600 focus:ring-brand-500 accent-brand-600 cursor-pointer"
                    />
                    <span className="text-sm">Tandai sebagai BARU (badge NEW)</span>
                  </label>
                  <label className="flex items-center gap-2 cursor-pointer font-medium text-gray-700">
                    <input
                      type="checkbox"
                      checked={formData.is_active}
                      onChange={e => setFormData({ ...formData, is_active: e.target.checked })}
                      className="rounded border-gray-300 text-brand-600 focus:ring-brand-500 accent-brand-600 cursor-pointer"
                    />
                    <span className="text-sm">Slider Aktif</span>
                  </label>
                </div>
              </div>
              <div className="flex justify-end gap-2 pt-4 border-t items-center mt-6">
                {isLimited && (
                  <span className="text-xs text-amber-600 font-semibold mr-auto flex items-center gap-1">
                    <i className="ph ph-timer text-sm" /> Tunggu {cooldown}s
                  </span>
                )}
                <button
                  type="button"
                  onClick={() => setIsFormOpen(false)}
                  disabled={isLimited}
                  className="px-4 py-2 border rounded-xl text-sm font-semibold hover:bg-slate-50 disabled:opacity-50 transition-colors"
                >
                  Batal
                </button>
                <button
                  type="submit"
                  disabled={isLimited}
                  className="px-5 py-2 bg-brand-600 text-white rounded-xl text-sm font-semibold hover:bg-brand-700 disabled:opacity-60 transition-colors flex items-center gap-2"
                >
                  Simpan
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </AdminLayout>
  )
}
