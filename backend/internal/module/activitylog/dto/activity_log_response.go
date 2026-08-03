package dto

type ActivityLogSummaryRes struct {
	TotalLogs     int64 `json:"total_logs"`
	HighRisk      int64 `json:"high_risk"`
	FailedLogin   int64 `json:"failed_login"`
	CMSAction     int64 `json:"cms_action"`
	FinanceAction int64 `json:"finance_action"`
}

type ActivityLogPaginationRes struct {
	CurrentPage int   `json:"current_page"`
	Limit       int   `json:"limit"`
	TotalData   int64 `json:"total_data"`
	TotalPages  int   `json:"total_pages"`
}

type ActivityLogItemRes struct {
	ID          uint64         `json:"id"`
	ActorID     *uint          `json:"actor_id"`
	ActorName   string         `json:"actor_name"`
	ActorRole   string         `json:"actor_role"`
	Action      string         `json:"action"`
	EntityType  string         `json:"entity_type"`
	EntityID    *uint          `json:"entity_id"`
	EntityLabel string         `json:"entity_label"`
	RiskLevel   string         `json:"risk_level"`
	Description string         `json:"description"`
	IPAddress   string         `json:"ip_address"`
	UserAgent   string         `json:"user_agent"`
	Metadata    map[string]any `json:"metadata"`
	CreatedAt   string         `json:"created_at"`
}

type ActivityLogDetailRes struct {
	ID          uint64         `json:"id"`
	ActorID     *uint          `json:"actor_id"`
	ActorName   string         `json:"actor_name"`
	ActorRole   string         `json:"actor_role"`
	Action      string         `json:"action"`
	EntityType  string         `json:"entity_type"`
	EntityID    *uint          `json:"entity_id"`
	EntityLabel string         `json:"entity_label"`
	RiskLevel   string         `json:"risk_level"`
	Description string         `json:"description"`
	IPAddress   string         `json:"ip_address"`
	UserAgent   string         `json:"user_agent"`
	Metadata    map[string]any `json:"metadata"`
	CreatedAt   string         `json:"created_at"`
}

type ActivityLogListRes struct {
	Summary ActivityLogSummaryRes    `json:"summary"`
	Items   []ActivityLogItemRes     `json:"items"`
	Meta    ActivityLogPaginationRes `json:"meta"`
}
