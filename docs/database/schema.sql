-- ============================================
-- DATABASE: dpp_gradasi
-- Schema Version: 2.0 (Disesuaikan dengan API Contracts)
-- Terakhir diperbarui: 2026-07-30
-- ============================================

CREATE DATABASE IF NOT EXISTS dpp_gradasi CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE dpp_gradasi;

-- ============================================
-- 1. ROLES (Tabel Peran - Baru)
-- ============================================
-- Menggantikan ENUM role pada tabel users lama.
-- Kontrak API menggunakan role_id + role_name + role_display_name.
CREATE TABLE roles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,           -- slug: super_admin, admin_keuangan, dll
    display_name VARCHAR(100) NOT NULL,          -- tampilan: Super Administrator, Admin Keuangan
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Seed data roles default
INSERT INTO roles (name, display_name, description) VALUES
('super_admin', 'Super Administrator', 'Akses penuh ke seluruh sistem'),
('admin', 'Administrator', 'Mengelola konten dan pengguna'),
('editor', 'Editor', 'Mengelola konten berita dan kegiatan'),
('viewer', 'Viewer', 'Hanya bisa melihat data');

-- ============================================
-- 2. USERS (Admin & Pengelola)
-- ============================================
-- Disesuaikan dengan kontrak API users.jsonc
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    role_id INT NOT NULL,                                 -- FK ke tabel roles (menggantikan ENUM)
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150) NOT NULL UNIQUE,
    password VARCHAR(255) DEFAULT NULL,                   -- Nullable: user baru belum set password
    photo_path VARCHAR(500),                              -- Renamed dari avatar_url
    status ENUM('active', 'inactive') NOT NULL DEFAULT 'inactive',  -- Menggantikan is_active BOOLEAN
    email_verified_at TIMESTAMP NULL,                     -- BARU: waktu verifikasi email
    last_login_at TIMESTAMP NULL,                         -- Renamed dari last_login
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,                            -- BARU: soft delete support

    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE RESTRICT,
    INDEX idx_role_id (role_id),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_email (email)
);

-- ============================================
-- 3. REFRESH TOKENS (Baru — untuk auth flow)
-- ============================================
-- Menyimpan refresh token untuk JWT auth dengan rotation.
-- Sesuai dengan kontrak auth.jsonc (refresh, logout, change_password).
CREATE TABLE refresh_tokens (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    token_hash VARCHAR(255) NOT NULL UNIQUE,    -- Hash dari refresh token (jangan simpan plain text)
    device_info VARCHAR(255),                   -- User-Agent / device identifier
    ip_address VARCHAR(45),                     -- IPv4 atau IPv6
    expires_at TIMESTAMP NOT NULL,              -- Waktu kadaluarsa token
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_expires_at (expires_at)
);

-- ============================================
-- 4. PASSWORD RESET TOKENS (Baru)
-- ============================================
-- Menyimpan token untuk forgot password flow.
-- Sesuai dengan kontrak auth.jsonc (forgot_password, validate_reset_token, reset_password).
CREATE TABLE password_reset_tokens (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    token_hash VARCHAR(255) NOT NULL UNIQUE,    -- Hash dari reset token
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP NULL,                     -- NULL jika belum digunakan
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_token_hash (token_hash)
);

-- ============================================
-- 5. ACTIVATION TOKENS (Baru)
-- ============================================
-- Menyimpan token untuk aktivasi akun user baru.
-- Sesuai dengan kontrak users.jsonc (activate, resend_notification).
CREATE TABLE activation_tokens (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    channel ENUM('email', 'whatsapp', 'all') NOT NULL DEFAULT 'email',
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_token_hash (token_hash)
);

-- ============================================
-- 6. SLIDERS (Hero Banner di Beranda)
-- ============================================
CREATE TABLE sliders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    subtitle TEXT,
    tag VARCHAR(50),
    is_new BOOLEAN DEFAULT FALSE,
    event_date VARCHAR(100),
    location VARCHAR(200),
    image_url VARCHAR(500) NOT NULL,
    link_url VARCHAR(500),
    sort_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_by INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
    INDEX idx_sort_order (sort_order),
    INDEX idx_is_active (is_active)
);

-- ============================================
-- 7. BERITA (Informasi / News)
-- ============================================
CREATE TABLE berita (
    id INT AUTO_INCREMENT PRIMARY KEY,
    slug VARCHAR(250) NOT NULL UNIQUE,
    title VARCHAR(300) NOT NULL,
    category VARCHAR(100) DEFAULT 'Berita Organisasi',
    published_date DATE NOT NULL,
    author_id INT,
    image_url VARCHAR(500),
    excerpt TEXT,
    content LONGTEXT,
    is_featured BOOLEAN DEFAULT FALSE,
    is_published BOOLEAN DEFAULT TRUE,
    views INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,                  -- BARU: soft delete support

    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE SET NULL,
    INDEX idx_published_date (published_date DESC),
    INDEX idx_category (category),
    INDEX idx_is_published (is_published),
    INDEX idx_deleted_at (deleted_at),
    FULLTEXT INDEX ft_berita (title, content)
);

-- ============================================
-- 8. BERITA_TAGS
-- ============================================
CREATE TABLE berita_tags (
    id INT AUTO_INCREMENT PRIMARY KEY,
    berita_id INT NOT NULL,
    tag VARCHAR(100) NOT NULL,

    FOREIGN KEY (berita_id) REFERENCES berita(id) ON DELETE CASCADE,
    INDEX idx_berita_id (berita_id),
    INDEX idx_tag (tag)
);

