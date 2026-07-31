export function formatDate(date) {
  if (!date) return '-'
  const parsed = new Date(date)
  if (Number.isNaN(parsed.getTime())) return date
  return parsed.toLocaleDateString('id-ID', {
    day: '2-digit',
    month: 'long',
    year: 'numeric',
  })
}

export function stripHtml(html = '') {
  return html.replace(/<[^>]+>/g, '')
}
