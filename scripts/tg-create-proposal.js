#!/usr/bin/env node

const path = require('path');
const {
  createProposalSkeleton,
  parseCliArgs,
  validateProposalContext,
  loadProposalContext,
} = require('./lib/tg-workflow');

function toArray(value) {
  if (value === undefined) return [];
  return Array.isArray(value) ? value : [value];
}

function parseArchiveTarget(raw, index) {
  const parts = String(raw).split(':');
  if (parts.length < 2 || parts.length > 3) {
    throw new Error(`invalid --archive-target value at item ${index + 1}: ${raw}`);
  }

  const [type, targetPath, role = index === 0 ? 'primary' : 'secondary'] = parts;
  return { type, path: targetPath, role };
}

function usage() {
  console.error('Usage: scripts/tg-create-proposal.js --title <title> --source-exploration <path> --archive-target <type:path[:role]> [--archive-target ...] [--owner <owner>] [--task-backend beads|github|linear|none] [--tag <tag>] [--slug <slug>] [--status draft|proposed|approved]');
}

function main() {
  const args = parseCliArgs(process.argv.slice(2));
  const title = args.title;
  const sourceExploration = args['source-exploration'];
  const archiveTargetValues = toArray(args['archive-target']);

  if (!title || !sourceExploration || archiveTargetValues.length === 0) {
    usage();
    process.exit(2);
  }

  try {
    const archiveTargets = archiveTargetValues.map(parseArchiveTarget);
    const created = createProposalSkeleton({
      title,
      slug: args.slug,
      owner: args.owner,
      sourceExploration,
      taskBackend: args['task-backend'],
      status: args.status,
      tags: toArray(args.tag),
      archiveTargets,
    }, process.cwd());

    const context = loadProposalContext(created.proposalDir, process.cwd());
    const validation = validateProposalContext(context, process.cwd());

    for (const warning of validation.warnings) {
      console.log(`WARN ${warning}`);
    }
    for (const error of validation.errors) {
      console.error(`ERROR ${error}`);
    }

    if (validation.errors.length > 0) {
      process.exit(1);
    }

    console.log(JSON.stringify({
      id: created.id,
      slug: created.slug,
      proposal_dir: path.relative(process.cwd(), created.proposalDir),
      status: created.meta.status,
      task_backend: created.meta.task_backend,
      archive_targets: created.meta.archive_targets,
    }, null, 2));
  } catch (error) {
    console.error(`ERROR ${error.message}`);
    process.exit(1);
  }
}

main();
