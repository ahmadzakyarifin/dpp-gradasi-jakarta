import { useState, useEffect, useRef, useCallback } from 'react'
import { Link, useLocation } from 'react-router-dom'
import PublicLayout from '../layouts/PublicLayout'
import { useSettings } from '../context/SettingsContext'
import { resolveAssetUrl } from '../utils/assetUrl'
import { slidersService } from '../services/slidersService'
import { kegiatanService } from '../services/kegiatanService'
import { beritaService } from '../services/beritaService'
import { kontakService } from '../services/kontakService'
import useReveal from '../hooks/useReveal'

export default function Home() {
  const { settings } = useSettings()
  const revealVisimisi = useReveal()
  const revealKegiatan = useReveal()
  const revealBerita = useReveal()
  const revealKontak = useReveal()

  const [sliders, setSliders] = useState([
    {
      title: 'Musyawarah Nasional Ke-II GRADASI',
      subtitle: 'Kolaborasi Membangun Negeri Menuju Indonesia Emas 2045',
      tag: 'HEADLINE EVENT',
      is_new: true,
      event_date: '10 - 12 Agustus 2026',
      location: 'Jakarta Convention Center',
      image_url: 'https://gradasi.org/uploads/img/slider/1749385864.jpg',
      link_url: '/kegiatan/munas-ke-ii'
    },
    {
      title: 'Pelatihan Digital Marketing UMKM',
      subtitle: 'Program Transformasi Berbasis Literasi Digital untuk Memajukan Ekonomi Daerah',
      tag: 'PELATIHAN',
      is_new: false,
      event_date: '15 November 2025',
      location: 'Online via Zoom',
      image_url: 'https://gradasi.org/uploads/img/slider/1746600828.jpg',
      link_url: '/kegiatan/pelatihan-digital-marketing'
    },
    {
      title: 'Audiensi & Sinergi Strategis',
      subtitle: 'Memperkuat Kemitraan Bersama Pemangku Kepentingan Nasional',
      tag: 'KEMITRAAN',
      is_new: false,
      event_date: 'September 2025',
      location: 'Kementerian Kominfo',
      image_url: 'https://gradasi.org/uploads/img/slider/1746600520.png',
      link_url: '#'
    }
  ])

  const [featuredKegiatan, setFeaturedKegiatan] = useState([
    {
      id: 1,
      title: 'Penyaluran Bantuan Kemanusiaan',
      slug: 'penyaluran-bantuan-kemanusiaan',
      category: 'Nasional',
      is_new: true,
      event_date: '31 December 2025',
      image_url: 'https://gradasi.org/uploads/img/event/1767154719.jpg',
      excerpt: 'Dewan Pimpinan Pusat (DPP) GRADASI menyalurkan bantuan kemanusiaan kepada korban bencana alam...'
    },
    {
      id: 2,
      title: 'Pelatihan Digital Marketing UMKM',
      slug: 'pelatihan-digital-marketing-umkm',
      category: 'Jawa Timur',
      is_new: false,
      event_date: '31 December 2025',
      image_url: 'https://gradasi.org/uploads/img/event/1767154619.jpg',
      excerpt: 'Dewan Pimpinan Pusat (DPP) GRADASI menyalurkan bantuan kemanusiaan kepada korban bencana b...'
    },
    {
      id: 3,
      title: 'Konsolidasi Pengurus DPP & DPD',
      slug: 'konsolidasi-pengurus-dpp-dpd',
      category: 'Lampung',
      is_new: false,
      event_date: '31 December 2025',
      image_url: 'https://gradasi.org/uploads/img/event/1767154397.jpg',
      excerpt: 'Dewan Pimpinan Pusat (DPP) GRADASI menyalurkan bantuan kemanusiaan kepada korban bencana b...'
    }
  ])

  // Default Fallback Berita matching index.html
  const [recentBerita, setRecentBerita] = useState([
    {
      id: 1,
      title: 'Rapat Kerja Daerah Jatim',
      slug: 'rapat-kerja-daerah-jatim',
      category: 'Berita Daerah',
      is_new: true,
      published_date: '11 February 2026',
      image_url: 'https://gradasi.org/uploads/img/berita/17708152730.jpg',
      excerpt: 'SURABAYA, Generasi Digital Indonesia (GRADASI) Jawa Timur bersiap menggelar Rapat Kerja Da...'
    },
    {
      id: 2,
      title: 'Peningkatan Kompetensi SDM',
      slug: 'peningkatan-kompetensi-sdm-pendidikan',
      category: 'Edukasi',
      is_new: false,
      published_date: '2 November 2025',
      image_url: 'https://gradasi.org/uploads/img/berita/17620765070.jpg',
      excerpt: 'Inisiatif GRADASI Mendorong Peningkatan Kompetensi SDM Pendidikan dalam Memanfaatkan Kecer...'
    },
    {
      id: 3,
      title: 'Rumusan Kunci Kebijakan',
      slug: 'rumusan-kunci-kebijakan-literasi-digital',
      category: 'Berita Utama',
      is_new: false,
      published_date: '31 October 2025',
      image_url: 'https://gradasi.org/uploads/img/berita/17618789900.jpg',
      excerpt: '#Ketua Dewan Pakar GRADASI, Damar Juniarto, Paparkan Lima Rumusan Kunci Kebijakan untuk Me...'
    }
  ])

  const [currentSlide, setCurrentSlide] = useState(0)
  const [progress, setProgress] = useState(0)
  const [aboutTab, setAboutTab] = useState('selayang')

  // Contact Form State
  const [contactForm, setContactForm] = useState({ nama: '', email: '', subjek: '', pesan: '' })
  const [contactSuccess, setContactSuccess] = useState(false)
  const [contactLoading, setContactLoading] = useState(false)

  // Floating Headset Contact Toggle
  const [floatingOpen, setFloatingOpen] = useState(false)

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
      }).catch(() => {})

    // Load Kegiatan
    kegiatanService.list()
      .then(res => {
        if (res.success && res.data) {
          const list = Array.isArray(res.data) ? res.data : (res.data.kegiatan || [])
          setFeaturedKegiatan(list.map(k => ({ ...k, image_url: k.image_path || k.image_url })))
        }
      }).catch(() => {})

    // Load Berita
    beritaService.list()
      .then(res => {
        if (res.success && res.data) {
          const list = Array.isArray(res.data) ? res.data : (res.data.berita || [])
          setRecentBerita(list.map(b => ({ ...b, image_url: b.image_path || b.image_url })))
        }
      }).catch(() => {})
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

  // Auto-play timer logic for hero slider
  useEffect(() => {
    if (sliders.length === 0) return
    const interval = setInterval(() => {
      setProgress(prev => {
        if (prev >= 100) {
          setCurrentSlide(curr => (curr + 1) % sliders.length)
          return 0
        }
        return prev + 2
      })
    }, 100)
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

  const handleContactSubmit = async (e) => {
    e.preventDefault()
    if (!contactForm.nama || !contactForm.email || !contactForm.pesan) return
    setContactLoading(true)
    try {
      await kontakService.submit(contactForm)
      setContactSuccess(true)
      setContactForm({ nama: '', email: '', subjek: '', pesan: '' })
      setTimeout(() => setContactSuccess(false), 5000)
    } catch (err) {
      alert('Gagal mengirim pesan: ' + (err.response?.data?.message || err.message))
    } finally {
      setContactLoading(false)
    }
  }

  const activeSlide = sliders[currentSlide] || sliders[0]

  return (
    <PublicLayout>
      {/* 1. HERO CAROUSEL WITH 3D FLYER FOCUS */}
      <section className="relative bg-brand-950 overflow-hidden pt-28 pb-20 md:pt-40 md:pb-24 border-b border-brand-900 min-h-screen flex items-center">
        {/* Slide Background Overlay */}
        <div className="absolute inset-0 z-0">
          <img 
            src={resolveAssetUrl(activeSlide.image_url)} 
            alt={activeSlide.title} 
            className="w-full h-full object-cover transition-all duration-700 blur-sm scale-110 opacity-40"
          />
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
            <div className="w-full lg:w-1/2 text-left order-2 lg:order-1 space-y-6">
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

              <h1 className="font-heading text-4xl sm:text-5xl lg:text-[4rem] font-black text-white tracking-tight leading-tight drop-shadow-2xl">
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

              <div className="pt-2">
                <a 
                  href={activeSlide.link_url || '#'} 
                  className="inline-flex items-center gap-2 bg-brand-600 hover:bg-brand-700 text-white px-8 py-3.5 rounded-xl font-bold text-sm transition-all shadow-sm hover:shadow-md transform hover:-translate-y-1"
                >
                  Lihat Detail Event <i className="ph-bold ph-arrow-right text-lg" />
                </a>
              </div>
            </div>

            {/* Right: 3D Stack Flyer Container */}
            <div className="w-full lg:w-1/2 order-1 lg:order-2 flex justify-center lg:justify-end relative">
              <div className="relative w-full max-w-sm lg:max-w-md aspect-[4/5] group">
                <div className="absolute inset-0 bg-brand-800/40 backdrop-blur-xl border border-white/10 rounded-3xl transform -rotate-6 scale-95 transition-transform duration-700 group-hover:-rotate-12 shadow-2xl origin-bottom-left" />
                <div className="absolute inset-0 bg-gradient-to-br from-brand-600/30 to-amber-500/20 backdrop-blur-md border border-white/20 rounded-3xl transform rotate-3 scale-100 transition-transform duration-700 group-hover:rotate-6 shadow-xl origin-bottom-right" />
                <div className="absolute inset-0 rounded-3xl overflow-hidden shadow-[0_30px_60px_rgba(0,0,0,0.6)] border-[4px] border-white/10 transform transition-transform duration-700 hover:scale-[1.03] bg-brand-950">
                  <img 
                    src={resolveAssetUrl(activeSlide.image_url)} 
                    alt={activeSlide.title} 
                    className="w-full h-full object-cover"
                  />
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Carousel Controls */}
        <div className="absolute bottom-8 left-0 right-0 z-20">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 flex justify-between items-center">
            <div className="flex gap-3">
              <button 
                onClick={() => { setCurrentSlide((currentSlide - 1 + sliders.length) % sliders.length); setProgress(0); }}
                className="w-12 h-12 rounded-full bg-white/10 hover:bg-white/20 border border-white/20 text-white flex items-center justify-center transition backdrop-blur-sm"
              >
                <i className="ph-bold ph-caret-left text-xl" />
              </button>
              <button 
                onClick={() => { setCurrentSlide((currentSlide + 1) % sliders.length); setProgress(0); }}
                className="w-12 h-12 rounded-full bg-white/10 hover:bg-white/20 border border-white/20 text-white flex items-center justify-center transition backdrop-blur-sm"
              >
                <i className="ph-bold ph-caret-right text-xl" />
              </button>
            </div>
            
            <div className="flex flex-col items-end gap-3">
              <div className="flex gap-2">
                {sliders.map((_, idx) => (
                  <button 
                    key={idx}
                    onClick={() => { setCurrentSlide(idx); setProgress(0); }}
                    className={`h-2.5 rounded-full transition-all duration-500 ${currentSlide === idx ? 'w-8 bg-amber-400' : 'w-2.5 bg-white/30 hover:bg-white/60'}`}
                  />
                ))}
              </div>
              <div className="w-32 h-1 bg-white/10 rounded-full overflow-hidden">
                <div className="h-full bg-brand-500 transition-all duration-100 ease-linear" style={{ width: `${progress}%` }} />
              </div>
            </div>
          </div>
        </div>
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
                  src={resolveAssetUrl(settings.greeting_image_url) || 'https://gradasi.org/uploads/img/event-terkini/1767154211.jpg'}
                  alt="Poster Sambutan"
                  className="w-full h-auto object-contain"
                />
              </div>
              <div className="absolute -bottom-6 -right-6 bg-white p-4 rounded-2xl shadow-xl border border-slate-100 flex items-center gap-4">
                <div className="w-12 h-12 bg-amber-100 rounded-full flex items-center justify-center text-amber-500 text-2xl">
                  <i className="ph-fill ph-star" />
                </div>
                <div>
                  <p className="font-bold text-slate-800 text-sm">{settings.greeting_title || 'Tahun Baru 2026'}</p>
                  <p className="text-[10px] text-slate-500 uppercase tracking-widest font-semibold">{settings.greeting_subtitle || 'Resolusi & Harapan'}</p>
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
                  {settings.greeting_title || 'Memperjuangkan'} <span className="text-transparent bg-clip-text bg-gradient-to-r from-brand-600 to-amber-500">{settings.greeting_subtitle || 'Kedaulatan Digital'}</span> Bangsa
                </h2>
              </div>

              {/* Elegant Quote Card */}
              <div className="relative bg-slate-50 rounded-2xl p-8 border border-slate-200/60 shadow-inner">
                <i className="ph-fill ph-quotes absolute -top-4 -left-2 text-6xl text-brand-200/50 transform -rotate-12" />
                <p className="relative z-10 text-slate-700 text-lg leading-relaxed font-quote italic">
                  "{settings.greeting_content || 'Tiada hadiah termahal di akhir tahun ini selain doa dan dukungan yang senantiasa mengiringi setiap langkah perjuangan kita. Tantangan teknologi bukan beban, melainkan proses pendewasaan agar kita tegak berdiri.'}"
                </p>
              </div>

              {/* Signature */}
              <div className="flex items-center gap-4 pt-4 border-t border-slate-100">
                <div className="flex -space-x-4">
                  <img src={resolveAssetUrl(settings.greeting_image_url) || "https://gradasi.org/uploads/img/s-anggota/ketua/1735027418.jpg"} alt="Upi" className="w-12 h-12 rounded-full border-2 border-white shadow-md relative z-20" />
                  <div className="w-12 h-12 rounded-full bg-brand-100 border-2 border-white shadow-md flex items-center justify-center text-brand-700 font-bold text-xs relative z-10">JS</div>
                </div>
                <div>
                  <p className="font-heading font-bold text-slate-900 text-sm">{settings.site_name || 'DPP GRADASI'}</p>
                  <p className="text-xs text-brand-600 font-medium">Upi Asmaradhana & Junaidi, S.Kom</p>
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
            <div className="w-full lg:w-1/3 flex justify-center">
              <div className="relative group">
                <div className="w-72 aspect-[4/5] rounded-3xl overflow-hidden shadow-2xl border-4 border-white mx-auto relative z-10 transform transition duration-500 group-hover:-translate-y-2">
                  <img 
                    src={resolveAssetUrl(settings.greeting_image_url) || "https://gradasi.org/uploads/img/s-anggota/ketua/1735027418.jpg"}
                    alt="Upi Asmaradhana - Ketua Umum"
                    className="w-full h-full object-cover"
                  />
                  <div className="absolute bottom-0 inset-x-0 bg-gradient-to-t from-slate-900 via-slate-900/60 to-transparent p-6 pt-12 text-center">
                    <h4 className="font-heading font-bold text-white text-lg">Upi Asmaradhana</h4>
                    <p className="text-[11px] text-brand-400 font-bold uppercase tracking-widest mt-1">Ketua Umum DPP</p>
                  </div>
                </div>
              </div>
            </div>

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
                      <p>{settings.history || 'Perkumpulan Generasi Digital Indonesia (GRADASI) didirikan sebagai wadah kolaborasi untuk meningkatkan kecakapan digital bangsa secara merata.'}</p>
                      <p>{settings.about_tutorial || 'Diinisiasi oleh para pegiat teknologi, GRADASI (Generasi Digital Indonesia) didirikan sebagai wadah kolaborasi untuk meningkatkan kecakapan digital bangsa secara merata.'}</p>
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
                        <h4 className="font-bold text-slate-900 text-lg mb-2">{settings.about_formation_date || '4 Februari 2019'}</h4>
                        <p className="text-slate-600 text-sm leading-relaxed mb-4">Secara resmi GRADASI disahkan melalui Surat Keputusan (SK) Kementerian Hukum dan HAM Republik Indonesia.</p>
                        <div className="inline-flex items-center gap-2 bg-white px-4 py-2 rounded-lg border border-brand-200 shadow-sm text-xs font-mono font-bold text-brand-800">
                          <i className="ph-bold ph-certificate" /> {settings.about_no_sk || 'AHU – 0000151.AH.01.07.2019'}
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
                      <h4 className="font-bold text-slate-900 text-base mb-3">Kantor Sekretariat DPP GRADASI</h4>
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

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 auto-rows-min">
            {/* VISI UTAMA Box */}
            <div className="md:col-span-2 lg:col-span-2 bg-gradient-to-br from-brand-700 to-brand-900 rounded-3xl p-8 sm:p-10 text-white relative overflow-hidden group shadow-xl">
              <i className="ph-fill ph-quotes absolute -bottom-10 -right-10 text-[180px] text-white/5 -rotate-12 pointer-events-none" />
              <div className="relative z-10 h-full flex flex-col justify-center">
                <div className="flex items-center gap-3 mb-6">
                  <div className="w-10 h-10 bg-white/20 backdrop-blur-sm rounded-full flex items-center justify-center border border-white/30">
                    <i className="ph-fill ph-eye text-xl" />
                  </div>
                  <span className="text-xs font-bold uppercase tracking-widest text-brand-200">Visi Utama</span>
                </div>
                <p className="text-2xl sm:text-3xl font-heading font-bold leading-tight drop-shadow-md italic">
                  "{settings.about_vision}"
                </p>
              </div>
            </div>

            {/* Logo Display Box */}
            <div className="bg-white rounded-3xl p-8 flex items-center justify-center border border-slate-100 shadow-md">
              <img src={resolveAssetUrl(settings.logo_url)} alt="Logo" className="w-48 h-auto object-contain filter drop-shadow-lg" />
            </div>

            {/* Dynamic Missions */}
            {getMissions().map((mission, idx) => (
              <div 
                key={idx} 
                className={`bg-white rounded-3xl p-8 border border-slate-100 shadow-md group transition-all duration-300 ${
                  idx % 3 === 0 ? 'hover:bg-blue-50' : idx % 3 === 1 ? 'hover:bg-teal-50' : 'hover:bg-amber-50'
                }`}
              >
                <div className={`w-14 h-14 rounded-2xl flex items-center justify-center text-2xl mb-6 shadow-sm ${
                  idx % 3 === 0 ? 'bg-blue-100 text-blue-600' : idx % 3 === 1 ? 'bg-teal-100 text-teal-600' : 'bg-amber-100 text-amber-600'
                }`}>
                  <i className={`ph-fill ${idx % 3 === 0 ? 'ph-users-three' : idx % 3 === 1 ? 'ph-chalkboard-teacher' : 'ph-handshake'}`} />
                </div>
                <h3 className="font-bold text-slate-900 mb-2 font-heading">Misi {idx + 1}</h3>
                <p className="text-slate-600 text-sm leading-relaxed">{mission}</p>
              </div>
            ))}
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
                      <button className="flex items-center gap-1.5 text-xs font-bold text-slate-400 hover:text-brand-600 transition">
                        <i className="ph-bold ph-share-network text-sm" /> Bagikan
                      </button>
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
                <source src={resolveAssetUrl(settings.video_profile_url) || 'https://gradasi.org/assets/video/gradasi.mp4'} type="video/mp4" />
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
                      <button className="flex items-center gap-1.5 text-xs font-bold text-slate-400 hover:text-brand-600 transition">
                        <i className="ph-bold ph-share-network text-sm" /> Bagikan
                      </button>
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
                  src={resolveAssetUrl(settings.maps_embed_url) || "https://maps.google.com/maps?q=The%20Bellagio%20Mall%20Mega%20Kuningan%20Jakarta&t=&z=16&ie=UTF8&iwloc=&output=embed"}
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
                        onChange={e => setContactForm({ ...contactForm, nama: e.target.value })}
                        className="w-full px-4 py-3 rounded-lg border border-slate-200 text-sm focus:border-brand-500 focus:outline-none bg-white transition"
                      />
                    </div>
                    <div>
                      <label className="block text-xs font-semibold text-slate-700 mb-1.5 uppercase tracking-wide">Email</label>
                      <input 
                        type="email" 
                        required
                        value={contactForm.email}
                        onChange={e => setContactForm({ ...contactForm, email: e.target.value })}
                        className="w-full px-4 py-3 rounded-lg border border-slate-200 text-sm focus:border-brand-500 focus:outline-none bg-white transition"
                      />
                    </div>
                  </div>
                  <div>
                    <label className="block text-xs font-semibold text-slate-700 mb-1.5 uppercase tracking-wide">Subjek / Perihal</label>
                    <input 
                      type="text" 
                      value={contactForm.subjek}
                      onChange={e => setContactForm({ ...contactForm, subjek: e.target.value })}
                      className="w-full px-4 py-3 rounded-lg border border-slate-200 text-sm focus:border-brand-500 focus:outline-none bg-white transition"
                    />
                  </div>
                  <div>
                    <label className="block text-xs font-semibold text-slate-700 mb-1.5 uppercase tracking-wide">Pesan Anda</label>
                    <textarea 
                      rows="4" 
                      required
                      value={contactForm.pesan}
                      onChange={e => setContactForm({ ...contactForm, pesan: e.target.value })}
                      className="w-full px-4 py-3 rounded-lg border border-slate-200 text-sm focus:border-brand-500 focus:outline-none bg-white transition resize-none"
                    />
                  </div>
                  <button 
                    type="submit" 
                    disabled={contactLoading}
                    className="bg-brand-600 hover:bg-brand-700 text-white px-8 py-3.5 rounded-xl text-sm font-bold transition shadow-sm hover:shadow flex items-center justify-center gap-2"
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

      {/* 9. FLOATING CONTACT WIDGET */}
      <div className="fixed bottom-6 right-6 z-50 flex flex-col items-end">
        {floatingOpen && (
          <div className="mb-4 space-y-3 flex flex-col items-end animate-fadeIn">
            <a 
              href={`mailto:${settings.contact_email}`}
              className="flex items-center gap-4 bg-white/90 backdrop-blur-md px-5 py-3 rounded-2xl shadow-lg border border-white/50 hover:bg-white transition"
            >
              <span className="text-sm font-bold text-slate-700">{settings.contact_email}</span>
              <div className="w-10 h-10 rounded-full bg-brand-50 flex items-center justify-center text-brand-600">
                <i className="ph-bold ph-envelope-simple text-xl" />
              </div>
            </a>
            <a 
              href={`tel:${settings.contact_phone}`}
              className="flex items-center gap-4 bg-white/90 backdrop-blur-md px-5 py-3 rounded-2xl shadow-lg border border-white/50 hover:bg-white transition"
            >
              <span className="text-sm font-bold text-slate-700">{settings.contact_phone}</span>
              <div className="w-10 h-10 rounded-full bg-brand-50 flex items-center justify-center text-brand-600">
                <i className="ph-bold ph-phone-call text-xl" />
              </div>
            </a>
          </div>
        )}
        <button 
          onClick={() => setFloatingOpen(!floatingOpen)}
          className={`w-16 h-16 rounded-full shadow-2xl flex items-center justify-center text-3xl text-white transition-all duration-300 ${
            floatingOpen ? 'bg-slate-800 rotate-[135deg]' : 'bg-brand-700 hover:bg-brand-600'
          }`}
        >
          <i className={`ph-bold ${floatingOpen ? 'ph-plus' : 'ph-headset'}`} />
        </button>
      </div>
    </PublicLayout>
  )
}
