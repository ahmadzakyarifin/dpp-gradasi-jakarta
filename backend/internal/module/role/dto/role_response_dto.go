package dto

import "time"

type RoleRes struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	DisplayName string     `json:"display_name"`
	IsSystem    bool       `json:"is_system"`
	IsActive    bool       `json:"is_active"`
	UserCount   int        `json:"user_count"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
