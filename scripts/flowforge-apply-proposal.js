#!/usr/bin/env node

const {
  ensureBeadsTasks,
  ensureNotesFile,
  loadProposalContext,
  transitionProposalStatus,
  validateProposalContext,
} = require('./lib/flowforge');

function main() {
  const target = process.argv[2];
  if (!target) {
    console.error('Usage: scripts/flowforge-apply-proposal.js <proposal-id|proposal-dir>');
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

    if (!['approved', 'active'].includes(context.meta.status)) {
      console.error(`ERROR proposal status must be approved before apply, got ${context.meta.status}`);
      process.exit(1);
    }

    let backendSummary = null;
    if (context.meta.task_backend === 'beads') {
      backendSummary = ensureBeadsTasks(context, process.cwd());
    }

    const notesCreated = ensureNotesFile(context, process.cwd());

    if (context.meta.status === 'approved') {
      transitionProposalStatus(context, 'active');
    }

    console.log(JSON.stringify({
      id: context.meta.id,
      status: context.meta.status,
      proposal_dir: context.proposalDir,
      notes_created: notesCreated,
      task_backend: context.meta.task_backend,
      task_epic_id: context.meta.task_epic_id,
      backend: backendSummary,
    }, null, 2));
  } catch (error) {
    console.error(`ERROR ${error.message}`);
    process.exit(1);
  }
}

main();
