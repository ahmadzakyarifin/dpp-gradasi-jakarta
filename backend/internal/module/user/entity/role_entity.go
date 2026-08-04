package entity

import "time"

// Role adalah representasi domain dari tabel roles.
// Disalin dari modul role (dipertahankan karena dipakai join user->role).
type Role struct {
	ID          uint
	Name        string
	DisplayName string
	IsSystem    bool
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
