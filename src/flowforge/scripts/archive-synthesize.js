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
const meta = loadMeta(proposalLocation.proposalDir);
const domainGroups = scanDomainGroups(proposalLocation.proposalDir);
const wikiRootFull = path.join(projectRoot, proposalLocation.wikiRoot);

const plan = {
  proposal: {
    id: meta ? meta.id : 'unknown',
    title: meta ? meta.title : 'unknown',
    status: meta ? meta.status : 'unknown',
    path: path.relative(projectRoot, proposalLocation.proposalDir),
    wikiRoot: proposalLocation.wikiRoot,
    projectId: proposalLocation.projectId,
    isCompleted: proposalLocation.proposalDir.includes('/completed/')
  },
  libraryRules: {},
  targets: []
};

if (activeProject && activeProject.rules && activeProject.rules.library) {
  plan.libraryRules = {
    requireReview: activeProject.rules.library.requireReview,
    autoUpdateHistory: activeProject.rules.library.autoUpdateHistory,
    strategy: activeProject.rules.library.strategy || null
  };
}

if (activeProject && activeProject.rules && activeProject.rules.archive) {
  plan.archiveStrategy = activeProject.rules.archive.strategy || null;
}

for (const group of domainGroups) {
  const libFullPath = path.join(wikiRootFull, group.archivePath);
  const state = detectLibraryState(libFullPath);
  const synthesisAction = classifySynthesis(group, state);

  const target = {
    archivePath: group.archivePath,
    domain: group.domain,
    action: synthesisAction.action,
    reason: synthesisAction.reason,
    sourceFiles: group.files.map(f => ({
      relPath: f.relPath,
      title: f.title,
      docType: f.docType
    })),
    libraryState: {
      status: state.status,
      size: state.size || 0,
      hasArchiveNotes: state.hasArchiveNotes || false,
      existingFiles: state.existingFiles || []
    },
    instructions: synthesisAction.instructions
  };

  plan.targets.push(target);
}

console.log(JSON.stringify(plan, null, 2));

// ============================================================
// Synthesis Classification
// ============================================================

function classifySynthesis(group, state) {
  if (state.status === 'not_exists') {
    return {
      action: 'create',
      reason: 'library 中尚无对应文档',
      instructions: {
        steps: [
          '按 doc_type 对应的 writing guide 格式化来源文档内容',
          '设置 frontmatter（doc_type, title, status=active, domain, created, updated）',
          '写入目标文件，必要时创建父目录',
          '运行 validate-doc.js 校验'
        ]
      }
    };
  }

  if (state.status === 'directory_exists') {
    const sourceFilenames = group.files.map(f => path.basename(f.relPath));
    const newFiles = sourceFilenames.filter(f => !state.existingFiles.includes(f));
    const existingFiles = sourceFilenames.filter(f => state.existingFiles.includes(f));

    return {
      action: 'mixed',
      reason: `模块目录已存在，${newFiles.length} 个文件需创建，${existingFiles.length} 个文件需检查合并`,
      instructions: {
        newFiles: newFiles.map(f => ({
          file: f,
          action: 'create',
          detail: 'library 中不存在此文件，按 writing guide 创建'
        })),
        existingFiles: existingFiles.map(f => ({
          file: f,
          action: 'merge_or_replace',
          detail: 'library 中已存在此文件，需读取对比后决定合并或替换章节'
        }))
      }
    };
  }

  if (state.hasArchiveNotes) {
    return {
      action: 'replace',
      reason: 'library 文档仅含过时的 Archived proposal notes 摘要段，需用完整设计替换',
      instructions: {
        steps: [
          '读取 library 已有文档，保留 Ownership summary、Reading order 等入口章节',
          '将 Current focus 段替换为提案中的最新设计描述',
          '将 Archived proposal notes 段替换为提案中的完整设计内容（按章节组织）',
          '对于 architecture、lifecycle、constraints 等内容，拆分为独立子文档',
          '更新 frontmatter 的 updated 时间戳',
          '运行 validate-doc.js 校验'
        ]
      }
    };
  }

  return {
    action: 'merge',
    reason: 'library 已有完整设计文档，追加提案中的新增内容',
    instructions: {
      steps: [
        '读取 library 已有文档全文',
        '对比提案内容，识别新增/变更的章节',
        '追加新章节到已有文档（不覆盖已有内容）',
        '对于冲突内容（同一主题的不同描述），以提案为准替换对应章节',
        '更新 frontmatter 的 updated 时间戳',
        '运行 validate-doc.js 校验'
      ]
    }
  };
}

// ============================================================
// Domain Scanning (same logic as archive-context.js)
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

// ============================================================
// Library State Detection
// ============================================================

function detectLibraryState(libFullPath) {
  if (!fs.existsSync(libFullPath)) {
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
    hasArchiveNotes
  };
}

// ============================================================
// Proposal Discovery
// ============================================================

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
