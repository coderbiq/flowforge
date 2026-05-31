#!/usr/bin/env node
'use strict';

const { loadMainConfig, findProposalDir } = require('./lib/config');
const { createAdapter } = require('./lib/adapters');

const projectRoot = process.argv[2] || process.cwd();
const proposalId = process.argv[3];
const taskId = process.argv[4];
const summary = process.argv[5] || '';

if (!proposalId || !taskId) {
  console.error('Usage: task-done.js <projectRoot> <proposalId> <taskId> [summary]');
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
    console.log(JSON.stringify({ done: false, error: `proposal ${proposalId} not found` }));
    process.exit(0);
  }

  const result = await adapter.completeTask(proposalDir, taskId, summary);
  console.log(JSON.stringify(result));
}

main().catch(e => {
  console.error(e.message);
  process.exit(1);
});
