package dto

import "time"

type ActiveClassCreateReq struct {
	AcademicYearID      uint    `json:"academic_year_id" binding:"required"`
	ClassTemplateID     uint    `json:"class_template_id" binding:"required"`
	Name                string  `json:"name" binding:"required,max=120"`
	HomeroomNumber      *string `json:"homeroom_number"`
	HomeroomTeacherName *string `json:"homeroom_teacher_name"`
	Room                *string `json:"room"`
	Capacity            int     `json:"capacity"`
	IsActive            bool    `json:"is_active"`
}

type ActiveClassUpdateReq struct {
	ClassTemplateID     uint    `json:"class_template_id" binding:"required"`
	Name                string  `json:"name" binding:"required,max=120"`
	HomeroomNumber      *string `json:"homeroom_number"`
	HomeroomTeacherName *string `json:"homeroom_teacher_name"`
	Room                *string `json:"room"`
	Capacity            int     `json:"capacity"`
	IsActive            bool    `json:"is_active"`
}

type BulkUpsertItem struct {
	ID                  uint    `json:"-"`
	ClassTemplateID     uint    `json:"class_template_id" binding:"required"`
	Name                string  `json:"name" binding:"required,max=120"`
	HomeroomNumber      *string `json:"homeroom_number"`
	HomeroomTeacherName *string `json:"homeroom_teacher_name"`
	Room                *string `json:"room"`
	Capacity            int     `json:"capacity"`
	IsActive            bool    `json:"is_active"`
}

type ActiveClassRes struct {
	ID                  uint       `json:"id"`
	AcademicYearID      uint       `json:"academic_year_id"`
	AcademicYearName    string     `json:"academic_year_name,omitempty"`
	ClassTemplateID     uint       `json:"class_template_id"`
	ClassTemplateName   string     `json:"class_template_name,omitempty"`
	Name                string     `json:"name"`
	HomeroomNumber      *string    `json:"homeroom_number"`
	HomeroomTeacherName *string    `json:"homeroom_teacher_name"`
	Room                *string    `json:"room"`
	Capacity            int        `json:"capacity"`
	StudentCount        int        `json:"student_count"`
	IsActive            bool       `json:"is_active"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at,omitempty"`
}

type ActiveClassQueryReq struct {
	Page           int    `form:"page,default=1"`
	Limit          int    `form:"limit,default=10"`
	Search         string `form:"search"`
	Status         string `form:"status"`
	AcademicYearID uint   `form:"academic_year_id"`
	Sort           string `form:"sort"`
}
