#!/usr/bin/env node
'use strict';

const fs = require('fs');
const path = require('path');

const projectRoot = process.argv[2] || process.cwd();

// ── 加载 config.yaml 获取 wiki root ──
const configPath = path.join(projectRoot, '.flowforge', 'config.yaml');
if (!fs.existsSync(configPath)) {
  console.error('ERROR: .flowforge/config.yaml 不存在');
  process.exit(0);
}
const configContent = fs.readFileSync(configPath, 'utf8');
const wikiRoot = readRootValue(configContent, 'root') || 'ff-wiki';

const proposalsDir = path.join(projectRoot, wikiRoot, 'workspace', 'proposals');
if (!fs.existsSync(proposalsDir)) {
  console.error(`ERROR: proposals 目录不存在: ${proposalsDir}`);
  process.exit(0);
}

// ── 扫描所有子目录中的 proposal ──
const activeStatuses = ['draft', 'active', 'implemented'];
const completedStatuses = ['archived', 'rejected'];

const activeProposals = [];
const completedProposals = [];

for (const subdir of ['active', 'completed']) {
  const dir = path.join(proposalsDir, subdir);
  if (!fs.existsSync(dir)) continue;
  const entries = fs.readdirSync(dir, { withFileTypes: true }).filter(d => d.isDirectory());
  for (const entry of entries) {
    const meta = readProposalMeta(path.join(dir, entry.name));
    if (!meta) continue;
    if (activeStatuses.includes(meta.status)) {
      activeProposals.push({ ...meta, relPath: `${subdir}/${entry.name}` });
    } else if (completedStatuses.includes(meta.status)) {
      completedProposals.push({ ...meta, relPath: `${subdir}/${entry.name}` });
    }
  }
}

// ── 按 updated_at 降序排列 ──
const sortByUpdated = (a, b) => (b.updated_at || '').localeCompare(a.updated_at || '');
activeProposals.sort(sortByUpdated);
completedProposals.sort(sortByUpdated);

const total = activeProposals.length + completedProposals.length;
const now = new Date();
const timestamp = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')} ${String(now.getHours()).padStart(2, '0')}:${String(now.getMinutes()).padStart(2, '0')}`;

// ── 生成 INDEX.md ──
let md = '';
md += '# 📋 Proposals Index\n\n';
md += `> 自动生成于 ${timestamp} · 共 ${total} 提案（🟢 ${activeProposals.length} 进行中 · 📦 ${completedProposals.length} 已完成）\n\n`;

md += '## 🔄 进行中\n\n';
if (activeProposals.length === 0) {
  md += '（无进行中提案）\n\n';
} else {
  md += '| ID | 标题 | 规模 | 状态 | 模块 | 更新 | 最新进度 |\n';
  md += '|----|------|------|------|------|------|----------|\n';
  for (const p of activeProposals) {
    const idLink = `[${p.id}](./${p.relPath}/)`;
    const size = p.size_class || '—';
    const status = `\`${p.status}\``;
    const modules = p.ownership || '—';
    const updated = formatDate(p.updated_at);
    const progress = p.latest_progress || '—';
    md += `| ${idLink} | ${p.title} | ${size} | ${status} | ${modules} | ${updated} | ${progress} |\n`;
  }
  md += '\n';
}

md += '## 📦 已完成\n\n';
if (completedProposals.length === 0) {
  md += '（无已完成提案）\n\n';
} else {
  md += '| ID | 标题 | 规模 | 状态 | 模块 | 更新 | 最新进度 |\n';
  md += '|----|------|------|------|------|------|----------|\n';
  for (const p of completedProposals) {
    const idLink = `[${p.id}](./${p.relPath}/)`;
    const size = p.size_class || '—';
    const status = `\`${p.status}\``;
    const modules = p.ownership || '—';
    const updated = formatDate(p.updated_at);
    const progress = p.latest_progress || '—';
    md += `| ${idLink} | ${p.title} | ${size} | ${status} | ${modules} | ${updated} | ${progress} |\n`;
  }
  md += '\n';
}

