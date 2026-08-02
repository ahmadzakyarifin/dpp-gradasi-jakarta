-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN date_of_birth DATE NULL AFTER phone,
    ADD COLUMN country_code VARCHAR(2) NULL AFTER date_of_birth;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
    DROP COLUMN country_code,
    DROP COLUMN date_of_birth;
-- +goose StatementEnd
