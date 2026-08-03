import { apiRequest } from '../api'

// GET /api/v1/roles — daftar role (super_admin, admin)
export const roleService = {
  list() {
    return apiRequest('/roles')
  },
}
