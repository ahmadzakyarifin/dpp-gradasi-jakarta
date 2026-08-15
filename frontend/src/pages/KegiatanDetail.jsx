import { useState, useEffect } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import PublicLayout from '../layouts/PublicLayout'
import { kegiatanService } from '../services/kegiatanService'
import { resolveAssetUrl } from '../utils/assetUrl'
import { useSettings } from '../context/useSettings'
import { shareContent, getShareUrl, copyToClipboard } from '../utils/share'
import ArticleSource from '../components/ArticleSource'
import ArticleCaption from '../components/ArticleCaption'

export default function KegiatanDetail() {
  const { settings } = useSettings()
  const { slug } = useParams()
  const [kegiatan, setKegiatan] = useState(null)
  const [loading, setLoading] = useState(true)
  const navigate = useNavigate()
  const [error, setError] = useState(null)
  const [copied, setCopied] = useState(false)
  const [activeImageIndex, setActiveImageIndex] = useState(null)

  useEffect(() => {
    setLoading(true)
    kegiatanService.detailBySlug(slug)
      .then(res => {
        if (res.success && res.data) {
          setKegiatan(res.data)
        } else {
          setError('Kegiatan tidak ditemukan')
        }
      })
      .catch(() => setError('Terjadi kesalahan koneksi'))
      .finally(() => setLoading(false))
  }, [slug])

  useEffect(() => {
    if (kegiatan?.title) {
      document.title = `${kegiatan.title} - DPP GRADASI`
    }
    return () => {
      document.title = 'DPP GRADASI - Generasi Digital Indonesia'
    }
  }, [kegiatan])


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

  // Cek apakah data benar-benar masih "baru" (diupdate dalam 7 hari terakhir)
  const isActuallyNew = (item) => {
    if (!item || !item.is_new || !item.updated_at) return false
    const updateTime = new Date(item.updated_at).getTime()
    const now = new Date().getTime()
    const diffDays = (now - updateTime) / (1000 * 60 * 60 * 24)
    return diffDays <= 7
  }


  return (
    <PublicLayout>
      <div className="fixed top-0 left-0 h-[3px] bg-gradient-to-r from-brand-600 to-brand-400 z-[9999] w-full" />
      <section className="pt-28 pb-0 bg-white border-b border-slate-100">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 py-4">
            <nav className="flex items-center gap-2 text-sm">
              <Link to="/" className="text-slate-400 hover:text-brand-600 transition flex items-center gap-1">
                <i className="ph-bold ph-house text-xs" /> Beranda
              </Link>
              <i className="ph-bold ph-caret-right text-[10px] text-slate-300" />
              <Link to="/kegiatan" className="text-slate-400 hover:text-brand-600 transition">Kegiatan</Link>
              <i className="ph-bold ph-caret-right text-[10px] text-slate-300" />
              <span className="text-brand-600 font-semibold truncate max-w-[250px]">{kegiatan?.title || 'Detail Kegiatan'}</span>
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

      {loading && <section className="bg-white py-24 text-center text-slate-500">Memuat detail kegiatan...</section>}
      {error && (
        <section className="bg-white py-24 text-center">
          <p className="text-red-600 font-medium mb-4">{error}</p>
          <Link to="/kegiatan" className="text-brand-600 hover:underline">Kembali ke Daftar Kegiatan</Link>
        </section>
      )}

      {!loading && !error && kegiatan && (
        <>
          <section className="bg-white pt-8 pb-6 md:pt-12 md:pb-8">
            <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
              <div className="flex flex-wrap items-center gap-3 mb-5">
                <span className="inline-flex items-center gap-1.5 bg-brand-50 text-brand-700 text-[11px] font-bold px-3.5 py-1.5 rounded-full border border-brand-100/80 uppercase tracking-wider">
                  <i className="ph-bold ph-calendar-blank" /> {kegiatan.category || 'Kegiatan'}
                </span>
                <span className="flex items-start gap-1.5 text-slate-400 text-xs font-medium max-w-full">
                  <i className="ph-bold ph-calendar shrink-0 mt-0.5" /> <span className="break-all">{kegiatan.event_date}</span>
                </span>
                <span className="flex items-start gap-1.5 text-slate-400 text-xs font-medium max-w-full">
                  <i className="ph-bold ph-map-pin shrink-0 mt-0.5" /> <span className="break-all">{kegiatan.location}</span>
                </span>
                <span className="flex items-start gap-1.5 text-slate-400 text-xs font-medium max-w-full">
                  <i className="ph-bold ph-buildings shrink-0 mt-0.5" /> <span className="break-all">Penyelenggara: {kegiatan.organizer}</span>
                </span>
                {kegiatan.author_name && (
                  <span className="flex items-start gap-1.5 text-slate-400 text-xs font-medium max-w-full">
                    <i className="ph-bold ph-user shrink-0 mt-0.5" /> <span className="break-all">Penulis: {kegiatan.author_name}</span>
                  </span>
                )}
              </div>

              <h1 className="font-heading text-3xl md:text-4xl lg:text-[2.75rem] font-extrabold text-slate-900 leading-tight tracking-tight mb-6 flex flex-wrap items-center gap-3">
                <span className="min-w-0 break-all">{kegiatan.title}</span>
                {isActuallyNew(kegiatan) && (
                  <span className="bg-red-500 text-white text-[12px] font-extrabold px-3 py-1.5 rounded animate-pulse shadow-[0_0_10px_rgba(239,68,68,0.5)] flex items-center gap-1.5 shrink-0 translate-y-[-2px] uppercase tracking-wider">
                    <i className="ph-fill ph-fire" /> TERBARU
                  </span>
                )}
              </h1>

              {kegiatan.image_url && (
                <div className="mb-8">
                  <div className="relative overflow-hidden rounded-2xl">
                    <img src={resolveAssetUrl(kegiatan.image_url)} alt={kegiatan.title} className="w-full h-64 md:h-96 lg:h-[480px] object-cover rounded-2xl" />
                    <div className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/30 to-transparent h-24 rounded-b-2xl" />
                  </div>
                      <ArticleCaption source={kegiatan.image_source} />
                </div>
              )}
            </div>
          </section>

          <section className="pb-16 bg-white">
            <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
              <div className="flex flex-col lg:flex-row gap-10">
                <article className="w-full lg:w-2/3 min-w-0 space-y-8">
                  {kegiatan.excerpt && (
                    <p className="text-slate-500 font-medium italic border-l-4 border-brand-500 pl-4 py-1 text-lg [word-break:break-word]">
                      {kegiatan.excerpt}
                    </p>
                  )}

                  <div
                    className="prose max-w-none text-slate-700 leading-relaxed [word-break:break-word] whitespace-pre-line"
                    dangerouslySetInnerHTML={{ __html: kegiatan.content }}
                  />

                  <ArticleSource footnote={kegiatan.footnote} />

                  {/* Gallery Section */}
                  {kegiatan.gallery && kegiatan.gallery.length > 0 && (
                    <div className="space-y-5 pt-8 border-t border-slate-100">
                      <div>
                        <h3 className="font-heading font-bold text-lg text-slate-900">Galeri Foto Kegiatan</h3>
                        <p className="text-xs text-slate-400 font-medium">Klik gambar untuk melihat dalam ukuran penuh.</p>
                      </div>
                      <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
                        {kegiatan.gallery.map((img, idx) => (
                          <div
                            key={img.id}
                            onClick={() => setActiveImageIndex(idx)}
                            className="group relative rounded-2xl overflow-hidden border border-slate-200 shadow-sm aspect-video bg-slate-50 cursor-pointer"
                          >
                            <img
                              src={resolveAssetUrl(img.image_url || img.image_path)}
                              alt={img.caption || kegiatan.title}
                              className="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500"
                            />
                            <div className="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity duration-300 flex items-end p-3.5">
                              <span className="text-[11px] text-white font-medium truncate w-full flex items-center gap-1.5 transform translate-y-2 group-hover:translate-y-0 transition-transform duration-300">
                                <i className="ph ph-magnifying-glass-plus text-sm" />
                                {img.caption || 'Lihat foto'}
                              </span>
                            </div>
                          </div>
                        ))}
                      </div>

                      {/* Premium Lightbox Modal */}
                      {activeImageIndex !== null && (
                        <div className="fixed inset-0 z-[9999] flex items-center justify-center p-4 animate-fade-in">
                          <div
                            className="fixed inset-0 bg-slate-950/80 backdrop-blur-md transition-opacity"
                            onClick={() => setActiveImageIndex(null)}
                          />
                          <div className="relative max-w-4xl w-full max-h-[85vh] flex flex-col justify-center items-center z-10 animate-scale-up">
                            {/* Close Button */}
                            <button
                              onClick={() => setActiveImageIndex(null)}
                              className="absolute -top-12 right-0 w-9 h-9 rounded-full bg-white/10 hover:bg-white/20 text-white flex items-center justify-center transition border border-white/10 shadow"
                              title="Tutup"
                            >
                              <i className="ph-bold ph-x text-lg" />
                            </button>

                            {/* Main Image */}
                            <div className="relative rounded-2xl overflow-hidden shadow-2xl border border-white/5 bg-slate-900 max-h-[75vh]">
                              <img
                                src={resolveAssetUrl(kegiatan.gallery[activeImageIndex].image_url || kegiatan.gallery[activeImageIndex].image_path)}
                                alt={kegiatan.gallery[activeImageIndex].caption || 'Detail foto'}
                                className="w-full h-auto max-h-[75vh] object-contain mx-auto"
                              />
                              {kegiatan.gallery[activeImageIndex].caption && (
                                <div className="absolute bottom-0 left-0 right-0 bg-black/60 backdrop-blur-xs text-white text-center py-3.5 px-6 text-sm font-semibold border-t border-white/5">
                                  {kegiatan.gallery[activeImageIndex].caption}
                                </div>
                              )}
                            </div>

                            {/* Prev / Next controls */}
                            {kegiatan.gallery.length > 1 && (
                              <div className="absolute top-1/2 -translate-y-1/2 left-0 right-0 flex justify-between pointer-events-none px-4">
                                <button
                                  type="button"
                                  onClick={(e) => {
                                    e.stopPropagation()
                                    setActiveImageIndex(prev => (prev > 0 ? prev - 1 : kegiatan.gallery.length - 1))
                                  }}
                                  className="pointer-events-auto w-11 h-11 rounded-full bg-white/15 hover:bg-white/25 text-white flex items-center justify-center transition border border-white/10 shadow-lg"
                                >
                                  <i className="ph-bold ph-caret-left text-lg" />
                                </button>
                                <button
                                  type="button"
                                  onClick={(e) => {
                                    e.stopPropagation()
                                    setActiveImageIndex(prev => (prev < kegiatan.gallery.length - 1 ? prev + 1 : 0))
                                  }}
                                  className="pointer-events-auto w-11 h-11 rounded-full bg-white/15 hover:bg-white/25 text-white flex items-center justify-center transition border border-white/10 shadow-lg"
                                >
                                  <i className="ph-bold ph-caret-right text-lg" />
                                </button>
                              </div>
                            )}
                          </div>
                        </div>
                      )}
                    </div>
                  )}

                  {/* Author Card */}
                  <div className="mt-10 pt-8 border-t border-slate-100">
                    <div className="bg-slate-50 rounded-2xl p-6 border border-slate-100">
                      <p className="text-xs font-bold text-brand-600 uppercase tracking-widest mb-3">SALAM KOLABORASI</p>
                      <div className="flex items-center gap-4">
                        <img src={resolveAssetUrl(settings.logo_path)} alt={settings.site_name} className="w-12 h-12 rounded-xl object-contain bg-white p-1.5 border border-slate-200 shadow-sm" />
                        <div>
                          <p className="font-heading font-bold text-slate-900 text-sm">DPP Gradasi</p>
                          <p className="text-slate-500 text-xs font-medium">Tim Redaksi & Publikasi DPP Gradasi</p>
                        </div>
                      </div>
                    </div>
                  </div>

                  {/* Tags */}
                  {kegiatan.tags && kegiatan.tags.length > 0 && (
                    <div className="mt-8 flex flex-wrap gap-2">
                      {kegiatan.tags.map((tag) => (
                        <span key={tag} className="inline-flex items-start gap-1 text-xs font-semibold text-slate-500 bg-slate-100 px-3 py-1.5 rounded-2xl hover:bg-brand-50 hover:text-brand-600 transition cursor-pointer max-w-full">
                          <i className="ph-bold ph-hash text-[10px] shrink-0 mt-1" /> <span className="break-all whitespace-normal leading-relaxed">{tag}</span>
                        </span>
                      ))}
                    </div>
                  )}
                </article>

                <aside className="w-full lg:w-1/3">
                  <div className="sticky top-28 space-y-6">
                    <div className="bg-white rounded-2xl p-5 border border-slate-100 shadow-sm animate-fade-in">
                      <h4 className="font-heading font-bold text-slate-900 text-sm mb-4 flex items-center gap-2">
                        <i className="ph-bold ph-share-network text-brand-600" /> Bagikan Kegiatan
                      </h4>
                      <div className="grid grid-cols-4 gap-2">
                        <button onClick={handleInstagramShare} className="flex items-center justify-center w-full h-11 rounded-xl bg-[#E1306C]/10 text-[#E1306C] hover:bg-[#E1306C] hover:text-white transition" title="Instagram"><i className="ph-fill ph-instagram-logo text-lg" /></button>
                        <a href={getShareUrl('whatsapp', { title: kegiatan.title, text: kegiatan.excerpt })} target="_blank" rel="noopener noreferrer" className="flex items-center justify-center w-full h-11 rounded-xl bg-[#25D366]/10 text-[#25D366] hover:bg-[#25D366] hover:text-white transition" title="WhatsApp"><i className="ph-fill ph-whatsapp-logo text-lg" /></a>
                        <a href={getShareUrl('facebook', { title: kegiatan.title, text: kegiatan.excerpt })} target="_blank" rel="noopener noreferrer" className="flex items-center justify-center w-full h-11 rounded-xl bg-[#1877F2]/10 text-[#1877F2] hover:bg-[#1877F2] hover:text-white transition" title="Facebook"><i className="ph-fill ph-facebook-logo text-lg" /></a>
                        <button
                          onClick={handleCopyLink}
                          className={`flex items-center justify-center w-full h-11 rounded-xl transition ${copied ? 'bg-emerald-500 text-white' : 'bg-slate-100 text-slate-600 hover:bg-slate-600 hover:text-white'}`}
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

