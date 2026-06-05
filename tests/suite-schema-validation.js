'use strict';

const path = require('path');
const fs = require('fs');

const SCHEMA_DIR = 'src/flowforge/schema';

function run(root) {
  let passed = 0;
  let failed = 0;
  const errors = [];

  const schemaDir = path.join(root, SCHEMA_DIR);
  if (!fs.existsSync(schemaDir)) {
    return { passed: 0, failed: 1, errors: ['Schema directory not found: ' + schemaDir] };
  }

  const schemaFiles = fs.readdirSync(schemaDir).filter(f => f.endsWith('.json'));
  if (schemaFiles.length === 0) {
    return { passed: 0, failed: 1, errors: ['No JSON schema files found'] };
  }

  for (const file of schemaFiles) {
    const filePath = path.join(schemaDir, file);
    let content;
    try {
      content = JSON.parse(fs.readFileSync(filePath, 'utf8'));
    } catch (e) {
      failed++;
      errors.push(file + ': invalid JSON - ' + e.message);
      continue;
    }
    passed++;

    if (!content.$schema) {
      failed++;
      errors.push(file + ': missing $schema field');
    } else {
      passed++;
    }

    if (!content.type || content.type !== 'object') {
      failed++;
      errors.push(file + ': root type must be "object"');
    } else {
      passed++;
    }

    if (!content.properties || Object.keys(content.properties).length === 0) {
      failed++;
      errors.push(file + ': no properties defined');
    } else {
      passed++;
    }

    if (content.required && !Array.isArray(content.required)) {
      failed++;
      errors.push(file + ': "required" must be an array');
    } else {
      passed++;
    }
  }

  const validateScripts = fs.readdirSync(path.join(root, 'src/cli/scripts'))
    .filter(f => f.startsWith('validate-') && f.endsWith('.js'));

  for (const script of validateScripts) {
    const content = fs.readFileSync(path.join(root, 'src/cli/scripts', script), 'utf8');
    if (content.includes('.schema.json') || content.includes('ajv') || content.includes('jsonschema')) {
      passed++;
    }
  }
  passed++;

  return { passed, failed, errors };
}

module.exports = { run };
