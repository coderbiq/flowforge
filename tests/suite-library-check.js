'use strict';

const path = require('path');
const fs = require('fs');
const os = require('os');
const { execSync } = require('child_process');

const SCRIPTS_DIR = 'src/cli/scripts';

function run(root) {
  let passed = 0;
  let failed = 0;
  const errors = [];

  const checkScript = path.join(root, SCRIPTS_DIR, 'library-check.js');
  if (!fs.existsSync(checkScript)) {
    return { passed: 0, failed: 1, errors: ['library-check.js not found'] };
  }

  const tmpDir = fs.mkdtempSync(path.join(os.tmpdir(), 'ff-test-lib-'));
  const libRoot = path.join(tmpDir, 'ff-wiki', 'library');
  fs.mkdirSync(path.join(libRoot, 'architecture'), { recursive: true });
  fs.mkdirSync(path.join(libRoot, 'conventions'), { recursive: true });

  const now = new Date();
  const oldDate = new Date(now.getTime() - 200 * 86400000);

  fs.writeFileSync(path.join(libRoot, 'architecture', 'fresh.md'), `---
doc_type: architecture
title: Fresh Doc
status: active
updated: "${now.toISOString()}"
review_interval: 180
domain:
  scope: system
  type: design
  importance: should
  maturity: growing
---
# Fresh Doc
Content here.
`);

  fs.writeFileSync(path.join(libRoot, 'architecture', 'stale.md'), `---
doc_type: architecture
title: Stale Doc
status: active
updated: "${oldDate.toISOString()}"
review_interval: 180
domain:
  scope: system
  type: design
  importance: should
  maturity: growing
---
# Stale Doc
Old content.
`);

  fs.writeFileSync(path.join(libRoot, 'conventions', 'with-broken-ref.md'), `---
doc_type: convention
title: Broken Ref Doc
status: active
enforcement: should
updated: "${now.toISOString()}"
domain:
  scope: system
  type: convention
  importance: should
  maturity: growing
related:
  - ref: "../architecture/nonexistent.md"
  - ref: "../architecture/fresh.md"
---
# Broken Ref
See [nonexistent](../architecture/nonexistent.md)
`);

  function runCheck(flags) {
    try {
      return JSON.parse(execSync(`node "${checkScript}" "${tmpDir}" ${flags}`, {
        encoding: 'utf8', stdio: 'pipe', timeout: 10000
      }).trim());
    } catch (e) {
      return { error: e.message };
    }
  }

  const stalenessResult = runCheck('--staleness');
  if (stalenessResult.staleness && stalenessResult.staleness.length === 1 && stalenessResult.staleness[0].path.includes('stale')) {
    passed++;
  } else {
    failed++;
    errors.push(`staleness check: expected 1 stale doc, got ${JSON.stringify(stalenessResult.staleness)}`);
  }

  if (stalenessResult.staleness && stalenessResult.staleness[0] && stalenessResult.staleness[0].daysSince >= 200) {
    passed++;
  } else {
    failed++;
    errors.push('staleness check: daysSince should be >= 200');
  }

  const brokenResult = runCheck('--broken-refs');
  if (brokenResult.brokenRefs && brokenResult.brokenRefs.length === 1) {
    passed++;
  } else {
    failed++;
    errors.push(`broken ref check: expected 1 broken ref, got ${JSON.stringify(brokenResult.brokenRefs)}`);
  }

  if (brokenResult.brokenRefs && brokenResult.brokenRefs[0] && brokenResult.brokenRefs[0].brokenRef.includes('nonexistent')) {
    passed++;
  } else {
    failed++;
    errors.push('broken ref check: should detect nonexistent.md');
  }

  const dupeResult = runCheck('--duplicates');
  passed++;

  const orphanResult = runCheck('--orphans');
  if (orphanResult.orphans && orphanResult.orphans.length >= 2) {
    passed++;
  } else {
    failed++;
    errors.push('orphan check: expected >=2 orphans (no incoming refs)');
  }

  const fallbackResult = runCheck('');
  if (fallbackResult.staleness !== undefined) {
    passed++;
  } else {
    failed++;
    errors.push('default check: should run staleness when no flags');
  }

  fs.rmSync(tmpDir, { recursive: true, force: true });

  return { passed, failed, errors };
}

module.exports = { run };
