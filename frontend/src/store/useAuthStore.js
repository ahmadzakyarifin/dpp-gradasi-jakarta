import { create } from 'zustand'

export const useAuthStore = create((set) => ({
  token: localStorage.getItem('access_token') || null,
  user: null,
  setToken: (token) => {
    if (token) {
      localStorage.setItem('access_token', token)
    } else {
      localStorage.removeItem('access_token')
    }
    set({ token })
  },
  setUser: (user) => set({ user }),
  logout: () => {
    localStorage.removeItem('access_token')
    set({ token: null, user: null })
  }
}))
