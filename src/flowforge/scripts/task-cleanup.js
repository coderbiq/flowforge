#!/usr/bin/env node
'use strict';

const { loadMainConfig, findProposalDir } = require('./lib/config');
const { createAdapter } = require('./lib/adapters');

const projectRoot = process.argv[2] || process.cwd();
const proposalId = process.argv[3];

if (!proposalId) {
  console.error('Usage: task-cleanup.js <projectRoot> <proposalId>');
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
    console.log(JSON.stringify({ clean: true, note: `proposal ${proposalId} not found` }));
    process.exit(0);
  }

  // 先同步 beads 状态到 YAML（beads 中已关闭的任务可能没更新 YAML）
  const adapterType = config.taskBackend?.adapter || 'yaml';
  if (adapterType === 'beads') {
    try {
      const syncResult = await adapter.sync(proposalDir, 'beads-to-yaml');
      if (syncResult.summary && syncResult.summary.updated > 0) {
        console.log(JSON.stringify({ sync: `updated ${syncResult.summary.updated} tasks from beads` }));
      }
    } catch (_) { /* sync 失败不阻塞，继续 YAML 检查 */ }
  }

  const result = await adapter.cleanup(proposalDir);
  console.log(JSON.stringify(result));

  if (!result.clean) {
    process.exit(1);
  }
}

main().catch(e => {
  console.error(e.message);
  process.exit(1);
});
