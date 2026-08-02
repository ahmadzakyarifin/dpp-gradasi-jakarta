package dto

type MajorCreateReq struct {
	Code        string  `json:"code"`
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}

type MajorUpdateReq struct {
	Code        string  `json:"code"`
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}

type MajorStatusReq struct {
	IsActive bool `json:"is_active"`
}

type MajorQueryReq struct {
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
	Search string `form:"search"`
	Status string `form:"status"`
	Sort   string `form:"sort"`
}

func (q *MajorQueryReq) Normalize() {
	if q.Page <= 0 {
		q.Page = 1
	}

	if q.Limit <= 0 {
		q.Limit = 10
	}

	if q.Limit > 100 {
		q.Limit = 100
	}
}

type BulkDeleteReq struct {
	IDs []uint `json:"ids" binding:"required,min=1,max=100"`
}
