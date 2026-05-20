#!/usr/bin/env node

const { loadProposalContext, validateProposalContext } = require('./lib/tg-workflow');

function main() {
  const target = process.argv[2];
  if (!target) {
    console.error('Usage: scripts/tg-validate-proposal.js <proposal-id|proposal-dir>');
    process.exit(2);
  }

  try {
    const context = loadProposalContext(target, process.cwd());
    const result = validateProposalContext(context);

    if (result.errors.length === 0 && result.warnings.length === 0) {
      console.log(`OK ${context.meta.id} ${context.proposalDir}`);
      return;
    }

    for (const warning of result.warnings) {
      console.log(`WARN ${warning}`);
    }
    for (const error of result.errors) {
      console.error(`ERROR ${error}`);
    }

    if (result.errors.length > 0) {
      process.exit(1);
    }
  } catch (error) {
    console.error(`ERROR ${error.message}`);
    process.exit(1);
  }
}

main();
