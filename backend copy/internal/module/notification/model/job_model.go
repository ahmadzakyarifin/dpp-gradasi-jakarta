package model

import "time"

// BackgroundJobModel memetakan tabel background_jobs (antrian job asynq/worker).
// Menggunakan GORM agar konsisten dengan modul lain (menggantikan entity bun-based).
type BackgroundJobModel struct {
	ID           uint       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	JobType      string     `gorm:"column:job_type;not null" json:"job_type"`
	QueueName    string     `gorm:"column:queue_name;not null;default:default" json:"queue_name"`
	Payload      string     `gorm:"column:payload;type:JSON;not null" json:"payload"`
	Status       string     `gorm:"column:status;not null;default:queued" json:"status"`
	Priority     int        `gorm:"column:priority;not null;default:0" json:"priority"`
	AttemptCount int        `gorm:"column:attempt_count;not null;default:0" json:"attempt_count"`
	MaxAttempts  int        `gorm:"column:max_attempts;not null;default:3" json:"max_attempts"`
	LockedUntil  *time.Time `gorm:"column:locked_until" json:"locked_until,omitempty"`
	AvailableAt  time.Time  `gorm:"column:available_at;not null;default:CURRENT_TIMESTAMP" json:"available_at"`
	ProcessedAt  *time.Time `gorm:"column:processed_at" json:"processed_at,omitempty"`
	FailedAt     *time.Time `gorm:"column:failed_at" json:"failed_at,omitempty"`
	ErrorMessage *string    `gorm:"column:error_message" json:"error_message,omitempty"`
	CreatedAt    time.Time  `gorm:"column:created_at;not null;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;not null;autoUpdateTime" json:"updated_at"`
}

func (BackgroundJobModel) TableName() string { return "background_jobs" }
