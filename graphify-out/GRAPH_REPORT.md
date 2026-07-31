# Graph Report - dpp_gradasi  (2026-07-31)

## Corpus Check
- 142 files · ~81,900 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 919 nodes · 1801 edges · 68 communities (47 shown, 21 thin omitted)
- Extraction: 77% EXTRACTED · 23% INFERRED · 0% AMBIGUOUS · INFERRED: 407 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `8b7b7015`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- ErrorResponse
- KegiatanRepo
- BeritaRepo
- NewServiceError
- auth_dto.go
- PengurusService
- rate_limiter_middleware.go
- activityLogRepository
- devDependencies
- KontakRepo
- SlidersRepo
- response.go
- AuthRepo
- Config
- Login.jsx
- ActivityLogService
- NewApp
- main
- ActivityLog
- Mailer
- dependencies
- ActivityLog
- scratch.js
- worker.go
- schema.sql
- AuthMiddleware
- README.md
- test-ui.js
- github.com/ahmadzakyarifin/dpp-gradasi/backend
- AuthService
- Errors
- 00003_create_auth_tokens_tables.sql
- Laporan Sinkronisasi 4 Lapisan — dpp_gradasi
- user_dto.go
- berita
- 00006_create_kegiatan_tables.sql
- users
- sliders
- activity_logs
- 00001_create_roles_table.sql
- 00007_create_pengurus_table.sql
- 00008_create_pesan_kontak_table.sql
- 00009_create_settings_table.sql
- pesan_kontak
- sliders
- GenerateResetToken
- .ChangePassword
- ServiceError
- .GetSummary
- React + Vite
- NewUserHandler
- kegiatan
- activation_tokens
- users

## God Nodes (most connected - your core abstractions)
1. `ErrorResponse()` - 81 edges
2. `SuccessResponse()` - 78 edges
3. `NewServiceError()` - 55 edges
4. `GetAuditMeta()` - 43 edges
5. `NewApp()` - 39 edges
6. `ValidationErrorResponse()` - 36 edges
7. `AuthRepo` - 24 edges
8. `Config` - 22 edges
9. `UserService` - 21 edges
10. `ActivityLogService` - 20 edges

## Surprising Connections (you probably didn't know these)
- `main()` --calls--> `MustLoad()`  [INFERRED]
  backend/cmd/api/main.go → backend/config/config.go
- `main()` --calls--> `NewApp()`  [INFERRED]
  backend/cmd/api/main.go → backend/internal/app/app.go
- `main()` --calls--> `RegisterValidator()`  [INFERRED]
  backend/cmd/api/main.go → backend/internal/validator/validator.go
- `main()` --calls--> `MustLoad()`  [INFERRED]
  backend/cmd/seeder/main.go → backend/config/config.go
- `NewApp()` --calls--> `NewMailer()`  [INFERRED]
  backend/internal/app/app.go → backend/internal/infrastructure/mail.go

## Import Cycles
- None detected.

## Communities (68 total, 21 thin omitted)

### Community 0 - "ErrorResponse"
Cohesion: 0.07
Nodes (26): GetAuditMeta(), Context, ErrorResponse(), Context, SuccessResponse(), ValidationErrorResponse(), GetUserID(), Context (+18 more)

### Community 1 - "KegiatanRepo"
Cohesion: 0.06
Nodes (28): BulkRequest, PaginationMeta, PaginationMeta, DeletedAt, Time, DB, maxInt(), NewKegiatanRepo() (+20 more)

### Community 2 - "BeritaRepo"
Cohesion: 0.07
Nodes (22): BulkRequest, PaginationMeta, PaginationMeta, DeletedAt, Time, DB, NewBeritaRepo(), formatDate() (+14 more)

### Community 3 - "NewServiceError"
Cohesion: 0.20
Nodes (8): NewServiceError(), generateRandomNumber(), ChangePasswordRequest, Client, NewUserService(), toResponse(), UserResponse, UserService

### Community 4 - "auth_dto.go"
Cohesion: 0.15
Nodes (14): GenerateAccessToken(), GenerateRefreshToken(), RegisteredClaims, Time, ValidateToken(), ChangePasswordRequest, AuthResponse, AuthUserResponse (+6 more)

