package dto

type RoleCreateReq struct {
	Name        string `json:"name" binding:"required,max=50"`
	DisplayName string `json:"display_name" binding:"required,max=100"`
	IsActive    bool   `json:"is_active"`
}

type RoleUpdateReq struct {
	Name        string `json:"name" binding:"required,max=50"`
	DisplayName string `json:"display_name" binding:"required,max=100"`
	IsActive    bool   `json:"is_active"`
}

type RoleStatusReq struct {
	IsActive bool `json:"is_active"`
}

type RoleQueryReq struct {
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
	Search string `form:"search"`
	Status string `form:"status"`
	Sort   string `form:"sort"`
}

func (q *RoleQueryReq) Normalize() {
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
	IDs []uint `json:"ids" binding:"required,min=1"`
}
