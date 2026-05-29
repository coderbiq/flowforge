#!/usr/bin/env node
'use strict';

const fs = require('fs');
const path = require('path');

const projectRoot = process.argv[2] || process.cwd();
const proposalId = process.argv[3] || null;

const configPath = path.join(projectRoot, '.flowforge', 'config.yaml');
if (!fs.existsSync(configPath)) {
  console.error('ERROR: .flowforge/config.yaml 不存在');
  process.exit(0);
}
const configContent = fs.readFileSync(configPath, 'utf8');
const wikiRoot = readRootValue(configContent, 'root') || 'ff-wiki';

const librarySection = extractSection(configContent, 'library');
const modulesSection = extractSection(configContent, 'modules');

const workspaceRoot = path.join(projectRoot, wikiRoot, 'workspace');

const proposalDir = proposalId
  ? findProposalById(workspaceRoot, proposalId)
  : findProposal(projectRoot, workspaceRoot, ['implemented', 'archived']);

console.log('# Archive Context\n');

if (librarySection) {
  console.log('## Library Rules\n');
  console.log(librarySection);
}

if (modulesSection) {
  console.log('\n## Module Registry\n');
  console.log(modulesSection);
}

if (proposalDir) {
  console.log('\n## Current Proposal\n');
  console.log(`路径: ${path.relative(projectRoot, proposalDir)}`);
  const meta = readProposalMeta(proposalDir);
  if (meta) {
    console.log(`ID: ${meta.id}`);
    console.log(`状态: ${meta.status}`);
    console.log(`标题: ${meta.title}`);
  }
  const archiveTargets = meta ? extractArchiveTargets(path.join(proposalDir, 'meta.yaml')) : null;
  if (archiveTargets) {
    console.log('\n### Archive Targets\n');
    console.log(archiveTargets);
  }

  for (const f of ['proposal.md', 'design.md', 'task-map.md']) {
    const fp = path.join(proposalDir, f);
    if (fs.existsSync(fp)) {
      console.log(`\n### ${f}\n`);
      console.log(fs.readFileSync(fp, 'utf8'));
    }
  }
}

function findProposalById(workspaceRoot, proposalId) {
  for (const subdir of ['active', 'completed']) {
    const candidate = path.join(workspaceRoot, 'proposals', subdir, proposalId);
    if (fs.existsSync(candidate)) return candidate;
  }
  return path.join(workspaceRoot, 'proposals', 'active', proposalId);
}

function findProposal(root, workspaceRoot, statuses) {
  for (const subdir of ['active', 'completed']) {
    const proposalsSubDir = path.join(workspaceRoot, 'proposals', subdir);
    if (!fs.existsSync(proposalsSubDir)) continue;
    const dirs = fs.readdirSync(proposalsSubDir, { withFileTypes: true }).filter(d => d.isDirectory());
    for (const d of dirs) {
      const meta = readProposalMeta(path.join(proposalsSubDir, d.name));
      if (meta && statuses.includes(meta.status)) return path.join(proposalsSubDir, d.name);
    }
  }
  return null;
}

function readProposalMeta(dir) {
  const metaPath = path.join(dir, 'meta.yaml');
  if (!fs.existsSync(metaPath)) return null;
  const content = fs.readFileSync(metaPath, 'utf8');
  return { id: readValue(content, 'id'), title: readValue(content, 'title'), status: readValue(content, 'status') };
}

function extractArchiveTargets(metaPath) {
  const content = fs.readFileSync(metaPath, 'utf8');
  const section = extractFromContent(content, 'archive_targets');
  return section || null;
}

function readRootValue(content, key) {
  const m = content.match(new RegExp(`^\\s*${key}\\s*:\\s*["']?([^"'\n#]+)["']?`, 'm'));
  return m ? m[1].trim() : null;
}

function readValue(content, key) {
  const m = content.match(new RegExp(`^\\s*${key}\\s*:\\s*["']?([^"'\n#]+)`, 'm'));
  return m ? m[1].trim() : null;
}

function extractSection(content, sectionKey) {
  const lines = content.split('\n');
  let inSection = false;
  let sectionIndent = -1;
  let inBlock = false;
  const result = [];
  for (const line of lines) {
    const indent = line.search(/\S/);
    const trimmed = line.trim();
    if (inSection && !inBlock) {
      if (indent >= 0 && indent <= sectionIndent) break;
      if (trimmed.endsWith('|')) { inBlock = true; result.push(line); continue; }
      result.push(line);
      continue;
    }
    if (inSection && inBlock) {
      if (indent >= 0 && indent < sectionIndent + 2) { inBlock = false; break; }
      result.push(line);
      continue;
    }
    if (trimmed.startsWith(sectionKey + ':')) {
      inSection = true;
      sectionIndent = indent;
      if (trimmed.endsWith('|')) { inBlock = true; }
    }
  }
  return result.map(l => l.replace(/^(\s*)#.*$/, '$1')).filter(l => l.trim()).join('\n');
}

function extractFromContent(content, key) {
  const lines = content.split('\n');
  let inSection = false;
  let sectionIndent = -1;
  const result = [];
  for (const line of lines) {
    const indent = line.search(/\S/);
    const trimmed = line.trim();
    if (inSection) {
      if (indent >= 0 && indent <= sectionIndent) break;
      result.push(line);
      continue;
    }
    if (trimmed.startsWith(key + ':')) {
      inSection = true;
      sectionIndent = indent;
    }
  }
  return result.map(l => l.replace(/^(\s*)#.*$/, '$1')).filter(l => l.trim()).join('\n');
}
