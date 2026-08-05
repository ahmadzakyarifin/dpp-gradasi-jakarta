package dto

type DashboardMessage struct {
	ID        uint   `json:"id"`
	Nama      string `json:"nama"`
	Subjek    string `json:"subjek"`
	CreatedAt string `json:"created_at"`
	IsRead    bool   `json:"is_read"`
}

type DashboardSummaryRes struct {
	TotalBerita    int64              `json:"total_berita"`
	TotalKegiatan  int64              `json:"total_kegiatan"`
	TotalPengurus  int64              `json:"total_pengurus"`
	TotalKontak    int64              `json:"total_kontak"`
	UnreadKontak   int64              `json:"unread_kontak"`
	TotalAdmin     int64              `json:"total_admin"`
	PendingAdmin   int64              `json:"pending_admin"`
	ActivityLogs   int64              `json:"activity_logs"`
	HighRiskLogs   int64              `json:"high_risk_logs"`
	FailedLogin    int64              `json:"failed_login"`
	CMSActions     int64              `json:"cms_actions"`
	LatestMessages []DashboardMessage `json:"latest_messages"`
}
