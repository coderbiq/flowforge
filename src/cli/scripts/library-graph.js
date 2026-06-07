#!/usr/bin/env node
'use strict';

const fs = require('fs');
const path = require('path');

const projectRoot = process.argv[2] || '.';
const args = process.argv.slice(3);
const subcommand = args[0];
const rest = args.slice(1);

const libRoot = findLibraryRoot(projectRoot);
if (!fs.existsSync(libRoot)) {
  console.log(JSON.stringify({ error: 'Library not found' }));
  process.exit(0);
}

function buildGraph() {
  const nodes = {};
  const adjIn = {};
  const adjOut = {};

  (function walk(dir) {
    for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
      if (entry.isDirectory()) { walk(path.join(dir, entry.name)); continue; }
      if (!entry.name.endsWith('.md')) continue;
      const full = path.join(dir, entry.name);
      const rel = path.relative(libRoot, full);
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
        if (currentKey === 'related' && line.match(/^\s*-\s+ref:\s*(.*)/)) {
          const refMatch = line.match(/^\s*-\s+ref:\s*(.*)/);
          if (!fm.related) fm.related = [];
          fm.related.push(refMatch[1].trim().replace(/^["']|["']$/g, ''));
          continue;
        }
        const kv = line.match(/^(\w+)\s*:\s*(.*)/);
        if (kv) { currentKey = kv[1]; fm[kv[1]] = kv[2].trim().replace(/^["']|["']$/g, ''); }
        else { currentKey = null; }
      }

      nodes[rel] = {
        title: fm.title || rel,
        doc_type: fm.doc_type || '',
        importance: fm.domain?.importance || 'should',
        maturity: fm.domain?.maturity || 'growing',
        topics: fm.topics || [],
      };
      adjIn[rel] = [];
      adjOut[rel] = [];
    }
  })(libRoot);

  for (const [rel, node] of Object.entries(nodes)) {
    const full = path.join(libRoot, rel);
    const content = fs.readFileSync(full, 'utf8');
    const m = content.match(/^---\n([\s\S]*?)\n---/);
    if (!m) continue;
    const refs = [];
    let inRelated = false;
    for (const line of m[1].split('\n')) {
      if (line.match(/^\s*related\s*:/)) { inRelated = true; continue; }
      if (inRelated && line.match(/^\s*-\s+ref:\s*(.*)/)) {
        refs.push(line.match(/^\s*-\s+ref:\s*(.*)/)[1].trim().replace(/^["']|["']$/g, ''));
        continue;
      }
      if (inRelated && !line.match(/^\s/)) inRelated = false;
    }
    for (const ref of refs) {
      const target = path.normalize(path.join(path.dirname(rel), ref)).replace(/\\/g, '/');
      if (nodes[target]) {
        adjOut[rel].push({ target, type: 'explicit-ref' });
        adjIn[target].push({ source: rel, type: 'explicit-ref' });
      }
    }
  }

  return { nodes, adjIn, adjOut };
}

const graph = buildGraph();

switch (subcommand) {
  case 'backlinks': {
    const target = rest[0];
    console.log(JSON.stringify({
      path: target,
      backlinks: (graph.adjIn[target] || []).map(e => ({
        source: e.source,
        title: graph.nodes[e.source]?.title || e.source,
      })),
    }, null, 2));
    break;
  }
  case 'refs': {
    const source = rest[0];
    console.log(JSON.stringify({
      path: source,
      refs: (graph.adjOut[source] || []).map(e => ({
        target: e.target,
        title: graph.nodes[e.target]?.title || e.target,
      })),
    }, null, 2));
    break;
  }
  case 'orphans': {
    const orphans = Object.keys(graph.nodes).filter(n => (graph.adjIn[n] || []).length === 0);
    console.log(JSON.stringify({ orphans: orphans.map(n => ({ path: n, title: graph.nodes[n].title })) }, null, 2));
    break;
  }
  case 'hubs': {
    const top = parseInt(rest[0]) || 10;
    const ranked = Object.entries(graph.adjIn)
      .sort((a, b) => b[1].length - a[1].length)
      .slice(0, top)
      .map(([n, edges]) => ({ path: n, title: graph.nodes[n]?.title, incoming: edges.length }));
    console.log(JSON.stringify({ hubs: ranked }, null, 2));
    break;
  }
  case 'blast-radius': {
    const source = rest[0];
    const depth = parseInt(rest[1]) || 2;
    const visited = new Set([source]);
    const queue = [[source, 0]];
    const result = { source, depth, affected: [] };
    while (queue.length > 0) {
      const [current, d] = queue.shift();
      if (d >= depth) continue;
      for (const edge of (graph.adjIn[current] || [])) {
        if (!visited.has(edge.source)) {
          visited.add(edge.source);
          queue.push([edge.source, d + 1]);
          result.affected.push({ path: edge.source, depth: d + 1, title: graph.nodes[edge.source]?.title });
        }
      }
    }
    console.log(JSON.stringify(result, null, 2));
    break;
  }
  default:
    console.log(JSON.stringify({
      usage: 'flowforge library graph <subcommand> [args]',
      subcommands: ['backlinks <path>', 'refs <path>', 'orphans', 'hubs [top]', 'blast-radius <path> [depth]'],
      stats: { nodes: Object.keys(graph.nodes).length, edges: Object.values(graph.adjOut).reduce((s, e) => s + e.length, 0) },
    }, null, 2));
}

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
