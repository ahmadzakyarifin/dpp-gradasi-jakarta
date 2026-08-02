package dto

import "time"

type MajorRes struct {
	ID          uint    `json:"id"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`

	IsActive bool `json:"is_active"`

	AcademicYearCount int `json:"academic_year_count"`
	ClassCount        int `json:"class_count"`
	StudentCount      int `json:"student_count"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
