package main

import (
	"log"
	"os"

	"github.com/ahmadzakyarifin/schoolpay/backend/config"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/infrastructure"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	if _, err := os.Stat("../.env"); err == nil {
		os.Chdir("..")
	}

	log.Println("Memuat konfigurasi...")
	cfg := config.MustLoad()

	// Override host dan port karena seeder dijalankan dari host
	cfg.Database.Host = "127.0.0.1"
	cfg.Database.Port = "3308"

	log.Println("Mengkoneksikan ke database...")
	db, err := infrastructure.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Gagal connect DB: %v", err)
	}

	log.Println("Membuat password hash untuk admin dan parent...")

	// Hash untuk admin
	adminHash, err := bcrypt.GenerateFromPassword([]byte(cfg.Dev.Admin.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Gagal hash password admin: %v", err)
	}

	// Hash untuk parent
	parentHash, err := bcrypt.GenerateFromPassword([]byte(cfg.Dev.Parent.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Gagal hash password parent: %v", err)
	}

	log.Println("Menjalankan seeder Roles...")
	err = db.Exec(`
		INSERT INTO roles (id, name, display_name, is_system, is_active, created_at, updated_at) 
		VALUES 
			(1, 'super_admin', 'Super Administrator', 1, 1, NOW(), NOW()),
			(2, 'parent', 'Orang Tua / Wali', 1, 1, NOW(), NOW())
		ON DUPLICATE KEY UPDATE is_active = 1;
	`).Error
	if err != nil {
		log.Fatalf("Gagal insert roles: %v", err)
	}

	log.Println("Menjalankan seeder Permissions (34 permission)...")
	err = db.Exec(`
		INSERT INTO permissions (name, display_name, module, description, created_at, updated_at)
		VALUES
			-- Data Referensi Akademik
			('major.view',              'Lihat Jurusan',            'major',          'Melihat daftar dan detail jurusan', NOW(), NOW()),
			('major.manage',            'Kelola Jurusan',           'major',          'Membuat, mengubah, menghapus, dan mengelola jurusan', NOW(), NOW()),
			('academicyear.view',       'Lihat Tahun Ajaran',       'academicyear',   'Melihat daftar dan detail tahun ajaran', NOW(), NOW()),
			('academicyear.manage',     'Kelola Tahun Ajaran',      'academicyear',   'Membuat, mengubah, menghapus tahun ajaran', NOW(), NOW()),
			('semester.view',           'Lihat Semester',           'semester',       'Melihat daftar dan detail semester', NOW(), NOW()),
			('semester.manage',         'Kelola Semester',          'semester',       'Membuat, mengubah, menghapus semester', NOW(), NOW()),
			('cohort.view',             'Lihat Angkatan',           'cohort',         'Melihat daftar dan detail angkatan', NOW(), NOW()),
			('cohort.manage',           'Kelola Angkatan',          'cohort',         'Membuat, mengubah, menghapus angkatan', NOW(), NOW()),
			('classtemplate.view',      'Lihat Template Kelas',     'classtemplate',  'Melihat daftar dan detail template kelas', NOW(), NOW()),
			('classtemplate.manage',    'Kelola Template Kelas',    'classtemplate',  'Membuat, mengubah, menghapus template kelas', NOW(), NOW()),
			('activeclass.view',        'Lihat Kelas Aktif',        'activeclass',    'Melihat daftar dan detail kelas aktif', NOW(), NOW()),
			('activeclass.manage',      'Kelola Kelas Aktif',       'activeclass',    'Mengelola kelas aktif per tahun ajaran', NOW(), NOW()),
			('guardian.view',           'Lihat Wali',               'guardian',       'Melihat daftar dan detail wali', NOW(), NOW()),
			('guardian.manage',         'Kelola Wali',              'guardian',       'Membuat, mengubah, menghapus wali', NOW(), NOW()),

			-- Keanggotaan Kelas
			('classmembership.view',    'Lihat Anggota Kelas',      'classmembership', 'Melihat daftar anggota kelas', NOW(), NOW()),
			('classmembership.manage',  'Kelola Anggota Kelas',     'classmembership', 'Mendaftarkan, memindahkan, mengubah status anggota kelas', NOW(), NOW()),

			-- Siswa
			('student.view',            'Lihat Siswa',              'student',        'Melihat daftar dan detail data siswa', NOW(), NOW()),
			('student.create',          'Tambah Siswa',             'student',        'Menambah data siswa baru', NOW(), NOW()),
			('student.update',          'Ubah Siswa',               'student',        'Mengubah data, status, dan memulihkan siswa', NOW(), NOW()),
			('student.delete',          'Hapus Siswa',              'student',        'Menghapus data siswa', NOW(), NOW()),
			('student.graduate',        'Wisuda Siswa',             'student',        'Wisuda massal siswa', NOW(), NOW()),
			('student.promote',         'Naik Kelas Siswa',         'student',        'Kenaikan kelas massal siswa', NOW(), NOW()),
			('student.export',          'Export Data Siswa',        'student',        'Export data siswa ke Excel', NOW(), NOW()),

			-- Pengguna & Role
			('user.view',               'Lihat Pengguna',           'user',           'Melihat daftar dan detail pengguna', NOW(), NOW()),
			('user.create',             'Tambah Pengguna',          'user',           'Menambah pengguna baru', NOW(), NOW()),
			('user.update',             'Ubah Pengguna',            'user',           'Mengubah data, status, dan notifikasi pengguna', NOW(), NOW()),
			('user.delete',             'Hapus Pengguna',           'user',           'Menghapus pengguna', NOW(), NOW()),
			('user.export',             'Export Data Pengguna',     'user',           'Export data pengguna ke Excel', NOW(), NOW()),
			('role.view',               'Lihat Role',               'role',           'Melihat daftar dan detail role serta permissions', NOW(), NOW()),
			('role.create',             'Tambah Role',              'role',           'Membuat role baru', NOW(), NOW()),
			('role.update',             'Ubah Role',                'role',           'Mengubah role dan statusnya', NOW(), NOW()),
			('role.delete',             'Hapus Role',               'role',           'Menghapus role', NOW(), NOW()),

			-- Monitoring
			('notification.view',       'Lihat Notifikasi',         'notification',   'Melihat efikasi pengiriman notifikasi', NOW(), NOW()),
			('activitylog.view',        'Lihat Activity Log',       'activitylog',    'Melihat log aktivitas sistem', NOW(), NOW())
		ON DUPLICATE KEY UPDATE display_name = VALUES(display_name), description = VALUES(description);
	`).Error
	if err != nil {
		log.Fatalf("Gagal insert permissions: %v", err)
	}

	log.Println("Menjalankan seeder Users (Admin & Parent)...")

	// Seeder Admin
	err = db.Exec(`
		INSERT INTO users (id, role_id, name, email, phone, password_hash, status, created_at, updated_at) 
		VALUES (1, 1, ?, ?, ?, ?, 'active', NOW(), NOW())
		ON DUPLICATE KEY UPDATE 
			name = VALUES(name), email = VALUES(email), phone = VALUES(phone), 
			password_hash = VALUES(password_hash), status = 'active';
	`, cfg.Dev.Admin.Name, cfg.Dev.Admin.Email, cfg.Dev.Admin.Phone, string(adminHash)).Error

	if err != nil {
		log.Fatalf("Gagal insert admin user: %v", err)
	}

	// Seeder Parent
	err = db.Exec(`
		INSERT INTO users (id, role_id, name, email, phone, password_hash, status, created_at, updated_at) 
		VALUES (2, 2, ?, ?, ?, ?, 'active', NOW(), NOW())
		ON DUPLICATE KEY UPDATE 
			name = VALUES(name), email = VALUES(email), phone = VALUES(phone), 
			password_hash = VALUES(password_hash), status = 'active';
	`, cfg.Dev.Parent.Name, cfg.Dev.Parent.Email, cfg.Dev.Parent.Phone, string(parentHash)).Error

	if err != nil {
		log.Fatalf("Gagal insert parent user: %v", err)
	}

	log.Println("==================================================")
	log.Println("Seeder berhasil dijalankan!")
	log.Println("[Admin]")
	log.Println("Email Login :", cfg.Dev.Admin.Email)
	log.Println("Password    :", cfg.Dev.Admin.Password)
	log.Println("---")
	log.Println("[Parent]")
	log.Println("Email Login :", cfg.Dev.Parent.Email)
	log.Println("Password    :", cfg.Dev.Parent.Password)
	log.Println("==================================================")
}
