import { useState, useEffect, useCallback } from 'react'
import AdminLayout from '../../layouts/AdminLayout'
import { pengurusService } from '../../services/pengurusService'
import ConfirmDialog from '../../components/admin/ConfirmDialog'
import ToastNotification from '../../components/admin/ToastNotification'
import { resolveAssetUrl } from '../../utils/assetUrl'
import { useFormErrors, useRateLimitCooldown } from '../../utils/parseApiError'

const PAGE_SIZE = 30

export default function PengurusAdmin() {
  const [items, setItems] = useState([]) // untuk tampilanNg paginated
  const [allItems, setAllItems] = useState([]) // semua item yang difilter (tanpa pagination)
  const [loading, setLoading] = useState(false)
  const [meta, setMeta] = useState({ total_data: 0, total_pages: 1, current_page: 1, limit: PAGE_SIZE })

  const [currentTab, setCurrentTab] = useState('active') // active | inactive | trash
  const [currentPage, setCurrentPage] = useState(1)
  const [searchQuery, setSearchQuery] = useState('')
  const [filterLevel, setFilterLevel] = useState('')
  const [selectedItems, setSelectedItems] = useState([])

  const [isFormOpen, setIsFormOpen] = useState(false)
  const [formMode, setFormMode] = useState('create')
  const [formData, setFormData] = useState({
    id: null,
    name: '',
    level: 'Pengurus Pusat',
    role: '',
    department: '',
    periode: '2025 - 2030',
    provinsi: '',
    kabupaten: '',
    facebook_url: '',
    instagram_url: '',
    linkedin_url: '',
    whatsapp: '',
    email: '',
    pekerjaan: '',
    bio: '',
    pendidikan: '',
    sertifikasi: '',
    sort_order: 1,
    image_path: '',
    image: null,
    cv_path: '',
    cv: null,
    is_active: true,
  })

  const [formErrors, setFormErrors] = useState({})
  const [touched, setTouched] = useState({})

  const validateForm = useCallback((data = formData) => {
    const errors = {}
    if (!data.name || !data.name.trim()) {
      errors.name = 'Nama lengkap wajib diisi.'
    }
    if (!data.level) {
      errors.level = 'Tingkat struktur wajib dipilih.'
    }
    if (!data.role || !data.role.trim()) {
      errors.role = 'Jabatan resmi wajib diisi.'
    }
    if (!data.periode || !data.periode.trim()) {
      errors.periode = 'Periode wajib diisi.'
    }
    if ((data.level === 'Pengurus Provinsi' || data.level === 'Pengurus Kab/Kota') && (!data.provinsi || !data.provinsi.trim())) {
      errors.provinsi = 'Provinsi harus diisi untuk pengurus daerah/cabang'
    }
    if (data.level === 'Pengurus Kab/Kota' && (!data.kabupaten || !data.kabupaten.trim())) {
      errors.kabupaten = 'Kabupaten/Kota wajib diisi.'
    }
    if (formMode === 'create' && !data.image) {
      errors.image = 'Foto profil wajib diunggah.'
    }
    return errors
  }, [formData, formMode])

  const [confirm, setConfirm] = useState({ isOpen: false, type: 'danger', title: '', message: '', action: null })
  const [toast, setToast] = useState({ show: false, message: '', type: 'success' })
  // Error dari backend: field errors inline + countdown rate limit
  const { applyError, clearFieldError, resetFieldErrors } = useFormErrors()
  const { cooldown, isLimited, applyRateLimit } = useRateLimitCooldown()

  const isProvinceRequired = formData.level === 'Pengurus Provinsi' || formData.level === 'Pengurus Kab/Kota'
  const isKabupatenRequired = formData.level === 'Pengurus Kab/Kota'

  const showToast = useCallback((message, type = 'success') => {
    setToast({ show: true, message, type })
    setTimeout(() => setToast(prev => ({ ...prev, show: false })), 3000)
  }, [])

  // Load SEMUA data (tanpa pagination) untuk keperluan reorder global
  const loadAllPengurus = useCallback(() => {
    // Kita perlu load semua data tanpa pagination
    // Contact backend untuk menyediakan endpoint listAdmin tanpa page/limit? Atau kita hit dengan page=1&limit=9999?
    // Untuk saat ini, kita akan hit dengan limit besar. Alternatif: endpoint khusus /admin/pengurus/all
    // Tapi kontrak tidak ada. Jadi kita pake listAdmin dengan limit besar.
    // Implementation: call listAdmin with a very high limit
    const params = {
      page: 1,
      limit: 9999, // harap tidak Page Not Found
      search: searchQuery,
      status: currentTab === 'trash' ? 'all' : (currentTab === 'inactive' ? 'inactive' : 'active'),
      trashed: currentTab === 'trash',
      level: filterLevel || undefined,
      sort: 'sort_order',
    }
    return pengurusService.listAdmin(params)
  }, [currentTab, searchQuery, filterLevel])

  const loadPengurus = useCallback(() => {
    setLoading(true)
    loadAllPengurus()
      .then(res => {
        if (res && res.data) {
          let list = []
          if (Array.isArray(res.data)) {
            list = res.data
          } else if (res.data.data) {
            list = res.data.data
          }
          setAllItems(list) // simpan semua data
          // Untuk items (ditampilkan), kita akan ambil slice berdasarkan currentPage dari filteredAllItems nanti
          updateMetaFromResponse(res.data)
        } else {
          setAllItems([])
          setItems([])
        }
      })
      .catch(() => {
        setAllItems([])
        setItems([])
      })
      .finally(() => setLoading(false))
  }, [loadAllPengurus])

  useEffect(() => {
    loadPengurus()
  }, [loadPengurus])

  useEffect(() => {
    // Ketika filter/search/tab berubah, reset page & pilihan
    setCurrentPage(1)
    setSelectedItems([])
  }, [currentTab, searchQuery, filterLevel])

  // Update items (page view) berdasarkan allItems yang difilter + pagination
  useEffect(() => {
    // allItems sudah terfilter oleh response (karena kita kirim parameter filter)
    // Tapi allItems berisi semua data, kita paginate di frontend
    const total = allItems.length
    const totalPages = Math.ceil(total / PAGE_SIZE) || 1
    const start = (currentPage - 1) * PAGE_SIZE
    const end = start + PAGE_SIZE
    const paginatedList = allItems.slice(start, end)
    setItems(paginatedList)
    setMeta({
      total_data: total,
      total_pages: totalPages,
      current_page: currentPage,
      limit: PAGE_SIZE,
    })
  }, [allItems, currentPage])

  function updateMetaFromResponse(data) {
    if (data && data.meta) {
      setMeta(data.meta)
    }
  }

  const openForm = (item = null) => {
    setFormErrors({})
    setTouched({})
    resetFieldErrors()
    if (item) {
      setFormMode('edit')
      setFormData({
        id: item.id,
        name: item.name,
        level: item.level,
        role: item.role,
        department: item.department || '',
        periode: item.periode || '2025 - 2030',
        provinsi: item.provinsi || '',
        kabupaten: item.kabupaten || '',
        facebook_url: item.facebook_url || '',
        instagram_url: item.instagram_url || '',
        linkedin_url: item.linkedin_url || '',
        whatsapp: item.whatsapp || '',
        email: item.email || '',
        pekerjaan: item.pekerjaan || '',
        bio: item.bio || '',
        pendidikan: item.pendidikan || '',
        sertifikasi: item.sertifikasi || '',
        sort_order: item.sort_order || 1,
        image_path: item.image_path || item.image_url || '',
        image: null,
        cv_path: item.cv_path || '',
        cv: null,
        is_active: item.is_active !== false,
      })
    } else {
      setFormMode('create')
      setFormData({
        id: null,
        name: '',
        level: 'Pengurus Pusat',
        role: '',
        department: '',
        periode: '2025 - 2030',
        provinsi: '',
        kabupaten: '',
        facebook_url: '',
        instagram_url: '',
        linkedin_url: '',
        whatsapp: '',
        email: '',
        pekerjaan: '',
        bio: '',
        pendidikan: '',
        sertifikasi: '',
        sort_order: (allItems.length || 0) + 1,
        image_path: '',
        image: null,
        cv_path: '',
        cv: null,
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

    const payload = { ...formData }
    // Backend butuh multipart; buang id (meta) — image_path dikirim sebagai URL string
    if (formMode === 'create') {
      try {
        await pengurusService.create(payload)
        showToast('Pengurus berhasil ditambahkan.')
        setIsFormOpen(false)
        loadPengurus()
      } catch (err) {
        const parsed = applyError(err)
        applyRateLimit(err)
        setFormErrors(prev => ({ ...prev, ...parsed.fieldErrors }))
        setTouched(prev => ({ ...prev, ...Object.keys(parsed.fieldErrors).reduce((acc, k) => ({ ...acc, [k]: true }), {}) }))
        if (Object.keys(parsed.fieldErrors).length === 0) {
          showToast(parsed.message || 'Gagal menyimpan pengurus', 'error')
        }
      }
    } else {
      try {
        await pengurusService.update(formData.id, payload)
        showToast('Pengurus berhasil diperbarui.')
        setIsFormOpen(false)
        loadPengurus()
      } catch (err) {
        const parsed = applyError(err)
        applyRateLimit(err)
        setFormErrors(prev => ({ ...prev, ...parsed.fieldErrors }))
        setTouched(prev => ({ ...prev, ...Object.keys(parsed.fieldErrors).reduce((acc, k) => ({ ...acc, [k]: true }), {}) }))
        if (Object.keys(parsed.fieldErrors).length === 0) {
          showToast(parsed.message || 'Gagal menyimpan pengurus', 'error')
        }
      }
    }
  }

  const handleMove = async (globalIndex, direction) => {
    // globalIndex adalah index dalam allItems (setelah filtered)
    const targetGlobalIndex = direction === 'up' ? globalIndex - 1 : globalIndex + 1
    if (targetGlobalIndex < 0 || targetGlobalIndex >= allItems.length) return

    const currentItem = allItems[globalIndex]
    const targetItem = allItems[targetGlobalIndex]

    // Swap IDs in the array
    const newAllItems = [...allItems]
    newAllItems[globalIndex] = targetItem
    newAllItems[targetGlobalIndex] = currentItem

    // Update local state immediately for responsive UI
    setAllItems(newAllItems)

    // Send reorder request with all IDs in new order
    const reorderedIds = newAllItems.map(item => item.id)
    try {
      await pengurusService.reorder(reorderedIds)
      showToast('Urutan pengurus diperbarui.')
      // Refresh data from server to ensure consistency
      loadPengurus()
    } catch (err) {
      // Rollback on error
      setAllItems(allItems)
      showToast(err.message || 'Gagal mengubah urutan', 'error')
    }
  }

  const confirmAction = (type, item = null) => {
    const name = item ? item.name : ''
    const configs = {
      delete: {
        type: 'danger',
        title: 'Hapus Pengurus',
        message: `Anda akan memindahkan "${name}" ke Sampah. Lanjutkan?`,
        action: async () => {
          await pengurusService.remove(item.id)
          showToast('Pengurus berhasil dihapus.')
          loadPengurus()
        },
      },
      restore: {
        type: 'info',
        title: 'Pulihkan Pengurus',
        message: `Anda akan memulihkan "${name}" dari Sampah. Lanjutkan?`,
        action: async () => {
          await pengurusService.restore(item.id)
          showToast('Pengurus berhasil dipulihkan.')
          loadPengurus()
        },
      },
      toggle_status: {
        type: item?.is_active ? 'warning' : 'info',
        title: item?.is_active ? 'Non-aktifkan Pengurus' : 'Aktifkan Pengurus',
        message: item?.is_active
          ? `Anda akan menonaktifkan "${name}". Pengurus tidak akan tampil di halaman publik. Lanjutkan?`
          : `Anda akan mengaktifkan "${name}". Pengurus akan tampil di halaman publik. Lanjutkan?`,
        action: async () => {
          try {
            await pengurusService.update(item.id, { ...item, is_active: !item.is_active, image: null, cv: null })
            // Update local state immediately
            setAllItems(prev => prev.map(i => i.id === item.id ? { ...i, is_active: !item.is_active } : i))
            showToast(item.is_active ? 'Pengurus dinonaktifkan.' : 'Pengurus berhasil diaktifkan!')
          } catch (err) {
            showToast(err?.message || 'Gagal mengubah status pengurus.', 'error')
          }
        },
      },
      bulk_delete: {
        type: 'danger',
        title: 'Hapus Massal',
        message: `Anda akan memindahkan ${selectedItems.length} item ke Sampah. Lanjutkan?`,
        action: async () => {
          await pengurusService.bulkDelete(selectedItems)
          setSelectedItems([])
          showToast('Pengurus berhasil dihapus massal.')
          loadPengurus()
        },
      },
      bulk_restore: {
        type: 'info',
        title: 'Pulihkan Massal',
        message: `Anda akan memulihkan ${selectedItems.length} item dari Sampah. Lanjutkan?`,
        action: async () => {
          await pengurusService.bulkRestore(selectedItems)
          setSelectedItems([])
          showToast('Pengurus berhasil dipulihkan massal.')
          loadPengurus()
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

  // selectedItems selection based on current page view
  const isAllSelected = items.length > 0 && selectedItems.length === items.length
  function toggleAll() {
    setSelectedItems(isAllSelected ? [] : items.map(i => i.id))
  }
  function toggleOne(id) {
    setSelectedItems(prev => prev.includes(id) ? prev.filter(x => x !== id) : [...prev, id])
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

  const totalPages = Math.max(1, meta.total_pages || 1)
  const totalData = meta.total_data ?? allItems.length



  function resetFilter() {
    setSearchQuery('')
    setFilterLevel('')
    setCurrentPage(1)
    showToast('Filter direset.', 'success')
  }

  const headerContent = (
    <div className="flex items-center gap-2 w-full max-w-2xl animate-fade-in-up">
      <div className="relative w-full">
        <i className="ph ph-magnifying-glass absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
        <input
          type="text"
          value={searchQuery}
          onChange={e => setSearchQuery(e.target.value)}
          placeholder="Cari pengurus..."
          className="w-full pl-9 pr-4 py-2 bg-gray-50 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-brand-500 focus:bg-white transition-colors"
        />
      </div>
      <select
        value={filterLevel}
        onChange={e => setFilterLevel(e.target.value)}
        className="shrink-0 bg-gray-50 border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-brand-500 focus:bg-white transition-colors cursor-pointer"
      >
        <option value="">Semua Level</option>
        <option value="Ketua Umum">Ketua Umum</option>
        <option value="Pengurus Pusat">Pengurus Pusat</option>
        <option value="Pengurus Provinsi">Pengurus Provinsi</option>
        <option value="Pengurus Kab/Kota">Pengurus Kab/Kota</option>
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
    <AdminLayout title="Kelola Pengurus" headerContent={headerContent}>
      <ConfirmDialog {...confirm} onClose={() => setConfirm(prev => ({ ...prev, isOpen: false, action: null }))} onConfirm={executeConfirm} />
      <ToastNotification show={toast.show} message={toast.message} type={toast.type} />

      <div className="space-y-6 animate-fade-in-up">
        {/* Navigation Tabs */}
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
          <div className="flex bg-white rounded-lg p-1 border border-gray-200 shadow-sm">
            <button
              onClick={() => setCurrentTab('active')}
              className={`px-4 py-1.5 rounded-md text-sm transition-all duration-200 ${currentTab === 'active' ? 'bg-brand-50 text-brand-600 shadow-sm font-medium' : 'text-gray-500 hover:text-gray-700'}`}
            >
              Pengurus Aktif
            </button>
            <button
              onClick={() => setCurrentTab('inactive')}
              className={`px-4 py-1.5 rounded-md text-sm transition-all duration-200 ${currentTab === 'inactive' ? 'bg-brand-50 text-brand-600 shadow-sm font-medium' : 'text-gray-500 hover:text-gray-700'}`}
            >
              Non-aktif
            </button>
            <button
              onClick={() => setCurrentTab('trash')}
              className={`px-4 py-1.5 rounded-md text-sm transition-all duration-200 flex items-center gap-2 ${currentTab === 'trash' ? 'bg-red-50 text-red-600 shadow-sm font-medium' : 'text-gray-500 hover:text-gray-700'}`}
            >
              <i className="ph ph-trash" /> Sampah (History)
            </button>
          </div>
          {currentTab !== 'trash' && (
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
            <span className="text-sm font-semibold text-brand-700">{selectedItems.length} dipilih</span>
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

        <div className="bg-white rounded-2xl border border-gray-200 overflow-hidden shadow-sm">
          {loading ? (
            <div className="py-16 text-center text-slate-500">Memuat pengurus...</div>
          ) : items.length === 0 ? (
            <div className="py-16 text-center text-slate-500">
              <i className="ph-bold ph-users-three text-4xl text-slate-300 mb-2 block" />
              {currentTab === 'trash' ? 'Sampah kosong' : 'Tidak ada pengurus ditemukan'}
            </div>
          ) : (
            <>
              <div className="overflow-x-auto">
                <table className="w-full text-left text-sm text-slate-700">
                  <thead className="bg-slate-50 border-b border-gray-200 font-semibold text-xs uppercase tracking-wider text-slate-500">
                    <tr>
                      <th className="p-4 w-10">
                        <input type="checkbox" checked={isAllSelected} onChange={toggleAll} className="accent-brand-600" />
                      </th>
                      <th className="p-4">Nama Lengkap</th>
                      <th className="p-4">Tingkat</th>
                      <th className="p-4">Jabatan</th>
                      <th className="p-4">Wilayah</th>
                      <th className="p-4">Status</th>
                      <th className="p-4 w-24">Urutan</th>
                      <th className="p-4 text-right">Aksi</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-100">
                    {items.map((item) => {
                      // global index dalam allItems untuk determined posisi global
                      const globalIndex = allItems.findIndex(i => i.id === item.id)
                      return (
                        <tr key={item.id} className="hover:bg-slate-50/50 admin-row">
                          <td className="p-4">
                            <input type="checkbox" checked={selectedItems.includes(item.id)} onChange={() => toggleOne(item.id)} className="accent-brand-600 border-gray-300 rounded" />
                          </td>
                          <td className="p-4 font-medium text-slate-900">
                            <div className="flex items-center gap-3">
                              {item.image_url ? (
                                <img src={resolveAssetUrl(item.image_url)} alt="" className="w-10 h-10 rounded-full object-cover border border-gray-200 previewable-image" />
                              ) : (
                                <div className="w-10 h-10 rounded-full bg-brand-100 text-brand-600 flex items-center justify-center font-bold text-sm">
                                  {item.name.charAt(0).toUpperCase()}
                                </div>
                              )}
                              <div>
                                <div>{item.name}</div>
                                <div className="text-xs text-slate-400">{item.periode}</div>
                              </div>
                            </div>
                          </td>
                          <td className="p-4 uppercase font-bold text-xs text-brand-600">{item.level}</td>
                          <td className="p-4">{item.role}</td>
                          <td className="p-4 text-slate-500 text-xs">
                            {item.provinsi ? `${item.provinsi}${item.kabupaten ? `, ${item.kabupaten}` : ''}` : '-'}
                          </td>
                          <td className="p-4">
                            {currentTab === 'trash' ? (
                              <span className="inline-flex items-center gap-1.5 text-xs font-semibold text-red-500">
                                <i className="ph ph-trash" /> Di Sampah
                              </span>
                            ) : (
                              <button onClick={() => confirmAction('toggle_status', item)} className="inline-flex items-center gap-2 cursor-pointer">
                                <div className={`relative w-9 h-5 rounded-full transition-colors ${item.is_active ? 'bg-brand-600' : 'bg-slate-200'}`}>
                                  <div className={`absolute top-[2px] left-[2px] w-4 h-4 bg-white rounded-full transition-transform ${item.is_active ? 'translate-x-4' : 'translate-x-0'}`} />
                                </div>
                                <span className={`text-xs font-semibold ${item.is_active ? 'text-brand-600' : 'text-slate-400'}`}>
                                  {item.is_active ? 'Aktif' : 'Non-aktif'}
                                </span>
                              </button>
                            )}
                          </td>
                          <td className="p-4">
                            <div className="flex items-center gap-3">
                              <button
                                type="button"
                                disabled={globalIndex === 0}
                                onClick={() => handleMove(globalIndex, 'up')}
                                className="p-1 hover:bg-slate-100 text-slate-500 hover:text-slate-700 rounded disabled:opacity-20 transition-all"
                                title="Pindahkan ke atas"
                              >
                                <i className="ph-bold ph-arrow-up text-base" />
                              </button>
                              <button
                                type="button"
                                disabled={globalIndex === allItems.length - 1}
                                onClick={() => handleMove(globalIndex, 'down')}
                                className="p-1 hover:bg-slate-100 text-slate-500 hover:text-slate-700 rounded disabled:opacity-20 transition-all"
                                title="Pindahkan ke bawah"
                              >
                                <i className="ph-bold ph-arrow-down text-base" />
                              </button>
                            </div>
                          </td>
                          <td className="p-4 text-right">
                            <div className="flex justify-end gap-2">
                              {currentTab === 'trash' ? (
                                <button onClick={() => confirmAction('restore', item)} title="Pulihkan" className="p-1.5 text-gray-500 hover:text-emerald-600 hover:bg-emerald-50 rounded">
                                  <i className="ph ph-arrow-counter-clockwise text-lg" />
                                </button>
                              ) : (
                                <>
                                  <button onClick={() => openForm(item)} className="p-1.5 text-gray-500 hover:text-brand-600 hover:bg-brand-50 rounded" title="Edit">
                                    <i className="ph ph-pencil-simple text-lg" />
                                  </button>
                                  <button onClick={() => confirmAction('delete', item)} className="p-1.5 text-gray-500 hover:text-red-600 hover:bg-red-50 rounded" title="Hapus (Soft Delete)">
                                    <i className="ph ph-trash text-lg" />
                                  </button>
                                </>
                              )}
                            </div>
                          </td>
                        </tr>
                      )
                    })}
                  </tbody>
                </table>
              </div>

              {/* Pagination */}
              <div className="flex items-center justify-between px-4 py-3 border-t border-slate-200">
                <span className="text-xs text-slate-500">Hal {currentPage} dari {totalPages} · {totalData} data</span>
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
            </>
          )}
        </div>
      </div>

        {/* Form Modal */}
        {isFormOpen && (
          <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div className="fixed inset-0 bg-black/40 backdrop-blur-[2px]" onClick={() => setIsFormOpen(false)} />
            <div className="relative bg-white rounded-2xl shadow-2xl max-w-4xl w-full max-h-[90vh] flex flex-col overflow-hidden z-10">
              <div className="border-b border-slate-200 px-6 py-4 flex items-center justify-between">
                <h3 className="font-heading font-bold text-slate-900 text-lg">
                  {formMode === 'create' ? 'Tambah Pengurus Baru' : 'Edit Pengurus'}
                </h3>
                <button onClick={() => setIsFormOpen(false)} className="text-slate-400 hover:text-slate-600 p-1">
                  <i className="ph-bold ph-x text-lg" />
                </button>
              </div>
              <form onSubmit={handleSubmit} noValidate className="p-6 overflow-y-auto max-h-[calc(90vh-120px)]">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6 animate-fade-in-up">
                  {/* Left Column (Main Data) */}
                  <div className="space-y-5">
                    <div>
                      <label className="block text-xs font-semibold text-gray-500 mb-1">Nama Lengkap <span className="text-red-500">*</span></label>
                      <input
                        type="text"
                        value={formData.name}
                        onChange={e => {
                          setFormData({ ...formData, name: e.target.value })
                          clearFieldError('name')
                          if (touched.name) {
                            const errs = validateForm({ ...formData, name: e.target.value })
                            setFormErrors(prev => ({ ...prev, name: errs.name }))
                          }
                        }}
                        onBlur={() => {
                          setTouched(prev => ({ ...prev, name: true }))
                          const errs = validateForm()
                          setFormErrors(prev => ({ ...prev, name: errs.name }))
                        }}
                        className={`w-full px-3.5 py-2.5 border rounded-xl text-sm outline-none bg-white transition-colors ${touched.name && formErrors.name ? 'border-red-400 focus:ring-2 focus:ring-red-100' : 'border-slate-300 focus:border-brand-500 focus:ring-2 focus:ring-brand-100'}`}
                      />
                      {touched.name && formErrors.name && (
                        <p className="text-red-500 text-[11px] font-semibold mt-1 flex items-center gap-1">
                          <i className="ph-bold ph-warning-circle text-xs" /> {formErrors.name}
                        </p>
                      )}
                    </div>

                    <div className="grid grid-cols-2 gap-4">
                      <div>
                        <label className="block text-xs font-semibold text-gray-500 mb-1">Tingkat Struktur <span className="text-red-500">*</span></label>
                        <select
                          value={formData.level}
                          onChange={e => {
                            const val = e.target.value
                            setFormData({ ...formData, level: val })
                            clearFieldError('level')
                            if (touched.level) {
                              const errs = validateForm({ ...formData, level: val })
                              setFormErrors(prev => ({ ...prev, level: errs.level }))
                            }
                          }}
                          onBlur={() => {
                            setTouched(prev => ({ ...prev, level: true }))
                            const errs = validateForm()
                            setFormErrors(prev => ({ ...prev, level: errs.level }))
                          }}
                          className="w-full px-3.5 py-2.5 border border-slate-300 rounded-xl text-sm outline-none bg-white focus:border-brand-500 focus:ring-2 focus:ring-brand-100 transition-colors cursor-pointer"
                        >
                          <option value="Ketua Umum">Ketua Umum</option>
                          <option value="Pengurus Pusat">Pengurus Pusat</option>
                          <option value="Pengurus Provinsi">Pengurus Provinsi</option>
                          <option value="Pengurus Kab/Kota">Pengurus Kab/Kota</option>
                        </select>
                      </div>
                      <div>
                        <label className="block text-xs font-semibold text-gray-500 mb-1">Jabatan Resmi <span className="text-red-500">*</span></label>
                        <input
                          type="text"
                          value={formData.role}
                          onChange={e => {
                            setFormData({ ...formData, role: e.target.value })
                            clearFieldError('role')
                            if (touched.role) {
                              const errs = validateForm({ ...formData, role: e.target.value })
                              setFormErrors(prev => ({ ...prev, role: errs.role }))
                            }
                          }}
                          onBlur={() => {
                            setTouched(prev => ({ ...prev, role: true }))
                            const errs = validateForm()
                            setFormErrors(prev => ({ ...prev, role: errs.role }))
                          }}
                          placeholder="Misal: Ketua Bidang Organisasi"
                          className={`w-full px-3.5 py-2.5 border rounded-xl text-sm outline-none bg-white transition-colors ${touched.role && formErrors.role ? 'border-red-400 focus:ring-2 focus:ring-red-100' : 'border-slate-300 focus:border-brand-500 focus:ring-2 focus:ring-brand-100'}`}
                        />
                        {touched.role && formErrors.role && (
                          <p className="text-red-500 text-[11px] font-semibold mt-1 flex items-center gap-1">
                            <i className="ph-bold ph-warning-circle text-xs" /> {formErrors.role}
                          </p>
                        )}
                      </div>
                    </div>

                    <div className="grid grid-cols-2 gap-4">
                      <div>
                        <label className="block text-xs font-semibold text-gray-500 mb-1">Periode <span className="text-red-500">*</span></label>
                        <input
                          type="text"
                          value={formData.periode}
                          onChange={e => {
                            setFormData({ ...formData, periode: e.target.value })
                            clearFieldError('periode')
                            if (touched.periode) {
                              const errs = validateForm({ ...formData, periode: e.target.value })
                              setFormErrors(prev => ({ ...prev, periode: errs.periode }))
                            }
                          }}
                          onBlur={() => {
                            setTouched(prev => ({ ...prev, periode: true }))
                            const errs = validateForm()
                            setFormErrors(prev => ({ ...prev, periode: errs.periode }))
                          }}
                          placeholder="Contoh: 2025 - 2030"
                          className={`w-full px-3.5 py-2.5 border rounded-xl text-sm outline-none bg-white transition-colors ${touched.periode && formErrors.periode ? 'border-red-400 focus:ring-2 focus:ring-red-100' : 'border-slate-300 focus:border-brand-500 focus:ring-2 focus:ring-brand-100'}`}
                        />
                        {touched.periode && formErrors.periode && (
                          <p className="text-red-500 text-[11px] font-semibold mt-1 flex items-center gap-1">
                            <i className="ph-bold ph-warning-circle text-xs" /> {formErrors.periode}
                          </p>
                        )}
                      </div>
                      <div>
                        <label className="block text-xs font-semibold text-gray-500 mb-1">Departemen <span className="text-gray-400 font-normal">(opsional)</span></label>
                        <input type="text" value={formData.department} onChange={e => setFormData({ ...formData, department: e.target.value })} placeholder="Misal: Departemen IT" className="w-full px-3.5 py-2.5 border border-slate-300 rounded-xl text-sm outline-none bg-white focus:border-brand-500 focus:ring-2 focus:ring-brand-100 transition-colors" />
                      </div>
                    </div>

                    <div className="grid grid-cols-2 gap-4">
                      <div>
                        <label className="block text-xs font-semibold text-gray-500 mb-1">Provinsi <span className={isProvinceRequired ? 'text-red-500' : 'text-gray-400 font-normal'}>{isProvinceRequired ? '*' : '(opsional)'}</span></label>
                        <input
                          type="text"
                          value={formData.provinsi}
                          onChange={e => {
                            const val = e.target.value
                            setFormData({ ...formData, provinsi: val })
                            clearFieldError('provinsi')
                            if (touched.provinsi) {
                              const errs = validateForm({ ...formData, provinsi: val })
                              setFormErrors(prev => ({ ...prev, provinsi: errs.provinsi }))
                            }
                          }}
                          onBlur={() => {
                            setTouched(prev => ({ ...prev, provinsi: true }))
                            const errs = validateForm()
                            setFormErrors(prev => ({ ...prev, provinsi: errs.provinsi }))
                          }}
                          placeholder="Wajib jika Provinsi/Kab/Kota"
                          className={`w-full px-3.5 py-2.5 border rounded-xl text-sm outline-none bg-white transition-colors ${touched.provinsi && formErrors.provinsi ? 'border-red-400 focus:ring-2 focus:ring-red-100' : 'border-slate-300 focus:border-brand-500 focus:ring-2 focus:ring-brand-100'}`}
                        />
                        {touched.provinsi && formErrors.provinsi && (
                          <p className="text-red-500 text-[11px] font-semibold mt-1 flex items-center gap-1">
                            <i className="ph-bold ph-warning-circle text-xs" /> {formErrors.provinsi}
                          </p>
                        )}
                      </div>
                      <div>
                        <label className="block text-xs font-semibold text-gray-500 mb-1">Kabupaten/Kota <span className={isKabupatenRequired ? 'text-red-500' : 'text-gray-400 font-normal'}>{isKabupatenRequired ? '*' : '(opsional)'}</span></label>
                        <input
                          type="text"
                          value={formData.kabupaten}
                          onChange={e => {
                            const val = e.target.value
                            setFormData({ ...formData, kabupaten: val })
                            clearFieldError('kabupaten')
                            if (touched.kabupaten) {
                              const errs = validateForm({ ...formData, kabupaten: val })
                              setFormErrors(prev => ({ ...prev, kabupaten: errs.kabupaten }))
                            }
                          }}
                          onBlur={() => {
                            setTouched(prev => ({ ...prev, kabupaten: true }))
                            const errs = validateForm()
                            setFormErrors(prev => ({ ...prev, kabupaten: errs.kabupaten }))
                          }}
                          placeholder="Wajib jika Kab/Kota"
                          className={`w-full px-3.5 py-2.5 border rounded-xl text-sm outline-none bg-white transition-colors ${touched.kabupaten && formErrors.kabupaten ? 'border-red-400 focus:ring-2 focus:ring-red-100' : 'border-slate-300 focus:border-brand-500 focus:ring-2 focus:ring-brand-100'}`}
                        />
                        {touched.kabupaten && formErrors.kabupaten && (
                          <p className="text-red-500 text-[11px] font-semibold mt-1 flex items-center gap-1">
                            <i className="ph-bold ph-warning-circle text-xs" /> {formErrors.kabupaten}
                          </p>
                        )}
                      </div>
                    </div>

                    <div>
                      <label className="block text-xs font-semibold text-gray-500 mb-1">Foto Profil <span className={formMode === 'create' ? 'text-red-500' : 'text-gray-400 font-normal'}>{formMode === 'create' ? '*' : '(opsional)'}</span></label>
                      <div className="flex items-center gap-3">
                        {(formData.image_path || formData.image) && (
                          <img
                            src={formData.image instanceof File ? URL.createObjectURL(formData.image) : resolveAssetUrl(formData.image_path)}
                            alt="Profil Preview"
                            className="w-16 h-16 rounded-full object-cover border border-slate-200 shrink-0"
                          />
                        )}
                        <label className="inline-flex items-center gap-2 px-4 py-2 bg-brand-600 text-white rounded-lg hover:bg-brand-700 text-sm font-semibold cursor-pointer transition shrink-0 shadow-sm">
                          <i className="ph-bold ph-upload-simple" />
                          {formData.image || formData.image_path ? 'Ganti Foto' : 'Upload Foto'}
                          <input
                            type="file"
                            accept="image/png,image/jpeg,image/webp"
                            className="hidden"
                            onChange={e => {
                              const file = e.target.files?.[0]
                              if (file) setFormData(prev => ({ ...prev, image: file, image_path: '' }))
                            }}
                          />
                        </label>
                        {(formData.image || formData.image_path) && (
                          <button
                            type="button"
                            onClick={() => setFormData({ ...formData, image: null, image_path: '' })}
                            className="text-xs text-red-500 hover:text-red-700 font-medium"
                          >
                            Hapus
                          </button>
                        )}
                      </div>
                      {formMode === 'create' && touched.image && formErrors.image && (
                        <p className="text-red-500 text-[11px] font-semibold mt-1.5 flex items-center gap-1">
                          <i className="ph-bold ph-warning-circle text-xs" /> {formErrors.image}
                        </p>
                      )}
                    </div>

                    <div className="col-span-1 md:col-span-2">
                      <label className="block text-sm font-bold text-slate-700 mb-2">Upload CV (Opsional - PDF)</label>
                      <div className="flex items-center gap-4">
                        <label className="flex items-center gap-2 px-5 py-2.5 bg-slate-50 border border-slate-200 text-slate-700 text-sm font-semibold rounded-xl cursor-pointer hover:bg-slate-100 transition shadow-sm">
                          <i className="ph-bold ph-upload-simple" />
                          {formData.cv || formData.cv_path ? 'Ganti CV' : 'Upload CV'}
                          <input
                            type="file"
                            className="hidden"
                            accept="application/pdf"
                            onChange={e => {
                              const file = e.target.files[0]
                              if (file) setFormData(prev => ({ ...prev, cv: file, cv_path: '' }))
                            }}
                          />
                        </label>
                        {(formData.cv || formData.cv_path) && (
                          <div className="flex items-center gap-2">
                            <span className="text-sm font-medium text-brand-600 bg-brand-50 px-3 py-1.5 rounded-lg border border-brand-100 flex items-center gap-2">
                              <i className="ph-fill ph-file-pdf text-brand-500"></i>
                              {formData.cv ? formData.cv.name : 'CV Tersimpan'}
                            </span>
                            <button
                              type="button"
                              onClick={() => setFormData({ ...formData, cv: null, cv_path: '' })}
                              className="w-8 h-8 flex items-center justify-center bg-red-50 text-red-500 rounded-lg hover:bg-red-100 transition border border-red-100 shadow-sm"
                              title="Hapus CV"
                            >
                              <i className="ph-bold ph-trash" />
                            </button>
                          </div>
                        )}
                      </div>
                      <p className="text-[10px] text-gray-400 mt-1.5">PNG / JPG / WEBP · maks 5MB. Foto wajib diunggah saat menambah pengurus.</p>
                    </div>

                    <div className="flex items-center gap-2 pt-2">
                      <label className="flex items-center gap-2 text-sm cursor-pointer font-medium text-slate-700">
                        <input type="checkbox" checked={formData.is_active} onChange={e => setFormData({ ...formData, is_active: e.target.checked })} className="accent-brand-600 rounded border-slate-300 text-brand-600 focus:ring-brand-500 cursor-pointer" />
                        Status Aktif
                      </label>
                    </div>
                  </div>

                  {/* Right Column (Social Media) */}
                  <div className="space-y-5">
                    <div>
                      <label className="block text-xs font-semibold text-gray-500 mb-1">Facebook URL <span className="text-gray-400 font-normal">(opsional)</span></label>
                      <input type="text" value={formData.facebook_url} onChange={e => setFormData({ ...formData, facebook_url: e.target.value })} placeholder="https://facebook.com/username" className="w-full px-3.5 py-2.5 border border-slate-300 rounded-xl text-sm outline-none bg-white focus:border-brand-500 focus:ring-2 focus:ring-brand-100 transition-colors" />
                    </div>

                    <div>
                      <label className="block text-xs font-semibold text-gray-500 mb-1">Instagram URL <span className="text-gray-400 font-normal">(opsional)</span></label>
                      <input type="text" value={formData.instagram_url} onChange={e => setFormData({ ...formData, instagram_url: e.target.value })} placeholder="https://instagram.com/username" className="w-full px-3.5 py-2.5 border border-slate-300 rounded-xl text-sm outline-none bg-white focus:border-brand-500 focus:ring-2 focus:ring-brand-100 transition-colors" />
                    </div>

                    <div>
                      <label className="block text-xs font-semibold text-gray-500 mb-1">LinkedIn URL <span className="text-gray-400 font-normal">(opsional)</span></label>
                      <input type="text" value={formData.linkedin_url} onChange={e => setFormData({ ...formData, linkedin_url: e.target.value })} placeholder="https://linkedin.com/in/username" className="w-full px-3.5 py-2.5 border border-slate-300 rounded-xl text-sm outline-none bg-white focus:border-brand-500 focus:ring-2 focus:ring-brand-100 transition-colors" />
                    </div>

                    <div>
                      <label className="block text-xs font-semibold text-gray-500 mb-1">WhatsApp <span className="text-gray-400 font-normal">(opsional)</span></label>
                      <input type="text" value={formData.whatsapp} onChange={e => setFormData({ ...formData, whatsapp: e.target.value })} placeholder="Contoh: 08123456789" className="w-full px-3.5 py-2.5 border border-slate-300 rounded-xl text-sm outline-none bg-white focus:border-brand-500 focus:ring-2 focus:ring-brand-100 transition-colors" />
                    </div>

                    <div>
                      <label className="block text-xs font-semibold text-gray-500 mb-1">Email <span className="text-gray-400 font-normal">(opsional)</span></label>
                      <input type="email" value={formData.email} onChange={e => setFormData({ ...formData, email: e.target.value })} placeholder="email@contoh.com" className="w-full px-3.5 py-2.5 border border-slate-300 rounded-xl text-sm outline-none bg-white focus:border-brand-500 focus:ring-2 focus:ring-brand-100 transition-colors" />
                    </div>

                    <div>
                      <label className="block text-xs font-semibold text-gray-500 mb-1">Pekerjaan <span className="text-gray-400 font-normal">(opsional)</span></label>
                      <input type="text" value={formData.pekerjaan} onChange={e => setFormData({ ...formData, pekerjaan: e.target.value })} placeholder="Misal: Pegawai Swasta" className="w-full px-3.5 py-2.5 border border-slate-300 rounded-xl text-sm outline-none bg-white focus:border-brand-500 focus:ring-2 focus:ring-brand-100 transition-colors" />
                    </div>
                  </div>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mt-6 animate-fade-in-up">
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Biografi <span className="text-gray-400 font-normal">(opsional)</span></label>
                    <textarea rows="4" value={formData.bio} onChange={e => setFormData({ ...formData, bio: e.target.value })} placeholder="Riwayat hidup singkat..." className="w-full px-3.5 py-2.5 border border-slate-300 rounded-xl text-sm outline-none bg-white focus:border-brand-500 focus:ring-2 focus:ring-brand-100 transition-colors" />
                  </div>
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Pendidikan <span className="text-gray-400 font-normal">(opsional)</span></label>
                    <textarea rows="4" value={formData.pendidikan} onChange={e => setFormData({ ...formData, pendidikan: e.target.value })} placeholder="Riwayat pendidikan... (Gunakan Enter untuk baris baru)" className="w-full px-3.5 py-2.5 border border-slate-300 rounded-xl text-sm outline-none bg-white focus:border-brand-500 focus:ring-2 focus:ring-brand-100 transition-colors" />
                  </div>
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Sertifikasi <span className="text-gray-400 font-normal">(opsional)</span></label>
                    <textarea rows="4" value={formData.sertifikasi} onChange={e => setFormData({ ...formData, sertifikasi: e.target.value })} placeholder="Daftar sertifikasi... (Gunakan Enter untuk baris baru)" className="w-full px-3.5 py-2.5 border border-slate-300 rounded-xl text-sm outline-none bg-white focus:border-brand-500 focus:ring-2 focus:ring-brand-100 transition-colors" />
                  </div>
                </div>

                <div className="flex justify-end gap-2 pt-4 border-t items-center mt-6">
                  {isLimited && (
                    <span className="text-xs text-amber-600 font-semibold mr-auto flex items-center gap-1">
                      <i className="ph ph-timer text-sm" /> Tunggu {cooldown}s
                    </span>
                  )}
                  <button type="button" onClick={() => setIsFormOpen(false)} disabled={cooldown > 0} className="px-4 py-2 border rounded-xl text-sm font-semibold hover:bg-slate-50 transition-colors">Batal</button>
                  <button type="submit" disabled={cooldown > 0} className="px-5 py-2 bg-brand-600 text-white rounded-xl text-sm font-semibold hover:bg-brand-700 transition-colors">Simpan</button>
                </div>
              </form>
            </div>
          </div>
        )}
      </AdminLayout>
  )
}
