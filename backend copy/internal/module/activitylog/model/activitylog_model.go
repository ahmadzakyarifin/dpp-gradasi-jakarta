package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

type JSONMap map[string]any

func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		str, ok := value.(string)
		if !ok {
			return errors.New("type assertion to []byte or string failed")
		}
		bytes = []byte(str)
	}
	return json.Unmarshal(bytes, j)
}

type ActivityLog struct {
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement"`

	ActorID   *uint  `gorm:"column:actor_id"`
	ActorName string `gorm:"column:actor_name"`
	ActorRole string `gorm:"column:actor_role"`

	Action string `gorm:"column:action"`

	EntityType  string `gorm:"column:entity_type"`
	EntityID    *uint  `gorm:"column:entity_id"`
	EntityLabel string `gorm:"column:entity_label"`

	RiskLevel string `gorm:"column:risk_level"`

	Description string `gorm:"column:description"`

	IPAddress string  `gorm:"column:ip_address"`
	UserAgent string  `gorm:"column:user_agent"`
	Metadata  JSONMap `gorm:"column:metadata;type:json"`

	CreatedAt time.Time `gorm:"column:created_at;not null;autoCreateTime"`
}

func (ActivityLog) TableName() string { return "activity_logs" }

func (a *ActivityLog) BeforeCreate(tx *gorm.DB) error {
	return nil
}
