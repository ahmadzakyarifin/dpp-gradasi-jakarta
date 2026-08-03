import { Navigate, useLocation } from 'react-router-dom'
import { useEffect } from 'react'
import { useAuthStore } from '../store/useAuthStore'

/**
 * Proteksi route berdasarkan:
 * 1. token (belum login → redirect ke /login)
 * 2. role (tidak punya akses → redirect ke /dashboard)
 *
 * Saat token ada tapi role belum dimuat (baru refresh halaman), komponen ini
 * memanggil fetchMe() dulu dan menampilkan loading — TIDAK menampilkan konten
 * maupun redirect, supaya tidak terjadi flash menu sidebar yang salah.
 */
export default function ProtectedRoute({ children, allowedRoles = null }) {
  const token = useAuthStore((state) => state.token)
  const role = useAuthStore((state) => state.role)
  const fetchMe = useAuthStore((state) => state.fetchMe)
  const location = useLocation()

  useEffect(() => {
    // Token ada tapi user/role belum terisi (refresh/reload) → ambil profil
    if (token && !role) {
      fetchMe()
    }
  }, [token, role, fetchMe])

  if (!token) {
    return <Navigate to="/login" replace state={{ from: location }} />
  }

  // Masih memuat profil (role belum diketahui) → jangan render konten dulu
  if (!role) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="flex flex-col items-center gap-3">
          <div className="w-10 h-10 border-4 border-brand-600 border-t-transparent rounded-full animate-spin" />
          <p className="text-sm text-gray-500 font-medium">Memuat profil...</p>
        </div>
      </div>
    )
  }

  // Role tidak diizinkan untuk route ini
  if (allowedRoles && !allowedRoles.includes(role)) {
    return <Navigate to="/dashboard" replace />
  }

  return children
}
