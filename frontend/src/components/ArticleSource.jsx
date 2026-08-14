/**
 * ArticleSource - Reusable footnote / source information component
 * for article details (Berita / Kegiatan).
 * Parses line-by-line formatted string: "Label: Value"
 */
export default function ArticleSource({ footnote }) {
  if (!footnote) return null

  const lines = footnote.split('\n').map(l => l.trim()).filter(Boolean)
  if (lines.length === 0) return null

  return (
    <div 
      style={{ backgroundColor: '#F5F6F8' }}
      className="mt-8 mb-6 py-3 px-4 border border-slate-200/50 rounded-lg max-w-fit text-[13.5px] text-[#6B7280] flex flex-col gap-2.5 shadow-xs"
      aria-label="Informasi Sumber dan Kontributor"
    >
      {lines.map((line, idx) => {
        const colonIndex = line.indexOf(':')
        let label = 'Sumber'
        let value = line

        if (colonIndex !== -1) {
          label = line.substring(0, colonIndex).trim()
          value = line.substring(colonIndex + 1).trim()
        }

        // Determine icon based on label type
        const labelLower = label.toLowerCase()
        let iconClass = 'ph-bold ph-link'
        let iconLabel = 'Tautan'

        if (labelLower.includes('penulis') || labelLower.includes('kontributor')) {
          iconClass = 'ph-bold ph-user-circle'
          iconLabel = 'Penulis'
        } else if (labelLower.includes('foto') || labelLower.includes('dokumentasi') || labelLower.includes('kamera') || labelLower.includes('gambar')) {
          iconClass = 'ph-bold ph-camera'
          iconLabel = 'Foto'
        } else if (labelLower.includes('editor') || labelLower.includes('redaksi') || labelLower.includes('penyunting')) {
          iconClass = 'ph-bold ph-note-pencil'
          iconLabel = 'Editor'
        } else if (labelLower.includes('sumber') || labelLower.includes('referensi') || labelLower.includes('kutip')) {
          iconClass = 'ph-bold ph-quotes'
          iconLabel = 'Sumber'
        }

        return (
          <div key={idx} className="flex items-center gap-2.5">
            <i 
              className={`${iconClass} text-[#6B7280] shrink-0 text-sm`} 
              title={iconLabel}
              aria-label={`Ikon ${iconLabel}`} 
            />
            <span className="font-semibold text-slate-800">{label}:</span>
            <span className="text-[#6B7280] font-medium">{value}</span>
          </div>
        )
      })}
    </div>
  )
}
