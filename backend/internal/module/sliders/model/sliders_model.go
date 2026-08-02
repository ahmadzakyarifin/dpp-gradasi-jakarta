package model

import (
	"time"

	"gorm.io/gorm"
)

type Slider struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Title     string         `gorm:"size:200;not null" json:"title"`
	Subtitle  *string        `gorm:"type:text" json:"subtitle,omitempty"`
	Tag       *string        `gorm:"size:50" json:"tag,omitempty"`
	IsNew     bool           `gorm:"default:false" json:"is_new"`
	EventDate *string        `gorm:"size:100" json:"event_date,omitempty"`
	Location  *string        `gorm:"size:200" json:"location,omitempty"`
	ImagePath string         `gorm:"size:500;not null" json:"image_path"`
	LinkURL   *string        `gorm:"size:500" json:"link_url,omitempty"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedBy *uint          `json:"created_by,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Slider) TableName() string {
	return "sliders"
}
