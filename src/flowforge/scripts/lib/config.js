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

function findProposalDir(projectRoot, config, proposalId) {
  if (!config || !config.projects) return null;
  for (const ref of config.projects) {
    const pc = loadProjectConfig(projectRoot, ref);
    if (!pc) continue;
    for (const sub of ['active', 'completed']) {
      const cand = path.join(projectRoot, pc.wikiRoot, 'workspace', 'proposals', sub, proposalId);
      if (fs.existsSync(cand)) return cand;
    }
  }
  return null;
}

module.exports = { readYamlFile, loadMainConfig, loadProjectConfig, getProjects, getProject, loadMeta, findProposalDir };
