import { useState, useEffect } from 'react'
import AdminLayout from '../../layouts/AdminLayout'
import { pengurusService } from '../../services/pengurusService'

export default function PengurusAdmin() {
  const [items, setItems] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  
  const [currentTab, setCurrentTab] = useState('active') // active, trash (Trash handled locally if backend doesn't support soft delete for pengurus, or backend-driven)
  const [currentPage, setCurrentPage] = useState(1)
  const pageSize = 5

  const [searchQuery, setSearchQuery] = useState('')
  const [filterLevel, setFilterLevel] = useState('') // ketua, dpp, dpd, dpc
  const [filterSort, setFilterSort] = useState('sort_order')

  const [isFormOpen, setIsFormOpen] = useState(false)
  const [formMode, setFormMode] = useState('create')
  
  const [formData, setFormData] = useState({
    id: null,
    name: '',
    level: 'dpp',
    role: '',
    department: '',
    periode: '2025 - 2030',
    provinsi: '',
    kabupaten: '',
    sort_order: 1,
    image_url: '',
    linkedin_url: '',
    facebook_url: '',
    instagram_url: '',
    whatsapp: '',
    is_active: true
  })

  const loadPengurus = () => {
    setLoading(true)
    pengurusService.list()
      .then(res => {
        if (res.success && res.data) {
          setItems(res.data)
        } else {
          setError('Gagal memuat pengurus')
        }
      })
      .catch(() => setError('Kesalahan koneksi ke server'))
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    loadPengurus()
  }, [])

  // Filter & Sort
  const filteredItems = items.filter(item => {
    // Check tab active / trash
    const matchesTab = currentTab === 'active' ? item.is_active : !item.is_active
    
    const matchesSearch = item.name.toLowerCase().includes(searchQuery.toLowerCase()) || 
                          item.role.toLowerCase().includes(searchQuery.toLowerCase())
    
    const matchesLevel = !filterLevel || item.level === filterLevel

    return matchesTab && matchesSearch && matchesLevel
  }).sort((a, b) => {
    if (filterSort === 'sort_order') return a.sort_order - b.sort_order
    return a.name.localeCompare(b.name)
  })

  const paginatedItems = filteredItems.slice((currentPage - 1) * pageSize, currentPage * pageSize)
  const totalPages = Math.ceil(filteredItems.length / pageSize) || 1

  const openForm = (item = null) => {
    if (item) {
      setFormMode('edit')
      setFormData({
        id: item.id,
        name: item.name,
        level: item.level,
        role: item.role,
        department: item.department || '',
        periode: item.periode || '2025 - 2030',
        provinsi: item.provinsi || '',
        kabupaten: item.kabupaten || '',
        sort_order: item.sort_order || 1,
        image_url: item.image_url,
        linkedin_url: item.linkedin_url || '',
        facebook_url: item.facebook_url || '',
        instagram_url: item.instagram_url || '',
        whatsapp: item.whatsapp || '',
        is_active: item.is_active
      })
    } else {
      setFormMode('create')
      setFormData({
        id: null,
        name: '',
        level: 'dpp',
        role: '',
        department: '',
        periode: '2025 - 2030',
        provinsi: '',
        kabupaten: '',
        sort_order: items.length + 1,
        image_url: '',
        linkedin_url: '',
        facebook_url: '',
        instagram_url: '',
        whatsapp: '',
        is_active: true
      })
    }
    setIsFormOpen(true)
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    
    // Validation for regional fields
    if ((formData.level === 'dpd' || formData.level === 'dpc') && !formData.provinsi) {
      alert('Provinsi wajib diisi untuk level DPD/DPC.')
      return
    }
    if (formData.level === 'dpc' && !formData.kabupaten) {
      alert('Kabupaten/Kota wajib diisi untuk level DPC.')
      return
    }

    try {
      if (formMode === 'create') {
        await pengurusService.create(formData)
      } else {
        await pengurusService.update(formData.id, formData)
      }
      setIsFormOpen(false)
      loadPengurus()
      alert('Data pengurus berhasil disimpan!')
    } catch {
      alert('Gagal menyimpan data pengurus')
    }
  }

  const handleDelete = async (id) => {
    if (window.confirm('Yakin ingin menghapus pengurus ini?')) {
      try {
        await pengurusService.remove(id)
        loadPengurus()
        alert('Data pengurus berhasil dihapus!')
      } catch {
        alert('Gagal menghapus pengurus')
      }
    }
  }

  const handleToggleActive = async (item) => {
    const updatedStatus = !item.is_active
    try {
      await pengurusService.update(item.id, {
        ...item,
        is_active: updatedStatus
      })
      loadPengurus()
    } catch {
      alert('Gagal mengubah status pengurus')
    }
  }

  return (
    <AdminLayout title="Kelola Pengurus">
      <div className="space-y-6">
        
        {/* Navigation Tabs */}
        <div className="flex justify-between items-center bg-white p-4 rounded-xl border border-gray-200">
          <div className="flex gap-2">
            <button 
              onClick={() => { setCurrentTab('active'); setCurrentPage(1); }}
              className={`px-4 py-2 text-xs font-semibold rounded-lg ${currentTab === 'active' ? 'bg-brand-600 text-white' : 'bg-gray-100 text-slate-700'}`}
            >
              Pengurus Aktif
            </button>
            <button 
              onClick={() => { setCurrentTab('trash'); setCurrentPage(1); }}
              className={`px-4 py-2 text-xs font-semibold rounded-lg ${currentTab === 'trash' ? 'bg-brand-600 text-white' : 'bg-gray-100 text-slate-700'}`}
            >
              Non-aktif
            </button>
          </div>
          {currentTab === 'active' && (
            <button 
              onClick={() => openForm()}
              className="px-4 py-2 bg-brand-600 text-white rounded-lg hover:bg-brand-700 text-xs font-semibold"
            >
              + Tambah Pengurus Baru
            </button>
          )}
        </div>

        {/* Filter Bar */}
        <div className="flex flex-col md:flex-row gap-4 justify-between items-center bg-white p-4 rounded-xl border border-gray-200">
          <input 
            type="text" 
            placeholder="Cari pengurus..." 
            value={searchQuery}
            onChange={e => { setSearchQuery(e.target.value); setCurrentPage(1); }}
            className="w-full md:max-w-md px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none"
          />
          <div className="flex gap-2 w-full md:w-auto">
            <select 
              value={filterLevel}
              onChange={e => { setFilterLevel(e.target.value); setCurrentPage(1); }}
              className="px-3 py-2 border border-gray-300 rounded-lg text-sm text-slate-600 outline-none w-full sm:w-auto"
            >
              <option value="">Semua Tingkatan</option>
              <option value="ketua">Ketua Umum</option>
              <option value="dpp">Pusat (DPP)</option>
              <option value="dpd">Provinsi (DPD)</option>
              <option value="dpc">Kab/Kota (DPC)</option>
            </select>
            <select 
              value={filterSort}
              onChange={e => { setFilterSort(e.target.value); setCurrentPage(1); }}
              className="px-3 py-2 border border-gray-300 rounded-lg text-sm text-slate-600 outline-none w-full sm:w-auto"
            >
              <option value="sort_order">Urutan Tampil</option>
              <option value="name">Abjad Nama</option>
            </select>
          </div>
        </div>

        {loading && <div className="text-slate-500 py-10 text-center">Memuat pengurus...</div>}
        {error && <div className="text-red-600 py-10 text-center font-medium">{error}</div>}

        {!loading && !error && (
          <div className="bg-white rounded-2xl border border-gray-200 overflow-hidden shadow-sm">
            <table className="w-full text-left text-sm text-slate-700">
              <thead className="bg-slate-50 border-b border-gray-200 font-semibold">
                <tr>
                  <th className="p-4">Nama Lengkap</th>
                  <th className="p-4">Tingkat</th>
                  <th className="p-4">Jabatan</th>
                  <th className="p-4">Wilayah</th>
                  <th className="p-4">Aksi</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-100">
                {paginatedItems.map(item => (
                  <tr key={item.id} className="hover:bg-slate-50/50">
                    <td className="p-4 font-medium text-slate-900 flex items-center gap-3">
                      <img src={item.image_url.startsWith('http') ? item.image_url : `http://127.0.0.1:8080${item.image_url}`} alt="" className="w-8 h-8 rounded-full object-cover border" />
                      {item.name}
                    </td>
                    <td className="p-4 uppercase">{item.level}</td>
                    <td className="p-4">{item.role}</td>
                    <td className="p-4">{item.provinsi ? `${item.provinsi}${item.kabupaten ? `, ${item.kabupaten}` : ''}` : '-'}</td>
                    <td className="p-4 flex gap-2">
                      <button onClick={() => handleToggleActive(item)} className="text-xs text-brand-600 font-semibold hover:underline">
                        {item.is_active ? 'Non-aktifkan' : 'Aktifkan'}
                      </button>
                      <button onClick={() => openForm(item)} className="text-xs text-slate-600 font-semibold hover:underline">Edit</button>
                      <button onClick={() => handleDelete(item.id)} className="text-xs text-red-600 font-semibold hover:underline">Hapus</button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        {/* Form Modal */}
        {isFormOpen && (
          <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
            <div className="bg-white rounded-2xl p-6 max-w-lg w-full shadow-2xl max-h-[90vh] overflow-y-auto space-y-4">
              <h3 className="font-heading font-bold text-lg text-slate-900">{formMode === 'create' ? 'Tambah Pengurus Baru' : 'Edit Pengurus'}</h3>
              
              <form onSubmit={handleSubmit} className="space-y-4">
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Nama Lengkap</label>
                  <input type="text" value={formData.name} onChange={e => setFormData({...formData, name: e.target.value})} required className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Tingkat Struktur</label>
                    <select value={formData.level} onChange={e => setFormData({...formData, level: e.target.value})} className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none bg-white">
                      <option value="ketua">Ketua Umum</option>
                      <option value="dpp">Pusat (DPP)</option>
                      <option value="dpd">Provinsi (DPD)</option>
                      <option value="dpc">Kab/Kota (DPC)</option>
                    </select>
                  </div>
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Jabatan Resmi</label>
                    <input type="text" value={formData.role} onChange={e => setFormData({...formData, role: e.target.value})} required placeholder="Misal: Ketua Bidang Organisasi" className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                  </div>
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Departemen / Divisi (Opsional)</label>
                    <input type="text" value={formData.department} onChange={e => setFormData({...formData, department: e.target.value})} className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                  </div>
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Periode Jabatan</label>
                    <input type="text" value={formData.periode} onChange={e => setFormData({...formData, periode: e.target.value})} required className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                  </div>
                </div>

                {(formData.level === 'dpd' || formData.level === 'dpc') && (
                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <label className="block text-xs font-semibold text-gray-500 mb-1">Provinsi</label>
                      <input type="text" value={formData.provinsi} onChange={e => setFormData({...formData, provinsi: e.target.value})} required className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                    </div>
                    {formData.level === 'dpc' && (
                      <div>
                        <label className="block text-xs font-semibold text-gray-500 mb-1">Kabupaten/Kota</label>
                        <input type="text" value={formData.kabupaten} onChange={e => setFormData({...formData, kabupaten: e.target.value})} required className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                      </div>
                    )}
                  </div>
                )}

                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">URL Foto Profil</label>
                    <input type="text" value={formData.image_url} onChange={e => setFormData({...formData, image_url: e.target.value})} required className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                  </div>
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Urutan Sortir</label>
                    <input type="number" value={formData.sort_order} onChange={e => setFormData({...formData, sort_order: parseInt(e.target.value)})} required className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none" />
                  </div>
                </div>

                <div className="border p-3 rounded-lg space-y-2 bg-slate-50/50">
                  <span className="block text-[11px] font-bold text-slate-500 uppercase tracking-widest">Media Sosial & Kontak</span>
                  <div className="grid grid-cols-2 gap-3">
                    <input type="text" value={formData.facebook_url} onChange={e => setFormData({...formData, facebook_url: e.target.value})} placeholder="URL Facebook" className="px-3 py-1.5 border border-gray-300 rounded-lg text-xs outline-none" />
                    <input type="text" value={formData.instagram_url} onChange={e => setFormData({...formData, instagram_url: e.target.value})} placeholder="URL Instagram" className="px-3 py-1.5 border border-gray-300 rounded-lg text-xs outline-none" />
                    <input type="text" value={formData.linkedin_url} onChange={e => setFormData({...formData, linkedin_url: e.target.value})} placeholder="URL LinkedIn" className="px-3 py-1.5 border border-gray-300 rounded-lg text-xs outline-none" />
                    <input type="text" value={formData.whatsapp} onChange={e => setFormData({...formData, whatsapp: e.target.value})} placeholder="No. WhatsApp (+62...)" className="px-3 py-1.5 border border-gray-300 rounded-lg text-xs outline-none" />
                  </div>
                </div>

                <div className="flex justify-end gap-2 pt-4 border-t border-gray-100">
                  <button type="button" onClick={() => setIsFormOpen(false)} className="px-4 py-2 border border-gray-300 rounded-lg text-xs font-semibold text-gray-700">Batal</button>
                  <button type="submit" className="px-4 py-2 bg-brand-600 text-white rounded-lg hover:bg-brand-700 text-xs font-semibold">Simpan</button>
                </div>
              </form>
            </div>
          </div>
        )}
      </div>
    </AdminLayout>
  )
}
