package repository

import (
	"strings"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kontak/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kontak/entity"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kontak/mapper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kontak/model"
	"gorm.io/gorm"
)

type KontakRepo interface {
	FindAll(q dto.KontakQuery) ([]entity.PesanKontak, int64, error)
	FindByID(id uint) (*entity.PesanKontak, error)
	// FindAnyByID mengabaikan soft delete — dipakai untuk detail & restore pesan di Sampah.
	FindAnyByID(id uint) (*entity.PesanKontak, error)
	Create(p *entity.PesanKontak) error
	MarkAsRead(id uint) error
	Delete(id uint) error
	Restore(id uint) error
	// Bulk* mengembalikan jumlah baris yang benar-benar berubah.
	BulkSoftDelete(ids []uint) (int64, error)
	BulkRestore(ids []uint) (int64, error)
}

type kontakRepo struct{ db *gorm.DB }

func NewKontakRepo(db *gorm.DB) KontakRepo {
	return &kontakRepo{db: db}
}

func (r *kontakRepo) FindAll(q dto.KontakQuery) ([]entity.PesanKontak, int64, error) {
	var items []model.PesanKontak
	var total int64

	trashed := q.IsTrashed()

	tx := r.db.Model(&model.PesanKontak{})

	if trashed {
		// Sampah: matikan soft-delete scope GORM, ambil hanya baris yang sudah dihapus.
		tx = tx.Unscoped().Where("deleted_at IS NOT NULL")
	} else {
		// Kotak Masuk: soft-delete scope GORM sudah menambahkan "deleted_at IS NULL".
		// Filter dibaca/belum dibaca hanya relevan di kotak masuk.
		switch q.Status {
		case "unread":
			tx = tx.Where("is_read = ?", false)
		case "read":
			tx = tx.Where("is_read = ?", true)
		}
	}

	// Pencarian berlaku di kedua tab. Kurungan wajib eksplisit supaya OR tidak
	// "menelan" kondisi status/soft-delete di atas.
	if s := strings.TrimSpace(q.Search); s != "" {
		like := "%" + s + "%"
		tx = tx.Where(
			"(nama LIKE ? OR email LIKE ? OR subjek LIKE ? OR pesan LIKE ?)",
			like, like, like, like,
		)
	}

	// Session() membuat statement di-clone tiap finisher, jadi Count tidak
	// meninggalkan SQL "SELECT count(*)" yang terpakai ulang oleh Find.
	tx = tx.Session(&gorm.Session{})

	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []entity.PesanKontak{}, 0, nil
	}

	// Di Sampah, yang relevan adalah kapan pesan dibuang, bukan kapan dibuat.
	orderColumn := "created_at"
	if trashed {
		orderColumn = "deleted_at"
	}
	direction := "DESC"
	if q.Sort == "oldest" {
		direction = "ASC"
	}

	_, limit, offset := q.Pagination()

	err := tx.Order(orderColumn + " " + direction).
		Limit(limit).
		Offset(offset).
		Find(&items).Error
	if err != nil {
		return nil, 0, err
	}

	return mapper.ModelListToEntity(items), total, nil
}

func (r *kontakRepo) FindByID(id uint) (*entity.PesanKontak, error) {
	var p model.PesanKontak
	err := r.db.First(&p, id).Error
	if err != nil {
		return nil, err
	}
	return mapper.ModelToEntity(&p), nil
}

func (r *kontakRepo) FindAnyByID(id uint) (*entity.PesanKontak, error) {
	var p model.PesanKontak
	err := r.db.Unscoped().First(&p, id).Error
	if err != nil {
		return nil, err
	}
	return mapper.ModelToEntity(&p), nil
}

func (r *kontakRepo) Create(p *entity.PesanKontak) error {
	return r.db.Create(mapper.EntityToModel(p)).Error
}

func (r *kontakRepo) MarkAsRead(id uint) error {
	return r.db.Model(&model.PesanKontak{}).Where("id = ?", id).
		Update("is_read", true).Error
}

func (r *kontakRepo) Delete(id uint) error {
	return r.db.Delete(&model.PesanKontak{}, id).Error
}

func (r *kontakRepo) Restore(id uint) error {
	return r.db.Unscoped().Model(&model.PesanKontak{}).
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Update("deleted_at", nil).Error
}

// BulkSoftDelete hanya menghitung baris yang tadinya masih di kotak masuk —
// scope soft-delete GORM otomatis menambahkan "deleted_at IS NULL".
func (r *kontakRepo) BulkSoftDelete(ids []uint) (int64, error) {
	tx := r.db.Where("id IN ?", ids).Delete(&model.PesanKontak{})
	return tx.RowsAffected, tx.Error
}

// BulkRestore hanya menyentuh baris yang benar-benar ada di Sampah, supaya
// pesan yang masih di kotak masuk tidak ikut ter-update. Mengembalikan jumlah
// baris yang berhasil dipulihkan.
func (r *kontakRepo) BulkRestore(ids []uint) (int64, error) {
	tx := r.db.Unscoped().Model(&model.PesanKontak{}).
		Where("id IN ? AND deleted_at IS NOT NULL", ids).
		Update("deleted_at", nil)
	return tx.RowsAffected, tx.Error
}
