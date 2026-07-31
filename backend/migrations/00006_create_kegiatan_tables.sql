-- +goose Up
-- +goose StatementBegin
CREATE TABLE kegiatan (
    id INT AUTO_INCREMENT PRIMARY KEY,
    slug VARCHAR(250) NOT NULL UNIQUE,
    title VARCHAR(300) NOT NULL,
    category VARCHAR(100) DEFAULT 'Kegiatan',
    event_date DATE,
    location VARCHAR(200),
    organizer VARCHAR(200),
    image_url VARCHAR(500),
    excerpt TEXT,
    content LONGTEXT,
    is_published BOOLEAN DEFAULT TRUE,
    views INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,

    INDEX idx_event_date (event_date DESC),
    INDEX idx_is_published (is_published),
    INDEX idx_deleted_at (deleted_at),
    FULLTEXT INDEX ft_kegiatan (title, content)
);

CREATE TABLE kegiatan_tags (
    id INT AUTO_INCREMENT PRIMARY KEY,
    kegiatan_id INT NOT NULL,
    tag VARCHAR(100) NOT NULL,

    FOREIGN KEY (kegiatan_id) REFERENCES kegiatan(id) ON DELETE CASCADE,
    INDEX idx_kegiatan_id (kegiatan_id),
    INDEX idx_tag (tag)
);

CREATE TABLE kegiatan_gallery (
    id INT AUTO_INCREMENT PRIMARY KEY,
    kegiatan_id INT NOT NULL,
    image_url VARCHAR(500) NOT NULL,
    caption VARCHAR(200),
    sort_order INT DEFAULT 0,

    FOREIGN KEY (kegiatan_id) REFERENCES kegiatan(id) ON DELETE CASCADE,
    INDEX idx_kegiatan_id (kegiatan_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS kegiatan_gallery;
DROP TABLE IF EXISTS kegiatan_tags;
DROP TABLE IF EXISTS kegiatan;
-- +goose StatementEnd
