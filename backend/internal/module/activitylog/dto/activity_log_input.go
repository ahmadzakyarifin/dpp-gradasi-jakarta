package dto

type ActivityLogInput struct {
	ActorID   *uint
	ActorName string
	ActorRole string

	Action string
	Status string

	EntityType  string
	EntityID    *uint
	EntityLabel string

	Description string

	IPAddress string
	UserAgent string

	Metadata map[string]any
}
