import { useState, useEffect, useMemo } from 'react'
import { useSearchParams, Link } from 'react-router-dom'
import PublicLayout from '../layouts/PublicLayout'
import { pengurusService } from '../services/pengurusService'
import { resolveAssetUrl } from '../utils/assetUrl'
import { getProvinces, getRegencies } from 'kode-wilayah-id'

export default function Kepengurusan() {
  const [searchParams, setSearchParams] = useSearchParams()
  const activeTab = searchParams.get('tab') || 'Ketua Umum'
  
  const [allPengurus, setAllPengurus] = useState([])
  const [searchInput, setSearchInput] = useState('')
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedProvinsi, setSelectedProvinsi] = useState('')
  const [selectedKabupaten, setSelectedKabupaten] = useState('')
  const [currentPage, setCurrentPage] = useState(1)
  const itemsPerPage = 24
  
  // Wilayah Data
  const [indonesiaProvinces, setIndonesiaProvinces] = useState([])
  const [indonesiaKabupaten, setIndonesiaKabupaten] = useState([])

  const provOptions = useMemo(() => {
    return [{ value: '', label: 'Semua Provinsi' }, ...indonesiaProvinces.map(p => ({ value: p, label: p }))]
  }, [indonesiaProvinces])

  const kabOptions = useMemo(() => {
    return [{ value: '', label: 'Semua Kab/Kota' }, ...indonesiaKabupaten.map(k => ({ value: k, label: k }))]
  }, [indonesiaKabupaten])

  // Load Regions
  useEffect(() => {
    try {
      const provs = getProvinces()
      setIndonesiaProvinces(provs.map(p => p.name).sort())
      const kabs = getRegencies()
      // Format to Title Case like "Kabupaten Bogor"
      const formattedKabs = kabs.map(k => k.name.toLowerCase().replace(/\b\w/g, l => l.toUpperCase())).sort()
      setIndonesiaKabupaten(formattedKabs)
    } catch (err) {
      console.error("Failed to load region data", err)
    }
  }, [])

  useEffect(() => {
    pengurusService.list({ limit: 100 })
      .then(res => {
        if (res.success && res.data) {
          const list = Array.isArray(res.data) ? res.data : (res.data.data || res.data.pengurus || [])
          // Normalisasi: API kirim image_path → FE pakai image_url di JSX
          setAllPengurus(list.map(p => ({ ...p, image_url: p.image_path || p.image_url })))
        }
      })
      .catch(() => {})
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

  const totalPages = Math.ceil(filteredList.length / itemsPerPage) || 1
  const paginatedList = filteredList.slice((currentPage - 1) * itemsPerPage, currentPage * itemsPerPage)
  const pageNumbers = Array.from({ length: totalPages }, (_, i) => i + 1)
  const showArrows = totalPages > 1 && filteredList.length >= 3

  const handleTabChange = (tab) => {
    setSearchParams({ tab })
    setSearchInput('')
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
            {activeTab === 'Ketua Umum' ? 'Ketua Umum' :
             activeTab === 'Pengurus Pusat' ? 'Pengurus Pusat' : 
             activeTab === 'Pengurus Provinsi' ? 'Pengurus Provinsi' : 'Pengurus Kab/Kota'}
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
            onClick={() => handleTabChange('Ketua Umum')}
            className={`px-5 py-2.5 text-sm font-bold rounded-xl transition ${activeTab === 'Ketua Umum' ? 'bg-brand-600 text-white shadow-md' : 'text-slate-600 hover:bg-slate-200/60'}`}
          >
            Ketua Umum
          </button>
          <button 
            onClick={() => handleTabChange('Pengurus Pusat')}
            className={`px-5 py-2.5 text-sm font-bold rounded-xl transition ${activeTab === 'Pengurus Pusat' ? 'bg-brand-600 text-white shadow-md' : 'text-slate-600 hover:bg-slate-200/60'}`}
          >
            Pengurus Pusat
          </button>
          <button 
            onClick={() => handleTabChange('Pengurus Provinsi')}
            className={`px-5 py-2.5 text-sm font-bold rounded-xl transition ${activeTab === 'Pengurus Provinsi' ? 'bg-brand-600 text-white shadow-md' : 'text-slate-600 hover:bg-slate-200/60'}`}
          >
            Pengurus Provinsi
          </button>
          <button 
            onClick={() => handleTabChange('Pengurus Kab/Kota')}
            className={`px-5 py-2.5 text-sm font-bold rounded-xl transition ${activeTab === 'Pengurus Kab/Kota' ? 'bg-brand-600 text-white shadow-md' : 'text-slate-600 hover:bg-slate-200/60'}`}
          >
            Pengurus Kab/Kota
          </button>
        </div>
      </div>

      <section className="py-12 md:py-16 bg-slate-50 min-h-[60vh]">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">

          {/* Clean Search & Filter Bar */}
          {activeTab !== 'Ketua Umum' && (
            <div className={`mx-auto mb-10 bg-white p-2 sm:p-2.5 rounded-2xl border border-slate-200/60 shadow-[0_4px_25px_rgba(0,0,0,0.03)] flex flex-col md:flex-row gap-3 items-center w-full ${activeTab === 'Pengurus Pusat' ? 'max-w-3xl' : 'max-w-5xl'}`}>
              
              {/* Inputs Group */}
              <div className="flex flex-col md:flex-row gap-2 w-full">
                {/* Text Search */}
                <div className="relative w-full">
                  <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                    <i className="ph-bold ph-magnifying-glass text-slate-400 text-lg" />
                  </div>
                  <input
                    type="text"
                    placeholder="Cari nama atau jabatan pengurus..."
                    value={searchInput}
                    onChange={(e) => setSearchInput(e.target.value)}
                    onKeyDown={(e) => { if (e.key === 'Enter') setSearchQuery(searchInput) }}
                    className={`w-full pl-10 pr-4 py-2.5 ${activeTab === 'Pengurus Pusat' ? 'bg-transparent shadow-none border-none' : 'bg-slate-50 border border-slate-200/80'} focus:bg-white focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 rounded-xl text-xs sm:text-sm text-slate-700 transition-all outline-none`}
                  />
                </div>

                {/* Dropdowns */}
                {activeTab === 'Pengurus Provinsi' && (
                  <div className="w-full md:w-auto flex-shrink-0 z-50">
                    <select
                      value={selectedProvinsi}
                      onChange={(e) => { setSelectedProvinsi(e.target.value); setSelectedKabupaten(''); }}
                      className="w-full md:w-auto bg-slate-50 border border-slate-200/80 focus:bg-white focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 rounded-xl text-xs sm:text-sm text-slate-700 transition-all outline-none py-2.5 px-4 pr-10 cursor-pointer appearance-none bg-[url('data:image/svg+xml;charset=US-ASCII,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20width%3D%22292.4%22%20height%3D%22292.4%22%3E%3Cpath%20fill%3D%22%2394a3b8%22%20d%3D%22M287%2069.4a17.6%2017.6%200%200%200-13-5.4H18.4c-5%200-9.3%201.8-12.9%205.4A17.6%2017.6%200%200%200%200%2082.2c0%205%201.8%209.3%205.4%2012.9l128%20127.9c3.6%203.6%207.8%205.4%2012.8%205.4s9.2-1.8%2012.8-5.4L287%2095c3.5-3.5%205.4-7.8%205.4-12.8%200-5-1.9-9.2-5.5-12.8z%22%2F%3E%3C%2Fsvg%3E')] bg-[length:0.7rem_auto] bg-no-repeat bg-[position:right_1rem_center]"
                    >
                      {provOptions.map(o => (
                        <option key={o.value} value={o.value}>{o.label}</option>
                      ))}
                    </select>
                  </div>
                )}

                {activeTab === 'Pengurus Kab/Kota' && (
                  <div className="w-full md:w-auto flex-shrink-0 z-50">
                    <select
                      value={selectedKabupaten}
                      onChange={(e) => setSelectedKabupaten(e.target.value)}
                      className="w-full md:w-auto bg-slate-50 border border-slate-200/80 focus:bg-white focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 rounded-xl text-xs sm:text-sm text-slate-700 transition-all outline-none py-2.5 px-4 pr-10 cursor-pointer appearance-none bg-[url('data:image/svg+xml;charset=US-ASCII,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20width%3D%22292.4%22%20height%3D%22292.4%22%3E%3Cpath%20fill%3D%22%2394a3b8%22%20d%3D%22M287%2069.4a17.6%2017.6%200%200%200-13-5.4H18.4c-5%200-9.3%201.8-12.9%205.4A17.6%2017.6%200%200%200%200%2082.2c0%205%201.8%209.3%205.4%2012.9l128%20127.9c3.6%203.6%207.8%205.4%2012.8%205.4s9.2-1.8%2012.8-5.4L287%2095c3.5-3.5%205.4-7.8%205.4-12.8%200-5-1.9-9.2-5.5-12.8z%22%2F%3E%3C%2Fsvg%3E')] bg-[length:0.7rem_auto] bg-no-repeat bg-[position:right_1rem_center]"
                    >
                      {kabOptions.map(o => (
                        <option key={o.value} value={o.value}>{o.label}</option>
                      ))}
                    </select>
                  </div>
                )}
              </div>

              {/* Action Buttons */}
              <div className="flex gap-2 w-full md:w-auto shrink-0">
                <button 
                  onClick={() => setSearchQuery(searchInput)} 
                  className="w-full md:w-auto bg-brand-700 hover:bg-brand-800 text-white px-5 sm:px-6 py-2.5 rounded-xl font-bold text-sm transition-colors shadow-sm whitespace-nowrap"
                >
                  Cari
                </button>
                {(searchQuery || selectedProvinsi || selectedKabupaten) && (
                  <button 
                    onClick={() => { 
                      setSearchInput(''); 
                      setSearchQuery(''); 
                      setSelectedProvinsi(''); 
                      setSelectedKabupaten(''); 
                    }} 
                    className="w-full md:w-auto bg-slate-100 hover:bg-slate-200 text-slate-700 px-5 sm:px-6 py-2.5 rounded-xl font-bold text-sm transition-colors shadow-sm whitespace-nowrap"
                  >
                    Reset
                  </button>
                )}
              </div>
            </div>
          )}

          {/* KETUA UMUM SPECIAL HERO CARD */}
          {activeTab === 'Ketua Umum' && (
            <div className="flex justify-center my-6">
              {allPengurus.filter(item => item.level === 'Ketua Umum').length === 0 ? (
                <div className="bg-white rounded-3xl p-8 md:p-12 border border-slate-200/80 shadow-[0_15px_40px_rgba(0,0,0,0.03)] flex flex-col items-center text-center max-w-md w-full">
                  <i className="ph-bold ph-user-focus text-4xl text-slate-300 mb-3 block" />
                  <h3 className="font-heading font-bold text-slate-800 text-lg mb-1">Ketua Umum belum diatur</h3>
                  <p className="text-xs text-slate-400 font-semibold">Data Ketua Umum belum dimasukkan ke database.</p>
                </div>
              ) : (
                allPengurus.filter(item => item.level === 'Ketua Umum').map(item => (
                  <div key={item.id} className="bg-white rounded-3xl p-8 md:p-12 border border-slate-200/80 shadow-[0_15px_40px_rgba(0,0,0,0.06)] flex flex-col items-center text-center max-w-md w-full">
                    <div className="w-40 h-40 rounded-full overflow-hidden mb-6 border-4 border-brand-100 p-1 shadow-md bg-slate-50">
                      <img src={resolveAssetUrl(item.image_url || item.image_path)} alt={item.name} className="w-full h-full object-cover rounded-full" />
                    </div>
                    <h3 className="font-heading font-extrabold text-slate-900 text-2xl md:text-3xl mb-1 tracking-tight hover:text-brand-600 transition">
                      <Link to={`/kepengurusan/profile/${item.id}`}>{item.name}</Link>
                    </h3>
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
                      {item.linkedin_url && (
                        <a href={item.linkedin_url} target="_blank" rel="noreferrer" className="w-10 h-10 rounded-xl bg-slate-50 border border-slate-200 text-slate-500 hover:bg-brand-600 hover:text-white flex items-center justify-center transition shadow-xs">
                          <i className="ph-fill ph-linkedin-logo text-lg" />
                        </a>
                      )}
                      {item.twitter_url && (
                        <a href={item.twitter_url} target="_blank" rel="noreferrer" className="w-10 h-10 rounded-xl bg-slate-50 border border-slate-200 text-slate-500 hover:bg-brand-600 hover:text-white flex items-center justify-center transition shadow-xs">
                          <i className="ph-fill ph-x-logo text-lg" />
                        </a>
                      )}

                      {item.whatsapp && (
                        <a href={`https://wa.me/${item.whatsapp.replace(/\D/g, '')}`} target="_blank" rel="noreferrer" className="w-10 h-10 rounded-xl bg-slate-50 border border-slate-200 text-slate-500 hover:bg-brand-600 hover:text-white flex items-center justify-center transition shadow-xs">
                          <i className="ph-fill ph-whatsapp-logo text-lg" />
                        </a>
                      )}
                    </div>
                  </div>
                ))
              )}
            </div>
          )}

          {/* GRID CARD FOR DPP, DPD, DPC */}
          {activeTab !== 'Ketua Umum' && (
            filteredList.length === 0 ? (
              <div className="py-16 text-center text-slate-500 bg-white rounded-2xl border border-slate-200/60 p-8 w-full max-w-lg mx-auto">
                <i className="ph-bold ph-users text-4xl text-slate-300 mb-2 block" />
                <p className="font-semibold text-slate-700">Tidak ada pengurus ditemukan</p>
                <p className="text-xs text-slate-400 mt-1">Data pengurus untuk kategori ini belum tersedia atau tidak cocok dengan filter pencarian.</p>
              </div>
            ) : (
              <div className="flex flex-wrap justify-center gap-6">
                {paginatedList.map(item => (
                  <div key={item.id} className="w-full sm:w-[calc(50%-12px)] md:w-[calc(33.333%-16px)] lg:w-[calc(25%-18px)] max-w-sm card-lift bg-white rounded-2xl p-6 border border-slate-100 flex flex-col items-center text-center hover:border-brand-200 hover:shadow-[0_8px_30px_rgba(37,99,235,0.08)] transition-all duration-300 group">
                    <div className="w-28 h-28 rounded-full overflow-hidden mb-5 border-4 border-slate-50 group-hover:border-brand-100 transition duration-300 shadow-sm">
                      <img src={resolveAssetUrl(item.image_url || item.image_path)} alt={item.name} className="w-full h-full object-cover group-hover:scale-110 transition duration-500" />
                    </div>
                    <div className="flex-grow flex flex-col w-full">
                      <h4 className="font-heading font-bold text-slate-900 text-base mb-1 group-hover:text-brand-600 transition">
                        <Link to={`/kepengurusan/profile/${item.id}`}>{item.name}</Link>
                      </h4>
                      <p className="text-xs font-semibold text-brand-600 mb-2">{item.role}</p>
                      {item.provinsi && <p className="text-[11px] text-slate-400 font-medium mb-4">{item.provinsi} {item.kabupaten ? `• ${item.kabupaten}` : ''}</p>}
                      
                      <div className="flex justify-center gap-3 mt-auto pt-4 border-t border-slate-100">

                        {item.facebook_url && (
                          <a href={item.facebook_url} target="_blank" rel="noreferrer" className="w-8 h-8 rounded-lg bg-slate-50 text-slate-400 hover:text-brand-600 flex items-center justify-center transition">
                            <i className="ph-fill ph-facebook-logo" />
                          </a>
                        )}
                        {item.instagram_url && (
                          <a href={item.instagram_url} target="_blank" rel="noreferrer" className="w-8 h-8 rounded-lg bg-slate-50 text-slate-400 hover:text-brand-600 flex items-center justify-center transition">
                            <i className="ph-fill ph-instagram-logo" />
                          </a>
                        )}
                        {item.linkedin_url && (
                          <a href={item.linkedin_url} target="_blank" rel="noreferrer" className="w-8 h-8 rounded-lg bg-slate-50 text-slate-400 hover:text-brand-600 flex items-center justify-center transition">
                            <i className="ph-fill ph-linkedin-logo" />
                          </a>
                        )}
                        {item.twitter_url && (
                          <a href={item.twitter_url} target="_blank" rel="noreferrer" className="w-8 h-8 rounded-lg bg-slate-50 text-slate-400 hover:text-brand-600 flex items-center justify-center transition">
                            <i className="ph-fill ph-x-logo" />
                          </a>
                        )}

                        {item.whatsapp && (
                          <a href={`https://wa.me/${item.whatsapp.replace(/\D/g, '')}`} target="_blank" rel="noreferrer" className="w-8 h-8 rounded-lg bg-slate-50 text-slate-400 hover:text-brand-600 flex items-center justify-center transition">
                            <i className="ph-fill ph-whatsapp-logo" />
                          </a>
                        )}
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )
          )}

          {/* ALWAYS VISIBLE PAGINATION FOR DPP, DPD, DPC */}
          {activeTab !== 'Ketua Umum' && filteredList.length > 0 && (
            <div className="mt-16 flex justify-center">
              <nav className="flex items-center gap-2">
                {showArrows && (
                  <button 
                    disabled={currentPage <= 1}
                    onClick={() => setCurrentPage(prev => Math.max(prev - 1, 1))} 
                    className="w-10 h-10 flex items-center justify-center rounded-xl border border-slate-200 text-slate-400 hover:border-brand-500 hover:text-brand-600 hover:bg-brand-50 transition disabled:opacity-40"
                  >
                    <i className="ph-bold ph-caret-left" />
                  </button>
                )}
                
                {pageNumbers.map(page => (
                  <button 
                    key={page} 
                    onClick={() => setCurrentPage(page)} 
                    className={`w-10 h-10 flex items-center justify-center rounded-xl font-bold transition ${currentPage === page ? 'bg-brand-700 text-white shadow-md shadow-brand-500/20' : 'border border-slate-200 text-slate-600 hover:border-brand-500 hover:text-brand-600 hover:bg-brand-50'}`}
                  >
                    {page}
                  </button>
                ))}

                {showArrows && (
                  <button 
                    disabled={currentPage >= totalPages}
                    onClick={() => setCurrentPage(prev => Math.min(prev + 1, totalPages))} 
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
