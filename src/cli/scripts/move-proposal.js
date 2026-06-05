#!/usr/bin/env node
'use strict';

const fs = require('fs');
const path = require('path');
const { loadMainConfig, loadProjectConfig, loadMeta } = require('./lib/config');

const projectRoot = process.argv[2] || process.cwd();
const proposalId = process.argv[3];

if (!proposalId) {
  console.error('用法: move-proposal.js <projectRoot> <proposalId>');
  process.exit(0);
}

const config = loadMainConfig(projectRoot);
if (!config) {
  console.error('ERROR: .flowforge/config.yaml 不存在或格式错误');
  process.exit(0);
}

const projectRefs = config.projects || [];
if (projectRefs.length === 0) {
  console.error('ERROR: config.yaml 中未定义 projects');
  process.exit(1);
}

const allProjects = [];
for (const ref of projectRefs) {
  const p = loadProjectConfig(projectRoot, ref);
  if (p) allProjects.push(p);
}

let proposalLocation = null;
let activeProject = null;

for (const p of allProjects) {
  const ws = path.join(projectRoot, p.wikiRoot, 'workspace');
  for (const sub of ['active', 'completed']) {
    const subDir = path.join(ws, 'proposals', sub);
    if (!fs.existsSync(subDir)) continue;
    const dirs = fs.readdirSync(subDir, { withFileTypes: true }).filter(d => d.isDirectory());
    for (const d of dirs) {
      if (d.name === proposalId || d.name.startsWith(proposalId + '-')) {
        proposalLocation = { proposalDir: path.join(subDir, d.name), projectId: p.id, wikiRoot: p.wikiRoot, currentSub: sub };
        activeProject = p;
        break;
      }
    }
    if (proposalLocation) break;
  }
  if (proposalLocation) break;
}

if (!proposalLocation) {
  console.error(`ERROR: 未找到 proposal: ${proposalId}`);
  process.exit(1);
}

const meta = loadMeta(proposalLocation.proposalDir);
if (!meta) {
  console.error('ERROR: meta.yaml 不存在或格式错误');
  process.exit(1);
}

const result = { proposalId, steps: [] };

const metaPath = path.join(proposalLocation.proposalDir, 'meta.yaml');
let metaContent = fs.readFileSync(metaPath, 'utf8');
const now = new Date().toISOString();

metaContent = metaContent.replace(/^(\s*status\s*:\s*).*/m, '$1"archived"');
metaContent = metaContent.replace(/^(\s*updated_at\s*:\s*).*/m, `$1${now}`);

fs.writeFileSync(metaPath, metaContent, 'utf8');
result.steps.push({ step: 'update_meta', status: 'done', detail: 'status → archived, updated_at 已刷新' });

if (proposalLocation.currentSub === 'active') {
  const wikiWs = path.join(projectRoot, proposalLocation.wikiRoot, 'workspace');
  const activeDir = path.join(wikiWs, 'proposals', 'active', proposalId);
  const completedDir = path.join(wikiWs, 'proposals', 'completed');
  const targetDir = path.join(completedDir, proposalId);

  if (!fs.existsSync(completedDir)) {
    fs.mkdirSync(completedDir, { recursive: true });
    result.steps.push({ step: 'create_completed_dir', status: 'done', detail: completedDir });
  }

  if (fs.existsSync(targetDir)) {
    result.steps.push({ step: 'move_directory', status: 'skipped', detail: `目标已存在: ${targetDir}` });
  } else {
    fs.renameSync(activeDir, targetDir);
    result.steps.push({ step: 'move_directory', status: 'done', detail: `active → completed: ${path.relative(projectRoot, targetDir)}` });
  }
} else {
  result.steps.push({ step: 'move_directory', status: 'skipped', detail: `已在 ${proposalLocation.currentSub}/ 中，无需移动` });
}

if (activeProject && activeProject.rules && activeProject.rules.library && activeProject.rules.library.autoUpdateHistory) {
  const archiveTargets = meta.archive_targets || [];
  for (const target of archiveTargets) {
    const key = typeof target === 'string' ? target : target.key;
    if (!key) continue;

    const moduleName = extractModuleName(key);
    if (moduleName) {
      const historyPath = path.join(projectRoot, proposalLocation.wikiRoot, 'library', 'modules', moduleName, 'HISTORY.md');
      const entry = `| ${now} | ${meta.id || proposalId} | ${meta.title || ''} | archived |\n`;
      appendHistory(historyPath, entry);
      result.steps.push({ step: 'update_history', module: moduleName, status: 'done', detail: `已追加归档记录到 ${path.relative(projectRoot, historyPath)}` });
    }
  }
}

console.log(JSON.stringify(result, null, 2));

function extractModuleName(key) {
  const patterns = [
    /^module:(.+)$/,
    /^key:\s*(.+)-module$/,
    /^(.+)-module$/
  ];
  for (const p of patterns) {
    const m = key.match(p);
    if (m) return m[1];
  }
  return null;
}

function appendHistory(historyPath, entry) {
  const dir = path.dirname(historyPath);
  if (!fs.existsSync(dir)) {
    fs.mkdirSync(dir, { recursive: true });
  }

  if (!fs.existsSync(historyPath)) {
    fs.writeFileSync(historyPath, `# Module History\n\n| Date | Proposal | Title | Action |\n|------|----------|-------|--------|\n${entry}`, 'utf8');
    return;
  }

  let content = fs.readFileSync(historyPath, 'utf8');
  if (content.includes(entry.trim())) return;

  const tableEnd = content.lastIndexOf('\n');
  if (tableEnd > 0) {
    content = content.substring(0, tableEnd + 1) + entry;
    fs.writeFileSync(historyPath, content, 'utf8');
  }
}
