-- +goose Up
ALTER TABLE `sliders` DROP COLUMN `link_url`;

-- +goose Down
ALTER TABLE `sliders` ADD COLUMN `link_url` VARCHAR(255) DEFAULT NULL;
