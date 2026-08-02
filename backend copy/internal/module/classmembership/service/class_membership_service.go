package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/helper"
	activitylogdto "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activitylog/dto"
	activitylogservice "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activitylog/service"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/classmembership/dto"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/classmembership/entity"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/classmembership/mapper"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/classmembership/repository"
	"gorm.io/gorm"
)

type ClassMembershipService interface {
	GetAll(ctx context.Context, req dto.ClassMembershipQueryReq) ([]dto.ClassMembershipRes, int, error)
	Enroll(ctx context.Context, req dto.EnrollReq) (*dto.ClassMembershipRes, error)
	Move(ctx context.Context, id uint, req dto.MoveReq) (*dto.ClassMembershipRes, error)
	SetStatus(ctx context.Context, id uint, req dto.SetStatusReq) (*dto.ClassMembershipRes, error)
}

type classMembershipService struct {
	db       *gorm.DB
	repo     repository.ClassMembershipRepo
	auditSvc activitylogservice.ActivityLogService
}

func NewClassMembershipService(db *gorm.DB, repo repository.ClassMembershipRepo, auditSvc activitylogservice.ActivityLogService) ClassMembershipService {
	return &classMembershipService{db: db, repo: repo, auditSvc: auditSvc}
}

func (s *classMembershipService) log(ctx context.Context, input *activitylogdto.ActivityLogInput) {
	if s.auditSvc == nil {
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
	_ = s.auditSvc.Log(ctx, s.db, input)
}

func (s *classMembershipService) GetAll(ctx context.Context, req dto.ClassMembershipQueryReq) ([]dto.ClassMembershipRes, int, error) {
	models, total, err := s.repo.FindAll(ctx, req.Page, req.Limit, req.ActiveClassID, req.StudentID, req.Status)
	if err != nil {
		return nil, 0, err
	}
	entities := mapper.ModelListToEntity(models)
	s.fillNames(ctx, entities)
	return mapper.EntitiesToRes(entities), total, nil
}

func (s *classMembershipService) fillNames(ctx context.Context, entities []entity.ClassMembership) {
	type row struct {
		ID   uint
		Name string
	}
	var stuIDs, clsIDs []uint
	for i := range entities {
		stuIDs = append(stuIDs, entities[i].StudentID)
		clsIDs = append(clsIDs, entities[i].ActiveClassID)
	}
	stuMap := s.mapNames(ctx, "students", stuIDs)
	clsMap := s.mapNames(ctx, "active_classes", clsIDs)
	for i := range entities {
		if n, ok := stuMap[entities[i].StudentID]; ok {
			entities[i].Student = &entity.Student{ID: entities[i].StudentID, Name: n}
		}
		if n, ok := clsMap[entities[i].ActiveClassID]; ok {
			entities[i].ActiveClass = &entity.ActiveClass{ID: entities[i].ActiveClassID, Name: n}
		}
	}
}

func (s *classMembershipService) mapNames(ctx context.Context, table string, ids []uint) map[uint]string {
	out := make(map[uint]string)
	if len(ids) == 0 {
		return out
	}
	var rows []struct {
		ID   uint
		Name string
	}
	if err := s.db.WithContext(ctx).Table(table).Select("id, name").Where("id IN ?", ids).Scan(&rows).Error; err != nil {
		return out
	}
	for _, r := range rows {
		out[r.ID] = r.Name
	}
	return out
}

func (s *classMembershipService) Enroll(ctx context.Context, req dto.EnrollReq) (*dto.ClassMembershipRes, error) {
	v := helper.NewValidationError()
	if req.StudentID == 0 {
		v.Add("student_id", "Siswa wajib dipilih")
	}
	if req.ActiveClassID == 0 {
		v.Add("active_class_id", "Kelas aktif wajib dipilih")
	}
	if len(v.Errors) > 0 {
		return nil, v
	}

	// Cek siswa sudah punya membership aktif? Jika ya, tolak (harus move dulu).
	if existing, _ := s.repo.FindActiveByStudent(ctx, req.StudentID); existing != nil {
		v.Add("student_id", "Siswa sudah terdaftar di kelas lain secara aktif, lakukan perpindahan")
		return nil, v
	}

	// Ambil academic_year_id dari active_class.
	var ayID uint
	if err := s.db.WithContext(ctx).Table("active_classes").Select("academic_year_id").
		Where("id = ? AND deleted_at IS NULL", req.ActiveClassID).Scan(&ayID).Error; err != nil || ayID == 0 {
		v.Add("active_class_id", "Kelas aktif tidak ditemukan")
		return nil, v
	}

	m := mapper.EnrollReqToModel(&req, ayID)
	now := time.Now()
	if m.StartDate == nil {
		m.StartDate = &now
	}
	if err := s.repo.Create(ctx, m); err != nil {
		return nil, err
	}
	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "class_membership.create",
		EntityType:  "class_membership",
		EntityID:    &m.ID,
		EntityLabel: fmt.Sprintf("siswa #%d -> kelas #%d", m.StudentID, m.ActiveClassID),
		Description: "Mendaftarkan siswa ke kelas aktif",
	})
	return s.getRes(ctx, m.ID)
}

func (s *classMembershipService) Move(ctx context.Context, id uint, req dto.MoveReq) (*dto.ClassMembershipRes, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, helper.NewNotFoundError("Keanggotaan kelas tidak ditemukan")
	}
	v := helper.NewValidationError()
	if req.ActiveClassID == 0 {
		v.Add("active_class_id", "Kelas aktif tujuan wajib dipilih")
	}
	if len(v.Errors) > 0 {
		return nil, v
	}
	// Tutup membership lama, buka baru di kelas tujuan (riwayat terjaga).
	now := time.Now()
	if err := s.repo.UpdateActiveClass(ctx, id, req.ActiveClassID, req.SemesterID, req.AttendanceNumber, &now, req.Note); err != nil {
		return nil, err
	}
	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "class_membership.update",
		EntityType:  "class_membership",
		EntityID:    &id,
		EntityLabel: fmt.Sprintf("pindah ke kelas #%d", req.ActiveClassID),
		Description: "Memindahkan siswa ke kelas aktif lain",
	})
	return s.getRes(ctx, id)
}

func (s *classMembershipService) SetStatus(ctx context.Context, id uint, req dto.SetStatusReq) (*dto.ClassMembershipRes, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, helper.NewNotFoundError("Keanggotaan kelas tidak ditemukan")
	}
	if err := s.repo.UpdateStatus(ctx, id, req.Status, req.EndDate, req.Note); err != nil {
		return nil, err
	}
	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "class_membership.update",
		EntityType:  "class_membership",
		EntityID:    &id,
		EntityLabel: req.Status,
		Description: fmt.Sprintf("Mengubah status keanggotaan kelas menjadi: %s", req.Status),
	})
	return s.getRes(ctx, id)
}

func (s *classMembershipService) getRes(ctx context.Context, id uint) (*dto.ClassMembershipRes, error) {
	m, err := s.repo.FindByID(ctx, id)
	if err != nil || m == nil {
		return nil, helper.NewNotFoundError("Keanggotaan kelas tidak ditemukan")
	}
	entities := []entity.ClassMembership{*mapper.ModelToEntity(m)}
	s.fillNames(ctx, entities)
	res := mapper.EntityToRes(&entities[0])
	return &res, nil
}
