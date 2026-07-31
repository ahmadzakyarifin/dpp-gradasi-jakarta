import { useState } from 'react'
import { Link, useLocation } from 'react-router-dom'

const logoUrl = 'https://gradasi.org/uploads/img/logo/1737187847.png'

const socialLinks = [
  { label: 'Facebook', icon: 'ph-facebook-logo', url: 'https://www.facebook.com/gradasiofficial.id' },
  { label: 'Instagram', icon: 'ph-instagram-logo', url: 'https://www.instagram.com/dppgradasi' },
  { label: 'YouTube', icon: 'ph-youtube-logo', url: 'https://www.youtube.com/channel/UCwdjB4LkqcF4Kw5-PoyOb5A' },
]

export default function PublicLayout({ children }) {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false)
  const [kegiatanDropdown, setKegiatanDropdown] = useState(false)
  const [informasiDropdown, setInformasiDropdown] = useState(false)
  const [pengurusDropdown, setPengurusDropdown] = useState(false)

  const location = useLocation()

  return (
    <div className="font-sans bg-slate-50 min-h-screen">
      {/* Navigation Header */}
      <header className="fixed top-0 w-full z-50 transition-all duration-300 bg-white py-4 border-b border-slate-200 shadow-xs">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center">
            {/* Logo Branding */}
            <Link to="/" className="flex-shrink-0 flex items-center gap-3 group">
              <img src={logoUrl} alt="Logo GRADASI" className="h-10 w-auto object-contain transition duration-300 group-hover:scale-105" />
              <div className="flex flex-col leading-tight">
                <span className="font-heading font-extrabold text-base text-slate-900 tracking-tight group-hover:text-brand-700 transition">DPP GRADASI</span>
                <span className="text-[9px] font-bold text-brand-600 tracking-wider uppercase">Generasi Digital Indonesia</span>
              </div>
            </Link>

            {/* Desktop Navigation */}
            <nav className="hidden lg:flex items-center space-x-1">
              <Link 
                to="/" 
                className={`px-4 py-2 text-sm font-semibold rounded-lg transition ${location.pathname === '/' ? 'text-brand-700 bg-brand-50' : 'text-slate-600 hover:text-brand-700 hover:bg-slate-50'}`}
              >
                Beranda
              </Link>
              
              <a 
                href="/#tentang" 
                className="px-4 py-2 text-sm font-medium text-slate-600 hover:text-brand-700 hover:bg-slate-50 transition rounded-lg"
              >
                Tentang Kami
              </a>

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
                    <div className="bg-white rounded-xl shadow-xl py-2 border border-slate-100 animate-fadeIn">
                      <Link 
                        to="/kegiatan" 
                        onClick={() => setKegiatanDropdown(false)}
                        className="block px-4 py-2.5 text-sm font-bold text-brand-700 bg-brand-50 hover:bg-brand-100 transition"
                      >
                        <i className="ph-bold ph-squares-four mr-2" /> Semua Kegiatan
                      </Link>
                      <a 
                        href="/#kegiatan" 
                        onClick={() => setKegiatanDropdown(false)}
                        className="block px-4 py-2.5 text-sm text-slate-600 hover:bg-brand-50 hover:text-brand-700 font-medium transition"
                      >
                        <i className="ph-bold ph-calendar-blank mr-2" /> Agenda Terkini
                      </a>
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
                    <div className="bg-white rounded-xl shadow-xl py-2 border border-slate-100 animate-fadeIn">
                      <Link 
                        to="/berita" 
                        onClick={() => setInformasiDropdown(false)}
                        className="block px-4 py-2.5 text-sm font-bold text-brand-700 bg-brand-50 hover:bg-brand-100 transition"
                      >
                        <i className="ph-bold ph-newspaper mr-2" /> Semua Informasi
                      </Link>
                      <a 
                        href="/#informasi" 
                        onClick={() => setInformasiDropdown(false)}
                        className="block px-4 py-2.5 text-sm text-slate-600 hover:bg-brand-50 hover:text-brand-700 font-medium transition"
                      >
                        <i className="ph-bold ph-article mr-2" /> Berita Terkini
                      </a>
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

              <a href="/#kontak" className="px-4 py-2 text-sm font-medium text-slate-600 hover:text-brand-700 hover:bg-slate-50 transition rounded-lg">
                Kontak
              </a>
            </nav>

            {/* Login Button */}
            <div className="hidden lg:flex items-center gap-3">
              <Link 
                to="/login" 
                className="group relative inline-flex items-center justify-center gap-2 px-6 py-2.5 text-sm font-bold text-white bg-brand-600 hover:bg-brand-700 rounded-xl overflow-hidden shadow-sm transition-all duration-300 hover:shadow-md hover:-translate-y-0.5"
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
              <a href="/#tentang" onClick={() => setMobileMenuOpen(false)} className="block px-4 py-3 text-sm font-medium text-slate-700 hover:bg-slate-50 rounded-lg">Tentang Kami</a>
              <Link to="/kegiatan" onClick={() => setMobileMenuOpen(false)} className="block px-4 py-3 text-sm font-medium text-slate-700 hover:bg-slate-50 rounded-lg">Kegiatan</Link>
              <Link to="/berita" onClick={() => setMobileMenuOpen(false)} className="block px-4 py-3 text-sm font-medium text-slate-700 hover:bg-slate-50 rounded-lg">Informasi & Berita</Link>
              <Link to="/kepengurusan" onClick={() => setMobileMenuOpen(false)} className="block px-4 py-3 text-sm font-medium text-slate-700 hover:bg-slate-50 rounded-lg">Kepengurusan</Link>
              <a href="/#kontak" onClick={() => setMobileMenuOpen(false)} className="block px-4 py-3 text-sm font-medium text-slate-700 hover:bg-slate-50 rounded-lg">Kontak</a>
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
                <img src={logoUrl} alt="Logo" className="w-full h-full object-contain" />
              </div>
              <div>
                <span className="text-white font-heading font-bold text-lg block leading-tight">DPP GRADASI</span>
                <span className="text-[10px] tracking-widest uppercase text-brand-300 font-semibold">Generasi Digital Indonesia</span>
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
              <p>© 2026 Perkumpulan Generasi Digital Indonesia. Dilindungi Undang-Undang.</p>
              <p className="text-[11px] text-brand-300/50">Office Park OL3-IZA The Bellagio Mall, Mega Kuningan, Jakarta Selatan.</p>
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
