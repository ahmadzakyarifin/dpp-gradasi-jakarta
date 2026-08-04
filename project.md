# Project DPP Gradasi Jakarta — Dokumentasi Pembelajaran

> Dibuat untuk dipelajari ulang. Sumber: analisis menyeluruh codebase (backend Go + frontend React + kontrak API + migrasi DB), audit FE↔BE, graphify (994 node / 47 komunitas), dan verifikasi live via curl + mysql.
> Terakhir diperbarui: 2026-08-03 (setelah sinkronisasi FE dengan backend baru teman "zaky").

---

## 1. Gambaran Umum

Website organisasi **DPP GRADASI** (DPD/DPC di 27 pengurus) — portal berita, kegiatan, pengurus, sliders, kontak, plus panel admin dengan role & activity log.

| Bagian | Teknologi | Lokasi |
|---|---|---|
| Backend | Go (Gin + GORM), MySQL 8, goose migrasi, JWT | `backend/` |
| Frontend | React 19, Vite (rolldown), Tailwind, Zustand | `frontend/` |
| Kontrak | JSONC per-module (source of truth FE↔BE) | `docs-final/api/*.jsonc` |
| Seed data | Go program `cmd/seeder` | `backend/cmd/seeder/` |

**Pola arsitektur backend:** modular monolith — tiap modul punya `routes.go` + `controller/` (handler) + `service/` + `repository/` + `model/`. 11 modul: `auth, user, role, berita, kegiatan, pengurus, sliders, kontak, settings, dashboard, activitylog`.

**Pola frontend:** pages (admin + publik) → `services/*Service.js` (satu fungsi per endpoint) → `api/index.js` (axios + interceptor 401-refresh) → backend. State global Zustand di `store/`.

---

## 2. Cara Menjalankan

### Backend (Go)
```bash
cd backend
# butuh .env (lihat catatan) — jangan di-commit
go run ./cmd/api        # atau: air (hot reload), port 8080, prefix /api/v1
```
Migrasi TIDAK jalan otomatis — jalankan manual:
```bash
goose -dir backend/migrations mysql "root:PASSWORD@tcp(127.0.0.1:3306)/dpp_gradasi?parseTime=true" up
```
Seed data (HATI-HATI: truncate DB dulu):
```bash
go run ./cmd/seeder    # make seed — menghapus isi tabel, bukan cuma insert
```
**Catatan goose (pitfall nyata):** MySQL 8 menolak >1 statement per query (Error 1064) → setiap SQL statement WAJIB punya blok `-- +goose StatementBegin/End` sendiri. `ADD COLUMN IF NOT EXISTS` juga TIDAK didukung MySQL 8 (hanya MariaDB) → gunakan cek `INFORMATION_SCHEMA.COLUMNS`.

### Frontend (React)
```bash
cd frontend
npm install
npm run dev      # dev server
npm run build    # build production
npm run lint     # ESLint
```
`.env` punya `VITE_API_URL` (default `http://127.0.0.1:8080`). Captcha toggle di `.env` (ada `.env.bak-captcha`).

### Login super admin (dev)
`admin@gradasi.org` / `password123` (role super_admin id 9)

---

## 3. Struktur & Alur (hasil graphify)

### Alur request backend (pola seragam di semua modul)
```
route (routes.go) → middleware (auth/role) → handler (controller) → service → repository → GORM → MySQL
                                                      ↓
                                response helper: SuccessResponse / ErrorResponse / HandleServiceError
```
Semua response envelope: `{ success, code, message, data }`. Error via `ServiceError` (NewServiceError → HandleServiceError). Ini God Nodes graph: `HandleServiceError` (99 edges), `SuccessResponse` (98), `ErrorResponse` (81), `NewApp` (42 — wiring semua handler).

### Alur login (spesifik, dari graph query)
```
POST /auth/login → LoginRequest (validate) → auth service .Login()
  → cek user by email → bcrypt compare → generate access JWT + refresh token
  → simpan refresh token (tabel refresh_tokens) → activity log
  → response: user + access_token + refresh_token
```

