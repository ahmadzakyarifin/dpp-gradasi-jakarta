import { Link } from 'react-router-dom'
import { authContent } from '../../content/authContent'
import { useSettings } from '../../context/useSettings'
import { resolveAssetUrl } from '../../utils/assetUrl'

export default function AuthBrandPanel({ variant = 'login' }) {
  const { settings } = useSettings()
  const data = variant === 'reset' ? authContent.reset : authContent.login

  const logoUrl = settings?.logo_path ? resolveAssetUrl(settings.logo_path) : ''
  const heroTitle = variant === 'reset'
    ? (settings?.site_name || '')
    : (settings?.login_hero_title || '')

  const heroDesc = variant === 'reset'
    ? data.heroDescription
    : (settings?.login_hero_description || '')

  const socialLinks = [
    { label: 'Facebook', icon: 'ph-facebook-logo', url: settings?.facebook_url || '#' },
    { label: 'Instagram', icon: 'ph-instagram-logo', url: settings?.instagram_url || '#' },
    { label: 'YouTube', icon: 'ph-youtube-logo', url: settings?.youtube_url || '#' },
  ]

  return (
    <div className="hidden md:flex md:w-1/2 lg:w-1/2 relative bg-brand-950 overflow-hidden flex-col justify-between p-10 lg:p-20 xl:p-24">
      <div className="absolute inset-0 z-0">
        <img src={authContent.backgroundUrl} alt="Background" className="w-full h-full object-cover opacity-30 mix-blend-overlay" />
        <div className="absolute inset-0 bg-gradient-to-br from-brand-950/95 via-brand-900/90 to-brand-800/80" />
        <div className="absolute inset-0 bg-texture-dots opacity-40" />
      </div>

      <div className="relative z-10">
        <Link to={authContent.homePath} className="inline-flex items-center gap-2 text-white/80 hover:text-white transition-all duration-300 ease-in-out group bg-white/5 hover:bg-white/10 px-5 py-2.5 rounded-full backdrop-blur-md border border-white/10 shadow-lg w-fit mb-24 lg:mb-32">
          <i className="ph-bold ph-arrow-left group-hover:-translate-x-1 transition-transform duration-300 ease-in-out" />
          <span className="text-sm font-semibold tracking-wide">Kembali ke Beranda</span>
        </Link>

        {logoUrl && (
          <div className="w-20 h-20 bg-white rounded-2xl p-3 shadow-2xl mb-8 transform transition hover:scale-105">
            <img src={logoUrl} alt="Logo" className="w-full h-full object-contain" />
          </div>
        )}

        <h1 className="font-heading text-4xl lg:text-5xl xl:text-6xl font-extrabold text-white leading-tight mb-5 tracking-tight">
          {heroTitle}
        </h1>

        <p className="text-lg text-brand-100/80 max-w-md font-light leading-relaxed">
          {heroDesc}
        </p>
      </div>

      <div className="relative z-10 flex items-center justify-between mt-12">
        <p className="text-white/40 text-sm font-medium">© 2026 {settings?.site_name || ''}.</p>
        <div className="flex gap-4">
          {socialLinks.map((item) => (
            <a key={item.label} href={item.url} target="_blank" rel="noreferrer" aria-label={item.label} className="w-10 h-10 rounded-full bg-white/5 border border-white/10 flex items-center justify-center text-white/60 hover:text-white hover:bg-white/20 hover:border-white/30 transition backdrop-blur-sm">
              <i className={`ph-fill ${item.icon} text-lg`} />
            </a>
          ))}
        </div>
      </div>
    </div>
  )
}
