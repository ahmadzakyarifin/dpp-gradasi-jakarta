package dto

type SettingsResponse struct {
	ID                 uint     `json:"id"`
	SiteName           string   `json:"site_name"`
	Tagline            string   `json:"tagline"`
	LogoURL            string   `json:"logo_url"`
	ContactEmail       string   `json:"contact_email"`
	ContactPhone       string   `json:"contact_phone"`
	Address            string   `json:"address"`
	MapsEmbedURL       string   `json:"maps_embed_url"`
	FacebookURL        string   `json:"facebook_url"`
	InstagramURL       string   `json:"instagram_url"`
	YoutubeURL         string   `json:"youtube_url"`
	VideoProfileURL    string   `json:"video_profile_url"`
	History            string   `json:"history"`
	AboutTutorial      string   `json:"about_tutorial"`
	AboutFormationDate string   `json:"about_formation_date"`
	AboutNoSK          string   `json:"about_no_sk"`
	AboutVision        string   `json:"about_vision"`
	AboutMission       []string `json:"about_mission"`
	GreetingTitle      string   `json:"greeting_title"`
	GreetingSubtitle   string   `json:"greeting_subtitle"`
	GreetingDate       string   `json:"greeting_date"`
	GreetingContent    string   `json:"greeting_content"`
	GreetingImageURL   string   `json:"greeting_image_url"`
	CreatedAt          string   `json:"created_at"`
	UpdatedAt          string   `json:"updated_at,omitempty"`
	UpdatedBy          *uint    `json:"updated_by,omitempty"`
}
