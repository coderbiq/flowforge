#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const cp = require('child_process');

const ROOT = process.cwd();
const GUIDES_DIR = path.join(ROOT, 'workflow', 'guides');
const BANNED_HEADING_PATTERNS = [
  /(^|\s)(why|background|history|tool positioning|design principles|design philosophy|explanation|overview|summary|recommend(ed)? structure)\b/i,
  /(为什么|背景|历史|说明|工具定位|设计原则|设计哲学|解释|概述|总结|摘要|推荐结构|优先级|归属边界|人类可读)/,
];

const DIRECTIVE_PATTERNS = [
  /(必须|不能|不要|应该|需要|先|再|如果|当|确认|检查|创建|更新|加载|路由|切换|保留|遵守|禁止|验证)/,
  /(输入|输出|输入项|输出项|场景|动作|约束|规则|退出条件|适用场景|必做|Steps|Actions|Inputs|Outputs|Exit|Contract|Guidance)/,
];

function sh(command) {
  return cp.execSync(command, { cwd: ROOT, encoding: 'utf8', stdio: ['ignore', 'pipe', 'pipe'] });
}

function collectGuideFiles() {
  const changed = new Set();
  try {
    for (const line of sh('git diff --name-only -- workflow/guides').split('\n')) {
      if (line.trim()) changed.add(line.trim());
    }
  } catch {}

  try {
    for (const line of sh('git ls-files --others --exclude-standard -- workflow/guides').split('\n')) {
      if (line.trim()) changed.add(line.trim());
    }
  } catch {}

  return [...changed]
    .filter((file) => file.endsWith('.md'))
    .map((file) => path.resolve(ROOT, file));
}

function headingText(line) {
  return line.replace(/^#{1,6}\s*/, '').trim();
}

function extractHeadings(content) {
  return content
    .split('\n')
    .filter((line) => /^#{2,6}\s+\S/.test(line))
    .map(headingText);
}

function hasAllowedHeading(headings) {
  return headings.some((heading) =>
    DIRECTIVE_PATTERNS.some((pattern) => pattern.test(heading))
  );
}

function findBannedHeadings(headings) {
  const hits = [];
  for (const heading of headings) {
    if (BANNED_HEADING_PATTERNS.some((pattern) => pattern.test(heading))) {
      hits.push(heading);
    }
  }
  return hits;
}

function validateFile(filePath) {
  const rel = path.relative(ROOT, filePath);
  const content = fs.readFileSync(filePath, 'utf8');
  const headings = extractHeadings(content);
  const errors = [];

  if (!/^#\s+\S/.test(content)) {
    errors.push('missing top-level title heading');
  }

  if (!hasAllowedHeading(headings) && !DIRECTIVE_PATTERNS.some((pattern) => pattern.test(content))) {
    errors.push('missing directive language; guide text must be action-oriented, not explanatory');
  }

  const banned = findBannedHeadings(headings);
  if (banned.length > 0) {
    errors.push(`contains explanation-heavy headings: ${banned.join(', ')}`);
  }

  return { rel, errors };
}

function main() {
  const targets = collectGuideFiles();

  if (targets.length === 0) {
    console.log('OK no modified workflow/guides files');
    return;
  }

  const findings = targets.map(validateFile);
  let failed = false;

  for (const finding of findings) {
    if (finding.errors.length === 0) continue;
    failed = true;
    for (const error of finding.errors) {
      console.error(`ERROR ${finding.rel}: ${error}`);
    }
  }

  if (!failed) {
    console.log(`OK ${targets.length} workflow/guides file(s)`);
    return;
  }

  process.exit(1);
}

main();
