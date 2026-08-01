-- +goose Up
-- Normalisasi nama role ke snake_case (kontrak middleware & JWT claims pakai
-- "super_admin" / "admin", bukan "Super Admin" / "Admin" dengan spasi).
--
-- Latar belakang: migration 00013 (seed admin_berita/admin_kegiatan) sudah
-- membuat role id 5/6 dengan nama snake_case. Sementara role id 1/2/3/4 dari
-- migration 00001 masih bernama "Super Admin"/"Admin"/"Admin Berita"/"Admin Kegiatan"
-- (dengan spasi) yang TIDAK cocok dengan RoleMiddleware("super_admin","admin").
--
-- Strategi:
--   1. Rename role 1 & 2 ke snake_case (super_admin, admin).
--   2. Hapus role 3 & 4 (Admin Berita/Admin Kegiatan versi lama) karena tidak
--      ada user yang memakainya, lalu ganti referensinya ke role 5 & 6 yang
--      sudah snake_case.
UPDATE roles SET name = 'super_admin' WHERE name = 'Super Admin';
UPDATE roles SET name = 'admin' WHERE name = 'Admin';
DELETE FROM roles WHERE name IN ('Admin Berita', 'Admin Kegiatan');

-- +goose Down
-- Kembalikan nama lama untuk role yang masih ada.
UPDATE roles SET name = 'Super Admin' WHERE name = 'super_admin';
UPDATE roles SET name = 'Admin' WHERE name = 'admin';
