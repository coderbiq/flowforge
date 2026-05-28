#!/usr/bin/env node
'use strict';

const fs = require('fs');
const path = require('path');

const projectRoot = process.argv[2] || process.cwd();

// ── 加载 config.yaml ──
const configPath = path.join(projectRoot, '.flowforge', 'config.yaml');
if (!fs.existsSync(configPath)) {
  outputError('.flowforge/config.yaml 不存在');
  process.exit(0);
}
const configContent = fs.readFileSync(configPath, 'utf8');

// ── 提取 scenes 配置 ──
const scenes = extractScenes(configContent);

// ── 加载 proposal 状态 ──
const wikiRoot = readWikiRoot(configContent) || 'ff-wiki';
const proposalsDir = path.join(projectRoot, wikiRoot, 'proposals');
let proposalSummaries = [];
if (fs.existsSync(proposalsDir)) {
  proposalSummaries = fs.readdirSync(proposalsDir, { withFileTypes: true })
    .filter(d => d.isDirectory())
    .map(d => readProposalMeta(path.join(proposalsDir, d.name)))
    .filter(Boolean);
}

// ── 输出 ──
console.log('# Workflow Context\n');
console.log('## Scenes (from config.yaml)\n');
console.log(scenes);
console.log('\n## Active Proposals\n');
if (proposalSummaries.length === 0) {
  console.log('（无活跃 proposal）');
} else {
  for (const p of proposalSummaries) {
    console.log(`- **${p.id}** ${p.title} (_${p.status}_) → ${p.dir}`);
  }
}

// ── Helpers ──

function outputError(msg) {
  console.error(`ERROR: ${msg}`);
}

/** 从 config.yaml 中提取 rules.workflow.scenes 片段 */
function extractScenes(content) {
  const lines = content.split('\n');
  let inScenes = false;
  let sceneIndent = -1;
  const result = [];

  for (const line of lines) {
    const indent = line.search(/\S/);

    if (inScenes) {
      if (indent >= 0 && indent <= sceneIndent) {
        break;
      }
      result.push(line);
      continue;
    }

    if (/^\s*rules\s*:/.test(line) ||
        /^\s*workflow\s*:/.test(line) ||
        /^\s*scenes\s*:/.test(line)) {
      if (/^\s*scenes\s*:/.test(line)) {
        inScenes = true;
        sceneIndent = indent;
      }
    }
  }

  // 去掉注释行，保留内容
  return result
    .map(l => l.replace(/^(\s*)#.*$/, '$1'))
    .filter(l => l.trim())
    .join('\n');
}

/** 从 config.yaml 中读取 wiki.root */
function readWikiRoot(content) {
  const m = content.match(/^\s*root\s*:\s*["']?([^"'\s#]+)["']?\s*(?:#.*)?$/m);
  return m ? m[1] : null;
}

/** 读取单个 proposal 的 meta.yaml */
function readProposalMeta(proposalDir) {
  const metaPath = path.join(proposalDir, 'meta.yaml');
  if (!fs.existsSync(metaPath)) return null;

  const content = fs.readFileSync(metaPath, 'utf8');
  const id = readYamlValue(content, 'id');
  const title = readYamlValue(content, 'title');
  const status = readYamlValue(content, 'status');

  return {
    dir: path.relative(projectRoot, proposalDir),
    id: id || path.basename(proposalDir),
    title: title || '（无标题）',
    status: status || 'unknown'
  };
}

/** 简单 YAML 值读取：key: value */
function readYamlValue(content, key) {
  const m = content.match(new RegExp(`^\\s*${key}\\s*:\\s*["']?([^"'\n#]+)`, 'm'));
  return m ? m[1].trim() : null;
}
