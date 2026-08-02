package dto

type GuardianQueryReq struct {
	Page   int    `form:"page,default=1"`
	Limit  int    `form:"limit,default=10"`
	Search string `form:"search"`
}

type GuardianCreateReq struct {
	Name        string `json:"name" binding:"required,min=2"`
	Phone       string `json:"phone" binding:"omitempty"`
	Email       string `json:"email" binding:"omitempty,email"`
	NIK         string `json:"nik" binding:"omitempty"`
	Education   string `json:"education"`
	Occupation  string `json:"occupation"`
	IncomeRange string `json:"income_range"`
}

type GuardianUpdateReq struct {
	Name        string `json:"name" binding:"required,min=2"`
	Phone       string `json:"phone"`
	Email       string `json:"email" binding:"omitempty,email"`
	NIK         string `json:"nik"`
	Education   string `json:"education"`
	Occupation  string `json:"occupation"`
	IncomeRange string `json:"income_range"`
}
