ALTER TABLE pengurus MODIFY COLUMN level VARCHAR(50) NOT NULL;

UPDATE pengurus SET level = 'Ketua Umum' WHERE level = 'ketua';
UPDATE pengurus SET level = 'Pengurus Pusat' WHERE level = 'dpp';
UPDATE pengurus SET level = 'Pengurus Provinsi' WHERE level = 'dpd';
UPDATE pengurus SET level = 'Pengurus Kab/Kota' WHERE level = 'dpc';

ALTER TABLE pengurus MODIFY COLUMN level ENUM('Ketua Umum', 'Pengurus Pusat', 'Pengurus Provinsi', 'Pengurus Kab/Kota') NOT NULL;
