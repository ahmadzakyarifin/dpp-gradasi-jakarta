package middleware

import "time"

type LimitKind string

const (
	KindIP      LimitKind = "ip"
	KindEmail   LimitKind = "email"
	KindIPEmail LimitKind = "ip_email"
	KindUser    LimitKind = "user"
)

type Rule struct {
	Name   string
	Kind   LimitKind
	Limit  int64
	Period time.Duration
}

func IP(limit int64, period time.Duration) Rule {
	return newRule("ip", KindIP, limit, period)
}

func Email(limit int64, period time.Duration) Rule {
	return newRule("email", KindEmail, limit, period)
}

func EmailNamed(name string, limit int64, period time.Duration) Rule {
	return newRule(name, KindEmail, limit, period)
}

func IPEmail(limit int64, period time.Duration) Rule {
	return newRule("ip_email", KindIPEmail, limit, period)
}

func User(limit int64, period time.Duration) Rule {
	return newRule("user", KindUser, limit, period)
}

func UserNamed(name string, limit int64, period time.Duration) Rule {
	return newRule(name, KindUser, limit, period)
}

func newRule(
	name string,
	kind LimitKind,
	limit int64,
	period time.Duration,
) Rule {
	if limit <= 0 {
		limit = 1
	}

	if period <= 0 {
		period = time.Minute
	}

	return Rule{
		Name:   name,
		Kind:   kind,
		Limit:  limit,
		Period: period,
	}
}
