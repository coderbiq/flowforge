#!/usr/bin/env node
'use strict';

const { loadMainConfig, findProposalDir } = require('./lib/config');
const { createAdapter } = require('./lib/adapters');

const projectRoot = process.argv[2] || process.cwd();
const proposalId = process.argv[3];
const parentTaskId = process.argv[4];
const title = process.argv[5];
const description = process.argv[6] || '';

if (!proposalId || !parentTaskId || !title) {
  console.error('Usage: task-discover.js <projectRoot> <proposalId> <parentTaskId> <title> [description]');
  process.exit(1);
}

const config = loadMainConfig(projectRoot);
if (!config) {
  console.error('ERROR: .flowforge/config.yaml not found');
  process.exit(1);
}

const adapter = createAdapter(config, projectRoot);

async function main() {
  const caps = adapter.getCapabilities();
  if (!caps.discoveredFrom) {
    console.log(JSON.stringify({ created: false, reason: 'adapter does not support discoverTask' }));
    process.exit(0);
  }

  const proposalDir = findProposalDir(projectRoot, config, proposalId);
  if (!proposalDir) {
    console.log(JSON.stringify({ created: false, error: `proposal ${proposalId} not found` }));
    process.exit(0);
  }

  const result = await adapter.discoverTask(proposalDir, parentTaskId, { title, description });
  console.log(JSON.stringify(result));
}

main().catch(e => {
  console.error(e.message);
  process.exit(1);
});
