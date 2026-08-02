package dto

import (
	"time"
)

type ClassTemplateCreateReq struct {
	Name        string  `json:"name" binding:"required,max=100"`
	MajorID     *uint   `json:"major_id"`
	GradeLevel  int     `json:"grade_level" binding:"required,min=1,max=12"`
	Description *string `json:"description"`
}

type ClassTemplateUpdateReq struct {
	Name        string  `json:"name" binding:"required,max=100"`
	MajorID     *uint   `json:"major_id"`
	GradeLevel  int     `json:"grade_level" binding:"required,min=1,max=12"`
	Description *string `json:"description"`
}

type ClassTemplateRes struct {
	ID          uint       `json:"id"`
	MajorID     *uint      `json:"major_id"`
	Name        string     `json:"name"`
	GradeLevel  int        `json:"grade_level"`
	Description *string    `json:"description"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`

	// Relasi
	MajorName *string `json:"major_name,omitempty"`
}

type ClassTemplateQueryReq struct {
	Page       int    `form:"page,default=1"`
	Limit      int    `form:"limit,default=10"`
	Search     string `form:"search"`
	Status     string `form:"status"` // "active", "inactive", "deleted"
	MajorID    string `form:"major_id"`
	GradeLevel string `form:"grade_level"`
	Sort       string `form:"sort"` // "name_asc", "name_desc", "created_at_desc", dll
}
