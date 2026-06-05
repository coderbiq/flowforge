'use strict';

const path = require('path');
const fs = require('fs');

const AGENTS_PATH = 'src/AGENTS.md';
const SKILLS_DIR = 'src/agents';

function run(root) {
  let passed = 0;
  let failed = 0;
  const errors = [];

  const agentsContent = fs.readFileSync(path.join(root, AGENTS_PATH), 'utf8');
  const agentsDir = path.join(root, SKILLS_DIR);

  const skillDirs = fs.readdirSync(agentsDir, { withFileTypes: true })
    .filter(d => d.isDirectory())
    .map(d => d.name);

  // Check 1: Every SKILL directory must be referenced in AGENTS.md
  for (const name of skillDirs) {
    if (agentsContent.includes(name)) {
      passed++;
    } else {
      failed++;
      errors.push(`AGENTS.md: SKILL "${name}" not referenced in AGENTS.md skill routing section`);
    }
  }

  // Check 2: AGENTS.md version must be consistent
  const versionMatch = agentsContent.match(/v:(\d+\.\d+\.\d+)/);
  if (!versionMatch) {
    failed++;
    errors.push('AGENTS.md: version comment (v:x.y.z) not found');
  } else {
    passed++;
    const agentsVersion = versionMatch[1];

    const metaPath = path.join(root, 'src', 'flowforge', 'meta.yaml');
    if (fs.existsSync(metaPath)) {
      const metaContent = fs.readFileSync(metaPath, 'utf8');
      const metaVersionMatch = metaContent.match(/version:\s*"?(\d+\.\d+\.\d+)"?/);
      if (metaVersionMatch) {
        if (metaVersionMatch[1] !== agentsVersion) {
          failed++;
          errors.push(`Version mismatch: AGENTS.md (${agentsVersion}) vs meta.yaml (${metaVersionMatch[1]})`);
        } else {
          passed++;
        }
      }
    }
  }

  // Check 3: AGENTS.md must have task operation rules section
  if (!agentsContent.includes('任务操作规则') && !agentsContent.includes('Task Operations')) {
    failed++;
    errors.push('AGENTS.md: missing task operation rules section');
  } else {
    passed++;
  }

  // Check 4: AGENTS.md must have CLI entry section
  if (!agentsContent.includes('CLI 入口') && !agentsContent.includes('CLI Entry')) {
    failed++;
    errors.push('AGENTS.md: missing CLI entry section');
  } else {
    passed++;
  }

  // Check 5: AGENTS.md must mention --proposal requirement
  if (!agentsContent.includes('--proposal')) {
    failed++;
    errors.push('AGENTS.md: missing --proposal flag usage');
  } else {
    passed++;
  }

  // Check 6: FlowForge comment wrapper integrity
  if (!agentsContent.includes('<!-- BEGIN FLOWFORGE') || !agentsContent.includes('<!-- END FLOWFORGE')) {
    failed++;
    errors.push('AGENTS.md: missing or corrupted FlowForge comment wrapper');
  } else {
    passed++;
  }

  // Check 7: Must not have stale task-map.yaml references
  if (agentsContent.includes('task-map.yaml')) {
    const lines = agentsContent.split('\n');
    for (const line of lines) {
      if (line.includes('task-map.yaml') && !line.includes('废弃') && !line.includes('deprecated') && !line.includes('禁止')) {
        failed++;
        errors.push('AGENTS.md: task-map.yaml referenced without deprecation warning');
        break;
      }
    }
    passed++;
  }

  return { passed, failed, errors };
}

module.exports = { run };
