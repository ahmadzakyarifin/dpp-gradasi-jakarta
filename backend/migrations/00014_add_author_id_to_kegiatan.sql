-- +goose Up
ALTER TABLE kegiatan ADD COLUMN author_id INT NULL;
ALTER TABLE kegiatan ADD CONSTRAINT fk_kegiatan_author FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE SET NULL;

-- +goose Down
ALTER TABLE kegiatan DROP FOREIGN KEY fk_kegiatan_author;
ALTER TABLE kegiatan DROP COLUMN author_id;