-- ============================================
-- 9. KEGIATAN (Events / Activities)
-- ============================================
CREATE TABLE kegiatan (
    id INT AUTO_INCREMENT PRIMARY KEY,
    slug VARCHAR(250) NOT NULL UNIQUE,
    title VARCHAR(300) NOT NULL,
    category VARCHAR(100) DEFAULT 'Kegiatan',
    event_date DATE,
    location VARCHAR(200),
    organizer VARCHAR(200),
    author_id INT,                              -- BARU: Relasi ke penulis (user)
    image_url VARCHAR(500),
    excerpt TEXT,
    content LONGTEXT,
    is_published BOOLEAN DEFAULT TRUE,
    views INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,                  -- BARU: soft delete support

    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE SET NULL,
    INDEX idx_event_date (event_date DESC),
    INDEX idx_is_published (is_published),
    INDEX idx_deleted_at (deleted_at),
    FULLTEXT INDEX ft_kegiatan (title, content)
);

-- ============================================
-- 10. KEGIATAN_TAGS
-- ============================================
CREATE TABLE kegiatan_tags (
    id INT AUTO_INCREMENT PRIMARY KEY,
    kegiatan_id INT NOT NULL,
    tag VARCHAR(100) NOT NULL,

    FOREIGN KEY (kegiatan_id) REFERENCES kegiatan(id) ON DELETE CASCADE,
    INDEX idx_kegiatan_id (kegiatan_id),
    INDEX idx_tag (tag)
);

-- ============================================
-- 11. KEGIATAN_GALLERY
-- ============================================
CREATE TABLE kegiatan_gallery (
    id INT AUTO_INCREMENT PRIMARY KEY,
    kegiatan_id INT NOT NULL,
    image_url VARCHAR(500) NOT NULL,
    caption VARCHAR(200),
    sort_order INT DEFAULT 0,

    FOREIGN KEY (kegiatan_id) REFERENCES kegiatan(id) ON DELETE CASCADE,
    INDEX idx_kegiatan_id (kegiatan_id)
);

-- ============================================
-- 12. PENGURUS (Kepengurusan Organisasi)
-- ============================================
CREATE TABLE pengurus (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    role VARCHAR(200) NOT NULL,               -- Jabatan: Ketua Umum DPP, dll
    department VARCHAR(100),                   -- Departemen/bidang
    level ENUM('ketua', 'dpp', 'dpd', 'dpc') NOT NULL,
    provinsi VARCHAR(100),                     -- Wajib jika level = dpd/dpc
    kabupaten VARCHAR(100),                    -- Wajib jika level = dpc
    image_url VARCHAR(500),
    facebook_url VARCHAR(500),
    instagram_url VARCHAR(500),
    linkedin_url VARCHAR(500),
    whatsapp VARCHAR(20),
    sort_order INT DEFAULT 0,
    periode VARCHAR(50),                       -- Contoh: "2025 - 2030"
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,                 -- BARU: soft delete support

    INDEX idx_level (level),
    INDEX idx_provinsi (provinsi),
    INDEX idx_kabupaten (kabupaten),
    INDEX idx_is_active (is_active),
    INDEX idx_sort_order (sort_order),
    INDEX idx_deleted_at (deleted_at)
);

-- ============================================
-- 13. PESAN_KONTAK (Contact Form Messages)
-- ============================================
CREATE TABLE pesan_kontak (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nama VARCHAR(100) NOT NULL,
    email VARCHAR(150) NOT NULL,
    subjek VARCHAR(200) NOT NULL,
    pesan TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    response_note TEXT,                        -- Catatan admin setelah merespon
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_is_read (is_read),
    INDEX idx_created_at (created_at DESC)
);

-- ============================================
-- 14. SETTINGS (Key-Value Configuration)
-- ============================================
CREATE TABLE settings (
    id INT AUTO_INCREMENT PRIMARY KEY,
    `key` VARCHAR(100) NOT NULL UNIQUE,
    value TEXT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_key (`key`)
);

-- Seed data settings default
INSERT INTO settings (`key`, value) VALUES
('site_name', 'DPP GRADASI'),
('site_tagline', 'Generasi Digital Indonesia'),
('logo_url', '/uploads/logo.png'),
('contact_email', 'dpp@gradasi.org'),
('maps_embed_url', ''),
('facebook_url', 'https://www.facebook.com/gradasiofficial.id'),
('instagram_url', 'https://www.instagram.com/dppgradasi'),
('youtube_url', 'https://www.youtube.com/channel/UCwdjB4LkqcF4Kw5-PoyOb5A'),
('tentang_tujuan', 'Pada tahun 2018 di Yogyakarta, sebuah inisiatif besar lahir untuk membentuk organisasi penggerak digital nasional.'),
('tentang_tanggal_terbentuk', '4 Februari 2019'),
('tentang_no_sk', 'AHU – 0000151.AH.01.07.2019'),
('visi', 'Mewujudkan masyarakat Indonesia yang cerdas dan berdaulat di era digital.'),
('misi', '["Membangun ekosistem literasi digital", "Meningkatkan kecakapan digital masyarakat"]');

-- ============================================
-- SEED: Default Super Admin User
-- ============================================
-- Password: password123 (hashed dengan bcrypt)
(1, 'Administrator', 'admin@gradasi.org', '6285279880008', '$2a$12$placeholder_bcrypt_hash_here', 'active', NOW(), NOW());
