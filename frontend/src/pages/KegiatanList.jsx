import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import PublicLayout from '../layouts/PublicLayout'
import { kegiatanService } from '../services/kegiatanService'

export default function KegiatanList() {
  const [items, setItems] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  
  const [searchQuery, setSearchQuery] = useState('')
  const [categoryFilter, setCategoryFilter] = useState('')
  const [sortBy, setSortBy] = useState('newest')

  useEffect(() => {
    setLoading(true)
    kegiatanService.list()
      .then(res => {
        if (res.success && res.data) {
          setItems(res.data)
        } else {
          setError('Gagal memuat kegiatan')
        }
      })
      .catch(() => setError('Terjadi kesalahan koneksi'))
      .finally(() => setLoading(false))
  }, [])

  // Filter & Sort Logic
  const filteredItems = items.filter(item => {
    const matchesSearch = item.title.toLowerCase().includes(searchQuery.toLowerCase()) || 
                          item.location.toLowerCase().includes(searchQuery.toLowerCase())
    const matchesCategory = !categoryFilter || item.category === categoryFilter
    return matchesSearch && matchesCategory && item.is_published
  }).sort((a, b) => {
    if (sortBy === 'newest') return new Date(b.created_at) - new Date(a.created_at)
    if (sortBy === 'oldest') return new Date(a.created_at) - new Date(b.created_at)
    return 0
  })

  // Get unique categories for filter dropdown
  const categories = [...new Set(items.map(item => item.category))].filter(Boolean)

  return (
    <PublicLayout>
      <section className="pt-36 pb-16 md:pt-40 md:pb-20 bg-white border-b border-slate-100">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <span className="inline-flex items-center justify-center gap-2 text-xs font-bold text-brand-600 bg-brand-50 px-4 py-1.5 rounded-full mb-4">
             Agenda & Event
          </span>
          <h1 className="font-heading text-3xl md:text-4xl font-extrabold text-slate-900 mb-4 tracking-tight">
            Kegiatan DPP GRADASI
          </h1>
          <p className="text-base text-slate-500 max-w-2xl mx-auto font-medium">
            Ikuti berbagai program kerja, webinar, pelatihan, dan musyawarah nasional yang kami selenggarakan.
          </p>
        </div>
      </section>

      <section className="py-12 md:py-16 bg-slate-50">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
          
          {/* Filter Bar */}
          <div className="mb-10 flex flex-col md:flex-row gap-4 justify-between items-center bg-white p-4 rounded-2xl border border-slate-100 shadow-sm w-full">
            <div className="relative w-full md:max-w-md">
              <input 
                type="text" 
                placeholder="Cari kegiatan atau lokasi..." 
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-full pl-4 pr-10 py-2.5 bg-slate-50 border border-slate-200 rounded-xl text-sm focus:bg-white focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 transition-all outline-none"
              />
            </div>
            
            <div className="flex flex-col sm:flex-row gap-4 w-full md:w-auto">
              <select 
                value={categoryFilter} 
                onChange={(e) => setCategoryFilter(e.target.value)}
                className="px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl text-sm text-slate-600 outline-none"
              >
                <option value="">Semua Kategori</option>
                {categories.map(cat => (
                  <option key={cat} value={cat}>{cat}</option>
                ))}
              </select>

              <select 
                value={sortBy} 
                onChange={(e) => setSortBy(e.target.value)}
                className="px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl text-sm text-slate-600 outline-none"
              >
                <option value="newest">Terbaru</option>
                <option value="oldest">Terlama</option>
              </select>
            </div>
          </div>

          {loading && <div className="py-16 text-center text-slate-500">Memuat kegiatan...</div>}
          {error && <div className="py-16 text-center text-red-600 font-medium">{error}</div>}
          {!loading && !error && filteredItems.length === 0 && (
            <div className="py-16 text-center text-slate-500">Tidak ada kegiatan yang ditemukan.</div>
          )}

          {!loading && !error && filteredItems.length > 0 && (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {filteredItems.map(item => (
                <Link to={`/kegiatan/${item.slug}`} key={item.id} className="group bg-white rounded-2xl overflow-hidden shadow-sm hover:shadow-[0_8px_30px_rgba(37,99,235,0.08)] transition-all duration-300 border border-slate-100 flex flex-col h-full">
                  <div className="h-48 relative overflow-hidden bg-slate-100">
                    <img 
                      src={item.image_url ? (item.image_url.startsWith('http') ? item.image_url : `http://127.0.0.1:8080${item.image_url}`) : 'https://images.unsplash.com/photo-1540575467063-178a50c2df87?q=80&w=600'} 
                      alt={item.title} 
                      className="w-full h-full object-cover transform group-hover:scale-105 transition duration-500"
                    />
                  </div>
                  <div className="p-5 flex flex-col flex-grow">
                    <span className="text-brand-600 text-[10px] font-bold tracking-wider uppercase mb-2 block">{item.category}</span>
                    <h3 className="font-heading font-bold text-slate-900 text-lg group-hover:text-brand-600 transition mb-3 line-clamp-2">{item.title}</h3>
                    <p className="text-slate-500 text-sm mb-4 line-clamp-3 flex-grow">{item.excerpt}</p>
                    
                    <div className="border-t border-slate-100 pt-4 mt-auto flex items-center justify-between text-xs text-slate-400 font-medium">
                      <span className="flex items-center gap-1"><span className="text-[14px]">📅</span> {item.event_date}</span>
                      <span className="flex items-center gap-1"><span className="text-[14px]">📍</span> {item.location}</span>
                    </div>
                  </div>
                </Link>
              ))}
            </div>
          )}

        </div>
      </section>
    </PublicLayout>
  )
}
