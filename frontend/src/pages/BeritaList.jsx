import { useEffect, useState } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import PublicLayout from '../layouts/PublicLayout'
import { beritaContent } from '../content/beritaContent'
import { beritaService } from '../services/beritaService'
import { formatDate } from '../utils/format'

export default function BeritaList() {
  const [searchParams, setSearchParams] = useSearchParams()
  const page = Number(searchParams.get('page') || 1)
  const [items, setItems] = useState([])
  const [meta, setMeta] = useState({ current_page: 1, total_pages: 1, total_data: 0 })
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    async function load() {
      setLoading(true)
      setError('')
      try {
        const response = await beritaService.list({ page, limit: 9, sort: 'newest' })
        setItems(response.data.berita || [])
        setMeta(response.data.meta || { current_page: 1, total_pages: 1, total_data: 0 })
      } catch (err) {
        setError(err.message || 'Gagal mengambil daftar berita')
      } finally {
        setLoading(false)
      }
    }
    load()
  }, [page])

  function goToPage(nextPage) {
    setSearchParams({ page: String(nextPage) })
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }

  return (
    <PublicLayout>
      <section className="pt-36 pb-16 md:pt-40 md:pb-20 bg-white border-b border-slate-100">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <div className="inline-flex items-center justify-center gap-2 text-xs font-bold text-brand-600 bg-brand-50 px-4 py-1.5 rounded-full mb-4">
            <i className="ph-bold ph-newspaper" /> {beritaContent.publicList.badge}
          </div>
          <h1 className="font-heading text-3xl md:text-4xl font-extrabold text-slate-900 mb-4 tracking-tight">{beritaContent.publicList.title}</h1>
          <p className="text-base text-slate-500 max-w-2xl mx-auto font-medium">{beritaContent.publicList.subtitle}</p>
        </div>
      </section>

      <section className="py-12 md:py-16 bg-slate-50">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
          {loading && <div className="py-16 text-center text-slate-500">Memuat berita...</div>}
          {error && <div className="py-16 text-center text-red-600 font-medium">{error}</div>}
          {!loading && !error && items.length === 0 && <div className="py-16 text-center text-slate-500">{beritaContent.publicList.empty}</div>}

          {!loading && !error && items.length > 0 && (
            <>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {items.map((item) => (
                  <article key={item.id} className="group cursor-pointer bg-white rounded-2xl overflow-hidden shadow-sm hover:shadow-[0_8px_30px_rgba(37,99,235,0.08)] transition-all duration-300 border border-slate-100 flex flex-col">
                    <div className="h-44 relative overflow-hidden bg-slate-100">
                      {item.image_url && <img src={item.image_url} alt={item.title} className="w-full h-full object-cover transform group-hover:scale-105 transition duration-500" />}
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

              <div className="mt-16 flex justify-center">
                <nav className="flex items-center gap-2">
                  <button type="button" disabled={page <= 1} onClick={() => goToPage(page - 1)} className="w-10 h-10 flex items-center justify-center rounded-xl border border-slate-200 text-slate-400 hover:border-brand-500 hover:text-brand-600 hover:bg-brand-50 transition disabled:opacity-40">
                    <i className="ph-bold ph-caret-left" />
                  </button>
                  {Array.from({ length: meta.total_pages || 1 }, (_, index) => index + 1).map((number) => (
                    <button key={number} type="button" onClick={() => goToPage(number)} className={`w-10 h-10 flex items-center justify-center rounded-xl font-bold transition ${number === page ? 'bg-brand-700 text-white shadow-md shadow-brand-500/20' : 'border border-slate-200 text-slate-600 hover:border-brand-500 hover:text-brand-600 hover:bg-brand-50'}`}>
                      {number}
                    </button>
                  ))}
                  <button type="button" disabled={page >= (meta.total_pages || 1)} onClick={() => goToPage(page + 1)} className="w-10 h-10 flex items-center justify-center rounded-xl border border-slate-200 text-slate-400 hover:border-brand-500 hover:text-brand-600 hover:bg-brand-50 transition disabled:opacity-40">
                    <i className="ph-bold ph-caret-right" />
                  </button>
                </nav>
              </div>
            </>
          )}
        </div>
      </section>
    </PublicLayout>
  )
}
