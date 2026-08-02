package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/helper"
	activitylogdto "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activitylog/dto"
	activitylogservice "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activitylog/service"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/major/dto"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/major/entity"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/major/mapper"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/major/repository"
	"gorm.io/gorm"
)

type MajorService interface {
	Create(ctx context.Context, j *entity.Major) error
	GetAll(ctx context.Context, req dto.MajorQueryReq) ([]entity.Major, int, error)
	Update(ctx context.Context, id uint, req *dto.MajorUpdateReq) (*entity.Major, error)
	Delete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error
	UpdateStatus(ctx context.Context, id uint, req *dto.MajorStatusReq) error
	BulkDelete(ctx context.Context, ids []uint) error
	BulkRestore(ctx context.Context, ids []uint) error
	GetByID(ctx context.Context, id uint) (*entity.Major, error)
	GetDependencyInfo(ctx context.Context, id uint) (map[string]interface{}, error)
}

type majorService struct {
	db    *gorm.DB
	repo  repository.MajorRepo
	audit activitylogservice.ActivityLogService
}

func NewMajorService(db *gorm.DB, repo repository.MajorRepo, audit activitylogservice.ActivityLogService) MajorService {
	return &majorService{db: db, repo: repo, audit: audit}
}

func (s *majorService) log(ctx context.Context, db *gorm.DB, input *activitylogdto.ActivityLogInput) {
	if s.audit == nil {
		return
	}
	userID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(ctx)
	var uID *uint
	if userID > 0 {
		uID = &userID
	}
	input.ActorID = uID
	input.ActorName = userName
	input.ActorRole = role
	input.IPAddress = ipAddress
	input.UserAgent = userAgent

	_ = s.audit.Log(ctx, db, input)
}

func majorDependencyMessage(classCount, studentCount, ayCount int) string {
	var messages []string
	if classCount > 0 {
		messages = append(messages, fmt.Sprintf("%d kelas aktif", classCount))
	}
	if studentCount > 0 {
		messages = append(messages, fmt.Sprintf("%d siswa aktif", studentCount))
	}
	if ayCount > 0 {
		messages = append(messages, fmt.Sprintf("%d angkatan", ayCount))
	}
	return strings.Join(messages, ", ")
}

func (s *majorService) dependencyCounts(ctx context.Context, id uint) (int, int, int, error) {
	classCount, err := s.repo.CountClasses(ctx, id)
	if err != nil {
		return 0, 0, 0, err
	}
	studentCount, err := s.repo.CountStudents(ctx, id)
	if err != nil {
		return 0, 0, 0, err
	}
	ayCount, err := s.repo.CountAcademicYears(ctx, id)
	if err != nil {
		return 0, 0, 0, err
	}
	return classCount, studentCount, ayCount, nil
}

func (s *majorService) Create(ctx context.Context, j *entity.Major) error {
	j.Name = strings.TrimSpace(j.Name)

	v := helper.NewValidationError()
	if j.Name == "" {
		v.Add("name", "Nama jurusan wajib diisi")
	}

	j.Code = strings.ToUpper(strings.TrimSpace(j.Code))
	if j.Code != "" {
		if len(j.Code) < 2 || len(j.Code) > 12 {
			v.Add("code", "Kode jurusan harus 2-12 karakter")
		}
	}

	if len(v.Errors) > 0 {
		return v
	}

	err := s.repo.Create(ctx, j)
	if err != nil {
		if strings.Contains(err.Error(), "1062") || strings.Contains(err.Error(), "Duplicate entry") {
			if strings.Contains(err.Error(), "uq_majors_code") {
				v.Add("code", fmt.Sprintf("Kode jurusan '%s' sudah digunakan", j.Code))
				return v
			}
			if strings.Contains(err.Error(), "uq_majors_name") {
				v.Add("name", fmt.Sprintf("Nama jurusan '%s' sudah terdaftar", j.Name))
				return v
			}
		}
		return err
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "CREATE",
		EntityType:  "majors",
		EntityID:    &j.ID,
		EntityLabel: j.Name,
		Description: fmt.Sprintf("Membuat jurusan baru: %s", j.Name),
		Metadata: map[string]any{
			"new_values": map[string]interface{}{"name": j.Name, "code": j.Code, "is_active": j.IsActive},
		},
	})

	return nil
}

