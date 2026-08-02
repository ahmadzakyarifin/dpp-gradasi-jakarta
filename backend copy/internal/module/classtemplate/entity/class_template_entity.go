package entity

import "time"

type ClassTemplate struct {
	ID          uint
	MajorID     *uint
	Name        string
	GradeLevel  int
	Description *string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time

	// Relasi
	MajorName *string
}
