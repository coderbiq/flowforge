#!/usr/bin/env node
'use strict';

const fs = require('fs');
const path = require('path');
const { loadMainConfig, loadProjectConfig, loadMeta } = require('./lib/config');

const projectRoot = process.argv[2] || process.cwd();
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
  proposalLocation = findProposal(projectRoot, allProjects, ['implemented', 'archived']);
}

if (!proposalLocation) {
  console.error('ERROR: 未找到 implemented 或 archived 状态的 proposal');
  process.exit(1);
}

const activeProject = allProjects.find(p => p.id === proposalLocation.projectId);
const r = activeProject && activeProject.rules ? activeProject.rules : null;

const meta = loadMeta(proposalLocation.proposalDir);

console.log('# Archive Context\n');

console.log('## Current Proposal\n');
console.log(`路径: ${path.relative(projectRoot, proposalLocation.proposalDir)}`);
console.log(`project: ${proposalLocation.projectId}`);
console.log(`wikiRoot: ${proposalLocation.wikiRoot}`);
if (meta) {
  if (meta.id) console.log(`ID: ${meta.id}`);
  if (meta.status) console.log(`状态: ${meta.status}`);
  if (meta.title) console.log(`标题: ${meta.title}`);
}

if (meta && meta.archive_targets && meta.archive_targets.length > 0) {
  console.log('\n### Archive Targets\n');
  for (const t of meta.archive_targets) {
    console.log(`- type: ${t.type}, ref: ${t.ref}${t.role ? ', role: ' + t.role : ''}`);
  }
} else {
  console.log('\narchive_targets: 未配置');
}

if (r && r.library) {
  console.log('\n## Library Rules\n');
  if (r.library.requireReview !== undefined) {
    console.log(`requireReview: ${r.library.requireReview}`);
  }
  if (r.library.autoUpdateHistory !== undefined) {
    console.log(`autoUpdateHistory: ${r.library.autoUpdateHistory}`);
  }
}

if (activeProject.modules && Object.keys(activeProject.modules).length > 0) {
  console.log('\n## Module Registry\n');
  for (const [name, mod] of Object.entries(activeProject.modules)) {
    const aliases = mod.aliases && mod.aliases.length > 0 ? ` (别名: ${mod.aliases.join(', ')})` : '';
    console.log(`- ${name} → ${mod.path}${aliases}`);
  }
}

for (const f of ['proposal.md', 'design.md', 'task-map.md']) {
  const fp = path.join(proposalLocation.proposalDir, f);
  if (fs.existsSync(fp)) {
    console.log(`\n### ${f}\n`);
    console.log(fs.readFileSync(fp, 'utf8'));
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

function findProposal(projectRoot, projects, statuses) {
  for (const p of projects) {
    const wikiWs = path.join(projectRoot, p.wikiRoot, 'workspace');
    for (const sub of ['active', 'completed']) {
      const subDir = path.join(wikiWs, 'proposals', sub);
      if (!fs.existsSync(subDir)) continue;
      const dirs = fs.readdirSync(subDir, { withFileTypes: true }).filter(d => d.isDirectory());
      for (const d of dirs) {
        const pd = path.join(subDir, d.name);
        const meta = loadMeta(pd);
        if (meta && statuses.includes(meta.status)) {
          return { proposalDir: pd, projectId: p.id, wikiRoot: p.wikiRoot };
        }
      }
    }
  }
  return null;
}
