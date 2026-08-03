-- +goose Up
-- +goose StatementBegin
ALTER TABLE settings 
ADD COLUMN IF NOT EXISTS login_hero_title VARCHAR(255) NULL,
ADD COLUMN IF NOT EXISTS login_hero_description TEXT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE settings 
DROP COLUMN IF EXISTS login_hero_title,
DROP COLUMN IF EXISTS login_hero_description;
-- +goose StatementEnd
