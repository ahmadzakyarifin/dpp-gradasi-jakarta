import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom'
import { useEffect } from 'react'
import Login from './pages/Login'
import ResetPassword from './pages/ResetPassword'
import Dashboard from './pages/Dashboard'
import BeritaList from './pages/BeritaList'
import BeritaDetail from './pages/BeritaDetail'
import KegiatanList from './pages/KegiatanList'
import KegiatanDetail from './pages/KegiatanDetail'
import Kepengurusan from './pages/Kepengurusan'
import Home from './pages/Home'

// Admin Pages
import ProfileAdmin from './pages/admin/ProfileAdmin'
import SlidersAdmin from './pages/admin/SlidersAdmin'
import SettingsAdmin from './pages/admin/SettingsAdmin'
import BeritaAdmin from './pages/admin/BeritaAdmin'
import KegiatanAdmin from './pages/admin/KegiatanAdmin'
import PengurusAdmin from './pages/admin/PengurusAdmin'
import KontakAdmin from './pages/admin/KontakAdmin'
import UsersAdmin from './pages/admin/UsersAdmin'
import RoleAdmin from './pages/admin/RoleAdmin'
import ActivityLogAdmin from './pages/admin/ActivityLogAdmin'

import { useAuthStore } from './store/useAuthStore'
import { SettingsProvider } from './context/SettingsContext'
import ProtectedRoute from './components/ProtectedRoute'

// Mapping role per route — cocokkan dengan middleware RoleMiddleware di backend:
//   dashboard & activitylog & user : super_admin
//   berita  : super_admin, admin, admin_berita
//   kegiatan: super_admin, admin, editor
//   settings/sliders/kontak/pengurus : super_admin, admin
//   role (manajemen role) : super_admin, admin
//   profile : semua role yang sudah login
const ROLES = {
  DASHBOARD: ['super_admin'],
  BERITA: ['super_admin', 'admin', 'admin_berita'],
  KEGIATAN: ['super_admin', 'admin', 'editor'],
  PENGURUS: ['super_admin', 'admin'],
  SLIDERS: ['super_admin', 'admin'],
  KONTAK: ['super_admin', 'admin'],
  USERS: ['super_admin'],
  ACTIVITY_LOG: ['super_admin'],
  SETTINGS: ['super_admin', 'admin'],
}

function App() {
  const token = useAuthStore((state) => state.token)
  const role = useAuthStore((state) => state.role)
  const fetchMe = useAuthStore((state) => state.fetchMe)

  // Global: saat token ada tapi role belum terisi (refresh halaman admin),
  // ambil profil user supaya role tersedia sebelum sidebar/rute dirender.
  useEffect(() => {
    if (token && !role) {
      fetchMe()
    }
  }, [token, role, fetchMe])

  return (
    <SettingsProvider>
      <Router>
        <Routes>
          {/* Public Homepage */}
          <Route path="/" element={<Home />} />

          {/* Auth */}
          <Route path="/login" element={!token ? <Login /> : <Navigate to="/dashboard" />} />
          <Route path="/reset-password" element={<ResetPassword />} />

          {/* Public Berita */}
          <Route path="/berita" element={<BeritaList />} />
          <Route path="/berita/:slug" element={<BeritaDetail />} />

          {/* Public Kegiatan */}
          <Route path="/kegiatan" element={<KegiatanList />} />
          <Route path="/kegiatan/:slug" element={<KegiatanDetail />} />

          {/* Public Kepengurusan */}
          <Route path="/kepengurusan" element={<Kepengurusan />} />

          {/* Dashboard */}
          <Route
            path="/dashboard"
            element={
              <ProtectedRoute allowedRoles={ROLES.DASHBOARD}>
                <Dashboard />
              </ProtectedRoute>
            }
          />

          {/* Admin Dashboard Pages */}
          <Route
            path="/admin/profile"
            element={
              <ProtectedRoute>
                <ProfileAdmin />
              </ProtectedRoute>
            }
          />
          <Route
            path="/admin/sliders"
            element={
              <ProtectedRoute allowedRoles={ROLES.SLIDERS}>
                <SlidersAdmin />
              </ProtectedRoute>
            }
          />
          <Route
            path="/admin/settings"
            element={
              <ProtectedRoute allowedRoles={ROLES.SETTINGS}>
                <SettingsAdmin />
              </ProtectedRoute>
            }
          />
          <Route
            path="/admin/berita"
            element={
              <ProtectedRoute allowedRoles={ROLES.BERITA}>
                <BeritaAdmin />
              </ProtectedRoute>
            }
          />
          <Route
            path="/admin/kegiatan"
            element={
              <ProtectedRoute allowedRoles={ROLES.KEGIATAN}>
                <KegiatanAdmin />
              </ProtectedRoute>
            }
          />
          <Route
            path="/admin/pengurus"
            element={
              <ProtectedRoute allowedRoles={ROLES.PENGURUS}>
                <PengurusAdmin />
              </ProtectedRoute>
            }
          />
          <Route
            path="/admin/kontak"
            element={
              <ProtectedRoute allowedRoles={ROLES.KONTAK}>
                <KontakAdmin />
              </ProtectedRoute>
            }
          />
          <Route
            path="/admin/users"
            element={
              <ProtectedRoute allowedRoles={ROLES.USERS}>
                <UsersAdmin />
              </ProtectedRoute>
            }
          />
          <Route
            path="/admin/roles"
            element={
              <ProtectedRoute allowedRoles={ROLES.USERS}>
                <RoleAdmin />
              </ProtectedRoute>
            }
          />
          <Route
            path="/admin/activity-log"
            element={
              <ProtectedRoute allowedRoles={ROLES.ACTIVITY_LOG}>
                <ActivityLogAdmin />
              </ProtectedRoute>
            }
          />

          {/* Root Fallback */}
          <Route path="*" element={<Navigate to="/" />} />
        </Routes>
      </Router>
    </SettingsProvider>
  )
}

export default App
