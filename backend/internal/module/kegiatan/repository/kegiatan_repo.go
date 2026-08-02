package repository

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kegiatan/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kegiatan/model"
	"gorm.io/gorm"
)

type KegiatanRepo interface {
	FindPublished(q dto.KegiatanQuery) ([]model.Kegiatan, int64, error)
	FindAll(q dto.KegiatanQuery) ([]model.Kegiatan, int64, error)
	FindBySlug(slug string) (*model.Kegiatan, error)
	FindUniqueSlug(slug string) (string, error)
	ExistsTitle(title string, excludeID uint) (bool, error)
	FindByID(id uint) (*model.Kegiatan, error)
	Create(k *model.Kegiatan) error
	Update(k *model.Kegiatan) error
	SoftDelete(id uint) error
	Restore(id uint) error
	BulkSoftDelete(ids []uint) error
	BulkRestore(ids []uint) error
	IncrementViews(id uint) error
	SaveTags(kegiatanID uint, tags []string) error
	SaveGallery(kegiatanID uint, items []dto.GalleryInput) error
	DeleteGalleryImage(galleryID uint) error
	GetCategories() ([]string, error)
}

type kegiatanRepo struct{ db *gorm.DB }

func NewKegiatanRepo(db *gorm.DB) KegiatanRepo {
	return &kegiatanRepo{db: db}
}

func (r *kegiatanRepo) FindPublished(q dto.KegiatanQuery) ([]model.Kegiatan, int64, error) {
	return r.query(true, q)
}

func (r *kegiatanRepo) FindAll(q dto.KegiatanQuery) ([]model.Kegiatan, int64, error) {
	return r.query(false, q)
}

func (r *kegiatanRepo) query(publishedOnly bool, q dto.KegiatanQuery) ([]model.Kegiatan, int64, error) {
	var items []model.Kegiatan
	var total int64

	// Session baru tanpa auto soft-delete: handle deleted_at manual (JOIN users bikin ambigu)
	db := r.db.Session(&gorm.Session{}).Model(&model.Kegiatan{}).
		Select("kegiatan.*, users.name as author_name").
		Joins("LEFT JOIN users ON users.id = kegiatan.author_id")

	if publishedOnly {
		db = db.Where("is_published = ?", true)
	}

	// Filter status: published | draft | trashed | all (admin mode)
	if !publishedOnly {
		switch q.Status {
		case "published":
			db = db.Where("is_published = ? AND kegiatan.deleted_at IS NULL", true)
		case "draft":
			db = db.Where("is_published = ? AND kegiatan.deleted_at IS NULL", false)
		case "trashed":
			db = db.Unscoped().Model(&model.Kegiatan{}).Where("kegiatan.deleted_at IS NOT NULL")
		default: // "all" atau kosong
			db = db.Where("kegiatan.deleted_at IS NULL")
		}
	} else if q.Status == "trashed" {
		db = db.Unscoped().Model(&model.Kegiatan{}).Where("kegiatan.deleted_at IS NOT NULL")
	} else {
		db = db.Where("kegiatan.deleted_at IS NULL")
	}

	if q.Search != "" {
		s := "%" + strings.TrimSpace(q.Search) + "%"
		db = db.Where("title LIKE ? OR content LIKE ?", s, s)
	}
	if q.Category != "" {
		db = db.Where("category = ?", q.Category)
	}

	db.Count(&total)

	switch q.Sort {
	case "oldest":
		db = db.Order("event_date ASC, created_at ASC")
	case "most_viewed":
		db = db.Order("views DESC, event_date DESC")
	default:
		db = db.Order("event_date DESC, created_at DESC")
	}

	page := maxInt(q.Page, 1)
	limit := q.Limit
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	err := db.Preload("Tags").Limit(limit).Offset(offset).Find(&items).Error
	return items, total, err
}

func (r *kegiatanRepo) FindBySlug(slug string) (*model.Kegiatan, error) {
	var k model.Kegiatan
	err := r.db.Model(&model.Kegiatan{}).
		Select("kegiatan.*, users.name as author_name").
		Joins("LEFT JOIN users ON users.id = kegiatan.author_id").
		Where("kegiatan.slug = ? AND kegiatan.deleted_at IS NULL AND kegiatan.is_published = ?", slug, true).
		Preload("Tags").Preload("Gallery", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC, id ASC")
	}).First(&k).Error
	if err != nil {
		return nil, err
	}
	return &k, nil
}

