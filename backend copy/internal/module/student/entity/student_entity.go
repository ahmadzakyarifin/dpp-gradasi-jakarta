package entity

import "time"

type Student struct {
	ID             uint       `json:"id"`
	UserID         *uint      `json:"user_id,omitempty"`
	CohortID       *uint      `json:"cohort_id,omitempty"`
	CurrentClassID *uint      `json:"current_active_class_id,omitempty"`
	CurrentMajorID *uint      `json:"current_major_id,omitempty"`
	NIS            string     `json:"nis"`
	NISN           *string    `json:"nisn,omitempty"`
	NIK            *string    `json:"nik,omitempty"`
	Name           string     `json:"name"`
	Gender         *string    `json:"gender,omitempty"`
	Religion       *string    `json:"religion,omitempty"`
	BirthPlace     *string    `json:"birth_place,omitempty"`
	BirthDate      *time.Time `json:"birth_date,omitempty"`
	Address        *string    `json:"address,omitempty"`
	RT             *string    `json:"rt,omitempty"`
	RW             *string    `json:"rw,omitempty"`
	Village        *string    `json:"village,omitempty"`
	District       *string    `json:"district,omitempty"`
	City           *string    `json:"city,omitempty"`
	Province       *string    `json:"province,omitempty"`
	Email          *string    `json:"email,omitempty"`
	Phone          *string    `json:"phone,omitempty"`
	EntryYear      *uint16    `json:"entry_year,omitempty"`
	Status         string     `json:"status"`
	PhotoPath      *string    `json:"photo_path,omitempty"`
	Description    *string    `json:"description,omitempty"`
	DepositBalance float64    `json:"deposit_balance"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`

	// Convenience fields (mapped from UserID → ParentID, CurrentClassID → ClassID, CurrentMajorID → MajorID)
	ParentID uint `json:"parent_id,omitempty"`
	ClassID  uint `json:"class_id,omitempty"`
	MajorID  uint `json:"major_id,omitempty"`

	// Joins
	ParentName string  `json:"parent_name,omitempty"`
	ClassName  *string `json:"class_name,omitempty"`
	MajorName  *string `json:"major_name,omitempty"`
}

type StudentClass struct {
	ID        uint      `json:"id"`
	StudentID uint      `json:"student_id"`
	ClassID   uint      `json:"class_id"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type ClassHistory struct {
	ID           uint      `json:"id"`
	StudentID    uint      `json:"student_id"`
	ClassID      uint      `json:"class_id"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	ClassName    string    `json:"class_name"`
	Grade        string    `json:"grade"`
	AcademicYear string    `json:"academic_year"`
}
