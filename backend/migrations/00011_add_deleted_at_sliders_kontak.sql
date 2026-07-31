ALTER TABLE sliders ADD COLUMN deleted_at TIMESTAMP NULL;
CREATE INDEX idx_sliders_deleted_at ON sliders(deleted_at);

ALTER TABLE pesan_kontak ADD COLUMN deleted_at TIMESTAMP NULL;
CREATE INDEX idx_pesan_kontak_deleted_at ON pesan_kontak(deleted_at);
