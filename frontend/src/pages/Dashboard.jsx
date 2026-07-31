import { useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuthStore } from '../store/useAuthStore'

export default function Dashboard() {
  const { token, logout, user, setUser } = useAuthStore()
  const navigate = useNavigate()

  useEffect(() => {
    if (!token) {
      navigate('/login')
      return
    }

    fetch('http://127.0.0.1:8080/api/v1/profile', {
      headers: { Authorization: `Bearer ${token}` }
    })
      .then((r) => r.json())
      .then((res) => {
        if (res.success) {
          setUser(res.data)
        } else {
          logout()
          navigate('/login')
        }
      })
      .catch(() => {
        logout()
        navigate('/login')
      })
  }, [token, navigate, logout, setUser])

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
        <div className="bg-white p-6 rounded-2xl shadow-sm border border-gray-200">
          <h2 className="text-2xl font-bold text-gray-900 font-heading mb-4">Selamat datang di Panel Admin!</h2>
          <p className="text-gray-600">Proyek React + Tailwind + Zustand + React Router ini siap Anda kembangkan.</p>
        </div>
      </main>
    </div>
  )
}
