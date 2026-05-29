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

const exploreSection = extractSection(configContent, 'exploration');
const designSection = extractSection(configContent, 'design');
const intakeSection = extractSection(configContent, 'intake');

const workspaceRoot = path.join(projectRoot, wikiRoot, 'workspace');

const intakeDir = path.join(workspaceRoot, 'intake');
const intakeFiles = fs.existsSync(intakeDir)
  ? fs.readdirSync(intakeDir).filter(f => f.endsWith('.md'))
  : [];

const proposalDir = proposalId
  ? path.join(workspaceRoot, 'proposals', proposalId)
  : findActiveProposal(projectRoot, workspaceRoot);

console.log('# Design Context\n');

if (intakeFiles.length > 0) {
  console.log('## Intake Material\n');
  console.log(`目录: ${path.relative(projectRoot, intakeDir)}`);
  console.log(`文件: ${intakeFiles.join(', ')}`);
  if (intakeSection) {
    console.log('\n分析步骤（来自 config）：\n');
    console.log(intakeSection);
  }
}

if (exploreSection) {
  console.log('\n## Exploration Strategy\n');
  console.log(exploreSection);
}

if (designSection) {
  console.log('\n## Design Rules\n');
  console.log(designSection);
}

if (proposalDir) {
  console.log('\n## Current Proposal\n');
  console.log(`路径: ${path.relative(projectRoot, proposalDir)}`);
  const meta = readProposalMeta(proposalDir);
  if (meta) {
    console.log(`状态: ${meta.status}`);
    console.log(`标题: ${meta.title}`);
  }
  const taskMap = path.join(proposalDir, 'task-map.md');
  if (fs.existsSync(taskMap)) {
    console.log('task-map: 已有');
  }
}

function findActiveProposal(root, workspaceRoot) {
  const proposalsDir = path.join(workspaceRoot, 'proposals');
  if (!fs.existsSync(proposalsDir)) return null;
  const dirs = fs.readdirSync(proposalsDir, { withFileTypes: true })
    .filter(d => d.isDirectory());
  for (const d of dirs) {
    const meta = readProposalMeta(path.join(proposalsDir, d.name));
    if (meta && (meta.status === 'draft' || meta.status === 'active')) {
      return path.join(proposalsDir, d.name);
    }
  }
  return null;
}

function readProposalMeta(dir) {
  const metaPath = path.join(dir, 'meta.yaml');
  if (!fs.existsSync(metaPath)) return null;
  const content = fs.readFileSync(metaPath, 'utf8');
  return {
    id: readValue(content, 'id'),
    title: readValue(content, 'title'),
    status: readValue(content, 'status')
  };
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

  return result
    .map(l => l.replace(/^(\s*)#.*$/, '$1'))
    .filter(l => l.trim())
    .join('\n');
}
