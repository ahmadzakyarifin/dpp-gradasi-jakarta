import { Link, useNavigate, useLocation } from 'react-router-dom'
import { useState } from 'react'
import { useAuthStore } from '../store/useAuthStore'
import { resolveAssetUrl } from '../utils/assetUrl'

const sidebarLinks = [
  { path: '/dashboard', label: 'Dashboard', icon: 'ph-squares-four' },
  { path: '/admin/berita', label: 'Berita', icon: 'ph-article', roles: ['super_admin', 'admin'] },
  { path: '/admin/kegiatan', label: 'Kegiatan', icon: 'ph-calendar-check', roles: ['super_admin', 'admin'] },
  { path: '/admin/pengurus', label: 'Pengurus', icon: 'ph-users-three' },
  { path: '/admin/sliders', label: 'Sliders', icon: 'ph-image' },
  { path: '/admin/kontak', label: 'Pesan Kontak', icon: 'ph-envelope-simple' },
  { path: '/admin/users', label: 'Manajemen Admin', icon: 'ph-user-gear', roles: ['super_admin'] },
  { path: '/admin/activity-log', label: 'Activity Log', icon: 'ph-clock-counter-clockwise', roles: ['super_admin'] },
  { path: '/admin/settings', label: 'Pengaturan Website', icon: 'ph-gear', roles: ['super_admin'] },
]

