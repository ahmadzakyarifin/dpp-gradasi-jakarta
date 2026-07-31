package service

var highRisk = map[string]struct{}{
	// Auth
	"auth.login_failed": {},

	// Users & Roles
	"user.delete":          {},
	"user.change_password": {},
}

var mediumRisk = map[string]struct{}{
	// Auth
	"auth.forgot_password": {},
	"auth.reset_password":  {},

	// Users & Roles
	"user.create":         {},
	"user.update_profile": {},
	"user.verify_email":   {},

	// Settings
	"settings.update": {},

	// CMS - Deletions (Bisa dianggap medium karena menghilangkan data dari publik)
	"berita.delete":        {},
	"berita.bulk_delete":   {},
	"kegiatan.delete":      {},
	"kegiatan.bulk_delete": {},
	"slider.delete":        {},
	"slider.bulk_delete":   {},
	"pengurus.delete":      {},
	"pengurus.bulk_delete": {},
	"kontak.delete":        {},
	"kontak.bulk_delete":   {},
	"user.bulk_delete":     {},
}

var lowRisk = map[string]struct{}{
	// Auth
	"auth.login":   {},
	"auth.logout":  {},
	"auth.refresh": {},

	// CMS - Create & Update (Rutinitas normal)
	"berita.create":         {},
	"berita.update":         {},
	"berita.restore":        {},
	"berita.bulk_restore":   {},
	"kegiatan.create":       {},
	"kegiatan.update":       {},
	"kegiatan.restore":      {},
	"kegiatan.bulk_restore": {},
	"slider.create":         {},
	"slider.update":         {},
	"slider.restore":        {},
	"slider.bulk_restore":   {},
	"pengurus.create":       {},
	"pengurus.update":       {},
	"pengurus.restore":      {},
	"pengurus.bulk_restore": {},
	"kontak.restore":        {},
	"kontak.bulk_restore":   {},
	"user.restore":          {},
	"user.bulk_restore":     {},
}

func determineRisk(action string) string {
	if _, ok := highRisk[action]; ok {
		return "high"
	}
	if _, ok := mediumRisk[action]; ok {
		return "medium"
	}
	if _, ok := lowRisk[action]; ok {
		return "low"
	}

	// Fallback untuk aksi lain yang belum terdaftar
	return "low"
}