// FindUniqueSlug mengembalikan slug unik dengan suffix -2, -3, dst jika slug sudah dipakai.
// Memakai Unscoped agar slug di soft-delete (trash) juga dianggap terpakai.
func (r *kegiatanRepo) FindUniqueSlug(slug string) (string, error) {
	base := slug
	candidate := slug
	for i := 2; ; i++ {
		var count int64
		err := r.db.Unscoped().Model(&model.Kegiatan{}).Where("slug = ?", candidate).Count(&count).Error
		if err != nil {
			return "", err
		}
		if count == 0 {
			return candidate, nil
		}
		candidate = fmt.Sprintf("%s-%d", base, i)
	}
}

// ExistsTitle mengembalikan true jika judul sudah dipakai kegiatan lain (soft-delete ikut dihitung).
// excludeID dipakai saat update — mengabaikan kegiatan dengan id tersebut (dirinya sendiri).
func (r *kegiatanRepo) ExistsTitle(title string, excludeID uint) (bool, error) {
	var count int64
	q := r.db.Unscoped().Model(&model.Kegiatan{}).Where("title = ?", title)
	if excludeID > 0 {
		q = q.Where("id != ?", excludeID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *kegiatanRepo) FindByID(id uint) (*model.Kegiatan, error) {
	var k model.Kegiatan
	err := r.db.Model(&model.Kegiatan{}).
		Select("kegiatan.*, users.name as author_name").
		Joins("LEFT JOIN users ON users.id = kegiatan.author_id").
		Preload("Tags").Preload("Gallery", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC, id ASC")
	}).First(&k, id).Error
	if err != nil {
		return nil, err
	}
	return &k, nil
}

func (r *kegiatanRepo) Create(k *model.Kegiatan) error {
	return r.db.Create(k).Error
}

func (r *kegiatanRepo) Update(k *model.Kegiatan) error {
	return r.db.Save(k).Error
}

func (r *kegiatanRepo) SoftDelete(id uint) error {
	return r.db.Delete(&model.Kegiatan{}, id).Error
}

func (r *kegiatanRepo) Restore(id uint) error {
	return r.db.Unscoped().Model(&model.Kegiatan{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *kegiatanRepo) BulkSoftDelete(ids []uint) error {
	return r.db.Where("id IN ?", ids).Delete(&model.Kegiatan{}).Error
}

func (r *kegiatanRepo) BulkRestore(ids []uint) error {
	return r.db.Unscoped().Model(&model.Kegiatan{}).Where("id IN ?", ids).Update("deleted_at", nil).Error
}

func (r *kegiatanRepo) IncrementViews(id uint) error {
	return r.db.Model(&model.Kegiatan{}).Where("id = ?", id).
		UpdateColumn("views", gorm.Expr("views + 1")).Error
}

func (r *kegiatanRepo) SaveTags(kegiatanID uint, tags []string) error {
	r.db.Where("kegiatan_id = ?", kegiatanID).Delete(&model.KegiatanTag{})
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		r.db.Create(&model.KegiatanTag{KegiatanID: kegiatanID, Tag: tag})
	}
	return nil
}

func (r *kegiatanRepo) SaveGallery(kegiatanID uint, items []dto.GalleryInput) error {
	r.db.Where("kegiatan_id = ?", kegiatanID).Delete(&model.KegiatanGallery{})
	for _, item := range items {
		if item.ImagePath == "" {
			continue
		}
		r.db.Create(&model.KegiatanGallery{
			KegiatanID: kegiatanID,
			ImagePath:  item.ImagePath,
			Caption:    item.Caption,
			SortOrder:  item.SortOrder,
		})
	}
	return nil
}

func (r *kegiatanRepo) DeleteGalleryImage(galleryID uint) error {
	return r.db.Delete(&model.KegiatanGallery{}, galleryID).Error
}

// GetCategories mengembalikan daftar kategori unik (non-kosong) dari data kegiatan aktif.
func (r *kegiatanRepo) GetCategories() ([]string, error) {
	var categories []string
	err := r.db.Model(&model.Kegiatan{}).
		Where("deleted_at IS NULL AND category IS NOT NULL AND category != ''").
		Distinct().
		Order("category ASC").
		Pluck("category", &categories).Error
	return categories, err
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ParseGalleryJSON parses gallery JSON string into GalleryInput slice
func ParseGalleryJSON(s string) ([]dto.GalleryInput, error) {
	if s == "" {
		return nil, nil
	}
	var items []dto.GalleryInput
	if err := json.Unmarshal([]byte(s), &items); err != nil {
		return nil, err
	}
	return items, nil
}
