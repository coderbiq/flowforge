#!/usr/bin/env node
'use strict';

const fs = require('fs');
const path = require('path');
const { loadMainConfig, loadProjectConfig, loadMeta } = require('./lib/config');

const projectRoot = process.argv[2] || process.cwd();
const crId = process.argv[3];
const type = process.argv[4];    // bug | finding | knowledge | missing-requirement | design-flaw
const title = process.argv[5];
const content = process.argv[6] || '';

if (!crId || !type || !title) {
  console.error('用法: feedback-capture.js <projectRoot> <CR-id> <type> <title> [content]');
  console.error('  type: bug | finding | knowledge | missing-requirement | design-flaw');
  process.exit(1);
}

if (!['bug', 'finding', 'knowledge', 'missing-requirement', 'design-flaw'].includes(type)) {
  console.error(`ERROR: 无效的 type "${type}"，有效值: bug, finding, knowledge, missing-requirement, design-flaw`);
  process.exit(1);
}

// 定位 proposal
const config = loadMainConfig(projectRoot);
if (!config) {
  console.error('ERROR: .flowforge/config.yaml 不存在或格式错误');
  process.exit(1);
}

const projectRefs = config.projects || [];
const allProjects = [];
for (const ref of projectRefs) {
  const p = loadProjectConfig(projectRoot, ref);
  if (p) allProjects.push(p);
}

const proposalLocation = findProposalById(projectRoot, allProjects, crId);
if (!proposalLocation) {
  console.error(`ERROR: 未找到 proposal "${crId}"`);
  process.exit(1);
}

const meta = loadMeta(proposalLocation.proposalDir);
const wikiRoot = path.join(projectRoot, proposalLocation.wikiRoot);
const proposalDir = proposalLocation.proposalDir;

// 处理不同类型的反馈
switch (type) {
  case 'bug':
    handleBug(proposalDir, crId, title, content);
    break;
  case 'finding':
    handleFinding(proposalDir, wikiRoot, meta, crId, title, content);
    break;
  case 'knowledge':
    handleKnowledge(proposalDir, title, content);
    break;
  case 'missing-requirement':
    handleMissingRequirement(crId, title, content);
    break;
  case 'design-flaw':
    handleDesignFlaw(crId, title, content);
    break;
  default:
    console.error(`ERROR: 未知 type "${type}"`);
    process.exit(1);
}

// --- Handler functions ---

