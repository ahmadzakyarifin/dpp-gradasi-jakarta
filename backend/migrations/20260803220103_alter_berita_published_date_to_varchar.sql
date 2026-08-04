-- +goose Up
-- +goose StatementBegin
ALTER TABLE berita MODIFY COLUMN published_date VARCHAR(100) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE berita MODIFY COLUMN published_date DATE NOT NULL;
-- +goose StatementEnd
