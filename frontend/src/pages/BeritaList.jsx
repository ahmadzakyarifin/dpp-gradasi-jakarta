import { useEffect, useState } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import PublicLayout from '../layouts/PublicLayout'
import { beritaContent } from '../content/beritaContent'
import { beritaService } from '../services/beritaService'
import { formatDate } from '../utils/format'

export default function BeritaList() {
  const [searchParams, setSearchParams] = useSearchParams()
  const page = Number(searchParams.get('page') || 1)
  const [items, setItems] = useState([
    {
      id: 1,
      title: 'Rapat Kerja Daerah Jatim',
      slug: 'rapat-kerja-daerah-jatim',
      category: 'Berita Daerah',
      published_date: '2026-02-11',
      image_url: 'https://gradasi.org/uploads/img/berita/17708152730.jpg',
      excerpt: 'SURABAYA, Generasi Digital Indonesia (GRADASI) Jawa Timur bersiap menggelar Rapat Kerja Daerah untuk menyelaraskan program kerja digitalisasi UMKM.'
    },
    {
      id: 2,
      title: 'Peningkatan Kompetensi SDM',
      slug: 'peningkatan-kompetensi-sdm-pendidikan',
      category: 'Edukasi',
      published_date: '2025-11-02',
      image_url: 'https://gradasi.org/uploads/img/berita/17620765070.jpg',
      excerpt: 'Inisiatif GRADASI Mendorong Peningkatan Kompetensi SDM Pendidikan dalam Memanfaatkan Kecerdasan Buatan (AI) secara bijak.'
    },
    {
      id: 3,
      title: 'Rumusan Kunci Kebijakan',
      slug: 'rumusan-kunci-kebijakan-literasi-digital',
      category: 'Berita Utama',
      published_date: '2025-10-31',
      image_url: 'https://gradasi.org/uploads/img/berita/17618789900.jpg',
      excerpt: '#Ketua Dewan Pakar GRADASI, Damar Juniarto, Paparkan Lima Rumusan Kunci Kebijakan untuk Mempercepat Transformasi Digital.'
    },
    {
      id: 4,
      title: 'Rapat Strategis Pengurus Pusat',
      slug: 'rapat-strategis-pengurus-pusat',
      category: 'Nasional',
      published_date: '2025-10-20',
      image_url: 'https://gradasi.org/uploads/img/berita/17708152730.jpg',
      excerpt: 'Pusat pelaporan kegiatan dalam rangka mempersiapkan agenda strategis organisasi untuk tahun 2026.'
    },
    {
      id: 5,
      title: 'Audiensi Lanjutan Kementerian',
      slug: 'audiensi-lanjutan-kementerian',
      category: 'Kemitraan',
      published_date: '2025-09-15',
      image_url: 'https://gradasi.org/uploads/img/berita/17620765070.jpg',
      excerpt: 'Laporan singkat dari hasil pertemuan pimpinan pusat bersama kementerian untuk program literasi tahap dua.'
    },
    {
      id: 6,
      title: 'MOU Ekosistem Digital Kota',
      slug: 'mou-ekosistem-digital-kota',
      category: 'Kerjasama',
      published_date: '2025-08-05',
      image_url: 'https://gradasi.org/uploads/img/berita/17618789900.jpg',
      excerpt: 'Meningkatkan sinergi antara DPD setempat dengan pemerintah kota dalam menyepakati ekosistem digital bersama.'
    }
  ])
  const [meta, setMeta] = useState({ current_page: 1, total_pages: 3, total_data: 18 })
  const [currentPage, setCurrentPage] = useState(page)

  useEffect(() => {
    beritaService.list({ page: currentPage, limit: 6, sort: 'newest' })
      .then(res => {
        if (res.data && res.data.berita && res.data.berita.length > 0) {
          setItems(res.data.berita)
          setMeta(res.data.meta || { current_page: currentPage, total_pages: Math.max(3, res.data.meta?.total_pages || 1), total_data: res.data.berita.length })
        }
      })
      .catch(() => {})
  }, [currentPage])

  function goToPage(nextPage) {
    setCurrentPage(nextPage)
    setSearchParams({ page: String(nextPage) })
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }

  const itemsPerPage = 3
  const totalPages = Math.max(3, Math.ceil(items.length / itemsPerPage))
  const paginatedItems = items.slice((currentPage - 1) * itemsPerPage, currentPage * itemsPerPage)

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
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {paginatedItems.map((item) => (
              <article key={item.id} className="group cursor-pointer bg-white rounded-2xl overflow-hidden shadow-sm hover:shadow-[0_8px_30px_rgba(37,99,235,0.08)] transition-all duration-300 border border-slate-100 flex flex-col">
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

          {/* ALWAYS VISIBLE PAGINATION MATCHING BERITA.HTML */}
          <div className="mt-16 flex justify-center">
            <nav className="flex items-center gap-2">
              <button 
                type="button" 
                disabled={currentPage <= 1} 
                onClick={() => goToPage(currentPage - 1)} 
                className="w-10 h-10 flex items-center justify-center rounded-xl border border-slate-200 text-slate-400 hover:border-brand-500 hover:text-brand-600 hover:bg-brand-50 transition disabled:opacity-40"
              >
                <i className="ph-bold ph-caret-left" />
              </button>
              
              {[1, 2, 3].map((number) => (
                <button 
                  key={number} 
                  type="button" 
                  onClick={() => goToPage(number)} 
                  className={`w-10 h-10 flex items-center justify-center rounded-xl font-bold transition ${number === currentPage ? 'bg-brand-700 text-white shadow-md shadow-brand-500/20' : 'border border-slate-200 text-slate-600 hover:border-brand-500 hover:text-brand-600 hover:bg-brand-50'}`}
                >
                  {number}
                </button>
              ))}

              <button 
                type="button" 
                disabled={currentPage >= totalPages} 
                onClick={() => goToPage(currentPage + 1)} 
                className="w-10 h-10 flex items-center justify-center rounded-xl border border-slate-200 text-slate-400 hover:border-brand-500 hover:text-brand-600 hover:bg-brand-50 transition disabled:opacity-40"
              >
                <i className="ph-bold ph-caret-right" />
              </button>
            </nav>
          </div>
        </div>
      </section>
    </PublicLayout>
  )
}
