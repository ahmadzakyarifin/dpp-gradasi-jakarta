import { useState, useEffect, useCallback, useRef } from 'react'
import AdminLayout from '../../layouts/AdminLayout'
import { beritaService } from '../../services/beritaService'
import { beritaContent } from '../../content/beritaContent'
import ConfirmDialog from '../../components/admin/ConfirmDialog'
import ToastNotification from '../../components/admin/ToastNotification'
import { useFormErrors, useRateLimitCooldown } from '../../utils/parseApiError'
import { resolveAssetUrl } from '../../utils/assetUrl'

const getTodayIndonesian = () => {
  const months = ['Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni', 'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'];
  const d = new Date();
  return `${d.getDate()} ${months[d.getMonth()]} ${d.getFullYear()}`;
}

const PAGE_SIZE = 30

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
  const categoryDropdownRef = useRef(null)
  const [showCategoryDropdown, setShowCategoryDropdown] = useState(false)
  const [isAddingCategory, setIsAddingCategory] = useState(false)
  const [newCategoryName, setNewCategoryName] = useState('')
  const [editingCategory, setEditingCategory] = useState(null)
  const [editCategoryName, setEditCategoryName] = useState('')

  useEffect(() => {
    function handleClickOutside(event) {
      if (categoryDropdownRef.current && !categoryDropdownRef.current.contains(event.target)) {
        setShowCategoryDropdown(false)
        setIsAddingCategory(false)
        setNewCategoryName('')
        setEditingCategory(null)
        setEditCategoryName('')
      }
    }
    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])
  const [isFormOpen, setIsFormOpen] = useState(false)
  const [formMode, setFormMode] = useState('create')
  const [formLoading, setFormLoading] = useState(false)
  const [formData, setFormData] = useState({
    id: null,
    title: '',
    category: 'Berita Nasional',
    published_date: getTodayIndonesian(),
    image_path: '',
    excerpt: '',
    content: '',
    tags: '',
    footnote: '',
    image_source: '',
    is_published: true,
  })
  const [formErrors, setFormErrors] = useState({})
  const [touched, setTouched] = useState({})

  // Modal Detail (read-only)
  const [isDetailOpen, setIsDetailOpen] = useState(false)
  const [detailData, setDetailData] = useState(null)
  const [detailLoading, setDetailLoading] = useState(false)

  // Validasi custom — Bahasa Indonesia, per-field
  const validateForm = useCallback((data = formData) => {
    const errors = {}
    if (!data.title || !data.title.trim()) {
      errors.title = 'Judul berita wajib diisi.'
    } else if (data.title.trim().length < 5) {
      errors.title = 'Judul minimal 5 karakter.'
    } else if (data.title.trim().length > 250) {
      errors.title = 'Judul maksimal 250 karakter.'
    }
    if (!data.category || !data.category.trim()) {
      errors.category = 'Kategori wajib dipilih.'
    }
    if (!data.published_date || !data.published_date.trim()) {
      errors.published_date = 'Tanggal terbit wajib diisi.'
    }
    if (!data.content || !data.content.trim()) {
      errors.content = 'Konten lengkap wajib diisi.'
    }
    if (data.tags && data.tags.trim()) {
      if (data.tags.trim().length > 200) {
        errors.tags = 'Tags maksimal 200 karakter keseluruhan.'
      }
    }
    if (data.image_source && data.image_source.trim().length > 150) {
      errors.image_source = 'Keterangan foto maksimal 150 karakter.'
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
  // Error dari backend: field errors inline + countdown rate limit
  const { applyError, clearFieldError } = useFormErrors()
  const { cooldown, isLimited, applyRateLimit } = useRateLimitCooldown()
  // Upload foto cover
  const [imageUploading, setImageUploading] = useState(false)

  const showToast = useCallback((message, type = 'success') => {
    setToast({ show: true, message, type })
  }, [])

  // Upload gambar cover berita → set formData.image_path dari path relatif backend
  const handleImageUpload = async (file) => {
    if (!file) return
    setImageUploading(true)
    try {
      const res = await beritaService.uploadImage(file)
      if (res?.data?.image_path) {
        setFormData(prev => ({ ...prev, image_path: res.data.image_path }))
        showToast('Gambar berhasil diunggah.')
      } else {
        showToast('Gagal mengunggah gambar.', 'error')
      }
    } catch (err) {
      showToast(err?.message || 'Gagal mengunggah gambar.', 'error')
    } finally {
      setImageUploading(false)
    }
  }

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

  // Helper untuk mengecek apakah berita masih baru (7 hari dari update terakhir)
  const isActuallyNew = (item) => {
    if (!item || !item.is_new || !item.updated_at) return false
    const updateTime = new Date(item.updated_at).getTime()
    const now = new Date().getTime()
    const diffDays = (now - updateTime) / (1000 * 60 * 60 * 24)
    return diffDays <= 7
  }

  async function openForm(item = null) {
    setFormErrors({})
    setTouched({})
    if (item) {
      setFormMode('edit')
      // List admin (tabel) tidak membawa field content & tags, jadi ambil detail
      // lengkap dari server supaya form edit terisi penuh (termasuk konten & tags).
      setFormLoading(true)
      try {
        const res = await beritaService.detailById(item.id)
        const d = res?.data || item
        setFormData({
          id: d.id ?? item.id,
          title: d.title ?? item.title ?? '',
          category: d.category ?? item.category ?? '',
          published_date: d.published_date ?? item.published_date ?? '',
          image_path: d.image_path ?? item.image_path ?? '',
          excerpt: d.excerpt ?? item.excerpt ?? '',
          content: d.content ?? '',
          tags: Array.isArray(d.tags) ? d.tags.join(', ') : (d.tags ?? ''),
          footnote: d.footnote ?? item.footnote ?? '',
          image_source: d.image_source ?? item.image_source ?? '',
          is_new: d.is_new && isActuallyNew(d),
          is_published: typeof d.is_published === 'boolean' ? d.is_published : !!item.is_published,
        })
        setIsFormOpen(true)
      } catch (err) {
        showToast(err?.message || 'Gagal memuat detail berita.', 'error')
      } finally {
        setFormLoading(false)
      }
    } else {
      setFormMode('create')
      const defaultCat = categories[0] || 'Berita Nasional'
      setFormData({
        id: null,
        title: '',
        category: defaultCat,
        published_date: getTodayIndonesian(),
        image_path: '',
        excerpt: '',
        content: '',
        tags: '',
        footnote: '',
        image_source: '',
        is_new: true, // Default to true (Otomatis pudar dalam 7 hari)
        is_published: false, // Default false sesuai request
      })
      setIsFormOpen(true)
    }
  }

  async function openDetail(item) {
    setIsDetailOpen(true)
    setDetailLoading(true)
    setDetailData(item)
    try {
      const res = await beritaService.detailById(item.id)
      if (res?.data) setDetailData(res.data)
    } catch (err) {
      showToast(err?.message || 'Gagal memuat detail berita.', 'error')
    } finally {
      setDetailLoading(false)
    }
  }

  async function handleSaveEditCategory(oldName) {
    const trimmed = editCategoryName.trim()
    if (!trimmed || trimmed === oldName) {
      setEditingCategory(null)
      return
    }
    try {
      await beritaService.renameCategory(oldName, trimmed)
      showToast('Kategori berhasil diubah.', 'success')
      setCategories(prev => prev.map(c => c === oldName ? trimmed : c))
      if (formData.category === oldName) {
        setFormData(prev => ({ ...prev, category: trimmed }))
      }
      setEditingCategory(null)
      loadBerita() // Reload to update table
    } catch (err) {
      showToast(err?.message || 'Gagal mengubah kategori.', 'error')
    }
  }

  function confirmDeleteCategory(name) {
    setConfirm({
      isOpen: true,
      type: 'danger',
      title: 'Hapus Kategori',
      message: `Anda yakin ingin menghapus kategori "${name}"? Berita dengan kategori ini akan dipindahkan ke "Berita Organisasi".`,
      action: async () => {
        try {
          await beritaService.deleteCategory(name)
          showToast('Kategori berhasil dihapus.', 'success')
          setCategories(prev => prev.filter(c => c !== name))
          if (formData.category === name) {
            setFormData(prev => ({ ...prev, category: 'Berita Organisasi' }))
          }
          setConfirm({ isOpen: false })
          loadBerita() // Reload to update table
        } catch (err) {
          showToast(err?.message || 'Gagal menghapus kategori.', 'error')
        }
      }
    })
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
    const actualCategory = formData.category?.trim() || ''
    if (!actualCategory) {
      showToast('Kategori tidak boleh kosong.', 'error')
      return
    }

    const payload = {
      ...formData,
      category: actualCategory
    }

    setFormLoading(true)
    try {
      if (formMode === 'create') {
        await beritaService.create(payload)
        showToast('Berita berhasil dibuat.')
      } else {
        await beritaService.update(formData.id, payload)
        showToast('Berita berhasil diperbarui.')
      }
      setIsFormOpen(false)
      if (actualCategory && !categories.includes(actualCategory)) {
        setCategories(prev => [...prev, actualCategory])
      }
      await loadBerita()
    } catch (err) {
      // Error validasi/bisnis dari backend → field errors inline + toast ringkas
      const parsed = applyError(err)
      applyRateLimit(err)
      setFormErrors(prev => ({ ...prev, ...parsed.fieldErrors }))
      setTouched(prev => ({ ...prev, ...Object.keys(parsed.fieldErrors).reduce((acc, k) => ({ ...acc, [k]: true }), {}) }))
      if (Object.keys(parsed.fieldErrors).length === 0) {
        showToast(parsed.message || 'Gagal menyimpan berita.', 'error')
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
    showToast('Filter direset.', 'success')
  }

  const getVisiblePageNumbers = (currentPage, totalPagesCount, maxVisible = 5) => {
    const totalP = totalPagesCount || 1
    if (totalP <= maxVisible) {
      return Array.from({ length: totalP }, (_, i) => i + 1)
    }
    let start = Math.max(1, currentPage - Math.floor(maxVisible / 2))
    let end = start + maxVisible - 1
    if (end > totalP) {
      end = totalP
      start = Math.max(1, end - maxVisible + 1)
    }
    return Array.from({ length: end - start + 1 }, (_, i) => start + i)
  }

  // Filter/search ditangani backend via query param — lihat loadBerita().
  const paginatedItems = items
  const totalPages = meta.total_pages || 1

  const headerContent = (
    <div className="flex items-center gap-2 w-full max-w-3xl animate-fade-in-up">
      <div className="relative w-full">
        <i className="ph ph-magnifying-glass absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
        <input
          type="text"
          value={searchQuery}
          onChange={e => setSearchQuery(e.target.value)}
          placeholder={beritaContent.admin.searchPlaceholder}
          className="w-full pl-9 pr-4 py-2 bg-gray-50 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-brand-500 focus:bg-white transition-colors"
        />
      </div>
      {currentTab !== 'trash' && (
        <select
          value={filterStatus}
          onChange={e => setFilterStatus(e.target.value)}
          className="shrink-0 bg-gray-50 border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-brand-500 focus:bg-white transition-colors cursor-pointer"
        >
          <option value="">{beritaContent.admin.allStatus}</option>
          <option value="published">Terbit</option>
          <option value="draft">Draft</option>
        </select>
      )}
      {currentTab !== 'trash' && (
        <select
          value={filterSort}
          onChange={e => setFilterSort(e.target.value)}
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
    <AdminLayout title="Kelola Berita" headerContent={headerContent}>
      {toast.show && <ToastNotification message={toast.message} type={toast.type} onClose={() => setToast({ ...toast, show: false })} />}
      <ConfirmDialog {...confirm} onClose={() => setConfirm(prev => ({ ...prev, isOpen: false, action: null }))} onConfirm={executeConfirm} />

      <div className="space-y-6 animate-fade-in-up">
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
          <div className="flex bg-white rounded-lg p-1 border border-gray-200 shadow-sm">
            <button
              onClick={() => setCurrentTab('active')}
              className={`px-4 py-1.5 rounded-md text-sm transition-all duration-200 ${currentTab === 'active' ? 'bg-brand-50 text-brand-600 shadow-sm font-medium' : 'text-gray-500 hover:text-gray-700'}`}
            >
              Aktif & Draft
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
              className="bg-brand-600 hover:bg-brand-700 text-white px-4 py-2 rounded-lg text-sm font-medium flex items-center gap-2 transition-all btn-press shadow-sm"
            >
              <i className="ph ph-plus-circle text-lg" /> Tambah
            </button>
          )}
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
                  <th className="p-4">Berita</th>
                  <th className="p-4">Kategori</th>
                  <th className="p-4">Penulis</th>
                  <th className="p-4">Tanggal Terbit</th>
                  <th className="p-4">Status</th>
                  <th className="p-4 text-right">Aksi</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100 text-sm text-gray-700">
                {paginatedItems.map(item => (
                  <tr key={item.id} className="hover:bg-slate-50/70 transition-colors group admin-row">
                    <td className="p-4 text-center">
                      <input type="checkbox" checked={selectedItems.includes(item.id)} onChange={() => toggleItem(item.id)} className="rounded border-slate-300 text-brand-600 focus:ring-brand-500 cursor-pointer" />
                    </td>
                    <td className="p-4">
                      <div className="flex items-center gap-3">
                        {item.image_url || item.image_path ? (
                          <img src={resolveAssetUrl(item.image_url || item.image_path)} alt="" className="w-12 h-9 rounded object-cover border border-slate-200 shrink-0 previewable-image" />
                        ) : (
                          <div className="w-12 h-9 rounded bg-slate-100 border border-slate-200 shrink-0 flex items-center justify-center text-slate-300">
                            <i className="ph ph-image text-base" />
                          </div>
                        )}
                        <div className="min-w-0 flex flex-col items-start">
                          <span className="font-semibold text-slate-900 truncate max-w-[500px] xl:max-w-[750px] block" title={item.title}>
                            {item.title}
                          </span>
                          {item.is_new && isActuallyNew(item) && (
                            <span className="inline-block mt-1 bg-brand-50 text-brand-600 text-[10px] px-2 py-0.5 rounded-full font-bold uppercase tracking-wide">TERBARU</span>
                          )}
                        </div>
                      </div>
                    </td>
                    <td className="p-4">
                      <span className="inline-block bg-slate-100 text-slate-600 text-[10px] font-semibold px-2.5 py-1 rounded-full">
                        {item.category}
                      </span>
                    </td>
                    <td className="p-4 text-slate-600 font-medium text-xs">
                      <div className="flex items-center gap-1.5">
                        <div className="w-6 h-6 rounded-full bg-brand-100 text-brand-700 flex items-center justify-center font-bold text-[10px]">
                          {(item.author_name || 'A').charAt(0).toUpperCase()}
                        </div>
                        {item.author_name || 'Admin'}
                      </div>
                    </td>
                    <td className="p-4 text-slate-500 text-xs whitespace-nowrap">{item.published_date}</td>
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
                      <div className="flex items-center justify-end gap-1">
                        <button onClick={() => openDetail(item)} className="p-2 text-slate-400 hover:text-brand-600 rounded-lg" title="Detail">
                          <i className="ph ph-eye text-base" />
                        </button>
                        {currentTab === 'trash' ? (
                          <button onClick={() => confirmAction('restore', item.id)} className="p-2 text-slate-400 hover:text-emerald-600 rounded-lg" title="Pulihkan">
                            <i className="ph ph-arrow-counter-clockwise text-base" />
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

          {/* PAGINATION (server-driven) */}
          {totalPages >= 1 && items.length > 0 && (
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
                {getVisiblePageNumbers(currentPage, totalPages, 5).map((n) => (
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
            <form onSubmit={handleSubmit} noValidate className="p-6 overflow-y-auto max-h-[calc(90vh-120px)]">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                
                {/* Left Column: Meta & Media */}
                <div className="space-y-4">
                  <div>
                    <label className="flex justify-between items-center text-xs font-semibold text-slate-500 mb-1">
                      <span>Judul Berita <span className="text-red-500">*</span></span>
                      <span className={formData.title.length > 250 ? 'text-red-500' : 'text-slate-400'}>
                        {formData.title.length}/250
                      </span>
                    </label>
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
                      placeholder="Masukkan judul berita..."
                      maxLength={250}
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
                      <label className="block text-xs font-semibold text-slate-500 mb-1">Kategori <span className="text-red-500">*</span></label>
                      <div className="relative" ref={categoryDropdownRef}>
                        <button
                          type="button"
                          onClick={() => {
                            setShowCategoryDropdown(!showCategoryDropdown)
                            setIsAddingCategory(false)
                          }}
                          className={`w-full px-3.5 py-2.5 border rounded-xl text-sm text-left outline-none bg-white transition-colors flex justify-between items-center ${touched.category && formErrors.category ? 'border-red-400 focus:ring-2 focus:ring-red-100' : 'border-slate-300 focus:border-brand-500 focus:ring-2 focus:ring-brand-100'}`}
                        >
                          <span className={formData.category ? 'text-slate-900' : 'text-slate-400'}>
                            {formData.category || 'Pilih kategori...'}
                          </span>
                          <i className={`ph-bold ph-caret-down text-slate-400 transition-transform ${showCategoryDropdown ? 'rotate-180' : ''}`} />
                        </button>
                        
                        {showCategoryDropdown && (
                          <div className="absolute left-0 z-50 mt-1 w-[320px] bg-white border border-slate-200 rounded-xl shadow-xl flex flex-col overflow-hidden">
                            <div className="max-h-48 overflow-y-auto py-1">
                              {categories.map(c => (
                                <div key={c} className="group relative w-full flex flex-col">
                                  {editingCategory === c ? (
                                    <div className="flex gap-2 p-2 bg-brand-50 border-y border-brand-100">
                                      <input
                                        type="text"
                                        autoFocus
                                        value={editCategoryName}
                                        onChange={e => setEditCategoryName(e.target.value)}
                                        className="flex-1 min-w-0 px-2 py-1.5 border border-brand-300 rounded text-sm outline-none focus:ring-1 focus:ring-brand-500"
                                        onKeyDown={e => {
                                          if (e.key === 'Enter') {
                                            e.preventDefault()
                                            handleSaveEditCategory(c)
                                          } else if (e.key === 'Escape') {
                                            setEditingCategory(null)
                                          }
                                        }}
                                      />
                                      <button
                                        type="button"
                                        onClick={() => handleSaveEditCategory(c)}
                                        className="px-2.5 py-1.5 bg-brand-600 hover:bg-brand-700 text-white rounded text-sm font-bold shrink-0"
                                      >
                                        Save
                                      </button>
                                      <button
                                        type="button"
                                        onClick={() => setEditingCategory(null)}
                                        className="px-2 py-1.5 bg-slate-200 hover:bg-slate-300 text-slate-700 rounded text-sm font-bold shrink-0"
                                      >
                                        Batal
                                      </button>
                                    </div>
                                  ) : (
                                    <>
                                      <button
                                        type="button"
                                        onClick={() => {
                                          setFormData(prev => ({ ...prev, category: c }))
                                          setShowCategoryDropdown(false)
                                          setIsAddingCategory(false)
                                          clearFieldError('category')
                                        }}
                                        className={`w-full px-4 py-2.5 text-left text-sm transition-colors font-medium ${formData.category === c ? 'bg-brand-50 text-brand-700' : 'text-slate-700 hover:bg-slate-50'}`}
                                      >
                                        {c}
                                      </button>
                                      <div className="absolute right-2 top-1/2 -translate-y-1/2 opacity-0 group-hover:opacity-100 flex gap-1 transition-opacity">
                                        <button
                                          type="button"
                                          onClick={(e) => {
                                            e.stopPropagation()
                                            setEditingCategory(c)
                                            setEditCategoryName(c)
                                          }}
                                          title="Ubah Nama Kategori Global"
                                          className="p-1.5 text-brand-600 hover:bg-white bg-slate-100 hover:text-brand-700 rounded shadow-sm border border-transparent hover:border-brand-200"
                                        >
                                          <i className="ph-bold ph-pencil-simple" />
                                        </button>
                                        {c !== 'Berita Organisasi' && (
                                          <button
                                            type="button"
                                            onClick={(e) => {
                                              e.stopPropagation()
                                              confirmDeleteCategory(c)
                                            }}
                                            title="Hapus Kategori Global"
                                            className="p-1.5 text-red-500 hover:bg-white bg-slate-100 hover:text-red-600 rounded shadow-sm border border-transparent hover:border-red-200"
                                          >
                                            <i className="ph-bold ph-trash" />
                                          </button>
                                        )}
                                      </div>
                                    </>
                                  )}
                                </div>
                              ))}
                            </div>
                            
                            <div className="border-t border-slate-100 p-2 bg-slate-50">
                              {!isAddingCategory ? (
                                <button
                                  type="button"
                                  onClick={() => setIsAddingCategory(true)}
                                  className="w-full px-3 py-2 text-sm text-brand-600 hover:bg-brand-100/50 rounded-lg font-bold transition-colors flex items-center justify-center gap-1.5"
                                >
                                  <i className="ph-bold ph-plus" /> Tambah Kategori Baru
                                </button>
                              ) : (
                                <div className="flex gap-2">
                                  <input
                                    type="text"
                                    autoFocus
                                    value={newCategoryName}
                                    onChange={e => setNewCategoryName(e.target.value)}
                                    placeholder="Kategori baru..."
                                    className="flex-1 min-w-0 px-3 py-2 border border-slate-300 rounded-lg text-sm outline-none focus:border-brand-500 focus:ring-1 focus:ring-brand-500"
                                    onKeyDown={e => {
                                      if (e.key === 'Enter') {
                                        e.preventDefault()
                                        const trimmed = newCategoryName.trim()
                                        if (trimmed) {
                                          setFormData(prev => ({ ...prev, category: trimmed }))
                                          if (!categories.includes(trimmed)) setCategories(prev => [...prev, trimmed])
                                          setShowCategoryDropdown(false)
                                          setIsAddingCategory(false)
                                          setNewCategoryName('')
                                          clearFieldError('category')
                                        }
                                      }
                                    }}
                                  />
                                  <button
                                    type="button"
                                    onClick={() => {
                                      const trimmed = newCategoryName.trim()
                                      if (trimmed) {
                                        setFormData(prev => ({ ...prev, category: trimmed }))
                                        if (!categories.includes(trimmed)) setCategories(prev => [...prev, trimmed])
                                        setShowCategoryDropdown(false)
                                        setIsAddingCategory(false)
                                        setNewCategoryName('')
                                        clearFieldError('category')
                                      }
                                    }}
                                    className="px-3 py-2 bg-brand-600 hover:bg-brand-700 text-white rounded-lg text-sm font-bold transition-colors shrink-0"
                                  >
                                    Add
                                  </button>
                                </div>
                              )}
                            </div>
                          </div>
                        )}
                      </div>
                      {touched.category && formErrors.category && (
                        <p className="text-red-500 text-[11px] font-semibold mt-1.5 flex items-center gap-1">
                          <i className="ph-bold ph-warning-circle text-xs" /> {formErrors.category}
                        </p>
                      )}
                    </div>
                    <div>
                      <label className="flex justify-between items-center text-xs font-semibold text-slate-500 mb-1">
                        <span>Tanggal Terbit <span className="text-red-500">*</span></span>
                        <span className={(formData.published_date?.length || 0) > 200 ? 'text-red-500' : 'text-slate-400'}>
                          {(formData.published_date?.length || 0)}/200
                        </span>
                      </label>
                      <input
                        type="text"
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
                        maxLength={200}
                        placeholder="Contoh: 30 Desember 2026"
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
                    <label className="block text-xs font-semibold text-slate-500 mb-1">Foto Cover <span className="text-gray-400 font-normal">(opsional)</span></label>
                    <div className="flex items-center gap-3">
                      {formData.image_path && (
                        <img src={resolveAssetUrl(formData.image_path)} alt="Cover" className="w-20 h-14 rounded-lg object-cover border border-slate-200 shrink-0" />
                      )}
                      <label className="inline-flex items-center gap-2 px-4 py-2 bg-brand-600 text-white rounded-lg hover:bg-brand-700 text-sm font-semibold cursor-pointer transition shrink-0">
                        <i className="ph-bold ph-upload-simple" />
                        {imageUploading ? 'Mengunggah...' : (formData.image_path ? 'Ganti Foto Cover' : 'Upload Foto Cover')}
                        <input
                          type="file"
                          accept="image/png,image/jpeg,image/webp"
                          className="hidden"
                          disabled={imageUploading}
                          onChange={e => {
                            const file = e.target.files?.[0]
                            if (file) handleImageUpload(file)
                            e.target.value = ''
                          }}
                        />
                      </label>
                      {formData.image_path && (
                        <button
                          type="button"
                          onClick={() => setFormData({ ...formData, image_path: '', image_source: '' })}
                          className="text-xs text-red-500 hover:text-red-700 font-medium"
                        >
                          Hapus
                        </button>
                      )}
                    </div>
                    
                    <div className="mt-4">
                        <div className="flex justify-between items-center mb-1">
                          <label className="block text-[11px] font-semibold text-slate-500">Sumber / Kredit Foto Cover</label>
                          <span className={`text-[10px] font-semibold ${(formData.image_source || '').length > 150 ? 'text-red-500' : 'text-slate-400'}`}>
                            {(formData.image_source || '').length}/150
                          </span>
                        </div>
                        <input
                          type="text"
                          value={formData.image_source || ''}
                          onChange={e => setFormData({ ...formData, image_source: e.target.value })}
                          placeholder="Contoh: Foto: Humas / Unsplash"
                          maxLength={150}
                          className="w-full px-3 py-1.5 border border-slate-300 rounded-lg text-xs outline-none focus:border-brand-500 focus:ring-1 focus:ring-brand-500 transition-colors"
                        />
                        {(formData.image_source || '').length > 150 && (
                          <p className="text-red-500 text-[10px] mt-1">Keterangan foto maksimal 150 karakter.</p>
                        )}
                      </div>
                    <p className="text-[10px] text-slate-400 mt-1.5">PNG / JPG / WEBP · Maks 5MB.</p>
                  </div>

                  <div>
                    <div className="flex justify-between items-center mb-1">
                      <label className="block text-[11px] font-semibold text-slate-500">Tags (pisahkan dengan koma) <span className="text-gray-400 font-normal">(opsional)</span></label>
                      <span className={`text-[10px] font-semibold ${(formData.tags || '').length > 200 ? 'text-red-500' : 'text-slate-400'}`}>
                        {(formData.tags || '').length}/200
                      </span>
                    </div>
                    <input type="text" maxLength={200} value={formData.tags} onChange={e => { setFormData({ ...formData, tags: e.target.value }); clearFieldError('tags') }} placeholder="teknologi, pendidikan, digital" className={`w-full px-3.5 py-2.5 border rounded-xl text-sm outline-none transition-colors ${formErrors.tags ? 'border-red-400 focus:ring-2 focus:ring-red-100' : 'border-slate-300 focus:border-brand-500 focus:ring-2 focus:ring-brand-100'}`} />
                    {formErrors.tags && (
                      <p className="text-red-500 text-[11px] font-semibold mt-1.5 flex items-center gap-1">
                        <i className="ph-bold ph-warning-circle text-xs" /> {formErrors.tags}
                      </p>
                    )}
                  </div>

                  <div className="pt-2">
                    <label className="flex items-center gap-2 text-sm cursor-pointer font-medium text-slate-700">
                      <input type="checkbox" checked={formData.is_published} onChange={e => setFormData({ ...formData, is_published: e.target.checked })} className="accent-brand-600 rounded border-slate-300 text-brand-600 focus:ring-brand-500 cursor-pointer" />
                      Terbitkan langsung (Published)
                    </label>
                  </div>
                </div>

                {/* Right Column: Descriptions & Content */}
                <div className="space-y-4">
                  <div>
                    <div className="flex justify-between items-center mb-1">
                      <label className="block text-[11px] font-semibold text-slate-500">Ringkasan (Excerpt) <span className="text-gray-400 font-normal">(opsional)</span></label>
                      <span className={`text-[10px] font-semibold ${(formData.excerpt || '').length >= 500 ? 'text-red-500' : 'text-slate-400'}`}>
                        {(formData.excerpt || '').length}/500
                      </span>
                    </div>
                    <textarea rows={3} maxLength={500} value={formData.excerpt} onChange={e => setFormData({ ...formData, excerpt: e.target.value })} placeholder="Tulis ringkasan berita singkat..." className="w-full px-3.5 py-2.5 border border-slate-300 rounded-xl text-sm outline-none focus:border-brand-500 focus:ring-2 focus:ring-brand-100 transition-colors overflow-y-auto resize-y min-h-[80px]" />
                  </div>

                  <div>
                    <div className="flex justify-between items-center mb-1">
                      <label className="block text-[11px] font-semibold text-slate-500">Sumber / Footnote <span className="text-gray-400 font-normal">(opsional)</span></label>
                      <span className={`text-[10px] font-semibold ${(formData.footnote || '').length >= 500 ? 'text-red-500' : 'text-slate-400'}`}>
                        {(formData.footnote || '').length}/500
                      </span>
                    </div>
                    <input type="text" maxLength={500} value={formData.footnote} onChange={e => setFormData({ ...formData, footnote: e.target.value })} placeholder="Contoh: Humas DPP GRADASI, DetikNews, dll." className="w-full px-3.5 py-2.5 border border-slate-300 rounded-xl text-sm outline-none focus:border-brand-500 focus:ring-2 focus:ring-brand-100 transition-colors" />
                  </div>

                  <div>
                    <div className="flex justify-between items-center mb-1">
                      <label className="block text-[11px] font-semibold text-slate-500">Konten Lengkap <span className="text-red-500">*</span></label>
                      <span className="text-[10px] font-semibold text-slate-400">
                        {(formData.content || '').length} karakter
                      </span>
                    </div>
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
                      placeholder="Isi berita lengkap di sini..."
                      className={`w-full px-3.5 py-2.5 border rounded-xl text-sm outline-none transition-colors overflow-y-auto resize-y min-h-[160px] ${touched.content && formErrors.content ? 'border-red-400 focus:ring-2 focus:ring-red-100' : 'border-slate-300 focus:border-brand-500 focus:ring-2 focus:ring-brand-100'}`}
                    />
                    {touched.content && formErrors.content && (
                      <p className="text-red-500 text-[11px] font-semibold mt-1.5 flex items-center gap-1">
                        <i className="ph-bold ph-warning-circle text-xs" /> {formErrors.content}
                      </p>
                    )}
                  </div>

                  <div className="mt-6 p-4 bg-slate-50 border border-slate-100 rounded-xl space-y-4">
                    <label className="relative inline-flex items-center cursor-pointer">
                      <input
                        type="checkbox"
                        checked={formData.is_new}
                        onChange={e => setFormData({ ...formData, is_new: e.target.checked })}
                        className="sr-only peer"
                      />
                      <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-brand-500"></div>
                      <span className="ml-3 text-sm font-medium text-gray-700">Tandai sebagai "TERBARU"</span>
                    </label>
                    <p className="text-xs text-gray-500 leading-relaxed flex gap-1.5 items-start">
                      <i className="ph-fill ph-info text-blue-500 text-sm shrink-0 mt-0.5" />
                      Tanda "TERBARU" akan otomatis hilang setelah 7 hari. Jika sudah hilang, Anda bisa memunculkannya kembali dengan mencentang kotak ini lalu klik Simpan.
                    </p>
                  </div>
                </div>

              </div>

              <div className="flex justify-end gap-2 pt-4 border-t items-center mt-6">
                {isLimited && (
                  <span className="text-xs text-amber-600 font-semibold mr-auto flex items-center gap-1">
                    <i className="ph ph-timer text-sm" /> Terlalu banyak percobaan. Tunggu {cooldown}s
                  </span>
                )}
                <button type="button" onClick={() => setIsFormOpen(false)} disabled={formLoading || isLimited} className="px-4 py-2 border rounded-xl text-sm font-semibold hover:bg-slate-50 disabled:opacity-50 transition-colors">Batal</button>
                <button type="submit" disabled={formLoading || isLimited} className="px-5 py-2 bg-brand-600 text-white rounded-xl text-sm font-semibold hover:bg-brand-700 disabled:opacity-60 disabled:cursor-not-allowed transition-colors flex items-center gap-2">
                  {formLoading && <i className="ph-bold ph-circle-notch animate-spin text-sm" />}
                  {formLoading ? 'Menyimpan...' : (isLimited ? `Tunggu ${cooldown}s` : 'Simpan')}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* DETAIL MODAL (read-only) */}
      {isDetailOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
          <div className="fixed inset-0 bg-slate-900/40 backdrop-blur-sm transition-opacity" onClick={() => setIsDetailOpen(false)} />
          <div className="relative bg-white rounded-2xl shadow-2xl max-w-2xl w-full max-h-[90vh] flex flex-col overflow-hidden z-10 animate-fade-in-up border border-slate-100">
            
            {/* Header (Sticky) */}
            <div className="border-b border-slate-100 px-6 py-4 flex items-center justify-between shrink-0 bg-white z-20">
              <h3 className="font-heading font-bold text-slate-800 text-lg flex items-center gap-2">
                <i className="ph-fill ph-newspaper text-brand-500 text-xl" />
                Detail Berita
              </h3>
              <button onClick={() => setIsDetailOpen(false)} className="w-8 h-8 flex items-center justify-center rounded-full bg-slate-100 text-slate-500 hover:bg-slate-200 hover:text-slate-800 transition-colors">
                <i className="ph-bold ph-x" />
              </button>
            </div>
            
            {/* Scrollable Content */}
            <div className="overflow-y-auto flex-1 relative bg-slate-50/30">
              {detailLoading ? (
                <div className="py-24 flex flex-col items-center justify-center text-slate-400 gap-3">
                   <i className="ph-bold ph-circle-notch animate-spin text-3xl text-brand-500" />
                   <span className="text-sm font-medium">Memuat data berita...</span>
                </div>
              ) : detailData ? (
                <div className="p-6 space-y-6">
                  
                  {/* Image Section */}
                  {(detailData.image_url || detailData.image_path) && (
                    <div className="relative rounded-xl overflow-hidden shadow-sm group bg-slate-100">
                      <img src={resolveAssetUrl(detailData.image_url || detailData.image_path)} alt={detailData.title} className="w-full aspect-[21/9] object-cover group-hover:scale-105 transition-transform duration-700" />
                      <div className="absolute inset-0 bg-gradient-to-t from-slate-900/60 via-transparent to-transparent opacity-90" />
                      {detailData.image_source && (
                        <p className="absolute bottom-3 right-4 text-[10px] text-white/90 italic font-medium drop-shadow-md">
                          <i className="ph-fill ph-camera mr-1" /> {detailData.image_source}
                        </p>
                      )}
                    </div>
                  )}

                  {/* Header Title & Meta */}
                  <div className="space-y-3">
                    <div className="flex flex-wrap gap-2">
                      <span className="px-2.5 py-1 bg-brand-50 text-brand-700 text-[10px] font-bold uppercase tracking-wider rounded-md border border-brand-100">
                        {detailData.category || 'Berita Organisasi'}
                      </span>
                      {detailData.is_published ? (
                        <span className="px-2.5 py-1 bg-emerald-50 text-emerald-700 text-[10px] font-bold uppercase tracking-wider rounded-md border border-emerald-100">
                          Published
                        </span>
                      ) : (
                        <span className="px-2.5 py-1 bg-amber-50 text-amber-700 text-[10px] font-bold uppercase tracking-wider rounded-md border border-amber-100">
                          Draft
                        </span>
                      )}
                      {detailData.is_new && (
                        <span className="px-2.5 py-1 bg-blue-50 text-blue-700 text-[10px] font-bold uppercase tracking-wider rounded-md border border-blue-100">
                          Terbaru
                        </span>
                      )}
                    </div>
                    
                    <h1 className="font-heading font-black text-slate-900 text-2xl md:text-3xl leading-snug break-all">
                      {detailData.title}
                    </h1>

                    <div className="flex flex-wrap items-center gap-4 text-xs font-medium text-slate-500 pt-2">
                      <div className="flex items-center gap-1.5 min-w-0" title="Tanggal Terbit">
                        <i className="ph-fill ph-calendar-blank text-slate-400 text-sm shrink-0" />
                        <span className="break-all">{detailData.published_date || '-'}</span>
                      </div>
                      <div className="flex items-center gap-1.5 min-w-0" title="Penulis">
                        <i className="ph-fill ph-user text-slate-400 text-sm shrink-0" />
                        <span className="break-all">{detailData.author_name || 'Admin'}</span>
                      </div>
                      <div className="flex items-center gap-1.5" title="Dilihat">
                        <i className="ph-fill ph-eye text-slate-400 text-sm" />
                        {detailData.views ?? 0} kali
                      </div>
                    </div>
                  </div>

                  {/* Ringkasan */}
                  {detailData.excerpt && (
                    <div className="p-4 bg-white border border-slate-200 rounded-xl shadow-sm relative overflow-hidden">
                      <div className="absolute top-0 left-0 w-1 h-full bg-brand-500" />
                      <h4 className="text-[10px] font-bold text-slate-400 uppercase tracking-widest mb-1.5">Ringkasan</h4>
                      <p className="text-slate-700 text-sm leading-relaxed italic break-all">"{detailData.excerpt}"</p>
                    </div>
                  )}

                  {/* Main Content */}
                  <div className="bg-white p-5 sm:p-6 rounded-xl border border-slate-200 shadow-sm relative overflow-hidden">
                    <h4 className="text-[10px] font-bold text-slate-400 uppercase tracking-widest mb-3 border-b border-slate-100 pb-2">Konten Lengkap</h4>
                    <div className="text-slate-700 text-sm leading-loose whitespace-pre-line break-all text-justify">
                      {detailData.content || '-'}
                    </div>
                  </div>

                  {/* Tags & Footnote */}
                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 bg-white p-5 rounded-xl border border-slate-200 shadow-sm">
                    {detailData.footnote && (
                      <div className="col-span-1 sm:col-span-2 min-w-0">
                        <span className="block text-[10px] font-bold text-slate-400 uppercase tracking-wider mb-1">Sumber / Footnote</span>
                        <p className="text-slate-700 text-sm italic break-all">{detailData.footnote}</p>
                      </div>
                    )}
                    {Array.isArray(detailData.tags) && detailData.tags.length > 0 && (
                      <div className="col-span-1 sm:col-span-2 mt-2">
                        <span className="block text-[10px] font-bold text-slate-400 uppercase tracking-wider mb-2">Tags</span>
                        <div className="flex flex-wrap gap-1.5">
                          {detailData.tags.map(tag => (
                            <span key={tag} className="bg-slate-100 text-slate-600 text-[11px] font-semibold px-2.5 py-1 rounded-md border border-slate-200 break-all">{tag}</span>
                          ))}
                        </div>
                      </div>
                    )}
                  </div>

                  {/* Safe padding bottom for absolute certainty */}
                  <div className="h-2"></div>
                </div>
              ) : (
                <div className="py-16 text-center text-slate-500">Data tidak ditemukan</div>
              )}
            </div>

            {/* Sticky Footer */}
            <div className="border-t border-slate-200 px-6 py-4 bg-slate-50 shrink-0 flex justify-end">
              <button onClick={() => setIsDetailOpen(false)} className="px-5 py-2.5 bg-white border border-slate-300 rounded-xl text-sm font-bold text-slate-700 hover:bg-slate-100 hover:text-slate-900 transition-colors shadow-sm focus:ring-2 focus:ring-slate-200 outline-none">
                Tutup Detail
              </button>
            </div>

          </div>
        </div>
      )}
    </AdminLayout>
  )
}
