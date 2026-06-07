#!/usr/bin/env node
'use strict';

const fs = require('fs');
const path = require('path');

const projectRoot = process.argv[2] || '.';
const args = process.argv.slice(3);

const flags = {
  staleness: args.includes('--staleness') || args.includes('--all'),
  brokenRefs: args.includes('--broken-refs') || args.includes('--all'),
  duplicates: args.includes('--duplicates') || args.includes('--all'),
  orphans: args.includes('--orphans') || args.includes('--all'),
  validateAll: args.includes('--validate-all') || args.includes('--all'),
  quality: args.includes('--quality') || args.includes('--all'),
  reviewList: args.includes('--review-list'),
};

const reviewTarget = args.includes('--review') ? args[args.indexOf('--review') + 1] : null;

if (!flags.staleness && !flags.brokenRefs && !flags.duplicates && !flags.orphans && !flags.validateAll && !flags.quality && !flags.reviewList && !reviewTarget) {
  flags.staleness = true;
}

const libRoot = findLibraryRoot(projectRoot);
if (!libRoot) {
  console.log(JSON.stringify({ error: 'Library not found' }));
  process.exit(0);
}

if (reviewTarget) {
  const docPath = path.join(libRoot, reviewTarget);
  if (!fs.existsSync(docPath)) {
    console.log(JSON.stringify({ error: 'Document not found', path: reviewTarget }));
    process.exit(0);
  }
  const content = fs.readFileSync(docPath, 'utf8');
  const m = content.match(/^---\n([\s\S]*?)\n---/);
  if (m) {
    const now = new Date().toISOString();
    const hasReview = /^last_reviewed:/m.test(m[1]);
    const updated = hasReview
      ? m[1].replace(/^last_reviewed:.*/m, `last_reviewed: ${now}`)
      : m[1] + `\nlast_reviewed: ${now}`;
    fs.writeFileSync(docPath, `---\n${updated}\n---` + content.slice(m[0].length));
    console.log(JSON.stringify({ reviewed: true, path: reviewTarget, reviewedAt: now }));
  }
  process.exit(0);
}

const results = { staleness: [], brokenRefs: [], duplicates: [], orphans: [], validateAll: [], quality: [], reviewList: [] };

const docs = [];
(function walk(dir) {
  for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
    if (entry.isDirectory()) { walk(path.join(dir, entry.name)); continue; }
    if (!entry.name.endsWith('.md')) continue;
    docs.push(path.join(dir, entry.name));
  }
})(libRoot);