### Community 5 - "PengurusService"
Cohesion: 0.07
Nodes (20): BulkRequest, PaginationMeta, FileHeader, PaginationMeta, DeletedAt, Time, DB, NewPengurusRepo() (+12 more)

### Community 6 - "rate_limiter_middleware.go"
Cohesion: 0.09
Nodes (38): Context, RestoreBody(), emailFromBody(), getUserID(), Client, Context, HandlerFunc, makeRateLimitKey() (+30 more)

### Community 7 - "activityLogRepository"
Cohesion: 0.09
Nodes (26): EntityToDetailResponse(), EntityToModel(), EntityToResponse(), ActivityLog, InputToEntity(), ModelToEntity(), firstNonEmpty(), ActivityLog (+18 more)

### Community 8 - "devDependencies"
Cohesion: 0.05
Nodes (39): esbuild-wasm, dependencies, react, react-dom, react-router-dom, zustand, devDependencies, autoprefixer (+31 more)

### Community 9 - "KontakRepo"
Cohesion: 0.08
Nodes (20): BulkRequest, PaginationMeta, PaginationMeta, DeletedAt, Time, DB, maxInt(), NewKontakRepo() (+12 more)

### Community 10 - "SlidersRepo"
Cohesion: 0.09
Nodes (14): BulkRequest, DeletedAt, Time, DB, NewSlidersRepo(), NewSlidersService(), toResponse(), ReorderRequest (+6 more)

### Community 11 - "response.go"
Cohesion: 0.09
Nodes (21): GenerateValidationMessage(), GetPaginationMeta(), IsValidEmail(), IsValidURL(), RateLimitResponse(), NewSettingsHandler(), DB, NewSettingsRepo() (+13 more)

### Community 12 - "AuthRepo"
Cohesion: 0.05
Nodes (16): Time, DB, NewAuthRepo(), Time, DeletedAt, Time, DB, NewUserRepo() (+8 more)

### Community 13 - "Config"
Cohesion: 0.11
Nodes (10): MustLoad(), AppConfig, Config, CookieConfig, DatabaseConfig, DevConfig, JWTConfig, RedisConfig (+2 more)

### Community 14 - "Login.jsx"
Cohesion: 0.12
Nodes (21): plugins, rules, react/only-export-components, react/rules-of-hooks, $schema, apiRequest(), App(), AuthBrandPanel() (+13 more)

### Community 15 - "ActivityLogService"
Cohesion: 0.18
Nodes (9): NewActivityLogHandler(), ActivityLogService, NewAuthHandler(), NewDashboardHandler(), DB, NewDashboardService(), NewPengurusHandler(), DashboardHandler (+1 more)

### Community 16 - "NewApp"
Cohesion: 0.16
Nodes (10): App, Client, DB, NewApp(), Client, registerRoutes(), NewBeritaHandler(), NewKontakHandler() (+2 more)

### Community 17 - "main"
Cohesion: 0.18
Nodes (7): main(), main(), ConnectDB(), DB, ConnectRedis(), Client, RegisterValidator()

### Community 18 - "ActivityLog"
Cohesion: 0.22
Nodes (5): DB, Time, ActivityLog, JSONMap, Value

### Community 19 - "Mailer"
Cohesion: 0.38
Nodes (3): NewMailer(), Dialer, Mailer

### Community 20 - "dependencies"
Cohesion: 0.13
Nodes (14): dependencies, puppeteer, react-router-dom, zustand, devDependencies, autoprefixer, postcss, tailwindcss (+6 more)

### Community 24 - "schema.sql"
Cohesion: 0.23
Nodes (14): activation_tokens, berita, berita_tags, kegiatan, kegiatan_gallery, kegiatan_tags, password_reset_tokens, pengurus (+6 more)

### Community 25 - "AuthMiddleware"
Cohesion: 0.06
Nodes (31): AuthMiddleware(), GetEmail(), GetRoleID(), Context, HandlerFunc, HandlerFunc, RoleMiddleware(), Client (+23 more)

### Community 30 - "AuthService"
Cohesion: 0.19
Nodes (7): HashToken(), Client, NewAuthService(), ActivateAccountRequest, ForgotPasswordRequest, ResetPasswordRequest, AuthService

