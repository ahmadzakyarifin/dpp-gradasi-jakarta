import { useEffect, useState, useCallback } from 'react'
import { settingsService } from '../services/settingsService'
import { SettingsContext } from './settingsContextObject'

const DEFAULT_SETTINGS = {
  site_name: 'DPP GRADASI',
  tagline: 'Generasi Digital Indonesia',
  logo_path: 'https://gradasi.org/uploads/img/logo/1737187847.png',
  contact_email: 'dpp@gradasi.org',
  contact_phone: '+6281234567890',
  address: 'Office Park OL3-IZA The Bellagio Mall, Mega Kuningan, Jakarta Selatan',
  maps_embed_url: 'https://maps.google.com/maps?q=The%20Bellagio%20Mall%20Mega%20Kuningan%20Jakarta&t=&z=16&ie=UTF8&iwloc=&output=embed',
  facebook_url: 'https://www.facebook.com/gradasiofficial.id',
  instagram_url: 'https://www.instagram.com/dppgradasi',
  youtube_url: 'https://www.youtube.com/channel/UCwdjB4LkqcF4Kw5-PoyOb5A',
  video_profile_path: 'https://gradasi.org/assets/video/gradasi.mp4',
  history: 'Perkumpulan Generasi Digital Indonesia (GRADASI) didirikan pada 4 Februari 2019 sebagai organisasi independen yang berfokus pada pengembangan literasi digital, pemberdayaan UMKM, dan transformasi teknologi di Indonesia.',
  about_tutorial: 'Pengesahan Badan Hukum Kemenkumham RI.',
  about_formation_date: '4 Februari 2019',
  about_no_sk: 'AHU-0000151.AH.01.07.2019',
  about_vision: 'Mewujudkan masyarakat Indonesia yang cerdas, kreatif, dan berdaulat di era digital.',
  about_mission: '["Membangun ekosistem literasi digital yang inklusif di seluruh daerah Indonesia.","Mengakselerasi transformasi digital bagi UMKM dan generasi muda.","Mendorong inovasi dan kolaborasi antar pemangku kepentingan industri kreatif digital."]',
  greeting_title: 'Tahun Baru 2026',
  greeting_subtitle: 'Resolusi & Harapan',
  greeting_date: '11 Februari 2026',
  greeting_content: 'Memasuki tahun 2026, GRADASI menetapkan pilar utama perjuangan: memastikan setiap masyarakat memiliki kecakapan digital (digital skills), serta mengembangkan program literasi yang berdampak nyata bagi pertumbuhan ekonomi lokal.',
  greeting_image_path: 'https://gradasi.org/uploads/img/event-terkini/1767154211.jpg',
  // Status CAPTCHA dari backend (single source of truth) — diisi oleh fetch,
  // default false supaya aman kalau backend belum menyediakan.
  captcha_enabled: false,
  captcha_site_key: '',
}

export function SettingsProvider({ children }) {
  const [settings, setSettings] = useState(DEFAULT_SETTINGS)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  const refresh = useCallback(async () => {
    try {
      const res = await settingsService.get()
      if (res?.success && res?.data) {
        setSettings({ ...DEFAULT_SETTINGS, ...res.data })
      }
    } catch (err) {
      setError(err)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    refresh()
  }, [refresh])

  return (
    <SettingsContext.Provider value={{ settings, loading, error, refresh, setSettings }}>
      {children}
    </SettingsContext.Provider>
  )
}
