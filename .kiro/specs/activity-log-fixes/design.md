# Activity Log Fixes Bugfix Design

## Overview

This is a batch bugfix covering 12 confirmed defects in the Activity Log feature and its dashboard consumer. The defects fall into three categories:

1. **Contract desync** (frontend reads fields/paths that don't exist on the real API response, backend documents params it doesn't bind) — items A, B, K.
2. **String literal drift** (typo'd or inconsistent action/entity-type strings that silently break filters, counts, and audit-trail granularity) — items C, D, E, F, J, and the dashboard `FinanceAction`/`CMSAction` mixup found during verification.
3. **Structural/formatting defects** (double-wrapped JSON envelope, mislabeled timestamp timezone, misspelled package directory) — items G, H, I.

Because this batch spans many independent, mostly-deterministic string/field defects rather than one bug over a continuous input space, this design groups the 12 items into 9 independently-testable clusters (each gets its own Bug Condition / Preservation property pair, Properties 1-18) plus 2 structural items verified by build/review rather than PBT (mapper directory rename, documentation-only drift). Each cluster is fixed at its root cause: the side that is wrong is corrected to match the side that is right (verified against the DB schema, the documented contract, and the pattern used by every other already-correct module).

## Glossary

- **Bug_Condition (C)**: The condition that identifies a defective code path for a given cluster — e.g., "the frontend reads `log.time`" or "the repository filters on `auth.failed_login`".
- **Property (P)**: The desired behavior once the fix is applied — e.g., "the frontend reads `log.created_at`" or "the repository filters on `auth.login_failed`".
- **Preservation**: Behavior that must be byte-for-byte identical before and after the fix for inputs outside the bug condition (e.g., the `id`/`action`/`description` fields, which already matched correctly).
- **ActivityLogItemRes**: The backend response shape (`backend/internal/module/activitylog/dto/activity_log_response.go`) — the source of truth for field names, confirmed to match `docs-final/api/activity_logs.jsonc`.
- **entity_type taxonomy**: The set of string values used across modules to tag which kind of entity an activity log row is about (e.g., `"user"`, `"berita"`, `"sliders"`).
- **actor_role**: The stored audit column populated from `role.Name` (the slug, e.g. `"super_admin"`), not `role.DisplayName`.
- **F**: The original (unfixed) function/file. **F'**: The fixed function/file.

## Bug Details

### Bug Condition

Each cluster below is an independent, mostly-deterministic defect (a specific wrong literal, wrong field path, or wrong wrapping), rather than a bug that manifests across a continuous input space. The dispatch function identifies which cluster an observation belongs to:

**Formal Specification:**
```
FUNCTION isBugCondition(X)
  INPUT: X of type ActivityLogObservation (tagged by cluster)
  OUTPUT: boolean

  RETURN CASE X.cluster OF
    "frontend_field_mapping"   => X.source IN {ActivityLogAdmin.jsx row render, CSV export}
                                  AND X.fieldRead IN {time, actor, role, entity, entityLabel, ip, device, risk}
    "pagination_meta"          => X.source = "fetchLogs" AND X.pathRead IN {pagination.total, pagination.totalPages}
    "summary_action_typo"      => X.query = "GetSummary" AND X.literal = "auth.failed_login"
    "entity_type_taxonomy"     => X.entityType IN {"auth" (from user-unidentified auth call sites), "slider" (singular)}
                                  OR (X.source = "frontend dropdown" AND X.option = "database")
    "role_filter_value"        => X.source = "frontend role dropdown" AND X.optionValue IN {"Super Admin", "Admin"}
    "bulk_action_code"         => X.action IN {"roles.delete" WHERE X.callSite = BulkDelete,
                                                "roles.restore" WHERE X.callSite = BulkRestore}
    "handler_double_wrap"      => X.handler IN {Detail, EntityLogs} AND X.responseData = gin.H{"data": res}
    "timestamp_timezone"       => X.formatter USES literal "Z" AND X.dbLoc <> "UTC"
    "risk_classification_gap"  => X.action IN {"users.update" WHERE roleChanged = true,
                                                "roles.update" WHERE onlyStatusChanged = true}
                                  OR X.riskMapEntry = "auth.forgot_password_spam" AND NOT EXISTS emitter(X.riskMapEntry)
    "dashboard_field_mixup"    => X.source = "DashboardService.GetSummary" AND X.fieldRead = "FinanceAction"
                                  WHERE X.targetStat = "CMSActions"
  END CASE
END FUNCTION
```

### Examples

- **Frontend field mapping**: `ActivityLogItemRes.ActorName` serializes to JSON key `actor_name`. `ActivityLogAdmin.jsx` reads `log.actor` (undefined) instead of `log.actor_name`. Table renders "System" (the `|| 'System'` fallback) for every row regardless of actual actor.
- **Pagination meta**: Backend returns `data.meta.total_data = 47`. Frontend reads `res.data.pagination.total`, which is `undefined`, so `setTotal(res.data.pagination.total || 0)` sets state to `0` forever. Footer shows "Menampilkan 1 sampai 10 dari 0 hasil" even though 47 rows exist.
- **Summary action typo**: A real failed-login row is stored with `action = "auth.login_failed"`. `GetSummary`'s query `WHERE action = "auth.failed_login"` matches zero rows. `summary.failed_login` is always `0`; the inverted query `WHERE action <> "auth.failed_login"` matches the failed-login row too, inflating `cms_action`.
- **Entity type taxonomy**: `auth_service.go` line ~74 (`Login`, email-not-found branch) and the `ResetPassword` invalid-token branch both log `EntityType: "auth"`, while the other 5 auth call sites (including 3 other login-failure branches) log `EntityType: "user"`. Selecting "Auth" in the frontend filter returns only those 2 rows out of ~7 auth-related actions; selecting "Sliders" (plural, what the dropdown offers) matches zero rows because `sliders_service.go` stores `"slider"` (singular).
- **Role filter value**: DB stores `actor_role = "super_admin"`. Frontend sends `role=Super+Admin`. Repository does `WHERE actor_role = 'Super Admin'` → 0 rows, always.
- **Bulk action code**: Deleting 5 roles at once logs one row with `Action: "roles.delete"` — identical to what a single role deletion logs. The audit trail cannot distinguish "admin deleted 1 role" from "admin deleted 5 roles."
- **Handler double-wrap**: `Detail` returns `{"data": {"data": {...actual log...}}}` instead of `{"data": {...actual log...}}}`. Any future consumer following the documented contract (`data.id`, `data.action`, ...) would read `undefined`.
- **Timestamp timezone**: DB connection DSN is `...?charset=utf8mb4&parseTime=True&loc=Local` (confirmed in `backend/internal/infrastructure/database.go`). A row created at 13:35 local time (e.g. WIB, UTC+7) is formatted as `"...13:35:00Z"` — a correct RFC3339 parser reads this as 13:35 **UTC**, 7 hours off from the real instant.
- **Risk classification gap**: Admin changes a user's `role_id` via `PUT /admin/users/:id`. `UserService.Update` logs `Action: "users.update"` (mediumRisk) — identical to editing someone's display name. The `"users.change_role"` highRisk entry in `risk.go` is never reached.
- **Dashboard field mixup**: `summary.CMSAction = 86` (real CMS activity count), `summary.FinanceAction = 0` (always zero, per the summary-typo/finance-taxonomy issue). `dashboard_service.go` assigns `res.CMSActions = summary.FinanceAction`, so the dashboard's "CMS Actions" stat always shows `0` instead of `86`.

## Expected Behavior

### Preservation Requirements

**Unchanged Behaviors:**
- The `id`, `action`, and `description` fields in the activity log table/CSV, which already read correctly, must continue to display the same values.
- Search, actor filter, date-range, sort, and risk-level filtering on the activity log list must continue to work exactly as before — none of these query paths are touched.
- `entity_type` values already correct and consistent (berita, kegiatan, kontak, pengurus, settings) must not change.
- Action codes and risk levels for auth actions that are not entity-type-mislabeled (login success, refresh, forgot-password, change-password, logout) must not change — only the `entity_type` of the two "no identified user yet" call sites changes.
- `List` and `Summary` handler responses, which are already correctly unwrapped, must not change shape.
- Role-based access control on all activity-log endpoints (`super_admin` only) is untouched.
- Metadata (`old_values`/`new_values`/etc.) attached to log entries must continue to be stored and returned unchanged.
- Already-correct risk classifications (e.g. `users.delete` = high, `berita.create` = medium) must not change.
- All other dashboard summary stats (`total_berita`, `total_kegiatan`, `total_pengurus`, `total_kontak`, `unread_kontak`, `total_admin`, `pending_admin`, `activity_logs`, `high_risk_logs`, `failed_login`) must not change.

**Scope:**
Each cluster's fix is scoped to the specific literal/field/wrapping identified in that cluster. Inputs/rows/requests that don't touch the specific wrong literal or field path are unaffected — for example, fixing `sliders_service.go`'s `entity_type` does not touch `berita_service.go`, and fixing the `Detail` handler's wrapping does not touch the `List` handler.

## Hypothesized Root Cause

1. **Frontend built against an assumed shape, never reconciled with the real DTO** (A, B): `ActivityLogAdmin.jsx` was likely written from an earlier draft contract or copied from a different module's pagination shape (`pagination.total`/`totalPages` doesn't match this module's `meta.total_data`/`total_pages` pattern used by berita/kegiatan/pengurus/kontak either — it appears to be a one-off guess).
2. **Copy-paste of a similarly-named but differently-spelled action string** (C): `"auth.failed_login"` vs `"auth.login_failed"` — classic word-order transposition typo, written once in the repository and never cross-checked against the service layer that actually emits the string.
3. **No shared constants for action/entity-type strings** (C, D, E, F): every service hand-writes string literals for `Action` and `EntityType`. Without a single source of truth (constants or an enum), typos and inconsistent pluralization (`slider` vs `sliders`, `roles.delete` reused for bulk) are structurally likely and undetectable at compile time.
4. **Handler-level copy-paste from a different helper convention** (G): `Detail`/`EntityLogs` wrap in `gin.H{"data": res}` while `List`/`Summary` pass `res` directly — likely copied from a different project/helper pattern where `SuccessResponse` did not already nest under `"data"`.
5. **Timestamp formatting written assuming a UTC-configured DB, without checking the actual DSN** (H): `Format("2006-01-02T15:04:05Z")` is a common copy-pasted pattern, but nobody checked that `loc=Local` is actually configured for this connection (also present in `auth_mapper.go`, confirming this is a copy-pasted pattern, not a one-off).
6. **Directory created with a typo early on, never caught because Go doesn't validate directory names against a convention** (I): `maper` vs `mapper` — a simple typo that compiles fine and was propagated into two import paths.
7. **`risk.go` was written speculatively ahead of the features that would need it** (J): entries for `users.change_role`, `roles.status_update`, and `auth.forgot_password_spam` were added anticipating functionality that was never wired up at the call site (role-change detection, status-update-specific logging, and forgot-password rate-limit-triggered audit logging).
8. **Field selected by name-proximity rather than by verifying the correct one** (dashboard mixup): `FinanceAction` and `CMSAction` are adjacent fields in `ActivityLogSummaryRes`; `CMSActions` (dashboard) and `FinanceAction` (activity log) merely share no name at all, suggesting a copy-paste of the wrong source field, likely while wiring up the dashboard against an earlier version of `ActivityLogSummaryRes` that only had a `FinanceAction`-named field for what is now `CMSAction`.

## Correctness Properties

Property 1: Bug Condition - Frontend Reads Wrong Response Field Names and Pagination Path

_For any_ activity log API response object (with arbitrary `actor_name`, `actor_role`, `entity_type`, `entity_label`, `ip_address`, `user_agent`, `risk_level`, `created_at` values) and arbitrary `meta.total_data`/`meta.total_pages` values, the fixed frontend row-mapping and pagination-extraction logic SHALL extract every field and every pagination value from its real path, rendering non-blank values and non-default pagination state.

**Validates: Requirements 2.1, 2.2**

Property 2: Preservation - Already-Correct Fields and Filters Unaffected

_For any_ activity log response object, the `id`, `action`, and `description` fields SHALL continue to render identically before and after the fix, and search/actor/date/sort/risk-level filter parameters SHALL continue to be sent identically to the backend.

**Validates: Requirements 3.1, 3.2, 3.9, 3.10**

Property 3: Bug Condition - Summary Action String Matches Real Emitted Value

_For any_ set of stored activity log rows containing some count N of rows with `action = "auth.login_failed"` and some count M of rows with other actions, the fixed `GetSummary` SHALL report `failed_login = N` and SHALL exclude exactly those N rows from `cms_action` (i.e. `cms_action` counts only the M non-failed-login rows, subject to the existing query filters).

**Validates: Requirements 2.3**

Property 4: Preservation - Other Summary Fields Unaffected

_For any_ set of stored activity log rows, `summary.total_logs` and `summary.high_risk` SHALL be computed identically before and after the fix (only the `action` literal used for `failed_login`/`cms_action` changes).

**Validates: Requirements 3.2**

Property 5: Bug Condition - entity_type Taxonomy Is Consistent Per Module

_For any_ of the 5 auth-service call sites that log with no identified user yet or with an identified user (login success/failure variants, refresh, forgot/reset/change password, logout), the fixed code SHALL log `EntityType: "user"`; for any sliders service mutation (create/update/delete/restore/bulk_delete/bulk_restore/reorder), the fixed code SHALL log `EntityType: "sliders"`; and the frontend entity filter dropdown SHALL NOT offer the `"database"` option.

**Validates: Requirements 2.4**

Property 6: Preservation - Non-Auth, Non-Sliders entity_type Values Unaffected

_For any_ activity log write from berita, kegiatan, kontak, pengurus, settings, users, or roles services, the `entity_type` value SHALL be identical before and after the fix.

**Validates: Requirements 3.3, 3.4**

Property 7: Bug Condition - Role Filter Dropdown Sends Matchable Slug Values

_For any_ selection of the role filter dropdown ("Super Admin" or "Admin" display option), the fixed dropdown SHALL submit the value `"super_admin"` or `"admin"` respectively, matching the slugs actually stored in `actor_role`.

**Validates: Requirements 2.5**

Property 8: Preservation - Other Filter Dropdowns Unaffected

_For any_ selection of the entity or risk filter dropdowns, the submitted value SHALL be unchanged by this fix (only the role dropdown's option values change).

**Validates: Requirements 3.2**

Property 9: Bug Condition - Bulk Role Actions Use Distinct Action Codes

_For any_ call to `RoleService.BulkDelete` or `BulkRestore` with a non-empty ID list, the fixed code SHALL log `Action: "roles.bulk_delete"` or `Action: "roles.bulk_restore"` respectively, distinct from the single-item `"roles.delete"`/`"roles.restore"` codes.

**Validates: Requirements 2.6**

Property 10: Preservation - Single-Item Role Actions Unaffected

_For any_ call to `RoleService.DeleteRole` or `RestoreRole` (single-item, not bulk), the logged action code SHALL remain `"roles.delete"`/`"roles.restore"` exactly as before.

**Validates: Requirements 3.8**

Property 11: Bug Condition - Detail and EntityLogs Return Flat Response Bodies

_For any_ valid request to the Detail or EntityLogs endpoints, the fixed handler SHALL produce a JSON body where `data` is the log object/array directly (matching the documented contract), not `{"data": {"data": ...}}`.

**Validates: Requirements 2.7**

Property 12: Preservation - List and Summary Wrapping Unaffected

_For any_ valid request to the List or Summary endpoints, the response body shape SHALL be identical before and after the fix.

**Validates: Requirements 3.5**

Property 13: Bug Condition - Timestamp Formatting Reflects Real Instant

_For any_ `time.Time` value read back from the database connection (configured with `loc=Local`), the fixed formatter SHALL produce an RFC3339 string that, when parsed by a correct RFC3339 parser, yields the same instant as the original `time.Time` value (i.e. the timezone designator is never a false `"Z"` for a non-UTC value).

**Validates: Requirements 2.8**

Property 14: Preservation - Non-Timestamp Fields and Existing UTC-Correct Cases Unaffected

_For any_ activity log detail/list response, all fields other than `created_at` SHALL be unchanged, and for any `time.Time` value whose location genuinely is UTC, the formatted output SHALL still correctly read as UTC.

**Validates: Requirements 3.1, 3.7**

Property 15: Bug Condition - Sensitive Actions Get Correct Risk Classification

_For any_ `UserService.Update` call where the resulting `RoleID` differs from the prior `RoleID`, the fixed code SHALL log `Action: "users.change_role"` (highRisk) instead of `"users.update"`; for any `RoleService.UpdateStatus` call, the fixed code SHALL log `Action: "roles.status_update"` (mediumRisk) instead of `"roles.update"`.

**Validates: Requirements 2.10**

Property 16: Preservation - Non-Role-Changing Updates and Existing Classifications Unaffected

_For any_ `UserService.Update` call where `RoleID` is unchanged, the logged action SHALL remain `"users.update"`; all risk classifications for action codes not touched by this cluster (e.g. `users.delete`, `berita.create`) SHALL remain unchanged.

**Validates: Requirements 3.8**

Property 17: Bug Condition - Dashboard CMS Actions Stat Reads the Correct Summary Field

_For any_ activity log summary where `CMSAction` and `FinanceAction` hold distinct values, the fixed `DashboardService.GetSummary` SHALL assign `res.CMSActions = summary.CMSAction`.

**Validates: Requirements 2.12**

Property 18: Preservation - Other Dashboard Stats Unaffected

_For any_ database state, `total_berita`, `total_kegiatan`, `total_pengurus`, `total_kontak`, `unread_kontak`, `total_admin`, `pending_admin`, `activity_logs`, `high_risk_logs`, and `failed_login` on the dashboard summary SHALL be computed identically before and after the fix.

**Validates: Requirements 3.11**

## Fix Implementation

### Changes Required

**Cluster 1 — Frontend field mapping + pagination meta**
**File**: `frontend/src/pages/admin/ActivityLogAdmin.jsx`
1. In `fetchLogs`, read `res.data.meta.total_data` / `res.data.meta.total_pages` instead of `res.data.pagination.total` / `res.data.pagination.totalPages`.
2. In `downloadCSV` and the table body, read `log.created_at`, `log.actor_name`, `log.actor_role`, `log.entity_type`, `log.entity_label`, `log.ip_address`, `log.user_agent`, `log.risk_level` instead of `log.time`, `log.actor`, `log.role`, `log.entity`, `log.entityLabel`, `log.ip`, `log.device`, `log.risk`.
3. The existing `(log.time || '').split(', ')` date-split logic assumed a pre-formatted display string; `created_at` is an ISO 8601 string, so the split logic is replaced with a small date/time formatter for the two display lines (date part, time part).

**Cluster 2 — Backend GetSummary action-string typo**
**File**: `backend/internal/module/activitylog/repository/activity_log_repo.go`
1. Change both `"auth.failed_login"` literals in `GetSummary` to `"auth.login_failed"`.

**Cluster 3 — entity_type taxonomy**
**Files**: `backend/internal/module/auth/service/auth_service.go`, `backend/internal/module/sliders/service/sliders_service.go`, `frontend/src/pages/admin/ActivityLogAdmin.jsx`
1. In `auth_service.go`, change the two `EntityType: "auth"` occurrences (in `Login`'s email-not-found branch, and `ResetPassword`'s invalid-token branch) to `EntityType: "user"`.
2. In `sliders_service.go`, change all six `EntityType: "slider"` occurrences (and the six `Action: "slider.*"` literals stay as-is per the existing convention documented in conventions.md — only `EntityType` is renamed, not the action prefix, since the report and conventions memory only flag `EntityType` plurality as the taxonomy signal the frontend filters on) to `EntityType: "sliders"`.
3. In `ActivityLogAdmin.jsx`, remove the `<option value="database">Database</option>` entity filter option.

**Cluster 4 — Role filter dropdown values**
**File**: `frontend/src/pages/admin/ActivityLogAdmin.jsx`
1. Change `<option value="Super Admin">Super Admin</option>` to `<option value="super_admin">Super Admin</option>` and `<option value="Admin">Admin</option>` to `<option value="admin">Admin</option>` (display text unchanged, only the submitted `value` changes).

**Cluster 5 — Role bulk action codes**
**File**: `backend/internal/module/role/service/role_service.go`
1. In `BulkDelete`, change `Action: "roles.delete"` to `Action: "roles.bulk_delete"`.
2. In `BulkRestore`, change `Action: "roles.restore"` to `Action: "roles.bulk_restore"`.

**Cluster 6 — Handler double-wrap**
**File**: `backend/internal/module/activitylog/handler/activity_log_handler.go`
1. In `Detail`, change `helper.SuccessResponse(c, http.StatusOK, "ACTIVITY_LOG_DETAIL", "...", gin.H{"data": res}, nil)` to pass `res` directly.
2. In `EntityLogs`, change `helper.SuccessResponse(c, http.StatusOK, "ACTIVITY_LOG_ENTITY_LISTED", "...", gin.H{"data": res}, nil)` to pass `res` directly.

**Cluster 7 — Timestamp UTC correctness**
**File**: `backend/internal/module/activitylog/maper/model_mapper.go` (path becomes `mapper/model_mapper.go` after the directory rename below)
1. In `EntityToResponse` and `EntityToDetailResponse`, change `e.CreatedAt.Format("2006-01-02T15:04:05Z")` to `e.CreatedAt.UTC().Format(time.RFC3339)`. Since the DB connection is confirmed `loc=Local` (not UTC), calling `.UTC()` first converts the local-time value to its correct UTC instant before formatting with the standard-library RFC3339 layout (which appends the correct `Z` only because the value is now genuinely UTC).
2. Add `"time"` to the file's imports.
3. This fix is scoped to the activitylog module only. The identical pattern in `backend/internal/module/auth/mapper/auth_mapper.go` is flagged in the summary as a follow-up recommendation, not fixed here (out of scope per the user's explicit scope note).

**Cluster 8 — Risk classification gaps**
**Files**: `backend/internal/module/user/service/user_service.go`, `backend/internal/module/role/service/role_service.go`, `backend/internal/module/activitylog/service/risk.go`
1. In `UserService.Update`, capture `oldRoleID := user.RoleID` before applying the update; after `mapper.UpdateReqToEntity`, if `user.RoleID != oldRoleID`, log `Action: "users.change_role"` with metadata including `old_role_id`/`new_role_id`; otherwise log `Action: "users.update"` as before (single log call, action chosen conditionally).
2. In `RoleService.UpdateStatus`, change `Action: "roles.update"` to `Action: "roles.status_update"`.
3. In `risk.go`, remove the `"auth.forgot_password_spam": {}` entry from `highRisk` (confirmed no rate-limiter hook or any code path emits this action string; `rate_limiter_fixed.go`/`rate_rules.go` only implement generic request throttling with no audit-log integration for spam detection).

**Cluster 9 — Dashboard CMSActions/FinanceAction mixup**
**File**: `backend/internal/module/dashboard/service/dashboard_service.go`
1. Change `res.CMSActions = summary.FinanceAction` to `res.CMSActions = summary.CMSAction`.

**Non-PBT structural item — mapper folder rename**
**Directory**: `backend/internal/module/activitylog/maper/` → `mapper/`
1. Rename the directory (via file move, which updates the two import paths in `activity_log_repo.go` and `activity_log_service.go` automatically).

**Non-PBT item — Documentation drift**
**File**: `docs-final/api/activity_logs.jsonc`
1. Change the `list.parameters.query` keys `entity_type`/`risk_level` to `entity`/`risk`.
2. Add `"finance_action": 0` to the `summary` response example (chosen over removing the field, since Cluster 8 keeps `finance_action` as a reserved-for-future-finance-module field rather than removing it — add a `notes` entry clarifying this).

## Testing Strategy

### Validation Approach

The testing strategy follows a two-phase approach per cluster: first, write and run exploration tests on the unfixed code to surface counterexamples confirming each defect, then implement the fix and re-run both the exploration test (now expected to pass) and a preservation test (observed on unfixed code first, then re-verified after the fix).

No new test dependencies are added. Backend property-based tests use Go's native fuzz testing (`func FuzzXxx(f *testing.F)`, part of the standard `testing` package since Go 1.18, zero new `go.mod` entries — this repo targets Go 1.25). Frontend property-based tests use Node's built-in `node:test` and `node:assert` modules with manual randomized case generation, against small pure logic helpers extracted only where needed for testability (no new `package.json` dependencies).

### Exploratory Bug Condition Checking

**Goal**: Surface counterexamples that demonstrate each of the 9 PBT clusters BEFORE implementing any fix (the mapper rename and documentation drift are structural/non-behavioral and are verified separately by build and review, not by exploration tests).

**Test Plan**: For each cluster, write a test (Go `_test.go` file or Node `node:test` file) that exercises the unfixed code path and asserts the expected-after-fix behavior. Run on unfixed code; each should fail (or, for frontend clusters with no harness today, be confirmed by direct code reading since no test runner currently exists for `.jsx` files).

**Test Cases** (one per cluster, matching Properties 1/3/5/7/9/11/13/15/17):
1. Frontend field/pagination mapping — fails against real backend-shaped fixture (Property 1)
2. Backend summary action typo — fails to count seeded `auth.login_failed` rows (Property 3)
3. entity_type taxonomy — auth/sliders call sites emit inconsistent values (Property 5)
4. Role filter dropdown value — submitted value doesn't match slug (Property 7)
5. Bulk role action code — bulk and single-item codes collide (Property 9)
6. Handler double-wrap — Detail/EntityLogs body has nested `data.data` (Property 11)
7. Timestamp timezone — formatted string misrepresents a known non-UTC instant (Property 13)
8. Risk classification gap — role change logs `users.update` not `users.change_role` (Property 15)
9. Dashboard field mixup — `CMSActions` reads `FinanceAction` (Property 17)

**Expected Counterexamples**: each test above fails on unfixed code with the exact mismatch described in Bug Analysis §Current Behavior of `bugfix.md`.

### Fix Checking

**Goal**: Verify that for all inputs where each cluster's bug condition holds, the fixed function produces the expected behavior.

**Pseudocode:**
```
FOR ALL cluster IN [1..9] DO
  FOR ALL input WHERE isBugCondition(input, cluster) DO
    result := fixedFunction(input)
    ASSERT expectedBehavior(result, cluster)
  END FOR
END FOR
```

### Preservation Checking

**Goal**: Verify that for all inputs where each cluster's bug condition does NOT hold, the fixed function produces the same result as the original function.

**Pseudocode:**
```
FOR ALL cluster IN [1..9] DO
  FOR ALL input WHERE NOT isBugCondition(input, cluster) DO
    ASSERT originalFunction(input) = fixedFunction(input)
  END FOR
END FOR
```

**Testing Approach**: Property-based testing (via Go's native fuzz testing, and manual randomized generation with `node:test` in JS) is used for clusters 1, 3, 7, and 9, where the input space is naturally a range of values (arbitrary field values, arbitrary counts, arbitrary timestamps, arbitrary summary numbers). Clusters 2, 4, 5, 6, 8 are deterministic single-literal defects with a small enumerable set of call sites; these are scoped to concrete-case tests (per "Scoped PBT Approach" for deterministic bugs) covering every call site enumerated in Bug Details.

**Test Plan**: Observe behavior on unfixed code for non-bug-condition inputs (e.g. already-correct `berita`/`kegiatan` entity_type writes, already-correct `List`/`Summary` wrapping, already-correct `users.delete` risk classification), record it, then write tests asserting that exact observed behavior, and confirm they pass both before and after the fix.

**Test Cases**:
1. **Non-affected field preservation**: `id`/`action`/`description` continue to map correctly (Property 2)
2. **Other summary fields preservation**: `total_logs`/`high_risk` computed identically (Property 4)
3. **Other modules' entity_type preservation**: berita/kegiatan/kontak/pengurus/settings/users/roles unaffected (Property 6)
4. **Other dropdown preservation**: entity/risk dropdown values unchanged (Property 8)
5. **Single-item role action preservation**: `DeleteRole`/`RestoreRole` codes unchanged (Property 10)
6. **List/Summary wrapping preservation**: shapes unchanged (Property 12)
7. **Non-timestamp field / genuine-UTC preservation** (Property 14)
8. **Non-role-changing update / other risk classification preservation** (Property 16)
9. **Other dashboard stats preservation** (Property 18)

### Unit Tests

- Go: one `_test.go` per backend cluster (2, 3, 5, 6, 7, 8, 9-dashboard) covering the concrete call sites enumerated in Bug Details.
- Node (`node:test`): pure-function tests for the frontend row-mapping and pagination-extraction logic (Cluster 1) and for the dropdown option lists (Clusters 3, 4), extracted as small pure helpers where the current inline JSX logic needs a testable seam.
- Manual/build-verified only (no automated test, per the nature of the change): mapper directory rename (verified by `go build`/`go vet`), documentation drift (verified by review against `activity_log_query.go` and `activity_log_response.go`).

### Property-Based Tests

- Property 1/2 (Go-free, Node `node:test` with manual randomized cases): generate N random `ActivityLogItemRes`-shaped objects and N random `meta` objects; assert the fixed mapping function extracts every field/pagination value correctly for all of them, and that `id`/`action`/`description` are unaffected.
- Property 3/4 (Go native fuzzing, `FuzzGetSummaryFailedLoginCount`): fuzz random non-negative counts of `auth.login_failed` vs. other-action rows; assert `summary.failed_login`/`summary.cms_action` match, and `total_logs`/`high_risk` are unaffected by the literal change.
- Property 13/14 (Go native fuzzing, `FuzzTimestampFormatting`): fuzz random `time.Time`-constructing inputs (unix seconds + fixed-offset minutes, including offset 0 for genuine UTC); assert the fixed formatter's output round-trips through a standard RFC3339 parser to the same instant.
- Property 17/18 (Go native fuzzing, `FuzzDashboardCMSActionsField`): fuzz random distinct `CMSAction`/`FinanceAction` int64 values; assert `res.CMSActions` equals `CMSAction`, and the other 10 dashboard fields are untouched by the change.

### Integration Tests

- Backend: `go build ./...` and `go vet ./...` across the whole module after all fixes, to catch any import-path breakage from the mapper directory rename and any signature mismatches from Cluster 8's conditional logging change.
- Frontend: `npm run build` (Vite build) after all `.jsx`/`.js` fixes, to catch any syntax/reference errors; clean up the generated `dist/` output afterward since it's only needed to prove the build succeeds.
- No live database or manual click-through is performed as part of this spec's automated verification — this is explicitly noted as unverified in the final summary, consistent with the user's stated verification expectations.
