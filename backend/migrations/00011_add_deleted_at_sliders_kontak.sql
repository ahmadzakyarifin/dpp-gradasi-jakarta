-- +goose Up
ALTER TABLE sliders ADD COLUMN deleted_at TIMESTAMP NULL;
CREATE INDEX idx_sliders_deleted_at ON sliders(deleted_at);

-- +goose Down
ALTER TABLE sliders DROP INDEX idx_sliders_deleted_at;
ALTER TABLE sliders DROP COLUMN deleted_at;
