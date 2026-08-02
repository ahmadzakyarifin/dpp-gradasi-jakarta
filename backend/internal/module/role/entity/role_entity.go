package entity

import "time"

// Permission adalah representasi domain dari tabel permissions.
type Permission struct {
	ID          uint
	Name        string
	DisplayName string
	Module      string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Role adalah representasi domain dari tabel roles.
type Role struct {
	ID          uint
	Name        string
	DisplayName string
	IsSystem    bool
	IsActive    bool
	Permissions []Permission
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
