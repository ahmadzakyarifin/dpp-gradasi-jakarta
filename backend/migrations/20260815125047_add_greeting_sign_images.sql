-- +goose Up
ALTER TABLE settings ADD COLUMN greeting_sign_image_1 VARCHAR(500) DEFAULT '';
ALTER TABLE settings ADD COLUMN greeting_sign_image_2 VARCHAR(500) DEFAULT '';

-- +goose Down
ALTER TABLE settings DROP COLUMN greeting_sign_image_1;
ALTER TABLE settings DROP COLUMN greeting_sign_image_2;
