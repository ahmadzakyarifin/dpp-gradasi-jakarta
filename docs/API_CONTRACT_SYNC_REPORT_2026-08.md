# Laporan Final: Sinkronisasi API Contract ↔ Backend (FINALIZED untuk Project Baru)

**Tanggal:** 2026-08-01
**Project:** DPP GRADASI (`/home/ahmadzakyarifin/Documents/poter/dpp-gradasi-jakarta`)
**Tujuan:** Finalisasi `docs/api/*.jsonc` sebagai source of truth valid untuk membangun project baru serupa dari awal.

---

## ✅ Hasil Akhir

| Metrik | Nilai |
|--------|-------|
| File JSONC valid | **10/10 (ALL_VALID)** — validator state-machine |
| Build backend `go build ./...` | **BUILD_OK** |
| Total endpoint contract | **79** |
| Total route BE | **79** |
| Endpoint sinkron contract ↔ BE | **79/79 (100%)** |
| Referensi `phone` (selain contact_phone settings) | **0 (bersih)** |
| camelCase keys anomali | **0** (hanya `entityLabel`/`totalPages` yang memang sesuai DTO BE) |

---

## Perubahan yang Dilakukan (Finalisasi Contract)

### 1. `settings.jsonc` — method + endpoint baru
- ✅ `POST /api/v1/admin/settings` → **`PUT /api/v1/admin/settings`** (sinkron dengan BE route `admin.PUT("")`)
- ✅ Tambah endpoint baru: **`POST /api/v1/admin/settings/logo`** (upload logo, multipart, role super_admin/admin) — sebelumnya undocumented

### 2. `users.jsonc` — role_id + endpoint baru
- ✅ `role_id: required,oneof=2 3 4` → **`required,oneof=2 5 6`** (sinkron dengan DTO BE `AdminCreateRequest`; role 3/4 dihapus migration 00016, diganti 5/6)
- ✅ Tambah endpoint baru: **`PUT /api/v1/profile/password`** (ganti password sendiri) — sebelumnya undocumented

### 3. `auth.jsonc` — hapus phone + tambah must_change_password
- ✅ Hapus `"phone": "6285279880008"` dari response `GET /auth/me` (backend TIDAK punya field phone; keputusan user: tidak pakai phone)
- ✅ Tambah `"must_change_password": false` di response `GET /auth/me` (sinkron dengan DTO `AuthUserResponse`)

### 4. `berita.jsonc` — query param
- ✅ Tambah `status` di `GET /berita` (published | draft | trashed) — BE repo `query()` memakai q.Status

### 5. `kegiatan.jsonc` — query param
- ✅ Tambah `status` di `GET /kegiatan` (published | draft | trashed) — BE repo memakai q.Status

### 6. `pengurus.jsonc` — query param
- ✅ Tambah `kabupaten` di `GET /admin/pengurus` — BE repo memakai q.Kabupaten

### 7. `activity_logs.jsonc` — response shape disinkronkan ke BE
- ✅ Detail `GET /activity-logs/:id`: shape lama `actor_id/actor_name/actor_role/status/entity_type/risk_level/ip_address/user_agent/created_at` → **shape BE `id/time/actor/role/action/entity/entityLabel/description/ip/device/risk/metadata`** (sesuai `ActivityLogDetailRes` + mapper)
- ✅ Entity logs `GET /activity-logs/entity/:entityType/:entityID`: shape lama → **shape BE `[]ActivityLogItemRes`** (id/time/actor/role/entity/entityLabel/...)

---

## Catatan Penting untuk Project Baru

### Status BE yang masih perlu diperbaiki (bukan contract — contract sudah benar):
1. **`user_repo.go` `FindAllAdmins`** pakai `role_id IN (1, 2, 3, 4)` — role 3/4 sudah TIDAK ada di DB (diganti 5/6) → **admin_berita/admin_kegiatan tidak muncul di list admin**. Fix: `role_id IN (1, 2, 5, 6)`.
2. **Kegiatan masih SATU DTO** `KegiatanRequest` untuk create+update (`binding:"required"` di title/content) → toggle publish `PUT /kegiatan/:id {"is_published":false}` = 400. Fix: split DTO (pola berita yang sudah benar).
3. **Seeder pakai role name spasi** ("Super Admin"/"Admin") — bertentangan dengan normalisasi 00016 (snake_case). Seeder ulang di env baru → semua endpoint admin 403.
4. **DB drift**: kolom `phone`/`phone_verified_at` ada di DB (goose 18/19) tapi tidak ada di migration folder. Untuk project baru: **JANGAN buat kolom phone** (sudah diputuskan tidak dipakai).

### Migration yang direkomendasikan untuk project baru (sinkron dengan contract final):
- `pesan_kontak` perlu `deleted_at TIMESTAMP NULL` (model GORM pakai soft-delete, tapi migration 00008 tidak punya) — tambahkan seperti 00011 untuk sliders.

---

## Matriks Endpoint Final (79 endpoint, 100% sinkron)

| Modul | Jumlah | Status |
|-------|--------|--------|
| activity_logs | 4 | ✅ |
| auth | 10 | ✅ |
| berita | 11 | ✅ |
| dashboard | 1 | ✅ |
| kegiatan | 12 | ✅ |
| kontak | 7 | ✅ |
| pengurus | 9 | ✅ |
| settings | 3 | ✅ |
| sliders | 10 | ✅ |
| users | 12 | ✅ |
| **TOTAL** | **79** | **100%** |

---

*Dokumen ini adalah versi final contract yang sudah disinkronkan dengan implementasi BE. Cocok dijadikan acuan untuk membangun project baru serupa dari awal.*
