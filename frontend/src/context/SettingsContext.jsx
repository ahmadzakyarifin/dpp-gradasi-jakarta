import { useEffect, useState, useCallback } from 'react'
import { settingsService } from '../services/settingsService'
import { SettingsContext } from './settingsContextObject'

const DEFAULT_SETTINGS = {
  site_name: '',
  tagline: '',
  logo_path: '',
  contact_email: '',
  contact_phone: '',
  address: '',
  maps_embed_url: '',
  facebook_url: '',
  instagram_url: '',
  youtube_url: '',
  video_profile_path: '',
  history: '',
  about_tutorial: '',
  about_formation_date: '',
  about_no_sk: '',
  about_vision: '',
  about_mission: '[]',
  greeting_title: '',
  greeting_subtitle: '',
  greeting_date: '',
  greeting_content: '',
  greeting_image_path: '',
  login_hero_title: '',
  login_hero_description: '',
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

  useEffect(() => {
    if (settings?.site_name) {
      document.title = `${settings.site_name}${settings.tagline ? ' - ' + settings.tagline : ''}`
    }
  }, [settings])

  return (
    <SettingsContext.Provider value={{ settings, loading, error, refresh, setSettings }}>
      {children}
    </SettingsContext.Provider>
  )
}
