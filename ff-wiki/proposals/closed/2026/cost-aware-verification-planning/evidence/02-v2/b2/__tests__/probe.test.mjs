import test from 'node:test';
import assert from 'node:assert/strict';
import { readFileSync, existsSync } from 'node:fs';
import { resolve } from 'node:path';
assert.ok(process.env.B2_ROOT, 'B2_ROOT must name the fixed fixture source directory');
const root = resolve(process.env.B2_ROOT);
test('permit must be unused before the pre-use gate', () => {
  assert.equal(existsSync(resolve(root, 'used.marker')), false, 'IRREVERSIBLE_USE_ALREADY_OCCURRED');
});
test('remote permit schema must be 2 before irreversible use', () => {
  const capability = JSON.parse(readFileSync(resolve(root, 'capability.json'), 'utf8'));
  assert.equal(capability.permit, 'once-02');
  assert.equal(capability.schema, 2, 'PRE_USE_SCHEMA_MISMATCH');
});
