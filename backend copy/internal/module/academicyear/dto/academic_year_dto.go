package dto

import "time"

type AcademicYearCreateReq struct {
	Name      string    `json:"name" binding:"required,max=50"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
}

type AcademicYearUpdateReq struct {
	Name      string    `json:"name" binding:"required,max=50"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
}

type AcademicYearRes struct {
	ID        uint       `json:"id"`
	Name      string     `json:"name"`
	StartDate time.Time  `json:"start_date"`
	EndDate   time.Time  `json:"end_date"`
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	SemesterCount    int `json:"semester_count"`
	ActiveClassCount int `json:"active_class_count"`
	StudentCount     int `json:"student_count"`
	BillingRuleCount int `json:"billing_rule_count"`
	InvoiceCount     int `json:"invoice_count"`
}

type AcademicYearQueryReq struct {
	Page   int    `form:"page,default=1"`
	Limit  int    `form:"limit,default=10"`
	Search string `form:"search"`
	Status string `form:"status"`
	Sort   string `form:"sort"`
}
