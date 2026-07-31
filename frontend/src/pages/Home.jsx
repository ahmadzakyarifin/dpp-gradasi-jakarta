import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import PublicLayout from '../layouts/PublicLayout'
import { slidersService } from '../services/slidersService'
import { settingsService } from '../services/settingsService'
import { kegiatanService } from '../services/kegiatanService'
import { beritaService } from '../services/beritaService'

export default function Home() {
  const [sliders, setSliders] = useState([])
  const [settings, setSettings] = useState({
    site_name: 'DPP GRADASI',
    tagline: 'Generasi Digital Indonesia',
    logo_url: '/uploads/logo.png',
    contact_email: 'dpp@gradasi.org',
    contact_phone: '+6285279880008',
    address: 'Jl. Jenderal Sudirman No.1, Jakarta Pusat',
    maps_embed_url: '',
    facebook_url: '',
    instagram_url: '',
    youtube_url: '',
    video_profile_url: '',
    history: '',
    about_formation_date: '',
    about_no_sk: '',
    about_vision: '',
    about_mission: '[]',
    greeting_title: '',
    greeting_subtitle: '',
    greeting_date: '',
    greeting_content: '',
    greeting_image_url: ''
  })
  
  const [featuredKegiatan, setFeaturedKegiatan] = useState([])
  const [recentBerita, setRecentBerita] = useState([])

  const [currentSlide, setCurrentSlide] = useState(0)
  const [aboutTab, setAboutTab] = useState('selayang')

  useEffect(() => {
    // Load Sliders
    slidersService.list(true)
      .then(res => {
        if (res.success && res.data && res.data.sliders) {
          setSliders(res.data.sliders)
        }
      }).catch(() => {})

    // Load Settings
    settingsService.get()
      .then(res => {
        if (res.success && res.data) {
          setSettings(res.data)
        }
      }).catch(() => {})

    // Load Kegiatan
    kegiatanService.list()
      .then(res => {
        if (res.success && res.data) {
          setFeaturedKegiatan(res.data.slice(0, 3))
        }
      }).catch(() => {})

    // Load Berita
    beritaService.list()
      .then(res => {
        if (res.success && res.data) {
          setRecentBerita(res.data.slice(0, 3))
        }
      }).catch(() => {})
  }, [])

  // Auto play slider logic
  useEffect(() => {
    if (sliders.length === 0) return
    const interval = setInterval(() => {
      setCurrentSlide(prev => (prev + 1) % sliders.length)
    }, 5000)
    return () => clearInterval(interval)
  }, [sliders])

  const getMissions = () => {
    try {
      if (typeof settings.about_mission === 'string') {
        return JSON.parse(settings.about_mission)
      }
      return settings.about_mission || []
    } catch {
      return (settings.about_mission || '').split('\n').map(s => s.trim()).filter(Boolean)
    }
  }

  return (
    <PublicLayout>
      {/* 1. HERO CAROUSEL */}
      {sliders.length > 0 && (
        <section className="relative bg-brand-950 overflow-hidden pt-28 pb-20 md:pt-40 md:pb-24 min-h-[90vh] flex items-center">
          <div className="max-w-6xl mx-auto px-4 w-full grid grid-cols-1 lg:grid-cols-12 gap-8 items-center">
            <div className="lg:col-span-7 space-y-6 text-white">
              {sliders[currentSlide].tag && (
                <span className="inline-block px-3 py-1 bg-brand-600/50 backdrop-blur-sm border border-brand-500/50 rounded-lg text-xs font-bold uppercase tracking-wider">
                  {sliders[currentSlide].tag} {sliders[currentSlide].is_new && <span className="ml-1 text-[10px] text-yellow-300">NEW</span>}
                </span>
              )}
              <h1 className="font-heading text-4xl md:text-5xl lg:text-6xl font-extrabold leading-tight">
                {sliders[currentSlide].title}
              </h1>
              <p className="text-lg text-brand-100/80 leading-relaxed max-w-xl">
                {sliders[currentSlide].subtitle}
              </p>
              {sliders[currentSlide].event_date && (
                <div className="flex flex-wrap gap-4 text-sm text-brand-200">
                  <span>📅 {sliders[currentSlide].event_date}</span>
                  {sliders[currentSlide].location && <span>📍 {sliders[currentSlide].location}</span>}
                </div>
              )}
              {sliders[currentSlide].link_url && (
                <a href={sliders[currentSlide].link_url} className="inline-block bg-brand-600 hover:bg-brand-700 px-6 py-3 rounded-xl font-bold transition">
                  Pelajari Lebih Lanjut
                </a>
              )}
            </div>
            
            <div className="lg:col-span-5 relative flex justify-center">
              <div className="w-80 h-80 md:w-96 md:h-96 rounded-3xl overflow-hidden shadow-2xl border border-white/10">
                <img 
                  src={sliders[currentSlide].image_url.startsWith('http') ? sliders[currentSlide].image_url : `http://127.0.0.1:8080${sliders[currentSlide].image_url}`} 
                  alt={sliders[currentSlide].title}
                  className="w-full h-full object-cover"
                />
              </div>
            </div>
          </div>

          {/* Dots Indicator */}
          <div className="absolute bottom-8 left-1/2 transform -translate-x-1/2 flex gap-2">
            {sliders.map((_, idx) => (
              <button 
                key={idx}
                onClick={() => setCurrentSlide(idx)}
                className={`w-3 h-3 rounded-full transition ${currentSlide === idx ? 'bg-white' : 'bg-white/30'}`}
              />
            ))}
          </div>
        </section>
      )}

      {/* 2. SAMBUTAN GREETING */}
      {settings.greeting_title && (
        <section className="py-20 bg-white border-b border-slate-100">
          <div className="max-w-6xl mx-auto px-4 grid grid-cols-1 md:grid-cols-12 gap-8 items-center">
            {settings.greeting_image_url && (
              <div className="md:col-span-5">
                <div className="rounded-3xl overflow-hidden shadow-md">
                  <img src={settings.greeting_image_url} alt="Sambutan" className="w-full h-auto object-cover" />
                </div>
              </div>
            )}
            <div className="md:col-span-7 space-y-4">
              <span className="text-xs font-bold text-brand-600 uppercase tracking-widest block">
                {settings.greeting_subtitle}
              </span>
              <h2 className="font-heading text-3xl font-extrabold text-slate-900">
                {settings.greeting_title}
              </h2>
              <p className="text-slate-600 leading-relaxed">
                {settings.greeting_content}
              </p>
            </div>
          </div>
        </section>
      )}

      {/* 3. TENTANG KAMI */}
      <section id="tentang" className="py-20 bg-slate-50 border-b border-slate-200">
        <div className="max-w-6xl mx-auto px-4 grid grid-cols-1 lg:grid-cols-12 gap-12 items-center">
          <div className="lg:col-span-7 space-y-6">
            <h2 className="font-heading text-3xl md:text-4xl font-extrabold text-slate-900">Tentang Kami</h2>
            
            {/* Tabs */}
            <div className="flex gap-2 border-b border-slate-200 pb-2">
              <button 
                onClick={() => setAboutTab('selayang')}
                className={`pb-2 text-sm font-bold border-b-2 transition ${aboutTab === 'selayang' ? 'border-brand-600 text-brand-600' : 'border-transparent text-slate-500 hover:text-slate-700'}`}
              >
                Selayang Pandang
              </button>
              <button 
                onClick={() => setAboutTab('sejarah')}
                className={`pb-2 text-sm font-bold border-b-2 transition ${aboutTab === 'sejarah' ? 'border-brand-600 text-brand-600' : 'border-transparent text-slate-500 hover:text-slate-700'}`}
              >
                Sejarah
              </button>
            </div>

            <div className="text-slate-600 leading-relaxed min-h-[150px]">
              {aboutTab === 'selayang' ? (
                <div className="space-y-4">
                  <p>{settings.history}</p>
                </div>
              ) : (
                <div className="space-y-2">
                  <p><strong>Legalitas Resmi:</strong> SK Kemenkumham No. {settings.about_no_sk}</p>
                  <p><strong>Tanggal Berdiri:</strong> {settings.about_formation_date}</p>
                  <p>{settings.about_tutorial}</p>
                </div>
              )}
            </div>
          </div>

          <div className="lg:col-span-5 flex justify-center">
            <div className="p-8 bg-white border border-slate-100 rounded-3xl shadow-md text-center max-w-sm w-full">
              <h3 className="font-heading font-extrabold text-xl text-slate-900 mb-2">Visi Kami</h3>
              <p className="text-slate-600 italic">"{settings.about_vision}"</p>
            </div>
          </div>
        </div>
      </section>

      {/* 4. VISI & MISI */}
      <section className="py-20 bg-white border-b border-slate-100">
        <div className="max-w-6xl mx-auto px-4">
          <div className="text-center mb-16">
            <h2 className="font-heading text-3xl font-extrabold text-slate-900">Misi Organisasi</h2>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {getMissions().map((mission, idx) => (
              <div 
                key={idx} 
                className={`p-8 rounded-3xl border border-slate-100 shadow-sm transition duration-300 hover:-translate-y-1 ${
                  idx % 3 === 0 ? 'bg-blue-50/50 hover:bg-blue-50' : 
                  idx % 3 === 1 ? 'bg-teal-50/50 hover:bg-teal-50' : 'bg-amber-50/50 hover:bg-amber-50'
                }`}
              >
                <h3 className="font-heading font-bold text-slate-900 text-lg mb-2">Misi {idx + 1}</h3>
                <p className="text-slate-600 text-sm leading-relaxed">{mission}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* 5. KEGIATAN TERBARU */}
      {featuredKegiatan.length > 0 && (
        <section className="py-20 bg-slate-50 border-b border-slate-200">
          <div className="max-w-6xl mx-auto px-4">
            <div className="flex justify-between items-end mb-12">
              <div>
                <h2 className="font-heading text-3xl font-extrabold text-slate-900">Kegiatan Terkini</h2>
                <p className="text-slate-500 mt-1 text-sm">Ikuti berbagai kegiatan terbaru dari GRADASI</p>
              </div>
              <Link to="/kegiatan" className="text-brand-600 font-bold hover:underline text-sm">Semua Kegiatan →</Link>
            </div>
            
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              {featuredKegiatan.map(item => (
                <Link to={`/kegiatan/${item.slug}`} key={item.id} className="group bg-white rounded-2xl overflow-hidden shadow-sm border border-slate-100 flex flex-col h-full hover:shadow-md transition">
                  <div className="h-44 relative bg-slate-100 overflow-hidden">
                    <img 
                      src={item.image_url ? (item.image_url.startsWith('http') ? item.image_url : `http://127.0.0.1:8080${item.image_url}`) : 'https://images.unsplash.com/photo-1540575467063-178a50c2df87?q=80&w=600'} 
                      alt={item.title} 
                      className="w-full h-full object-cover transform group-hover:scale-102 transition duration-300"
                    />
                  </div>
                  <div className="p-5 flex flex-col flex-grow">
                    <span className="text-[10px] font-bold tracking-wider text-brand-600 uppercase mb-2 block">{item.category}</span>
                    <h3 className="font-heading font-bold text-slate-900 mb-2 group-hover:text-brand-600 transition line-clamp-2">{item.title}</h3>
                    <p className="text-slate-500 text-xs mb-4 line-clamp-3">{item.excerpt}</p>
                    <div className="mt-auto pt-4 border-t border-slate-100 flex justify-between text-[11px] text-slate-400 font-medium">
                      <span>📅 {item.event_date}</span>
                      <span>📍 {item.location}</span>
                    </div>
                  </div>
                </Link>
              ))}
            </div>
          </div>
        </section>
      )}

      {/* 6. BERITA TERBARU */}
      {recentBerita.length > 0 && (
        <section className="py-20 bg-white">
          <div className="max-w-6xl mx-auto px-4">
            <div className="flex justify-between items-end mb-12">
              <div>
                <h2 className="font-heading text-3xl font-extrabold text-slate-900">Informasi & Berita</h2>
                <p className="text-slate-500 mt-1 text-sm">Baca artikel dan berita perkembangan digital terkini</p>
              </div>
              <Link to="/berita" className="text-brand-600 font-bold hover:underline text-sm">Semua Berita →</Link>
            </div>
            
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              {recentBerita.map(item => (
                <Link to={`/berita/${item.slug}`} key={item.id} className="group bg-white rounded-2xl overflow-hidden shadow-sm border border-slate-100 flex flex-col h-full hover:shadow-md transition">
                  <div className="h-44 relative bg-slate-100 overflow-hidden">
                    <img 
                      src={item.image_url ? (item.image_url.startsWith('http') ? item.image_url : `http://127.0.0.1:8080${item.image_url}`) : 'https://images.unsplash.com/photo-1504711434969-e33886168f5c?q=80&w=600'} 
                      alt={item.title} 
                      className="w-full h-full object-cover transform group-hover:scale-102 transition duration-300"
                    />
                  </div>
                  <div className="p-5 flex flex-col flex-grow">
                    <span className="text-[10px] font-bold tracking-wider text-brand-600 uppercase mb-2 block">{item.category}</span>
                    <h3 className="font-heading font-bold text-slate-900 mb-2 group-hover:text-brand-600 transition line-clamp-2">{item.title}</h3>
                    <p className="text-slate-500 text-xs mb-4 line-clamp-3">{item.excerpt}</p>
                    <div className="mt-auto pt-4 border-t border-slate-100 text-[11px] text-slate-400 font-medium">
                      <span>👤 {item.author_name || 'Admin'} • 📅 {new Date(item.created_at).toLocaleDateString('id-ID')}</span>
                    </div>
                  </div>
                </Link>
              ))}
            </div>
          </div>
        </section>
      )}
    </PublicLayout>
  )
}
