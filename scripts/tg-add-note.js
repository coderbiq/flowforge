#!/usr/bin/env node

const {
  appendImplementationNote,
  loadProposalContext,
} = require('./lib/tg-workflow');

function main() {
  const target = process.argv[2];
  const note = process.argv.slice(3).join(' ').trim();

  if (!target || !note) {
    console.error('Usage: scripts/tg-add-note.js <proposal-id|proposal-dir> <note text>');
    process.exit(2);
  }

  try {
    const context = loadProposalContext(target, process.cwd());
    const result = appendImplementationNote(context, note, process.cwd());

    console.log(JSON.stringify({
      id: context.meta.id,
      proposal_dir: context.proposalDir,
      notes_path: result.notes_path,
      timestamp: result.timestamp,
    }, null, 2));
  } catch (error) {
    console.error(`ERROR ${error.message}`);
    process.exit(1);
  }
}

main();
