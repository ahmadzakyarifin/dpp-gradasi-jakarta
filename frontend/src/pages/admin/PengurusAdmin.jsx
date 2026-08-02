import { useState, useEffect, useCallback } from 'react'
import AdminLayout from '../../layouts/AdminLayout'
import { pengurusService } from '../../services/pengurusService'
import ConfirmDialog from '../../components/admin/ConfirmDialog'
import ToastNotification from '../../components/admin/ToastNotification'
import { resolveAssetUrl } from '../../utils/assetUrl'

const PAGE_SIZE = 10

export default function PengurusAdmin() {
  const [items, setItems] = useState([])
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
    level: 'dpp',
    role: '',
    department: '',
    periode: '2025 - 2030',
    provinsi: '',
    kabupaten: '',
    facebook_url: '',
    instagram_url: '',
    linkedin_url: '',
    whatsapp: '',
    sort_order: 1,
    image_url: '',
    image: null,
    is_active: true,
  })

  const [confirm, setConfirm] = useState({ isOpen: false, type: 'danger', title: '', message: '', action: null })
  const [toast, setToast] = useState({ show: false, message: '', type: 'success' })

  const showToast = useCallback((message, type = 'success') => {
    setToast({ show: true, message, type })
    setTimeout(() => setToast(prev => ({ ...prev, show: false })), 3000)
  }, [])

  const loadPengurus = useCallback(() => {
    setLoading(true)
    const params = {
      page: currentPage,
      limit: PAGE_SIZE,
      search: searchQuery,
      status: currentTab === 'trash' ? 'all' : (currentTab === 'inactive' ? 'inactive' : 'active'),
      trashed: currentTab === 'trash',
      level: filterLevel || undefined,
      sort: 'sort_order',
    }
    pengurusService.listAdmin(params)
      .then(res => {
        if (res && res.data) {
          const list = Array.isArray(res.data) ? res.data : (res.data.data || [])
          setItems(list)
          if (res.data.meta) setMeta(res.data.meta)
        } else {
          setItems([])
        }
      })
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [currentTab, currentPage, searchQuery, filterLevel])

  useEffect(() => {
    loadPengurus()
  }, [loadPengurus])

  useEffect(() => {
    setCurrentPage(1)
    setSelectedItems([])
  }, [currentTab, searchQuery, filterLevel])

  const openForm = (item = null) => {
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
        sort_order: item.sort_order || 1,
        image_url: item.image_url || '',
        image: null,
        is_active: item.is_active !== false,
      })
    } else {
      setFormMode('create')
      setFormData({
        id: null,
        name: '',
        level: 'dpp',
        role: '',
        department: '',
        periode: '2025 - 2030',
        provinsi: '',
        kabupaten: '',
        facebook_url: '',
        instagram_url: '',
        linkedin_url: '',
        whatsapp: '',
        sort_order: (meta.total_data || 0) + 1,
        image_url: '',
        image: null,
        is_active: true,
      })
    }
    setIsFormOpen(true)
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    const payload = { ...formData }
    // Backend butuh multipart; buang id (meta) — image_url dikirim sebagai URL string
    if (formMode === 'create') {
      try {
        await pengurusService.create(payload)
        showToast('Pengurus berhasil ditambahkan.')
        setIsFormOpen(false)
        loadPengurus()
      } catch (err) {
        showToast(err.message || 'Gagal menyimpan pengurus', 'error')
      }
    } else {
      try {
        await pengurusService.update(formData.id, payload)
        showToast('Pengurus berhasil diperbarui.')
        setIsFormOpen(false)
        loadPengurus()
      } catch (err) {
        showToast(err.message || 'Gagal menyimpan pengurus', 'error')
      }
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

  const isAllSelected = items.length > 0 && selectedItems.length === items.length
  function toggleAll() {
    setSelectedItems(isAllSelected ? [] : items.map(i => i.id))
  }
  function toggleOne(id) {
    setSelectedItems(prev => prev.includes(id) ? prev.filter(x => x !== id) : [...prev, id])
  }

  const totalPages = Math.max(1, meta.total_pages || 1)
  const totalData = meta.total_data ?? items.length

  const levelLabel = { ketua: 'Ketua Umum', dpp: 'Pusat (DPP)', dpd: 'Provinsi (DPD)', dpc: 'Kab/Kota (DPC)' }

  return (
    <AdminLayout title="Kelola Pengurus">
      <ConfirmDialog {...confirm} onClose={() => setConfirm(prev => ({ ...prev, isOpen: false, action: null }))} onConfirm={executeConfirm} />
      <ToastNotification show={toast.show} message={toast.message} type={toast.type} />

      <div className="space-y-6">
        {/* Navigation Tabs */}
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
          <div className="flex bg-white rounded-xl p-1 border border-gray-200 shadow-sm">
            <button
              onClick={() => setCurrentTab('active')}
              className={`px-4 py-2 rounded-lg text-sm font-semibold transition-all ${currentTab === 'active' ? 'bg-brand-50 text-brand-600 shadow-sm' : 'text-slate-500 hover:text-slate-700'}`}
            >
              Pengurus Aktif
            </button>
            <button
              onClick={() => setCurrentTab('inactive')}
              className={`px-4 py-2 rounded-lg text-sm font-semibold transition-all ${currentTab === 'inactive' ? 'bg-brand-50 text-brand-600 shadow-sm' : 'text-slate-500 hover:text-slate-700'}`}
            >
              Non-aktif
            </button>
            <button
              onClick={() => setCurrentTab('trash')}
              className={`px-4 py-2 rounded-lg text-sm font-semibold transition-all ${currentTab === 'trash' ? 'bg-brand-50 text-brand-600 shadow-sm' : 'text-slate-500 hover:text-slate-700'}`}
            >
              Sampah
            </button>
          </div>
          {currentTab !== 'trash' && (
            <button
              onClick={() => openForm()}
              className="px-4 py-2.5 bg-brand-600 text-white rounded-xl hover:bg-brand-700 text-sm font-semibold shadow-sm"
            >
              + Tambah Pengurus Baru
            </button>
          )}
        </div>

        {/* Filter Bar */}
        <div className="flex flex-col md:flex-row gap-4 justify-between items-center bg-white p-4 rounded-xl border border-gray-200 shadow-sm">
          <input
            type="text"
            placeholder="Cari nama atau jabatan pengurus..."
            value={searchQuery}
            onChange={e => setSearchQuery(e.target.value)}
            className="w-full md:max-w-md px-3.5 py-2.5 bg-slate-50 border border-gray-200 rounded-xl text-sm outline-none"
          />
          <div className="flex gap-2 w-full md:w-auto">
            <select
              value={filterLevel}
              onChange={e => setFilterLevel(e.target.value)}
              className="px-3.5 py-2.5 bg-slate-50 border border-gray-200 rounded-xl text-sm text-slate-600 outline-none w-full sm:w-auto"
            >
              <option value="">Semua Tingkatan</option>
              <option value="ketua">Ketua Umum</option>
              <option value="dpp">Pusat (DPP)</option>
              <option value="dpd">Provinsi (DPD)</option>
              <option value="dpc">Kab/Kota (DPC)</option>
            </select>
          </div>
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
                      <th className="p-4 text-right">Aksi</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-100">
                    {items.map(item => (
                      <tr key={item.id} className="hover:bg-slate-50/50">
                        <td className="p-4">
                          <input type="checkbox" checked={selectedItems.includes(item.id)} onChange={() => toggleOne(item.id)} className="accent-brand-600" />
                        </td>
                        <td className="p-4 font-medium text-slate-900">
                          <div className="flex items-center gap-3">
                            {item.image_url && (
                              <img src={resolveAssetUrl(item.image_url)} alt="" className="w-8 h-8 rounded-full object-cover border" />
                            )}
                            <div>
                              <div>{item.name}</div>
                              <div className="text-xs text-slate-400">{item.periode}</div>
                            </div>
                          </div>
                        </td>
                        <td className="p-4 uppercase font-bold text-xs text-brand-600">{levelLabel[item.level] || item.level}</td>
                        <td className="p-4">{item.role}</td>
                        <td className="p-4 text-slate-500 text-xs">
                          {item.provinsi ? `${item.provinsi}${item.kabupaten ? `, ${item.kabupaten}` : ''}` : '-'}
                        </td>
                        <td className="p-4 text-right">
                          <div className="flex justify-end gap-3">
                            {currentTab === 'trash' ? (
                              <button onClick={() => confirmAction('restore', item)} title="Pulihkan" className="text-xs text-emerald-600 font-semibold hover:underline">
                                Pulihkan
                              </button>
                            ) : (
                              <>
                                <button onClick={() => openForm(item)} className="text-xs text-slate-600 font-semibold hover:underline">Edit</button>
                                <button onClick={() => confirmAction('delete', item)} className="text-xs text-red-600 font-semibold hover:underline">Hapus</button>
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
                  <span className="text-xs text-slate-500">Hal {currentPage} dari {totalPages} · {totalData} data</span>
                  <div className="flex items-center gap-1.5">
                    <button
                      disabled={currentPage <= 1}
                      onClick={() => setCurrentPage(p => p - 1)}
                      className="w-8 h-8 flex items-center justify-center rounded-lg border border-gray-200 text-slate-500 hover:border-brand-500 hover:text-brand-600 disabled:opacity-40"
                    >
                      <i className="ph-bold ph-caret-left" />
                    </button>
                    {Array.from({ length: totalPages }, (_, i) => i + 1).map(page => (
                      <button
                        key={page}
                        onClick={() => setCurrentPage(page)}
                        className={`w-8 h-8 flex items-center justify-center rounded-lg text-xs font-bold ${page === currentPage ? 'bg-brand-700 text-white' : 'border border-gray-200 text-slate-600 hover:border-brand-500'}`}
                      >
                        {page}
                      </button>
                    ))}
                    <button
                      disabled={currentPage >= totalPages}
                      onClick={() => setCurrentPage(p => p + 1)}
                      className="w-8 h-8 flex items-center justify-center rounded-lg border border-gray-200 text-slate-500 hover:border-brand-500 hover:text-brand-600 disabled:opacity-40"
                    >
                      <i className="ph-bold ph-caret-right" />
                    </button>
                  </div>
                </div>
              )}
            </>
          )}
        </div>

        {/* Form Modal */}
        {isFormOpen && (
          <div className="fixed inset-0 bg-black/40 backdrop-blur-[2px] flex items-center justify-center z-50 p-4">
            <div className="bg-white rounded-2xl p-6 max-w-lg w-full shadow-2xl space-y-4 max-h-[90vh] overflow-y-auto">
              <h3 className="font-heading font-bold text-lg text-slate-900">{formMode === 'create' ? 'Tambah Pengurus Baru' : 'Edit Pengurus'}</h3>
              <form onSubmit={handleSubmit} className="space-y-4">
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Nama Lengkap *</label>
                  <input type="text" value={formData.name} onChange={e => setFormData({ ...formData, name: e.target.value })} required className="w-full px-3 py-2 border border-gray-300 rounded-xl text-sm outline-none" />
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Tingkat Struktur *</label>
                    <select value={formData.level} onChange={e => setFormData({ ...formData, level: e.target.value })} className="w-full px-3 py-2 border border-gray-300 rounded-xl text-sm outline-none bg-white">
                      <option value="ketua">Ketua Umum</option>
                      <option value="dpp">Pusat (DPP)</option>
                      <option value="dpd">Provinsi (DPD)</option>
                      <option value="dpc">Kab/Kota (DPC)</option>
                    </select>
                  </div>
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Jabatan Resmi *</label>
                    <input type="text" value={formData.role} onChange={e => setFormData({ ...formData, role: e.target.value })} required placeholder="Misal: Ketua Bidang Organisasi" className="w-full px-3 py-2 border border-gray-300 rounded-xl text-sm outline-none" />
                  </div>
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Periode *</label>
                    <input type="text" value={formData.periode} onChange={e => setFormData({ ...formData, periode: e.target.value })} required placeholder="2025 - 2030" className="w-full px-3 py-2 border border-gray-300 rounded-xl text-sm outline-none" />
                  </div>
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Departemen</label>
                    <input type="text" value={formData.department} onChange={e => setFormData({ ...formData, department: e.target.value })} className="w-full px-3 py-2 border border-gray-300 rounded-xl text-sm outline-none" />
                  </div>
                </div>
                {formData.level === 'dpd' || formData.level === 'dpc' ? (
                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <label className="block text-xs font-semibold text-gray-500 mb-1">Provinsi {formData.level === 'dpc' ? '*' : ''}</label>
                      <input type="text" value={formData.provinsi} onChange={e => setFormData({ ...formData, provinsi: e.target.value })} required={formData.level === 'dpd' || formData.level === 'dpc'} className="w-full px-3 py-2 border border-gray-300 rounded-xl text-sm outline-none" />
                    </div>
                    {formData.level === 'dpc' && (
                      <div>
                        <label className="block text-xs font-semibold text-gray-500 mb-1">Kabupaten/Kota *</label>
                        <input type="text" value={formData.kabupaten} onChange={e => setFormData({ ...formData, kabupaten: e.target.value })} required className="w-full px-3 py-2 border border-gray-300 rounded-xl text-sm outline-none" />
                      </div>
                    )}
                  </div>
                ) : (
                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <label className="block text-xs font-semibold text-gray-500 mb-1">Provinsi</label>
                      <input type="text" value={formData.provinsi} onChange={e => setFormData({ ...formData, provinsi: e.target.value })} className="w-full px-3 py-2 border border-gray-300 rounded-xl text-sm outline-none" />
                    </div>
                    <div>
                      <label className="block text-xs font-semibold text-gray-500 mb-1">Kabupaten/Kota</label>
                      <input type="text" value={formData.kabupaten} onChange={e => setFormData({ ...formData, kabupaten: e.target.value })} className="w-full px-3 py-2 border border-gray-300 rounded-xl text-sm outline-none" />
                    </div>
                  </div>
                )}
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Foto Profil {formMode === 'create' ? '*' : ''}</label>
                  <input type="file" accept="image/*" required={formMode === 'create'} onChange={e => {
                    const file = e.target.files && e.target.files[0]
                    if (file) setFormData(prev => ({ ...prev, image: file, image_url: '' }))
                  }} className="w-full px-3 py-2 border border-gray-300 rounded-xl text-sm outline-none" />
                  {formMode === 'create' && (
                    <p className="text-xs text-slate-400 mt-1">Foto wajib diunggah saat menambah pengurus.</p>
                  )}
                  {formMode === 'edit' && formData.image_url && !formData.image && (
                    <p className="text-xs text-slate-400 mt-1">Foto saat ini: {formData.image_url.split('/').pop()}</p>
                  )}
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Facebook URL</label>
                    <input type="text" value={formData.facebook_url} onChange={e => setFormData({ ...formData, facebook_url: e.target.value })} className="w-full px-3 py-2 border border-gray-300 rounded-xl text-sm outline-none" />
                  </div>
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Instagram URL</label>
                    <input type="text" value={formData.instagram_url} onChange={e => setFormData({ ...formData, instagram_url: e.target.value })} className="w-full px-3 py-2 border border-gray-300 rounded-xl text-sm outline-none" />
                  </div>
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">LinkedIn URL</label>
                    <input type="text" value={formData.linkedin_url} onChange={e => setFormData({ ...formData, linkedin_url: e.target.value })} className="w-full px-3 py-2 border border-gray-300 rounded-xl text-sm outline-none" />
                  </div>
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">WhatsApp</label>
                    <input type="text" value={formData.whatsapp} onChange={e => setFormData({ ...formData, whatsapp: e.target.value })} className="w-full px-3 py-2 border border-gray-300 rounded-xl text-sm outline-none" />
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <label className="flex items-center gap-2 text-sm cursor-pointer">
                    <input type="checkbox" checked={formData.is_active} onChange={e => setFormData({ ...formData, is_active: e.target.checked })} className="accent-brand-600" />
                    Aktif
                  </label>
                </div>
                <div className="flex justify-end gap-2 pt-4 border-t">
                  <button type="button" onClick={() => setIsFormOpen(false)} className="px-4 py-2 border rounded-xl text-xs font-semibold">Batal</button>
                  <button type="submit" className="px-4 py-2 bg-brand-600 text-white rounded-xl hover:bg-brand-700 text-xs font-semibold">Simpan</button>
                </div>
              </form>
            </div>
          </div>
        )}
      </div>
    </AdminLayout>
  )
}
