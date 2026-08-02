package entity

import "time"

type ClassMembership struct {
	ID               uint
	StudentID        uint
	Student          *Student
	ActiveClassID    uint
	ActiveClass      *ActiveClass
	AcademicYearID   uint
	SemesterID       *uint
	AttendanceNumber *int
	StartDate        *time.Time
	EndDate          *time.Time
	Status           string
	Note             *string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
}

type Student struct {
	ID   uint
	Name string
}

type ActiveClass struct {
	ID   uint
	Name string
}
