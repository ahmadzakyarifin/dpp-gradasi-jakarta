# Bugfix Requirements Document

## Introduction

The Activity Log feature (backend module `activitylog`, admin page `ActivityLogAdmin.jsx`, and its consumer in the Dashboard summary) has drifted out of sync across the frontend, backend, and documentation layers. Verification against the current codebase confirmed 12 distinct defects spanning: a frontend/backend field-name mismatch that leaves most table columns blank, a pagination metadata mismatch that freezes the pager, a typo'd action string that zeroes out a dashboard stat, an inconsistent `entity_type` taxonomy that breaks entity filtering, a role filter that can never match stored data, bulk role actions that collide with single-item action codes in the audit trail, two endpoints that double-wrap their JSON response, a timestamp formatter that mislabels local time as UTC, a misspelled package directory, three risk classifications that are declared but never reachable (under-classifying genuinely risky actions), a documentation/implementation mismatch for query parameter names, and a dashboard stat that reads the wrong summary field. Two of these (the second stray `entity_type: "auth"` call site in `ResetPassword`, and the dashboard `FinanceAction`/`CMSAction` mixup) were found during this investigation's own verification pass, in addition to the ones already reported.

This document defines the defective behavior, the corrected behavior, and the behavior that must remain unchanged for each of the 12 issues.

## Bug Analysis

### Current Behavior (Defect)

1.1 WHEN the admin activity log table renders rows or the CSV export is generated (`ActivityLogAdmin.jsx`) THEN the system reads `log.time`, `log.actor`, `log.role`, `log.entity`, `log.entityLabel`, `log.ip`, `log.device`, and `log.risk` — none of which exist on the API response object — so the Waktu, Aktor, Role, IP, Device, Entity, and Risiko columns render blank or undefined.

1.2 WHEN the activity log list response includes pagination metadata THEN the frontend reads `res.data.pagination.total` and `res.data.pagination.totalPages`, which do not exist on the response (the backend returns `res.data.meta.total_data` and `res.data.meta.total_pages`), so the pagination footer text and page buttons never update from their initial defaults (0 and 1).

1.3 WHEN the activity log summary statistics are computed (`GetSummary`) THEN the repository filters on the action string `"auth.failed_login"`, but every actual login-failure call site logs `"auth.login_failed"` (reversed word order), so `summary.failed_login` is always 0 and `summary.cms_action` (which excludes `"auth.failed_login"`) wrongly includes every failed-login attempt.

1.4 WHEN activity logs are written across modules THEN `entity_type` values are inconsistent: `auth_service.go` logs `"user"` for most auth actions but logs `"auth"` at two call sites where no user has been identified yet (login with an unrecognized email, and password reset with an invalid/expired token); `sliders_service.go` logs `"slider"` (singular) while the frontend filter dropdown and every other CMS module use the plural module name; and the frontend entity filter dropdown also offers a `"database"` option that no code path ever logs. As a result, filtering by Sliders, Auth, or Database in the admin UI returns incomplete or permanently empty results.

1.5 WHEN the admin selects a role in the activity log's role filter dropdown THEN the dropdown submits the display-style values `"Super Admin"` or `"Admin"`, but `actor_role` is stored in the database as the role slug (`"super_admin"`, `"admin"`) and the repository performs an exact match, so selecting either role filter option always returns zero rows.

1.6 WHEN roles are deleted or restored in bulk THEN `RoleService.BulkDelete` and `BulkRestore` log the action codes `"roles.delete"` and `"roles.restore"` — identical to the codes used for single-item delete/restore — instead of `"roles.bulk_delete"` and `"roles.bulk_restore"`, so bulk operations are indistinguishable from single-item operations in the audit trail, and the existing `risk.go` entries for `"roles.bulk_delete"` (highRisk) and `"roles.bulk_restore"` (mediumRisk) are never reachable.

