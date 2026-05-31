#!/usr/bin/env node
'use strict';

const { loadMainConfig, findProposalDir } = require('./lib/config');
const { createAdapter } = require('./lib/adapters');

const projectRoot = process.argv[2] || process.cwd();
const proposalId = process.argv[3];

if (!proposalId) {
  console.error('Usage: task-create.js <projectRoot> <proposalId>');
  process.exit(1);
}

const config = loadMainConfig(projectRoot);
if (!config) {
  console.error('ERROR: .flowforge/config.yaml not found');
  process.exit(1);
}

const adapter = createAdapter(config, projectRoot);

async function main() {
  const avail = await adapter.checkAvailability(projectRoot);

  const proposalDir = findProposalDir(projectRoot, config, proposalId);
  if (!proposalDir) {
    console.log(JSON.stringify({ error: `proposal ${proposalId} not found`, created: 0 }));
    process.exit(0);
  }

  const result = await adapter.createFromTaskMap(proposalDir);
  result.adapter = config.taskBackend?.adapter || 'yaml';
  result.adapterAvailable = avail.available;

  console.log(JSON.stringify(result));
}

main().catch(e => {
  console.error(e.message);
  process.exit(1);
});
