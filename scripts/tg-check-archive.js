#!/usr/bin/env node

const { beadTaskSummary, loadProposalContext, validateProposalContext } = require('./lib/tg-workflow');

function main() {
  const target = process.argv[2];
  if (!target) {
    console.error('Usage: scripts/tg-check-archive.js <proposal-id|proposal-dir>');
    process.exit(2);
  }

  try {
    const context = loadProposalContext(target, process.cwd());
    const validation = validateProposalContext(context);
    const failures = [...validation.errors];
    const warnings = [...validation.warnings];

    if (context.meta.status !== 'implemented') {
      failures.push(`proposal status must be implemented before archive, got ${context.meta.status}`);
    }

    if (context.meta.task_backend === 'beads') {
      const beadSummary = beadTaskSummary(context.meta.id, process.cwd());
      if (!beadSummary.available) {
        failures.push(`cannot verify Beads tasks: ${beadSummary.error}`);
      } else if (beadSummary.openTasks.length > 0) {
        failures.push(`proposal still has ${beadSummary.openTasks.length} open Beads tasks`);
      }
    }

    const primaryTarget = (context.meta.archive_targets || []).find((target) => target.role === 'primary');
    if (!primaryTarget) {
      failures.push('proposal must define a primary archive target');
    }

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
