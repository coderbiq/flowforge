'use strict';

const path = require('path');
const fs = require('fs');

const META_PATH = 'src/flowforge/meta.yaml';
const AGENTS_PATH = 'src/AGENTS.md';
const CLI_PATH = 'src/cli/flowforge';

function run(root) {
  let passed = 0;
  let failed = 0;
  const errors = [];

  let metaVersion = null;
  let agentsVersion = null;
  let cliVersion = null;

  if (fs.existsSync(path.join(root, META_PATH))) {
    const meta = fs.readFileSync(path.join(root, META_PATH), 'utf8');
    const m = meta.match(/version:\s*"?(\d+\.\d+\.\d+)"?/);
    if (m) {
      metaVersion = m[1];
      passed++;
    } else {
      failed++;
      errors.push('meta.yaml: version field not found');
    }
  }

  if (fs.existsSync(path.join(root, AGENTS_PATH))) {
    const agents = fs.readFileSync(path.join(root, AGENTS_PATH), 'utf8');
    const m = agents.match(/v:(\d+\.\d+\.\d+)/);
    if (m) {
      agentsVersion = m[1];
      passed++;
    } else {
      failed++;
      errors.push('AGENTS.md: version comment (v:x.y.z) not found');
    }
  }

  if (fs.existsSync(path.join(root, CLI_PATH))) {
    const cli = fs.readFileSync(path.join(root, CLI_PATH), 'utf8');
    const m = cli.match(/FlowForge v(\d+\.\d+(?:\.\d+)?)/);
    if (m) {
      cliVersion = m[1];
      passed++;
    } else {
      failed++;
      errors.push('CLI: version string not found');
    }
  }

  if (metaVersion && agentsVersion && metaVersion !== agentsVersion) {
    failed++;
    errors.push('meta.yaml (' + metaVersion + ') != AGENTS.md (' + agentsVersion + ')');
  } else if (metaVersion && agentsVersion) {
    passed++;
  }

  if (metaVersion && cliVersion) {
    const cliMajorMinor = cliVersion.split('.').slice(0, 2).join('.');
    const metaMajorMinor = metaVersion.split('.').slice(0, 2).join('.');
    if (cliMajorMinor !== metaMajorMinor) {
      failed++;
      errors.push('meta.yaml (' + metaVersion + ') != CLI (' + cliVersion + ')');
    } else {
      passed++;
    }
  }

  return { passed, failed, errors };
}

module.exports = { run };
