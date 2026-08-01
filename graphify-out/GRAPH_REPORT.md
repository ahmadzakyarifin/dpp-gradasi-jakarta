# Graph Report - dpp-gradasi-jakarta  (2026-08-01)

## Corpus Check
- 179 files · ~117,727 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1142 nodes · 2212 edges · 82 communities (63 shown, 19 thin omitted)
- Extraction: 81% EXTRACTED · 19% INFERRED · 0% AMBIGUOUS · INFERRED: 413 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `72d372ca`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- ErrorResponse
- App.jsx
- KegiatanRepo
- activityLogRepository
- AuthRepo
- BeritaRepo
- response.go
- PengurusService
- Rule
- RoleMiddleware
- KontakRepo
- devDependencies
- AuthMiddleware
- Config
- NewServiceError
- 2. PUT /api/v1/admin/settings
- rate_limiter_middleware.go
- UserService
- schema.sql
- dependencies
- Laporan Sinkronisasi 4 Lapisan — dpp_gradasi
- FITUR 3 — COMPANY PROFILE (Halaman Beranda Dinamis) — LAPORAN AUDIT (2026-08-01)
- .oxlintrc.json
- NewApp
- RateLimitPerUser
- 00003_create_auth_tokens_tables.sql
- berita
- 00006_create_kegiatan_tables.sql
- React + Vite
- ActivityLog
- users
- sliders
- activity_logs
- worker.go
- 00001_create_roles_table.sql
- 00007_create_pengurus_table.sql
- 00008_create_pesan_kontak_table.sql
- 00009_create_settings_table.sql
- kegiatan
- README.md
- activation_tokens
- users
- FITUR 2: MANAJEMEN ADMIN (USER)
- github.com/ahmadzakyarifin/dpp-gradasi/backend
- sliders
- main
- settings
- FITUR 1: FULL AUTH (Login + Forgot Password + Reset Password)
- CAPTCHA LOGIN FIX — ROOT CAUSE & SOLUSI (2026-08-01)
- 🟡 High
- FITUR 3 — COMPANY PROFILE DINAMIS (SELESAI 2026-08-01)
- LOG FIX — FITUR 2 (2026-08-01, approved Abang)
- UPDATE 2026-08-01 (verifikasi lanjutan Abang)
- DashboardService
- UPDATE 2 (2026-08-01) — SMTP real aktif + rate limiter bug fix besar
- UPDATE 3 (2026-08-01) — CAPTCHA aktif + countdown bersih + dead code dihapus
- SlidersRepo
- Mailer
- Errors
- context.go
- users
- ROOT CAUSE SESUNGGUHNYA CAPTCHA GAGAL — BUG STORE (2026-08-01, 17:11)
- ConnectDB
- .GetSummary
- incrFixedWindow
- RateLimiter
- RegisterRoutes
- RegisterRoutes
- RegisterRoutes

## God Nodes (most connected - your core abstractions)
1. `ErrorResponse()` - 84 edges
2. `SuccessResponse()` - 81 edges
3. `NewServiceError()` - 55 edges
4. `GetAuditMeta()` - 43 edges
5. `NewApp()` - 39 edges
6. `ValidationErrorResponse()` - 38 edges
7. `react` - 26 edges
8. `AuthRepo` - 24 edges
9. `Config` - 22 edges
10. `AuthService` - 21 edges

## Surprising Connections (you probably didn't know these)
- `main()` --calls--> `MustLoad()`  [INFERRED]
  backend/cmd/api/main.go → backend/config/config.go
- `main()` --calls--> `NewApp()`  [INFERRED]
  backend/cmd/api/main.go → backend/internal/app/app.go
- `main()` --calls--> `ConnectDB()`  [INFERRED]
  backend/cmd/api/main.go → backend/internal/infrastructure/database.go
- `main()` --calls--> `RegisterValidator()`  [INFERRED]
  backend/cmd/api/main.go → backend/internal/validator/validator.go
- `main()` --calls--> `MustLoad()`  [INFERRED]
  backend/cmd/seeder/main.go → backend/config/config.go

## Import Cycles
- None detected.

## Communities (82 total, 19 thin omitted)

