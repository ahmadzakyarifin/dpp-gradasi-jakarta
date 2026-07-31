import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import AdminLayout from '../layouts/AdminLayout'
import { dashboardService } from '../services/dashboardService'

export default function Dashboard() {
  const [summary, setSummary] = useState({
    total_berita: 3,
    total_kegiatan: 3,
    total_pengurus: 14,
    total_kontak: 2,
    latest_logs: [
      { id: 1, action: 'CREATE', module: 'BERITA', details: 'Menambahkan Berita Rapat Kerja Daerah Jatim', created_at: '2026-07-31 16:00:00' },
      { id: 2, action: 'UPDATE', module: 'SETTINGS', details: 'Memperbarui Pengaturan Website', created_at: '2026-07-31 15:30:00' }
    ]
  })
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    dashboardService.getSummary()
      .then(res => {
        if (res.data) {
          setSummary(res.data)
        }
      })
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [])

  const stats = [
    { label: 'Total Berita', value: summary.total_berita ?? 3, color: 'bg-blue-50 text-blue-600', icon: 'ph-article', link: '/admin/berita' },
    { label: 'Total Kegiatan', value: summary.total_kegiatan ?? 3, color: 'bg-emerald-50 text-emerald-600', icon: 'ph-calendar-check', link: '/admin/kegiatan' },
    { label: 'Total Pengurus', value: summary.total_pengurus ?? 14, color: 'bg-indigo-50 text-indigo-600', icon: 'ph-users-three', link: '/admin/pengurus' },
    { label: 'Pesan Kontak Baru', value: summary.total_kontak ?? 2, color: 'bg-amber-50 text-amber-600', icon: 'ph-envelope-simple', link: '/admin/kontak' },
  ]

  return (
    <AdminLayout title="Dashboard Overview">
      <div className="space-y-8">
        
        {/* Welcome Banner */}
        <div className="bg-gradient-to-r from-brand-900 to-brand-700 rounded-2xl p-8 text-white shadow-lg flex flex-col md:flex-row justify-between items-start md:items-center gap-6">
          <div>
            <h2 className="font-heading text-2xl font-bold mb-2">Selamat Datang di Panel Admin DPP GRADASI</h2>
            <p className="text-brand-100 text-sm max-w-xl">
              Kelola seluruh informasi berita, kegiatan, susunan pengurus, dan pengaturan website secara terpusat dari satu dashboard.
            </p>
          </div>
          <Link
            to="/admin/settings"
            className="bg-white/10 hover:bg-white/20 text-white border border-white/20 px-5 py-2.5 rounded-xl font-semibold text-sm transition backdrop-blur-sm shrink-0 flex items-center gap-2"
          >
            <i className="ph ph-gear text-lg" /> Pengaturan Website
          </Link>
        </div>

        {/* Stats Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {stats.map((stat) => (
            <Link
              key={stat.label}
              to={stat.link}
              className="bg-white p-6 rounded-2xl border border-gray-100 shadow-sm hover:shadow-md transition group flex flex-col justify-between"
            >
              <div className="flex items-center justify-between mb-4">
                <span className="text-gray-500 font-semibold text-xs tracking-wider uppercase">{stat.label}</span>
                <div className={`w-12 h-12 rounded-2xl ${stat.color} flex items-center justify-center text-2xl group-hover:scale-110 transition-transform`}>
                  <i className={`ph ${stat.icon}`} />
                </div>
              </div>
              <div>
                <h3 className="text-3xl font-heading font-extrabold text-gray-900">{loading ? '...' : stat.value}</h3>
              </div>
            </Link>
          ))}
        </div>

        {/* Quick Links & Activity Log Summary */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          
          {/* Quick Actions */}
          <div className="bg-white p-6 rounded-2xl border border-gray-100 shadow-sm space-y-4">
            <h3 className="font-heading font-bold text-gray-900 text-lg flex items-center gap-2">
              <i className="ph ph-lightning text-brand-600" /> Akses Cepat Admin
            </h3>
            <div className="space-y-2">
              <Link to="/admin/berita" className="flex items-center justify-between p-3 rounded-xl bg-gray-50 hover:bg-brand-50 hover:text-brand-600 text-gray-700 text-sm font-semibold transition">
                <span className="flex items-center gap-2.5"><i className="ph ph-plus-circle text-lg" /> Tambah Berita Baru</span>
                <i className="ph ph-caret-right" />
              </Link>
              <Link to="/admin/kegiatan" className="flex items-center justify-between p-3 rounded-xl bg-gray-50 hover:bg-brand-50 hover:text-brand-600 text-gray-700 text-sm font-semibold transition">
                <span className="flex items-center gap-2.5"><i className="ph ph-calendar-plus text-lg" /> Tambah Kegiatan / Event</span>
                <i className="ph ph-caret-right" />
              </Link>
              <Link to="/admin/pengurus" className="flex items-center justify-between p-3 rounded-xl bg-gray-50 hover:bg-brand-50 hover:text-brand-600 text-gray-700 text-sm font-semibold transition">
                <span className="flex items-center gap-2.5"><i className="ph ph-user-plus text-lg" /> Tambah Data Pengurus</span>
                <i className="ph ph-caret-right" />
              </Link>
            </div>
          </div>

          {/* Activity Logs */}
          <div className="lg:col-span-2 bg-white p-6 rounded-2xl border border-gray-100 shadow-sm space-y-4">
            <div className="flex items-center justify-between">
              <h3 className="font-heading font-bold text-gray-900 text-lg flex items-center gap-2">
                <i className="ph ph-clock-counter-clockwise text-brand-600" /> Log Aktivitas Terakhir
              </h3>
              <Link to="/admin/activity-log" className="text-xs font-bold text-brand-600 hover:text-brand-700">
                Lihat Semua
              </Link>
            </div>
            
            <div className="divide-y divide-gray-100">
              {(summary.latest_logs || []).slice(0, 4).map((log) => (
                <div key={log.id} className="py-3 flex items-center justify-between text-sm">
                  <div className="flex items-center gap-3">
                    <span className="px-2.5 py-1 rounded-lg bg-blue-50 text-blue-700 font-mono text-[11px] font-bold">
                      {log.action}
                    </span>
                    <span className="text-gray-700 font-medium">{log.details || log.module}</span>
                  </div>
                  <span className="text-xs text-gray-400 font-mono">{log.created_at}</span>
                </div>
              ))}
            </div>
          </div>

        </div>

      </div>
    </AdminLayout>
  )
}
