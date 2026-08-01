package helper

// contextKey adalah tipe unik untuk key di context.Context (hindari collision
// dengan key string biasa — aman dari staticcheck SA1029).
type contextKey string

// Keys context untuk data autentikasi. Didefinisikan di sini (bukan di
// middleware) supaya bisa dipakai helper & middleware tanpa circular import.
const (
	ContextUserID    contextKey = "user_id"
	ContextRoleID    contextKey = "role_id"
	ContextEmail     contextKey = "email"
	ContextRoleName  contextKey = "role_name"
	ContextUserName  contextKey = "user_name"
	ContextIPAddress contextKey = "ip_address"
	ContextUserAgent contextKey = "user_agent"
)
