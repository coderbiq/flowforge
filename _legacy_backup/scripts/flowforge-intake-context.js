#!/usr/bin/env node

const { loadIntakeContext, parseCliArgs } = require('./lib/flowforge');

function renderMarkdown(context) {
  const lines = [];
  lines.push('# FlowForge Intake Context');
  lines.push('');
  lines.push(`- intake_dir: ${context.intakeDir}`);
  lines.push(`- assets: ${context.assets.length > 0 ? context.assets.join(', ') : 'none'}`);
  lines.push('');

  for (const file of context.files) {
    lines.push(`## ${file.file_name}`);
    lines.push('');
    lines.push(file.content.trimEnd());
    lines.push('');
  }

  return lines.join('\n').trimEnd() + '\n';
}

function main() {
  const args = parseCliArgs(process.argv.slice(2));
  const target = args._[0];
  const jsonMode = Boolean(args.json);

  if (!target) {
    console.error('Usage: scripts/flowforge-intake-context.js <intake-slug|intake-dir> [--json]');
    process.exit(2);
  }

  try {
    const context = loadIntakeContext(target, process.cwd());
    if (jsonMode) {
      console.log(JSON.stringify(context, null, 2));
      return;
    }

    process.stdout.write(renderMarkdown(context));
  } catch (error) {
    console.error(`ERROR ${error.message}`);
    process.exit(1);
  }
}

main();
