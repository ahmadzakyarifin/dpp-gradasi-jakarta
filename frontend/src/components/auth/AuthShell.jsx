import AuthBrandPanel from './AuthBrandPanel'

export default function AuthShell({ children, variant = 'login' }) {
  return (
    <div className="font-sans antialiased bg-white hide-scrollbar">
      <div className="flex flex-col md:flex-row min-h-screen">
        <AuthBrandPanel variant={variant} />
        <div className="w-full md:w-1/2 lg:w-1/2 min-h-screen flex items-center justify-center p-6 sm:p-12 lg:p-16 bg-white relative shadow-[-10px_0_30px_-15px_rgba(0,0,0,0.1)] z-10">
          {children}
        </div>
      </div>
    </div>
  )
}
