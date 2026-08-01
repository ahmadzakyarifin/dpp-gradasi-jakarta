import { create } from 'zustand'
import { authService } from '../services/authService'

export const useAuthStore = create((set, get) => ({
  token: localStorage.getItem('access_token') || null,
  user: null,
  loading: false,
  error: null,

  setToken: (token) => {
    if (token) {
      localStorage.setItem('access_token', token)
    } else {
      localStorage.removeItem('access_token')
    }
    set({ token })
  },

  setUser: (user) => set({ user }),

  login: async ({ email, password, rememberMe = false, captcha_token = '' }) => {
    set({ loading: true, error: null })
    try {
      const response = await authService.login({
        email,
        password,
        remember_me: rememberMe,
        captcha_token,
      })
      const token = response?.data?.access_token || response?.data?.token || response?.access_token
      const user = response?.data?.user || response?.user

      if (!token) {
        set({ loading: false, error: 'Login gagal: token tidak ditemukan.' })
        throw new Error('Token tidak ditemukan pada response login.')
      }

      localStorage.setItem('access_token', token)
      set({
        token,
        user,
        loading: false,
        error: null,
      })
      return response
    } catch (err) {
      set({ loading: false, error: err?.message || 'Email atau password salah.' })
      throw err // re-throw error asli agar retryAfter/code ikut terbawa ke UI
    }
  },

  fetchMe: async () => {
    if (!get().token) return null
    try {
      const response = await authService.me()
      if (response && response.data) {
        set({ user: response.data })
        return response.data
      }
    } catch {
      // Return local user fallback
    }
    return get().user
  },

  logoutLocal: () => {
    localStorage.removeItem('access_token')
    set({ token: null, user: null, error: null })
  },

  logout: async () => {
    try {
      await authService.logout()
    } catch {
      // Tetap hapus sesi lokal walaupun backend gagal/expired.
    } finally {
      get().logoutLocal()
    }
  },
}))
