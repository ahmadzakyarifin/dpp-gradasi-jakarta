import { useState, useEffect } from 'react'
import { useSearchParams } from 'react-router-dom'
import PublicLayout from '../layouts/PublicLayout'
import { pengurusService } from '../services/pengurusService'

export default function Kepengurusan() {
  const [searchParams, setSearchParams] = useSearchParams()
  const activeTab = searchParams.get('tab') || 'ketua'
  
  const [allPengurus, setAllPengurus] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedProvinsi, setSelectedProvinsi] = useState('')
  const [selectedKabupaten, setSelectedKabupaten] = useState('')
  const [currentPage, setCurrentPage] = useState(1)
  const itemsPerPage = 8

  useEffect(() => {
    setLoading(true)
    pengurusService.list()
      .then(res => {
        if (res.success && res.data) {
          setAllPengurus(res.data)
        } else {
          setError('Gagal memuat data pengurus')
        }
      })
      .catch(() => setError('Terjadi kesalahan koneksi'))
      .finally(() => setLoading(false))
  }, [])

  useEffect(() => {
    setCurrentPage(1)
  }, [activeTab, searchQuery, selectedProvinsi, selectedKabupaten])

  // Filter lists based on selected tab level
  const activeDataList = allPengurus.filter(item => item.level === activeTab && item.is_active)

  // Get unique lists for filters
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

  // Filter items by search query and dropdown selections
  const filteredList = activeDataList.filter(item => {
    const matchesSearch = item.name.toLowerCase().includes(searchQuery.toLowerCase()) || 
                          item.role.toLowerCase().includes(searchQuery.toLowerCase())
    
    const matchesProvinsi = !selectedProvinsi || item.provinsi === selectedProvinsi
    const matchesKabupaten = !selectedKabupaten || item.kabupaten === selectedKabupaten

    return matchesSearch && matchesProvinsi && matchesKabupaten
  }).sort((a, b) => a.sort_order - b.sort_order)

  // Pagination calculations
  const totalPages = Math.ceil(filteredList.length / itemsPerPage) || 1
  const paginatedList = filteredList.slice((currentPage - 1) * itemsPerPage, currentPage * itemsPerPage)

  const handleTabChange = (tab) => {
    setSearchParams({ tab })
    setSearchQuery('')
    setSelectedProvinsi('')
    setSelectedKabupaten('')
  }

  return (
    <PublicLayout>
      {/* Navigation submenu inside header equivalent */}
      <div className="bg-slate-50 border-b border-slate-100 py-3 sticky top-16 z-40">
        <div className="max-w-6xl mx-auto px-4 flex flex-wrap gap-2 justify-center">
          <button 
            onClick={() => handleTabChange('ketua')}
            className={`px-4 py-2 text-sm font-semibold rounded-lg transition ${activeTab === 'ketua' ? 'bg-brand-600 text-white shadow-sm' : 'text-slate-600 hover:bg-slate-100'}`}
          >
            Ketua Umum
          </button>
          <button 
            onClick={() => handleTabChange('dpp')}
            className={`px-4 py-2 text-sm font-semibold rounded-lg transition ${activeTab === 'dpp' ? 'bg-brand-600 text-white shadow-sm' : 'text-slate-600 hover:bg-slate-100'}`}
          >
            Pengurus Pusat (DPP)
          </button>
          <button 
            onClick={() => handleTabChange('dpd')}
            className={`px-4 py-2 text-sm font-semibold rounded-lg transition ${activeTab === 'dpd' ? 'bg-brand-600 text-white shadow-sm' : 'text-slate-600 hover:bg-slate-100'}`}
          >
            Pengurus Daerah (DPD)
          </button>
          <button 
            onClick={() => handleTabChange('dpc')}
            className={`px-4 py-2 text-sm font-semibold rounded-lg transition ${activeTab === 'dpc' ? 'bg-brand-600 text-white shadow-sm' : 'text-slate-600 hover:bg-slate-100'}`}
          >
            Pengurus Cabang (DPC)
          </button>
        </div>
      </div>

      <section className="pt-20 pb-20 bg-white min-h-[80vh]">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          
          <div className="mb-10">
            <h2 className="font-heading text-3xl md:text-4xl font-extrabold text-slate-900 mb-4 tracking-tight">
              {activeTab === 'ketua' ? 'Ketua Umum' :
               activeTab === 'dpp' ? 'Kepengurusan Bagian Pusat' : 
               activeTab === 'dpd' ? 'Kepengurusan Bagian Provinsi' : 'Kepengurusan Bagian Kab/Kota'}
            </h2>
            <p className="text-slate-500 text-sm max-w-3xl mx-auto leading-relaxed">
              Kami sudah hadir di 38 provinsi di Indonesia. Kenali individu-individu yang berdedikasi dan bersemangat, 
              siap untuk bekerja sama dan mendorong inovasi ke puncak pencapaian.
            </p>
          </div>

          {loading && <div className="py-16 text-slate-500">Memuat data pengurus...</div>}
          {error && <div className="py-16 text-red-600 font-medium">{error}</div>}

          {!loading && !error && (
            <>
              {/* Search & Filter Bar */}
              {activeTab !== 'ketua' && (
                <div className="max-w-4xl mx-auto mb-12 flex flex-col md:flex-row gap-4">
                  <div className="relative flex-grow">
                    <input 
                      type="text" 
                      placeholder="Cari nama atau jabatan..." 
                      value={searchQuery}
                      onChange={(e) => setSearchQuery(e.target.value)}
                      className="w-full pl-4 pr-10 py-3 bg-slate-50 border border-slate-200 rounded-xl text-sm focus:bg-white focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 transition-all outline-none text-slate-700"
                    />
                    {searchQuery && (
                      <button onClick={() => setSearchQuery('')} className="absolute inset-y-0 right-0 pr-4 flex items-center text-slate-400">
                        ✕
                      </button>
                    )}
                  </div>

                  {(activeTab === 'dpd' || activeTab === 'dpc') && (
                    <div className="flex flex-col sm:flex-row gap-4">
                      <select 
                        value={selectedProvinsi} 
                        onChange={(e) => { setSelectedProvinsi(e.target.value); setSelectedKabupaten(''); }}
                        className="px-4 py-3 bg-slate-50 border border-slate-200 rounded-xl text-sm text-slate-600 outline-none min-w-[160px]"
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
                          className="px-4 py-3 bg-slate-50 border border-slate-200 rounded-xl text-sm text-slate-600 outline-none min-w-[160px]"
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

              {/* Ketua Umum Special Card */}
              {activeTab === 'ketua' && paginatedList.length > 0 && (
                <div className="flex justify-center my-6">
                  {paginatedList.map(item => (
                    <div key={item.id} className="bg-white rounded-2xl p-8 md:p-10 border border-slate-200/60 shadow-[0_4px_25px_rgba(0,0,0,0.03)] flex flex-col items-center text-center max-w-sm w-full">
                      <div className="w-32 h-32 md:w-36 md:h-36 rounded-full overflow-hidden mb-6 border-2 border-slate-100 p-1 shadow-sm bg-slate-50">
                        <img src={item.image_url.startsWith('http') ? item.image_url : `http://127.0.0.1:8080${item.image_url}`} alt={item.name} className="w-full h-full object-cover rounded-full" />
                      </div>
                      <h3 class="font-heading font-bold text-slate-900 text-2xl mb-1 tracking-tight">{item.name}</h3>
                      <p className="text-xs font-semibold text-brand-600 tracking-wider uppercase mb-1">{item.role}</p>
                      <p className="text-xs text-slate-400 font-medium mb-6">Periode {item.periode}</p>
                      
                      <div className="flex justify-center gap-3 pt-5 border-t border-slate-100 w-full">
                        {item.facebook_url && <a href={item.facebook_url} target="_blank" className="w-9 h-9 rounded-lg bg-slate-50 text-slate-400 hover:text-brand-600 flex items-center justify-center transition">FB</a>}
                        {item.instagram_url && <a href={item.instagram_url} target="_blank" className="w-9 h-9 rounded-lg bg-slate-50 text-slate-400 hover:text-brand-600 flex items-center justify-center transition">IG</a>}
                        {item.linkedin_url && <a href={item.linkedin_url} target="_blank" className="w-9 h-9 rounded-lg bg-slate-50 text-slate-400 hover:text-brand-600 flex items-center justify-center transition">IN</a>}
                      </div>
                    </div>
                  ))}
                </div>
              )}

              {/* Regular Card Grid for DPP, DPD, DPC */}
              {activeTab !== 'ketua' && paginatedList.length > 0 && (
                <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
                  {paginatedList.map(item => (
                    <div key={item.id} className="bg-white rounded-2xl p-6 border border-slate-100 flex flex-col items-center text-center hover:border-brand-200 hover:shadow-[0_8px_30px_rgba(37,99,235,0.08)] transition-all duration-300 group">
                      <div className="w-24 h-24 rounded-full overflow-hidden mb-5 border-4 border-slate-50 group-hover:border-brand-100 transition duration-300 shadow-sm">
                        <img src={item.image_url.startsWith('http') ? item.image_url : `http://127.0.0.1:8080${item.image_url}`} alt={item.name} className="w-full h-full object-cover group-hover:scale-110 transition duration-500" />
                      </div>
                      <div className="flex-grow flex flex-col w-full">
                        <h4 className="font-heading font-bold text-slate-900 text-[15px] mb-1 group-hover:text-brand-600 transition">{item.name}</h4>
                        <p className="text-xs font-medium text-slate-500 mb-5 flex-grow">{item.role}</p>
                        
                        <div className="flex justify-center gap-3 mt-auto pt-4 border-t border-slate-50">
                          {item.facebook_url && <a href={item.facebook_url} target="_blank" className="text-slate-400 hover:text-brand-600 transition">FB</a>}
                          {item.instagram_url && <a href={item.instagram_url} target="_blank" className="text-slate-400 hover:text-brand-600 transition">IG</a>}
                          {item.linkedin_url && <a href={item.linkedin_url} target="_blank" className="text-slate-400 hover:text-brand-600 transition">LN</a>}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              )}

              {/* Empty state */}
              {activeTab !== 'ketua' && filteredList.length === 0 && (
                <div className="py-20 flex flex-col items-center justify-center text-slate-400 border border-dashed border-slate-200 rounded-2xl">
                  <p className="font-medium text-lg text-slate-600">Pengurus tidak ditemukan</p>
                </div>
              )}

              {/* Pagination */}
              {activeTab !== 'ketua' && totalPages > 1 && (
                <div className="mt-12 flex items-center justify-center gap-4 border-t border-slate-100 pt-8">
                  <button 
                    onClick={() => setCurrentPage(prev => Math.max(prev - 1, 1))} 
                    disabled={currentPage === 1}
                    className={`px-4 py-2 text-sm font-medium rounded-md transition ${currentPage === 1 ? 'text-slate-300 cursor-not-allowed' : 'text-slate-600 hover:bg-slate-50 hover:text-brand-600'}`}
                  >
                    Sebelumnya
                  </button>
                  <div className="flex items-center gap-2">
                    {Array.from({ length: totalPages }, (_, i) => i + 1).map(page => (
                      <button 
                        key={page} 
                        onClick={() => setCurrentPage(page)} 
                        className={`w-8 h-8 flex items-center justify-center rounded-md text-sm font-semibold transition ${currentPage === page ? 'bg-brand-600 text-white shadow-sm' : 'text-slate-500 hover:bg-slate-50 hover:text-brand-600'}`}
                      >
                        {page}
                      </button>
                    ))}
                  </div>
                  <button 
                    onClick={() => setCurrentPage(prev => Math.min(prev + 1, totalPages))} 
                    disabled={currentPage === totalPages}
                    className={`px-4 py-2 text-sm font-medium rounded-md transition ${currentPage === totalPages ? 'text-slate-300 cursor-not-allowed' : 'text-slate-600 hover:bg-slate-50 hover:text-brand-600'}`}
                  >
                    Selanjutnya
                  </button>
                </div>
              )}
            </>
          )}

        </div>
      </section>
    </PublicLayout>
  )
}