### Community 0 - "ErrorResponse"
Cohesion: 0.07
Nodes (30): GetAuditMeta(), Context, ErrorResponse(), Context, SuccessResponse(), ValidationErrorResponse(), GetUserID(), Context (+22 more)

### Community 1 - "App.jsx"
Cohesion: 0.05
Nodes (56): apiRequest(), App(), ConfirmDialog(), ToastNotification(), AuthBrandPanel(), AuthShell(), CaptchaWidget(), authContent (+48 more)

### Community 2 - "KegiatanRepo"
Cohesion: 0.06
Nodes (28): BulkRequest, PaginationMeta, PaginationMeta, DeletedAt, Time, DB, maxInt(), NewKegiatanRepo() (+20 more)

### Community 3 - "activityLogRepository"
Cohesion: 0.07
Nodes (31): EntityToDetailResponse(), EntityToModel(), EntityToResponse(), ActivityLog, InputToEntity(), ModelToEntity(), DB, Time (+23 more)

### Community 4 - "AuthRepo"
Cohesion: 0.08
Nodes (9): Time, DB, NewAuthRepo(), Time, ActivationToken, PasswordResetToken, RefreshToken, Role (+1 more)

### Community 5 - "BeritaRepo"
Cohesion: 0.07
Nodes (22): BulkRequest, PaginationMeta, PaginationMeta, DeletedAt, Time, DB, NewBeritaRepo(), formatDate() (+14 more)

### Community 6 - "response.go"
Cohesion: 0.09
Nodes (20): GenerateValidationMessage(), GetPaginationMeta(), IsValidEmail(), IsValidURL(), RateLimitResponse(), Context, NewSettingsHandler(), DB (+12 more)

### Community 7 - "PengurusService"
Cohesion: 0.07
Nodes (20): BulkRequest, PaginationMeta, FileHeader, PaginationMeta, DeletedAt, Time, DB, NewPengurusRepo() (+12 more)

### Community 8 - "Rule"
Cohesion: 0.24
Nodes (14): normalizeRule(), normalizeScope(), abortRateLimitError(), abortRateLimitUnauthorized(), abortTooManyRequests(), Email(), Context, Duration (+6 more)

### Community 9 - "RoleMiddleware"
Cohesion: 0.13
Nodes (11): HandlerFunc, RoleMiddleware(), Client, RouterGroup, RegisterRoutes(), Client, RouterGroup, RegisterRoutes() (+3 more)

### Community 10 - "KontakRepo"
Cohesion: 0.08
Nodes (20): BulkRequest, PaginationMeta, PaginationMeta, DeletedAt, Time, DB, maxInt(), NewKontakRepo() (+12 more)

### Community 11 - "devDependencies"
Cohesion: 0.05
Nodes (39): esbuild-wasm, dependencies, react, react-dom, react-router-dom, zustand, devDependencies, autoprefixer (+31 more)

### Community 12 - "AuthMiddleware"
Cohesion: 0.16
Nodes (11): AuthMiddleware(), GetEmail(), GetRoleID(), Context, HandlerFunc, Client, RouterGroup, RegisterRoutes() (+3 more)

### Community 13 - "Config"
Cohesion: 0.11
Nodes (10): MustLoad(), AppConfig, Config, CookieConfig, DatabaseConfig, DevConfig, JWTConfig, RedisConfig (+2 more)

### Community 14 - "NewServiceError"
Cohesion: 0.06
Nodes (37): GenerateAccessToken(), GenerateRefreshToken(), RegisteredClaims, Time, ValidateToken(), CheckPassword(), HashPassword(), NewServiceError() (+29 more)

### Community 15 - "2. PUT /api/v1/admin/settings"
Cohesion: 0.11
Nodes (17): 1. GET /api/v1/settings, 2. PUT /api/v1/admin/settings, 3. POST /api/v1/admin/settings/logo, Alur Sinkronisasi FE-BE, API Contract — Dynamic Content Management (Site Settings), Catatan field, Request, Request (multipart/form-data) (+9 more)

### Community 16 - "rate_limiter_middleware.go"
Cohesion: 0.40
Nodes (8): emailFromBody(), getUserID(), Context, makeRateLimitKey(), needsEmail(), newRequestInfo(), rateLimitPath(), requestInfo

