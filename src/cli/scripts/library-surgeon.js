#!/usr/bin/env node
'use strict';

const fs = require('fs');
const path = require('path');

const projectRoot = process.argv[2] || '.';
const args = process.argv.slice(3);
const operation = args[0];
const rest = args.slice(1);

const libRoot = findLibraryRoot(projectRoot);
if (!libRoot) { console.log(JSON.stringify({ error: 'Library not found' })); process.exit(1); }

const dryRun = rest.includes('--dry-run');
const autoConfirm = rest.includes('--auto-confirm');
const cleanArgs = rest.filter(a => a !== '--dry-run' && a !== '--auto-confirm');

switch (operation) {
  case 'merge': handleMerge(cleanArgs); break;
  case 'deprecate': handleDeprecate(cleanArgs); break;
  case 'upgrade': handleUpgrade(cleanArgs); break;
  default:
    console.log(JSON.stringify({
      usage: 'flowforge library surgeon <operation> [args] [--dry-run]',
      operations: ['merge <src1> <src2> --target <dst>', 'deprecate <path>', 'upgrade <path> --importance must']
    }, null, 2));
}

function handleMerge(srcPaths) {
  const targetIdx = srcPaths.indexOf('--target');
  const target = targetIdx >= 0 ? srcPaths[targetIdx + 1] : null;
  const sources = targetIdx >= 0 ? srcPaths.slice(0, targetIdx) : srcPaths;

  if (sources.length < 2 || !target) {
    console.log(JSON.stringify({ error: 'Usage: merge <src1> <src2> --target <dst>' }));
    process.exit(1);
  }

  const plan = { operation: 'merge', sources: [], target, dryRun, actions: [] };
  for (const src of sources) {
    const full = path.join(libRoot, src);
    if (!fs.existsSync(full)) { console.log(JSON.stringify({ error: `Not found: ${src}` })); process.exit(1); }
    plan.sources.push({ path: src, exists: true });
  }

  const targetFull = path.join(libRoot, target);
  if (fs.existsSync(targetFull)) {
    plan.actions.push({ type: 'merge_into', file: target, note: 'Target exists, content will be appended' });
  } else {
    plan.actions.push({ type: 'create', file: target });
  }

  for (const src of sources) {
    if (src === target) continue;
    plan.actions.push({
      type: 'deprecate',
      file: src,
      action: `mark maturity=deprecated, add related.ref → ${target}`
    });
  }

  if (dryRun) { console.log(JSON.stringify(plan, null, 2)); process.exit(0); }

  let merged = '';
  for (const src of sources) {
    const content = fs.readFileSync(path.join(libRoot, src), 'utf8');
    merged += content.replace(/^---\n[\s\S]*?\n---\n?/, '') + '\n\n---\n\n';
  }

  const firstSrc = sources[0];
  const firstFm = extractFrontmatter(fs.readFileSync(path.join(libRoot, firstSrc), 'utf8'));
  if (firstFm) {
    firstFm.title = firstFm.title ? `${firstFm.title} (merged)` : 'Merged Document';
    firstFm.updated = new Date().toISOString();
  }

  const targetDir = path.dirname(targetFull);
  if (!fs.existsSync(targetDir)) fs.mkdirSync(targetDir, { recursive: true });
  fs.writeFileSync(targetFull, formatFrontmatter(firstFm) + merged);

  for (const src of sources) {
    if (src === target) continue;
    const srcFull = path.join(libRoot, src);
    const content = fs.readFileSync(srcFull, 'utf8');
    const m = content.match(/^---\n([\s\S]*?)\n---/);
    if (m) {
      const updated = m[1]
        .replace(/^maturity:.*/m, `maturity: deprecated`)
        .replace(/^status:.*/m, `status: superseded`)
        .replace(/^updated:.*/m, `updated: ${new Date().toISOString()}`);
      fs.writeFileSync(srcFull, `---\n${updated}\n---` + content.slice(m[0].length));
    }
  }

  console.log(JSON.stringify({ merged: true, target, sources, dryRun: false }));
}

