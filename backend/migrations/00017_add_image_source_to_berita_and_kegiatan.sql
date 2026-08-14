-- +goose Up
-- +goose StatementBegin
ALTER TABLE berita ADD COLUMN image_source VARCHAR(250) NULL;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE kegiatan ADD COLUMN image_source VARCHAR(250) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE berita DROP COLUMN image_source;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE kegiatan DROP COLUMN image_source;
-- +goose StatementEnd
