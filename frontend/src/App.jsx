import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom'
import Login from './pages/Login'
import ResetPassword from './pages/ResetPassword'
import Dashboard from './pages/Dashboard'
import BeritaList from './pages/BeritaList'
import BeritaDetail from './pages/BeritaDetail'
import { useAuthStore } from './store/useAuthStore'

function App() {
  const token = useAuthStore((state) => state.token)

  return (
    <Router>
      <Routes>
        {/* Auth */}
        <Route path="/login" element={!token ? <Login /> : <Navigate to="/dashboard" />} />
        <Route path="/reset-password" element={<ResetPassword />} />
        
        {/* Public Berita */}
        <Route path="/berita" element={<BeritaList />} />
        <Route path="/berita/:slug" element={<BeritaDetail />} />

        {/* Dashboard */}
        <Route path="/dashboard" element={token ? <Dashboard /> : <Navigate to="/login" />} />
        
        {/* Root */}
        <Route path="*" element={<Navigate to={token ? "/dashboard" : "/login"} />} />
      </Routes>
    </Router>
  )
}

export default App
