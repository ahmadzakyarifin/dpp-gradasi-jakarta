-- +goose Up
-- +goose StatementBegin
-- Selaraskan tabel refresh_tokens dengan model RefreshTokenModel (auth module).
-- Migration 00003 membuat device_info, tetapi model & repo memakai:
--   device_name (SaveRefreshToken), user_agent (SaveRefreshToken),
--   revoked_at  (FindUserByRefreshToken / DeleteRefreshToken /
--                DeleteAllUserRefreshTokens).
-- Tanpa kolom ini, login & logout gagal dengan SQL error
-- "Unknown column 'device_name'/'revoked_at'".
ALTER TABLE refresh_tokens
  CHANGE COLUMN device_info device_name VARCHAR(255) NULL,
  ADD COLUMN user_agent VARCHAR(255) NULL AFTER device_name,
  ADD COLUMN revoked_at TIMESTAMP NULL AFTER expires_at;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE refresh_tokens
  DROP COLUMN revoked_at,
  DROP COLUMN user_agent,
  CHANGE COLUMN device_name device_info VARCHAR(255) NULL;
-- +goose StatementEnd
