package jobs

type WhatsAppJob struct {
	NotificationID *uint  `json:"notification_id,omitempty"`
	ChatID         string `json:"chat_id"`
	Text           string `json:"text"`
}
