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

// 扫描 proposal 目录下所有 .md 文件，提取 domain
const domainGroups = scanDomainGroups(proposalLocation.proposalDir);
if (domainGroups.length > 0) {
  console.log('\n## 归档目标（从文档 domain 推导）\n');
  for (const group of domainGroups) {
    console.log(`### → ${group.archivePath}`);
    console.log(`   domain: { scope: ${group.domain.scope}${group.domain.module ? ', module: ' + group.domain.module : ''}, type: ${group.domain.type} }`);
    console.log(`   来源文件:`);
    for (const f of group.files) {
      console.log(`    - ${f.relPath} (${f.title || '无标题'})`);
    }
    console.log('');
  }
} else {
  console.log('\n归档目标: 未在文档 frontmatter 中找到 domain 字段');
}

if (r && r.library) {
  console.log('## Library Rules\n');
  if (r.library.requireReview !== undefined) {
    console.log(`requireReview: ${r.library.requireReview}`);
  }
  if (r.library.autoUpdateHistory !== undefined) {
    console.log(`autoUpdateHistory: ${r.library.autoUpdateHistory}`);
  }
  console.log('');
  if (r.library.strategy) {
    console.log('## Library Strategy\n');
    console.log(r.library.strategy.trim());
    console.log('');
  }
}

if (r && r.archive && r.archive.strategy) {
  console.log('\n## Archive Strategy\n');
  console.log(r.archive.strategy.trim());
}

// 输出 library 现状（每个归档目标的已有文件状态）
if (domainGroups.length > 0) {
  console.log('\n## Library 现状（与归档目标对比）\n');
  for (const group of domainGroups) {
    const libFullPath = path.join(proposalLocation.wikiRoot, group.archivePath);
    const state = detectLibraryState(libFullPath);
    console.log(`### → ${group.archivePath}`);
    console.log(`   状态: ${state.status}`);
    if (state.status === 'exists') {
      console.log(`   已有大小: ${state.size} bytes`);
      if (state.hasArchiveNotes) {
        console.log(`   内容类型: 过时摘要（仅含 Archived proposal notes 段）`);
      } else {
        console.log(`   内容类型: 已有完整设计`);
      }
    } else if (state.status === 'directory_exists') {
      console.log(`   已有文件: ${state.existingFiles.join(', ') || '(空目录)'}`);
    } else {
      console.log(`   建议动作: 首次创建`);
    }
    console.log('');
  }
}

// 扫描 notes.md 中待提取的 knowledge 记录
const notesPath = path.join(proposalLocation.proposalDir, 'notes.md');
if (fs.existsSync(notesPath)) {
  const knowledgeEntries = scanNotesKnowledge(notesPath);
  if (knowledgeEntries.length > 0) {
    console.log('\n## notes.md 中待提取的 Knowledge 记录\n');
    for (const entry of knowledgeEntries) {
      console.log(`- [${entry.date}] ${entry.summary}`);
      if (entry.domain) {
        console.log(`  domain: { scope: ${entry.domain.scope}${entry.domain.module ? ', module: ' + entry.domain.module : ''}, type: ${entry.domain.type} }`);
      }
    }
  console.log(`\n共 ${knowledgeEntries.length} 条记录待归档时提取到 library。`);
  }
}

// 输出文档全文
for (const f of ['proposal.md', 'design.md']) {
  const fp = path.join(proposalLocation.proposalDir, f);
  if (fs.existsSync(fp)) {
    console.log(`\n### ${f}\n`);
    console.log(fs.readFileSync(fp, 'utf8'));
  }
}

// 任务快照仅提示存在，不输出内容——Agent 使用 flowforge task status 查询
const snapshotPath = path.join(proposalLocation.proposalDir, 'tasks.snapshot.md');
if (fs.existsSync(snapshotPath)) {
  console.log('\ntasks.snapshot.md: 已存在（使用 flowforge task status --proposal <id> 查看任务状态）');
}

// 输出 design/ 目录下的文件
const designDir = path.join(proposalLocation.proposalDir, 'design');
if (fs.existsSync(designDir) && fs.statSync(designDir).isDirectory()) {
  const designFiles = fs.readdirSync(designDir).filter(f => f.endsWith('.md'));
  for (const f of designFiles) {
    console.log(`\n### design/${f}\n`);
    console.log(fs.readFileSync(path.join(designDir, f), 'utf8'));
  }
}

// ============================================================
// Helpers
// ============================================================

function detectLibraryState(libFullPath) {
  if (!fs.existsSync(libFullPath)) {
    // 检查父目录是否存在（module 目录可能已存在但该文件还没创建）
    const parentDir = path.dirname(libFullPath);
    if (fs.existsSync(parentDir) && fs.statSync(parentDir).isDirectory()) {
      const files = fs.readdirSync(parentDir).filter(f => f.endsWith('.md'));
      return { status: 'directory_exists', existingFiles: files };
    }
    return { status: 'not_exists' };
  }

  const stat = fs.statSync(libFullPath);
  if (stat.isDirectory()) {
    const files = fs.readdirSync(libFullPath).filter(f => f.endsWith('.md'));
    return { status: 'directory_exists', existingFiles: files };
  }

  const content = fs.readFileSync(libFullPath, 'utf8');
  const hasArchiveNotes = /Archived\s+proposal\s+notes/i.test(content);
  return {
    status: 'exists',
    size: stat.size,
    hasArchiveNotes,
    contentSnippet: content.substring(0, 500)
  };
}

