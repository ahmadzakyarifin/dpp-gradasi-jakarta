# Graph Report - dpp-gradasi-jakarta  (2026-08-02)

## Corpus Check
- 180 files · ~115,903 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1070 nodes · 2168 edges · 70 communities (51 shown, 19 thin omitted)
- Extraction: 81% EXTRACTED · 19% INFERRED · 0% AMBIGUOUS · INFERRED: 413 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `bc7cbc91`
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
- User
- .oxlintrc.json
- Perubahan yang Dilakukan (Finalisasi Contract)
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
- user_dto.go
- github.com/ahmadzakyarifin/dpp-gradasi/backend
- sliders
- UserRepo
- settings
- pesan_kontak
- DashboardService
- SlidersRepo
- Mailer
- context.go
- users
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
6. `ValidationErrorResponse()` - 37 edges
7. `react` - 28 edges
8. `AuthRepo` - 24 edges
9. `KegiatanRepo` - 23 edges
10. `Config` - 22 edges

## Surprising Connections (you probably didn't know these)
- `main()` --calls--> `NewApp()`  [INFERRED]
  backend/cmd/api/main.go → backend/internal/app/app.go
- `main()` --calls--> `RegisterValidator()`  [INFERRED]
  backend/cmd/api/main.go → backend/internal/validator/validator.go
- `NewApp()` --calls--> `NewMailer()`  [INFERRED]
  backend/internal/app/app.go → backend/internal/infrastructure/mail.go
- `NewApp()` --calls--> `NewRedisRateLimiter()`  [INFERRED]
  backend/internal/app/app.go → backend/internal/middleware/rate_limiter_middleware.go
- `NewApp()` --calls--> `SetDefaultRateLimiter()`  [INFERRED]
  backend/internal/app/app.go → backend/internal/middleware/rate_limiter_middleware.go

## Import Cycles
- None detected.

## Communities (70 total, 19 thin omitted)

### Community 0 - "ErrorResponse"
Cohesion: 0.06
Nodes (40): App, Client, DB, NewApp(), Client, registerRoutes(), GetAuditMeta(), Context (+32 more)

### Community 1 - "App.jsx"
Cohesion: 0.05
Nodes (55): apiRequest(), App(), ConfirmDialog(), ToastNotification(), AuthBrandPanel(), AuthShell(), CaptchaWidget(), NavbarBeritaSearch() (+47 more)

### Community 2 - "KegiatanRepo"
Cohesion: 0.05
Nodes (29): BulkRequest, PaginationMeta, PaginationMeta, DeletedAt, Time, DB, maxInt(), NewKegiatanRepo() (+21 more)

### Community 3 - "activityLogRepository"
Cohesion: 0.07
Nodes (31): EntityToDetailResponse(), EntityToModel(), EntityToResponse(), ActivityLog, InputToEntity(), ModelToEntity(), DB, Time (+23 more)

### Community 4 - "AuthRepo"
Cohesion: 0.11
Nodes (7): Time, DB, NewAuthRepo(), ActivationToken, PasswordResetToken, RefreshToken, AuthRepo

### Community 5 - "BeritaRepo"
Cohesion: 0.06
Nodes (23): BulkRequest, PaginationMeta, PaginationMeta, DeletedAt, Time, DB, NewBeritaRepo(), formatDate() (+15 more)

### Community 6 - "response.go"
Cohesion: 0.10
Nodes (17): GenerateValidationMessage(), GetPaginationMeta(), IsValidEmail(), IsValidURL(), RateLimitResponse(), DB, NewSettingsRepo(), FileHeader (+9 more)

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
Cohesion: 0.07
Nodes (20): main(), main(), MustLoad(), ConnectDB(), DB, ConnectRedis(), Client, validationMessage() (+12 more)

### Community 14 - "NewServiceError"
Cohesion: 0.05
Nodes (38): GenerateAccessToken(), GenerateRefreshToken(), RegisteredClaims, Time, ValidateToken(), CheckPassword(), HashPassword(), NewServiceError() (+30 more)

### Community 15 - "2. PUT /api/v1/admin/settings"
Cohesion: 0.11
Nodes (17): 1. GET /api/v1/settings, 2. PUT /api/v1/admin/settings, 3. POST /api/v1/admin/settings/logo, Alur Sinkronisasi FE-BE, API Contract — Dynamic Content Management (Site Settings), Catatan field, Request, Request (multipart/form-data) (+9 more)

### Community 16 - "rate_limiter_middleware.go"
Cohesion: 0.40
Nodes (8): emailFromBody(), getUserID(), Context, makeRateLimitKey(), needsEmail(), newRequestInfo(), rateLimitPath(), requestInfo

### Community 17 - "UserService"
Cohesion: 0.17
Nodes (9): generateRandomNumber(), generateRandomPassword(), ChangePasswordRequest, Client, NewUserService(), toResponse(), AdminCreateRequest, UserResponse (+1 more)

