#!/usr/bin/env node
'use strict';

const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');
const { loadMainConfig, loadProjectConfig } = require('./lib/config');

const projectRoot = require('./lib/config').findProjectRoot(process.argv[2] || process.cwd());
const args = process.argv.slice(3);
const isFull = args.includes('--full');
const isMcp = args.includes('--mcp');

const config = loadMainConfig(projectRoot);
const projectRefs = config?.projects || [];
let metaVersion = '0.0.0';
try {
  const yaml = require('./vendor/js-yaml');
  const metaPath = path.join(projectRoot, '.flowforge', 'meta.yaml');
  if (fs.existsSync(metaPath)) {
    const meta = yaml.load(fs.readFileSync(metaPath, 'utf8'));
    metaVersion = meta.version || metaVersion;
  }
} catch (_) {}

if (isMcp) {
  outputMcp(projectRoot, config, metaVersion);
} else {
  outputFull(projectRoot, config, projectRefs, metaVersion);
}

function outputMcp(root, config, version) {
  const activeInfo = getActiveProposalInfo(root, config, 1);
  const proposalLine = activeInfo.length > 0
    ? ` | 活跃: ${activeInfo[0].id} [${activeInfo[0].progress}]`
    : '';
  console.log(`FlowForge v${version} | task: status/ready/claim/done | SKILL: design/implement/feedback/archive${proposalLine} | 禁止直接使用 bd`);
}

function outputFull(root, config, projectRefs, version) {
  console.log(`# FlowForge v${version}`);
  console.log();

  const activeInfo = getActiveProposalInfo(root, config, 3);
  if (activeInfo.length > 0) {
    console.log('## 活跃 Proposal');
    console.log();
    for (const p of activeInfo) {
      console.log(`**${p.id}** ${p.title}`);
      console.log(`  进度: ${p.progress}`);
      console.log(`  命令: \`flowforge task status --proposal ${p.id}\``);
      console.log();
    }
  }

  console.log('## SKILL 路由');
  console.log('- 新需求/分析/设计 → `flowforge-design`');
  console.log('- 执行任务/继续推进 → `flowforge-implement`');
  console.log('- 实施中发现/新认知 → `flowforge-feedback`');
  console.log('- 归档沉淀 → `flowforge-archive`');
  console.log();

  console.log('## 任务操作（全部通过 flowforge task CLI）');
  console.log('```bash');
  console.log('flowforge task status --proposal <CR-id>     # 任务状态');
  console.log('flowforge task ready --proposal <CR-id>      # 就绪任务');
  console.log('flowforge task claim --proposal <CR-id> <id> # 认领');
  console.log('flowforge task done --proposal <CR-id> <id>  # 完成');
  console.log('```');
  console.log();

  console.log('## 禁止');
  console.log('- 禁止直接使用 bd 命令操作任务');
  console.log('- 禁止读写 tasks.snapshot.md');
}

function getActiveProposalInfo(root, config, limit) {
  const results = [];
  if (!config?.projects) return results;

  for (const ref of config.projects) {
    const pc = loadProjectConfig(root, ref);
    if (!pc?.wikiRoot) continue;
    const activeDir = path.join(root, pc.wikiRoot, 'workspace', 'proposals', 'active');
    if (!fs.existsSync(activeDir)) continue;

    const dirs = fs.readdirSync(activeDir, { withFileTypes: true })
      .filter(d => d.isDirectory());

    for (const d of dirs) {
      const metaPath = path.join(activeDir, d.name, 'meta.yaml');
      if (!fs.existsSync(metaPath)) continue;
      try {
        const yaml = require('./vendor/js-yaml');
        const meta = yaml.load(fs.readFileSync(metaPath, 'utf8'));
        if (!meta?.id) continue;
        results.push({
          id: meta.id,
          title: meta.title || '无标题',
          updated_at: meta.updated_at || '',
          dir: path.join(activeDir, d.name)
        });
      } catch (_) {}
    }
  }

  results.sort((a, b) => {
    const au = a.updated_at ? String(a.updated_at) : '';
    const bu = b.updated_at ? String(b.updated_at) : '';
    return bu.localeCompare(au);
  });
  const top = results.slice(0, limit);

  for (const r of top) {
    try {
      const taskOutput = execSync(
        `flowforge task status --proposal ${r.id}`,
        { cwd: root, encoding: 'utf8', timeout: 15000, stdio: 'pipe' }
      );
      const status = JSON.parse(taskOutput);
      const done = status.byStatus?.done || 0;
      const total = status.total || 0;
      r.progress = `${done}/${total} done`;
    } catch (_) {
      r.progress = '?/? done';
    }
  }

  return top;
}
