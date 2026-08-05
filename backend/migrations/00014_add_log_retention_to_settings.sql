-- +goose Up
ALTER TABLE `settings` ADD COLUMN `log_retention_days` INT DEFAULT 30;

-- +goose Down
ALTER TABLE `settings` DROP COLUMN `log_retention_days`;
