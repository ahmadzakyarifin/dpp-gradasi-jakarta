import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import PublicLayout from '../layouts/PublicLayout'
import { useSettings } from '../context/useSettings'
import { resolveAssetUrl } from '../utils/assetUrl'
import { beritaContent } from '../content/beritaContent'
import { beritaService } from '../services/beritaService'
import { formatDate } from '../utils/format'

export default function BeritaDetail() {
  const { settings } = useSettings()
  const { slug } = useParams()
  const [article, setArticle] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    async function load() {
      setLoading(true)
      setError('')
      try {
        const response = await beritaService.detailBySlug(slug)
        setArticle(response.data)
      } catch (err) {
        setError(err.message || beritaContent.detail.notFound)
      } finally {
        setLoading(false)
      }
    }
    load()
  }, [slug])

  function shareTo(platform) {
    const url = window.location.href
    const text = article?.title || beritaContent.detail.defaultTitle
    if (platform === 'instagram') {
      navigator.clipboard.writeText(url).then(() => {
        alert('Link berita berhasil disalin ke clipboard! Silakan bagikan di Instagram.')
      }).catch(() => {
        alert('Gagal menyalin link. Silakan salin URL di address bar browser Anda.')
      })
      return
    }
    const links = {
      facebook: `https://www.facebook.com/sharer/sharer.php?u=${encodeURIComponent(url)}`,
      whatsapp: `https://api.whatsapp.com/send?text=${encodeURIComponent(`${text} ${url}`)}`,
    }
    window.open(links[platform], '_blank', 'noopener,noreferrer')
  }

  return (
    <PublicLayout>
      <div className="fixed top-0 left-0 h-[3px] bg-gradient-to-r from-brand-600 to-brand-400 z-[9999] w-full" />
      <section className="pt-28 pb-0 bg-white border-b border-slate-100">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <nav className="flex items-center gap-2 text-sm py-4">
            <Link to="/" className="text-slate-400 hover:text-brand-600 transition flex items-center gap-1">
              <i className="ph-bold ph-house text-xs" /> {beritaContent.detail.breadcrumbHome}
            </Link>
            <i className="ph-bold ph-caret-right text-[10px] text-slate-300" />
            <Link to="/berita" className="text-slate-400 hover:text-brand-600 transition">{beritaContent.detail.breadcrumbInfo}</Link>
            <i className="ph-bold ph-caret-right text-[10px] text-slate-300" />
            <span className="text-brand-600 font-semibold truncate max-w-[250px]">{article?.title || beritaContent.detail.defaultTitle}</span>
          </nav>
        </div>
      </section>

      {loading && <section className="bg-white py-24 text-center text-slate-500">Memuat detail berita...</section>}
      {error && <section className="bg-white py-24 text-center text-red-600 font-medium">{error}</section>}

      {!loading && !error && article && (
        <>
          <section className="bg-white pt-8 pb-6 md:pt-12 md:pb-8">
            <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
              <div className="flex flex-wrap items-center gap-3 mb-5">
                <span className="inline-flex items-center gap-1.5 bg-brand-50 text-brand-700 text-[11px] font-bold px-3.5 py-1.5 rounded-full border border-brand-100/80 uppercase tracking-wider">
                  <i className="ph-bold ph-newspaper" /> {article.category || 'Informasi'}
                </span>
                <span className="inline-flex items-center gap-1.5 text-slate-400 text-xs font-medium">
                  <i className="ph-bold ph-calendar-blank" /> {formatDate(article.published_date)}
                </span>
                <span className="inline-flex items-center gap-1.5 text-slate-400 text-xs font-medium">
                  <i className="ph-bold ph-user" /> {article.author_name || beritaContent.detail.authorName}
                </span>
              </div>

              <h1 className="font-heading text-3xl md:text-4xl lg:text-[2.75rem] font-extrabold text-slate-900 leading-tight tracking-tight mb-6">
                {article.title}
              </h1>

              {article.image_url && (
                <div className="relative overflow-hidden rounded-2xl mb-8">
                  <img src={article.image_url} alt={article.title} className="w-full h-64 md:h-96 lg:h-[480px] object-cover rounded-2xl" />
                  <div className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/30 to-transparent h-24 rounded-b-2xl" />
                </div>
              )}
            </div>
          </section>

          <section className="pb-16 bg-white">
            <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
              <div className="flex flex-col lg:flex-row gap-10">
                <article className="w-full lg:w-2/3">
                  <div className="article-content text-[15px] sm:text-base" dangerouslySetInnerHTML={{ __html: article.content || article.excerpt || '' }} />

                  <div className="mt-10 pt-8 border-t border-slate-100">
                    <div className="bg-slate-50 rounded-2xl p-6 border border-slate-100">
                      <p className="text-xs font-bold text-brand-600 uppercase tracking-widest mb-3">{beritaContent.detail.authorBadge}</p>
                      <div className="flex items-center gap-4">
                        <img src={resolveAssetUrl(settings.logo_path)} alt={settings.site_name} className="w-12 h-12 rounded-xl object-contain bg-white p-1.5 border border-slate-200 shadow-sm" />
                        <div>
                          <p className="font-heading font-bold text-slate-900 text-sm">{beritaContent.detail.authorName}</p>
                          <p className="text-slate-500 text-xs font-medium">{beritaContent.detail.authorDescription}</p>
                        </div>
                      </div>
                    </div>
                  </div>

                  {article.tags?.length > 0 && (
                    <div className="mt-8 flex flex-wrap gap-2">
                      {article.tags.map((tag) => (
                        <span key={tag} className="inline-flex items-center gap-1 text-xs font-semibold text-slate-500 bg-slate-100 px-3 py-1.5 rounded-full hover:bg-brand-50 hover:text-brand-600 transition cursor-pointer">
                          <i className="ph-bold ph-hash text-[10px]" /> {tag}
                        </span>
                      ))}
                    </div>
                  )}
                </article>

                <aside className="w-full lg:w-1/3">
                  <div className="sticky top-28 space-y-6">
                    <div className="bg-white rounded-2xl p-5 border border-slate-100 shadow-sm">
                      <h4 className="font-heading font-bold text-slate-900 text-sm mb-4 flex items-center gap-2">
                        <i className="ph-bold ph-share-network text-brand-600" /> {beritaContent.detail.shareTitle}
                      </h4>
                      <div className="grid grid-cols-3 gap-2">
                        <button onClick={() => shareTo('facebook')} className="flex items-center justify-center w-full h-11 rounded-xl bg-[#1877F2]/10 text-[#1877F2] hover:bg-[#1877F2] hover:text-white transition" title="Facebook"><i className="ph-fill ph-facebook-logo text-lg" /></button>
                        <button onClick={() => shareTo('whatsapp')} className="flex items-center justify-center w-full h-11 rounded-xl bg-[#25D366]/10 text-[#25D366] hover:bg-[#25D366] hover:text-white transition" title="WhatsApp"><i className="ph-fill ph-whatsapp-logo text-lg" /></button>
                        <button onClick={() => shareTo('instagram')} className="flex items-center justify-center w-full h-11 rounded-xl bg-[#E1306C]/10 text-[#E1306C] hover:bg-[#E1306C] hover:text-white transition" title="Instagram"><i className="ph-fill ph-instagram-logo text-lg" /></button>
                      </div>
                    </div>
                  </div>
                </aside>
              </div>
            </div>
          </section>
        </>
      )}
    </PublicLayout>
  )
}