### Community 17 - "UserService"
Cohesion: 0.06
Nodes (25): BulkRequest, ChangePasswordRequest, FileHeader, DeletedAt, Time, DB, NewUserRepo(), generateRandomNumber() (+17 more)

### Community 18 - "schema.sql"
Cohesion: 0.23
Nodes (14): activation_tokens, berita, berita_tags, kegiatan, kegiatan_gallery, kegiatan_tags, password_reset_tokens, pengurus (+6 more)

### Community 19 - "dependencies"
Cohesion: 0.13
Nodes (14): dependencies, puppeteer, react-router-dom, zustand, devDependencies, autoprefixer, postcss, tailwindcss (+6 more)

### Community 20 - "Laporan Sinkronisasi 4 Lapisan — dpp_gradasi"
Cohesion: 0.15
Nodes (12): 10. Verifikasi, 1. users.status ENUM (SELESAI), 2. Settings Update rusak (SELESAI), 3. Activity Log (SELESAI), 4. Kegiatan (SELESAI), 5. Berita (SELESAI), 6. Pengurus (SELESAI), 7. Duplikasi change-password (SELESAI) (+4 more)

### Community 21 - "FITUR 3 — COMPANY PROFILE (Halaman Beranda Dinamis) — LAPORAN AUDIT (2026-08-01)"
Cohesion: 0.20
Nodes (10): FITUR 3 — COMPANY PROFILE (Halaman Beranda Dinamis) — LAPORAN AUDIT (2026-08-01), ✅ SUDAH BENAR, ✅ SUDAH DINAMIS (Home.jsx sudah fetch dari API), 🔴 TEMUAN 1 — HARDCODE di Home.jsx (harus dinamis), 🟡 TEMUAN 2 — Admin: 2 field settings tidak ada input di SettingsAdmin.jsx, 🟡 TEMUAN 3 — Kategori berita/kegiatan statis, 🔵 TEMUAN 4 — NetworkError rate limiter (bukan bug kode), 🔴🔴 TEMUAN 5 (KRITIS, DIFIX LANGSUNG 2026-08-01) — Berita & Kegiatan publik SELALU 500 (+2 more)

### Community 22 - ".oxlintrc.json"
Cohesion: 0.25
Nodes (7): plugins, rules, react/only-export-components, react/rules-of-hooks, $schema, oxc, warn

### Community 23 - "NewApp"
Cohesion: 0.31
Nodes (7): App, Client, DB, NewApp(), Client, registerRoutes(), Engine

### Community 24 - "RateLimitPerUser"
Cohesion: 0.20
Nodes (10): HandlerFunc, RateLimiterMiddleware(), RateLimitPerUser(), RateLimitRules(), Client, RouterGroup, RegisterRoutes(), Client (+2 more)

### Community 25 - "00003_create_auth_tokens_tables.sql"
Cohesion: 0.60
Nodes (4): activation_tokens, password_reset_tokens, refresh_tokens, users

### Community 26 - "berita"
Cohesion: 0.67
Nodes (3): berita, berita_tags, users

### Community 27 - "00006_create_kegiatan_tables.sql"
Cohesion: 0.83
Nodes (3): kegiatan, kegiatan_gallery, kegiatan_tags

### Community 28 - "React + Vite"
Cohesion: 0.50
Nodes (3): Expanding the Oxlint configuration, React Compiler, React + Vite

### Community 50 - "FITUR 2: MANAJEMEN ADMIN (USER)"
Cohesion: 0.29
Nodes (7): #1 Contract `users.jsonc` TIDAK VALID JSON — ✅ FIXED, #2 FE-BE list mismatch — halaman UsersAdmin TIDAK pernah tampil — ✅ FIXED, #3 `ProfileAdmin.jsx` hardcode `http://127.0.0.1:8080` — ✅ FIXED, 🔴 Critical, FITUR 2: MANAJEMEN ADMIN (USER), Ringkasan, ✅ Sudah benar (tidak diubah)

### Community 53 - "main"
Cohesion: 0.33
Nodes (4): main(), ConnectRedis(), Client, RegisterValidator()

