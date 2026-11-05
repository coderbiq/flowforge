// Preset by the planning owner. Exercises the real PI loader/runner; no tool writes execute.
import assert from 'node:assert/strict';
import { mkdtempSync, mkdirSync, writeFileSync, readFileSync, rmSync, symlinkSync, realpathSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { resolve, join } from 'node:path';
import { randomUUID } from 'node:crypto';
import { pathToFileURL } from 'node:url';

const piRoot = process.argv[2];
assert.ok(piRoot, 'pass the installed PI package root');
const { loadExtensions } = await import(pathToFileURL(join(piRoot, 'dist/core/extensions/loader.js')));
const { ExtensionRunner } = await import(pathToFileURL(join(piRoot, 'dist/core/extensions/runner.js')));
const originalCwd = process.cwd();
const source = resolve(process.argv[3] || 'assets/pi/flowforge.ts');
const scratch = realpathSync(mkdtempSync(join(tmpdir(), 'flowforge-authority-probe-')));
const envKeys = ['PI_SUBAGENT_CHILD', 'PI_SUBAGENT_RUNNER_PID', 'PI_SUBAGENT_EXTENSION_BINDINGS'];
const savedEnv = Object.fromEntries(envKeys.map(k => [k, process.env[k]]));
let checks = 0;
const runtimes = [];

async function start(project, sessionId = randomUUID(), child = true) {
  mkdirSync(project, { recursive: true });
  process.chdir(project);
  if (child) {
    process.env.PI_SUBAGENT_CHILD = '1';
    process.env.PI_SUBAGENT_RUNNER_PID = String(process.pid);
  } else {
    delete process.env.PI_SUBAGENT_CHILD;
    delete process.env.PI_SUBAGENT_RUNNER_PID;
  }
  delete process.env.PI_SUBAGENT_EXTENSION_BINDINGS;
  const loaded = await loadExtensions([source], project);
  assert.deepEqual(loaded.errors, [], 'real extension must load');
  const runner = new ExtensionRunner(loaded.extensions, loaded.runtime, project, { getSessionId: () => sessionId }, undefined);
  const errors = [];
  runner.onError(error => errors.push(error));
  await runner.emit({ type: 'session_start', reason: 'startup' });
  assert.deepEqual(errors, [], 'session startup must not throw');
  const statusTool = loaded.extensions.flatMap(e => [...e.tools.values()]).find(t => t.definition.name === 'flowforge_test_authority');
  assert.ok(statusTool, 'missing read-only flowforge_test_authority entry');
  const state = await statusTool.definition.execute('status', {}, undefined, undefined, runner.createContext());
  assert.equal(state.isError, undefined);
  const status = state.details;
  assert.equal(status.sessionId, sessionId);
  assert.equal(status.projectRoot, realpathSync(project));
  assert.ok(typeof status.instanceId === 'string' && status.instanceId.length > 10);
  assert.equal(status.eligible, child);
  assert.equal(status.grantPath, join(realpathSync(project), '.flowforge/test-author-grants', sessionId + '.json'));
  const result = { runner, status, project, errors };
  runtimes.push(result);
  return result;
}
function grant(runtime, changes = {}) {
  const { status } = runtime;
  const value = {
    version: 1, sessionId: status.sessionId, instanceId: status.instanceId,
    projectRoot: status.projectRoot, role: 'plan', paths: ['allowed_test.go'],
    expiresAt: new Date(Date.now() + 3600000).toISOString(), ...changes,
  };
  mkdirSync(resolve(status.grantPath, '..'), { recursive: true });
  writeFileSync(status.grantPath, JSON.stringify(value));
}
async function expect(runtime, target, blocked, label, toolName = 'write') {
  const result = await runtime.runner.emitToolCall({
    type: 'tool_call', toolName, toolCallId: randomUUID(),
    input: { path: target, content: 'never executed', edits: [{ oldText: 'old', newText: 'new' }] },
  });
  assert.deepEqual(runtime.errors, [], 'guard must return denial rather than throw');
  assert.equal(Boolean(result?.block), blocked, label);
  checks++;
  console.log(`PASS ${checks}: ${label}`);
}
try {
  const a = await start(join(scratch, 'a'));
  await expect(a, 'allowed_test.go', true, 'default write denied');
  await expect(a, 'allowed_test.go', true, 'default edit denied', 'edit');
  await expect(a, 'source.go', false, 'ordinary source preserved');
  await expect(a, a.status.grantPath, true, 'child cannot create its own grant');
  await expect(a, '.flowforge/test-author-grants', true, 'child cannot overwrite grant directory');
  await expect(a, '.flowforge', true, 'child cannot overwrite grant ancestor');
  grant(a);
  await expect(a, 'allowed_test.go', false, 'authorized Plan write');
  await expect(a, 'allowed_test.go', false, 'authorized Plan edit', 'edit');
  await expect(a, join(a.project, 'allowed_test.go'), false, 'absolute exact target');
  await expect(a, 'other_test.go', true, 'unlisted protected target denied');
  await expect(a, '../allowed_test.go', true, 'outside project denied');
  await expect(a, a.status.grantPath, true, 'authorized child cannot enlarge grant');
  grant(a, { role: 'refine' });
  await expect(a, 'allowed_test.go', false, 'authorized Refine write');
  const invalid = [
    ['implementer role', { role: 'implement' }], ['unknown field', { extra: true }],
    ['empty paths', { paths: [] }], ['wildcard paths', { paths: ['**/*_test.go'] }],
    ['absolute grant path', { paths: [join(a.project, 'allowed_test.go')] }],
    ['traversal grant', { paths: ['../allowed_test.go'] }], ['dot segment', { paths: ['./allowed_test.go'] }],
    ['wrong instance', { instanceId: randomUUID() }], ['wrong session', { sessionId: randomUUID() }],
    ['wrong project', { projectRoot: scratch }], ['expired', { expiresAt: '2000-01-01T00:00:00Z' }],
    ['invalid expiry', { expiresAt: 'later' }], ['version mismatch', { version: 2 }],
  ];
  for (const [label, fields] of invalid) {
    grant(a, fields);
    await expect(a, 'allowed_test.go', true, label + ' denied');
  }
  writeFileSync(a.status.grantPath, '{invalid');
  await expect(a, 'allowed_test.go', true, 'malformed grant denied');
  grant(a);
  rmSync(a.status.grantPath);
  await expect(a, 'allowed_test.go', true, 'revocation immediate');
  const outside = join(scratch, 'outside');
  mkdirSync(outside);
  symlinkSync(outside, join(a.project, 'escape'));
  grant(a, { paths: ['escape/target_test.go'] });
  await expect(a, 'escape/target_test.go', true, 'symlink parent denied');
  writeFileSync(join(outside, 'target_test.go'), 'original');
  symlinkSync(join(outside, 'target_test.go'), join(a.project, 'allowed_test.go'));
  grant(a);
  await expect(a, 'allowed_test.go', true, 'symlink target denied');
  symlinkSync(join(outside, 'target_test.go'), join(a.project, 'alias.go'));
  await expect(a, 'alias.go', true, 'non-test alias cannot reach protected target');
  symlinkSync(join(a.project, '.flowforge/test-author-grants'), join(a.project, 'grant-alias'));
  await expect(a, 'grant-alias/new.json', true, 'grant directory alias denied');
  assert.equal(readFileSync(join(outside, 'target_test.go'), 'utf8'), 'original');
  rmSync(join(a.project, 'allowed_test.go'));
  grant(a);
  const b = await start(join(scratch, 'b'));
  grant(b);
  // A keeps its factory identity despite process-global env/cwd changes from B.
  delete process.env.PI_SUBAGENT_CHILD;
  delete process.env.PI_SUBAGENT_RUNNER_PID;
  await expect(a, 'allowed_test.go', false, 'parallel environment changes do not alter A');
  await expect(b, 'other_test.go', true, 'parallel B cannot inherit broader rights');
  await a.runner.emit({ type: 'session_shutdown', reason: 'quit' });
  await expect(a, 'allowed_test.go', true, 'shutdown revokes instance');
  const resumed = await start(a.project, a.status.sessionId);
  assert.notEqual(resumed.status.instanceId, a.status.instanceId);
  await expect(resumed, 'allowed_test.go', true, 'same session restart rejects old grant');
  const controller = await start(join(scratch, 'controller'), randomUUID(), false);
  await expect(controller, controller.status.grantPath, false, 'controller can issue grant records');
  grant(controller);
  await expect(controller, 'allowed_test.go', true, 'controller does not gain child-author exemption');
  console.log(`PASS: ${checks} guard observations; no built-in write/edit executed`);
} finally {
  process.chdir(originalCwd);
  for (const key of envKeys) {
    if (savedEnv[key] === undefined) delete process.env[key]; else process.env[key] = savedEnv[key];
  }
  rmSync(scratch, { recursive: true, force: true });
}
