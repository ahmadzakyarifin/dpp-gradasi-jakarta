package model

import (
	"time"

	"gorm.io/gorm"
)

type Student struct {
	ID             uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         *uint          `gorm:"column:user_id" json:"user_id"`
	CohortID       *uint          `gorm:"column:cohort_id" json:"cohort_id"`
	CurrentClassID *uint          `gorm:"column:current_active_class_id" json:"current_active_class_id"`
	CurrentMajorID *uint          `gorm:"column:current_major_id" json:"current_major_id"`
	NIS            string         `gorm:"column:nis;type:varchar(50);not null;uniqueIndex" json:"nis"`
	NISN           *string        `gorm:"column:nisn;type:varchar(50);uniqueIndex" json:"nisn"`
	NIK            *string        `gorm:"column:nik;type:varchar(30);uniqueIndex" json:"nik"`
	Name           string         `gorm:"column:name;type:varchar(150);not null" json:"name"`
	Gender         *string        `gorm:"column:gender;type:enum('male','female')" json:"gender"`
	Religion       *string        `gorm:"column:religion;type:varchar(50)" json:"religion"`
	BirthPlace     *string        `gorm:"column:birth_place;type:varchar(100)" json:"birth_place"`
	BirthDate      *time.Time     `gorm:"column:birth_date;type:date" json:"birth_date"`
	Address        *string        `gorm:"column:address;type:text" json:"address"`
	RT             *string        `gorm:"column:rt;type:varchar(10)" json:"rt"`
	RW             *string        `gorm:"column:rw;type:varchar(10)" json:"rw"`
	Village        *string        `gorm:"column:village;type:varchar(100)" json:"village"`
	District       *string        `gorm:"column:district;type:varchar(100)" json:"district"`
	City           *string        `gorm:"column:city;type:varchar(100)" json:"city"`
	Province       *string        `gorm:"column:province;type:varchar(100)" json:"province"`
	Email          *string        `gorm:"column:email;type:varchar(150)" json:"email"`
	Phone          *string        `gorm:"column:phone;type:varchar(30)" json:"phone"`
	EntryYear      *uint16        `gorm:"column:entry_year;type:smallint unsigned" json:"entry_year"`
	Status         string         `gorm:"column:status;type:enum('active','inactive','graduated','transferred','dropped_out');not null;default:active" json:"status"`
	PhotoPath      *string        `gorm:"column:photo_path;type:varchar(255)" json:"photo_path"`
	Description    *string        `gorm:"column:description;type:text" json:"description"`
	DepositBalance float64        `gorm:"column:deposit_balance;type:decimal(14,2);not null;default:0" json:"deposit_balance"`
	CreatedAt      time.Time      `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(3)" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`

	// Scan-only join fields
	ClassName  string `gorm:"-" json:"class_name,omitempty"`
	MajorName  string `gorm:"-" json:"major_name,omitempty"`
	ParentName string `gorm:"-" json:"parent_name,omitempty"`
}

func (Student) TableName() string {
	return "students"
}
