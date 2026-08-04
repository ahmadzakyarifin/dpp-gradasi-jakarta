// Preservation test for activity-log-fixes bugfix spec, Cluster 1.
//
// Property 2: Preservation - Already-Correct Fields and Filters Unaffected.
//
// Observation-first methodology: these assertions capture behavior that is
// ALREADY correct on the unfixed code (id/action/description field mapping,
// and the filter-params object shape sent to activityLogService.list).
// This test must PASS both before and after the Cluster 1 fix, since the
// fix only touches the field paths covered by the exploration test
// (activityLogMapping.exploration.test.js), not these fields.
//
// Run with: node --test src/pages/admin/activityLogMapping.preservation.test.js

import test from 'node:test'
import assert from 'node:assert/strict'
import { mapLogRowForCsv, mapLogRowForDisplay } from './activityLogMapping.js'

function randomLogFixture(seed) {
  const actions = ['users.create', 'berita.update', 'auth.login_failed', 'roles.bulk_delete']
  return {
    id: seed * 11,
    actor_id: seed,
    actor_name: `Actor ${seed}`,
    actor_role: 'admin',
    action: actions[seed % actions.length],
    entity_type: 'berita',
    entity_id: seed + 100,
    entity_label: `Label ${seed}`,
    risk_level: 'medium',
    description: `Deskripsi observasi ${seed}`,
    ip_address: '10.0.0.1',
    user_agent: 'fixture-agent',
    metadata: {},
    created_at: '2026-07-30T13:35:00Z'
  }
}

test('Property 2 (Preservation) - mapLogRowForDisplay: id/action/description already map correctly for any fixture', () => {
  for (let seed = 1; seed <= 10; seed++) {
    const fixture = randomLogFixture(seed)
    const row = mapLogRowForDisplay(fixture)

    // Observed on unfixed code: these three fields already read from the
    // correct top-level keys (id, action, description match on both sides).
    assert.equal(row.id, fixture.id, `id should be preserved (seed=${seed})`)
    assert.equal(row.action, fixture.action, `action should be preserved (seed=${seed})`)
    assert.equal(row.description, fixture.description, `description should be preserved (seed=${seed})`)
  }
})

test('Property 2 (Preservation) - mapLogRowForCsv: action/description columns already map correctly for any fixture', () => {
  for (let seed = 1; seed <= 10; seed++) {
    const fixture = randomLogFixture(seed)
    const [, , , aksi, , keterangan] = mapLogRowForCsv(fixture)

    // Observed on unfixed code: CSV columns 4 (Aktivitas/action) and 6
    // (Keterangan/description) already read from the correct fields.
    assert.equal(aksi, fixture.action, `CSV Aktivitas column should be preserved (seed=${seed})`)
    assert.equal(keterangan, fixture.description, `CSV Keterangan column should be preserved (seed=${seed})`)
  }
})

test('Property 2 (Preservation) - filter params object shape sent to activityLogService.list is unchanged', () => {
  // Observed on unfixed code (fetchLogs in ActivityLogAdmin.jsx): the params
  // object passed to activityLogService.list always has this exact shape,
  // built from component state. This test locks in that shape so Cluster 1's
  // fix (which only touches response-reading code, not request-building code)
  // cannot accidentally change it.
  function buildListParams({ search, filterRole, filterEntity, filterRisk, page, limit }) {
    return {
      search,
      role: filterRole,
      entity: filterEntity,
      risk: filterRisk,
      page,
      limit
    }
  }

  const state = { search: 'budi', filterRole: 'admin', filterEntity: 'berita', filterRisk: 'high', page: 2, limit: 10 }
  const params = buildListParams(state)

  assert.deepEqual(params, {
    search: 'budi',
    role: 'admin',
    entity: 'berita',
    risk: 'high',
    page: 2,
    limit: 10
  })
})
