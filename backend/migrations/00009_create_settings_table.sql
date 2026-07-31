-- +goose Up
-- +goose StatementBegin
CREATE TABLE settings (
    id INT AUTO_INCREMENT PRIMARY KEY,
    site_name VARCHAR(100) NOT NULL,
    tagline VARCHAR(200),
    logo_url VARCHAR(200),
    contact_email VARCHAR(100) NOT NULL,
    contact_phone VARCHAR(20) NOT NULL,
    address VARCHAR(255),
    maps_embed_url VARCHAR(500),
    facebook_url VARCHAR(200),
    instagram_url VARCHAR(200),
    youtube_url VARCHAR(200),
    video_profile_url VARCHAR(500),
    history TEXT,
    about_tutorial TEXT,
    about_formation_date VARCHAR(50),
    about_no_sk VARCHAR(100),
    about_vision TEXT,
    about_mission TEXT,
    greeting_title VARCHAR(255),
    greeting_subtitle VARCHAR(255),
    greeting_date VARCHAR(100),
    greeting_content TEXT,
    greeting_image_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

INSERT INTO settings (id, site_name, tagline, logo_url, contact_email, contact_phone, about_formation_date, about_no_sk, about_vision, about_mission, greeting_title, greeting_subtitle, greeting_date, greeting_content, greeting_image_url) 
VALUES (1, 'DPP GRADASI', 'Generasi Digital Indonesia', '/uploads/logo.png', 'dpp@gradasi.org', '+6285279880008', '4 Februari 2019', 'AHU – 0000151.AH.01.07.2019', 'Mewujudkan masyarakat Indonesia yang cerdas dan berdaulat di era digital.', '["Membangun ekosistem literasi digital", "Meningkatkan kecakapan digital masyarakat"]', 'Tahun Baru 2026', 'Resolusi & Harapan', '11 February 2026', 'Memasuki tahun 2026, GRADASI menetapkan pilar utama perjuangan...', 'https://gradasi.org/uploads/img/event-terkini/1767154211.jpg');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS settings;
-- +goose StatementEnd
