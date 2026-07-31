-- +goose Up
-- +goose StatementBegin
-- Seed role admin_berita & admin_kegiatan (dipakai dropdown CreateAdmin di UI:
-- role_id 3 = Admin Berita, role_id 4 = Admin Kegiatan).
-- Sebelum migrasi ini, role 3 & 4 TIDAK ada di tabel roles sehingga
-- CreateAdmin dengan role_id 3/4 gagal (foreign key violation).
INSERT INTO roles (name, display_name, description) VALUES
('admin_berita', 'Admin Berita', 'Mengelola konten berita'),
('admin_kegiatan', 'Admin Kegiatan', 'Mengelola konten kegiatan')
ON DUPLICATE KEY UPDATE display_name = VALUES(display_name), description = VALUES(description);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM roles WHERE name IN ('admin_berita', 'admin_kegiatan');
-- +goose StatementEnd
