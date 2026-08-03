import { useState } from 'react'
import { Link, useLocation } from 'react-router-dom'
import { useSettings } from '../context/SettingsContext'
import { resolveAssetUrl } from '../utils/assetUrl'
import NavbarBeritaSearch from '../components/NavbarBeritaSearch'

export default function PublicLayout({ children }) {
  const { settings } = useSettings()
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false)
  const [kegiatanDropdown, setKegiatanDropdown] = useState(false)
  const [informasiDropdown, setInformasiDropdown] = useState(false)
  const [pengurusDropdown, setPengurusDropdown] = useState(false)

  const location = useLocation()

  const socialLinks = [
    { label: 'Facebook', icon: 'ph-facebook-logo', url: settings.facebook_url },
    { label: 'Instagram', icon: 'ph-instagram-logo', url: settings.instagram_url },
    { label: 'YouTube', icon: 'ph-youtube-logo', url: settings.youtube_url },
  ].filter((item) => item.url)

  return (
    <div className="font-sans bg-slate-50 min-h-screen">
      {/* Navigation Header — glassmorphism + scroll effect */}
      <header className="fixed top-0 w-full z-50 transition-all duration-300 bg-white/80 backdrop-blur-xl border-b border-white/40 shadow-[0_1px_20px_-5px_rgba(37,99,235,0.08)]">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-4">
            {/* Logo Branding */}
            <Link to="/" className="flex-shrink-0 flex items-center gap-3 group">
              <img src={resolveAssetUrl(settings.logo_path)} alt={`Logo ${settings.site_name}`} className="h-10 w-auto object-contain transition duration-300 group-hover:scale-105" />
              <div className="flex flex-col leading-tight">
                <span className="font-heading font-extrabold text-base text-slate-900 tracking-tight group-hover:text-brand-700 transition">{settings.site_name}</span>
                <span className="text-[9px] font-bold text-brand-600 tracking-wider uppercase">{settings.tagline}</span>
              </div>
            </Link>

            {/* Berita Search (khusus halaman /berita — tengah navbar) */}
            {location.pathname.startsWith('/berita') && (
              <div className="hidden lg:flex flex-1 justify-center px-4">
                <NavbarBeritaSearch />
              </div>
            )}

            {/* Desktop Navigation */}
            <nav className="hidden lg:flex items-center space-x-1">
              <Link 
                to="/" 
                className={`px-4 py-2 text-sm font-semibold rounded-lg transition ${location.pathname === '/' ? 'text-brand-700 bg-brand-50' : 'text-slate-600 hover:text-brand-700 hover:bg-slate-50'}`}
              >
                Beranda
              </Link>
              
              <Link 
                to="/#tentang" 
                className="px-4 py-2 text-sm font-medium text-slate-600 hover:text-brand-700 hover:bg-slate-50 transition rounded-lg"
              >
                Tentang Kami
              </Link>

              {/* Kegiatan Dropdown */}
              <div 
                className="relative"
                onMouseEnter={() => setKegiatanDropdown(true)}
                onMouseLeave={() => setKegiatanDropdown(false)}
              >
                <Link 
                  to="/kegiatan"
                  className={`px-4 py-2 text-sm font-semibold rounded-lg transition flex items-center gap-1 ${location.pathname.startsWith('/kegiatan') ? 'text-brand-700 bg-brand-50' : 'text-slate-600 hover:text-brand-700 hover:bg-slate-50'}`}
                >
                  Kegiatan <i className={`ph-bold ph-caret-down text-[10px] transition-transform duration-300 ${kegiatanDropdown ? 'rotate-180' : ''}`} />
                </Link>
                {kegiatanDropdown && (
                  <div className="absolute left-0 top-full pt-1 w-52 z-50">
                    <div className="glass animate-scale-in rounded-xl shadow-xl py-2 border border-slate-100">
                      <Link 
                        to="/kegiatan" 
                        onClick={() => setKegiatanDropdown(false)}
                        className="block px-4 py-2.5 text-sm font-bold text-brand-700 bg-brand-50 hover:bg-brand-100 transition"
                      >
                        <i className="ph-bold ph-squares-four mr-2" /> Semua Kegiatan
                      </Link>
                      <Link 
                        to="/#kegiatan" 
                        onClick={() => setKegiatanDropdown(false)}
                        className="block px-4 py-2.5 text-sm text-slate-600 hover:bg-brand-50 hover:text-brand-700 font-medium transition"
                      >
                        <i className="ph-bold ph-calendar-blank mr-2" /> Agenda Terkini
                      </Link>
                    </div>
                  </div>
                )}
              </div>

              {/* Informasi Dropdown */}
              <div 
                className="relative"
                onMouseEnter={() => setInformasiDropdown(true)}
                onMouseLeave={() => setInformasiDropdown(false)}
              >
                <Link 
                  to="/berita"
                  className={`px-4 py-2 text-sm font-semibold rounded-lg transition flex items-center gap-1 ${location.pathname.startsWith('/berita') ? 'text-brand-700 bg-brand-50' : 'text-slate-600 hover:text-brand-700 hover:bg-slate-50'}`}
                >
                  Informasi <i className={`ph-bold ph-caret-down text-[10px] transition-transform duration-300 ${informasiDropdown ? 'rotate-180' : ''}`} />
                </Link>
                {informasiDropdown && (
                  <div className="absolute left-0 top-full pt-1 w-52 z-50">
                    <div className="glass animate-scale-in rounded-xl shadow-xl py-2 border border-slate-100">
                      <Link 
                        to="/berita" 
                        onClick={() => setInformasiDropdown(false)}
                        className="block px-4 py-2.5 text-sm font-bold text-brand-700 bg-brand-50 hover:bg-brand-100 transition"
                      >
                        <i className="ph-bold ph-newspaper mr-2" /> Semua Informasi
                      </Link>
                      <Link 
                        to="/#informasi" 
                        onClick={() => setInformasiDropdown(false)}
                        className="block px-4 py-2.5 text-sm text-slate-600 hover:bg-brand-50 hover:text-brand-700 font-medium transition"
                      >
                        <i className="ph-bold ph-article mr-2" /> Berita Terkini
                      </Link>
                    </div>
                  </div>
                )}
              </div>

              <Link 
                to="/kepengurusan" 
                className={`px-4 py-2 text-sm font-semibold rounded-lg transition ${location.pathname.startsWith('/kepengurusan') ? 'text-brand-700 bg-brand-50' : 'text-slate-600 hover:text-brand-700 hover:bg-slate-50'}`}
              >
                Kepengurusan
              </Link>

              <Link 
                to="/#kontak" 
                className="px-4 py-2 text-sm font-medium text-slate-600 hover:text-brand-700 hover:bg-slate-50 transition rounded-lg"
              >
                Kontak
              </Link>
            </nav>

            {/* Login Button */}
            <div className="hidden lg:flex items-center gap-3">
              <Link 
                to="/login" 
                className="btn-press group relative inline-flex items-center justify-center gap-2 px-6 py-2.5 text-sm font-bold text-white bg-gradient-brand animate-gradient rounded-xl overflow-hidden glow-brand transition-all duration-300"
              >
                <i className="ph-bold ph-sign-in text-lg transition-transform duration-300 group-hover:translate-x-0.5" />
                <span>Masuk</span>
              </Link>
            </div>

            {/* Mobile menu toggle */}
            <div className="lg:hidden flex items-center">
              <button 
                onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
                className="text-slate-600 hover:text-brand-700 focus:outline-none p-2 bg-slate-50 rounded-lg"
              >
                <i className={`ph-bold text-2xl ${mobileMenuOpen ? 'ph-x' : 'ph-list'}`} />
              </button>
            </div>
          </div>
        </div>

        {/* Mobile Navigation Drawer */}
        {mobileMenuOpen && (
          <div className="lg:hidden bg-white border-t border-slate-100 absolute w-full shadow-2xl max-h-[85vh] overflow-y-auto z-50">
            <div className="px-4 py-4 space-y-2">
              <Link to="/" onClick={() => setMobileMenuOpen(false)} className="block px-4 py-3 text-sm font-bold text-brand-700 bg-brand-50 rounded-lg">Beranda</Link>
              <Link to="/#tentang" onClick={() => setMobileMenuOpen(false)} className="block px-4 py-3 text-sm font-medium text-slate-700 hover:bg-slate-50 rounded-lg">Tentang Kami</Link>
              <Link to="/kegiatan" onClick={() => setMobileMenuOpen(false)} className="block px-4 py-3 text-sm font-medium text-slate-700 hover:bg-slate-50 rounded-lg">Kegiatan</Link>
              <Link to="/berita" onClick={() => setMobileMenuOpen(false)} className="block px-4 py-3 text-sm font-medium text-slate-700 hover:bg-slate-50 rounded-lg">Informasi & Berita</Link>
              <Link to="/kepengurusan" onClick={() => setMobileMenuOpen(false)} className="block px-4 py-3 text-sm font-medium text-slate-700 hover:bg-slate-50 rounded-lg">Kepengurusan</Link>
              <Link to="/#kontak" onClick={() => setMobileMenuOpen(false)} className="block px-4 py-3 text-sm font-medium text-slate-700 hover:bg-slate-50 rounded-lg">Kontak</Link>
              <div className="border-t border-slate-100 pt-3 mt-3">
                <Link to="/login" onClick={() => setMobileMenuOpen(false)} className="flex items-center justify-center gap-2 px-4 py-3 text-sm font-bold text-white bg-brand-600 rounded-xl shadow-md">
                  <i className="ph-bold ph-sign-in text-lg" />
                  Masuk Admin
                </Link>
              </div>
            </div>
          </div>
        )}
      </header>

      {/* Main Content */}
      <main>
        {children}
      </main>

      {/* Footer */}
      <footer className="bg-brand-950 border-t-4 border-brand-700 text-slate-400 text-sm relative overflow-hidden">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10 pt-16 pb-8">
          <div className="flex flex-col md:flex-row justify-between items-center gap-8 border-b border-brand-800/50 pb-8 mb-8">
            <div className="flex items-center gap-3">
              <div className="w-12 h-12 bg-white rounded-lg flex items-center justify-center p-1 shadow-inner">
                <img src={resolveAssetUrl(settings.logo_path)} alt={`Logo ${settings.site_name}`} className="w-full h-full object-contain" />
              </div>
              <div>
                <span className="text-white font-heading font-bold text-lg block leading-tight">{settings.site_name}</span>
                <span className="text-[10px] tracking-widest uppercase text-brand-300 font-semibold">{settings.tagline}</span>
              </div>
            </div>
            <div className="flex gap-3">
              {socialLinks.map((item) => (
                <a key={item.label} href={item.url} target="_blank" rel="noreferrer" aria-label={item.label} className="w-10 h-10 rounded-lg bg-brand-900/50 border border-brand-800 flex items-center justify-center hover:bg-brand-700 hover:border-brand-700 hover:text-white transition shadow-sm">
                  <i className={`ph-fill ${item.icon} text-lg`} />
                </a>
              ))}
            </div>
          </div>
          <div className="flex flex-col md:flex-row justify-between items-center text-[13px] text-brand-200/70 border-t border-brand-800/50 pt-8 mt-8">
            <div className="flex flex-col text-center md:text-left gap-1 mb-4 md:mb-0">
              <p>© 2026 {settings.site_name}. Dilindungi Undang-Undang.</p>
              <p className="text-[11px] text-brand-300/50">{settings.address}</p>
            </div>
            <div className="flex gap-6 font-medium">
              <a href="#" className="hover:text-white transition">Kebijakan Privasi</a>
              <a href="#" className="hover:text-white transition">Syarat & Ketentuan</a>
            </div>
          </div>
        </div>
      </footer>
    </div>
  )
}
