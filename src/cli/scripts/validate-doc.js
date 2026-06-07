#!/usr/bin/env node
'use strict';

const fs = require('fs');
const path = require('path');

// 兼容 CLI 模式（argv[2]=projectRoot, argv[3]=docPath）和直接调用模式（argv[2]=docPath）
const docPath = process.argv[3] || process.argv[2];
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
  'intake', 'finding', 'decision', 'journal',
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

const domainOptionalTypes = ['proposal', 'task-map', 'notes', 'journal', 'intake'];
if (frontmatter.domain) {
  if (!frontmatter.domain.scope || !['system', 'module'].includes(frontmatter.domain.scope)) {
    errors.push(`domain.scope 无效: ${frontmatter.domain.scope || '缺失'}`);
  }
  if (frontmatter.domain.scope === 'module' && !frontmatter.domain.module) {
    errors.push('domain.scope=module 但缺少 domain.module');
  }
  if (!frontmatter.domain.type || !['design', 'decision', 'convention'].includes(frontmatter.domain.type)) {
    errors.push(`domain.type 无效: ${frontmatter.domain.type || '缺失'}`);
  }
} else if (!domainOptionalTypes.includes(frontmatter.doc_type)) {
  errors.push('缺少 domain 字段（可归档文档必须设置 domain）');
}

// L2: doc_type 专属字段校验
const TYPE_SPECIFIC = {
  finding: {
    required: ['source'],
    validStatus: ['active'],
    defaultImportance: 'info',
    defaultMaturity: 'seed'
  },
  convention: {
    required: ['enforcement'],
    validStatus: ['active', 'superseded', 'deprecated'],
    defaultImportance: 'should',
    defaultMaturity: 'growing'
  },
  architecture: {
    required: [],
    validStatus: ['draft', 'active', 'deprecated'],
    defaultImportance: 'should',
    defaultMaturity: 'growing'
  },
  decision: {
    required: ['decision_status'],
    validStatus: ['accepted', 'rejected', 'superseded'],
    defaultImportance: 'should',
    defaultMaturity: 'growing'
  },
  adr: {
    required: ['adr_id'],
    validStatus: ['proposed', 'accepted', 'rejected', 'superseded', 'deprecated'],
    defaultImportance: 'should',
    defaultMaturity: 'growing'
  },
  module: {
    required: [],
    validStatus: ['draft', 'active', 'deprecated'],
    defaultImportance: 'should',
    defaultMaturity: 'growing'
  }
};

const docRules = TYPE_SPECIFIC[frontmatter.doc_type];
if (docRules) {
  for (const field of docRules.required) {
    if (!frontmatter[field]) {
      errors.push(`缺少 ${frontmatter.doc_type} 专属必填字段: ${field}`);
    }
  }
  if (frontmatter.status && !docRules.validStatus.includes(frontmatter.status)) {
    errors.push(`${frontmatter.doc_type} 的 status 无效: ${frontmatter.status}（合法值: ${docRules.validStatus.join(', ')}）`);
  }
}

// L1+: importance/maturity 枚举校验
const validImportance = ['must', 'should', 'may', 'info'];
const validMaturity = ['seed', 'growing', 'stable', 'deprecated'];
if (frontmatter.domain?.importance && !validImportance.includes(frontmatter.domain.importance)) {
  errors.push(`domain.importance 无效: ${frontmatter.domain.importance}`);
}
if (frontmatter.domain?.maturity && !validMaturity.includes(frontmatter.domain.maturity)) {
  errors.push(`domain.maturity 无效: ${frontmatter.domain.maturity}`);
}

// 其他可选字段格式校验
if (frontmatter.review_interval && isNaN(Number(frontmatter.review_interval))) {
  errors.push(`review_interval 格式错误: ${frontmatter.review_interval}（期望数字）`);
}
if (frontmatter.last_reviewed && !isISODate(frontmatter.last_reviewed)) {
  errors.push(`last_reviewed 格式错误: ${frontmatter.last_reviewed}（期望 ISO-8601）`);
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

function isISODate(str) {
  return /^\d{4}-\d{2}-\d{2}(T\d{2}:\d{2}:\d{2}(\.\d+)?(Z|[+-]\d{2}:\d{2})?)?$/.test(str);
}
