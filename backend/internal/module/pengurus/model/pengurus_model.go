package model

import (
	"time"

	"gorm.io/gorm"
)

type Pengurus struct {
	ID           uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string         `gorm:"size:150;not null" json:"name"`
	Role         string         `gorm:"size:200;not null" json:"role"`
	Department   *string        `gorm:"size:100" json:"department,omitempty"`
	Level        string         `gorm:"size:50;not null" json:"level"` // ketua, dpp, dpd, dpc
	Provinsi     *string        `gorm:"size:100" json:"provinsi,omitempty"`
	Kabupaten    *string        `gorm:"size:100" json:"kabupaten,omitempty"`
	ImageURL     string         `gorm:"size:500;not null" json:"image_url"`
	FacebookURL  *string        `gorm:"size:500" json:"facebook_url,omitempty"`
	InstagramURL *string        `gorm:"size:500" json:"instagram_url,omitempty"`
	LinkedinURL  *string        `gorm:"size:500" json:"linkedin_url,omitempty"`
	Whatsapp     *string        `gorm:"size:20" json:"whatsapp,omitempty"`
	Periode      string         `gorm:"size:50;not null" json:"periode"`
	SortOrder    int            `gorm:"default:0" json:"sort_order"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Pengurus) TableName() string {
	return "pengurus"
}
