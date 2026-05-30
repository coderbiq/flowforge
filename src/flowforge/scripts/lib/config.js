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

module.exports = { readYamlFile, loadMainConfig, loadProjectConfig, getProjects, getProject, loadMeta };
