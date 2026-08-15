import { useState, useEffect } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import PublicLayout from '../layouts/PublicLayout'
import { pengurusService } from '../services/pengurusService'
import { resolveAssetUrl } from '../utils/assetUrl'
import { copyToClipboard } from '../utils/share'

export default function KepengurusanDetail() {
  const { id } = useParams()
  const [pengurus, setPengurus] = useState(null)
  const [loading, setLoading] = useState(true)
  const navigate = useNavigate()
  const [error, setError] = useState(null)
  const [copied, setCopied] = useState(false)

  useEffect(() => {
    setLoading(true)
    pengurusService.detailById(id)
      .then(res => {
        if (res.success && res.data) {
          setPengurus(res.data)
        } else {
          setError('Data pengurus tidak ditemukan')
        }
      })
      .catch(() => setError('Terjadi kesalahan saat memuat profil pengurus.'))
      .finally(() => setLoading(false))
  }, [id])

  useEffect(() => {
    if (pengurus?.name) {
      document.title = `${pengurus.name} — Kepengurusan DPP GRADASI`
    }
    return () => { document.title = 'DPP GRADASI — Generasi Digital Indonesia' }
  }, [pengurus])

  const handleShare = () => {
    copyToClipboard(window.location.href).then(ok => {
      if (ok) { setCopied(true); setTimeout(() => setCopied(false), 2500) }
    })
  }

  const getLevelLabel = (level, pengurus) => {
    if (!level) return 'Pengurus';
    let base = level;
    let suffix = '';
    switch (level.toLowerCase()) {
      case 'ketua umum': base = 'Ketua Umum'; break;
      case 'pengurus pusat': base = 'Pengurus Pusat'; break;
      case 'pengurus provinsi': base = 'Pengurus Provinsi'; break;
      case 'pengurus kab/kota': base = 'Pengurus Kab/Kota'; break;
    }
    
    if (level.toLowerCase() === 'pengurus provinsi' && pengurus?.provinsi) {
      suffix = ` ${pengurus.provinsi}`;
    }
    if (level.toLowerCase() === 'pengurus kab/kota' && pengurus?.kabupaten) {
      suffix = ` ${pengurus.kabupaten}`;
    }
    return base + suffix;
  }

  const socials = pengurus ? [
    { url: pengurus.cv_path ? resolveAssetUrl(pengurus.cv_path) : null, icon: 'ph-file-text', label: 'CV' },
    { url: pengurus.facebook_url, icon: 'ph-facebook-logo', label: 'Facebook' },
    { url: pengurus.instagram_url, icon: 'ph-instagram-logo', label: 'Instagram' },
    { url: pengurus.linkedin_url, icon: 'ph-linkedin-logo', label: 'LinkedIn' },
    { url: pengurus.twitter_url, icon: 'ph-x-logo', label: 'X (Twitter)' },
    { url: pengurus.whatsapp ? `https://wa.me/${pengurus.whatsapp.replace(/\D/g, '')}` : null, icon: 'ph-whatsapp-logo', label: 'WhatsApp' },
  ].filter(s => s.url) : []

  return (
    <PublicLayout>
      <div className="min-h-screen bg-slate-50 pt-24 pb-20 selection:bg-brand-500 selection:text-white">
        
        {/* Loading State */}
        {loading && (
          <div className="flex flex-col items-center justify-center min-h-[60vh]">
            <i className="ph ph-spinner-gap animate-spin text-4xl text-brand-500 mb-4" />
            <p className="text-sm font-medium text-slate-500 tracking-wide uppercase">Memuat Profil...</p>
          </div>
        )}

        {/* Error State */}
        {error && (
          <div className="flex flex-col items-center justify-center min-h-[60vh] px-4 text-center">
            <div className="bg-white p-8 rounded-2xl shadow-sm border border-slate-100 max-w-md w-full">
              <div className="w-16 h-16 bg-red-50 text-red-500 rounded-full flex items-center justify-center mx-auto mb-6">
                <i className="ph ph-warning-circle text-3xl" />
              </div>
              <h3 className="text-lg font-bold text-slate-900 mb-2">Terjadi Kesalahan</h3>
              <p className="text-slate-500 mb-8 leading-relaxed">{error}</p>
              <button onClick={() => navigate('/kepengurusan')} className="w-full bg-slate-900 hover:bg-slate-800 text-white font-medium py-3 px-6 rounded-xl transition-colors">
                Kembali ke Daftar
              </button>
            </div>
          </div>
        )}

        {/* Main Content */}
        {!loading && !error && pengurus && (
          <main className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
            
            {/* Breadcrumb Navigation */}
            <nav className="flex items-center gap-2 text-xs font-medium text-slate-400 mb-10 tracking-wide">
              <Link to="/" className="hover:text-slate-900 transition-colors">BERANDA</Link>
              <span className="text-slate-300">/</span>
              <Link to="/kepengurusan" className="hover:text-slate-900 transition-colors">KEPENGURUSAN</Link>
              <span className="text-slate-300">/</span>
              <span className="text-brand-600">PROFIL</span>
            </nav>

            <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 lg:gap-12 items-start">
              
              {/* Left Column: Photo & Quick Actions */}
              <div className="lg:col-span-4 lg:sticky lg:top-32">
                <div className="bg-white rounded-3xl p-3 shadow-sm border border-slate-100 mb-6">
                  <div className="aspect-[3/4] rounded-2xl overflow-hidden bg-slate-100 relative group">
                    <img 
                      src={resolveAssetUrl(pengurus.image_url || pengurus.image_path)} 
                      alt={pengurus.name}
                      className="w-full h-full object-cover object-top transition-transform duration-700 group-hover:scale-105"
                    />
                    <div className="absolute inset-0 bg-gradient-to-t from-slate-900/60 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500" />
                  </div>
                </div>

                <div className="flex flex-col gap-3">
                  <button 
                    onClick={handleShare}
                    className={`w-full py-4 px-6 rounded-2xl font-medium text-sm flex items-center justify-center gap-3 transition-all duration-300 ${
                      copied 
                        ? 'bg-emerald-50 text-emerald-600 border border-emerald-200' 
                        : 'bg-white text-slate-700 border border-slate-200 hover:border-brand-300 hover:text-brand-600 hover:shadow-sm'
                    }`}
                  >
                    <i className={`text-lg ${copied ? 'ph-fill ph-check-circle' : 'ph ph-share-network'}`} />
                    {copied ? 'Tautan Tersalin' : 'Bagikan Profil'}
                  </button>

                  {socials.length > 0 && (
                    <div className="flex items-center justify-center gap-3 p-4 bg-white rounded-2xl border border-slate-100 shadow-sm">
                      {socials.map((s, i) => (
                        <a
                          key={i}
                          href={s.url}
                          target="_blank"
                          rel="noreferrer"
                          title={s.label}
                          className="w-10 h-10 flex items-center justify-center rounded-xl bg-slate-50 text-slate-400 hover:bg-brand-50 hover:text-brand-600 transition-colors"
                        >
                          {s.label === 'CV' ? (
                            <span className="text-sm font-bold font-heading tracking-wide">CV</span>
                          ) : (
                            <i className={`text-xl ph-fill ${s.icon}`} />
                          )}
                        </a>
                      ))}
                    </div>
                  )}
                </div>
              </div>

              {/* Right Column: Details */}
              <div className="lg:col-span-8">
                
                {/* Header Info */}
                <div className="bg-white rounded-3xl p-8 md:p-12 shadow-sm border border-slate-100 mb-8">
                  <div className="inline-flex items-center gap-2 px-3 py-1.5 rounded-lg bg-slate-50 border border-slate-100 text-xs font-semibold text-slate-500 tracking-wide uppercase mb-6">
                    <span className="w-1.5 h-1.5 rounded-full bg-brand-500" />
                    {getLevelLabel(pengurus.level, pengurus)}
                  </div>
                  
                  <h1 className="font-heading text-4xl md:text-5xl font-extrabold text-slate-900 mb-4 tracking-tight leading-tight">
                    {pengurus.name}
                  </h1>
                  <h2 className="text-xl md:text-2xl text-brand-600 font-medium mb-8 pb-8 border-b border-slate-100">
                    {pengurus.role}
                  </h2>

                  <div className="flex flex-col gap-5 text-sm md:text-base">
                    <div className="grid grid-cols-12 gap-4">
                      <div className="col-span-5 md:col-span-4 lg:col-span-3 font-bold text-slate-900">Email</div>
                      <div className="col-span-7 md:col-span-8 lg:col-span-9 text-slate-600 font-medium">{pengurus.email || '-'}</div>
                    </div>
                    <div className="grid grid-cols-12 gap-4">
                      <div className="col-span-5 md:col-span-4 lg:col-span-3 font-bold text-slate-900">No Hp</div>
                      <div className="col-span-7 md:col-span-8 lg:col-span-9 text-slate-600 font-medium">{pengurus.whatsapp || '-'}</div>
                    </div>
                    <div className="grid grid-cols-12 gap-4">
                      <div className="col-span-5 md:col-span-4 lg:col-span-3 font-bold text-slate-900">Tingkat Struktur</div>
                      <div className="col-span-7 md:col-span-8 lg:col-span-9 text-slate-600 font-medium capitalize">{getLevelLabel(pengurus.level, pengurus)}</div>
                    </div>
                    {pengurus.kepengurusan && (
                      <div className="grid grid-cols-12 gap-4">
                        <div className="col-span-5 md:col-span-4 lg:col-span-3 font-bold text-slate-900">Kepengurusan</div>
                        <div className="col-span-7 md:col-span-8 lg:col-span-9 text-slate-600 font-medium">{pengurus.kepengurusan}</div>
                      </div>
                    )}
                    {pengurus.periode && (
                      <div className="grid grid-cols-12 gap-4">
                        <div className="col-span-5 md:col-span-4 lg:col-span-3 font-bold text-slate-900">Periode</div>
                        <div className="col-span-7 md:col-span-8 lg:col-span-9 text-slate-600 font-medium">{pengurus.periode}</div>
                      </div>
                    )}
                    {pengurus.department && (
                      <div className="grid grid-cols-12 gap-4">
                        <div className="col-span-5 md:col-span-4 lg:col-span-3 font-bold text-slate-900">Departemen</div>
                        <div className="col-span-7 md:col-span-8 lg:col-span-9 text-slate-600 font-medium">{pengurus.department}</div>
                      </div>
                    )}

                    <div className="grid grid-cols-12 gap-4">
                      <div className="col-span-5 md:col-span-4 lg:col-span-3 font-bold text-slate-900">Pekerjaan</div>
                      <div className="col-span-7 md:col-span-8 lg:col-span-9 text-slate-600 font-medium">{pengurus.pekerjaan || '-'}</div>
                    </div>
                  </div>
                </div>

                {/* Biografi */}
                <div className="mb-12">
                  <h3 className="flex items-center gap-3 text-lg font-bold text-slate-900 mb-6">
                    <span className="w-8 h-8 rounded-lg bg-slate-100 text-slate-500 flex items-center justify-center">
                      <i className="ph-fill ph-identification-card" />
                    </span>
                    Biografi
                  </h3>
                  {pengurus.bio ? (
                    <div className="prose prose-slate prose-lg max-w-none text-slate-600 leading-relaxed">
                      {pengurus.bio.split('\\n').map((paragraph, index) => (
                        <p key={index} className="mb-4">{paragraph}</p>
                      ))}
                    </div>
                  ) : (
                    <p className="text-slate-400 italic">Belum ada informasi biografi yang ditambahkan.</p>
                  )}
                </div>

                <div className="grid md:grid-cols-2 gap-12">
                  {/* Pendidikan */}
                  <div>
                    <h3 className="flex items-center gap-3 text-lg font-bold text-slate-900 mb-6">
                      <span className="w-8 h-8 rounded-lg bg-slate-100 text-slate-500 flex items-center justify-center">
                        <i className="ph-fill ph-graduation-cap" />
                      </span>
                      Pendidikan
                    </h3>
                    {pengurus.pendidikan ? (
                      <div className="prose prose-slate max-w-none text-slate-600 leading-relaxed text-sm">
                        {pengurus.pendidikan.split(/\n\s*\n/).filter(item => item.trim()).map((item, index) => (
                          <div key={index} className="flex items-start gap-3 mb-3">
                            <i className="ph-fill ph-caret-right text-brand-500 mt-1 flex-shrink-0" />
                            <p className="m-0 whitespace-pre-line">{item.trim()}</p>
                          </div>
                        ))}
                      </div>
                    ) : (
                      <p className="text-slate-400 italic text-sm">Informasi pendidikan terlampir pada CV.</p>
                    )}
                  </div>

                  {/* Sertifikasi */}
                  <div>
                    <h3 className="flex items-center gap-3 text-lg font-bold text-slate-900 mb-6">
                      <span className="w-8 h-8 rounded-lg bg-slate-100 text-slate-500 flex items-center justify-center">
                        <i className="ph-fill ph-certificate" />
                      </span>
                      Sertifikasi
                    </h3>
                    {pengurus.sertifikasi ? (
                      <div className="prose prose-slate max-w-none text-slate-600 leading-relaxed text-sm">
                        {pengurus.sertifikasi.split(/\n\s*\n/).filter(item => item.trim()).map((item, index) => (
                          <div key={index} className="flex items-start gap-3 mb-3">
                            <i className="ph-fill ph-star text-amber-400 mt-1 flex-shrink-0" />
                            <p className="m-0 whitespace-pre-line">{item.trim()}</p>
                          </div>
                        ))}
                      </div>
                    ) : (
                      <p className="text-slate-400 italic text-sm">Informasi sertifikasi terlampir pada CV.</p>
                    )}
                  </div>
                </div>

              </div>
            </div>
          </main>
        )}
      </div>
    </PublicLayout>
  )
}
