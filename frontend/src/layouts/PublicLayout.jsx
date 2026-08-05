import { useState, useEffect } from 'react'
import { Link, useLocation } from 'react-router-dom'
import { useSettings } from '../context/useSettings'
import { resolveAssetUrl } from '../utils/assetUrl'
import NavbarBeritaSearch from '../components/NavbarBeritaSearch'

export default function PublicLayout({ children }) {
  const { settings } = useSettings()
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false)
  const [kegiatanDropdown, setKegiatanDropdown] = useState(false)
  const [informasiDropdown, setInformasiDropdown] = useState(false)
  const [scrolled, setScrolled] = useState(false)

  const location = useLocation()
  const isHome = location.pathname === '/'

  useEffect(() => {
    const handleScroll = () => {
      if (window.scrollY > 20) {
        setScrolled(true)
      } else {
        setScrolled(false)
      }
    }
    window.addEventListener('scroll', handleScroll)
    return () => window.removeEventListener('scroll', handleScroll)
  }, [])

  const navTheme = isHome && !scrolled ? 'transparent' : 'glass'

  const socialLinks = [
    { label: 'Facebook', icon: 'ph-facebook-logo', url: settings.facebook_url || '#' },
    { label: 'Instagram', icon: 'ph-instagram-logo', url: settings.instagram_url || '#' },
    { label: 'YouTube', icon: 'ph-youtube-logo', url: settings.youtube_url || '#' },
  ]

  const headerClass = navTheme === 'transparent'
    ? 'fixed top-0 w-full z-50 transition-all duration-300 bg-transparent border-transparent py-5'
    : 'fixed top-0 w-full z-50 transition-all duration-300 bg-white/90 backdrop-blur-xl border-b border-slate-200/80 shadow-[0_10px_30px_-10px_rgba(0,0,0,0.05)] py-3'

  const titleClass = navTheme === 'transparent'
    ? 'font-heading font-extrabold text-base text-white tracking-tight group-hover:text-amber-400 transition'
    : 'font-heading font-extrabold text-base text-slate-900 tracking-tight group-hover:text-brand-700 transition'

  const taglineClass = navTheme === 'transparent'
    ? 'text-[9px] font-bold text-brand-300 tracking-wider uppercase'
    : 'text-[9px] font-bold text-brand-600 tracking-wider uppercase'

  const getLinkClass = (path, matchStart = false) => {
    const isActive = matchStart ? location.pathname.startsWith(path) : location.pathname === path
    if (navTheme === 'transparent') {
      return `px-4 py-2 text-sm font-semibold rounded-lg transition-all duration-200 ${
        isActive ? 'text-white bg-white/15' : 'text-white/80 hover:text-white hover:bg-white/10'
      }`
    } else {
      return `px-4 py-2 text-sm font-semibold rounded-lg transition-all duration-200 ${
        isActive ? 'text-brand-700 bg-brand-50' : 'text-slate-600 hover:text-brand-700 hover:bg-slate-50'
      }`
    }
  }

  const caretClass = navTheme === 'transparent' ? 'text-white/60' : 'text-slate-400'

  const loginBtnClass = navTheme === 'transparent'
    ? 'btn-press group relative inline-flex items-center justify-center gap-2 px-6 py-2.5 text-sm font-bold text-brand-900 bg-white hover:bg-slate-100 rounded-xl transition-all duration-300 shadow-lg shadow-white/5'
    : 'btn-press group relative inline-flex items-center justify-center gap-2 px-6 py-2.5 text-sm font-bold text-white bg-gradient-brand animate-gradient rounded-xl overflow-hidden glow-brand transition-all duration-300'

  const burgerClass = navTheme === 'transparent'
    ? 'text-white hover:text-white/80 focus:outline-none p-2 bg-white/10 hover:bg-white/20 rounded-lg transition'
    : 'text-slate-600 hover:text-brand-700 focus:outline-none p-2 bg-slate-50 hover:bg-slate-100 rounded-lg transition'

  return (
    <div className="font-sans bg-slate-50 min-h-screen">
      {/* Navigation Header — glassmorphism + scroll effect */}
      <header className={headerClass}>
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center">
            {/* Logo Branding */}
            <Link to="/" className="flex-shrink-0 flex items-center gap-3 group">
              <img src={resolveAssetUrl(settings.logo_path)} alt={`Logo ${settings.site_name}`} className="h-10 w-auto object-contain transition duration-300 group-hover:scale-105" />
              <div className="flex flex-col leading-tight">
                <span className={titleClass}>{settings.site_name}</span>
                <span className={taglineClass}>{settings.tagline}</span>
              </div>
            </Link>

            {/* Desktop Navigation */}
            <nav className="hidden lg:flex items-center space-x-1">
              <Link 
                to="/" 
                className={getLinkClass('/')}
              >
                Beranda
              </Link>
              
              <Link 
                to="/#tentang" 
                className={getLinkClass('/#tentang')}
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
                  className={getLinkClass('/kegiatan', true) + ' flex items-center gap-1'}
                >
                  Kegiatan <i className={`ph-bold ph-caret-down text-[10px] transition-transform duration-300 ${caretClass} ${kegiatanDropdown ? 'rotate-180' : ''}`} />
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
                  className={getLinkClass('/berita', true) + ' flex items-center gap-1'}
                >
                  Informasi <i className={`ph-bold ph-caret-down text-[10px] transition-transform duration-300 ${caretClass} ${informasiDropdown ? 'rotate-180' : ''}`} />
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
                className={getLinkClass('/kepengurusan', true)}
              >
                Kepengurusan
              </Link>

              <Link 
                to="/#kontak" 
                className={getLinkClass('/#kontak')}
              >
                Kontak
              </Link>
            </nav>

            {/* Login Button */}
            <div className="hidden lg:flex items-center gap-3">
              <Link 
                to="/login" 
                className={loginBtnClass}
              >
                <i className="ph-bold ph-sign-in text-lg transition-transform duration-300 group-hover:translate-x-0.5" />
                <span>Masuk</span>
              </Link>
            </div>

            {/* Mobile menu toggle */}
            <div className="lg:hidden flex items-center">
              <button 
                onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
                className={burgerClass}
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
              <p>© <span className="text-amber-500 font-semibold">{settings.site_name}</span>, Dilindungi undang-undang.</p>
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
