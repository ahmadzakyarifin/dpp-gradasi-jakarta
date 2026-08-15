-- +goose Up
ALTER TABLE pengurus ADD COLUMN kepengurusan ENUM('Ketua', 'Anggota') NOT NULL DEFAULT 'Anggota' AFTER role;

-- +goose Down
ALTER TABLE pengurus DROP COLUMN kepengurusan;
