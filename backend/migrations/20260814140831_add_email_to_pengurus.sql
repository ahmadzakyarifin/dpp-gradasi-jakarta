-- +goose Up
-- +goose StatementBegin
ALTER TABLE pengurus
    ADD COLUMN email VARCHAR(150);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE pengurus
    DROP COLUMN email;
-- +goose StatementEnd
