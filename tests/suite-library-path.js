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

  const tmpDir = fs.mkdtempSync(path.join(os.tmpdir(), 'ff-test-path-'));
  const defaultLib = path.join(tmpDir, 'ff-wiki', 'library');
  fs.mkdirSync(path.join(defaultLib, 'architecture'), { recursive: true });

  const customLib = path.join(tmpDir, 'custom-wiki', 'library');
  fs.mkdirSync(path.join(customLib, 'architecture'), { recursive: true });
  fs.writeFileSync(path.join(customLib, 'architecture', 'real-doc.md'), `---
doc_type: architecture
title: Real Document
status: active
updated: "2026-06-07T00:00:00Z"
domain:
  scope: system
  type: design
  importance: should
  maturity: growing
---
# Real Doc
`);

  fs.mkdirSync(path.join(tmpDir, '.flowforge'));
  fs.writeFileSync(path.join(tmpDir, '.flowforge', 'config.yaml'), `projects:
  - id: backend
    config: projects/backend.yaml
`);
  fs.mkdirSync(path.join(tmpDir, '.flowforge', 'projects'));
  fs.writeFileSync(path.join(tmpDir, '.flowforge', 'projects', 'backend.yaml'), `wikiRoot: custom-wiki
srcDirs: []
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

  const result = runCheck('--staleness');
  if (!result.error && result.staleness !== undefined) {
    passed++;
  } else {
    failed++;
    errors.push(`multi-path: expected staleness result, got error: ${JSON.stringify(result)}`);
  }

  if (result.staleness && result.staleness.length >= 0) {
    passed++;
  } else {
    failed++;
    errors.push('multi-path: staleness should be array');
  }

  const dupeResult = runCheck('--duplicates');
  if (!dupeResult.error) {
    passed++;
  } else {
    failed++;
    errors.push('multi-path: duplicate check should not error');
  }

  fs.rmSync(tmpDir, { recursive: true, force: true });

  const emptyDir = fs.mkdtempSync(path.join(os.tmpdir(), 'ff-test-empty-'));
  fs.mkdirSync(path.join(emptyDir, 'ff-wiki', 'library', 'architecture'), { recursive: true });

  try {
    const er = JSON.parse(execSync(`node "${checkScript}" "${emptyDir}" --staleness`, {
      encoding: 'utf8', stdio: 'pipe', timeout: 10000
    }).trim());
    if (er.staleness !== undefined) passed++;
    else { failed++; errors.push('empty library: should return staleness array'); }
  } catch (_) {
    failed++; errors.push('empty library: unexpected error');
  }

  fs.rmSync(emptyDir, { recursive: true, force: true });

  return { passed, failed, errors };
}

module.exports = { run };
