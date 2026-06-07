#!/usr/bin/env node
'use strict';

const fs = require('fs');
const path = require('path');

const projectRoot = process.argv[2] || '.';
const args = process.argv.slice(3);
const shouldRefresh = args.includes('--refresh');

const libRoot = findLibraryRoot(projectRoot);
if (!fs.existsSync(libRoot)) {
  console.log(JSON.stringify({ error: 'Library not found' }));
  process.exit(0);
}

const docs = [];
(function walk(dir) {
  for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
    if (entry.isDirectory()) { walk(path.join(dir, entry.name)); continue; }
    if (!entry.name.endsWith('.md')) continue;
    docs.push(path.join(dir, entry.name));
  }
})(libRoot);

const importanceOrder = { must: 0, should: 1, may: 2, info: 3 };
const icons = { must: '⚠️', should: '📌', may: '💡', info: '📄' };

const stats = { total: 0, byType: {}, byMaturity: {}, byDir: {}, must: [], deprecated: [] };

for (const docPath of docs) {
  const content = fs.readFileSync(docPath, 'utf8');
  const m = content.match(/^---\n([\s\S]*?)\n---/);
  if (!m) continue;
  const fm = {};
  let currentKey = null;
  for (const line of m[1].split('\n')) {
    const nested = line.match(/^  (\w+)\s*:\s*(.*)/);
    if (nested && currentKey === 'domain') {
      if (!fm.domain) fm.domain = {};
      fm.domain[nested[1]] = nested[2].trim().replace(/^["']|["']$/g, '');
      continue;
    }
    const kv = line.match(/^(\w+)\s*:\s*(.*)/);
    if (kv) { currentKey = kv[1]; fm[kv[1]] = kv[2].trim().replace(/^["']|["']$/g, ''); }
    else { currentKey = null; }
  }
  if (!fm.title) continue;

  const rel = path.relative(libRoot, docPath);
  const dir = rel.split('/')[0];
  const docType = fm.doc_type || 'unknown';
  const maturity = fm.domain?.maturity || 'growing';
  const importance = fm.domain?.importance || 'should';

  stats.total++;
  stats.byType[docType] = (stats.byType[docType] || 0) + 1;
  stats.byMaturity[maturity] = (stats.byMaturity[maturity] || 0) + 1;
  stats.byDir[dir] = (stats.byDir[dir] || 0) + 1;

  if (importance === 'must') stats.must.push({ rel, title: fm.title });
  if (maturity === 'deprecated') stats.deprecated.push({ rel, title: fm.title });
}

let md = '# Library Index\n\n';
md += `> 摘要看板 — ${stats.total} 篇文档 | 更新于 ${new Date().toISOString().replace('T', ' ').slice(0, 16)}\n\n`;

md += '## 📊 概况\n\n';
md += '| 目录 | 文档数 |\n';
md += '|------|--------|\n';
const dirOrder = ['architecture', 'conventions', 'decisions', 'modules'];
for (const d of dirOrder) {
  if (stats.byDir[d]) md += `| ${d}/ | ${stats.byDir[d]} |\n`;
}
for (const d of Object.keys(stats.byDir).sort()) {
  if (!dirOrder.includes(d)) md += `| ${d}/ | ${stats.byDir[d]} |\n`;
}
md += '\n';

md += '| 类型 | 文档数 |\n';
md += '|------|--------|\n';
for (const [t, n] of Object.entries(stats.byType).sort((a, b) => b[1] - a[1])) {
  md += `| ${t} | ${n} |\n`;
}
md += '\n';

md += '| 成熟度 | 文档数 |\n';
md += '|--------|--------|\n';
const matOrder = ['stable', 'growing', 'seed', 'deprecated'];
for (const m of matOrder) {
  if (stats.byMaturity[m]) md += `| ${m} | ${stats.byMaturity[m]} |\n`;
}
md += '\n';

if (stats.must.length > 0) {
  md += '## ⚠️ 铁律\n\n';
  for (const d of stats.must) {
    md += `- [${d.title}](${d.rel})\n`;
  }
  md += '\n';
}

if (stats.deprecated.length > 0) {
  md += '## 🗑️ 待清理（已废弃）\n\n';
  for (const d of stats.deprecated) {
    md += `- [${d.title}](${d.rel})\n`;
  }
  md += '\n';
}

md += '## 🔍 查找\n\n';
md += '```bash\n';
md += 'flowforge library search "keyword"    # 全文搜索\n';
md += 'flowforge library list --type design  # 按类型过滤\n';
md += 'flowforge library list --module X     # 按模块过滤\n';
md += 'flowforge library check --staleness   # 过期检测\n';
md += '```\n';

const indexPath = path.join(libRoot, 'INDEX.md');
fs.writeFileSync(indexPath, md);
console.log(JSON.stringify({ refreshed: true, path: path.relative(projectRoot, indexPath), count: stats.total }));

function findLibraryRoot(root) {
  const paths = [];

  const configPath = path.join(root, '.flowforge', 'config.yaml');
  if (fs.existsSync(configPath)) {
    try {
      const yaml = fs.readFileSync(configPath, 'utf8');
      for (const line of yaml.split('\n')) {
        const m = line.match(/^\s*-\s*id:\s*(\S+)/);
        if (!m) continue;
        const projYaml = path.join(root, '.flowforge', 'projects', `${m[1]}.yaml`);
        if (!fs.existsSync(projYaml)) continue;
        const pc = fs.readFileSync(projYaml, 'utf8');
        const wm = pc.match(/wikiRoot:\s*(\S+)/);
        if (wm) paths.push(path.join(root, wm[1], 'library'));
      }
    } catch (_) {}
  }

  paths.push(path.join(root, 'ff-wiki', 'library'));

  for (const p of paths) {
    if (!fs.existsSync(p)) continue;
    const hasContent = fs.readdirSync(p, { recursive: true }).some(f => f.endsWith('.md'));
    if (hasContent) return p;
  }

  return paths.find(p => fs.existsSync(p)) || paths[paths.length - 1];
}
