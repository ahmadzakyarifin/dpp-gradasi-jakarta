package model

type Settings struct {
	ID                 uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	SiteName           string `gorm:"size:100;not null" json:"site_name"`
	Tagline            string `gorm:"size:200" json:"tagline"`
	LogoPath           string `gorm:"size:200" json:"logo_path"`
	ContactEmail       string `gorm:"size:100;not null" json:"contact_email"`
	ContactPhone       string `gorm:"size:20;not null" json:"contact_phone"`
	Address            string `gorm:"size:255" json:"address"`
	MapsEmbedURL       string `gorm:"size:500" json:"maps_embed_url"`
	FacebookURL        string `gorm:"size:200" json:"facebook_url"`
	InstagramURL       string `gorm:"size:200" json:"instagram_url"`
	YoutubeURL         string `gorm:"size:200" json:"youtube_url"`
	VideoProfilePath   string `gorm:"size:500" json:"video_profile_path"`
	History            string `gorm:"type:text" json:"history"`
	AboutTutorial      string `gorm:"type:text" json:"about_tutorial"`
	AboutFormationDate string `gorm:"size:50" json:"about_formation_date"`
	AboutNoSK          string `gorm:"size:100" json:"about_no_sk"`
	AboutVision        string `gorm:"type:text" json:"about_vision"`
	AboutMission       string `gorm:"type:text" json:"about_mission"`
	GreetingTitle      string `gorm:"size:255" json:"greeting_title"`
	GreetingSubtitle   string `gorm:"size:255" json:"greeting_subtitle"`
	GreetingDate       string `gorm:"size:100" json:"greeting_date"`
	GreetingContent    string `gorm:"type:text" json:"greeting_content"`
	GreetingImagePath  string `gorm:"size:500" json:"greeting_image_path"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
	UpdatedBy          *uint  `gorm:"column:updated_by" json:"updated_by"`
}
