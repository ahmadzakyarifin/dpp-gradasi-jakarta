package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/ahmadzakyarifin/schoolpay/backend/config"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/helper"
	activitylogdto "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activitylog/dto"
	activitylogservice "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activitylog/service"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/student/dto"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/student/entity"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/student/mapper"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/student/model"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/student/repository"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/validator"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type StudentService interface {
	GetPaginated(ctx context.Context, req dto.StudentQueryReq) ([]dto.StudentResponse, int64, error)
	GetAcademicFilters(ctx context.Context) (map[string]interface{}, error)
	GetByID(ctx context.Context, id uint) (*dto.StudentResponse, error)
	Create(ctx context.Context, req dto.StudentCreateReq) (*dto.StudentResponse, error)
	Update(ctx context.Context, id uint, req dto.StudentUpdateReq) (*dto.StudentResponse, error)
	Delete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error
	BulkDelete(ctx context.Context, ids []uint) error
	BulkRestore(ctx context.Context, ids []uint) error
	ToggleStatus(ctx context.Context, id uint) error
	BulkGraduate(ctx context.Context, classID uint, ids []uint) error
	BulkPromote(ctx context.Context, sourceClassID, targetClassID uint, ids []uint) error
	GetClassHistory(ctx context.Context, id uint) ([]entity.ClassHistory, error)
	GetDependencyInfo(ctx context.Context, id uint) (map[string]interface{}, error)
	CheckUnique(ctx context.Context, field, value string, excludeID uint) (bool, error)
	ExportExcel(ctx context.Context, req dto.StudentQueryReq) ([]byte, error)
}

type studentService struct {
	db    *gorm.DB
	repo  repository.StudentRepo
	audit activitylogservice.ActivityLogService
	cfg   *config.Config
}

func NewStudentService(
	db *gorm.DB,
	repo repository.StudentRepo,
	audit activitylogservice.ActivityLogService,
	cfg *config.Config,
) StudentService {
	return &studentService{
		db:    db,
		repo:  repo,
		audit: audit,
		cfg:   cfg,
	}
}

func (s *studentService) log(ctx context.Context, tx *gorm.DB, input *activitylogdto.ActivityLogInput) {
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
	_ = s.audit.Log(ctx, tx, input)
}

func (s *studentService) GetPaginated(ctx context.Context, req dto.StudentQueryReq) ([]dto.StudentResponse, int64, error) {
	models, total, err := s.repo.FindAllPaginated(ctx, req.Page, req.Limit, req.Search, req.Filter, req.Status, req.EntryYear, req.ClassID, req.MajorID, req.Sort)
	if err != nil {
		return nil, 0, err
	}
	entities := make([]entity.Student, len(models))
	for i, m := range models {
		entities[i] = mapper.ModelToEntity(m)
	}
	return mapper.EntityListToResponse(entities), total, nil
}

func (s *studentService) GetByID(ctx context.Context, id uint) (*dto.StudentResponse, error) {
	m, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, helper.NewNotFoundError("siswa tidak ditemukan")
	}
	e := mapper.ModelToEntity(*m)
	resp := mapper.EntityToResponse(e)
	return &resp, nil
}

