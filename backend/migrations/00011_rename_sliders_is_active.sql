-- +goose Up
-- +goose StatementBegin
ALTER TABLE sliders CHANGE COLUMN is_active is_published BOOLEAN DEFAULT TRUE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE sliders CHANGE COLUMN is_published is_active BOOLEAN DEFAULT TRUE;
-- +goose StatementEnd
