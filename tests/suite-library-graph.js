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

  const graphScript = path.join(root, SCRIPTS_DIR, 'library-graph.js');
  if (!fs.existsSync(graphScript)) {
    return { passed: 0, failed: 1, errors: ['library-graph.js not found'] };
  }

  const tmpDir = fs.mkdtempSync(path.join(os.tmpdir(), 'ff-test-graph-'));
  const libRoot = path.join(tmpDir, 'ff-wiki', 'library');
  fs.mkdirSync(path.join(libRoot, 'architecture'), { recursive: true });
  fs.mkdirSync(path.join(libRoot, 'decisions'), { recursive: true });

  fs.writeFileSync(path.join(libRoot, 'architecture', 'arch-doc.md'), `---
doc_type: architecture
title: Architecture Overview
status: active
domain:
  scope: system
  type: design
  importance: should
  maturity: stable
related:
  - ref: "../decisions/decision-doc.md"
---
# Architecture
References the decision doc.
`);

  fs.writeFileSync(path.join(libRoot, 'decisions', 'decision-doc.md'), `---
doc_type: decision
title: Key Decision
status: active
decision_status: accepted
domain:
  scope: system
  type: decision
  importance: should
  maturity: stable
related:
  - ref: "../architecture/arch-doc.md"
---
# Decision
References the architecture doc.
`);

  fs.writeFileSync(path.join(libRoot, 'architecture', 'orphan-doc.md'), `---
doc_type: architecture
title: Orphan Document
status: active
domain:
  scope: system
  type: design
  importance: info
  maturity: seed
---
# Orphan
No one references this document.
`);

  function runGraph(args) {
    try {
      return JSON.parse(execSync(`node "${graphScript}" "${tmpDir}" ${args}`, {
        encoding: 'utf8', stdio: 'pipe', timeout: 10000
      }).trim());
    } catch (e) {
      return { error: e.message };
    }
  }

  const defaultResult = runGraph('');
  if (defaultResult.stats && defaultResult.stats.nodes === 3) {
    passed++;
  } else {
    failed++;
    errors.push(`graph stats nodes: expected 3, got ${defaultResult.stats?.nodes}`);
  }

  const backlinksResult = runGraph('backlinks decisions/decision-doc.md');
  if (backlinksResult.backlinks && backlinksResult.backlinks.length === 1) {
    passed++;
  } else {
    failed++;
    errors.push('backlinks: expected 1 backlink from arch-doc');
  }

  const refsResult = runGraph('refs architecture/arch-doc.md');
  if (refsResult.refs && refsResult.refs.length === 1) {
    passed++;
  } else {
    failed++;
    errors.push('refs: expected 1 ref to decision-doc');
  }

  const orphansResult = runGraph('orphans');
  if (orphansResult.orphans && orphansResult.orphans.length === 1) {
    passed++;
  } else {
    failed++;
    errors.push(`orphans: expected 1 orphan doc, got ${orphansResult.orphans?.length}`);
  }

  const hubsResult = runGraph('hubs 5');
  if (hubsResult.hubs && hubsResult.hubs.length >= 1) {
    passed++;
  } else {
    failed++;
    errors.push('hubs: expected at least 1 hub');
  }

  const blastResult = runGraph('blast-radius decisions/decision-doc.md 2');
  if (blastResult.affected !== undefined) {
    passed++;
  } else {
    failed++;
    errors.push('blast-radius: missing affected field');
  }

  fs.rmSync(tmpDir, { recursive: true, force: true });

  return { passed, failed, errors };
}

module.exports = { run };
