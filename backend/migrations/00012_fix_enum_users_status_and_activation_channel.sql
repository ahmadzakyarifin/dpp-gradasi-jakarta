-- +goose Up
-- +goose StatementBegin
-- 1) Tambah nilai 'pending_activation' ke ENUM users.status
--    CreateAdmin (invite) menyimpan status 'pending_activation' sebelum admin
--    mengaktifkan akun via link email. Sebelum migrasi ini, INSERT akan gagal
--    dengan MySQL error 1265 (Data truncated for column 'status').
ALTER TABLE users
    MODIFY COLUMN status ENUM('active', 'inactive', 'pending_activation') NOT NULL DEFAULT 'inactive';

-- 2) Aktivasi hanya via email (keputusan: WAHA & channel 'all' dihapus)
ALTER TABLE activation_tokens
    MODIFY COLUMN channel ENUM('email') NOT NULL DEFAULT 'email';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE activation_tokens
    MODIFY COLUMN channel ENUM('email', 'whatsapp', 'all') NOT NULL DEFAULT 'email';

-- Catatan: data dengan status 'pending_activation' akan dipaksa ke 'inactive'
-- saat rollback (perilaku ENUM MySQL). Aman karena 'pending_activation' hanya
-- dipakai oleh akun undangan yang belum aktif.
ALTER TABLE users
    MODIFY COLUMN status ENUM('active', 'inactive') NOT NULL DEFAULT 'inactive';
-- +goose StatementEnd
