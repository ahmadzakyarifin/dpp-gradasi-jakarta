import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import AdminLayout from '../layouts/AdminLayout'
import { useSettings } from '../context/useSettings'
import { dashboardService } from '../services/dashboardService'

export default function Dashboard() {
  const { settings } = useSettings()
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
            <h2 className="font-heading text-2xl font-bold mb-2">Selamat Datang di Panel Admin {settings.site_name}</h2>
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
              className="card-lift bg-white p-6 rounded-2xl border border-gray-100 shadow-sm hover:shadow-md transition group flex flex-col justify-between"
            >
              <div className="flex items-center justify-between mb-4">
                <span className="text-gray-500 font-semibold text-xs tracking-wider uppercase">{stat.label}</span>
                <div className={`w-12 h-12 rounded-full ${stat.color} flex items-center justify-center text-2xl group-hover:scale-110 group-hover:rotate-3 transition-transform shadow-sm`}>
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
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 items-stretch">
          
          {/* Quick Actions */}
          <div className="bg-white p-6 rounded-2xl border border-gray-100 shadow-sm flex flex-col justify-between">
            <div className="space-y-5">
              <h3 className="font-heading font-bold text-gray-900 text-lg flex items-center gap-2">
                <i className="ph ph-lightning text-brand-600" /> Akses Cepat Admin
              </h3>
              <div className="space-y-3">
                <Link to="/admin/berita" className="flex items-center justify-between p-4 rounded-xl bg-gray-50 hover:bg-brand-50 hover:text-brand-600 text-gray-700 text-sm font-semibold transition-all hover:translate-x-1 duration-200">
                  <span className="flex items-center gap-3"><i className="ph ph-plus-circle text-xl text-brand-500" /> Tambah Berita Baru</span>
                  <i className="ph ph-caret-right text-gray-400" />
                </Link>
                <Link to="/admin/kegiatan" className="flex items-center justify-between p-4 rounded-xl bg-gray-50 hover:bg-brand-50 hover:text-brand-600 text-gray-700 text-sm font-semibold transition-all hover:translate-x-1 duration-200">
                  <span className="flex items-center gap-3"><i className="ph ph-calendar-plus text-xl text-emerald-500" /> Tambah Kegiatan / Event</span>
                  <i className="ph ph-caret-right text-gray-400" />
                </Link>
                <Link to="/admin/pengurus" className="flex items-center justify-between p-4 rounded-xl bg-gray-50 hover:bg-brand-50 hover:text-brand-600 text-gray-700 text-sm font-semibold transition-all hover:translate-x-1 duration-200">
                  <span className="flex items-center gap-3"><i className="ph ph-user-plus text-xl text-indigo-500" /> Tambah Data Pengurus</span>
                  <i className="ph ph-caret-right text-gray-400" />
                </Link>
              </div>
            </div>
          </div>

          {/* Latest Messages */}
          <div className="lg:col-span-2 bg-white p-6 rounded-2xl border border-gray-100 shadow-sm flex flex-col justify-between">
            <div className="space-y-4">
              <div className="flex items-center justify-between mb-2">
                <h3 className="font-heading font-bold text-gray-900 text-lg flex items-center gap-2">
                  <i className="ph ph-envelope-simple text-brand-600" /> Pesan Masuk Terbaru
                </h3>
                <Link to="/admin/kontak" className="text-xs font-bold text-brand-600 hover:text-brand-700 bg-brand-50 hover:bg-brand-100 px-3 py-1.5 rounded-lg transition-colors">
                  Lihat Semua
                </Link>
              </div>
              
              <div className="divide-y divide-gray-100">
                {(summary.latest_messages || []).map((msg) => (
                  <div key={msg.id} className="py-3.5 flex items-center justify-between text-sm hover:bg-gray-50/50 px-2 rounded-xl transition-colors">
                    <div className="flex items-center gap-3">
                      {msg.is_read ? (
                        <span className="w-2.5 h-2.5 rounded-full bg-gray-300 shrink-0" title="Sudah Dibaca" />
                      ) : (
                        <span className="w-2.5 h-2.5 rounded-full bg-brand-500 animate-pulse shrink-0" title="Pesan Baru" />
                      )}
                      <div className="flex flex-col">
                        <span className="text-gray-900 font-semibold text-xs leading-tight">{msg.nama}</span>
                        <span className="text-gray-500 font-medium text-xs mt-0.5 line-clamp-1 max-w-[280px] md:max-w-[420px]" title={msg.subjek}>{msg.subjek}</span>
                      </div>
                    </div>
                    <span className="text-xs text-gray-400 font-mono hidden sm:inline">{msg.created_at}</span>
                  </div>
                ))}
                {(!summary.latest_messages || summary.latest_messages.length === 0) && (
                  <div className="py-12 text-center text-gray-400 text-sm flex flex-col items-center justify-center">
                    <i className="ph ph-envelope-open text-3xl text-gray-300 mb-2" />
                    Tidak ada pesan masuk terbaru
                  </div>
                )}
              </div>
            </div>
          </div>

        </div>

      </div>
    </AdminLayout>
  )
}
