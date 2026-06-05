'use strict';

const path = require('path');
const fs = require('fs');

const INTERFACE_PATH = 'src/cli/scripts/lib/backends/interface.js';
const BEADS_PATH = 'src/cli/scripts/lib/backends/beads.js';

const REQUIRED_METHODS = [
  'getCapabilities', 'checkAvailability',
  'init', 'teardown',
  'addTask', 'addTasks', 'discoverTask', 'cancelTask',
  'claimTask', 'completeTask', 'blockTask', 'unclaimTask', 'reopenTask',
  'addDependency', 'removeDependency',
  'getReadyTasks', 'getStatus', 'getBlockedTasks', 'getTask', 'isAllDone',
  'exportSnapshot',
  'migrateFromYaml', 'cleanupOrphans',
];

function run(root) {
  let passed = 0;
  let failed = 0;
  const errors = [];

  const interfaceContent = fs.readFileSync(path.join(root, INTERFACE_PATH), 'utf8');
  const beadsContent = fs.readFileSync(path.join(root, BEADS_PATH), 'utf8');

  if (!interfaceContent.includes('TaskBackend')) {
    failed++; errors.push('interface.js: TaskBackend class not exported');
  } else { passed++; }

  if (!interfaceContent.includes('BackendCapabilities')) {
    failed++; errors.push('interface.js: BackendCapabilities class not exported');
  } else { passed++; }

  if (!interfaceContent.includes('TASK_STATUS')) {
    failed++; errors.push('interface.js: TASK_STATUS constants not exported');
  } else { passed++; }

  if (!interfaceContent.includes('TASK_TYPE')) {
    failed++; errors.push('interface.js: TASK_TYPE constants not exported');
  } else { passed++; }

  for (const method of REQUIRED_METHODS) {
    const inIface = interfaceContent.includes('async ' + method + '(') || interfaceContent.includes(method + '(');
    const inBeads = beadsContent.includes('async ' + method + '(') || beadsContent.includes(method + '(');

    if (!inIface) { failed++; errors.push('interface.js: method ' + method + ' not defined'); }
    else { passed++; }

    if (!inBeads) { failed++; errors.push('beads.js: method ' + method + ' not implemented'); }
    else { passed++; }
  }

  if (!beadsContent.includes('extends TaskBackend')) {
    failed++; errors.push('beads.js: BeadsBackend must extend TaskBackend');
  } else { passed++; }

  if (!beadsContent.includes('super(projectRoot)') && !beadsContent.includes('super(this._projectRoot)')) {
    failed++; errors.push('beads.js: BeadsBackend constructor must call super');
  } else { passed++; }

  const ifaceSigs = interfaceContent.match(/async\s+\w+\s*\([^)]*\)/g) || [];
  let hasProposalDir = false;
  for (const sig of ifaceSigs) {
    if (sig.includes('proposalDir') && !sig.includes('proposalId')) {
      failed++;
      errors.push('interface.js: signature uses proposalDir without proposalId: ' + sig.trim());
      hasProposalDir = true;
    }
  }
  if (!hasProposalDir) { passed++; }

  return { passed, failed, errors };
}

module.exports = { run };
