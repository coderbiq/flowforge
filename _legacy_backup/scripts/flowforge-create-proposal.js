#!/usr/bin/env node

const path = require('path');
const {
  createProposalSkeleton,
  parseCliArgs,
  parseOwnershipEntry,
  loadExplorationContext,
  validateProposalContext,
  loadProposalContext,
} = require('./lib/flowforge');

function toArray(value) {
  if (value === undefined) return [];
  return Array.isArray(value) ? value : [value];
}

function parseArchiveTarget(raw, index) {
  const parts = String(raw).split(':');
  if (parts.length < 2) {
    throw new Error(`invalid --archive-target value at item ${index + 1}: ${raw}`);
  }

  const [typeSpec, ref, maybeRole, maybeKey] = parts;
  const [type, workspace] = typeSpec.includes('@') ? typeSpec.split('@', 2) : [typeSpec, undefined];
  if (!type || !ref) {
    throw new Error(`invalid --archive-target value at item ${index + 1}: ${raw}`);
  }

  const role = maybeRole === 'primary' || maybeRole === 'secondary'
    ? maybeRole
    : index === 0
      ? 'primary'
      : 'secondary';

  const key = maybeKey && maybeKey !== role ? maybeKey : undefined;
  return {
    type,
    workspace,
    ref,
    role,
    key,
  };
}

function parseCanonicalCorpus(raw, index) {
  const parts = String(raw).split(':');
  if (parts.length < 2) {
    throw new Error(`invalid --canonical-corpus value at item ${index + 1}: ${raw}`);
  }

  const [typeSpec, ref, maybeRole] = parts;
  const [type, workspace] = typeSpec.includes('@') ? typeSpec.split('@', 2) : [typeSpec, undefined];
  if (!type || !ref) {
    throw new Error(`invalid --canonical-corpus value at item ${index + 1}: ${raw}`);
  }

  const role = maybeRole === 'primary' ? 'primary' : 'secondary';
  return {
    type,
    workspace,
    ref,
    role,
  };
}

function usage() {
  console.error('Usage: scripts/flowforge-create-proposal.js --title <title> --source-exploration <ref> --archive-target <type[@workspace]:ref[:role[:key]]> [--archive-target ...] [--canonical-corpus <type[@workspace]:ref[:role]>] [--canonical-corpus ...] [--workspace <workspace>] [--scope workspace|cross-workspace|monorepo] [--size-class small|medium|large] [--design-layout single|split] [--ownership <type:target[:role]> ...] [--reusable-rule <title[:summary]> ...] [--owner <owner>] [--task-backend beads|github|linear|none] [--tag <tag>] [--slug <slug>] [--status draft|proposed|approved]');
}

function main() {
  const args = parseCliArgs(process.argv.slice(2));
  const title = args.title;
  const sourceExploration = args['source-exploration'];
  const archiveTargetValues = toArray(args['archive-target']);
  const canonicalCorpusValues = toArray(args['canonical-corpus']);
  const ownershipValues = toArray(args.ownership);
  const reusableRuleValues = toArray(args['reusable-rule']);

  if (!title || !sourceExploration || archiveTargetValues.length === 0) {
    usage();
    process.exit(2);
  }

  try {
    const archiveTargets = archiveTargetValues.map(parseArchiveTarget);
    const canonicalCorpus = canonicalCorpusValues.map(parseCanonicalCorpus);
    const ownership = ownershipValues.map(parseOwnershipEntry);
    const reusableRules = reusableRuleValues.map((raw) => {
      const [titleText, ...rest] = String(raw).split(':');
      const summary = rest.join(':').trim();
      return {
        title: titleText.trim(),
        summary: summary || undefined,
      };
    }).filter((rule) => rule.title);
    const explorationContext = loadExplorationContext(sourceExploration, process.cwd());
    const inheritedReusableRules = reusableRules.length > 0
      ? reusableRules
      : (explorationContext.parsed?.reusable_rules || []);
    const created = createProposalSkeleton({
      title,
      slug: args.slug,
      owner: args.owner,
      workspace: args.workspace,
      scope: args.scope,
      sizeClass: args['size-class'],
      designLayout: args['design-layout'],
      ownership: ownership.length > 0 ? ownership : undefined,
      reusableRules: inheritedReusableRules.length > 0 ? inheritedReusableRules : undefined,
      sourceExploration,
      taskBackend: args['task-backend'],
      status: args.status,
      tags: toArray(args.tag),
      archiveTargets,
      canonicalCorpus: canonicalCorpus.length > 0 ? canonicalCorpus : undefined,
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
      size_class: created.meta.size_class,
      ownership: created.meta.ownership,
      design_layout: created.meta.links.design === 'design/README.md' ? 'split' : 'single',
      canonical_corpus: created.meta.canonical_corpus,
      archive_targets: created.meta.archive_targets,
    }, null, 2));
  } catch (error) {
    console.error(`ERROR ${error.message}`);
    process.exit(1);
  }
}

main();
