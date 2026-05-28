#!/usr/bin/env node

const { loadProjectRuleBundle, parseCliArgs } = require('./lib/flowforge');

function renderMarkdown(bundle) {
  const lines = [];
  lines.push('# FlowForge Project Rules Context');
  lines.push('');
  lines.push(`- workspace: ${bundle.workspace.name}`);
  lines.push(`- rules_root: ${bundle.rulesRoot}`);
  lines.push(`- available: ${bundle.available ? 'yes' : 'no'}`);
  lines.push(`- missing_files: ${bundle.missing_files.length > 0 ? bundle.missing_files.join(', ') : 'none'}`);
  lines.push('');

  for (const file of bundle.files) {
    lines.push(`## ${file.file_name}`);
    lines.push('');
    if (!file.exists) {
      lines.push('_missing_');
      lines.push('');
      continue;
    }
    lines.push(file.content.trimEnd());
    lines.push('');
  }

  return lines.join('\n').trimEnd() + '\n';
}

function main() {
  const args = parseCliArgs(process.argv.slice(2));
  const workspace = args.workspace || null;
  const jsonMode = Boolean(args.json);

  try {
    const bundle = loadProjectRuleBundle(process.cwd(), workspace);
    if (jsonMode) {
      console.log(JSON.stringify(bundle, null, 2));
      return;
    }

    process.stdout.write(renderMarkdown(bundle));
  } catch (error) {
    console.error(`ERROR ${error.message}`);
    process.exit(1);
  }
}

main();
