#!/usr/bin/env node
'use strict';

const fs = require('fs');
const path = require('path');

const proposalDir = process.argv[2];
if (!proposalDir) {
  console.error('用法: validate-proposal.js <proposal目录>');
  process.exit(0);
}

if (!fs.existsSync(proposalDir)) {
  console.error(`ERROR: 目录不存在: ${proposalDir}`);
  process.exit(0);
}

const errors = [];

const metaPath = path.join(proposalDir, 'meta.yaml');
if (!fs.existsSync(metaPath)) {
  errors.push('缺少 meta.yaml');
} else {
  const meta = parseYaml(fs.readFileSync(metaPath, 'utf8'));
  const requiredFields = ['id', 'title', 'status', 'created_at', 'updated_at'];
  for (const field of requiredFields) {
    if (!meta[field]) errors.push(`meta.yaml 缺少必填字段: ${field}`);
  }
  if (meta.id && !/^[A-Z]*\d{6}\d{2}$/.test(meta.id)) {
    errors.push(`meta.yaml id 格式疑似错误: ${meta.id}（期望 前缀+YYMMDDNN）`);
  }
  if (meta.status && !['draft', 'active', 'implemented', 'archived', 'rejected'].includes(meta.status)) {
    errors.push(`meta.yaml status 无效: ${meta.status}`);
  }
}

if (!fs.existsSync(path.join(proposalDir, 'proposal.md'))) {
  errors.push('缺少 proposal.md');
}

const designDir = path.join(proposalDir, 'design');
const designFile = path.join(proposalDir, 'design.md');
if (!fs.existsSync(designDir) && !fs.existsSync(designFile)) {
  errors.push('缺少 design/ 目录或 design.md');
}

if (!fs.existsSync(path.join(proposalDir, 'task-map.yaml'))) {
  errors.push('缺少 task-map.yaml');
}

if (errors.length === 0) {
  console.log(`PASS: ${path.basename(proposalDir)}`);
} else {
  console.log(`FAIL: ${path.basename(proposalDir)}`);
  for (const e of errors) console.log(`  - ${e}`);
}

function parseYaml(content) {
  const result = {};
  for (const line of content.split('\n')) {
    const m = line.match(/^\s*([a-zA-Z_]+)\s*:\s*(.*)/);
    if (m) result[m[1]] = m[2].trim().replace(/^["']|["']$/g, '');
  }
  return result;
}
