-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN email_pending VARCHAR(255) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN email_pending;
-- +goose StatementEnd
