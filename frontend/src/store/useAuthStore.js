import { create } from 'zustand'
import { authService } from '../services/authService'

export const useAuthStore = create((set, get) => ({
  token: localStorage.getItem('access_token') || null,
  user: null,
  loading: false,
  error: null,

  // role name dari user.role (di-set saat login/fetchMe); dipakai untuk filter menu FE
  role: null,

  setToken: (token) => {
    if (token) {
      localStorage.setItem('access_token', token)
    } else {
      localStorage.removeItem('access_token')
    }
    set({ token })
  },

  // Setelah inisialisasi awal selesai (fetchMe global di App), flag ini true.
  // Dipakai untuk menghindari flash: AdminLayout menampilkan skeleton sidebar
  // selama role belum diketahui, bukan fallback "tampilkan semua menu".
  appReady: false,
  markAppReady: () => set({ appReady: true }),

  setUser: (user) => {
    // Terima bentuk apa pun: user langsung atau wrapper {user: {...}}
    const u = user?.user || user
    const role = u?.role?.name || u?.role_name || null
    if (role) localStorage.setItem('user_role', role)
    set({ user: u, role })
  },

  logoutLocal: () => {
    localStorage.removeItem('access_token')
    localStorage.removeItem('user_role')
    set({ token: null, user: null, role: null, appReady: true, error: null })
  },

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
      const role = user?.role?.name || user?.role_name || null
      if (role) localStorage.setItem('user_role', role)
      set({
        token,
        user,
        role,
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
      // apiRequest mengembalikan wrapper {success, data: {user: {...}}}
      const user = response?.data?.user || response?.data || response?.user || null
      if (user) {
        const role = user?.role?.name || user?.role_name || null
        if (role) localStorage.setItem('user_role', role)
        set({ user, role })
        return user
      }
    } catch {
      // Token invalid/expired → bersihkan sesi lokal
      localStorage.removeItem('access_token')
      localStorage.removeItem('user_role')
      set({ token: null, user: null, role: null })
    }
    return get().user
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
