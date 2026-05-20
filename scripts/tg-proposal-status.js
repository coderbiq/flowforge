#!/usr/bin/env node

const { beadTaskSummary, loadProposalContext, validateProposalContext } = require('./lib/tg-workflow');

function main() {
  const target = process.argv[2];
  if (!target) {
    console.error('Usage: scripts/tg-proposal-status.js <proposal-id|proposal-dir>');
    process.exit(2);
  }

  try {
    const context = loadProposalContext(target, process.cwd());
    const validation = validateProposalContext(context);
    const summary = {
      id: context.meta.id,
      title: context.meta.title,
      status: context.meta.status,
      proposal_dir: context.proposalDir,
      task_backend: context.meta.task_backend,
      task_count: context.taskMap.tasks.length,
      validation_errors: validation.errors.length,
      validation_warnings: validation.warnings.length,
      archive_targets: context.meta.archive_targets || [],
    };

    if (context.meta.task_backend === 'beads') {
      const beadSummary = beadTaskSummary(context.meta.id, process.cwd());
      summary.backend_available = beadSummary.available;
      summary.backend_error = beadSummary.error || null;
      summary.backend_open_tasks = beadSummary.openTasks ? beadSummary.openTasks.length : null;
    }

    console.log(JSON.stringify(summary, null, 2));
    if (validation.errors.length > 0) {
      process.exit(1);
    }
  } catch (error) {
    console.error(`ERROR ${error.message}`);
    process.exit(1);
  }
}

main();
