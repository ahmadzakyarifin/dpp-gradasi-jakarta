import { useState, useEffect, useCallback } from 'react'
import AdminLayout from '../../layouts/AdminLayout'
import { userService } from '../../services/userService'
import { useFormErrors, useRateLimitCooldown } from '../../utils/parseApiError'

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

  // Modal Form (Create) — hanya 2 role: super_admin (1) & admin (2)
  const [isFormOpen, setIsFormOpen] = useState(false)
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    role_id: '2' // default: Admin
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
    type: 'success'
  })
  const [formErrors, setFormErrors] = useState({})
  const [touched, setTouched] = useState({})

  // Reset Password State
  const [isResetPwdOpen, setIsResetPwdOpen] = useState(false)
  const [resetUserId, setResetUserId] = useState(null)
  const [resetUserName, setResetUserName] = useState('')
  const [newPassword, setNewPassword] = useState('')
  const [resetLoading, setResetLoading] = useState(false)
  const [resetErrors, setResetErrors] = useState({})

  const validateForm = (data = formData) => {
    const errors = {}
    if (!data.name || !data.name.trim()) {
      errors.name = 'Nama lengkap wajib diisi.'
    }
    if (!data.email || !data.email.trim()) {
      errors.email = 'Email wajib diisi.'
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(data.email)) {
      errors.email = 'Format email tidak valid.'
    }
    if (!data.role_id) {
      errors.role_id = 'Role akses wajib dipilih.'
    }
    return errors
  }

  // Error backend: pesan error dari helper + countdown rate limit
  const { fieldErrors, applyError, clearFieldError, resetFieldErrors } = useFormErrors()
  const { cooldown, isLimited, applyRateLimit } = useRateLimitCooldown()

  const showToast = (message, type = 'success') => {
    setToast({ show: true, message, type })
    setTimeout(() => {
      setToast(prev => ({ ...prev, show: false }))
    }, 3000)
  }

  const fetchUsers = useCallback(() => {
    setLoading(true)
    const params = {
      search: search.trim(),
      page,
      limit
    }

    if (currentTab === 'active') {
      params.status = 'active'
      params.trashed = false
    } else if (currentTab === 'pending') {
      params.status = 'inactive'
      params.trashed = false
    } else if (currentTab === 'trash') {
      params.status = ''
      params.trashed = true
    }

    userService.list(params)
      .then(res => {
        if (res.success && res.data) {
          setUsers(res.data.items || res.data.users || [])
          const meta = res.data.meta || {}
          if (meta.total !== undefined) setTotal(meta.total)
          if (meta.total_pages !== undefined) setTotalPages(meta.total_pages)
          setError(null)
        } else {
          setError('Gagal memuat data admin')
        }
      })
      .catch((err) => {
        setError(err.message || 'Kesalahan koneksi ke server')
      })
      .finally(() => setLoading(false))
  }, [currentTab, search, page, limit])

  useEffect(() => {
    setSelectedIds([])
    setPage(1)
  }, [currentTab])

  // Muat ulang data setiap page/currentTab/search berubah (pagination & filter)
  useEffect(() => {
    fetchUsers()
  }, [fetchUsers])

  const handleSearchSubmit = (e) => {
    e.preventDefault()
    setPage(1)
  }

  const handleReset = () => {
    setSearch('')
    setFilterRole('')
    setPage(1)
  }

  // Row Selection logic
  const selectableUsers = users.filter(
    item => !item.is_system && !(item.status === 'inactive' && !item.has_password)
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
    const errors = validateForm()
    if (Object.keys(errors).length > 0) {
      setFormErrors(errors)
      setTouched(Object.keys(errors).reduce((acc, k) => ({ ...acc, [k]: true }), {}))
      return
    }
    setFormErrors({})
    resetFieldErrors()

    try {
      const payload = {
        name: formData.name,
        email: formData.email,
        role: formData.role_id === '1' ? 'super_admin' : 'admin'
      }
      await userService.create(payload)
      setIsFormOpen(false)
      setFormData({ name: '', email: '', role_id: '2' })
      setTouched({})
      setFormErrors({})
      showToast('Admin berhasil dibuat. Kredensial login dikirim ke email admin!', 'success')
      fetchUsers()
    } catch (err) {
      const parsed = applyError(err)
      applyRateLimit(err)
      setFormErrors(prev => ({ ...prev, ...parsed.fieldErrors }))
      setTouched(prev => ({ ...prev, ...Object.keys(parsed.fieldErrors).reduce((acc, k) => ({ ...acc, [k]: true }), {}) }))
      if (Object.keys(parsed.fieldErrors).length === 0) {
        showToast(parsed.message || 'Gagal mengirim undangan admin baru.', 'error')
      }
    }
  }

  const handleResetPasswordSubmit = async (e) => {
    e.preventDefault()
    if (!newPassword || newPassword.trim().length < 6) {
      setResetErrors({ password: 'Password baru minimal 6 karakter.' })
      return
    }
    setResetErrors({})
    setResetLoading(true)
    try {
      await userService.resetPassword(resetUserId, newPassword.trim())
      setIsResetPwdOpen(false)
      showToast(`Password untuk ${resetUserName} berhasil direset!`, 'success')
      fetchUsers()
    } catch (err) {
      const parsed = applyError(err)
      setResetErrors(parsed.fieldErrors || {})
      if (!parsed.fieldErrors || Object.keys(parsed.fieldErrors).length === 0) {
        showToast(parsed.message || 'Gagal mereset password.', 'error')
      }
    } finally {
      setResetLoading(false)
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
    } else if (type === 'activate') {
      title = 'Aktifkan Admin'
      message = 'Akun admin ini akan diaktifkan kembali. Lanjutkan?'
      icon = 'ph-check-circle'
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
      } else if (type === 'activate') {
        await userService.setStatus(id, 'active')
        showToast('Akun admin berhasil diaktifkan.', 'success')
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

  const headerContent = (
    <div className="flex items-center gap-2 w-full max-w-3xl animate-fade-in-up">
      <div className="relative w-full">
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
        onClick={handleReset}
        className="shrink-0 bg-gray-50 border border-gray-200 text-gray-700 hover:bg-gray-100 px-3 py-2 rounded-lg text-sm font-medium flex items-center gap-2 transition-all btn-press"
      >
        <i className="ph ph-arrows-counter-clockwise text-lg" /> Reset
      </button>
    </div>
  )

  return (
    <AdminLayout title="Manajemen Admin" headerContent={headerContent}>
      <div className="max-w-7xl mx-auto space-y-6 animate-fade-in-up">
        {/* Header Actions & Filter */}
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
          {/* Tabs */}
          <div className="flex bg-white rounded-lg p-1 border border-gray-200 shadow-sm flex-wrap sm:flex-nowrap">
            <button
              onClick={() => setCurrentTab('active')}
              className={`px-4 py-1.5 rounded-md text-sm transition-all duration-200 ${
                currentTab === 'active'
                  ? 'bg-brand-50 text-brand-600 shadow-sm font-medium'
                  : 'text-gray-500 hover:text-gray-700'
              }`}
            >
              Admin Aktif
            </button>
            <button
              onClick={() => setCurrentTab('pending')}
              className={`px-4 py-1.5 rounded-md text-sm transition-all duration-200 flex items-center gap-2 ${
                currentTab === 'pending'
                  ? 'bg-amber-50 text-amber-600 shadow-sm font-medium'
                  : 'text-gray-500 hover:text-gray-700'
              }`}
            >
              <i className="ph ph-clock" /> Menunggu Aktivasi
            </button>
            <button
              onClick={() => setCurrentTab('trash')}
              className={`px-4 py-1.5 rounded-md text-sm transition-all duration-200 flex items-center gap-2 ${
                currentTab === 'trash'
                  ? 'bg-red-50 text-red-600 shadow-sm font-medium'
                  : 'text-gray-500 hover:text-gray-700'
              }`}
            >
              <i className="ph ph-trash" /> Akun Terhapus
            </button>
          </div>

          {currentTab === 'active' && (
            <button
              onClick={() => {
                setFormData({ name: '', email: '', role_id: '2' })
                setFormErrors({})
                setTouched({})
                resetFieldErrors()
                setIsFormOpen(true)
              }}
              className="shrink-0 bg-brand-600 hover:bg-brand-700 text-white px-4 py-2 rounded-lg text-sm font-medium flex items-center gap-2 transition-all btn-press shadow-sm"
            >
              <i className="ph ph-user-plus text-lg" /> Tambah
            </button>
          )}
        </div>

        {/* Bulk Actions Bar */}
        {selectedIds.length > 0 && (
          <div className="flex flex-wrap items-center gap-2 bg-brand-50/60 border border-brand-100 rounded-xl px-4 py-2.5 shadow-sm animate-fade-in">
            <span className="text-sm font-semibold text-brand-700">{selectedIds.length} akun terpilih</span>
            <div className="flex gap-2 ml-auto">
              {currentTab === 'active' ? (
                <button
                  onClick={() => triggerConfirm('bulk_delete')}
                  className="bg-red-600 hover:bg-red-700 text-white px-3.5 py-1.5 rounded-lg text-xs font-semibold flex items-center gap-1.5 transition-colors"
                >
                  <i className="ph ph-trash" /> Hapus Massal
                </button>
              ) : (
                <button
                  onClick={() => triggerConfirm('bulk_restore')}
                  className="bg-emerald-600 hover:bg-emerald-700 text-white px-3.5 py-1.5 rounded-lg text-xs font-semibold flex items-center gap-1.5 transition-colors"
                >
                  <i className="ph ph-arrow-counter-clockwise" /> Pulihkan Massal
                </button>
              )}
              <button
                onClick={() => setSelectedIds([])}
                className="bg-white border border-slate-200 text-slate-600 hover:bg-slate-50 px-3.5 py-1.5 rounded-lg text-xs font-semibold transition-colors"
              >
                Batal
              </button>
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
              <div className="py-16 text-center text-slate-500 flex flex-col items-center justify-center">
                <i className="ph ph-users text-gray-300 text-5xl mb-4" />
                <p className="font-medium text-gray-500">Tidak ada data admin untuk ditampilkan.</p>
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
                    const isProtected = item.is_system || item.role === 'super_admin'
                    const isPending = item.status === 'inactive' && !item.has_password
                    const canSelect = !isProtected && !isPending

                    return (
                      <tr key={item.id} className="hover:bg-gray-50 transition-colors group admin-row">
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
                            <span>{item.role === 'super_admin' ? 'Super Admin' : (item.role || 'Admin')}</span>
                          </span>
                          {item.must_change_password && !isProtected && (
                            <span className="inline-flex items-center gap-1 px-2 py-0.5 ml-2 rounded bg-amber-50 text-amber-700 border border-amber-200 text-[11px] font-medium">
                              <i className="ph ph-lock-key" /> Ganti pwd
                            </span>
                          )}
                        </td>
                        <td className="p-4">
                          {isPending ? (
                            <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-md text-xs font-medium border bg-amber-50 text-amber-700 border-amber-200">
                              <i className="ph ph-clock" />
                              <span>Menunggu Aktivasi</span>
                            </span>
                          ) : (
                            <button
                              disabled={isProtected}
                              onClick={() => triggerConfirm(item.status === 'active' ? 'deactivate' : 'activate', item.id)}
                              className={`inline-flex items-center gap-2 ${isProtected ? 'cursor-not-allowed opacity-80' : 'cursor-pointer'}`}
                            >
                              <div className={`relative w-9 h-5 rounded-full transition-colors ${item.status === 'active' ? 'bg-brand-600' : 'bg-slate-200'}`}>
                                <div className={`absolute top-[2px] left-[2px] w-4 h-4 bg-white rounded-full transition-transform ${item.status === 'active' ? 'translate-x-4' : 'translate-x-0'}`} />
                              </div>
                              <span className={`text-xs font-semibold ${item.status === 'active' ? 'text-brand-600' : 'text-slate-400'}`}>
                                {item.status === 'active' ? 'Aktif' : 'Nonaktif'}
                              </span>
                            </button>
                          )}
                        </td>
                        <td className="p-4 text-right">
                          <div className="flex items-center justify-end gap-2">
                            {!isProtected ? (
                              <>
                                {currentTab === 'active' && (
                                  <div className="flex gap-2">
                                    <button
                                      onClick={() => {
                                        setResetUserId(item.id)
                                        setResetUserName(item.name)
                                        setNewPassword('')
                                        setResetErrors({})
                                        setIsResetPwdOpen(true)
                                      }}
                                      className="p-1.5 text-gray-500 hover:text-indigo-600 hover:bg-indigo-50 rounded"
                                      title="Reset Password Admin"
                                    >
                                      <i className="ph ph-key text-lg" />
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
          </div>          {/* Pagination */}
          {!loading && !error && users.length > 0 && (
            <div className="bg-white border-t border-slate-200 px-4 py-3 flex items-center justify-between sm:px-6 rounded-b-xl">
              <span className="text-xs text-slate-500">
                Hal {page} dari {totalPages} · {total} data
              </span>
              <div className="flex items-center gap-1.5">
                <button
                  type="button"
                  disabled={page <= 1}
                  onClick={() => setPage(p => Math.max(p - 1, 1))}
                  className="w-8 h-8 flex items-center justify-center rounded-lg border border-slate-200 text-slate-500 hover:border-brand-500 hover:text-brand-600 disabled:opacity-40 transition"
                >
                  <i className="ph-bold ph-caret-left text-sm" />
                </button>
                {getVisiblePageNumbers(page, totalPages, 5).map((n) => (
                  <button
                    key={n}
                    type="button"
                    onClick={() => setPage(n)}
                    className={`w-8 h-8 flex items-center justify-center rounded-lg text-sm font-semibold transition ${n === page ? 'bg-brand-600 text-white' : 'border border-slate-200 text-slate-500 hover:border-brand-500 hover:text-brand-600'}`}
                  >
                    {n}
                  </button>
                ))}
                <button
                  type="button"
                  disabled={page >= totalPages}
                  onClick={() => setPage(p => Math.min(p + 1, totalPages))}
                  className="w-8 h-8 flex items-center justify-center rounded-lg border border-slate-200 text-slate-500 hover:border-brand-500 hover:text-brand-600 disabled:opacity-40 transition"
                >
                  <i className="ph-bold ph-caret-right text-sm" />
                </button>
              </div>
            </div>
          )}
        </div>
      </div>      {/* FORM MODAL (Create) */}
      {isFormOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
          <div className="fixed inset-0 bg-black/40 backdrop-blur-[2px]" onClick={() => setIsFormOpen(false)} />
          <div className="relative bg-white rounded-2xl shadow-2xl max-w-lg w-full max-h-[90vh] flex flex-col overflow-hidden z-10 animate-scale-up">
            <div className="border-b border-slate-200 px-6 py-4 flex items-center justify-between">
              <h3 className="font-heading font-bold text-slate-900 text-lg">Tambah Admin Baru</h3>
              <button onClick={() => setIsFormOpen(false)} className="text-slate-400 hover:text-slate-600 p-1">
                <i className="ph-bold ph-x text-lg" />
              </button>
            </div>
            <form onSubmit={handleFormSubmit} noValidate id="createAdminForm" className="p-6 overflow-y-auto max-h-[calc(90vh-120px)] space-y-4">
              <div>
                <label className="block text-xs font-semibold text-gray-500 mb-1">Nama Lengkap <span className="text-red-500">*</span></label>
                <input
                  type="text"
                  value={formData.name}
                  onChange={(e) => {
                    setFormData(prev => ({ ...prev, name: e.target.value }))
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
                  className={`w-full px-3 py-2 border rounded-xl focus:ring-brand-500 focus:border-brand-500 text-sm outline-none transition-colors ${touched.name && formErrors.name ? 'border-red-400 focus:ring-2 focus:ring-red-100' : 'border-gray-300'}`}
                />
                {touched.name && formErrors.name && (
                  <p className="text-red-500 text-[11px] font-semibold mt-1 flex items-center gap-1">
                    <i className="ph-bold ph-warning-circle text-xs" /> {formErrors.name}
                  </p>
                )}
              </div>
              <div>
                <label className="block text-xs font-semibold text-gray-500 mb-1">Email <span className="text-red-500">*</span></label>
                <input
                  type="email"
                  value={formData.email}
                  onChange={(e) => {
                    setFormData(prev => ({ ...prev, email: e.target.value }))
                    clearFieldError('email')
                    if (touched.email) {
                      const errs = validateForm({ ...formData, email: e.target.value })
                      setFormErrors(prev => ({ ...prev, email: errs.email }))
                    }
                  }}
                  onBlur={() => {
                    setTouched(prev => ({ ...prev, email: true }))
                    const errs = validateForm()
                    setFormErrors(prev => ({ ...prev, email: errs.email }))
                  }}
                  className={`w-full px-3 py-2 border rounded-xl focus:ring-brand-500 focus:border-brand-500 text-sm outline-none transition-colors ${touched.email && formErrors.email ? 'border-red-400 focus:ring-2 focus:ring-red-100' : 'border-gray-300'}`}
                />
                {touched.email && formErrors.email && (
                  <p className="text-red-500 text-[11px] font-semibold mt-1 flex items-center gap-1">
                    <i className="ph-bold ph-warning-circle text-xs" /> {formErrors.email}
                  </p>
                )}
                <p className="text-[10px] text-gray-400 mt-1">Kredensial login (email & password default) akan dikirim ke email ini.</p>
              </div>
              <div>
                <label className="block text-xs font-semibold text-gray-500 mb-1">Role Akses <span className="text-red-500">*</span></label>
                <select
                  value={formData.role_id}
                  onChange={(e) => {
                    const val = e.target.value
                    setFormData(prev => ({ ...prev, role_id: val }))
                    clearFieldError('role_id')
                    if (touched.role_id) {
                      const errs = validateForm({ ...formData, role_id: val })
                      setFormErrors(prev => ({ ...prev, role_id: errs.role_id }))
                    }
                  }}
                  onBlur={() => {
                    setTouched(prev => ({ ...prev, role_id: true }))
                    const errs = validateForm()
                    setFormErrors(prev => ({ ...prev, role_id: errs.role_id }))
                  }}
                  className={`w-full px-3 py-2 border rounded-xl focus:ring-brand-500 focus:border-brand-500 text-sm bg-white outline-none transition-colors ${touched.role_id && formErrors.role_id ? 'border-red-400 focus:ring-2 focus:ring-red-100' : 'border-gray-300'}`}
                >
                  <option value="1">Super Admin</option>
                  <option value="2">Admin</option>
                </select>
                {touched.role_id && formErrors.role_id && (
                  <p className="text-red-500 text-[11px] font-semibold mt-1 flex items-center gap-1">
                    <i className="ph-bold ph-warning-circle text-xs" /> {formErrors.role_id}
                  </p>
                )}
                <p className="text-[10px] text-gray-400 mt-1">Super Admin punya akses penuh (termasuk manajemen admin & log). Admin mengelola konten.</p>
              </div>
              <div className="flex justify-end gap-2 pt-4 border-t items-center mt-4">
                <button
                  type="button"
                  onClick={() => setIsFormOpen(false)}
                  className="px-4 py-2 border rounded-xl text-sm font-semibold hover:bg-slate-50 transition-colors"
                >
                  Batal
                </button>
                <button
                  type="submit"
                  className="px-5 py-2 bg-brand-600 text-white rounded-xl text-sm font-semibold hover:bg-brand-700 transition-colors flex items-center gap-2"
                >
                  <i className="ph ph-paper-plane-right" /> Buat & Kirim Kredensial
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* MODAL RESET PASSWORD */}
      {isResetPwdOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
          <div className="fixed inset-0 bg-black/40 backdrop-blur-[2px]" onClick={() => setIsResetPwdOpen(false)} />
          <div className="relative bg-white rounded-2xl shadow-2xl max-w-md w-full flex flex-col overflow-hidden z-10 animate-scale-up">
            <div className="border-b border-slate-200 px-6 py-4 flex items-center justify-between">
              <h3 className="font-heading font-bold text-slate-900 text-lg">Reset Password Admin</h3>
              <button onClick={() => setIsResetPwdOpen(false)} className="text-slate-400 hover:text-slate-600 p-1">
                <i className="ph-bold ph-x text-lg" />
              </button>
            </div>
            <form onSubmit={handleResetPasswordSubmit} noValidate className="p-6 space-y-4">
              <p className="text-xs text-gray-500 leading-relaxed">
                Anda akan mereset password untuk admin: <strong>{resetUserName}</strong>. Admin tersebut akan otomatis ter-logout dari semua perangkat demi keamanan.
              </p>
              <div>
                <label className="block text-xs font-semibold text-gray-500 mb-1">Password Baru <span className="text-red-500">*</span></label>
                <input
                  type="password"
                  placeholder="Minimal 6 karakter..."
                  value={newPassword}
                  onChange={(e) => {
                    setNewPassword(e.target.value)
                    if (resetErrors.password) setResetErrors({})
                  }}
                  className={`w-full px-3 py-2 border rounded-xl focus:ring-brand-500 focus:border-brand-500 text-sm outline-none transition-colors ${resetErrors.password ? 'border-red-400 focus:ring-2 focus:ring-red-100' : 'border-gray-300'}`}
                />
                {resetErrors.password && (
                  <p className="text-red-500 text-[11px] font-semibold mt-1 flex items-center gap-1">
                    <i className="ph-bold ph-warning-circle text-xs" /> {resetErrors.password}
                  </p>
                )}
              </div>
              <div className="flex justify-end gap-2 pt-4 border-t items-center mt-4">
                <button
                  type="button"
                  onClick={() => setIsResetPwdOpen(false)}
                  className="px-4 py-2 border rounded-xl text-sm font-semibold hover:bg-slate-50 transition-colors"
                >
                  Batal
                </button>
                <button
                  type="submit"
                  disabled={resetLoading}
                  className="px-5 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-xl text-sm font-semibold transition-colors flex items-center gap-2"
                >
                  <i className="ph ph-key" /> {resetLoading ? 'Menyimpan...' : 'Reset Password'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* CONFIRMATION MODAL */}
      {confirm.isOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
          <div className="fixed inset-0 bg-black/40 backdrop-blur-[2px]" onClick={() => setConfirm(prev => ({ ...prev, isOpen: false }))} />
          <div className="relative bg-white rounded-2xl shadow-2xl max-w-md w-full overflow-hidden z-10 p-6 animate-scale-up">
            <div className="sm:flex sm:items-start">
              <div className={`mx-auto flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-full sm:mx-0 sm:h-10 sm:w-10 ${
                confirm.type.includes('delete') || confirm.type === 'deactivate' ? 'bg-red-100 text-red-600' : 'bg-emerald-100 text-emerald-600'
              }`}>
                <i className={`text-2xl ph ${confirm.icon}`} />
              </div>
              <div className="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left">
                <h3 className="text-lg leading-6 font-bold text-slate-900 font-heading">{confirm.title}</h3>
                <div className="mt-2">
                  <p className="text-sm text-gray-500">{confirm.message}</p>
                </div>
              </div>
            </div>
            <div className="mt-6 flex justify-end gap-2 border-t pt-4">
              <button
                type="button"
                onClick={() => setConfirm(prev => ({ ...prev, isOpen: false }))}
                className="px-4 py-2 border rounded-xl text-sm font-semibold hover:bg-slate-50 transition-colors"
              >
                Batal
              </button>
              <button
                type="button"
                onClick={executeAction}
                className={`px-5 py-2 text-white rounded-xl text-sm font-semibold transition-colors ${
                  confirm.type.includes('delete') || confirm.type === 'deactivate' ? 'bg-red-600 hover:bg-red-700' : 'bg-emerald-600 hover:bg-emerald-700'
                }`}
              >
                Ya, Lanjutkan
              </button>
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
