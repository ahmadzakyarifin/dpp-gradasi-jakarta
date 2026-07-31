-- +goose Up
-- +goose StatementBegin
CREATE TABLE pengurus (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    role VARCHAR(200) NOT NULL,
    department VARCHAR(100),
    level ENUM('ketua', 'dpp', 'dpd', 'dpc') NOT NULL,
    provinsi VARCHAR(100),
    kabupaten VARCHAR(100),
    image_url VARCHAR(500),
    facebook_url VARCHAR(500),
    instagram_url VARCHAR(500),
    linkedin_url VARCHAR(500),
    whatsapp VARCHAR(20),
    sort_order INT DEFAULT 0,
    periode VARCHAR(50),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,

    INDEX idx_level (level),
    INDEX idx_provinsi (provinsi),
    INDEX idx_kabupaten (kabupaten),
    INDEX idx_is_active (is_active),
    INDEX idx_sort_order (sort_order),
    INDEX idx_deleted_at (deleted_at)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pengurus;
-- +goose StatementEnd
