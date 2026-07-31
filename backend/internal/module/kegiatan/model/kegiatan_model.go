package model

import (
	"time"

	"gorm.io/gorm"
)

type Kegiatan struct {
	ID          uint              `gorm:"primaryKey;autoIncrement" json:"id"`
	Slug        string            `gorm:"uniqueIndex;size:250;not null" json:"slug"`
	Title       string            `gorm:"size:300;not null" json:"title"`
	Category    string            `gorm:"size:100;default:'Kegiatan'" json:"category"`
	EventDate   string            `gorm:"size:20" json:"event_date"`
	Location    string            `gorm:"size:200" json:"location"`
	Organizer   string            `gorm:"size:200" json:"organizer"`
	AuthorID    *uint             `json:"author_id,omitempty"`
	AuthorName  string            `gorm:"->;-:migration" json:"author_name,omitempty"`
	ImageURL    *string           `gorm:"size:500" json:"image_url,omitempty"`
	Excerpt     *string           `gorm:"type:text" json:"excerpt,omitempty"`
	Content     *string           `gorm:"type:longtext" json:"content,omitempty"`
	IsPublished bool              `gorm:"default:true" json:"is_published"`
	Views       int               `gorm:"default:0" json:"views"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	DeletedAt   gorm.DeletedAt    `gorm:"index" json:"-"`
	Tags        []KegiatanTag     `gorm:"foreignKey:KegiatanID" json:"tags,omitempty"`
	Gallery     []KegiatanGallery `gorm:"foreignKey:KegiatanID" json:"gallery,omitempty"`
}

func (Kegiatan) TableName() string { return "kegiatan" }

type KegiatanTag struct {
	ID         uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	KegiatanID uint   `gorm:"not null;index" json:"kegiatan_id"`
	Tag        string `gorm:"size:100;not null" json:"tag"`
}

func (KegiatanTag) TableName() string { return "kegiatan_tags" }

type KegiatanGallery struct {
	ID         uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	KegiatanID uint   `gorm:"not null;index" json:"kegiatan_id"`
	ImageURL   string `gorm:"size:500;not null" json:"image_url"`
	Caption    string `gorm:"size:200" json:"caption"`
	SortOrder  int    `gorm:"default:0" json:"sort_order"`
}

func (KegiatanGallery) TableName() string { return "kegiatan_gallery" }
