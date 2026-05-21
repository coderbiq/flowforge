#!/usr/bin/env node

const {
  archiveProposal,
  getArchiveReadiness,
  loadProposalContext,
} = require('./lib/flowforge');

function main() {
  const target = process.argv[2];
  if (!target) {
    console.error('Usage: scripts/flowforge-archive-proposal.js <proposal-id|proposal-dir>');
    process.exit(2);
  }

  try {
    const context = loadProposalContext(target, process.cwd());
    const readiness = getArchiveReadiness(context, process.cwd());

    for (const warning of readiness.warnings) {
      console.log(`WARN ${warning}`);
    }
    if (readiness.failures.length > 0) {
      for (const failure of readiness.failures) {
        console.error(`ERROR ${failure}`);
      }
      process.exit(1);
    }

    const result = archiveProposal(context, process.cwd());
    console.log(JSON.stringify(result, null, 2));
  } catch (error) {
    if (error.readiness?.failures) {
      for (const warning of error.readiness.warnings || []) {
        console.log(`WARN ${warning}`);
      }
      for (const failure of error.readiness.failures) {
        console.error(`ERROR ${failure}`);
      }
      process.exit(1);
    }

    console.error(`ERROR ${error.message}`);
    process.exit(1);
  }
}

main();
