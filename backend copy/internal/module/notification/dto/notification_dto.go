package dto

// SendEmailRequest adalah payload notifikasi email dengan konten HTML matang.
type SendEmailRequest struct {
	To      string
	Subject string
	HTML    string
	UserID  *uint
}

// SendWhatsAppRequest adalah payload notifikasi WhatsApp dengan teks matang.
type SendWhatsAppRequest struct {
	To     string
	Text   string
	UserID *uint
}

// NotificationResponse adalah response untuk endpoint list notifikasi (efikasi).
type NotificationResponse struct {
	ID        uint   `json:"id"`
	CreatedAt string `json:"createdAt"`
	Channel   string `json:"channel"`
	Recipient string `json:"recipient"`
	Contact   string `json:"contact"`
	Subject   string `json:"title"`
	Message   string `json:"message"`
	Error     string `json:"error,omitempty"`
	Status    string `json:"status"`
}

// PaginatedResponse wrapper untuk response pagination.
type PaginatedResponse struct {
	Data       any    `json:"data"`
	Total      int64  `json:"total"`
	Page       int    `json:"page"`
	PageSize   int    `json:"pageSize"`
	TotalPages int    `json:"totalPages"`
	Status     string `json:"status"`
	Message    string `json:"message"`
}
