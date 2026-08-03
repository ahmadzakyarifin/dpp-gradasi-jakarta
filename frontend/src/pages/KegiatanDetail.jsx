import { useState, useEffect } from 'react'
import { useParams, Link } from 'react-router-dom'
import PublicLayout from '../layouts/PublicLayout'
import { kegiatanService } from '../services/kegiatanService'
import { resolveAssetUrl } from '../utils/assetUrl'

export default function KegiatanDetail() {
  const { slug } = useParams()
  const [kegiatan, setKegiatan] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

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

  const copyToClipboard = () => {
    navigator.clipboard.writeText(window.location.href)
    alert('Tautan disalin ke papan klip!')
  }

  return (
    <PublicLayout>
      <section className="pt-28 pb-16 md:pt-36 bg-slate-50 min-h-[90vh]">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          
          {loading && <div className="py-20 text-center text-slate-500">Memuat detail kegiatan...</div>}
          {error && (
            <div className="py-20 text-center">
              <p className="text-red-600 font-medium mb-4">{error}</p>
              <Link to="/kegiatan" className="text-brand-600 hover:underline">Kembali ke Daftar Kegiatan</Link>
            </div>
          )}

          {kegiatan && (
            <article className="bg-white rounded-3xl overflow-hidden border border-slate-100 shadow-sm p-6 md:p-10 space-y-8">
              {/* Header */}
              <div className="space-y-4">
                <span className="inline-block px-3 py-1 bg-brand-50 text-brand-700 rounded-lg text-xs font-bold uppercase tracking-wider">
                  {kegiatan.category}
                </span>
                <h1 className="font-heading text-2xl md:text-4xl font-extrabold text-slate-900 leading-tight">
                  {kegiatan.title}
                </h1>
                
                <div className="flex flex-wrap gap-4 text-xs text-slate-400 font-medium border-y border-slate-100 py-3">
                  <span>📅 {kegiatan.event_date}</span>
                  <span>📍 {kegiatan.location}</span>
                  <span>🏢 Penyelenggara: {kegiatan.organizer}</span>
                  {kegiatan.author_name && <span>✍️ Penulis: {kegiatan.author_name}</span>}
                </div>
              </div>

              {/* Main Image */}
              {kegiatan.image_url && (
                <div className="rounded-2xl overflow-hidden max-h-[450px]">
                  <img 
                    src={resolveAssetUrl(kegiatan.image_url)} 
                    alt={kegiatan.title} 
                    className="w-full h-full object-cover"
                  />
                </div>
              )}

              {/* Excerpt */}
              {kegiatan.excerpt && (
                <p className="text-slate-500 font-medium italic border-l-4 border-brand-500 pl-4 py-1 text-lg">
                  {kegiatan.excerpt}
                </p>
              )}

              {/* Content */}
              <div 
                className="prose max-w-none text-slate-700 leading-relaxed"
                dangerouslySetInnerHTML={{ __html: kegiatan.content }}
              />

              {/* Tags */}
              {kegiatan.tags && kegiatan.tags.length > 0 && (
                <div className="flex flex-wrap gap-2 pt-4 border-t border-slate-100">
                  <span className="text-xs font-bold text-slate-400 self-center mr-2">TAGS:</span>
                  {kegiatan.tags.map(tagItem => (
                    <span key={tagItem} className="px-3 py-1 bg-slate-100 text-slate-600 rounded-lg text-xs font-semibold">
                      #{tagItem}
                    </span>
                  ))}
                </div>
              )}

              {/* Gallery Section */}
              {kegiatan.gallery && kegiatan.gallery.length > 0 && (
                <div className="space-y-4 pt-6 border-t border-slate-100">
                  <h3 className="font-heading font-bold text-lg text-slate-900">Galeri Foto Kegiatan</h3>
                  <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
                    {kegiatan.gallery.map(img => (
                      <div key={img.id} className="rounded-xl overflow-hidden border border-slate-200 shadow-sm aspect-video bg-slate-50">
                        <img 
                          src={resolveAssetUrl(img.image_url || img.image_path)} 
                          alt={img.caption || kegiatan.title} 
                          className="w-full h-full object-cover hover:scale-105 transition duration-300"
                        />
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {/* Share Options */}
              <div className="flex flex-wrap gap-3 items-center justify-between pt-6 border-t border-slate-100">
                <span className="text-xs font-bold text-slate-400">BAGIKAN ARTIKEL:</span>
                <div className="flex gap-2">
                  <a 
                    href={`https://api.whatsapp.com/send?text=${encodeURIComponent(kegiatan.title + ' ' + window.location.href)}`} 
                    target="_blank" 
                    rel="noreferrer"
                    className="px-3 py-1.5 bg-green-50 text-green-700 rounded-lg text-xs font-semibold hover:bg-green-100 transition"
                  >
                    WhatsApp
                  </a>
                  <a 
                    href={`https://www.facebook.com/sharer/sharer.php?u=${encodeURIComponent(window.location.href)}`} 
                    target="_blank" 
                    rel="noreferrer"
                    className="px-3 py-1.5 bg-blue-50 text-blue-700 rounded-lg text-xs font-semibold hover:bg-blue-100 transition"
                  >
                    Facebook
                  </a>
                  <a 
                    href={`https://www.linkedin.com/sharing/share-offsite/?url=${encodeURIComponent(window.location.href)}`} 
                    target="_blank" 
                    rel="noreferrer"
                    className="px-3 py-1.5 bg-indigo-50 text-indigo-700 rounded-lg text-xs font-semibold hover:bg-indigo-100 transition"
                  >
                    LinkedIn
                  </a>
                  <button 
                    onClick={copyToClipboard}
                    className="px-3 py-1.5 bg-slate-50 text-slate-700 rounded-lg text-xs font-semibold hover:bg-slate-100 transition"
                  >
                    Salin Tautan
                  </button>
                </div>
              </div>

            </article>
          )}

        </div>
      </section>
    </PublicLayout>
  )
}
