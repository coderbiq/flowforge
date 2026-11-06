import test from 'node:test';
import assert from 'node:assert/strict';
import { resolve } from 'node:path';
import { pathToFileURL } from 'node:url';
assert.ok(process.env.B1_ROOT, 'B1_ROOT must name the fixed fixture source directory');
const { help } = await import(pathToFileURL(resolve(process.env.B1_ROOT, 'help.mjs')));
test('help text spells Quote correctly', () => {
  assert.equal(help, 'Quote an order.');
});
