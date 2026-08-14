import { useState, useEffect, useCallback } from 'react'
import AdminLayout from '../../layouts/AdminLayout'
import { kegiatanService } from '../../services/kegiatanService'
import ConfirmDialog from '../../components/admin/ConfirmDialog'
import ToastNotification from '../../components/admin/ToastNotification'
import { useFormErrors, useRateLimitCooldown } from '../../utils/parseApiError'
import { resolveAssetUrl } from '../../utils/assetUrl'

const PAGE_SIZE = 10

const getTodayDateString = () => {
  const months = [
    'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
    'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
  ]
  const today = new Date()
  const day = String(today.getDate()).padStart(2, '0')
  const month = months[today.getMonth()]
  const year = today.getFullYear()
  return `${day} ${month} ${year}`
}

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

  const [categories, setCategories] = useState(['Kegiatan', 'Munas', 'Pelatihan'])
  const [categorySearch, setCategorySearch] = useState('')
  const [showCategoryDropdown, setShowCategoryDropdown] = useState(false)

  const [formData, setFormData] = useState({
    id: null,
    title: '',
    category: 'Kegiatan',
    organizer: 'DPP GRADASI',
    eventDate: getTodayDateString(),
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
  // Error backend: field errors inline + countdown rate limit
  const { fieldErrors, applyError, clearFieldError } = useFormErrors()
  const { cooldown, isLimited, applyRateLimit } = useRateLimitCooldown()

  const [formErrors, setFormErrors] = useState({})
  const [touched, setTouched] = useState({})

  // Modal Detail (read-only)
  const [isDetailOpen, setIsDetailOpen] = useState(false)
  const [detailData, setDetailData] = useState(null)
  const [detailLoading, setDetailLoading] = useState(false)

  const validateForm = useCallback((data = formData) => {
    const errors = {}
    if (!data.title || !data.title.trim()) {
      errors.title = 'Judul kegiatan wajib diisi.'
    } else if (data.title.trim().length < 5) {
      errors.title = 'Judul minimal 5 karakter.'
    } else if (data.title.trim().length > 300) {
      errors.title = 'Judul maksimal 300 karakter.'
    }
    if (!data.eventDate || !data.eventDate.trim()) {
      errors.eventDate = 'Tanggal event wajib diisi.'
    }
    if (!data.content || !data.content.trim()) {
      errors.content = 'Konten lengkap wajib diisi.'
    }
    return errors
  }, [formData])

  // Galeri gambar kegiatan
  const [gallery, setGallery] = useState([]) // [{id?, image_path, caption, sort_order}]
  const [galleryUploading, setGalleryUploading] = useState(false)
  const [coverUploading, setCoverUploading] = useState(false)

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
    setFormErrors({})
    setTouched({})
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
        tags: Array.isArray(item.tags) ? item.tags.join(', ') : (item.tags ?? ''),
        isPublished: item.is_published !== false
      })
      // list_admin tidak menyertakan content/tags/gallery — ambil detail penuh dulu
      kegiatanService.detailById(item.id)
        .then(res => {
          if (res?.data) {
            const d = res.data
            console.log("Detail Kegiatan dari API:", d)
            setFormData(prev => ({
              ...prev,
              content: d.content ?? prev.content,
              excerpt: d.excerpt ?? prev.excerpt,
              eventDate: d.event_date ?? prev.eventDate,
              location: d.location ?? prev.location,
              organizer: d.organizer ?? prev.organizer,
              image: d.image_path || d.image_url || prev.image,
              tags: Array.isArray(d.tags) ? d.tags.join(', ') : (d.tags ?? prev.tags),
            }))
            setCategorySearch(d.category ?? item.category ?? '')
            // Isi galeri dari detail (id sudah ada → bisa dihapus per item via API)
            if (Array.isArray(d.gallery)) {
              setGallery(d.gallery.map((g, idx) => ({
                id: g.id,
                image_path: g.image_path,
                caption: g.caption || '',
                sort_order: g.sort_order ?? idx,
              })))
            }
          }
        })
        .catch(() => {}) // gagal ambil detail → tetap pakai data list
      setCategorySearch(item.category || 'Kegiatan')
    } else {
      setFormMode('create')
      setFormData({
        id: null,
        title: '',
        category: 'Kegiatan',
        organizer: 'DPP GRADASI',
        eventDate: getTodayDateString(),
        location: '',
        image: '',
        excerpt: '',
        content: '',
        tags: '',
        isPublished: true
      })
      setCategorySearch('Kegiatan')
      setGallery([])
    }
    setIsFormOpen(true)
  }

  async function openDetail(item) {
    setIsDetailOpen(true)
    setDetailLoading(true)
    setDetailData(item)
    try {
      const res = await kegiatanService.detailById(item.id)
      if (res?.data) setDetailData(res.data)
    } catch (err) {
      showToast(err?.message || 'Gagal memuat detail kegiatan.', 'error')
    } finally {
      setDetailLoading(false)
    }
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

    const actualCategory = formData.category?.trim() || ''
    if (!actualCategory) {
      showToast('Kategori tidak boleh kosong.', 'error')
      return
    }

    const payload = {
      title: formData.title,
      category: actualCategory,
      organizer: formData.organizer,
      event_date: formData.eventDate,
      location: formData.location,
      image_path: formData.image,
      excerpt: formData.excerpt,
      content: formData.content,
      tags: formData.tags,
      is_published: formData.isPublished,
      // Galeri: item baru (tanpa id) dikirim, item lama dengan id tetap dipertahankan
      gallery: JSON.stringify(gallery.map((g, idx) => ({
        image_path: g.image_path,
        caption: g.caption || '',
        sort_order: g.sort_order ?? idx,
      })))
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
      setGallery([])
      if (actualCategory && !categories.includes(actualCategory)) {
        setCategories(prev => [...prev, actualCategory])
      }
      loadKegiatan()
    } catch (err) {
      const parsed = applyError(err)
      applyRateLimit(err)
      setFormErrors(prev => ({ ...prev, ...parsed.fieldErrors }))
      setTouched(prev => ({ ...prev, ...Object.keys(parsed.fieldErrors).reduce((acc, k) => ({ ...acc, [k]: true }), {}) }))
      if (Object.keys(parsed.fieldErrors).length === 0) {
        showToast(parsed.message || 'Gagal menyimpan kegiatan.', 'error')
      }
    }
  }

  // Upload satu gambar (cover atau galeri)
  const handleImageUpload = async (file, type) => {
    if (!file) return
    if (type === 'cover') setCoverUploading(true)
    else setGalleryUploading(true)
    try {
      const res = await kegiatanService.uploadImage(file)
      if (res?.data?.image_path) {
        const path = res.data.image_path
        if (type === 'cover') {
          setFormData(prev => ({ ...prev, image: path }))
        } else {
          setGallery(prev => [...prev, { id: null, image_path: path, caption: '', sort_order: prev.length }])
        }
        return path
      }
      showToast('Gagal mengunggah gambar.', 'error')
    } catch (err) {
      showToast(err?.message || 'Gagal mengunggah gambar.', 'error')
    } finally {
      if (type === 'cover') setCoverUploading(false)
      else setGalleryUploading(false)
    }
    return null
  }

  // Hapus item galeri: kalau sudah punya id di DB → panggil API delete gallery
  const handleRemoveGallery = async (index) => {
    const item = gallery[index]
    if (!item) return
    if (item.id) {
      try {
        await kegiatanService.removeGallery(item.id)
      } catch (err) {
        showToast(err?.message || 'Gagal menghapus gambar galeri.', 'error')
        return
      }
    }
    setGallery(prev => prev.filter((_, i) => i !== index))
    showToast('Gambar galeri dihapus.')
  }

  function confirmAction(type, id = null, extraData = null) {
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
      toggle_publish: {
        type: extraData ? 'warning' : 'info',
        title: extraData ? 'Jadikan Draft' : 'Terbitkan Kegiatan',
        message: extraData ? 'Kegiatan ini akan diubah menjadi draft. Lanjutkan?' : 'Kegiatan ini akan diterbitkan. Lanjutkan?',
        action: async () => {
          try {
            await kegiatanService.update(id, { is_published: !extraData })
            setItems(prev => prev.map(i => i.id === id ? { ...i, is_published: !extraData } : i))
            showToast(extraData ? 'Kegiatan dijadikan draft.' : 'Kegiatan berhasil diterbitkan!')
          } catch (err) {
            showToast(err?.message || 'Gagal mengubah status kegiatan.', 'error')
          }
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
      {currentTab !== 'trash' && (
        <select
          value={filterStatus}
          onChange={e => { setFilterStatus(e.target.value); setCurrentPage(1); }}
          className="shrink-0 bg-gray-50 border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-brand-500 focus:bg-white transition-colors cursor-pointer"
        >
          <option value="">Semua Status</option>
          <option value="published">Terbit</option>
          <option value="draft">Draft</option>
        </select>
      )}
      {currentTab !== 'trash' && (
        <select
          value={filterSort}
          onChange={e => { setFilterSort(e.target.value); setCurrentPage(1); }}
          className="shrink-0 bg-gray-50 border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-brand-500 focus:bg-white transition-colors cursor-pointer"
        >
          <option value="newest">Terbaru</option>
          <option value="oldest">Terlama</option>
        </select>
      )}
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
                        <img src={resolveAssetUrl(item.image_url)} alt="" className="w-16 h-12 rounded-lg object-cover border border-slate-200 shrink-0 previewable-image" />
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
                        <span className="inline-flex items-center gap-1.5 text-xs font-semibold text-red-500">
                          <i className="ph ph-trash" /> Di Sampah
                        </span>
                      ) : (
                        <button onClick={() => confirmAction('toggle_publish', item.id, item.is_published)} className="inline-flex items-center gap-2 cursor-pointer">
                          <div className={`relative w-9 h-5 rounded-full transition-colors ${item.is_published ? 'bg-brand-500' : 'bg-slate-200'}`}>
                            <div className={`absolute top-[2px] left-[2px] w-4 h-4 bg-white rounded-full transition-transform ${item.is_published ? 'translate-x-4' : 'translate-x-0'}`} />
                          </div>
                          <span className={`text-xs font-semibold ${item.is_published ? 'text-brand-600' : 'text-slate-400'}`}>
                            {item.is_published ? 'Terbit' : 'Draft'}
                          </span>
                        </button>
                      )}
                    </td>
                    <td className="p-4 text-right">
                      <div className="flex justify-end gap-2">
                        <button onClick={() => openDetail(item)} className="p-2 text-slate-400 hover:text-brand-600 rounded-lg" title="Detail">
                          <i className="ph ph-eye text-base" />
                        </button>
                        {currentTab === 'trash' ? (
                          <button onClick={() => confirmAction('restore', item.id)} className="p-2 text-slate-400 hover:text-emerald-600 rounded-lg" title="Pulihkan">
                            <i className="ph ph-arrow-counter-clockwise text-base" /> Pulihkan
                          </button>
                        ) : (
                          <>
                            <button onClick={() => openForm(item)} className="p-2 text-slate-400 hover:text-brand-600 rounded-lg" title="Edit">
                              <i className="ph ph-pencil-simple text-base" />
                            </button>
                            <button onClick={() => confirmAction('delete', item.id)} className="p-2 text-slate-400 hover:text-red-600 rounded-lg" title="Hapus">
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
                {Array.from({ length: totalPages || 1 }, (_, i) => i + 1).map((n) => (
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
          </>
          )}
        </div>
      </div>
         {/* Form Modal */}
      {isFormOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
          <div className="fixed inset-0 bg-black/40 backdrop-blur-[2px]" onClick={() => setIsFormOpen(false)} />
          <div className="relative bg-white rounded-2xl shadow-2xl max-w-3xl w-full max-h-[90vh] flex flex-col overflow-hidden z-10">
            <div className="border-b border-slate-200 px-6 py-4 flex items-center justify-between">
              <h3 className="font-heading font-bold text-slate-900 text-lg">
                {formMode === 'create' ? 'Tambah Kegiatan Baru' : 'Edit Kegiatan'}
              </h3>
              <button onClick={() => setIsFormOpen(false)} className="text-slate-400 hover:text-slate-600 p-1">
                <i className="ph-bold ph-x text-lg" />
              </button>
            </div>
            <form onSubmit={handleSubmit} noValidate className="p-6 overflow-y-auto max-h-[calc(90vh-120px)]">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6 animate-fade-in-up">
                
                {/* Left Column: Meta & Media */}
                <div className="space-y-4">
                  <div>
                    <label className="block text-xs font-semibold text-slate-500 mb-1">Judul Kegiatan <span className="text-red-500">*</span></label>
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
                      className={`w-full px-3.5 py-2.5 border rounded-xl text-sm outline-none transition-colors ${touched.title && formErrors.title ? 'border-red-400 focus:ring-2 focus:ring-red-100' : 'border-slate-300 focus:border-brand-500 focus:ring-2 focus:ring-brand-100'}`}
                    />
                    {touched.title && formErrors.title && (
                      <p className="text-red-500 text-[11px] font-semibold mt-1.5 flex items-center gap-1">
                        <i className="ph-bold ph-warning-circle text-xs" /> {formErrors.title}
                      </p>
                    )}
                  </div>

                  <div className="grid grid-cols-2 gap-4">
                    <div className="col-span-2 sm:col-span-1 relative">
                      <label className="block text-xs font-semibold text-slate-500 mb-1">Kategori <span className="text-gray-400 font-normal">(opsional)</span></label>
                      <input
                        type="text"
                        value={categorySearch}
                        onChange={e => {
                          setCategorySearch(e.target.value)
                          setShowCategoryDropdown(true)
                          setFormData(prev => ({ ...prev, category: e.target.value }))
                        }}
                        onFocus={() => setShowCategoryDropdown(true)}
                        onBlur={() => {
                          setTimeout(() => setShowCategoryDropdown(false), 200)
                          setTouched(prev => ({ ...prev, category: true }))
                        }}
                        onKeyDown={e => {
                          if (e.key === 'Enter') {
                            e.preventDefault()
                            const trimmed = categorySearch.trim()
                            if (trimmed) {
                              setFormData(prev => ({ ...prev, category: trimmed }))
                              setCategorySearch(trimmed)
                              if (!categories.includes(trimmed)) {
                                setCategories(prev => [...prev, trimmed])
                              }
                              setShowCategoryDropdown(false)
                            }
                          }
                        }}
                        placeholder="Ketik atau pilih kategori..."
                        className="w-full px-3.5 py-2.5 border border-slate-300 rounded-xl text-sm outline-none bg-white focus:border-brand-500 focus:ring-2 focus:ring-brand-100 transition-colors"
                      />
                      
                      {showCategoryDropdown && (
                        <div className="absolute left-0 right-0 z-50 mt-1 max-h-48 overflow-y-auto bg-white border border-slate-200 rounded-xl shadow-xl">
                          {categories
                            .filter(c => c.toLowerCase().includes(categorySearch.toLowerCase()))
                            .map(c => (
                              <button
                                key={c}
                                type="button"
                                onMouseDown={() => {
                                  setFormData(prev => ({ ...prev, category: c }))
                                  setCategorySearch(c)
                                  setShowCategoryDropdown(false)
                                }}
                                className="w-full px-4 py-2.5 text-left text-sm text-slate-700 hover:bg-slate-50 transition-colors font-medium"
                              >
                                {c}
                              </button>
                            ))
                          }
                          {categorySearch.trim() && !categories.map(c => c.toLowerCase()).includes(categorySearch.trim().toLowerCase()) && (
                            <button
                              type="button"
                              onMouseDown={() => {
                                const trimmed = categorySearch.trim()
                                setFormData(prev => ({ ...prev, category: trimmed }))
                                setCategorySearch(trimmed)
                                setCategories(prev => [...prev, trimmed])
                                setShowCategoryDropdown(false)
                              }}
                              className="w-full px-4 py-2.5 text-left text-sm text-brand-600 hover:bg-brand-50 font-bold transition-colors border-t border-slate-100"
                            >
                              + Tambah "{categorySearch.trim()}"
                            </button>
                          )}
                        </div>
                      )}
                    </div>

                    <div className="col-span-2 sm:col-span-1">
                      <label className="block text-xs font-semibold text-slate-500 mb-1">Tanggal Event <span className="text-red-500">*</span></label>
                      <input
                        type="text"
                        value={formData.eventDate}
                        placeholder="Contoh: 31 Desember 2025"
                        onChange={e => {
                          setFormData({ ...formData, eventDate: e.target.value })
                          clearFieldError('event_date')
                          if (touched.eventDate) {
                            const errs = validateForm({ ...formData, eventDate: e.target.value })
                            setFormErrors(prev => ({ ...prev, eventDate: errs.eventDate }))
                          }
                        }}
                        onBlur={() => {
                          setTouched(prev => ({ ...prev, eventDate: true }))
                          const errs = validateForm()
                          setFormErrors(prev => ({ ...prev, eventDate: errs.eventDate }))
                        }}
                        className={`w-full px-3.5 py-2.5 border rounded-xl text-sm outline-none transition-colors ${touched.eventDate && formErrors.eventDate ? 'border-red-400 focus:ring-2 focus:ring-red-100' : 'border-slate-300 focus:border-brand-500 focus:ring-2 focus:ring-brand-100'}`}
                      />
                      {touched.eventDate && formErrors.eventDate && (
                        <p className="text-red-500 text-[11px] font-semibold mt-1.5 flex items-center gap-1">
                          <i className="ph-bold ph-warning-circle text-xs" /> {formErrors.eventDate}
                        </p>
                      )}
                    </div>

                    <div className="col-span-2">
                      <label className="block text-xs font-semibold text-slate-500 mb-1">Lokasi <span className="text-gray-400 font-normal">(opsional)</span></label>
                      <input type="text" value={formData.location} onChange={e => setFormData({ ...formData, location: e.target.value })} placeholder="Jakarta" className="w-full px-3.5 py-2.5 border border-slate-300 rounded-xl text-sm outline-none focus:border-brand-500 focus:ring-2 focus:ring-brand-100 transition-colors" />
                    </div>
                  </div>

                  <div>
                    <label className="block text-xs font-semibold text-slate-500 mb-1">Gambar Cover <span className="text-gray-400 font-normal">(opsional)</span></label>
                    <div className="flex items-center gap-3">
                      {formData.image && (
                        <img src={resolveAssetUrl(formData.image)} alt="Cover" className="w-24 h-16 rounded-lg object-cover border border-slate-200 shrink-0" />
                      )}
                      <label className="inline-flex items-center gap-2 px-4 py-2 bg-brand-600 text-white rounded-lg hover:bg-brand-700 text-sm font-semibold cursor-pointer transition shrink-0">
                        <i className="ph-bold ph-upload-simple" />
                        {coverUploading ? 'Mengunggah...' : (formData.image ? 'Ganti Cover' : 'Upload Cover')}
                        <input
                          type="file"
                          accept="image/png,image/jpeg,image/webp"
                          className="hidden"
                          disabled={coverUploading}
                          onChange={e => {
                            const file = e.target.files?.[0]
                            if (file) handleImageUpload(file, 'cover')
                            e.target.value = ''
                          }}
                        />
                      </label>
                      {formData.image && (
                        <button
                          type="button"
                          onClick={() => setFormData({ ...formData, image: '' })}
                          className="text-xs text-red-500 hover:text-red-700 font-medium"
                        >
                          Hapus
                        </button>
                      )}
                    </div>
                    {fieldErrors.image && (
                      <p className="text-red-500 text-[11px] font-semibold mt-1.5 flex items-center gap-1">
                        <i className="ph-bold ph-warning-circle text-xs" /> {fieldErrors.image}
                      </p>
                    )}
                  </div>
<br/>
                  <div>
                    <label className="block text-xs font-semibold text-slate-500 mb-1">Tags (pisahkan dengan koma) <span className="text-gray-400 font-normal">(opsional)</span></label>
                    <input type="text" value={formData.tags || ''} onChange={e => setFormData({ ...formData, tags: e.target.value })} placeholder="kegiatan, pemuda, digital" className="w-full px-3.5 py-2.5 border border-slate-300 rounded-xl text-sm outline-none focus:border-brand-500 focus:ring-2 focus:ring-brand-100 transition-colors" />
                  </div>

                  {/* Galeri Gambar Kegiatan */}
                  <div>
                    <label className="block text-xs font-semibold text-slate-500 mb-2">Galeri Gambar <span className="text-gray-400 font-normal">(opsional)</span> ({gallery.length})</label>
                    {gallery.length > 0 && (
                      <div className="grid grid-cols-3 gap-3 mb-3">
                        {gallery.map((g, idx) => (
                          <div key={g.id || `new-${idx}`} className="relative group rounded-xl overflow-hidden border border-slate-200">
                            <img src={resolveAssetUrl(g.image_path)} alt={`Galeri ${idx + 1}`} className="w-full h-24 object-cover" />
                            <input
                              type="text"
                              value={g.caption || ''}
                              placeholder="Caption (opsional)"
                              onChange={e => {
                                const next = [...gallery]
                                next[idx] = { ...g, caption: e.target.value }
                                setGallery(next)
                              }}
                              className="absolute bottom-0 left-0 right-0 px-2 py-1 bg-black/60 text-white text-[10px] outline-none placeholder:text-white/60"
                            />
                            <button
                              type="button"
                              onClick={() => handleRemoveGallery(idx)}
                              className="absolute top-1.5 right-1.5 w-6 h-6 flex items-center justify-center rounded-full bg-red-600 text-white text-xs transition-opacity shadow"
                              title="Hapus gambar"
                            >
                              <i className="ph-bold ph-x" />
                            </button>
                          </div>
                        ))}
                      </div>
                    )}
                    <label className="inline-flex items-center gap-2 px-4 py-2 border border-dashed border-slate-300 text-slate-600 hover:border-brand-500 hover:text-brand-600 rounded-lg text-sm font-semibold cursor-pointer transition">
                      <i className="ph-bold ph-plus" />
                      {galleryUploading ? 'Mengunggah...' : 'Tambah Gambar Galeri'}
                      <input
                        type="file"
                        accept="image/png,image/jpeg,image/webp"
                        className="hidden"
                        disabled={galleryUploading}
                        onChange={e => {
                          const file = e.target.files?.[0]
                          if (file) handleImageUpload(file, 'gallery')
                          e.target.value = ''
                        }}
                      />
                    </label>
                    <p className="text-[11px] text-slate-400 mt-1.5">PNG / JPG / WEBP · maks 5MB. Gambar baru otomatis tersimpan saat kegiatan disimpan.</p>
                  </div>
                </div>

                {/* Right Column: Descriptions & Content */}
                <div className="space-y-4">
                  <div>
                    <label className="block text-xs font-semibold text-slate-500 mb-1">Ringkasan <span className="text-gray-400 font-normal">(opsional)</span></label>
                    <textarea rows={3} value={formData.excerpt} onChange={e => { setFormData({ ...formData, excerpt: e.target.value }); clearFieldError('excerpt') }} className="w-full px-3.5 py-2.5 border border-slate-300 rounded-xl text-sm outline-none focus:border-brand-500 focus:ring-2 focus:ring-brand-100 transition-colors overflow-y-auto resize-y min-h-[80px]" />
                  </div>

                  <div>
                    <label className="block text-xs font-semibold text-slate-500 mb-1">Konten Lengkap <span className="text-red-500">*</span></label>
                    <textarea
                      rows={8}
                      value={formData.content}
                      onChange={e => {
                        setFormData({ ...formData, content: e.target.value })
                        clearFieldError('content')
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
                      className={`w-full px-3.5 py-2.5 border rounded-xl text-sm outline-none transition-colors overflow-y-auto resize-y min-h-[160px] ${touched.content && formErrors.content ? 'border-red-400 focus:ring-2 focus:ring-red-100' : 'border-slate-300 focus:border-brand-500 focus:ring-2 focus:ring-brand-100'}`}
                    />
                    {touched.content && formErrors.content && (
                      <p className="text-red-500 text-[11px] font-semibold mt-1.5 flex items-center gap-1">
                        <i className="ph-bold ph-warning-circle text-xs" /> {formErrors.content}
                      </p>
                    )}
                  </div>

                  <div>
                    <label className="block text-xs font-semibold text-slate-500 mb-2">Status Publikasi</label>
                    <div className="flex gap-4">
                      <label className="flex items-center gap-2 text-sm cursor-pointer font-medium text-slate-700">
                        <input type="radio" name="isPublished" checked={formData.isPublished} onChange={() => setFormData({ ...formData, isPublished: true })} className="accent-brand-600 rounded border-slate-300 text-brand-600 focus:ring-brand-500 cursor-pointer" />
                        Terbit
                      </label>
                      <label className="flex items-center gap-2 text-sm cursor-pointer font-medium text-slate-700">
                        <input type="radio" name="isPublished" checked={!formData.isPublished} onChange={() => setFormData({ ...formData, isPublished: false })} className="accent-slate-400 rounded border-slate-300 text-brand-600 focus:ring-brand-500 cursor-pointer" />
                        Draft
                      </label>
                    </div>
                  </div>
                </div>

              </div>

              <div className="flex justify-end gap-2 pt-4 border-t items-center mt-6">
                {isLimited && (
                  <span className="text-xs text-amber-600 font-semibold mr-auto flex items-center gap-1">
                    <i className="ph ph-timer text-sm" /> Terlalu banyak percobaan. Tunggu {cooldown}s
                  </span>
                )}
                <button type="button" onClick={() => setIsFormOpen(false)} disabled={coverUploading || galleryUploading || isLimited} className="px-4 py-2 border rounded-xl text-sm font-semibold hover:bg-slate-50 disabled:opacity-50 transition-colors">Batal</button>
                <button type="submit" disabled={coverUploading || galleryUploading || isLimited} className="px-5 py-2 bg-brand-600 text-white rounded-xl text-sm font-semibold hover:bg-brand-700 disabled:opacity-60 transition-colors flex items-center gap-2">
                  Simpan
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* DETAIL MODAL (read-only) */}
      {isDetailOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
          <div className="fixed inset-0 bg-black/40 backdrop-blur-[2px]" onClick={() => setIsDetailOpen(false)} />
          <div className="relative bg-white rounded-2xl shadow-2xl max-w-2xl w-full max-h-[90vh] flex flex-col overflow-hidden z-10">
            <div className="border-b border-slate-200 px-6 py-4 flex items-center justify-between shrink-0">
              <h3 className="font-heading font-bold text-slate-900 text-lg">Detail Kegiatan</h3>
              <button onClick={() => setIsDetailOpen(false)} className="text-slate-400 hover:text-slate-600 p-1">
                <i className="ph-bold ph-x text-lg" />
              </button>
            </div>
            <div className="p-6 overflow-y-auto">
              {detailLoading ? (
                <div className="py-16 text-center text-slate-500">Memuat detail...</div>
              ) : detailData ? (
                <div className="space-y-4">
                  {(detailData.image_url || detailData.image_path) && (
                    <img src={resolveAssetUrl(detailData.image_url || detailData.image_path)} alt="" className="w-full h-48 object-cover rounded-xl border border-slate-200" />
                  )}
                  <div>
                    <span className="text-[11px] font-semibold text-slate-400 uppercase tracking-wider">Judul</span>
                    <p className="font-bold text-slate-900 text-lg leading-snug">{detailData.title}</p>
                  </div>
                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <span className="text-[11px] font-semibold text-slate-400 uppercase tracking-wider">Kategori</span>
                      <p className="text-slate-700 text-sm">{detailData.category || '-'}</p>
                    </div>
                    <div>
                      <span className="text-[11px] font-semibold text-slate-400 uppercase tracking-wider">Tanggal Event</span>
                      <p className="text-slate-700 text-sm">{detailData.event_date || '-'}</p>
                    </div>
                    <div>
                      <span className="text-[11px] font-semibold text-slate-400 uppercase tracking-wider">Lokasi</span>
                      <p className="text-slate-700 text-sm">{detailData.location || '-'}</p>
                    </div>
                    <div>
                      <span className="text-[11px] font-semibold text-slate-400 uppercase tracking-wider">Penyelenggara</span>
                      <p className="text-slate-700 text-sm">{detailData.organizer || '-'}</p>
                    </div>
                    <div>
                      <span className="text-[11px] font-semibold text-slate-400 uppercase tracking-wider">Penulis</span>
                      <p className="text-slate-700 text-sm">{detailData.author_name || 'Admin'}</p>
                    </div>
                    <div>
                      <span className="text-[11px] font-semibold text-slate-400 uppercase tracking-wider">Dilihat</span>
                      <p className="text-slate-700 text-sm">{detailData.views ?? 0} kali</p>
                    </div>
                  </div>
                  {detailData.excerpt && (
                    <div>
                      <span className="text-[11px] font-semibold text-slate-400 uppercase tracking-wider">Ringkasan</span>
                      <p className="text-slate-700 text-sm leading-relaxed whitespace-pre-line">{detailData.excerpt}</p>
                    </div>
                  )}
                  <div>
                    <span className="text-[11px] font-semibold text-slate-400 uppercase tracking-wider">Konten Lengkap</span>
                    <p className="text-slate-700 text-sm leading-relaxed whitespace-pre-line break-words">{detailData.content || '-'}</p>
                  </div>
                  {Array.isArray(detailData.gallery) && detailData.gallery.length > 0 && (
                    <div>
                      <span className="text-[11px] font-semibold text-slate-400 uppercase tracking-wider">Galeri Foto</span>
                      <div className="grid grid-cols-3 gap-2 mt-1.5">
                        {detailData.gallery.map(g => (
                          <img key={g.id} src={resolveAssetUrl(g.image_path)} alt={g.caption || ''} className="w-full h-20 object-cover rounded-lg border border-slate-200" />
                        ))}
                      </div>
                    </div>
                  )}
                  {Array.isArray(detailData.tags) && detailData.tags.length > 0 && (
                    <div>
                      <span className="text-[11px] font-semibold text-slate-400 uppercase tracking-wider">Tags</span>
                      <div className="flex flex-wrap gap-1.5 mt-1">
                        {detailData.tags.map(tag => (
                          <span key={tag} className="bg-slate-100 text-slate-600 text-[11px] font-semibold px-2.5 py-1 rounded-full">{tag}</span>
                        ))}
                      </div>
                    </div>
                  )}
                  <div>
                    <span className="text-[11px] font-semibold text-slate-400 uppercase tracking-wider">Status</span>
                    <p className="text-slate-700 text-sm">{detailData.is_published ? 'Terbit' : 'Draft'}</p>
                  </div>
                </div>
              ) : (
                <div className="py-16 text-center text-slate-500">Data tidak ditemukan.</div>
              )}
            </div>
          </div>
        </div>
      )}
    </AdminLayout>
  )
}
