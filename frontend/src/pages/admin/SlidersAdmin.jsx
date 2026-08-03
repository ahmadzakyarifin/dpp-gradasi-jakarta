import { useState, useEffect, useCallback } from 'react'
import AdminLayout from '../../layouts/AdminLayout'
import { slidersService } from '../../services/slidersService'
import ConfirmDialog from '../../components/admin/ConfirmDialog'
import ToastNotification from '../../components/admin/ToastNotification'
import { resolveAssetUrl } from '../../utils/assetUrl'

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

  const [confirm, setConfirm] = useState({ isOpen: false, type: 'danger', title: '', message: '', action: null })
  const [toast, setToast] = useState({ show: false, message: '', type: 'success' })

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
      showToast(err.message || 'Gagal menyimpan slider', 'error')
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
                        <div className="flex items-center justify-end gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
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
                <div className="py-12 text-center">
                  <i className="ph ph-image text-4xl text-gray-300 mb-2" />
                  <p className="text-gray-500 text-sm">Tidak ada data slider untuk ditampilkan.</p>
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

      {/* FORM MODAL (Create/Edit) — mirip sliders.html */}
      {isFormOpen && (
        <div className="fixed inset-0 z-50 overflow-y-auto" role="dialog" aria-modal="true">
          <div className="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
            <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" onClick={() => setIsFormOpen(false)} />
            <span className="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>
            <div className="inline-block align-bottom bg-white rounded-xl text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-2xl sm:w-full">
              <div className="bg-white">
                <div className="border-b border-gray-200 px-6 py-4 flex items-center justify-between">
                  <h3 className="text-lg leading-6 font-heading font-semibold text-gray-900">{formMode === 'create' ? 'Tambah Slider' : 'Edit Slider'}</h3>
                  <button onClick={() => setIsFormOpen(false)} className="text-gray-400 hover:text-gray-500">
                    <i className="ph-bold ph-x text-xl" />
                  </button>
                </div>
                <form onSubmit={handleSubmit} className="px-6 py-4 space-y-5 max-h-[60vh] overflow-y-auto">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
                    <div className="md:col-span-2">
                      <label className="block text-sm font-medium text-gray-700 mb-1">Judul Utama *</label>
                      <input
                        type="text"
                        value={formData.title}
                        onChange={e => setFormData({ ...formData, title: e.target.value })}
                        required
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-brand-500 focus:border-brand-500 text-sm outline-none"
                      />
                    </div>
                    <div className="md:col-span-2">
                      <label className="block text-sm font-medium text-gray-700 mb-1">Sub-judul (Opsional)</label>
                      <input
                        type="text"
                        value={formData.subtitle}
                        onChange={e => setFormData({ ...formData, subtitle: e.target.value })}
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-brand-500 focus:border-brand-500 text-sm outline-none"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">Badge Tag (Opsional)</label>
                      <input
                        type="text"
                        value={formData.tag}
                        onChange={e => setFormData({ ...formData, tag: e.target.value })}
                        placeholder="Misal: Webinar, Event"
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-brand-500 focus:border-brand-500 text-sm outline-none"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">Urutan (Sort Order) *</label>
                      <input
                        type="number"
                        value={formData.sort_order}
                        onChange={e => setFormData({ ...formData, sort_order: Number(e.target.value) })}
                        required
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-brand-500 focus:border-brand-500 text-sm outline-none"
                      />
                    </div>
                    <div className="md:col-span-2">
                      <label className="block text-sm font-medium text-gray-700 mb-1">Gambar Slider (URL) *</label>
                      <input
                        type="url"
                        value={formData.image_path}
                        onChange={e => setFormData({ ...formData, image_path: e.target.value })}
                        required
                        placeholder="https://example.com/banner.jpg"
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-brand-500 focus:border-brand-500 text-sm outline-none"
                      />
                      <p className="text-xs text-gray-500 mt-1">Rekomendasi ukuran: 1920x600 px (Rasio lebar)</p>
                    </div>
                    <div className="md:col-span-2">
                      <label className="block text-sm font-medium text-gray-700 mb-1">Link URL (Jika diklik mengarah kemana)</label>
                      <input
                        type="url"
                        value={formData.link_url}
                        onChange={e => setFormData({ ...formData, link_url: e.target.value })}
                        placeholder="https://..."
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-brand-500 focus:border-brand-500 text-sm outline-none"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">Tanggal Kegiatan (Opsional di slider)</label>
                      <input
                        type="text"
                        value={formData.event_date}
                        onChange={e => setFormData({ ...formData, event_date: e.target.value })}
                        placeholder="20 Okt 2024"
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-brand-500 focus:border-brand-500 text-sm outline-none"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">Lokasi (Opsional di slider)</label>
                      <input
                        type="text"
                        value={formData.location}
                        onChange={e => setFormData({ ...formData, location: e.target.value })}
                        placeholder="Jakarta"
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-brand-500 focus:border-brand-500 text-sm outline-none"
                      />
                    </div>
                    <div className="md:col-span-2 flex gap-6">
                      <label className="flex items-center gap-2 cursor-pointer">
                        <input
                          type="checkbox"
                          checked={formData.is_new}
                          onChange={e => setFormData({ ...formData, is_new: e.target.checked })}
                          className="rounded border-gray-300 text-brand-600 focus:ring-brand-500 accent-brand-600"
                        />
                        <span className="text-sm font-medium text-gray-700">Tandai sebagai BARU (badge NEW)</span>
                      </label>
                      <label className="flex items-center gap-2 cursor-pointer">
                        <input
                          type="checkbox"
                          checked={formData.is_active}
                          onChange={e => setFormData({ ...formData, is_active: e.target.checked })}
                          className="rounded border-gray-300 text-brand-600 focus:ring-brand-500 accent-brand-600"
                        />
                        <span className="text-sm font-medium text-gray-700">Slider Aktif</span>
                      </label>
                    </div>
                  </div>
                </form>
              </div>
              <div className="bg-gray-50 px-6 py-4 border-t border-gray-200 flex justify-end gap-3">
                <button
                  type="button"
                  onClick={() => setIsFormOpen(false)}
                  className="px-4 py-2 bg-white border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50"
                >
                  Batal
                </button>
                <button
                  type="button"
                  onClick={handleSubmit}
                  className="px-4 py-2 bg-brand-600 border border-transparent rounded-lg text-sm font-medium text-white hover:bg-brand-700"
                >
                  Simpan Data
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </AdminLayout>
  )
}
