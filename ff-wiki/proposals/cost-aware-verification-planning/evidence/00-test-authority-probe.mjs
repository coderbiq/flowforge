// Preset by the planning owner. Exercises the real PI loader/runner; no tool writes execute.
import assert from 'node:assert/strict';
import { mkdtempSync, mkdirSync, writeFileSync, readFileSync, rmSync, symlinkSync, realpathSync, existsSync, lstatSync } from 'node:fs';
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
// Repair cases aggregate assertion failures only; loader/setup errors still abort.
// The original expect helper and all 38 original observations below stay unchanged.
async function runRepairMatrix() {
  assert.equal(checks, 38, 'all 38 original observations must pass before repair cases');
  console.log(`HOST: Node ${process.version}; Pi ${JSON.parse(readFileSync(join(piRoot, 'package.json'), 'utf8')).version}; source ${source}`);
  const failures = [];
  const groups = {};
  async function observe(group, runtime, target, blocked, label, toolName) {
    const count = groups[group] ??= { total: 0, passed: 0, failed: 0 };
    count.total++;
    try {
      await expect(runtime, target, blocked, `03/${group}/${label}/${toolName}`, toolName);
      count.passed++;
    } catch (error) {
      if (!(error instanceof assert.AssertionError)) throw error;
      count.failed++;
      const failure = { group, label, toolName, target, expected: error.expected, actual: error.actual, message: error.message };
      failures.push(failure);
      console.error(`FAIL: ${JSON.stringify(failure)}`);
    }
  }
  const tools = ['write', 'edit'];
  for (const component of ['final', 'parent']) {
    for (const location of ['inside', 'outside']) {
      const name = `dangling-${component}-${location}`;
      const runtime = await start(join(scratch, name));
      const destinationRoot = location === 'inside' ? runtime.project : join(scratch, `${name}-outside`);
      const destination = join(destinationRoot, 'uncreated', component === 'final' ? 'target_test.go' : 'parent');
      const link = join(runtime.project, component === 'final' ? 'dangling_test.go' : 'dangling-parent');
      symlinkSync(destination, link);
      assert.ok(lstatSync(link).isSymbolicLink(), 'fixture must contain a real dangling symlink');
      assert.equal(existsSync(destination), false, 'symlink destination must not exist');
      const relativeTarget = component === 'final' ? 'dangling_test.go' : 'dangling-parent/leaf_test.go';
      for (const hasGrant of [false, true]) {
        if (hasGrant) grant(runtime, { paths: [relativeTarget] });
        for (const form of ['relative', 'absolute']) {
          const target = form === 'relative' ? relativeTarget : join(runtime.project, relativeTarget);
          for (const tool of tools) {
            await observe('dangling', runtime, target, true, `${component}/${location}/${hasGrant ? 'grant' : 'no-grant'}/${form}`, tool);
          }
        }
      }
      assert.equal(existsSync(destination), false, 'emitting events must not create symlink destinations');
    }
  }

  for (const initiallyExists of [true, false]) {
    for (const hasGrant of [false, true]) {
      const name = `zone-${initiallyExists ? 'existing' : 'missing'}-${hasGrant ? 'grant' : 'no-grant'}`;
      const project = join(scratch, name);
      const physicalParent = join(scratch, `${name}-outside`, 'parent');
      const physicalRoot = join(physicalParent, 'config');
      const physicalZone = join(physicalRoot, 'test-author-grants');
      mkdirSync(project);
      if (initiallyExists) mkdirSync(physicalRoot, { recursive: true });
      symlinkSync(physicalRoot, join(project, '.flowforge'));
      // These aliases can themselves be dangling. No tool implementation is invoked.
      symlinkSync(physicalZone, join(project, 'zone-alias'));
      symlinkSync(physicalRoot, join(project, 'root-alias'));
      symlinkSync(physicalParent, join(project, 'parent-alias'));
      const runtime = await start(project);
      assert.equal(existsSync(physicalRoot), initiallyExists, 'verify pre-grant physical fixture state');
      if (hasGrant) {
        // mkdir through a dangling link cannot create its destination; fixture setup does so explicitly.
        mkdirSync(physicalRoot, { recursive: true });
        grant(runtime);
      }
      assert.equal(existsSync(runtime.status.grantPath), hasGrant, 'verify actual grant presence');
      assert.equal(existsSync(physicalRoot), initiallyExists || hasGrant, 'verify live physical fixture state');
      const record = runtime.status.sessionId + '.json';
      const targets = [
        ['lexical-record', runtime.status.grantPath],
        ['lexical-missing-record', '.flowforge/test-author-grants/new.json'],
        ['lexical-zone', '.flowforge/test-author-grants'],
        ['lexical-parent', '.flowforge'],
        ['lexical-project-ancestor', project],
        ['physical-record', join(physicalZone, record)],
        ['physical-missing-record', join(physicalZone, 'new.json')],
        ['physical-zone', physicalZone],
        ['physical-parent', physicalRoot],
        ['physical-grandparent', physicalParent],
        ['alias-record', 'zone-alias/' + record],
        ['alias-missing-record', 'zone-alias/new.json'],
        ['alias-zone', 'zone-alias'],
        ['root-alias-record', 'root-alias/test-author-grants/' + record],
        ['root-alias-zone', 'root-alias/test-author-grants'],
        ['root-alias-parent', 'root-alias'],
        ['parent-alias-record', 'parent-alias/config/test-author-grants/' + record],
        ['parent-alias-ancestor', 'parent-alias'],
      ];
      const originalRecord = hasGrant ? readFileSync(runtime.status.grantPath, 'utf8') : undefined;
      for (const [label, target] of targets) {
        for (const tool of tools) await observe('grant-zone', runtime, target, true, `${name}/${label}`, tool);
      }
      for (const tool of tools) {
        await observe('grant-zone', runtime, 'source.go', false, `${name}/ordinary-source`, tool);
        await observe('grant-zone', runtime, 'allowed_test.go', !hasGrant, `${name}/exact-missing-file`, tool);
      }
      assert.equal(existsSync(join(physicalZone, 'new.json')), false, 'no synthetic tool call may issue a record');
      assert.equal(existsSync(join(project, 'allowed_test.go')), false, 'no protected test write may execute');
      assert.equal(existsSync(physicalRoot), initiallyExists || hasGrant, 'events must not materialize missing parents');
      if (hasGrant) assert.equal(readFileSync(runtime.status.grantPath, 'utf8'), originalRecord, 'grant fixture must remain unchanged');
    }
  }

  const clock = await start(join(scratch, 'utc-expiry'));
  const futureYear = new Date().getUTCFullYear() + 4;
  let nonLeapYear = futureYear;
  while (nonLeapYear % 4 === 0) nonLeapYear++;
  let leapYear = futureYear;
  while (leapYear % 4 !== 0 || (leapYear % 100 === 0 && leapYear % 400 !== 0)) leapYear++;
  const timestamps = [
    ['utc-seconds', `${futureYear}-06-15T12:34:56Z`, false],
    ['utc-milliseconds', `${futureYear}-06-15T12:34:56.123Z`, false],
    ['utc-leap-date', `${leapYear}-02-29T12:34:56Z`, false],
    ['expired-utc', '2000-01-01T00:00:00Z', true],
    ['year-only', `${futureYear}`, true],
    ['date-only', `${futureYear}-06-15`, true],
    ['missing-Z', `${futureYear}-06-15T12:34:56`, true],
    ['missing-Z-milliseconds', `${futureYear}-06-15T12:34:56.123`, true],
    ['non-UTC-offset', `${futureYear}-06-15T12:34:56+08:00`, true],
    ['zero-offset-without-Z', `${futureYear}-06-15T12:34:56+00:00`, true],
    ['missing-seconds', `${futureYear}-06-15T12:34Z`, true],
    ['impossible-Feb-30', `${futureYear}-02-30T12:34:56Z`, true],
    ['non-leap-Feb-29', `${nonLeapYear}-02-29T12:34:56Z`, true],
    ['impossible-Apr-31', `${futureYear}-04-31T12:34:56.123Z`, true],
    ['impossible-month-13', `${futureYear}-13-15T12:34:56Z`, true],
    ['impossible-day-zero', `${futureYear}-06-00T12:34:56Z`, true],
    ['normalized-hour-24', `${futureYear}-06-15T24:00:00Z`, true],
    ['impossible-minute-60', `${futureYear}-06-15T12:60:00Z`, true],
  ];
  for (const [label, expiresAt, blocked] of timestamps) {
    grant(clock, { expiresAt, paths: ['ordinary/new/allowed_test.go'] });
    for (const tool of tools) {
      await observe('UTC', clock, 'ordinary/new/allowed_test.go', blocked, label, tool);
    }
  }
  assert.equal(existsSync(join(clock.project, 'ordinary')), false, 'exact missing ordinary parents are never actually written');
  for (const [group, counts] of Object.entries(groups)) console.log(`GROUP ${group}: ${JSON.stringify(counts)}`);
  console.log(`SUMMARY: original=38/38; repair=${Object.values(groups).reduce((sum, g) => sum + g.total, 0)}; failures=${failures.length}; passed=${checks}; no built-in write/edit executed`);
  assert.equal(failures.length, 0, 'repair boundary observations must all pass (aggregated failures above)');
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
  await runRepairMatrix();
} finally {
  process.chdir(originalCwd);
  for (const key of envKeys) {
    if (savedEnv[key] === undefined) delete process.env[key]; else process.env[key] = savedEnv[key];
  }
  rmSync(scratch, { recursive: true, force: true });
}
