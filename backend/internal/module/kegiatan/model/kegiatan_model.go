package model

import (
	"time"

	"gorm.io/gorm"
)

type Kegiatan struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Slug        string `gorm:"uniqueIndex;size:250;not null"`
	Title       string `gorm:"size:300;not null"`
	Category    string `gorm:"size:100;default:'Kegiatan'"`
	EventDate   string `gorm:"size:200;not null"`
	Location    string `gorm:"size:200"`
	Organizer   string `gorm:"size:200"`
	AuthorID    *uint
	AuthorName  string  `gorm:"->;-:migration"`
	ImagePath   *string `gorm:"size:500"`
	Excerpt     *string `gorm:"type:text"`
	Content     *string `gorm:"type:longtext"`
	IsPublished bool
	IsNew       bool    `gorm:"default:false"`
	Views       int     `gorm:"default:0"`
	Footnote    *string `gorm:"size:500"`
	ImageSource *string `gorm:"size:250"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt    `gorm:"index"`
	Tags        []KegiatanTag     `gorm:"foreignKey:KegiatanID"`
	Gallery     []KegiatanGallery `gorm:"foreignKey:KegiatanID"`
}

func (Kegiatan) TableName() string { return "kegiatan" }

type KegiatanTag struct {
	ID         uint   `gorm:"primaryKey;autoIncrement"`
	KegiatanID uint   `gorm:"not null;index"`
	Tag        string `gorm:"size:100;not null"`
}

func (KegiatanTag) TableName() string { return "kegiatan_tags" }

type KegiatanGallery struct {
	ID         uint   `gorm:"primaryKey;autoIncrement"`
	KegiatanID uint   `gorm:"not null;index"`
	ImagePath  string `gorm:"size:500;not null"`
	Caption    string `gorm:"size:200"`
	SortOrder  int    `gorm:"default:0"`
}

func (KegiatanGallery) TableName() string { return "kegiatan_gallery" }
