import AuthBrandPanel from './AuthBrandPanel'

export default function AuthShell({ children, variant = 'login' }) {
  return (
    <div className="font-sans antialiased bg-white hide-scrollbar animate-fade-in-up duration-300">
      <div className="flex flex-col md:flex-row min-h-screen">
        <AuthBrandPanel variant={variant} />
        <div className="w-full md:w-1/2 lg:w-1/2 min-h-screen flex items-center justify-center p-6 sm:p-12 lg:p-16 bg-gradient-to-tr from-slate-50/50 via-white to-slate-50/20 relative shadow-[-10px_0_30px_-15px_rgba(0,0,0,0.05)] z-10 overflow-hidden">
          {/* Soft decorative background glows */}
          <div className="absolute top-0 right-0 w-80 h-80 bg-brand-500/5 blur-[120px] rounded-full pointer-events-none -mr-20 -mt-20 z-0" />
          <div className="absolute bottom-0 left-0 w-80 h-80 bg-emerald-500/5 blur-[120px] rounded-full pointer-events-none -ml-20 -mb-20 z-0" />
          
          <div className="w-full max-w-md relative z-10 animate-fade-in-up" style={{ animationDelay: '100ms' }}>
            {children}
          </div>
        </div>
      </div>
    </div>
  )
}
