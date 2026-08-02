package dto

import "time"

type SemesterCreateReq struct {
	AcademicYearID uint      `json:"academic_year_id" binding:"required"`
	Name           string    `json:"name" binding:"required,oneof=ganjil genap Ganjil Genap"`
	StartDate      time.Time `json:"start_date" binding:"required"`
	EndDate        time.Time `json:"end_date" binding:"required"`
}

type SemesterUpdateReq struct {
	AcademicYearID uint      `json:"academic_year_id" binding:"required"`
	Name           string    `json:"name" binding:"required,oneof=ganjil genap Ganjil Genap"`
	StartDate      time.Time `json:"start_date" binding:"required"`
	EndDate        time.Time `json:"end_date" binding:"required"`
}

type SemesterRes struct {
	ID               uint       `json:"id"`
	AcademicYearID   uint       `json:"academic_year_id"`
	AcademicYearName string     `json:"academic_year_name,omitempty"`
	Name             string     `json:"name"`
	StartDate        time.Time  `json:"start_date"`
	EndDate          time.Time  `json:"end_date"`
	IsActive         bool       `json:"is_active"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`

	ClassMembershipCount int `json:"class_membership_count"`
	BillingRuleCount     int `json:"billing_rule_count"`
	InvoiceCount         int `json:"invoice_count"`
	BatchCount           int `json:"batch_count"`
}

type SemesterQueryReq struct {
	Page           int    `form:"page,default=1"`
	Limit          int    `form:"limit,default=10"`
	Search         string `form:"search"`
	Status         string `form:"status"`
	AcademicYearID uint   `form:"academic_year_id"`
	Sort           string `form:"sort"`
}
