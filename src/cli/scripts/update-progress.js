#!/usr/bin/env node
'use strict';

const fs = require('fs');
const path = require('path');
const { loadMainConfig, loadProjectConfig, loadMeta } = require('./lib/config');

// 兼容 CLI 模式（argv[2]=projectRoot, argv[3]=proposalPath, argv[4]=text）
// 和直接调用模式（argv[2]=proposalPath, argv[3]=text）
const proposalPath = process.argv[3] || process.argv[2];
const progressText = process.argv[4] || '';

if (!proposalPath) {
  console.error('用法: update-progress.js <proposal路径> "[进度总结]"');
  process.exit(0);
}

const metaPath = path.join(proposalPath, 'meta.yaml');
if (!fs.existsSync(metaPath)) {
  console.error(`ERROR: meta.yaml 不存在: ${metaPath}`);
  process.exit(0);
}

const projectRoot = findProjectRoot(proposalPath);
if (!projectRoot) {
  console.error('ERROR: 无法定位项目根路径（未找到 .flowforge/）');
  process.exit(0);
}

const meta = loadMeta(proposalPath);
const projectId = meta && meta.project;

let content = fs.readFileSync(metaPath, 'utf8');
const now = new Date().toISOString();

content = content.replace(/^(\s*updated_at\s*:\s*).*/m, `$1${now}`);

if (/^\s*latest_progress\s*:/m.test(content)) {
  content = content.replace(
    /^(\s*latest_progress\s*:\s*).*/m,
    `$1"${escapeYaml(progressText || '')}"`
  );
} else {
  content = content.replace(
    /^(\s*updated_at\s*:.*)$/m,
    `$1\nlatest_progress: "${escapeYaml(progressText || '')}"`
  );
}

fs.writeFileSync(metaPath, content, 'utf8');
console.log(`meta.yaml 已更新: latest_progress="${progressText}"`);

const refreshScript = path.join(projectRoot, '.flowforge', 'scripts', 'refresh-index.js');
if (fs.existsSync(refreshScript)) {
  try {
    const { execSync } = require('child_process');
    const args = projectId ? `"${projectRoot}" "${projectId}"` : `"${projectRoot}"`;
    const output = execSync(`node "${refreshScript}" ${args}`, { encoding: 'utf8' });
    console.log(output.trim());
  } catch (e) {
    console.error(`WARNING: refresh-index.js 执行失败: ${e.message}`);
  }
}

function findProjectRoot(startPath) {
  let dir = fs.lstatSync(startPath).isDirectory() ? startPath : path.dirname(startPath);
  while (true) {
    if (fs.existsSync(path.join(dir, '.flowforge', 'config.yaml'))) return dir;
    const parent = path.dirname(dir);
    if (parent === dir) return null;
    dir = parent;
  }
}

function escapeYaml(str) {
  return str.replace(/\\/g, '\\\\').replace(/"/g, '\\"');
}
