// Exploration test for activity-log-fixes bugfix spec, Cluster 1.
//
// Property 1: Bug Condition - Frontend Reads Wrong Response Field Names and
// Pagination Path.
//
// IMPORTANT: This test is written and run BEFORE the fix (task 3). It is
// expected to FAIL on the current (unfixed) mapping functions in
// activityLogMapping.js, because those functions read field names/paths
// that do not exist on the real backend response shape
// (ActivityLogItemRes / ActivityLogPaginationRes).
//
// Run with: node --test src/pages/admin/activityLogMapping.exploration.test.js

import test from 'node:test'
import assert from 'node:assert/strict'
import { mapLogRowForCsv, mapLogRowForDisplay, extractPaginationMeta } from './activityLogMapping.js'

// Generates a realistic, backend-shaped ActivityLogItemRes fixture with
// randomized-ish field values so the property holds "for any" such object,
// not just one hand-picked example.
function randomLogFixture(seed) {
  const actorNames = ['Budi Santoso', 'Administrator Utama', 'Siti Aminah', 'Joko Widodo']
  const actorRoles = ['super_admin', 'admin']
  const entityTypes = ['user', 'berita', 'kegiatan', 'sliders']
  const riskLevels = ['low', 'medium', 'high']
  const actions = ['users.create', 'berita.update', 'auth.login_failed']

  return {
    id: seed,
    actor_id: seed,
    actor_name: actorNames[seed % actorNames.length],
    actor_role: actorRoles[seed % actorRoles.length],
    action: actions[seed % actions.length],
    entity_type: entityTypes[seed % entityTypes.length],
    entity_id: seed + 100,
    entity_label: `Entity Label ${seed}`,
    risk_level: riskLevels[seed % riskLevels.length],
    description: `Deskripsi aktivitas nomor ${seed}`,
    ip_address: `192.168.1.${seed % 255}`,
    user_agent: `Mozilla/5.0 (fixture ${seed})`,
    metadata: {},
    created_at: `2026-07-${String((seed % 28) + 1).padStart(2, '0')}T13:35:00Z`
  }
}

function randomPaginationFixture(seed) {
  return {
    meta: {
      current_page: 1,
      limit: 15,
      total_data: seed * 7 + 3,
      total_pages: seed + 1
    }
  }
}

test('Property 1 (Bug Condition) - mapLogRowForDisplay should extract real field values for any fixture', () => {
  for (let seed = 1; seed <= 10; seed++) {
    const fixture = randomLogFixture(seed)
    const row = mapLogRowForDisplay(fixture)

    // These assertions encode the EXPECTED (post-fix) behavior.
    // They are expected to FAIL on the current buggy implementation,
    // which reads log.actor/log.role/log.entity/etc. instead of the
    // real actor_name/actor_role/entity_type/etc. fields.
    assert.equal(row.actor, fixture.actor_name, `actor should come from actor_name (seed=${seed})`)
    assert.equal(row.role, fixture.actor_role, `role should come from actor_role (seed=${seed})`)
    assert.equal(row.ip, fixture.ip_address, `ip should come from ip_address (seed=${seed})`)
    assert.equal(row.device, fixture.user_agent, `device should come from user_agent (seed=${seed})`)
    assert.equal(row.entity, fixture.entity_type, `entity should come from entity_type (seed=${seed})`)
    assert.equal(row.entityLabel, fixture.entity_label, `entityLabel should come from entity_label (seed=${seed})`)
    assert.equal(row.risk, fixture.risk_level, `risk should come from risk_level (seed=${seed})`)
  }
})

test('Property 1 (Bug Condition) - mapLogRowForCsv should extract real field values for any fixture', () => {
  for (let seed = 1; seed <= 10; seed++) {
    const fixture = randomLogFixture(seed)
    const [waktu, aktor, role, , , , ip, device, risk] = mapLogRowForCsv(fixture)

    assert.equal(waktu, fixture.created_at, `CSV Waktu column should come from created_at (seed=${seed})`)
    assert.equal(aktor, fixture.actor_name, `CSV Aktor column should come from actor_name (seed=${seed})`)
    assert.equal(role, fixture.actor_role, `CSV Role column should come from actor_role (seed=${seed})`)
    assert.equal(ip, fixture.ip_address, `CSV IP column should come from ip_address (seed=${seed})`)
    assert.equal(device, fixture.user_agent, `CSV Device column should come from user_agent (seed=${seed})`)
    assert.equal(risk, fixture.risk_level, `CSV Risiko column should come from risk_level (seed=${seed})`)
  }
})

test('Property 1 (Bug Condition) - extractPaginationMeta should read meta.total_data/meta.total_pages for any fixture', () => {
  for (let seed = 1; seed <= 10; seed++) {
    const fixture = randomPaginationFixture(seed)
    const result = extractPaginationMeta(fixture)

    assert.ok(result, `extractPaginationMeta should return a value when meta is present (seed=${seed})`)
    assert.equal(result.total, fixture.meta.total_data, `total should come from meta.total_data (seed=${seed})`)
    assert.equal(result.totalPages, fixture.meta.total_pages, `totalPages should come from meta.total_pages (seed=${seed})`)
  }
})
