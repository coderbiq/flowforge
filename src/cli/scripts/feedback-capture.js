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
  console.error('用法: flowforge feedback-capture <CR-id> <type> <title> [content]');
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
    execSync(
      `flowforge task add --proposal "${crId}" implementation "${repairTitle}" --desc "${repairDesc}"`,
      { encoding: 'utf8', stdio: 'pipe', timeout: 5000, cwd: projectRoot }
    );
  } catch (e) {
    console.log('(flowforge task add 执行失败，已记录 bug 到 notes.md）');
  }

  console.log(`[bug] "${title}" 已写入 notes.md，修复任务已创建`);
}

function handleFinding(proposalDir, wikiRoot, meta, crId, title, content) {
  // 推断 domain
  let domainScope = 'system';
  let domainModule = '';
  let domainType = 'design';
  if (meta && meta.modules && meta.modules.length > 0) {
    domainScope = 'module';
    domainModule = meta.modules[0];
  }

  // 推导 library 目标路径
  const topic = title
    .replace(/[^a-zA-Z0-9\u4e00-\u9fff_-]/g, '-')
    .replace(/-+/g, '-')
    .replace(/^-|-$/g, '')
    .toLowerCase();

  let archivePath;
  if (domainScope === 'system') {
    archivePath = `library/architecture/${topic}.md`;
  } else {
    archivePath = `library/modules/${domainModule}/findings/${topic}.md`;
  }

  const libFullPath = path.join(wikiRoot, archivePath);
  const libDir = path.dirname(libFullPath);
  if (!fs.existsSync(libDir)) {
    fs.mkdirSync(libDir, { recursive: true });
  }

  const findingContent = `---
doc_type: finding
title: ${title}
status: active
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
`;

  fs.writeFileSync(libFullPath, findingContent);

  console.log(`[finding] "${title}" → ${path.relative(projectRoot, libFullPath)}`);
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
  console.log(`  2. 将新发现的设计事实写入 library/ 对应路径`);
  console.log(`  3. 通过 flowforge task cancel 废弃受影响任务`);
  console.log(`  4. 通过 flowforge task add 添加新任务`);
}

function handleDesignFlaw(crId, title, content) {
  console.log(`[design-flaw] "${title}"`);
  if (content) console.log(`  描述: ${content}`);
  console.log('');
  console.log('## 路由指引');
  console.log(`激活 flowforge-design 修改方案：`);
  console.log(`  1. 说明缺陷所在的任务和需要修正的设计点`);
  console.log(`  2. 在 design/ 下修改对应的设计文档`);
  console.log(`  3. 通过 flowforge task cancel 废弃受影响任务`);
  console.log(`  4. 通过 flowforge task add 添加修正任务`);
}

// --- Helpers ---

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
