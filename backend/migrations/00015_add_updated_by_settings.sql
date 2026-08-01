-- +goose Up
ALTER TABLE settings ADD COLUMN updated_by INT NULL AFTER updated_at;
ALTER TABLE settings ADD CONSTRAINT fk_settings_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL;

-- +goose Down
ALTER TABLE settings DROP FOREIGN KEY fk_settings_updated_by;
ALTER TABLE settings DROP COLUMN updated_by;
