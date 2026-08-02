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
	log.Println("=== Memulai Database Seeder DPP GRADASI ===")

	// Load .env dari root folder atau backend directory
	err := godotenv.Load("../.env")
	if err != nil {
		err = godotenv.Load(".env")
		if err != nil {
			log.Println("Perhatian: .env tidak ditemukan, menggunakan variabel lingkungan bawaan.")
		}
	}

	cfg := config.MustLoad()

	db, err := infrastructure.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Gagal koneksi ke database: %v", err)
	}

	// 1. AutoMigrate & reset schema perbaikan jika ada
	log.Println("[1/10] Melakukan migrasi & perataan skema...")
	db.Exec("ALTER TABLE sliders ADD COLUMN IF NOT EXISTS deleted_at DATETIME(3) NULL")
	db.Exec("ALTER TABLE pesan_kontak ADD COLUMN IF NOT EXISTS deleted_at DATETIME(3) NULL")
	db.Exec("ALTER TABLE kegiatan ADD COLUMN IF NOT EXISTS author_id INT NULL")
	db.Exec("ALTER TABLE kegiatan MODIFY COLUMN event_date VARCHAR(200) NULL")
	db.Exec("DROP TABLE IF EXISTS settings;")
	err = db.AutoMigrate(&settingsModel.Settings{})
	if err != nil {
		log.Fatalf("Gagal automigrate settings: %v", err)
	}

	// 2. Truncate Tables
	log.Println("[2/10] Menghapus data lama (Truncate)...")
	tables := []string{
		"refresh_tokens", "password_reset_tokens", "activation_tokens",
		"kegiatan_gallery", "kegiatan_tags", "berita_tags",
		"users", "roles", "berita", "kegiatan", "pengurus", "sliders", "pesan_kontak", "activity_logs",
	}
	db.Exec("SET FOREIGN_KEY_CHECKS = 0;")
	for _, t := range tables {
		db.Exec("TRUNCATE TABLE " + t)
	}
	db.Exec("SET FOREIGN_KEY_CHECKS = 1;")

	// Helper for string pointers
	strPtr := func(s string) *string { return &s }

	// 3. Insert Roles (nama snake_case sesuai normalisasi 00016 & kontrak middleware)
	log.Println("[3/10] Menyisipkan data Roles...")
	roles := []roleModel.RoleModel{
		{Name: "super_admin", DisplayName: "Super Administrator", IsSystem: true, IsActive: true},
		{Name: "admin", DisplayName: "Admin", IsSystem: false, IsActive: true},
		{Name: "admin_berita", DisplayName: "Admin Berita", IsSystem: false, IsActive: true},
		{Name: "admin_kegiatan", DisplayName: "Admin Kegiatan", IsSystem: false, IsActive: true},
	}
	db.Create(&roles)

	// 4. Insert Super Admin dari .env root
	adminName := cfg.Dev.SuperAdminName
	if adminName == "" {
		adminName = "Super Admin"
	}
	adminEmail := cfg.Dev.SuperAdminEmail
	if adminEmail == "" {
		adminEmail = "admin@gradasi.org"
	}
	adminPassword := cfg.Dev.SuperAdminPassword
	if adminPassword == "" {
		adminPassword = "password123"
	}

	log.Printf("[4/10] Menyisipkan Super Admin (%s | %s)...", adminName, adminEmail)
	hash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Gagal enkripsi password super admin: %v", err)
	}

	adminUser := userModel.UserModel{
		RoleID:       roles[0].ID,
		Name:         adminName,
		Email:        adminEmail,
		PasswordHash: string(hash),
		Status:       "active",
		PhotoPath:    strPtr("https://ui-avatars.com/api/?name=Super+Admin&background=0D8ABC&color=fff"),
	}
	db.Create(&adminUser)

	// 5. Insert Settings (dari index.html)
	log.Println("[5/10] Menyisipkan Settings website (Visi, Misi, Legalitas, Kontak)...")
	setting := settingsModel.Settings{
		SiteName:           "DPP GRADASI",
		Tagline:            "Generasi Digital Indonesia",
		LogoPath:           "https://gradasi.org/uploads/img/logo/1737187847.png",
		ContactEmail:       "dpp@gradasi.org",
		ContactPhone:       "+6285279880008",
		Address:            "Office Park OL3-IZA The Bellagio Mall, Mega Kuningan, Jakarta Selatan",
		MapsEmbedURL:       "https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d126920.24075052957!2d106.75871239853926!3d-6.22974649774619!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x2e69f3e945e34b9d%3A0x100c5e82dd4b820!2sJakarta!5e0!3m2!1sid!2sid!4v1700000000000!5m2!1sid!2sid",
		FacebookURL:        "https://www.facebook.com/gradasiofficial.id",
		InstagramURL:       "https://www.instagram.com/dppgradasi",
		YoutubeURL:         "https://www.youtube.com/channel/UCwdjB4LkqcF4Kw5-PoyOb5A",
		VideoProfilePath:   "https://www.youtube.com/embed/dQw4w9WgXcQ",
		History:            "Perkumpulan Generasi Digital Indonesia (GRADASI) didirikan pada 4 Februari 2019 sebagai organisasi independen yang berfokus pada pengembangan literasi digital, pemberdayaan UMKM, dan transformasi teknologi di Indonesia.",
		AboutTutorial:      "Pengesahan Badan Hukum Kemenkumham RI.",
		AboutFormationDate: "4 Februari 2019",
		AboutNoSK:          "AHU-0000151.AH.01.07.2019",
		AboutVision:        "Mewujudkan masyarakat Indonesia yang cerdas, kreatif, dan berdaulat di era digital.",
		AboutMission:       `["Membangun ekosistem literasi digital yang inklusif di seluruh daerah Indonesia.","Mengakselerasi transformasi digital bagi UMKM dan generasi muda.","Mendorong inovasi dan kolaborasi antar pemangku kepentingan industri kreatif digital."]`,
		GreetingTitle:      "Tahun Baru 2026",
		GreetingSubtitle:   "Resolusi & Harapan",
		GreetingDate:       "11 Februari 2026",
		GreetingContent:    "Memasuki tahun 2026, GRADASI menetapkan pilar utama perjuangan: memastikan setiap masyarakat memiliki kecakapan digital (digital skills), serta mengembangkan program literasi yang berdampak nyata bagi pertumbuhan ekonomi lokal.",
		GreetingImagePath:  "https://gradasi.org/uploads/img/event-terkini/1767154211.jpg",
	}
	db.Create(&setting)

	// 6. Insert Sliders Carousel (dari index.html)
	log.Println("[6/10] Menyisipkan Banner Sliders (Header Carousel)...")
	sliders := []slidersModel.Slider{
		{
			Title:     "Musyawarah Nasional Ke-II GRADASI",
			Subtitle:  strPtr("Kolaborasi Membangun Negeri Menuju Indonesia Emas 2045"),
			Tag:       strPtr("HEADLINE EVENT"),
			ImagePath: "https://images.unsplash.com/photo-1550751827-4bd374c3f58b?auto=format&fit=crop&q=80&w=1920",
			SortOrder: 1,
			IsActive:  true,
			IsNew:     true,
			EventDate: strPtr("10 - 12 Agustus 2026"),
			Location:  strPtr("Jakarta Convention Center"),
			LinkURL:   strPtr("/kegiatan/munas-ke-ii"),
		},
		{
			Title:     "Program 1 Juta Talenta Digital",
			Subtitle:  strPtr("Menyiapkan SDM Unggul Siap Kerja dan Berdaya Saing Global"),
			Tag:       strPtr("PROGRAM UNGGULAN"),
			ImagePath: "https://images.unsplash.com/photo-1531482615713-2afd69097998?auto=format&fit=crop&q=80&w=1920",
			SortOrder: 2,
			IsActive:  true,
			IsNew:     false,
			EventDate: strPtr("Sepanjang 2026"),
			Location:  strPtr("Seluruh Indonesia"),
			LinkURL:   strPtr("/kegiatan/1-juta-talenta-digital"),
		},
		{
			Title:     "Pelatihan Instruktur IT & Literasi",
			Subtitle:  strPtr("Tingkatkan Kualitas Pengajar Digital di Seluruh Pelosok Negeri"),
			Tag:       strPtr("INFO KEGIATAN"),
			ImagePath: "https://images.unsplash.com/photo-1515162816999-a0c47dc192f7?auto=format&fit=crop&q=80&w=1920",
			SortOrder: 3,
			IsActive:  true,
			IsNew:     false,
			EventDate: strPtr("20 September 2026"),
			Location:  strPtr("Bandung, Jawa Barat"),
			LinkURL:   strPtr("/kegiatan/pelatihan-instruktur-it"),
		},
	}
	db.Create(&sliders)

	// 7. Insert Berita Default (dari berita.html & index.html)
	log.Println("[7/10] Menyisipkan Data Berita Default...")
	beritaList := []beritaModel.Berita{
		{
			Title:         "Rapat Kerja Daerah Jatim",
			Slug:          "rapat-kerja-daerah-jatim",
			Category:      "Berita Daerah",
			PublishedDate: "2026-02-11",
			AuthorID:      &adminUser.ID,
			ImagePath:     strPtr("https://gradasi.org/uploads/img/berita/17708152730.jpg"),
			Excerpt:       strPtr("SURABAYA, Generasi Digital Indonesia (GRADASI) Jawa Timur bersiap menggelar Rapat Kerja Daerah..."),
			Content:       strPtr("<p>SURABAYA, Generasi Digital Indonesia (GRADASI) Jawa Timur bersiap menggelar Rapat Kerja Daerah untuk menyelaraskan program kerja digitalisasi UMKM dan literasi masyarakat di wilayah Jawa Timur.</p>"),
			IsPublished:   true,
			Views:         1250,
		},
		{
			Title:         "Peningkatan Kompetensi SDM Pendidikan",
			Slug:          "peningkatan-kompetensi-sdm-pendidikan",
			Category:      "Edukasi",
			PublishedDate: "2025-11-02",
			AuthorID:      &adminUser.ID,
			ImagePath:     strPtr("https://gradasi.org/uploads/img/berita/17620765070.jpg"),
			Excerpt:       strPtr("Inisiatif GRADASI Mendorong Peningkatan Kompetensi SDM Pendidikan dalam Memanfaatkan Kecerdasan Buatan..."),
			Content:       strPtr("<p>Inisiatif GRADASI Mendorong Peningkatan Kompetensi SDM Pendidikan dalam Memanfaatkan Kecerdasan Buatan (AI) secara bijak dan adaptif.</p>"),
			IsPublished:   true,
			Views:         890,
		},
		{
			Title:         "Rumusan Kunci Kebijakan Literasi Digital",
			Slug:          "rumusan-kunci-kebijakan-literasi-digital",
			Category:      "Berita Utama",
			PublishedDate: "2025-10-31",
			AuthorID:      &adminUser.ID,
			ImagePath:     strPtr("https://gradasi.org/uploads/img/berita/17618789900.jpg"),
			Excerpt:       strPtr("Ketua Dewan Pakar GRADASI memaparkan lima rumusan kunci kebijakan untuk mempercepat transformasi digital..."),
			Content:       strPtr("<p>Ketua Dewan Pakar GRADASI memaparkan lima rumusan kunci kebijakan untuk mempercepat transformasi digital nasional yang inklusif.</p>"),
			IsPublished:   true,
			Views:         450,
		},
	}
	for i := range beritaList {
		db.Create(&beritaList[i])
		db.Create(&beritaModel.BeritaTag{BeritaID: beritaList[i].ID, Tag: "gradasi"})
		db.Create(&beritaModel.BeritaTag{BeritaID: beritaList[i].ID, Tag: "digital"})
	}

	// 8. Insert Kegiatan Default (dari kegiatan.html & index.html)
	log.Println("[8/10] Menyisipkan Data Kegiatan Default & Galeri...")
	kegiatanList := []kegiatanModel.Kegiatan{
		{
			Title:       "Penyaluran Bantuan Kemanusiaan",
			Slug:        "penyaluran-bantuan-kemanusiaan",
			Category:    "Nasional",
			EventDate:   "31 December 2025",
			Location:    "Aceh, Indonesia",
			Organizer:   "DPP GRADASI",
			AuthorID:    &adminUser.ID,
			ImagePath:   strPtr("https://gradasi.org/uploads/img/event/1767154719.jpg"),
			Excerpt:     strPtr("Dewan Pimpinan Pusat (DPP) GRADASI menyalurkan bantuan kemanusiaan kepada korban bencana alam."),
			Content:     strPtr("<p>Dewan Pimpinan Pusat (DPP) GRADASI menyalurkan bantuan kemanusiaan kepada korban bencana alam sebagai wujud kepedulian sosial organisasi.</p>"),
			IsPublished: true,
			Views:       1500,
		},
		{
			Title:       "Pelatihan Digital Marketing UMKM",
			Slug:        "pelatihan-digital-marketing-umkm",
			Category:    "Jawa Timur",
			EventDate:   "31 December 2025",
			Location:    "Surabaya, Jawa Timur",
			Organizer:   "DPD Jawa Timur",
			AuthorID:    &adminUser.ID,
			ImagePath:   strPtr("https://gradasi.org/uploads/img/event/1767154619.jpg"),
			Excerpt:     strPtr("Program pendampingan dan pelatihan pemasaran digital gratis bagi pelaku UMKM Jawa Timur."),
			Content:     strPtr("<p>Program pendampingan dan pelatihan pemasaran digital gratis bagi pelaku UMKM Jawa Timur agar produk lokal dapat bersaing nasional.</p>"),
			IsPublished: true,
			Views:       980,
		},
		{
			Title:       "Konsolidasi Pengurus DPP & DPD",
			Slug:        "konsolidasi-pengurus-dpp-dpd",
			Category:    "Lampung",
			EventDate:   "31 December 2025",
			Location:    "Bandar Lampung",
			Organizer:   "DPD Lampung",
			AuthorID:    &adminUser.ID,
			ImagePath:   strPtr("https://gradasi.org/uploads/img/event/1767154397.jpg"),
			Excerpt:     strPtr("Rapat konsolidasi pengurus pusat dan pengurus daerah untuk memantapkan peta jalan program kerja."),
			Content:     strPtr("<p>Rapat konsolidasi pengurus pusat dan pengurus daerah untuk memantapkan peta jalan program kerja lima tahun ke depan.</p>"),
			IsPublished: true,
			Views:       620,
		},
	}
	for i := range kegiatanList {
		db.Create(&kegiatanList[i])
	}

	// Add gallery items for first activity
	if len(kegiatanList) > 0 {
		db.Create(&kegiatanModel.KegiatanGallery{
			KegiatanID: kegiatanList[0].ID,
			ImagePath:  "https://images.unsplash.com/photo-1511578314322-379afb476865?q=80&w=600",
			Caption:    "Pembukaan Munas II",
			SortOrder:  1,
		})
		db.Create(&kegiatanModel.KegiatanGallery{
			KegiatanID: kegiatanList[0].ID,
			ImagePath:  "https://images.unsplash.com/photo-1540575467063-178a50c2df87?q=80&w=600",
			Caption:    "Sidang Pleno Organisasi",
			SortOrder:  2,
		})
	}

	// 9. Insert Pengurus Default (dari kepengurusan.html)
	log.Println("[9/10] Menyisipkan Data Pengurus (Ketua, DPP, DPD, DPC)...")
	pengurusList := []pengurusModel.Pengurus{
		// Ketua Umum
		{
			Name:         "Upi Asmaradhana",
			Role:         "Ketua Umum Generasi Digital Indonesia 2025 - 2030",
			Department:   strPtr("Pimpinan Pusat"),
			Level:        "ketua",
			ImagePath:    "https://gradasi.org/uploads/img/s-anggota/ketua/1735027418.jpg",
			Periode:      "2025 - 2030",
			SortOrder:    1,
			IsActive:     true,
			FacebookURL:  strPtr("https://www.facebook.com/gradasiofficial.id"),
			InstagramURL: strPtr("https://www.instagram.com/dppgradasi"),
		},
		// DPP
		{Name: "Dr. Susi Susanti, M.Pd", Role: "Wakil Ketua I", Department: strPtr("Pengurus Harian"), Level: "dpp", ImagePath: "https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?q=80&w=200", SortOrder: 2, Periode: "2025 - 2030", IsActive: true},
		{Name: "Ir. Budi Santoso", Role: "Wakil Ketua II", Department: strPtr("Pengurus Harian"), Level: "dpp", ImagePath: "https://images.unsplash.com/photo-1519085360753-af0119f7cbe7?q=80&w=200", SortOrder: 3, Periode: "2025 - 2030", IsActive: true},
		{Name: "Junaidi, S.Kom", Role: "Sekretaris Jenderal", Department: strPtr("Kesekretariatan"), Level: "dpp", ImagePath: "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?q=80&w=200", SortOrder: 4, Periode: "2025 - 2030", IsActive: true},
		{Name: "Dina Mariana, S.ST", Role: "Wakil Sekjen 1", Department: strPtr("Kesekretariatan"), Level: "dpp", ImagePath: "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?q=80&w=200", SortOrder: 5, Periode: "2025 - 2030", IsActive: true},
		{Name: "Rina Wijaya, M.Sc", Role: "Bendahara Umum", Department: strPtr("Kebendaharaan"), Level: "dpp", ImagePath: "https://images.unsplash.com/photo-1580489944761-15a19d654956?q=80&w=200", SortOrder: 6, Periode: "2025 - 2030", IsActive: true},
		{Name: "Muhammad Hertiyadi Alfaqy S.Kom", Role: "Koordinator Dept 02 IT & Digital", Department: strPtr("Departemen 02"), Level: "dpp", ImagePath: "https://images.unsplash.com/photo-1599566150163-29194dcaad36?q=80&w=200", SortOrder: 7, Periode: "2025 - 2030", IsActive: true},
		// DPD
		{Name: "Kusbeni Abdulloh, S.Kom.", Role: "Ketua DPD Gradasi Jawa Timur", Level: "dpd", Provinsi: strPtr("Jawa Timur"), ImagePath: "https://images.unsplash.com/photo-1560250097-0b93528c311a?q=80&w=200", SortOrder: 8, Periode: "2025 - 2030", IsActive: true},
		{Name: "SHANTY OCTAVIA UTAMI, ST", Role: "Sekretaris DPD Gradasi Jawa Timur", Level: "dpd", Provinsi: strPtr("Jawa Timur"), ImagePath: "https://images.unsplash.com/photo-1580489944761-15a19d654956?q=80&w=200", SortOrder: 9, Periode: "2025 - 2030", IsActive: true},
		{Name: "Ridona", Role: "Ketua DPD Gradasi Provinsi Riau", Level: "dpd", Provinsi: strPtr("Riau"), ImagePath: "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?q=80&w=200", SortOrder: 10, Periode: "2025 - 2030", IsActive: true},
		{Name: "Safrial", Role: "Ketua Gradasi DPD Sumatera Utara", Level: "dpd", Provinsi: strPtr("Sumatera Utara"), ImagePath: "https://images.unsplash.com/photo-1492562080023-ab3db95bfbce?q=80&w=200", SortOrder: 11, Periode: "2025 - 2030", IsActive: true},
		// DPC
		{Name: "Budi Pratama, S.T.", Role: "Ketua DPC Gradasi Kab. Malang", Level: "dpc", Provinsi: strPtr("Jawa Timur"), Kabupaten: strPtr("Kabupaten Malang"), ImagePath: "https://images.unsplash.com/photo-1519085360753-af0119f7cbe7?q=80&w=200", SortOrder: 12, Periode: "2025 - 2030", IsActive: true},
		{Name: "Siti Rahmawati", Role: "Sekretaris DPC Gradasi Kota Surabaya", Level: "dpc", Provinsi: strPtr("Jawa Timur"), Kabupaten: strPtr("Kota Surabaya"), ImagePath: "https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?q=80&w=200", SortOrder: 13, Periode: "2025 - 2030", IsActive: true},
		{Name: "Ahmad Fauzi", Role: "Ketua DPC Gradasi Kab. Bogor", Level: "dpc", Provinsi: strPtr("Jawa Barat"), Kabupaten: strPtr("Kabupaten Bogor"), ImagePath: "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?q=80&w=200", SortOrder: 14, Periode: "2025 - 2030", IsActive: true},
		// ===== Tambahan dari dummy frontend (Kepengurusan.jsx) agar tidak ada data hilang =====
		// DPP (lanjutan)
		{Name: "Sudarwati", Role: "Wakil Sekjen 2", Department: strPtr("Kesekretariatan"), Level: "dpp", ImagePath: "https://images.unsplash.com/photo-1544005313-94ddf0286df2?q=80&w=200", SortOrder: 15, Periode: "2025 - 2030", IsActive: true},
		{Name: "Yoseph Budi", Role: "Wakil Bendahara", Department: strPtr("Kebendaharaan"), Level: "dpp", ImagePath: "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?q=80&w=200", SortOrder: 16, Periode: "2025 - 2030", IsActive: true},
		{Name: "Dwi Purnomo, S.Kom", Role: "Koordinator Dept 01 Organisasi", Department: strPtr("Departemen 01"), Level: "dpp", ImagePath: "https://images.unsplash.com/photo-1500648767791-00dcc994a43e?q=80&w=200", SortOrder: 17, Periode: "2025 - 2030", IsActive: true},
		// DPD (lanjutan)
		{Name: "Drs. H. Ahmad Fauzi", Role: "Ketua DPD Jawa Barat", Level: "dpd", Provinsi: strPtr("Jawa Barat"), ImagePath: "https://images.unsplash.com/photo-1560250097-0b93528c311a?q=80&w=200", SortOrder: 18, Periode: "2025 - 2030", IsActive: true},
		{Name: "Bambang Irawan, S.T", Role: "Ketua DPD Jawa Timur", Level: "dpd", Provinsi: strPtr("Jawa Timur"), ImagePath: "https://images.unsplash.com/photo-1519085360753-af0119f7cbe7?q=80&w=200", SortOrder: 19, Periode: "2025 - 2030", IsActive: true},
		{Name: "Siti Aminah, M.Si", Role: "Ketua DPD Jawa Tengah", Level: "dpd", Provinsi: strPtr("Jawa Tengah"), ImagePath: "https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?q=80&w=200", SortOrder: 20, Periode: "2025 - 2030", IsActive: true},
		{Name: "Hendra Gunawan", Role: "Ketua DPD DKI Jakarta", Level: "dpd", Provinsi: strPtr("DKI Jakarta"), ImagePath: "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?q=80&w=200", SortOrder: 21, Periode: "2025 - 2030", IsActive: true},
		{Name: "Tri Wahyudi", Role: "Ketua DPD Banten", Level: "dpd", Provinsi: strPtr("Banten"), ImagePath: "https://images.unsplash.com/photo-1500648767791-00dcc994a43e?q=80&w=200", SortOrder: 22, Periode: "2025 - 2030", IsActive: true},
		{Name: "Eko Prasetyo", Role: "Ketua DPD DI Yogyakarta", Level: "dpd", Provinsi: strPtr("DI Yogyakarta"), ImagePath: "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?q=80&w=200", SortOrder: 23, Periode: "2025 - 2030", IsActive: true},
		// DPC (lanjutan)
		{Name: "Syamsul Bahri", Role: "Ketua DPC Kota Bandung", Level: "dpc", Provinsi: strPtr("Jawa Barat"), Kabupaten: strPtr("Kota Bandung"), ImagePath: "https://images.unsplash.com/photo-1500648767791-00dcc994a43e?q=80&w=200", SortOrder: 24, Periode: "2025 - 2030", IsActive: true},
		{Name: "Herman Wijaya", Role: "Ketua DPC Kab. Bogor", Level: "dpc", Provinsi: strPtr("Jawa Barat"), Kabupaten: strPtr("Kabupaten Bogor"), ImagePath: "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?q=80&w=200", SortOrder: 25, Periode: "2025 - 2030", IsActive: true},
		{Name: "Ridwan Malik", Role: "Ketua DPC Kota Surabaya", Level: "dpc", Provinsi: strPtr("Jawa Timur"), Kabupaten: strPtr("Kota Surabaya"), ImagePath: "https://images.unsplash.com/photo-1599566150163-29194dcaad36?q=80&w=200", SortOrder: 26, Periode: "2025 - 2030", IsActive: true},
		{Name: "Anita Rahayu", Role: "Ketua DPC Kab. Malang", Level: "dpc", Provinsi: strPtr("Jawa Timur"), Kabupaten: strPtr("Kabupaten Malang"), ImagePath: "https://images.unsplash.com/photo-1544005313-94ddf0286df2?q=80&w=200", SortOrder: 27, Periode: "2025 - 2030", IsActive: true},
	}
	db.Create(&pengurusList)

	// 10. Insert Pesan Kontak Contoh
	log.Println("[10/10] Menyisipkan Pesan Kontak Sampel...")
	kontakList := []kontakModel.PesanKontak{
		{
			Nama:   "Ahmad Zaky",
			Email:  "zaky@example.com",
			Subjek: "Permohonan Kerjasama Literasi Digital",
			Pesan:  "Halo DPP GRADASI, kami dari komunitas IT daerah ingin mengajukan kolaborasi program literasi digital bagi UMKM lokal.",
			IsRead: false,
		},
		{
			Nama:   "Siti Rahma",
			Email:  "siti@example.com",
			Subjek: "Informasi Keanggotaan DPD",
			Pesan:  "Mohon informasi persyaratan dan prosedur pendaftaran pengurus DPD Jawa Tengah. Terima kasih.",
			IsRead: true,
		},
	}
	db.Create(&kontakList)

	log.Println("===========================================")
	log.Println("🚀 SEEDING DATABASE BERHASIL 100%!")
	log.Printf("Super Admin: %s | Pass: %s\n", adminEmail, adminPassword)
	log.Println("===========================================")
}
