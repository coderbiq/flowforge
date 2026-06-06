'use strict';

const fs = require('fs');
const path = require('path');
const os = require('os');

const CONFIG_PATH = 'src/cli/scripts/lib/config.js';

function run(root) {
  let passed = 0;
  let failed = 0;
  const errors = [];

  let config;
  try {
    config = require(path.join(root, CONFIG_PATH));
    passed++;
  } catch (e) {
    failed++;
    errors.push(`Cannot require config.js: ${e.message}`);
    return { passed, failed, errors };
  }

  // --- verify all expected exports ---
  const expectedExports = [
    'readYamlFile', 'loadMainConfig', 'loadProjectConfig',
    'getProjects', 'getProject', 'loadMeta',
    'findProjectRoot', 'findProposalDir',
  ];
  for (const name of expectedExports) {
    if (typeof config[name] === 'function') {
      passed++;
    } else {
      failed++;
      errors.push(`config.${name} is not a function (got ${typeof config[name]})`);
    }
  }

  // --- findProjectRoot tests ---
  const tmpBase = path.join(os.tmpdir(), `ff-test-${Date.now()}`);
  fs.mkdirSync(tmpBase, { recursive: true });

  try {
    // Scenario: no .flowforge directory → returns null
    const noConfig = config.findProjectRoot(tmpBase);
    if (noConfig === null) {
      passed++;
    } else {
      failed++;
      errors.push(`findProjectRoot(no config) should return null, got "${noConfig}"`);
    }

    // Scenario: .flowforge/config.yaml present → returns project root
    const projectDir = path.join(tmpBase, 'my-project');
    const flowforgeDir = path.join(projectDir, '.flowforge');
    fs.mkdirSync(flowforgeDir, { recursive: true });
    fs.writeFileSync(path.join(flowforgeDir, 'config.yaml'), 'projects: []\n');

    const found = config.findProjectRoot(projectDir);
    if (found === projectDir) {
      passed++;
    } else {
      failed++;
      errors.push(`findProjectRoot(from root) should return projectDir, got "${found}", expected "${projectDir}"`);
    }

    // Scenario: from a deep subdirectory → walks up
    const deepDir = path.join(projectDir, 'src', 'components', 'auth');
    fs.mkdirSync(deepDir, { recursive: true });
    const foundDeep = config.findProjectRoot(deepDir);
    if (foundDeep === projectDir) {
      passed++;
    } else {
      failed++;
      errors.push(`findProjectRoot(from deep) should return projectDir, got "${foundDeep}", expected "${projectDir}"`);
    }

    // Scenario: from a file path (not directory) → extracts dirname then walks up
    const fileInDeep = path.join(deepDir, 'login.ts');
    fs.writeFileSync(fileInDeep, '// stub');
    const foundFile = config.findProjectRoot(fileInDeep);
    if (foundFile === projectDir) {
      passed++;
    } else {
      failed++;
      errors.push(`findProjectRoot(from file) should return projectDir, got "${foundFile}", expected "${projectDir}"`);
    }

    // Scenario: sibling project (no .flowforge) → reaches fs root → returns null
    const siblingDir = path.join(tmpBase, 'other-project');
    fs.mkdirSync(siblingDir, { recursive: true });
    const foundSibling = config.findProjectRoot(siblingDir);
    if (foundSibling === null) {
      passed++;
    } else {
      failed++;
      errors.push(`findProjectRoot(sibling w/o config) should return null, got "${foundSibling}"`);
    }

    // --- loadMainConfig tests ---
    const mainConfig = config.loadMainConfig(projectDir);
    if (mainConfig && Array.isArray(mainConfig.projects)) {
      passed++;
    } else {
      failed++;
      errors.push(`loadMainConfig should return config with projects array, got ${JSON.stringify(mainConfig)}`);
    }

    const noMainConfig = config.loadMainConfig(siblingDir);
    if (noMainConfig === null) {
      passed++;
    } else {
      failed++;
      errors.push(`loadMainConfig(no config) should return null, got ${JSON.stringify(noMainConfig)}`);
    }

    // --- getProjects tests ---
    const projects = config.getProjects(projectDir);
    if (Array.isArray(projects)) {
      passed++;
    } else {
      failed++;
      errors.push(`getProjects should return array, got ${typeof projects}`);
    }

    const noProjects = config.getProjects(siblingDir);
    if (Array.isArray(noProjects) && noProjects.length === 0) {
      passed++;
    } else {
      failed++;
      errors.push(`getProjects(no config) should return [], got ${JSON.stringify(noProjects)}`);
    }

    // --- getProject tests ---
    const nonExistentProject = config.getProject(projectDir, 'nonexistent');
    if (nonExistentProject === null) {
      passed++;
    } else {
      failed++;
      errors.push(`getProject(nonexistent) should return null, got ${JSON.stringify(nonExistentProject)}`);
    }

    // --- readYamlFile tests ---
    const configPath = path.join(flowforgeDir, 'config.yaml');
    const yamlContent = config.readYamlFile(configPath);
    if (yamlContent && Array.isArray(yamlContent.projects)) {
      passed++;
    } else {
      failed++;
      errors.push(`readYamlFile should parse YAML, got ${JSON.stringify(yamlContent)}`);
    }

    const missingYaml = config.readYamlFile(path.join(tmpBase, 'nope.yaml'));
    if (missingYaml === null) {
      passed++;
    } else {
      failed++;
      errors.push(`readYamlFile(missing) should return null, got ${JSON.stringify(missingYaml)}`);
    }

    // --- loadMeta tests ---
    const metaContent = config.loadMeta(projectDir);
    if (metaContent === null) {
      passed++;
    } else {
      failed++;
      errors.push(`loadMeta(no meta) should return null, got ${JSON.stringify(metaContent)}`);
    }

  } finally {
    // cleanup
    fs.rmSync(tmpBase, { recursive: true, force: true });
  }

  return { passed, failed, errors };
}

module.exports = { run };
