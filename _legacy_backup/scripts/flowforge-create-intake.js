#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const {
  getIntakeRoot,
  nowIso,
  parseCliArgs,
  slugify,
} = require('./lib/flowforge');

function usage() {
  console.error('Usage: scripts/flowforge-create-intake.js --title <title> [--slug <slug>] [--workspace <workspace>] [--question <question>]');
}

function frontmatterBlock(title, slug, workspace, question) {
  const timestamp = nowIso();
  return `---\ndoc_type: note\ntitle: ${JSON.stringify(title)}\nstatus: draft\nworkspace: ${JSON.stringify(workspace)}\nmodule_scope: []\nsystem_scope: []\nconvention_scope: []\nownership: []\ninformation_class: note\ntopics: []\nrelated_docs: []\narchive_target: none\ncreated: ${JSON.stringify(timestamp)}\nupdated: ${JSON.stringify(timestamp)}\nintake_slug: ${JSON.stringify(slug)}\nquestion: ${JSON.stringify(question || '')}\n---\n`;
}

function main() {
  const args = parseCliArgs(process.argv.slice(2));
  const title = args.title;
  const slug = slugify(args.slug || title);
  const workspace = args.workspace || 'default';
  const question = args.question || '';

  if (!title) {
    usage();
    process.exit(2);
  }

  try {
    const intakeRoot = getIntakeRoot(process.cwd(), workspace);
    const intakeDir = path.join(intakeRoot, slug);
    if (fs.existsSync(intakeDir)) {
      throw new Error(`intake package already exists: ${intakeDir}`);
    }

    fs.mkdirSync(path.join(intakeDir, 'assets'), { recursive: true });
    fs.writeFileSync(path.join(intakeDir, 'index.md'), `${frontmatterBlock(title, slug, workspace, question)}\n# ${title}\n\n## Ownership summary\n\n- Primary module: none\n- System / architecture targets: none\n- Convention targets: none\n- Canonical reading path: this intake package\n\n## Problem statement\n\nDescribe the requested change in plain language.\n\n## Goals\n\n- Goal\n\n## Non-goals\n\n- Non-goal\n\n## Constraints\n\n- Constraint\n\n## References\n\n- Related issue, doc, screenshot, or link\n\n## Open questions\n\n- Question that must be answered during exploration\n`, 'utf8');
    fs.writeFileSync(path.join(intakeDir, 'references.md'), `---\ndoc_type: note\ntitle: Intake References\nstatus: draft\nworkspace: ${JSON.stringify(workspace)}\nmodule_scope: []\nsystem_scope: []\nconvention_scope: []\nownership: []\ninformation_class: note\ntopics: []\nrelated_docs: []\narchive_target: none\ncreated: ${JSON.stringify(nowIso())}\nupdated: ${JSON.stringify(nowIso())}\nintake_slug: ${JSON.stringify(slug)}\n---\n\n# References\n\n- Link or screenshot reference\n`, 'utf8');
    fs.writeFileSync(path.join(intakeDir, 'questions.md'), `---\ndoc_type: note\ntitle: Intake Questions\nstatus: draft\nworkspace: ${JSON.stringify(workspace)}\nmodule_scope: []\nsystem_scope: []\nconvention_scope: []\nownership: []\ninformation_class: note\ntopics: []\nrelated_docs: []\narchive_target: none\ncreated: ${JSON.stringify(nowIso())}\nupdated: ${JSON.stringify(nowIso())}\nintake_slug: ${JSON.stringify(slug)}\n---\n\n# Questions\n\n- Open question\n`, 'utf8');
    fs.writeFileSync(path.join(intakeDir, 'assets', 'README.md'), '# Assets\n\nPut screenshots, sketches, and other non-Markdown reference material here.\n', 'utf8');

    console.log(JSON.stringify({
      title,
      slug,
      workspace,
      intake_dir: path.relative(process.cwd(), intakeDir),
    }, null, 2));
  } catch (error) {
    console.error(`ERROR ${error.message}`);
    process.exit(1);
  }
}

main();
