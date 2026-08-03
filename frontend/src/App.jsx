import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom'
import { lazy, Suspense, useEffect } from 'react'

// Code-splitting: halaman dimuat per-route (React.lazy) supaya bundle awal
// kecil — halaman admin tidak ikut dimuat saat user di halaman publik.
const Login = lazy(() => import('./pages/Login'))
const ResetPassword = lazy(() => import('./pages/ResetPassword'))
const Dashboard = lazy(() => import('./pages/Dashboard'))
const BeritaList = lazy(() => import('./pages/BeritaList'))
const BeritaDetail = lazy(() => import('./pages/BeritaDetail'))
const KegiatanList = lazy(() => import('./pages/KegiatanList'))
const KegiatanDetail = lazy(() => import('./pages/KegiatanDetail'))
const Kepengurusan = lazy(() => import('./pages/Kepengurusan'))
const Home = lazy(() => import('./pages/Home'))

// Admin Pages
const ProfileAdmin = lazy(() => import('./pages/admin/ProfileAdmin'))
const SlidersAdmin = lazy(() => import('./pages/admin/SlidersAdmin'))
const SettingsAdmin = lazy(() => import('./pages/admin/SettingsAdmin'))
const BeritaAdmin = lazy(() => import('./pages/admin/BeritaAdmin'))
const KegiatanAdmin = lazy(() => import('./pages/admin/KegiatanAdmin'))
const PengurusAdmin = lazy(() => import('./pages/admin/PengurusAdmin'))
const KontakAdmin = lazy(() => import('./pages/admin/KontakAdmin'))
const UsersAdmin = lazy(() => import('./pages/admin/UsersAdmin'))
const RoleAdmin = lazy(() => import('./pages/admin/RoleAdmin'))
const ActivityLogAdmin = lazy(() => import('./pages/admin/ActivityLogAdmin'))

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

function PageLoader() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-white">
      <div className="flex flex-col items-center gap-3">
        <div className="w-8 h-8 border-4 border-brand-200 border-t-brand-600 rounded-full animate-spin" />
        <p className="text-xs text-slate-400 font-medium">Memuat...</p>
      </div>
    </div>
  )
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
        <Suspense fallback={<PageLoader />}>
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
        </Suspense>
      </Router>
    </SettingsProvider>
  )
}

export default App
