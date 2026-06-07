'use strict';

const fs = require('fs');
const path = require('path');
const yaml = require('../vendor/js-yaml');

const CONFIG_DIR = '.flowforge';

function readYamlFile(filePath) {
  if (!fs.existsSync(filePath)) return null;
  return yaml.load(fs.readFileSync(filePath, 'utf8'));
}

function loadMainConfig(projectRoot) {
  const configPath = path.join(projectRoot, CONFIG_DIR, 'config.yaml');
  if (!fs.existsSync(configPath)) return null;
  return readYamlFile(configPath);
}

function loadProjectConfig(projectRoot, projectRef) {
  const configPath = path.join(projectRoot, CONFIG_DIR, projectRef.config);
  if (!fs.existsSync(configPath)) return null;
  const data = readYamlFile(configPath);
  if (!data) return null;
  return { id: projectRef.id, name: projectRef.name, ...data };
}

function getProjects(projectRoot) {
  const config = loadMainConfig(projectRoot);
  return (config && config.projects) ? config.projects : [];
}

function getProject(projectRoot, projectId) {
  const config = loadMainConfig(projectRoot);
  if (!config) return null;
  const ref = config.projects.find(p => p.id === projectId);
  if (!ref) return null;
  return loadProjectConfig(projectRoot, ref);
}

function loadMeta(proposalDir) {
  const metaPath = path.join(proposalDir, 'meta.yaml');
  return readYamlFile(metaPath);
}

function findProjectRoot(startPath) {
  let dir = fs.lstatSync(startPath).isDirectory() ? startPath : path.dirname(startPath);
  while (true) {
    if (fs.existsSync(path.join(dir, CONFIG_DIR, 'config.yaml'))) return dir;
    const parent = path.dirname(dir);
    if (parent === dir) return null;
    dir = parent;
  }
}

function findProposalDir(projectRoot, config, proposalId) {
  if (!config || !config.projects) return null;
  for (const ref of config.projects) {
    const pc = loadProjectConfig(projectRoot, ref);
    if (!pc) continue;
    for (const sub of ['active', 'completed']) {
      const subDir = path.join(projectRoot, pc.wikiRoot, 'workspace', 'proposals', sub);
      if (!fs.existsSync(subDir)) continue;
      const dirs = fs.readdirSync(subDir, { withFileTypes: true }).filter(d => d.isDirectory());
      for (const d of dirs) {
        if (d.name === proposalId || d.name.startsWith(proposalId + '-')) {
          return path.join(subDir, d.name);
        }
      }
    }
  }
  return null;
}

/**
 * 检查指定 proposal ID 是否已被其他 proposal 占用。
 * 扫描所有 project 的 active/ 和 completed/ 目录，匹配 d.name === id || d.name.startsWith(id + '-')。
 *
 * @param {string} projectRoot
 * @param {object} config — loadMainConfig 的返回值
 * @param {string} proposalId — CR-id，如 CR26060701
 * @returns {{ id: string, exists: boolean, conflicts: Array<{dir: string, status: string, project: string}> }}
 */
function checkProposalId(projectRoot, config, proposalId) {
  const conflicts = [];
  if (!config || !config.projects) return { id: proposalId, exists: false, conflicts };

  for (const ref of config.projects) {
    const pc = loadProjectConfig(projectRoot, ref);
    if (!pc) continue;
    for (const sub of ['active', 'completed']) {
      const subDir = path.join(projectRoot, pc.wikiRoot, 'workspace', 'proposals', sub);
      if (!fs.existsSync(subDir)) continue;
      const dirs = fs.readdirSync(subDir, { withFileTypes: true }).filter(d => d.isDirectory());
      for (const d of dirs) {
        if (d.name === proposalId || d.name.startsWith(proposalId + '-')) {
          conflicts.push({
            dir: path.relative(projectRoot, path.join(subDir, d.name)),
            status: sub,
            project: ref.id,
          });
        }
      }
    }
  }

  return { id: proposalId, exists: conflicts.length > 0, conflicts };
}

/**
 * 根据当日已存在的 proposal，计算下一个可用的 NN 序号。
 * 扫描所有 project 的 active/ 和 completed/ 目录中当日前缀的目录，取最大 NN + 1。
 *
 * @param {string} projectRoot
 * @param {object} config — loadMainConfig 的返回值
 * @param {string} [prefix='CR'] — ID 前缀
 * @returns {string} — 完整的建议 ID，如 CR26060702
 */
function suggestProposalId(projectRoot, config, prefix = 'CR') {
  const now = new Date();
  const yymmdd =
    String(now.getFullYear()).slice(2) +
    String(now.getMonth() + 1).padStart(2, '0') +
    String(now.getDate()).padStart(2, '0');
  const dayPrefix = prefix + yymmdd;

  const existingNNs = [];
  if (!config || !config.projects) return dayPrefix + '01';

  for (const ref of config.projects) {
    const pc = loadProjectConfig(projectRoot, ref);
    if (!pc) continue;
    for (const sub of ['active', 'completed']) {
      const subDir = path.join(projectRoot, pc.wikiRoot, 'workspace', 'proposals', sub);
      if (!fs.existsSync(subDir)) continue;
      for (const d of fs.readdirSync(subDir, { withFileTypes: true })) {
        if (!d.isDirectory()) continue;
        // 匹配 CR{YYMMDD}{NN}[-...] 或 CR{YYMMDD}{NN}
        const escapedPrefix = prefix.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
        const m = d.name.match(new RegExp('^' + escapedPrefix + '(\\d{6})(\\d{2})(-|$)'));
        if (m && m[1] === yymmdd) existingNNs.push(parseInt(m[2], 10));
      }
    }
  }

  const nextNN = existingNNs.length === 0 ? 1 : Math.max(...existingNNs) + 1;
  if (nextNN > 99) {
    throw new Error(`当日 proposal 数量已达上限 (99)，无法建议新 ID`);
  }
  return dayPrefix + String(nextNN).padStart(2, '0');
}

module.exports = {
  readYamlFile, loadMainConfig, loadProjectConfig, getProjects, getProject,
  loadMeta, findProjectRoot, findProposalDir,
  checkProposalId, suggestProposalId,
};
