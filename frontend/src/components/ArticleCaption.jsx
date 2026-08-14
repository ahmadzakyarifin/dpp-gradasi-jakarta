import { useState } from 'react'

/**
 * ArticleCaption - Reusable caption/source text under cover image
 * with support for dynamic prefixes, text wrapping, and expansion.
 */
export default function ArticleCaption({ source }) {
  const [isExpanded, setIsExpanded] = useState(false)
  
  if (!source) return null
  const trimmed = source.trim()
  if (!trimmed) return null

  // Ensure standard prefix format if not already present
  let formatted = trimmed
  const hasPrefix = /^[a-zA-ZÀ-ÿ0-9\s/]+:/.test(trimmed)
  if (!hasPrefix) {
    formatted = `Gambar: ${trimmed}`
  }

  // Check if text is long enough to warrant a "Read More" button
  // Typically more than 100 characters or contains newlines
  const isLongText = formatted.length > 100

  return (
    <div className="w-full text-center mt-2 px-4 mx-auto max-w-[85%] sm:max-w-[70%]">
      <div className="inline-flex flex-col items-center">
        <p 
          style={{ 
            wordBreak: 'break-word', 
            overflowWrap: 'break-word',
            whiteSpace: 'normal',
          }}
          className={`text-[12.5px] text-[#6B7280] font-normal italic leading-relaxed text-center ${
            !isExpanded && isLongText ? 'line-clamp-2' : ''
          }`}
        >
          {formatted}
        </p>
        
        {isLongText && (
          <button
            type="button"
            onClick={() => setIsExpanded(!isExpanded)}
            className="text-[11px] font-bold text-brand-600 hover:text-brand-700 transition mt-1 not-italic"
            aria-label={isExpanded ? "Sembunyikan detail keterangan" : "Lihat selengkapnya detail keterangan"}
          >
            {isExpanded ? 'Sembunyikan' : 'Lihat Selengkapnya'}
          </button>
        )}
      </div>
    </div>
  )
}
