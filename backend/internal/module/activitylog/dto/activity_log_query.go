package dto

type ActivityLogQueryReq struct {
	Search    string `form:"search"`
	Action    string `form:"action"`
	Entity    string `form:"entity"`
	Role      string `form:"role"`
	Risk      string `form:"risk"`
	ActorID   uint   `form:"actor_id"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	SortBy    string `form:"sort_by"`
	Order     string `form:"order"`

	Page  int `form:"page,default=1"`
	Limit int `form:"limit,default=15"`
}

var ActivityLogSortWhitelist = map[string]bool{
	"created_at":  true,
	"id":          true,
	"actor_name":  true,
	"action":      true,
	"entity_type": true,
	"risk_level":  true,
}
