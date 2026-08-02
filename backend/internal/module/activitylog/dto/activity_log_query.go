package dto

type ActivityLogQueryReq struct {
	Search string `form:"search"`
	Action string `form:"action"`
	Entity string `form:"entity"`
	Role   string `form:"role"`
	Risk   string `form:"risk"`

	Page  int `form:"page,default=1"`
	Limit int `form:"limit,default=15"`
}
