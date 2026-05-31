#!/usr/bin/env node
'use strict';

const { loadMainConfig, findProposalDir } = require('./lib/config');
const { createAdapter } = require('./lib/adapters');

const projectRoot = process.argv[2] || process.cwd();
const proposalId = process.argv[3];
const taskId = process.argv[4];

if (!proposalId || !taskId) {
  console.error('Usage: task-claim.js <projectRoot> <proposalId> <taskId>');
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
    console.log(JSON.stringify({ claimed: false, conflict: `proposal ${proposalId} not found` }));
    process.exit(0);
  }

  const result = await adapter.claimTask(proposalDir, taskId);
  console.log(JSON.stringify(result));

  if (!result.claimed && result.conflict) {
    process.exit(1);
  }
}

main().catch(e => {
  console.error(e.message);
  process.exit(1);
});