func (s *studentService) Create(ctx context.Context, req dto.StudentCreateReq) (*dto.StudentResponse, error) {
	student := mapper.CreateReqToEntity(req)
	s.normalize(&student)

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		v := helper.NewValidationError()

		if err := s.validateBiz(ctx, &student, v); err != nil {
			return err
		}

		m := mapper.EntityToModel(student)

		if err := s.repo.Create(ctx, tx, &m); err != nil {
			return err
		}
		student.ID = m.ID

		if student.ClassID > 0 {
			if err := s.repo.AddClassHistory(ctx, tx, student.ID, student.ClassID); err != nil {
				return err
			}
		}

		s.log(ctx, tx, &activitylogdto.ActivityLogInput{
			Action:      "students.create",
			EntityType:  "students",
			EntityID:    &student.ID,
			EntityLabel: student.Name,
			Description: fmt.Sprintf("Membuat siswa baru: %s", student.Name),
			Metadata: map[string]any{
				"new_values": studentAuditValues(&student),
			},
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	resp := mapper.EntityToResponse(student)
	return &resp, nil
}

func (s *studentService) Update(ctx context.Context, id uint, req dto.StudentUpdateReq) (*dto.StudentResponse, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, helper.NewNotFoundError("siswa tidak ditemukan")
	}

	student := mapper.UpdateReqToEntity(existing.ID, req)
	s.normalize(&student)
	existingEntity := mapper.ModelToEntity(*existing)

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		v := helper.NewValidationError()

		if err := s.validateBiz(ctx, &student, v); err != nil {
			return err
		}

		student.PhotoPath = existing.PhotoPath
		m := mapper.EntityToModel(student)

		if err := s.repo.Update(ctx, tx, &m); err != nil {
			return err
		}

		if existing.CurrentClassID != nil && *existing.CurrentClassID != student.ClassID {
			if err := s.repo.UpdateActiveHistory(ctx, tx, student.ID, student.ClassID); err != nil {
				return err
			}
		}

		s.log(ctx, tx, &activitylogdto.ActivityLogInput{
			Action:      "students.update",
			EntityType:  "students",
			EntityID:    &student.ID,
			EntityLabel: student.Name,
			Description: fmt.Sprintf("Memperbarui data siswa: %s", student.Name),
			Metadata: map[string]any{
				"old_values": studentAuditValues(&existingEntity),
				"new_values": studentAuditValues(&student),
			},
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	resp := mapper.EntityToResponse(student)
	return &resp, nil
}

func (s *studentService) Delete(ctx context.Context, id uint) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return helper.NewNotFoundError("siswa tidak ditemukan")
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.repo.Delete(ctx, tx, id); err != nil {
			return err
		}
		s.log(ctx, tx, &activitylogdto.ActivityLogInput{
			Action:      "students.delete",
			EntityType:  "students",
			EntityID:    &id,
			EntityLabel: existing.Name,
			Description: fmt.Sprintf("Menghapus siswa: %s", existing.Name),
			Metadata: map[string]any{
				"old_values": map[string]interface{}{"status": existing.Status},
				"new_values": map[string]interface{}{"status": "deleted"},
			},
		})
		return nil
	})
}

func (s *studentService) Restore(ctx context.Context, id uint) error {
	existing, err := s.repo.FindByIDUnscoped(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return helper.NewNotFoundError("siswa tidak ditemukan di riwayat penghapusan")
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.repo.Restore(ctx, tx, id); err != nil {
			return err
		}

		s.log(ctx, tx, &activitylogdto.ActivityLogInput{
			Action:      "students.restore",
			EntityType:  "students",
			EntityID:    &existing.ID,
			EntityLabel: existing.Name,
			Description: fmt.Sprintf("Memulihkan siswa: %s", existing.Name),
			Metadata: map[string]any{
				"old_values": map[string]interface{}{
					"status": "deleted",
				},
				"new_values": map[string]interface{}{
					"status": "active",
				},
			},
		})
		return nil
	})
}

func (s *studentService) BulkDelete(ctx context.Context, ids []uint) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, id := range ids {
			if err := s.repo.Delete(ctx, tx, id); err != nil {
				return fmt.Errorf("gagal menghapus id %d: %w", id, err)
			}
		}
		return nil
	})
}

func (s *studentService) BulkRestore(ctx context.Context, ids []uint) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return s.repo.BulkRestore(ctx, tx, ids)
	})
}

func (s *studentService) ToggleStatus(ctx context.Context, id uint) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return helper.NewNotFoundError("siswa tidak ditemukan")
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.repo.ToggleStatus(ctx, tx, id); err != nil {
			return err
		}
		newStatus := "inactive"
		if existing.Status != "active" {
			newStatus = "active"
		}
		s.log(ctx, tx, &activitylogdto.ActivityLogInput{
			Action:      "students.update",
			EntityType:  "students",
			EntityID:    &id,
			EntityLabel: existing.Name,
			Description: fmt.Sprintf("Mengubah status siswa %s menjadi %s", existing.Name, newStatus),
			Metadata: map[string]any{
				"old_values": map[string]interface{}{"status": existing.Status},
				"new_values": map[string]interface{}{"status": newStatus},
			},
		})
		return nil
	})
}

