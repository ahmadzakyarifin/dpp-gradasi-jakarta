-- +goose Up
-- +goose StatementBegin
-- Perbaiki schema drift: DB yang dibuat sebelum Fase 1 tidak punya
-- kolom must_change_password (file 00002 sudah di-update di Fase 1,
-- tapi goose tidak re-run file yang sudah applied — migration immutable).
-- IF NOT EXISTS: idempoten — fresh DB (kolom sudah ada dari 00002) aman,
-- DB drift lama tetap ter-patch. Didukung MariaDB 10.0.2+.
ALTER TABLE users
  ADD COLUMN IF NOT EXISTS must_change_password TINYINT(1) NOT NULL DEFAULT 0
  AFTER status;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
  DROP COLUMN IF EXISTS must_change_password;
-- +goose StatementEnd
