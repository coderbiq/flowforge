#!/usr/bin/env node
'use strict';

const fs = require('fs');
const path = require('path');
const { loadMainConfig, loadProjectConfig, loadMeta } = require('./lib/config');

const projectRoot = process.argv[2] || process.cwd();
const projectId = process.argv[3] || null;

const config = loadMainConfig(projectRoot);
if (!config) {
  console.error('ERROR: .flowforge/config.yaml 不存在或格式错误');
  process.exit(0);
}

const projectRefs = config.projects || [];
const targets = projectId
  ? projectRefs.filter(p => p.id === projectId)
  : projectRefs;

if (targets.length === 0) {
  console.error(`ERROR: project '${projectId}' 未在 config.yaml 中定义`);
  process.exit(0);
}

for (const ref of targets) {
  const p = loadProjectConfig(projectRoot, ref);
  if (!p) continue;
  generateIndex(projectRoot, p);
}

function generateIndex(projectRoot, project) {
  const proposalsDir = path.join(projectRoot, project.wikiRoot, 'workspace', 'proposals');
  if (!fs.existsSync(proposalsDir)) {
    console.log(`[${project.id}] proposals 目录不存在，跳过`);
    return;
  }

  const activeStatuses = ['draft', 'active', 'implemented'];
  const completedStatuses = ['archived', 'rejected'];
  const activeProposals = [];
  const completedProposals = [];

  for (const subdir of ['active', 'completed']) {
    const dir = path.join(proposalsDir, subdir);
    if (!fs.existsSync(dir)) continue;
    const entries = fs.readdirSync(dir, { withFileTypes: true }).filter(d => d.isDirectory());
    for (const entry of entries) {
      const pd = path.join(dir, entry.name);
      const meta = loadMeta(pd);
      if (!meta) continue;
      const item = { ...meta, _relPath: `${subdir}/${entry.name}` };
      if (activeStatuses.includes(meta.status)) {
        activeProposals.push(item);
      } else if (completedStatuses.includes(meta.status)) {
        completedProposals.push(item);
      }
    }
  }

  const sortByUpdated = (a, b) => (b.updated_at || '').localeCompare(a.updated_at || '');
  activeProposals.sort(sortByUpdated);
  completedProposals.sort(sortByUpdated);

  const total = activeProposals.length + completedProposals.length;
  const now = new Date();
  const ts = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')} ${String(now.getHours()).padStart(2, '0')}:${String(now.getMinutes()).padStart(2, '0')}`;

  let md = '';
  md += `# 📋 Proposals Index\n\n`;
  if (project.name) md += `**Project**: ${project.name} (${project.id})\n\n`;
  md += `> ${ts} · ${total} 提案（🟢 ${activeProposals.length} / 📦 ${completedProposals.length}）\n\n`;

  md += '## 🔄 进行中\n\n';
  if (activeProposals.length === 0) {
    md += '（无）\n\n';
  } else {
    md += '| ID | 标题 | 规模 | 状态 | 模块 | 更新 | 进度 |\n';
    md += '|----|------|------|------|------|------|------|\n';
    for (const p of activeProposals) {
      md += row(p);
    }
    md += '\n';
  }

  md += '## 📦 已完成\n\n';
  if (completedProposals.length === 0) {
    md += '（无）\n\n';
  } else {
    md += '| ID | 标题 | 规模 | 状态 | 模块 | 更新 | 进度 |\n';
    md += '|----|------|------|------|------|------|------|\n';
    for (const p of completedProposals) {
      md += row(p);
    }
    md += '\n';
  }

  const indexPath = path.join(proposalsDir, 'INDEX.md');
  fs.writeFileSync(indexPath, md, 'utf8');
  console.log(`[${project.id}] INDEX.md → ${path.relative(projectRoot, indexPath)} (${activeProposals.length} active / ${completedProposals.length} completed)`);
}

function row(p) {
  const idLink = `[${p.id}](./${p._relPath}/)`;
  const size = p.size_class || '—';
  const status = `\`${p.status}\``;
  const modules = formatModules(p.modules);
  const updated = formatDate(p.updated_at);
  const progress = p.latest_progress || '—';
  return `| ${idLink} | ${p.title || '—'} | ${size} | ${status} | ${modules} | ${updated} | ${progress} |\n`;
}

function formatModules(modules) {
  if (!modules || !Array.isArray(modules) || modules.length === 0) return '—';
  return modules.map(m => `\`${m}\``).join(' ');
}

function formatDate(isoStr) {
  if (!isoStr) return '—';
  const m = String(isoStr).match(/(\d{4})-(\d{2})-(\d{2})/);
  return m ? `${m[2]}-${m[3]}` : String(isoStr).substring(0, 10);
}
