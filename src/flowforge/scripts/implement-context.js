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

const implementSection = extractSection(configContent, 'implement');
const taskRulesSection = extractSection(configContent, 'task_rules');

const workspaceRoot = path.join(projectRoot, wikiRoot, 'workspace');

const proposalDir = proposalId
  ? findProposalById(workspaceRoot, proposalId)
  : findActiveProposal(projectRoot, workspaceRoot);

const taskMap = proposalDir ? path.join(proposalDir, 'task-map.md') : null;
const notes = proposalDir ? path.join(proposalDir, 'notes.md') : null;

console.log('# Implement Context\n');

if (implementSection) {
  console.log('## Implement Rules\n');
  console.log(implementSection);
}

if (taskRulesSection) {
  console.log('\n## Task Rules\n');
  console.log(taskRulesSection);
}

if (proposalDir) {
  console.log('\n## Current Proposal\n');
  console.log(`路径: ${path.relative(projectRoot, proposalDir)}`);
  const meta = readProposalMeta(proposalDir);
  if (meta) {
    console.log(`状态: ${meta.status}`);
    console.log(`标题: ${meta.title}`);
  }
  if (taskMap && fs.existsSync(taskMap)) {
    console.log('\n### task-map.md\n');
    console.log(fs.readFileSync(taskMap, 'utf8'));
  } else {
    console.log('\ntask-map.md: 不存在');
  }
  if (notes && fs.existsSync(notes)) {
    console.log('\n### notes.md\n');
    console.log(fs.readFileSync(notes, 'utf8'));
  }
}

function findProposalById(workspaceRoot, proposalId) {
  for (const subdir of ['active', 'completed']) {
    const candidate = path.join(workspaceRoot, 'proposals', subdir, proposalId);
    if (fs.existsSync(candidate)) return candidate;
  }
  return path.join(workspaceRoot, 'proposals', 'active', proposalId);
}

function findActiveProposal(root, workspaceRoot) {
  const activeDir = path.join(workspaceRoot, 'proposals', 'active');
  if (!fs.existsSync(activeDir)) return null;
  const dirs = fs.readdirSync(activeDir, { withFileTypes: true }).filter(d => d.isDirectory());
  for (const d of dirs) {
    const meta = readProposalMeta(path.join(activeDir, d.name));
    if (meta && (meta.status === 'active')) return path.join(activeDir, d.name);
  }
  return null;
}

function readProposalMeta(dir) {
  const metaPath = path.join(dir, 'meta.yaml');
  if (!fs.existsSync(metaPath)) return null;
  const content = fs.readFileSync(metaPath, 'utf8');
  return { id: readValue(content, 'id'), title: readValue(content, 'title'), status: readValue(content, 'status') };
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
