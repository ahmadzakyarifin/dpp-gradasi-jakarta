import { Link } from 'react-router-dom'

const logoUrl = 'https://gradasi.org/uploads/img/logo/1737187847.png'
const socialLinks = [
  { label: 'Facebook', icon: 'ph-facebook-logo', url: 'https://www.facebook.com/gradasiofficial.id' },
  { label: 'Instagram', icon: 'ph-instagram-logo', url: 'https://www.instagram.com/dppgradasi' },
  { label: 'YouTube', icon: 'ph-youtube-logo', url: 'https://www.youtube.com/channel/UCwdjB4LkqcF4Kw5-PoyOb5A' },
]

export default function PublicLayout({ children }) {
  return (
    <div className="font-sans bg-slate-50 min-h-screen">
      <header className="fixed top-0 w-full z-50 transition-all duration-300 bg-white py-5 border-b border-slate-100">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center">
            <Link to="/" className="flex-shrink-0 flex items-center gap-3 group">
              <img src={logoUrl} alt="Logo GRADASI" className="h-10 w-auto object-contain transition duration-300 group-hover:scale-105" />
              <div className="flex flex-col leading-tight">
                <span className="font-heading font-extrabold text-base text-slate-900 tracking-tight group-hover:text-brand-700 transition">DPP GRADASI</span>
                <span className="text-[9px] font-bold text-brand-600 tracking-wider uppercase">Generasi Digital Indonesia</span>
              </div>
            </Link>

            <nav className="hidden lg:flex items-center space-x-2">
              <Link to="/" className="px-4 py-2 text-sm font-medium text-slate-600 hover:text-brand-700 hover:bg-slate-50 transition rounded-lg">Beranda</Link>
              <Link to="/berita" className="px-4 py-2 text-sm font-semibold text-brand-700 bg-brand-50 transition rounded-lg">Informasi</Link>
              <Link to="/login" className="group relative inline-flex items-center justify-center gap-2 px-6 py-2.5 text-sm font-bold text-white bg-brand-600 hover:bg-brand-700 rounded-xl overflow-hidden shadow-sm transition-all duration-300 hover:shadow-md hover:-translate-y-0.5">
                <i className="ph-bold ph-sign-in text-lg" />
                <span>Masuk</span>
              </Link>
            </nav>
          </div>
        </div>
      </header>

      {children}

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
