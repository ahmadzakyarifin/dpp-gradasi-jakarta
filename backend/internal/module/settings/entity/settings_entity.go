package entity

// Settings adalah representasi domain dari tabel settings (single row).
type Settings struct {
	ID                 uint
	SiteName           string
	Tagline            string
	LogoPath           string
	ContactEmail       string
	ContactPhone       string
	Address            string
	MapsEmbedURL       string
	FacebookURL        string
	InstagramURL       string
	YoutubeURL         string
	VideoProfilePath   string
	History            string
	AboutTutorial      string
	AboutFormationDate string
	AboutNoSK          string
	AboutVision        string
	AboutMission       string // JSON string di DB, di-unmarshal ke []string saat response
	GreetingTitle      string
	GreetingSubtitle   string
	GreetingDate       string
	GreetingContent    string
	GreetingImagePath  string
	LoginHeroTitle     string
	LoginHeroDescription string
	LogRetentionDays     int
	CreatedAt            string
	UpdatedAt            string
	UpdatedBy            *uint
}
