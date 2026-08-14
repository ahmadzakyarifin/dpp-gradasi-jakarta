import { useState, useEffect, useRef, useCallback } from 'react'
import { Link, useLocation } from 'react-router-dom'
import PublicLayout from '../layouts/PublicLayout'
import { useSettings } from '../context/useSettings'
import CaptchaWidget from '../components/auth/CaptchaWidget'
import { resolveAssetUrl } from '../utils/assetUrl'
import { slidersService } from '../services/slidersService'
import { kegiatanService } from '../services/kegiatanService'
import { beritaService } from '../services/beritaService'
import { kontakService } from '../services/kontakService'
import { pengurusService } from '../services/pengurusService'
import useReveal from '../hooks/useReveal'
import { shareContent, getShareUrl, copyToClipboard } from '../utils/share'
import ToastNotification from '../components/admin/ToastNotification'
import { parseApiError } from '../utils/parseApiError'

export default function Home() {
  const { settings } = useSettings()
  const revealVisimisi = useReveal()
  const revealKegiatan = useReveal()
  const revealBerita = useReveal()
  const revealKontak = useReveal()

  const [sliders, setSliders] = useState([])

  const [featuredKegiatan, setFeaturedKegiatan] = useState([])
  const [toast, setToast] = useState({ show: false, message: '', type: 'success' })

  const handleInstagramShare = (url) => {
    copyToClipboard(url).then(success => {
      if (success) {
        setToast({ show: true, message: 'Tautan disalin! Membuka Instagram...', type: 'success' });
        window.open('https://www.instagram.com/', '_blank');
      }
    });
  };

  const handleCopyLink = (url, label = 'Tautan') => {
    copyToClipboard(url).then(success => {
      if (success) {
        setToast({ show: true, message: `${label} berhasil disalin!`, type: 'success' });
      }
    });
  };

  const [recentBerita, setRecentBerita] = useState([])

  const [ketuaUmum, setKetuaUmum] = useState(null)
  const [sekjen, setSekjen] = useState(null)

  const [currentSlide, setCurrentSlide] = useState(0)
  const [isTransitioning, setIsTransitioning] = useState(false)
  const [activeShareId, setActiveShareId] = useState(null) // ID format: 'kegiatan-ID' atau 'berita-ID'
  const [progress, setProgress] = useState(0)
  const [aboutTab, setAboutTab] = useState('selayang')

  // Contact Form State
  const [contactForm, setContactForm] = useState({ nama: '', email: '', subjek: '', pesan: '' })
  const [contactSuccess, setContactSuccess] = useState(false)
  const [contactLoading, setContactLoading] = useState(false)
  const [contactErrors, setContactErrors] = useState({})
  const [contactTouched, setContactTouched] = useState({})
  const captchaEnabled = !!settings?.captcha_enabled

  const validateContactForm = (data = contactForm) => {
    const errors = {}
    if (!data.nama || !data.nama.trim()) {
      errors.nama = 'Nama lengkap wajib diisi.'
    }
    if (!data.email || !data.email.trim()) {
      errors.email = 'Alamat email wajib diisi.'
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(data.email)) {
      errors.email = 'Format email tidak valid, contoh: nama@email.com'
    }
    if (!data.pesan || !data.pesan.trim()) {
      errors.pesan = 'Pesan Anda wajib diisi.'
    } else if (data.pesan.trim().length < 10) {
      errors.pesan = 'Pesan minimal harus 10 karakter.'
    }
    return errors
  }

  const handleContactChange = (field, val) => {
    const next = { ...contactForm, [field]: val }
    setContactForm(next)
    
    if (contactTouched[field]) {
      const errs = validateContactForm(next)
      setContactErrors(prev => ({ ...prev, [field]: errs[field] }))
    }
  }

  const handleContactBlur = (field) => {
    setContactTouched(prev => ({ ...prev, [field]: true }))
    const errs = validateContactForm()
    setContactErrors(prev => ({ ...prev, [field]: errs[field] }))
  }


  // Video Playing State
  const [isPlayingVideo, setIsPlayingVideo] = useState(false)
  const videoRef = useRef(null)

  // Refs for horizontal slider scrolling
  const kegiatanSliderRef = useRef(null)
  const beritaSliderRef = useRef(null)

  const loadData = useCallback(() => {
    // Load Sliders
    slidersService.list()
      .then(res => {
        if (res.success && res.data) {
          const list = Array.isArray(res.data) ? res.data : (res.data.sliders || [])
          // Normalisasi: API kirim image_path → FE pakai image_url di JSX
          setSliders(list.map(s => ({ ...s, image_url: s.image_path || s.image_url })))
        }
      }).catch(() => { })

    // Load Kegiatan
    kegiatanService.list()
      .then(res => {
        if (res.success && res.data) {
          const list = Array.isArray(res.data) ? res.data : (res.data.kegiatan || [])
          setFeaturedKegiatan(list.map(k => ({ ...k, image_url: k.image_path || k.image_url })))
        }
      }).catch(() => { })

    // Load Berita
    beritaService.list()
      .then(res => {
        if (res.success && res.data) {
          const list = Array.isArray(res.data) ? res.data : (res.data.berita || [])
          setRecentBerita(list.map(b => ({ ...b, image_url: b.image_path || b.image_url })))
        }
      }).catch(() => { })

    // Load Pengurus untuk Tanda Tangan Sambutan Ketua Umum
    pengurusService.list({ limit: 100 })
      .then(res => {
        if (res.success && res.data) {
          const list = Array.isArray(res.data) ? res.data : (res.data.data || res.data.pengurus || [])
          const ketua = list.find(p => p.level === 'ketua')
          const sek = list.find(p => (p.role || '').toLowerCase().includes('sekretaris jenderal') || (p.role || '').toLowerCase().includes('sekjen'))
          setKetuaUmum(ketua)
          setSekjen(sek)
        }
      }).catch(() => { })
  }, [])

  useEffect(() => {
    loadData()
    // Re-fetch saat tab/window publik dapat fokus lagi (mis. abis ubah konten di admin)
    window.addEventListener('focus', loadData)
    return () => window.removeEventListener('focus', loadData)
  }, [loadData])

  // Scroll ke section berdasarkan hash URL (#tentang, #kegiatan, #informasi, #kontak)
  // Dipicu saat Home mount atau hash berubah (navigasi dari navbar halaman lain).
  const location = useLocation()
  useEffect(() => {
    const hash = location.hash
    if (!hash) {
      // Tanpa hash: kembali ke atas (perilaku standar saat buka Beranda)
      window.scrollTo({ top: 0, behavior: 'smooth' })
      return
    }
    const id = hash.replace('#', '')
    // Beri waktu render section & reveal animation selesai
    const timer = setTimeout(() => {
      const el = document.getElementById(id)
      if (el) {
        el.scrollIntoView({ behavior: 'smooth', block: 'start' })
      }
    }, 60)
    return () => clearTimeout(timer)
  }, [location.hash])

  // Global click listener to close popovers
  useEffect(() => {
    const handleGlobalClick = () => {
      setActiveShareId(null)
    }
    window.addEventListener('click', handleGlobalClick)
    return () => window.removeEventListener('click', handleGlobalClick)
  }, [])

  const triggerSlideChange = useCallback((nextIndex) => {
    setIsTransitioning(true)
    setTimeout(() => {
      setCurrentSlide(nextIndex)
      setProgress(0)
      setIsTransitioning(false)
    }, 250)
  }, [sliders])

  // Auto-play timer logic for hero slider
  useEffect(() => {
    if (sliders.length === 0) return
    const interval = setInterval(() => {
      setProgress(prev => {
        if (prev >= 100) {
          triggerSlideChange((currentSlide + 1) % sliders.length)
          return 0
        }
        return prev + 1.5
      })
    }, 100)
    return () => clearInterval(interval)
  }, [sliders, currentSlide, triggerSlideChange])

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

  const handleContactSubmit = async (e) => {
    e.preventDefault()
    
    const touchedAll = { nama: true, email: true, subjek: true, pesan: true }
    setContactTouched(touchedAll)
    
    const errors = validateContactForm()
    if (Object.keys(errors).length > 0) {
       setContactErrors(errors)
       setToast({ show: true, message: 'Validasi gagal. Periksa kembali isian form.', type: 'error' })
       return
    }
    
    setContactLoading(true)
    try {
      await kontakService.submit(contactForm)
      setContactSuccess(true)
      setContactForm({ nama: '', email: '', subjek: '', pesan: '' })
      setContactTouched({})
      setContactErrors({})
      setToast({ show: true, message: 'Pesan Anda berhasil terkirim!', type: 'success' })
      setTimeout(() => setContactSuccess(false), 5000)
    } catch (err) {
      const parsed = parseApiError(err)
      if (parsed.fieldErrors && Object.keys(parsed.fieldErrors).length > 0) {
        setContactErrors(parsed.fieldErrors)
      }
      setToast({ show: true, message: parsed.message || 'Gagal mengirim pesan.', type: 'error' })
    } finally {
      setContactLoading(false)
    }
  }

  const activeSlide = (sliders && sliders.length > 0) ? (sliders[currentSlide] || sliders[0]) : {
    title: '',
    subtitle: '',
    tag: '',
    image_url: '',
    is_new: false,
    event_date: '',
    location: '',
    link_url: ''
  }

  return (
    <PublicLayout>
      {/* 1. HERO CAROUSEL WITH 3D FLYER FOCUS */}
      <section className="relative bg-brand-950 overflow-hidden pt-28 pb-20 md:pt-40 md:pb-24 border-b border-brand-900 min-h-screen flex items-center">
        <div className="absolute inset-0 z-0">
          {activeSlide.image_url && (
            <img
              src={resolveAssetUrl(activeSlide.image_url)}
              alt={activeSlide.title}
              className={`w-full h-full object-cover transition-all duration-[750ms] blur-sm scale-110 opacity-40 ${isTransitioning ? 'opacity-0' : 'opacity-40'}`}
            />
          )}
          <div className="absolute inset-0 bg-gradient-to-r from-brand-950 via-brand-950/90 to-brand-900/60" />
          <div className="absolute inset-0 bg-texture-dots opacity-20 mix-blend-overlay" />
          {/* Floating glow orbs — hidup */}
          <div className="absolute top-[15%] left-[10%] w-72 h-72 bg-brand-500/30 blur-[120px] rounded-full pointer-events-none orb-float" />
          <div className="absolute bottom-[10%] right-[8%] w-80 h-80 bg-amber-500/15 blur-[100px] rounded-full pointer-events-none orb-float-slow" />
          <div className="absolute top-[45%] right-[30%] w-40 h-40 bg-purple-500/25 blur-[90px] rounded-full pointer-events-none orb-float" />
        </div>

        {/* Hero Content Grid */}
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10 w-full">
          <div className="flex flex-col lg:flex-row items-center gap-12 lg:gap-16">

            {/* Left: Information */}
            <div className={`w-full lg:w-1/2 text-left order-2 lg:order-1 space-y-6 transition-all duration-500 transform ${isTransitioning ? 'opacity-0 translate-y-2' : 'opacity-100 translate-y-0'}`}>
              <div className="flex flex-wrap items-center gap-3">
                {activeSlide.is_new && (
                  <span className="inline-flex items-center gap-1.5 bg-gradient-to-r from-red-600 to-red-500 text-white text-[11px] font-extrabold px-4 py-1.5 rounded-full shadow-[0_0_20px_rgba(239,68,68,0.6)] animate-pulse uppercase tracking-wider border border-red-400/50">
                    <i className="ph-fill ph-fire" /> TERBARU
                  </span>
                )}
                {activeSlide.tag && (
                  <span className="inline-flex items-center gap-1.5 bg-white/10 backdrop-blur-md border border-white/20 text-brand-100 text-[11px] font-bold px-4 py-1.5 rounded-full uppercase tracking-wider shadow-sm">
                    <i className="ph-bold ph-tag" /> {activeSlide.tag}
                  </span>
                )}
              </div>

              <h1 className="font-heading text-4xl sm:text-5xl lg:text-[3.5rem] font-extrabold text-white tracking-tight leading-[1.1] drop-shadow-md">
                {activeSlide.title}
              </h1>

              <p className="text-lg md:text-xl text-white/80 font-medium leading-relaxed max-w-xl">
                {activeSlide.subtitle}
              </p>

              {/* Event Meta Cards */}
              <div className="flex flex-col sm:flex-row gap-4">
                {activeSlide.event_date && (
                  <div className="flex items-center gap-4 bg-white/5 backdrop-blur-lg border border-white/10 rounded-2xl p-4 pr-8 shadow-sm">
                    <div className="w-12 h-12 rounded-full bg-brand-500/20 flex items-center justify-center border border-brand-400/30 text-amber-400 text-xl">
                      <i className="ph-bold ph-calendar-blank" />
                    </div>
                    <div>
                      <p className="text-[10px] text-white/60 font-bold uppercase tracking-widest mb-0.5">Waktu Event</p>
                      <p className="text-sm font-bold text-white">{activeSlide.event_date}</p>
                    </div>
                  </div>
                )}
                {activeSlide.location && (
                  <div className="flex items-center gap-4 bg-white/5 backdrop-blur-lg border border-white/10 rounded-2xl p-4 pr-8 shadow-sm">
                    <div className="w-12 h-12 rounded-full bg-brand-500/20 flex items-center justify-center border border-brand-400/30 text-amber-400 text-xl">
                      <i className="ph-bold ph-map-pin" />
                    </div>
                    <div>
                      <p className="text-[10px] text-white/60 font-bold uppercase tracking-widest mb-0.5">Lokasi</p>
                      <p className="text-sm font-bold text-white">{activeSlide.location}</p>
                    </div>
                  </div>
                )}
              </div>

              {activeSlide.link_url && (
                <div className="pt-2">
                  <a
                    href={activeSlide.link_url || '#'}
                    className="inline-flex items-center gap-2 bg-brand-600 hover:bg-brand-700 text-white px-8 py-3.5 rounded-xl font-bold text-sm transition-all shadow-sm hover:shadow-md transform hover:-translate-y-1"
                  >
                    Lihat Detail Event <i className="ph-bold ph-arrow-right text-lg" />
                  </a>
                </div>
              )}
            </div>

            {/* Right: 3D Stack Flyer Container */}
            <div className={`w-full lg:w-1/2 order-1 lg:order-2 flex justify-center lg:justify-end relative transition-all duration-500 transform ${isTransitioning ? 'opacity-0 scale-98 rotate-1' : 'opacity-100 scale-100'}`}>
              <div className="relative w-full max-w-sm lg:max-w-md aspect-[4/5] group">
                <div className="absolute inset-0 bg-brand-800/40 backdrop-blur-xl border border-white/10 rounded-3xl transform -rotate-6 scale-95 transition-transform duration-700 group-hover:-rotate-12 shadow-2xl origin-bottom-left" />
                <div className="absolute inset-0 bg-gradient-to-br from-brand-600/30 to-amber-500/20 backdrop-blur-md border border-white/20 rounded-3xl transform rotate-3 scale-100 transition-transform duration-700 group-hover:rotate-6 shadow-xl origin-bottom-right" />
                <div className="absolute inset-0 rounded-3xl overflow-hidden shadow-[0_30px_60px_rgba(0,0,0,0.6)] border-[4px] border-white/10 transform transition-transform duration-700 hover:scale-[1.03] bg-brand-950">
                  {activeSlide.image_url && (
                    <img
                      src={resolveAssetUrl(activeSlide.image_url)}
                      alt={activeSlide.title}
                      className={`w-full h-full object-cover transition-transform duration-[8000ms] ease-out ${isTransitioning ? 'scale-100' : 'scale-105'}`}
                    />
                  )}
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Carousel Controls */}
        {sliders.length > 1 && (
          <div className="absolute bottom-8 left-0 right-0 z-20">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 flex justify-between items-center">
              <div className="flex gap-3">
                <button
                  onClick={() => triggerSlideChange((currentSlide - 1 + sliders.length) % sliders.length)}
                  className="w-12 h-12 rounded-full bg-white/10 hover:bg-white/20 border border-white/20 text-white flex items-center justify-center transition backdrop-blur-sm"
                >
                  <i className="ph-bold ph-caret-left text-xl" />
                </button>
                <button
                  onClick={() => triggerSlideChange((currentSlide + 1) % sliders.length)}
                  className="w-12 h-12 rounded-full bg-white/10 hover:bg-white/20 border border-white/20 text-white flex items-center justify-center transition backdrop-blur-sm"
                >
                  <i className="ph-bold ph-caret-right text-xl" />
                </button>
              </div>
              <div className="flex items-center gap-3">
                {sliders.map((_, idx) => {
                  const isActive = currentSlide === idx
                  return (
                    <button
                      key={idx}
                      onClick={() => triggerSlideChange(idx)}
                      className="group relative h-2.5 rounded-full bg-white/20 overflow-hidden transition-all duration-300"
                      style={{ width: isActive ? '48px' : '10px' }}
                      title={`Ke slide ${idx + 1}`}
                    >
                      {isActive && (
                        <div
                          className="absolute inset-y-0 left-0 bg-gradient-to-r from-amber-400 to-amber-300 transition-all duration-100 ease-linear"
                          style={{ width: `${progress}%` }}
                        />
                      )}
                      <div className="absolute inset-0 bg-white/10 opacity-0 group-hover:opacity-100 transition-opacity" />
                    </button>
                  )
                })}
              </div>
            </div>
          </div>
        )}
      </section>

      {/* 2. SAMBUTAN KETUA UMUM */}
      <section className="py-24 bg-white border-b border-slate-200 relative overflow-hidden">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex flex-col lg:flex-row gap-16 items-center">

            {/* Poster Image */}
            <div className="w-full lg:w-5/12 relative">
              <div className="absolute inset-0 bg-brand-500/20 blur-[80px] rounded-full" />
              <div className="relative rounded-3xl overflow-hidden shadow-[0_20px_50px_rgba(0,0,0,0.15)] border border-white transform transition hover:-translate-y-2 duration-500">
                <img
                  src={resolveAssetUrl(settings.greeting_image_path) || ''}
                  alt="Poster Sambutan"
                  className="w-full h-auto object-contain"
                />
              </div>
              <div className="absolute -bottom-6 -right-6 bg-white p-4 rounded-2xl shadow-xl border border-slate-100 flex items-center gap-4">
                <div className="w-12 h-12 bg-amber-100 rounded-full flex items-center justify-center text-amber-500 text-2xl">
                  <i className="ph-fill ph-star" />
                </div>
                <div>
                  <p className="font-bold text-slate-800 text-sm">{settings.greeting_title || ''}</p>
                  <p className="text-[10px] text-slate-500 uppercase tracking-widest font-semibold">{settings.greeting_subtitle || ''}</p>
                </div>
              </div>
            </div>

            {/* Content Text */}
            <div className="w-full lg:w-7/12 space-y-8">
              <div>
                <span className="inline-flex items-center gap-2 bg-brand-50 text-brand-700 text-xs font-bold px-4 py-2 rounded-full uppercase tracking-widest mb-4 border border-brand-100">
                  <i className="ph-fill ph-quotes" /> Refleksi Resmi
                </span>
                <h2 className="font-heading text-4xl lg:text-5xl font-black text-slate-900 tracking-tight leading-tight">
                  {settings.greeting_title || ''} <span className="text-transparent bg-clip-text bg-gradient-to-r from-brand-600 to-amber-500">{settings.greeting_subtitle || ''}</span>{settings.greeting_title ? ' Bangsa' : ''}
                </h2>
              </div>

              {/* Elegant Quote Card */}
              <div className="relative bg-slate-50 rounded-2xl p-8 border border-slate-200/60 shadow-inner">
                <i className="ph-fill ph-quotes absolute -top-4 -left-2 text-6xl text-brand-200/50 transform -rotate-12" />
                <p className="relative z-10 text-slate-700 text-lg leading-relaxed font-quote italic">
                  {settings.greeting_content ? `"${settings.greeting_content}"` : ""}
                </p>
              </div>

              {/* Signature */}
              <div className="flex items-center gap-4 pt-4 border-t border-slate-100">
                <div className="flex -space-x-4">
                  {ketuaUmum?.image_path && (
                    <img src={resolveAssetUrl(ketuaUmum.image_path)} alt={ketuaUmum.name} className="w-12 h-12 rounded-full border-2 border-white shadow-md relative z-20 object-cover" />
                  )}
                  {sekjen?.name && !settings.greeting_sign_subtitle && (
                    <div className="w-12 h-12 rounded-full bg-brand-100 border-2 border-white shadow-md flex items-center justify-center text-brand-700 font-bold text-xs relative z-10">
                      {sekjen.name.split(' ').map(n => n[0]).join('').substring(0, 2).toUpperCase()}
                    </div>
                  )}
                </div>
                <div>
                  <p className="font-heading font-bold text-slate-900 text-sm">
                    {settings.greeting_sign_name || settings.site_name || ''}
                  </p>
                  <p className="text-xs text-brand-600 font-medium">
                    {settings.greeting_sign_subtitle || (
                      <>
                        {ketuaUmum?.name || ''} {sekjen?.name ? `& ${sekjen.name}` : ''}
                      </>
                    )}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* 3. TENTANG KAMI */}
      <section id="tentang" className="py-24 bg-slate-50 border-b border-slate-200 relative scroll-mt-24">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
          <div className="text-center mb-16">
            <span className="inline-block px-4 py-1.5 rounded-full bg-white border border-slate-200 text-xs font-bold text-brand-600 tracking-widest uppercase mb-4 shadow-sm">
              Profil Organisasi
            </span>
            <h2 className="font-heading text-4xl md:text-5xl font-black text-slate-900">Tentang Kami</h2>
          </div>

          <div className="flex flex-col lg:flex-row gap-16 items-center">
            {/* Foto Pimpinan */}
            {ketuaUmum && (
              <div className="w-full lg:w-1/3 flex flex-col items-center">
                <div className="relative group">
                  <div className="w-72 aspect-[4/5] rounded-3xl overflow-hidden shadow-2xl border-4 border-white mx-auto relative z-10 transform transition duration-500 group-hover:-translate-y-2">
                    <img
                      src={resolveAssetUrl(ketuaUmum.image_path) || ""}
                      alt={ketuaUmum.name}
                      className="w-full h-full object-cover"
                    />
                  </div>
                </div>
                <div className="mt-4 text-center">
                  <h4 className="font-heading font-extrabold text-slate-800 text-lg">{ketuaUmum.name}</h4>
                  {(() => {
                    const roleText = ketuaUmum.role || '';
                    const yearRegex = /(\d{4}\s*-\s*\d{4}|\d{4})/g;
                    const match = roleText.match(yearRegex);
                    let mainRole = roleText;
                    let yearPart = '';
                    if (match) {
                      yearPart = match[0];
                      mainRole = roleText.replace(yearPart, '').trim();
                    }
                    return (
                      <>
                        <p className="text-[11px] text-brand-600 font-black uppercase tracking-wider leading-snug max-w-[260px] mx-auto">{mainRole}</p>
                        {yearPart && (
                          <p className="text-[10px] text-slate-400 font-bold uppercase tracking-widest mt-0.5">{yearPart}</p>
                        )}
                      </>
                    );
                  })()}
                </div>
              </div>
            )}

            {/* Tabs Area */}
            <div className="w-full lg:w-2/3">
              {/* PILL TABS */}
              <div className="inline-flex bg-white rounded-full p-1.5 shadow-sm border border-slate-200 mb-8 max-w-full overflow-x-auto">
                <button
                  onClick={() => setAboutTab('selayang')}
                  className={`px-6 py-2.5 rounded-full font-bold text-sm transition-all whitespace-nowrap ${aboutTab === 'selayang' ? 'bg-brand-600 text-white shadow-md' : 'text-slate-500 hover:text-slate-800 hover:bg-slate-50'}`}
                >
                  Selayang Pandang
                </button>
                <button
                  onClick={() => setAboutTab('tanggal')}
                  className={`px-6 py-2.5 rounded-full font-bold text-sm transition-all whitespace-nowrap ${aboutTab === 'tanggal' ? 'bg-brand-600 text-white shadow-md' : 'text-slate-500 hover:text-slate-800 hover:bg-slate-50'}`}
                >
                  Sejarah Terbentuk
                </button>
                <button
                  onClick={() => setAboutTab('lokasi')}
                  className={`px-6 py-2.5 rounded-full font-bold text-sm transition-all whitespace-nowrap ${aboutTab === 'lokasi' ? 'bg-brand-600 text-white shadow-md' : 'text-slate-500 hover:text-slate-800 hover:bg-slate-50'}`}
                >
                  Lokasi Sekretariat
                </button>
              </div>

              {/* Tab Content */}
              <div className="bg-white/80 backdrop-blur-xl p-8 sm:p-10 rounded-3xl shadow-[0_10px_40px_rgba(0,0,0,0.04)] border border-slate-100 min-h-[300px] relative overflow-hidden">
                {aboutTab === 'selayang' && (
                  <div>
                    <h3 className="font-heading text-2xl font-bold text-slate-900 mb-6 flex items-center gap-3">
                      <i className="ph-fill ph-book-open-text text-brand-500" /> Tujuan Utama
                    </h3>
                    <div className="space-y-4 text-slate-600 text-sm sm:text-base leading-relaxed">
                      <p>{settings.history || ''}</p>
                      <p>{settings.about_tutorial || ''}</p>
                    </div>
                  </div>
                )}

                {aboutTab === 'tanggal' && (
                  <div>
                    <h3 className="font-heading text-2xl font-bold text-slate-900 mb-6 flex items-center gap-3">
                      <i className="ph-fill ph-stamp text-brand-500" /> Legalitas Resmi
                    </h3>
                    <div className="bg-brand-50 rounded-2xl p-6 border border-brand-100 flex gap-5 items-start">
                      <div className="w-14 h-14 rounded-full bg-white shadow-sm flex items-center justify-center flex-shrink-0">
                        <i className="ph-fill ph-calendar-check text-2xl text-brand-600" />
                      </div>
                      <div>
                        <h4 className="font-bold text-slate-900 text-lg mb-2">{settings.about_formation_date || ''}</h4>
                        <p className="text-slate-600 text-sm leading-relaxed mb-4">Secara resmi GRADASI disahkan melalui Surat Keputusan (SK) Kementerian Hukum dan HAM Republik Indonesia.</p>
                        <div className="inline-flex items-center gap-2 bg-white px-4 py-2 rounded-lg border border-brand-200 shadow-sm text-xs font-mono font-bold text-brand-800">
                          <i className="ph-bold ph-certificate" /> {settings.about_no_sk || ''}
                        </div>
                      </div>
                    </div>
                  </div>
                )}

                {aboutTab === 'lokasi' && (
                  <div>
                    <h3 className="font-heading text-2xl font-bold text-slate-900 mb-6 flex items-center gap-3">
                      <i className="ph-fill ph-map-pin-line text-brand-500" /> Markas Pusat
                    </h3>
                    <div className="bg-slate-50 rounded-2xl p-6 border border-slate-200">
                      <h4 className="font-bold text-slate-900 text-base mb-3">Kantor Sekretariat {settings.site_name || ''}</h4>
                      <div className="flex gap-3 text-slate-600 text-sm leading-relaxed">
                        <i className="ph-bold ph-map-pin text-brand-600 mt-1" />
                        <p>{settings.address}</p>
                      </div>
                    </div>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* 4. VISI & MISI BENTO GRID */}
      <section className="py-24 bg-slate-50 border-b border-slate-200">
        <div ref={revealVisimisi} className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 reveal">
          <div className="text-center mb-16">
            <span className="inline-block px-4 py-1.5 rounded-full bg-brand-100 border border-brand-200 text-xs font-bold text-brand-700 tracking-widest uppercase mb-4 shadow-sm">
              Tujuan & Arah Organisasi
            </span>
            <h2 className="font-heading text-4xl md:text-5xl font-black text-slate-900">Visi & Misi</h2>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
            {/* Left: Visi & Logo */}
            <div className="lg:col-span-5 space-y-6 lg:sticky lg:top-28">
              {/* VISI UTAMA Box */}
              <div className="bg-gradient-to-br from-brand-700 to-brand-900 rounded-3xl p-8 sm:p-10 text-white relative overflow-hidden group shadow-xl">
                <i className="ph-fill ph-quotes absolute -bottom-10 -right-10 text-[180px] text-white/5 -rotate-12 pointer-events-none" />
                <div className="relative z-10">
                  <div className="flex items-center gap-3 mb-6">
                    <div className="w-10 h-10 bg-white/20 backdrop-blur-sm rounded-full flex items-center justify-center border border-white/30">
                      <i className="ph-fill ph-eye text-xl" />
                    </div>
                    <span className="text-xs font-bold uppercase tracking-widest text-brand-200">Visi Utama</span>
                  </div>
                  <p className="text-xl sm:text-2xl font-heading font-bold leading-relaxed drop-shadow-md italic">
                    "{settings.about_vision}"
                  </p>
                </div>
              </div>

              {/* Logo Display Box */}
              <div className="bg-white rounded-3xl p-8 flex items-center justify-center border border-slate-100 shadow-md">
                <img src={resolveAssetUrl(settings.logo_path)} alt="Logo" className="w-44 h-auto object-contain filter drop-shadow-lg" />
              </div>
            </div>

            {/* Right: Dynamic Missions */}
            <div className="lg:col-span-7">
              <div className="bg-white rounded-3xl p-6 sm:p-8 border border-slate-100 shadow-md">
                <h3 className="font-heading text-xl font-extrabold text-slate-900 mb-8 flex items-center gap-2">
                  <i className="ph-bold ph-list-numbers text-brand-600 text-2xl" /> Misi Organisasi
                </h3>

                <div className="space-y-6">
                  {getMissions().map((mission, idx) => (
                    <div
                      key={idx}
                      className="flex gap-4 items-start group relative"
                    >
                      {/* Stepper Circle */}
                      <div className="flex flex-col items-center flex-shrink-0">
                        <div className={`w-10 h-10 rounded-full flex items-center justify-center font-bold text-sm shadow-sm transition-all duration-300 group-hover:scale-110 ${idx % 3 === 0 ? 'bg-blue-50 text-blue-600 border border-blue-200' : idx % 3 === 1 ? 'bg-teal-50 text-teal-600 border border-teal-200' : 'bg-amber-50 text-amber-600 border border-amber-200'
                          }`}>
                          {idx + 1}
                        </div>
                        {idx < getMissions().length - 1 && (
                          <div className="w-[1.5px] bg-slate-100 flex-grow h-12 mt-2" />
                        )}
                      </div>

                      {/* Mission Card */}
                      <div className="flex-1 bg-slate-50/40 hover:bg-slate-50 group-hover:translate-x-1 border border-slate-200/40 rounded-2xl p-5 transition-all duration-300">
                        <h4 className="font-bold text-slate-800 text-xs mb-1.5 uppercase tracking-wider">Misi Ke-{idx + 1}</h4>
                        <p className="text-slate-600 text-sm leading-relaxed">{mission}</p>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* 5. BERBAGAI EVENT MENARIK */}
      <section id="kegiatan" className="py-20 bg-slate-50 border-b border-slate-200 scroll-mt-24">
        <div ref={revealKegiatan} className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 reveal">
          <div className="text-center mb-12">
            <span className="text-xs font-bold text-brand-600 tracking-widest uppercase mb-2 block">Aktivitas Terbaru</span>
            <h2 className="font-heading text-3xl font-bold text-slate-900">Berbagai Event Menarik</h2>
          </div>

          <div className="max-w-[1098px] mx-auto relative group">
            <button
              onClick={() => kegiatanSliderRef.current?.scrollBy({ left: -380, behavior: 'smooth' })}
              className="absolute -left-12 lg:-left-16 top-1/2 -translate-y-1/2 z-20 w-12 h-12 bg-white/95 backdrop-blur rounded-full shadow-[0_5px_15px_rgba(0,0,0,0.15)] border border-slate-100 text-brand-600 hover:bg-brand-50 hidden md:flex items-center justify-center transition hover:scale-110"
            >
              <i className="ph-bold ph-caret-left text-xl" />
            </button>
            <button
              onClick={() => kegiatanSliderRef.current?.scrollBy({ left: 380, behavior: 'smooth' })}
              className="absolute -right-12 lg:-right-16 top-1/2 -translate-y-1/2 z-20 w-12 h-12 bg-white/95 backdrop-blur rounded-full shadow-[0_5px_15px_rgba(0,0,0,0.15)] border border-slate-100 text-brand-600 hover:bg-brand-50 hidden md:flex items-center justify-center transition hover:scale-110"
            >
              <i className="ph-bold ph-caret-right text-xl" />
            </button>

            <div ref={kegiatanSliderRef} className="flex overflow-x-auto gap-6 pb-6 snap-x snap-mandatory hide-scrollbar">
              {featuredKegiatan.map((item, idx) => (
                <div key={item.id || idx} className="w-[85vw] sm:w-[350px] flex-shrink-0 snap-center card-minimal overflow-hidden flex flex-col group cursor-pointer">
                  <div className="image-zoom-container h-56 relative">
                    <img
                      src={item.image_url ? resolveAssetUrl(item.image_url) : 'https://images.unsplash.com/photo-1540575467063-178a50c2df87?q=80&w=600'}
                      alt={item.title}
                      className="w-full h-full object-cover"
                    />
                    {idx === 0 && (
                      <div className="absolute top-4 right-4 bg-red-500 text-white text-[10px] font-extrabold px-3 py-1.5 rounded animate-pulse shadow-[0_0_10px_rgba(239,68,68,0.5)] z-10 tracking-wider flex items-center gap-1">
                        <i className="ph-fill ph-fire" /> TERBARU
                      </div>
                    )}
                    <div className="absolute top-4 left-4 bg-white/90 backdrop-blur-sm text-brand-700 text-xs font-bold px-3 py-1.5 rounded">
                      {item.category || 'Event'}
                    </div>
                  </div>
                  <div className="p-6 flex flex-col flex-grow bg-white">
                    <p className="text-slate-500 text-[11px] font-semibold tracking-wider uppercase mb-3 flex items-center gap-1.5">
                      <i className="ph-bold ph-calendar-blank text-brand-500" /> {item.event_date || '31 December 2025'}
                    </p>
                    <Link to={`/kegiatan/${item.slug}`} className="font-heading text-lg font-bold text-slate-900 mb-2 group-hover:text-brand-600 transition line-clamp-2">
                      {item.title}
                    </Link>
                    <p className="text-slate-600 text-sm flex-grow line-clamp-3 mb-4 leading-relaxed">{item.excerpt}</p>
                    <div className="border-t border-slate-100 pt-4 flex justify-between items-center mt-auto">
                      <div className="relative">
                        <button
                          onClick={(e) => {
                            e.preventDefault()
                            e.stopPropagation()
                            setActiveShareId(activeShareId === `kegiatan-${item.id}` ? null : `kegiatan-${item.id}`)
                          }}
                          className={`w-8 h-8 rounded-full flex items-center justify-center transition-all duration-300 ${activeShareId === `kegiatan-${item.id}` ? 'bg-brand-600 text-white' : 'bg-slate-100 text-slate-500 hover:bg-brand-50 hover:text-brand-600'}`}
                          title="Bagikan"
                        >
                          <i className="ph-bold ph-share-network text-sm" />
                        </button>

                        {activeShareId === `kegiatan-${item.id}` && (
                          <div
                            className="absolute bottom-full left-0 mb-2 bg-white/95 backdrop-blur-md border border-slate-200/80 p-2.5 rounded-2xl shadow-xl z-30 flex gap-2 animate-scale-in"
                            onClick={(e) => e.stopPropagation()}
                          >
                            <button
                              onClick={(e) => {
                                e.preventDefault();
                                e.stopPropagation();
                                handleInstagramShare(`${window.location.origin}/kegiatan/${item.slug}`);
                              }}
                              className="w-8 h-8 rounded-xl bg-[#E1306C]/10 text-[#E1306C] hover:bg-[#E1306C] hover:text-white flex items-center justify-center transition"
                              title="Instagram"
                            >
                              <i className="ph-fill ph-instagram-logo text-sm" />
                            </button>
                            <a
                              href={getShareUrl('whatsapp', { title: item.title, text: item.excerpt, url: `${window.location.origin}/kegiatan/${item.slug}` })}
                              target="_blank"
                              rel="noopener noreferrer"
                              onClick={(e) => e.stopPropagation()}
                              className="w-8 h-8 rounded-xl bg-[#25D366]/10 text-[#25D366] hover:bg-[#25D366] hover:text-white flex items-center justify-center transition"
                              title="WhatsApp"
                            >
                              <i className="ph-fill ph-whatsapp-logo text-sm" />
                            </a>
                            <a
                              href={getShareUrl('facebook', { title: item.title, text: item.excerpt, url: `${window.location.origin}/kegiatan/${item.slug}` })}
                              target="_blank"
                              rel="noopener noreferrer"
                              onClick={(e) => e.stopPropagation()}
                              className="w-8 h-8 rounded-xl bg-[#1877F2]/10 text-[#1877F2] hover:bg-[#1877F2] hover:text-white flex items-center justify-center transition"
                              title="Facebook"
                            >
                              <i className="ph-fill ph-facebook-logo text-sm" />
                            </a>
                            <button
                              onClick={(e) => {
                                e.preventDefault();
                                e.stopPropagation();
                                handleCopyLink(`${window.location.origin}/kegiatan/${item.slug}`, 'Tautan kegiatan');
                              }}
                              className="w-8 h-8 rounded-xl bg-slate-100 text-slate-600 hover:bg-slate-600 hover:text-white flex items-center justify-center transition"
                              title="Salin Tautan"
                            >
                              <i className="ph-bold ph-link text-sm" />
                            </button>
                          </div>
                        )}
                      </div>
                      <Link to={`/kegiatan/${item.slug}`} className="flex items-center gap-1.5 text-xs font-bold text-brand-600 hover:text-brand-800 transition">
                        Baca Selengkapnya <i className="ph-bold ph-arrow-right" />
                      </Link>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>

          <div className="text-center mt-8">
            <Link to="/kegiatan" className="inline-flex items-center gap-2 bg-white border border-slate-300 hover:border-brand-500 hover:text-brand-700 text-slate-700 px-8 py-3 rounded-lg font-semibold text-sm transition shadow-sm">
              Lihat Seluruh Kegiatan <i className="ph-bold ph-arrow-right" />
            </Link>
          </div>
        </div>
      </section>

      {/* 6. VIDEO PROFILE ORGANISASI (CINEMATIC STYLE) */}
      <section className="py-28 bg-slate-900 border-b border-slate-800 relative overflow-hidden flex items-center justify-center">
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-3/4 h-3/4 bg-brand-600/20 blur-[150px] rounded-full pointer-events-none z-0" />
        <div className="max-w-6xl w-full mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
          <div className="text-center mb-12">
            <span className="inline-flex items-center gap-2 px-4 py-1.5 rounded-full bg-slate-800 border border-slate-700 text-xs font-bold text-brand-400 tracking-widest uppercase mb-4 shadow-sm">
              <i className="ph-fill ph-film-strip" /> Dokumentasi Visual
            </span>
            <h2 className="font-heading text-4xl md:text-5xl font-black text-white">Video Profile Organisasi</h2>
          </div>

          <div className="relative group mx-auto">
            <div className="aspect-video bg-black rounded-3xl overflow-hidden shadow-[0_30px_60px_rgba(0,0,0,0.8)] relative ring-1 ring-white/10">
              {!isPlayingVideo && (
                <div className="absolute inset-0 z-10 flex flex-col items-center justify-center bg-black/60 transition-opacity duration-500">
                  <img
                    src="https://images.unsplash.com/photo-1475721025871-872ba5bb619a?q=80&w=2070&auto=format&fit=crop"
                    alt="Video Cover"
                    className="absolute inset-0 w-full h-full object-cover mix-blend-overlay opacity-60"
                  />
                  <button
                    onClick={() => { if (videoRef.current) { videoRef.current.play(); setIsPlayingVideo(true); } }}
                    className="relative z-20 w-24 h-24 bg-brand-600 hover:bg-brand-500 text-white rounded-full flex items-center justify-center text-4xl shadow-[0_0_40px_rgba(37,99,235,0.8)] transform transition-all hover:scale-110 border-4 border-white/20"
                  >
                    <i className="ph-fill ph-play ml-2" />
                  </button>
                </div>
              )}
              <video
                ref={videoRef}
                onPause={() => setIsPlayingVideo(false)}
                onPlay={() => setIsPlayingVideo(true)}
                className="w-full h-full object-cover"
                controls
                preload="metadata"
              >
                <source src={resolveAssetUrl(settings.video_profile_path) || ''} type="video/mp4" />
                Browser Anda tidak mendukung pemutar video.
              </video>
            </div>
          </div>
        </div>
      </section>

      {/* 7. BERITA TERBARU */}
      <section id="informasi" className="py-20 bg-slate-50 border-b border-slate-200 scroll-mt-24">
        <div ref={revealBerita} className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 reveal">
          <div className="text-center mb-12">
            <span className="text-xs font-bold text-brand-600 tracking-widest uppercase mb-2 block">Pusat Informasi</span>
            <h2 className="font-heading text-3xl font-bold text-slate-900">Berita Terbaru</h2>
          </div>

          <div className="max-w-[1098px] mx-auto relative group">
            <button
              onClick={() => beritaSliderRef.current?.scrollBy({ left: -380, behavior: 'smooth' })}
              className="absolute -left-12 lg:-left-16 top-1/2 -translate-y-1/2 z-20 w-12 h-12 bg-white/95 backdrop-blur rounded-full shadow-[0_5px_15px_rgba(0,0,0,0.15)] border border-slate-100 text-brand-600 hover:bg-brand-50 hidden md:flex items-center justify-center transition hover:scale-110"
            >
              <i className="ph-bold ph-caret-left text-xl" />
            </button>
            <button
              onClick={() => beritaSliderRef.current?.scrollBy({ left: 380, behavior: 'smooth' })}
              className="absolute -right-12 lg:-right-16 top-1/2 -translate-y-1/2 z-20 w-12 h-12 bg-white/95 backdrop-blur rounded-full shadow-[0_5px_15px_rgba(0,0,0,0.15)] border border-slate-100 text-brand-600 hover:bg-brand-50 hidden md:flex items-center justify-center transition hover:scale-110"
            >
              <i className="ph-bold ph-caret-right text-xl" />
            </button>

            <div ref={beritaSliderRef} className="flex overflow-x-auto gap-6 pb-6 snap-x snap-mandatory hide-scrollbar">
              {recentBerita.map((item, idx) => (
                <div key={item.id || idx} className="w-[85vw] sm:w-[350px] flex-shrink-0 snap-center card-minimal overflow-hidden flex flex-col group cursor-pointer">
                  <div className="image-zoom-container h-48 relative">
                    <img
                      src={item.image_url ? resolveAssetUrl(item.image_url) : 'https://images.unsplash.com/photo-1504711434969-e33886168f5c?q=80&w=600'}
                      alt={item.title}
                      className="w-full h-full object-cover"
                    />
                    {idx === 0 && (
                      <div className="absolute top-4 right-4 bg-red-500 text-white text-[10px] font-extrabold px-3 py-1.5 rounded animate-pulse shadow-[0_0_10px_rgba(239,68,68,0.5)] z-10 tracking-wider flex items-center gap-1">
                        <i className="ph-fill ph-fire" /> TERBARU
                      </div>
                    )}
                  </div>
                  <div className="p-6 flex flex-col flex-grow bg-white">
                    <p className="text-brand-600 text-[11px] font-bold tracking-widest uppercase mb-3">
                      {item.published_date || '11 February 2026'}
                    </p>
                    <Link to={`/berita/${item.slug}`} className="font-heading text-lg font-bold text-slate-900 mb-3 group-hover:text-brand-600 transition line-clamp-2">
                      {item.title}
                    </Link>
                    <p className="text-slate-600 text-sm flex-grow line-clamp-2 mb-4 leading-relaxed">{item.excerpt}</p>
                    <div className="border-t border-slate-100 pt-4 flex justify-between items-center mt-auto">
                      <div className="relative">
                        <button
                          onClick={(e) => {
                            e.preventDefault()
                            e.stopPropagation()
                            setActiveShareId(activeShareId === `berita-${item.id}` ? null : `berita-${item.id}`)
                          }}
                          className={`w-8 h-8 rounded-full flex items-center justify-center transition-all duration-300 ${activeShareId === `berita-${item.id}` ? 'bg-brand-600 text-white' : 'bg-slate-100 text-slate-500 hover:bg-brand-50 hover:text-brand-600'}`}
                          title="Bagikan"
                        >
                          <i className="ph-bold ph-share-network text-sm" />
                        </button>

                        {activeShareId === `berita-${item.id}` && (
                          <div
                            className="absolute bottom-full left-0 mb-2 bg-white/95 backdrop-blur-md border border-slate-200/80 p-2.5 rounded-2xl shadow-xl z-30 flex gap-2 animate-scale-in"
                            onClick={(e) => e.stopPropagation()}
                          >
                            <button
                              onClick={(e) => {
                                e.preventDefault();
                                e.stopPropagation();
                                handleInstagramShare(`${window.location.origin}/berita/${item.slug}`);
                              }}
                              className="w-8 h-8 rounded-xl bg-[#E1306C]/10 text-[#E1306C] hover:bg-[#E1306C] hover:text-white flex items-center justify-center transition"
                              title="Instagram"
                            >
                              <i className="ph-fill ph-instagram-logo text-sm" />
                            </button>
                            <a
                              href={getShareUrl('whatsapp', { title: item.title, text: item.excerpt, url: `${window.location.origin}/berita/${item.slug}` })}
                              target="_blank"
                              rel="noopener noreferrer"
                              onClick={(e) => e.stopPropagation()}
                              className="w-8 h-8 rounded-xl bg-[#25D366]/10 text-[#25D366] hover:bg-[#25D366] hover:text-white flex items-center justify-center transition"
                              title="WhatsApp"
                            >
                              <i className="ph-fill ph-whatsapp-logo text-sm" />
                            </a>
                            <a
                              href={getShareUrl('facebook', { title: item.title, text: item.excerpt, url: `${window.location.origin}/berita/${item.slug}` })}
                              target="_blank"
                              rel="noopener noreferrer"
                              onClick={(e) => e.stopPropagation()}
                              className="w-8 h-8 rounded-xl bg-[#1877F2]/10 text-[#1877F2] hover:bg-[#1877F2] hover:text-white flex items-center justify-center transition"
                              title="Facebook"
                            >
                              <i className="ph-fill ph-facebook-logo text-sm" />
                            </a>
                            <button
                              onClick={(e) => {
                                e.preventDefault();
                                e.stopPropagation();
                                handleCopyLink(`${window.location.origin}/berita/${item.slug}`, 'Tautan berita');
                              }}
                              className="w-8 h-8 rounded-xl bg-slate-100 text-slate-600 hover:bg-slate-600 hover:text-white flex items-center justify-center transition"
                              title="Salin Tautan"
                            >
                              <i className="ph-bold ph-link text-sm" />
                            </button>
                          </div>
                        )}
                      </div>
                      <Link to={`/berita/${item.slug}`} className="flex items-center gap-1.5 text-xs font-bold text-brand-600 hover:text-brand-800 transition">
                        Baca Berita <i className="ph-bold ph-arrow-right" />
                      </Link>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>

          <div className="text-center mt-8">
            <Link to="/berita" className="inline-flex items-center gap-2 bg-white border border-slate-300 hover:border-brand-500 hover:text-brand-700 text-slate-700 px-8 py-3 rounded-lg font-semibold text-sm transition shadow-sm">
              Informasi Lainnya <i className="ph-bold ph-newspaper" />
            </Link>
          </div>
        </div>
      </section>

      {/* 8. KONTAK (MAP & FORM) */}
      <section id="kontak" className="py-20 bg-white scroll-mt-24">
        <div ref={revealKontak} className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 reveal">
          <div className="bg-slate-50 p-6 md:p-8 rounded-2xl border border-slate-200 shadow-sm">
            <div className="grid lg:grid-cols-2 gap-10">
              {/* Google Maps Embed */}
              <div className="h-[400px] relative bg-white border border-slate-200 rounded-xl overflow-hidden shadow-inner">
                <iframe
                  src={resolveAssetUrl(settings.maps_embed_url) || ""}
                  width="100%"
                  height="100%"
                  style={{ border: 0 }}
                  allowFullScreen=""
                  loading="lazy"
                />
              </div>

              {/* Contact Form */}
              <div className="flex flex-col justify-center">
                <div className="mb-6">
                  <h3 className="font-heading text-2xl font-bold text-slate-900 mb-2">Kirim Pesan Langsung</h3>
                  <p className="text-sm text-slate-500">Gunakan formulir di bawah ini untuk mengirimkan pertanyaan atau permohonan kerjasama secara resmi.</p>
                </div>

                {contactSuccess && (
                  <div className="mb-4 p-4 bg-emerald-50 border border-emerald-200 text-emerald-700 rounded-xl text-sm font-bold flex items-center gap-2">
                    <i className="ph-bold ph-check-circle text-lg" /> Pesan Anda telah berhasil terkirim! Terima kasih.
                  </div>
                )}

                <form onSubmit={handleContactSubmit} className="space-y-5">
                  <div className="grid grid-cols-2 gap-5">
                    <div>
                      <label className="block text-xs font-semibold text-slate-700 mb-1.5 uppercase tracking-wide">Nama Lengkap</label>
                      <input
                        type="text"
                        required
                        value={contactForm.nama}
                        onChange={e => handleContactChange('nama', e.target.value)}
                        onBlur={() => handleContactBlur('nama')}
                        className={`w-full px-4 py-3 rounded-lg border ${
                          contactTouched.nama && contactErrors.nama ? 'border-red-500 focus:ring-1 focus:ring-red-500' : 'border-slate-200 focus:border-brand-500'
                        } text-sm focus:outline-none bg-white transition`}
                      />
                      {contactTouched.nama && contactErrors.nama && (
                        <p className="text-red-500 text-[11px] font-semibold mt-1 flex items-center gap-1">
                          <i className="ph-bold ph-warning-circle text-xs" /> {contactErrors.nama}
                        </p>
                      )}
                    </div>
                    <div>
                      <label className="block text-xs font-semibold text-slate-700 mb-1.5 uppercase tracking-wide">Email</label>
                      <input
                        type="email"
                        required
                        value={contactForm.email}
                        onChange={e => handleContactChange('email', e.target.value)}
                        onBlur={() => handleContactBlur('email')}
                        className={`w-full px-4 py-3 rounded-lg border ${
                          contactTouched.email && contactErrors.email ? 'border-red-500 focus:ring-1 focus:ring-red-500' : 'border-slate-200 focus:border-brand-500'
                        } text-sm focus:outline-none bg-white transition`}
                      />
                      {contactTouched.email && contactErrors.email && (
                        <p className="text-red-500 text-[11px] font-semibold mt-1 flex items-center gap-1">
                          <i className="ph-bold ph-warning-circle text-xs" /> {contactErrors.email}
                        </p>
                      )}
                    </div>
                  </div>
                  <div>
                    <label className="block text-xs font-semibold text-slate-700 mb-1.5 uppercase tracking-wide">Subjek / Perihal <span className="text-gray-400 font-normal text-[10px]">(opsional)</span></label>
                    <input
                      type="text"
                      value={contactForm.subjek}
                      onChange={e => handleContactChange('subjek', e.target.value)}
                      className="w-full px-4 py-3 rounded-lg border border-slate-200 text-sm focus:border-brand-500 focus:outline-none bg-white transition"
                    />
                  </div>
                  <div>
                    <label className="block text-xs font-semibold text-slate-700 mb-1.5 uppercase tracking-wide">Pesan Anda</label>
                    <textarea
                      rows="4"
                      required
                      value={contactForm.pesan}
                      onChange={e => handleContactChange('pesan', e.target.value)}
                      onBlur={() => handleContactBlur('pesan')}
                      className={`w-full px-4 py-3 rounded-lg border ${
                        contactTouched.pesan && contactErrors.pesan ? 'border-red-500 focus:ring-1 focus:ring-red-500' : 'border-slate-200 focus:border-brand-500'
                      } text-sm focus:outline-none bg-white transition resize-none`}
                    />
                    {contactTouched.pesan && contactErrors.pesan && (
                      <p className="text-red-500 text-[11px] font-semibold mt-1 flex items-center gap-1">
                        <i className="ph-bold ph-warning-circle text-xs" /> {contactErrors.pesan}
                      </p>
                    )}
                  </div>
                  <button
                    type="submit"
                    disabled={contactLoading || !contactForm.nama?.trim() || !contactForm.email?.trim() || !contactForm.pesan?.trim()}
                    className="bg-brand-600 hover:bg-brand-700 text-white px-8 py-3.5 rounded-xl text-sm font-bold transition shadow-sm hover:shadow flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    <i className="ph-bold ph-paper-plane-right text-lg" />
                    <span>{contactLoading ? 'Mengirim...' : 'Kirim Pesan Sekarang'}</span>
                  </button>
                </form>
              </div>
            </div>
          </div>
        </div>
      </section>

      <ToastNotification
        show={toast.show}
        message={toast.message}
        type={toast.type}
        onClose={() => setToast(prev => ({ ...prev, show: false }))}
      />
    </PublicLayout>
  )
}
