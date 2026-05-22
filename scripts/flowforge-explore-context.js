#!/usr/bin/env node

const { loadIntakeContext, loadProjectRuleBundle, parseCliArgs } = require('./lib/flowforge');

function renderMarkdown(rules, intake) {
  const lines = [];
  lines.push('# FlowForge Exploration Seed Context');
  lines.push('');
  lines.push('## Rule Bundle');
  lines.push('');
  lines.push(`- workspace: ${rules.workspace.name}`);
  lines.push(`- rules_root: ${rules.rulesRoot}`);
  lines.push(`- available: ${rules.available ? 'yes' : 'no'}`);
  lines.push(`- missing_files: ${rules.missing_files.length > 0 ? rules.missing_files.join(', ') : 'none'}`);
  lines.push('');

  for (const file of rules.files) {
    lines.push(`### ${file.file_name}`);
    lines.push('');
    if (!file.exists) {
      lines.push('_missing_');
      lines.push('');
      continue;
    }
    lines.push(file.content.trimEnd());
    lines.push('');
  }

  lines.push('## Intake Package');
  lines.push('');
  if (intake) {
    lines.push(`- intake_dir: ${intake.intakeDir}`);
    lines.push(`- assets: ${intake.assets.length > 0 ? intake.assets.join(', ') : 'none'}`);
    lines.push('');

    for (const file of intake.files) {
      lines.push(`### ${file.file_name}`);
      lines.push('');
      lines.push(file.content.trimEnd());
      lines.push('');
    }
  } else {
    lines.push('- intake_dir: none');
    lines.push('- assets: none');
    lines.push('');
    lines.push('_No intake package was provided. Continue with the rule bundle and any in-chat request material, and create or attach an intake package before expanding the exploration if the topic needs durable evidence._');
    lines.push('');
  }

  lines.push('## Exploration Instructions');
  lines.push('');
  lines.push('- Read the rule bundle before drafting the exploration skeleton.');
  lines.push('- Read the intake package as evidence, not as a fixed outline.');
  lines.push('- Generate a new exploration structure that cites the intake package sources used.');
  lines.push('- If intake material is missing or stale, note the gap and continue with explicit assumptions.');
  lines.push('');

  return lines.join('\n').trimEnd() + '\n';
}

function main() {
  const args = parseCliArgs(process.argv.slice(2));
  const intakeTarget = args['intake-package'] || args._[0];
  const workspace = args.workspace || null;
  const jsonMode = Boolean(args.json);

  try {
    const rules = loadProjectRuleBundle(process.cwd(), workspace);
    const intake = intakeTarget ? loadIntakeContext(intakeTarget, process.cwd()) : null;
    if (jsonMode) {
      console.log(JSON.stringify({ rules, intake }, null, 2));
      return;
    }

    process.stdout.write(renderMarkdown(rules, intake));
  } catch (error) {
    console.error(`ERROR ${error.message}`);
    process.exit(1);
  }
}

main();
