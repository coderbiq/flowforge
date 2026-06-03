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
  process.exit(1);
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

// 定位 proposal
let proposalLocation = null;
if (proposalId) {
  proposalLocation = findProposalById(projectRoot, allProjects, proposalId);
} else {
  proposalLocation = findActiveProposal(projectRoot, allProjects);
}

if (!proposalLocation) {
  console.error('ERROR: 未找到活跃状态的 proposal');
  process.exit(1);
}

const meta = loadMeta(proposalLocation.proposalDir);

const activeProject = allProjects.find(p => p.id === proposalLocation.projectId);
const r = activeProject && activeProject.rules ? activeProject.rules : null;

console.log('# Feedback Context\n');

if (r && r.feedback && r.feedback.strategy) {
  console.log('## Feedback Strategy\n');
  console.log(r.feedback.strategy.trim());
  console.log('');
}

// 当前 proposal
console.log('## Current Proposal\n');
console.log(`路径: ${path.relative(projectRoot, proposalLocation.proposalDir)}`);
console.log(`project: ${proposalLocation.projectId}`);
console.log(`wikiRoot: ${proposalLocation.wikiRoot}`);
if (meta) {
  if (meta.status) console.log(`状态: ${meta.status}`);
  if (meta.title) console.log(`标题: ${meta.title}`);
  if (meta.id) console.log(`ID: ${meta.id}`);
}
console.log('');

// Blocked 任务
console.log('## Blocked Tasks\n');
const taskMapPath = path.join(proposalLocation.proposalDir, 'task-map.yaml');
let blockedTasks = [];
if (fs.existsSync(taskMapPath)) {
  try {
    const yaml = require('./vendor/js-yaml');
    const taskMap = yaml.load(fs.readFileSync(taskMapPath, 'utf8'));
    if (taskMap && taskMap.tasks) {
      blockedTasks = taskMap.tasks.filter(t => t.status === 'blocked');
    }
  } catch (e) {
    console.log('(task-map.yaml 解析失败)');
  }
}

if (blockedTasks.length === 0) {
  console.log('无 blocked 任务\n');
} else {
  for (const t of blockedTasks) {
    console.log(`- T${t.id}: ${t.title}`);
    if (t.blocked_reason) console.log(`  原因: ${t.blocked_reason}`);
  }
  console.log('');
}

// Notes 中的问题和 blocked 记录
console.log('## Notes Summary\n');
const notesPath = path.join(proposalLocation.proposalDir, 'notes.md');
if (fs.existsSync(notesPath)) {
  const notesContent = fs.readFileSync(notesPath, 'utf8');
  // 提取 blocked / bug / 问题相关行
  const lines = notesContent.split('\n');
  const relevantLines = [];
  let foundBlocked = false;
  for (const line of lines) {
    const trimmed = line.trim();
    if (trimmed.includes('| blocked |') || trimmed.includes('| bug |') ||
        trimmed.includes('TODO:') || trimmed.includes('FIXME:') ||
        trimmed.includes('问题') || trimmed.includes('失败') ||
        trimmed.includes('错误')) {
      relevantLines.push(trimmed);
      if (trimmed.includes('| blocked |') || trimmed.includes('| bug |')) {
        foundBlocked = true;
      }
    }
  }
  if (relevantLines.length === 0) {
    console.log('notes.md 中无 blocked/bug/问题 记录\n');
  } else {
    for (const line of relevantLines) {
      console.log(`  ${line}`);
    }
    console.log('');
  }

  // 统计 note_kind
  const bugCount = (notesContent.match(/\| bug \|/g) || []).length;
  const blockedCount = (notesContent.match(/\| blocked \|/g) || []).length;
  if (bugCount > 0 || blockedCount > 0) {
    console.log(`统计: ${blockedCount} 条 blocked 记录, ${bugCount} 条 bug 记录\n`);
  }

  // 检查是否有未消费的 knowledge 标记
  const knowledgeCount = (notesContent.match(/\| knowledge \|/g) || []).length;
  if (knowledgeCount > 0) {
    console.log(`待提取的 knowledge 记录: ${knowledgeCount} 条\n`);
  }
} else {
  console.log('notes.md 不存在\n');
}

// 关联的 explorations
console.log('## Associated Explorations\n');
const wikiRoot = path.join(projectRoot, proposalLocation.wikiRoot);
const explorationsDir = path.join(wikiRoot, 'workspace', 'explorations');

if (meta && meta.source_explorations && meta.source_explorations.length > 0) {
  for (const src of meta.source_explorations) {
    const expDir = path.join(explorationsDir, src.ref);
    if (fs.existsSync(expDir)) {
      const findingsDir = path.join(expDir, 'findings');
      const decisionsDir = path.join(expDir, 'decisions');
      const findings = fs.existsSync(findingsDir)
        ? fs.readdirSync(findingsDir).filter(f => f.endsWith('.md'))
        : [];
      const decisions = fs.existsSync(decisionsDir)
        ? fs.readdirSync(decisionsDir).filter(f => f.endsWith('.md'))
        : [];
      console.log(`- ${src.ref}`);
      console.log(`  findings: ${findings.length > 0 ? findings.join(', ') : '(空)'}`);
      console.log(`  decisions: ${decisions.length > 0 ? decisions.join(', ') : '(空)'}`);
    } else {
      console.log(`- ${src.ref} (目录不存在)`);
    }
  }
}

// 如果没有关联 exploration，检查是否有匹配的
if (!meta || !meta.source_explorations || meta.source_explorations.length === 0) {
  if (fs.existsSync(explorationsDir)) {
    const allExplorations = fs.readdirSync(explorationsDir, { withFileTypes: true })
      .filter(d => d.isDirectory());
    if (allExplorations.length > 0) {
      console.log('无关联 exploration，但 workspace/explorations/ 下存在以下目录：');
      for (const exp of allExplorations) {
        console.log(`  - ${exp.name}`);
      }
      console.log('如需关联，在 meta.yaml 的 source_explorations 中添加引用。');
    } else {
      console.log('无关联 exploration，workspace/explorations/ 下也无现有目录。');
    }
  } else {
    console.log('无关联 exploration。');
  }
}
console.log('');

// 建议的反馈项
console.log('## Suggested Feedback Items\n');
let suggestions = 0;

// 基于 blocked 任务
for (const t of blockedTasks) {
  suggestions++;
  console.log(`[${suggestions}] [design-flaw/bug] T${t.id} 被阻塞: ${t.title}`);
  if (t.blocked_reason) console.log(`    原因: ${t.blocked_reason}`);
}

// 基于 notes.md blocked 记录（未关联到具体任务的）
if (fs.existsSync(notesPath)) {
  const notesContent = fs.readFileSync(notesPath, 'utf8');
  const blockedPattern = /\d{2}:\d{2}\s*\|\s*blocked\s*\|(.+)/g;
  let match;
  while ((match = blockedPattern.exec(notesContent)) !== null) {
    suggestions++;
    console.log(`[${suggestions}] [finding/knowledge] blocked 记录: ${match[1].trim()}`);
  }
}

if (suggestions === 0) {
  console.log('当前无自动检测到的反馈项。');
  console.log('如果你在实施/测试中发现了值得记录的东西，手动指定类型和内容运行 feedback-capture.js。');
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
