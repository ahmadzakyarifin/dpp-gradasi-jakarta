import { createContext } from 'react'

// Context object dipisah dari provider supaya file provider hanya berisi komponen
// (fast-refresh React tetap bekerja) dan hook konsumen (useSettings) punya
// referensi context yang stabil.
export const SettingsContext = createContext(null)
