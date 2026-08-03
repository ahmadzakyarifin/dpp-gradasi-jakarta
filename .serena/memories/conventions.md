# Konvensi Backend dpp-gradasi-jakarta (backend/)

- Bahasa komunikasi: Bahasa Indonesia.
- Layer order: dto → entity → model → mapper → handler → service → repository.
- DTO: tag `json` + `binding` (contoh `Name string \`json:"name" binding:"required,min=2"\``).
- Entity: KOSONG tanpa tag (domain murni). Model: tag `gorm` untuk column mapping.
- Mapper 2 arah: `CreateReqToEntity`, `UpdateReqToEntity(req, existing)` mutasi, `EntityToModel`, `ModelToEntity`, `EntityToResponse`.
- Handler TIPIS: constructor 1 arg `NewXHandler(svc)`, bind `ShouldBindJSON`, `validator.Errors(err)` → 422, `helper.HandleServiceError(c, err)`, `helper.SuccessResponse(c, code, msg, data, meta)`. Handler TIDAK pegang logSvc.
- AUDIT DI SERVICE: struct service punya `audit activitylogservice.ActivityLogService` + `db *gorm.DB`; method `log(ctx, input)` isi ActorID/ActorName/ActorRole/IPAddress/UserAgent dari `helper.GetAuditMeta(ctx)` (guard `if s.audit == nil`); panggil `s.audit.Log(ctx, s.db, input)` di SETIAP method mutasi. Action naming `<modul>.<aksi>`.
- Service constructor: `NewXService(db, repo, audit)`.
- Repository: satu-satunya tempat query GORM.
- Typed errors: `helper/service_error.go` — ValidationError{Fields,Errors}, AuthenticationError, NotFoundError, ServiceError. `HandleServiceError` switch semuanya.
- Rate limit: in-memory fixed-window (`internal/middleware/rate_limiter_fixed.go`), tanpa Redis. Rate rules per-path di `rate_rules.go`.
- Risk map: `internal/module/activitylog/service/risk.go` — highRisk (auth.login_failed, auth.forgot_password_spam, auth.reset_password, auth.change_password, auth.*_failed, users.delete/bulk_delete/toggle_status/change_role, roles.update/delete/bulk_delete), mediumRisk (sisanya). Sliders pakai prefix `slider.*` (bukan `sliders.*`).
- Email: `infrastructure.Mailer` (gomail sync) langsung dari service, bukan job/queue.