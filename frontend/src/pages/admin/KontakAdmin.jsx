import { useState, useEffect, useCallback, useMemo, useRef } from 'react'
import AdminLayout from '../../layouts/AdminLayout'
import { kontakService } from '../../services/kontakService'
import ConfirmDialog from '../../components/admin/ConfirmDialog'
import ToastNotification from '../../components/admin/ToastNotification'
import { parseApiError } from '../../utils/parseApiError'

const PAGE_SIZE = 10
const SEARCH_DEBOUNCE_MS = 350

// Tanggal + jam — untuk pesan masuk, jam ikut relevan (beda dengan berita/kegiatan).
function formatDateTime(value) {
  if (!value) return '-'
  const parsed = new Date(value)
  if (Number.isNaN(parsed.getTime())) return value
  return parsed.toLocaleString('id-ID', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

export default function KontakAdmin() {
  const [messages, setMessages] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  const [currentTab, setCurrentTab] = useState('active') // active | trash

  // searchInput = yang diketik user, search = nilai ter-debounce yang dikirim ke API.
  const [searchInput, setSearchInput] = useState('')
  const [search, setSearch] = useState('')
  const [filterStatus, setFilterStatus] = useState('all') // all | unread | read
  const [filterSort, setFilterSort] = useState('newest') // newest | oldest

  // Pagination (server-side)
  const [page, setPage] = useState(1)
  const [total, setTotal] = useState(0)
  const [totalPages, setTotalPages] = useState(1)

  // Selection (hanya baris di halaman yang sedang tampil)
  const [selectedIds, setSelectedIds] = useState([])

  // Modal detail
  const [viewData, setViewData] = useState(null)
  const [isViewOpen, setIsViewOpen] = useState(false)
  const [viewLoading, setViewLoading] = useState(false)

  const [confirm, setConfirm] = useState({
    isOpen: false,
    type: 'danger',
    title: '',
    message: '',
    action: null,
  })

  const [toast, setToast] = useState({ show: false, message: '', type: 'success' })

  const selectAllRef = useRef(null)
  // Kata kunci yang sudah benar-benar dikirim ke API — pembanding supaya debounce
  // tidak memicu reset halaman/refetch saat nilainya tidak berubah.
  const appliedSearchRef = useRef('')
  // Token urutan request: klik Prev/Next cepat bisa membuat respons datang tidak
  // berurutan. Hanya respons dari request terakhir yang boleh dipakai.
  const listTokenRef = useRef(0)
  const detailTokenRef = useRef(0)

  const showToast = useCallback((message, type = 'success') => {
    setToast({ show: true, message, type })
  }, [])

  // --- Debounce pencarian: satu request setelah user berhenti mengetik ---
  useEffect(() => {
    const timer = setTimeout(() => {
      const next = searchInput.trim()
      if (appliedSearchRef.current === next) return
      appliedSearchRef.current = next
      setSearch(next)
      setPage(1)
    }, SEARCH_DEBOUNCE_MS)
    return () => clearTimeout(timer)
  }, [searchInput])

  const fetchMessages = useCallback(() => {
    const token = ++listTokenRef.current
    setLoading(true)
    setError(null)

    const params = {
      tab: currentTab,
      page,
      limit: PAGE_SIZE,
      sort: filterSort,
    }
    if (search) params.search = search
    // Filter dibaca/belum dibaca hanya berlaku di Kotak Masuk.
    if (currentTab === 'active' && filterStatus !== 'all') params.status = filterStatus

    kontakService.list(params)
      .then(res => {
        if (token !== listTokenRef.current) return // respons basi
        const data = res?.data || {}
        // Backend mengembalikan { kontak: [...], meta: {...} } (lihat KontakListResponse).
        const list = Array.isArray(data) ? data : (data.kontak || data.items || [])
        setMessages(list)

        // Buang id terpilih yang sudah tidak ada di halaman ini (mis. dihapus
        // dari sesi lain), supaya aksi massal tidak mengirim id hantu.
        setSelectedIds(prev => {
          if (prev.length === 0) return prev
          const visible = new Set(list.map(m => m.id))
          const next = prev.filter(id => visible.has(id))
          return next.length === prev.length ? prev : next
        })

        const meta = data.meta || {}
        const totalData = Number(meta.total_data ?? meta.total ?? list.length) || 0
        const pages = Number(meta.total_pages) || (totalData > 0 ? Math.ceil(totalData / PAGE_SIZE) : 1)

        setTotal(totalData)
        setTotalPages(pages)
        // Jaring pengaman: halaman di luar jangkauan (mis. data menyusut dari
        // sesi lain) ditarik kembali ke halaman terakhir yang valid.
        if (page > pages) setPage(pages)
      })
      .catch(err => {
        if (token !== listTokenRef.current) return
        setMessages([])
        setTotal(0)
        setTotalPages(1)
        setError(parseApiError(err).message || 'Gagal memuat pesan masuk.')
      })
      .finally(() => {
        if (token !== listTokenRef.current) return
        setLoading(false)
      })
  }, [currentTab, search, filterStatus, filterSort, page])

  useEffect(() => {
    fetchMessages()
  }, [fetchMessages])

  // Pilihan tidak boleh terbawa antar tab maupun antar halaman.
  useEffect(() => {
    setSelectedIds([])
  }, [currentTab, filterStatus, filterSort, search, page])

  // Ganti tab/filter selalu dibarengi reset halaman di handler-nya (bukan lewat
  // effect terpisah) supaya tidak ada request ganda dengan page yang basi.
  function changeTab(tab) {
    if (tab === currentTab) return
    setCurrentTab(tab)
    setPage(1)
  }

  function changeStatus(value) {
    setFilterStatus(value)
    setPage(1)
  }

  function changeSort(value) {
    setFilterSort(value)
    setPage(1)
  }

  function resetFilter() {
    appliedSearchRef.current = ''
    setSearchInput('')
    setSearch('')
    setFilterStatus('all')
    setFilterSort('newest')
    setPage(1)
    showToast('Filter direset.', 'info')
  }

  // Urutan bukan filter, dan filter status tidak dikirim di tab Sampah — keduanya
  // tidak boleh membuat daftar diklaim "terfilter".
  const hasActiveFilter = Boolean(search) || (currentTab === 'active' && filterStatus !== 'all')

  // --- Selection ---
  const isAllSelected = messages.length > 0 && selectedIds.length === messages.length
  const isPartiallySelected = selectedIds.length > 0 && !isAllSelected

  useEffect(() => {
    // messages.length ikut jadi dependency: checkbox header baru ada di DOM setelah
    // tabel ter-render, jadi state indeterminate perlu dipasang ulang tiap data berubah.
    if (selectAllRef.current) selectAllRef.current.indeterminate = isPartiallySelected
  }, [isPartiallySelected, messages.length])

  const toggleAll = () => {
    setSelectedIds(isAllSelected ? [] : messages.map(m => m.id))
  }

  const toggleOne = (id) => {
    setSelectedIds(prev => prev.includes(id) ? prev.filter(x => x !== id) : [...prev, id])
  }

  // --- Detail ---
  const openViewModal = async (item) => {
    const token = ++detailTokenRef.current
    setViewData(item)
    setIsViewOpen(true)
    setViewLoading(true)
    try {
      // Detail dipanggil ke server: isinya lengkap sekaligus menandai pesan sudah dibaca.
      const res = await kontakService.getById(item.id)
      if (token !== detailTokenRef.current) return // baris lain sudah diklik
      const detail = res?.data
      if (detail) {
        setViewData(detail)
        if (detail.is_read) {
          setMessages(prev => prev.map(m => (m.id === item.id ? { ...m, is_read: true } : m)))
        }
      }
    } catch (err) {
      if (token !== detailTokenRef.current) return
      showToast(parseApiError(err).message || 'Gagal memuat detail pesan.', 'error')
    } finally {
      if (token === detailTokenRef.current) setViewLoading(false)
    }
  }

  const closeViewModal = () => {
    setIsViewOpen(false)
    // Membuka pesan otomatis menandainya "dibaca". Saat filter "Belum Dibaca" aktif,
    // daftar perlu dimuat ulang supaya tidak menampilkan pesan yang sudah tak cocok.
    if (currentTab === 'active' && filterStatus === 'unread') fetchMessages()
  }

  // --- Confirm actions (pola sama dengan modul admin lain) ---

  // Setelah sejumlah pesan hilang dari halaman ini: kalau halaman jadi kosong dan
  // bukan halaman pertama, mundur satu halaman. setPage memicu refetch lewat
  // dependency fetchMessages, jadi cukup satu request — bukan render kosong dulu
  // lalu refetch lagi.
  const refreshAfterRemoval = (removedCount) => {
    if (page > 1 && messages.length - removedCount <= 0) {
      setPage(p => Math.max(p - 1, 1))
    } else {
      fetchMessages()
    }
  }

  // Backend melaporkan jumlah baris yang benar-benar terpengaruh; kalau ada yang
  // sudah diproses lewat sesi lain, angkanya bisa lebih kecil dari yang dipilih.
  const affectedCount = (res, fallback) => Number(res?.data?.affected ?? fallback) || fallback

  function confirmAction(type, id = null) {
    const selectedCount = selectedIds.length

    const configs = {
      delete: {
        type: 'danger',
        title: 'Hapus Pesan',
        message: 'Pesan ini akan dipindahkan ke Sampah. Lanjutkan?',
        action: async () => {
          await kontakService.remove(id)
          setIsViewOpen(false)
          setSelectedIds(prev => prev.filter(x => x !== id))
          refreshAfterRemoval(1)
          showToast('Pesan berhasil dipindahkan ke Sampah.')
        },
      },
      restore: {
        type: 'success',
        title: 'Pulihkan Pesan',
        message: 'Pesan ini akan dikembalikan ke Kotak Masuk. Lanjutkan?',
        action: async () => {
          await kontakService.restore(id)
          setIsViewOpen(false)
          setSelectedIds(prev => prev.filter(x => x !== id))
          refreshAfterRemoval(1)
          showToast('Pesan berhasil dipulihkan.')
        },
      },
      bulk_delete: {
        type: 'danger',
        title: 'Hapus Massal',
        message: `${selectedCount} pesan akan dipindahkan ke Sampah. Lanjutkan?`,
        action: async () => {
          const res = await kontakService.bulkDelete(selectedIds)
          const n = affectedCount(res, selectedCount)
          setSelectedIds([])
          refreshAfterRemoval(n)
          showToast(`${n} pesan berhasil dipindahkan ke Sampah.`)
        },
      },
      bulk_restore: {
        type: 'success',
        title: 'Pulihkan Massal',
        message: `${selectedCount} pesan akan dikembalikan ke Kotak Masuk. Lanjutkan?`,
        action: async () => {
          const res = await kontakService.bulkRestore(selectedIds)
          const n = affectedCount(res, selectedCount)
          setSelectedIds([])
          refreshAfterRemoval(n)
          showToast(`${n} pesan berhasil dipulihkan.`)
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
    if (!action) return
    try {
      await action()
    } catch (err) {
      showToast(parseApiError(err).message || 'Gagal melakukan aksi.', 'error')
    }
  }

  // --- Pagination ---
  const pageNumbers = useMemo(() => {
    const maxButtons = 5
    if (totalPages <= maxButtons) {
      return Array.from({ length: totalPages }, (_, i) => i + 1)
    }
    const end = Math.min(totalPages, Math.max(page + 2, maxButtons))
    const start = Math.max(1, end - maxButtons + 1)
    return Array.from({ length: end - start + 1 }, (_, i) => start + i)
  }, [page, totalPages])

  const showingFrom = total === 0 ? 0 : (page - 1) * PAGE_SIZE + 1
  const showingTo = Math.min(page * PAGE_SIZE, total)

  // Status pesan di modal dibaca dari payload, bukan dari tab yang sedang dibuka:
  // kalau fetch detail gagal, UI tidak boleh mengklaim pesan sudah ditandai dibaca.
  const viewIsTrashed = Boolean(viewData?.deleted_at) || currentTab === 'trash'
  const viewIsRead = Boolean(viewData?.is_read)

  const emptyStateText = search
    ? `Tidak ada pesan yang cocok dengan "${search}".`
    : currentTab === 'trash'
      ? 'Sampah kosong.'
      : filterStatus === 'unread'
        ? 'Semua pesan sudah dibaca.'
        : 'Belum ada pesan masuk.'

  const headerContent = (
    <div className="flex items-center gap-2 w-full max-w-3xl animate-fade-in-up">
      <div className="relative w-full">
        <i className="ph ph-magnifying-glass absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
        <input
          type="text"
          value={searchInput}
          onChange={(e) => setSearchInput(e.target.value)}
          placeholder="Cari nama / email / subjek / isi pesan..."
          aria-label="Cari pesan masuk"
          className="w-full pl-9 pr-4 py-2 bg-gray-50 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-brand-500 focus:bg-white transition-colors"
        />
      </div>

      {currentTab === 'active' && (
        <select
          value={filterStatus}
          onChange={(e) => changeStatus(e.target.value)}
          aria-label="Filter status dibaca"
          className="shrink-0 bg-gray-50 border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-brand-500 focus:bg-white transition-colors cursor-pointer"
        >
          <option value="all">Semua Status</option>
          <option value="unread">Belum Dibaca</option>
          <option value="read">Sudah Dibaca</option>
        </select>
      )}

      <select
        value={filterSort}
        onChange={(e) => changeSort(e.target.value)}
        aria-label="Urutkan pesan"
        className="shrink-0 bg-gray-50 border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-brand-500 focus:bg-white transition-colors cursor-pointer"
      >
        <option value="newest">Terbaru</option>
        <option value="oldest">Terlama</option>
      </select>

      <button
        type="button"
        onClick={resetFilter}
        className="shrink-0 bg-gray-50 border border-gray-200 text-gray-700 hover:bg-gray-100 px-3 py-2 rounded-lg text-sm font-medium flex items-center gap-2 transition-all btn-press"
      >
        <i className="ph ph-arrows-counter-clockwise text-lg" /> Reset
      </button>
    </div>
  )

  return (
    <AdminLayout title="Kelola Pesan Kontak" headerContent={headerContent}>
      <ToastNotification
        show={toast.show}
        message={toast.message}
        type={toast.type}
        onClose={() => setToast(prev => ({ ...prev, show: false }))}
      />
      <ConfirmDialog
        isOpen={confirm.isOpen}
        type={confirm.type}
        title={confirm.title}
        message={confirm.message}
        onConfirm={executeConfirm}
        onClose={() => setConfirm(prev => ({ ...prev, isOpen: false, action: null }))}
      />

      <div className="max-w-7xl mx-auto space-y-6 animate-fade-in-up">
        {/* Tabs */}
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
          <div className="flex bg-white rounded-lg p-1 border border-gray-200 shadow-sm">
            <button
              onClick={() => changeTab('active')}
              className={`px-4 py-1.5 rounded-md text-sm transition-all duration-200 ${
                currentTab === 'active'
                  ? 'bg-brand-50 text-brand-600 shadow-sm font-medium'
                  : 'text-gray-500 hover:text-gray-700'
              }`}
            >
              Kotak Masuk
            </button>
            <button
              onClick={() => changeTab('trash')}
              className={`px-4 py-1.5 rounded-md text-sm transition-all duration-200 flex items-center gap-2 ${
                currentTab === 'trash'
                  ? 'bg-red-50 text-red-600 shadow-sm font-medium'
                  : 'text-gray-500 hover:text-gray-700'
              }`}
            >
              <i className="ph ph-trash" /> Sampah
            </button>
          </div>

          {!loading && !error && (
            <span className="text-sm text-gray-500">
              {total} pesan{hasActiveFilter ? ' (terfilter)' : ''}
            </span>
          )}
        </div>

        {/* Bulk Actions Bar */}
        {selectedIds.length > 0 && (
          <div className="flex flex-wrap items-center gap-2 bg-brand-50/60 border border-brand-100 rounded-xl px-4 py-2.5 shadow-sm animate-fade-in">
            <span className="text-sm font-semibold text-brand-700">{selectedIds.length} pesan terpilih</span>
            <div className="flex gap-2 ml-auto">
              {currentTab === 'active' ? (
                <button
                  onClick={() => confirmAction('bulk_delete')}
                  className="bg-red-600 hover:bg-red-700 text-white px-3.5 py-1.5 rounded-lg text-xs font-semibold flex items-center gap-1.5 transition-colors"
                >
                  <i className="ph ph-trash" /> Hapus Massal
                </button>
              ) : (
                <button
                  onClick={() => confirmAction('bulk_restore')}
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
              <div className="py-12 text-center text-gray-500" role="status" aria-live="polite">
                <i className="ph ph-spinner text-2xl animate-spin block mx-auto mb-2" />
                Memuat pesan...
              </div>
            ) : error ? (
              <div className="py-12 text-center" role="alert">
                <i className="ph ph-warning-circle text-4xl text-red-300 mb-2 block mx-auto" />
                <p className="text-red-600">{error}</p>
                <button
                  onClick={fetchMessages}
                  className="mt-3 inline-flex items-center gap-2 px-3 py-1.5 rounded-lg border border-gray-200 text-sm text-gray-700 hover:bg-gray-50 transition-colors"
                >
                  <i className="ph ph-arrows-clockwise" /> Coba lagi
                </button>
              </div>
            ) : messages.length === 0 ? (
              <div className="py-12 text-center text-gray-500">
                <i className="ph ph-envelope-open text-4xl text-gray-300 mb-2 block mx-auto" />
                {emptyStateText}
                {hasActiveFilter && (
                  <button
                    onClick={resetFilter}
                    className="mt-3 block mx-auto text-sm text-brand-600 hover:text-brand-700 font-medium"
                  >
                    Reset filter
                  </button>
                )}
              </div>
            ) : (
              <table className="w-full text-left border-collapse">
                <thead>
                  <tr className="bg-gray-50 border-b border-gray-200 text-xs uppercase tracking-wider text-gray-500 font-semibold">
                    <th scope="col" className="p-4 w-12 text-center">
                      <input
                        ref={selectAllRef}
                        type="checkbox"
                        onChange={toggleAll}
                        checked={isAllSelected}
                        aria-label="Pilih semua pesan di halaman ini"
                        className="rounded border-gray-300 text-brand-600 focus:ring-brand-500"
                      />
                    </th>
                    <th scope="col" className="p-4">Pengirim</th>
                    <th scope="col" className="p-4">Subjek</th>
                    <th scope="col" className="p-4 w-48">{currentTab === 'trash' ? 'Dihapus' : 'Tanggal'}</th>
                    <th scope="col" className="p-4 text-right">Aksi</th>
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
                            onChange={() => toggleOne(item.id)}
                            aria-label={`Pilih pesan dari ${item.nama || 'pengirim'}`}
                            className="rounded border-gray-300 text-brand-600 focus:ring-brand-500"
                          />
                        </td>
                        <td className="p-4">
                          <div className="flex flex-col">
                            <span className="truncate max-w-[200px]">{item.nama || '-'}</span>
                            <span className="text-xs font-normal text-gray-500">{item.email || '-'}</span>
                          </div>
                        </td>
                        <td className="p-4">
                          <div className="flex items-center gap-2">
                            {isUnread && <span className="w-2 h-2 rounded-full bg-brand-500 shrink-0" title="Belum dibaca" />}
                            <span className="truncate max-w-xs">{item.subjek || '-'}</span>
                          </div>
                        </td>
                        <td className="p-4">
                          <span className="text-gray-500 font-normal">
                            {formatDateTime(currentTab === 'trash' ? (item.deleted_at || item.created_at) : item.created_at)}
                          </span>
                        </td>
                        <td className="p-4 text-right" onClick={(e) => e.stopPropagation()}>
                          <div className="flex items-center justify-end gap-2">
                            <button
                              onClick={() => openViewModal(item)}
                              className="p-1.5 text-gray-500 hover:text-brand-600 hover:bg-brand-50 rounded"
                              title="Lihat Detail"
                              aria-label="Lihat detail pesan"
                            >
                              <i className="ph ph-eye text-lg" />
                            </button>
                            {currentTab === 'active' ? (
                              <button
                                onClick={() => confirmAction('delete', item.id)}
                                className="p-1.5 text-gray-500 hover:text-red-600 hover:bg-red-50 rounded"
                                title="Hapus"
                                aria-label="Hapus pesan"
                              >
                                <i className="ph ph-trash text-lg" />
                              </button>
                            ) : (
                              <button
                                onClick={() => confirmAction('restore', item.id)}
                                className="p-1.5 text-gray-500 hover:text-emerald-600 hover:bg-emerald-50 rounded"
                                title="Pulihkan"
                                aria-label="Pulihkan pesan"
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
                {Array.from({ length: totalPages || 1 }, (_, i) => i + 1).map((n) => (
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
      </div>

      {/* DETAIL MODAL */}
      {isViewOpen && viewData && (
        // z-50 (bukan z-[60]) — ConfirmDialog ada di z-[60] dan harus tetap di atas
        // modal ini, kalau tidak overlay-nya menelan klik tombol konfirmasi.
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4" role="dialog" aria-modal="true" aria-label="Detail pesan">
          <div className="fixed inset-0 bg-black/40 backdrop-blur-[2px]" onClick={closeViewModal} />

          <div className="relative bg-white rounded-2xl shadow-2xl max-w-2xl w-full max-h-[90vh] flex flex-col overflow-hidden">
            <div className="border-b border-gray-200 px-6 py-4 flex items-center justify-between">
              <h3 className="text-lg leading-6 font-heading font-semibold text-gray-900">Detail Pesan</h3>
              <div className="flex items-center gap-2">
                {viewIsTrashed ? (
                  <button
                    onClick={() => confirmAction('restore', viewData.id)}
                    className="text-gray-400 hover:text-emerald-600 p-1"
                    title="Pulihkan"
                    aria-label="Pulihkan pesan"
                  >
                    <i className="ph ph-arrow-counter-clockwise text-xl" />
                  </button>
                ) : (
                  <button
                    onClick={() => confirmAction('delete', viewData.id)}
                    className="text-gray-400 hover:text-red-500 p-1"
                    title="Hapus"
                    aria-label="Hapus pesan"
                  >
                    <i className="ph ph-trash text-xl" />
                  </button>
                )}
                <button
                  onClick={closeViewModal}
                  className="text-gray-400 hover:text-gray-600 p-1"
                  aria-label="Tutup"
                >
                  <i className="ph ph-x text-xl" />
                </button>
              </div>
            </div>

            <div className="px-6 py-5 space-y-6 overflow-y-auto">
              <div className="flex justify-between items-start gap-4">
                <div className="min-w-0">
                  <h4 className="text-xl font-bold text-gray-900 break-words">{viewData.subjek || '(tanpa subjek)'}</h4>
                  <div className="mt-2 flex items-center gap-2">
                    <div className="w-10 h-10 rounded-full bg-brand-100 text-brand-600 flex items-center justify-center font-bold text-sm shrink-0">
                      {(viewData.nama || 'K').charAt(0).toUpperCase()}
                    </div>
                    <div className="min-w-0">
                      <p className="text-sm font-medium text-gray-900 truncate">{viewData.nama || '-'}</p>
                      <p className="text-xs text-gray-500 truncate">{viewData.email || '-'}</p>
                    </div>
                  </div>
                </div>
                <div className="text-right shrink-0">
                  <span className="text-xs text-gray-500 block">{formatDateTime(viewData.created_at)}</span>
                  {viewIsTrashed ? (
                    <span className="mt-1 inline-flex items-center gap-1 text-xs font-medium text-red-600 bg-red-50 px-2 py-0.5 rounded-full">
                      <i className="ph ph-trash" /> Di Sampah
                    </span>
                  ) : viewIsRead ? (
                    <span className="mt-1 inline-flex items-center gap-1 text-xs font-medium text-emerald-600 bg-emerald-50 px-2 py-0.5 rounded-full">
                      <i className="ph ph-check-circle" /> Sudah Dibaca
                    </span>
                  ) : (
                    <span className="mt-1 inline-flex items-center gap-1 text-xs font-medium text-amber-600 bg-amber-50 px-2 py-0.5 rounded-full">
                      <i className="ph ph-envelope" /> Belum Dibaca
                    </span>
                  )}
                </div>
              </div>

              <div className="bg-gray-50 rounded-lg p-4 border border-gray-100 text-sm text-gray-800 whitespace-pre-wrap break-words">
                {viewLoading && !viewData.pesan
                  ? 'Memuat isi pesan...'
                  : (viewData.pesan || '(pesan kosong)')}
              </div>

              {viewData.response_note && (
                <div>
                  <p className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-1">Catatan Tindak Lanjut</p>
                  <div className="bg-amber-50 border border-amber-100 rounded-lg p-3 text-sm text-amber-900 whitespace-pre-wrap break-words">
                    {viewData.response_note}
                  </div>
                </div>
              )}

              {viewData.email && (
                <div className="border-t border-gray-200 pt-4">
                  <a
                    href={`mailto:${viewData.email}?subject=${encodeURIComponent(`Re: ${viewData.subjek || ''}`)}`}
                    className="inline-flex items-center gap-2 px-4 py-2 bg-white border border-gray-300 shadow-sm rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
                  >
                    <i className="ph ph-paper-plane-tilt text-lg" /> Balas via Email
                  </a>
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </AdminLayout>
  )
}
