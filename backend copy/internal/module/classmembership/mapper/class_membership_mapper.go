package mapper

import (
	"time"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/classmembership/dto"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/classmembership/entity"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/classmembership/model"
	"gorm.io/gorm"
)

func toGormDeletedAt(t *time.Time) gorm.DeletedAt {
	if t == nil {
		return gorm.DeletedAt{}
	}
	return gorm.DeletedAt{Time: *t, Valid: true}
}

func toDeletedAtPtr(d gorm.DeletedAt) *time.Time {
	if !d.Valid {
		return nil
	}
	t := d.Time
	return &t
}

func EnrollReqToModel(req *dto.EnrollReq, academicYearID uint) *model.ClassMembership {
	if req == nil {
		return nil
	}
	return &model.ClassMembership{
		StudentID:        req.StudentID,
		ActiveClassID:    req.ActiveClassID,
		AcademicYearID:   academicYearID,
		SemesterID:       req.SemesterID,
		AttendanceNumber: req.AttendanceNumber,
		StartDate:        req.StartDate,
		Status:           "active",
		Note:             req.Note,
	}
}

func ModelToEntity(m *model.ClassMembership) *entity.ClassMembership {
	if m == nil {
		return nil
	}
	return &entity.ClassMembership{
		ID:               m.ID,
		StudentID:        m.StudentID,
		ActiveClassID:    m.ActiveClassID,
		AcademicYearID:   m.AcademicYearID,
		SemesterID:       m.SemesterID,
		AttendanceNumber: m.AttendanceNumber,
		StartDate:        m.StartDate,
		EndDate:          m.EndDate,
		Status:           m.Status,
		Note:             m.Note,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
		DeletedAt:        toDeletedAtPtr(m.DeletedAt),
	}
}

func ModelListToEntity(models []model.ClassMembership) []entity.ClassMembership {
	list := make([]entity.ClassMembership, 0, len(models))
	for i := range models {
		list = append(list, *ModelToEntity(&models[i]))
	}
	return list
}

func EntityToRes(e *entity.ClassMembership) dto.ClassMembershipRes {
	if e == nil {
		return dto.ClassMembershipRes{}
	}
	sName, cName := "", ""
	if e.Student != nil {
		sName = e.Student.Name
	}
	if e.ActiveClass != nil {
		cName = e.ActiveClass.Name
	}
	return dto.ClassMembershipRes{
		ID:               e.ID,
		StudentID:        e.StudentID,
		StudentName:      sName,
		ActiveClassID:    e.ActiveClassID,
		ActiveClassName:  cName,
		AcademicYearID:   e.AcademicYearID,
		SemesterID:       e.SemesterID,
		AttendanceNumber: e.AttendanceNumber,
		StartDate:        e.StartDate,
		EndDate:          e.EndDate,
		Status:           e.Status,
		Note:             e.Note,
		CreatedAt:        e.CreatedAt,
		UpdatedAt:        e.UpdatedAt,
		DeletedAt:        e.DeletedAt,
	}
}

func EntitiesToRes(entities []entity.ClassMembership) []dto.ClassMembershipRes {
	list := make([]dto.ClassMembershipRes, 0, len(entities))
	for i := range entities {
		list = append(list, EntityToRes(&entities[i]))
	}
	return list
}
