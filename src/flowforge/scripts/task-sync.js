#!/usr/bin/env node
'use strict';

const { execSync } = require('child_process');
const { loadMainConfig, findProposalDir } = require('./lib/config');
const { createAdapter } = require('./lib/adapters');

const projectRoot = process.argv[2] || process.cwd();
const proposalId = process.argv[3];
let direction = 'merge';
let checkOnly = false;

for (let i = 4; i < process.argv.length; i++) {
  if (process.argv[i] === '--from') {
    const v = process.argv[i + 1];
    if (v === 'yaml') direction = 'yaml-to-beads';
    else if (v === 'beads') direction = 'beads-to-yaml';
    i++;
  } else if (process.argv[i] === '--check') {
    checkOnly = true;
  }
}

if (!proposalId) {
  console.error('Usage: task-sync.js <projectRoot> <proposalId> [--from yaml|beads] [--check]');
  console.error('  --from yaml   以 task-map.yaml 为源，覆盖 beads');
  console.error('  --from beads  以 beads 为源，更新 task-map.yaml 状态');
  console.error('  默认           智能合并：yaml 管定义，beads 管状态');
  console.error('  --check        只检查不修改，报告不一致的内容');
  process.exit(1);
}

const config = loadMainConfig(projectRoot);
if (!config) {
  console.error('ERROR: .flowforge/config.yaml not found');
  process.exit(1);
}

const adapter = createAdapter(config, projectRoot);

async function main() {
  const proposalDir = findProposalDir(projectRoot, config, proposalId);
  if (!proposalDir) {
    console.log(JSON.stringify({ synced: false, error: `proposal ${proposalId} not found` }));
    process.exit(0);
  }

  if (checkOnly) {
    const result = await runCheck(proposalDir);
    console.log(JSON.stringify(result));
    if (result.issues && result.issues.length > 0) process.exit(1);
  } else {
    const result = await adapter.sync(proposalDir, direction);
    console.log(JSON.stringify(result));
  }
}

async function runCheck(proposalDir) {
  const adapterType = config.taskBackend?.adapter || 'yaml';
  if (adapterType !== 'beads') {
    return { consistent: true, issues: [], note: 'yaml backend, no sync needed' };
  }

  const yaml = require('./vendor/js-yaml');
  const fs = require('fs');
  const path = require('path');
  const taskMapPath = path.join(proposalDir, 'task-map.yaml');

  if (!fs.existsSync(taskMapPath)) {
    return { consistent: true, issues: [], note: 'no task-map.yaml found' };
  }

  const data = yaml.load(fs.readFileSync(taskMapPath, 'utf8'));
  const yamlTasks = data.tasks || [];
  const issues = [];

  let beadTasks = [];
  try {
    const proposalId = data.proposal_id;
    if (proposalId) {
      const result = execSync(`bd query spec=${proposalId} --json`, {
        cwd: projectRoot, encoding: 'utf8', stdio: 'pipe', timeout: 5000
      });
      beadTasks = JSON.parse(result);
      if (!Array.isArray(beadTasks)) beadTasks = [];
    }
  } catch (_) {
    return { consistent: true, issues: [], note: 'beads unavailable, skipped check' };
  }

  const beadById = {};
  for (const bt of beadTasks) beadById[bt.id] = bt;

  for (const yt of yamlTasks) {
    if (yt._beadId && !beadById[yt._beadId]) {
      issues.push(`[${yt.id}] ${yt.title}: beads issue ${yt._beadId} missing`);
    } else if (!yt._beadId) {
      issues.push(`[${yt.id}] ${yt.title}: 未关联 beads issue`);
    } else if (yt._beadId && beadById[yt._beadId]) {
      const bt = beadById[yt._beadId];
      const btStatus = _mapBeadStatus(bt.status);
      if (btStatus && btStatus !== yt.status) {
        issues.push(`[${yt.id}] ${yt.title}: yaml=${yt.status} beads=${btStatus}`);
      }
    }
  }

  for (const bt of beadTasks) {
    const match = yamlTasks.find(yt => yt._beadId === bt.id);
    if (!match && bt.issue_type !== 'epic') {
      issues.push(`[beads:${bt.id}] ${bt.title || 'unknown'}: 在 task-map.yaml 中无对应任务`);
    }
  }

  return { consistent: issues.length === 0, issues };
}

function _mapBeadStatus(beadStatus) {
  const mapping = {
    'open': 'pending', 'pending': 'pending',
    'in_progress': 'in_progress', 'in progress': 'in_progress',
    'done': 'done', 'closed': 'done', 'completed': 'done',
    'blocked': 'blocked', 'cancelled': 'cancelled',
  };
  return mapping[beadStatus] || null;
}

main().catch(e => {
  console.error(e.message);
  process.exit(1);
});