### Alur frontend
```
page → service (normalizeImage) → api/index.js (axios + interceptors)
  → interceptor 401: coba refresh token → kalau gagal logout
  → response { data: ... } → page set state → render
```

### Migrasi & schema (14 migrasi)
```
00001 roles → 00002 users → 00003 auth tokens (refresh_tokens, activation_tokens)
00004 sliders → 00005 berita → 00006 kegiatan → 00007 pengurus → 00008 pesan_kontak
00009 settings → 00010 activity_logs
00011 seed roles (super_admin/admin/admin_berita/editor)
00012 roles +is_system +deleted_at
00013 refresh_tokens: device_info→device_name, +user_agent, +revoked_at
00014 users +must_change_password
```
**Role id (setelah reseed):** 9=super_admin, 10=admin, 11=editor, 12=admin_berita. (Sebelumnya 1-4; seed tidak truncate → AUTO_INCREMENT lanjut → id 9-12.)

---

## 4. Kontrak API (docs-final — source of truth)

Tiap file `docs-final/api/*.jsonc` = spesifikasi 1 modul: endpoint, method, request, response, validasi per-field (422), error codes. **FE harus membaca ini sebelum menyentuh service.**

Poin penting (dari audit):
- **Response list**: `data.items` + `data.meta.{total, page, limit, total_pages}` (users, kegiatan, berita — admin & publik). **BUKAN** `data.pagination` (key lama yang sudah hilang).
- **Gambar**: response selalu `image_path` (path relatif `/uploads/...`); FE pakai `normalizeImage()` → tambah `image_url` absolut. Request create/update WAJIB `image_path` (bukan `image_url`).
- **Pengurus**: multipart/form-data; create wajib `image` (File); kalau edit tanpa ganti foto, kirim `image_path` string lama.
- **Roles**: `GET /api/v1/roles` → `data[]` `{id, name, display_name, is_system, is_active, user_count}`. Role super_admin punya `is_system: true`.
- **Validasi 422 per-field**: `data.errors[]` `{field, tag, message}` — FE harus tampilkan per-field.
- **Kontak**: rate limit 5/menit per IP, header `X-RateLimit-*`, code `KONTAK_RATE_LIMITED`.
- **Auth**: login `AUTH_LOGIN_SUCCESS`; 401 dengan body `{code: "TOKEN_EXPIRED"|...}` → FE refresh.
- **Activity logs**: `data.summary` + `data.items` + `data.meta`, snake_case.

---

## 5. Gap FE↔BE yang Pernah Ada & Fix-nya (pembelajaran!)

Ini inti pelajaran dari sesi sinkronisasi 2026-08-03:

| Gap | Akar masalah | Fix |
|---|---|---|
| Login 500 | Migrasi 11-14 belum di-apply → tabel roles/refresh_tokens kurang kolom; kode baru query kolom yang tak ada | `goose up` (schema sinkron) — bukan ubah kode |
| 401 semua auth | Middleware set `contextKey` (type) tapi handler baca `c.Get("user_id")` literal → key beda | Pakai helper `middleware.GetUserID(c)` konsisten (commit `abbdf67`) |
| Role dropdown hardcode 2/5/6 | Role id berubah jadi 9-12; kontrak lama | Fetch `GET /roles` dinamis + filter super_admin |
| Pagination users "selalu 1 halaman" | FE baca `data.pagination.*`; backend kirim `data.meta.*` | Baca `data.meta.{total,total_pages}` |
| Gambar tak tersimpan | FE kirim `image_url`, kontrak `image_path` | Rename key di payload (berita/kegiatan/sliders/pengurus) |
| Edit kegiatan kehilangan content | list_admin tidak sertakan content → form isi = excerpt | `openForm` fetch `detailById(id)` dulu |
| Gambar publik 404 | Hardcode `http://127.0.0.1:8080` di img src | Helper `resolveAssetUrl()` |
| Berita/kegiatan tak muncul | Semua data ter-soft-delete (deleted_at) | `POST /berita/bulk-restore` + `/kegiatan/bulk-restore` |

