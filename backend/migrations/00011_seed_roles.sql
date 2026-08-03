-- +goose Up
-- +goose StatementBegin
-- Seed data untuk tabel `roles`. Nama role ini dirujuk langsung oleh seluruh
-- dokumentasi API di docs-final/api/*.jsonc (lihat setiap "security.roles"),
-- jadi WAJIB ada sebelum sistem otorisasi berbasis role bisa berfungsi.
INSERT INTO roles (name, display_name, description, is_active) VALUES
    ('super_admin', 'Super Administrator', 'Akses penuh ke seluruh modul, termasuk manajemen admin lain dan activity logs.', TRUE),
    ('admin', 'Administrator', 'Akses ke modul konten (berita, kegiatan, pengurus, sliders, kontak, settings).', TRUE),

-- +goose Down
-- +goose StatementBegin
DELETE FROM roles WHERE name IN ('super_admin', 'admin', 'editor', 'admin_berita');
-- +goose StatementEnd