function handleBug(proposalDir, crId, title, content) {
  const notesPath = path.join(proposalDir, 'notes.md');
  const now = new Date();
  const dateStr = now.toISOString().split('T')[0];
  const timeStr = now.toTimeString().slice(0, 5);

  const entry = `\n${timeStr} | bug | ${title}` +
    (content ? `\n     | 根因: ${content}` : '') + '\n';

  // 追加到 notes.md
  if (fs.existsSync(notesPath)) {
    const existingContent = fs.readFileSync(notesPath, 'utf8');
    if (existingContent.includes(`## ${dateStr}`)) {
      // 插入到当天段落的末尾
      const sections = existingContent.split(/(?=^## \d{4}-\d{2}-\d{2})/m);
      let found = false;
      for (let i = 0; i < sections.length; i++) {
        if (sections[i].startsWith(`## ${dateStr}`)) {
          sections[i] = sections[i].trimEnd() + entry;
          found = true;
          break;
        }
      }
      if (found) {
        fs.writeFileSync(notesPath, sections.join(''));
      } else {
        fs.appendFileSync(notesPath, `\n## ${dateStr}${entry}`);
      }
    } else {
      fs.appendFileSync(notesPath, `\n## ${dateStr}${entry}`);
    }
  } else {
    // notes.md 不存在，创建
    const frontmatter = `---
doc_type: notes
title: ${meta ? meta.title : crId} 实施日志
status: active
note_kind: progress
---

`;
    fs.writeFileSync(notesPath, frontmatter + `## ${dateStr}${entry}`);
  }

  // 创建修复任务
  const repairTitle = `[修复] ${title}`;
  const repairDesc = content || title;
  try {
    const { execSync } = require('child_process');
    const discoverScript = path.join(__dirname, 'task-discover.js');
    execSync(
      `node "${discoverScript}" "${projectRoot}" "${crId}" "0" "${repairTitle}" "${repairDesc}"`,
      { encoding: 'utf8', stdio: 'pipe', timeout: 5000 }
    );
  } catch (e) {
    console.log('(task-discover.js 执行失败，已记录 bug 到 notes.md）');
  }

  console.log(`[bug] "${title}" 已写入 notes.md，修复任务已创建`);
}

function handleFinding(proposalDir, wikiRoot, meta, crId, title, content) {
  // 查找关联的 exploration
  const explorationsDir = path.join(wikiRoot, 'workspace', 'explorations');
  let targetExpDir = null;
  let targetExpSlug = null;

  if (meta && meta.source_explorations && meta.source_explorations.length > 0) {
    // 优先使用第一个关联的 exploration
    for (const src of meta.source_explorations) {
      const expDir = path.join(explorationsDir, src.ref);
      if (fs.existsSync(expDir)) {
        targetExpDir = expDir;
        targetExpSlug = src.ref;
        break;
      }
    }
  }

  if (!targetExpDir) {
    // 没有关联 exploration，检查是否有任何 exploration
    if (fs.existsSync(explorationsDir)) {
      const existing = fs.readdirSync(explorationsDir, { withFileTypes: true })
        .filter(d => d.isDirectory());
      if (existing.length > 0) {
        targetExpDir = path.join(explorationsDir, existing[0].name);
        targetExpSlug = existing[0].name;
      }
    }
  }

  if (!targetExpDir) {
    console.log('[finding] 无关联 exploration 目录，需要在 flowforge-design 中先创建 exploration 再写入。');
    console.log(`  建议的 exploration slug: ${slugify(title)}`);
    console.log('  创建 exploration 后重新运行此命令。');
    return;
  }

  // 确定 finding ID
  const findingsDir = path.join(targetExpDir, 'findings');
  if (!fs.existsSync(findingsDir)) {
    fs.mkdirSync(findingsDir, { recursive: true });
  }

  const existingFindings = fs.readdirSync(findingsDir)
    .filter(f => f.match(/^F-\d+\.md$/))
    .map(f => parseInt(f.match(/F-(\d+)/)[1], 10));
  const nextNum = existingFindings.length > 0 ? Math.max(...existingFindings) + 1 : 1;
  const findingId = `F-${String(nextNum).padStart(3, '0')}`;

  // 推断 domain
  let domainScope = 'system';
  let domainModule = '';
  let domainType = 'design';
  if (meta && meta.modules && meta.modules.length > 0) {
    domainScope = 'module';
    domainModule = meta.modules[0];
  }

  const findingContent = `---
doc_type: finding
title: ${title}
status: active
finding_id: ${findingId}
source: implementation
source_proposal: ${crId}
domain:
  scope: ${domainScope}
  module: ${domainModule || ''}
  type: ${domainType}
---

# ${title}

## 发现

${content || title}

## 证据

- 来自提案 ${crId} 的实施过程
- 发现于 ${targetExpSlug} 关联的代码区域
`;

  const findingPath = path.join(findingsDir, `${findingId}.md`);
  fs.writeFileSync(findingPath, findingContent);

  // 检查并更新 exploration 的 status（如果已归档则重新激活）
  const explorationIndexPath = path.join(targetExpDir, 'index.md');
  if (fs.existsSync(explorationIndexPath)) {
    let indexContent = fs.readFileSync(explorationIndexPath, 'utf8');
    if (indexContent.includes('status: archived')) {
      indexContent = indexContent.replace('status: archived', 'status: active');
      fs.writeFileSync(explorationIndexPath, indexContent);
      console.log(`  已将 exploration "${targetExpSlug}" 状态从 archived 改为 active`);
    }
  }

  console.log(`[finding] "${title}" → ${path.relative(projectRoot, findingPath)}`);
  console.log(`  finding_id: ${findingId}, exploration: ${targetExpSlug}`);
}

function handleKnowledge(proposalDir, title, content) {
  const notesPath = path.join(proposalDir, 'notes.md');
  const now = new Date();
  const dateStr = now.toISOString().split('T')[0];
  const timeStr = now.toTimeString().slice(0, 5);

  const entry = `\n${timeStr} | knowledge | ${title}` +
    (content ? `\n     | ${content}` : '') +
    '\n     | note: 待 flowforge-archive 提取到 library\n';

  if (fs.existsSync(notesPath)) {
    const existingContent = fs.readFileSync(notesPath, 'utf8');
    if (existingContent.includes(`## ${dateStr}`)) {
      const sections = existingContent.split(/(?=^## \d{4}-\d{2}-\d{2})/m);
      let found = false;
      for (let i = 0; i < sections.length; i++) {
        if (sections[i].startsWith(`## ${dateStr}`)) {
          sections[i] = sections[i].trimEnd() + entry;
          found = true;
          break;
        }
      }
      if (found) {
        fs.writeFileSync(notesPath, sections.join(''));
      } else {
        fs.appendFileSync(notesPath, `\n## ${dateStr}${entry}`);
      }
    } else {
      fs.appendFileSync(notesPath, `\n## ${dateStr}${entry}`);
    }
  } else {
    const frontmatter = `---
doc_type: notes
title: ${meta ? meta.title : crId} 实施日志
status: active
note_kind: progress
---

`;
    fs.writeFileSync(notesPath, frontmatter + `## ${dateStr}${entry}`);
  }

  console.log(`[knowledge] "${title}" 已写入 notes.md`);
  console.log('  待 flowforge-archive 时提取到 library');
}

function handleMissingRequirement(crId, title, content) {
  console.log(`[missing-requirement] "${title}"`);
  if (content) console.log(`  描述: ${content}`);
  console.log('');
  console.log('## 路由指引');
  console.log(`激活 flowforge-design，在 proposal "${crId}" 中补充设计：`);
  console.log(`  1. 在 design/ 下补充对应模块的设计文档`);
  console.log(`  2. 如需新的探索，创建新的 exploration`);
  console.log(`  3. 通过 task-cancel.js 废弃受影响任务`);
  console.log(`  4. 通过 task-add.js 添加新任务`);
}

function handleDesignFlaw(crId, title, content) {
  console.log(`[design-flaw] "${title}"`);
  if (content) console.log(`  描述: ${content}`);
  console.log('');
  console.log('## 路由指引');
  console.log(`激活 flowforge-design 修改方案：`);
  console.log(`  1. 说明缺陷所在的任务和需要修正的设计点`);
  console.log(`  2. 在 design/ 下修改对应的设计文档`);
  console.log(`  3. 通过 task-cancel.js 废弃受影响任务`);
  console.log(`  4. 通过 task-add.js 添加修正任务`);
}

// --- Helpers ---

function slugify(text) {
  return text
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '')
    .slice(0, 50);
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
