'use strict';

const path = require('path');
const fs = require('fs');

const SCRIPTS_DIR = 'src/cli/scripts';

const CONTEXT_SCRIPTS = {
  'design-context.js': ['# Design Context', '## Projects', '## Domain'],
  'implement-context.js': ['# Implement Context', '## Current Proposal', '## Task Status', '## Implement Rules'],
  'feedback-context.js': ['# Feedback Context', '## Current Proposal', '## Blocked Tasks'],
  'archive-context.js': ['# Archive Context', '## Current Proposal', '## 归档目标'],
};

const REQUIRED_SCRIPTS = [
  'design-context.js', 'implement-context.js', 'feedback-context.js',
  'archive-context.js', 'validate-proposal.js', 'validate-doc.js',
  'update-progress.js', 'refresh-index.js', 'docs-guide.js',
  'move-proposal.js', 'archive-synthesize.js', 'feedback-capture.js',
];

function run(root) {
  let passed = 0;
  let failed = 0;
  const errors = [];

  const scriptsDir = path.join(root, SCRIPTS_DIR);

  // Check 1: All required scripts exist
  for (const script of REQUIRED_SCRIPTS) {
    if (fs.existsSync(path.join(scriptsDir, script))) {
      passed++;
    } else {
      failed++;
      errors.push(`Missing script: ${script}`);
    }
  }

  // Check 2: Context scripts output expected sections
  for (const [script, requiredSections] of Object.entries(CONTEXT_SCRIPTS)) {
    const content = fs.readFileSync(path.join(scriptsDir, script), 'utf8');

    for (const section of requiredSections) {
      if (content.includes(section)) {
        passed++;
      } else {
        failed++;
        errors.push(`${script}: missing expected output section "${section}"`);
      }
    }

    // Must have findActiveProposal fallback
    if (content.includes('findActiveProposal') || content.includes('findProposal')) {
      passed++;
    } else {
      failed++;
      errors.push(`${script}: missing proposal auto-detection (findActiveProposal/findProposal)`);
    }
  }

  // Check 3: Scripts must handle argv correctly for CLI mode
  for (const script of REQUIRED_SCRIPTS) {
    const content = fs.readFileSync(path.join(scriptsDir, script), 'utf8');

    const usesArgv2Directly = content.match(/process\.argv\[2\]\s*(?![|=]\s*process)/);
    const usesArgv3Fallback = content.includes('process.argv[3] || process.argv[2]');

    if (script.startsWith('validate-') || script === 'update-progress.js') {
      if (usesArgv3Fallback) {
        passed++;
      } else if (usesArgv2Directly) {
        failed++;
        errors.push(`${script}: uses process.argv[2] directly without process.argv[3] fallback (may break in CLI mode)`);
      }
    } else if (script.match(/^(design|implement|feedback|archive)-context\.js$/)) {
      if (content.includes('process.argv[2] || process.cwd()')) {
        passed++;
      } else {
        failed++;
        errors.push(`${script}: missing projectRoot=process.argv[2]||process.cwd() pattern`);
      }
    }
  }

  // Check 4: design-context.js must handle --check-id and --suggest-id
  {
    const content = fs.readFileSync(path.join(scriptsDir, 'design-context.js'), 'utf8');
    if (content.includes('--check-id') && content.includes('checkProposalId')) {
      passed++;
    } else {
      failed++;
      errors.push('design-context.js: missing --check-id flag handling');
    }
    if (content.includes('--suggest-id') && content.includes('suggestProposalId')) {
      passed++;
    } else {
      failed++;
      errors.push('design-context.js: missing --suggest-id flag handling');
    }
  }

  // Check 5: validate-proposal.js must include ID uniqueness check
  {
    const content = fs.readFileSync(path.join(scriptsDir, 'validate-proposal.js'), 'utf8');
    if (content.includes('checkProposalId')) {
      passed++;
    } else {
      failed++;
      errors.push('validate-proposal.js: missing ID uniqueness check (checkProposalId)');
    }
    if (content.includes('otherConflicts') || content.includes('conflicts.filter')) {
      passed++;
    } else {
      failed++;
      errors.push('validate-proposal.js: missing self-exclusion logic for uniqueness check');
    }
  }

  return { passed, failed, errors };
}

module.exports = { run };
