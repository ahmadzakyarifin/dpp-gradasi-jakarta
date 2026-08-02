package dto

import "mime/multipart"

type StudentQueryReq struct {
	Page      int    `form:"page" binding:"omitempty,min=1"`
	Limit     int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Search    string `form:"search"`
	Filter    string `form:"filter"`
	Status    string `form:"status"`
	EntryYear int    `form:"entry_year" binding:"omitempty,min=1900"`
	ClassID   uint   `form:"class_id"`
	MajorID   uint   `form:"major_id"`
	Sort      string `form:"sort"`
}

type StudentCreateReq struct {
	NIK         string                `form:"nik" binding:"omitempty,numeric,len=16"`
	NIS         string                `form:"nis" binding:"required,min=3,max=20"`
	NISN        string                `form:"nisn" binding:"omitempty,min=10,max=10"`
	ParentID    uint                  `form:"parent_id" binding:"omitempty"`
	ClassID     uint                  `form:"class_id" binding:"omitempty"`
	MajorID     uint                  `form:"major_id" binding:"omitempty"`
	Name        string                `form:"name" binding:"required,min=2"`
	Gender      string                `form:"gender" binding:"omitempty"`
	BirthPlace  string                `form:"birth_place" binding:"omitempty"`
	BirthDate   string                `form:"birth_date" binding:"omitempty"`
	Religion    string                `form:"religion" binding:"omitempty"`
	Address     string                `form:"address" binding:"omitempty"`
	RT          string                `form:"rt" binding:"omitempty,min=1,max=5"`
	RW          string                `form:"rw" binding:"omitempty,min=1,max=5"`
	Village     string                `form:"village" binding:"omitempty"`
	District    string                `form:"district" binding:"omitempty"`
	City        string                `form:"city" binding:"omitempty"`
	Province    string                `form:"province" binding:"omitempty"`
	PhoneNumber string                `form:"phone_number" binding:"omitempty,min=9,max=20"`
	EntryYear   int                   `form:"entry_year" binding:"omitempty"`
	Email       string                `form:"email" binding:"omitempty,email"`
	Status      string                `form:"status" binding:"omitempty"`
	Description string                `form:"description" binding:"omitempty"`
	ImageFile   *multipart.FileHeader `form:"image_path"`
}

type StudentUpdateReq struct {
	ID          uint                  `form:"-"`
	NIK         string                `form:"nik" binding:"omitempty,numeric,len=16"`
	NIS         string                `form:"nis" binding:"required,min=3,max=20"`
	NISN        string                `form:"nisn" binding:"omitempty,min=10,max=10"`
	ParentID    uint                  `form:"parent_id" binding:"omitempty"`
	ClassID     uint                  `form:"class_id" binding:"omitempty"`
	MajorID     uint                  `form:"major_id" binding:"omitempty"`
	Name        string                `form:"name" binding:"required,min=2"`
	Gender      string                `form:"gender" binding:"omitempty"`
	BirthPlace  string                `form:"birth_place" binding:"omitempty"`
	BirthDate   string                `form:"birth_date" binding:"omitempty"`
	Religion    string                `form:"religion" binding:"omitempty"`
	Address     string                `form:"address" binding:"omitempty"`
	RT          string                `form:"rt" binding:"omitempty,min=1,max=5"`
	RW          string                `form:"rw" binding:"omitempty,min=1,max=5"`
	Village     string                `form:"village" binding:"omitempty"`
	District    string                `form:"district" binding:"omitempty"`
	City        string                `form:"city" binding:"omitempty"`
	Province    string                `form:"province" binding:"omitempty"`
	PhoneNumber string                `form:"phone_number" binding:"omitempty,min=9,max=20"`
	EntryYear   int                   `form:"entry_year" binding:"omitempty"`
	Email       string                `form:"email" binding:"omitempty,email"`
	Status      string                `form:"status" binding:"omitempty"`
	Description string                `form:"description" binding:"omitempty"`
	ImageFile   *multipart.FileHeader `form:"image_path"`
}

type StudentResponse struct {
	ID             uint    `json:"id"`
	NIK            string  `json:"nik"`
	NIS            string  `json:"nis"`
	NISN           string  `json:"nisn"`
	ParentID       uint    `json:"parent_id"`
	ClassID        uint    `json:"class_id"`
	MajorID        uint    `json:"major_id"`
	Name           string  `json:"name"`
	Gender         string  `json:"gender"`
	BirthPlace     string  `json:"birth_place"`
	BirthDate      string  `json:"birth_date"`
	Religion       string  `json:"religion"`
	Address        string  `json:"address"`
	RT             string  `json:"rt"`
	RW             string  `json:"rw"`
	Village        string  `json:"village"`
	District       string  `json:"district"`
	City           string  `json:"city"`
	Province       string  `json:"province"`
	PhoneNumber    string  `json:"phone_number"`
	EntryYear      int     `json:"entry_year"`
	Email          string  `json:"email"`
	ImagePath      *string `json:"image_path"`
	Status         string  `json:"status"`
	Description    string  `json:"description"`
	DepositBalance float64 `json:"deposit_balance"`
	ClassName      *string `json:"class_name"`
	MajorName      *string `json:"major_name"`
	ParentName     string  `json:"parent_name"`
}

type StudentBulkDeleteReq struct {
	IDs []uint `json:"ids" binding:"required,min=1"`
}

type StudentBulkRestoreReq struct {
	IDs []uint `json:"ids" binding:"required"`
}

type StudentBulkGraduateReq struct {
	ClassID    uint   `json:"class_id"`
	StudentIDs []uint `json:"student_ids"`
}

type StudentBulkPromoteReq struct {
	SourceClassID uint   `json:"source_class_id" binding:"required"`
	TargetClassID uint   `json:"target_class_id" binding:"required"`
	StudentIDs    []uint `json:"student_ids" binding:"required,min=1"`
}

type AcademicFilterResponse struct {
	Years   []interface{} `json:"years"`
	Majors  []interface{} `json:"majors"`
	Classes []interface{} `json:"classes"`
}

type CheckUniqueReq struct {
	Field     string `form:"field" binding:"required"`
	Value     string `form:"value" binding:"required"`
	ExcludeID uint   `form:"exclude_id"`
}
