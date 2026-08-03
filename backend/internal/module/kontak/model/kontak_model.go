package model

import (
	"time"

	"gorm.io/gorm"
)

type PesanKontak struct {
	ID           uint    `gorm:"primaryKey;autoIncrement"`
	Nama         string  `gorm:"size:100;not null"`
	Email        string  `gorm:"size:150;not null"`
	Subjek       string  `gorm:"size:200;not null"`
	Pesan        string  `gorm:"type:text;not null"`
	IsRead       bool    `gorm:"default:false"`
	ResponseNote *string `gorm:"type:text"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (PesanKontak) TableName() string { return "pesan_kontak" }
