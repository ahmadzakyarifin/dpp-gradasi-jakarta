-- +goose Up
-- +goose StatementBegin
CREATE TABLE settings (
    id INT AUTO_INCREMENT PRIMARY KEY,
    site_name VARCHAR(100) NOT NULL,
    tagline VARCHAR(200),
    logo_path VARCHAR(500),
    contact_email VARCHAR(100) NOT NULL,
    contact_phone VARCHAR(20) NOT NULL,
    address VARCHAR(255),
    maps_embed_url VARCHAR(500),
    facebook_url VARCHAR(200),
    instagram_url VARCHAR(200),
    youtube_url VARCHAR(200),
    video_profile_path VARCHAR(500),
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
    greeting_image_path VARCHAR(500),
    updated_by INT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS settings;
-- +goose StatementEnd
