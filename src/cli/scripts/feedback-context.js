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
let blockedTasks = [];
const backend = config.taskBackend?.adapter || 'yaml';

if (backend !== 'yaml') {
  try {
    const { createBackend } = require('./lib/backends');
    const be = createBackend(config, projectRoot);
    const meta = loadMeta(proposalLocation.proposalDir);
    const proposalId = meta ? meta.id : null;
    if (proposalId) {
      blockedTasks = await be.getBlockedTasks(proposalId);
    }
  } catch (_) { console.log('(无法查询后端 blocked 任务)'); }
}

if (blockedTasks.length === 0) {
  console.log('无 blocked 任务\n');
} else {
  for (const t of blockedTasks) {
    const label = t.id || `T${t.id}`;
    console.log(`- ${label}: ${t.title}`);
    if (t.blockReason || t.blocked_reason) console.log(`  原因: ${t.blockReason || t.blocked_reason}`);
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

// 关联的 library 文档
console.log('## Related Library Documents\n');
const wikiRoot = path.join(projectRoot, proposalLocation.wikiRoot);
const libDir = path.join(wikiRoot, 'library');

if (fs.existsSync(libDir)) {
  const modules = meta && meta.modules ? meta.modules : [];
  const libModulesDir = path.join(libDir, 'modules');

  for (const mod of modules) {
    const modDir = path.join(libModulesDir, mod);
    if (fs.existsSync(modDir)) {
      const modFiles = fs.readdirSync(modDir, { recursive: true })
        .filter(f => f.endsWith('.md'))
        .map(f => `modules/${mod}/${f}`);
      console.log(`- modules/${mod}/ 已有文档: ${modFiles.length > 0 ? modFiles.join(', ') : '(空)'}`);
    } else {
      console.log(`- modules/${mod}/ (目录不存在，新发现将写入此路径)`);
    }
  }

  // 列出 library 顶层目录
  for (const sub of ['architecture', 'decisions', 'conventions']) {
    const subDir = path.join(libDir, sub);
    if (fs.existsSync(subDir)) {
      const files = fs.readdirSync(subDir).filter(f => f.endsWith('.md'));
      if (files.length > 0) {
        console.log(`- ${sub}/ 已有 ${files.length} 个文档`);
      }
    }
  }

  if (modules.length === 0) {
    console.log('(无关联模块，系统级发现将写入 library/architecture/)');
  }
} else {
  console.log('library/ 目录不存在，新发现将自动创建。');
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