function scanNotesKnowledge(notesPath) {
  const content = fs.readFileSync(notesPath, 'utf8');
  const entries = [];
  const lines = content.split('\n');

  let currentDate = null;
  for (const line of lines) {
    const dateMatch = line.match(/^##\s+(\d{4}-\d{2}-\d{2})/);
    if (dateMatch) {
      currentDate = dateMatch[1];
      continue;
    }

    // 匹配 `| knowledge | <内容> |` 格式
    const knowledgeMatch = line.match(/\|\s*knowledge\s*\|\s*(.+?)(?:\s*\|.*)?$/);
    if (knowledgeMatch && currentDate) {
      const summary = knowledgeMatch[1].trim();
      const domain = extractDomainFromLine(line);
      entries.push({ date: currentDate, summary, domain });
    }
  }

  return entries;
}

function extractDomainFromLine(line) {
  const m = line.match(/domain\s*:\s*\{([^}]+)\}/);
  if (!m) return null;
  const parts = m[1].split(',').map(p => p.trim());
  const domain = {};
  for (const part of parts) {
    const kv = part.match(/^(\w+)\s*:\s*(.+)/);
    if (kv) domain[kv[1]] = kv[2].trim().replace(/^["']|["']$/g, '');
  }
  if (!domain.scope || !domain.type) return null;
  return domain;
}

function scanDomainGroups(proposalDir) {
  const files = collectMdFiles(proposalDir);
  const groups = {};

  for (const f of files) {
    const content = fs.readFileSync(f.absPath, 'utf8');
    const fm = extractFrontmatter(content);
    if (!fm || !fm.domain) continue;

    const domain = parseDomain(fm.domain);
    if (!domain) continue;

    const archivePath = deriveArchivePath(domain, f.relPath, fm.title || path.basename(f.absPath, '.md'));
    if (!groups[archivePath]) {
      groups[archivePath] = { domain, archivePath, files: [] };
    }
    groups[archivePath].files.push({
      relPath: f.relPath,
      title: fm.title || null,
      absPath: f.absPath,
      docType: fm.doc_type || null
    });
  }

  return Object.values(groups);
}

function collectMdFiles(dir, baseDir) {
  if (!baseDir) baseDir = dir;
  const results = [];
  const entries = fs.readdirSync(dir, { withFileTypes: true });

  for (const entry of entries) {
    const fullPath = path.join(dir, entry.name);
    if (entry.isDirectory() && !entry.name.startsWith('.')) {
      results.push(...collectMdFiles(fullPath, baseDir));
    } else if (entry.isFile() && entry.name.endsWith('.md')) {
      results.push({
        absPath: fullPath,
        relPath: path.relative(baseDir, fullPath)
      });
    }
  }

  return results;
}

function parseDomain(domainStr) {
  // domain 可能是 YAML 多行字符串或已解析的对象
  if (typeof domainStr === 'object') return domainStr;

  const lines = domainStr.split('\n').map(l => l.trim()).filter(Boolean);
  const result = {};
  for (const line of lines) {
    const m = line.match(/^(\w+)\s*:\s*(.+)/);
    if (m) result[m[1]] = m[2].trim();
  }
  if (!result.scope || !result.type) return null;
  return result;
}

function deriveArchivePath(domain, sourceRelPath, title) {
  const topic = title.replace(/[^a-zA-Z0-9\u4e00-\u9fff_-]/g, '-').replace(/-+/g, '-').replace(/^-|-$/g, '').toLowerCase();

  if (domain.scope === 'system') {
    switch (domain.type) {
      case 'design':
        return `library/architecture/${topic}.md`;
      case 'decision':
        return `library/decisions/${topic}.md`;
      case 'convention':
        return `library/conventions/${topic}.md`;
    }
  }

  if (domain.scope === 'module' && domain.module) {
    // 使用来源文件的相对路径派生 library 路径，保持目录结构
    // e.g. design/architecture.md → library/modules/data-service/architecture.md
    // e.g. model/DpDataSource.md → library/modules/data-service/model/DpDataSource.md
    return `library/modules/${domain.module}/${sourceRelPath}`;
  }

  return null;
}

function extractFrontmatter(text) {
  const m = text.match(/^---\n([\s\S]*?)\n---/);
  if (!m) return null;
  const result = {};
  let currentKey = null;

  for (const line of m[1].split('\n')) {
    // 嵌套属性（domain 的子字段）
    const nested = line.match(/^\s{2}(\w+)\s*:\s*(.*)/);
    if (nested && currentKey === 'domain') {
      if (!result.domain) result.domain = {};
      result.domain[nested[1]] = nested[2].trim().replace(/^["']|["']$/g, '');
      continue;
    }

    const kv = line.match(/^\s*([a-zA-Z_]+)\s*:\s*(.*)/);
    if (kv) {
      currentKey = kv[1];
      const val = kv[2].trim().replace(/^["']|["']$/g, '');
      result[currentKey] = val || '';
    } else {
      currentKey = null;
    }
  }
  return result;
}

function findProposalById(projectRoot, projects, id) {
  for (const p of projects) {
    const ws = path.join(projectRoot, p.wikiRoot, 'workspace');
    for (const sub of ['active', 'completed']) {
      const subDir = path.join(ws, 'proposals', sub);
      if (!fs.existsSync(subDir)) continue;
      const dirs = fs.readdirSync(subDir, { withFileTypes: true }).filter(d => d.isDirectory());
      for (const d of dirs) {
        if (d.name === id || d.name.startsWith(id + '-')) {
          return { proposalDir: path.join(subDir, d.name), projectId: p.id, wikiRoot: p.wikiRoot };
        }
      }
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
