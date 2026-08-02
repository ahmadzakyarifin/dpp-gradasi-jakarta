package model

import "time"

// Notification mencerminkan baris di tabel notifications.
type Notification struct {
	ID                uint       `gorm:"primaryKey;autoIncrement"`
	EventCode         string     `gorm:"column:event_code;type:varchar(100);not null;default:''"`
	StudentID         *uint      `gorm:"column:student_id"`
	GuardianID        *uint      `gorm:"column:guardian_id"`
	RecipientUserID   *uint      `gorm:"column:recipient_user_id"`
	Channel           string     `gorm:"column:channel;type:enum('whatsapp','email','system');not null"`
	Destination       string     `gorm:"column:destination;type:varchar(200);not null"`
	Subject           string     `gorm:"column:subject;type:varchar(200)"`
	Message           string     `gorm:"column:message;type:text;not null"`
	Payload           *string    `gorm:"column:payload;type:json"`
	Status            string     `gorm:"column:status;type:enum('pending','sent','failed','cancelled');not null;default:'pending'"`
	ProviderMessageID string     `gorm:"column:provider_message_id;type:varchar(200)"`
	ProviderError     string     `gorm:"column:provider_error;type:text"`
	ErrorMessage      string     `gorm:"column:error_message;type:text"`
	IsRead            bool       `gorm:"column:is_read;type:tinyint(1);not null;default:0"`
	ReadAt            *time.Time `gorm:"column:read_at"`
	SentAt            *time.Time `gorm:"column:sent_at"`
	CreatedAt         time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(3)"`
	UpdatedAt         time.Time  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(3);onUpdate:CURRENT_TIMESTAMP(3)"`
}

func (Notification) TableName() string {
	return "notifications"
}

type NotificationRow struct {
	Notification
	RecipientName string `gorm:"column:recipient_name"`
}