### Community 18 - "schema.sql"
Cohesion: 0.23
Nodes (14): activation_tokens, berita, berita_tags, kegiatan, kegiatan_gallery, kegiatan_tags, password_reset_tokens, pengurus (+6 more)

### Community 19 - "dependencies"
Cohesion: 0.13
Nodes (14): dependencies, puppeteer, react-router-dom, zustand, devDependencies, autoprefixer, postcss, tailwindcss (+6 more)

### Community 20 - "Laporan Sinkronisasi 4 Lapisan — dpp_gradasi"
Cohesion: 0.15
Nodes (12): 10. Verifikasi, 1. users.status ENUM (SELESAI), 2. Settings Update rusak (SELESAI), 3. Activity Log (SELESAI), 4. Kegiatan (SELESAI), 5. Berita (SELESAI), 6. Pengurus (SELESAI), 7. Duplikasi change-password (SELESAI) (+4 more)

### Community 21 - "User"
Cohesion: 0.14
Nodes (5): Time, DeletedAt, Time, Role, User

### Community 22 - ".oxlintrc.json"
Cohesion: 0.25
Nodes (7): plugins, rules, react/only-export-components, react/rules-of-hooks, $schema, oxc, warn

### Community 23 - "Perubahan yang Dilakukan (Finalisasi Contract)"
Cohesion: 0.13
Nodes (14): 1. `settings.jsonc` — method + endpoint baru, 2. `users.jsonc` — role_id + endpoint baru, 3. `auth.jsonc` — hapus phone + tambah must_change_password, 4. `berita.jsonc` — query param, 5. `kegiatan.jsonc` — query param, 6. `pengurus.jsonc` — query param, 7. `activity_logs.jsonc` — response shape disinkronkan ke BE, Catatan Penting untuk Project Baru (+6 more)

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

### Community 50 - "user_dto.go"
Cohesion: 0.16
Nodes (9): BulkRequest, ChangePasswordRequest, FileHeader, AdminStatusRequest, ListUsersQuery, Pagination, ProfileUpdateRequest, UserListResponse (+1 more)

### Community 53 - "UserRepo"
Cohesion: 0.22
Nodes (4): DB, NewUserRepo(), FileHeader, UserRepo

### Community 63 - "DashboardService"
Cohesion: 0.25
Nodes (7): NewDashboardHandler(), Context, DB, NewDashboardService(), DashboardSummaryRes, DashboardHandler, DashboardService

### Community 66 - "SlidersRepo"
Cohesion: 0.09
Nodes (14): BulkRequest, DeletedAt, Time, DB, NewSlidersRepo(), NewSlidersService(), toResponse(), ReorderRequest (+6 more)

### Community 67 - "Mailer"
Cohesion: 0.40
Nodes (3): NewMailer(), Dialer, Mailer

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
- **102 isolated node(s):** `WorkerConfig`, `github.com/ahmadzakyarifin/dpp-gradasi/backend`, `contextKey`, `Response`, `TurnstileVerifyResult` (+97 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **19 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `NewApp()` connect `ErrorResponse` to `KegiatanRepo`, `activityLogRepository`, `Mailer`, `AuthRepo`, `BeritaRepo`, `PengurusService`, `response.go`, `SlidersRepo`, `KontakRepo`, `Config`, `RateLimiter`, `NewServiceError`, `UserService`, `UserRepo`, `DashboardService`?**
  _High betweenness centrality (0.188) - this node is a cross-community bridge._
- **Why does `NewServiceError()` connect `NewServiceError` to `ErrorResponse`, `KegiatanRepo`, `SlidersRepo`, `BeritaRepo`, `PengurusService`, `KontakRepo`, `UserService`, `user_dto.go`?**
  _High betweenness centrality (0.178) - this node is a cross-community bridge._
- **Why does `ErrorResponse()` connect `ErrorResponse` to `response.go`, `Rule`, `RoleMiddleware`, `AuthMiddleware`, `NewServiceError`?**
  _High betweenness centrality (0.064) - this node is a cross-community bridge._
- **Are the 82 inferred relationships involving `ErrorResponse()` (e.g. with `AuthMiddleware()` and `abortRateLimitError()`) actually correct?**
  _`ErrorResponse()` has 82 INFERRED edges - model-reasoned connections that need verification._
- **Are the 79 inferred relationships involving `SuccessResponse()` (e.g. with `.Detail()` and `.EntityLogs()`) actually correct?**
  _`SuccessResponse()` has 79 INFERRED edges - model-reasoned connections that need verification._
- **Are the 53 inferred relationships involving `NewServiceError()` (e.g. with `.ActivateAccount()` and `.ChangePassword()`) actually correct?**
  _`NewServiceError()` has 53 INFERRED edges - model-reasoned connections that need verification._
- **Are the 41 inferred relationships involving `GetAuditMeta()` (e.g. with `.ChangePassword()` and `.Logout()`) actually correct?**
  _`GetAuditMeta()` has 41 INFERRED edges - model-reasoned connections that need verification._