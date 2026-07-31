import { useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuthStore } from '../store/useAuthStore'
import { useDashboardStore } from '../store/useDashboardStore'

export default function Dashboard() {
  const { token, logout, user } = useAuthStore()
  const { summary, fetchSummary, loading } = useDashboardStore()
  const navigate = useNavigate()

  useEffect(() => {
    if (!token) {
      navigate('/login')
      return
    }
    fetchSummary()
  }, [token, navigate, fetchSummary])

  if (loading || !summary) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-gray-500">Memuat data...</div>
      </div>
    )
  }

  const stats = [
    { label: 'Total Berita', value: summary.total_berita, color: 'text-blue-600', icon: 'ph-article' },
    { label: 'Total Kegiatan', value: summary.total_kegiatan, color: 'text-indigo-600', icon: 'ph-calendar-check' },
    { label: 'Total Pengurus', value: summary.total_pengurus, color: 'text-orange-600', icon: 'ph-users-three' },
    { label: 'Pesan Kontak Baru', value: summary.total_kontak, color: 'text-rose-600', icon: 'ph-envelope-simple' },
  ]

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col">
      <header className="bg-white shadow-sm border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
          <h1 className="text-xl font-bold text-gray-900 font-heading">GRADASI Admin</h1>
          <div className="flex items-center gap-4">
            <span className="text-sm text-gray-700">Halo, {user?.name || 'Admin'}</span>
            <button
              onClick={() => {
                logout()
                navigate('/login')
              }}
              className="text-sm font-semibold text-red-600 hover:text-red-500"
            >
              Keluar
            </button>
          </div>
        </div>
      </header>
      <main className="flex-1 max-w-7xl w-full mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <h2 className="text-2xl font-bold text-gray-900 font-heading mb-6">Dashboard Ringkasan</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {stats.map((stat) => (
            <div key={stat.label} className="bg-white p-6 rounded-xl border border-gray-200 shadow-sm flex flex-col">
              <div className={`w-12 h-12 bg-gray-50 ${stat.color} rounded-full flex items-center justify-center mb-4`}>
                <i className={`ph ${stat.icon} text-2xl`}></i>
              </div>
              <h3 className="text-3xl font-heading font-bold text-gray-900 mb-1">{stat.value}</h3>
              <p className="text-gray-500 font-medium text-sm">{stat.label}</p>
            </div>
          ))}
        </div>
      </main>
    </div>
  )
}
