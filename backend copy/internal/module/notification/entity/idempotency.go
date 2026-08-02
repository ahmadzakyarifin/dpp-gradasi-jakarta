package entity

import "time"

type IdempotencyKey struct {
	ID              uint       `gorm:"primaryKey;autoIncrement"`
	Key             string     `gorm:"column:key;type:varchar(200);not null;uniqueIndex"`
	UserID          *uint      `gorm:"column:user_id"`
	RequestMethod   string     `gorm:"column:request_method;type:varchar(20);not null"`
	RequestPath     string     `gorm:"column:request_path;type:varchar(500);not null"`
	RequestHash     string     `gorm:"column:request_hash;type:char(64);not null"`
	ResponseStatus  *int       `gorm:"column:response_status"`
	ResponseBody    string     `gorm:"column:response_body;type:json"`
	ResponsePayload string     `gorm:"-" json:"-"`
	Status          string     `gorm:"column:status;type:enum('processing','completed','failed');not null;default:'processing'"`
	LockedUntil     *time.Time `gorm:"column:locked_until"`
	ExpiresAt       time.Time  `gorm:"column:expires_at;not null"`
	CreatedAt       time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(3)"`
	UpdatedAt       time.Time  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(3);onUpdate:CURRENT_TIMESTAMP(3)"`
}

const (
	IdempotencyStatusProcessing = "processing"
	IdempotencyStatusCompleted  = "completed"
	IdempotencyStatusFailed     = "failed"
)

func (IdempotencyKey) TableName() string {
	return "idempotency_keys"
}
