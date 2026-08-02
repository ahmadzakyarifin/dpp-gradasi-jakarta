package notification

// Channel notifikasi
const (
	ChannelEmail    = "email"
	ChannelWhatsApp = "whatsapp"
	ChannelSystem   = "system"
)

// Status notifikasi (cocok dengan enum di tabel notifications)
const (
	StatusPending   = "pending"
	StatusSent      = "sent"
	StatusFailed    = "failed"
	StatusCancelled = "cancelled"
)

// Event code notifikasi
const (
	EventAuthAccountActivation = "AUTH_ACCOUNT_ACTIVATION"
	EventAuthForgotPassword    = "AUTH_FORGOT_PASSWORD"
)

// EventToSubject mengembalikan subject email default berdasarkan event code.
func EventToSubject(eventCode string) string {
	switch eventCode {
	case EventAuthAccountActivation:
		return "Aktivasi Akun — SchoolPay"
	case EventAuthForgotPassword:
		return "Reset Password — SchoolPay"
	default:
		return "Notifikasi SchoolPay"
	}
}
