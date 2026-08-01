import { useState, useEffect } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import PublicLayout from '../layouts/PublicLayout'
import { kegiatanService } from '../services/kegiatanService'

export default function KegiatanList() {
  const [searchParams] = useSearchParams()
  const categoryParam = searchParams.get('category') || ''

  const [items, setItems] = useState([
    {
      id: 1,
      title: 'Penyaluran Bantuan Kemanusiaan oleh DPP GRADASI',
      slug: 'penyaluran-bantuan-kemanusiaan',
      category: 'Nasional',
      event_date: '31 Desember 2025',
      location: 'Jakarta',
      image_url: 'https://gradasi.org/uploads/img/event/1767154719.jpg',
      excerpt: 'Dewan Pimpinan Pusat (DPP) GRADASI turun langsung menyalurkan bantuan kemanusiaan kepada masyarakat yang terdampak bencana alam sebagai wujud kepedulian sosial.',
      is_published: true
    },
    {
      id: 2,
      title: 'Pelatihan Digital Marketing UMKM Go Online',
      slug: 'pelatihan-digital-marketing-umkm',
      category: 'Jawa Timur',
      event_date: '15 November 2025',
      location: 'Surabaya',
      image_url: 'https://gradasi.org/uploads/img/event/1767154619.jpg',
      excerpt: 'Program pelatihan intensif bagi pelaku Usaha Mikro Kecil Menengah (UMKM) untuk memasarkan produknya secara digital demi menjangkau pasar yang lebih luas.',
      is_published: true
    },
    {
      id: 3,
      title: 'Konsolidasi Pengurus DPP & Penyerahan SK Daerah',
      slug: 'konsolidasi-pengurus-dpp-dpd',
      category: 'Lampung',
      event_date: '02 Oktober 2025',
      location: 'Bandar Lampung',
      image_url: 'https://gradasi.org/uploads/img/event/1767154397.jpg',
      excerpt: 'Acara konsolidasi pengurus tingkat pusat serta penyerahan Surat Keputusan (SK) kepada perwakilan pengurus daerah demi memperkuat struktur organisasi di seluruh nusantara.',
      is_published: true
    }
  ])
  const [loading, setLoading] = useState(true)

  const [searchQuery, setSearchQuery] = useState('')
  const [categoryFilter, setCategoryFilter] = useState(categoryParam)
  const [sortBy, setSortBy] = useState('newest')
  const [currentPage, setCurrentPage] = useState(1)
  const itemsPerPage = 6

  useEffect(() => {
    if (categoryParam) setCategoryFilter(categoryParam)
  }, [categoryParam])

  useEffect(() => {
    kegiatanService.list()
      .then(res => {
        if (res.success && res.data) {
          const list = Array.isArray(res.data) ? res.data : (res.data.kegiatan || [])
          setItems(list)
        }
      })
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [])

  useEffect(() => {
    setCurrentPage(1)
  }, [searchQuery, categoryFilter, sortBy])

  // Filter & Sort Logic
  const filteredItems = items.filter(item => {
    const matchesSearch = item.title.toLowerCase().includes(searchQuery.toLowerCase()) || 
                          (item.location || '').toLowerCase().includes(searchQuery.toLowerCase())
    const matchesCategory = !categoryFilter || item.category === categoryFilter
    return matchesSearch && matchesCategory && (item.is_published !== false)
  }).sort((a, b) => {
    if (sortBy === 'newest') return new Date(b.created_at || '2026-01-01') - new Date(a.created_at || '2026-01-01')
    if (sortBy === 'oldest') return new Date(a.created_at || '2026-01-01') - new Date(b.created_at || '2026-01-01')
    return 0
  })

  const totalPages = Math.ceil(filteredItems.length / itemsPerPage) || 1
  const paginatedItems = filteredItems.slice((currentPage - 1) * itemsPerPage, currentPage * itemsPerPage)
  // Panah < > hanya tampil saat minimal 3 data & lebih dari 1 halaman
  const showArrows = totalPages > 1 && filteredItems.length >= 3

  const categories = [...new Set(items.map(item => item.category))].filter(Boolean)

  return (
    <PublicLayout>
      {/* HEADER BANNER */}
      <section className="pt-36 pb-16 md:pt-40 md:pb-20 bg-white border-b border-slate-100">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <div className="inline-flex items-center justify-center gap-2 text-xs font-bold text-brand-600 bg-brand-50 px-4 py-1.5 rounded-full mb-4">
            <i className="ph-bold ph-calendar-star" /> Kegiatan GRADASI
          </div>
          <h1 className="font-heading text-3xl md:text-4xl font-extrabold text-slate-900 mb-4 tracking-tight">
            Dokumentasi & Aktivitas
          </h1>
          <p className="text-base text-slate-500 max-w-2xl mx-auto font-medium">
            Ikuti terus berbagai kegiatan inspiratif, program kerja, dan langkah nyata pemuda dalam membangun literasi digital nasional.
          </p>
        </div>
      </section>

      {/* MAIN CONTENT GRID */}
      <section className="py-12 md:py-16 bg-slate-50">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
          
          {/* Filter Bar */}
          <div className="mb-10 flex flex-col md:flex-row gap-4 justify-between items-center bg-white p-4 rounded-2xl border border-slate-100 shadow-sm w-full">
            <div className="relative w-full md:max-w-md">
              <i className="ph-bold ph-magnifying-glass absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400 text-lg" />
              <input 
                type="text" 
                placeholder="Cari kegiatan atau lokasi..." 
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-full pl-10 pr-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl text-sm focus:bg-white focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 transition-all outline-none"
              />
            </div>
            
            <div className="flex flex-col sm:flex-row gap-4 w-full md:w-auto">
              <select 
                value={categoryFilter} 
                onChange={(e) => setCategoryFilter(e.target.value)}
                className="px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl text-sm text-slate-600 outline-none cursor-pointer"
              >
                <option value="">Semua Kategori</option>
                {categories.map(cat => (
                  <option key={cat} value={cat}>{cat}</option>
                ))}
              </select>

              <select 
                value={sortBy} 
                onChange={(e) => setSortBy(e.target.value)}
                className="px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl text-sm text-slate-600 outline-none cursor-pointer"
              >
                <option value="newest">Urutkan: Terbaru</option>
                <option value="oldest">Urutkan: Terlama</option>
              </select>
            </div>
          </div>

          {/* Items Grid */}
          {filteredItems.length === 0 ? (
            <div className="py-16 text-center text-slate-500 bg-white rounded-2xl border border-slate-100 p-8">
              <i className="ph-bold ph-calendar-x text-4xl text-slate-300 mb-2 block" />
              <p className="font-semibold text-slate-700">Tidak ada kegiatan yang ditemukan</p>
              <p className="text-xs text-slate-400 mt-1">Coba gunakan kata kunci pencarian atau kategori lain.</p>
            </div>
          ) : (
            <>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {paginatedItems.map(item => (
                  <div 
                    key={item.id} 
                    className="group cursor-pointer bg-white rounded-2xl overflow-hidden shadow-sm hover:shadow-[0_8px_30px_rgba(37,99,235,0.08)] transition-all duration-300 border border-slate-100 flex flex-col h-full"
                  >
                    <div className="h-44 relative overflow-hidden bg-slate-100">
                      <img 
                        src={item.image_url ? (item.image_url.startsWith('http') ? item.image_url : `http://127.0.0.1:8080${item.image_url}`) : 'https://images.unsplash.com/photo-1540575467063-178a50c2df87?q=80&w=600'} 
                        alt={item.title} 
                        className="w-full h-full object-cover transform group-hover:scale-105 transition duration-500"
                      />
                    </div>
                    <div className="p-5 flex flex-col flex-grow">
                      <div className="flex items-center gap-3 mb-2 flex-wrap">
                        <p className="text-brand-600 text-[10px] font-bold tracking-wider uppercase flex items-center gap-1.5">
                          <i className="ph-bold ph-calendar-blank text-sm" /> {item.event_date || '31 Desember 2025'}
                        </p>
                        {item.location && (
                          <p className="text-slate-500 text-[10px] font-bold tracking-wider uppercase flex items-center gap-1.5 border-l border-slate-200 pl-3">
                            <i className="ph-bold ph-map-pin text-sm" /> {item.location}
                          </p>
                        )}
                      </div>
                      <Link to={`/kegiatan/${item.slug}`} className="font-heading text-lg font-bold text-slate-900 mb-2 group-hover:text-brand-600 transition leading-snug">
                        {item.title}
                      </Link>
                      <p className="text-slate-500 text-[13px] flex-grow line-clamp-2 mb-4 leading-relaxed">{item.excerpt}</p>
                      <div className="pt-4 flex justify-between items-center border-t border-slate-100 mt-auto">
                        <button className="flex items-center gap-1.5 text-xs font-bold text-slate-400 hover:text-brand-600 transition group/btn">
                          <div className="w-7 h-7 rounded-full bg-slate-50 flex items-center justify-center group-hover/btn:bg-brand-50 transition-colors">
                            <i className="ph-bold ph-share-network text-sm" />
                          </div>
                          <span className="hidden sm:inline">Bagikan</span>
                        </button>
                        <Link to={`/kegiatan/${item.slug}`} className="flex items-center gap-1.5 text-xs font-bold text-brand-600 hover:text-brand-800 transition">
                          Baca Selengkapnya <i className="ph-bold ph-arrow-right" />
                        </Link>
                      </div>
                    </div>
                  </div>
                ))}
              </div>

              {/* PAGINATION COMPONENT */}
              <div className="mt-16 flex justify-center">
                <nav className="flex items-center gap-2">
                  {showArrows && (
                  <button 
                    disabled={currentPage <= 1}
                    onClick={() => setCurrentPage(prev => Math.max(prev - 1, 1))}
                    className="w-10 h-10 flex items-center justify-center rounded-xl border border-slate-200 text-slate-400 hover:border-brand-500 hover:text-brand-600 hover:bg-brand-50 transition font-bold disabled:opacity-40"
                  >
                    <i className="ph-bold ph-caret-left" />
                  </button>
                  )}
                  {Array.from({ length: totalPages }, (_, i) => i + 1).map(page => (
                    <button 
                      key={page} 
                      onClick={() => setCurrentPage(page)} 
                      className={`w-10 h-10 flex items-center justify-center rounded-xl font-bold transition ${currentPage === page ? 'bg-brand-700 text-white shadow-md shadow-brand-500/20' : 'border border-slate-200 text-slate-600 hover:border-brand-500 hover:text-brand-600 hover:bg-brand-50'}`}
                    >
                      {page}
                    </button>
                  ))}
                  {showArrows && (
                  <button 
                    disabled={currentPage >= totalPages}
                    onClick={() => setCurrentPage(prev => Math.min(prev + 1, totalPages))}
                    className="w-10 h-10 flex items-center justify-center rounded-xl border border-slate-200 text-slate-400 hover:border-brand-500 hover:text-brand-600 hover:bg-brand-50 transition font-bold disabled:opacity-40"
                  >
                    <i className="ph-bold ph-caret-right" />
                  </button>
                  )}
                </nav>
              </div>
            </>
          )}

        </div>
      </section>
    </PublicLayout>
  )
}
