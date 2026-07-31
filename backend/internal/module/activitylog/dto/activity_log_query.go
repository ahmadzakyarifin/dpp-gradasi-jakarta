package dto

type ActivityLogQueryReq struct {
	Search string `form:"search"`
	Action string `form:"action"`
	Entity string `form:"entity"`
	Role   string `form:"role"`
	Risk   string `form:"risk"`

	// Alias contract: entity_type & risk_level (JSONC pakai nama ini)
	EntityType string `form:"entity_type"`
	RiskLevel  string `form:"risk_level"`

	ActorID   *uint  `form:"actor_id"`
	Status    string `form:"status"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	SortBy    string `form:"sort_by"`
	Order     string `form:"order"`

	Page  int `form:"page,default=1"`
	Limit int `form:"limit,default=15"`
	// Alias contract: per_page
	PerPage int `form:"per_page"`
}
