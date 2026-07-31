import { apiRequest } from '../api'

export const dashboardService = {
  async getSummary() {
    return apiRequest('/admin/dashboard/summary')
  },
}
