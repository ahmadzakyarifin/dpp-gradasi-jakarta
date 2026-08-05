import { useEffect, useState } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import PublicLayout from '../layouts/PublicLayout'
import { beritaContent } from '../content/beritaContent'
import { beritaService } from '../services/beritaService'
import { formatDate } from '../utils/format'
import { resolveAssetUrl } from '../utils/assetUrl'
import { shareContent, getShareUrl, copyToClipboard } from '../utils/share'
import ToastNotification from '../components/admin/ToastNotification'

export default function BeritaList() {
  const [searchParams, setSearchParams] = useSearchParams()
  const page = Number(searchParams.get('page') || 1)
  const searchQuery = searchParams.get('q') || ''
  const filterCategory = searchParams.get('category') || ''
  const sortBy = searchParams.get('sort') || 'newest'
  
  const [items, setItems] = useState([])
  const [meta, setMeta] = useState({ current_page: 1, total_pages: 1, total_data: 0 })
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [categories, setCategories] = useState([])
  const [toast, setToast] = useState({ show: false, message: '', type: 'success' })
  const currentPage = page

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

  useEffect(() => {
    beritaService.getCategories()
      .then(res => {
        if (res && res.data && Array.isArray(res.data)) {
          setCategories(res.data)
        }
      })
      .catch(() => {})
  }, [])

  useEffect(() => {
    setLoading(true)
    setError(null)
    const params = { page: currentPage, limit: 6, sort: sortBy }
    if (searchQuery.trim()) params.search = searchQuery.trim()
    if (filterCategory) params.category = filterCategory
    beritaService.list(params)
      .then(res => {
        if (res.data && res.data.berita) {
          setItems(res.data.berita)
          setMeta(res.data.meta || { current_page: currentPage, total_pages: 1, total_data: res.data.berita.length })
        }
      })
      .catch(err => {
        setError(err?.message || 'Gagal memuat berita')
        setItems([])
      })
      .finally(() => setLoading(false))
  }, [currentPage, searchQuery, filterCategory, sortBy])

  function goToPage(nextPage) {
    const params = new URLSearchParams(searchParams)
    params.set('page', String(nextPage))
    setSearchParams(params)
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }

  function updateFilter(next) {
    const params = new URLSearchParams(searchParams)
    Object.entries(next).forEach(([k, v]) => {
      if (v) params.set(k, v)
      else params.delete(k)
    })
    params.set('page', '1')
    setSearchParams(params)
  }

  const totalPages = meta.total_pages || 1
  const pageNumbers = Array.from({ length: totalPages }, (_, i) => i + 1)
  // Panah < > hanya tampil saat ada minimal 3 data & lebih dari 1 halaman; seluruh pagination hilang saat kosong
  const showArrows = totalPages > 1 && meta.total_data >= 3

  return (
    <PublicLayout>
      <section className="pt-36 pb-16 md:pt-40 md:pb-20 bg-white border-b border-slate-100">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <div className="inline-flex items-center justify-center gap-2 text-xs font-bold text-brand-600 bg-brand-50 px-4 py-1.5 rounded-full mb-4">
            <i className="ph-bold ph-newspaper" /> {beritaContent.publicList.badge}
          </div>
          <h1 className="font-heading text-3xl md:text-4xl font-extrabold text-slate-900 mb-4 tracking-tight">{beritaContent.publicList.title}</h1>
          <p className="text-base text-slate-500 max-w-2xl mx-auto font-medium">{beritaContent.publicList.subtitle}</p>
          {(searchQuery || filterCategory) && (
            <p className="mt-4 inline-flex items-center gap-2 text-xs font-bold text-brand-600 bg-brand-50 px-4 py-1.5 rounded-full">
              <i className="ph-bold ph-funnel" />
              {searchQuery && <>Hasil untuk "{searchQuery}"</>}
              {searchQuery && filterCategory && ' · '}
              {filterCategory && <>Kategori: {filterCategory}</>}
            </p>
          )}
        </div>
      </section>

      <section className="py-12 md:py-16 bg-slate-50">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
          
          {/* Filter Bar */}
          <div className="mb-10 flex flex-col md:flex-row gap-4 justify-between items-center bg-white p-4 rounded-2xl border border-slate-100 shadow-sm w-full">
            <div className="relative w-full md:max-w-md">
              <i className="ph-bold ph-magnifying-glass absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400 text-lg" />
              <input 
                type="text" 
                placeholder="Cari berita atau informasi..." 
                value={searchQuery}
                onChange={(e) => updateFilter({ q: e.target.value })}
                className="w-full pl-10 pr-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl text-sm focus:bg-white focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 transition-all outline-none"
              />
            </div>
            
            <div className="flex flex-col sm:flex-row gap-4 w-full md:w-auto">
              <select 
                value={filterCategory} 
                onChange={(e) => updateFilter({ category: e.target.value })}
                className="px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl text-sm text-slate-600 outline-none cursor-pointer"
              >
                <option value="">Semua Kategori</option>
                {categories.map(cat => (
                  <option key={cat} value={cat}>{cat}</option>
                ))}
              </select>

              <select 
                value={sortBy} 
                onChange={(e) => updateFilter({ sort: e.target.value })}
                className="px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl text-sm text-slate-600 outline-none cursor-pointer"
              >
                <option value="newest">Urutkan: Terbaru</option>
                <option value="oldest">Urutkan: Terlama</option>
              </select>
            </div>
          </div>

          {loading && <div className="py-20 text-center text-slate-500">Memuat berita...</div>}
          {!loading && error && (
            <div className="py-20 text-center text-red-600 font-medium">
              <i className="ph-bold ph-warning-circle text-2xl mb-2 block mx-auto" /> {error}
            </div>
          )}
          {!loading && !error && items.length === 0 && (
            <div className="py-16 text-center text-slate-500 bg-white rounded-2xl border border-slate-100 p-8">
              <i className="ph-bold ph-newspaper-clipping text-4xl text-slate-300 mb-2 block" />
              <p className="font-semibold text-slate-700">Tidak ada berita yang ditemukan</p>
              <p className="text-xs text-slate-400 mt-1">Belum ada berita atau informasi yang dipublikasikan.</p>
            </div>
          )}
          {!loading && !error && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {items.map((item) => (
              <article key={item.id} className="group cursor-pointer card-lift bg-white rounded-2xl overflow-hidden shadow-sm hover:shadow-[0_8px_30px_rgba(37,99,235,0.08)] transition-all duration-300 border border-slate-100 flex flex-col">
                <div className="h-44 relative overflow-hidden bg-slate-100">
                  <img 
                    src={item.image_url ? resolveAssetUrl(item.image_url) : 'https://images.unsplash.com/photo-1504711434969-e33886168f5c?q=80&w=600'} 
                    alt={item.title} 
                    className="w-full h-full object-cover transform group-hover:scale-105 transition duration-500" 
                  />
                </div>
                <div className="p-5 flex flex-col flex-grow">
                  <p className="text-brand-600 text-[10px] font-bold tracking-wider uppercase mb-2 flex items-center gap-1.5">
                    <i className="ph-bold ph-calendar-blank text-sm" /> {formatDate(item.published_date)}
                  </p>
                  <Link to={`/berita/${item.slug}`} className="font-heading text-lg font-bold text-slate-900 mb-2 group-hover:text-brand-600 transition line-clamp-2">
                    {item.title}
                  </Link>
                  <p className="text-slate-500 text-[13px] flex-grow line-clamp-2 mb-4 leading-relaxed">{item.excerpt || '-'}</p>
                  <div className="pt-4 flex justify-between items-center border-t border-slate-100 mt-auto">
                      <div className="flex items-center gap-1.5">
                        <button 
                          onClick={(e) => {
                            e.preventDefault();
                            e.stopPropagation();
                            handleInstagramShare(`${window.location.origin}/berita/${item.slug}`);
                          }} 
                          className="w-7 h-7 rounded-full bg-[#E1306C]/10 text-[#E1306C] hover:bg-[#E1306C] hover:text-white flex items-center justify-center transition" 
                          title="Instagram"
                        >
                          <i className="ph-fill ph-instagram-logo text-xs" />
                        </button>
                        <a 
                          href={getShareUrl('whatsapp', { title: item.title, text: item.excerpt, url: `${window.location.origin}/berita/${item.slug}` })} 
                          target="_blank" 
                          rel="noopener noreferrer" 
                          onClick={(e) => e.stopPropagation()}
                          className="w-7 h-7 rounded-full bg-[#25D366]/10 text-[#25D366] hover:bg-[#25D366] hover:text-white flex items-center justify-center transition" 
                          title="WhatsApp"
                        >
                          <i className="ph-fill ph-whatsapp-logo text-xs" />
                        </a>
                        <a 
                          href={getShareUrl('facebook', { title: item.title, text: item.excerpt, url: `${window.location.origin}/berita/${item.slug}` })} 
                          target="_blank" 
                          rel="noopener noreferrer" 
                          onClick={(e) => e.stopPropagation()}
                          className="w-7 h-7 rounded-full bg-[#1877F2]/10 text-[#1877F2] hover:bg-[#1877F2] hover:text-white flex items-center justify-center transition" 
                          title="Facebook"
                        >
                          <i className="ph-fill ph-facebook-logo text-xs" />
                        </a>
                        <button 
                          onClick={(e) => {
                            e.preventDefault();
                            e.stopPropagation();
                            handleCopyLink(`${window.location.origin}/berita/${item.slug}`, 'Tautan berita');
                          }} 
                          className="w-7 h-7 rounded-full bg-slate-100 text-slate-600 hover:bg-slate-600 hover:text-white flex items-center justify-center transition" 
                          title="Salin Tautan"
                        >
                          <i className="ph-bold ph-link text-xs" />
                        </button>
                      </div>
                    <Link to={`/berita/${item.slug}`} className="flex items-center gap-1.5 text-xs font-bold text-brand-600 hover:text-brand-800 transition">
                      {beritaContent.publicList.readMore} <i className="ph-bold ph-arrow-right" />
                    </Link>
                  </div>
                </div>
              </article>
            ))}
          </div>
          )}

          {/* PAGINATION — hilang total saat data kosong */}
          {items.length > 0 && (
          <div className="mt-16 flex justify-center">
            <nav className="flex items-center gap-2">
              {showArrows && (
              <button 
                type="button" 
                disabled={currentPage <= 1} 
                onClick={() => goToPage(currentPage - 1)} 
                className="w-10 h-10 flex items-center justify-center rounded-xl border border-slate-200 text-slate-400 hover:border-brand-500 hover:text-brand-600 hover:bg-brand-50 transition disabled:opacity-40"
              >
                <i className="ph-bold ph-caret-left" />
              </button>
              )}
              
              {pageNumbers.map((number) => (
                <button 
                  key={number} 
                  type="button" 
                  onClick={() => goToPage(number)} 
                  className={`w-10 h-10 flex items-center justify-center rounded-xl font-bold transition ${number === currentPage ? 'bg-brand-700 text-white shadow-md shadow-brand-500/20' : 'border border-slate-200 text-slate-600 hover:border-brand-500 hover:text-brand-600 hover:bg-brand-50'}`}
                >
                  {number}
                </button>
              ))}

              {showArrows && (
              <button 
                type="button" 
                disabled={currentPage >= totalPages} 
                onClick={() => goToPage(currentPage + 1)} 
                className="w-10 h-10 flex items-center justify-center rounded-xl border border-slate-200 text-slate-400 hover:border-brand-500 hover:text-brand-600 hover:bg-brand-50 transition disabled:opacity-40"
              >
                <i className="ph-bold ph-caret-right" />
              </button>
              )}
            </nav>
          </div>
          )}
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