function handleDeprecate(args) {
  const target = args[0];
  if (!target) { console.log(JSON.stringify({ error: 'Usage: deprecate <path>' })); process.exit(1); }

  const full = path.join(libRoot, target);
  if (!fs.existsSync(full)) { console.log(JSON.stringify({ error: `Not found: ${target}` })); process.exit(1); }

  if (dryRun) {
    console.log(JSON.stringify({ operation: 'deprecate', path: target, action: 'set maturity=deprecated, status=superseded', dryRun }));
    process.exit(0);
  }

  const content = fs.readFileSync(full, 'utf8');
  const m = content.match(/^---\n([\s\S]*?)\n---/);
  if (m) {
    const updated = m[1]
      .replace(/^maturity:.*/m, `maturity: deprecated`)
      .replace(/^status:.*/m, `status: superseded`)
      .replace(/^updated:.*/m, `updated: ${new Date().toISOString()}`);
    const hasMaturity = /^  maturity:/m.test(updated);
    const finalUpdated = hasMaturity ? updated : updated.replace(/^(  type:.*)$/m, `$1\n  maturity: deprecated`);
    fs.writeFileSync(full, `---\n${finalUpdated}\n---` + content.slice(m[0].length));
    console.log(JSON.stringify({ deprecated: true, path: target }));
  }
}

function handleUpgrade(args) {
  const target = args[0];
  const impIdx = args.indexOf('--importance');
  const importance = impIdx >= 0 ? args[impIdx + 1] : null;

  if (!target) { console.log(JSON.stringify({ error: 'Usage: upgrade <path> [--importance must|should]' })); process.exit(1); }

  const full = path.join(libRoot, target);
  if (!fs.existsSync(full)) { console.log(JSON.stringify({ error: `Not found: ${target}` })); process.exit(1); }

  if (dryRun) {
    console.log(JSON.stringify({ operation: 'upgrade', path: target, importance: importance || '(unchanged)', dryRun }));
    process.exit(0);
  }

  const content = fs.readFileSync(full, 'utf8');
  const m = content.match(/^---\n([\s\S]*?)\n---/);
  if (m) {
    let updated = m[1].replace(/^updated:.*/m, `updated: ${new Date().toISOString()}`);
    if (importance) {
      if (/^  importance:/m.test(updated)) {
        updated = updated.replace(/^  importance:.*/m, `  importance: ${importance}`);
      } else if (/^  type:/m.test(updated)) {
        updated = updated.replace(/^(  type:.*)$/m, `$1\n  importance: ${importance}`);
      }
    }
    fs.writeFileSync(full, `---\n${updated}\n---` + content.slice(m[0].length));
    console.log(JSON.stringify({ upgraded: true, path: target, importance: importance || '(unchanged)' }));
  }
}

function extractFrontmatter(text) {
  const m = text.match(/^---\n([\s\S]*?)\n---/);
  if (!m) return null;
  const result = {};
  let currentKey = null;
  for (const line of m[1].split('\n')) {
    const nested = line.match(/^  (\w+)\s*:\s*(.*)/);
    if (nested && currentKey === 'domain') {
      if (!result.domain) result.domain = {};
      result.domain[nested[1]] = nested[2].trim().replace(/^["']|["']$/g, '');
      continue;
    }
    const kv = line.match(/^(\w+)\s*:\s*(.*)/);
    if (kv) { currentKey = kv[1]; result[kv[1]] = kv[2].trim().replace(/^["']|["']$/g, ''); }
    else { currentKey = null; }
  }
  return result;
}

function formatFrontmatter(fm) {
  if (!fm) return '---\n---\n';
  let yaml = '---\n';
  const topKeys = ['doc_type', 'title', 'status', 'created', 'updated'];
  for (const k of topKeys) { if (fm[k]) yaml += `${k}: ${fm[k]}\n`; }
  if (fm.domain) {
    yaml += 'domain:\n';
    if (fm.domain.scope) yaml += `  scope: ${fm.domain.scope}\n`;
    if (fm.domain.module) yaml += `  module: ${fm.domain.module}\n`;
    if (fm.domain.type) yaml += `  type: ${fm.domain.type}\n`;
    if (fm.domain.importance) yaml += `  importance: ${fm.domain.importance}\n`;
    if (fm.domain.maturity) yaml += `  maturity: ${fm.domain.maturity}\n`;
  }
  yaml += '---\n\n';
  return yaml;
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
