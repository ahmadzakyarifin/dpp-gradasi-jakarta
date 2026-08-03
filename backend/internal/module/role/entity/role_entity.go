package entity

import "time"

// Role adalah representasi domain dari tabel roles.
type Role struct {
	ID          uint
	Name        string
	DisplayName string
	IsSystem    bool
	IsActive    bool
	UserCount   int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
