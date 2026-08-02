package entity

import "time"

type AcademicYear struct {
	ID        uint
	Name      string
	StartDate time.Time
	EndDate   time.Time
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	SemesterCount    int
	ActiveClassCount int
	StudentCount     int
	BillingRuleCount int
	InvoiceCount     int
}
