package mapper

import (
	"time"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/student/dto"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/student/entity"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/student/model"
)

func CreateReqToEntity(req dto.StudentCreateReq) entity.Student {
	var birthDate *time.Time
	if req.BirthDate != "" {
		if parsed, err := time.Parse("02/01/2006", req.BirthDate); err == nil {
			birthDate = &parsed
		}
	}
	return entity.Student{
		NIK:         nilIfEmpty(req.NIK),
		NIS:         req.NIS,
		NISN:        nilIfEmpty(req.NISN),
		ParentID:    req.ParentID,
		ClassID:     req.ClassID,
		MajorID:     req.MajorID,
		Name:        req.Name,
		Gender:      nilIfEmpty(req.Gender),
		BirthPlace:  nilIfEmpty(req.BirthPlace),
		BirthDate:   birthDate,
		Religion:    nilIfEmpty(req.Religion),
		Address:     nilIfEmpty(req.Address),
		RT:          nilIfEmpty(req.RT),
		RW:          nilIfEmpty(req.RW),
		Village:     nilIfEmpty(req.Village),
		District:    nilIfEmpty(req.District),
		City:        nilIfEmpty(req.City),
		Province:    nilIfEmpty(req.Province),
		Phone:       nilIfEmpty(req.PhoneNumber),
		EntryYear:   intNilIfZero(req.EntryYear),
		Email:       nilIfEmpty(req.Email),
		Status:      req.Status,
		Description: nilIfEmpty(req.Description),
	}
}

func UpdateReqToEntity(id uint, req dto.StudentUpdateReq) entity.Student {
	var birthDate *time.Time
	if req.BirthDate != "" {
		if parsed, err := time.Parse("02/01/2006", req.BirthDate); err == nil {
			birthDate = &parsed
		}
	}
	return entity.Student{
		ID:          id,
		NIK:         nilIfEmpty(req.NIK),
		NIS:         req.NIS,
		NISN:        nilIfEmpty(req.NISN),
		ParentID:    req.ParentID,
		ClassID:     req.ClassID,
		MajorID:     req.MajorID,
		Name:        req.Name,
		Gender:      nilIfEmpty(req.Gender),
		BirthPlace:  nilIfEmpty(req.BirthPlace),
		BirthDate:   birthDate,
		Religion:    nilIfEmpty(req.Religion),
		Address:     nilIfEmpty(req.Address),
		RT:          nilIfEmpty(req.RT),
		RW:          nilIfEmpty(req.RW),
		Village:     nilIfEmpty(req.Village),
		District:    nilIfEmpty(req.District),
		City:        nilIfEmpty(req.City),
		Province:    nilIfEmpty(req.Province),
		Phone:       nilIfEmpty(req.PhoneNumber),
		EntryYear:   intNilIfZero(req.EntryYear),
		Email:       nilIfEmpty(req.Email),
		Status:      req.Status,
		Description: nilIfEmpty(req.Description),
	}
}

func EntityToModel(e entity.Student) model.Student {
	userID := e.UserID
	if userID == nil && e.ParentID > 0 {
		uid := e.ParentID
		userID = &uid
	}
	currentClassID := e.CurrentClassID
	if currentClassID == nil && e.ClassID > 0 {
		cid := e.ClassID
		currentClassID = &cid
	}
	currentMajorID := e.CurrentMajorID
	if currentMajorID == nil && e.MajorID > 0 {
		mid := e.MajorID
		currentMajorID = &mid
	}

	return model.Student{
		ID:             e.ID,
		UserID:         userID,
		CohortID:       e.CohortID,
		CurrentClassID: currentClassID,
		CurrentMajorID: currentMajorID,
		NIS:            e.NIS,
		NISN:           e.NISN,
		NIK:            e.NIK,
		Name:           e.Name,
		Gender:         e.Gender,
		Religion:       e.Religion,
		BirthPlace:     e.BirthPlace,
		BirthDate:      e.BirthDate,
		Address:        e.Address,
		RT:             e.RT,
		RW:             e.RW,
		Village:        e.Village,
		District:       e.District,
		City:           e.City,
		Province:       e.Province,
		Email:          e.Email,
		Phone:          e.Phone,
		EntryYear:      e.EntryYear,
		Status:         e.Status,
		PhotoPath:      e.PhotoPath,
		Description:    e.Description,
		DepositBalance: e.DepositBalance,
	}
}

