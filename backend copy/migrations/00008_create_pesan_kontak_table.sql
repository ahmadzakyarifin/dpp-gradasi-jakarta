-- +goose Up
-- +goose StatementBegin
CREATE TABLE pesan_kontak (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nama VARCHAR(100) NOT NULL,
    email VARCHAR(150) NOT NULL,
    subjek VARCHAR(200) NOT NULL,
    pesan TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    response_note TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,

    INDEX idx_is_read (is_read),
    INDEX idx_created_at (created_at DESC),
    INDEX idx_deleted_at (deleted_at)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pesan_kontak;
-- +goose StatementEnd
