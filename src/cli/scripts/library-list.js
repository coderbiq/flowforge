#!/usr/bin/env node
'use strict';

const fs = require('fs');
const path = require('path');

const projectRoot = process.argv[2] || '.';
const args = process.argv.slice(3);

const filters = {};
for (let i = 0; i < args.length; i++) {
  if (args[i] === '--type') filters.type = args[++i];
  if (args[i] === '--scope') filters.scope = args[++i];
  if (args[i] === '--importance') filters.importance = args[++i];
  if (args[i] === '--maturity') filters.maturity = args[++i];
  if (args[i] === '--module') filters.module = args[++i];
}

const libRoot = findLibraryRoot(projectRoot);
if (!libRoot) { console.log(JSON.stringify([])); process.exit(0); }

const results = [];
walk(libRoot);
function walk(dir) {
  const entries = fs.readdirSync(dir, { withFileTypes: true });
  for (const entry of entries) {
    if (entry.isDirectory()) { walk(path.join(dir, entry.name)); continue; }
    if (!entry.name.endsWith('.md')) continue;
    const full = path.join(dir, entry.name);
    const content = fs.readFileSync(full, 'utf8');
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

    const rel = path.relative(libRoot, full);
    const doc = {
      path: rel,
      title: fm.title || '',
      doc_type: fm.doc_type || '',
      status: fm.status || '',
      importance: fm.domain?.importance || 'should',
      maturity: fm.domain?.maturity || 'growing',
      topics: fm.topics || [],
    };

    if (filters.type && doc.doc_type !== filters.type) continue;
    if (filters.scope && fm.domain?.scope !== filters.scope) continue;
    if (filters.importance && doc.importance !== filters.importance) continue;
    if (filters.maturity && doc.maturity !== filters.maturity) continue;
    if (filters.module && fm.domain?.module !== filters.module) continue;

    results.push(doc);
  }
}

results.sort((a, b) => a.path.localeCompare(b.path));
console.log(JSON.stringify(results, null, 2));

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
