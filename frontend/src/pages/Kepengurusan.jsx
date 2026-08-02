import { useState, useEffect } from 'react'
import { useSearchParams } from 'react-router-dom'
import PublicLayout from '../layouts/PublicLayout'
import { pengurusService } from '../services/pengurusService'

export default function Kepengurusan() {
  const [searchParams, setSearchParams] = useSearchParams()
  const activeTab = searchParams.get('tab') || 'ketua'
  
  const [allPengurus, setAllPengurus] = useState([
    // KETUA UMUM
    { id: 1, name: 'Upi Asmaradhana', role: 'Ketua Umum DPP GRADASI', level: 'ketua', is_active: true, periode: '2024 - 2029', image_url: 'https://gradasi.org/uploads/img/s-anggota/ketua/1735027418.jpg', sort_order: 1, facebook_url: 'https://facebook.com', instagram_url: 'https://instagram.com' },

    // DPP
    { id: 2, name: 'Dr. Susi Susanti, M.Pd', role: 'Wakil Ketua I', level: 'dpp', is_active: true, image_url: 'https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?q=80&w=200', sort_order: 1 },
    { id: 3, name: 'Ir. Budi Santoso', role: 'Wakil Ketua II', level: 'dpp', is_active: true, image_url: 'https://images.unsplash.com/photo-1519085360753-af0119f7cbe7?q=80&w=200', sort_order: 2 },
    { id: 4, name: 'Junaidi, S.Kom', role: 'Sekretaris Jenderal', level: 'dpp', is_active: true, image_url: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?q=80&w=200', sort_order: 3 },
    { id: 5, name: 'Dina Mariana, S.ST', role: 'Wakil Sekjen 1', level: 'dpp', is_active: true, image_url: 'https://images.unsplash.com/photo-1438761681033-6461ffad8d80?q=80&w=200', sort_order: 4 },
    { id: 6, name: 'Sudarwati', role: 'Wakil Sekjen 2', level: 'dpp', is_active: true, image_url: 'https://images.unsplash.com/photo-1544005313-94ddf0286df2?q=80&w=200', sort_order: 5 },
    { id: 7, name: 'Rina Wijaya, M.Sc', role: 'Bendahara Umum', level: 'dpp', is_active: true, image_url: 'https://images.unsplash.com/photo-1580489944761-15a19d654956?q=80&w=200', sort_order: 6 },
    { id: 8, name: 'Yoseph Budi', role: 'Wakil Bendahara', level: 'dpp', is_active: true, image_url: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?q=80&w=200', sort_order: 7 },
    { id: 9, name: 'Dwi Purnomo, S.Kom', role: 'Koordinator Dept 01 Organisasi', level: 'dpp', is_active: true, image_url: 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?q=80&w=200', sort_order: 8 },
    { id: 10, name: 'Muhammad Hertiyadi Alfaqy S.Kom', role: 'Koordinator Dept 02 IT & Digital', level: 'dpp', is_active: true, image_url: 'https://images.unsplash.com/photo-1599566150163-29194dcaad36?q=80&w=200', sort_order: 9 },

    // DPD
    { id: 11, name: 'Drs. H. Ahmad Fauzi', role: 'Ketua DPD Jawa Barat', level: 'dpd', provinsi: 'Jawa Barat', is_active: true, image_url: 'https://images.unsplash.com/photo-1560250097-0b93528c311a?q=80&w=200', sort_order: 1 },
    { id: 12, name: 'Bambang Irawan, S.T', role: 'Ketua DPD Jawa Timur', level: 'dpd', provinsi: 'Jawa Timur', is_active: true, image_url: 'https://images.unsplash.com/photo-1519085360753-af0119f7cbe7?q=80&w=200', sort_order: 2 },
    { id: 13, name: 'Siti Aminah, M.Si', role: 'Ketua DPD Jawa Tengah', level: 'dpd', provinsi: 'Jawa Tengah', is_active: true, image_url: 'https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?q=80&w=200', sort_order: 3 },
    { id: 14, name: 'Hendra Gunawan', role: 'Ketua DPD DKI Jakarta', level: 'dpd', provinsi: 'DKI Jakarta', is_active: true, image_url: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?q=80&w=200', sort_order: 4 },
    { id: 15, name: 'Tri Wahyudi', role: 'Ketua DPD Banten', level: 'dpd', provinsi: 'Banten', is_active: true, image_url: 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?q=80&w=200', sort_order: 5 },
    { id: 16, name: 'Eko Prasetyo', role: 'Ketua DPD DI Yogyakarta', level: 'dpd', provinsi: 'DI Yogyakarta', is_active: true, image_url: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?q=80&w=200', sort_order: 6 },

    // DPC
    { id: 17, name: 'Syamsul Bahri', role: 'Ketua DPC Kota Bandung', level: 'dpc', provinsi: 'Jawa Barat', kabupaten: 'Kota Bandung', is_active: true, image_url: 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?q=80&w=200', sort_order: 1 },
    { id: 18, name: 'Herman Wijaya', role: 'Ketua DPC Kab. Bogor', level: 'dpc', provinsi: 'Jawa Barat', kabupaten: 'Kabupaten Bogor', is_active: true, image_url: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?q=80&w=200', sort_order: 2 },
    { id: 19, name: 'Ridwan Malik', role: 'Ketua DPC Kota Surabaya', level: 'dpc', provinsi: 'Jawa Timur', kabupaten: 'Kota Surabaya', is_active: true, image_url: 'https://images.unsplash.com/photo-1599566150163-29194dcaad36?q=80&w=200', sort_order: 3 },
    { id: 20, name: 'Anita Rahayu', role: 'Ketua DPC Kab. Malang', level: 'dpc', provinsi: 'Jawa Timur', kabupaten: 'Kabupaten Malang', is_active: true, image_url: 'https://images.unsplash.com/photo-1544005313-94ddf0286df2?q=80&w=200', sort_order: 4 }
  ])
  const [loading, setLoading] = useState(true)
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedProvinsi, setSelectedProvinsi] = useState('')
  const [selectedKabupaten, setSelectedKabupaten] = useState('')
  const [currentPage, setCurrentPage] = useState(1)
  const itemsPerPage = 4

  useEffect(() => {
    pengurusService.list()
      .then(res => {
        if (res.success && res.data) {
          const list = Array.isArray(res.data) ? res.data : (res.data.pengurus || [])
          // Normalisasi: API kirim image_path → FE pakai image_url di JSX
          setAllPengurus(list.map(p => ({ ...p, image_url: p.image_path || p.image_url })))
        }
      })
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [])

  useEffect(() => {
    setCurrentPage(1)
  }, [activeTab, searchQuery, selectedProvinsi, selectedKabupaten])

  // Filter lists based on selected tab level
  const activeDataList = allPengurus.filter(item => item.level === activeTab && (item.is_active !== false))

  const uniqueProvinsi = [...new Set(
    allPengurus
      .filter(item => item.level === activeTab && item.provinsi)
      .map(item => item.provinsi)
  )].sort()

  const uniqueKabupaten = [...new Set(
    allPengurus
      .filter(item => item.level === activeTab && item.kabupaten && (!selectedProvinsi || item.provinsi === selectedProvinsi))
      .map(item => item.kabupaten)
  )].sort()

  const filteredList = activeDataList.filter(item => {
    const matchesSearch = item.name.toLowerCase().includes(searchQuery.toLowerCase()) || 
                          (item.role || '').toLowerCase().includes(searchQuery.toLowerCase())
    
    const matchesProvinsi = !selectedProvinsi || item.provinsi === selectedProvinsi
    const matchesKabupaten = !selectedKabupaten || item.kabupaten === selectedKabupaten

    return matchesSearch && matchesProvinsi && matchesKabupaten
  }).sort((a, b) => (a.sort_order || 0) - (b.sort_order || 0))

  const totalPages = Math.max(3, Math.ceil(filteredList.length / itemsPerPage))
  const paginatedList = filteredList.slice((currentPage - 1) * itemsPerPage, currentPage * itemsPerPage)

  const handleTabChange = (tab) => {
    setSearchParams({ tab })
    setSearchQuery('')
    setSelectedProvinsi('')
    setSelectedKabupaten('')
  }

  return (
    <PublicLayout>
      {/* HEADER BANNER */}
      <section className="pt-36 pb-12 md:pt-40 md:pb-16 bg-white border-b border-slate-100">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <div className="inline-flex items-center justify-center gap-2 text-xs font-bold text-brand-600 bg-brand-50 px-4 py-1.5 rounded-full mb-4">
            <i className="ph-bold ph-users-three" /> Struktur Kepengurusan
          </div>
          <h1 className="font-heading text-3xl md:text-4xl font-extrabold text-slate-900 mb-4 tracking-tight">
            {activeTab === 'ketua' ? 'Ketua Umum DPP' :
             activeTab === 'dpp' ? 'Pengurus Pusat (DPP)' : 
             activeTab === 'dpd' ? 'Pengurus Daerah (DPD)' : 'Pengurus Cabang (DPC)'}
          </h1>
          <p className="text-base text-slate-500 max-w-2xl mx-auto font-medium">
            GRADASI telah hadir di 38 provinsi di Indonesia. Kenali individu-individu berdedikasi yang siap mendorong akselerasi digital nasional.
          </p>
        </div>
      </section>

      {/* SUB-MENU TABS */}
      <div className="bg-slate-50 border-b border-slate-200 py-3 sticky top-16 z-30 shadow-xs">
        <div className="max-w-6xl mx-auto px-4 flex flex-wrap gap-2 justify-center">
          <button 
            onClick={() => handleTabChange('ketua')}
            className={`px-5 py-2.5 text-sm font-bold rounded-xl transition ${activeTab === 'ketua' ? 'bg-brand-600 text-white shadow-md' : 'text-slate-600 hover:bg-slate-200/60'}`}
          >
            Ketua Umum
          </button>
          <button 
            onClick={() => handleTabChange('dpp')}
            className={`px-5 py-2.5 text-sm font-bold rounded-xl transition ${activeTab === 'dpp' ? 'bg-brand-600 text-white shadow-md' : 'text-slate-600 hover:bg-slate-200/60'}`}
          >
            Pengurus Pusat (DPP)
          </button>
          <button 
            onClick={() => handleTabChange('dpd')}
            className={`px-5 py-2.5 text-sm font-bold rounded-xl transition ${activeTab === 'dpd' ? 'bg-brand-600 text-white shadow-md' : 'text-slate-600 hover:bg-slate-200/60'}`}
          >
            Pengurus Daerah (DPD)
          </button>
          <button 
            onClick={() => handleTabChange('dpc')}
            className={`px-5 py-2.5 text-sm font-bold rounded-xl transition ${activeTab === 'dpc' ? 'bg-brand-600 text-white shadow-md' : 'text-slate-600 hover:bg-slate-200/60'}`}
          >
            Pengurus Cabang (DPC)
          </button>
        </div>
      </div>

      <section className="py-12 md:py-16 bg-slate-50 min-h-[60vh]">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">

          {/* Clean Search & Filter Bar */}
          {activeTab !== 'ketua' && (
            <div className="max-w-4xl mx-auto mb-10 bg-white p-3 sm:p-3.5 rounded-2xl border border-slate-200/60 shadow-[0_4px_25px_rgba(0,0,0,0.03)] flex flex-col md:flex-row gap-3 justify-between items-center w-full">
              <div className="relative w-full md:w-80">
                <i className="ph-bold ph-magnifying-glass absolute left-4 top-1/2 -translate-y-1/2 text-slate-400 text-base" />
                <input
                  type="text"
                  placeholder="Cari nama atau jabatan pengurus..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="w-full pl-11 pr-4 py-2.5 bg-slate-50 border border-slate-200/80 rounded-xl text-xs sm:text-sm text-slate-700 focus:bg-white focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 transition-all outline-none"
                />
              </div>

              {(activeTab === 'dpd' || activeTab === 'dpc') && (
                <div className="flex flex-col sm:flex-row gap-3 w-full md:w-auto">
                  <select
                    value={selectedProvinsi}
                    onChange={(e) => { setSelectedProvinsi(e.target.value); setSelectedKabupaten(''); }}
                    className="px-4 py-2.5 bg-slate-50 border border-slate-200/80 rounded-xl text-xs sm:text-sm text-slate-600 outline-none cursor-pointer hover:bg-slate-100/60 transition"
                  >
                    <option value="">Semua Provinsi</option>
                    {uniqueProvinsi.map(prov => (
                      <option key={prov} value={prov}>{prov}</option>
                    ))}
                  </select>

                  {activeTab === 'dpc' && (
                    <select
                      value={selectedKabupaten}
                      onChange={(e) => setSelectedKabupaten(e.target.value)}
                      className="px-4 py-2.5 bg-slate-50 border border-slate-200/80 rounded-xl text-xs sm:text-sm text-slate-600 outline-none cursor-pointer hover:bg-slate-100/60 transition"
                    >
                      <option value="">Semua Kab/Kota</option>
                      {uniqueKabupaten.map(kab => (
                        <option key={kab} value={kab}>{kab}</option>
                      ))}
                    </select>
                  )}
                </div>
              )}
            </div>
          )}

          {/* KETUA UMUM SPECIAL HERO CARD */}
          {activeTab === 'ketua' && (
            <div className="flex justify-center my-6">
              {allPengurus.filter(item => item.level === 'ketua').map(item => (
                <div key={item.id} className="bg-white rounded-3xl p-8 md:p-12 border border-slate-200/80 shadow-[0_15px_40px_rgba(0,0,0,0.06)] flex flex-col items-center text-center max-w-md w-full">
                  <div className="w-40 h-40 rounded-full overflow-hidden mb-6 border-4 border-brand-100 p-1 shadow-md bg-slate-50">
                    <img src={item.image_url.startsWith('http') ? item.image_url : `http://127.0.0.1:8080${item.image_url}`} alt={item.name} className="w-full h-full object-cover rounded-full" />
                  </div>
                  <h3 className="font-heading font-extrabold text-slate-900 text-2xl md:text-3xl mb-1 tracking-tight">{item.name}</h3>
                  <p className="text-sm font-bold text-brand-600 tracking-wider uppercase mb-1">{item.role}</p>
                  <p className="text-xs text-slate-400 font-semibold mb-6">Masa Bakti {item.periode || '2024 - 2029'}</p>
                  
                  <div className="flex justify-center gap-4 pt-6 border-t border-slate-100 w-full">
                    {item.facebook_url && (
                      <a href={item.facebook_url} target="_blank" rel="noreferrer" className="w-10 h-10 rounded-xl bg-slate-50 border border-slate-200 text-slate-500 hover:bg-brand-600 hover:text-white flex items-center justify-center transition shadow-xs">
                        <i className="ph-fill ph-facebook-logo text-lg" />
                      </a>
                    )}
                    {item.instagram_url && (
                      <a href={item.instagram_url} target="_blank" rel="noreferrer" className="w-10 h-10 rounded-xl bg-slate-50 border border-slate-200 text-slate-500 hover:bg-brand-600 hover:text-white flex items-center justify-center transition shadow-xs">
                        <i className="ph-fill ph-instagram-logo text-lg" />
                      </a>
                    )}
                  </div>
                </div>
              ))}
            </div>
          )}

          {/* GRID CARD FOR DPP, DPD, DPC */}
          {activeTab !== 'ketua' && (
            <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
              {paginatedList.map(item => (
                <div key={item.id} className="card-lift bg-white rounded-2xl p-6 border border-slate-100 flex flex-col items-center text-center hover:border-brand-200 hover:shadow-[0_8px_30px_rgba(37,99,235,0.08)] transition-all duration-300 group">
                  <div className="w-28 h-28 rounded-full overflow-hidden mb-5 border-4 border-slate-50 group-hover:border-brand-100 transition duration-300 shadow-sm">
                    <img src={item.image_url.startsWith('http') ? item.image_url : `http://127.0.0.1:8080${item.image_url}`} alt={item.name} className="w-full h-full object-cover group-hover:scale-110 transition duration-500" />
                  </div>
                  <div className="flex-grow flex flex-col w-full">
                    <h4 className="font-heading font-bold text-slate-900 text-base mb-1 group-hover:text-brand-600 transition">{item.name}</h4>
                    <p className="text-xs font-semibold text-brand-600 mb-2">{item.role}</p>
                    {item.provinsi && <p className="text-[11px] text-slate-400 font-medium mb-4">{item.provinsi} {item.kabupaten ? `• ${item.kabupaten}` : ''}</p>}
                    
                    <div className="flex justify-center gap-3 mt-auto pt-4 border-t border-slate-100">
                      <a href="#" className="w-8 h-8 rounded-lg bg-slate-50 text-slate-400 hover:text-brand-600 flex items-center justify-center transition">
                        <i className="ph-fill ph-facebook-logo" />
                      </a>
                      <a href="#" className="w-8 h-8 rounded-lg bg-slate-50 text-slate-400 hover:text-brand-600 flex items-center justify-center transition">
                        <i className="ph-fill ph-instagram-logo" />
                      </a>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}

          {/* ALWAYS VISIBLE PAGINATION FOR DPP, DPD, DPC */}
          {activeTab !== 'ketua' && (
            <div className="mt-16 flex justify-center">
              <nav className="flex items-center gap-2">
                <button 
                  disabled={currentPage <= 1}
                  onClick={() => setCurrentPage(prev => Math.max(prev - 1, 1))} 
                  className="w-10 h-10 flex items-center justify-center rounded-xl border border-slate-200 text-slate-400 hover:border-brand-500 hover:text-brand-600 hover:bg-brand-50 transition disabled:opacity-40"
                >
                  <i className="ph-bold ph-caret-left" />
                </button>
                
                {[1, 2, 3].map(page => (
                  <button 
                    key={page} 
                    onClick={() => setCurrentPage(page)} 
                    className={`w-10 h-10 flex items-center justify-center rounded-xl font-bold transition ${currentPage === page ? 'bg-brand-700 text-white shadow-md shadow-brand-500/20' : 'border border-slate-200 text-slate-600 hover:border-brand-500 hover:text-brand-600 hover:bg-brand-50'}`}
                  >
                    {page}
                  </button>
                ))}

                <button 
                  disabled={currentPage >= totalPages}
                  onClick={() => setCurrentPage(prev => Math.min(prev + 1, totalPages))} 
                  className="w-10 h-10 flex items-center justify-center rounded-xl border border-slate-200 text-slate-400 hover:border-brand-500 hover:text-brand-600 hover:bg-brand-50 transition disabled:opacity-40"
                >
                  <i className="ph-bold ph-caret-right" />
                </button>
              </nav>
            </div>
          )}

        </div>
      </section>
    </PublicLayout>
  )
}
