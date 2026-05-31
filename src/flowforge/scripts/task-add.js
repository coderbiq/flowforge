#!/usr/bin/env node
'use strict';

const { loadMainConfig, findProposalDir } = require('./lib/config');
const { createAdapter } = require('./lib/adapters');

const projectRoot = process.argv[2] || process.cwd();
const proposalId = process.argv[3];
const title = process.argv[4];
const description = process.argv[5] || '';
const depIds = process.argv.slice(6);

if (!proposalId || !title) {
  console.error('Usage: task-add.js <projectRoot> <proposalId> <title> <description> [depId1 depId2 ...]');
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
    console.log(JSON.stringify({ added: false, error: `proposal ${proposalId} not found` }));
    process.exit(0);
  }

  const result = await adapter.addTask(proposalDir, {
    title,
    description,
    dependencies: depIds
  });
  console.log(JSON.stringify(result));
}

main().catch(e => {
  console.error(e.message);
  process.exit(1);
});
