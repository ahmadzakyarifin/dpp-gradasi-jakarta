-- +goose Up
-- +goose StatementBegin

CREATE TABLE import_batches (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    file_attachment_id BIGINT UNSIGNED NULL,
    import_type ENUM('students', 'guardians', 'class_templates', 'active_classes', 'invoices', 'payments', 'mixed') NOT NULL,
    status ENUM('uploaded', 'validating', 'ready', 'importing', 'completed', 'failed', 'cancelled') NOT NULL DEFAULT 'uploaded',
    total_rows INT UNSIGNED NOT NULL DEFAULT 0,
    success_rows INT UNSIGNED NOT NULL DEFAULT 0,
    failed_rows INT UNSIGNED NOT NULL DEFAULT 0,
    uploaded_by BIGINT UNSIGNED NULL,
    started_at DATETIME(3) NULL,
    completed_at DATETIME(3) NULL,
    error_message TEXT NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    CONSTRAINT fk_import_batches_file_attachment FOREIGN KEY (file_attachment_id) REFERENCES file_attachments(id) ON DELETE SET NULL,
    CONSTRAINT fk_import_batches_uploaded_by FOREIGN KEY (uploaded_by) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_import_batches_file_attachment_id ON import_batches(file_attachment_id);
CREATE INDEX idx_import_batches_import_type ON import_batches(import_type);
CREATE INDEX idx_import_batches_status ON import_batches(status);

CREATE TABLE import_rows (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    import_batch_id BIGINT UNSIGNED NOT NULL,
    `row_number` INT UNSIGNED NOT NULL,
    raw_data JSON NOT NULL,
    mapped_data JSON NULL,
    status ENUM('pending', 'valid', 'invalid', 'imported', 'failed', 'skipped') NOT NULL DEFAULT 'pending',
    error_messages JSON NULL,
    target_type VARCHAR(100) NULL,
    target_id BIGINT UNSIGNED NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    CONSTRAINT fk_import_rows_import_batch FOREIGN KEY (import_batch_id) REFERENCES import_batches(id) ON DELETE CASCADE,
    CONSTRAINT uq_import_rows_batch_row UNIQUE (import_batch_id, `row_number`),
    CONSTRAINT chk_import_rows_raw_data CHECK (JSON_VALID(raw_data)),
    CONSTRAINT chk_import_rows_mapped_data CHECK (mapped_data IS NULL OR JSON_VALID(mapped_data)),
    CONSTRAINT chk_import_rows_error_messages CHECK (error_messages IS NULL OR JSON_VALID(error_messages))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_import_rows_import_batch_id ON import_rows(import_batch_id);
CREATE INDEX idx_import_rows_status ON import_rows(status);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS import_rows;
DROP TABLE IF EXISTS import_batches;
-- +goose StatementEnd
