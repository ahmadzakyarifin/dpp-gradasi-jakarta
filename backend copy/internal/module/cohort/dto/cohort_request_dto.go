package dto

import (
	"time"
)

type CohortCreateReq struct {
	Name        string    `json:"name" binding:"required,max=50"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required"`
	Description *string   `json:"description"`
}

type CohortUpdateReq struct {
	Name        string    `json:"name" binding:"required,max=50"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required"`
	Description *string   `json:"description"`
}

type CohortRes struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     time.Time  `json:"end_date"`
	Description *string    `json:"description"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`

	StudentCount int `json:"student_count"`
}

type CohortQueryReq struct {
	Page   int    `form:"page,default=1"`
	Limit  int    `form:"limit,default=10"`
	Search string `form:"search"`
	Status string `form:"status"` // "active", "inactive", "deleted"
	Sort   string `form:"sort"`
}
