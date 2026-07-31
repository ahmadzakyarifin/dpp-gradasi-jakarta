import { useState, useEffect } from 'react'
import AdminLayout from '../../layouts/AdminLayout'
import { pengurusService } from '../../services/pengurusService'

const DEFAULT_PENGURUS = [
  { id: 1, name: 'Upi Asmaradhana', role: 'Ketua Umum DPP GRADASI', level: 'ketua', is_active: true, periode: '2024 - 2029', image_url: 'https://gradasi.org/uploads/img/s-anggota/ketua/1735027418.jpg', sort_order: 1 },
  { id: 2, name: 'Dr. Susi Susanti, M.Pd', role: 'Wakil Ketua I', level: 'dpp', is_active: true, image_url: 'https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?q=80&w=200', sort_order: 1 },
  { id: 3, name: 'Ir. Budi Santoso', role: 'Wakil Ketua II', level: 'dpp', is_active: true, image_url: 'https://images.unsplash.com/photo-1519085360753-af0119f7cbe7?q=80&w=200', sort_order: 2 },
  { id: 4, name: 'Junaidi, S.Kom', role: 'Sekretaris Jenderal', level: 'dpp', is_active: true, image_url: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?q=80&w=200', sort_order: 3 },
  { id: 5, name: 'Drs. H. Ahmad Fauzi', role: 'Ketua DPD Jawa Barat', level: 'dpd', provinsi: 'Jawa Barat', is_active: true, image_url: 'https://images.unsplash.com/photo-1560250097-0b93528c311a?q=80&w=200', sort_order: 1 },
  { id: 6, name: 'Bambang Irawan, S.T', role: 'Ketua DPD Jawa Timur', level: 'dpd', provinsi: 'Jawa Timur', is_active: true, image_url: 'https://images.unsplash.com/photo-1519085360753-af0119f7cbe7?q=80&w=200', sort_order: 2 }
]

export default function PengurusAdmin() {
  const [items, setItems] = useState(DEFAULT_PENGURUS)
  const [loading, setLoading] = useState(false)
  
  const [currentTab, setCurrentTab] = useState('active')
  const [currentPage, setCurrentPage] = useState(1)
  const pageSize = 5

  const [searchQuery, setSearchQuery] = useState('')
  const [filterLevel, setFilterLevel] = useState('')
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
    is_active: true
  })

  const loadPengurus = () => {
    pengurusService.list()
      .then(res => {
        if (res && res.data) {
          const list = Array.isArray(res.data) ? res.data : (res.data.pengurus || [])
          if (list.length > 0) setItems(list)
        }
      })
      .catch(() => {})
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    loadPengurus()
  }, [])

  const filteredItems = items.filter(item => {
    const matchesTab = currentTab === 'active' ? item.is_active !== false : item.is_active === false
    const matchesSearch = !searchQuery || item.name.toLowerCase().includes(searchQuery.toLowerCase()) || 
                          (item.role || '').toLowerCase().includes(searchQuery.toLowerCase())
    const matchesLevel = !filterLevel || item.level === filterLevel

    return matchesTab && matchesSearch && matchesLevel
  }).sort((a, b) => (a.sort_order || 0) - (b.sort_order || 0))

  const paginatedItems = filteredItems.slice((currentPage - 1) * pageSize, currentPage * pageSize)
  const totalPages = Math.max(1, Math.ceil(filteredItems.length / pageSize))

  const openForm = (item = null) => {
    if (item) {
      setFormMode('edit')
      setFormData({ ...item })
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
        is_active: true
      })
    }
    setIsFormOpen(true)
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    if (formMode === 'create') {
      const newObj = { ...formData, id: Date.now() }
      setItems(prev => [newObj, ...prev])
      try { await pengurusService.create(formData) } catch {}
    } else {
      setItems(prev => prev.map(i => i.id === formData.id ? { ...i, ...formData } : i))
      try { await pengurusService.update(formData.id, formData) } catch {}
    }
    setIsFormOpen(false)
  }

  const handleDelete = async (id) => {
    if (window.confirm('Yakin ingin menghapus pengurus ini?')) {
      setItems(prev => prev.filter(i => i.id !== id))
      try { await pengurusService.remove(id) } catch {}
    }
  }

  const handleToggleActive = async (item) => {
    const updatedStatus = !item.is_active
    setItems(prev => prev.map(i => i.id === item.id ? { ...i, is_active: updatedStatus } : i))
    try {
      await pengurusService.update(item.id, { ...item, is_active: updatedStatus })
    } catch {}
  }

  return (
    <AdminLayout title="Kelola Pengurus">
      <div className="space-y-6">
        
        {/* Navigation Tabs */}
        <div className="flex justify-between items-center bg-white p-4 rounded-xl border border-gray-200 shadow-sm">
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
              className="px-4 py-2.5 bg-brand-600 text-white rounded-xl hover:bg-brand-700 text-xs font-semibold shadow-sm"
            >
              + Tambah Pengurus Baru
            </button>
          )}
        </div>

        {/* Filter Bar */}
        <div className="flex flex-col md:flex-row gap-4 justify-between items-center bg-white p-4 rounded-xl border border-gray-200 shadow-sm">
          <input 
            type="text" 
            placeholder="Cari nama atau jabatan pengurus..." 
            value={searchQuery}
            onChange={e => { setSearchQuery(e.target.value); setCurrentPage(1); }}
            className="w-full md:max-w-md px-3.5 py-2.5 bg-slate-50 border border-gray-200 rounded-xl text-sm outline-none"
          />
          <div className="flex gap-2 w-full md:w-auto">
            <select 
              value={filterLevel}
              onChange={e => { setFilterLevel(e.target.value); setCurrentPage(1); }}
              className="px-3.5 py-2.5 bg-slate-50 border border-gray-200 rounded-xl text-sm text-slate-600 outline-none w-full sm:w-auto"
            >
              <option value="">Semua Tingkatan</option>
              <option value="ketua">Ketua Umum</option>
              <option value="dpp">Pusat (DPP)</option>
              <option value="dpd">Provinsi (DPD)</option>
              <option value="dpc">Kab/Kota (DPC)</option>
            </select>
          </div>
        </div>

        <div className="bg-white rounded-2xl border border-gray-200 overflow-hidden shadow-sm">
          <table className="w-full text-left text-sm text-slate-700">
            <thead className="bg-slate-50 border-b border-gray-200 font-semibold text-xs uppercase tracking-wider text-slate-500">
              <tr>
                <th className="p-4">Nama Lengkap</th>
                <th className="p-4">Tingkat</th>
                <th className="p-4">Jabatan</th>
                <th className="p-4">Wilayah</th>
                <th className="p-4 text-right">Aksi</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100">
              {paginatedItems.map(item => (
                <tr key={item.id} className="hover:bg-slate-50/50">
                  <td className="p-4 font-medium text-slate-900 flex items-center gap-3">
                    <img src={item.image_url} alt="" className="w-8 h-8 rounded-full object-cover border" />
                    {item.name}
                  </td>
                  <td className="p-4 uppercase font-bold text-xs text-brand-600">{item.level}</td>
                  <td className="p-4">{item.role}</td>
                  <td className="p-4 text-slate-500 text-xs">{item.provinsi ? `${item.provinsi}${item.kabupaten ? `, ${item.kabupaten}` : ''}` : '-'}</td>
                  <td className="p-4 text-right">
                    <div className="flex justify-end gap-3">
                      <button onClick={() => handleToggleActive(item)} className="text-xs text-brand-600 font-semibold hover:underline">
                        {item.is_active ? 'Non-aktifkan' : 'Aktifkan'}
                      </button>
                      <button onClick={() => openForm(item)} className="text-xs text-slate-600 font-semibold hover:underline">Edit</button>
                      <button onClick={() => handleDelete(item.id)} className="text-xs text-red-600 font-semibold hover:underline">Hapus</button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {/* Form Modal */}
        {isFormOpen && (
          <div className="fixed inset-0 bg-black/40 backdrop-blur-[2px] flex items-center justify-center z-50 p-4">
            <div className="bg-white rounded-2xl p-6 max-w-lg w-full shadow-2xl space-y-4 max-h-[90vh] overflow-y-auto">
              <h3 className="font-heading font-bold text-lg text-slate-900">{formMode === 'create' ? 'Tambah Pengurus Baru' : 'Edit Pengurus'}</h3>
              
              <form onSubmit={handleSubmit} className="space-y-4">
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">Nama Lengkap</label>
                  <input type="text" value={formData.name} onChange={e => setFormData({...formData, name: e.target.value})} required className="w-full px-3 py-2 border border-gray-300 rounded-xl text-sm outline-none" />
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Tingkat Struktur</label>
                    <select value={formData.level} onChange={e => setFormData({...formData, level: e.target.value})} className="w-full px-3 py-2 border border-gray-300 rounded-xl text-sm outline-none bg-white">
                      <option value="ketua">Ketua Umum</option>
                      <option value="dpp">Pusat (DPP)</option>
                      <option value="dpd">Provinsi (DPD)</option>
                      <option value="dpc">Kab/Kota (DPC)</option>
                    </select>
                  </div>
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 mb-1">Jabatan Resmi</label>
                    <input type="text" value={formData.role} onChange={e => setFormData({...formData, role: e.target.value})} required placeholder="Misal: Ketua Bidang Organisasi" className="w-full px-3 py-2 border border-gray-300 rounded-xl text-sm outline-none" />
                  </div>
                </div>
                <div>
                  <label className="block text-xs font-semibold text-gray-500 mb-1">URL Foto Profil</label>
                  <input type="text" value={formData.image_url} onChange={e => setFormData({...formData, image_url: e.target.value})} required className="w-full px-3 py-2 border border-gray-300 rounded-xl text-sm outline-none" />
                </div>
                <div className="flex justify-end gap-2 pt-4 border-t">
                  <button type="button" onClick={() => setIsFormOpen(false)} className="px-4 py-2 border rounded-xl text-xs font-semibold">Batal</button>
                  <button type="submit" className="px-4 py-2 bg-brand-600 text-white rounded-xl hover:bg-brand-700 text-xs font-semibold">Simpan</button>
                </div>
              </form>
            </div>
          </div>
        )}
      </div>
    </AdminLayout>
  )
}
