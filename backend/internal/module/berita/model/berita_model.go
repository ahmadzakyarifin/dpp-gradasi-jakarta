package model

import (
	"time"

	"gorm.io/gorm"
)

type Berita struct {
	ID            uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Slug          string         `gorm:"uniqueIndex;size:250;not null" json:"slug"`
	Title         string         `gorm:"size:300;not null" json:"title"`
	Category      string         `gorm:"size:100;default:'Berita Organisasi'" json:"category"`
	PublishedDate string         `gorm:"size:20;not null" json:"published_date"`
	AuthorID      *uint          `json:"author_id,omitempty"`
	AuthorName    string         `gorm:"->;-:migration" json:"author_name,omitempty"`
	ImagePath     *string        `gorm:"size:500" json:"image_path,omitempty"`
	Excerpt       *string        `gorm:"type:text" json:"excerpt,omitempty"`
	Content       *string        `gorm:"type:longtext" json:"content,omitempty"`
	IsFeatured    bool           `gorm:"default:false" json:"is_featured"`
	IsPublished   bool           `gorm:"default:true" json:"is_published"`
	Views         int            `gorm:"default:0" json:"views"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Tags          []BeritaTag    `gorm:"foreignKey:BeritaID" json:"tags,omitempty"`
}

func (Berita) TableName() string {
	return "berita"
}

type BeritaTag struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	BeritaID uint   `gorm:"not null;index" json:"berita_id"`
	Tag      string `gorm:"size:100;not null" json:"tag"`
}

func (BeritaTag) TableName() string {
	return "berita_tags"
}