func (s *majorService) GetAll(ctx context.Context, req dto.MajorQueryReq) ([]entity.Major, int, error) {
	return s.repo.FindAll(ctx, req.Page, req.Limit, req.Search, req.Status, req.Sort)
}

func (s *majorService) GetByID(ctx context.Context, id uint) (*entity.Major, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, helper.NewNotFoundError("Jurusan tidak ditemukan")
	}
	return existing, nil
}

func (s *majorService) Update(ctx context.Context, id uint, req *dto.MajorUpdateReq) (*entity.Major, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Code = strings.ToUpper(strings.TrimSpace(req.Code))

	v := helper.NewValidationError()

	if req.Name == "" {
		v.Add("name", "Nama jurusan wajib diisi")
	}
	if req.Code != "" {
		if len(req.Code) < 2 || len(req.Code) > 12 {
			v.Add("code", "Kode jurusan harus 2-12 karakter")
		}
	}

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, helper.NewNotFoundError("Jurusan tidak ditemukan")
	}

	nameChanged := !strings.EqualFold(strings.TrimSpace(existing.Name), req.Name)
	if nameChanged {
		classCount, studentCount, ayCount, err := s.dependencyCounts(ctx, id)
		if err != nil {
			return nil, err
		}
		if classCount+studentCount+ayCount > 0 {
			msg := majorDependencyMessage(classCount, studentCount, ayCount)
			v.Add("name", fmt.Sprintf("Nama jurusan tidak dapat diubah karena masih terhubung dengan %s. Nonaktifkan jurusan jika hanya ingin menyembunyikan dari form", msg))
		}
	}

	if len(v.Errors) > 0 {
		return nil, v
	}

	oldVals := map[string]interface{}{"name": existing.Name, "code": existing.Code, "is_active": existing.IsActive}

	mapper.UpdateReqToEntity(req, existing)

	err = s.repo.Update(ctx, existing)
	if err != nil {
		if strings.Contains(err.Error(), "1062") || strings.Contains(err.Error(), "Duplicate entry") {
			if strings.Contains(err.Error(), "uq_majors_code") {
				v.Add("code", fmt.Sprintf("Kode jurusan '%s' sudah digunakan", req.Code))
				return nil, v
			}
			if strings.Contains(err.Error(), "uq_majors_name") {
				v.Add("name", fmt.Sprintf("Nama jurusan '%s' sudah digunakan", req.Name))
				return nil, v
			}
		}
		return nil, err
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "UPDATE",
		EntityType:  "majors",
		EntityID:    &existing.ID,
		EntityLabel: existing.Name,
		Description: fmt.Sprintf("Memperbarui jurusan: %s", existing.Name),
		Metadata: map[string]any{
			"old_values": oldVals,
			"new_values": map[string]interface{}{"name": existing.Name, "code": existing.Code, "is_active": existing.IsActive},
		},
	})

	return existing, nil
}

func (s *majorService) Delete(ctx context.Context, id uint) error {
	classCount, studentCount, ayCount, err := s.dependencyCounts(ctx, id)
	if err != nil {
		return err
	}
	if classCount+studentCount+ayCount > 0 {
		msg := majorDependencyMessage(classCount, studentCount, ayCount)
		v := helper.NewValidationError()
		v.Add("general", fmt.Sprintf("Jurusan tidak dapat dihapus karena masih terhubung dengan %s", msg))
		return v
	}

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return helper.NewNotFoundError("Jurusan tidak ditemukan")
	}

	err = s.repo.Delete(ctx, id)
	if err == nil {
		s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
			Action:      "DELETE",
			EntityType:  "majors",
			EntityID:    &id,
			EntityLabel: existing.Name,
			Description: fmt.Sprintf("Menghapus jurusan: %s", existing.Name),
			Metadata: map[string]any{
				"old_values": map[string]interface{}{"is_active": existing.IsActive, "status": "active"},
				"new_values": map[string]interface{}{"is_active": false, "status": "deleted"},
			},
		})
	}
	return err
}

