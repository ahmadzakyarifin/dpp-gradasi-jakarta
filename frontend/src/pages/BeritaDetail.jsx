import { useEffect, useState } from 'react'
import { Link, useParams, useNavigate } from 'react-router-dom'
import PublicLayout from '../layouts/PublicLayout'
import { useSettings } from '../context/useSettings'
import { resolveAssetUrl } from '../utils/assetUrl'
import { beritaContent } from '../content/beritaContent'
import { beritaService } from '../services/beritaService'
import { formatDate } from '../utils/format'
import { shareContent, getShareUrl, copyToClipboard } from '../utils/share'
import ArticleSource from '../components/ArticleSource'
import ArticleCaption from '../components/ArticleCaption'

export default function BeritaDetail() {
  const { settings } = useSettings()
  const { slug } = useParams()
  const [article, setArticle] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [copied, setCopied] = useState(false)
  const navigate = useNavigate()

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

  useEffect(() => {
    if (article?.title) {
      document.title = `${article.title} - DPP GRADASI`
    }
    return () => {
      document.title = 'DPP GRADASI - Generasi Digital Indonesia'
    }
  }, [article])


  const handleInstagramShare = () => {
    copyToClipboard(window.location.href).then(success => {
      if (success) {
        setCopied(true)
        setTimeout(() => setCopied(false), 2000)
        window.open('https://www.instagram.com/', '_blank')
      }
    })
  }

  const handleCopyLink = () => {
    copyToClipboard(window.location.href).then(success => {
      if (success) {
        setCopied(true)
        setTimeout(() => setCopied(false), 2000)
      }
    })
  }

  return (
    <PublicLayout>
      <div className="fixed top-0 left-0 h-[3px] bg-gradient-to-r from-brand-600 to-brand-400 z-[9999] w-full" />
      <section className="pt-28 pb-0 bg-white border-b border-slate-100">
        <div className="px-4 mx-auto max-w-4xl sm:px-6 lg:px-8">
          <div className="flex flex-col gap-3 py-4 sm:flex-row sm:items-center sm:justify-between">
            <nav className="flex gap-2 items-center text-sm">
              <Link to="/" className="flex gap-1 items-center transition text-slate-400 hover:text-brand-600">
                <i className="text-xs ph-bold ph-house" /> {beritaContent.detail.breadcrumbHome}
              </Link>
              <i className="ph-bold ph-caret-right text-[10px] text-slate-300" />
              <Link to="/berita" className="transition text-slate-400 hover:text-brand-600">{beritaContent.detail.breadcrumbInfo}</Link>
              <i className="ph-bold ph-caret-right text-[10px] text-slate-300" />
              <span className="text-brand-600 font-semibold truncate max-w-[250px]">{article?.title || beritaContent.detail.defaultTitle}</span>
            </nav>
            <button
              onClick={() => navigate('/')}
              className="inline-flex items-center gap-2 text-xs font-semibold text-slate-600 hover:text-brand-600 transition bg-slate-100 hover:bg-brand-50 px-3.5 py-2 rounded-xl border border-slate-200/60 w-fit cursor-pointer btn-press"
            >
              <i className="ph-bold ph-arrow-left text-[11px]" /> Kembali ke Beranda
            </button>
          </div>
        </div>
      </section>

      {loading && <section className="py-24 text-center bg-white text-slate-500">Memuat detail berita...</section>}
      {error && <section className="py-24 font-medium text-center text-red-600 bg-white">{error}</section>}

      {!loading && !error && article && (
        <>
          <section className="pt-8 pb-6 bg-white md:pt-12 md:pb-8">
            <div className="px-4 mx-auto max-w-4xl sm:px-6 lg:px-8">
              <div className="flex flex-wrap gap-3 items-center mb-5">
                <span className="inline-flex items-center gap-1.5 bg-brand-50 text-brand-700 text-[11px] font-bold px-3.5 py-1.5 rounded-full border border-brand-100/80 uppercase tracking-wider">
                  <i className="ph-bold ph-newspaper" /> {article.category || 'Informasi'}
                </span>
                <span className="flex items-start gap-1.5 text-slate-400 text-xs font-medium max-w-full">
                  <i className="ph-bold ph-calendar-blank shrink-0 mt-0.5" /> <span className="break-all">{formatDate(article.published_date)}</span>
                </span>
                <span className="flex items-start gap-1.5 text-slate-400 text-xs font-medium max-w-full">
                  <i className="ph-bold ph-user shrink-0 mt-0.5" /> <span className="break-all">{article.author_name || beritaContent.detail.authorName}</span>
                </span>
              </div>

              <h1 className="font-heading text-3xl md:text-4xl lg:text-[2.75rem] font-extrabold text-slate-900 leading-tight tracking-tight mb-6 break-all">
                {article.title}
              </h1>

              {article.image_url && (
                <div className="mb-8">
                  <div className="overflow-hidden relative rounded-2xl">
                    <img src={resolveAssetUrl(article.image_url)} alt={article.title} className="w-full h-64 md:h-96 lg:h-[480px] object-cover rounded-2xl" />
                    <div className="absolute right-0 bottom-0 left-0 h-24 bg-gradient-to-t to-transparent rounded-b-2xl from-black/30" />
                  </div>
                      <ArticleCaption source={article.image_source} />
                </div>
              )}
            </div>
          </section>

          <section className="pb-16 bg-white">
            <div className="px-4 mx-auto max-w-4xl sm:px-6 lg:px-8">
              <div className="flex flex-col gap-10 lg:flex-row">
                <article className="w-full min-w-0 lg:w-2/3">
                  <p className="article-content text-[15px] sm:text-base text-slate-700 leading-relaxed whitespace-pre-line break-all">
                    {article.content || article.excerpt || ''}
                  </p>

                  <ArticleSource footnote={article.footnote} />

                  <div className="pt-8 mt-10 border-t border-slate-100">
                    <div className="p-6 rounded-2xl border bg-slate-50 border-slate-100">
                      <p className="mb-3 text-xs font-bold tracking-widest uppercase text-brand-600">{beritaContent.detail.authorBadge}</p>
                      <div className="flex gap-4 items-center">
                        <img src={resolveAssetUrl(settings.logo_path)} alt={settings.site_name} className="w-12 h-12 rounded-xl object-contain bg-white p-1.5 border border-slate-200 shadow-sm" />
                        <div>logic terbaru engg mnucl 
                          <p className="text-sm font-bold font-heading text-slate-900">{beritaContent.detail.authorName}</p>
                          <p className="text-xs font-medium text-slate-500">{beritaContent.detail.authorDescription}</p>
                        </div>
                      </div>
                    </div>
                  </div>

                  {article.tags?.length > 0 && (
                    <div className="flex flex-wrap gap-2 mt-8">
                      {article.tags.map((tag) => (
                        <span key={tag} className="inline-flex items-start gap-1 text-xs font-semibold text-slate-500 bg-slate-100 px-3 py-1.5 rounded-2xl hover:bg-brand-50 hover:text-brand-600 transition cursor-pointer max-w-full">
                          <i className="ph-bold ph-hash text-[10px] shrink-0 mt-1" /> <span className="break-all whitespace-normal leading-relaxed">{tag}</span>
                        </span>
                      ))}
                    </div>
                  )}
                </article>

                <aside className="w-full lg:w-1/3">
                  <div className="sticky top-28 space-y-6">
                    <div className="p-5 bg-white rounded-2xl border shadow-sm border-slate-100 animate-fade-in">
                      <h4 className="flex gap-2 items-center mb-4 text-sm font-bold font-heading text-slate-900">
                        <i className="ph-bold ph-share-network text-brand-600" /> {beritaContent.detail.shareTitle}
                      </h4>
                      <div className="grid grid-cols-4 gap-2">
                        <button onClick={handleInstagramShare} className="flex items-center justify-center w-full h-11 rounded-xl bg-[#E1306C]/10 text-[#E1306C] hover:bg-[#E1306C] hover:text-white transition" title="Instagram"><i className="text-lg ph-fill ph-instagram-logo" /></button>
                        <a href={getShareUrl('whatsapp', { title: article.title, text: article.excerpt })} target="_blank" rel="noopener noreferrer" className="flex items-center justify-center w-full h-11 rounded-xl bg-[#25D366]/10 text-[#25D366] hover:bg-[#25D366] hover:text-white transition" title="WhatsApp"><i className="text-lg ph-fill ph-whatsapp-logo" /></a>
                        <a href={getShareUrl('facebook', { title: article.title, text: article.excerpt })} target="_blank" rel="noopener noreferrer" className="flex items-center justify-center w-full h-11 rounded-xl bg-[#1877F2]/10 text-[#1877F2] hover:bg-[#1877F2] hover:text-white transition" title="Facebook"><i className="text-lg ph-fill ph-facebook-logo" /></a>
                        <button
                          onClick={handleCopyLink}
                          className={`flex items-center justify-center w-full h-11 rounded-xl transition ${copied ? 'text-white bg-emerald-500' : 'bg-slate-100 text-slate-600 hover:bg-slate-600 hover:text-white'}`}
                          title={copied ? "Tautan Berhasil Disalin" : "Salin Tautan"}
                        >
                          <i className={`text-lg ${copied ? 'ph-bold ph-check' : 'ph-bold ph-link'}`} />
                        </button>
                      </div>
                      {copied && (
                        <p className="text-[10px] text-emerald-600 font-semibold mt-2 text-center animate-pulse">
                          Tautan berhasil disalin ke papan klip!
                        </p>
                      )}
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
