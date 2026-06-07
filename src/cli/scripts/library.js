#!/usr/bin/env node
'use strict';

const path = require('path');

const projectRoot = process.argv[3] ? process.argv[2] : process.cwd();
const subcommand = process.argv[3] || process.argv[2];
const rest = process.argv[3] ? process.argv.slice(4) : process.argv.slice(3);

function delegateToScript(scriptName, scriptArgs) {
  require('child_process').execSync(
    `node "${path.join(__dirname, scriptName)}" "${projectRoot}" ${scriptArgs.map(a => `"${a}"`).join(' ')}`,
    { stdio: 'inherit', cwd: projectRoot }
  );
}

if (!subcommand || subcommand === '--help') {
  console.log('flowforge library <subcommand> [args]');
  console.log('');
  console.log('Subcommands:');
  console.log('  check     Audit & review (--all, --staleness, --quality, --review-list, --review <path>)');
  console.log('  list      List documents (--type, --scope, --importance, --maturity, --module)');
  console.log('  graph     Document relationship graph (backlinks, refs, orphans, blast-radius, hubs)');
  console.log('  index     Refresh INDEX.md (--refresh)');
  console.log('  surgeon   Content maintenance (merge, deprecate, upgrade)');
  console.log('  init      Initialize library with seed templates (--template minimal|full)');
  process.exit(0);
}

switch (subcommand) {
  case 'check':
    delegateToScript('library-check.js', rest);
    break;
  case 'list':
    delegateToScript('library-list.js', rest);
    break;
  case 'index':
    delegateToScript('library-index.js', rest);
    break;
  case 'graph':
    delegateToScript('library-graph.js', rest);
    break;
  case 'surgeon':
    delegateToScript('library-surgeon.js', rest);
    break;
  default:
    console.error(`Unknown library subcommand: ${subcommand}`);
    console.error('Run flowforge library --help for usage.');
    process.exit(1);
}
