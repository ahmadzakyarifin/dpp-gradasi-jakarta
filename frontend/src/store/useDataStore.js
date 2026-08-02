import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { beritaService } from '../services/beritaService'
import { kegiatanService } from '../services/kegiatanService'
import { pengurusService } from '../services/pengurusService'

export const useDataStore = create(
  persist(
    (set, get) => ({
      // INITIAL SEED DATA
      berita: [
        {
          id: 1,
          title: 'Rapat Kerja Daerah Jatim',
          slug: 'rapat-kerja-daerah-jatim',
          category: 'Berita Daerah',
          published_date: '2026-02-11',
          image_url: 'https://gradasi.org/uploads/img/berita/17708152730.jpg',
          excerpt: 'SURABAYA, Generasi Digital Indonesia (GRADASI) Jawa Timur bersiap menggelar Rapat Kerja Daerah untuk menyelaraskan program kerja digitalisasi UMKM.',
          content: 'SURABAYA, Generasi Digital Indonesia (GRADASI) Jawa Timur bersiap menggelar Rapat Kerja Daerah...',
          views: 142,
          is_published: true
        },
        {
          id: 2,
          title: 'Peningkatan Kompetensi SDM',
          slug: 'peningkatan-kompetensi-sdm-pendidikan',
          category: 'Edukasi',
          published_date: '2025-11-02',
          image_url: 'https://gradasi.org/uploads/img/berita/17620765070.jpg',
          excerpt: 'Inisiatif GRADASI Mendorong Peningkatan Kompetensi SDM Pendidikan dalam Memanfaatkan Kecerdasan Buatan (AI) secara bijak.',
          content: 'Inisiatif GRADASI Mendorong Peningkatan Kompetensi SDM Pendidikan...',
          views: 98,
          is_published: true
        },
        {
          id: 3,
          title: 'Rumusan Kunci Kebijakan',
          slug: 'rumusan-kunci-kebijakan-literasi-digital',
          category: 'Berita Utama',
          published_date: '2025-10-31',
          image_url: 'https://gradasi.org/uploads/img/berita/17618789900.jpg',
          excerpt: '#Ketua Dewan Pakar GRADASI, Damar Juniarto, Paparkan Lima Rumusan Kunci Kebijakan untuk Mempercepat Transformasi Digital.',
          content: '#Ketua Dewan Pakar GRADASI, Damar Juniarto, Paparkan Lima Rumusan Kunci Kebijakan...',
          views: 215,
          is_published: true
        }
      ],

      kegiatan: [
        {
          id: 1,
          title: 'Penyaluran Bantuan Kemanusiaan oleh DPP GRADASI',
          slug: 'penyaluran-bantuan-kemanusiaan',
          category: 'Nasional',
          organizer: 'DPP GRADASI',
          event_date: '31 Desember 2025',
          location: 'Jakarta',
          image_url: 'https://gradasi.org/uploads/img/event/1767154719.jpg',
          excerpt: 'Dewan Pimpinan Pusat (DPP) GRADASI turun langsung menyalurkan bantuan kemanusiaan kepada masyarakat yang terdampak bencana alam sebagai wujud kepedulian sosial.',
          content: 'Dewan Pimpinan Pusat (DPP) GRADASI turun langsung menyalurkan bantuan...',
          is_published: true
        },
        {
          id: 2,
          title: 'Pelatihan Digital Marketing UMKM Go Online',
          slug: 'pelatihan-digital-marketing-umkm',
          category: 'Jawa Timur',
          organizer: 'DPD GRADASI Jatim',
          event_date: '15 November 2025',
          location: 'Surabaya',
          image_url: 'https://gradasi.org/uploads/img/event/1767154619.jpg',
          excerpt: 'Program pelatihan intensif bagi pelaku Usaha Mikro Kecil Menengah (UMKM) untuk memasarkan produknya secara digital demi menjangkau pasar yang lebih luas.',
          content: 'Program pelatihan intensif bagi pelaku Usaha Mikro Kecil Menengah...',
          is_published: true
        },
        {
          id: 3,
          title: 'Konsolidasi Pengurus DPP & Penyerahan SK Daerah',
          slug: 'konsolidasi-pengurus-dpp-dpd',
          category: 'Lampung',
          organizer: 'DPP GRADASI',
          event_date: '02 Oktober 2025',
          location: 'Bandar Lampung',
          image_url: 'https://gradasi.org/uploads/img/event/1767154397.jpg',
          excerpt: 'Acara konsolidasi pengurus tingkat pusat serta penyerahan Surat Keputusan (SK) kepada perwakilan pengurus daerah demi memperkuat struktur organisasi di seluruh nusantara.',
          content: 'Acara konsolidasi pengurus tingkat pusat serta penyerahan SK...',
          is_published: true
        }
      ],

      pengurus: [
        { id: 1, name: 'Upi Asmaradhana', role: 'Ketua Umum DPP GRADASI', level: 'ketua', is_active: true, periode: '2024 - 2029', image_url: 'https://gradasi.org/uploads/img/s-anggota/ketua/1735027418.jpg', sort_order: 1 },
        { id: 2, name: 'Dr. Susi Susanti, M.Pd', role: 'Wakil Ketua I', level: 'dpp', is_active: true, image_url: 'https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?q=80&w=200', sort_order: 1 },
        { id: 3, name: 'Ir. Budi Santoso', role: 'Wakil Ketua II', level: 'dpp', is_active: true, image_url: 'https://images.unsplash.com/photo-1519085360753-af0119f7cbe7?q=80&w=200', sort_order: 2 },
        { id: 4, name: 'Junaidi, S.Kom', role: 'Sekretaris Jenderal', level: 'dpp', is_active: true, image_url: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?q=80&w=200', sort_order: 3 },
        { id: 5, name: 'Drs. H. Ahmad Fauzi', role: 'Ketua DPD Jawa Barat', level: 'dpd', provinsi: 'Jawa Barat', is_active: true, image_url: 'https://images.unsplash.com/photo-1560250097-0b93528c311a?q=80&w=200', sort_order: 1 },
        { id: 6, name: 'Bambang Irawan, S.T', role: 'Ketua DPD Jawa Timur', level: 'dpd', provinsi: 'Jawa Timur', is_active: true, image_url: 'https://images.unsplash.com/photo-1519085360753-af0119f7cbe7?q=80&w=200', sort_order: 2 }
      ],

      // ACTIONS FOR BERITA
      addBerita: (item) => set(state => ({ berita: [item, ...state.berita] })),
      updateBerita: (id, payload) => set(state => ({ berita: state.berita.map(i => i.id === id ? { ...i, ...payload } : i) })),
      deleteBerita: (id) => set(state => ({ berita: state.berita.filter(i => i.id !== id) })),

      // ACTIONS FOR KEGIATAN
      addKegiatan: (item) => set(state => ({ kegiatan: [item, ...state.kegiatan] })),
      updateKegiatan: (id, payload) => set(state => ({ kegiatan: state.kegiatan.map(i => i.id === id ? { ...i, ...payload } : i) })),
      deleteKegiatan: (id) => set(state => ({ kegiatan: state.kegiatan.filter(i => i.id !== id) })),

      // ACTIONS FOR PENGURUS
      addPengurus: (item) => set(state => ({ pengurus: [item, ...state.pengurus] })),
      updatePengurus: (id, payload) => set(state => ({ pengurus: state.pengurus.map(i => i.id === id ? { ...i, ...payload } : i) })),
      deletePengurus: (id) => set(state => ({ pengurus: state.pengurus.filter(i => i.id !== id) })),

      // SYNC DATA FROM SERVER API IF AVAILABLE
      fetchInitialData: async () => {
        try {
          const [bRes, kRes, pRes] = await Promise.all([
            beritaService.list().catch(() => null),
            kegiatanService.list().catch(() => null),
            pengurusService.list().catch(() => null),
          ])

          if (bRes?.data?.berita && bRes.data.berita.length > 0) set({ berita: bRes.data.berita.map(b => ({ ...b, image_url: b.image_path || b.image_url })) })
          if (kRes?.data?.kegiatan && kRes.data.kegiatan.length > 0) set({ kegiatan: kRes.data.kegiatan.map(k => ({ ...k, image_url: k.image_path || k.image_url })) })
          if (pRes?.data?.pengurus && pRes.data.pengurus.length > 0) set({ pengurus: pRes.data.pengurus })
        } catch {}
      }
    }),
    {
      name: 'gradasi_app_store',
    }
  )
)
