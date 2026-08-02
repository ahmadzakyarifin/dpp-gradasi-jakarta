-- +goose Up
-- +goose StatementBegin

CREATE TABLE file_attachments (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    uploaded_by BIGINT UNSIGNED NULL,
    attachable_type VARCHAR(100) NULL,
    attachable_id BIGINT UNSIGNED NULL,
    original_name VARCHAR(255) NOT NULL,
    stored_name VARCHAR(255) NOT NULL,
    file_path TEXT NOT NULL,
    mime_type VARCHAR(150) NULL,
    file_size BIGINT UNSIGNED NULL,
    disk VARCHAR(50) NOT NULL DEFAULT 'local',
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    CONSTRAINT fk_file_attachments_uploaded_by FOREIGN KEY (uploaded_by) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_file_attachments_uploaded_by ON file_attachments(uploaded_by);
CREATE INDEX idx_file_attachments_attachable ON file_attachments(attachable_type, attachable_id);

ALTER TABLE customer_support_messages
    ADD CONSTRAINT fk_customer_support_messages_attachment
    FOREIGN KEY (attachment_id) REFERENCES file_attachments(id)
    ON DELETE SET NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE customer_support_messages DROP FOREIGN KEY fk_customer_support_messages_attachment;
DROP TABLE IF EXISTS file_attachments;
-- +goose StatementEnd
