import { useState, useEffect, useCallback } from 'react'
import AdminLayout from '../../layouts/AdminLayout'
import { kontakService } from '../../services/kontakService'
import { useFormErrors } from '../../utils/parseApiError'

export default function KontakAdmin() {
  const [messages, setMessages] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  const [currentTab, setCurrentTab] = useState('active') // active, trash
  const [search, setSearch] = useState('')
  const [filterStatus, setFilterStatus] = useState('all') // all, read, unread

  // Pagination
  const [page, setPage] = useState(1)
  const [limit] = useState(10)
  const [total, setTotal] = useState(0)
  const [totalPages, setTotalPages] = useState(1)

  // Selection
  const [selectedIds, setSelectedIds] = useState([])

  // Modal Detail
  const [viewData, setViewData] = useState(null)
  const [isViewOpen, setIsViewOpen] = useState(false)

  // Confirm Dialog
  const [confirm, setConfirm] = useState({
    isOpen: false,
    type: '', // delete, restore, bulk_delete, bulk_restore
    id: null,
    title: '',
    message: '',
    icon: ''
  })

  // Toast
  const [toast, setToast] = useState({
    show: false,
    message: '',
    type: 'success'
  })
  // Error backend: field errors inline
  const { applyError } = useFormErrors()

  const showToast = (message, type = 'success') => {
    setToast({ show: true, message, type })
    setTimeout(() => {
      setToast(prev => ({ ...prev, show: false }))
    }, 3000)
  }

  const fetchMessages = useCallback(() => {
    setLoading(true)
    kontakService.list({
      tab: currentTab,
      search,
      status: currentTab === 'active' ? filterStatus : undefined,
      page,
      limit
    })
      .then(res => {
        if (res.success && res.data) {
          setMessages(res.data.items || [])
          if (res.data.pagination) {
            setTotal(res.data.pagination.total || 0)
            setTotalPages(res.data.pagination.totalPages || 1)
          }
          setError(null)
        } else {
          setError('Gagal memuat pesan masuk')
        }
      })
      .catch((err) => {
        setError(err.message || 'Kesalahan koneksi ke server')
      })
      .finally(() => setLoading(false))
  }, [currentTab, search, filterStatus, page, limit])

  useEffect(() => {
    setSelectedIds([])
    setPage(1)
  }, [currentTab, filterStatus])

  // Muat ulang data setiap page/search/filter berubah (pagination & filter)
  useEffect(() => {
    fetchMessages()
  }, [fetchMessages])

  // Handles manual search submit
  const handleSearchSubmit = (e) => {
    e.preventDefault()
    setPage(1)
    // fetchMessages() otomatis dipanggil via useEffect [fetchMessages] saat search/page berubah
  }

  const handleReset = () => {
    setSearch('')
    setFilterStatus('all')
    setPage(1)
  }

  // Row Selection logic
  const handleSelectAll = (e) => {
    if (e.target.checked) {
      setSelectedIds(messages.map(m => m.id))
    } else {
      setSelectedIds([])
    }
  }

  const handleSelectItem = (id) => {
    setSelectedIds(prev => 
      prev.includes(id) ? prev.filter(item => item !== id) : [...prev, id]
    )
  }

  const isAllSelected = messages.length > 0 && selectedIds.length === messages.length

  // View modal detail (marks message as read if active tab)
  const openViewModal = (item) => {
    if (currentTab === 'trash') return
    setViewData(item)
    setIsViewOpen(true)
    
    // Call detail API to mark it read on the backend (or detail fetches read state update)
    kontakService.getById(item.id)
      .then(res => {
        if (res.success) {
          // Update local status of is_read
          setMessages(prev => prev.map(m => m.id === item.id ? { ...m, is_read: true } : m))
        }
      })
      .catch(() => {})
  }

  // Trigger Confirm Modal
  const triggerConfirm = (type, id = null) => {
    let title = ''
    let message = ''
    let icon = ''

    if (type === 'delete') {
      title = 'Hapus Pesan'
      message = 'Pesan ini akan dipindahkan ke Sampah.'
      icon = 'ph-trash'
    } else if (type === 'restore') {
      title = 'Pulihkan Pesan'
      message = 'Pesan ini akan dikembalikan ke Kotak Masuk.'
      icon = 'ph-arrow-counter-clockwise'
    } else if (type === 'bulk_delete') {
      title = 'Hapus Massal'
      message = `Anda akan memindahkan ${selectedIds.length} pesan ke Sampah. Lanjutkan?`
      icon = 'ph-trash'
    } else if (type === 'bulk_restore') {
      title = 'Pulihkan Massal'
      message = `Anda akan memulihkan ${selectedIds.length} pesan. Lanjutkan?`
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
    setIsViewOpen(false)

    try {
      if (type === 'delete') {
        await kontakService.remove(id)
        showToast('Pesan berhasil dihapus.', 'success')
      } else if (type === 'restore') {
        await kontakService.restore(id)
        showToast('Pesan berhasil dipulihkan.', 'success')
      } else if (type === 'bulk_delete') {
        await kontakService.bulkDelete(selectedIds)
        setSelectedIds([])
        showToast('Pesan massal berhasil dihapus.', 'success')
      } else if (type === 'bulk_restore') {
        await kontakService.bulkRestore(selectedIds)
        setSelectedIds([])
        showToast('Pesan massal berhasil dipulihkan.', 'success')
      }
      fetchMessages()
    } catch (err) {
      const parsed = applyError(err)
      if (Object.keys(parsed.fieldErrors).length === 0) {
        showToast(parsed.message || 'Gagal melakukan aksi.', 'error')
      }
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

  const headerContent = (
    <form onSubmit={handleSearchSubmit} className="flex items-center gap-2 w-full max-w-2xl animate-fade-in-up">
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
      
      {currentTab === 'active' && (
        <select
          value={filterStatus}
          onChange={(e) => {
            setFilterStatus(e.target.value)
            setPage(1)
          }}
          className="shrink-0 bg-gray-50 border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-brand-500 focus:bg-white transition-colors cursor-pointer"
        >
          <option value="all">Semua Status</option>
          <option value="unread">Belum Dibaca</option>
          <option value="read">Sudah Dibaca</option>
        </select>
      )}

      <button
        type="submit"
        className="bg-brand-600 hover:bg-brand-700 text-white px-3 py-2 rounded-lg text-sm font-medium transition-all btn-press shrink-0"
      >
        Cari
      </button>
      <button
        type="button"
        onClick={handleReset}
        className="shrink-0 bg-gray-50 border border-gray-200 text-gray-700 hover:bg-gray-100 px-3 py-2 rounded-lg text-sm font-medium flex items-center gap-2 transition-all btn-press"
        title="Reset Filter"
      >
        <i className="ph ph-arrows-counter-clockwise text-lg" />
      </button>
    </form>
  )

  return (
    <AdminLayout title="Kelola Pesan Kontak" headerContent={headerContent}>
      <div className="max-w-7xl mx-auto space-y-6 animate-fade-in-up">
        {/* Header Actions & Filter */}
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
          {/* Tabs */}
          <div className="flex bg-white rounded-lg p-1 border border-gray-200 shadow-sm">
            <button
              onClick={() => setCurrentTab('active')}
              className={`px-4 py-1.5 rounded-md text-sm transition-all duration-200 ${
                currentTab === 'active'
                  ? 'bg-brand-50 text-brand-600 shadow-sm font-medium'
                  : 'text-gray-500 hover:text-gray-700'
              }`}
            >
              Kotak Masuk
            </button>
            <button
              onClick={() => setCurrentTab('trash')}
              className={`px-4 py-1.5 rounded-md text-sm transition-all duration-200 flex items-center gap-2 ${
                currentTab === 'trash'
                  ? 'bg-red-50 text-red-600 shadow-sm font-medium'
                  : 'text-gray-500 hover:text-gray-700'
              }`}
            >
              <i className="ph ph-trash" /> Sampah
            </button>
          </div>
        </div>

        {/* Bulk Actions Bar */}
        {selectedIds.length > 0 && (
          <div className="bg-indigo-50 border border-indigo-100 rounded-lg p-3 flex items-center justify-between shadow-sm animate-fade-in">
            <span className="text-sm text-indigo-800 font-medium">{selectedIds.length} pesan terpilih</span>
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
              <div className="py-12 text-center text-gray-500">Memuat pesan...</div>
            ) : error ? (
              <div className="py-12 text-center text-red-500">{error}</div>
            ) : messages.length === 0 ? (
              <div className="py-12 text-center text-gray-500">
                <i className="ph ph-envelope-open text-4xl text-gray-300 mb-2 block mx-auto" />
                Tidak ada pesan untuk ditampilkan.
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
                        className="rounded border-gray-300 text-brand-600 focus:ring-brand-500"
                      />
                    </th>
                    <th className="p-4">Pengirim</th>
                    <th className="p-4">Subjek</th>
                    <th className="p-4 w-40">Tanggal</th>
                    <th className="p-4 text-right">Aksi</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-100 text-sm">
                  {messages.map((item) => {
                    const isUnread = !item.is_read && currentTab === 'active'
                    return (
                      <tr
                        key={item.id}
                        onClick={() => openViewModal(item)}
                        className={`hover:bg-gray-50 transition-colors group cursor-pointer admin-row ${
                          isUnread ? 'font-semibold text-gray-900 bg-brand-50/20' : 'text-gray-700'
                        }`}
                      >
                        <td className="p-4 text-center" onClick={(e) => e.stopPropagation()}>
                          <input
                            type="checkbox"
                            checked={selectedIds.includes(item.id)}
                            onChange={() => handleSelectItem(item.id)}
                            className="rounded border-gray-300 text-brand-600 focus:ring-brand-500"
                          />
                        </td>
                        <td className="p-4">
                          <div className="flex flex-col">
                            <span className="truncate max-w-[200px]">{item.name}</span>
                            <span className="text-xs font-normal text-gray-500">{item.email}</span>
                          </div>
                        </td>
                        <td className="p-4">
                          <div className="flex items-center gap-2">
                            {isUnread && <span className="w-2 h-2 rounded-full bg-brand-500 shrink-0" />}
                            <span className="truncate max-w-xs">{item.subject}</span>
                          </div>
                        </td>
                        <td className="p-4">
                          <span className="text-gray-500 font-normal">{item.date}</span>
                        </td>
                        <td className="p-4 text-right" onClick={(e) => e.stopPropagation()}>
                          <div className="flex items-center justify-end gap-2">
                            {currentTab === 'active' ? (
                              <>
                                <button
                                  onClick={() => openViewModal(item)}
                                  className="p-1.5 text-gray-500 hover:text-brand-600 hover:bg-brand-50 rounded"
                                  title="Lihat Detail"
                                >
                                  <i className="ph ph-eye text-lg" />
                                </button>
                                <button
                                  onClick={() => triggerConfirm('delete', item.id)}
                                  className="p-1.5 text-gray-500 hover:text-red-600 hover:bg-red-50 rounded opacity-0 group-hover:opacity-100 transition-opacity"
                                  title="Hapus"
                                >
                                  <i className="ph ph-trash text-lg" />
                                </button>
                              </>
                            ) : (
                              <button
                                onClick={() => triggerConfirm('restore', item.id)}
                                className="p-1.5 text-gray-500 hover:text-emerald-600 hover:bg-emerald-50 rounded opacity-0 group-hover:opacity-100 transition-opacity"
                                title="Pulihkan"
                              >
                                <i className="ph ph-arrow-counter-clockwise text-lg" />
                              </button>
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
          {!loading && !error && messages.length > 0 && (
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

      {/* VIEW DETAIL MODAL */}
      {isViewOpen && viewData && (
        <div className="fixed inset-0 z-50 overflow-y-auto" role="dialog" aria-modal="true">
          <div className="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
            <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" onClick={() => setIsViewOpen(false)} />
            <span className="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>
            
            <div className="inline-block align-bottom bg-white rounded-xl text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-2xl sm:w-full">
              <div className="bg-white">
                <div className="border-b border-gray-200 px-6 py-4 flex items-center justify-between">
                  <h3 className="text-lg leading-6 font-heading font-semibold text-gray-900">Detail Pesan</h3>
                  <div className="flex items-center gap-2">
                    <button
                      onClick={() => triggerConfirm('delete', viewData.id)}
                      className="text-gray-400 hover:text-red-500 p-1"
                      title="Hapus"
                    >
                      <i className="ph ph-trash text-xl" />
                    </button>
                    <button
                      onClick={() => setIsViewOpen(false)}
                      className="text-gray-400 hover:text-gray-500 p-1"
                    >
                      <i className="ph ph-x text-xl" />
                    </button>
                  </div>
                </div>
                
                <div className="px-6 py-5 space-y-6">
                  <div className="flex justify-between items-start">
                    <div>
                      <h4 className="text-xl font-bold text-gray-900">{viewData.subject}</h4>
                      <div className="mt-2 flex items-center gap-2">
                        <div className="w-10 h-10 rounded-full bg-brand-100 text-brand-600 flex items-center justify-center font-bold text-sm">
                          {(viewData.name || 'C').charAt(0).toUpperCase()}
                        </div>
                        <div>
                          <p className="text-sm font-medium text-gray-900">{viewData.name}</p>
                          <p className="text-xs text-gray-500">{viewData.email}</p>
                        </div>
                      </div>
                    </div>
                    <span className="text-xs text-gray-500">{viewData.date}</span>
                  </div>
                  
                  <div className="bg-gray-50 rounded-lg p-4 border border-gray-100 text-sm text-gray-800 whitespace-pre-wrap font-sans">
                    {viewData.message}
                  </div>

                  <div className="border-t border-gray-200 pt-4">
                    <a
                      href={`mailto:${viewData.email}?subject=Re: ${encodeURIComponent(viewData.subject)}`}
                      className="inline-flex items-center gap-2 px-4 py-2 bg-white border border-gray-300 shadow-sm rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
                    >
                      <i className="ph ph-paper-plane-tilt text-lg" /> Balas via Email
                    </a>
                  </div>
                </div>
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
                    confirm.type.includes('delete') ? 'bg-red-100 text-red-600' : 'bg-emerald-100 text-emerald-600'
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
                    confirm.type.includes('delete') ? 'bg-red-600 hover:bg-red-700' : 'bg-emerald-600 hover:bg-emerald-700'
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
