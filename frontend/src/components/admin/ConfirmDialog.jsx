import { useEffect, useRef } from 'react'

export default function ConfirmDialog({ isOpen, title, message, type = 'danger', onConfirm, onCancel, onClose }) {
  const dialogRef = useRef(null)
  const handleCancel = onCancel || onClose

  useEffect(() => {
    if (isOpen) {
      dialogRef.current?.focus()
    }
  }, [isOpen])

  useEffect(() => {
    function handleEsc(e) {
      if (e.key === 'Escape' && isOpen) handleCancel?.()
    }
    document.addEventListener('keydown', handleEsc)
    return () => document.removeEventListener('keydown', handleEsc)
  }, [isOpen, handleCancel])

  if (!isOpen) return null

  const iconMap = {
    danger: { icon: 'ph-bold ph-trash', bg: 'bg-red-100', text: 'text-red-600' },
    success: { icon: 'ph-bold ph-arrow-counter-clockwise', bg: 'bg-emerald-100', text: 'text-emerald-600' },
    warning: { icon: 'ph-bold ph-warning', bg: 'bg-amber-100', text: 'text-amber-600' },
    info: { icon: 'ph-bold ph-paper-plane-right', bg: 'bg-brand-100', text: 'text-brand-600' },
  }

  const btnMap = {
    danger: 'bg-red-600 hover:bg-red-700 focus:ring-red-500/30',
    success: 'bg-emerald-600 hover:bg-emerald-700 focus:ring-emerald-500/30',
    warning: 'bg-amber-600 hover:bg-amber-700 focus:ring-amber-500/30',
    info: 'bg-brand-600 hover:bg-brand-700 focus:ring-brand-500/30',
  }

  const style = iconMap[type] || iconMap.danger
  const btnStyle = btnMap[type] || btnMap.danger

  return (
    <div className="fixed inset-0 z-[60] flex items-center justify-center p-4" role="dialog" aria-modal="true">
      {/* Overlay */}
      <div
        className="fixed inset-0 bg-black/40 backdrop-blur-[2px] animate-[fadeIn_150ms_ease-out]"
        onClick={handleCancel}
      />

      {/* Dialog */}
      <div
        ref={dialogRef}
        tabIndex={-1}
        className="relative bg-white rounded-2xl shadow-2xl max-w-sm w-full animate-[scaleIn_200ms_ease-out] overflow-hidden outline-none"
      >
        <div className="p-6">
          <div className="flex items-start gap-4">
            <div className={`flex-shrink-0 w-11 h-11 rounded-xl ${style.bg} ${style.text} flex items-center justify-center`}>
              <i className={`${style.icon} text-xl`} />
            </div>
            <div className="flex-1 min-w-0">
              <h3 className="font-heading font-bold text-slate-900 text-base">{title}</h3>
              <p className="mt-1.5 text-sm text-slate-500 leading-relaxed">{message}</p>
            </div>
          </div>
        </div>

        <div className="bg-slate-50 border-t border-slate-100 px-6 py-4 flex justify-end gap-2.5">
          <button
            type="button"
            onClick={handleCancel}
            className="px-4 py-2 bg-white border border-slate-200 rounded-xl text-sm font-semibold text-slate-600 hover:bg-slate-50 hover:border-slate-300 transition-colors shadow-sm"
          >
            Batal
          </button>
          <button
            type="button"
            onClick={onConfirm}
            className={`px-4 py-2 rounded-xl text-sm font-semibold text-white transition-colors shadow-sm focus:outline-none focus:ring-2 ${btnStyle}`}
          >
            Ya, Lanjutkan
          </button>
        </div>
      </div>

      <style>{`
        @keyframes fadeIn { from { opacity: 0; } to { opacity: 1; } }
        @keyframes scaleIn { from { opacity: 0; transform: scale(0.95); } to { opacity: 1; transform: scale(1); } }
      `}</style>
    </div>
  )
}
