package model

import (
	"time"

	"gorm.io/gorm"
)

type Slider struct {
	ID        uint    `gorm:"primaryKey;autoIncrement"`
	Title     string  `gorm:"size:200;not null"`
	Subtitle  *string `gorm:"type:text"`
	Tag       *string `gorm:"size:50"`
	IsNew     bool    `gorm:"default:false"`
	EventDate *string `gorm:"size:100"`
	Location  *string `gorm:"size:200"`
	ImagePath string  `gorm:"size:500;not null"`
	SortOrder int     `gorm:"default:0"`
	IsPublished bool    `gorm:"default:true;column:is_published"`
	CreatedBy *uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Slider) TableName() string {
	return "sliders"
}
