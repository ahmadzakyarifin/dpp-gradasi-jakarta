package model

import (
	"time"

	"gorm.io/gorm"
)

type GuardianModel struct {
	ID          uint   `gorm:"column:id;primaryKey;autoIncrement"`
	UserID      *uint  `gorm:"column:user_id"`
	Name        string `gorm:"column:name;not null"`
	Phone       string `gorm:"column:phone"`
	Email       string `gorm:"column:email"`
	NIK         string `gorm:"column:nik"`
	Education   string `gorm:"column:education"`
	Occupation  string `gorm:"column:occupation"`
	IncomeRange string `gorm:"column:income_range"`

	CreatedAt time.Time      `gorm:"column:created_at;not null;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (GuardianModel) TableName() string { return "guardians" }