func (s *studentService) BulkGraduate(ctx context.Context, classID uint, ids []uint) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		students, err := s.repo.FindActiveByIDs(ctx, ids)
		if err != nil {
			return err
		}
		if len(students) == 0 {
			return fmt.Errorf("tidak ada siswa aktif yang cocok untuk diluluskan")
		}

		var graduateIDs []uint
		for _, stu := range students {
			graduateIDs = append(graduateIDs, stu.ID)
			sid := stu.ID
			s.log(ctx, tx, &activitylogdto.ActivityLogInput{
				Action:      "students.update",
				EntityType:  "students",
				EntityID:    &sid,
				EntityLabel: stu.Name,
				Description: fmt.Sprintf("Meluluskan siswa: %s", stu.Name),
				Metadata: map[string]any{
					"old_values": map[string]interface{}{"status": "active"},
					"new_values": map[string]interface{}{"status": "graduated"},
				},
			})
		}

		if err := s.repo.DeactivateHistory(ctx, tx, graduateIDs); err != nil {
			return err
		}
		return s.repo.BulkGraduate(ctx, tx, graduateIDs)
	})
}

func (s *studentService) BulkPromote(ctx context.Context, sourceClassID, targetClassID uint, ids []uint) error {
	if sourceClassID == targetClassID {
		return fmt.Errorf("kelas asal dan tujuan tidak boleh sama")
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		students, err := s.repo.FindActiveByIDs(ctx, ids)
		if err != nil {
			return err
		}
		if len(students) == 0 {
			return fmt.Errorf("tidak ada siswa aktif yang cocok untuk dipindahkan")
		}

		var promoteIDs []uint
		for _, stu := range students {
			promoteIDs = append(promoteIDs, stu.ID)
		}

		if err := s.repo.DeactivateHistory(ctx, tx, promoteIDs); err != nil {
			return err
		}
		if err := s.repo.BulkPromote(ctx, tx, targetClassID, promoteIDs); err != nil {
			return err
		}

		for _, id := range promoteIDs {
			if err := s.repo.AddClassHistory(ctx, tx, id, targetClassID); err != nil {
				return err
			}
		}

		for _, stu := range students {
			sid := stu.ID
			s.log(ctx, tx, &activitylogdto.ActivityLogInput{
				Action:      "students.update",
				EntityType:  "students",
				EntityID:    &sid,
				EntityLabel: stu.Name,
				Description: fmt.Sprintf("Memindahkan siswa %s ke kelas %d", stu.Name, targetClassID),
				Metadata: map[string]any{
					"old_values": map[string]interface{}{"class_id": sourceClassID},
					"new_values": map[string]interface{}{"class_id": targetClassID},
				},
			})
		}
		return nil
	})
}

func (s *studentService) GetClassHistory(ctx context.Context, id uint) ([]entity.ClassHistory, error) {
	return s.repo.GetClassHistory(ctx, id)
}

func (s *studentService) GetDependencyInfo(ctx context.Context, id uint) (map[string]interface{}, error) {
	m, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, helper.NewNotFoundError("siswa tidak ditemukan")
	}
	return map[string]interface{}{
		"id":              m.ID,
		"name":            m.Name,
		"status":          m.Status,
		"deposit_balance": m.DepositBalance,
	}, nil
}

func (s *studentService) CheckUnique(ctx context.Context, field, value string, excludeID uint) (bool, error) {
	switch field {
	case "nik":
		existing, err := s.repo.FindByNIK(ctx, value)
		if err != nil || existing == nil {
			return true, nil
		}
		return existing.ID == excludeID, nil
	case "nisn":
		existing, err := s.repo.FindByNISN(ctx, value)
		if err != nil || existing == nil {
			return true, nil
		}
		return existing.ID == excludeID, nil
	case "nis":
		existing, err := s.repo.FindByNIS(ctx, value)
		if err != nil || existing == nil {
			return true, nil
		}
		return existing.ID == excludeID, nil
	case "email":
		existing, err := s.repo.FindByEmail(ctx, value)
		if err != nil || existing == nil {
			return true, nil
		}
		return existing.ID == excludeID, nil
	case "phone":
		existing, err := s.repo.FindByPhone(ctx, value)
		if err != nil || existing == nil {
			return true, nil
		}
		return existing.ID == excludeID, nil
	default:
		return false, fmt.Errorf("field '%s' tidak didukung", field)
	}
}

