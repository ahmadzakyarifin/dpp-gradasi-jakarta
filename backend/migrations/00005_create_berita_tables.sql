-- +goose Up
-- +goose StatementBegin
CREATE TABLE berita (
    id INT AUTO_INCREMENT PRIMARY KEY,
    slug VARCHAR(250) NOT NULL UNIQUE,
    title VARCHAR(300) NOT NULL,
    category VARCHAR(100) DEFAULT 'Berita Organisasi',
    published_date DATE NOT NULL,
    author_id INT,
    image_path VARCHAR(500),
    excerpt TEXT,
    content LONGTEXT,
    is_featured BOOLEAN DEFAULT FALSE,
    is_published BOOLEAN DEFAULT TRUE,
    views INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,

    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE SET NULL,
    INDEX idx_published_date (published_date DESC),
    INDEX idx_category (category),
    INDEX idx_is_published (is_published),
    INDEX idx_deleted_at (deleted_at),
    FULLTEXT INDEX ft_berita (title, content)
);

CREATE TABLE berita_tags (
    id INT AUTO_INCREMENT PRIMARY KEY,
    berita_id INT NOT NULL,
    tag VARCHAR(100) NOT NULL,

    FOREIGN KEY (berita_id) REFERENCES berita(id) ON DELETE CASCADE,
    INDEX idx_berita_id (berita_id),
    INDEX idx_tag (tag)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS berita_tags;
DROP TABLE IF EXISTS berita;
-- +goose StatementEnd
