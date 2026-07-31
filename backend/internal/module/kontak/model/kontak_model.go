package model

import (
	"time"

	"gorm.io/gorm"
)

type PesanKontak struct {
	ID           uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Nama         string         `gorm:"size:100;not null" json:"nama"`
	Email        string         `gorm:"size:150;not null" json:"email"`
	Subjek       string         `gorm:"size:200;not null" json:"subjek"`
	Pesan        string         `gorm:"type:text;not null" json:"pesan"`
	IsRead       bool           `gorm:"default:false" json:"is_read"`
	ResponseNote *string        `gorm:"type:text" json:"response_note,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (PesanKontak) TableName() string { return "pesan_kontak" }
