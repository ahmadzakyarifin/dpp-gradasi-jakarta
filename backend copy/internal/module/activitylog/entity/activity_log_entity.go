package entity

import "time"

type ActivityLog struct {
	ID uint64

	ActorID   *uint
	ActorName string
	ActorRole string

	Action string

	EntityType  string
	EntityID    *uint
	EntityLabel string

	RiskLevel string

	Description string

	IPAddress string
	UserAgent string

	Metadata  map[string]any
	CreatedAt time.Time
}
