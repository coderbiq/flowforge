#!/usr/bin/env node
'use strict';

const fs = require('fs');
const path = require('path');
const { loadMainConfig, loadProjectConfig, loadMeta } = require('./lib/config');

const projectRoot = require('./lib/config').findProjectRoot(process.argv[2] || process.cwd());
const proposalId = process.argv[3] || null;

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
if (proposalId) {
  proposalLocation = findProposalById(projectRoot, allProjects, proposalId);
} else {
  proposalLocation = findActiveProposal(projectRoot, allProjects);
}

if (!proposalLocation) {
  console.error('ERROR: 未找到活跃状态的 proposal。请确认 active/ 目录下存在 proposal。');
  process.exit(1);
}

const activeProject = allProjects.find(p => p.id === proposalLocation.projectId);
const r = activeProject && activeProject.rules ? activeProject.rules : null;

console.log('# Implement Context\n');

if (r && r.implement) {
  console.log('## Implement Rules\n');
  if (r.implement.task_states) {
    console.log(`task_states: ${r.implement.task_states.join(', ')}`);
  }
  if (r.implement.notes && r.implement.notes.fields) {
    console.log(`notes.fields: ${r.implement.notes.fields.join(', ')}`);
  }
  console.log('');
  if (r.implement.strategy) {
    console.log('## Implement Strategy\n');
    console.log(r.implement.strategy.trim());
    console.log('');
  }
}

if (r && r.design && r.design.task_rules) {
  console.log('## Task Rules\n');
  const tr = r.design.task_rules;
  if (tr.fields) {
    console.log(`fields: ${tr.fields.join(', ')}`);
  }
  if (tr.time_estimate) {
    console.log(`time_estimate: ${tr.time_estimate}`);
  }
  console.log('');
}

console.log('## Current Proposal\n');
console.log(`路径: ${path.relative(projectRoot, proposalLocation.proposalDir)}`);
console.log(`project: ${proposalLocation.projectId}`);
console.log(`wikiRoot: ${proposalLocation.wikiRoot}`);

const meta = loadMeta(proposalLocation.proposalDir);
if (meta) {
  if (meta.status) console.log(`状态: ${meta.status}`);
  if (meta.title) console.log(`标题: ${meta.title}`);
}

console.log('\n## Task Status\n');

const backend = config.taskBackend?.adapter || 'yaml';
console.log(`backend: ${backend}`);

if (backend === 'beads' || backend !== 'yaml') {
  _printBackendTaskStatus(config, projectRoot, meta, proposalLocation);
} else {
  console.log('\n(task backend unavailable, use flowforge task status)');
}

console.log('');

const notesPath = path.join(proposalLocation.proposalDir, 'notes.md');
if (fs.existsSync(notesPath)) {
  console.log('\n### notes.md\n');
  console.log(fs.readFileSync(notesPath, 'utf8'));
}

async function _printBackendTaskStatus(config, projectRoot, meta, loc) {
  const { createBackend } = require('./lib/backends');
  const backend = createBackend(config, projectRoot);
  const proposalId = meta ? meta.id : null;
  if (!proposalId) {
    console.log('(无法获取 proposal ID)');
    return;
  }

  try {
    const caps = backend.getCapabilities();
    console.log(`atomicClaim: ${caps.atomicClaim}, dependencySort: ${caps.dependencySort}`);

    const status = await backend.getStatus(proposalId);
    console.log(`\n总任务: ${status.total} | 完成: ${status.byStatus.done || 0} | 进行中: ${status.byStatus.in_progress || 0} | 待处理: ${status.byStatus.pending || 0} | 阻塞: ${status.byStatus.blocked || 0}`);

    if (status.byType && Object.keys(status.byType).length > 0) {
      console.log('');
      for (const [type, stats] of Object.entries(status.byType)) {
        console.log(`  ${type}: ${stats.done || 0}/${stats.total} done (in_progress: ${stats.inProgress || 0}, pending: ${stats.pending || 0}, blocked: ${stats.blocked || 0})`);
      }
    }

    const ready = await backend.getReadyTasks(proposalId);
    if (ready.length > 0) {
      console.log('\n## Ready Tasks');
      for (const t of ready) {
        console.log(`- [${t.id}] ${t.title} (${t.type})`);
      }
    }

    const blocked = await backend.getBlockedTasks(proposalId);
    if (blocked.length > 0) {
      console.log('\n## Blocked Tasks');
      for (const t of blocked) {
        console.log(`- [${t.id}] ${t.title} — ${t.blockReason || 'no reason'}`);
      }
    }
  } catch (e) {
    console.log(`(Backend query failed: ${e.message})`);
  }
}

function findProposalById(projectRoot, projects, id) {
  for (const p of projects) {
    const ws = path.join(projectRoot, p.wikiRoot, 'workspace');
    for (const sub of ['active', 'completed']) {
      const cand = path.join(ws, 'proposals', sub, id);
      if (fs.existsSync(cand)) return { proposalDir: cand, projectId: p.id, wikiRoot: p.wikiRoot };
    }
  }
  return null;
}

function findActiveProposal(projectRoot, projects) {
  for (const p of projects) {
    const activeDir = path.join(projectRoot, p.wikiRoot, 'workspace', 'proposals', 'active');
    if (!fs.existsSync(activeDir)) continue;
    const dirs = fs.readdirSync(activeDir, { withFileTypes: true }).filter(d => d.isDirectory());
    for (const d of dirs) {
      const pd = path.join(activeDir, d.name);
      const meta = loadMeta(pd);
      if (meta) {
        return { proposalDir: pd, projectId: p.id, wikiRoot: p.wikiRoot };
      }
    }
  }
  return null;
}
