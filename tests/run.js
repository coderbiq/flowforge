#!/usr/bin/env node
'use strict';

const path = require('path');
const fs = require('fs');

const ROOT = path.resolve(__dirname, '..');
const suites = [
  { name: 'js-syntax', file: 'suite-js-syntax.js' },
  { name: 'js-runtime', file: 'suite-js-runtime.js' },
  { name: 'config', file: 'suite-config.js' },
  { name: 'skill-frontmatter', file: 'suite-skill-frontmatter.js' },
  { name: 'skill-descriptions', file: 'suite-skill-descriptions.js' },
  { name: 'cli-cross-ref', file: 'suite-cli-cross-ref.js' },
  { name: 'agents-md-cross-ref', file: 'suite-agents-md-cross-ref.js' },
  { name: 'context-output', file: 'suite-context-output.js' },
  { name: 'backend-interface', file: 'suite-backend-interface.js' },
  { name: 'schema-validation', file: 'suite-schema-validation.js' },
  { name: 'version-consistency', file: 'suite-version-consistency.js' },
];

let totalPassed = 0;
let totalFailed = 0;
const failures = [];

console.log('FlowForge Test Suite\n');
console.log('=' .repeat(60));

for (const suite of suites) {
  const suitePath = path.join(__dirname, suite.file);
  if (!fs.existsSync(suitePath)) {
    console.log(`\n[${suite.name}] SKIP: file not found`);
    continue;
  }

  try {
    const mod = require(suitePath);
    if (typeof mod.run !== 'function') {
      console.log(`\n[${suite.name}] SKIP: no run() exported`);
      continue;
    }

    const result = mod.run(ROOT);
    totalPassed += result.passed;
    totalFailed += result.failed;

    console.log(`\n[${suite.name}] ${result.passed} passed, ${result.failed} failed`);
    for (const err of (result.errors || [])) {
      console.log(`  ✗ ${err}`);
      failures.push({ suite: suite.name, message: err });
    }
  } catch (e) {
    totalFailed++;
    console.log(`\n[${suite.name}] CRASH: ${e.message}`);
    failures.push({ suite: suite.name, message: `CRASH: ${e.message}` });
  }
}

console.log('\n' + '='.repeat(60));
console.log(`\nTotal: ${totalPassed} passed, ${totalFailed} failed`);

if (failures.length > 0) {
  console.log('\nFailures:');
  for (const f of failures) {
    console.log(`  [${f.suite}] ${f.message}`);
  }
}

process.exit(totalFailed > 0 ? 1 : 0);
