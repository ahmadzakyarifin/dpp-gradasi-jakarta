package model

type Settings struct {
	ID                 uint   `gorm:"primaryKey;autoIncrement"`
	SiteName           string `gorm:"size:100;not null"`
	Tagline            string `gorm:"size:200"`
	LogoPath           string `gorm:"size:200"`
	ContactEmail       string `gorm:"size:100;not null"`
	ContactPhone       string `gorm:"size:20;not null"`
	Address            string `gorm:"size:255"`
	MapsEmbedURL       string `gorm:"size:500"`
	FacebookURL        string `gorm:"size:200"`
	InstagramURL       string `gorm:"size:200"`
	YoutubeURL         string `gorm:"size:200"`
	VideoProfilePath   string `gorm:"size:500"`
	History            string `gorm:"type:text"`
	AboutTutorial      string `gorm:"type:text"`
	AboutFormationDate string `gorm:"size:50"`
	AboutNoSK          string `gorm:"size:100"`
	AboutVision        string `gorm:"type:text"`
	AboutMission       string `gorm:"type:text"`
	GreetingTitle      string `gorm:"size:255"`
	GreetingSubtitle   string `gorm:"size:255"`
	GreetingDate       string `gorm:"size:100"`
	GreetingContent    string `gorm:"type:text"`
	GreetingImagePath    string `gorm:"size:500"`
	LoginHeroTitle       string `gorm:"size:255"`
	LoginHeroDescription string `gorm:"type:text"`
	LogRetentionDays     int    `gorm:"column:log_retention_days;default:30"`
	CreatedAt            string
	UpdatedAt            string
	UpdatedBy            *uint `gorm:"column:updated_by"`
}