1.7 WHEN the Detail or EntityLogs activity log endpoints return a successful response THEN the handler wraps the payload as `gin.H{"data": res}` before passing it to `helper.SuccessResponse`, which independently nests whatever it receives under the response's own `"data"` key, producing a doubly-nested `{ "data": { "data": {...} } }` body instead of the documented flat `{ "data": {...} }` shape.

1.8 WHEN an activity log's `created_at` timestamp is formatted for an API response THEN `EntityToResponse` and `EntityToDetailResponse` unconditionally append a literal `"Z"` UTC designator via `Format("2006-01-02T15:04:05Z")`, even though the database connection is opened with `loc=Local` (`backend/internal/infrastructure/database.go`), so the emitted timestamp is falsely labeled as UTC when it is actually in the server's local time.

1.9 WHEN a developer or import statement references the activitylog module's mapper package THEN the directory is named `maper` (missing a "p") instead of `mapper`, unlike every other module in the codebase, and this misspelling is baked into the Go import path in two files.

1.10 WHEN certain sensitive actions occur THEN they are not classified at the risk level the `risk.go` map was designed to assign: changing a user's role goes through `UserService.Update`, which only ever logs the generic `"users.update"` (mediumRisk) — never the highRisk `"users.change_role"` entry that already exists in the map; `RoleService.UpdateStatus` logs `"roles.update"` instead of the mediumRisk `"roles.status_update"` entry that already exists in the map; and `"auth.forgot_password_spam"` is declared highRisk but no code path anywhere emits it.

1.11 WHEN `docs-final/api/activity_logs.jsonc` documents the List endpoint's query parameters THEN it names them `entity_type` and `risk_level`, but the actual backend query binding (`ActivityLogQueryReq`) and the frontend service both use `entity` and `risk`, so the published contract does not match the real implementation. The documented summary response example also omits the `finance_action` field that `ActivityLogSummaryRes` actually returns.

1.12 WHEN the admin dashboard summary is computed (`DashboardService.GetSummary`) THEN the system assigns `res.CMSActions = summary.FinanceAction` — reading the always-zero finance field from the activity log summary instead of `summary.CMSAction` — so the dashboard's `cms_actions` statistic is always 0 regardless of actual CMS activity.

### Expected Behavior (Correct)

2.1 WHEN the admin activity log table renders rows or the CSV export is generated THEN the system SHALL read the actual response field names (`created_at`, `actor_name`, `actor_role`, `entity_type`, `entity_label`, `ip_address`, `user_agent`, `risk_level`) so every column displays real data.

2.2 WHEN the activity log list response includes pagination metadata THEN the frontend SHALL read `res.data.meta.total_data` and `res.data.meta.total_pages` so the pagination footer text and page buttons reflect the real result count and page count.

2.3 WHEN the activity log summary statistics are computed THEN the repository SHALL filter on the action string `"auth.login_failed"` for both the failed-login count and the CMS-action exclusion, so `summary.failed_login` reflects real failed logins and `summary.cms_action` correctly excludes them.

2.4 WHEN activity logs are written across modules THEN `entity_type` SHALL be consistent per module: all `auth.*` actions in `auth_service.go` (including the two call sites with no identified user) SHALL use `"user"`; the sliders module SHALL log `"sliders"` (plural, matching the module name and the frontend filter); and the frontend entity filter dropdown SHALL only list entity types that are actually ever logged, removing the dead `"database"` option.

2.5 WHEN the admin selects a role in the activity log's role filter THEN the dropdown SHALL submit the actual role slugs stored in the database (`"super_admin"`, `"admin"`) so the filter returns matching rows.

2.6 WHEN roles are deleted or restored in bulk THEN `RoleService.BulkDelete` and `BulkRestore` SHALL log `"roles.bulk_delete"` and `"roles.bulk_restore"` respectively, distinguishing bulk operations from single-item operations in the audit trail and making the existing `risk.go` classifications for those codes reachable.

