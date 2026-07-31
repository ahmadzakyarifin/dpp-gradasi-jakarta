-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    role_id INT NOT NULL,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150) NOT NULL UNIQUE,
    password VARCHAR(255) DEFAULT NULL,
    photo_path VARCHAR(500),
    status ENUM('active', 'inactive') NOT NULL DEFAULT 'inactive',
    email_verified_at TIMESTAMP NULL,
    last_login_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,

    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE RESTRICT,
    INDEX idx_role_id (role_id),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_email (email)
);

(1, 'Administrator', 'admin@gradasi.org', '6285279880008', '$2a$10$WYeiTAB4gHLq1jxprM47.OIlQLC1xxOEQgdov6/6K2pV8mBYPP.a6', 'active', NOW(), NOW());
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
