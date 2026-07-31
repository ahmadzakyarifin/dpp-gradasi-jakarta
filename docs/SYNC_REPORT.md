# Laporan Sinkronisasi 4 Lapisan — dpp_gradasi

> Status: **SELESAI DIPERBAIKI** (per 31 Juli 2026)
> Lapisan: DB (migrasi goose) → Backend Go (Gin+GORM) → Contract (`docs/api/*.jsonc`) → UI statis (`admin/*.html`, halaman publik).

---

## 1. users.status ENUM (SELESAI)

- **Masalah**: DB cuma `('active','inactive')`, backend & contract pakai `pending_activation` → CreateAdmin gagal (MySQL 1265).
- **Fix**:
  - Migration baru `backend/migrations/00012_fix_enum_users_status_and_activation_channel.sql`:
    - `users.status` → `ENUM('active','inactive','pending_activation') NOT NULL DEFAULT 'inactive'`
    - `activation_tokens.channel` → `ENUM('email')` (keputusan aktivasi hanya via email; WAHA & 'all' dihapus).
  - GORM model `user_model.go` disinkronkan.
  - `AdminCreateRequest.RoleID` kini `required,oneof=2 3 4` (super_admin tidak bisa dibuat via undangan).
  - Email undangan dikirim **sinkron** (`mailer.Send`): jika SMTP gagal → user + token dihapus (rollback), response `MAIL_SEND_FAILED`.

## 2. Settings Update rusak (SELESAI)

- **Masalah**: handler pakai `reflect FieldByName("site_name")` → selalu 400 "Field tidak dikenal".
- **Fix**: `settings_service.go` dibangun ulang dengan pemetaan **JSON tag snake_case → field Go**; key tak dikenal → 400; validasi `_email`/`_url`; `about_mission` diterima sebagai array string (disimpan JSON).
- `SettingsResponse` DTO dilengkapi: `address`, `video_profile_url`, `history`, `greeting_title/subtitle/date/content/image_url` (+ semua field contract).
- `settings_repo.go`: `Update` no-op bila map kosong.
- Contract `settings.jsonc` disinkronkan (`about_mission` array).

## 3. Activity Log (SELESAI)

- DTO query kini mendukung contract: `per_page`, `actor_id`, `status`, `start_date`, `end_date`, `sort_by`, `order` + alias `entity_type`/`risk_level`.
- Repo: normalisasi `limit`/`per_page`, filter actor_id/status/rentang tanggal, sorting whitelist via `sortClause`.
- **`risk_level` hanya `low|medium|high`** (nilai `CRITICAL` tidak valid di DB) — contract & UI admin/activity-log.html sudah memakai 3 nilai tersebut.
- Endpoint baru `GET /api/v1/activity-logs/summary` untuk dashboard.
- Contract `activity_logs.jsonc`: response list mengikuti bentuk backend (`summary`+`pagination`+`items`).

## 4. Kegiatan (SELESAI)

- `list()` query ganda dihapus → panggil `FindPublished` XOR `FindAll` sesuai mode.
- Update **tidak** mengubah slug lagi (konsisten dengan Berita — link lama tidak 404).
- Filter status `published|draft|trashed|all` di repo (admin).

## 5. Berita (SELESAI)

- Filter status `published|draft|trashed|all` di `berita_repo.go` (admin mode).

## 6. Pengurus (SELESAI)

- Validasi `required_if` (provinsi wajib utk `dpd`/`dpc`, kabupaten utk `dpc`) di service (`ValidateRegionRules`) — sesuai contract.
- Form admin/pengurus.html sudah punya field `department`, `periode`, `provinsi`, `kabupaten`, `sort_order`, `role` (nama field disinkronkan: `position`→`role`, `sortOrder`→`sort_order`, `linkedin`→`linkedin_url`).
- `handleUpload` kini `MkdirAll` direktori upload sebelum menulis file.

## 7. Duplikasi change-password (SELESAI)

- Diputuskan: **`POST /api/v1/auth/change-password`** (auth) yang dipakai.
- `PUT /api/v1/profile/password` dihapus dari user routes.

## 8. File kontak_settings.jsonc (SELESAI)

- Dihapus (versi lama, bentrok dengan `kontak.jsonc`).

## 9. Dashboard (BARU)

- Modul `internal/module/dashboard` (handler+service+DTO) dengan `GET /api/v1/admin/dashboard/summary` (super_admin):
  `total_berita, total_kegiatan, total_pengurus, total_kontak, unread_kontak, total_admin, pending_admin, activity_logs, high_risk_logs, failed_login, cms_actions`.
- Contract `docs/api/dashboard.jsonc`; UI `admin/index.html` di-fetch ke endpoint (fallback 0 bila tanpa token/error).

## 10. Verifikasi

- `go build ./...` ✅, `go vet ./...` ✅, `gofmt -l .` kosong ✅.
- Semua 10 file `docs/api/*.jsonc` valid (stripper JSONC state-machine — komentar `//` dalam URL tidak dianggap komentar). ✅
- `admin/activity-log.html` & contract: tidak ada lagi `CRITICAL`; dropdown risiko = high/medium/low. ✅
- `admin/users.html`: tab "Menunggu Aktivasi" (pending_activation) + aksi "Kirim Ulang Undangan" + nonaktifkan; dropdown role = role_id 2/3/4. ✅

## Catatan UI statis

- `admin/*.html` tetap UI gambaran (mock data) sesuai arahan: React nanti konsumsi contract `docs/api/*.jsonc`. Hanya `index.html` (dashboard) & `reset-password.html` yang sudah connect API.
- `pending_activation` = "akun dibuat via undangan, belum aktivasi lewat email" — **bukan** nonaktif. `inactive` = super admin nonaktifkan manual.
