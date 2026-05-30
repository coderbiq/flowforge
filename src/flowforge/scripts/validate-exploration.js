#!/usr/bin/env node
'use strict';

const fs = require('fs');
const path = require('path');

const explorationDir = process.argv[2];
if (!explorationDir) {
  console.error('用法: validate-exploration.js <exploration目录>');
  process.exit(0);
}

if (!fs.existsSync(explorationDir)) {
  console.error(`ERROR: 目录不存在: ${explorationDir}`);
  process.exit(0);
}

const errors = [];

const indexPath = path.join(explorationDir, 'index.md');
if (!fs.existsSync(indexPath)) {
  errors.push('缺少 index.md');
} else {
  const content = fs.readFileSync(indexPath, 'utf8');
  const fm = extractFrontmatter(content);
  if (!fm) {
    errors.push('index.md 缺少 frontmatter');
  } else {
    const requiredFields = ['title', 'status', 'question', 'created', 'updated'];
    for (const field of requiredFields) {
      if (!fm[field]) errors.push(`index.md frontmatter 缺少: ${field}`);
    }
    if (fm.status && !['active', 'archived', 'rejected'].includes(fm.status)) {
      errors.push(`index.md status 无效: ${fm.status}`);
    }
    if (fm.confidence && !['high', 'medium', 'low'].includes(fm.confidence)) {
      errors.push(`index.md confidence 无效: ${fm.confidence}`);
    }
    if (!fm.domain) {
      errors.push('index.md 缺少 domain 字段（需设置 scope、type）');
    } else {
      const dm = parseDomain(fm.domain);
      if (!dm || !dm.scope || !dm.type) {
        errors.push('index.md domain 字段不完整（需要 scope 和 type）');
      } else {
        if (!['system', 'module'].includes(dm.scope)) {
          errors.push(`index.md domain.scope 无效: ${dm.scope}`);
        }
        if (dm.scope === 'module' && !dm.module) {
          errors.push('index.md domain.scope=module 但缺少 module 字段');
        }
        if (!['design', 'decision', 'convention'].includes(dm.type)) {
          errors.push(`index.md domain.type 无效: ${dm.type}`);
        }
      }
    }
  }
}

if (!fs.existsSync(path.join(explorationDir, 'findings'))) {
  errors.push('缺少 findings/ 目录');
}
if (!fs.existsSync(path.join(explorationDir, 'journal'))) {
  errors.push('缺少 journal/ 目录');
}

if (errors.length === 0) {
  console.log(`PASS: ${path.basename(explorationDir)}`);
} else {
  console.log(`FAIL: ${path.basename(explorationDir)}`);
  for (const e of errors) console.log(`  - ${e}`);
}

function extractFrontmatter(text) {
  const m = text.match(/^---\n([\s\S]*?)\n---/);
  if (!m) return null;
  const result = {};
  let currentKey = null;

  for (const line of m[1].split('\n')) {
    const nested = line.match(/^\s{2}(\w+)\s*:\s*(.*)/);
    if (nested && currentKey === 'domain') {
      if (!result.domain) result.domain = {};
      result.domain[nested[1]] = nested[2].trim().replace(/^["']|["']$/g, '');
      continue;
    }
    const kv = line.match(/^\s*([a-zA-Z_]+)\s*:\s*(.*)/);
    if (kv) {
      currentKey = kv[1];
      result[kv[1]] = kv[2].trim().replace(/^["']|["']$/g, '');
    } else {
      currentKey = null;
    }
  }
  return result;
}

function parseDomain(dm) {
  if (typeof dm === 'object') return dm;
  return null;
}
