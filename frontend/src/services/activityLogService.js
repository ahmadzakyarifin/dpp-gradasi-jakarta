import { apiRequest } from '../api'

function toQuery(params = {}) {
  const query = new URLSearchParams()
  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== null && value !== '') {
      query.set(key, value)
    }
  })
  const text = query.toString()
  return text ? `?${text}` : ''
}

export const activityLogService = {
  list(params) {
    return apiRequest(`/activity-logs${toQuery(params)}`)
  },

  summary(params) {
    return apiRequest(`/activity-logs/summary${toQuery(params)}`)
  },

  detail(id) {
    return apiRequest(`/activity-logs/${id}`)
  },

  entityLogs(entityType, entityId) {
    return apiRequest(`/activity-logs/entity/${encodeURIComponent(entityType)}/${entityId}`)
  }
}
