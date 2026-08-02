package entity

import "time"

type Major struct {
	ID   uint
	Code string
	Name string

	IsActive bool

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	// Read-only2024 aggregate fields
	AcademicYearCount int
	ClassCount        int
	StudentCount      int
}
