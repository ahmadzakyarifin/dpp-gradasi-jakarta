package dto

import "time"

type EnrollReq struct {
	StudentID        uint       `json:"student_id" binding:"required"`
	ActiveClassID    uint       `json:"active_class_id" binding:"required"`
	SemesterID       *uint      `json:"semester_id"`
	AttendanceNumber *int       `json:"attendance_number"`
	StartDate        *time.Time `json:"start_date"`
	Note             *string    `json:"note"`
}

type MoveReq struct {
	ActiveClassID    uint       `json:"active_class_id" binding:"required"`
	SemesterID       *uint      `json:"semester_id"`
	AttendanceNumber *int       `json:"attendance_number"`
	EndDate          *time.Time `json:"end_date"`
	Note             *string    `json:"note"`
}

type SetStatusReq struct {
	Status  string     `json:"status" binding:"required,oneof=active moved completed graduated inactive"`
	EndDate *time.Time `json:"end_date"`
	Note    *string    `json:"note"`
}

type ClassMembershipRes struct {
	ID               uint       `json:"id"`
	StudentID        uint       `json:"student_id"`
	StudentName      string     `json:"student_name,omitempty"`
	ActiveClassID    uint       `json:"active_class_id"`
	ActiveClassName  string     `json:"active_class_name,omitempty"`
	AcademicYearID   uint       `json:"academic_year_id"`
	SemesterID       *uint      `json:"semester_id"`
	AttendanceNumber *int       `json:"attendance_number"`
	StartDate        *time.Time `json:"start_date"`
	EndDate          *time.Time `json:"end_date"`
	Status           string     `json:"status"`
	Note             *string    `json:"note"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
}

type ClassMembershipQueryReq struct {
	Page          int    `form:"page,default=1"`
	Limit         int    `form:"limit,default=50"`
	ActiveClassID uint   `form:"active_class_id"`
	StudentID     uint   `form:"student_id"`
	Status        string `form:"status"`
	Sort          string `form:"sort"`
}