func (s *studentService) GetAcademicFilters(ctx context.Context) (map[string]interface{}, error) {
	type AcademicYear struct {
		ID   uint   `json:"id"`
		Year int    `json:"year"`
		Name string `json:"name"`
	}
	type Major struct {
		ID       uint   `json:"id"`
		Name     string `json:"name"`
		IsActive bool   `json:"is_active"`
	}
	type Class struct {
		ID             uint   `json:"id"`
		Name           string `json:"name"`
		AcademicYearID *uint  `json:"academic_year_id"`
		MajorID        *uint  `json:"major_id"`
	}

	var years []AcademicYear
	var majors []Major
	var classes []Class

	s.db.WithContext(ctx).Raw("SELECT id, year, name FROM academic_years WHERE deleted_at IS NULL ORDER BY year DESC").Find(&years)
	s.db.WithContext(ctx).Raw("SELECT id, name, is_active FROM majors WHERE deleted_at IS NULL ORDER BY name").Find(&majors)
	s.db.WithContext(ctx).Raw("SELECT ac.id, ac.name, ac.academic_year_id, ct.major_id FROM active_classes ac JOIN class_templates ct ON ac.class_template_id = ct.id WHERE ac.deleted_at IS NULL ORDER BY ac.name").Find(&classes)

	return map[string]interface{}{
		"years":   years,
		"majors":  majors,
		"classes": classes,
	}, nil
}

