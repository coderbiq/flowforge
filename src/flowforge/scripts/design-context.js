#!/usr/bin/env node
'use strict';

const fs = require('fs');
const path = require('path');
const { loadMainConfig, loadProjectConfig, loadMeta } = require('./lib/config');

const projectRoot = process.argv[2] || process.cwd();
let proposalId = null;
let explicitProjectId = null;

for (let i = 3; i < process.argv.length; i++) {
  const arg = process.argv[i];
  if (arg === '--project') {
    explicitProjectId = process.argv[i + 1];
    i++;
  } else if (!arg.startsWith('--')) {
    proposalId = arg;
  }
}

const config = loadMainConfig(projectRoot);
if (!config) {
  console.error('ERROR: .flowforge/config.yaml 不存在或格式错误');
  process.exit(0);
}

const projectRefs = config.projects || [];
if (projectRefs.length === 0) {
  console.error('ERROR: config.yaml 中未定义 projects（至少需要一个 project）');
  process.exit(1);
}

const allProjects = [];
for (const ref of projectRefs) {
  const p = loadProjectConfig(projectRoot, ref);
  if (p) allProjects.push(p);
}

let activeProjectId = explicitProjectId;
let proposalLocation = null;

if (proposalId) {
  proposalLocation = findProposalById(projectRoot, allProjects, proposalId);
  if (proposalLocation && !activeProjectId) {
    activeProjectId = proposalLocation.projectId;
  }
} else if (!explicitProjectId) {
  proposalLocation = findActiveProposal(projectRoot, allProjects);
  if (proposalLocation) {
    activeProjectId = proposalLocation.projectId;
  }
}

const activeProject = activeProjectId ? allProjects.find(p => p.id === activeProjectId) : null;

console.log('# Design Context\n');

console.log('## Projects\n');
console.log(`共 ${allProjects.length} 个 project 配置：\n`);

for (const p of allProjects) {
  console.log(`### ${p.id}${p.name ? ` (${p.name})` : ''}`);
  console.log(`- wikiRoot: ${p.wikiRoot}`);
  if (p.srcDirs && p.srcDirs.length > 0) {
    console.log(`- srcDirs:`);
    for (const d of p.srcDirs) console.log(`  - ${d}`);
  }
  if (p.description) console.log(`- description: ${p.description}`);
  if (p.keywords && p.keywords.length > 0) {
    console.log(`- keywords: ${p.keywords.join(', ')}`);
  }
  console.log('');
}

if (!explicitProjectId) {
  const intakeEntries = collectIntake(projectRoot, allProjects);
  if (intakeEntries.length > 0) {
    console.log('## Intake Material\n');
    for (const entry of intakeEntries) {
      console.log(`- [project: ${entry.projectId}] ${entry.relPath}`);
      for (const f of entry.files) console.log(`    - ${f}`);
    }
    if (activeProject && activeProject.rules && activeProject.rules.intake) {
      console.log('\n分析步骤（来自当前 project）：\n');
      outputIntakeSteps(activeProject.rules.intake);
    }
    console.log('');
  }
}

if (activeProject && activeProject.rules) {
  const r = activeProject.rules;

  if (r.exploration && r.exploration.strategy) {
    console.log('## Exploration Strategy\n');
    console.log(r.exploration.strategy.trim());
    console.log('');
  }

  if (r.design) {
    console.log('## Design Rules\n');
    if (r.design.naming) {
      console.log('### Naming');
      console.log(`- proposal_id: ${r.design.naming.proposal_id}`);
      console.log(`- exploration_slug: ${r.design.naming.exploration_slug}`);
      console.log('');
    }
    if (r.design.task_rules) {
      console.log('### Task Rules');
      if (r.design.task_rules.fields) {
        console.log(`- fields: ${r.design.task_rules.fields.join(', ')}`);
      }
      if (r.design.task_rules.time_estimate) {
        console.log(`- time_estimate: ${r.design.task_rules.time_estimate}`);
      }
      console.log('');
    }
  }

  if (r.implement) {
    console.log('## Implement Rules\n');
    if (r.implement.task_states) {
      console.log(`- task_states: ${r.implement.task_states.join(', ')}`);
    }
    if (r.implement.notes && r.implement.notes.fields) {
      console.log(`- notes.fields: ${r.implement.notes.fields.join(', ')}`);
    }
    console.log('');
  }

  if (r.library) {
    console.log('## Library Rules\n');
    if (r.library.requireReview !== undefined) {
      console.log(`- requireReview: ${r.library.requireReview}`);
    }
    if (r.library.autoUpdateHistory !== undefined) {
      console.log(`- autoUpdateHistory: ${r.library.autoUpdateHistory}`);
    }
    console.log('');
  }

  if (activeProject.modules && Object.keys(activeProject.modules).length > 0) {
    console.log('## Modules\n');
    for (const [name, mod] of Object.entries(activeProject.modules)) {
      console.log(`- ${name} → ${mod.path}`);
    }
    console.log('');
  }
}

if (proposalLocation) {
  console.log('## Current Proposal\n');
  console.log(`路径: ${path.relative(projectRoot, proposalLocation.proposalDir)}`);
  console.log(`project: ${proposalLocation.projectId}`);
  console.log(`wikiRoot: ${proposalLocation.wikiRoot}`);
  console.log('(已锁定，无需重新决策 project 归属)');
  const meta = loadMeta(proposalLocation.proposalDir);
  if (meta) {
    if (meta.status) console.log(`状态: ${meta.status}`);
    if (meta.title) console.log(`标题: ${meta.title}`);
  }
  const taskMap = path.join(proposalLocation.proposalDir, 'task-map.md');
  if (fs.existsSync(taskMap)) console.log('task-map: 已有');
}

function outputIntakeSteps(intakeRules) {
  if (!intakeRules.steps || intakeRules.steps.length === 0) return;
  for (const s of intakeRules.steps) {
    console.log(`  - ${s.id}: ${s.action}`);
  }
}

function collectIntake(projectRoot, projects) {
  const result = [];
  for (const p of projects) {
    const intakeDir = path.join(projectRoot, p.wikiRoot, 'workspace', 'intake');
    if (!fs.existsSync(intakeDir)) continue;
    const files = fs.readdirSync(intakeDir).filter(f => f.endsWith('.md'));
    if (files.length === 0) continue;
    result.push({ projectId: p.id, relPath: path.relative(projectRoot, intakeDir), files });
  }
  return result;
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
      if (meta && (meta.status === 'draft' || meta.status === 'active')) {
        return { proposalDir: pd, projectId: p.id, wikiRoot: p.wikiRoot };
      }
    }
  }
  return null;
}