// ── 写入 INDEX.md ──
const indexPath = path.join(proposalsDir, 'INDEX.md');
fs.writeFileSync(indexPath, md, 'utf8');
console.log(`INDEX.md 已生成: ${indexPath}`);
console.log(`  进行中: ${activeProposals.length}, 已完成: ${completedProposals.length}`);

// ── Helpers ──

function readProposalMeta(dir) {
  const metaPath = path.join(dir, 'meta.yaml');
  if (!fs.existsSync(metaPath)) return null;
  const content = fs.readFileSync(metaPath, 'utf8');
  const id = readYamlValue(content, 'id');
  const title = readYamlValue(content, 'title');
  const status = readYamlValue(content, 'status');
  const updated_at = readYamlValue(content, 'updated_at');
  const size_class = readYamlValue(content, 'size_class');
  const latest_progress = readYamlValue(content, 'latest_progress');
  const ownership = extractPrimaryOwnership(content);
  return {
    id: id || path.basename(dir),
    title: title || '（无标题）',
    status: status || 'unknown',
    updated_at: updated_at || '',
    size_class: size_class || '',
    latest_progress: latest_progress || '',
    ownership: ownership || ''
  };
}

function extractPrimaryOwnership(content) {
  const lines = content.split('\n');
  let inOwnership = false;
  let ownershipBaseIndent = -1;
  let currentItem = null;
  const primaries = [];

  for (const line of lines) {
    const indent = line.search(/\S/);
    if (indent < 0) continue;
    const trimmed = line.trim();

    if (/^ownership\s*:/.test(trimmed)) {
      inOwnership = true;
      ownershipBaseIndent = indent;
      continue;
    }

    if (inOwnership) {
      if (indent <= ownershipBaseIndent) {
        if (currentItem && currentItem.role === 'primary' && currentItem.target) {
          primaries.push(currentItem.target.replace(/^.*\//, ''));
        }
        break;
      }

      if (/^-\s/.test(trimmed)) {
        if (currentItem && currentItem.role === 'primary' && currentItem.target) {
          primaries.push(currentItem.target.replace(/^.*\//, ''));
        }
        currentItem = { target: '', role: '', type: '' };
        const inline = trimmed.replace(/^-\s*/, '');
        const inlineMatch = inline.match(/^(target|role|type)\s*:\s*(.*)/);
        if (inlineMatch) {
          currentItem[inlineMatch[1]] = inlineMatch[2].replace(/^["']|["']$/g, '').trim();
        }
      } else if (currentItem) {
        const m = trimmed.match(/^(target|role|type)\s*:\s*(.*)/);
        if (m) {
          currentItem[m[1]] = m[2].replace(/^["']|["']$/g, '').trim();
        }
      }
    }
  }

  if (inOwnership && currentItem && currentItem.role === 'primary' && currentItem.target) {
    primaries.push(currentItem.target.replace(/^.*\//, ''));
  }

  return primaries.length > 0 ? primaries.map(p => `\`${p}\``).join(' ') : '';
}

function formatDate(isoStr) {
  if (!isoStr) return '—';
  const m = isoStr.match(/(\d{4})-(\d{2})-(\d{2})/);
  if (m) return `${m[2]}-${m[3]}`;
  return isoStr.substring(0, 10);
}

function readRootValue(content, key) {
  const m = content.match(new RegExp(`^\\s*${key}\\s*:\\s*["']?([^"'\n#]+)["']?`, 'm'));
  return m ? m[1].trim() : null;
}

function readYamlValue(content, key) {
  const m = content.match(new RegExp(`^\\s*${key}\\s*:\\s*["']?([^"'\n#]+)`, 'm'));
  return m ? m[1].trim() : null;
}
