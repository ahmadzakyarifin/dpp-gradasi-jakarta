# Graph Report - dpp_gradasi  (2026-07-31)

## Corpus Check
- 110 files · ~67,256 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 724 nodes · 1547 edges · 30 communities (22 shown, 8 thin omitted)
- Extraction: 75% EXTRACTED · 25% INFERRED · 0% AMBIGUOUS · INFERRED: 389 edges (avg confidence: 0.8)
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
- PengurusService
- rate_limiter_middleware.go
- activityLogRepository
- AuthMiddleware
- KontakRepo
- NewServiceError
- response.go
- AuthRepo
- Config
- SlidersRepo
- ActivityLogService
- NewApp
- main
- ActivityLog
- Mailer
- dependencies
- ActivityLog
- scratch.js
- worker.go
- NewAuthHandler
- NewUserHandler
- README.md
- test-ui.js
- github.com/ahmadzakyarifin/dpp-gradasi/backend

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

## Communities (30 total, 8 thin omitted)

### Community 0 - "ErrorResponse"
Cohesion: 0.10
Nodes (20): GetAuditMeta(), Context, ErrorResponse(), Context, RateLimitResponse(), SuccessResponse(), ValidationErrorResponse(), Context (+12 more)

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

### Community 5 - "PengurusService"
Cohesion: 0.08
Nodes (20): BulkRequest, PaginationMeta, FileHeader, PaginationMeta, DeletedAt, Time, DB, NewPengurusRepo() (+12 more)

### Community 6 - "rate_limiter_middleware.go"
Cohesion: 0.09
Nodes (38): Context, RestoreBody(), emailFromBody(), getUserID(), Client, Context, HandlerFunc, makeRateLimitKey() (+30 more)

### Community 7 - "activityLogRepository"
Cohesion: 0.09
Nodes (24): EntityToDetailResponse(), EntityToModel(), EntityToResponse(), ActivityLog, InputToEntity(), ModelToEntity(), ActivityLog, Context (+16 more)

### Community 8 - "AuthMiddleware"
Cohesion: 0.06
Nodes (31): AuthMiddleware(), GetEmail(), GetRoleID(), Context, HandlerFunc, HandlerFunc, RoleMiddleware(), Client (+23 more)

### Community 9 - "KontakRepo"
Cohesion: 0.08
Nodes (20): BulkRequest, PaginationMeta, PaginationMeta, DeletedAt, Time, DB, maxInt(), NewKontakRepo() (+12 more)

### Community 10 - "NewServiceError"
Cohesion: 0.11
Nodes (13): NewServiceError(), GetUserID(), Context, BulkRequest, NewSlidersService(), toResponse(), ReorderRequest, SliderListResponse (+5 more)

### Community 11 - "response.go"
Cohesion: 0.08
Nodes (20): GenerateValidationMessage(), GetPaginationMeta(), IsValidEmail(), IsValidURL(), Context, NewSettingsHandler(), DB, NewSettingsRepo() (+12 more)

### Community 12 - "AuthRepo"
Cohesion: 0.08
Nodes (9): Time, DB, NewAuthRepo(), Time, ActivationToken, PasswordResetToken, RefreshToken, Role (+1 more)

### Community 13 - "Config"
Cohesion: 0.11
Nodes (10): MustLoad(), AppConfig, Config, CookieConfig, DatabaseConfig, DevConfig, JWTConfig, RedisConfig (+2 more)

### Community 14 - "SlidersRepo"
Cohesion: 0.17
Nodes (6): DeletedAt, Time, DB, NewSlidersRepo(), Slider, SlidersRepo

### Community 15 - "ActivityLogService"
Cohesion: 0.18
Nodes (7): NewActivityLogHandler(), ActivityLogService, NewBeritaHandler(), NewKegiatanHandler(), NewKontakHandler(), NewSlidersHandler(), ActivityLogHandler

### Community 16 - "NewApp"
Cohesion: 0.24
Nodes (8): App, Client, DB, NewApp(), Client, registerRoutes(), NewPengurusHandler(), Engine

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

## Knowledge Gaps
- **23 isolated node(s):** `WorkerConfig`, `github.com/ahmadzakyarifin/dpp-gradasi/backend`, `Response`, `RegisterRequest`, `ChangePasswordRequest` (+18 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **8 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `NewApp()` connect `NewApp` to `KegiatanRepo`, `BeritaRepo`, `UserService`, `AuthService`, `PengurusService`, `rate_limiter_middleware.go`, `activityLogRepository`, `KontakRepo`, `NewServiceError`, `response.go`, `AuthRepo`, `Config`, `SlidersRepo`, `ActivityLogService`, `main`, `Mailer`, `NewAuthHandler`, `NewUserHandler`?**
  _High betweenness centrality (0.343) - this node is a cross-community bridge._
- **Why does `NewServiceError()` connect `NewServiceError` to `ErrorResponse`, `KegiatanRepo`, `BeritaRepo`, `UserService`, `AuthService`, `PengurusService`, `KontakRepo`?**
  _High betweenness centrality (0.278) - this node is a cross-community bridge._
- **Why does `ErrorResponse()` connect `ErrorResponse` to `AuthMiddleware`, `NewServiceError`, `response.go`, `rate_limiter_middleware.go`?**
  _High betweenness centrality (0.100) - this node is a cross-community bridge._
- **Are the 75 inferred relationships involving `ErrorResponse()` (e.g. with `AuthMiddleware()` and `abortRateLimitError()`) actually correct?**
  _`ErrorResponse()` has 75 INFERRED edges - model-reasoned connections that need verification._
- **Are the 72 inferred relationships involving `SuccessResponse()` (e.g. with `.Detail()` and `.EntityLogs()`) actually correct?**
  _`SuccessResponse()` has 72 INFERRED edges - model-reasoned connections that need verification._
- **Are the 51 inferred relationships involving `NewServiceError()` (e.g. with `.ActivateAccount()` and `.ChangePassword()`) actually correct?**
  _`NewServiceError()` has 51 INFERRED edges - model-reasoned connections that need verification._
- **Are the 39 inferred relationships involving `GetAuditMeta()` (e.g. with `.ChangePassword()` and `.Logout()`) actually correct?**
  _`GetAuditMeta()` has 39 INFERRED edges - model-reasoned connections that need verification._