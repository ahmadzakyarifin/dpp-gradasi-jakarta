package repository

import (
	"fmt"
	"strings"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/berita/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/berita/entity"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/berita/mapper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/berita/model"
	"gorm.io/gorm"
)

type BeritaRepo interface {
	FindPublished(query dto.BeritaQuery) ([]entity.Berita, int64, error)
	FindAll(query dto.BeritaQuery) ([]entity.Berita, int64, error)
	FindBySlug(slug string) (*entity.Berita, error)
	FindUniqueSlug(slug string) (string, error)
	ExistsTitle(title string, excludeID uint) (bool, error)
	FindByID(id uint) (*entity.Berita, error)
	Create(berita *entity.Berita) error
	Update(berita *entity.Berita) error
	SoftDelete(id uint) error
	Restore(id uint) error
	BulkSoftDelete(ids []uint) error
	BulkRestore(ids []uint) error
	IncrementViews(id uint) error
	SaveTags(beritaID uint, tags []string) error
	GetCategories() ([]string, error)
}

type beritaRepo struct {
	db *gorm.DB
}

func NewBeritaRepo(db *gorm.DB) BeritaRepo {
	return &beritaRepo{db: db}
}

func (r *beritaRepo) FindPublished(query dto.BeritaQuery) ([]entity.Berita, int64, error) {
	return r.query(true, query)
}

func (r *beritaRepo) FindAll(query dto.BeritaQuery) ([]entity.Berita, int64, error) {
	return r.query(false, query)
}

func (r *beritaRepo) query(publishedOnly bool, q dto.BeritaQuery) ([]entity.Berita, int64, error) {
	var beritas []model.Berita
	var total int64

	// Session baru tanpa auto soft-delete: handle deleted_at manual (JOIN users bikin ambigu)
	db := r.db.Session(&gorm.Session{}).Model(&model.Berita{}).
		Select("berita.*, users.name as author_name").
		Joins("LEFT JOIN users ON users.id = berita.author_id")

	if publishedOnly {
		db = db.Where("is_published = ?", true)
	}

	// Filter status: published | draft | trashed | all (admin mode)
	if !publishedOnly {
		switch q.Status {
		case "published":
			db = db.Where("is_published = ? AND berita.deleted_at IS NULL", true)
		case "draft":
			db = db.Where("is_published = ? AND berita.deleted_at IS NULL", false)
		case "trashed":
			db = db.Unscoped().Where("berita.deleted_at IS NOT NULL")
		default: // "all" atau kosong
			db = db.Where("berita.deleted_at IS NULL")
		}
	} else if q.Status == "trashed" {
		db = db.Unscoped().Where("berita.deleted_at IS NOT NULL")
	} else {
		db = db.Where("berita.deleted_at IS NULL")
	}

	if q.Search != "" {
		search := strings.TrimSpace(q.Search)
		db = db.Where("MATCH(title, content) AGAINST (? IN BOOLEAN MODE)", search+"*")
	}

	if q.Category != "" {
		db = db.Where("category = ?", q.Category)
	}

	if q.Tag != "" {
		db = db.Where("EXISTS (SELECT 1 FROM berita_tags WHERE berita_tags.berita_id = berita.id AND berita_tags.tag = ?)", q.Tag)
	}

	// Count total
	db.Count(&total)

	// Sort
	switch q.Sort {
	case "oldest":
		db = db.Order("published_date ASC, created_at ASC")
	case "most_viewed":
		db = db.Order("views DESC, published_date DESC")
	default:
		db = db.Order("published_date DESC, created_at DESC")
	}

	// Pagination
	page := q.Page
	limit := q.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	err := db.Preload("Tags").Limit(limit).Offset(offset).Find(&beritas).Error
	return mapper.ModelListToEntity(beritas), total, err
}

func (r *beritaRepo) FindBySlug(slug string) (*entity.Berita, error) {
	var b model.Berita
	err := r.db.Model(&model.Berita{}).
		Select("berita.*, users.name as author_name").
		Joins("LEFT JOIN users ON users.id = berita.author_id").
		Where("berita.slug = ? AND berita.is_published = ?", slug, true).
		Preload("Tags").First(&b).Error
	if err != nil {
		return nil, err
	}
	return mapper.ModelToEntity(&b), nil
}

// FindUniqueSlug mengembalikan slug unik dengan suffix -2, -3, dst jika slug sudah dipakai.
// Memakai Unscoped agar slug di soft-delete (trash) juga dianggap terpakai.
func (r *beritaRepo) FindUniqueSlug(slug string) (string, error) {
	base := slug
	candidate := slug
	for i := 2; ; i++ {
		var count int64
		err := r.db.Unscoped().Model(&model.Berita{}).Where("slug = ?", candidate).Count(&count).Error
		if err != nil {
			return "", err
		}
		if count == 0 {
			return candidate, nil
		}
		candidate = fmt.Sprintf("%s-%d", base, i)
	}
}

// ExistsTitle mengembalikan true jika judul sudah dipakai berita lain (soft-delete ikut dihitung).
// excludeID dipakai saat update — mengabaikan berita dengan id tersebut (dirinya sendiri).
func (r *beritaRepo) ExistsTitle(title string, excludeID uint) (bool, error) {
	var count int64
	q := r.db.Unscoped().Model(&model.Berita{}).Where("title = ?", title)
	if excludeID > 0 {
		q = q.Where("id != ?", excludeID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *beritaRepo) FindByID(id uint) (*entity.Berita, error) {
	var b model.Berita
	err := r.db.Model(&model.Berita{}).
		Select("berita.*, users.name as author_name").
		Joins("LEFT JOIN users ON users.id = berita.author_id").
		Preload("Tags").First(&b, id).Error
	if err != nil {
		return nil, err
	}
	return mapper.ModelToEntity(&b), nil
}

func (r *beritaRepo) Create(berita *entity.Berita) error {
	return r.db.Create(mapper.EntityToModel(berita)).Error
}

func (r *beritaRepo) Update(berita *entity.Berita) error {
	return r.db.Save(mapper.EntityToModel(berita)).Error
}

func (r *beritaRepo) SoftDelete(id uint) error {
	return r.db.Delete(&model.Berita{}, id).Error
}

func (r *beritaRepo) Restore(id uint) error {
	return r.db.Unscoped().Model(&model.Berita{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *beritaRepo) BulkSoftDelete(ids []uint) error {
	return r.db.Where("id IN ?", ids).Delete(&model.Berita{}).Error
}

func (r *beritaRepo) BulkRestore(ids []uint) error {
	return r.db.Unscoped().Model(&model.Berita{}).Where("id IN ?", ids).Update("deleted_at", nil).Error
}

func (r *beritaRepo) IncrementViews(id uint) error {
	return r.db.Model(&model.Berita{}).Where("id = ?", id).
		UpdateColumn("views", gorm.Expr("views + 1")).Error
}

func (r *beritaRepo) SaveTags(beritaID uint, tags []string) error {
	// Delete existing tags
	r.db.Where("berita_id = ?", beritaID).Delete(&model.BeritaTag{})

	// Insert new tags
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		r.db.Create(&model.BeritaTag{
			BeritaID: beritaID,
			Tag:      tag,
		})
	}
	return nil
}

// GetCategories mengembalikan daftar kategori unik (non-kosong) dari data berita aktif.
func (r *beritaRepo) GetCategories() ([]string, error) {
	var categories []string
	err := r.db.Model(&model.Berita{}).
		Where("deleted_at IS NULL AND category IS NOT NULL AND category != ''").
		Distinct().
		Order("category ASC").
		Pluck("category", &categories).Error
	return categories, err
}
