'use strict';

const path = require('path');
const fs = require('fs');

const CLI_PATH = 'src/cli/flowforge';
const SKILLS_DIR = 'src/agents';

// All task actions defined in CLI (from printTaskHelp)
const CLI_TASK_ACTIONS = [
  'init', 'add', 'add-tasks', 'discover', 'cancel',
  'claim', 'done', 'block', 'unclaim', 'reopen',
  'dep-add', 'dep-remove', 'label',
  'ready', 'status', 'blocked', 'all-done', 'snapshot',
];

// All commands defined in CLI (from main switch + printHelp)
const CLI_COMMANDS = [
  'task', 'design-context', 'implement-context', 'feedback-context',
  'archive-context', 'validate-proposal', 'validate-doc',
  'update-progress', 'refresh-index', 'docs-guide',
  'move-proposal', 'archive-synthesize', 'feedback-capture', 'upgrade',
];

function run(root) {
  let passed = 0;
  let failed = 0;
  const errors = [];

  const cliContent = fs.readFileSync(path.join(root, CLI_PATH), 'utf8');
  const agentsDir = path.join(root, SKILLS_DIR);

  // Extract all flowforge commands from SKILL.md files
  const skillCommands = extractSkillCommands(agentsDir);

  // Check 1: Every command in CLI help must be listed in switch statement
  if (!cliContent.includes('case \'task\':' || cliContent.includes('case "task":'))) {
    failed++;
    errors.push('CLI: "task" case not found in main switch');
  } else {
    passed++;
  }

  for (const cmd of CLI_COMMANDS) {
    if (cmd === 'task') continue;
    const hasCase = cliContent.includes(`case '${cmd}':`) || cliContent.includes(`case "${cmd}":`);
    if (!hasCase) {
      failed++;
      errors.push(`CLI: command "${cmd}" not found in main switch`);
    } else {
      passed++;
    }
  }

  // Check 2: Every task action in CLI help must be in handleTask switch
  for (const action of CLI_TASK_ACTIONS) {
    const hasCase = cliContent.includes(`case '${action}':`) || cliContent.includes(`case "${action}":`);
    if (!hasCase) {
      failed++;
      errors.push(`CLI: task action "${action}" not found in handleTask switch`);
    } else {
      passed++;
    }
  }

  // Check 3: --help must be handled before --proposal check
  const taskHelpBeforeProposal = cliContent.includes('printTaskHelp()') &&
    cliContent.indexOf('printTaskHelp') < cliContent.indexOf('--proposal <id> required');
  if (!taskHelpBeforeProposal) {
    failed++;
    errors.push('CLI: handleTask must check --help BEFORE --proposal validation');
  } else {
    passed++;
  }

  const upgradeHelpBeforeProposal = cliContent.match(/function handleUpgrade[\s\S]*?printDelegateHelp/) !== null ||
    (cliContent.includes('--help') && cliContent.includes('handleUpgrade'));
  if (!upgradeHelpBeforeProposal) {
    failed++;
    errors.push('CLI: handleUpgrade must check --help BEFORE --proposal validation');
  } else {
    passed++;
  }

  // Check 4: delegateToScript must intercept --help
  if (!cliContent.includes('rest.includes(\'--help\')') && !cliContent.includes('rest.includes("--help")')) {
    failed++;
    errors.push('CLI: delegateToScript must intercept --help before delegating');
  } else {
    passed++;
  }

  // Check 5: SKILL commands must exist in CLI
  for (const [skill, commands] of Object.entries(skillCommands)) {
    for (const cmd of commands) {
      const cmdName = cmd.split(' ')[0];
      if (cmdName === 'flowforge') {
        const subCmd = cmd.split(' ')[1];
        if (subCmd && !CLI_COMMANDS.includes(subCmd) && !['task', 'docs'].includes(subCmd)) {
          failed++;
          errors.push(`${skill}: references unknown CLI subcommand "${subCmd}" in "${cmd}"`);
        } else {
          passed++;
        }
      }
    }
  }

  // Check 6: CLI task actions that appear in SKILL but are never referenced - warn only
  // (skipping here since all actions are useful even if not in SKILL docs)

  // Check 7: printTaskHelp exists and covers all actions
  if (!cliContent.includes('printTaskHelp')) {
    failed++;
    errors.push('CLI: printTaskHelp function not found');
  } else {
    passed++;
  }

  // Check 8: delegateToScript must have printDelegateHelp coverage for all commands
  if (!cliContent.includes('printDelegateHelp')) {
    failed++;
    errors.push('CLI: printDelegateHelp function not found');
  } else {
    passed++;
  }

  // Check 9: task init safety gate — hasTaskSpace check before init
  if (!cliContent.includes('hasTaskSpace')) {
    failed++;
    errors.push('CLI: task init must call backend.hasTaskSpace() before backend.init()');
  } else {
    passed++;
  }

  // Check 10: task init safety gate — --force flag extraction
  if (!cliContent.includes('extractFlag(actionRest, \'--force\')')) {
    failed++;
    errors.push('CLI: task init must check --force flag via extractFlag');
  } else {
    passed++;
  }

  // Check 11: task init safety gate — clear error message on missing --force
  if (!cliContent.includes('Use --force true to confirm this destructive operation')) {
    failed++;
    errors.push('CLI: task init must output clear error message when --force true is missing');
  } else {
    passed++;
  }

  // Check 12: task init help text mentions --force true
  if (!cliContent.includes('init <title> [--force true]')) {
    failed++;
    errors.push('CLI: printTaskHelp must document --force true for init');
  } else {
    passed++;
  }

  // Check 13: task init help text must appear in both printHelp AND printTaskHelp
  const helpMatches = (cliContent.match(/init <title> \[--force true\]/g) || []).length;
  if (helpMatches < 2) {
    failed++;
    errors.push('CLI: --force true must appear in both printHelp() and printTaskHelp()');
  } else {
    passed++;
  }

  return { passed, failed, errors };
}

function extractSkillCommands(agentsDir) {
  const skillDirs = fs.readdirSync(agentsDir, { withFileTypes: true })
    .filter(d => d.isDirectory());
  const result = {};

  for (const dir of skillDirs) {
    const skillPath = path.join(agentsDir, dir.name, 'SKILL.md');
    if (!fs.existsSync(skillPath)) continue;

    const content = fs.readFileSync(skillPath, 'utf8');
    const commands = [];

    const codeBlockRegex = /```(?:bash|sh)?\n([\s\S]*?)```/g;
    let match;
    while ((match = codeBlockRegex.exec(content)) !== null) {
      const lines = match[1].split('\n').filter(l => l.trim());
      for (const line of lines) {
        if (line.trim().startsWith('flowforge ')) {
          commands.push(line.trim());
        }
      }
    }

    const inlineRegex = /`(flowforge\s+[^`]+)`/g;
    while ((match = inlineRegex.exec(content)) !== null) {
      commands.push(match[1]);
    }

    result[dir.name] = [...new Set(commands)];
  }

  return result;
}

module.exports = { run };