func (s *majorService) Restore(ctx context.Context, id uint) error {
	existing, err := s.repo.FindByIDUnscoped(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil || existing.DeletedAt == nil {
		return helper.NewNotFoundError("Jurusan tidak ditemukan di riwayat penghapusan")
	}

	err = s.repo.Restore(ctx, id)
	if err == nil {
		s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
			Action:      "RESTORE",
			EntityType:  "majors",
			EntityID:    &id,
			EntityLabel: existing.Name,
			Description: fmt.Sprintf("Memulihkan jurusan: %s", existing.Name),
			Metadata: map[string]any{
				"old_values": map[string]interface{}{"is_active": existing.IsActive, "status": "deleted"},
				"new_values": map[string]interface{}{"is_active": true, "status": "active"},
			},
		})
	}
	return err
}

func (s *majorService) UpdateStatus(ctx context.Context, id uint, req *dto.MajorStatusReq) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return helper.NewNotFoundError("Jurusan tidak ditemukan")
	}

	oldVals := map[string]interface{}{"is_active": existing.IsActive}

	mapper.StatusReqToEntity(req, existing)

	err = s.repo.Update(ctx, existing)
	if err == nil {
		s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
			Action:      "UPDATE",
			EntityType:  "majors",
			EntityID:    &id,
			EntityLabel: existing.Name,
			Description: fmt.Sprintf("Mengubah status aktif/nonaktif jurusan: %s", existing.Name),
			Metadata: map[string]any{
				"old_values": oldVals,
				"new_values": map[string]interface{}{"is_active": existing.IsActive},
			},
		})
	}
	return err
}

func (s *majorService) BulkDelete(ctx context.Context, ids []uint) error {
	for _, id := range ids {
		classCount, studentCount, ayCount, err := s.dependencyCounts(ctx, id)
		if err != nil {
			return err
		}
		if classCount+studentCount+ayCount > 0 {
			msg := majorDependencyMessage(classCount, studentCount, ayCount)
			v := helper.NewValidationError()
			v.Add("general", fmt.Sprintf("Beberapa jurusan tidak dapat dihapus karena masih terhubung dengan %s", msg))
			return v
		}
	}

	err := s.repo.BulkDelete(ctx, ids)
	if err == nil {
		nameMap, _ := s.repo.FindNamesByIDs(ctx, ids)
		for _, id := range ids {
			majorID := id
			name := nameMap[id]
			s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
				Action:      "DELETE",
				EntityType:  "majors",
				EntityID:    &majorID,
				EntityLabel: name,
				Description: fmt.Sprintf("Menghapus jurusan (bulk): %s", name),
				Metadata: map[string]any{
					"old_values": map[string]interface{}{"status": "active"},
					"new_values": map[string]interface{}{"status": "deleted"},
				},
			})
		}
	}
	return err
}

func (s *majorService) BulkRestore(ctx context.Context, ids []uint) error {
	err := s.repo.BulkRestore(ctx, ids)
	if err == nil {
		nameMap, _ := s.repo.FindNamesByIDs(ctx, ids)
		for _, id := range ids {
			majorID := id
			name := nameMap[id]
			s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
				Action:      "RESTORE",
				EntityType:  "majors",
				EntityID:    &majorID,
				EntityLabel: name,
				Description: fmt.Sprintf("Memulihkan jurusan (bulk): %s", name),
				Metadata: map[string]any{
					"old_values": map[string]interface{}{"status": "deleted"},
					"new_values": map[string]interface{}{"status": "active"},
				},
			})
		}
	}
	return err
}

func (s *majorService) GetDependencyInfo(ctx context.Context, id uint) (map[string]interface{}, error) {
	classCount, studentCount, ayCount, err := s.dependencyCounts(ctx, id)
	if err != nil {
		return nil, err
	}

	message := majorDependencyMessage(classCount, studentCount, ayCount)

	return map[string]interface{}{
		"has_dependencies": message != "",
		"message":          message,
		"counts": map[string]int{
			"classes":        classCount,
			"students":       studentCount,
			"academic_years": ayCount,
		},
	}, nil
}
