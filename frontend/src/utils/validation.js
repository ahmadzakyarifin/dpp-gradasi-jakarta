export function isValidEmail(email) {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)
}

export function validateEmail(email, requiredMessage = 'Email wajib diisi') {
  if (!email) return requiredMessage
  if (!isValidEmail(email)) return 'Format email tidak valid'
  return ''
}

export function validatePassword(password, label = 'Password') {
  if (!password) return `${label} wajib diisi.`
  if (password.length < 6) return `${label} minimal 6 karakter.`
  return ''
}
