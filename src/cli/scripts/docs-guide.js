#!/usr/bin/env node
'use strict';

const fs = require('fs');
const path = require('path');
const { findProjectRoot } = require('./lib/config');

let projectRoot = process.argv[2] || process.cwd();
let docType = process.argv[3];

// 兼容单参数调用：docs-guide.js <doc_type>
if (!docType && projectRoot !== process.cwd() && !projectRoot.startsWith('/')) {
  docType = projectRoot;
  projectRoot = process.cwd();
}

const resolved = findProjectRoot(projectRoot);
if (resolved) projectRoot = resolved;

const guidesDir = path.join(projectRoot, '.flowforge', 'guides');

if (!docType) {
  if (!fs.existsSync(guidesDir)) {
    console.log('暂无写作指南。');
    process.exit(0);
  }
  const files = fs.readdirSync(guidesDir).filter(f => f.endsWith('.md'));
  console.log('# Registered Document Types\n');
  console.log('| doc_type | 默认位置 |');
  console.log('|----------|---------|');
  for (const f of files) {
    const guideContent = fs.readFileSync(path.join(guidesDir, f), 'utf8');
    const location = extractLocation(guideContent);
    const name = f.replace('.md', '');
    console.log(`| \`${name}\` | ${location || '—'} |`);
  }
  process.exit(0);
}

const guidePath = path.join(guidesDir, `${docType}.md`);
if (!fs.existsSync(guidePath)) {
  console.log(`## ${docType}\n`);
  console.log('该类型暂无写作指南。');
  process.exit(0);
}

console.log(fs.readFileSync(guidePath, 'utf8'));

function extractLocation(content) {
  const m = content.match(/## 位置\n\n`([^`]+)`/);
  return m ? m[1] : null;
}
