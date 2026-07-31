# Graph Report - dpp_gradasi  (2026-07-31)

## Corpus Check
- 110 files · ~67,256 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 772 nodes · 1595 edges · 47 communities (32 shown, 15 thin omitted)
- Extraction: 76% EXTRACTED · 24% INFERRED · 0% AMBIGUOUS · INFERRED: 389 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `5354c2e9`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- ErrorResponse
- KegiatanRepo
- BeritaRepo
- UserService
- AuthService
- NewServiceError
- Rule
- activityLogRepository
- RoleMiddleware
- KontakRepo
- SlidersRepo
- response.go
- AuthRepo
- Config
- rate_limiter_middleware.go
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
- RateLimitPerUser
- Errors
- 00003_create_auth_tokens_tables.sql
- RegisterRoutes
- RegisterRoutes
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

## God Nodes (most connected - your core abstractions)
1. `ErrorResponse()` - 77 edges
2. `SuccessResponse()` - 74 edges
3. `NewServiceError()` - 53 edges
4. `GetAuditMeta()` - 41 edges
5. `NewApp()` - 37 edges
6. `ValidationErrorResponse()` - 35 edges
7. `AuthRepo` - 24 edges
8. `Config` - 22 edges
9. `AuthService` - 20 edges
10. `KegiatanRepo` - 20 edges

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

## Communities (47 total, 15 thin omitted)

### Community 0 - "ErrorResponse"
Cohesion: 0.08
Nodes (23): GetAuditMeta(), Context, ErrorResponse(), Context, SuccessResponse(), ValidationErrorResponse(), GetUserID(), Context (+15 more)

### Community 1 - "KegiatanRepo"
Cohesion: 0.06
Nodes (28): BulkRequest, PaginationMeta, PaginationMeta, DeletedAt, Time, DB, maxInt(), NewKegiatanRepo() (+20 more)

### Community 2 - "BeritaRepo"
Cohesion: 0.07
Nodes (22): BulkRequest, PaginationMeta, PaginationMeta, DeletedAt, Time, DB, NewBeritaRepo(), formatDate() (+14 more)

### Community 3 - "UserService"
Cohesion: 0.07
Nodes (20): BulkRequest, ChangePasswordRequest, FileHeader, DeletedAt, Time, DB, NewUserRepo(), generateRandomNumber() (+12 more)

### Community 4 - "AuthService"
Cohesion: 0.07
Nodes (29): GenerateAccessToken(), GenerateRefreshToken(), RegisteredClaims, Time, ValidateToken(), CheckPassword(), HashPassword(), GenerateActivationToken() (+21 more)

### Community 5 - "NewServiceError"
Cohesion: 0.07
Nodes (22): NewServiceError(), BulkRequest, PaginationMeta, FileHeader, PaginationMeta, DeletedAt, Time, DB (+14 more)

### Community 6 - "Rule"
Cohesion: 0.22
Nodes (16): normalizeRule(), normalizeScope(), abortRateLimitError(), abortRateLimitUnauthorized(), abortTooManyRequests(), Email(), Context, IP() (+8 more)

### Community 7 - "activityLogRepository"
Cohesion: 0.09
Nodes (24): EntityToDetailResponse(), EntityToModel(), EntityToResponse(), ActivityLog, InputToEntity(), ModelToEntity(), ActivityLog, Context (+16 more)

### Community 8 - "RoleMiddleware"
Cohesion: 0.13
Nodes (11): HandlerFunc, RoleMiddleware(), Client, RouterGroup, RegisterRoutes(), Client, RouterGroup, RegisterRoutes() (+3 more)

### Community 9 - "KontakRepo"
Cohesion: 0.08
Nodes (20): BulkRequest, PaginationMeta, PaginationMeta, DeletedAt, Time, DB, maxInt(), NewKontakRepo() (+12 more)

### Community 10 - "SlidersRepo"
Cohesion: 0.09
Nodes (14): BulkRequest, DeletedAt, Time, DB, NewSlidersRepo(), NewSlidersService(), toResponse(), ReorderRequest (+6 more)

### Community 11 - "response.go"
Cohesion: 0.10
Nodes (17): GenerateValidationMessage(), GetPaginationMeta(), IsValidEmail(), IsValidURL(), RateLimitResponse(), NewSettingsHandler(), DB, NewSettingsRepo() (+9 more)

### Community 12 - "AuthRepo"
Cohesion: 0.08
Nodes (9): Time, DB, NewAuthRepo(), Time, ActivationToken, PasswordResetToken, RefreshToken, Role (+1 more)

### Community 13 - "Config"
Cohesion: 0.11
Nodes (10): MustLoad(), AppConfig, Config, CookieConfig, DatabaseConfig, DevConfig, JWTConfig, RedisConfig (+2 more)

