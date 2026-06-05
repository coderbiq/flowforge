'use strict';

const path = require('path');
const fs = require('fs');
const { execSync } = require('child_process');

const JS_DIRS = [
  'src/cli',
  'src/cli/scripts',
  'src/cli/scripts/lib',
  'src/cli/scripts/lib/backends',
  'src/cli/scripts/vendor/js-yaml/lib',
  'tests',
];

function run(root) {
  let passed = 0;
  let failed = 0;
  const errors = [];

  for (const dir of JS_DIRS) {
    const fullDir = path.join(root, dir);
    if (!fs.existsSync(fullDir)) continue;

    const files = fs.readdirSync(fullDir, { recursive: true })
      .map(f => path.join(fullDir, f))
      .filter(f => f.endsWith('.js') && fs.statSync(f).isFile())
      .filter(f => !f.includes('node_modules') && !f.includes('.backup'));

    for (const file of files) {
      try {
        execSync(`node -c "${file}"`, { stdio: 'pipe', encoding: 'utf8', timeout: 5000 });
        passed++;
      } catch (e) {
        failed++;
        const relPath = path.relative(root, file);
        const errMsg = e.stderr ? e.stderr.toString().trim().split('\n')[0] : e.message;
        errors.push(`${relPath}: ${errMsg}`);
      }
    }
  }

  return { passed, failed, errors };
}

module.exports = { run };
