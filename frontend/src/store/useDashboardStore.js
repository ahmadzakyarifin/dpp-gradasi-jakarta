import { create } from 'zustand'
import { dashboardService } from '../services/dashboardService'

export const useDashboardStore = create((set) => ({
  summary: null,
  loading: false,
  error: null,

  fetchSummary: async () => {
    set({ loading: true, error: null })
    try {
      const response = await dashboardService.getSummary()
      set({ summary: response.data, loading: false })
    } catch (error) {
      set({ error: error.message, loading: false })
    }
  },
}))
