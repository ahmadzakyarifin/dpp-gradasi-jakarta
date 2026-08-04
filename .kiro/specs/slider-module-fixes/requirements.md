# Requirements — Perbaikan Modul Slider

## Pendahuluan

Modul Slider mengelola banner carousel yang tampil di hero section halaman beranda publik. Saat ini modul ini punya beberapa cacat yang membuat admin tidak bisa mengelola slider dengan benar: gambar tidak bisa diunggah, slider non-aktif tetap tampil ke publik, dan semantik "non-aktif" tercampur dengan "terhapus" sehingga fitur Sampah/Pulihkan tidak berfungsi.

Dokumen ini mendefinisikan perilaku yang diharapkan untuk setiap perbaikan. Ruang lingkupnya mencakup backend (endpoint, service, repository, DTO, validasi) dan frontend (halaman admin slider, tampilan beranda publik).

### Istilah

- **Aktif / Non-aktif** — status publikasi slider, disimpan di kolom `is_active`. Slider non-aktif tidak tampil ke publik tetapi tetap ada dan bisa diaktifkan kembali oleh admin.
- **Terhapus (Sampah)** — slider yang di-soft-delete, ditandai kolom `deleted_at` terisi. Tidak tampil ke publik dan tidak muncul di daftar aktif, tetapi bisa dipulihkan.

Kedua konsep ini berbeda dan tidak boleh saling menggantikan.

---

## Requirement 1 — Unggah Gambar Slider dari Perangkat

**User story:** Sebagai admin, saya ingin mengunggah file gambar slider langsung dari komputer saya, supaya saya tidak perlu mencari hosting gambar eksternal terlebih dahulu.

### Acceptance Criteria

1. WHEN admin membuka form Tambah atau Edit slider, THEN sistem SHALL menyediakan tombol unggah file gambar.
2. WHEN admin memilih file gambar, THEN sistem SHALL mengunggah file tersebut ke server dan menyimpan path relatif hasil unggahan sebagai nilai `image_path` slider.
3. WHEN file gambar berhasil diunggah, THEN sistem SHALL menampilkan pratinjau gambar tersebut di dalam form.
4. IF file yang dipilih berukuran lebih dari 5 MB, THEN sistem SHALL menolak unggahan dan menampilkan pesan bahwa ukuran maksimal 5 MB.
5. IF file yang dipilih bukan bertipe PNG, JPG, atau WEBP, THEN sistem SHALL menolak unggahan dan menampilkan pesan format yang diizinkan.
6. WHEN admin sudah mengunggah gambar, THEN sistem SHALL menyediakan opsi untuk mengganti atau menghapus gambar tersebut.
7. WHERE slider sudah memiliki gambar dari URL eksternal (data lama), THE sistem SHALL tetap menampilkan dan mempertahankan URL tersebut tanpa memaksa admin mengunggah ulang.

---

## Requirement 2 — Gambar Slider Tampil Benar di Admin dan Beranda

**User story:** Sebagai admin, saya ingin gambar slider yang sudah saya simpan benar-benar tampil di tabel admin dan di beranda publik, supaya saya yakin perubahan saya berhasil.

### Acceptance Criteria

1. WHEN slider memiliki `image_path` berupa path relatif hasil unggahan, THEN sistem SHALL menampilkan gambar tersebut dengan menggabungkan path dengan base URL server.
2. WHEN slider memiliki `image_path` berupa URL absolut, THEN sistem SHALL menampilkan gambar tersebut apa adanya.
3. WHEN gambar slider gagal dimuat atau `image_path` kosong, THEN sistem SHALL menampilkan placeholder alih-alih gambar rusak.
4. WHEN admin menyimpan slider baru atau perubahan slider, THEN daftar slider di halaman admin SHALL menampilkan data terbaru tanpa perlu memuat ulang halaman secara manual.

---

## Requirement 3 — Beranda Publik Hanya Menampilkan Slider Aktif

**User story:** Sebagai pengelola situs, saya ingin slider yang saya non-aktifkan langsung berhenti tampil di beranda, supaya saya bisa menyembunyikan banner tanpa menghapusnya.

### Acceptance Criteria

1. WHEN pengunjung membuka halaman beranda, THEN sistem SHALL hanya menampilkan slider yang berstatus aktif dan tidak terhapus.
2. WHEN admin menonaktifkan sebuah slider, THEN slider tersebut SHALL berhenti tampil di beranda publik.
3. WHEN admin menghapus sebuah slider, THEN slider tersebut SHALL berhenti tampil di beranda publik.
4. WHEN slider aktif ditampilkan di beranda, THEN sistem SHALL mengurutkannya berdasarkan nilai urutan (`sort_order`) dari kecil ke besar.

---

## Requirement 4 — Pemisahan Status Non-aktif dan Terhapus

**User story:** Sebagai admin, saya ingin membedakan dengan jelas antara slider yang saya non-aktifkan dan slider yang saya hapus, supaya saya tidak bingung mencari slider yang hilang dari daftar.

### Acceptance Criteria

1. WHEN admin membuka tab daftar utama, THEN sistem SHALL menampilkan semua slider yang belum terhapus, baik yang aktif maupun yang non-aktif.
2. WHEN admin membuka tab Sampah, THEN sistem SHALL menampilkan hanya slider yang sudah terhapus.
3. WHEN admin menonaktifkan sebuah slider, THEN slider tersebut SHALL tetap berada di tab daftar utama dengan penanda status non-aktif.
4. WHEN admin menonaktifkan sebuah slider, THEN slider tersebut SHALL TIDAK muncul di tab Sampah.
5. WHEN admin menghapus sebuah slider, THEN slider tersebut SHALL berpindah dari tab daftar utama ke tab Sampah.

---

## Requirement 5 — Memulihkan Slider dari Sampah

