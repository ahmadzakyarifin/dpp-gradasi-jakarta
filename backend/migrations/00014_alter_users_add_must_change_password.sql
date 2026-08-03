-- +goose Up
-- +goose StatementBegin
-- Perbaiki schema drift: DB yang dibuat sebelum Fase 1 tidak punya
-- kolom must_change_password (file 00002 sudah di-update di Fase 1,
-- tapi goose tidak re-run file yang sudah applied — migration immutable).
-- MySQL 8.0 TIDAK mendukung ADD COLUMN IF NOT EXISTS (hanya MariaDB),
-- jadi lakukan guard manual via stored procedure idempoten.
-- Fresh DB (kolom sudah ada dari 00002) aman, DB drift lama ter-patch.
-- +goose StatementEnd

-- +goose StatementBegin
SET @col_exists := (SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'users'
    AND COLUMN_NAME = 'must_change_password');
-- +goose StatementEnd

-- +goose StatementBegin
SET @ddl := IF(@col_exists = 0,
  'ALTER TABLE users ADD COLUMN must_change_password TINYINT(1) NOT NULL DEFAULT 0 AFTER status',
  'SELECT 1');
-- +goose StatementEnd

-- +goose StatementBegin
PREPARE stmt FROM @ddl;
-- +goose StatementEnd

-- +goose StatementBegin
EXECUTE stmt;
-- +goose StatementEnd

-- +goose StatementBegin
DEALLOCATE PREPARE stmt;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET @col_exists := (SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'users'
    AND COLUMN_NAME = 'must_change_password');
-- +goose StatementEnd

-- +goose StatementBegin
SET @ddl := IF(@col_exists = 1,
  'ALTER TABLE users DROP COLUMN must_change_password',
  'SELECT 1');
-- +goose StatementEnd

-- +goose StatementBegin
PREPARE stmt FROM @ddl;
-- +goose StatementEnd

-- +goose StatementBegin
EXECUTE stmt;
-- +goose StatementEnd

-- +goose StatementBegin
DEALLOCATE PREPARE stmt;
-- +goose StatementEnd
