#!/usr/bin/env node

const {
  loadProposalContext,
  transitionProposalStatus,
  validateProposalContext,
} = require('./lib/flowforge');

function main() {
  const target = process.argv[2];
  if (!target) {
    console.error('Usage: scripts/flowforge-approve-proposal.js <proposal-id|proposal-dir>');
    process.exit(2);
  }

  try {
    const context = loadProposalContext(target, process.cwd());
    const validation = validateProposalContext(context, process.cwd());

    for (const warning of validation.warnings) {
      console.log(`WARN ${warning}`);
    }
    if (validation.errors.length > 0) {
      for (const error of validation.errors) {
        console.error(`ERROR ${error}`);
      }
      process.exit(1);
    }

    if (!['draft', 'proposed', 'approved'].includes(context.meta.status)) {
      console.error(`ERROR proposal status must be draft or proposed before approve, got ${context.meta.status}`);
      process.exit(1);
    }

    if (context.meta.status !== 'approved') {
      transitionProposalStatus(context, 'approved');
    }

    console.log(JSON.stringify({
      id: context.meta.id,
      status: context.meta.status,
      proposal_dir: context.proposalDir,
    }, null, 2));
  } catch (error) {
    console.error(`ERROR ${error.message}`);
    process.exit(1);
  }
}

main();
