package service

var highRisk = map[string]struct{}{
	// Auth
	"auth.login_failed":         {},
	"auth.forgot_password_spam": {},

	// Users & Roles
	"users.delete":      {},
	"users.bulk_delete": {},
	"users.change_role": {},
	"roles.update":      {},
	"roles.delete":      {},
	"roles.bulk_delete": {},

	// Students
	"students.delete":      {},
	"students.bulk_delete": {},

	// Academic
	"majors.delete":               {},
	"majors.bulk_delete":          {},
	"classes.delete":              {},
	"classes.bulk_delete":         {},
	"class_templates.delete":      {},
	"class_templates.bulk_delete": {},
	"cohorts.delete":              {},
	"cohorts.bulk_delete":         {},
	"academic_years.delete":       {},
	"academic_years.bulk_delete":  {},
	"semesters.delete":            {},
	"semesters.bulk_delete":       {},

	// Finance
	"bills.write_off":      {},
	"billing_rules.delete": {},
	"bill_types.delete":    {},
	"payments.refund":      {},
	"payments.correction":  {},

	// System & Import
	"smart_import.import": {},
}

var mediumRisk = map[string]struct{}{
	// Auth
	"auth.reset_password": {},

	// Users & Roles
	"users.update":        {},
	"users.status_update": {},
	"users.restore":       {},
	"users.bulk_restore":  {},
	"roles.status_update": {},
	"roles.restore":       {},
	"roles.bulk_restore":  {},

	// Students
	"students.update":        {},
	"students.promote":       {},
	"students.graduate":      {},
	"students.status_update": {},
	"students.restore":       {},
	"students.bulk_restore":  {},

	// Academic
	"majors.update":                 {},
	"majors.status_update":          {},
	"majors.restore":                {},
	"majors.bulk_restore":           {},
	"classes.update":                {},
	"classes.status_update":         {},
	"classes.restore":               {},
	"classes.bulk_restore":          {},
	"class_templates.update":        {},
	"class_templates.status_update": {},
	"class_templates.restore":       {},
	"class_templates.bulk_restore":  {},
	"cohorts.update":                {},
	"cohorts.status_update":         {},
	"cohorts.restore":               {},
	"cohorts.bulk_restore":          {},
	"academic_years.update":         {},
	"academic_years.status_update":  {},
	"academic_years.restore":        {},
	"academic_years.bulk_restore":   {},
	"semesters.update":              {},
	"semesters.status_update":       {},
	"semesters.restore":             {},
	"semesters.bulk_restore":        {},
	"active_classes.update":         {},

	// Finance
	"bills.update":                {},
	"bills.waive":                 {},
	"billing_rules.update":        {},
	"billing_rules.status_update": {},
	"bill_types.update":           {},
	"bill_types.status_update":    {},

	// System
	"notifications.resend":     {},
	"support.takeover_ticket":  {},
	"support.close_ticket":     {},
	"support_templates.update": {},
	"support_templates.delete": {},
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