func (s *studentService) ExportExcel(ctx context.Context, req dto.StudentQueryReq) ([]byte, error) {
	req.Limit = 1000000
	models, _, err := s.repo.FindAllPaginated(ctx, 1, req.Limit, req.Search, req.Filter, req.Status, req.EntryYear, req.ClassID, req.MajorID, "")
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	sheet := "Data Siswa"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{
		"NIS", "NISN", "NAMA LENGKAP", "NIK", "L/P", "TEMPAT LAHIR", "TANGGAL LAHIR",
		"AGAMA", "EMAIL", "NO. WHATSAPP", "PROVINSI", "KOTA", "KECAMATAN", "KELURAHAN",
		"ALAMAT", "RT", "RW", "KELAS", "JURUSAN", "ANGKATAN", "NAMA WALI", "STATUS",
	}
	colWidths := []float64{18, 18, 28, 20, 8, 20, 20, 14, 28, 20, 20, 20, 18, 18, 30, 6, 6, 15, 18, 10, 28, 12}

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
		colName, _ := excelize.ColumnNumberToName(i + 1)
		w := 18.0
		if i < len(colWidths) {
			w = colWidths[i]
		}
		f.SetColWidth(sheet, colName, colName, w)
	}

	bodyStyle, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{Vertical: "center"},
	})

	for i, m := range models {
		row := i + 2
		birthDate := ""
		if m.BirthDate != nil {
			birthDate = m.BirthDate.Format("02/01/2006")
		}
		phone := ""
		if m.Phone != nil {
			phone = *m.Phone
		}
		if len(phone) > 0 && phone[0] != '+' {
			phone = "+" + phone
		}
		nik := ""
		if m.NIK != nil {
			nik = *m.NIK
		}
		email := ""
		if m.Email != nil {
			email = *m.Email
		}
		province := ""
		if m.Province != nil {
			province = *m.Province
		}
		city := ""
		if m.City != nil {
			city = *m.City
		}
		district := ""
		if m.District != nil {
			district = *m.District
		}
		village := ""
		if m.Village != nil {
			village = *m.Village
		}
		address := ""
		if m.Address != nil {
			address = *m.Address
		}
		rt := ""
		if m.RT != nil {
			rt = *m.RT
		}
		rw := ""
		if m.RW != nil {
			rw = *m.RW
		}
		entryYear := uint16(0)
		if m.EntryYear != nil {
			entryYear = *m.EntryYear
		}
		nisn := ""
		if m.NISN != nil {
			nisn = *m.NISN
		}
		gender := ""
		if m.Gender != nil {
			gender = *m.Gender
		}
		birthPlace := ""
		if m.BirthPlace != nil {
			birthPlace = *m.BirthPlace
		}
		religion := ""
		if m.Religion != nil {
			religion = *m.Religion
		}

		vals := []interface{}{
			m.NIS, nisn, m.Name, nik, gender,
			birthPlace, birthDate, religion,
			email, phone, province, city, district, village,
			address, rt, rw, m.ClassName, m.MajorName,
			entryYear, m.ParentName, strings.ToUpper(m.Status),
		}
		for ci, v := range vals {
			cell, _ := excelize.CoordinatesToCellName(ci+1, row)
			f.SetCellValue(sheet, cell, v)
		}
		startCell, _ := excelize.CoordinatesToCellName(1, row)
		endCell, _ := excelize.CoordinatesToCellName(len(headers), row)
		f.SetCellStyle(sheet, startCell, endCell, bodyStyle)
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Internal validators

func (s *studentService) normalize(student *entity.Student) {
	student.Name = strings.Title(strings.ToLower(strings.TrimSpace(student.Name)))
	if student.Email != nil {
		e := strings.ToLower(strings.TrimSpace(*student.Email))
		student.Email = &e
	}
	if student.NIK != nil {
		nik := strings.TrimSpace(*student.NIK)
		student.NIK = &nik
	}
	if student.BirthPlace != nil {
		bp := strings.TrimSpace(*student.BirthPlace)
		student.BirthPlace = &bp
	}
	if student.Religion != nil {
		r := strings.TrimSpace(*student.Religion)
		student.Religion = &r
	}
	if student.Address != nil {
		a := strings.TrimSpace(*student.Address)
		student.Address = &a
	}
	if student.RT != nil {
		r := strings.TrimSpace(*student.RT)
		student.RT = &r
	}
	if student.RW != nil {
		r := strings.TrimSpace(*student.RW)
		student.RW = &r
	}
	if student.Description != nil {
		d := strings.TrimSpace(*student.Description)
		student.Description = &d
	}
	if student.Phone != nil {
		p := validator.NormalizePhoneNumber(*student.Phone)
		student.Phone = &p
	}
}

func (s *studentService) validateBiz(ctx context.Context, student *entity.Student, v *helper.ValidationError) error {
	if student.Phone != nil && *student.Phone != "" {
		if err := validator.ValidatePhoneNumber(*student.Phone); err != nil {
			v.Add("phone_number", "Nomor telepon tidak valid")
		}
	}

	if student.MajorID > 0 {
		var isActive bool
		err := s.db.WithContext(ctx).Table("majors").Select("is_active").Where("id = ? AND deleted_at IS NULL", student.MajorID).Row().Scan(&isActive)
		if err != nil || !isActive {
			v.Add("major_id", "Jurusan tidak valid atau sedang tidak aktif")
		}
	}

	if student.NIK != nil && *student.NIK != "" {
		existing, _ := s.repo.FindByNIK(ctx, *student.NIK)
		if existing != nil && existing.ID != student.ID {
			v.Add("nik", fmt.Sprintf("NIK '%s' sudah terdaftar", *student.NIK))
		}
	}

	if student.NISN != nil && *student.NISN != "" {
		existing, _ := s.repo.FindByNISN(ctx, *student.NISN)
		if existing != nil && existing.ID != student.ID {
			v.Add("nisn", fmt.Sprintf("NISN '%s' sudah terdaftar", *student.NISN))
		}
	}

	existing, _ := s.repo.FindByNIS(ctx, student.NIS)
	if existing != nil && existing.ID != student.ID {
		v.Add("nis", fmt.Sprintf("NIS '%s' sudah terdaftar", student.NIS))
	}

	if student.Email != nil && *student.Email != "" {
		existing, _ := s.repo.FindByEmail(ctx, *student.Email)
		if existing != nil && existing.ID != student.ID {
			v.Add("email", fmt.Sprintf("Email '%s' sudah terdaftar", *student.Email))
		}
	}

	if student.Phone != nil && *student.Phone != "" {
		existing, _ := s.repo.FindByPhone(ctx, *student.Phone)
		if existing != nil && existing.ID != student.ID {
			v.Add("phone_number", fmt.Sprintf("Nomor telepon '%s' sudah terdaftar", *student.Phone))
		}
	}

	if !v.IsEmpty() {
		return v
	}
	return nil
}

func studentAuditValues(student *entity.Student) map[string]interface{} {
	if student == nil {
		return nil
	}
	birthDate := ""
	if student.BirthDate != nil && !student.BirthDate.IsZero() {
		birthDate = student.BirthDate.Format("2006-01-02")
	}
	nik := ""
	if student.NIK != nil {
		nik = *student.NIK
	}
	return map[string]interface{}{
		"name":         student.Name,
		"nik":          nik,
		"nis":          student.NIS,
		"nisn":         student.NISN,
		"parent_id":    student.ParentID,
		"gender":       student.Gender,
		"birth_place":  student.BirthPlace,
		"birth_date":   birthDate,
		"email":        student.Email,
		"phone_number": student.Phone,
		"class_id":     student.ClassID,
		"major_id":     student.MajorID,
		"entry_year":   student.EntryYear,
		"address":      student.Address,
		"province":     student.Province,
		"city":         student.City,
		"district":     student.District,
		"village":      student.Village,
		"rt":           student.RT,
		"rw":           student.RW,
		"status":       student.Status,
	}
}

// Ensure imports are used
var _ = &model.Student{}
var _ = validator.NormalizePhoneNumber