### Community 14 - "rate_limiter_middleware.go"
Cohesion: 0.18
Nodes (15): Context, RestoreBody(), emailFromBody(), getUserID(), Client, Context, makeRateLimitKey(), needsEmail() (+7 more)

### Community 15 - "ActivityLogService"
Cohesion: 0.15
Nodes (8): NewActivityLogHandler(), ActivityLogService, NewAuthHandler(), NewBeritaHandler(), NewKegiatanHandler(), NewPengurusHandler(), NewUserHandler(), ActivityLogHandler

### Community 16 - "NewApp"
Cohesion: 0.19
Nodes (9): App, Client, DB, NewApp(), Client, registerRoutes(), NewKontakHandler(), NewSlidersHandler() (+1 more)

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
Cohesion: 0.50
Nodes (3): dependencies, puppeteer, puppeteer

### Community 24 - "schema.sql"
Cohesion: 0.22
Nodes (14): activation_tokens, berita, berita_tags, kegiatan, kegiatan_gallery, kegiatan_tags, password_reset_tokens, pengurus (+6 more)

### Community 25 - "AuthMiddleware"
Cohesion: 0.16
Nodes (11): AuthMiddleware(), GetEmail(), GetRoleID(), Context, HandlerFunc, Client, RouterGroup, RegisterRoutes() (+3 more)

### Community 30 - "RateLimitPerUser"
Cohesion: 0.20
Nodes (10): HandlerFunc, RateLimiterMiddleware(), RateLimitPerUser(), RateLimitRules(), Client, RouterGroup, RegisterRoutes(), Client (+2 more)

### Community 31 - "Errors"
Cohesion: 0.50
Nodes (3): validationMessage(), Errors(), ValidationErrorItem

### Community 32 - "00003_create_auth_tokens_tables.sql"
Cohesion: 0.60
Nodes (4): activation_tokens, password_reset_tokens, refresh_tokens, users

### Community 33 - "RegisterRoutes"
Cohesion: 0.50
Nodes (3): Client, RouterGroup, RegisterRoutes()

### Community 34 - "RegisterRoutes"
Cohesion: 0.50
Nodes (3): Client, RouterGroup, RegisterRoutes()

### Community 35 - "berita"
Cohesion: 0.67
Nodes (3): berita, berita_tags, users

### Community 36 - "00006_create_kegiatan_tables.sql"
Cohesion: 0.83
Nodes (3): kegiatan, kegiatan_gallery, kegiatan_tags

## Knowledge Gaps
- **30 isolated node(s):** `WorkerConfig`, `github.com/ahmadzakyarifin/dpp-gradasi/backend`, `Response`, `RegisterRequest`, `ChangePasswordRequest` (+25 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **15 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `NewApp()` connect `NewApp` to `KegiatanRepo`, `BeritaRepo`, `UserService`, `AuthService`, `NewServiceError`, `activityLogRepository`, `KontakRepo`, `SlidersRepo`, `response.go`, `AuthRepo`, `Config`, `rate_limiter_middleware.go`, `ActivityLogService`, `main`, `Mailer`?**
  _High betweenness centrality (0.301) - this node is a cross-community bridge._
- **Why does `NewServiceError()` connect `NewServiceError` to `ErrorResponse`, `KegiatanRepo`, `BeritaRepo`, `UserService`, `AuthService`, `KontakRepo`, `SlidersRepo`?**
  _High betweenness centrality (0.244) - this node is a cross-community bridge._
- **Why does `ErrorResponse()` connect `ErrorResponse` to `RoleMiddleware`, `AuthMiddleware`, `response.go`, `Rule`?**
  _High betweenness centrality (0.088) - this node is a cross-community bridge._
- **Are the 75 inferred relationships involving `ErrorResponse()` (e.g. with `AuthMiddleware()` and `abortRateLimitError()`) actually correct?**
  _`ErrorResponse()` has 75 INFERRED edges - model-reasoned connections that need verification._
- **Are the 72 inferred relationships involving `SuccessResponse()` (e.g. with `.Detail()` and `.EntityLogs()`) actually correct?**
  _`SuccessResponse()` has 72 INFERRED edges - model-reasoned connections that need verification._
- **Are the 51 inferred relationships involving `NewServiceError()` (e.g. with `.ActivateAccount()` and `.ChangePassword()`) actually correct?**
  _`NewServiceError()` has 51 INFERRED edges - model-reasoned connections that need verification._
- **Are the 39 inferred relationships involving `GetAuditMeta()` (e.g. with `.ChangePassword()` and `.Logout()`) actually correct?**
  _`GetAuditMeta()` has 39 INFERRED edges - model-reasoned connections that need verification._