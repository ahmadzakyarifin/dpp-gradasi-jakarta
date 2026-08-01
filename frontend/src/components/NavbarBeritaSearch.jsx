import { useEffect, useState } from 'react'
import { useSearchParams } from 'react-router-dom'
import { beritaService } from '../services/beritaService'

// Search + filter kategori + reset berita, tampil di navbar (tengah) khusus
// halaman /berita. State disinkronkan lewat URL (?q=&category=&page=) sehingga
// BeritaList membaca dari satu sumber kebenaran.
export default function NavbarBeritaSearch() {
  const [searchParams, setSearchParams] = useSearchParams()
  const [categories, setCategories] = useState([])

  const q = searchParams.get('q') || ''
  const category = searchParams.get('category') || ''

  useEffect(() => {
    beritaService.getCategories()
      .then(res => {
        if (res && res.data && Array.isArray(res.data)) setCategories(res.data)
      }).catch(() => {})
  }, [])

  function update(next) {
    const params = new URLSearchParams(searchParams)
    Object.entries(next).forEach(([k, v]) => {
      if (v) params.set(k, v)
      else params.delete(k)
    })
    params.set('page', '1')
    setSearchParams(params)
  }

  function reset() {
    setSearchParams({})
  }

  return (
    <div className="hidden lg:flex items-center gap-2 bg-slate-50 border border-slate-200 rounded-full px-3 py-1.5">
      <div className="relative">
        <i className="ph-bold ph-magnifying-glass absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-sm" />
        <input
          type="text"
          value={q}
          onChange={(e) => update({ q: e.target.value })}
          placeholder="Cari berita..."
          className="w-40 xl:w-52 pl-8 pr-3 py-1.5 rounded-full border border-slate-200 text-sm bg-white focus:border-brand-500 focus:outline-none"
        />
      </div>
      <select
        value={category}
        onChange={(e) => update({ category: e.target.value })}
        className="w-32 xl:w-40 px-2.5 py-1.5 rounded-full border border-slate-200 text-sm bg-white focus:border-brand-500 focus:outline-none text-slate-600"
      >
        <option value="">Semua Kategori</option>
        {categories.map((c) => <option key={c} value={c}>{c}</option>)}
      </select>
      {(q || category) && (
        <button
          type="button"
          onClick={reset}
          title="Reset filter"
          className="w-8 h-8 flex items-center justify-center rounded-full text-slate-400 hover:text-brand-600 hover:bg-brand-50 transition"
        >
          <i className="ph-bold ph-x text-base" />
        </button>
      )}
    </div>
  )
}
