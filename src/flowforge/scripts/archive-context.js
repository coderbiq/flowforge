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
}

// 输出文档全文
for (const f of ['proposal.md', 'design.md', 'task-map.md']) {
  const fp = path.join(proposalLocation.proposalDir, f);
  if (fs.existsSync(fp)) {
    console.log(`\n### ${f}\n`);
    console.log(fs.readFileSync(fp, 'utf8'));
  }
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

function scanDomainGroups(proposalDir) {
  const files = collectMdFiles(proposalDir);
  const groups = {};

  for (const f of files) {
    const content = fs.readFileSync(f.absPath, 'utf8');
    const fm = extractFrontmatter(content);
    if (!fm || !fm.domain) continue;

    const domain = parseDomain(fm.domain);
    if (!domain) continue;

    const archivePath = deriveArchivePath(domain, fm.title || path.basename(f.absPath, '.md'));
    if (!groups[archivePath]) {
      groups[archivePath] = { domain, archivePath, files: [] };
    }
    groups[archivePath].files.push({
      relPath: f.relPath,
      title: fm.title || null,
      absPath: f.absPath
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

function deriveArchivePath(domain, title) {
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
    return `library/modules/${domain.module}/`;
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
