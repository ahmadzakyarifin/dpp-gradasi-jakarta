import { useEffect, useState } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import PublicLayout from '../layouts/PublicLayout'
import { beritaContent } from '../content/beritaContent'
import { beritaService } from '../services/beritaService'
import { formatDate } from '../utils/format'

export default function BeritaList() {
  const [searchParams, setSearchParams] = useSearchParams()
  const page = Number(searchParams.get('page') || 1)
  const searchQuery = searchParams.get('q') || ''
  const filterCategory = searchParams.get('category') || ''
  const [items, setItems] = useState([])
  const [meta, setMeta] = useState({ current_page: 1, total_pages: 1, total_data: 0 })
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const currentPage = page

  useEffect(() => {
    setLoading(true)
    setError(null)
    const params = { page: currentPage, limit: 6, sort: 'newest' }
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
  }, [currentPage, searchQuery, filterCategory])

  function goToPage(nextPage) {
    const params = new URLSearchParams(searchParams)
    params.set('page', String(nextPage))
    setSearchParams(params)
    window.scrollTo({ top: 0, behavior: 'smooth' })
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
          {loading && <div className="py-20 text-center text-slate-500">Memuat berita...</div>}
          {!loading && error && (
            <div className="py-20 text-center text-red-600 font-medium">
              <i className="ph-bold ph-warning-circle text-2xl mb-2 block mx-auto" /> {error}
            </div>
          )}
          {!loading && !error && items.length === 0 && (
            <div className="py-20 text-center text-slate-500">Belum ada berita.</div>
          )}
          {!loading && !error && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {items.map((item) => (
              <article key={item.id} className="group cursor-pointer card-lift bg-white rounded-2xl overflow-hidden shadow-sm hover:shadow-[0_8px_30px_rgba(37,99,235,0.08)] transition-all duration-300 border border-slate-100 flex flex-col">
                <div className="h-44 relative overflow-hidden bg-slate-100">
                  <img 
                    src={item.image_url ? (item.image_url.startsWith('http') ? item.image_url : `http://127.0.0.1:8080${item.image_url}`) : 'https://images.unsplash.com/photo-1504711434969-e33886168f5c?q=80&w=600'} 
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
                    <button type="button" className="flex items-center gap-1.5 text-xs font-bold text-slate-400 hover:text-brand-600 transition group/btn">
                      <div className="w-7 h-7 rounded-full bg-slate-50 flex items-center justify-center group-hover/btn:bg-brand-50 transition-colors">
                        <i className="ph-bold ph-share-network text-sm" />
                      </div>
                      <span className="hidden sm:inline">{beritaContent.publicList.share}</span>
                    </button>
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
    </PublicLayout>
  )
}
