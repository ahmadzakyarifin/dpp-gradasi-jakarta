-- +goose Up
ALTER TABLE pengurus ADD COLUMN twitter_url VARCHAR(500) DEFAULT NULL AFTER linkedin_url;

-- +goose Down
ALTER TABLE pengurus DROP COLUMN twitter_url;
