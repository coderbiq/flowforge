#!/usr/bin/env node

const { getArchiveReadiness, loadProposalContext } = require('./lib/tg-workflow');

function main() {
  const target = process.argv[2];
  if (!target) {
    console.error('Usage: scripts/tg-check-archive.js <proposal-id|proposal-dir>');
    process.exit(2);
  }

  try {
    const context = loadProposalContext(target, process.cwd());
    const readiness = getArchiveReadiness(context, process.cwd());
    const failures = [...readiness.failures];
    const warnings = [...readiness.warnings];

    for (const warning of warnings) {
      console.log(`WARN ${warning}`);
    }
    for (const failure of failures) {
      console.error(`ERROR ${failure}`);
    }

    if (failures.length > 0) {
      process.exit(1);
    }

    console.log(`OK archive-check ${context.meta.id}`);
  } catch (error) {
    console.error(`ERROR ${error.message}`);
    process.exit(1);
  }
}

main();
