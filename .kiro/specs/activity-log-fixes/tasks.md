# Implementation Plan

This batch covers 9 independently-testable clusters (Properties 1-18) plus 2 non-PBT structural items. Each cluster follows explore -> preserve -> implement -> verify. Work through clusters in order; each cluster's exploration/preservation/implementation/checkpoint group is self-contained.

---

## Cluster 1 — Frontend field mapping + pagination meta

- [x] 1. Write bug condition exploration test for frontend field/pagination mapping
  - **Property 1: Bug Condition** - Frontend Reads Wrong Response Field Names and Pagination Path
  - **IMPORTANT**: Write this test BEFORE implementing the fix
  - **GOAL**: Surface counterexamples showing the row-mapping and pagination-extraction logic reads nonexistent fields
  - Extract the row-mapping logic (CSV row builder + table cell values) and pagination-extraction logic (`fetchLogs`'s `setTotal`/`setTotalPages` calls) from `ActivityLogAdmin.jsx` into small pure functions if needed for testability (e.g. `mapLogRowForDisplay(log)`, `extractPaginationMeta(resData)`), OR test via a Node script that imports and exercises the existing inline logic if extraction is not feasible without behavior change
  - Using `node:test`, generate random `ActivityLogItemRes`-shaped fixture objects (arbitrary `actor_name`, `actor_role`, `entity_type`, `entity_label`, `ip_address`, `user_agent`, `risk_level`, `created_at` strings) and random `{meta: {total_data, total_pages}}` fixture objects
  - Assert that the mapped row contains the real field values (not undefined) and that extracted pagination equals `meta.total_data`/`meta.total_pages`
  - Run test against the UNFIXED code path (reading `log.time`/`log.actor`/etc. and `res.data.pagination.total`/etc.)
  - **EXPECTED OUTCOME**: Test FAILS — mapped fields are `undefined`/blank and pagination stays at defaults (0, 1)
  - Document the counterexample (e.g. "fixture with actor_name='Budi' produces log.actor=undefined, rendered as 'System'")
  - _Requirements: 1.1, 1.2_

- [x] 2. Write preservation property test for already-correct fields and untouched filters
  - **Property 2: Preservation** - Already-Correct Fields and Filters Unaffected
  - **IMPORTANT**: Follow observation-first methodology
  - Observe on UNFIXED code: `log.id`, `log.action`, `log.description` already map correctly from the response object (these field names match on both sides)
  - Observe on UNFIXED code: the query params object built in `fetchLogs` (`search`, `role`, `entity`, `risk`, `page`, `limit`) is passed to `activityLogService.list` unchanged
  - Write a `node:test` property test: for random fixture objects, `id`/`action`/`description` extraction is unaffected by the field-mapping fix, and the filter-params object shape sent to `activityLogService.list` is unchanged
  - Verify test PASSES on UNFIXED code
  - _Requirements: 3.1, 3.2, 3.9, 3.10_

- [ ] 3. Fix frontend field mapping and pagination meta in ActivityLogAdmin.jsx
  - [x] 3.1 Implement the fix
    - In `fetchLogs`, change `setTotal(res.data.pagination.total || 0)` / `setTotalPages(res.data.pagination.totalPages || 1)` to read `res.data.meta.total_data` / `res.data.meta.total_pages`
    - In `downloadCSV`, change the row-builder to read `log.created_at`, `log.actor_name`, `log.actor_role`, `log.entity_type`, `log.entity_label`, `log.ip_address`, `log.user_agent`, `log.risk_level`
    - In the table body (`logs.map`), update all `log.time`/`log.actor`/`log.role`/`log.entity`/`log.entityLabel`/`log.ip`/`log.device`/`log.risk` references to the real field names
    - Replace the `(log.time || '').split(', ')` date-split with logic that formats the `created_at` ISO string into a date part and a time part for display
    - _Bug_Condition: isBugCondition(X) where X.fieldRead IN {time, actor, role, entity, entityLabel, ip, device, risk} OR X.pathRead IN {pagination.total, pagination.totalPages}_
    - _Expected_Behavior: Property 1 (fixed field names, fixed pagination path)_
    - _Preservation: Property 2 (id/action/description and filter params unaffected)_
    - _Requirements: 2.1, 2.2_

  - [x] 3.2 Verify bug condition exploration test now passes
    - **Property 1: Expected Behavior** - Frontend Reads Wrong Response Field Names and Pagination Path
    - **IMPORTANT**: Re-run the SAME test from task 1 - do NOT write a new test
    - **EXPECTED OUTCOME**: Test PASSES (confirms all fields and pagination now read correctly)
    - _Requirements: 2.1, 2.2_

  - [x] 3.3 Verify preservation test still passes
    - **Property 2: Preservation** - Already-Correct Fields and Filters Unaffected
    - **IMPORTANT**: Re-run the SAME test from task 2 - do NOT write new tests
    - **EXPECTED OUTCOME**: Test PASSES (confirms no regression to id/action/description or filter params)
    - _Requirements: 3.1, 3.2, 3.9, 3.10_

---

## Cluster 2 — Backend GetSummary action-string typo

- [x] 4. Write bug condition exploration test for GetSummary action-string typo
  - **Property 3: Bug Condition** - Summary Action String Matches Real Emitted Value
  - **IMPORTANT**: Write this test BEFORE implementing the fix
  - **GOAL**: Surface the counterexample showing `GetSummary` never matches real `auth.login_failed` rows
  - **Scoped PBT Approach**: This is a deterministic string-literal bug (not a range bug), so scope to a concrete fixture: seed an in-memory/sqlite or mocked query builder (or, if a real test DB is unavailable, test at the query-construction level by asserting the exact SQL/`Where` argument passed) with N rows where `action = "auth.login_failed"` and M rows with other actions
  - Write a Go test in `backend/internal/module/activitylog/repository/activity_log_repo_test.go` that calls `GetSummary` (or, if DB access is unavailable in this environment, asserts on the literal string constants used in the `Where` calls directly via a small refactor-free string check) and asserts `summary.FailedLogin == N` and `summary.CMSAction == M`
  - Run test on UNFIXED code
  - **EXPECTED OUTCOME**: Test FAILS — `FailedLogin` is 0 regardless of N, and `CMSAction` incorrectly includes the N failed-login rows
  - Document the counterexample (e.g. "N=3 seeded auth.login_failed rows, GetSummary reports FailedLogin=0")
  - _Requirements: 1.3_

- [ ] 5. Write preservation property test for other summary fields
  - **Property 4: Preservation** - Other Summary Fields Unaffected
  - **IMPORTANT**: Follow observation-first methodology
  - Observe on UNFIXED code: `summary.TotalLogs` and `summary.HighRisk` are computed correctly today (they don't depend on the typo'd literal)
  - Using Go native fuzz testing (`FuzzGetSummaryFailedLoginCount` or a table test with randomized counts), generate random non-negative counts of total rows and high-risk rows; assert `TotalLogs`/`HighRisk` match expected counts regardless of the failed-login literal
  - Verify test PASSES on UNFIXED code
  - _Requirements: 3.2_

- [ ] 6. Fix GetSummary action-string typo
  - [ ] 6.1 Implement the fix
    - In `backend/internal/module/activitylog/repository/activity_log_repo.go`, change both `"auth.failed_login"` literals in `GetSummary` (the `FailedLogin` count query and the `CMSAction` exclusion query) to `"auth.login_failed"`
    - _Bug_Condition: isBugCondition(X) where X.query = "GetSummary" AND X.literal = "auth.failed_login"_
    - _Expected_Behavior: Property 3 (FailedLogin/CMSAction computed from the real action string)_
    - _Preservation: Property 4 (TotalLogs/HighRisk unaffected)_
    - _Requirements: 2.3_

  - [ ] 6.2 Verify bug condition exploration test now passes
    - **Property 3: Expected Behavior** - Summary Action String Matches Real Emitted Value
    - **IMPORTANT**: Re-run the SAME test from task 4 - do NOT write a new test
    - **EXPECTED OUTCOME**: Test PASSES
    - _Requirements: 2.3_

  - [ ] 6.3 Verify preservation test still passes
    - **Property 4: Preservation** - Other Summary Fields Unaffected
    - **IMPORTANT**: Re-run the SAME test from task 5 - do NOT write new tests
    - **EXPECTED OUTCOME**: Test PASSES
    - _Requirements: 3.2_

---

## Cluster 3 — entity_type taxonomy (auth "auth" stray values, sliders singular, dead "database" filter)

- [ ] 7. Write bug condition exploration test for entity_type taxonomy inconsistency
  - **Property 5: Bug Condition** - entity_type Taxonomy Is Consistent Per Module
  - **IMPORTANT**: Write this test BEFORE implementing the fix
  - **GOAL**: Surface counterexamples showing auth's 2 stray `"auth"` values and sliders' singular `"slider"` value
  - Write a Go test in `backend/internal/module/auth/service/auth_service_test.go` (using a fake/mock `ActivityLogService` that captures the `ActivityLogInput` passed to `Log`) that calls `Login` with an unrecognized email and calls `ResetPassword` with an invalid token, then asserts `EntityType == "user"` for both captured inputs
  - Write a Go test in `backend/internal/module/sliders/service/sliders_service_test.go` (same fake-capture approach) that calls `Create`/`Update`/`Delete`/`Restore`/`BulkDelete`/`BulkRestore`/`Reorder` and asserts `EntityType == "sliders"` for each
  - Write a `node:test` test asserting the entity filter dropdown options list in `ActivityLogAdmin.jsx` does not include `"database"`
  - Run tests on UNFIXED code
  - **EXPECTED OUTCOME**: Tests FAIL — auth's two call sites report `EntityType == "auth"`, sliders reports `EntityType == "slider"`, and the dropdown still includes `"database"`
  - Document counterexamples (exact call sites: `Login` email-not-found branch, `ResetPassword` invalid-token branch; all 7 sliders service methods)
  - _Requirements: 1.4_

- [ ] 8. Write preservation property test for other modules' entity_type values
  - **Property 6: Preservation** - Non-Auth, Non-Sliders entity_type Values Unaffected
  - **IMPORTANT**: Follow observation-first methodology
  - Observe on UNFIXED code: berita/kegiatan/kontak/pengurus/settings services log `EntityType` equal to their own module name (already consistent); `users.*` and `roles.*` actions log `EntityType: "users"`/`"roles"` (plural, already consistent); auth's OTHER 5 call sites (login success, wrong-password, inactive-account, refresh, forgot-password, change-password-failed, logout) already log `EntityType: "user"`
  - Write Go tests (fake-capture approach, same pattern as task 7) asserting these already-correct values for a representative call in each of berita/kegiatan/kontak/pengurus/settings/user/role services, and for auth's non-"auth"-labeled call sites
  - Verify tests PASS on UNFIXED code
  - _Requirements: 3.3, 3.4_

- [ ] 9. Fix entity_type taxonomy
  - [ ] 9.1 Implement the fix
    - In `backend/internal/module/auth/service/auth_service.go`, change `EntityType: "auth"` to `EntityType: "user"` in `Login`'s email-not-found branch and in `ResetPassword`'s invalid-token branch
    - In `backend/internal/module/sliders/service/sliders_service.go`, change all six `EntityType: "slider"` occurrences (Create, Update, Delete, Reorder, Restore, BulkDelete, BulkRestore) to `EntityType: "sliders"` — leave the `Action: "slider.*"` string prefixes unchanged
    - In `frontend/src/pages/admin/ActivityLogAdmin.jsx`, remove the `<option value="database">Database</option>` line from the entity filter dropdown
    - _Bug_Condition: isBugCondition(X) where X.entityType IN {"auth", "slider"} OR X.option = "database"_
    - _Expected_Behavior: Property 5 (consistent "user"/"sliders" values, no dead dropdown option)_
    - _Preservation: Property 6 (other modules' entity_type unaffected)_
    - _Requirements: 2.4_

  - [ ] 9.2 Verify bug condition exploration test now passes
    - **Property 5: Expected Behavior** - entity_type Taxonomy Is Consistent Per Module
    - **IMPORTANT**: Re-run the SAME tests from task 7 - do NOT write new tests
    - **EXPECTED OUTCOME**: Tests PASS
    - _Requirements: 2.4_

  - [ ] 9.3 Verify preservation test still passes
    - **Property 6: Preservation** - Non-Auth, Non-Sliders entity_type Values Unaffected
    - **IMPORTANT**: Re-run the SAME tests from task 8 - do NOT write new tests
    - **EXPECTED OUTCOME**: Tests PASS
    - _Requirements: 3.3, 3.4_

---

## Cluster 4 — Role filter dropdown values

- [ ] 10. Write bug condition exploration test for role filter dropdown values
  - **Property 7: Bug Condition** - Role Filter Dropdown Sends Matchable Slug Values
  - **IMPORTANT**: Write this test BEFORE implementing the fix
  - **GOAL**: Surface the counterexample showing the dropdown submits display text instead of slugs
  - Write a `node:test` test that reads the role filter `<select>` option values in `ActivityLogAdmin.jsx` (via a small extracted constant/array if needed for testability, e.g. `ROLE_FILTER_OPTIONS`) and asserts each option's `value` matches `/^[a-z_]+$/` (a slug pattern) and specifically equals `"super_admin"`/`"admin"`
  - Run test on UNFIXED code
  - **EXPECTED OUTCOME**: Test FAILS — option values are `"Super Admin"`/`"Admin"` (fail the slug pattern and the exact-match assertion)
  - Document counterexample ("option value 'Super Admin' does not match stored actor_role 'super_admin'")
  - _Requirements: 1.5_

- [ ] 11. Write preservation property test for other filter dropdowns
  - **Property 8: Preservation** - Other Filter Dropdowns Unaffected
  - **IMPORTANT**: Follow observation-first methodology
  - Observe on UNFIXED code: the entity filter dropdown values (`user`, `auth`, `berita`, `kegiatan`, `pengurus`, `sliders`, `kontak`, `settings` — excluding `database` per Cluster 3) and risk filter dropdown values (`high`, `medium`, `low`) are already lowercase slugs matching stored data
  - Write a `node:test` test asserting these entity/risk option values are unchanged
  - Verify test PASSES on UNFIXED code
  - _Requirements: 3.2_

- [ ] 12. Fix role filter dropdown values
  - [ ] 12.1 Implement the fix
    - In `frontend/src/pages/admin/ActivityLogAdmin.jsx`, change `<option value="Super Admin">Super Admin</option>` to `<option value="super_admin">Super Admin</option>` and `<option value="Admin">Admin</option>` to `<option value="admin">Admin</option>` (display label text stays the same, only the `value` attribute changes)
    - _Bug_Condition: isBugCondition(X) where X.optionValue IN {"Super Admin", "Admin"}_
    - _Expected_Behavior: Property 7 (option values are "super_admin"/"admin")_
    - _Preservation: Property 8 (entity/risk dropdowns unaffected)_
    - _Requirements: 2.5_

  - [ ] 12.2 Verify bug condition exploration test now passes
    - **Property 7: Expected Behavior** - Role Filter Dropdown Sends Matchable Slug Values
    - **IMPORTANT**: Re-run the SAME test from task 10 - do NOT write a new test
    - **EXPECTED OUTCOME**: Test PASSES
    - _Requirements: 2.5_

  - [ ] 12.3 Verify preservation test still passes
    - **Property 8: Preservation** - Other Filter Dropdowns Unaffected
    - **IMPORTANT**: Re-run the SAME test from task 11 - do NOT write new tests
    - **EXPECTED OUTCOME**: Test PASSES
    - _Requirements: 3.2_

---

## Cluster 5 — Role bulk action codes

- [ ] 13. Write bug condition exploration test for role bulk action codes
  - **Property 9: Bug Condition** - Bulk Role Actions Use Distinct Action Codes
  - **IMPORTANT**: Write this test BEFORE implementing the fix
  - **GOAL**: Surface the counterexample showing bulk and single-item actions collide
  - Write a Go test in `backend/internal/module/role/service/role_service_test.go` (fake-capture `ActivityLogService`) that calls `BulkDelete` with a non-empty ID list and asserts `Action == "roles.bulk_delete"`; and calls `BulkRestore` asserting `Action == "roles.bulk_restore"`
  - Run test on UNFIXED code
  - **EXPECTED OUTCOME**: Test FAILS — captured `Action` is `"roles.delete"`/`"roles.restore"` (identical to single-item codes)
  - Document counterexample ("BulkDelete([1,2,3]) logs Action='roles.delete', indistinguishable from DeleteRole(1)")
  - _Requirements: 1.6_

- [ ] 14. Write preservation property test for single-item role actions
  - **Property 10: Preservation** - Single-Item Role Actions Unaffected
  - **IMPORTANT**: Follow observation-first methodology
  - Observe on UNFIXED code: `DeleteRole(id)` logs `Action: "roles.delete"` and `RestoreRole(id)` logs `Action: "roles.restore"` — this is already correct and must not change
  - Write a Go test asserting these single-item action codes remain exactly as observed
  - Verify test PASSES on UNFIXED code
  - _Requirements: 3.8_

- [ ] 15. Fix role bulk action codes
  - [ ] 15.1 Implement the fix
    - In `backend/internal/module/role/service/role_service.go`, change `Action: "roles.delete"` to `Action: "roles.bulk_delete"` in `BulkDelete`, and `Action: "roles.restore"` to `Action: "roles.bulk_restore"` in `BulkRestore`
    - _Bug_Condition: isBugCondition(X) where X.action = "roles.delete" WHERE X.callSite = BulkDelete, OR X.action = "roles.restore" WHERE X.callSite = BulkRestore_
    - _Expected_Behavior: Property 9 (bulk actions use "roles.bulk_delete"/"roles.bulk_restore")_
    - _Preservation: Property 10 (single-item DeleteRole/RestoreRole unaffected)_
    - _Requirements: 2.6_

  - [ ] 15.2 Verify bug condition exploration test now passes
    - **Property 9: Expected Behavior** - Bulk Role Actions Use Distinct Action Codes
    - **IMPORTANT**: Re-run the SAME test from task 13 - do NOT write a new test
    - **EXPECTED OUTCOME**: Test PASSES
    - _Requirements: 2.6_

  - [ ] 15.3 Verify preservation test still passes
    - **Property 10: Preservation** - Single-Item Role Actions Unaffected
    - **IMPORTANT**: Re-run the SAME test from task 14 - do NOT write new tests
    - **EXPECTED OUTCOME**: Test PASSES
    - _Requirements: 3.8_

---

## Cluster 6 — Handler double-wrap (Detail, EntityLogs)

- [ ] 16. Write bug condition exploration test for handler double-wrap
  - **Property 11: Bug Condition** - Detail and EntityLogs Return Flat Response Bodies
  - **IMPORTANT**: Write this test BEFORE implementing the fix
  - **GOAL**: Surface the counterexample showing Detail/EntityLogs nest an extra "data" level
  - Write a Go test in `backend/internal/module/activitylog/handler/activity_log_handler_test.go` using `httptest` + a fake/mock `ActivityLogService` that returns a known `ActivityLogDetailRes`/`[]ActivityLogItemRes`, invoking `Detail` and `EntityLogs` via `gin`'s test context, then decoding the JSON response and asserting `body["data"]` is the log object/array directly (e.g. `body["data"].(map[string]any)["id"]` exists)
  - Run test on UNFIXED code
  - **EXPECTED OUTCOME**: Test FAILS — `body["data"]` is `{"data": {...}}`, so `body["data"].(map[string]any)["id"]` is absent (it's nested one level deeper at `body["data"].(map[string]any)["data"].(map[string]any)["id"]`)
  - Document counterexample ("Detail(id=1) response body is {\"data\":{\"data\":{\"id\":1,...}}} instead of {\"data\":{\"id\":1,...}}")
  - _Requirements: 1.7_

- [ ] 17. Write preservation property test for List/Summary wrapping
  - **Property 12: Preservation** - List and Summary Wrapping Unaffected
  - **IMPORTANT**: Follow observation-first methodology
  - Observe on UNFIXED code: `List` and `Summary` handlers already pass `res` directly to `helper.SuccessResponse`, producing a correctly flat `{"data": {...}}` body
  - Write a Go test (same `httptest` approach) asserting `List`/`Summary` response bodies are flat (not double-wrapped)
  - Verify test PASSES on UNFIXED code
  - _Requirements: 3.5_

- [ ] 18. Fix handler double-wrap
  - [ ] 18.1 Implement the fix
    - In `backend/internal/module/activitylog/handler/activity_log_handler.go`, change `Detail`'s call from `helper.SuccessResponse(c, http.StatusOK, "ACTIVITY_LOG_DETAIL", "Berhasil mengambil detail activity log.", gin.H{"data": res}, nil)` to pass `res` directly instead of `gin.H{"data": res}`
    - Change `EntityLogs`'s call similarly, passing `res` directly instead of `gin.H{"data": res}`
    - _Bug_Condition: isBugCondition(X) where X.handler IN {Detail, EntityLogs} AND X.responseData = gin.H{"data": res}_
    - _Expected_Behavior: Property 11 (flat response body)_
    - _Preservation: Property 12 (List/Summary unaffected)_
    - _Requirements: 2.7_

  - [ ] 18.2 Verify bug condition exploration test now passes
    - **Property 11: Expected Behavior** - Detail and EntityLogs Return Flat Response Bodies
    - **IMPORTANT**: Re-run the SAME test from task 16 - do NOT write a new test
    - **EXPECTED OUTCOME**: Test PASSES
    - _Requirements: 2.7_

  - [ ] 18.3 Verify preservation test still passes
    - **Property 12: Preservation** - List and Summary Wrapping Unaffected
    - **IMPORTANT**: Re-run the SAME test from task 17 - do NOT write new tests
    - **EXPECTED OUTCOME**: Test PASSES
    - _Requirements: 3.5_

---

## Cluster 7 — Timestamp UTC correctness (includes mapper directory rename)

- [ ] 19. Write bug condition exploration test for timestamp timezone mislabeling
  - **Property 13: Bug Condition** - Timestamp Formatting Reflects Real Instant
  - **IMPORTANT**: Write this test BEFORE implementing the fix
  - **GOAL**: Surface the counterexample showing a local-time value gets falsely labeled "Z" (UTC)
  - Write a Go native fuzz test `FuzzTimestampFormatting` in `backend/internal/module/activitylog/maper/model_mapper_test.go` (before the rename) that constructs `time.Time` values in a known non-UTC fixed-offset location (e.g. `time.FixedZone("WIB", 7*3600)`) from fuzzed unix-seconds input, calls `EntityToResponse`/`EntityToDetailResponse` with an entity carrying that timestamp, parses the returned `CreatedAt` string with `time.Parse(time.RFC3339, ...)`, and asserts the parsed instant equals the original instant (`.Equal(original)`)
  - Run test on UNFIXED code
  - **EXPECTED OUTCOME**: Test FAILS for any non-UTC fuzzed instant — the parsed-back instant is off by the zone's offset (e.g. 7 hours) because the literal "Z" falsely claims UTC
  - Document counterexample ("time.Time at 2026-07-30 13:35:00 +0700 formats to '2026-07-30T13:35:00Z', which parses back as 2026-07-30 13:35:00 UTC — 7 hours off from the real instant")
  - _Requirements: 1.8_

- [ ] 20. Write preservation property test for non-timestamp fields and genuine-UTC case
  - **Property 14: Preservation** - Non-Timestamp Fields and Existing UTC-Correct Cases Unaffected
  - **IMPORTANT**: Follow observation-first methodology
  - Observe on UNFIXED code: all fields other than `CreatedAt` (ID, ActorName, Action, EntityType, etc.) are copied through unchanged by `EntityToResponse`/`EntityToDetailResponse`; and for a `time.Time` value that IS genuinely UTC (`time.UTC` location), the current formatter happens to produce a correct result
  - Extend the `FuzzTimestampFormatting` test (or add a table test) asserting non-timestamp fields are unchanged, and that a genuinely-UTC input still round-trips correctly after the fix
  - Verify test PASSES on UNFIXED code (for the UTC-location case and for non-timestamp fields)
  - _Requirements: 3.1, 3.7_

- [ ] 21. Fix timestamp UTC correctness and rename mapper directory
  - [ ] 21.1 Rename the maper directory to mapper
    - Move `backend/internal/module/activitylog/maper/model_mapper.go` (and its test file from task 19/20) to `backend/internal/module/activitylog/mapper/model_mapper.go`, updating the two import paths in `activity_log_repo.go` and `activity_log_service.go` (this is the non-PBT structural fix bundled here since it touches the same file)
    - _Requirements: n/a (structural, no behavior change — verified by build)_

  - [ ] 21.2 Implement the timestamp fix
    - In the now-renamed `mapper/model_mapper.go`, change `e.CreatedAt.Format("2006-01-02T15:04:05Z")` to `e.CreatedAt.UTC().Format(time.RFC3339)` in both `EntityToResponse` and `EntityToDetailResponse`
    - Add `"time"` to the file's imports
    - _Bug_Condition: isBugCondition(X) where X.formatter USES literal "Z" AND X.dbLoc <> "UTC"_
    - _Expected_Behavior: Property 13 (formatted output round-trips to the same instant)_
    - _Preservation: Property 14 (non-timestamp fields and genuine-UTC case unaffected)_
    - _Requirements: 2.8_

  - [ ] 21.3 Verify bug condition exploration test now passes
    - **Property 13: Expected Behavior** - Timestamp Formatting Reflects Real Instant
    - **IMPORTANT**: Re-run the SAME test from task 19 (now at its new path) - do NOT write a new test
    - **EXPECTED OUTCOME**: Test PASSES for all fuzzed non-UTC instants
    - _Requirements: 2.8_

  - [ ] 21.4 Verify preservation test still passes
    - **Property 14: Preservation** - Non-Timestamp Fields and Existing UTC-Correct Cases Unaffected
    - **IMPORTANT**: Re-run the SAME test from task 20 - do NOT write new tests
    - **EXPECTED OUTCOME**: Test PASSES
    - _Requirements: 3.1, 3.7_

---

## Cluster 8 — Risk classification gaps

- [ ] 22. Write bug condition exploration test for risk classification gaps
  - **Property 15: Bug Condition** - Sensitive Actions Get Correct Risk Classification
  - **IMPORTANT**: Write this test BEFORE implementing the fix
  - **GOAL**: Surface counterexamples showing role changes and status updates are misclassified
  - Write a Go test in `backend/internal/module/user/service/user_service_test.go` (fake-capture `ActivityLogService`) that calls `Update` with a `req.RoleID` different from the existing user's `RoleID`, and asserts the captured `Action == "users.change_role"`
  - Write a Go test in `backend/internal/module/role/service/role_service_test.go` that calls `UpdateStatus` and asserts the captured `Action == "roles.status_update"`
  - Run tests on UNFIXED code
  - **EXPECTED OUTCOME**: Tests FAIL — both report `Action == "users.update"`/`"roles.update"` respectively (mediumRisk, not the intended highRisk/mediumRisk-but-distinct classification)
  - Document counterexamples ("Update(id, req with RoleID=2 when current RoleID=1) logs Action='users.update' instead of 'users.change_role'")
  - _Requirements: 1.10_

- [ ] 23. Write preservation property test for non-role-changing updates and other classifications
  - **Property 16: Preservation** - Non-Role-Changing Updates and Existing Classifications Unaffected
  - **IMPORTANT**: Follow observation-first methodology
  - Observe on UNFIXED code: `UserService.Update` with `req.RoleID` equal to the existing `RoleID` logs `Action: "users.update"`; risk classifications for untouched codes (`users.delete` = high, `berita.create` = medium, etc.) are already correct
  - Write Go tests asserting the non-role-changing `Update` call still logs `"users.update"`, and a determineRisk table test confirming untouched action codes keep their existing classification
  - Verify tests PASS on UNFIXED code
  - _Requirements: 3.8_

- [ ] 24. Fix risk classification gaps
  - [ ] 24.1 Implement the fix
    - In `backend/internal/module/user/service/user_service.go`'s `Update`, capture `oldRoleID := user.RoleID` before calling `mapper.UpdateReqToEntity`; after the update, log `Action: "users.change_role"` (with metadata including `old_role_id`/`new_role_id`) when `user.RoleID != oldRoleID`, otherwise log `Action: "users.update"` as before
    - In `backend/internal/module/role/service/role_service.go`'s `UpdateStatus`, change `Action: "roles.update"` to `Action: "roles.status_update"`
    - In `backend/internal/module/activitylog/service/risk.go`, remove the `"auth.forgot_password_spam": {}` entry from `highRisk` (confirmed zero call sites and no rate-limiter hook exists to emit it)
    - _Bug_Condition: isBugCondition(X) where X.action IN {"users.update" WHERE roleChanged=true, "roles.update" WHERE onlyStatusChanged=true} OR X.riskMapEntry="auth.forgot_password_spam"_
    - _Expected_Behavior: Property 15 (role changes -> users.change_role, status updates -> roles.status_update)_
    - _Preservation: Property 16 (non-role-changing updates and other classifications unaffected)_
    - _Requirements: 2.10_

  - [ ] 24.2 Verify bug condition exploration test now passes
    - **Property 15: Expected Behavior** - Sensitive Actions Get Correct Risk Classification
    - **IMPORTANT**: Re-run the SAME tests from task 22 - do NOT write new tests
    - **EXPECTED OUTCOME**: Tests PASS
    - _Requirements: 2.10_

  - [ ] 24.3 Verify preservation test still passes
    - **Property 16: Preservation** - Non-Role-Changing Updates and Existing Classifications Unaffected
    - **IMPORTANT**: Re-run the SAME tests from task 23 - do NOT write new tests
    - **EXPECTED OUTCOME**: Tests PASS
    - _Requirements: 3.8_

---

## Cluster 9 — Dashboard CMSActions/FinanceAction mixup

- [ ] 25. Write bug condition exploration test for dashboard field mixup
  - **Property 17: Bug Condition** - Dashboard CMS Actions Stat Reads the Correct Summary Field
  - **IMPORTANT**: Write this test BEFORE implementing the fix
  - **GOAL**: Surface the counterexample showing the dashboard reads the wrong summary field
  - Write a Go native fuzz test `FuzzDashboardCMSActionsField` in `backend/internal/module/dashboard/service/dashboard_service_test.go` (fake/mock `ActivityLogService.Summary` returning fuzzed distinct `CMSAction`/`FinanceAction` int64 values) that calls `GetSummary` and asserts `res.CMSActions == summary.CMSAction` (from the fuzzed fake)
  - Run test on UNFIXED code
  - **EXPECTED OUTCOME**: Test FAILS for any fuzzed case where `CMSAction != FinanceAction` — `res.CMSActions` equals the fuzzed `FinanceAction` value instead
  - Document counterexample ("Summary returns CMSAction=86, FinanceAction=0; dashboard reports CMSActions=0")
  - _Requirements: 1.12_

- [ ] 26. Write preservation property test for other dashboard stats
  - **Property 18: Preservation** - Other Dashboard Stats Unaffected
  - **IMPORTANT**: Follow observation-first methodology
  - Observe on UNFIXED code: `total_berita`, `total_kegiatan`, `total_pengurus`, `total_kontak`, `unread_kontak`, `total_admin`, `pending_admin`, `activity_logs`, `high_risk_logs`, `failed_login` are computed from independent DB counts / summary fields untouched by this bug
  - Extend the fuzz/table test to assert these 10 fields are unaffected by varying `CMSAction`/`FinanceAction` (i.e. they depend only on their own independent inputs)
  - Verify test PASSES on UNFIXED code
  - _Requirements: 3.11_

- [ ] 27. Fix dashboard field mixup
  - [ ] 27.1 Implement the fix
    - In `backend/internal/module/dashboard/service/dashboard_service.go`'s `GetSummary`, change `res.CMSActions = summary.FinanceAction` to `res.CMSActions = summary.CMSAction`
    - _Bug_Condition: isBugCondition(X) where X.source="DashboardService.GetSummary" AND X.fieldRead="FinanceAction" WHERE X.targetStat="CMSActions"_
    - _Expected_Behavior: Property 17 (CMSActions reads CMSAction)_
    - _Preservation: Property 18 (other 10 dashboard stats unaffected)_
    - _Requirements: 2.12_

  - [ ] 27.2 Verify bug condition exploration test now passes
    - **Property 17: Expected Behavior** - Dashboard CMS Actions Stat Reads the Correct Summary Field
    - **IMPORTANT**: Re-run the SAME test from task 25 - do NOT write a new test
    - **EXPECTED OUTCOME**: Test PASSES for all fuzzed cases
    - _Requirements: 2.12_

  - [ ] 27.3 Verify preservation test still passes
    - **Property 18: Preservation** - Other Dashboard Stats Unaffected
    - **IMPORTANT**: Re-run the SAME test from task 26 - do NOT write new tests
    - **EXPECTED OUTCOME**: Test PASSES
    - _Requirements: 3.11_

---

## Non-PBT structural item — Documentation drift

- [ ] 28. Fix docs-final/api/activity_logs.jsonc query parameter names and summary example
  - Change the `list.parameters.query` example keys `entity_type`/`risk_level` to `entity`/`risk` to match the real `ActivityLogQueryReq` form tags
  - Add a `"finance_action": 0` field to the `summary` response example, plus a short `notes` entry clarifying it is reserved for a future finance module (consistent with keeping the field per Cluster 8's decision not to remove it)
  - Verify by re-reading `backend/internal/module/activitylog/dto/activity_log_query.go` and `activity_log_response.go` to confirm the documented names now match exactly
  - _Requirements: 1.11, 2.11_

---

## Final Checkpoint

- [ ] 29. Checkpoint - Ensure all tests pass and builds are green
  - Run `go build ./...` in `backend/`
  - Run `go vet ./...` in `backend/`
  - Run all new Go tests (`go test ./...` scoped to touched packages, plus the fuzz tests in normal unit-test mode: `go test -run Fuzz...` or as part of the regular suite via their seed corpus)
  - Run all new `node:test` files for the frontend pure-logic extractions
  - Run `npm run build` in `frontend/` (Vite build) — clean up the generated `dist/` output afterward since it's only needed to prove the build succeeds
  - Confirm every exploration test now passes and every preservation test still passes
  - State explicitly what was verified (build/vet/test green) and what could NOT be verified (no live database or manual click-through in the running app, no visual confirmation of the rendered admin page) — flag these as outstanding manual verification for the user
  - Ask the user if any questions arise, and surface the two follow-up recommendations noted during design (the identical hardcoded-"Z" timestamp pattern in `backend/internal/module/auth/mapper/auth_mapper.go`, left out of scope per the user's explicit scope note; and the incomplete `00011_seed_roles.sql` migration file, which appears to be missing its closing INSERT rows and `-- +goose StatementEnd` marker — unrelated to activity log but noticed during verification)
