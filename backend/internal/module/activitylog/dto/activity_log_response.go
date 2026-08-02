package dto

type ActivityLogSummaryRes struct {
	TotalLogs     int64 `json:"totalLogs"`
	HighRisk      int64 `json:"highRisk"`
	FailedLogin   int64 `json:"failedLogin"`
	FinanceAction int64 `json:"financeAction"`
}

type ActivityLogPaginationRes struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
}

type ActivityLogItemRes struct {
	ID string `json:"id"`

	Time string `json:"time"`

	Actor string `json:"actor"`
	Role  string `json:"role"`

	Action string `json:"action"`

	Entity string `json:"entity"`

	EntityLabel string `json:"entityLabel"`

	Description string `json:"description"`

	IP string `json:"ip"`

	Device string `json:"device"`

	Risk string `json:"risk"`

	Metadata map[string]any `json:"metadata"`
}

type ActivityLogDetailRes struct {
	ID string `json:"id"`

	Time string `json:"time"`

	Actor string `json:"actor"`
	Role  string `json:"role"`

	Action string `json:"action"`

	Entity string `json:"entity"`

	EntityLabel string `json:"entityLabel"`

	Description string `json:"description"`

	IP string `json:"ip"`

	Device string `json:"device"`

	Risk string `json:"risk"`

	Metadata map[string]any `json:"metadata"`
}

type ActivityLogListRes struct {
	Summary ActivityLogSummaryRes `json:"summary"`

	Pagination ActivityLogPaginationRes `json:"pagination"`

	Items []ActivityLogItemRes `json:"items"`
}
