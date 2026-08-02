import { Link, useNavigate, useLocation } from 'react-router-dom'
import { useState } from 'react'
import { useAuthStore } from '../store/useAuthStore'

const sidebarLinks = [
  { path: '/dashboard', label: 'Dashboard', icon: 'ph-squares-four' },
  { path: '/admin/berita', label: 'Berita', icon: 'ph-article', roles: ['super_admin', 'admin', 'admin_berita'] },
  { path: '/admin/kegiatan', label: 'Kegiatan', icon: 'ph-calendar-check', roles: ['super_admin', 'admin', 'editor'] },
  { path: '/admin/pengurus', label: 'Pengurus', icon: 'ph-users-three' },
  { path: '/admin/sliders', label: 'Sliders', icon: 'ph-image' },
  { path: '/admin/kontak', label: 'Pesan Kontak', icon: 'ph-envelope-simple' },
  { path: '/admin/users', label: 'Manajemen Admin', icon: 'ph-user-gear', roles: ['super_admin'] },
  { path: '/admin/activity-log', label: 'Activity Log', icon: 'ph-clock-counter-clockwise', roles: ['super_admin'] },
  { path: '/admin/settings', label: 'Pengaturan Website', icon: 'ph-gear' },
]

export default function AdminLayout({ children, title = 'Admin Panel', headerContent }) {
  const navigate = useNavigate()
  const location = useLocation()
  const logout = useAuthStore((state) => state.logout)
  const role = useAuthStore((state) => state.role)
  const [userMenuOpen, setUserMenuOpen] = useState(false)

  const visibleLinks = role ? sidebarLinks.filter((l) => !l.roles || l.roles.includes(role)) : sidebarLinks

  function handleLogout() {
    logout()
    navigate('/login')
  }

  return (
    <div className="font-sans antialiased overflow-hidden flex h-screen bg-gray-50 text-gray-800">
      {/* Sidebar */}
      <aside className="w-64 bg-brand-900 text-white flex flex-col h-full shrink-0 shadow-lg z-10">
        <div className="h-16 flex items-center px-6 border-b border-white/10 shrink-0">
          <span className="font-heading font-bold text-xl tracking-wide">
            GRADASI<span className="text-brand-400">Admin</span>
          </span>
        </div>
        <nav className="flex-1 overflow-y-auto py-4 space-y-1 px-3">
          {visibleLinks.map((link) => {
            const isActive = location.pathname === link.path
            return (
              <Link
                key={link.path}
                to={link.path}
                className={`flex items-center px-3 py-2.5 rounded-lg transition-all duration-200 btn-press ${
                  isActive
                    ? 'bg-brand-600 text-white shadow-md font-medium'
                    : 'text-white/70 hover:bg-white/10 hover:text-white'
                }`}
              >
                <i className={`ph ${link.icon} text-xl mr-3`} /> {link.label}
              </Link>
            )
          })}
        </nav>
        <div className="p-4 border-t border-white/10 shrink-0">
          <button
            onClick={handleLogout}
            className="w-full flex items-center px-3 py-2.5 rounded-lg text-red-300 hover:bg-red-500/20 transition-all duration-200 btn-press font-medium"
          >
            <i className="ph ph-sign-out text-xl mr-3" /> Logout
          </button>
        </div>
      </aside>

      {/* Main Content */}
      <main className="flex-1 flex flex-col h-full overflow-hidden">
        {/* Header */}
        <header className="h-16 bg-white border-b border-gray-200 flex items-center justify-between px-8 shrink-0 shadow-sm z-10">
          <h1 className="font-heading font-semibold text-gray-800 text-lg shrink-0">{title}</h1>
          
          {/* Center Actions (Search, Filter, etc.) */}
          {headerContent && (
            <div className="flex-1 flex items-center justify-center px-6">
              {headerContent}
            </div>
          )}

          <div className="flex items-center gap-4 shrink-0">
            <div className="relative">
              <button
                onClick={() => setUserMenuOpen(!userMenuOpen)}
                className="flex items-center gap-2 hover:bg-gray-100 p-2 rounded-lg transition-all"
              >
                <img
                  src="https://ui-avatars.com/api/?name=Admin&background=0D8ABC&color=fff"
                  className="w-8 h-8 rounded-full border border-gray-200"
                  alt="Admin"
                />
                <span className="text-sm font-medium text-gray-700">Admin</span>
                <i className="ph ph-caret-down text-gray-500" />
              </button>
              {userMenuOpen && (
                <div className="absolute right-0 mt-2 w-48 bg-white rounded-xl shadow-lg border border-gray-100 py-1 z-50 animate-scale-in">
                  <Link
                    to="/admin/profile"
                    onClick={() => setUserMenuOpen(false)}
                    className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 hover:text-brand-600 transition-colors"
                  >
                    <i className="ph ph-user-circle mr-2" /> Profil Saya
                  </Link>
                  <div className="border-t border-gray-100 my-1" />
                  <button
                    onClick={handleLogout}
                    className="w-full text-left block px-4 py-2 text-sm text-red-600 hover:bg-red-50 transition-colors"
                  >
                    <i className="ph ph-sign-out mr-2" /> Logout
                  </button>
                </div>
              )}
            </div>
          </div>
        </header>

        {/* Content Area */}
        <div className="flex-1 overflow-auto bg-gray-50 p-8">{children}</div>
      </main>
    </div>
  )
}
