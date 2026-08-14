-- +goose Up
-- +goose StatementBegin
ALTER TABLE berita ADD COLUMN footnote VARCHAR(500) NULL;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE kegiatan ADD COLUMN footnote VARCHAR(500) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE berita DROP COLUMN footnote;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE kegiatan DROP COLUMN footnote;
-- +goose StatementEnd
