package service

// highRisk: aksi destruktif / sensitif — menghapus data, mengubah hak akses.
var highRisk = map[string]struct{}{
	// Auth
	"auth.login_failed":           {},
	"auth.forgot_password_spam":   {},
	"auth.reset_password":         {},
	"auth.change_password":        {},
	"auth.reset_password_failed":  {},
	"auth.change_password_failed": {},

	// Users
	"users.delete":        {},
	"users.bulk_delete":   {},
	"users.toggle_status": {},
	"users.change_role":   {},

	// Roles
	"roles.update":      {},
	"roles.delete":      {},
	"roles.bulk_delete": {},
}

// mediumRisk: aksi perubahan data yang berdampak sedang.
var mediumRisk = map[string]struct{}{
	// Auth
	"auth.login":           {},
	"auth.refresh_token":   {},
	"auth.logout":          {},
	"auth.forgot_password": {},

	// Users
	"users.create":            {},
	"users.update":            {},
	"users.restore":           {},
	"users.bulk_restore":      {},
	"users.activate":          {},
	"users.update_profile":    {},
	"users.resend_activation": {},

	// Roles
	"roles.create":        {},
	"roles.restore":       {},
	"roles.bulk_restore":  {},
	"roles.status_update": {},

	// CMS — Berita
	"berita.create":       {},
	"berita.update":       {},
	"berita.delete":       {},
	"berita.restore":      {},
	"berita.bulk_delete":  {},
	"berita.bulk_restore": {},

	// CMS — Kegiatan
	"kegiatan.create":       {},
	"kegiatan.update":       {},
	"kegiatan.delete":       {},
	"kegiatan.restore":      {},
	"kegiatan.bulk_delete":  {},
	"kegiatan.bulk_restore": {},

	// CMS — Kontak
	"kontak.update":       {},
	"kontak.delete":       {},
	"kontak.restore":      {},
	"kontak.bulk_delete":  {},
	"kontak.bulk_restore": {},

	// CMS — Pengurus
	"pengurus.create":       {},
	"pengurus.update":       {},
	"pengurus.delete":       {},
	"pengurus.restore":      {},
	"pengurus.bulk_delete":  {},
	"pengurus.bulk_restore": {},

	// CMS — Settings
	"settings.update": {},

	// CMS — Sliders (aksi di service memakai prefix "slider")
	"slider.create":       {},
	"slider.update":       {},
	"slider.delete":       {},
	"slider.restore":      {},
	"slider.bulk_delete":  {},
	"slider.bulk_restore": {},
}

func determineRisk(action string) string {

	if _, ok := highRisk[action]; ok {
		return "high"
	}

	if _, ok := mediumRisk[action]; ok {
		return "medium"
	}

	return "low"
}
