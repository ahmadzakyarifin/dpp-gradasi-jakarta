import { useEffect, useRef } from 'react'

/**
 * useReveal — scroll-reveal hook (AOS-like, tanpa library).
 * Tambahkan className "reveal" (+ varian "reveal-left|right|zoom") ke elemen,
 * lalu pasang ref dari hook ini. Elemen muncul saat masuk viewport.
 *
 * Usage:
 *   const revealRef = useReveal()
 *   <div ref={revealRef} className="reveal">...</div>
 */
export default function useReveal() {
  const ref = useRef(null)

  useEffect(() => {
    const el = ref.current
    if (!el) return

    // Respect reduced motion — langsung tampil
    if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
      el.classList.add('is-visible')
      return
    }

    let rafId = null
    const check = () => {
      const rect = el.getBoundingClientRect()
      const vh = window.innerHeight || document.documentElement.clientHeight
      // Elemen dianggap visible jika bagian atasnya sudah melewati 85% viewport
      if (rect.top < vh * 0.88 && rect.bottom > 0) {
        el.classList.add('is-visible')
        window.removeEventListener('scroll', onScroll, { passive: true })
        window.removeEventListener('resize', onScroll, { passive: true })
      }
    }
    const onScroll = () => {
      if (rafId) return
      rafId = requestAnimationFrame(() => {
        rafId = null
        check()
      })
    }

    // Cek awal (kalau elemen sudah di viewport saat mount)
    check()
    window.addEventListener('scroll', onScroll, { passive: true })
    window.addEventListener('resize', onScroll, { passive: true })

    return () => {
      window.removeEventListener('scroll', onScroll, { passive: true })
      window.removeEventListener('resize', onScroll, { passive: true })
      if (rafId) cancelAnimationFrame(rafId)
    }
  }, [])

  return ref
}
