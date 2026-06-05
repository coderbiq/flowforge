'use strict';

const path = require('path');
const fs = require('fs');

const SKILLS_DIR = 'src/agents';
const REQUIRED_FIELDS = ['name', 'description'];
const SKILL_NAME_PATTERN = /^flowforge-[a-z][a-z-]*[a-z]$/;

function run(root) {
  let passed = 0;
  let failed = 0;
  const errors = [];

  const agentsDir = path.join(root, SKILLS_DIR);
  if (!fs.existsSync(agentsDir)) {
    return { passed: 0, failed: 1, errors: [`SKILLS_DIR not found: ${agentsDir}`] };
  }

  const skillDirs = fs.readdirSync(agentsDir, { withFileTypes: true })
    .filter(d => d.isDirectory());

  if (skillDirs.length === 0) {
    return { passed: 0, failed: 1, errors: ['No SKILL directories found'] };
  }

  for (const dir of skillDirs) {
    const skillPath = path.join(agentsDir, dir.name, 'SKILL.md');
    if (!fs.existsSync(skillPath)) {
      failed++;
      errors.push(`${dir.name}/SKILL.md: file missing`);
      continue;
    }

    const content = fs.readFileSync(skillPath, 'utf8');
    const fm = extractFrontmatter(content);

    if (!fm) {
      failed++;
      errors.push(`${dir.name}/SKILL.md: no valid YAML frontmatter`);
      continue;
    }

    // Check required fields
    for (const field of REQUIRED_FIELDS) {
      if (!fm[field]) {
        failed++;
        errors.push(`${dir.name}/SKILL.md: missing required field "${field}"`);
      } else {
        passed++;
      }
    }

    // Validate name format
    if (fm.name) {
      if (!SKILL_NAME_PATTERN.test(fm.name)) {
        failed++;
        errors.push(`${dir.name}/SKILL.md: name "${fm.name}" must match flowforge-<verb> pattern`);
      } else {
        passed++;
      }

      // name must match directory name
      if (fm.name !== dir.name) {
        failed++;
        errors.push(`${dir.name}/SKILL.md: name "${fm.name}" must match directory "${dir.name}"`);
      } else {
        passed++;
      }
    }

    // Validate description structure
    if (fm.description) {
      const desc = String(fm.description).trim();

      if (desc.length < 50) {
        failed++;
        errors.push(`${dir.name}/SKILL.md: description too short (${desc.length} chars, min 50)`);
      } else {
        passed++;
      }

      if (desc.length > 2048) {
        failed++;
        errors.push(`${dir.name}/SKILL.md: description too long (${desc.length} chars, max 2048)`);
      } else {
        passed++;
      }

      // Must have activation section
      if (!desc.includes('激活') && !desc.includes('activate')) {
        failed++;
        errors.push(`${dir.name}/SKILL.md: description missing activation criteria`);
      } else {
        passed++;
      }

      // Must have anti-activation section
      if (!desc.includes('不要') && !desc.includes('不应') && !desc.includes('must not') && !desc.includes('should not')) {
        failed++;
        errors.push(`${dir.name}/SKILL.md: description missing deactivation criteria`);
      } else {
        passed++;
      }

      // Must not contain XML tags (common AI slop)
      if (/<\/?[a-z]+>/.test(desc)) {
        failed++;
        errors.push(`${dir.name}/SKILL.md: description contains XML tags`);
      } else {
        passed++;
      }
    }
  }

  return { passed, failed, errors };
}

function extractFrontmatter(text) {
  const m = text.match(/^---\n([\s\S]*?)\n---/);
  if (!m) return null;

  const result = {};
  const lines = m[1].split('\n');
  let currentKey = null;
  let currentValue = [];
  let keyIndent = 0;
  let isBlockScalar = false;

  for (const line of lines) {
    const indent = line.search(/\S/);
    if (indent === -1 && currentKey) {
      currentValue.push(line);
      continue;
    }

    const kv = line.match(/^\s*([a-zA-Z_][a-zA-Z0-9_]*)\s*:\s*(.*)/);
    if (kv && indent <= keyIndent && !isBlockScalar) {
      if (currentKey) {
        let val = currentValue.length > 1
          ? currentValue.join('\n').trim()
          : (currentValue[0] || '').trim();
        if (val.startsWith('|')) val = val.substring(1).trim();
        if (val.startsWith('>')) val = val.substring(1).trim();
        result[currentKey] = val;
      }
      currentKey = kv[1];
      currentValue = [kv[2].trim()];
      keyIndent = indent;
      isBlockScalar = kv[2].trim() === '|' || kv[2].trim() === '>';
    } else if (currentKey) {
      currentValue.push(line);
    }
  }

  if (currentKey) {
    let val = currentValue.length > 1
      ? currentValue.join('\n').trim()
      : (currentValue[0] || '').trim();
    if (val.startsWith('|')) val = val.substring(1).trim();
    if (val.startsWith('>')) val = val.substring(1).trim();
    result[currentKey] = val;
  }

  for (const key of Object.keys(result)) {
    const v = result[key];
    if (typeof v === 'string') result[key] = v.replace(/^["']|["']$/g, '');
  }

  return result;
}

module.exports = { run };
