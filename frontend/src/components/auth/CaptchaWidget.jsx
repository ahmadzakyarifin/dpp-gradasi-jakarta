import { useState, useEffect, useRef } from 'react'

export default function CaptchaWidget({ onVerify, hasError }) {
  const [captchaCode, setCaptchaCode] = useState('')
  const [userInput, setUserInput] = useState('')
  const canvasRef = useRef(null)

  // Generate random captcha string
  const generateCaptcha = () => {
    const chars = '23456789ABCDEFGHJKLMNPQRSTUVWXYZ'
    let code = ''
    for (let i = 0; i < 5; i++) {
      code += chars.charAt(Math.floor(Math.random() * chars.length))
    }
    setCaptchaCode(code)
    setUserInput('')
    if (onVerify) onVerify('')
  }

  // Draw captcha image on HTML5 canvas for security and aesthetics
  const drawCaptcha = (code) => {
    const canvas = canvasRef.current
    if (!canvas) return
    const ctx = canvas.getContext('2d')
    ctx.clearRect(0, 0, canvas.width, canvas.height)

    // Background gradient
    const grad = ctx.createLinearGradient(0, 0, canvas.width, canvas.height)
    grad.addColorStop(0, '#f1f5f9')
    grad.addColorStop(1, '#e2e8f0')
    ctx.fillStyle = grad
    ctx.fillRect(0, 0, canvas.width, canvas.height)

    // Draw noise lines
    for (let i = 0; i < 4; i++) {
      ctx.strokeStyle = `rgba(${Math.random() * 100}, ${Math.random() * 100}, ${Math.random() * 200}, 0.3)`
      ctx.lineWidth = 1 + Math.random() * 2
      ctx.beginPath()
      ctx.moveTo(Math.random() * canvas.width, Math.random() * canvas.height)
      ctx.lineTo(Math.random() * canvas.width, Math.random() * canvas.height)
      ctx.stroke()
    }

    // Draw characters with random rotation and colors
    ctx.font = 'bold 22px "Plus Jakarta Sans", sans-serif'
    ctx.textBaseline = 'middle'

    for (let i = 0; i < code.length; i++) {
      const char = code[i]
      const x = 16 + i * 22
      const y = canvas.height / 2 + (Math.random() - 0.5) * 6
      const angle = (Math.random() - 0.5) * 0.4

      ctx.save()
      ctx.translate(x, y)
      ctx.rotate(angle)
      ctx.fillStyle = ['#1e3a8a', '#2563eb', '#0284c7', '#0d9488', '#4f46e5'][i % 5]
      ctx.fillText(char, 0, 0)
      ctx.restore()
    }
  }

  useEffect(() => {
    generateCaptcha()
  }, [])

  useEffect(() => {
    if (captchaCode) {
      drawCaptcha(captchaCode)
    }
  }, [captchaCode])

  const handleInputChange = (e) => {
    const val = e.target.value.toUpperCase().trim()
    setUserInput(val)
    if (val === captchaCode) {
      if (onVerify) onVerify(`TOKEN_VERIFIED_${captchaCode}`)
    } else {
      if (onVerify) onVerify('')
    }
  }

  return (
    <div className="space-y-2">
      <label className="block text-xs font-semibold text-slate-600 uppercase tracking-wider">
        Verifikasi Keamanan (CAPTCHA)
      </label>
      
      <div className="flex items-center gap-3">
        <div className="relative border border-slate-300 rounded-xl overflow-hidden shadow-inner bg-slate-100 shrink-0">
          <canvas ref={canvasRef} width={130} height={42} className="block select-none" />
        </div>

        <button
          type="button"
          onClick={generateCaptcha}
          className="p-2 text-slate-500 hover:text-brand-600 hover:bg-slate-100 rounded-xl transition border border-slate-200 shadow-sm shrink-0"
          title="Acak CAPTCHA"
        >
          <i className="ph-bold ph-arrows-counter-clockwise text-lg" />
        </button>

        <div className="relative flex-1">
          <input
            type="text"
            maxLength={5}
            value={userInput}
            onChange={handleInputChange}
            placeholder="Kode CAPTCHA"
            className={`w-full px-3 py-2.5 bg-slate-50 border rounded-xl focus:outline-none text-sm font-bold uppercase tracking-widest ${
              hasError
                ? 'border-red-400 focus:ring-2 focus:ring-red-500/20 focus:border-red-500 text-red-900 placeholder-slate-400'
                : userInput === captchaCode && captchaCode !== ''
                ? 'border-emerald-500 focus:ring-2 focus:ring-emerald-500/20 text-emerald-700 bg-emerald-50/30'
                : 'border-slate-200 focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 text-slate-700 placeholder-slate-400'
            }`}
          />
          {userInput === captchaCode && captchaCode !== '' && (
            <i className="ph-fill ph-check-circle absolute right-3 top-1/2 -translate-y-1/2 text-emerald-500 text-lg" />
          )}
        </div>
      </div>
    </div>
  )
}
