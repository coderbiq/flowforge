import test from 'node:test';
import assert from 'node:assert/strict';
import { resolve } from 'node:path';
import { pathToFileURL } from 'node:url';
import { spawnSync } from 'node:child_process';

assert.ok(process.env.B1_ROOT, 'B1_ROOT must name the fixed fixture source directory');
const root = resolve(process.env.B1_ROOT);
const { validateOrder } = await import(pathToFileURL(resolve(root, 'contract.mjs')));
const { quote } = await import(pathToFileURL(resolve(root, 'quote.mjs')));
const order = { lines: [{ sku: 'B', quantity: 2 }, { sku: 'A', quantity: 1 }, { sku: 'B', quantity: 1 }] };
const expected = { items: [{ sku: 'B', quantity: 2, lineCents: 600 }, { sku: 'A', quantity: 1, lineCents: 125 }, { sku: 'B', quantity: 1, lineCents: 300 }], totalCents: 1025 };
const cli = input => spawnSync(process.execPath, [resolve(root, 'cli.mjs'), JSON.stringify(input)], { encoding: 'utf8' });

test('interface accepts the complete order and rejects malformed quantities', () => {
  assert.deepEqual(validateOrder(structuredClone(order)), order);
  for (const input of [null, {}, { sku: 'A' }, { lines: [] }, { lines: [null] }, { lines: [{ sku: 1, quantity: 1 }] }, ...[0, -1, 1.5, '2'].map(quantity => ({ lines: [{ sku: 'A', quantity }] }))]) {
    assert.throws(() => validateOrder(input), /^Error: INVALID_ORDER$/);
  }
});

test('implementation preserves order and duplicates and computes exact cents without mutation', () => {
  const input = structuredClone(order);
  assert.deepEqual(quote(input), expected);
  assert.deepEqual(input, order);
  assert.throws(() => quote({ lines: [{ sku: 'missing', quantity: 1 }] }), /^Error: UNKNOWN_SKU$/);
  assert.throws(() => quote({ lines: [{ sku: 'A', quantity: 0 }] }), /^Error: INVALID_ORDER$/);
});

test('caller emits the full quote, including ordered duplicate rows', () => {
  const run = cli(order);
  assert.equal(run.status, 0, run.stderr);
  assert.equal(run.stderr, '');
  assert.deepEqual(JSON.parse(run.stdout), expected);
});

test('caller failure emits no partial result and returns exit 2', () => {
  for (const [input, message] of [[{ lines: [{ sku: 'A', quantity: 1 }, { sku: 'missing', quantity: 1 }] }, 'UNKNOWN_SKU'], [{ lines: [{ sku: 'A', quantity: 0 }] }, 'INVALID_ORDER']]) {
    const run = cli(input);
    assert.equal(run.status, 2);
    assert.equal(run.stdout, '');
    assert.equal(run.stderr.trim(), message);
  }
});