export default function AdminLayout({ children, title = 'Admin Panel', headerContent }) {
  const navigate = useNavigate()
  const location = useLocation()
  const logout = useAuthStore((state) => state.logout)
  const user = useAuthStore((state) => state.user)
  const role = useAuthStore((state) => state.role)
  const [userMenuOpen, setUserMenuOpen] = useState(false)
  const [isSidebarOpen, setIsSidebarOpen] = useState(false)
  const [previewImage, setPreviewImage] = useState(null)

  // Selama role belum diketahui (baru refresh, fetchMe masih jalan),
  // jangan fallback ke "tampilkan semua menu" — tampilkan skeleton.
  if (!role) {
    return (
      <div className="font-sans antialiased overflow-hidden flex h-screen bg-slate-50/50 text-gray-800">
        <aside className="w-64 bg-slate-900 text-white flex flex-col h-full shrink-0 shadow-xl z-10">
          <div className="h-16 flex items-center px-6 border-b border-white/5 shrink-0">
            <span className="font-heading font-bold text-xl tracking-wide">
              GRADASI<span className="text-brand-400">Admin</span>
            </span>
          </div>
          <nav className="flex-1 overflow-y-auto py-4 space-y-2 px-3">
            {[...Array(6)].map((_, i) => (
              <div key={i} className="h-10 rounded-lg bg-white/5 animate-pulse" />
            ))}
          </nav>
          <div className="p-4 border-t border-white/5 shrink-0">
            <div className="h-10 rounded-lg bg-white/5 animate-pulse" />
          </div>
        </aside>
        <main className="flex-1 flex flex-col h-full overflow-hidden">
          <header className="h-16 bg-white border-b border-slate-100 flex items-center justify-between px-8 shrink-0 shadow-sm z-10">
            <h1 className="font-heading font-semibold text-gray-800 text-lg shrink-0">{title}</h1>
          </header>
          <div className="flex-1 overflow-auto bg-slate-50/50 p-8">
            <div className="max-w-7xl mx-auto space-y-6">
              <div className="h-40 rounded-2xl bg-white animate-pulse" />
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                {[...Array(4)].map((_, i) => (
                  <div key={i} className="h-32 rounded-2xl bg-white animate-pulse" />
                ))}
              </div>
            </div>
          </div>
        </main>
      </div>
    )
  }

  const visibleLinks = sidebarLinks.filter((l) => !l.roles || l.roles.includes(role))

  function handleLogout() {
    logout()
    navigate('/login')
  }

  // Intercept clicks on images with 'previewable-image' class
  const handleLayoutClick = (e) => {
    if (e.target.tagName === 'IMG' && e.target.classList.contains('previewable-image')) {
      setPreviewImage(e.target.src)
    }
  }

  return (
    <div 
      onClick={handleLayoutClick}
      className="font-sans antialiased overflow-hidden flex h-screen bg-[#f8fafc] text-gray-800 relative"
    >
      {/* Mobile Sidebar Backdrop */}
      {isSidebarOpen && (
        <div 
          className="fixed inset-0 z-40 bg-slate-900/60 backdrop-blur-md lg:hidden transition-opacity"
          onClick={() => setIsSidebarOpen(false)}
        />
      )}

      {/* Sidebar */}
      <aside className={`fixed inset-y-0 left-0 w-64 bg-brand-950 text-white flex flex-col h-full shrink-0 shadow-[4px_0_24px_rgba(23,37,84,0.15)] border-r border-brand-900/50 z-50 lg:z-10 lg:relative transform transition-transform duration-300 ease-in-out lg:translate-x-0 ${isSidebarOpen ? 'translate-x-0' : '-translate-x-full'}`}>
        <div className="h-16 flex items-center justify-between px-6 border-b border-brand-900/50 shrink-0 bg-brand-950/40">
          <span className="font-heading font-extrabold text-xl tracking-wide flex items-center gap-2">
            <span className="w-8 h-8 rounded-lg bg-brand-600 flex items-center justify-center text-white text-base font-black shadow-lg shadow-brand-500/30">G</span>
            <span>GRADASI<span className="text-brand-400 font-medium">.</span></span>
          </span>
          <button 
            onClick={() => setIsSidebarOpen(false)} 
            className="lg:hidden p-1.5 text-white/50 hover:text-white rounded-lg hover:bg-white/5 transition-colors"
          >
            <i className="ph ph-x text-xl" />
          </button>
        </div>
        <nav className="flex-1 overflow-y-auto py-6 space-y-1.5 px-4">
          {visibleLinks.map((link) => {
            const isActive = location.pathname === link.path
            return (
              <Link
                key={link.path}
                to={link.path}
                onClick={() => setIsSidebarOpen(false)}
                className={`flex items-center px-4 py-3 rounded-xl transition-all duration-200 btn-press group ${
                  isActive
                    ? 'bg-brand-600 text-white shadow-lg shadow-brand-600/30 font-bold'
                    : 'text-brand-200/70 hover:bg-brand-900/40 hover:text-white'
                }`}
              >
                <i className={`ph ${link.icon} text-lg mr-3.5 transition-transform duration-200 group-hover:scale-110`} /> 
                <span className="text-[13px] tracking-wide font-medium">{link.label}</span>
              </Link>
            )
          })}
        </nav>
        <div className="p-4 border-t border-brand-900/50 shrink-0 bg-brand-950/20">
          <button
            onClick={handleLogout}
            className="w-full flex items-center px-4 py-3 rounded-xl text-red-300 hover:bg-red-500/10 hover:text-red-200 transition-all duration-200 btn-press font-semibold text-[13px] tracking-wide"
          >
            <i className="ph ph-sign-out text-lg mr-3.5" /> Logout
          </button>
        </div>
      </aside>

      {/* Main Content */}
      <main className="flex-1 flex flex-col h-full overflow-hidden">
        {/* Header */}
        <header className="h-16 bg-white/80 backdrop-blur-md border-b border-slate-100 flex items-center justify-between px-6 sm:px-8 shrink-0 shadow-[0_2px_12px_-4px_rgba(0,0,0,0.02)] z-10 gap-4">
          <div className="flex items-center gap-3">
            <button
              onClick={() => setIsSidebarOpen(true)}
              className="lg:hidden p-1.5 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-lg transition-colors"
              title="Open Menu"
            >
              <i className="ph ph-list text-2xl" />
            </button>
            <h1 className="font-heading font-bold text-slate-800 text-base sm:text-lg tracking-tight truncate max-w-[150px] sm:max-w-none">{title}</h1>
          </div>
          
          {/* Center Actions (Search, Filter, etc.) */}
          {headerContent && (
            <div className="flex-1 hidden md:flex items-center justify-center px-6">
              {headerContent}
            </div>
          )}

          <div className="flex items-center gap-4 shrink-0">
            <div className="relative">
              <button
                onClick={() => setUserMenuOpen(!userMenuOpen)}
                className="flex items-center gap-2.5 hover:bg-slate-50 px-3 py-1.5 rounded-xl border border-slate-100 transition-all shadow-sm"
              >
                <img
                  src={user?.photo_path ? resolveAssetUrl(user.photo_path) : `https://ui-avatars.com/api/?name=${encodeURIComponent(user?.name || 'Admin')}&background=2563EB&color=fff`}
                  className="w-7 h-7 rounded-full border border-slate-100 object-cover"
                  alt={user?.name || 'Admin'}
                />
                <span className="text-xs font-semibold text-slate-700 hidden sm:inline">{user?.name || 'Admin'}</span>
                <i className="ph ph-caret-down text-gray-400 text-xs" />
              </button>
              {userMenuOpen && (
                <div className="absolute right-0 mt-2 w-48 bg-white rounded-xl shadow-xl border border-slate-100 py-1.5 z-50 animate-scale-in">
                  <Link
                    to="/admin/profile"
                    onClick={() => setUserMenuOpen(false)}
                    className="flex items-center px-4 py-2 text-xs font-bold text-slate-700 hover:bg-slate-50 hover:text-brand-600 transition-colors"
                  >
                    <i className="ph ph-user-circle mr-2 text-base" /> Profil Saya
                  </Link>
                </div>
              )}
            </div>
          </div>
        </header>

        {/* Content Area */}
        <div className="flex-1 overflow-auto bg-[#f8fafc] p-4 sm:p-8">{children}</div>
      </main>

      {/* Global Image Preview Lightbox Modal */}
      {previewImage && (
        <div 
          className="fixed inset-0 z-[999] flex items-center justify-center bg-slate-950/80 backdrop-blur-md animate-fade-in-up"
          onClick={() => setPreviewImage(null)}
        >
          <button 
            onClick={() => setPreviewImage(null)}
            className="absolute top-6 right-6 w-12 h-12 rounded-full bg-white/10 hover:bg-white/20 text-white flex items-center justify-center text-2xl transition btn-press border border-white/10"
            title="Tutup Preview"
          >
            <i className="ph ph-x" />
          </button>
          <div 
            className="max-w-[90vw] max-h-[85vh] p-2 bg-white rounded-2xl shadow-2xl overflow-hidden border border-white/10 animate-scale-in"
            onClick={e => e.stopPropagation()}
          >
            <img 
              src={previewImage} 
              alt="Preview" 
              className="max-w-full max-h-[80vh] rounded-xl object-contain"
            />
          </div>
        </div>
      )}
    </div>
  )
}
