#!/usr/bin/env node
'use strict';

const fs = require('fs');
const path = require('path');

const docPath = process.argv[2];
if (!docPath) {
  console.error('用法: validate-doc.js <文档路径>');
  process.exit(0);
}

if (!fs.existsSync(docPath)) {
  console.error(`ERROR: 文件不存在: ${docPath}`);
  process.exit(0);
}

const content = fs.readFileSync(docPath, 'utf8');

const frontmatter = extractFrontmatter(content);
if (!frontmatter) {
  console.log('FAIL: 未找到 frontmatter（缺少 --- 包裹的 YAML 头部）');
  process.exit(0);
}

const errors = [];

const requiredFields = ['doc_type', 'title', 'status', 'created', 'updated'];
for (const field of requiredFields) {
  if (!frontmatter[field]) {
    errors.push(`缺少必填字段: ${field}`);
  }
}

const validDocTypes = [
  'intake', 'exploration', 'finding', 'decision', 'journal',
  'proposal', 'design', 'model', 'task-map', 'notes',
  'module', 'architecture', 'convention', 'adr'
];
if (frontmatter.doc_type && !validDocTypes.includes(frontmatter.doc_type)) {
  errors.push(`无效的 doc_type: ${frontmatter.doc_type}`);
}

if (frontmatter.created && !isISODate(frontmatter.created)) {
  errors.push(`created 格式错误: ${frontmatter.created}（期望 ISO-8601）`);
}
if (frontmatter.updated && !isISODate(frontmatter.updated)) {
  errors.push(`updated 格式错误: ${frontmatter.updated}（期望 ISO-8601）`);
}

if (errors.length === 0) {
  console.log(`PASS: ${path.basename(docPath)} (doc_type: ${frontmatter.doc_type})`);
} else {
  console.log(`FAIL: ${path.basename(docPath)}`);
  for (const e of errors) console.log(`  - ${e}`);
}

function extractFrontmatter(text) {
  const m = text.match(/^---\n([\s\S]*?)\n---/);
  if (!m) return null;
  const result = {};
  for (const line of m[1].split('\n')) {
    const kv = line.match(/^\s*([a-zA-Z_]+)\s*:\s*(.*)/);
    if (kv) result[kv[1]] = kv[2].trim().replace(/^["']|["']$/g, '');
  }
  return result;
}

function isISODate(str) {
  return /^\d{4}-\d{2}-\d{2}(T\d{2}:\d{2}:\d{2}(\.\d+)?(Z|[+-]\d{2}:\d{2})?)?$/.test(str);
}
