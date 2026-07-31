import { useState, useEffect } from 'react'
import AdminLayout from '../../layouts/AdminLayout'
import { userService } from '../../services/userService'

export default function UsersAdmin() {
  const [users, setUsers] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  const [currentTab, setCurrentTab] = useState('active') // active, pending, trash
  const [search, setSearch] = useState('')

  // Pagination
  const [page, setPage] = useState(1)
  const [limit] = useState(10)
  const [total, setTotal] = useState(0)
  const [totalPages, setTotalPages] = useState(1)

  // Selection
  const [selectedIds, setSelectedIds] = useState([])

  // Modal Form (Create)
  const [isFormOpen, setIsFormOpen] = useState(false)
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    phone: '',
    role_id: '2' // 2: Admin
  })

  // Confirm Modal
  const [confirm, setConfirm] = useState({
    isOpen: false,
    type: '', // delete, deactivate, restore, resend, bulk_delete, bulk_restore
    id: null,
    title: '',
    message: '',
    icon: ''
  })

  // Toast Notification
  const [toast, setToast] = useState({
    show: false,
    message: '',
    type: 'success' // success, error
  })

  const showToast = (message, type = 'success') => {
    setToast({ show: true, message, type })
    setTimeout(() => {
      setToast(prev => ({ ...prev, show: false }))
    }, 3000)
  }

  const fetchUsers = () => {
    setLoading(true)
    userService.list({
      tab: currentTab,
      search,
      page,
      limit
    })
      .then(res => {
        if (res.success && res.data) {
          setUsers(res.data.items || res.data.users || [])
          if (res.data.pagination) {
            setTotal(res.data.pagination.total || 0)
            setTotalPages(res.data.pagination.totalPages || 1)
          }
          setError(null)
        } else {
          setError('Gagal memuat data admin')
        }
      })
      .catch((err) => {
        setError(err.message || 'Kesalahan koneksi ke server')
      })
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    setSelectedIds([])
    setPage(1)
    fetchUsers()
  }, [currentTab])

  const handleSearchSubmit = (e) => {
    e.preventDefault()
    setPage(1)
    fetchUsers()
  }

  const handleReset = () => {
    setSearch('')
    setPage(1)
    setLoading(true)
    userService.list({
      tab: currentTab,
      search: '',
      page: 1,
      limit
    })
      .then(res => {
        if (res.success && res.data) {
          setUsers(res.data.items || res.data.users || [])
          if (res.data.pagination) {
            setTotal(res.data.pagination.total || 0)
            setTotalPages(res.data.pagination.totalPages || 1)
          }
          setError(null)
        }
      })
      .catch((err) => setError(err.message || 'Kesalahan koneksi ke server'))
      .finally(() => setLoading(false))
  }

  // Row Selection logic
  const selectableUsers = users.filter(
    item => item.role !== 'Super Admin' && item.status !== 'pending_activation'
  )

  const isAllSelected = selectableUsers.length > 0 && selectedIds.length === selectableUsers.length

  const handleSelectAll = (e) => {
    if (e.target.checked) {
      setSelectedIds(selectableUsers.map(m => m.id))
    } else {
      setSelectedIds([])
    }
  }

  const handleSelectItem = (id) => {
    setSelectedIds(prev => 
      prev.includes(id) ? prev.filter(item => item !== id) : [...prev, id]
    )
  }

  // Handle Form Submit
  const handleFormSubmit = async (e) => {
    e.preventDefault()
    try {
      await userService.create(formData)
      setIsFormOpen(false)
      setFormData({ name: '', email: '', phone: '', role_id: '2' })
      showToast('Undangan berhasil dikirim ke email admin baru!', 'success')
      fetchUsers()
    } catch (err) {
      showToast(err.message || 'Gagal mengirim undangan admin baru.', 'error')
    }
  }

  // Trigger Confirmation Dialog
  const triggerConfirm = (type, id = null) => {
    let title = ''
    let message = ''
    let icon = ''

    if (type === 'delete') {
      title = 'Hapus Akses Admin'
      message = 'Akun admin ini akan dinonaktifkan dan masuk ke daftar sampah.'
      icon = 'ph-trash'
    } else if (type === 'deactivate') {
      title = 'Nonaktifkan Admin'
      message = 'Akun admin ini akan dinonaktifkan (inactive) dan tidak bisa login. Lanjutkan?'
      icon = 'ph-prohibit'
    } else if (type === 'restore') {
      title = 'Pulihkan Akses Admin'
      message = 'Akun admin ini akan kembali aktif.'
      icon = 'ph-arrow-counter-clockwise'
    } else if (type === 'resend') {
      const item = users.find(i => i.id === id)
      title = 'Kirim Ulang Undangan'
      message = `Email undangan aktivasi akan dikirim ulang ke ${item ? item.email : 'admin'}. Lanjutkan?`
      icon = 'ph-paper-plane-right'
    } else if (type === 'bulk_delete') {
      title = 'Hapus Massal'
      message = `Nonaktifkan ${selectedIds.length} akun admin terpilih?`
      icon = 'ph-trash'
    } else if (type === 'bulk_restore') {
      title = 'Pulihkan Massal'
      message = `Pulihkan ${selectedIds.length} akun admin terpilih?`
      icon = 'ph-arrow-counter-clockwise'
    }

    setConfirm({
      isOpen: true,
      type,
      id,
      title,
      message,
      icon
    })
  }

  const executeAction = async () => {
    const { type, id } = confirm
    setConfirm(prev => ({ ...prev, isOpen: false }))

    try {
      if (type === 'resend') {
        await userService.resendActivation(id)
        showToast('Undangan aktivasi berhasil dikirim ulang.', 'success')
      } else if (type === 'deactivate') {
        await userService.setStatus(id, 'inactive')
        showToast('Akun admin dinonaktifkan.', 'success')
      } else if (type === 'delete') {
        await userService.remove(id)
        showToast('Akun admin berhasil dihapus.', 'success')
      } else if (type === 'restore') {
        await userService.restore(id)
        showToast('Akun admin berhasil dipulihkan.', 'success')
      } else if (type === 'bulk_delete') {
        await userService.bulkDelete(selectedIds)
        setSelectedIds([])
        showToast('Akun admin terpilih berhasil dihapus.', 'success')
      } else if (type === 'bulk_restore') {
        await userService.bulkRestore(selectedIds)
        setSelectedIds([])
        showToast('Akun admin terpilih berhasil dipulihkan.', 'success')
      }
      fetchUsers()
    } catch (err) {
      showToast(err.message || 'Gagal melakukan aksi.', 'error')
    }
  }

  const renderPagination = () => {
    const pages = []
    for (let i = 1; i <= totalPages; i++) {
      pages.push(
        <button
          key={i}
          onClick={() => setPage(i)}
          className={`px-3 py-1 border rounded transition ${
            page === i
              ? 'border-brand-500 bg-brand-50 text-brand-600 font-medium'
              : 'border-gray-200 bg-white text-gray-500 hover:bg-gray-50'
          }`}
        >
          {i}
        </button>
      )
    }
    return pages
  }

  return (
    <AdminLayout title="Manajemen Admin">
      <div className="max-w-7xl mx-auto space-y-6">
        {/* Header Actions & Filter */}
        <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4 bg-white p-4 rounded-xl border border-gray-200 shadow-sm">
          {/* Tabs */}
          <div className="flex bg-gray-100 rounded-lg p-1 border border-gray-200 shadow-sm flex-wrap sm:flex-nowrap">
            <button
              onClick={() => setCurrentTab('active')}
              className={`px-4 py-1.5 rounded-md text-sm transition-all ${
                currentTab === 'active'
                  ? 'bg-white text-brand-600 shadow-sm font-medium'
                  : 'text-gray-500 hover:text-gray-700'
              }`}
            >
              Admin Aktif
            </button>
            <button
              onClick={() => setCurrentTab('pending')}
              className={`px-4 py-1.5 rounded-md text-sm transition-all flex items-center gap-2 ${
                currentTab === 'pending'
                  ? 'bg-white text-amber-600 shadow-sm font-medium'
                  : 'text-gray-500 hover:text-gray-700'
              }`}
            >
              <i className="ph ph-clock-countdown" /> Menunggu Aktivasi
            </button>
            <button
              onClick={() => setCurrentTab('trash')}
              className={`px-4 py-1.5 rounded-md text-sm transition-all flex items-center gap-2 ${
                currentTab === 'trash'
                  ? 'bg-white text-red-600 shadow-sm font-medium'
                  : 'text-gray-500 hover:text-gray-700'
              }`}
            >
              <i className="ph ph-trash" /> Akun Terhapus
            </button>
          </div>

          {/* Right Actions: Search & Add */}
          <div className="flex items-center gap-2 w-full md:w-auto">
            <form onSubmit={handleSearchSubmit} className="flex items-center gap-2 flex-1 md:flex-initial">
              <div className="relative w-full md:w-64">
                <i className="ph ph-magnifying-glass absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
                <input
                  type="text"
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                  placeholder="Cari nama / email..."
                  className="w-full pl-9 pr-4 py-2 bg-gray-50 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-brand-500 focus:bg-white transition-colors"
                />
              </div>
              <button
                type="submit"
                className="bg-brand-600 hover:bg-brand-700 text-white px-3 py-2 rounded-lg text-sm font-medium transition"
              >
                Cari
              </button>
              <button
                type="button"
                onClick={handleReset}
                className="shrink-0 bg-gray-50 border border-gray-200 text-gray-700 hover:bg-gray-100 px-3 py-2 rounded-lg text-sm font-medium flex items-center gap-2 transition"
                title="Reset Filter"
              >
                <i className="ph ph-arrows-counter-clockwise text-lg" />
              </button>
            </form>

            {currentTab === 'active' && (
              <button
                onClick={() => setIsFormOpen(true)}
                className="shrink-0 bg-brand-600 hover:bg-brand-700 text-white px-4 py-2 rounded-lg text-sm font-medium flex items-center gap-2 transition-colors shadow-sm"
              >
                <i className="ph ph-user-plus text-lg" /> Tambah
              </button>
            )}
          </div>
        </div>

        {/* Bulk Actions Bar */}
        {selectedIds.length > 0 && (
          <div className="bg-indigo-50 border border-indigo-100 rounded-lg p-3 flex items-center justify-between shadow-sm animate-fade-in">
            <span className="text-sm text-indigo-800 font-medium">{selectedIds.length} akun terpilih</span>
            <div className="flex gap-2">
              {currentTab === 'active' ? (
                <button
                  onClick={() => triggerConfirm('bulk_delete')}
                  className="bg-red-100 text-red-700 hover:bg-red-200 px-3 py-1.5 rounded-md text-sm font-medium flex items-center gap-2 transition-colors"
                >
                  <i className="ph ph-trash" /> Hapus Massal
                </button>
              ) : (
                <button
                  onClick={() => triggerConfirm('bulk_restore')}
                  className="bg-emerald-100 text-emerald-700 hover:bg-emerald-200 px-3 py-1.5 rounded-md text-sm font-medium flex items-center gap-2 transition-colors"
                >
                  <i className="ph ph-arrow-counter-clockwise" /> Pulihkan Massal
                </button>
              )}
            </div>
          </div>
        )}

        {/* Table Card */}
        <div className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
          <div className="overflow-x-auto">
            {loading ? (
              <div className="py-12 text-center text-gray-500">Memuat data admin...</div>
            ) : error ? (
              <div className="py-12 text-center text-red-500">{error}</div>
            ) : users.length === 0 ? (
              <div className="py-12 text-center text-gray-500">
                <i className="ph ph-users text-4xl text-gray-300 mb-2 block mx-auto" />
                Tidak ada data untuk ditampilkan.
              </div>
            ) : (
              <table className="w-full text-left border-collapse">
                <thead>
                  <tr className="bg-gray-50 border-b border-gray-200 text-xs uppercase tracking-wider text-gray-500 font-semibold">
                    <th className="p-4 w-12 text-center">
                      <input
                        type="checkbox"
                        onChange={handleSelectAll}
                        checked={isAllSelected}
                        disabled={selectableUsers.length === 0}
                        className="rounded border-gray-300 text-brand-600 focus:ring-brand-500"
                      />
                    </th>
                    <th className="p-4">Pengguna</th>
                    <th className="p-4">Role</th>
                    <th className="p-4">Status Email</th>
                    <th className="p-4 text-right">Aksi</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-100 text-sm text-gray-700">
                  {users.map((item) => {
                    const isProtected = item.role === 'Super Admin' || item.role_id === 1
                    const isPending = item.status === 'pending_activation'
                    const canSelect = !isProtected && !isPending

                    return (
                      <tr key={item.id} className="hover:bg-gray-50 transition-colors group">
                        <td className="p-4 text-center">
                          <input
                            type="checkbox"
                            checked={selectedIds.includes(item.id)}
                            onChange={() => handleSelectItem(item.id)}
                            disabled={!canSelect}
                            className="rounded border-gray-300 text-brand-600 focus:ring-brand-500 disabled:opacity-30"
                          />
                        </td>
                        <td className="p-4">
                          <div className="flex items-center gap-3">
                            <div className="w-10 h-10 rounded-full bg-brand-100 text-brand-600 flex items-center justify-center font-bold text-sm shrink-0">
                              {(item.name || 'A').charAt(0).toUpperCase()}
                            </div>
                            <div>
                              <p className="font-medium text-gray-900">{item.name}</p>
                              <p className="text-xs text-gray-500">{item.email}</p>
                            </div>
                          </div>
                        </td>
                        <td className="p-4">
                          <span
                            className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-md text-xs font-medium border ${
                              isProtected
                                ? 'bg-purple-50 text-purple-700 border-purple-200'
                                : 'bg-blue-50 text-blue-700 border-blue-200'
                            }`}
                          >
                            <i className={`ph ${isProtected ? 'ph-crown' : 'ph-shield-check'}`} />
                            <span>{item.role || 'Admin'}</span>
                          </span>
                        </td>
                        <td className="p-4">
                          <span
                            className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-md text-xs font-medium border ${
                              item.status === 'active'
                                ? 'bg-emerald-50 text-emerald-700 border-emerald-200'
                                : item.status === 'pending_activation'
                                ? 'bg-amber-50 text-amber-700 border-amber-200'
                                : 'bg-gray-100 text-gray-600 border-gray-200'
                            }`}
                          >
                            <i
                              className={`ph ${
                                item.status === 'active'
                                  ? 'ph-check-circle'
                                  : item.status === 'pending_activation'
                                  ? 'ph-clock'
                                  : 'ph-prohibit'
                              }`}
                            />
                            <span>
                              {item.status === 'active'
                                ? 'Aktif'
                                : item.status === 'pending_activation'
                                ? 'Menunggu Aktivasi'
                                : 'Nonaktif'}
                            </span>
                          </span>
                        </td>
                        <td className="p-4 text-right">
                          <div className="flex items-center justify-end gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                            {!isProtected ? (
                              <>
                                {currentTab === 'active' && (
                                  <div className="flex gap-2">
                                    <button
                                      onClick={() => triggerConfirm('deactivate', item.id)}
                                      className="p-1.5 text-gray-500 hover:text-amber-600 hover:bg-amber-50 rounded"
                                      title="Nonaktifkan Akun"
                                    >
                                      <i className="ph ph-prohibit text-lg" />
                                    </button>
                                    <button
                                      onClick={() => triggerConfirm('delete', item.id)}
                                      className="p-1.5 text-gray-500 hover:text-red-600 hover:bg-red-50 rounded"
                                      title="Hapus Akun"
                                    >
                                      <i className="ph ph-trash text-lg" />
                                    </button>
                                  </div>
                                )}
                                {currentTab === 'pending' && (
                                  <div className="flex gap-2">
                                    <button
                                      onClick={() => triggerConfirm('resend', item.id)}
                                      className="p-1.5 text-gray-500 hover:text-brand-600 hover:bg-brand-50 rounded"
                                      title="Kirim Ulang Undangan Aktivasi"
                                    >
                                      <i className="ph ph-paper-plane-right text-lg" />
                                    </button>
                                    <button
                                      onClick={() => triggerConfirm('delete', item.id)}
                                      className="p-1.5 text-gray-500 hover:text-red-600 hover:bg-red-50 rounded"
                                      title="Hapus Akun"
                                    >
                                      <i className="ph ph-trash text-lg" />
                                    </button>
                                  </div>
                                )}
                                {currentTab === 'trash' && (
                                  <div className="flex gap-2">
                                    <button
                                      onClick={() => triggerConfirm('restore', item.id)}
                                      className="p-1.5 text-gray-500 hover:text-emerald-600 hover:bg-emerald-50 rounded"
                                      title="Pulihkan Akun"
                                    >
                                      <i className="ph ph-arrow-counter-clockwise text-lg" />
                                    </button>
                                  </div>
                                )}
                              </>
                            ) : (
                              <span className="text-xs text-gray-400 italic">Protected</span>
                            )}
                          </div>
                        </td>
                      </tr>
                    )
                  })}
                </tbody>
              </table>
            )}
          </div>

          {/* Pagination */}
          {!loading && !error && users.length > 0 && (
            <div className="bg-gray-50 border-t border-gray-200 px-4 py-3 flex items-center justify-between sm:px-6 rounded-b-xl">
              <div className="text-sm text-gray-500">
                Menampilkan <span className="font-medium text-gray-900">{(page - 1) * limit + 1}</span> sampai <span className="font-medium text-gray-900">{Math.min(page * limit, total)}</span> dari <span className="font-medium text-gray-900">{total}</span> hasil
              </div>
              <div className="flex gap-1">
                <button
                  onClick={() => setPage(p => Math.max(p - 1, 1))}
                  disabled={page === 1}
                  className="px-3 py-1 border border-gray-200 bg-white text-gray-500 rounded hover:bg-gray-50 disabled:opacity-50 transition"
                >
                  Prev
                </button>
                
                {renderPagination()}

                <button
                  onClick={() => setPage(p => Math.min(p + 1, totalPages))}
                  disabled={page === totalPages}
                  className="px-3 py-1 border border-gray-200 bg-white text-gray-500 rounded hover:bg-gray-50 disabled:opacity-50 transition"
                >
                  Next
                </button>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* FORM MODAL (Create) */}
      {isFormOpen && (
        <div className="fixed inset-0 z-50 overflow-y-auto" role="dialog" aria-modal="true">
          <div className="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
            <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" onClick={() => setIsFormOpen(false)} />
            <span className="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>
            
            <div className="inline-block align-bottom bg-white rounded-xl text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-md sm:w-full">
              <div className="bg-white">
                <div className="border-b border-gray-200 px-6 py-4 flex items-center justify-between">
                  <h3 className="text-lg leading-6 font-heading font-semibold text-gray-900">Tambah Admin Baru</h3>
                  <button onClick={() => setIsFormOpen(false)} className="text-gray-400 hover:text-gray-500">
                    <i className="ph ph-x text-xl" />
                  </button>
                </div>
                <form onSubmit={handleFormSubmit} id="createAdminForm" className="px-6 py-4 space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Nama Lengkap *</label>
                    <input
                      type="text"
                      value={formData.name}
                      onChange={(e) => setFormData(prev => ({ ...prev, name: e.target.value }))}
                      required
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-brand-500 focus:border-brand-500 text-sm"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Email Aktif *</label>
                    <input
                      type="email"
                      value={formData.email}
                      onChange={(e) => setFormData(prev => ({ ...prev, email: e.target.value }))}
                      required
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-brand-500 focus:border-brand-500 text-sm"
                    />
                    <p className="text-xs text-gray-500 mt-1">Link verifikasi dan password sementara akan dikirimkan ke email ini.</p>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Role Akses *</label>
                    <select
                      value={formData.role_id}
                      onChange={(e) => setFormData(prev => ({ ...prev, role_id: e.target.value }))}
                      required
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-brand-500 focus:border-brand-500 text-sm bg-white"
                    >
                      <option value="2">Admin</option>
                      <option value="3">Admin Berita</option>
                      <option value="4">Admin Kegiatan</option>
                    </select>
                    <p className="text-xs text-gray-500 mt-1">Super Admin tidak bisa dibuat via undangan (role_id 1 khusus seed).</p>
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
                  type="submit"
                  form="createAdminForm"
                  className="px-4 py-2 bg-brand-600 border border-transparent rounded-lg text-sm font-medium text-white hover:bg-brand-700 flex items-center gap-2 transition"
                >
                  <i className="ph ph-paper-plane-right" /> Kirim Undangan
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* CONFIRMATION MODAL */}
      {confirm.isOpen && (
        <div className="fixed inset-0 z-50 overflow-y-auto" role="dialog" aria-modal="true">
          <div className="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
            <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" onClick={() => setConfirm(prev => ({ ...prev, isOpen: false }))} />
            <span className="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>
            
            <div className="inline-block align-bottom bg-white rounded-xl text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-sm sm:w-full">
              <div className="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
                <div className="sm:flex sm:items-start">
                  <div className={`mx-auto flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-full sm:mx-0 sm:h-10 sm:w-10 ${
                    confirm.type.includes('delete') || confirm.type === 'deactivate' ? 'bg-red-100 text-red-600' : 'bg-emerald-100 text-emerald-600'
                  }`}>
                    <i className={`text-2xl ph ${confirm.icon}`} />
                  </div>
                  <div className="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left">
                    <h3 className="text-lg leading-6 font-medium text-gray-900">{confirm.title}</h3>
                    <div className="mt-2">
                      <p className="text-sm text-gray-500">{confirm.message}</p>
                    </div>
                  </div>
                </div>
              </div>
              
              <div className="bg-gray-50 px-4 py-3 sm:px-6 flex justify-end gap-2">
                <button
                  type="button"
                  onClick={() => setConfirm(prev => ({ ...prev, isOpen: false }))}
                  className="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none sm:mt-0 sm:w-auto sm:text-sm"
                >
                  Batal
                </button>
                <button
                  type="button"
                  onClick={executeAction}
                  className={`w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 text-base font-medium text-white focus:outline-none sm:w-auto sm:text-sm ${
                    confirm.type.includes('delete') || confirm.type === 'deactivate' ? 'bg-red-600 hover:bg-red-700' : 'bg-emerald-600 hover:bg-emerald-700'
                  }`}
                >
                  Ya, Lanjutkan
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Toast Notification */}
      {toast.show && (
        <div className={`fixed bottom-4 right-4 z-50 flex items-center p-4 rounded-lg shadow-lg text-white transition-opacity duration-300 ${
          toast.type === 'success' ? 'bg-emerald-500' : 'bg-red-500'
        }`}>
          <i className={`text-xl mr-2 ph ${toast.type === 'success' ? 'ph-check-circle' : 'ph-warning-circle'}`} />
          <span className="text-sm font-medium">{toast.message}</span>
        </div>
      )}
    </AdminLayout>
  )
}