for (const docPath of docs) {
  const content = fs.readFileSync(docPath, 'utf8');
  const fm = extractFrontmatter(content);
  if (!fm) continue;

  const relPath = path.relative(libRoot, docPath);
  const body = content.replace(/^---\n[\s\S]*?\n---\n?/, '');
  const wordCount = body.replace(/[#*>\[\]`|=-]/g, ' ').split(/\s+/).filter(Boolean).length;
  const sections = (body.match(/^## /gm) || []).length;
  const codeBlocks = (body.match(/```/g) || []).length / 2;
  const refCount = (fm.related || []).length;

  if (flags.staleness) {
    const lastReviewed = fm.last_reviewed ? new Date(fm.last_reviewed) : null;
    const updated = fm.updated ? new Date(fm.updated) : null;
    const interval = Number(fm.review_interval) || 180;
    const refDate = lastReviewed || updated;

    if (refDate) {
      const daysSince = Math.floor((Date.now() - refDate.getTime()) / 86400000);
      if (daysSince > interval) {
        results.staleness.push({ path: relPath, title: fm.title, daysSince, interval, reviewed: !!lastReviewed });
      }
    } else {
      results.staleness.push({ path: relPath, title: fm.title, daysSince: null, interval, reviewed: false, note: '从未审查' });
    }
  }

  if (flags.brokenRefs) {
    for (const ref of (fm.related || [])) {
      const refPath = path.resolve(path.dirname(docPath), ref.ref || '');
      if (!fs.existsSync(refPath)) {
        results.brokenRefs.push({ path: relPath, brokenRef: ref.ref });
      }
    }
  }

  if (flags.validateAll) {
    const errors = [];
    if (!fm.doc_type || !fm.title || !fm.status) errors.push('missing required fields');
    if (fm.domain) {
      if (fm.domain.importance && !['must','should','may','info'].includes(fm.domain.importance))
        errors.push(`invalid importance: ${fm.domain.importance}`);
      if (fm.domain.maturity && !['seed','growing','stable','deprecated'].includes(fm.domain.maturity))
        errors.push(`invalid maturity: ${fm.domain.maturity}`);
    }
    if (errors.length) results.validateAll.push({ path: relPath, errors });
  }

  if (flags.quality) {
    const score = Math.min(100, Math.round(
      (wordCount > 200 ? 25 : wordCount / 8) +
      (sections > 3 ? 30 : sections * 10) +
      (codeBlocks > 1 ? 15 : codeBlocks * 7) +
      (refCount > 2 ? 20 : refCount * 10) +
      (fm.topics?.length > 2 ? 10 : (fm.topics?.length || 0) * 3)
    ));
    const tags = [];
    if (wordCount < 100) tags.push('thin');
    if (sections < 2) tags.push('flat');
    if (refCount === 0) tags.push('isolated');
    if (!fm.last_reviewed) tags.push('unreviewed');
    if (score < 30) tags.push('low-quality');

    results.quality.push({
      path: relPath, title: fm.title, score,
      wordCount, sections, codeBlocks, refCount,
      tags: tags.length ? tags : ['ok']
    });
  }

  if (flags.reviewList) {
    const interval = Number(fm.review_interval) || 180;
    const lastReviewed = fm.last_reviewed ? new Date(fm.last_reviewed) : null;
    const updated = fm.updated ? new Date(fm.updated) : null;
    const refDate = lastReviewed || updated;
    const daysSince = refDate ? Math.floor((Date.now() - refDate.getTime()) / 86400000) : null;
    const daysUntil = daysSince !== null ? interval - daysSince : -1;

    results.reviewList.push({
      path: relPath, title: fm.title,
      lastReviewed: fm.last_reviewed || null,
      interval,
      daysSince,
      daysUntil: Math.max(0, daysUntil),
      overdue: daysSince !== null && daysSince > interval,
      neverReviewed: !fm.last_reviewed,
      importance: fm.domain?.importance || 'should'
    });
  }
}

if (flags.duplicates) {
  const titles = {};
  for (const docPath of docs) {
    const fm = extractFrontmatter(fs.readFileSync(docPath, 'utf8'));
    if (!fm?.title) continue;
    const rel = path.relative(libRoot, docPath);
    if (!titles[fm.title]) titles[fm.title] = [];
    titles[fm.title].push(rel);
  }
  for (const [title, paths] of Object.entries(titles)) {
    if (paths.length > 1) results.duplicates.push({ title, paths });
  }
}

if (flags.orphans) {
  const allRefs = new Set();
  for (const docPath of docs) {
    const fm = extractFrontmatter(fs.readFileSync(docPath, 'utf8'));
    for (const ref of (fm?.related || [])) {
      allRefs.add(ref.ref);
    }
  }
  for (const docPath of docs) {
    const rel = path.relative(libRoot, docPath);
    if (!allRefs.has(rel)) results.orphans.push({ path: rel });
  }
}

if (flags.reviewList) {
  results.reviewList.sort((a, b) => {
    if (a.overdue !== b.overdue) return a.overdue ? -1 : 1;
    return (a.daysUntil || 999) - (b.daysUntil || 999);
  });
}

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
    if (currentKey === 'related' && line.match(/^\s*-\s+ref:\s*(.*)/)) {
      const refMatch = line.match(/^\s*-\s+ref:\s*(.*)/);
      if (!result.related) result.related = [];
      result.related.push({ ref: refMatch[1].trim().replace(/^["']|["']$/g, '') });
      continue;
    }
    const kv = line.match(/^(\w+)\s*:\s*(.*)/);
    if (kv) {
      currentKey = kv[1];
      result[kv[1]] = kv[2].trim().replace(/^["']|["']$/g, '');
    } else {
      currentKey = null;
    }
  }
  return result;
}
