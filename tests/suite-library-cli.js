'use strict';

const path = require('path');
const fs = require('fs');
const { execSync } = require('child_process');

const SCRIPTS_DIR = 'src/cli/scripts';

function run(root) {
  let passed = 0;
  let failed = 0;
  const errors = [];

  const scriptsDir = path.join(root, SCRIPTS_DIR);

  // Check 1: library dispatcher exists
  const libraryJs = path.join(scriptsDir, 'library.js');
  if (fs.existsSync(libraryJs)) { passed++; }
  else { failed++; errors.push('library.js: missing'); }

  // Check 2: sub-scripts exist
  const subs = ['library-check.js', 'library-index.js', 'library-graph.js'];
  for (const sub of subs) {
    if (fs.existsSync(path.join(scriptsDir, sub))) { passed++; }
    else { failed++; errors.push(`${sub}: missing`); }
  }

  // Check 3: library dispatcher handles subcommands
  const libContent = fs.readFileSync(libraryJs, 'utf8');
  for (const sub of ['check', 'index', 'graph']) {
    if (libContent.includes(`case '${sub}':`)) { passed++; }
    else { failed++; errors.push(`library.js: missing case '${sub}'`); }
  }

  // Check 4: library dispatcher has help
  if (libContent.includes('--help') && libContent.includes('Subcommands')) { passed++; }
  else { failed++; errors.push('library.js: missing help output'); }

  // Check 5: CLI entry registers library command
  const flowforgeCli = path.join(root, 'src/cli/flowforge');
  const cliContent = fs.readFileSync(flowforgeCli, 'utf8');
  if (cliContent.includes("case 'library':")) { passed++; }
  else { failed++; errors.push('flowforge: missing library case'); }

  // Check 6: library-check.js handles flags
  const checkContent = fs.readFileSync(path.join(scriptsDir, 'library-check.js'), 'utf8');
  for (const flag of ['--staleness', '--broken-refs', '--duplicates', '--orphans']) {
    if (checkContent.includes(flag)) { passed++; }
    else { failed++; errors.push(`library-check.js: missing flag ${flag}`); }
  }
  if (checkContent.includes('--all')) { passed++; }
  else { failed++; errors.push('library-check.js: missing --all flag'); }

  // Check 7: library-check.js outputs JSON
  if (checkContent.includes('JSON.stringify')) { passed++; }
  else { failed++; errors.push('library-check.js: no JSON output'); }

  // Check 8: library-index.js handles --refresh
  const indexContent = fs.readFileSync(path.join(scriptsDir, 'library-index.js'), 'utf8');
  if (indexContent.includes('--refresh')) { passed++; }
  else { failed++; errors.push('library-index.js: missing --refresh'); }

  // Check 9: library-graph.js handles subcommands
  const graphContent = fs.readFileSync(path.join(scriptsDir, 'library-graph.js'), 'utf8');
  for (const cmd of ['backlinks', 'refs', 'orphans', 'hubs', 'blast-radius']) {
    if (graphContent.includes(`case '${cmd}':`)) { passed++; }
    else { failed++; errors.push(`library-graph.js: missing case '${cmd}'`); }
  }

  // Check 10: library-check.js handles empty library
  if (checkContent.includes("'Library not found'") || checkContent.includes('Library not found')) { passed++; }
  else { failed++; errors.push('library-check.js: no empty library handling'); }

  // Check 11: library.js argv routing (Bug fix: projectRoot from argv[3]?argv[2]:cwd)
  if (libContent.includes('process.argv[3] ? process.argv[2] : process.cwd()')) { passed++; }
  else { failed++; errors.push('library.js: missing projectRoot argv[3]?argv[2]:cwd() pattern'); }

  // Check 12: library.js subcommand from argv[3]||argv[2]
  if (libContent.includes('process.argv[3] || process.argv[2]')) { passed++; }
  else { failed++; errors.push('library.js: missing subcommand argv[3]||argv[2] pattern'); }

  // Check 13-15: findLibraryRoot in all 3 scripts
  for (const script of ['library-check.js', 'library-index.js', 'library-graph.js']) {
    const content = fs.readFileSync(path.join(scriptsDir, script), 'utf8');
    if (content.includes('function findLibraryRoot')) { passed++; }
    else { failed++; errors.push(`${script}: missing findLibraryRoot function`); }
  }

  // Check 16: findLibraryRoot scans project configs
  if (checkContent.includes('wikiRoot') && checkContent.includes('.flowforge')) { passed++; }
  else { failed++; errors.push('library-check.js: findLibraryRoot missing project config scan'); }

  // Check 17: findLibraryRoot prefers .md content over empty dirs
  if (checkContent.includes('.md') && checkContent.includes('hasContent')) { passed++; }
  else { failed++; errors.push('library-check.js: findLibraryRoot missing content-based selection'); }

  return { passed, failed, errors };
}

module.exports = { run };