### Community 31 - "Errors"
Cohesion: 0.50
Nodes (3): validationMessage(), Errors(), ValidationErrorItem

### Community 32 - "00003_create_auth_tokens_tables.sql"
Cohesion: 0.60
Nodes (4): activation_tokens, password_reset_tokens, refresh_tokens, users

### Community 33 - "Laporan Sinkronisasi 4 Lapisan — dpp_gradasi"
Cohesion: 0.15
Nodes (12): 10. Verifikasi, 1. users.status ENUM (SELESAI), 2. Settings Update rusak (SELESAI), 3. Activity Log (SELESAI), 4. Kegiatan (SELESAI), 5. Berita (SELESAI), 6. Pengurus (SELESAI), 7. Duplikasi change-password (SELESAI) (+4 more)

### Community 34 - "user_dto.go"
Cohesion: 0.20
Nodes (7): BulkRequest, ChangePasswordRequest, FileHeader, AdminCreateRequest, AdminStatusRequest, ProfileUpdateRequest, VerifyEmailRequest

### Community 35 - "berita"
Cohesion: 0.67
Nodes (3): berita, berita_tags, users

### Community 36 - "00006_create_kegiatan_tables.sql"
Cohesion: 0.83
Nodes (3): kegiatan, kegiatan_gallery, kegiatan_tags

### Community 47 - "GenerateResetToken"
Cohesion: 0.53
Nodes (5): GenerateActivationToken(), GenerateResetToken(), RegisteredClaims, Time, ParseExpiresAt()

### Community 48 - ".ChangePassword"
Cohesion: 0.50
Nodes (3): CheckPassword(), HashPassword(), ChangePasswordRequest

### Community 51 - "React + Vite"
Cohesion: 0.50
Nodes (3): Expanding the Oxlint configuration, React Compiler, React + Vite

## Knowledge Gaps
- **74 isolated node(s):** `WorkerConfig`, `github.com/ahmadzakyarifin/dpp-gradasi/backend`, `Response`, `RegisterRequest`, `ChangePasswordRequest` (+69 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **21 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `NewApp()` connect `NewApp` to `ErrorResponse`, `KegiatanRepo`, `BeritaRepo`, `NewServiceError`, `PengurusService`, `rate_limiter_middleware.go`, `activityLogRepository`, `KontakRepo`, `SlidersRepo`, `response.go`, `AuthRepo`, `Config`, `ActivityLogService`, `main`, `Mailer`, `NewUserHandler`, `AuthService`?**
  _High betweenness centrality (0.231) - this node is a cross-community bridge._
- **Why does `NewServiceError()` connect `NewServiceError` to `ErrorResponse`, `KegiatanRepo`, `BeritaRepo`, `user_dto.go`, `auth_dto.go`, `PengurusService`, `KontakRepo`, `SlidersRepo`, `.ChangePassword`, `ServiceError`, `AuthService`?**
  _High betweenness centrality (0.180) - this node is a cross-community bridge._
- **Why does `ErrorResponse()` connect `ErrorResponse` to `AuthMiddleware`, `response.go`, `rate_limiter_middleware.go`?**
  _High betweenness centrality (0.068) - this node is a cross-community bridge._
- **Are the 79 inferred relationships involving `ErrorResponse()` (e.g. with `AuthMiddleware()` and `abortRateLimitError()`) actually correct?**
  _`ErrorResponse()` has 79 INFERRED edges - model-reasoned connections that need verification._
- **Are the 76 inferred relationships involving `SuccessResponse()` (e.g. with `.Detail()` and `.EntityLogs()`) actually correct?**
  _`SuccessResponse()` has 76 INFERRED edges - model-reasoned connections that need verification._
- **Are the 53 inferred relationships involving `NewServiceError()` (e.g. with `.ActivateAccount()` and `.ChangePassword()`) actually correct?**
  _`NewServiceError()` has 53 INFERRED edges - model-reasoned connections that need verification._
- **Are the 41 inferred relationships involving `GetAuditMeta()` (e.g. with `.ChangePassword()` and `.Logout()`) actually correct?**
  _`GetAuditMeta()` has 41 INFERRED edges - model-reasoned connections that need verification._