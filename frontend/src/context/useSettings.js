import { useContext } from 'react'
import { SettingsContext } from './SettingsContext'

// Dipisah dari SettingsContext.jsx supaya fast-refresh React tidak rusak
// (file dengan export context provider + hook dalam satu modul).
export function useSettings() {
  const ctx = useContext(SettingsContext)
  if (!ctx) throw new Error('useSettings harus dipakai di dalam SettingsProvider')
  return ctx
}
