// Pure helper functions for ActivityLogAdmin.jsx.
// Extracted so the field-mapping and pagination-extraction logic can be
// unit/property tested without a JSX/React test harness.
//
// Field names below match the real backend response shape:
// ActivityLogItemRes (backend/internal/module/activitylog/dto/activity_log_response.go)
// and ActivityLogPaginationRes, confirmed against docs-final/api/activity_logs.jsonc.

/**
 * Formats an ISO 8601 `created_at` timestamp into separate date/time display
 * strings for the table's "Waktu" column.
 * @param {string} isoString
 * @returns {{ datePart: string, hourPart: string }}
 */
function formatCreatedAt(isoString) {
  if (!isoString) {
    return { datePart: '', hourPart: '' }
  }
  let date = new Date(isoString)
  if (Number.isNaN(date.getTime())) {
    return { datePart: '', hourPart: '' }
  }
  // Backend (model_mapper.go) sekarang mengirim RFC3339 UTC jujur (…Z / +07:00).
  // Untuk data lama yang mungkin dikirim sebagai local-time tanpa offset
  // (…T11:07:58 tanpa Z), jangan biarkan JS menganggapnya UTC lalu +7 jam:
  // parse manual, dan kalau tidak ada offset/Z, perlakukan sebagai waktu lokal.
  if (!/Z|[+-]\d{2}:?\d{2}$/.test(isoString.trim())) {
    const m = isoString.trim().match(/^(\d{4})-(\d{2})-(\d{2})[T ](\d{2}):(\d{2})(?::(\d{2}))?/)
    if (m) {
      const [, y, mo, d, h, mi, s = '0'] = m
      date = new Date(+y, +mo - 1, +d, +h, +mi, +s)
    }
  }
  const datePart = date.toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric', timeZone: 'Asia/Jakarta' })
  const hourPart = date.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit', second: '2-digit', timeZone: 'Asia/Jakarta' })
  return { datePart, hourPart }
}

/**
 * Maps a single activity log row (as returned by the backend) into the
 * ordered array of CSV cell values used by downloadCSV().
 * @param {object} log
 * @returns {string[]}
 */
export function mapLogRowForCsv(log) {
  return [
    log.created_at || '',
    log.actor_name || '',
    log.actor_role || '',
    log.action || '',
    `${log.entity_type || ''} ${log.entity_label ? `(${log.entity_label})` : ''}`,
    log.description || '',
    log.ip_address || '',
    log.user_agent || '',
    log.risk_level || ''
  ]
}

/**
 * Maps a single activity log row into the fields used to render the table row.
 * @param {object} log
 * @returns {{ datePart: string, hourPart: string, id: *, actor: string, role: string, ip: string, device: string, action: *, entity: *, entityLabel: *, risk: string, description: * }}
 */
export function mapLogRowForDisplay(log) {
  const { datePart, hourPart } = formatCreatedAt(log.created_at)

  return {
    datePart,
    hourPart,
    id: log.id,
    actor: log.actor_name || 'System',
    role: log.actor_role || 'System',
    ip: log.ip_address || '127.0.0.1',
    device: log.user_agent || 'N/A',
    action: log.action,
    entity: log.entity_type,
    entityLabel: log.entity_label || '-',
    risk: log.risk_level,
    description: log.description
  }
}

/**
 * Extracts pagination totals from a list response body (res.data).
 * Returns null when the response has no pagination meta at all, so callers
 * can choose to leave existing state untouched (matching the original
 * `if (res.data.pagination) { ... }` guard's intent, but pointed at the
 * real `meta` key returned by the backend).
 * @param {object} resData - the `data` field of the activity log list response
 * @returns {{ total: number, totalPages: number } | null}
 */
export function extractPaginationMeta(resData) {
  if (!resData || !resData.meta) {
    return null
  }
  return {
    total: resData.meta.total_data || 0,
    totalPages: resData.meta.total_pages || 1
  }
}
