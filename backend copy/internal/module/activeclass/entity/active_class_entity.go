package entity

import "time"

type ActiveClass struct {
	ID                  uint
	AcademicYearID      uint
	AcademicYear        *AcademicYear
	ClassTemplateID     uint
	ClassTemplate       *ClassTemplate
	Name                string
	HomeroomNumber      *string
	HomeroomTeacherName *string
	Room                *string
	Capacity            int
	StudentCount        int
	IsActive            bool
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           *time.Time
}

type AcademicYear struct {
	ID   uint
	Name string
}

type ClassTemplate struct {
	ID   uint
	Name string
}
