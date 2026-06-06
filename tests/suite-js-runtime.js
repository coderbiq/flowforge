'use strict';

const path = require('path');

const BEADS_PATH = 'src/cli/scripts/lib/backends/beads.js';
const INTERFACE_PATH = 'src/cli/scripts/lib/backends/interface.js';

function run(root) {
  let passed = 0;
  let failed = 0;
  const errors = [];

  let BeadsBackend, TaskBackend, BackendCapabilities, TASK_STATUS, TASK_TYPE;

  try {
    const beadsMod = require(path.join(root, BEADS_PATH));
    BeadsBackend = beadsMod.BeadsBackend;
    passed++;
  } catch (e) {
    failed++;
    errors.push(`Cannot require beads.js: ${e.message}`);
    return { passed, failed, errors };
  }

  try {
    const ifaceMod = require(path.join(root, INTERFACE_PATH));
    TaskBackend = ifaceMod.TaskBackend;
    BackendCapabilities = ifaceMod.BackendCapabilities;
    TASK_STATUS = ifaceMod.TASK_STATUS;
    TASK_TYPE = ifaceMod.TASK_TYPE;
    passed++;
  } catch (e) {
    failed++;
    errors.push(`Cannot require interface.js: ${e.message}`);
    return { passed, failed, errors };
  }

  const backend = new BeadsBackend('/tmp/nonexistent');

  if (backend instanceof TaskBackend) {
    passed++;
  } else {
    failed++;
    errors.push('BeadsBackend is not an instance of TaskBackend');
  }

  const caps = backend.getCapabilities();
  if (caps instanceof BackendCapabilities) {
    passed++;
  } else {
    failed++;
    errors.push('getCapabilities() does not return BackendCapabilities');
  }

  if (caps.atomicClaim === true) { passed++; } else { failed++; errors.push('atomicClaim should be true'); }
  if (caps.dependencySort === true) { passed++; } else { failed++; errors.push('dependencySort should be true'); }
  if (caps.auditTrail === true) { passed++; } else { failed++; errors.push('auditTrail should be true'); }

  const availability = backend.checkAvailability();
  if (typeof availability.then === 'function') {
    passed++;
  } else {
    failed++;
    errors.push('checkAvailability() must return Promise');
  }

  const methods = {
    init: ['proposalId', 'title'],
    hasTaskSpace: ['proposalId'],
    addTask: ['proposalId', 'task'],
    addTasks: ['proposalId', 'tasks'],
    claimTask: ['proposalId', 'taskId'],
    completeTask: ['proposalId', 'taskId'],
    getStatus: ['proposalId'],
    getReadyTasks: ['proposalId'],
    exportSnapshot: ['proposalId'],
  };

  for (const [method, expectedParams] of Object.entries(methods)) {
    if (typeof backend[method] !== 'function') {
      failed++;
      errors.push(`BeadsBackend.${method} is not a function`);
      continue;
    }
    passed++;

    const fnStr = backend[method].toString();
    const params = fnStr.match(/\(([^)]*)\)/);
    if (params) {
      const paramList = params[1].split(',').map(p => p.trim()).filter(Boolean);
      for (const expected of expectedParams) {
        if (!paramList.some(p => p === expected)) {
          failed++;
          errors.push(`BeadsBackend.${method} missing parameter "${expected}"`);
        } else {
          passed++;
        }
      }
    }
  }

  return { passed, failed, errors };
}

module.exports = { run };
