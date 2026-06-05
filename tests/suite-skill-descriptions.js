'use strict';

const path = require('path');
const fs = require('fs');

const SKILLS_DIR = 'src/agents';

// Each SKILL must have unique trigger words not shared with others
const EXPECTED_TRIGGERS = {
  'flowforge-design': ['需求', '想法', '变更', '探索', '分析', '设计', '提案'],
  'flowforge-implement': ['实施', '执行', '任务', '推进'],
  'flowforge-feedback': ['测试', '失败', '问题', 'bug', 'block'],
  'flowforge-archive': ['归档', '沉淀', 'library', '知识'],
  'flowforge-docs': ['文档', 'frontmatter', '写作', 'doc_type'],
  'flowforge-progress': ['进度', 'latest_progress', 'INDEX'],
};

// Border triggers: each SKILL should define clear boundaries
const REQUIRED_BOUNDARIES = {
  'flowforge-design': ['flowforge-implement', 'flowforge-archive'],
  'flowforge-implement': ['flowforge-design', 'flowforge-feedback'],
  'flowforge-feedback': ['flowforge-implement', 'flowforge-design'],
  'flowforge-archive': [],
  'flowforge-docs': [],
  'flowforge-progress': [],
};

function run(root) {
  let passed = 0;
  let failed = 0;
  const errors = [];

  const agentsDir = path.join(root, SKILLS_DIR);
  const skillDirs = fs.readdirSync(agentsDir, { withFileTypes: true })
    .filter(d => d.isDirectory());

  const skills = {};
  const descriptions = {};

  for (const dir of skillDirs) {
    const skillPath = path.join(agentsDir, dir.name, 'SKILL.md');
    const content = fs.readFileSync(skillPath, 'utf8');
    const fm = extractFrontmatter(content);

    if (!fm || !fm.name || !fm.description) continue;

    skills[fm.name] = fm;
    descriptions[fm.name] = String(fm.description).trim();
  }

  // Check trigger word coverage
  for (const [name, triggers] of Object.entries(EXPECTED_TRIGGERS)) {
    const desc = descriptions[name];
    if (!desc) continue;

    const missing = [];
    for (const trigger of triggers) {
      if (!desc.includes(trigger)) {
        missing.push(trigger);
      }
    }
    if (missing.length > 0) {
      failed++;
      errors.push(`${name}: description missing trigger words: ${missing.join(', ')}`);
    } else {
      passed++;
    }
  }

  // Check mutual exclusivity: no two descriptions should share all trigger domains
  for (const nameA of Object.keys(descriptions)) {
    for (const nameB of Object.keys(descriptions)) {
      if (nameA >= nameB) continue;

      const aWords = extractKeywords(descriptions[nameA]);
      const bWords = extractKeywords(descriptions[nameB]);
      const overlap = aWords.filter(w => bWords.includes(w));

      if (overlap.length >= 3 && aWords.length > 0 && bWords.length > 0) {
        const overlapRatio = overlap.length / Math.min(aWords.length, bWords.length);
        if (overlapRatio > 0.4) {
          failed++;
          errors.push(`${nameA} / ${nameB}: high description overlap (${overlapRatio.toFixed(1)}) - shared keywords: ${overlap.join(', ')}`);
        } else {
          passed++;
        }
      }
    }
  }

  // Check boundary references
  for (const [name, boundaries] of Object.entries(REQUIRED_BOUNDARIES)) {
    const desc = descriptions[name];
    if (!desc) continue;

    for (const boundary of boundaries) {
      if (desc.includes(boundary)) {
        passed++;
      } else {
        failed++;
        errors.push(`${name}: description should reference boundary SKILL "${boundary}"`);
      }
    }
  }

  // Check general quality: each description must have multiple sentences
  for (const [name, desc] of Object.entries(descriptions)) {
    const sentences = desc.split(/[。！？.!?]\s*/).filter(s => s.trim().length > 5);
    if (sentences.length < 2) {
      failed++;
      errors.push(`${name}: description too brief (${sentences.length} substantive sentences, need >= 2)`);
    } else {
      passed++;
    }
  }

  return { passed, failed, errors };
}

function extractKeywords(text) {
  // Extract significant Chinese and English words
  const words = [];
  // Chinese 2-8 char sequences
  const cnMatches = text.match(/[\u4e00-\u9fff]{2,8}/g) || [];
  words.push(...cnMatches);
  // English words 3+ chars
  const enMatches = text.match(/[a-z]{3,}/gi) || [];
  words.push(...enMatches.map(w => w.toLowerCase()));
  return [...new Set(words)];
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