**Pelajaran utama:** (1) migrasi tidak otomatis → selalu cek `goose status` setelah pull; (2) kontrak `docs-final` adalah sumber kebenaran, baca sebelum coding; (3) normalize/asset resolution di service layer, bukan di page; (4) root cause sering di DB schema drift, bukan logika kode.

---

## 6. Modul Backend (detail singkat per modul)

| Modul | Endpoint utama | Catatan |
|---|---|---|
| auth | login, register?, me, refresh, activate, verify-email | JWT access + refresh; role dari users.role_id |
| user | profile CRUD, admin users (list/create/update_status/resend_activation) | `column:password` (bukan password_hash) — fix login |
| role | CRUD roles | super_admin is_system |
| berita | list publik, detail slug, admin CRUD, bulk-restore, toggle publish | FULLTEXT search; is_published; soft delete |
| kegiatan | list publik, detail, admin CRUD, gallery, bulk | event_date, location, organizer, tags/gallery string JSON |
| pengurus | list publik (filter level/provinsi/kabupaten), regions, admin CRUD multipart | image wajib saat create |
| sliders | list, admin CRUD, sort_order | image_path, link_url |
| kontak | submit (rate limit 5/menit), list admin, read status | |
| settings | site settings | image_path |
| dashboard | summary | super_admin/admin |
| activitylog | list + summary | snake_case; auto-log tiap aksi |

---

## 7. Pola & Konvensi yang Perlu Diingat

### Backend
- Envelope: `{success, code, message, data}` — code ALL_CAPS (mis. `BERITA_CREATED`).
- Error: `ServiceError` + `HandleServiceError` di controller.
- Validasi: per-field → 422 `data.errors[]`.
- GORM soft delete: `deleted_at`; query publik filter `deleted_at IS NULL` + `is_published`.
- Handler di `controller/`, logika di `service/`, DB di `repository/` — jangan campur.
- `NewApp()` di `cmd/api/main.go` me-wire semua handler (42 edges di graph).

### Frontend
- `src/services/*.js` — satu file per modul, fungsi per endpoint, `toQuery()` untuk params, `.then` normalize image.
- `normalizeImage` — response `image_path` → `image_url` (dipakai service layer).
- `resolveAssetUrl` — path relatif `/uploads/...` → URL absolut (dipakai page display).
- Interceptor 401-refresh di `api/index.js` — jangan bypass, kecuali endpoint publik.
- Naming: response snake_case (`is_published`, `event_date`, `image_path`) → form state camelCase (`isPublished`, `eventDate`) → payload snake_case lagi.

### Migrasi (goose)
- 1 statement per `StatementBegin/End` block (MySQL 8).
- Jangan pakai `ADD COLUMN IF NOT EXISTS` (MySQL 8).
- Seed via `00011_seed_roles.sql`; seeder Go `make seed` truncate (jangan jalankan di prod!).
- Setelah reseed roles, user id lama harus di-`UPDATE users SET role_id = <id baru>`.

---

## 8. State Terkini (2026-08-03)

- Backend: goose version 14, 0 pending; login super_admin OK; semua endpoint inti 200.
- Frontend: build ✅, lint ✅ (14 warning pre-existing, 0 error); semua gap FE↔BE dari tabel §5 sudah di-fix.
- Graph: `backend/graphify-out/` (graph.json 1.3MB, graph.html, GRAPH_REPORT.md) — query via `graphify query "..." --graph backend/graphify-out/graph.json`.
- **Masih terbuka:** UI gallery delete (`DELETE /kegiatan/gallery/:gallery_id`), UI verify-email (`POST /profile/verify-email`), input tags/gallery di form kegiatan.

---

## 9. Quick Reference

```bash
# backend
cd backend && go run ./cmd/api
goose -dir backend/migrations mysql "root:root@tcp(127.0.0.1:3306)/dpp_gradasi?parseTime=true" status
# frontend
cd frontend && npm run dev
# graph
graphify query "alur login" --graph backend/graphify-out/graph.json
```

**Kredensial dev:** admin@gradasi.org / password123 · MySQL root/root · DB `dpp_gradasi`
