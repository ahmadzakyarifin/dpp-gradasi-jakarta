package entity

import "time"

type Guardian struct {
	ID          uint       `json:"id"`
	UserID      *uint      `json:"user_id"`
	Name        string     `json:"name"`
	Phone       string     `json:"phone"`
	Email       string     `json:"email"`
	NIK         string     `json:"nik"`
	Education   string     `json:"education"`
	Occupation  string     `json:"occupation"`
	IncomeRange string     `json:"income_range"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
