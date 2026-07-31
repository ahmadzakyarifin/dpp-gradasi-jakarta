package main

import (
	"log"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/config"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/infrastructure"
	beritaModel "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/berita/model"
	kegiatanModel "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kegiatan/model"
	kontakModel "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kontak/model"
	pengurusModel "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/pengurus/model"
	roleModel "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/role/model"
	settingsModel "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/settings/model"
	slidersModel "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/sliders/model"
	userModel "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/model"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	log.Println("Memulai Database Seeder...")

	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("No .env file found in ../, relying on environment variables")
	}

	cfg := config.MustLoad()

	db, err := infrastructure.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Gagal konek DB: %v", err)
	}

	// 1. Drop Settings dan AutoMigrate untuk memastikan schema settings yang baru
	log.Println("Melakukan migrasi tabel Settings...")
	db.Exec("DROP TABLE IF EXISTS settings;")
	err = db.AutoMigrate(
		&settingsModel.Settings{},
	)
	if err != nil {
		log.Fatalf("Gagal automigrate settings: %v", err)
	}

	db.Exec("ALTER TABLE sliders ADD COLUMN IF NOT EXISTS deleted_at DATETIME(3) NULL")
	db.Exec("ALTER TABLE pesan_kontak ADD COLUMN IF NOT EXISTS deleted_at DATETIME(3) NULL")

	// 2. Truncate Tables
	log.Println("Menghapus data lama...")
	tables := []string{"refresh_tokens", "password_reset_tokens", "activation_tokens", "users", "roles", "berita", "kegiatan", "pengurus", "sliders", "pesan_kontak", "activity_logs"}
	for _, t := range tables {
		db.Exec("SET FOREIGN_KEY_CHECKS = 0;")
		db.Exec("TRUNCATE TABLE " + t)
		db.Exec("SET FOREIGN_KEY_CHECKS = 1;")
	}

	// 3. Insert Roles
	log.Println("Menyisipkan Roles...")
	roles := []roleModel.Role{
		{Name: "Super Admin", Description: "Akses penuh ke semua modul"},
		{Name: "Admin", Description: "Akses pengelolaan konten"},
	}
	db.Create(&roles)

	// 4. Insert Super Admin
	log.Println("Menyisipkan Super Admin...")
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	admin := userModel.User{
		RoleID:    roles[0].ID,
		Name:      "Super Admin",
		Email:     "admin@gradasi.org",
		Password:  string(hash),
		Status:    "active",
		PhotoPath: "https://ui-avatars.com/api/?name=Super+Admin&background=random",
	}
	db.Create(&admin)

	// 5. Insert Settings
	log.Println("Menyisipkan Settings...")
	setting := settingsModel.Settings{
		SiteName:           "DPP GRADASI",
		Tagline:            "Generasi Digital Indonesia",
		LogoURL:            "/uploads/logo.png",
		ContactEmail:       "dpp@gradasi.org",
		ContactPhone:       "+6285279880008",
		Address:            "Jl. Jenderal Sudirman No.1, Jakarta Pusat",
		MapsEmbedURL:       "https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d126920.24075052957!2d106.75871239853926!3d-6.22974649774619!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x2e69f3e945e34b9d%3A0x100c5e82dd4b820!2sJakarta!5e0!3m2!1sid!2sid!4v1700000000000!5m2!1sid!2sid",
		FacebookURL:        "https://www.facebook.com/gradasiofficial.id",
		InstagramURL:       "https://www.instagram.com/dppgradasi",
		YoutubeURL:         "https://www.youtube.com/channel/UCwdjB4LkqcF4Kw5-PoyOb5A",
		VideoProfileURL:    "https://www.youtube.com/embed/dQw4w9WgXcQ",
		History:            "Perkumpulan Generasi Digital Indonesia (GRADASI) didirikan pada tahun 2018 di Yogyakarta...",
		AboutTutorial:      "Pada tahun 2018 di Yogyakarta...",
		AboutFormationDate: "4 Februari 2019",
		AboutNoSK:          "AHU – 0000151.AH.01.07.2019",
		AboutVision:        "Mewujudkan masyarakat Indonesia yang cerdas dan berdaulat di era digital.",
		AboutMission:       "1. Membangun ekosistem literasi digital\n2. Meningkatkan kecakapan digital masyarakat",
		GreetingTitle:      "Tahun Baru 2026",
		GreetingSubtitle:   "Resolusi & Harapan",
		GreetingDate:       "11 February 2026",
		GreetingContent:    "Memasuki tahun 2026, GRADASI menetapkan pilar utama perjuangan: memastikan setiap masyarakat memiliki kecakapan digital (digital skills), serta mengembangkan program literasi yang berdampak nyata.",
		GreetingImageURL:   "https://gradasi.org/uploads/img/event-terkini/1767154211.jpg",
	}
	db.Create(&setting)

	strPtr := func(s string) *string { return &s }

	// 6. Insert Sliders
	log.Println("Menyisipkan Sliders...")
	sliders := []slidersModel.Slider{
		{Title: "Musyawarah Nasional Ke-II GRADASI", Subtitle: strPtr("Kolaborasi Membangun Negeri..."), Tag: strPtr("HEADLINE EVENT"), ImageURL: "https://images.unsplash.com/photo-1550751827-4bd374c3f58b?auto=format&fit=crop&q=80&w=1920", SortOrder: 1, IsActive: true, EventDate: strPtr("10 - 12 Agustus 2026"), Location: strPtr("Jakarta Convention Center"), LinkURL: strPtr("/kegiatan/munas-ke-ii")},
		{Title: "Program 1 Juta Talenta Digital", Subtitle: strPtr("Menyiapkan SDM Unggul"), Tag: strPtr("PROGRAM UNGGULAN"), ImageURL: "https://images.unsplash.com/photo-1531482615713-2afd69097998?auto=format&fit=crop&q=80&w=1920", SortOrder: 2, IsActive: true, EventDate: strPtr("Sepanjang 2026"), Location: strPtr("Seluruh Indonesia"), LinkURL: strPtr("#")},
		{Title: "Pelatihan Instruktur IT", Subtitle: strPtr("Tingkatkan Kualitas Pengajar"), Tag: strPtr("INFO KEGIATAN"), ImageURL: "https://images.unsplash.com/photo-1515162816999-a0c47dc192f7?auto=format&fit=crop&q=80&w=1920", SortOrder: 3, IsActive: true, EventDate: strPtr("20 September 2026"), Location: strPtr("Bandung"), LinkURL: strPtr("#")},
	}
	db.Create(&sliders)

	// 7. Insert Berita
	log.Println("Menyisipkan Berita...")
	berita := []beritaModel.Berita{
		{Title: "Rapat Kerja Daerah Jatim", Slug: "rapat-kerja-daerah-jatim", Category: "Berita Daerah", AuthorID: &admin.ID, Content: strPtr("SURABAYA, Generasi Digital Indonesia (GRADASI) Jawa Timur bersiap menggelar Rapat Kerja Daerah..."), Excerpt: strPtr("SURABAYA, Generasi Digital Indonesia (GRADASI) Jawa Timur bersiap menggelar Rapat Kerja Da..."), ImageURL: strPtr("https://images.unsplash.com/photo-1575320293158-b19e9334ccb8?auto=format&fit=crop&q=80&w=200"), Views: 120, IsPublished: true, PublishedDate: "2026-02-11"},
		{Title: "Rapat Strategis Pengurus Pusat", Slug: "rapat-strategis-pengurus-pusat", Category: "Berita Utama", AuthorID: &admin.ID, Content: strPtr("Pusat pelaporan kegiatan dalam rangka mempersiapkan agenda strategis organisasi untuk tahun 2026."), Excerpt: strPtr("Pusat pelaporan kegiatan dalam rangka mempersiapkan agenda strategis organisasi untuk tahun 2026."), ImageURL: strPtr("https://images.unsplash.com/photo-1517245386807-bb43f82c33c4?auto=format&fit=crop&q=80&w=200"), Views: 450, IsPublished: true, PublishedDate: "2025-10-20"},
	}
	db.Create(&berita)

	// 8. Insert Kegiatan
	log.Println("Menyisipkan Kegiatan...")
	kegiatan := []kegiatanModel.Kegiatan{
		{Title: "Penyaluran Bantuan Kemanusiaan", Slug: "penyaluran-bantuan-kemanusiaan", Location: "Aceh, Indonesia", ImageURL: strPtr("https://images.unsplash.com/photo-1511578314322-379afb476865?auto=format&fit=crop&q=80&w=200"), IsPublished: true, Category: "Sosial", Organizer: "DPP GRADASI", Content: strPtr("Penyaluran bantuan"), EventDate: "2026-08-15 08:00:00"},
		{Title: "Workshop Pemrograman Web Lanjut", Slug: "workshop-pemrograman-web-lanjut", Location: "Bandung", ImageURL: strPtr("https://images.unsplash.com/photo-1517048676732-d65bc937f952?auto=format&fit=crop&q=80&w=200"), IsPublished: true, Category: "Edukasi", Organizer: "DPD Jawa Barat", Content: strPtr("Workshop web programming"), EventDate: "2026-08-20 09:00:00"},
	}
	db.Create(&kegiatan)

	// 9. Insert Pengurus
	log.Println("Menyisipkan Pengurus...")
	pengurus := []pengurusModel.Pengurus{
		{Name: "Upi Asmaradhana", Role: "Ketua Umum DPP", Level: "ketua", ImageURL: "https://ui-avatars.com/api/?name=Upi+Asmaradhana&background=random", SortOrder: 1, IsActive: true, Periode: "2025 - 2030"},
		{Name: "Kusbeni Abdulloh, S.Kom.", Role: "Ketua DPD Gradasi Jawa Timur", Level: "dpd", ImageURL: "https://ui-avatars.com/api/?name=Kusbeni+Abdulloh&background=random", SortOrder: 2, IsActive: true, Periode: "2025 - 2030", Provinsi: strPtr("Jawa Timur")},
		{Name: "Andi Setiawan", Role: "Sekretaris Jenderal", Level: "dpp", ImageURL: "https://ui-avatars.com/api/?name=Andi+Setiawan&background=random", SortOrder: 3, IsActive: true, Periode: "2025 - 2030"},
	}
	db.Create(&pengurus)

	// 10. Insert Kontak
	log.Println("Menyisipkan Kontak...")
	kontak := []kontakModel.PesanKontak{
		{Nama: "Ahmad Zaky", Email: "zaky@example.com", Subjek: "Kerjasama", Pesan: "Halo DPP GRADASI, kami ingin mengajukan kerjasama...", IsRead: false},
		{Nama: "Siti Rahma", Email: "siti@example.com", Subjek: "Daftar Anggota", Pesan: "Bagaimana cara mendaftar menjadi anggota DPD?", IsRead: true},
	}
	db.Create(&kontak)

	log.Println("Seeding Database Selesai!")
}