### Community 57 - "FITUR 1: FULL AUTH (Login + Forgot Password + Reset Password)"
Cohesion: 0.06
Nodes (33): 1.1 Audit UI — React vs HTML statis (Login, Forgot, Reset), 1.2 Audit Validasi Form, 1.3 Audit Autentikasi & Session — Login, 1.4 Audit Forgot Password, 1.5 Audit Reset Password, 1.6 Audit CAPTCHA (Cloudflare Turnstile), 1.7 Audit Rate Limiter, 1.8 Audit Kode & Arsitektur (+25 more)

### Community 58 - "CAPTCHA LOGIN FIX — ROOT CAUSE & SOLUSI (2026-08-01)"
Cohesion: 0.29
Nodes (7): CAPTCHA LOGIN FIX — ROOT CAUSE & SOLUSI (2026-08-01), Catatan, Gejala, Perbaikan tambahan, Root cause (ditemukan via browser console), Solusi (industry standard — Cloudflare test keys), Verifikasi

### Community 59 - "🟡 High"
Cohesion: 0.33
Nodes (6): #4 Tab filter (pending/trash) tidak berfungsi — ✅ FIXED, #5 Bulk delete/restore tidak proteksi super_admin — ✅ FIXED, #6 Role name tidak ada di response — ✅ FIXED, #7 Upload foto profil tanpa validasi MIME/ukuran — ✅ FIXED, #8 Tidak ada rate limit di admin routes — ✅ FIXED, 🟡 High

### Community 60 - "FITUR 3 — COMPANY PROFILE DINAMIS (SELESAI 2026-08-01)"
Cohesion: 0.33
Nodes (6): Backend, Contract, FITUR 3 — COMPANY PROFILE DINAMIS (SELESAI 2026-08-01), Frontend, Migrasi, Verifikasi E2E

### Community 61 - "LOG FIX — FITUR 2 (2026-08-01, approved Abang)"
Cohesion: 0.33
Nodes (6): Bug pre-existing ditemukan & di-fix (di luar temuan audit), Catatan, E2E verification (2026-08-01), File diubah, LOG FIX — FITUR 2 (2026-08-01, approved Abang), Perubahan alur aktivasi akun (sesuai permintaan Abang)

### Community 62 - "UPDATE 2026-08-01 (verifikasi lanjutan Abang)"
Cohesion: 0.33
Nodes (6): ⚠️ SMTP Gmail masih `535 BadCredentials`, UPDATE 2026-08-01 (verifikasi lanjutan Abang), ✅ User management & aktivasi (dari E2E sebelumnya), ✅ Verifikasi alur forgot password (dengan SMTP mock lokal, bukti end-to-end), ✅ Verifikasi rate limiter login, ✅ Verifikasi reset inputan saat login gagal

### Community 63 - "DashboardService"
Cohesion: 0.48
Nodes (5): NewDashboardHandler(), DB, NewDashboardService(), DashboardHandler, DashboardService

### Community 64 - "UPDATE 2 (2026-08-01) — SMTP real aktif + rate limiter bug fix besar"
Cohesion: 0.40
Nodes (5): ✅ Alur aktivasi akun lengkap terverifikasi (SMTP real), 🔴🔴 BUG KRITIS DITEMUKAN & DIFIX: rate limiter TIDAK PERNAH blokir + retry_after selalu 1, ✅ Remember me terverifikasi (cookie max-age), ✅ SMTP Gmail REAL sudah aktif (app password baru Abang), UPDATE 2 (2026-08-01) — SMTP real aktif + rate limiter bug fix besar

### Community 65 - "UPDATE 3 (2026-08-01) — CAPTCHA aktif + countdown bersih + dead code dihapus"
Cohesion: 0.40
Nodes (5): 🧹 Bersih-bersih dead code (tanpa ubah logic), ✅ CAPTCHA Turnstile sekarang BENAR-BENAR berfungsi (bukan hiasan), ✅ Countdown rate limit tanpa duplikasi pesan + tombol disabled, ✅ Status auth & user VALID, UPDATE 3 (2026-08-01) — CAPTCHA aktif + countdown bersih + dead code dihapus

