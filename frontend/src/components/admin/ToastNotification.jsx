import { useEffect } from 'react'

export default function ToastNotification({ show, message, type = 'success', onClose, duration = 3000 }) {
  useEffect(() => {
    if (show && duration > 0) {
      const timer = setTimeout(() => onClose?.(), duration)
      return () => clearTimeout(timer)
    }
  }, [show, duration, onClose])

  if (!show) return null

  const styles = {
    success: {
      bg: 'bg-emerald-600',
      icon: 'ph-bold ph-check-circle',
    },
    error: {
      bg: 'bg-red-600',
      icon: 'ph-bold ph-warning-circle',
    },
    info: {
      bg: 'bg-brand-600',
      icon: 'ph-bold ph-info',
    },
  }

  const style = styles[type] || styles.success

  return (
    <div className="fixed bottom-5 right-5 z-[70] animate-[slideUp_300ms_ease-out]">
      <div className={`${style.bg} text-white px-5 py-3.5 rounded-xl shadow-lg flex items-center gap-3 min-w-[280px] max-w-md`}>
        <i className={`${style.icon} text-xl flex-shrink-0`} />
        <span className="text-sm font-medium flex-1">{message}</span>
        <button
          onClick={onClose}
          className="flex-shrink-0 p-0.5 hover:bg-white/20 rounded-lg transition-colors"
        >
          <i className="ph-bold ph-x text-sm" />
        </button>
      </div>

      <style>{`
        @keyframes slideUp {
          from { opacity: 0; transform: translateY(16px); }
          to { opacity: 1; transform: translateY(0); }
        }
      `}</style>
    </div>
  )
}
