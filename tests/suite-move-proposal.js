'use strict';

const path = require('path');
const fs = require('fs');
const os = require('os');
const { execSync } = require('child_process');

const MOVE_SCRIPT = 'src/cli/scripts/move-proposal.js';
const REFRESH_SCRIPT = 'src/cli/scripts/refresh-index.js';

function setupConfig(tmp, wikiRoot) {
  const ffDir = path.join(tmp, '.flowforge');
  fs.mkdirSync(path.join(ffDir, 'projects'), { recursive: true });
  fs.writeFileSync(path.join(ffDir, 'config.yaml'),
    'projects:\n  - id: default\n    config: projects/default.yaml\n');
  fs.writeFileSync(path.join(ffDir, 'projects', 'default.yaml'),
    `wikiRoot: ${wikiRoot}\n`);
}

function run(root) {
  let passed = 0;
  let failed = 0;
  const errors = [];

  const moveScript = path.join(root, MOVE_SCRIPT);
  const refreshScript = path.join(root, REFRESH_SCRIPT);

  if (!fs.existsSync(moveScript)) {
    return { passed: 0, failed: 1, errors: [`Missing script: ${MOVE_SCRIPT}`] };
  }

  // ========================================================
  // Test 1: Directory name differs from proposalId
  // ========================================================
  {
    const tmp = fs.mkdtempSync(path.join(os.tmpdir(), 'ff-test-move-'));
    try {
      setupConfig(tmp, 'ff-wiki');
      const activeDir = path.join(tmp, 'ff-wiki', 'workspace', 'proposals', 'active');
      const completedDir = path.join(tmp, 'ff-wiki', 'workspace', 'proposals', 'completed');
      const dirName = 'CRTEST01-my-proposal-title';
      fs.mkdirSync(path.join(activeDir, dirName), { recursive: true });
      fs.writeFileSync(path.join(activeDir, dirName, 'meta.yaml'),
        'id: CRTEST01\n' +
        'title: Test Proposal\n' +
        'status: draft\n' +
        'project: default\n' +
        'created_at: 2026-01-01T00:00:00Z\n' +
        'updated_at: 2026-01-01T00:00:00Z\n');

      const output = execSync(
        `node "${moveScript}" "${tmp}" "CRTEST01"`,
        { encoding: 'utf8', timeout: 10000 }
      );
      const result = JSON.parse(output);

      if (fs.existsSync(path.join(completedDir, dirName))) {
        passed++;
      } else {
        failed++;
        errors.push('move-proposal: should move directory using actual dirName');
      }

      if (!fs.existsSync(path.join(activeDir, dirName))) {
        passed++;
      } else {
        failed++;
        errors.push('move-proposal: active dir should be removed after move');
      }

      const movedMeta = fs.readFileSync(path.join(completedDir, dirName, 'meta.yaml'), 'utf8');
      const newTs = movedMeta.match(/updated_at:\s*(.+)/);
      if (newTs && newTs[1] !== '2026-01-01T00:00:00Z') {
        passed++;
      } else {
        failed++;
        errors.push('move-proposal: updated_at should be refreshed');
      }

      if (!movedMeta.includes('status:')) {
        passed++;
      } else {
        failed++;
        errors.push('move-proposal: status line should be removed');
      }

      if (result.steps.some(s => s.step === 'update_meta') &&
          result.steps.some(s => s.step === 'move_directory')) {
        passed++;
      } else {
        failed++;
        errors.push('move-proposal: should report update_meta + move_directory');
      }
    } catch (e) {
      failed += 5;
      errors.push(`move-proposal test 1 crash: ${e.message}`);
    } finally {
      fs.rmSync(tmp, { recursive: true, force: true });
    }
  }

  // ========================================================
  // Test 2: Already in completed -> skip move
  // ========================================================
  {
    const tmp = fs.mkdtempSync(path.join(os.tmpdir(), 'ff-test-move2-'));
    try {
      setupConfig(tmp, 'ff-wiki');
      const completedDir = path.join(tmp, 'ff-wiki', 'workspace', 'proposals', 'completed');
      const dirName = 'CRTEST02-already-done';
      fs.mkdirSync(path.join(completedDir, dirName), { recursive: true });
      fs.writeFileSync(path.join(completedDir, dirName, 'meta.yaml'),
        'id: CRTEST02\n' +
        'title: Already Done\n' +
        'status: archived\n' +
        'project: default\n' +
        'created_at: 2026-01-01T00:00:00Z\n' +
        'updated_at: 2026-01-01T00:00:00Z\n');

      const output = execSync(
        `node "${moveScript}" "${tmp}" "CRTEST02"`,
        { encoding: 'utf8', timeout: 10000 }
      );
      const result = JSON.parse(output);

      if (fs.existsSync(path.join(completedDir, dirName))) {
        passed++;
      } else {
        failed++;
        errors.push('move-proposal: should keep dir in completed/');
      }

      const moveStep = result.steps.find(s => s.step === 'move_directory');
      if (moveStep && moveStep.status === 'skipped') {
        passed++;
      } else {
        failed++;
        errors.push('move-proposal: move_directory should be skipped');
      }

      const meta = fs.readFileSync(path.join(completedDir, dirName, 'meta.yaml'), 'utf8');
      if (!meta.includes('status:')) {
        passed++;
      } else {
        failed++;
        errors.push('move-proposal: should remove status even in completed/');
      }
    } catch (e) {
      failed += 3;
      errors.push(`move-proposal test 2 crash: ${e.message}`);
    } finally {
      fs.rmSync(tmp, { recursive: true, force: true });
    }
  }

  // ========================================================
  // Test 3: refresh-index groups by directory, no status col
  // ========================================================
  {
    const tmp = fs.mkdtempSync(path.join(os.tmpdir(), 'ff-test-idx-'));
    try {
      setupConfig(tmp, 'ff-wiki');
      const activeDir = path.join(tmp, 'ff-wiki', 'workspace', 'proposals', 'active');
      const completedDir = path.join(tmp, 'ff-wiki', 'workspace', 'proposals', 'completed');
      fs.mkdirSync(activeDir, { recursive: true });
      fs.mkdirSync(completedDir, { recursive: true });

      fs.mkdirSync(path.join(activeDir, 'CR10-active'));
      fs.writeFileSync(path.join(activeDir, 'CR10-active', 'meta.yaml'),
        'id: CR10\n' +
        'title: Active One\n' +
        'project: default\n' +
        'created_at: 2026-01-01T00:00:00Z\n' +
        'updated_at: 2026-01-01T00:00:00Z\n');

      fs.mkdirSync(path.join(completedDir, 'CR20-done'));
      fs.writeFileSync(path.join(completedDir, 'CR20-done', 'meta.yaml'),
        'id: CR20\n' +
        'title: Completed One\n' +
        'project: default\n' +
        'created_at: 2026-01-01T00:00:00Z\n' +
        'updated_at: 2026-01-01T00:00:00Z\n');

      execSync(`node "${refreshScript}" "${tmp}"`, { encoding: 'utf8', timeout: 10000 });
      const indexPath = path.join(tmp, 'ff-wiki', 'workspace', 'proposals', 'INDEX.md');
      const indexContent = fs.readFileSync(indexPath, 'utf8');

      if (!indexContent.includes('| 状态 |')) {
        passed++;
      } else {
        failed++;
        errors.push('refresh-index: INDEX.md should NOT have status column');
      }

      if (indexContent.includes('🟢 1') && indexContent.includes('📦 1')) {
        passed++;
      } else {
        failed++;
        errors.push('refresh-index: should show correct counts (1 active, 1 completed)');
      }

      if (indexContent.includes('CR10') && indexContent.includes('CR20')) {
        passed++;
      } else {
        failed++;
        errors.push('refresh-index: should include both proposals');
      }
    } catch (e) {
      failed += 3;
      errors.push(`refresh-index test crash: ${e.message}`);
    } finally {
      fs.rmSync(tmp, { recursive: true, force: true });
    }
  }

  return { passed, failed, errors };
}

module.exports = { run };
