package model

import (
	"time"

	"gorm.io/gorm"
)

type Berita struct {
	ID            uint   `gorm:"primaryKey;autoIncrement"`
	Slug          string `gorm:"uniqueIndex;size:250;not null"`
	Title         string `gorm:"size:300;not null"`
	Category      string `gorm:"size:100;default:'Berita Organisasi'"`
	PublishedDate string `gorm:"size:20;not null"`
	AuthorID      *uint
	AuthorName    string  `gorm:"->;-:migration"`
	ImagePath     *string `gorm:"size:500"`
	Excerpt       *string `gorm:"type:text"`
	Content       *string `gorm:"type:longtext;not null"`
	IsFeatured    bool    `gorm:"default:false"`
	IsNew         bool    `gorm:"default:false"`
	IsPublished   bool    `gorm:"default:false"`
	Views         int     `gorm:"default:0"`
	Footnote      *string `gorm:"size:500"`
	ImageSource   *string `gorm:"size:250"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	Tags          []BeritaTag    `gorm:"foreignKey:BeritaID"`
}

func (Berita) TableName() string {
	return "berita"
}

type BeritaTag struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	BeritaID uint   `gorm:"not null;index"`
	Tag      string `gorm:"size:100;not null"`
}

func (BeritaTag) TableName() string {
	return "berita_tags"
}
