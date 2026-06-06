'use strict';

const path = require('path');
const fs = require('fs');
const os = require('os');

// Copy the fixed findProposalById implementation for testing
function findProposalById(projectRoot, projects, id) {
  for (const p of projects) {
    const ws = path.join(projectRoot, p.wikiRoot, 'workspace');
    for (const sub of ['active', 'completed']) {
      const subDir = path.join(ws, 'proposals', sub);
      if (!fs.existsSync(subDir)) continue;
      const dirs = fs.readdirSync(subDir, { withFileTypes: true }).filter(d => d.isDirectory());
      for (const d of dirs) {
        if (d.name === id || d.name.startsWith(id + '-')) {
          return { proposalDir: path.join(subDir, d.name), projectId: p.id, wikiRoot: p.wikiRoot };
        }
      }
    }
  }
  return null;
}

function run(root) {
  let passed = 0;
  let failed = 0;
  const errors = [];

  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), 'ff-test-find-'));
  try {
    const wikiRoot = path.join(tmp, 'ff-wiki');
    const activeDir = path.join(wikiRoot, 'workspace', 'proposals', 'active');
    const completedDir = path.join(wikiRoot, 'workspace', 'proposals', 'completed');
    fs.mkdirSync(activeDir, { recursive: true });
    fs.mkdirSync(completedDir, { recursive: true });

    const projects = [{ id: 'default', wikiRoot: 'ff-wiki' }];

    // Test 1: exact match works
    fs.mkdirSync(path.join(activeDir, 'CR01'));
    const r1 = findProposalById(tmp, projects, 'CR01');
    if (r1 && r1.proposalDir.endsWith('CR01')) {
      passed++;
    } else {
      failed++;
      errors.push('findProposalById: exact match (CR01) should work');
    }

    // Test 2: prefix match works (CR-id finds CR-id-suffix)
    fs.mkdirSync(path.join(activeDir, 'CR02-my-proposal-title'));
    const r2 = findProposalById(tmp, projects, 'CR02');
    if (r2 && r2.proposalDir.endsWith('CR02-my-proposal-title')) {
      passed++;
    } else {
      failed++;
      errors.push('findProposalById: prefix match (CR02 → CR02-my-proposal-title) should work');
    }

    // Test 3: prefix match in completed/ works
    fs.mkdirSync(path.join(completedDir, 'CR03-archived-proposal'));
    const r3 = findProposalById(tmp, projects, 'CR03');
    if (r3 && r3.proposalDir.includes('/completed/')) {
      passed++;
    } else {
      failed++;
      errors.push('findProposalById: prefix match in completed/ should work');
    }

    // Test 4: non-existent ID returns null
    const r4 = findProposalById(tmp, projects, 'CR99');
    if (r4 === null) {
      passed++;
    } else {
      failed++;
      errors.push('findProposalById: non-existent CR99 should return null');
    }

    // Test 5: partial prefix match should NOT match (CR0 should not match CR01)
    fs.mkdirSync(path.join(activeDir, 'CRX-something'));
    const r5 = findProposalById(tmp, projects, 'CR');
    if (r5 === null) {
      passed++;
    } else {
      failed++;
      errors.push('findProposalById: partial match (CR) should NOT match CRX-something');
    }

    // Test 6: active/ takes priority when same ID in both dirs
    fs.mkdirSync(path.join(activeDir, 'CR06-both'));
    fs.mkdirSync(path.join(completedDir, 'CR06-both-old'));
    const r6 = findProposalById(tmp, projects, 'CR06');
    if (r6 && r6.proposalDir.includes('/active/')) {
      passed++;
    } else {
      failed++;
      errors.push('findProposalById: active/ should take priority over completed/');
    }

    // Test 7: exact match before prefix match (CR07-b should match exactly)
    fs.mkdirSync(path.join(activeDir, 'CR07-b-exact'));
    fs.mkdirSync(path.join(activeDir, 'CR07-b'));
    const r7 = findProposalById(tmp, projects, 'CR07-b');
    if (r7 && r7.proposalDir.endsWith('CR07-b')) {
      passed++;
    } else {
      failed++;
      errors.push('findProposalById: exact match should take priority over prefix match');
    }

  } catch (e) {
    failed++;
    errors.push(`findProposalById test crash: ${e.message}`);
  } finally {
    fs.rmSync(tmp, { recursive: true, force: true });
  }

  return { passed, failed, errors };
}

module.exports = { run };
