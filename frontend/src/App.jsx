import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom'
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
import ActivityLogAdmin from './pages/admin/ActivityLogAdmin'

import { useAuthStore } from './store/useAuthStore'

function App() {
  const token = useAuthStore((state) => state.token)

  return (
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
        <Route path="/dashboard" element={token ? <Dashboard /> : <Navigate to="/login" />} />
        
        {/* Admin Dashboard Pages */}
        <Route path="/admin/profile" element={token ? <ProfileAdmin /> : <Navigate to="/login" />} />
        <Route path="/admin/sliders" element={token ? <SlidersAdmin /> : <Navigate to="/login" />} />
        <Route path="/admin/settings" element={token ? <SettingsAdmin /> : <Navigate to="/login" />} />
        <Route path="/admin/berita" element={token ? <BeritaAdmin /> : <Navigate to="/login" />} />
        <Route path="/admin/kegiatan" element={token ? <KegiatanAdmin /> : <Navigate to="/login" />} />
        <Route path="/admin/pengurus" element={token ? <PengurusAdmin /> : <Navigate to="/login" />} />
        <Route path="/admin/kontak" element={token ? <KontakAdmin /> : <Navigate to="/login" />} />
        <Route path="/admin/users" element={token ? <UsersAdmin /> : <Navigate to="/login" />} />
        <Route path="/admin/activity-log" element={token ? <ActivityLogAdmin /> : <Navigate to="/login" />} />
        
        {/* Root Fallback */}
        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
    </Router>
  )
}

export default App