**User story:** Sebagai admin, saya ingin memulihkan slider yang terhapus, supaya saya bisa membatalkan penghapusan yang tidak sengaja.

### Acceptance Criteria

1. WHEN admin membuka tab Sampah, THEN sistem SHALL mengambil dan menampilkan data slider yang terhapus dari server.
2. WHEN admin memulihkan sebuah slider dari tab Sampah, THEN sistem SHALL mengembalikan slider tersebut ke tab daftar utama.
3. WHEN admin memilih beberapa slider di tab Sampah dan memilih pulihkan massal, THEN sistem SHALL memulihkan semua slider yang dipilih.
4. WHEN slider berhasil dipulihkan, THEN sistem SHALL mempertahankan status aktif atau non-aktif yang dimilikinya sebelum dihapus.
5. WHEN tab Sampah tidak memiliki data, THEN sistem SHALL menampilkan pesan bahwa Sampah kosong.

---

## Requirement 6 — Mengubah Status Aktif Slider Kapan Saja

**User story:** Sebagai admin, saya ingin bisa mengaktifkan kembali slider yang sudah saya non-aktifkan, supaya saya bisa memakai ulang banner lama tanpa membuat data baru.

### Acceptance Criteria

1. WHEN slider ditampilkan di tab daftar utama, THEN sistem SHALL menyediakan kontrol untuk mengubah statusnya antara aktif dan non-aktif.
2. WHEN admin mengubah status slider menjadi aktif, THEN slider tersebut SHALL mulai tampil di beranda publik.
3. WHEN admin mengubah status slider, THEN sistem SHALL meminta konfirmasi sebelum menerapkan perubahan.
4. WHEN admin mengubah status slider, THEN sistem SHALL mempertahankan seluruh data slider lainnya tanpa perubahan.
5. WHILE slider berada di tab Sampah, THE sistem SHALL TIDAK menyediakan kontrol pengubahan status aktif.

---

## Requirement 7 — Filter dan Pencarian yang Konsisten

**User story:** Sebagai admin, saya ingin filter status dan pencarian bekerja sesuai harapan di dalam tab yang sedang saya buka, supaya saya tidak mendapat daftar kosong tanpa alasan jelas.

### Acceptance Criteria

1. WHEN admin memilih filter status Aktif di tab daftar utama, THEN sistem SHALL menampilkan hanya slider aktif yang belum terhapus.
2. WHEN admin memilih filter status Non-aktif di tab daftar utama, THEN sistem SHALL menampilkan hanya slider non-aktif yang belum terhapus.
3. WHEN admin memilih filter Semua Status di tab daftar utama, THEN sistem SHALL menampilkan seluruh slider yang belum terhapus.
4. WHEN admin membuka tab Sampah, THEN sistem SHALL TIDAK menampilkan kontrol filter status.
5. WHEN admin mengisi kata kunci pencarian, THEN sistem SHALL memfilter slider berdasarkan judul dan sub-judul di dalam tab yang sedang aktif.
6. WHEN admin menekan tombol Reset, THEN sistem SHALL mengosongkan kata kunci pencarian dan mengembalikan filter status ke Semua Status.
7. WHEN hasil filter atau pencarian kosong, THEN sistem SHALL menampilkan pesan bahwa tidak ada data yang cocok.

---

## Requirement 8 — Validasi Input Slider

**User story:** Sebagai admin, saya ingin mendapat pesan kesalahan yang jelas saat mengisi form slider dengan tidak benar, supaya saya tahu apa yang harus diperbaiki.

### Acceptance Criteria

1. WHEN admin menyimpan form tanpa mengisi judul, THEN sistem SHALL menolak penyimpanan dan menampilkan pesan bahwa judul wajib diisi.
2. WHEN admin menyimpan form tanpa gambar slider, THEN sistem SHALL menolak penyimpanan dan menampilkan pesan bahwa gambar wajib diisi.
3. IF judul melebihi 200 karakter, THEN sistem SHALL menolak penyimpanan dan menampilkan pesan batas panjang judul.
4. IF nilai urutan bukan bilangan bulat nol atau positif, THEN sistem SHALL menolak penyimpanan dan menampilkan pesan bahwa urutan harus berupa angka yang valid.
5. IF tautan pada kolom Link URL bukan URL yang valid, THEN sistem SHALL menolak penyimpanan dan menampilkan pesan format tautan tidak valid.
6. WHEN validasi gagal, THEN sistem SHALL menampilkan pesan kesalahan di bawah kolom yang bermasalah dan tidak mengirim data ke server.
7. WHEN server menolak data karena validasi, THEN sistem SHALL menampilkan pesan kesalahan dari server pada kolom yang bersangkutan.
8. WHEN admin memperbaiki isi kolom yang bermasalah, THEN sistem SHALL menghapus pesan kesalahan pada kolom tersebut.

---

## Requirement 9 — Pengaturan Urutan Slider

**User story:** Sebagai admin, saya ingin mengatur urutan tampil slider dengan hasil yang konsisten, supaya susunan banner di beranda sesuai rencana saya.

### Acceptance Criteria

1. WHEN admin mengubah nilai urutan sebuah slider dari tabel, THEN sistem SHALL menyimpan nilai tersebut dan memuat ulang daftar dengan urutan terbaru.
2. WHEN admin mengubah nilai urutan sebuah slider, THEN sistem SHALL mempertahankan seluruh data slider lainnya tanpa perubahan.
3. WHEN daftar slider ditampilkan di halaman admin, THEN sistem SHALL mengurutkannya berdasarkan nilai urutan dari kecil ke besar.
4. IF nilai urutan yang dimasukkan tidak valid, THEN sistem SHALL membatalkan perubahan dan mempertahankan nilai sebelumnya.
