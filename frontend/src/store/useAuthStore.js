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

  login: async ({ email, password, rememberMe = false }) => {
    set({ loading: true, error: null })
    try {
      const response = await authService.login({
        email,
        password,
        remember_me: rememberMe,
      })
      const token = response.data.access_token
      localStorage.setItem('access_token', token)
      set({
        token,
        user: response.data.user,
        loading: false,
        error: null,
      })
      return response
    } catch (error) {
      set({ loading: false, error: error.message })
      throw error
    }
  },

  fetchMe: async () => {
    if (!get().token) return null
    try {
      const response = await authService.me()
      set({ user: response.data })
      return response.data
    } catch (error) {
      get().logoutLocal()
      throw error
    }
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
