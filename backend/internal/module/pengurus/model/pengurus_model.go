package model

import (
	"time"

	"gorm.io/gorm"
)

type Pengurus struct {
	ID           uint    `gorm:"primaryKey;autoIncrement"`
	Name         string  `gorm:"size:150;not null"`
	Role         string  `gorm:"size:200;not null"`
	Department   *string `gorm:"size:100"`
	Level        string  `gorm:"size:50;not null"` // ketua, dpp, dpd, dpc
	Provinsi     *string `gorm:"size:100"`
	Kabupaten    *string `gorm:"size:100"`
	ImagePath    string  `gorm:"size:500;not null"`
	FacebookURL  *string `gorm:"size:500"`
	InstagramURL *string `gorm:"size:500"`
	LinkedinURL  *string `gorm:"size:500"`
	Whatsapp     *string `gorm:"size:20"`
	Periode      string  `gorm:"size:50;not null"`
	SortOrder    int     `gorm:"default:0"`
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (Pengurus) TableName() string {
	return "pengurus"
}
