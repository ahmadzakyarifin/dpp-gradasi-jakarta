-- +goose Up
-- +goose StatementBegin
ALTER TABLE pengurus
    ADD COLUMN pekerjaan VARCHAR(150),
    ADD COLUMN bio TEXT,
    ADD COLUMN pendidikan TEXT,
    ADD COLUMN sertifikasi TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE pengurus
    DROP COLUMN pekerjaan,
    DROP COLUMN bio,
    DROP COLUMN pendidikan,
    DROP COLUMN sertifikasi;
-- +goose StatementEnd
