#!/usr/bin/env node
'use strict';

const { loadMainConfig, findProposalDir } = require('./lib/config');
const { createAdapter } = require('./lib/adapters');

const projectRoot = process.argv[2] || process.cwd();
const proposalId = process.argv[3];

if (!proposalId) {
  console.error('Usage: task-context.js <projectRoot> <proposalId>');
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
  if (!caps.contextInjection) {
    process.exit(0);
  }

  const proposalDir = findProposalDir(projectRoot, config, proposalId);
  if (!proposalDir) {
    process.exit(0);
  }

  const ctx = await adapter.getContext(proposalDir);
  if (ctx) {
    console.log(ctx);
  }
}

main().catch(() => process.exit(0));