func ModelToEntity(m model.Student) entity.Student {
	parentID := uint(0)
	if m.UserID != nil {
		parentID = *m.UserID
	}
	classID := uint(0)
	if m.CurrentClassID != nil {
		classID = *m.CurrentClassID
	}
	majorID := uint(0)
	if m.CurrentMajorID != nil {
		majorID = *m.CurrentMajorID
	}

	return entity.Student{
		ID:             m.ID,
		UserID:         m.UserID,
		CohortID:       m.CohortID,
		CurrentClassID: m.CurrentClassID,
		CurrentMajorID: m.CurrentMajorID,
		ParentID:       parentID,
		ClassID:        classID,
		MajorID:        majorID,
		NIS:            m.NIS,
		NISN:           m.NISN,
		NIK:            m.NIK,
		Name:           m.Name,
		Gender:         m.Gender,
		Religion:       m.Religion,
		BirthPlace:     m.BirthPlace,
		BirthDate:      m.BirthDate,
		Address:        m.Address,
		RT:             m.RT,
		RW:             m.RW,
		Village:        m.Village,
		District:       m.District,
		City:           m.City,
		Province:       m.Province,
		Email:          m.Email,
		Phone:          m.Phone,
		EntryYear:      m.EntryYear,
		Status:         m.Status,
		PhotoPath:      m.PhotoPath,
		Description:    m.Description,
		DepositBalance: m.DepositBalance,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
		DeletedAt: func() *time.Time {
			if m.DeletedAt.Valid {
				return &m.DeletedAt.Time
			}
			return nil
		}(),
	}
}

func EntityToResponse(e entity.Student) dto.StudentResponse {
	birthDate := ""
	if e.BirthDate != nil && !e.BirthDate.IsZero() {
		birthDate = e.BirthDate.Format("02/01/2006")
	}
	nik := ""
	if e.NIK != nil {
		nik = *e.NIK
	}
	nisn := ""
	if e.NISN != nil {
		nisn = *e.NISN
	}
	gender := ""
	if e.Gender != nil {
		gender = *e.Gender
	}
	birthPlace := ""
	if e.BirthPlace != nil {
		birthPlace = *e.BirthPlace
	}
	religion := ""
	if e.Religion != nil {
		religion = *e.Religion
	}
	address := ""
	if e.Address != nil {
		address = *e.Address
	}
	rt := ""
	if e.RT != nil {
		rt = *e.RT
	}
	rw := ""
	if e.RW != nil {
		rw = *e.RW
	}
	village := ""
	if e.Village != nil {
		village = *e.Village
	}
	district := ""
	if e.District != nil {
		district = *e.District
	}
	city := ""
	if e.City != nil {
		city = *e.City
	}
	province := ""
	if e.Province != nil {
		province = *e.Province
	}
	phone := ""
	if e.Phone != nil {
		phone = *e.Phone
	}
	email := ""
	if e.Email != nil {
		email = *e.Email
	}
	desc := ""
	if e.Description != nil {
		desc = *e.Description
	}
	var entryYear int
	if e.EntryYear != nil {
		entryYear = int(*e.EntryYear)
	}

	return dto.StudentResponse{
		ID:             e.ID,
		NIK:            nik,
		NIS:            e.NIS,
		NISN:           nisn,
		ParentID:       e.ParentID,
		ClassID:        e.ClassID,
		MajorID:        e.MajorID,
		Name:           e.Name,
		Gender:         gender,
		BirthPlace:     birthPlace,
		BirthDate:      birthDate,
		Religion:       religion,
		Address:        address,
		RT:             rt,
		RW:             rw,
		Village:        village,
		District:       district,
		City:           city,
		Province:       province,
		PhoneNumber:    phone,
		EntryYear:      entryYear,
		Email:          email,
		ImagePath:      e.PhotoPath,
		Status:         e.Status,
		Description:    desc,
		DepositBalance: e.DepositBalance,
		ClassName:      e.ClassName,
		MajorName:      e.MajorName,
		ParentName:     e.ParentName,
	}
}

func EntityListToResponse(list []entity.Student) []dto.StudentResponse {
	result := make([]dto.StudentResponse, len(list))
	for i, e := range list {
		result[i] = EntityToResponse(e)
	}
	return result
}

// nilIfEmpty returns nil for empty string, pointer otherwise.
func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// intNilIfZero returns nil for 0, pointer otherwise.
func intNilIfZero(v int) *uint16 {
	if v == 0 {
		return nil
	}
	u := uint16(v)
	return &u
}