### Community 66 - "SlidersRepo"
Cohesion: 0.09
Nodes (14): BulkRequest, DeletedAt, Time, DB, NewSlidersRepo(), NewSlidersService(), toResponse(), ReorderRequest (+6 more)

### Community 67 - "Mailer"
Cohesion: 0.40
Nodes (3): NewMailer(), Dialer, Mailer

### Community 69 - "Errors"
Cohesion: 0.50
Nodes (3): validationMessage(), Errors(), ValidationErrorItem

### Community 74 - "ROOT CAUSE SESUNGGUHNYA CAPTCHA GAGAL — BUG STORE (2026-08-01, 17:11)"
Cohesion: 0.33
Nodes (6): Bug, Fix, Penemuan, Perbaikan pendukung (Login.jsx), ROOT CAUSE SESUNGGUHNYA CAPTCHA GAGAL — BUG STORE (2026-08-01, 17:11), Verifikasi E2E (browser nyata)

### Community 75 - "ConnectDB"
Cohesion: 0.40
Nodes (3): main(), ConnectDB(), DB

### Community 77 - "incrFixedWindow"
Cohesion: 0.27
Nodes (8): Context, Duration, Client, Context, Duration, incrFixedWindow(), toInt64(), fixedWindowResult

### Community 78 - "RateLimiter"
Cohesion: 0.36
Nodes (7): Client, newFixedWindowLimiter(), Client, NewRedisRateLimiter(), SetDefaultRateLimiter(), fixedWindowLimiter, RateLimiter

### Community 79 - "RegisterRoutes"
Cohesion: 0.50
Nodes (3): Client, RouterGroup, RegisterRoutes()

### Community 80 - "RegisterRoutes"
Cohesion: 0.50
Nodes (3): Client, RouterGroup, RegisterRoutes()

### Community 82 - "RegisterRoutes"
Cohesion: 0.50
Nodes (3): Client, RouterGroup, RegisterRoutes()

## Knowledge Gaps
- **171 isolated node(s):** `WorkerConfig`, `github.com/ahmadzakyarifin/dpp-gradasi/backend`, `contextKey`, `Response`, `TurnstileVerifyResult` (+166 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **19 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `NewApp()` connect `NewApp` to `ErrorResponse`, `KegiatanRepo`, `activityLogRepository`, `Mailer`, `AuthRepo`, `BeritaRepo`, `PengurusService`, `response.go`, `SlidersRepo`, `KontakRepo`, `Config`, `RateLimiter`, `NewServiceError`, `UserService`, `main`, `DashboardService`?**
  _High betweenness centrality (0.166) - this node is a cross-community bridge._
- **Why does `NewServiceError()` connect `NewServiceError` to `ErrorResponse`, `KegiatanRepo`, `SlidersRepo`, `BeritaRepo`, `PengurusService`, `KontakRepo`, `UserService`?**
  _High betweenness centrality (0.135) - this node is a cross-community bridge._
- **Why does `ErrorResponse()` connect `ErrorResponse` to `response.go`, `Rule`, `RoleMiddleware`, `AuthMiddleware`, `NewServiceError`?**
  _High betweenness centrality (0.055) - this node is a cross-community bridge._
- **Are the 82 inferred relationships involving `ErrorResponse()` (e.g. with `AuthMiddleware()` and `abortRateLimitError()`) actually correct?**
  _`ErrorResponse()` has 82 INFERRED edges - model-reasoned connections that need verification._
- **Are the 79 inferred relationships involving `SuccessResponse()` (e.g. with `.Detail()` and `.EntityLogs()`) actually correct?**
  _`SuccessResponse()` has 79 INFERRED edges - model-reasoned connections that need verification._
- **Are the 53 inferred relationships involving `NewServiceError()` (e.g. with `.ActivateAccount()` and `.ChangePassword()`) actually correct?**
  _`NewServiceError()` has 53 INFERRED edges - model-reasoned connections that need verification._
- **Are the 41 inferred relationships involving `GetAuditMeta()` (e.g. with `.ChangePassword()` and `.Logout()`) actually correct?**
  _`GetAuditMeta()` has 41 INFERRED edges - model-reasoned connections that need verification._