package main

import (
	"log"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/config"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/infrastructure"
	beritaModel "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/berita/model"
	kegiatanModel "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kegiatan/model"
	kontakModel "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kontak/model"
	pengurusModel "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/pengurus/model"
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
		"users", "berita", "kegiatan", "pengurus", "sliders", "pesan_kontak", "activity_logs",
	}
	db.Exec("SET FOREIGN_KEY_CHECKS = 0;")
	for _, t := range tables {
		db.Exec("TRUNCATE TABLE " + t)
	}
	db.Exec("SET FOREIGN_KEY_CHECKS = 1;")

	// Helper for string pointers
	strPtr := func(s string) *string { return &s }

	// 3. (Langkah ini dilewati karena tabel roles sudah dihapus dan diubah menjadi ENUM di tabel users)

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
		Role:      "super_admin",
		Name:      adminName,
		Email:     adminEmail,
		Password:  string(hash),
		Status:    "active",
		PhotoPath: strPtr("https://ui-avatars.com/api/?name=Super+Admin&background=0D8ABC&color=fff"),
	}
	db.Create(&adminUser)

	// 5. Insert Settings (dari index.html)
	log.Println("[5/10] Menyisipkan Settings website (Visi, Misi, Legalitas, Kontak)...")
	setting := settingsModel.Settings{
		SiteName:             "DPP GRADASI",
		Tagline:              "Generasi Digital Indonesia",
		LogoPath:             "https://gradasi.org/uploads/img/logo/1737187847.png",
		ContactEmail:         "dpp@gradasi.org",
		ContactPhone:         "+6285279880008",
		Address:              "Office Park OL3-12A The Bellagio Mall Jl. Kawasan Mega Kuningan Kav.E.4.3 Mega Kuningan, Kel. Kuningan Timur, Kec.Setiabudi, Jakarta Selatan 15810",
		MapsEmbedURL:         "https://www.google.com/maps/embed?pb=!1m14!1m8!1m3!1d4591.169749593668!2d106.824223!3d-6.227444!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x2e69f3e4ac2fa421%3A0x2de3d495cc84d79d!2sThe%20Bellagio%20Boutique%20Mall!5e1!3m2!1sid!2sus!4v1785759846547!5m2!1sid!2sus",
		FacebookURL:          "https://www.facebook.com/gradasiofficial.id",
		InstagramURL:         "https://www.instagram.com/dppgradasi",
		YoutubeURL:           "https://www.youtube.com/channel/UCwdjB4LkqcF4Kw5-PoyOb5A",
		VideoProfilePath:     "https://gradasi.org/assets/video/gradasi.mp4",
		History:              "Pada tahun 2018 di Yogyakarta, atas arahan Bang Oni kepada Kang Ditu, Cak Levy, juga diajak Mas Memet dan Bang Fuad membentuk sebuah organisasi penggerak digital di Indonesia sebagai wadah perkumpulan pegiat dan penggiat digital di tingkat nasional yang terbuka untuk umum, guna melakukan pendampingan dan pemberdayaan untuk semua kalangan masyarakat.\n\nMenindaklanjuti arahan tersebut, kemudian dibentuk tim inisiator Gradasi yang terdiri dari Muhammad Sidik Kaimuddin Tomsio, Fuad Rizaldy, Memet Toto Raharjo dan beberapa kawan lain seperti Heri Safrizal, Sari Erni Agustin, serta didukung oleh Nur Alfi Khabibah, Junaidi, dan Yunita di Yogyakarta pada bulan Februari 2019 untuk mendirikan Perkumpulan Generasi Digital Indonesia (Gradasi).",
		AboutTutorial:        "Pengesahan Badan Hukum Kemenkumham RI.",
		AboutFormationDate:   "4 Februari 2019",
		AboutNoSK:            "AHU-0000151.AH.01.07.2019",
		AboutVision:          "menjadikan GRADASI sebagai organisasi penggerak literasi digital nasional, yang hadir untuk mengedukasi masyarakat dalam percepatan transformasi literasi digital.",
		AboutMission:         `["Menghimpun praktisi, akademisi dan penggiat digital untuk berbagi pengetahuan dan keterampilan khususnya dalam bidang literasi berbasis digital","Mengedukasi masyarakat melalui kegiatan literasi berbasis digital","Menfasilitasi kegiatan yang menunjang pengembangan SDM Indonesia khususnya dalam bidang literasi berbasis digital melalui workshop, training, seminar, bimtek dan lainya","Bersinergi dengan stakeholder Pusat maupun Daerah dan insetusi lainnya untuk menyusun konsep","Peningkatan kualitas pendidikan berbasis digital"]`,
		GreetingTitle:        "SELAMAT TAHUN BARU 2026",
		GreetingSubtitle:     "Refleksi & Optimisme",
		GreetingDate:         "31 Desember 2025",
		GreetingContent:      "Dewan Pengurus Pusat Generasi Digital Indonesia (DPP GRADASI) secara resmi menyampaikan refleksi akhir tahun sekaligus pernyataan optimisme menyambut tahun baru 2026. GRADASI menegaskan komitmennya untuk tetap menjadi garda terdepan dalam memperjuangkan kedaulatan digital bangsa.\n\nTahun 2025 yang segera berakhir dinilai sebagai sebuah perjalanan yang penuh makna, sebuah \"Kisah Panjang\" yang diisi dengan berbagai dinamika transformasi teknologi. Perjalanan ini merupakan bentuk dedikasi seluruh elemen organisasi dalam menghadapi tantangan era digital, baik di tingkat pusat maupun daerah.\n\n\"Tiada hadiah termahal di akhir tahun ini selain doa dan dukungan yang senantiasa mengiringi setiap langkah perjuangan kita,\" ujar perwakilan DPP GRADASI. Organisasi ini memosisikan tantangan teknologi di tahun 2025 bukan sebagai beban, melainkan sebagai proses pendewasaan yang memungkinkan seluruh anggota untuk tetap tegak berdiri demi kemandirian digital Indonesia.\n\nMemasuki tahun 2026, GRADASI menetapkan tiga pilar utama perjuangan:\n1. Kedaulatan Digital: Terus memperjuangkan agar Indonesia mampu mandiri dan berdaulat atas ekosistem digitalnya sendiri.\n2. Kemandirian Generasi: Memastikan setiap pengurus dan anggota masyarakat Indonesia memiliki kecakapan digital (digital skills) untuk menjadi aktor utama di panggung global.\n3. Inovasi Berkelanjutan: Mengembangkan program-materi literasi digital yang relevan dan berdampak luas bagi kemajuan bangsa.\n4. Program-program yang sudah dijalankan baik di DPP maupun DPD Gradasi seluruh Indonesia.\n\nSelamat Tahun Baru 2026. Mari kita tutup babak \"Kisah Panjang\" tahun ini dengan rasa syukur dan memulai lembaran baru dengan semangat kolaborasi yang lebih kuat.\n\nSalam Kolaborasi\nGenerasi Membangun Negeri\n\nDPP Generasi Digital Indonesia\n\nKetua Umum\nUpi Asmaradhana\n\nSekretaris Jenderal\nJunaidi\n\nDPP Gradasi",
		GreetingImagePath:    "https://gradasi.org/uploads/img/event-terkini/1767154211.jpg",
		LoginHeroTitle:       "Kembangkan Potensi Digital Anda",
		LoginHeroDescription: "Bergabung bersama ribuan pemuda Indonesia lainnya dalam membangun ekosistem digital yang kuat, mandiri, dan inovatif untuk kemajuan bangsa.",
	}
	db.Create(&setting)

	// 6. Insert Sliders Carousel (dari index.html)
	log.Println("[6/10] Menyisipkan Banner Sliders (Header Carousel)...")
	sliders := []slidersModel.Slider{
		{
			Title:     "Foto Ketua Umum DPP GRADASI",
			Subtitle:  strPtr("Foto Ketua Umum DPP GRADASI"),
			Tag:       strPtr("FOTO RESMI"),
			ImagePath: "https://gradasi.org/uploads/img/slider/1749385864.jpg",
			SortOrder: 1,
			IsPublished: true,
			IsNew:     false,
			EventDate: nil,
			Location:  nil,
		},
		{
			Title:     "Musyawarah Nasional Ke-II GRADASI",
			Subtitle:  strPtr("Kolaborasi Membangun Negeri Menuju Indonesia Emas 2045"),
			Tag:       strPtr("HEADLINE EVENT"),
			ImagePath: "https://gradasi.org/uploads/img/slider/1746600520.png",
			SortOrder: 2,
			IsPublished: true,
			IsNew:     true,
			EventDate: strPtr("10 - 12 Agustus 2026"),
			Location:  strPtr("Jakarta Convention Center"),
		},
		{
			Title:     "UMKM KOTA PALANGKA RAYA",
			Subtitle:  strPtr("Pelatihan Digital UMKM"),
			Tag:       strPtr("INFO KEGIATAN"),
			ImagePath: "https://gradasi.org/uploads/img/slider/1746600828.jpg",
			SortOrder: 3,
			IsPublished: true,
			IsNew:     false,
			EventDate: strPtr("20 September 2026"),
			Location:  strPtr("Palangka Raya"),
		},
	}
	db.Create(&sliders)

	// 7. Insert Berita Default (dari berita.html & index.html)
	log.Println("[7/10] Menyisipkan Data Berita Default...")
	beritaList := []beritaModel.Berita{
		{
			Title:         "GRADASI Desak Dewan Perwakilan Rakyat Wujudkan Dekolonisasi Digital: Sampaikan Usulan Kebijakan Strategis di NCSC 2025",
			Slug:          "gradasi-desak-dpr-wujudkan-dekolonisasi-digital",
			Category:      "Berita Utama",
			PublishedDate: "30 Oktober 2025",
			AuthorID:      &adminUser.ID,
			ImagePath:     strPtr("https://gradasi.org/uploads/img/berita/17618789900.jpg"),
			Excerpt:       strPtr("Generasi Digital Indonesia (GRADASI) secara tegas menyuarakan perlunya penguatan kedaulatan digital dan mendesak reformasi kebijakan untuk mengatasi dominasi Big Tech global."),
			Content:       strPtr("<p><strong>Jakarta, 30 Oktober 2025</strong> – Generasi Digital Indonesia (GRADASI) secara tegas menyuarakan perlunya penguatan kedaulatan digital dan mendesak reformasi kebijakan untuk mengatasi dominasi Big Tech global. Posisi strategis ini disampaikan dalam ajang National Cybersecurity Connect 2025 (NCSC 2025) yang diselenggarakan di Bidakara Hotel, Jakarta.</p><p>GRADASI mengukuhkan perannya dalam forum tingkat nasional ini, di mana Ketua Umum Upi Asmaradhana diundang ke sesi Congress untuk memberikan pandangan kebijakan. Kontribusi pemikiran GRADASI dalam sesi Congress secara formal akan diresmikan menjadi Usulan Kebijakan Strategis (Recommendation Summary) bagi penguatan ketahanan siber nasional, yang selanjutnya akan disampaikan kepada pemangku kepentingan utama, termasuk Dewan Perwakilan Rakyat (DPR) dan Kementerian/Lembaga (K/L) terkait.</p><p><em>Damar Juniarto Sampaikan Rekomendasi Kunci</em></p><p>GRADASI melalui Ketua Dewan Pakar, Damar Juniarto, S.Sos, M.Kom, sebagai pembicara kunci di sesi Cyberstage. Ia menyampaikan presentasi yang berjudul \"Kedaulatan Digital Dalam Perspektif Kesinambungan Ekosistem: Belajar dari Pengalaman Negara-Bangsa Melawan Raksasa Digital\".</p><p>\"Kedaulatan digital bukan sekadar isu keamanan siber, melainkan tentang menuntut kesetaraan, keadilan, dan akuntabilitas dari platform digital. Kita harus berani mengambil langkah untuk menjaga kepentingan nasional,\" tegas Damar Juniarto.</p><p>Untuk mewujudkan hal tersebut, GRADASI mengajukan lima rekomendasi kebijakan krusial (tertuang dalam materi presentasi), meliputi:</p><ol><li>Dekolonisasi Digital (Rebut Kedaulatan Digital) dan Penguatan Infrastruktur Lokal.</li><li>Mendorong Ko-regulasi dengan Platform Digital.</li><li>Perkuat Pengawasan Monopoli dan Persaingan Usaha.</li><li>Diplomasi Digital, dan</li><li>Penguatan infrastruktur lokal.</li></ol><p>Partisipasi GRADASI di NCSC 2025 juga diperkuat dengan booth pameran yang menjadi sentra sinergi DPD GRADASI DKI Jakarta, LSK Siber Indonesia, dan TUK binaan.</p><p><strong>Tentang Generasi Digital Indonesia (GRADASI)</strong><br/>Generasi Digital Indonesia (GRADASI) adalah organisasi yang didirikan untuk mendorong literasi, inovasi, dan kedaulatan digital di Indonesia. GRADASI berkomitmen untuk menjadi mitra strategis pemerintah dan masyarakat dalam membangun ekosistem digital yang sehat, aman, dan berkeadilan, melalui semangat #Bekerja #Berkarya #Berdaya.</p><p>Narahubung:<br/>Adi Rasmiadi<br/>Koordinator Departemen Data Informasi, Dokumentasi dan Publikasi<br/>082352148375</p><p>#NCSC2025 #KedaulatanDigital #DekolonisasiDigital #GRADASI #KebijakanSiber</p>"),
			IsPublished:   true,
			Views:         1250,
		},
		{
			Title:         "GRADASI Berpartisipasi Penuh pada Kegiatan National Cybersecurity Connect 2025: Dorong Kedaulatan Digital di Forum Nasional",
			Slug:          "gradasi-berpartisipasi-penuh-ncsc-2025",
			Category:      "Berita Utama",
			PublishedDate: "29 Oktober 2025",
			AuthorID:      &adminUser.ID,
			ImagePath:     strPtr("https://gradasi.org/uploads/img/berita/17617487980.jpg"),
			Excerpt:       strPtr("Generasi Digital Indonesia (GRADASI) ikut berpartisipasi penuh dalam ajang terbesar cybersecurity di Indonesia, National Cybersecurity Connect 2025 (NCSC 2025), yang diselenggarakan pada 29-30 Oktober 2025 di Bidakara Hotel, Jakarta."),
			Content:       strPtr("<p><strong>Jakarta, 29 Oktober 2025</strong> – Generasi Digital Indonesia (GRADASI) ikut berpartisipasi penuh dalam ajang terbesar cybersecurity di Indonesia, National Cybersecurity Connect 2025 (NCSC 2025), yang diselenggarakan pada 29-30 Oktober 2025 di Bidakara Hotel, Jakarta.</p><p>NCSC 2025 merupakan forum tingkat nasional yang mempertemukan pemangku kepentingan dari sektor pemerintahan, industri, dan masyarakat. Kontribusi pemikiran dan pandangan GRADASI dalam sesi Congress akan dicatat secara strategis (Recommendation Summary) dan diresmikan menjadi usulan kebijakan untuk penguatan ketahanan siber nasional, yang kemudian akan diserahkan kepada para pemangku kepentingan utama, termasuk Dewan Perwakilan Rakyat (DPR) dan Kementerian/Lembaga (K/L).</p><p>GRADASI dipercaya menjadi salah satu pembicara di sesi Cyberstage hari kedua (30 Oktober), diwakili oleh Ketua Dewan Pakar, Damar Juniarto, S.Sos, M.Kom. Topik yang dibawakan sangat strategis:</p><p>\"Kedaulatan Digital Dalam Perspektif Kesinambungan Ekosistem: Belajar dari Pengalaman Negara-Bangsa Melawan Raksasa Digital\"</p><p>Komitmen GRADASI juga ditunjukkan dengan mengikuti pameran yang menjadi sentra sinergi organisasi, didukung oleh DPD GRADASI DKI Jakarta, LSK Siber Indonesia, serta dua TUK (Tempat Uji Kompetensi) binaan: LKP Kariermu Sekolahmu Jakarta dan LKP Digitalindo Academy Depok.</p><p>Ketua Umum DPP GRADASI Upi Asmaradhana mengatakan pihaknya siap terus berkontribusi dalam membentuk ekosistem keamanan siber yang kuat, adaptif, dan inklusif di Indonesia.</p>"),
			IsPublished:   true,
			Views:         890,
		},
		{
			Title:         "GRADASI Gelar Webinar 'AI Essentials': Bekali Guru dan EduPreneur Kemampuan Dasar AI untuk Produktivitas Digital",
			Slug:          "gradasi-gelar-webinar-ai-essentials",
			Category:      "Edukasi",
			PublishedDate: "2 November 2025",
			AuthorID:      &adminUser.ID,
			ImagePath:     strPtr("https://gradasi.org/uploads/img/berita/17620765070.jpg"),
			Excerpt:       strPtr("Inisiatif GRADASI Mendorong Peningkatan Kompetensi SDM Pendidikan dalam Memanfaatkan Kecerdasan Buatan (AI) di Era 4.0"),
			Content:       strPtr("<p><strong>Jakarta, 2 November 2025</strong> – Generasi Digital Indonesia (GRADASI), sebagai organisasi yang berfokus pada literasi and inovasi digital, mengumumkan penyelenggaraan Webinar AI Essentials: Langkah Awal Memahami dan Memanfaatkan AI di Era Digital. Webinar ini secara khusus ditujukan untuk segmen Guru Kreatif dan EduPreneur guna memastikan pemanfaatan teknologi AI yang optimal dalam sektor pendidikan.</p><p>Acara yang akan dilaksanakan pada Rabu, 5 November 2025, pukul 19.00 – 21.00 WIB melalui platform Zoom Meeting ini, merupakan respons GRADASI terhadap kebutuhan mendesak akan pemahaman dasar Kecerdasan Buatan (AI) yang kini menjadi keterampilan esensial.</p><p>\"AI bukan lagi teknologi masa depan, melainkan perangkat fundamental yang harus dikuasai oleh para pendidik dan pegiat edukasi untuk meningkatkan efektivitas pengajaran dan mengembangkan model bisnis pendidikan yang inovatif,\" ujar Upi Asmaradhana, Ketua Umum GRADASI.</p><p><em>Narasumber Kompeten dan Manfaat Eksklusif</em></p><p>Webinar ini akan diisi oleh tiga narasumber yang kompeten di bidangnya, yaitu:</p><ol><li>Ujang Hernawan, M.M., M.Pd.</li><li>Nuropik, S.M.</li><li>Sri Purwiasih, S.Pd., M.Pd.</li></ol><p>Peserta (terutama Guru Kreatif dan EduPreneur yang mendaftar secara Gratis) akan mendapatkan berbagai manfaat, termasuk ilmu yang aplikatif, Modul Pembelajaran, e-sertifikat, serta bonus berupa Buku 100 AI Rekomendasi dan Buku ChatGPT Mastery.</p><p>GRADASI mengajak seluruh pegiat pendidikan untuk mengambil kesempatan ini sebagai langkah awal dalam mengintegrasikan AI ke dalam ekosistem belajar-mengajar.</p><p>Pendaftaran dapat diakses melalui: bit.ly/webinargradasisesi1.</p><p><strong>Tentang Generasi Digital Indonesia (GRADASI)</strong><br/>Generasi Digital Indonesia (GRADASI) adalah organisasi yang didirikan untuk mendorong literasi, inovasi, dan kedaulatan digital di Indonesia. GRADASI berkomitmen untuk menjadi mitra strategis pemerintah dan masyarakat dalam membangun ekosistem digital yang sehat, aman, dan berkeadilan, melalui semangat #Bekerja #Berkarya #Berdaya.</p><p>Narahubung Media:<br/>Adi Rasmiadi<br/>Koordinator Departemen Data Informasi Publikasi<br/>082352148375</p><p>#GRADASI #WebinarAI #Pendidikan4_0 #GuruDigital #AIIndonesia #DigitalLearning</p>"),
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
			Title:       "Penyaluran Donasi Bencana Alam DPD Gradasi Sumatera Barat",
			Slug:        "penyaluran-donasi-bencana-dpd-gradasi-sumbar",
			Category:    "Kegiatan",
			EventDate:   "31 Desember 2025",
			Location:    "Sumatera Barat",
			Organizer:   "DPD Sumatera Barat",
			AuthorID:    &adminUser.ID,
			ImagePath:   strPtr("https://gradasi.org/uploads/img/event/1767154719.jpg"),
			Excerpt:     strPtr("Dewan Pimpinan Pusat (DPP) GRADASI menyalurkan bantuan kemanusiaan kepada korban bencana banjir bandang dan tanah longsor melalui DPD Gradasi Sumatera Barat."),
			Content:     strPtr("<p>Dewan Pimpinan Pusat (DPP) GRADASI menyalurkan bantuan kemanusiaan kepada korban bencana banjir bandang dan tanah longsor melalui DPD Gradasi Sumatera Barat. Bantuan ini merupakan wujud kepedulian dan solidaritas keluarga besar GRADASI terhadap masyarakat Sumbar yang terdampak musibah.</p><p>Melalui penyaluran bantuan ini, DPP GRADASI berharap dapat meringankan beban para korban serta memberikan semangat agar mereka segera bangkit dan pulih. DPP GRADASI berkomitmen untuk terus hadir dan berperan aktif dalam aksi-aksi kemanusiaan di berbagai daerah.</p><p>Salam Kolaborasi<br/>DPP Gradasi<br/>Tim Bantuan Kemanusian DPP Gradasi</p>"),
			IsPublished: true,
			Views:       1500,
		},
		{
			Title:       "Penyaluran Donasi Bencana Alam DPD Gradasi Sumatera Utara",
			Slug:        "penyaluran-donasi-bencana-dpd-gradasi-sumut",
			Category:    "Kegiatan",
			EventDate:   "31 Desember 2025",
			Location:    "Sumatera Utara",
			Organizer:   "DPD Sumatera Utara",
			AuthorID:    &adminUser.ID,
			ImagePath:   strPtr("https://gradasi.org/uploads/img/event/1767154397.jpg"),
			Excerpt:     strPtr("Dewan Pimpinan Pusat (DPP) GRADASI menyalurkan bantuan kemanusiaan kepada korban bencana banjir bandang dan tanah longsor melalui DPD Gradasi Sumut."),
			Content:     strPtr("<p>Dewan Pimpinan Pusat (DPP) GRADASI menyalurkan bantuan kemanusiaan kepada korban bencana banjir bandang dan tanah longsor melalui DPD Gradasi Sumut. Bantuan ini merupakan wujud kepedulian dan solidaritas keluarga besar GRADASI terhadap masyarakat Sumut yang terdampak musibah.</p><p>Melalui penyaluran bantuan ini, DPP GRADASI berharap dapat meringankan beban para korban serta memberikan semangat agar mereka segera bangkit dan pulih. DPP GRADASI berkomitmen untuk terus hadir dan berperan aktif dalam aksi-aksi kemanusiaan di berbagai daerah.</p><p>Salam Kolaborasi<br/>DPP Gradasi<br/>Tim Bantuan Kemanusian DPP Gradasi</p>"),
			IsPublished: true,
			Views:       980,
		},
		{
			Title:       "DPP GRADASI Gelar Upacara Virtual HUT Ke-80 RI, Kuatkan Komitmen Perjuangan Kedaulatan Digital Indonesia",
			Slug:        "dpp-gradasi-upacara-virtual-hut-ke-80-ri",
			Category:    "Kegiatan",
			EventDate:   "16 Agustus 2025",
			Location:    "Jakarta",
			Organizer:   "DPP GRADASI",
			AuthorID:    &adminUser.ID,
			ImagePath:   strPtr("https://gradasi.org/uploads/img/event/1755349199.png"),
			Excerpt:     strPtr("Dalam rangka memperingati Hari Ulang Tahun ke-80 Kemerdekaan Republik Indonesia, Dewan Pengurus Pusat Generasi Digital Indonesia (DPP GRADASI) secara resmi mengundang seluruh pengurus dan anggota untuk mengikuti Upacara Detik-Detik Proklamasi secara virtual."),
			Content:     strPtr("<p><strong>Jakarta, 16 Agustus 2025</strong> – Dalam rangka memperingati Hari Ulang Tahun ke-80 Kemerdekaan Republik Indonesia, Dewan Pengurus Pusat Generasi Digital Indonesia (DPP GRADASI) secara resmi mengundang seluruh pengurus dan anggota untuk mengikuti Upacara Detik-Detik Proklamasi secara virtual. Acara ini tidak hanya menjadi momentum perayaan kemerdekaan, tetapi juga penegasan komitmen GRADASI dalam memperjuangkan kedaulatan digital bangsa.</p><p>Upacara akan diselenggarakan pada Minggu, 17 Agustus 2025, mulai pukul 09.00 WIB hingga selesai. Sesuai dengan agenda rapat, upacara akan dipimpin oleh Inspektur Upacara yang juga merupakan Ketua Umum GRADASI, Upi Asmaradhana. Dengan mengusung semangat \"Mari Rebut Kedaulatan Digital Indonesia,\" upacara ini bertujuan untuk menumbuhkan rasa nasionalisme di kalangan penggiat digital dan menegaskan visi organisasi.</p><p>Acara akan mencakup serangkaian kegiatan, seperti pengibaran bendera Merah Putih, pembacaan teks Proklamasi Kemerdekaan oleh Tjipto Priyoto Prastowo, dan pembacaan Teks Pancasila oleh Pandu. Seluruh pengurus dan anggota diharapkan dapat hadir tepat waktu dan mengikuti upacara dengan khidmat.</p><p>Untuk menjaga keamanan dan kelancaran acara, detail teknis mengenai tautan Zoom tidak dipublikasikan secara umum. Informasi lengkap hanya akan disampaikan melalui kanal komunikasi internal organisasi.</p><p>Dewan Pengurus Pusat Generasi Digital Indonesia (DPP Gradasi)<br/>Ketua Umum - Upi Asmaradhana<br/>Sekretaris Jenderal - Junaidi</p>"),
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
