package repository

import (
	"context"
	"testing"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/dto"
	"gorm.io/gorm"
	"gorm.io/gorm/utils/tests"
)

// Exploration test for activity-log-fixes bugfix spec, Cluster 2.
//
// Property 3: Bug Condition - Summary Action String Matches Real Emitted
// Value.
//
// IMPORTANT: written and run BEFORE the fix. Expected to FAIL on the
// current (unfixed) GetSummary, which filters on the literal
// "auth.failed_login" while every real call site in auth_service.go logs
// "auth.login_failed" (reversed word order).
//
// This test exercises the REAL GetSummary method (not a reimplementation)
// against gorm's built-in utils/tests.DummyDialector in DryRun mode, using
// a query callback hook to capture the exact SQL/args issued for each
// sub-query. No live database connection or new test dependency is
// required: DummyDialector ships inside the already-required gorm.io/gorm
// module.
//
// Run with: go test ./internal/module/activitylog/repository/... -run Exploration -v

// newDryRunRepoWithCapture builds a repository whose internal *gorm.DB is in
// DryRun mode and records the SQL/vars of every query issued through it.
func newDryRunRepoWithCapture(t *testing.T) (*activityLogRepository, *[]capturedQuery) {
	t.Helper()
	db, err := gorm.Open(tests.DummyDialector{}, &gorm.Config{DryRun: true})
	if err != nil {
		t.Fatalf("failed to open dummy dialector db: %v", err)
	}

	captured := &[]capturedQuery{}
	err = db.Callback().Query().After("gorm:query").Register("capture_test_queries", func(tx *gorm.DB) {
		*captured = append(*captured, capturedQuery{
			SQL:  tx.Statement.SQL.String(),
			Vars: append([]interface{}{}, tx.Statement.Vars...),
		})
	})
	if err != nil {
		t.Fatalf("failed to register capture callback: %v", err)
	}

	return &activityLogRepository{db: db}, captured
}

type capturedQuery struct {
	SQL  string
	Vars []interface{}
}

// containsVar reports whether any of the captured query's bound args equals want.
func (c capturedQuery) containsVar(want string) bool {
	for _, v := range c.Vars {
		if s, ok := v.(string); ok && s == want {
			return true
		}
	}
	return false
}

func TestExploration_GetSummary_UsesRealEmittedFailedLoginLiteral(t *testing.T) {
	repo, captured := newDryRunRepoWithCapture(t)

	// Real GetSummary call, exactly as the service layer invokes it.
	_, err := repo.GetSummary(context.Background(), &dto.ActivityLogQueryReq{})
	if err != nil {
		t.Fatalf("GetSummary returned unexpected error: %v", err)
	}

	const realEmittedAction = "auth.login_failed" // confirmed via auth_service.go's 5 Log call sites

	var failedLoginQuery, cmsActionQuery *capturedQuery
	for i := range *captured {
		q := (*captured)[i]
		for _, v := range q.Vars {
			if s, ok := v.(string); ok {
				if s == "auth.failed_login" {
					// Whichever of the two typo'd queries we see first,
					// classify by the operator in the SQL text.
					if containsSQL(q.SQL, "<>") {
						cmsActionQuery = &q
					} else {
						failedLoginQuery = &q
					}
				}
			}
		}
	}

	t.Logf("All captured queries during GetSummary: %+v", *captured)

	// EXPECTED OUTCOME on unfixed code: FAILS. failedLoginQuery should not
	// exist once fixed (no query should filter on "auth.failed_login" at
	// all); on unfixed code we currently find it using that wrong literal
	// instead of the real "auth.login_failed".
	if failedLoginQuery != nil {
		t.Errorf("counterexample: GetSummary's FailedLogin count query uses literal 'auth.failed_login' (SQL=%s VARS=%v) instead of the real emitted action string %q - summary.failed_login will always be 0", failedLoginQuery.SQL, failedLoginQuery.Vars, realEmittedAction)
	}
	if cmsActionQuery != nil {
		t.Errorf("counterexample: GetSummary's CMSAction exclusion query excludes literal 'auth.failed_login' (SQL=%s VARS=%v) instead of excluding the real emitted action string %q - summary.cms_action will wrongly include failed-login rows", cmsActionQuery.SQL, cmsActionQuery.Vars, realEmittedAction)
	}

	// Also assert the positive expectation directly: some captured query
	// should filter on the REAL action string. On unfixed code, none does,
	// so this fails too (belt-and-suspenders with the above).
	foundRealLiteral := false
	for _, q := range *captured {
		if q.containsVar(realEmittedAction) {
			foundRealLiteral = true
		}
	}
	if !foundRealLiteral {
		t.Errorf("counterexample: none of GetSummary's queries reference the real emitted action string %q; captured=%+v", realEmittedAction, *captured)
	}
}

func containsSQL(sql, substr string) bool {
	for i := 0; i+len(substr) <= len(sql); i++ {
		if sql[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
