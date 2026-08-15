package dto

type SettingsResponse struct {
	ID                 uint     `json:"id"`
	SiteName           string   `json:"site_name"`
	Tagline            string   `json:"tagline"`
	LogoPath           string   `json:"logo_path"`
	ContactEmail       string   `json:"contact_email"`
	ContactPhone       string   `json:"contact_phone"`
	Address            string   `json:"address"`
	MapsEmbedURL       string   `json:"maps_embed_url"`
	FacebookURL        string   `json:"facebook_url"`
	InstagramURL       string   `json:"instagram_url"`
	YoutubeURL         string   `json:"youtube_url"`
	VideoProfilePath   string   `json:"video_profile_path"`
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
	GreetingImagePath  string   `json:"greeting_image_path"`
	GreetingSignName   string   `json:"greeting_sign_name"`
	GreetingSignSubtitle string `json:"greeting_sign_subtitle"`
	GreetingSignImage1   string `json:"greeting_sign_image_1"`
	GreetingSignImage2   string `json:"greeting_sign_image_2"`
	LoginHeroTitle     string   `json:"login_hero_title"`
	LoginHeroDescription string   `json:"login_hero_description"`
	LogRetentionDays     int      `json:"log_retention_days"`
	CreatedAt            string   `json:"created_at"`
	UpdatedAt            string   `json:"updated_at,omitempty"`
	UpdatedBy            *uint    `json:"updated_by,omitempty"`

	// Status CAPTCHA dari config backend — dibaca frontend (single source of truth)
	CaptchaEnabled bool   `json:"captcha_enabled,omitempty"`
	CaptchaSiteKey string `json:"captcha_site_key,omitempty"`
}
