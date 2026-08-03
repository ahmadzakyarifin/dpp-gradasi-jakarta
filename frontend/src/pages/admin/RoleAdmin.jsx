import { useState, useEffect, useCallback } from 'react'
import AdminLayout from '../../layouts/AdminLayout'
import { roleService } from '../../services/roleService'
import ConfirmDialog from '../../components/admin/ConfirmDialog'
import ToastNotification from '../../components/admin/ToastNotification'
import { useFormErrors, useRateLimitCooldown } from '../../utils/parseApiError'

export default function RoleAdmin() {
  const [roles, setRoles] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  const [currentTab, setCurrentTab] = useState('active') // active, trash
  const [searchQuery, setSearchQuery] = useState('')

  const [selectedItems, setSelectedItems] = useState([])

  const [isFormOpen, setIsFormOpen] = useState(false)
  const [formMode, setFormMode] = useState('create') // create, edit
  const [formData, setFormData] = useState({ name: '', display_name: '', is_active: true })
  const [formSaving, setFormSaving] = useState(false)
  const [confirmEditId, setConfirmEditId] = useState(null)

  // Dependency info sebelum delete/nonaktifkan
  const [depInfo, setDepInfo] = useState(null)

  const [confirm, setConfirm] = useState({
    isOpen: false,
    type: 'danger',
    title: '',
    message: '',
    action: null,
  })

  const [toast, setToast] = useState({ show: false, message: '', type: 'success' })
  const { fieldErrors, applyError, clearFieldError, resetFieldErrors } = useFormErrors()
  const { cooldown, isLimited, applyRateLimit } = useRateLimitCooldown()

  const showToast = useCallback((message, type = 'success') => {
    setToast({ show: true, message, type })
  }, [])

  const loadRoles = useCallback(() => {
    setLoading(true)
    setError(null)
    roleService.list({
      status: currentTab === 'trash' ? 'trashed' : undefined,
      search: searchQuery || undefined,
      limit: 100,
    })
      .then(res => {
        if (res?.data) {
          const list = Array.isArray(res.data) ? res.data : (res.data.roles || res.data.items || [])
          setRoles(list)
        }
      })
      .catch(err => setError(err?.message || 'Gagal memuat role'))
      .finally(() => setLoading(false))
  }, [currentTab, searchQuery])

  useEffect(() => {
    loadRoles()
  }, [loadRoles])

  useEffect(() => {
    setSelectedItems([])
  }, [currentTab])

  const openForm = (item = null) => {
    resetFieldErrors()
    if (item) {
      setFormMode('edit')
      setFormData({
        name: item.name || '',
        display_name: item.display_name || '',
        is_active: item.is_active !== false,
      })
    } else {
      setFormMode('create')
      setFormData({ name: '', display_name: '', is_active: true })
    }
    setIsFormOpen(true)
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setFormSaving(true)
    try {
      if (formMode === 'create') {
        await roleService.create(formData)
        showToast('Role berhasil dibuat.')
      } else {
        await roleService.update(confirmEditId, formData)
        showToast('Role berhasil diperbarui.')
      }
      setIsFormOpen(false)
      loadRoles()
    } catch (err) {
      const parsed = applyError(err)
      applyRateLimit(err)
      if (Object.keys(parsed.fieldErrors).length === 0) {
        showToast(parsed.message || 'Gagal menyimpan role.', 'error')
      }
    } finally {
      setFormSaving(false)
    }
  }

  // Simpan id yang sedang diedit (dipisah agar formData bersih)
  const askDependencyInfo = async (id, type) => {
    try {
      const res = await roleService.dependencyInfo(id)
      setDepInfo({ id, type, info: res?.data || {} })
    } catch {
      setDepInfo({ id, type, info: null })
    }
  }

  const doDelete = async (id) => {
    await roleService.remove(id)
    showToast('Role berhasil dihapus.')
    loadRoles()
  }

  const confirmAction = (type, id = null) => {
    const item = roles.find(r => r.id === id)
    const configs = {
      delete: {
        type: 'danger',
        title: 'Hapus Role',
        message: `Role "${item?.display_name || item?.name || ''}" akan dihapus. Lanjutkan?`,
        action: () => doDelete(id),
      },
      restore: {
        type: 'info',
        title: 'Pulihkan Role',
        message: `Role "${item?.display_name || item?.name || ''}" akan dipulihkan. Lanjutkan?`,
        action: async () => {
          await roleService.restore(id)
          showToast('Role berhasil dipulihkan.')
          loadRoles()
        },
      },
      bulk_delete: {
        type: 'danger',
        title: 'Hapus Massal',
        message: `${selectedItems.length} role akan dihapus. Lanjutkan?`,
        action: async () => {
          await roleService.bulkDelete(selectedItems)
          setSelectedItems([])
          showToast('Role berhasil dihapus massal.')
          loadRoles()
        },
      },
      bulk_restore: {
        type: 'info',
        title: 'Pulihkan Massal',
        message: `${selectedItems.length} role akan dipulihkan. Lanjutkan?`,
        action: async () => {
          await roleService.bulkRestore(selectedItems)
          setSelectedItems([])
          showToast('Role berhasil dipulihkan massal.')
          loadRoles()
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
      try { await action() } catch (err) {
        showToast(err?.message || 'Terjadi kesalahan, coba lagi.', 'error')
      }
    }
  }

  const isAllSelected = roles.length > 0 && selectedItems.length === roles.length
  const toggleAll = () => setSelectedItems(isAllSelected ? [] : roles.map(r => r.id))
  const toggleItem = (id) => setSelectedItems(prev => prev.includes(id) ? prev.filter(x => x !== id) : [...prev, id])

  const inputCls = "w-full px-3.5 py-2.5 border rounded-xl text-sm outline-none focus:border-brand-500 focus:ring-2 focus:ring-brand-100 transition-colors"

  return (
    <AdminLayout title="Manajemen Role">
      {toast.show && <ToastNotification message={toast.message} type={toast.type} onClose={() => setToast({ ...toast, show: false })} />}
      <ConfirmDialog {...confirm} onClose={() => setConfirm(prev => ({ ...prev, isOpen: false, action: null }))} onConfirm={executeConfirm} />

      <div className="space-y-6 animate-fade-in-up">
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
          <div className="flex bg-white rounded-lg p-1 border border-gray-200 shadow-sm">
            <button
              onClick={() => setCurrentTab('active')}
              className={`px-4 py-1.5 rounded-md text-sm transition-all duration-200 ${currentTab === 'active' ? 'bg-brand-50 text-brand-600 shadow-sm font-medium' : 'text-gray-500 hover:text-gray-700'}`}
            >
              Role Aktif
            </button>
            <button
              onClick={() => setCurrentTab('trash')}
              className={`px-4 py-1.5 rounded-md text-sm transition-all duration-200 flex items-center gap-2 ${currentTab === 'trash' ? 'bg-red-50 text-red-600 shadow-sm font-medium' : 'text-gray-500 hover:text-gray-700'}`}
            >
              <i className="ph ph-trash" /> Terhapus
            </button>
          </div>
          <div className="flex items-center gap-2 w-full max-w-md">
            <div className="relative w-full">
              <i className="ph ph-magnifying-glass absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
              <input
                type="text"
                value={searchQuery}
                onChange={e => setSearchQuery(e.target.value)}
                placeholder="Cari nama role..."
                className="w-full pl-9 pr-4 py-2 bg-gray-50 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-brand-500 focus:bg-white transition-colors"
              />
            </div>
            {currentTab === 'active' && (
              <button
                onClick={() => openForm()}
                className="shrink-0 bg-brand-600 hover:bg-brand-700 text-white px-4 py-2 rounded-lg text-sm font-medium flex items-center gap-2 transition-all btn-press shadow-sm"
              >
                <i className="ph ph-plus-circle text-lg" /> Tambah Role
              </button>
            )}
          </div>
        </div>

        {/* Bulk Actions */}
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

        {/* Table */}
        <div className="bg-white rounded-2xl border border-gray-200 overflow-hidden shadow-sm">
          {loading && <div className="py-16 text-center text-slate-500">Memuat role...</div>}
          {!loading && error && (
            <div className="py-16 text-center text-red-600 font-medium">
              <i className="ph-bold ph-warning-circle text-2xl mb-2 block mx-auto" /> {error}
            </div>
          )}
          {!loading && !error && (
            <div className="overflow-x-auto">
              <table className="w-full text-left text-sm text-slate-700">
                <thead className="bg-slate-50 border-b border-gray-200 font-semibold text-xs uppercase tracking-wider text-slate-500">
                  <tr>
                    <th className="p-4 w-12 text-center">
                      <input type="checkbox" onChange={toggleAll} checked={isAllSelected} className="rounded border-slate-300 text-brand-600 focus:ring-brand-500 cursor-pointer" />
                    </th>
                    <th className="p-4">Role</th>
                    <th className="p-4">Nama Tampilan</th>
                    <th className="p-4">Pengguna</th>
                    <th className="p-4">Status</th>
                    <th className="p-4 text-right">Aksi</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-100">
                  {roles.length === 0 && (
                    <tr><td colSpan={6} className="p-10 text-center text-slate-400">Tidak ada role untuk ditampilkan.</td></tr>
                  )}
                  {roles.map(item => (
                    <tr key={item.id} className="hover:bg-slate-50/60 transition admin-row">
                      <td className="p-4 text-center">
                        <input
                          type="checkbox"
                          checked={selectedItems.includes(item.id)}
                          onChange={() => toggleItem(item.id)}
                          disabled={item.is_system}
                          className="rounded border-slate-300 text-brand-600 focus:ring-brand-500 cursor-pointer disabled:opacity-30"
                        />
                      </td>
                      <td className="p-4">
                        <span className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-md text-xs font-bold border ${item.is_system ? 'bg-purple-50 text-purple-700 border-purple-200' : 'bg-blue-50 text-blue-700 border-blue-200'}`}>
                          {item.is_system && <i className="ph ph-crown" />}
                          {item.name}
                        </span>
                      </td>
                      <td className="p-4 font-medium text-slate-900">{item.display_name}</td>
                      <td className="p-4 text-slate-500">{item.user_count ?? 0} user</td>
                      <td className="p-4">
                        <span className={`text-xs font-semibold px-2.5 py-1 rounded-lg ${currentTab === 'trash' ? 'bg-red-50 text-red-500' : (item.is_active ? 'bg-emerald-50 text-emerald-600' : 'bg-slate-100 text-slate-500')}`}>
                          {currentTab === 'trash' ? 'Terhapus' : (item.is_active ? 'Aktif' : 'Nonaktif')}
                        </span>
                      </td>
                      <td className="p-4 text-right">
                        <div className="flex justify-end gap-2">
                          {currentTab === 'trash' ? (
                            <button onClick={() => confirmAction('restore', item.id)} className="p-2 text-slate-400 hover:text-emerald-600 rounded-lg" title="Pulihkan">
                              <i className="ph ph-arrow-counter-clockwise text-base" /> Pulihkan
                            </button>
                          ) : (
                            <>
                              <button
                                onClick={() => {
                                  setConfirmEditId(item.id)
                                  openForm(item)
                                }}
                                disabled={item.is_system}
                                className="p-2 text-slate-400 hover:text-brand-600 rounded-lg disabled:opacity-30"
                                title={item.is_system ? 'Role sistem tidak bisa diubah' : 'Edit'}
                              >
                                <i className="ph ph-pencil-simple text-base" />
                              </button>
                              <button
                                onClick={() => { askDependencyInfo(item.id, 'delete'); }}
                                disabled={item.is_system}
                                className="p-2 text-slate-400 hover:text-red-600 rounded-lg disabled:opacity-30"
                                title={item.is_system ? 'Role sistem tidak bisa dihapus' : 'Hapus'}
                              >
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
          )}
        </div>
      </div>

      {/* Dependency Info Modal */}
      {depInfo && (
        <div className="fixed inset-0 bg-black/40 backdrop-blur-[2px] flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl p-6 max-w-md w-full shadow-2xl space-y-4">
            <h3 className="font-heading font-bold text-lg text-slate-900">Cek Dependensi Role</h3>
            {depInfo.info ? (
              <>
                {depInfo.info.has_dependencies ? (
                  <p className="text-sm text-amber-700 bg-amber-50 border border-amber-200 rounded-xl p-3">
                    <i className="ph ph-warning-circle mr-1" />
                    {depInfo.info.message || `Role ini dipakai oleh ${depInfo.info.user_count} pengguna.`}
                  </p>
                ) : (
                  <p className="text-sm text-emerald-700 bg-emerald-50 border border-emerald-200 rounded-xl p-3">
                    <i className="ph ph-check-circle mr-1" />
                    {depInfo.info.message || 'Role ini tidak dipakai pengguna lain, aman untuk dihapus.'}
                  </p>
                )}
                <p className="text-xs text-slate-500">Jumlah pengguna memakai role ini: <b>{depInfo.info.user_count ?? 0}</b></p>
              </>
            ) : (
              <p className="text-sm text-slate-500">Gagal memuat info dependensi.</p>
            )}
            <div className="flex justify-end gap-2 pt-2 border-t">
              <button onClick={() => setDepInfo(null)} className="px-4 py-2 border rounded-xl text-sm font-semibold hover:bg-slate-50">Tutup</button>
              {depInfo.info && !depInfo.info.has_dependencies && depInfo.type === 'delete' && (
                <button
                  onClick={() => { setDepInfo(null); confirmAction('delete', depInfo.id) }}
                  className="px-4 py-2 bg-red-600 text-white rounded-xl text-sm font-semibold hover:bg-red-700"
                >
                  Hapus Role
                </button>
              )}
            </div>
          </div>
        </div>
      )}

      {/* Form Modal */}
      {isFormOpen && (
        <div className="fixed inset-0 bg-black/40 backdrop-blur-[2px] flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl p-6 max-w-md w-full shadow-2xl space-y-4">
            <h3 className="font-heading font-bold text-lg text-slate-900">{formMode === 'create' ? 'Tambah Role Baru' : 'Edit Role'}</h3>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-xs font-semibold text-gray-500 mb-1">Nama Role (slug) *</label>
                <input
                  type="text"
                  value={formData.name}
                  onChange={e => { setFormData({ ...formData, name: e.target.value }); clearFieldError('name') }}
                  placeholder="contoh: admin_berita"
                  className={`${inputCls} ${fieldErrors.name ? 'border-red-400' : 'border-gray-300'}`}
                />
                {fieldErrors.name && (
                  <p className="text-red-500 text-[11px] font-semibold mt-1.5 flex items-center gap-1">
                    <i className="ph-bold ph-warning-circle text-xs" /> {fieldErrors.name}
                  </p>
                )}
              </div>
              <div>
                <label className="block text-xs font-semibold text-gray-500 mb-1">Nama Tampilan *</label>
                <input
                  type="text"
                  value={formData.display_name}
                  onChange={e => { setFormData({ ...formData, display_name: e.target.value }); clearFieldError('display_name') }}
                  placeholder="contoh: Admin Berita"
                  className={`${inputCls} ${fieldErrors.display_name ? 'border-red-400' : 'border-gray-300'}`}
                />
                {fieldErrors.display_name && (
                  <p className="text-red-500 text-[11px] font-semibold mt-1.5 flex items-center gap-1">
                    <i className="ph-bold ph-warning-circle text-xs" /> {fieldErrors.display_name}
                  </p>
                )}
              </div>
              <div>
                <label className="flex items-center gap-2 text-sm cursor-pointer">
                  <input type="checkbox" checked={formData.is_active} onChange={e => setFormData({ ...formData, is_active: e.target.checked })} className="accent-brand-600" />
                  Aktif (bisa dipakai login)
                </label>
              </div>
              <div className="flex justify-end gap-2 pt-4 border-t items-center">
                {isLimited && (
                  <span className="text-xs text-amber-600 font-semibold mr-auto flex items-center gap-1">
                    <i className="ph ph-timer text-sm" /> Tunggu {cooldown}s
                  </span>
                )}
                <button type="button" onClick={() => setIsFormOpen(false)} disabled={formSaving || isLimited} className="px-4 py-2 border rounded-xl text-sm font-semibold">Batal</button>
                <button type="submit" disabled={formSaving || isLimited} className="px-5 py-2 bg-brand-600 text-white rounded-xl text-sm font-semibold">
                  {formSaving ? 'Menyimpan...' : (isLimited ? `Tunggu ${cooldown}s` : 'Simpan')}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </AdminLayout>
  )
}
