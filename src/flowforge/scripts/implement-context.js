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
  proposalLocation = findActiveProposal(projectRoot, allProjects);
}

if (!proposalLocation) {
  console.error('ERROR: 未找到活跃状态的 proposal。请先在 design SKILL 中将 proposal 状态设为 active。');
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

const taskMapYamlPath = path.join(proposalLocation.proposalDir, 'task-map.yaml');
const taskMapMdPath = path.join(proposalLocation.proposalDir, 'task-map.md');
if (fs.existsSync(taskMapYamlPath)) {
  console.log('\n### task-map.yaml\n');
  console.log(fs.readFileSync(taskMapYamlPath, 'utf8'));
} else if (fs.existsSync(taskMapMdPath)) {
  console.log('\n### task-map.md (旧格式)\n');
  console.log(fs.readFileSync(taskMapMdPath, 'utf8'));
} else {
  console.log('\ntask-map: 不存在');
}

if (config.taskBackend && config.taskBackend.adapter) {
  const { createAdapter } = require('./lib/adapters');
  const adapter = createAdapter(config, projectRoot);
  const caps = adapter.getCapabilities();
  console.log('\n## Task Backend\n');
  console.log(`adapter: ${config.taskBackend.adapter}`);

  if (config.taskBackend.adapter === 'beads') {
    try {
      const { execSync } = require('child_process');
      const syncScript = path.join(__dirname, 'task-sync.js');
      const meta = loadMeta(proposalLocation.proposalDir);
      const pid = meta ? meta.id : null;
      if (pid) {
        const checkResult = execSync(
          `node "${syncScript}" "${projectRoot}" "${pid}" --check`,
          { encoding: 'utf8', stdio: 'pipe', timeout: 8000 }
        );
        const parsed = JSON.parse(checkResult);
        if (parsed.issues && parsed.issues.length > 0) {
          console.log(`sync: ${parsed.issues.length} 处不一致`);
          for (const issue of parsed.issues) {
            console.log(`  - ${issue}`);
          }
          console.log(`修复: node scripts/task-sync.js <root> ${pid}`);
        } else {
          console.log('sync: 数据一致');
        }
      }
    } catch (e) {
      if (e.stdout) {
        try {
          const parsed = JSON.parse(e.stdout.toString());
          if (parsed.issues && parsed.issues.length > 0) {
            console.log(`sync: ${parsed.issues.length} 处不一致（运行 task-sync.js 修复）`);
          }
        } catch (_) {}
      }
    }
  }

  console.log('');
}

const notesPath = path.join(proposalLocation.proposalDir, 'notes.md');
if (fs.existsSync(notesPath)) {
  console.log('\n### notes.md\n');
  console.log(fs.readFileSync(notesPath, 'utf8'));
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
      if (meta && meta.status === 'active') {
        return { proposalDir: pd, projectId: p.id, wikiRoot: p.wikiRoot };
      }
    }
  }
  return null;
}
