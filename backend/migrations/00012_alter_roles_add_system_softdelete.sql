-- +goose Up
-- +goose StatementBegin
-- Tambahkan kolom is_system & deleted_at ke tabel `roles`.
-- is_system: menandai role bawaan sistem (super_admin) yang tidak boleh
-- dihapus/dinonaktifkan. deleted_at: soft delete role (dipakai model
-- RoleModel, repo Delete/Restore/BulkDelete, dan filter status "trash").
ALTER TABLE roles
  ADD COLUMN is_system TINYINT(1) NOT NULL DEFAULT 0 AFTER description,
  ADD COLUMN deleted_at TIMESTAMP NULL AFTER updated_at,
  ADD INDEX idx_roles_deleted_at (deleted_at);

-- +goose StatementEnd

-- +goose StatementBegin
-- Role inti super_admin tidak boleh dihapus/nonaktifkan.
UPDATE roles SET is_system = 1 WHERE name = 'super_admin';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE roles
  DROP INDEX idx_roles_deleted_at,
  DROP COLUMN deleted_at,
  DROP COLUMN is_system;
-- +goose StatementEnd