2.7 WHEN the Detail or EntityLogs activity log endpoints return a successful response THEN the handler SHALL pass the response object directly to `helper.SuccessResponse` (matching the pattern already used by the List and Summary handlers) so the body matches the documented flat `{ "data": {...} }` contract shape.

2.8 WHEN an activity log's `created_at` timestamp is formatted for an API response THEN the system SHALL ensure the emitted timezone designator accurately reflects the actual timezone of the underlying value (normalizing to UTC before formatting with a correct RFC3339 formatter), so consumers never misinterpret the timestamp.

2.9 WHEN a developer or import statement references the activitylog module's mapper package THEN the directory SHALL be named `mapper`, consistent with every other module, and all import paths SHALL reference the corrected path.

2.10 WHEN a user's role is changed via `UserService.Update` THEN the system SHALL emit `"users.change_role"` (highRisk) to reflect the sensitivity of privilege changes; WHEN `RoleService.UpdateStatus` changes a role's active status THEN the system SHALL emit `"roles.status_update"` (mediumRisk) instead of reusing `"roles.update"`; and the dead `"auth.forgot_password_spam"` classification entry SHALL be removed since no rate-limiting hook exists to emit it.

2.11 WHEN `docs-final/api/activity_logs.jsonc` documents the List endpoint's query parameters THEN it SHALL name them `entity` and `risk` to match the real implementation, and the documented summary response example SHALL include the `finance_action` field.

2.12 WHEN the admin dashboard summary is computed THEN `DashboardService.GetSummary` SHALL assign `res.CMSActions = summary.CMSAction` so the dashboard's `cms_actions` statistic reflects real CMS activity counts.

### Unchanged Behavior (Regression Prevention)

3.1 WHEN the activity log list response includes `id`, `action`, and `description` fields THEN the system SHALL CONTINUE TO render them correctly in the table and CSV export, since these fields already match between frontend and backend.

3.2 WHEN the admin applies search, actor filter, start/end date range, or sort parameters to the activity log list THEN the system SHALL CONTINUE TO apply them exactly as before, since this query logic is untouched by the fixes.

3.3 WHEN activity logs are written for berita, kegiatan, kontak, pengurus, and settings actions THEN the system SHALL CONTINUE TO use their existing, already-consistent `entity_type` and action-code values without modification.

3.4 WHEN a successful login, refresh-token, forgot-password, change-password, or logout occurs THEN the system SHALL CONTINUE TO log the same action codes and risk classifications as before (only the `entity_type` of the two "unknown user" auth call sites changes).

3.5 WHEN the List or Summary activity log endpoints return a successful response THEN the system SHALL CONTINUE TO return the response body unwrapped exactly as before — only Detail and EntityLogs change.

3.6 WHEN non-super_admin users attempt to access any activity log endpoint THEN the system SHALL CONTINUE TO be denied access via the existing role-based middleware.

3.7 WHEN metadata is attached to an activity log entry (for example old/new values on update actions) THEN the system SHALL CONTINUE TO store and return it unchanged.

3.8 WHEN risk levels are already correctly assigned to existing, actively-emitted action codes (for example `"users.delete"` = high, `"berita.create"` = medium) THEN the system SHALL CONTINUE TO classify them the same way.

3.9 WHEN the admin filters the activity log list by risk level (high/medium/low) THEN the system SHALL CONTINUE TO filter correctly, since risk-level filtering is unaffected by these fixes.

3.10 WHEN CSV export is triggered THEN the system SHALL CONTINUE TO produce a file with the same column headers, ordering, and download filename format as before — only the underlying field values read for each row change from blank to populated.

3.11 WHEN the dashboard summary reports `total_berita`, `total_kegiatan`, `total_pengurus`, `total_kontak`, `unread_kontak`, `total_admin`, `pending_admin`, `activity_logs`, `high_risk_logs`, and `failed_login` THEN the system SHALL CONTINUE TO compute them exactly as before — only `cms_actions` changes.
