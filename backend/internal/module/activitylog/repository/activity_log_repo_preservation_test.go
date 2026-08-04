package repository

import (
	"context"
	"strings"
	"testing"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/dto"
)

// Preservation test for activity-log-fixes bugfix spec, Cluster 2.
//
// Property 4: Preservation - Other Summary Fields Unaffected.
//
// Observation-first methodology: on the UNFIXED code, summary.TotalLogs and
// summary.HighRisk are already computed correctly (they never reference
// the typo'd "auth.failed_login"/"auth.login_failed" literal at all). This
// test locks in that structural independence using Go's native fuzz
// testing over the query-filter fields, so it must keep passing after
// Cluster 2's fix (which only changes the FailedLogin/CMSAction literals).
//
// Run as a normal test with: go test ./internal/module/activitylog/repository/... -run Preservation -v
// Run as a fuzz target with:  go test ./internal/module/activitylog/repository/... -fuzz FuzzGetSummary_TotalLogsAndHighRisk_Unaffected

const failedLoginTypo = "auth.failed_login"
const failedLoginReal = "auth.login_failed"

// TestPreservation_GetSummary_TotalLogsAndHighRisk_NeverReferenceActionLiteral
// is the seed/regression form of the fuzz property below, run by default
// under `go test` (fuzz targets only fuzz under `-fuzz`, but always run
// their seed corpus as regular subtests otherwise).
func TestPreservation_GetSummary_TotalLogsAndHighRisk_NeverReferenceActionLiteral(t *testing.T) {
	cases := []string{"", "some search", "budi santoso", "%", "'; DROP TABLE x; --", "unicode-研究"}
	for _, search := range cases {
		assertTotalLogsAndHighRiskUnaffected(t, search)
	}
}

func FuzzGetSummary_TotalLogsAndHighRisk_Unaffected(f *testing.F) {
	for _, seed := range []string{"", "search term", "café", "100% match"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, search string) {
		assertTotalLogsAndHighRiskUnaffected(t, search)
	})
}

func assertTotalLogsAndHighRiskUnaffected(t *testing.T, search string) {
	repo, captured := newDryRunRepoWithCapture(t)

	_, err := repo.GetSummary(context.Background(), &dto.ActivityLogQueryReq{Search: search})
	if err != nil {
		t.Fatalf("GetSummary returned unexpected error for search=%q: %v", search, err)
	}

	// The FIRST captured query is always TotalLogs (r.countQuery(ctx, r.buildQuery(req))),
	// and the SECOND is always HighRisk (adds "risk_level = 'high'"), per the
	// fixed call order in GetSummary. Both must never reference either the
	// typo'd or the real failed-login literal, for any fuzzed search string.
	if len(*captured) < 2 {
		t.Fatalf("expected at least 2 captured queries (TotalLogs, HighRisk), got %d: %+v", len(*captured), *captured)
	}

	totalLogsQuery := (*captured)[0]
	highRiskQuery := (*captured)[1]

	for _, q := range []capturedQuery{totalLogsQuery, highRiskQuery} {
		if q.containsVar(failedLoginTypo) || q.containsVar(failedLoginReal) {
			t.Errorf("search=%q: TotalLogs/HighRisk query unexpectedly references a failed-login literal: SQL=%s VARS=%v (these two queries must stay independent of the Cluster 2 fix)", search, q.SQL, q.Vars)
		}
	}

	// HighRisk's query must always filter on risk_level = 'high', regardless
	// of the fuzzed search string, and this must remain true after the fix.
	if !strings.Contains(highRiskQuery.SQL, "risk_level = ?") || !highRiskQuery.containsVar("high") {
		t.Errorf("search=%q: HighRisk query does not filter on risk_level='high' as expected: SQL=%s VARS=%v", search, highRiskQuery.SQL, highRiskQuery.Vars)
	}
}
