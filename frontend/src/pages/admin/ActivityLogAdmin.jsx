import { useState, useEffect, useCallback } from 'react'
import AdminLayout from '../../layouts/AdminLayout'
import { activityLogService } from '../../services/activityLogService'
import { mapLogRowForCsv, mapLogRowForDisplay, extractPaginationMeta } from './activityLogMapping'

export default function ActivityLogAdmin() {
  const [logs, setLogs] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  // Filters state
  const [search, setSearch] = useState('')
  const [filterRole, setFilterRole] = useState('')
  const [filterEntity, setFilterEntity] = useState('')
  const [filterRisk, setFilterRisk] = useState('')

  // Pagination state
  const [page, setPage] = useState(1)
  const [limit] = useState(10)
  const [total, setTotal] = useState(0)
  const [totalPages, setTotalPages] = useState(1)

  const fetchLogs = useCallback(() => {
    setLoading(true)
    activityLogService.list({
      search,
      role: filterRole,
      entity: filterEntity,
      risk: filterRisk,
      page,
      limit
    })
      .then(res => {
        if (res.success && res.data) {
          setLogs(res.data.items || [])
          const paginationMeta = extractPaginationMeta(res.data)
          if (paginationMeta) {
            setTotal(paginationMeta.total)
            setTotalPages(paginationMeta.totalPages)
          }
          setError(null)
        } else {
          setError('Gagal memuat log aktivitas')
        }
      })
      .catch((err) => {
        setError(err.message || 'Kesalahan koneksi ke server')
      })
      .finally(() => setLoading(false))
  }, [search, filterRole, filterEntity, filterRisk, page, limit])

  useEffect(() => {
    fetchLogs()
  }, [fetchLogs])

  const handleSearchSubmit = (e) => {
    e.preventDefault()
    setPage(1)
    // fetchLogs() otomatis dipanggil via useEffect [fetchLogs] saat search/page berubah
  }

  const handleReset = () => {
    setSearch('')
    setFilterRole('')
    setFilterEntity('')
    setFilterRisk('')
    setPage(1)
  }

  // Client-side CSV Download
  const downloadCSV = () => {
    if (logs.length === 0) return
    const headers = ['Waktu', 'Aktor', 'Role', 'Aktivitas', 'Entitas', 'Keterangan', 'IP Address', 'Device/User Agent', 'Risiko']
    const rows = logs.map(mapLogRowForCsv)

    const csvContent = 'data:text/csv;charset=utf-8,' 
      + [headers.join(','), ...rows.map(e => e.map(val => `"${val.replace(/"/g, '""')}"`).join(','))].join('\n')
    
    const encodedUri = encodeURI(csvContent)
    const link = document.createElement('a')
    link.setAttribute('href', encodedUri)
    link.setAttribute('download', `activity_log_${new Date().toISOString().slice(0, 10)}.csv`)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
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

  const getActionBadgeClass = (action = '') => {
    const act = action.toLowerCase()
    if (act.includes('delete') || act.includes('remove')) {
      return 'bg-red-50 text-red-600 border-red-200'
    }
    if (act.includes('restore') || act.includes('create') || act.includes('insert')) {
      return 'bg-emerald-50 text-emerald-600 border-emerald-200'
    }
    return 'bg-blue-50 text-blue-600 border-blue-200'
  }

  const getRiskBadgeClass = (risk = '') => {
    const r = risk.toLowerCase()
    if (r === 'high') return 'bg-red-50 text-red-700 border-red-200'
    if (r === 'medium') return 'bg-amber-50 text-amber-700 border-amber-200'
    return 'bg-emerald-50 text-emerald-700 border-emerald-200'
  }

  const headerContent = (
    <div className="flex items-center gap-2 w-full max-w-4xl justify-end animate-fade-in-up">
      <select
        value={filterRole}
        onChange={(e) => {
          setFilterRole(e.target.value)
          setPage(1)
        }}
        className="shrink-0 bg-gray-50 border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-brand-500 focus:bg-white transition-colors cursor-pointer"
      >
        <option value="">Semua Role</option>
        <option value="super_admin">Super Admin</option>
        <option value="admin">Admin</option>
      </select>

      <select
        value={filterEntity}
        onChange={(e) => {
          setFilterEntity(e.target.value)
          setPage(1)
        }}
        className="shrink-0 bg-gray-50 border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-brand-500 focus:bg-white transition-colors cursor-pointer"
      >
        <option value="">Semua Entitas</option>
        <option value="user">User</option>
        <option value="auth">Auth</option>
        <option value="berita">Berita</option>
        <option value="kegiatan">Kegiatan</option>
        <option value="pengurus">Pengurus</option>
        <option value="sliders">Sliders</option>
        <option value="kontak">Kontak</option>
        <option value="settings">Settings</option>
        <option value="database">Database</option>
      </select>

      <select
        value={filterRisk}
        onChange={(e) => {
          setFilterRisk(e.target.value)
          setPage(1)
        }}
        className="shrink-0 bg-gray-50 border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-brand-500 focus:bg-white transition-colors cursor-pointer"
      >
        <option value="">Semua Risiko</option>
        <option value="high">High</option>
        <option value="medium">Medium</option>
        <option value="low">Low</option>
      </select>

      <button
        onClick={handleReset}
        className="shrink-0 bg-gray-50 border border-gray-200 text-gray-700 hover:bg-gray-100 px-3 py-2 rounded-lg text-sm font-medium flex items-center gap-2 transition-all btn-press"
      >
        <i className="ph ph-arrows-counter-clockwise text-lg" /> Reset
      </button>
    </div>
  )

  return (
    <AdminLayout title="Catatan Aktivitas Sistem" headerContent={headerContent}>
      <div className="max-w-7xl mx-auto space-y-6 animate-fade-in-up">

        {/* Table Card */}
        <div className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
          <div className="p-4 border-b border-gray-200 flex justify-end items-center bg-gray-50">
            <button
              onClick={downloadCSV}
              disabled={logs.length === 0}
              className="text-sm text-gray-500 hover:text-gray-800 bg-white border border-gray-300 px-3 py-1.5 rounded flex items-center gap-2 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <i className="ph ph-download-simple" /> Download CSV
            </button>
          </div>
          
          <div className="overflow-x-auto">
            {loading ? (
              <div className="py-12 text-center text-gray-500">Memuat log aktivitas...</div>
            ) : error ? (
              <div className="py-12 text-center text-red-500">{error}</div>
            ) : logs.length === 0 ? (
              <div className="py-16 text-center text-slate-500 flex flex-col items-center justify-center">
                <i className="ph ph-clock-counter-clockwise text-gray-300 text-5xl mb-4" />
                <p className="font-medium text-gray-500">Tidak ada log aktivitas untuk ditampilkan.</p>
              </div>
            ) : (
              <table className="w-full text-left border-collapse">
                <thead>
                  <tr className="bg-white border-b border-gray-200 text-xs uppercase tracking-wider text-gray-500 font-semibold">
                    <th className="p-4 w-40">Waktu</th>
                    <th className="p-4">Pengguna (Aktor)</th>
                    <th className="p-4">Aktivitas & Entitas</th>
                    <th className="p-4">Tingkat Risiko</th>
                    <th className="p-4">Keterangan Tambahan</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-100 text-sm text-gray-700">
                  {logs.map((log) => {
                    const row = mapLogRowForDisplay(log)

                    return (
                      <tr key={row.id} className="hover:bg-gray-50 transition-colors admin-row">
                        {/* WAKTU */}
                        <td className="p-4">
                          <div className="text-gray-900 font-medium">{row.datePart}</div>
                          <div className="text-gray-500 text-xs mt-0.5">{row.hourPart}</div>
                          <div className="text-[10px] text-gray-400 mt-1 uppercase">Log ID: {row.id}</div>
                        </td>
                        
                        {/* AKTOR & DEVICE */}
                        <td className="p-4">
                          <div className="flex items-start gap-3">
                            <div className="w-8 h-8 rounded-full bg-brand-100 text-brand-600 flex items-center justify-center font-bold text-xs shrink-0">
                              {row.actor.charAt(0).toUpperCase()}
                            </div>
                            <div>
                              <div className="font-medium text-gray-900 leading-tight">
                                {row.actor}
                              </div>
                              <div className="text-[11px] text-gray-500 mt-0.5">
                                {row.role} • <span className="font-mono text-gray-400">{row.ip}</span>
                              </div>
                              <div className="text-[10px] text-gray-400 mt-1 leading-tight max-w-[200px]" title={row.device}>
                                {row.device}
                              </div>
                            </div>
                          </div>
                        </td>
                        
                        {/* ACTION & ENTITY */}
                        <td className="p-4">
                          <div className="flex flex-col gap-1.5 items-start">
                            <span className={`inline-flex items-center px-1.5 py-0.5 rounded text-[11px] font-bold font-mono tracking-wider border ${getActionBadgeClass(row.action)}`}>
                              {row.action}
                            </span>
                            
                            <div className="bg-gray-100 border border-gray-200 px-2 py-1 rounded text-xs w-full max-w-[220px]">
                              <span className="text-gray-500 font-medium capitalize">{row.entity}</span> 
                              <div className="font-medium text-gray-800 mt-0.5 truncate" title={row.entityLabel}>
                                {row.entityLabel}
                              </div>
                            </div>
                          </div>
                        </td>
                        
                        {/* RISIKO */}
                        <td className="p-4">
                          <span className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium border ${getRiskBadgeClass(row.risk)}`}>
                            <i className={`ph ${row.risk?.toLowerCase() === 'high' ? 'ph-warning-circle' : 'ph-info'}`} />
                            <span className="capitalize">{row.risk}</span>
                          </span>
                        </td>
                        
                        {/* KETERANGAN */}
                        <td className="p-4 text-gray-600 text-xs leading-relaxed max-w-[200px]">
                          {row.description}
                        </td>
                      </tr>
                    )
                  })}
                </tbody>
              </table>
            )}
          </div>

          {/* Pagination */}
          {!loading && !error && logs.length > 0 && (
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
    </AdminLayout>
  )
}
