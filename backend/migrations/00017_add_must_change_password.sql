-- +goose Up
-- +goose StatementBegin
-- Tambah kolom must_change_password: menandai admin yang harus mengganti
-- password default pada login pertama (alur undangan via email).
ALTER TABLE users
    ADD COLUMN must_change_password TINYINT(1) NOT NULL DEFAULT 0 AFTER status;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
    DROP COLUMN must_change_password;
-- +goose StatementEnd
