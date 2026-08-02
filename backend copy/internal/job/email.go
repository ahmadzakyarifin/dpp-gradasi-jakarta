package jobs

type EmailJob struct {
	NotificationID *uint  `json:"notification_id,omitempty"`
	To             string `json:"to"`
	Subject        string `json:"subject"`
	HTML           string `json:"html"`
}
