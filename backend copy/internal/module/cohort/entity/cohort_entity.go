package entity

import "time"

type Cohort struct {
	ID          uint
	Name        string
	StartDate   time.Time
	EndDate     time.Time
	Description *string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time

	// Read-only aggregate fields
	StudentCount int
}
