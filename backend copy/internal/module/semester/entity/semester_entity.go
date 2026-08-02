package entity

import "time"

type Semester struct {
	ID             uint
	AcademicYearID uint
	AcademicYear   *AcademicYear
	Name           string
	StartDate      time.Time
	EndDate        time.Time
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time

	ClassMembershipCount int
	BillingRuleCount     int
	InvoiceCount         int
	BatchCount           int
}

type AcademicYear struct {
	ID   uint
	Name string
}
