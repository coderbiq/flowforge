#!/usr/bin/env node
'use strict';

const fs = require('fs');
const path = require('path');

const proposalPath = process.argv[2];
const progressText = process.argv[3];

if (!proposalPath) {
  console.error('用法: update-progress.js <proposal完整路径> "[进度总结]"');
  process.exit(0);
}

if (!fs.existsSync(proposalPath)) {
  console.error(`ERROR: 目录不存在: ${proposalPath}`);
  process.exit(0);
}

const metaPath = path.join(proposalPath, 'meta.yaml');
if (!fs.existsSync(metaPath)) {
  console.error(`ERROR: meta.yaml 不存在: ${metaPath}`);
  process.exit(0);
}

// ── 更新 meta.yaml ──
let content = fs.readFileSync(metaPath, 'utf8');
const now = new Date().toISOString();

// 更新 updated_at
content = content.replace(
  /^(\s*updated_at\s*:\s*).*/m,
  `$1${now}`
);

// 更新或新增 latest_progress
if (/^\s*latest_progress\s*:/m.test(content)) {
  content = content.replace(
    /^(\s*latest_progress\s*:\s*).*/m,
    `$1"${escapeYaml(progressText || '')}"`
  );
} else {
  // 在 updated_at 行后插入 latest_progress
  content = content.replace(
    /^(\s*updated_at\s*:.*)$/m,
    `$1\nlatest_progress: "${escapeYaml(progressText || '')}"`
  );
}

fs.writeFileSync(metaPath, content, 'utf8');
console.log(`meta.yaml 已更新: latest_progress="${progressText}"`);

// ── 确定项目根路径并刷新 INDEX.md ──
// 从 proposal 路径向上推导项目根路径
// proposalPath 形如: .../ff-wiki/workspace/proposals/active/CRxxxxx
const parts = proposalPath.split(path.sep);
const proposalsIdx = parts.lastIndexOf('proposals');
if (proposalsIdx >= 3) {
  // proposals 前面是 workspace, ff-wiki, 项目根
  const projectRoot = parts.slice(0, proposalsIdx - 2).join(path.sep);
  const refreshScript = path.join(projectRoot, '.flowforge', 'scripts', 'refresh-index.js');
  if (fs.existsSync(refreshScript)) {
    try {
      const { execSync } = require('child_process');
      const output = execSync(`node "${refreshScript}" "${projectRoot}"`, { encoding: 'utf8' });
      console.log(output.trim());
    } catch (e) {
      console.error(`WARNING: refresh-index.js 执行失败: ${e.message}`);
      console.error('请手动运行: node .flowforge/scripts/refresh-index.js');
    }
  } else {
    console.log('refresh-index.js 不存在，跳过 INDEX.md 重建');
  }
} else {
  console.log('无法推导项目根路径，跳过 INDEX.md 重建');
}

// ── Helpers ──

function escapeYaml(str) {
  return str.replace(/\\/g, '\\\\').replace(/"/g, '\\"');
}
